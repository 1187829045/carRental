package handler

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"carRental/internal/model"
	"carRental/internal/service"

	"github.com/gin-gonic/gin"
)

const (
	defaultCarImage = "images/defaultcarimage.jpg"
	uploadTempTag   = "_temp"
)

func registerModules(authed *gin.RouterGroup, d Deps) {
	registerMenuRoutes(authed, d)
	registerSystemRoutes(authed, d)
	registerBusRoutes(authed, d)
	registerStatRoutes(authed, d)
	registerFileRoutes(authed, d)
}

func registerMenuRoutes(rg *gin.RouterGroup, d Deps) {
	loadMenuManagerLeftTree := func(c *gin.Context) {
		spread := c.Query("spread")
		available := 1
		menus, err := d.MenuService.QueryAllMenus(c.Request.Context(), available)
		if err != nil {
			c.JSON(http.StatusOK, DataGridView{Code: 0, Msg: "菜单树加载失败", Data: []model.TreeNode{}})
			return
		}
		data := make([]model.TreeNode, 0, len(menus))
		for _, m := range menus {
			data = append(data, model.TreeNode{
				ID:       m.ID,
				ParentID: m.Pid,
				Title:    m.Title,
				Icon:     m.Icon,
				Href:     m.Href,
				Spread:   m.Spread == 1 || spread == "1",
				Target:   m.Target,
				CheckArr: "0",
			})
		}
		c.JSON(http.StatusOK, NewData(data))
	}
	rg.GET("/menu/loadMenuManagerLeftTreeJson.action", loadMenuManagerLeftTree)
	rg.POST("/menu/loadMenuManagerLeftTreeJson.action", loadMenuManagerLeftTree)

	rg.GET("/menu/loadAllMenu.action", func(c *gin.Context) {
		page := intParam(c, "page")
		if page == 0 {
			page = 1
		}
		limit := intParam(c, "limit")
		if limit == 0 {
			limit = 10
		}
		var available *int
		if v := strings.TrimSpace(c.Query("available")); v != "" {
			n := intParam(c, "available")
			available = &n
		}
		var idp *int
		if v := strings.TrimSpace(c.Query("id")); v != "" {
			n := intParam(c, "id")
			idp = &n
		}
		count, data, err := d.MenuService.QueryMenuList(c.Request.Context(), c.Query("title"), available, idp, page, limit)
		if err != nil {
			fail(c, "查询菜单失败")
			return
		}
		c.JSON(http.StatusOK, NewPage(count, data))
	})

	rg.POST("/menu/addMenu.action", func(c *gin.Context) {
		err := d.MenuService.AddMenu(c.Request.Context(), model.Menu{
			Pid:       intParam(c, "pid"),
			Title:     c.PostForm("title"),
			Href:      ptrString(c.PostForm("href")),
			Spread:    intParam(c, "spread"),
			Target:    ptrString(c.PostForm("target")),
			Icon:      ptrString(c.PostForm("icon")),
			Available: intParam(c, "available"),
		})
		if err != nil {
			c.JSON(http.StatusOK, AddError)
			return
		}
		c.JSON(http.StatusOK, AddSuccess)
	})

	rg.POST("/menu/updateMenu.action", func(c *gin.Context) {
		err := d.MenuService.UpdateMenu(c.Request.Context(), model.Menu{
			ID:        intParam(c, "id"),
			Pid:       intParam(c, "pid"),
			Title:     c.PostForm("title"),
			Href:      ptrString(c.PostForm("href")),
			Spread:    intParam(c, "spread"),
			Target:    ptrString(c.PostForm("target")),
			Icon:      ptrString(c.PostForm("icon")),
			Available: intParam(c, "available"),
		})
		if err != nil {
			c.JSON(http.StatusOK, UpdateError)
			return
		}
		c.JSON(http.StatusOK, UpdateSuccess)
	})

	rg.POST("/menu/checkMenuHasChildren.action", func(c *gin.Context) {
		count, err := d.MenuService.CountByPid(c.Request.Context(), intParam(c, "id"))
		if err != nil || count <= 0 {
			c.JSON(http.StatusOK, StatusFalse)
			return
		}
		c.JSON(http.StatusOK, StatusTrue)
	})

	rg.POST("/menu/deleteMenu.action", func(c *gin.Context) {
		if err := d.MenuService.DeleteMenu(c.Request.Context(), intParam(c, "id")); err != nil {
			c.JSON(http.StatusOK, DeleteError)
			return
		}
		c.JSON(http.StatusOK, DeleteSuccess)
	})
}

func registerSystemRoutes(rg *gin.RouterGroup, d Deps) {
	rg.GET("/user/loadAllUser.action", func(c *gin.Context) {
		count, data, err := d.SystemService.QueryUsers(c.Request.Context(), flatQuery(c))
		if err != nil {
			fail(c, "查询用户失败")
			return
		}
		c.JSON(http.StatusOK, NewPage(count, data))
	})
	rg.POST("/user/addUser.action", func(c *gin.Context) {
		u := model.SysUser{
			LoginName: c.PostForm("loginname"),
			Identity:  ptrString(c.PostForm("identity")),
			RealName:  ptrString(c.PostForm("realname")),
			Sex:       ptrIntMaybe(c.PostForm("sex")),
			Address:   ptrString(c.PostForm("address")),
			Phone:     ptrString(c.PostForm("phone")),
			Position:  ptrString(c.PostForm("position")),
			Available: ptrIntMaybe(c.PostForm("available")),
		}
		if err := d.SystemService.AddUser(c.Request.Context(), u); err != nil {
			c.JSON(http.StatusOK, AddError)
			return
		}
		c.JSON(http.StatusOK, AddSuccess)
	})
	rg.POST("/user/updateUser.action", func(c *gin.Context) {
		u := model.SysUser{
			UserID:    intParam(c, "userid"),
			LoginName: c.PostForm("loginname"),
			Identity:  ptrString(c.PostForm("identity")),
			RealName:  ptrString(c.PostForm("realname")),
			Sex:       ptrIntMaybe(c.PostForm("sex")),
			Address:   ptrString(c.PostForm("address")),
			Phone:     ptrString(c.PostForm("phone")),
			Position:  ptrString(c.PostForm("position")),
			Type:      ptrIntMaybe(c.PostForm("type")),
			Available: ptrIntMaybe(c.PostForm("available")),
		}
		if err := d.SystemService.UpdateUser(c.Request.Context(), u); err != nil {
			c.JSON(http.StatusOK, UpdateError)
			return
		}
		c.JSON(http.StatusOK, UpdateSuccess)
	})
	rg.POST("/user/deleteUser.action", func(c *gin.Context) {
		if err := d.SystemService.DeleteUser(c.Request.Context(), intParam(c, "userid")); err != nil {
			c.JSON(http.StatusOK, DeleteError)
			return
		}
		c.JSON(http.StatusOK, DeleteSuccess)
	})
	rg.POST("/user/deleteBatchUser.action", func(c *gin.Context) {
		for _, id := range intArrayParam(c, "ids") {
			if err := d.SystemService.DeleteUser(c.Request.Context(), id); err != nil {
				c.JSON(http.StatusOK, DeleteError)
				return
			}
		}
		c.JSON(http.StatusOK, DeleteSuccess)
	})
	rg.POST("/user/resetUserPwd.action", func(c *gin.Context) {
		if err := d.SystemService.ResetUserPwd(c.Request.Context(), intParam(c, "userid")); err != nil {
			c.JSON(http.StatusOK, ResetError)
			return
		}
		c.JSON(http.StatusOK, ResetSuccess)
	})
	rg.GET("/user/initUserRole.action", func(c *gin.Context) {
		data, err := d.SystemService.QueryUserRoles(c.Request.Context(), intParam(c, "userid"))
		if err != nil {
			fail(c, "查询角色失败")
			return
		}
		c.JSON(http.StatusOK, NewData(data))
	})
	rg.POST("/user/saveUserRole.action", func(c *gin.Context) {
		if err := d.SystemService.SaveUserRole(c.Request.Context(), intParam(c, "userid"), intArrayParam(c, "ids")); err != nil {
			c.JSON(http.StatusOK, DispatchError)
			return
		}
		c.JSON(http.StatusOK, DispatchSuccess)
	})
	rg.POST("/user/changePassword.action", func(c *gin.Context) {
		u, ok := getSessionUser(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "未登录"})
			return
		}
		res := d.SystemService.ChangePassword(c.Request.Context(), u.UserID, c.PostForm("oldPassword"), c.PostForm("newPassword"), c.PostForm("confirmPassword"))
		c.JSON(http.StatusOK, res)
	})

	rg.GET("/role/loadAllRole.action", func(c *gin.Context) {
		count, data, err := d.SystemService.QueryRoles(c.Request.Context(), flatQuery(c))
		if err != nil {
			fail(c, "查询角色失败")
			return
		}
		c.JSON(http.StatusOK, NewPage(count, data))
	})
	rg.POST("/role/addRole.action", func(c *gin.Context) {
		if err := d.SystemService.AddRole(c.Request.Context(), model.Role{
			RoleName:  c.PostForm("rolename"),
			RoleDesc:  c.PostForm("roledesc"),
			Available: intParam(c, "available"),
		}); err != nil {
			c.JSON(http.StatusOK, AddError)
			return
		}
		c.JSON(http.StatusOK, AddSuccess)
	})
	rg.POST("/role/updateRole.action", func(c *gin.Context) {
		if err := d.SystemService.UpdateRole(c.Request.Context(), model.Role{
			RoleID:    intParam(c, "roleid"),
			RoleName:  c.PostForm("rolename"),
			RoleDesc:  c.PostForm("roledesc"),
			Available: intParam(c, "available"),
		}); err != nil {
			c.JSON(http.StatusOK, UpdateError)
			return
		}
		c.JSON(http.StatusOK, UpdateSuccess)
	})
	rg.POST("/role/deleteRole.action", func(c *gin.Context) {
		if err := d.SystemService.DeleteRole(c.Request.Context(), intParam(c, "roleid")); err != nil {
			c.JSON(http.StatusOK, DeleteError)
			return
		}
		c.JSON(http.StatusOK, DeleteSuccess)
	})
	rg.POST("/role/deleteBatchRole.action", func(c *gin.Context) {
		for _, id := range intArrayParam(c, "ids") {
			if err := d.SystemService.DeleteRole(c.Request.Context(), id); err != nil {
				c.JSON(http.StatusOK, DeleteError)
				return
			}
		}
		c.JSON(http.StatusOK, DeleteSuccess)
	})
	rg.GET("/role/initRoleMenuTreeJson.action", func(c *gin.Context) {
		data, err := d.SystemService.InitRoleMenuTree(c.Request.Context(), intParam(c, "roleid"))
		if err != nil {
			fail(c, "查询角色菜单失败")
			return
		}
		c.JSON(http.StatusOK, NewData(data))
	})
	rg.POST("/role/saveRoleMenu.action", func(c *gin.Context) {
		if err := d.SystemService.SaveRoleMenu(c.Request.Context(), intParam(c, "roleid"), intArrayParam(c, "ids")); err != nil {
			c.JSON(http.StatusOK, DispatchError)
			return
		}
		c.JSON(http.StatusOK, DispatchSuccess)
	})

	rg.GET("/logInfo/loadAllLogInfo.action", func(c *gin.Context) {
		count, data, err := d.SystemService.QueryLogInfos(c.Request.Context(), flatQuery(c))
		if err != nil {
			fail(c, "查询日志失败")
			return
		}
		c.JSON(http.StatusOK, NewPage(count, data))
	})
	rg.POST("/logInfo/deleteLogInfo.action", func(c *gin.Context) {
		if err := d.SystemService.DeleteLogInfo(c.Request.Context(), intParam(c, "id")); err != nil {
			c.JSON(http.StatusOK, DeleteError)
			return
		}
		c.JSON(http.StatusOK, DeleteSuccess)
	})
	rg.POST("/logInfo/deleteBatchLogInfo.action", func(c *gin.Context) {
		for _, id := range intArrayParam(c, "ids") {
			if err := d.SystemService.DeleteLogInfo(c.Request.Context(), id); err != nil {
				c.JSON(http.StatusOK, DeleteError)
				return
			}
		}
		c.JSON(http.StatusOK, DeleteSuccess)
	})

	rg.GET("/news/loadAllNews.action", func(c *gin.Context) {
		count, data, err := d.SystemService.QueryNews(c.Request.Context(), flatQuery(c))
		if err != nil {
			fail(c, "查询公告失败")
			return
		}
		c.JSON(http.StatusOK, NewPage(count, data))
	})
	rg.POST("/news/addNews.action", func(c *gin.Context) {
		u, _ := getSessionUser(c)
		if err := d.SystemService.AddNews(c.Request.Context(), c.PostForm("title"), c.PostForm("content"), u.RealName); err != nil {
			c.JSON(http.StatusOK, AddError)
			return
		}
		c.JSON(http.StatusOK, AddSuccess)
	})
	rg.POST("/news/updateNews.action", func(c *gin.Context) {
		if err := d.SystemService.UpdateNews(c.Request.Context(), intParam(c, "id"), c.PostForm("title"), c.PostForm("content")); err != nil {
			c.JSON(http.StatusOK, UpdateError)
			return
		}
		c.JSON(http.StatusOK, UpdateSuccess)
	})
	rg.POST("/news/deleteNews.action", func(c *gin.Context) {
		if err := d.SystemService.DeleteNews(c.Request.Context(), intParam(c, "id")); err != nil {
			c.JSON(http.StatusOK, DeleteError)
			return
		}
		c.JSON(http.StatusOK, DeleteSuccess)
	})
	rg.POST("/news/deleteBatchNews.action", func(c *gin.Context) {
		for _, id := range intArrayParam(c, "ids") {
			if err := d.SystemService.DeleteNews(c.Request.Context(), id); err != nil {
				c.JSON(http.StatusOK, DeleteError)
				return
			}
		}
		c.JSON(http.StatusOK, DeleteSuccess)
	})
	rg.GET("/news/loadNewsById.action", func(c *gin.Context) {
		x, err := d.SystemService.GetNewsByID(c.Request.Context(), intParam(c, "id"))
		if err != nil || x == nil {
			c.JSON(http.StatusOK, gin.H{})
			return
		}
		c.JSON(http.StatusOK, x)
	})

	rg.GET("/message/loadAllMessage.action", func(c *gin.Context) {
		count, data, err := d.SystemService.QueryMessages(c.Request.Context(), flatQuery(c))
		if err != nil {
			fail(c, "查询留言失败")
			return
		}
		c.JSON(http.StatusOK, NewPage(count, data))
	})
	rg.POST("/message/addMessage.action", func(c *gin.Context) {
		u, _ := getSessionUser(c)
		if err := d.SystemService.AddMessage(c.Request.Context(), c.PostForm("title"), c.PostForm("content"), u.RealName); err != nil {
			c.JSON(http.StatusOK, AddError)
			return
		}
		c.JSON(http.StatusOK, AddSuccess)
	})
	rg.POST("/message/updateMessage.action", func(c *gin.Context) {
		if err := d.SystemService.UpdateMessage(c.Request.Context(), intParam(c, "id"), c.PostForm("title"), c.PostForm("content")); err != nil {
			c.JSON(http.StatusOK, UpdateError)
			return
		}
		c.JSON(http.StatusOK, UpdateSuccess)
	})
	rg.POST("/message/deleteMessage.action", func(c *gin.Context) {
		if err := d.SystemService.DeleteMessage(c.Request.Context(), intParam(c, "id")); err != nil {
			c.JSON(http.StatusOK, DeleteError)
			return
		}
		c.JSON(http.StatusOK, DeleteSuccess)
	})
	rg.POST("/message/deleteBatchMessage.action", func(c *gin.Context) {
		for _, id := range intArrayParam(c, "ids") {
			if err := d.SystemService.DeleteMessage(c.Request.Context(), id); err != nil {
				c.JSON(http.StatusOK, DeleteError)
				return
			}
		}
		c.JSON(http.StatusOK, DeleteSuccess)
	})
	rg.GET("/message/loadMessageById.action", func(c *gin.Context) {
		x, err := d.SystemService.GetMessageByID(c.Request.Context(), intParam(c, "id"))
		if err != nil || x == nil {
			c.JSON(http.StatusOK, gin.H{})
			return
		}
		c.JSON(http.StatusOK, x)
	})
}

func registerBusRoutes(rg *gin.RouterGroup, d Deps) {
	rg.GET("/car/loadAllCar.action", func(c *gin.Context) {
		q := flatQuery(c)
		if u, ok := getSessionUser(c); ok && u.Type != 1 {
			q["operid"] = fmt.Sprint(u.UserID)
		}
		count, data, err := d.BusService.QueryCars(c.Request.Context(), q)
		if err != nil {
			fail(c, "查询车辆失败")
			return
		}
		c.JSON(http.StatusOK, NewPage(count, data))
	})
	rg.POST("/car/addCar.action", func(c *gin.Context) {
		operid := 1 // default to admin
		if u, ok := getSessionUser(c); ok && u.Type != 1 {
			operid = u.UserID
		}
		x := model.Car{
			CarNumber:   c.PostForm("carnumber"),
			CarType:     c.PostForm("cartype"),
			Color:       c.PostForm("color"),
			Price:       floatParam(c, "price"),
			RentPrice:   floatParam(c, "rentprice"),
			Deposit:     floatParam(c, "deposit"),
			IsRenting:   intParam(c, "isrenting"),
			Description: c.PostForm("description"),
			CarImg:      c.PostForm("carimg"),
			CreateTime:  time.Now(),
			OperId:      operid,
		}
		if x.CarImg == "" {
			x.CarImg = defaultCarImage
		}
		if err := d.BusService.AddCar(c.Request.Context(), x); err != nil {
			c.JSON(http.StatusOK, AddError)
			return
		}
		c.JSON(http.StatusOK, AddSuccess)
	})
	rg.POST("/car/updateCar.action", func(c *gin.Context) {
		carnumber := c.PostForm("carnumber")
		car, err := d.BusService.GetCar(c.Request.Context(), carnumber)
		if err != nil || car == nil {
			c.JSON(http.StatusOK, UpdateError)
			return
		}
		if u, ok := getSessionUser(c); ok && u.Type != 1 {
			if car.OperId != u.UserID {
				c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "无权限修改他人的车辆"})
				return
			}
		}
		x := model.Car{
			CarNumber:   carnumber,
			CarType:     c.PostForm("cartype"),
			Color:       c.PostForm("color"),
			Price:       floatParam(c, "price"),
			RentPrice:   floatParam(c, "rentprice"),
			Deposit:     floatParam(c, "deposit"),
			IsRenting:   intParam(c, "isrenting"),
			Description: c.PostForm("description"),
			CarImg:      c.PostForm("carimg"),
			OperId:      car.OperId, // Preserve original OperId
		}
		if err := d.BusService.UpdateCar(c.Request.Context(), x); err != nil {
			c.JSON(http.StatusOK, UpdateError)
			return
		}
		c.JSON(http.StatusOK, UpdateSuccess)
	})
	rg.POST("/car/deleteCar.action", func(c *gin.Context) {
		carnumber := c.PostForm("carnumber")
		if u, ok := getSessionUser(c); ok && u.Type != 1 {
			car, err := d.BusService.GetCar(c.Request.Context(), carnumber)
			if err != nil || car == nil || car.OperId != u.UserID {
				c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "无权限删除他人的车辆"})
				return
			}
		}
		if err := d.BusService.DeleteCar(c.Request.Context(), carnumber); err != nil {
			c.JSON(http.StatusOK, DeleteError)
			return
		}
		c.JSON(http.StatusOK, DeleteSuccess)
	})
	rg.POST("/car/deleteBatchCar.action", func(c *gin.Context) {
		ids := stringArrayParam(c, "ids")
		if u, ok := getSessionUser(c); ok && u.Type != 1 {
			for _, id := range ids {
				car, err := d.BusService.GetCar(c.Request.Context(), id)
				if err != nil || car == nil || car.OperId != u.UserID {
					c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "无权限删除他人的车辆"})
					return
				}
			}
		}
		for _, id := range ids {
			if err := d.BusService.DeleteCar(c.Request.Context(), id); err != nil {
				c.JSON(http.StatusOK, DeleteError)
				return
			}
		}
		c.JSON(http.StatusOK, DeleteSuccess)
	})

	rg.GET("/customer/loadAllCustomer.action", func(c *gin.Context) {
		count, data, err := d.BusService.QueryCustomers(c.Request.Context(), flatQuery(c))
		if err != nil {
			fail(c, "查询客户失败")
			return
		}
		c.JSON(http.StatusOK, NewPage(count, data))
	})
	rg.POST("/customer/addCustomer.action", func(c *gin.Context) {
		err := d.BusService.AddCustomer(c.Request.Context(), model.Customer{
			Identity:   c.PostForm("identity"),
			CustName:   c.PostForm("custname"),
			Sex:        intParam(c, "sex"),
			Address:    c.PostForm("address"),
			Phone:      c.PostForm("phone"),
			Career:     c.PostForm("career"),
			CreateTime: time.Now(),
		})
		if err != nil {
			c.JSON(http.StatusOK, AddError)
			return
		}
		c.JSON(http.StatusOK, AddSuccess)
	})
	rg.POST("/customer/updateCustomer.action", func(c *gin.Context) {
		err := d.BusService.UpdateCustomer(c.Request.Context(), model.Customer{
			Identity: c.PostForm("identity"),
			CustName: c.PostForm("custname"),
			Sex:      intParam(c, "sex"),
			Address:  c.PostForm("address"),
			Phone:    c.PostForm("phone"),
			Career:   c.PostForm("career"),
		})
		if err != nil {
			c.JSON(http.StatusOK, UpdateError)
			return
		}
		c.JSON(http.StatusOK, UpdateSuccess)
	})
	rg.POST("/customer/deleteCustomer.action", func(c *gin.Context) {
		if err := d.BusService.DeleteCustomer(c.Request.Context(), c.PostForm("identity")); err != nil {
			c.JSON(http.StatusOK, DeleteError)
			return
		}
		c.JSON(http.StatusOK, DeleteSuccess)
	})
	rg.POST("/customer/deleteBatchCustomer.action", func(c *gin.Context) {
		for _, id := range stringArrayParam(c, "ids") {
			if err := d.BusService.DeleteCustomer(c.Request.Context(), id); err != nil {
				c.JSON(http.StatusOK, DeleteError)
				return
			}
		}
		c.JSON(http.StatusOK, DeleteSuccess)
	})

	rg.POST("/rent/checkCustomerExist.action", func(c *gin.Context) {
		x, err := d.BusService.GetCustomer(c.Request.Context(), c.PostForm("identity"))
		if err != nil || x == nil {
			c.JSON(http.StatusOK, StatusFalse)
			return
		}
		c.JSON(http.StatusOK, StatusTrue)
	})
	rg.GET("/rent/initRentFrom.action", func(c *gin.Context) {
		x, err := d.BusService.NewRentForm(c.Request.Context(), c.Query("identity"))
		if err != nil || x == nil {
			c.JSON(http.StatusOK, gin.H{})
			return
		}
		if u, ok := getSessionUser(c); ok {
			x.OperId = u.UserID
			x.OperName = u.RealName
		}
		c.JSON(http.StatusOK, x)
	})
	rg.POST("/rent/saveRent.action", func(c *gin.Context) {
		operid := 1
		if u, ok := getSessionUser(c); ok && u.Type != 1 {
			operid = u.UserID
		}
		rt := model.Rent{
			RentID:     c.PostForm("rentid"),
			Price:      floatParam(c, "price"),
			BeginDate:  timeParamValue(c, "begindate", time.Now()),
			ReturnDate: timePtrParam(c, "returndate"),
			RentFlag:   2,
			Identity:   c.PostForm("identity"),
			CarNumber:  c.PostForm("carnumber"),
			OperId:     operid,
			CreateTime: time.Now(),
		}
		if err := d.BusService.SaveRent(c.Request.Context(), rt); err != nil {
			c.JSON(http.StatusOK, AddErrorRent)
			return
		}
		c.JSON(http.StatusOK, AddSuccessRent)
	})
	rg.POST("/rent/deleteRent.action", func(c *gin.Context) {
		rentid := c.PostForm("rentid")
		if u, ok := getSessionUser(c); ok && u.Type != 1 {
			rent, err := d.BusService.GetRent(c.Request.Context(), rentid)
			if err != nil || rent == nil || rent.OperId != u.UserID {
				c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "无权限删除他人的出租单"})
				return
			}
		}
		if err := d.BusService.DeleteRent(c.Request.Context(), rentid); err != nil {
			c.JSON(http.StatusOK, DeleteError)
			return
		}
		c.JSON(http.StatusOK, DeleteSuccess)
	})
	rg.POST("/rent/updateRent.action", func(c *gin.Context) {
		rentid := c.PostForm("rentid")
		rent, err := d.BusService.GetRent(c.Request.Context(), rentid)
		if err != nil || rent == nil {
			c.JSON(http.StatusOK, UpdateError)
			return
		}
		if u, ok := getSessionUser(c); ok && u.Type != 1 {
			if rent.OperId != u.UserID {
				c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "无权限修改他人的出租单"})
				return
			}
		}
		rt := model.Rent{
			RentID:     rentid,
			Price:      floatParam(c, "price"),
			BeginDate:  timeParamValue(c, "begindate", time.Now()),
			ReturnDate: timePtrParam(c, "returndate"),
			RentFlag:   intParam(c, "rentflag"),
			Identity:   c.PostForm("identity"),
			CarNumber:  c.PostForm("carnumber"),
			OperId:     rent.OperId, // Preserve OperId
		}
		if err := d.BusService.UpdateRent(c.Request.Context(), rt); err != nil {
			c.JSON(http.StatusOK, UpdateError)
			return
		}
		c.JSON(http.StatusOK, UpdateSuccess)
	})
	rg.POST("/rent/checkRent.action", func(c *gin.Context) {
		rentid := c.PostForm("rentid")
		if u, ok := getSessionUser(c); ok && u.Type != 1 {
			rent, err := d.BusService.GetRent(c.Request.Context(), rentid)
			if err != nil || rent == nil || rent.OperId != u.UserID {
				c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "无权限操作他人的出租单"})
				return
			}
		}
		if err := d.BusService.CheckRent(c.Request.Context(), rentid, c.PostForm("carnumber")); err != nil {
			c.JSON(http.StatusOK, CheckError)
			return
		}
		c.JSON(http.StatusOK, CheckSuccess)
	})
	rg.GET("/rent/loadAllRent.action", func(c *gin.Context) {
		q := flatQuery(c)
		if u, ok := getSessionUser(c); ok && u.Type != 1 {
			q["operid"] = fmt.Sprint(u.UserID)
		}
		count, data, err := d.BusService.QueryRents(c.Request.Context(), q)
		if err != nil {
			fail(c, "查询出租单失败")
			return
		}
		c.JSON(http.StatusOK, NewPage(count, data))
	})

	rg.GET("/check/checkRentExist.action", func(c *gin.Context) {
		x, err := d.BusService.GetRent(c.Request.Context(), c.Query("rentid"))
		if err != nil || x == nil {
			c.JSON(http.StatusOK, nil)
			return
		}
		if u, ok := getSessionUser(c); ok && u.Type != 1 {
			if x.OperId != u.UserID {
				c.JSON(http.StatusOK, nil)
				return
			}
		}
		c.JSON(http.StatusOK, x)
	})
	rg.GET("/check/initCheckFormData.action", func(c *gin.Context) {
		data, err := d.BusService.InitCheckFormData(c.Request.Context(), c.Query("rentid"))
		if err != nil || data == nil {
			c.JSON(http.StatusOK, gin.H{})
			return
		}
		if u, ok := getSessionUser(c); ok {
			if u.Type != 1 {
				if rent, ok := data["rent"].(*model.Rent); ok && rent.OperId != u.UserID {
					c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "无权限操作他人的出租单"})
					return
				}
			}
			if check, ok := data["check"].(model.Check); ok {
				check.OperId = u.UserID
				check.OperName = u.RealName
				data["check"] = check
			}
		}
		c.JSON(http.StatusOK, data)
	})
	rg.POST("/check/saveCheck.action", func(c *gin.Context) {
		operid := 1
		if u, ok := getSessionUser(c); ok && u.Type != 1 {
			operid = u.UserID
		}
		x := model.Check{
			CheckID:    c.PostForm("checkid"),
			CheckDate:  timeParamValue(c, "checkdate", time.Now()),
			CheckDesc:  c.PostForm("checkdesc"),
			Problem:    c.PostForm("problem"),
			PayMoney:   floatParam(c, "paymoney"),
			OperId:     operid,
			RentID:     c.PostForm("rentid"),
			CreateTime: time.Now(),
		}
		if err := d.BusService.SaveCheck(c.Request.Context(), x); err != nil {
			c.JSON(http.StatusOK, AddError)
			return
		}
		c.JSON(http.StatusOK, AddSuccess)
	})
	rg.GET("/check/loadAllCheck.action", func(c *gin.Context) {
		q := flatQuery(c)
		if u, ok := getSessionUser(c); ok && u.Type != 1 {
			q["operid"] = fmt.Sprint(u.UserID)
		}
		count, data, err := d.BusService.QueryChecks(c.Request.Context(), q)
		if err != nil {
			fail(c, "查询检查单失败")
			return
		}
		c.JSON(http.StatusOK, NewPage(count, data))
	})
	rg.POST("/check/deleteCheck.action", func(c *gin.Context) {
		checkid := c.PostForm("checkid")
		if u, ok := getSessionUser(c); ok && u.Type != 1 {
			// check ownership
			count, chks, err := d.BusService.QueryChecks(c.Request.Context(), map[string]string{"checkid": checkid})
			if err != nil || count == 0 || chks[0].OperId != u.UserID {
				c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "无权限删除他人的检查单"})
				return
			}
		}
		if err := d.BusService.DeleteCheck(c.Request.Context(), checkid); err != nil {
			c.JSON(http.StatusOK, DeleteError)
			return
		}
		c.JSON(http.StatusOK, DeleteSuccess)
	})
	rg.POST("/check/deleteBatchCheck.action", func(c *gin.Context) {
		ids := stringArrayParam(c, "ids")
		if u, ok := getSessionUser(c); ok && u.Type != 1 {
			for _, id := range ids {
				count, chks, err := d.BusService.QueryChecks(c.Request.Context(), map[string]string{"checkid": id})
				if err != nil || count == 0 || chks[0].OperId != u.UserID {
					c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "无权限删除他人的检查单"})
					return
				}
			}
		}
		for _, id := range ids {
			if err := d.BusService.DeleteCheck(c.Request.Context(), id); err != nil {
				c.JSON(http.StatusOK, DeleteError)
				return
			}
		}
		c.JSON(http.StatusOK, DeleteSuccess)
	})
	rg.POST("/check/updateCheck.action", func(c *gin.Context) {
		checkid := c.PostForm("checkid")
		count, chks, err := d.BusService.QueryChecks(c.Request.Context(), map[string]string{"checkid": checkid})
		if err != nil || count == 0 {
			c.JSON(http.StatusOK, UpdateError)
			return
		}
		if u, ok := getSessionUser(c); ok && u.Type != 1 {
			if chks[0].OperId != u.UserID {
				c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "无权限修改他人的检查单"})
				return
			}
		}
		x := model.Check{
			CheckID:   checkid,
			CheckDate: timeParamValue(c, "checkdate", time.Now()),
			CheckDesc: c.PostForm("checkdesc"),
			Problem:   c.PostForm("problem"),
			PayMoney:  floatParam(c, "paymoney"),
			RentID:    c.PostForm("rentid"),
			OperId:    chks[0].OperId, // Preserve OperId
		}
		if err := d.BusService.UpdateCheck(c.Request.Context(), x); err != nil {
			c.JSON(http.StatusOK, UpdateError)
			return
		}
		c.JSON(http.StatusOK, UpdateSuccess)
	})

	rg.GET("/franchisee/loadAllFranchisee.action", func(c *gin.Context) {
		count, data, err := d.BusService.QueryFranchisees(c.Request.Context(), flatQuery(c))
		if err != nil {
			fail(c, "查询加盟商失败")
			return
		}
		c.JSON(http.StatusOK, NewPage(count, data))
	})
	rg.POST("/franchisee/addFranchisee.action", func(c *gin.Context) {
		if err := d.BusService.AddFranchisee(c.Request.Context(), model.Franchisee{Name: c.PostForm("name"), Phone: c.PostForm("phone")}); err != nil {
			c.JSON(http.StatusOK, AddError)
			return
		}
		c.JSON(http.StatusOK, AddSuccess)
	})
	rg.POST("/franchisee/updateFranchisee.action", func(c *gin.Context) {
		if err := d.BusService.UpdateFranchisee(c.Request.Context(), model.Franchisee{ID: intParam(c, "id"), Name: c.PostForm("name"), Phone: c.PostForm("phone")}); err != nil {
			c.JSON(http.StatusOK, UpdateError)
			return
		}
		c.JSON(http.StatusOK, UpdateSuccess)
	})
	rg.POST("/franchisee/deleteFranchisee.action", func(c *gin.Context) {
		if err := d.BusService.DeleteFranchisee(c.Request.Context(), intParam(c, "id")); err != nil {
			c.JSON(http.StatusOK, DeleteError)
			return
		}
		c.JSON(http.StatusOK, DeleteSuccess)
	})
}

func registerStatRoutes(rg *gin.RouterGroup, d Deps) {
	rg.GET("/stat/loadCustomerAreaStatJson.action", func(c *gin.Context) {
		data, err := d.StatService.LoadCustomerAreaStat(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusOK, []any{})
			return
		}
		c.JSON(http.StatusOK, data)
	})
	rg.GET("/stat/loadCustomerAreaSexStatJson.action", func(c *gin.Context) {
		data, err := d.StatService.LoadCustomerAreaSexStat(c.Request.Context(), c.Query("area"))
		if err != nil {
			c.JSON(http.StatusOK, []any{})
			return
		}
		c.JSON(http.StatusOK, data)
	})
	rg.GET("/stat/loadOpernameYearGradeStatJson.action", func(c *gin.Context) {
		data, err := d.StatService.LoadOpernameYearGradeStat(c.Request.Context(), c.Query("year"))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"name": []string{}, "value": []float64{}})
			return
		}
		names := make([]string, 0, len(data))
		values := make([]float64, 0, len(data))
		for _, x := range data {
			names = append(names, x.Name)
			values = append(values, x.Value)
		}
		c.JSON(http.StatusOK, gin.H{"name": names, "value": values})
	})
	rg.GET("/stat/loadCompanyYearGradeStatJson.action", func(c *gin.Context) {
		data, err := d.StatService.LoadCompanyYearGradeStat(c.Request.Context(), c.Query("year"))
		if err != nil {
			c.JSON(http.StatusOK, []float64{})
			return
		}
		c.JSON(http.StatusOK, data)
	})
	rg.GET("/stat/exportCustomer.action", func(c *gin.Context) {
		data, err := d.BusService.ListCustomers(c.Request.Context(), flatQuery(c))
		if err != nil {
			c.String(http.StatusInternalServerError, "export error")
			return
		}
		content, err := service.ExportCustomersAsXLS(data, "客户数据")
		if err != nil {
			c.String(http.StatusInternalServerError, "export error")
			return
		}
		downloadBytes(c, content, "客户数据.xlsx")
	})
	rg.GET("/stat/exportRent.action", func(c *gin.Context) {
		rent, err := d.BusService.GetRent(c.Request.Context(), c.Query("rentid"))
		if err != nil || rent == nil {
			c.String(http.StatusNotFound, "rent not found")
			return
		}
		customer, _ := d.BusService.GetCustomer(c.Request.Context(), rent.Identity)
		name := "出租单.xls"
		if customer != nil {
			name = customer.CustName + "-的出租单.xls"
		}
		content, err := service.ExportRentAsXLS(rent, customer, "出租单")
		if err != nil {
			c.String(http.StatusInternalServerError, "export error")
			return
		}
		if strings.HasSuffix(name, ".xls") {
			name = strings.TrimSuffix(name, ".xls") + ".xlsx"
		}
		downloadBytes(c, content, name)
	})
}

func registerFileRoutes(rg *gin.RouterGroup, d Deps) {
	upload := func(c *gin.Context) {
		fh, err := c.FormFile("mf")
		if err != nil {
			fh, err = c.FormFile("file")
		}
		if err != nil {
			fail(c, "上传失败")
			return
		}
		dirName := service.DateDir(time.Now())
		dirPath := filepath.Join(d.FileService.UploadRoot, dirName)
		if err := service.EnsureDir(dirPath); err != nil {
			fail(c, "上传失败")
			return
		}
		newName := newUploadName(fh.Filename, true)
		rel := filepath.ToSlash(filepath.Join(dirName, newName))
		abs := filepath.Join(dirPath, newName)
		if err := c.SaveUploadedFile(fh, abs); err != nil {
			fail(c, "上传失败")
			return
		}
		c.JSON(http.StatusOK, NewData(gin.H{"src": rel}))
	}
	rg.POST("/file/uploadFile.action", upload)
	rg.POST("/file/uploadImage.action", upload)
	rg.GET("/file/downloadShowFile.action", func(c *gin.Context) {
		serveUploadedFile(c, d, c.Query("path"), false)
	})
	rg.GET("/file/downloadFile.action", func(c *gin.Context) {
		serveUploadedFile(c, d, c.Query("path"), true)
	})
}

func flatQuery(c *gin.Context) map[string]string {
	out := map[string]string{}
	for k, vals := range c.Request.URL.Query() {
		if len(vals) > 0 {
			out[k] = vals[0]
		}
	}
	return out
}

func ptrString(s string) *string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return &s
}
func ptrIntMaybe(s string) *int {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	if n, err := strconv.Atoi(s); err == nil {
		return &n
	}
	return nil
}
func floatParam(c *gin.Context, key string) float64 {
	v := strings.TrimSpace(c.PostForm(key))
	if v == "" {
		v = strings.TrimSpace(c.Query(key))
	}
	if n, err := strconv.ParseFloat(v, 64); err == nil {
		return n
	}
	return 0
}
func timePtrParam(c *gin.Context, key string) *time.Time {
	v := strings.TrimSpace(c.PostForm(key))
	if v == "" {
		v = strings.TrimSpace(c.Query(key))
	}
	t, err := service.ParseTimeFlexible(v)
	if err != nil || t == nil {
		return nil
	}
	return t
}
func timeParamValue(c *gin.Context, key string, def time.Time) time.Time {
	if t := timePtrParam(c, key); t != nil {
		return *t
	}
	return def
}
func newUploadName(old string, temp bool) string {
	ext := filepath.Ext(old)
	base := strings.TrimSuffix(filepath.Base(old), ext)
	name := fmt.Sprintf("%s_%s", time.Now().Format("20060102150405.000000"), base)
	name = strings.NewReplacer(" ", "", "/", "", "\\", "").Replace(name)
	if temp {
		return name + ext + uploadTempTag
	}
	return name + ext
}
func serveUploadedFile(c *gin.Context, d Deps, rel string, attachment bool) {
	abs, err := service.SafeJoin(d.FileService.UploadRoot, rel)
	if err != nil {
		c.String(http.StatusBadRequest, "文件不存在")
		return
	}
	f, err := os.Open(abs)
	if err != nil {
		c.String(http.StatusNotFound, "文件不存在")
		return
	}
	defer f.Close()
	st, _ := f.Stat()
	if attachment {
		name := filepath.Base(abs)
		setAttachmentFilename(c, name)
	}
	if ct := mime.TypeByExtension(filepath.Ext(abs)); ct != "" {
		c.Header("Content-Type", ct)
	}
	if st != nil {
		c.Header("Content-Length", fmt.Sprintf("%d", st.Size()))
	}
	_, _ = io.Copy(c.Writer, f)
}
func downloadBytes(c *gin.Context, body []byte, filename string) {
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	setAttachmentFilename(c, filename)
	c.Header("X-Content-Type-Options", "nosniff")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", body)
}

func setAttachmentFilename(c *gin.Context, filename string) {
	filename = strings.TrimSpace(filename)
	if filename == "" {
		filename = "download"
	}
	ascii := make([]rune, 0, len(filename))
	for _, r := range filename {
		if r >= 32 && r < 127 && r != '"' && r != '\\' {
			ascii = append(ascii, r)
		} else {
			ascii = append(ascii, '_')
		}
	}
	fallback := string(ascii)
	encoded := url.PathEscape(filename)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"; filename*=UTF-8''%s", fallback, encoded))
}
