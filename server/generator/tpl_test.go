package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"go-makeadmin/config"
	"go-makeadmin/model/gen"
)

func TestGeneratedCrudGoCodeCompiles(t *testing.T) {
	table, _, vars := demoCrudFixture()
	tplCodeMap := map[string]string{}
	for _, tplPath := range []string{
		"gocode/model.go.tpl",
		"gocode/schema.go.tpl",
		"gocode/service.go.tpl",
		"gocode/route.go.tpl",
	} {
		code, err := TemplateUtil.Render(tplPath, vars)
		if err != nil {
			t.Fatalf("render %s: %v", tplPath, err)
		}
		tplCodeMap[tplPath] = code
	}

	basePath := filepath.Join(config.Config.RootPath, ".tmp-codegen-compile")
	if err := os.RemoveAll(basePath); err != nil {
		t.Fatalf("remove old temp dir: %v", err)
	}
	defer os.RemoveAll(basePath)

	if err := TemplateUtil.GenCodeFiles(tplCodeMap, table.ModuleName, basePath); err != nil {
		t.Fatalf("write generated code: %v", err)
	}

	packageDir := filepath.Join(basePath, "gocode", table.ModuleName)
	cmd := exec.Command("go", "test", ".")
	cmd.Dir = packageDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("generated package does not compile:\n%s", string(output))
	}
}

func TestGeneratedCrudFrontendCodeTypeChecks(t *testing.T) {
	if os.Getenv("MAKEADMIN_CODEGEN_FRONTEND_CHECK") != "1" {
		t.Skip("set MAKEADMIN_CODEGEN_FRONTEND_CHECK=1 to type-check generated frontend code")
	}
	table, _, vars := demoCrudFixture()
	repoRoot := filepath.Dir(config.Config.RootPath)
	adminRoot := filepath.Join(repoRoot, "admin")
	apiPath := filepath.Join(adminRoot, "src", "api", table.ModuleName+".ts")
	viewDir := filepath.Join(adminRoot, "src", "views", table.ModuleName)
	if exists(apiPath) || exists(viewDir) {
		t.Fatalf("refusing to overwrite existing generated frontend fixture for module %s", table.ModuleName)
	}
	defer os.Remove(apiPath)
	defer os.RemoveAll(viewDir)

	renders := map[string]string{
		apiPath:                             "vue/api.ts.tpl",
		filepath.Join(viewDir, "index.vue"): "vue/index.vue.tpl",
		filepath.Join(viewDir, "edit.vue"):  "vue/edit.vue.tpl",
	}
	for target, tplPath := range renders {
		code, err := TemplateUtil.Render(tplPath, vars)
		if err != nil {
			t.Fatalf("render %s: %v", tplPath, err)
		}
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			t.Fatalf("mkdir %s: %v", filepath.Dir(target), err)
		}
		if err := os.WriteFile(target, []byte(code), 0644); err != nil {
			t.Fatalf("write %s: %v", target, err)
		}
	}

	cmd := exec.Command("npm", "run", "type-check")
	cmd.Dir = adminRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("generated frontend code does not type-check:\n%s", string(output))
	}
}

func TestGeneratedCrudCodeMatchesModuleManifest(t *testing.T) {
	manifestPath := os.Getenv("MAKEADMIN_CODEGEN_MANIFEST")
	if manifestPath == "" {
		t.Skip("set MAKEADMIN_CODEGEN_MANIFEST to verify a scaffold manifest")
	}
	table, _, vars, manifest := manifestCrudFixture(t, manifestPath)

	rendered := map[string]string{}
	for _, tplPath := range []string{
		"gocode/model.go.tpl",
		"gocode/schema.go.tpl",
		"gocode/service.go.tpl",
		"gocode/route.go.tpl",
		"vue/api.ts.tpl",
		"vue/index.vue.tpl",
		"vue/edit.vue.tpl",
	} {
		code, err := TemplateUtil.Render(tplPath, vars)
		if err != nil {
			t.Fatalf("render %s: %v", tplPath, err)
		}
		rendered[tplPath] = code
	}

	assertManifestRoutes(t, rendered["gocode/route.go.tpl"], rendered["vue/api.ts.tpl"], manifest)
	assertManifestPermissions(t, rendered["vue/index.vue.tpl"], manifest)

	tplCodeMap := map[string]string{
		"gocode/model.go.tpl":   rendered["gocode/model.go.tpl"],
		"gocode/schema.go.tpl":  rendered["gocode/schema.go.tpl"],
		"gocode/service.go.tpl": rendered["gocode/service.go.tpl"],
		"gocode/route.go.tpl":   rendered["gocode/route.go.tpl"],
	}
	basePath := filepath.Join(config.Config.RootPath, ".tmp-codegen-manifest")
	if err := os.RemoveAll(basePath); err != nil {
		t.Fatalf("remove old temp dir: %v", err)
	}
	defer os.RemoveAll(basePath)

	if err := TemplateUtil.GenCodeFiles(tplCodeMap, table.ModuleName, basePath); err != nil {
		t.Fatalf("write generated code: %v", err)
	}

	packageDir := filepath.Join(basePath, "gocode", table.ModuleName)
	cmd := exec.Command("go", "test", ".")
	cmd.Dir = packageDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("generated manifest package does not compile:\n%s", string(output))
	}
}

func demoCrudFixture() (gen.GenTable, []gen.GenTableColumn, TplVars) {
	return crudFixture("ma_demo_article", "DemoArticle", "article", "Demo article")
}

func crudFixture(tableName string, entityName string, moduleName string, functionName string) (gen.GenTable, []gen.GenTableColumn, TplVars) {
	table := gen.GenTable{
		ID:           1,
		TableName:    tableName,
		TableComment: functionName,
		AuthorName:   "codepeep",
		EntityName:   entityName,
		ModuleName:   moduleName,
		FunctionName: functionName,
		GenTpl:       GenConstants.TplCrud,
		GenPath:      "/",
	}
	columns := []gen.GenTableColumn{
		{
			ID:            1,
			TableID:       1,
			ColumnName:    "id",
			ColumnComment: "ID",
			JavaType:      GoConstants.TypeInt,
			JavaField:     "id",
			IsPk:          1,
			IsList:        1,
			Sort:          1,
		},
		{
			ID:            2,
			TableID:       1,
			ColumnName:    "title",
			ColumnComment: "Title",
			JavaType:      GoConstants.TypeString,
			JavaField:     "title",
			IsRequired:    1,
			IsInsert:      1,
			IsEdit:        1,
			IsList:        1,
			IsQuery:       1,
			QueryType:     GenConstants.QueryLike,
			HtmlType:      HtmlConstants.HtmlInput,
			Sort:          2,
		},
		{
			ID:            3,
			TableID:       1,
			ColumnName:    "status",
			ColumnComment: "Status",
			JavaType:      GoConstants.TypeInt,
			JavaField:     "status",
			IsInsert:      1,
			IsEdit:        1,
			IsList:        1,
			IsQuery:       1,
			QueryType:     GenConstants.QueryEq,
			HtmlType:      HtmlConstants.HtmlInput,
			Sort:          3,
		},
	}
	vars := TemplateUtil.PrepareVars(table, columns, gen.GenTableColumn{}, nil)
	return table, columns, vars
}

type moduleManifest struct {
	Module      string `json:"module"`
	Entity      string `json:"entity"`
	Table       string `json:"table"`
	Menu        struct {
		Name string `json:"name"`
	} `json:"menu"`
	Backend struct {
		Routes []struct {
			Method     string `json:"method"`
			Path       string `json:"path"`
			Permission string `json:"permission"`
		} `json:"routes"`
	} `json:"backend"`
	Permissions []struct {
		Code   string `json:"code"`
		Action string `json:"action"`
	} `json:"permissions"`
}

func manifestCrudFixture(t *testing.T, path string) (gen.GenTable, []gen.GenTableColumn, TplVars, moduleManifest) {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	var manifest moduleManifest
	if err := json.Unmarshal(content, &manifest); err != nil {
		t.Fatalf("decode manifest: %v", err)
	}
	if manifest.Module == "" || manifest.Entity == "" || manifest.Table == "" {
		t.Fatalf("manifest must include module, entity and table")
	}
	functionName := manifest.Menu.Name
	if functionName == "" {
		functionName = manifest.Module
	}
	table, columns, vars := crudFixture(manifest.Table, manifest.Entity, manifest.Module, functionName)
	return table, columns, vars, manifest
}

func assertManifestRoutes(t *testing.T, routeCode string, apiCode string, manifest moduleManifest) {
	t.Helper()
	if len(manifest.Backend.Routes) == 0 {
		t.Fatalf("manifest backend.routes must not be empty")
	}
	for _, route := range manifest.Backend.Routes {
		routeNeedle := fmt.Sprintf(`rg.%s("%s",`, strings.ToUpper(route.Method), route.Path)
		if !strings.Contains(routeCode, routeNeedle) {
			t.Fatalf("generated route missing %s", routeNeedle)
		}
		apiNeedle := fmt.Sprintf("url: '%s'", route.Path)
		if !strings.Contains(apiCode, apiNeedle) {
			t.Fatalf("generated api missing %s", apiNeedle)
		}
	}
}

func assertManifestPermissions(t *testing.T, indexCode string, manifest moduleManifest) {
	t.Helper()
	permissionByAction := map[string]string{}
	for _, permission := range manifest.Permissions {
		permissionByAction[permission.Action] = permission.Code
	}
	for _, action := range []string{"add", "edit", "del"} {
		code := permissionByAction[action]
		if code == "" {
			t.Fatalf("manifest permission missing action %s", action)
		}
		needle := fmt.Sprintf("v-perms=\"['%s']\"", code)
		if !strings.Contains(indexCode, needle) {
			t.Fatalf("generated index missing permission %s", code)
		}
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
