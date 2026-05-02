package handler

import (
	"encoding/gob"
	"image"
	"image/color"
	imagedraw "image/draw"
	"image/jpeg"
	"math/rand/v2"
	"net/http"
	"path"
	"strings"
	"time"

	"carRental/internal/config"
	"carRental/internal/model"
	"carRental/internal/service"
	"carRental/internal/view"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	xdraw "golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

type Deps struct {
	Cfg           config.Config
	AuthService   *service.AuthService
	MenuService   *service.MenuService
	SystemService *service.SystemService
	BusService    *service.BusService
	StatService   *service.StatService
	FileService   *service.FileService
}

func NewRouter(d Deps) *gin.Engine {
	gob.Register(model.User{})

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	basePath := normalizeBasePath(d.Cfg.BasePath)

	store := cookie.NewStore([]byte("carRental-secret-change-me"))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	r.Use(sessions.Sessions("carRental.sid", store))

	// Static assets: JSP uses `${yeqifu}/static/...`
	r.Static("/static", d.Cfg.StaticDir)
	if basePath != "" {
		r.Static(path.Join(basePath, "/static"), d.Cfg.StaticDir)
	}

	root := r.Group(basePath)

	druid := root.Group("/druid")
	druid.GET("/toLogin.action", func(c *gin.Context) {
		renderView(c, d.Cfg, "system/druid/druidLogin", map[string]any{"error": ""})
	})
	druid.POST("/login.action", func(c *gin.Context) {
		username := strings.TrimSpace(c.PostForm("username"))
		password := strings.TrimSpace(c.PostForm("password"))
		if username == "" || password == "" {
			renderView(c, d.Cfg, "system/druid/druidLogin", map[string]any{"error": "请输入账号和密码"})
			return
		}
		if username != d.Cfg.MonitorUser || password != d.Cfg.MonitorPass {
			renderView(c, d.Cfg, "system/druid/druidLogin", map[string]any{"error": "账号或密码错误"})
			return
		}
		sess := sessions.Default(c)
		sess.Set("druid_authed", true)
		_ = sess.Save()
		c.Redirect(http.StatusFound, withBase(basePath, "/druid/"))
	})
	druid.POST("/logout.action", func(c *gin.Context) {
		sess := sessions.Default(c)
		sess.Delete("druid_authed")
		_ = sess.Save()
		c.Redirect(http.StatusFound, withBase(basePath, "/druid/toLogin.action"))
	})

	// Lightweight replacement for the original druid monitor entry.
	druid.GET("/", func(c *gin.Context) {
		sess := sessions.Default(c)
		authed, _ := sess.Get("druid_authed").(bool)
		if !authed {
			c.Redirect(http.StatusFound, withBase(basePath, "/druid/toLogin.action"))
			return
		}
		if d.SystemService == nil {
			c.String(http.StatusOK, "druid replacement unavailable")
			return
		}
		stats, err := d.SystemService.DBStats(c.Request.Context())
		if err != nil {
			c.String(http.StatusInternalServerError, "load db stats error")
			return
		}
		html := "<html><head><meta charset=\"utf-8\"><title>DB Monitor</title></head><body>" +
			"<h2>Go DB Monitor</h2>" +
			"<table border=\"1\" cellpadding=\"8\">" +
			"<tr><th>OpenConnections</th><td>" + stats["open_connections"] + "</td></tr>" +
			"<tr><th>InUse</th><td>" + stats["in_use"] + "</td></tr>" +
			"<tr><th>Idle</th><td>" + stats["idle"] + "</td></tr>" +
			"<tr><th>WaitCount</th><td>" + stats["wait_count"] + "</td></tr>" +
			"<tr><th>WaitDuration</th><td>" + stats["wait_duration"] + "</td></tr>" +
			"<tr><th>MaxOpenConnections</th><td>" + stats["max_open_connections"] + "</td></tr>" +
			"</table></body></html>"
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	})
	druid.GET("/index.html", func(c *gin.Context) {
		c.Redirect(http.StatusFound, withBase(basePath, "/druid/"))
	})

	// Root: behave like original index.jsp forward, and use it as logout target.
	r.GET("/", func(c *gin.Context) {
		sess := sessions.Default(c)
		sess.Clear()
		_ = sess.Save()
		c.Redirect(http.StatusFound, withBase(basePath, "/login/toLogin.action"))
	})
	root.GET("/", func(c *gin.Context) {
		sess := sessions.Default(c)
		sess.Clear()
		_ = sess.Save()
		c.Redirect(http.StatusFound, withBase(basePath, "/login/toLogin.action"))
	})

	// Login flow
	root.GET("/login/toLogin.action", func(c *gin.Context) {
		renderView(c, d.Cfg, "system/main/login", map[string]any{
			"error": "",
		})
	})
	root.POST("/login/login.action", func(c *gin.Context) {
		loginName := c.PostForm("loginname")
		pwd := c.PostForm("pwd")
		code := c.PostForm("code")

		sess := sessions.Default(c)
		want, _ := sess.Get("code").(string)
		if want == "" || code != want {
			renderView(c, d.Cfg, "system/main/login", map[string]any{
				"error": "验证码错误",
			})
			return
		}

		u, err := d.AuthService.Login(c.Request.Context(), loginName, pwd)
		if err != nil {
			c.String(http.StatusInternalServerError, "login error")
			return
		}
		if u == nil {
			renderView(c, d.Cfg, "system/main/login", map[string]any{
				"error": "用户名或密码错误",
			})
			return
		}

		sess.Set("user", *u)
		_ = sess.Save()

		_ = d.AuthService.AddLoginLog(c.Request.Context(), u.RealName+"-"+u.LoginName, c.ClientIP(), time.Now())
		renderViewWithUser(c, d.Cfg, "system/main/index", *u, nil)
	})
	root.GET("/login/getCode.action", func(c *gin.Context) {
		code := randomDigits(4)
		sess := sessions.Default(c)
		sess.Set("code", code)
		_ = sess.Save()

		img := drawCaptcha(code, 180, 50)
		c.Header("Content-Type", "image/jpeg")
		_ = jpeg.Encode(c.Writer, img, &jpeg.Options{Quality: 85})
	})

	authed := root.Group("/")
	authed.Use(requireLogin(basePath))

	// Desk
	authed.GET("/desk/toDeskManager.action", func(c *gin.Context) {
		u, ok := getSessionUser(c)
		if !ok {
			c.Redirect(http.StatusFound, withBase(basePath, "/login/toLogin.action"))
			return
		}
		renderViewWithUser(c, d.Cfg, "system/main/deskManager", u, nil)
	})

	// Sys page routes
	authed.GET("/sys/toChangePassword.action", pageWithUser(d.Cfg, "system/user/changePassword"))
	authed.GET("/sys/toMenuManager.action", pageWithUser(d.Cfg, "system/menu/menuManager"))
	authed.GET("/sys/toMenuLeft.action", pageWithUser(d.Cfg, "system/menu/menuLeft"))
	authed.GET("/sys/toMenuRight.action", pageWithUser(d.Cfg, "system/menu/menuRight"))
	authed.GET("/sys/toRoleManager.action", pageWithUser(d.Cfg, "system/role/roleManager"))
	authed.GET("/sys/toUserManager.action", pageWithUser(d.Cfg, "system/user/userManager"))
	authed.GET("/sys/toLogInfoManager.action", pageWithUser(d.Cfg, "system/logInfo/logInfoManager"))
	authed.GET("/sys/toNewsManager.action", pageWithUser(d.Cfg, "system/news/newsManager"))
	authed.GET("/sys/toMessageManager.action", pageWithUser(d.Cfg, "system/message/messageManager"))

	// Bus page routes
	authed.GET("/bus/toCustomerManager.action", pageWithUser(d.Cfg, "business/customer/customerManager"))
	authed.GET("/bus/toCarManager.action", pageWithUser(d.Cfg, "business/car/carManager"))
	authed.GET("/bus/toRentCarManager.action", pageWithUser(d.Cfg, "business/rent/rentCarManager"))
	authed.GET("/bus/toRentManager.action", pageWithUser(d.Cfg, "business/rent/rentManager"))
	authed.GET("/bus/toCheckCarManager.action", pageWithUser(d.Cfg, "business/check/checkCarManager"))
	authed.GET("/bus/toCheckManager.action", pageWithUser(d.Cfg, "business/check/checkManager"))
	authed.GET("/bus/toFranchiseeManager.action", pageWithUser(d.Cfg, "business/franchisee/franchiseeManagers"))

	// Stat page routes
	authed.GET("/stat/toCustomerAreaStat.action", pageWithUser(d.Cfg, "stat/customerAreaStat"))
	authed.GET("/stat/toCustomerAreaSexStat.action", pageWithUser(d.Cfg, "stat/customerAreaSexStat"))
	authed.GET("/stat/toOpernameYearGradeStat.action", pageWithUser(d.Cfg, "stat/opernameYearGradeStat"))
	authed.GET("/stat/toCompanyYearGradeStat.action", pageWithUser(d.Cfg, "stat/companyYearGradeStat"))

	// APIs: menu (needed by index.js)
	authed.GET("/menu/loadIndexleftMenuJson.action", func(c *gin.Context) {
		u, ok := getSessionUser(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"code": -1, "msg": "not login"})
			return
		}
		available := 1
		var menus []model.Menu
		var err error
		if u.Type == 1 {
			menus, err = d.MenuService.QueryAllMenus(c.Request.Context(), available)
		} else {
			menus, err = d.MenuService.QueryMenusByUID(c.Request.Context(), available, u.UserID)
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "db error"})
			return
		}
		c.JSON(http.StatusOK, service.MenusToTree(menus, 1))
	})

	registerModules(authed, d)

	return r
}

func requireLogin(basePath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := c.Request.URL.Path
		basePath = normalizeBasePath(basePath)
		if basePath != "" && strings.HasPrefix(p, basePath) {
			p = strings.TrimPrefix(p, basePath)
			if p == "" {
				p = "/"
			}
		}
		if p == "/menu/loadMenuManagerLeftTreeJson.action" || p == "/menu/loadAllMenu.action" {
			c.Next()
			return
		}
		if strings.HasPrefix(p, "/login/") {
			c.Next()
			return
		}
		if _, ok := getSessionUser(c); !ok {
			c.Redirect(http.StatusFound, withBase(basePath, "/login/toLogin.action"))
			c.Abort()
			return
		}
		c.Next()
	}
}

func getSessionUser(c *gin.Context) (model.User, bool) {
	sess := sessions.Default(c)
	v := sess.Get("user")
	u, ok := v.(model.User)
	return u, ok
}

func pageWithUser(cfg config.Config, viewPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		u, ok := getSessionUser(c)
		if !ok {
			c.Redirect(http.StatusFound, withBase(normalizeBasePath(cfg.BasePath), "/login/toLogin.action"))
			return
		}
		renderViewWithUser(c, cfg, viewPath, u, nil)
	}
}

func renderView(c *gin.Context, cfg config.Config, viewPath string, extra map[string]any) {
	ctx := map[string]any{
		"yeqifu": normalizeBasePath(cfg.BasePath),
	}
	for k, v := range extra {
		ctx[k] = v
	}
	b, err := view.RenderJSPFile(cfg.ViewRoot, viewPath, ctx)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", b)
}

func renderViewWithUser(c *gin.Context, cfg config.Config, viewPath string, u model.User, extra map[string]any) {
	ctx := map[string]any{
		"yeqifu": normalizeBasePath(cfg.BasePath),
		"user": map[string]any{
			"userid":    u.UserID,
			"loginname": u.LoginName,
			"realname":  u.RealName,
			"type":      u.Type,
		},
	}
	for k, v := range extra {
		ctx[k] = v
	}
	b, err := view.RenderJSPFile(cfg.ViewRoot, viewPath, ctx)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", b)
}

func normalizeBasePath(p string) string {
	p = strings.TrimSpace(p)
	if p == "" || p == "/" {
		return ""
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return strings.TrimRight(p, "/")
}

func withBase(basePath, p string) string {
	basePath = normalizeBasePath(basePath)
	if basePath == "" {
		return p
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return basePath + p
}

func randomDigits(n int) string {
	const digits = "0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = digits[rand.IntN(len(digits))]
	}
	return string(b)
}

func drawCaptcha(code string, w, h int) image.Image {
	baseW, baseH := 116, 36
	small := image.NewRGBA(image.Rect(0, 0, baseW, baseH))
	imagedraw.Draw(small, small.Bounds(), &image.Uniform{C: color.White}, image.Point{}, imagedraw.Src)

	// Noise
	for i := 0; i < 160; i++ {
		x := rand.IntN(baseW)
		y := rand.IntN(baseH)
		small.Set(x, y, color.RGBA{uint8(rand.IntN(255)), uint8(rand.IntN(255)), uint8(rand.IntN(255)), 255})
	}
	for i := 0; i < 5; i++ {
		x1 := rand.IntN(baseW)
		y1 := rand.IntN(baseH)
		x2 := rand.IntN(baseW)
		y2 := rand.IntN(baseH)
		col := color.RGBA{uint8(rand.IntN(200)), uint8(rand.IntN(200)), uint8(rand.IntN(200)), 255}
		drawLine(small, x1, y1, x2, y2, col)
	}

	d := &font.Drawer{
		Dst:  small,
		Src:  image.NewUniform(color.Black),
		Face: basicfont.Face7x13,
		Dot:  fixed.P(10, 24),
	}
	d.DrawString(code)
	d.Dot = fixed.P(11, 24)
	d.DrawString(code)

	out := image.NewRGBA(image.Rect(0, 0, w, h))
	xdraw.NearestNeighbor.Scale(out, out.Bounds(), small, small.Bounds(), xdraw.Over, nil)
	return out
}

func drawLine(img *image.RGBA, x1, y1, x2, y2 int, c color.Color) {
	dx := abs(x2 - x1)
	sx := -1
	if x1 < x2 {
		sx = 1
	}
	dy := -abs(y2 - y1)
	sy := -1
	if y1 < y2 {
		sy = 1
	}
	err := dx + dy
	for {
		img.Set(x1, y1, c)
		if x1 == x2 && y1 == y2 {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x1 += sx
		}
		if e2 <= dx {
			err += dx
			y1 += sy
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
