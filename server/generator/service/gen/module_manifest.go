package gen

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"go-makeadmin/config"
	"go-makeadmin/core/response"
	"go-makeadmin/generator"
	"go-makeadmin/generator/schemas/req"
	"go-makeadmin/generator/schemas/resp"
	legacygen "go-makeadmin/model/gen"
	"go-makeadmin/model/makeadmin"
)

const maxModuleManifestBodyBytes = 256 * 1024

var permissionCodePattern = regexp.MustCompile(`^[a-z][a-z0-9_]*(?::[a-z][a-z0-9_]+){1,2}$`)

type moduleManifest struct {
	Version           int    `json:"version"`
	Module            string `json:"module"`
	Entity            string `json:"entity"`
	Table             string `json:"table"`
	BackendPackage    string `json:"backendPackage"`
	RuntimeRegistered bool   `json:"runtimeRegistered"`
	RequiresSchema    bool   `json:"requiresSchema"`
	Backend           struct {
		Routes []moduleManifestRoute `json:"routes"`
	} `json:"backend"`
	Frontend struct {
		API   string   `json:"api"`
		Views []string `json:"views"`
	} `json:"frontend"`
	Menu        moduleManifestMenu         `json:"menu"`
	Permissions []moduleManifestPermission `json:"permissions"`
	Codegen     moduleManifestCodegen      `json:"codegen"`
}

type moduleManifestRoute struct {
	Method     string `json:"method"`
	Path       string `json:"path"`
	Permission string `json:"permission"`
}

type moduleManifestMenu struct {
	Code       string `json:"code"`
	Parent     string `json:"parent"`
	Type       string `json:"type"`
	Name       string `json:"name"`
	RoutePath  string `json:"routePath"`
	RouteName  string `json:"routeName"`
	Component  string `json:"component"`
	Permission string `json:"permission"`
	Visible    bool   `json:"visible"`
	Sort       int    `json:"sort"`
}

type moduleManifestPermission struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	Module   string `json:"module"`
	Resource string `json:"resource"`
	Action   string `json:"action"`
}

type moduleManifestCodegen struct {
	Columns []moduleManifestColumn `json:"columns"`
}

type moduleManifestColumn struct {
	ColumnName    string `json:"columnName"`
	ColumnComment string `json:"columnComment"`
	ColumnLength  int    `json:"columnLength"`
	ColumnType    string `json:"columnType"`
	GoType        string `json:"goType"`
	GoField       string `json:"goField"`
	IsRequired    uint8  `json:"isRequired"`
	IsInsert      *uint8 `json:"isInsert"`
	IsEdit        *uint8 `json:"isEdit"`
	IsList        *uint8 `json:"isList"`
	IsQuery       *uint8 `json:"isQuery"`
	QueryType     string `json:"queryType"`
	HtmlType      string `json:"htmlType"`
	DictType      string `json:"dictType"`
	Sort          int    `json:"sort"`
}

type moduleColumnDefaults struct {
	ColumnType   string
	ColumnLength int
	GoType       string
	QueryType    string
}

var htmlColumnDefaults = map[string]moduleColumnDefaults{
	"input":    {ColumnType: "varchar", ColumnLength: 255, GoType: "string", QueryType: generator.GenConstants.QueryLike},
	"number":   {ColumnType: "int", ColumnLength: 0, GoType: "int", QueryType: generator.GenConstants.QueryEq},
	"textarea": {ColumnType: "text", ColumnLength: 0, GoType: "string", QueryType: generator.GenConstants.QueryLike},
	"select":   {ColumnType: "varchar", ColumnLength: 100, GoType: "string", QueryType: generator.GenConstants.QueryEq},
	"radio":    {ColumnType: "varchar", ColumnLength: 100, GoType: "string", QueryType: generator.GenConstants.QueryEq},
	"checkbox": {ColumnType: "varchar", ColumnLength: 255, GoType: "string", QueryType: generator.GenConstants.QueryEq},
	"datetime": {ColumnType: "datetime", ColumnLength: 0, GoType: "time.Time", QueryType: generator.GenConstants.QueryEq},
}

// PreviewModuleManifest converts a module manifest into the generator detail shape and rendered code preview.
func (genSrv generateService) PreviewModuleManifest(previewReq req.ModuleManifestPreviewReq) (res resp.ModuleManifestPreviewResp, e error) {
	manifest, source, err := loadModuleManifest(previewReq)
	if err != nil {
		return res, err
	}
	if err = validateModuleManifest(manifest); err != nil {
		return res, err
	}
	table, columns := moduleManifestLegacyConfig(manifest, previewReq, source)
	code, err := renderCodeByLegacyTable(table, columns)
	if err != nil {
		return res, err
	}

	makeadminTable := makeadminTableFromLegacy(table)
	makeadminTable.TenantID = manifestTenantID(previewReq)
	makeadminColumns := make([]makeadmin.CodegenColumn, 0, len(columns))
	for _, column := range columns {
		makeadminColumns = append(makeadminColumns, makeadminColumnFromLegacy(column))
	}
	return resp.ModuleManifestPreviewResp{
		Source:  source,
		Warning: "Preview only; this request did not connect to a database, did not write ma_codegen_* rows, and did not execute install SQL.",
		Manifest: resp.ModuleManifestSummaryResp{
			Module:         manifest.Module,
			Entity:         manifest.Entity,
			Table:          manifest.Table,
			MenuName:       manifest.Menu.Name,
			RequiresSchema: manifest.RequiresSchema,
		},
		Detail: genTableDetailRespFromMakeadmin(makeadminTable, makeadminColumns),
		Code:   code,
		Plan:   buildModuleManifestPlan(manifest, previewReq),
	}, nil
}

func loadModuleManifest(previewReq req.ModuleManifestPreviewReq) (moduleManifest, string, error) {
	var manifest moduleManifest
	manifestPath := strings.TrimSpace(previewReq.ManifestPath)
	manifestBody := strings.TrimSpace(previewReq.ManifestBody)
	if manifestPath == "" && manifestBody == "" {
		return manifest, "", responseArgumentError("manifestPath or manifestBody is required")
	}
	if manifestPath != "" && manifestBody != "" {
		return manifest, "", responseArgumentError("only one of manifestPath or manifestBody can be provided")
	}
	if manifestBody != "" {
		if len([]byte(manifestBody)) > maxModuleManifestBodyBytes {
			return manifest, "", responseArgumentError("manifestBody is too large")
		}
		if err := json.Unmarshal([]byte(manifestBody), &manifest); err != nil {
			return manifest, "", responseArgumentError("manifestBody must be valid JSON")
		}
		return manifest, "inline", nil
	}

	absPath, relPath, err := resolveManifestPath(manifestPath)
	if err != nil {
		return manifest, "", err
	}
	content, readErr := os.ReadFile(absPath)
	if readErr != nil {
		return manifest, "", responseArgumentError(fmt.Sprintf("manifestPath cannot be read: %s", relPath))
	}
	if len(content) > maxModuleManifestBodyBytes {
		return manifest, "", responseArgumentError("manifestPath file is too large")
	}
	if err := json.Unmarshal(content, &manifest); err != nil {
		return manifest, "", responseArgumentError("manifestPath must contain valid JSON")
	}
	return manifest, relPath, nil
}

func resolveManifestPath(manifestPath string) (string, string, error) {
	repoRoot := filepath.Dir(config.Config.RootPath)
	candidate := manifestPath
	if !filepath.IsAbs(candidate) {
		candidate = filepath.Join(repoRoot, candidate)
	}
	absPath, err := filepath.Abs(candidate)
	if err != nil {
		return "", "", responseArgumentError("manifestPath is invalid")
	}
	relPath, err := filepath.Rel(repoRoot, absPath)
	if err != nil || strings.HasPrefix(relPath, ".."+string(filepath.Separator)) || relPath == ".." || filepath.IsAbs(relPath) {
		return "", "", responseArgumentError("manifestPath must stay inside repository")
	}
	if filepath.Base(absPath) != "manifest.json" {
		return "", "", responseArgumentError("manifestPath must point to manifest.json")
	}
	return absPath, filepath.ToSlash(relPath), nil
}

func validateModuleManifest(manifest moduleManifest) error {
	if manifest.Version != 1 {
		return responseArgumentError("version must be 1")
	}
	if err := requireManifestText("module", manifest.Module); err != nil {
		return err
	}
	if err := requireManifestText("entity", manifest.Entity); err != nil {
		return err
	}
	if err := requireManifestText("table", manifest.Table); err != nil {
		return err
	}
	if err := requireManifestText("backendPackage", manifest.BackendPackage); err != nil {
		return err
	}
	if err := validateManifestPermissions(manifest.Permissions); err != nil {
		return err
	}
	permissionCodes := make(map[string]bool, len(manifest.Permissions))
	for _, permission := range manifest.Permissions {
		permissionCodes[permission.Code] = true
	}
	if err := validateManifestRoutes(manifest.Backend.Routes, permissionCodes); err != nil {
		return err
	}
	if err := validateManifestFrontend(manifest.Frontend.API, manifest.Frontend.Views); err != nil {
		return err
	}
	if err := validateManifestMenu(manifest.Menu, permissionCodes, manifest.Module); err != nil {
		return err
	}
	return validateManifestCodegen(manifest.Codegen)
}

func validateManifestPermissions(permissions []moduleManifestPermission) error {
	if len(permissions) == 0 {
		return responseArgumentError("permissions must be a non-empty list")
	}
	seen := map[string]bool{}
	for _, permission := range permissions {
		if err := requireManifestText("permission.code", permission.Code); err != nil {
			return err
		}
		if !permissionCodePattern.MatchString(permission.Code) {
			return responseArgumentError("permission.code is invalid")
		}
		if seen[permission.Code] {
			return responseArgumentError("permission.code must be unique")
		}
		seen[permission.Code] = true
		for key, value := range map[string]string{
			"permission.name":     permission.Name,
			"permission.module":   permission.Module,
			"permission.resource": permission.Resource,
			"permission.action":   permission.Action,
		} {
			if err := requireManifestText(key, value); err != nil {
				return err
			}
		}
	}
	return nil
}

func validateManifestRoutes(routes []moduleManifestRoute, permissionCodes map[string]bool) error {
	if len(routes) == 0 {
		return responseArgumentError("backend.routes must be a non-empty list")
	}
	seen := map[string]bool{}
	for _, route := range routes {
		method := strings.ToUpper(strings.TrimSpace(route.Method))
		if !map[string]bool{"GET": true, "POST": true, "PUT": true, "PATCH": true, "DELETE": true}[method] {
			return responseArgumentError("backend.routes method is unsupported")
		}
		if err := requireManifestText("backend.routes.path", route.Path); err != nil {
			return err
		}
		if !strings.HasPrefix(route.Path, "/") {
			return responseArgumentError("backend.routes.path must start with /")
		}
		routeKey := method + " " + route.Path
		if seen[routeKey] {
			return responseArgumentError("backend.routes entries must be unique")
		}
		seen[routeKey] = true
		if !permissionCodes[route.Permission] {
			return responseArgumentError("backend.routes.permission must be declared in permissions")
		}
	}
	return nil
}

func validateManifestFrontend(apiPath string, views []string) error {
	if err := requireManifestText("frontend.api", apiPath); err != nil {
		return err
	}
	if !strings.HasPrefix(apiPath, "admin/src/api/") {
		return responseArgumentError("frontend.api must live under admin/src/api")
	}
	if len(views) == 0 {
		return responseArgumentError("frontend.views must be a non-empty list")
	}
	for _, view := range views {
		if !strings.HasPrefix(view, "admin/src/views/") {
			return responseArgumentError("frontend.views must live under admin/src/views")
		}
	}
	return nil
}

func validateManifestMenu(menu moduleManifestMenu, permissionCodes map[string]bool, module string) error {
	for key, value := range map[string]string{
		"menu.code":       menu.Code,
		"menu.parent":     menu.Parent,
		"menu.type":       menu.Type,
		"menu.name":       menu.Name,
		"menu.routePath":  menu.RoutePath,
		"menu.routeName":  menu.RouteName,
		"menu.component":  menu.Component,
		"menu.permission": menu.Permission,
	} {
		if err := requireManifestText(key, value); err != nil {
			return err
		}
	}
	if !map[string]bool{"catalog": true, "page": true, "button": true}[menu.Type] {
		return responseArgumentError("menu.type is unsupported")
	}
	if !permissionCodes[menu.Permission] {
		return responseArgumentError("menu.permission must be declared in permissions")
	}
	if !strings.Contains(menu.Component, module) && !strings.Contains(menu.RouteName, module) {
		return responseArgumentError("menu should reference module in component or routeName")
	}
	return nil
}

func validateManifestCodegen(codegen moduleManifestCodegen) error {
	for _, column := range codegen.Columns {
		if err := requireManifestText("codegen.columns.columnName", column.ColumnName); err != nil {
			return err
		}
		if err := requireManifestText("codegen.columns.goField", column.GoField); err != nil {
			return err
		}
		if err := requireManifestText("codegen.columns.htmlType", column.HtmlType); err != nil {
			return err
		}
		if _, ok := htmlColumnDefaults[column.HtmlType]; !ok {
			return responseArgumentError("codegen.columns.htmlType is unsupported")
		}
	}
	return nil
}

func requireManifestText(key string, value string) error {
	if strings.TrimSpace(value) == "" {
		return responseArgumentError(key + " must be a non-empty string")
	}
	return nil
}

func moduleManifestLegacyConfig(manifest moduleManifest, previewReq req.ModuleManifestPreviewReq, source string) (legacygen.GenTable, []legacygen.GenTableColumn) {
	now := time.Now().Unix()
	authorName := strings.TrimSpace(previewReq.AuthorName)
	if authorName == "" {
		authorName = "codepeep"
	}
	table := legacygen.GenTable{
		TableName:    manifest.Table,
		TableComment: manifest.Menu.Name,
		AuthorName:   authorName,
		EntityName:   manifest.Entity,
		ModuleName:   manifest.Module,
		FunctionName: manifest.Menu.Name,
		GenTpl:       generator.GenConstants.TplCrud,
		GenType:      legacyGenTypeZip,
		GenPath:      defaultGenPath,
		Remarks:      "generated from " + source,
		CreateTime:   now,
		UpdateTime:   now,
	}
	return table, moduleManifestLegacyColumns(manifest.Codegen.Columns, now)
}

func moduleManifestLegacyColumns(configuredColumns []moduleManifestColumn, now int64) []legacygen.GenTableColumn {
	if len(configuredColumns) == 0 {
		return defaultManifestColumns(now)
	}
	columns := []legacygen.GenTableColumn{primaryManifestColumn(now)}
	for index, column := range configuredColumns {
		columns = append(columns, configuredManifestColumn(column, index+2, now))
	}
	return columns
}

func defaultManifestColumns(now int64) []legacygen.GenTableColumn {
	return []legacygen.GenTableColumn{
		primaryManifestColumn(now),
		{
			ColumnName:    "title",
			ColumnComment: "Title",
			ColumnLength:  200,
			ColumnType:    "varchar",
			JavaType:      generator.GoConstants.TypeString,
			JavaField:     "title",
			IsRequired:    1,
			IsInsert:      1,
			IsEdit:        1,
			IsList:        1,
			IsQuery:       1,
			QueryType:     generator.GenConstants.QueryLike,
			HtmlType:      generator.HtmlConstants.HtmlInput,
			Sort:          2,
			CreateTime:    now,
			UpdateTime:    now,
		},
		{
			ColumnName:    "status",
			ColumnComment: "Status",
			ColumnType:    "tinyint",
			JavaType:      generator.GoConstants.TypeInt,
			JavaField:     "status",
			IsInsert:      1,
			IsEdit:        1,
			IsList:        1,
			IsQuery:       1,
			QueryType:     generator.GenConstants.QueryEq,
			HtmlType:      generator.HtmlConstants.HtmlInput,
			Sort:          3,
			CreateTime:    now,
			UpdateTime:    now,
		},
	}
}

func primaryManifestColumn(now int64) legacygen.GenTableColumn {
	return legacygen.GenTableColumn{
		ColumnName:    "id",
		ColumnComment: "ID",
		ColumnType:    "bigint",
		JavaType:      "uint",
		JavaField:     "id",
		IsPk:          1,
		IsIncrement:   1,
		IsList:        1,
		QueryType:     generator.GenConstants.QueryEq,
		HtmlType:      generator.HtmlConstants.HtmlInput,
		Sort:          1,
		CreateTime:    now,
		UpdateTime:    now,
	}
}

func configuredManifestColumn(column moduleManifestColumn, sort int, now int64) legacygen.GenTableColumn {
	defaults := htmlColumnDefaults[column.HtmlType]
	columnType := firstNonEmpty(column.ColumnType, defaults.ColumnType)
	goType := firstNonEmpty(column.GoType, defaults.GoType)
	queryType := firstNonEmpty(column.QueryType, defaults.QueryType)
	columnLength := column.ColumnLength
	if columnLength == 0 {
		columnLength = defaults.ColumnLength
	}
	if column.Sort > 0 {
		sort = column.Sort
	}
	return legacygen.GenTableColumn{
		ColumnName:    column.ColumnName,
		ColumnComment: firstNonEmpty(column.ColumnComment, strings.Title(strings.ReplaceAll(column.ColumnName, "_", " "))),
		ColumnLength:  columnLength,
		ColumnType:    columnType,
		JavaType:      goType,
		JavaField:     column.GoField,
		IsRequired:    column.IsRequired,
		IsInsert:      valueOrDefault(column.IsInsert, 1),
		IsEdit:        valueOrDefault(column.IsEdit, 1),
		IsList:        valueOrDefault(column.IsList, 1),
		IsQuery:       valueOrDefault(column.IsQuery, 1),
		QueryType:     queryType,
		HtmlType:      column.HtmlType,
		DictType:      column.DictType,
		Sort:          sort,
		CreateTime:    now,
		UpdateTime:    now,
	}
}

func manifestTenantID(previewReq req.ModuleManifestPreviewReq) uint64 {
	if previewReq.TenantID > 0 {
		return previewReq.TenantID
	}
	return makeadmin.GlobalTenantID
}

func manifestRoleID(previewReq req.ModuleManifestPreviewReq) uint64 {
	if previewReq.RoleID > 0 {
		return previewReq.RoleID
	}
	return 1
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func valueOrDefault(value *uint8, fallback uint8) uint8 {
	if value == nil {
		return fallback
	}
	return *value
}

func responseArgumentError(message string) error {
	return response.AssertArgumentError.Make(message)
}
