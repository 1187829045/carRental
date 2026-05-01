package handler

import (
	"testing"

	"carRental/internal/config"
)

func TestRouterHasMenuRoutesUnderBasePath(t *testing.T) {
	r := NewRouter(Deps{Cfg: config.Config{BasePath: "/carRental", StaticDir: "./src/main/webapp/static", ViewRoot: "./src/main/webapp/WEB-INF/view"}})
	paths := make(map[string]bool)
	for _, rt := range r.Routes() {
		paths[rt.Path] = true
	}
	require := []string{
		"/carRental/sys/toMenuManager.action",
		"/carRental/sys/toMenuLeft.action",
		"/carRental/sys/toMenuRight.action",
		"/carRental/menu/loadMenuManagerLeftTreeJson.action",
		"/carRental/menu/loadAllMenu.action",
		"/carRental/menu/addMenu.action",
		"/carRental/menu/updateMenu.action",
		"/carRental/menu/deleteMenu.action",
		"/carRental/menu/checkMenuHasChildren.action",
		"/carRental/static/*filepath",
	}
	for _, p := range require {
		if !paths[p] {
			t.Fatalf("missing route %s", p)
		}
	}
}

