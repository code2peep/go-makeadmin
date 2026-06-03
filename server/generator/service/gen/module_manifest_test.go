package gen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go-makeadmin/config"
	"go-makeadmin/generator/schemas/req"
)

func TestPreviewModuleManifestFromInlineJSON(t *testing.T) {
	body := `{
  "version": 1,
  "module": "article",
  "entity": "DemoArticle",
  "table": "ma_demo_article",
  "backendPackage": "gencode",
  "backend": {
    "routes": [
      {"method": "GET", "path": "/article/list", "permission": "article:list"},
      {"method": "GET", "path": "/article/detail", "permission": "article:detail"},
      {"method": "POST", "path": "/article/add", "permission": "article:add"},
      {"method": "POST", "path": "/article/edit", "permission": "article:edit"},
      {"method": "POST", "path": "/article/del", "permission": "article:del"}
    ]
  },
  "frontend": {
    "api": "admin/src/api/article.ts",
    "views": ["admin/src/views/article/index.vue", "admin/src/views/article/edit.vue"]
  },
  "menu": {
    "code": "demo.article",
    "parent": "dev_tools",
    "type": "page",
    "name": "Demo Article",
    "routePath": "/demo/article",
    "routeName": "demo.article",
    "component": "article/index",
    "permission": "article:list",
    "visible": false,
    "sort": 10
  },
  "permissions": [
    {"code": "article:list", "name": "Article list", "module": "article", "resource": "article", "action": "list"},
    {"code": "article:detail", "name": "Article detail", "module": "article", "resource": "article", "action": "detail"},
    {"code": "article:add", "name": "Article add", "module": "article", "resource": "article", "action": "add"},
    {"code": "article:edit", "name": "Article edit", "module": "article", "resource": "article", "action": "edit"},
    {"code": "article:del", "name": "Article delete", "module": "article", "resource": "article", "action": "del"}
  ],
  "runtimeRegistered": false,
  "requiresSchema": false,
  "codegen": {
    "columns": [
      {"columnName": "title", "goField": "title", "htmlType": "input", "isRequired": 1},
      {"columnName": "status", "goField": "status", "htmlType": "number", "goType": "int", "queryType": "="}
    ]
  }
}`
	srv := generateService{}
	res, err := srv.PreviewModuleManifest(req.ModuleManifestPreviewReq{ManifestBody: body, TenantID: 7, RoleID: 3, AuthorName: "tester"})
	if err != nil {
		t.Fatalf("preview inline manifest: %v", err)
	}
	if res.Source != "inline" || res.Manifest.Module != "article" || res.Manifest.Entity != "DemoArticle" {
		t.Fatalf("unexpected manifest summary: %+v source=%s", res.Manifest, res.Source)
	}
	if res.Detail.Base.TableName != "ma_demo_article" || res.Detail.Gen.ModuleName != "article" {
		t.Fatalf("unexpected detail: %+v", res.Detail)
	}
	if len(res.Detail.Column) != 3 {
		t.Fatalf("columns length = %d, want 3", len(res.Detail.Column))
	}
	assertContains(t, res.Code["gocode/model.go"], "type DemoArticle struct")
	assertContains(t, strings.Join(strings.Fields(res.Code["gocode/model.go"]), " "), "Title string")
	assertContains(t, strings.Join(strings.Fields(res.Code["gocode/model.go"]), " "), "Status int")
	assertContains(t, res.Code["gocode/route.go"], `rg.GET("/article/list"`)
	assertContains(t, res.Code["vue/api.ts"], "url: '/article/list'")
	assertContains(t, res.Code["vue/index.vue"], `v-perms="['article:add']"`)
	if res.Plan.TenantID != 7 || res.Plan.RoleID != 3 {
		t.Fatalf("unexpected plan tenant/role: %+v", res.Plan)
	}
	assertContains(t, res.Plan.RegistrySQL, "INSERT INTO `ma_permission`")
	assertContains(t, res.Plan.RegistrySQL, "`route_name` = 'demo.article'")
	assertContains(t, res.Plan.RoleGrantSQL, "INSERT INTO `ma_role_permission`")
	assertContains(t, res.Plan.RoleGrantSQL, "SET @tenant_id = 7;")
	assertContains(t, res.Plan.RoleGrantSQL, "SET @role_id = 3;")
	assertContains(t, res.Plan.InstallSQL, res.Plan.RegistrySQL)
	assertContains(t, res.Plan.InstallSQL, res.Plan.RoleGrantSQL)
	assertContains(t, res.Plan.UninstallSQL, "DELETE FROM `ma_permission`")
	assertContains(t, res.Plan.UninstallSQL, "'article:add'")
	assertContains(t, res.Plan.RuntimeHint, "MAKEADMIN_ENABLE_DEMO_MODULE=1")
}

func TestPreviewModuleManifestFromRepositoryPath(t *testing.T) {
	manifestPath := filepath.Join(filepath.Dir(config.Config.RootPath), "examples", "demo", "manifest.json")
	if _, err := os.Stat(manifestPath); err != nil {
		t.Fatalf("stat demo manifest: %v", err)
	}

	srv := generateService{}
	res, err := srv.PreviewModuleManifest(req.ModuleManifestPreviewReq{ManifestPath: "examples/demo/manifest.json"})
	if err != nil {
		t.Fatalf("preview repository manifest: %v", err)
	}
	if res.Source != "examples/demo/manifest.json" {
		t.Fatalf("source = %q, want examples/demo/manifest.json", res.Source)
	}
	if strings.TrimSpace(res.Warning) == "" {
		t.Fatalf("warning must not be empty")
	}
	if res.Plan.RoleID != 1 {
		t.Fatalf("default role id = %d, want 1", res.Plan.RoleID)
	}
	assertContains(t, res.Code["vue/edit.vue"], "articleAdd")
}

func TestPreviewModuleManifestIncludesInstallPlan(t *testing.T) {
	srv := generateService{}
	res, err := srv.PreviewModuleManifest(req.ModuleManifestPreviewReq{
		ManifestPath: "examples/demo/manifest.json",
		TenantID:     0,
		RoleID:       2,
	})
	if err != nil {
		t.Fatalf("preview repository manifest: %v", err)
	}
	if res.Plan.TenantID != 0 || res.Plan.RoleID != 2 {
		t.Fatalf("unexpected plan ids: %+v", res.Plan)
	}
	assertContains(t, res.Plan.RegistrySQL, "SET @parent_route_name = 'dev_tools';")
	assertContains(t, res.Plan.RegistrySQL, "INSERT INTO `ma_menu_permission`")
	assertContains(t, res.Plan.RoleGrantSQL, "SET @role_id = 2;")
	assertContains(t, res.Plan.InstallSQL, "INSERT INTO `ma_permission`")
	assertContains(t, res.Plan.InstallSQL, "INSERT INTO `ma_role_permission`")
	assertContains(t, res.Plan.UninstallSQL, "DELETE rp FROM `ma_role_permission`")
	assertContains(t, res.Plan.UninstallSQL, "DELETE FROM `ma_menu`")
	assertContains(t, res.Plan.RuntimeHint, "MAKEADMIN_ENABLE_DEMO_MODULE=1")
}

func TestPreviewModuleManifestRejectsUnsafePath(t *testing.T) {
	srv := generateService{}
	if _, err := srv.PreviewModuleManifest(req.ModuleManifestPreviewReq{ManifestPath: "../manifest.json"}); err == nil {
		t.Fatalf("expected unsafe manifest path to fail")
	}
}

func TestModuleManifestInstallApplyGateRequiresEnv(t *testing.T) {
	srv := generateService{}
	res, err := srv.ApplyModuleManifestInstall(req.ModuleManifestInstallApplyReq{
		ManifestPath: "examples/demo/manifest.json",
	})
	if err == nil {
		t.Fatalf("expected missing env gate to fail")
	}
	assertContains(t, err.Error(), moduleManifestInstallApplyEnv+"=1 is required")
	assertContains(t, err.Error(), "no database access was attempted")
	if res.Status != "blocked" || res.RequiredEnv != moduleManifestInstallApplyEnv {
		t.Fatalf("unexpected gate response: %+v", res)
	}
	if len(res.Checks) != 1 || res.Checks[0].Status != "failed" {
		t.Fatalf("unexpected gate checks: %+v", res.Checks)
	}
}

func TestModuleManifestInstallApplyGateRequiresConfirmations(t *testing.T) {
	t.Setenv(moduleManifestInstallApplyEnv, "1")
	srv := generateService{}
	tenantID := uint64(0)
	roleID := uint64(1)

	res, err := srv.ApplyModuleManifestInstall(req.ModuleManifestInstallApplyReq{
		ManifestPath:  "examples/demo/manifest.json",
		ConfirmModule: "wrong",
	})
	if err == nil {
		t.Fatalf("expected confirmModule gate to fail")
	}
	assertContains(t, err.Error(), `confirmModule must be "article"`)
	assertContains(t, err.Error(), "no database access was attempted")
	if res.Checks[len(res.Checks)-1].Name != "confirmModule" {
		t.Fatalf("expected confirmModule check, got %+v", res.Checks)
	}

	res, err = srv.ApplyModuleManifestInstall(req.ModuleManifestInstallApplyReq{
		ManifestPath:    "examples/demo/manifest.json",
		ConfirmModule:   "article",
		ConfirmTenantID: uint64Ptr(9),
	})
	if err == nil {
		t.Fatalf("expected confirmTenantId gate to fail")
	}
	assertContains(t, err.Error(), "confirmTenantId must be 0")
	assertContains(t, err.Error(), "no database access was attempted")

	res, err = srv.ApplyModuleManifestInstall(req.ModuleManifestInstallApplyReq{
		ManifestPath:    "examples/demo/manifest.json",
		ConfirmModule:   "article",
		ConfirmTenantID: &tenantID,
		ConfirmRoleID:   uint64Ptr(9),
	})
	if err == nil {
		t.Fatalf("expected confirmRoleId gate to fail")
	}
	assertContains(t, err.Error(), "confirmRoleId must be 1")
	assertContains(t, err.Error(), "no database access was attempted")

	res, err = srv.ApplyModuleManifestInstall(req.ModuleManifestInstallApplyReq{
		ManifestPath:    "examples/demo/manifest.json",
		ConfirmModule:   "article",
		ConfirmTenantID: &tenantID,
		ConfirmRoleID:   &roleID,
	})
	if err == nil {
		t.Fatalf("expected confirmInstall gate to fail")
	}
	assertContains(t, err.Error(), "confirmInstall must be true")
	assertContains(t, err.Error(), "no database access was attempted")

	res, err = srv.ApplyModuleManifestInstall(req.ModuleManifestInstallApplyReq{
		ManifestBody:    moduleManifestBodyWithRequiresSchema(t),
		ConfirmModule:   "article",
		ConfirmTenantID: &tenantID,
		ConfirmRoleID:   &roleID,
		ConfirmInstall:  true,
	})
	if err == nil {
		t.Fatalf("expected confirmSchemaRisk gate to fail")
	}
	assertContains(t, err.Error(), "confirmSchemaRisk must be true")
	assertContains(t, err.Error(), "no database access was attempted")
}

func TestModuleManifestInstallApplyGateRequiresDatabaseWhenConfirmed(t *testing.T) {
	t.Setenv(moduleManifestInstallApplyEnv, "1")
	srv := generateService{}
	tenantID := uint64(0)
	roleID := uint64(1)

	res, err := srv.ApplyModuleManifestInstall(req.ModuleManifestInstallApplyReq{
		ManifestPath:    "examples/demo/manifest.json",
		ConfirmModule:   "article",
		ConfirmTenantID: &tenantID,
		ConfirmRoleID:   &roleID,
		ConfirmInstall:  true,
	})
	if err == nil {
		t.Fatalf("expected missing database to fail")
	}
	assertContains(t, err.Error(), "module install apply requires configured database")
	assertContains(t, err.Error(), "no database access was attempted")
	if res.Manifest.Module != "article" || res.Plan.InstallSQL == "" {
		t.Fatalf("unexpected install gate response: %+v", res)
	}
	if res.Checks[len(res.Checks)-1].Name != "database" || res.Checks[len(res.Checks)-1].Status != "failed" {
		t.Fatalf("unexpected database check: %+v", res.Checks)
	}
}

func moduleManifestBodyWithRequiresSchema(t *testing.T) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(filepath.Dir(config.Config.RootPath), "examples", "demo", "manifest.json"))
	if err != nil {
		t.Fatalf("read demo manifest: %v", err)
	}
	return strings.Replace(string(content), `"requiresSchema": false`, `"requiresSchema": true`, 1)
}

func uint64Ptr(value uint64) *uint64 {
	return &value
}
