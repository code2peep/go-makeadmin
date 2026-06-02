package generator

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"go-makeadmin/config"
	"go-makeadmin/model/gen"
)

func TestGeneratedCrudGoCodeCompiles(t *testing.T) {
	table := gen.GenTable{
		ID:           1,
		TableName:    "ma_demo_article",
		TableComment: "Demo article",
		AuthorName:   "codepeep",
		EntityName:   "DemoArticle",
		ModuleName:   "demo_article",
		FunctionName: "Demo article",
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
