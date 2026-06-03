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
	res, err := srv.PreviewModuleManifest(req.ModuleManifestPreviewReq{ManifestBody: body, TenantID: 7, AuthorName: "tester"})
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
	assertContains(t, res.Code["vue/edit.vue"], "articleAdd")
}

func TestPreviewModuleManifestRejectsUnsafePath(t *testing.T) {
	srv := generateService{}
	if _, err := srv.PreviewModuleManifest(req.ModuleManifestPreviewReq{ManifestPath: "../manifest.json"}); err == nil {
		t.Fatalf("expected unsafe manifest path to fail")
	}
}
