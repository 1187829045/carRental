package config

import (
	"bufio"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	mysql "github.com/go-sql-driver/mysql"
)

type Config struct {
	Addr      string // e.g. ":8080"
	BasePath  string // corresponds to JSP `${yeqifu}`; e.g. "/carRental" for legacy deployment
	ViewRoot  string // JSP root, e.g. "src/main/webapp/WEB-INF/view"
	StaticDir string // static root to serve at /static, e.g. "src/main/webapp/static"

	MySQLDSN string // root:pwd@tcp(host:port)/db?parseTime=true&loc=Local&charset=utf8

	MonitorUser string // druid replacement basic auth user
	MonitorPass string // druid replacement basic auth pass
}

func Load(projectRoot string) (Config, error) {
	cfg := Config{
		Addr:      envOr("ADDR", ":8080"),
		BasePath:  envOr("BASE_PATH", "/carRental"),
		ViewRoot:  filepath.Join(projectRoot, "src/main/webapp/WEB-INF/view"),
		StaticDir: filepath.Join(projectRoot, "src/main/webapp/static"),
	}

	propsPath := filepath.Join(projectRoot, "src/main/resources/db.properties")
	props, err := readProperties(propsPath)
	if err != nil {
		return Config{}, fmt.Errorf("read db.properties: %w", err)
	}

	user := strings.TrimSpace(props["jdbc.user"])
	pass := strings.TrimSpace(props["jdbc.password"])
	jdbcURL := strings.TrimSpace(props["url"])
	if user == "" || jdbcURL == "" {
		return Config{}, errors.New("db.properties missing jdbc.user or url")
	}

	dsn, err := jdbcToMySQLDSN(user, pass, jdbcURL)
	if err != nil {
		return Config{}, fmt.Errorf("parse jdbc url: %w", err)
	}
	cfg.MySQLDSN = envOr("MYSQL_DSN", dsn)
	cfg.MonitorUser = envOr("MONITOR_USER", user)
	cfg.MonitorPass = envOr("MONITOR_PASS", pass)
	return cfg, nil
}

func envOr(k, def string) string {
	if v := strings.TrimSpace(os.Getenv(k)); v != "" {
		return v
	}
	return def
}

func readProperties(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	m := make(map[string]string)
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		kv := strings.SplitN(line, "=", 2)
		if len(kv) != 2 {
			continue
		}
		m[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return m, nil
}

// Converts:
// jdbc:mysql://localhost:3306/carRental?useUnicode=true&characterEncoding=UTF-8&serverTimezone=UTC&useSSL=false
// to:
// user:pass@tcp(localhost:3306)/carRental?parseTime=true&loc=Local&charset=utf8
func jdbcToMySQLDSN(user, pass, jdbcURL string) (string, error) {
	const prefix = "jdbc:mysql://"
	if !strings.HasPrefix(jdbcURL, prefix) {
		return "", fmt.Errorf("unsupported jdbc url: %q", jdbcURL)
	}
	raw := strings.TrimPrefix(jdbcURL, "jdbc:")
	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}

	host := u.Host
	dbName := strings.TrimPrefix(u.Path, "/")
	if host == "" || dbName == "" {
		return "", fmt.Errorf("invalid jdbc url host/path: %q", jdbcURL)
	}

	q := u.Query()
	// Preserve charset intent; default to utf8 to match old schema.
	charset := "utf8"
	if enc := q.Get("characterEncoding"); enc != "" {
		if strings.EqualFold(enc, "utf-8") || strings.EqualFold(enc, "utf8") {
			charset = "utf8"
		}
	}

	// The legacy app implicitly parses DATETIME to java.util.Date.
	// In Go, parseTime is needed to scan time.Time.
	params := url.Values{}
	params.Set("parseTime", "true")
	params.Set("loc", "Local")
	params.Set("charset", charset)

	cfg := mysql.NewConfig()
	cfg.User = user
	cfg.Passwd = pass
	cfg.Net = "tcp"
	cfg.Addr = host
	cfg.DBName = dbName
	cfg.ParseTime = true
	cfg.Loc = time.Local
	cfg.Params = map[string]string{"charset": params.Get("charset")}
	for k, vs := range q {
		if len(vs) == 0 {
			continue
		}
		if k == "useUnicode" || k == "characterEncoding" || k == "serverTimezone" {
			continue
		}
		cfg.Params[k] = vs[0]
	}
	return cfg.FormatDSN(), nil
}
