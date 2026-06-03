package gen

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"go-makeadmin/core/request"
	genreq "go-makeadmin/generator/schemas/req"
	"go-makeadmin/generator/schemas/resp"
	"go-makeadmin/model/makeadmin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestCodegenConfigReadbackAndTemplateGenerationSmoke(t *testing.T) {
	if os.Getenv("MAKEADMIN_CODEGEN_READBACK_SMOKE") != "1" {
		t.Skip("set MAKEADMIN_CODEGEN_READBACK_SMOKE=1 to verify local ma_codegen_* readback")
	}
	db := openReadbackSmokeDB(t)
	srv := NewGenerateService(db)

	page, err := srv.List(request.PageReq{PageNo: 1, PageSize: 10}, genreq.ListTableReq{TableName: "ma_demo_article"})
	if err != nil {
		t.Fatalf("list codegen table: %v", err)
	}
	if page.Count != 1 {
		t.Fatalf("list count = %d, want 1", page.Count)
	}
	lists, ok := page.Lists.([]resp.GenTableResp)
	if !ok {
		t.Fatalf("list type = %T, want []resp.GenTableResp", page.Lists)
	}
	if len(lists) != 1 {
		t.Fatalf("list length = %d, want 1", len(lists))
	}
	item := lists[0]
	if item.TableName != "ma_demo_article" || item.GenType != legacyGenTypeZip {
		t.Fatalf("unexpected list item: %+v", item)
	}

	detail, err := srv.Detail(item.ID)
	if err != nil {
		t.Fatalf("detail codegen table: %v", err)
	}
	assertLegacyDetailShape(t, detail)

	preview, err := srv.PreviewCode(item.ID)
	if err != nil {
		t.Fatalf("preview code: %v", err)
	}
	assertRenderedCode(t, preview)

	zipBytes, err := srv.DownloadCode([]string{"ma_demo_article"})
	if err != nil {
		t.Fatalf("download code: %v", err)
	}
	assertDownloadedZip(t, zipBytes)
}

func openReadbackSmokeDB(t *testing.T) *gorm.DB {
	t.Helper()
	user := envOrDefault("MYSQL_USER", "root")
	password := os.Getenv("MYSQL_PASSWORD")
	host := envOrDefault("MYSQL_HOST", "127.0.0.1")
	port := envOrDefault("MYSQL_PORT", "3306")
	database := envOrDefault("MYSQL_DATABASE", "go_makeadmin")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open mysql: %v", err)
	}
	return db
}

func envOrDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func assertLegacyDetailShape(t *testing.T, detail resp.GenTableDetailResp) {
	t.Helper()
	if detail.Base.TableName != "ma_demo_article" || detail.Base.EntityName != "DemoArticle" {
		t.Fatalf("unexpected detail base: %+v", detail.Base)
	}
	if detail.Gen.GenTpl != makeadmin.CodegenTemplateCRUD || detail.Gen.GenType != legacyGenTypeZip || detail.Gen.GenPath != defaultGenPath {
		t.Fatalf("unexpected legacy gen shape: %+v", detail.Gen)
	}
	if detail.Gen.ModuleName != "article" || detail.Gen.FunctionName != "Demo Article" {
		t.Fatalf("unexpected legacy module shape: %+v", detail.Gen)
	}
	if len(detail.Column) != 3 {
		t.Fatalf("column length = %d, want 3", len(detail.Column))
	}
	columns := map[string]resp.GenColumnResp{}
	for _, column := range detail.Column {
		columns[column.ColumnName] = column
	}
	if columns["id"].JavaType != "uint" || columns["id"].JavaField != "id" {
		t.Fatalf("unexpected id column: %+v", columns["id"])
	}
	if columns["title"].JavaType != "string" || columns["title"].JavaField != "title" || columns["title"].QueryType != "LIKE" {
		t.Fatalf("unexpected title column: %+v", columns["title"])
	}
	if columns["status"].JavaType != "int" || columns["status"].JavaField != "status" || columns["status"].QueryType != "=" {
		t.Fatalf("unexpected status column: %+v", columns["status"])
	}
}

func assertRenderedCode(t *testing.T, preview map[string]string) {
	t.Helper()
	expectedKeys := []string{
		"gocode/model.go",
		"gocode/schema.go",
		"gocode/service.go",
		"gocode/route.go",
		"vue/api.ts",
		"vue/edit.vue",
		"vue/index.vue",
	}
	for _, key := range expectedKeys {
		if strings.TrimSpace(preview[key]) == "" {
			t.Fatalf("preview missing %s", key)
		}
	}
	modelCode := strings.Join(strings.Fields(preview["gocode/model.go"]), " ")
	assertContains(t, modelCode, "type DemoArticle struct")
	assertContains(t, modelCode, "Title string")
	assertContains(t, modelCode, "Status int")
	assertContains(t, preview["gocode/route.go"], `rg.GET("/article/list"`)
	assertContains(t, preview["vue/api.ts"], "url: '/article/list'")
	assertContains(t, preview["vue/index.vue"], "v-perms=\"['article:add']\"")
	assertContains(t, preview["vue/index.vue"], "handleAdd")
}

func assertDownloadedZip(t *testing.T, zipBytes []byte) {
	t.Helper()
	reader, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		t.Fatalf("open generated zip: %v", err)
	}
	entries := map[string]bool{}
	for _, file := range reader.File {
		entries[file.Name] = true
	}
	expectedEntries := []string{
		"gocode/article/model.go",
		"gocode/article/schema.go",
		"gocode/article/service.go",
		"gocode/article/route.go",
		"vue/article/api.ts",
		"vue/article/edit.vue",
		"vue/article/index.vue",
	}
	for _, entry := range expectedEntries {
		if !entries[entry] {
			t.Fatalf("zip missing %s; entries=%v", entry, entries)
		}
	}
}

func assertContains(t *testing.T, body string, needle string) {
	t.Helper()
	if !strings.Contains(body, needle) {
		t.Fatalf("expected generated code to contain %q", needle)
	}
}
