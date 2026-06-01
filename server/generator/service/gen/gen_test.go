package gen

import (
	"testing"

	legacygen "go-makeadmin/model/gen"
	"go-makeadmin/model/makeadmin"
)

func TestCodegenTableLegacyConversionPreservesOldFields(t *testing.T) {
	legacy := legacygen.GenTable{
		ID:           9,
		TableName:    "ma_article",
		TableComment: "文章表",
		SubTableName: "ma_article_item",
		SubTableFk:   "article_id",
		AuthorName:   "codepeep",
		EntityName:   "Article",
		ModuleName:   "article",
		FunctionName: "文章",
		TreePrimary:  "id",
		TreeParent:   "pid",
		TreeName:     "title",
		GenTpl:       "tree",
		GenType:      legacyGenTypePath,
		GenPath:      "/tmp/go-makeadmin-gen",
		Remarks:      "demo",
		CreateTime:   100,
		UpdateTime:   200,
	}

	table := makeadminTableFromLegacy(legacy)
	if table.SourceTable != legacy.TableName {
		t.Fatalf("SourceTable = %q, want %q", table.SourceTable, legacy.TableName)
	}
	if table.GenerateType != makeadminGenerateTypePath {
		t.Fatalf("GenerateType = %q, want %q", table.GenerateType, makeadminGenerateTypePath)
	}
	if table.TemplateType != makeadmin.CodegenTemplateTree {
		t.Fatalf("TemplateType = %q, want %q", table.TemplateType, makeadmin.CodegenTemplateTree)
	}

	roundTrip := legacyTableFromMakeadmin(table)
	if roundTrip.TableName != legacy.TableName || roundTrip.GenType != legacy.GenType || roundTrip.GenTpl != legacy.GenTpl {
		t.Fatalf("round trip table mismatch: got %+v", roundTrip)
	}
	if roundTrip.SubTableName != legacy.SubTableName || roundTrip.TreeParent != legacy.TreeParent {
		t.Fatalf("round trip options mismatch: got %+v", roundTrip)
	}
}

func TestCodegenColumnLegacyConversionPreservesOldFields(t *testing.T) {
	legacy := legacygen.GenTableColumn{
		ID:            11,
		TableID:       9,
		ColumnName:    "article_title",
		ColumnComment: "标题",
		ColumnLength:  120,
		ColumnType:    "varchar",
		JavaType:      "string",
		JavaField:     "article_title",
		IsPk:          0,
		IsRequired:    1,
		IsInsert:      1,
		IsEdit:        1,
		IsList:        1,
		IsQuery:       1,
		QueryType:     "LIKE",
		HtmlType:      "input",
		DictType:      "article_status",
		Sort:          3,
		CreateTime:    100,
		UpdateTime:    200,
	}

	column := makeadminColumnFromLegacy(legacy)
	if column.GoField != legacy.JavaField || column.HTMLType != legacy.HtmlType {
		t.Fatalf("converted column mismatch: got %+v", column)
	}
	if column.JSONField != "article_title" {
		t.Fatalf("JSONField = %q, want article_title", column.JSONField)
	}

	roundTrip := legacyColumnFromMakeadmin(column)
	if roundTrip.JavaField != legacy.JavaField || roundTrip.HtmlType != legacy.HtmlType || roundTrip.Sort != legacy.Sort {
		t.Fatalf("round trip column mismatch: got %+v", roundTrip)
	}
}
