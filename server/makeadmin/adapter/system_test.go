package adapter

import (
	"testing"

	makeadminsvc "go-makeadmin/makeadmin/service"
	"go-makeadmin/model/makeadmin"
	"go-makeadmin/util"
)

func TestRouteMenuMapsBuildsNestedChildren(t *testing.T) {
	menus := []makeadminsvc.RouteMenu{
		{ID: 500, MenuType: makeadmin.MenuTypeCatalog, Name: "系统设置", RoutePath: "/setting"},
		{ID: 501, ParentID: 500, MenuType: makeadmin.MenuTypeCatalog, Name: "网站设置", RoutePath: "/setting/website"},
		{ID: 502, ParentID: 501, MenuType: makeadmin.MenuTypePage, Name: "网站信息", RoutePath: "/setting/website/information"},
	}

	tree := util.ArrayUtil.ListToTree(routeMenuMaps(menus), "id", "pid", "children")
	if len(tree) != 1 {
		t.Fatalf("tree length = %d, want 1", len(tree))
	}

	root := tree[0].(map[string]interface{})
	children := root["children"].([]interface{})
	if len(children) != 1 {
		t.Fatalf("root children length = %d, want 1", len(children))
	}

	website := children[0].(map[string]interface{})
	if website["paths"] != "website" {
		t.Fatalf("website paths = %q, want website", website["paths"])
	}

	pages := website["children"].([]interface{})
	information := pages[0].(map[string]interface{})
	if information["paths"] != "information" {
		t.Fatalf("information paths = %q, want information", information["paths"])
	}
}
