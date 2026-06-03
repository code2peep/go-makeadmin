package gen

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"strings"
	"time"

	"go-makeadmin/config"
	"go-makeadmin/core"
	"go-makeadmin/core/request"
	"go-makeadmin/core/response"
	"go-makeadmin/generator"
	"go-makeadmin/generator/schemas/req"
	"go-makeadmin/generator/schemas/resp"
	legacygen "go-makeadmin/model/gen"
	"go-makeadmin/model/makeadmin"
	"go-makeadmin/util"
	"gorm.io/gorm"
)

const (
	legacyGenTypeZip  = 0
	legacyGenTypePath = 1

	makeadminGenerateTypeZip  = "zip"
	makeadminGenerateTypePath = "path"
	defaultGenPath            = "/"
)

type codegenTableOptions struct {
	TreePrimary  string `json:"treePrimary,omitempty"`
	TreeParent   string `json:"treeParent,omitempty"`
	TreeName     string `json:"treeName,omitempty"`
	SubTableName string `json:"subTableName,omitempty"`
	SubTableFk   string `json:"subTableFk,omitempty"`
}

type IGenerateService interface {
	DbTables(page request.PageReq, req req.DbTablesReq) (res response.PageResp, e error)
	List(page request.PageReq, listReq req.ListTableReq) (res response.PageResp, e error)
	Detail(id uint) (res resp.GenTableDetailResp, e error)
	ImportTable(tableNames []string) (e error)
	SyncTable(id uint) (e error)
	EditTable(editReq req.EditTableReq) (e error)
	DelTable(ids []uint) (e error)
	PreviewCode(id uint) (res map[string]string, e error)
	ListModuleRegistry() []resp.ModuleRegistryItemResp
	PreviewModuleManifest(previewReq req.ModuleManifestPreviewReq) (res resp.ModuleManifestPreviewResp, e error)
	ReadModuleManifestInstallStatus(previewReq req.ModuleManifestPreviewReq) (res resp.ModuleManifestInstallStatusResp, e error)
	ApplyModuleManifestInstall(applyReq req.ModuleManifestInstallApplyReq) (res resp.ModuleManifestInstallApplyResp, e error)
	ApplyModuleManifestUninstall(applyReq req.ModuleManifestUninstallApplyReq) (res resp.ModuleManifestUninstallApplyResp, e error)
	GenCode(tableName string) (e error)
	DownloadCode(tableNames []string) ([]byte, error)
}

// NewGenerateService 初始化
func NewGenerateService(db *gorm.DB) IGenerateService {
	return &generateService{db: db}
}

// GenerateService 代码生成服务实现类
type generateService struct {
	db *gorm.DB
}

// DbTables 库表列表
func (genSrv generateService) DbTables(page request.PageReq, dbReq req.DbTablesReq) (res response.PageResp, e error) {
	page = normalizePage(page)
	// 分页信息
	limit := page.PageSize
	offset := page.PageSize * (page.PageNo - 1)
	tbModel := generator.GenUtil.GetDbTablesQuery(genSrv.db, dbReq.TableName, dbReq.TableComment)
	// 总数
	var count int64
	err := tbModel.Count(&count).Error
	if e = response.CheckErr(err, "DbTables Count err"); e != nil {
		return
	}
	// 数据
	var tbResp []resp.DbTableResp
	err = tbModel.Limit(limit).Offset(offset).Find(&tbResp).Error
	if e = response.CheckErr(err, "DbTables Find err"); e != nil {
		return
	}
	return response.PageResp{
		PageNo:   page.PageNo,
		PageSize: page.PageSize,
		Count:    count,
		Lists:    tbResp,
	}, nil
}

// List 生成列表
func (genSrv generateService) List(page request.PageReq, listReq req.ListTableReq) (res response.PageResp, e error) {
	page = normalizePage(page)
	// 分页信息
	limit := page.PageSize
	offset := page.PageSize * (page.PageNo - 1)
	genModel := genSrv.db.Model(&makeadmin.CodegenTable{}).
		Where("tenant_id = ? AND delete_time = ?", makeadmin.GlobalTenantID, 0)
	if listReq.TableName != "" {
		genModel = genModel.Where("table_name like ?", "%"+listReq.TableName+"%")
	}
	if listReq.TableComment != "" {
		genModel = genModel.Where("table_comment like ?", "%"+listReq.TableComment+"%")
	}
	if !listReq.StartTime.IsZero() {
		genModel = genModel.Where("create_time >= ?", listReq.StartTime.Unix())
	}
	if !listReq.EndTime.IsZero() {
		genModel = genModel.Where("create_time <= ?", listReq.EndTime.Unix())
	}
	// 总数
	var count int64
	err := genModel.Count(&count).Error
	if e = response.CheckErr(err, "List Count err"); e != nil {
		return
	}
	// 数据
	var tables []makeadmin.CodegenTable
	err = genModel.Order("id desc").Limit(limit).Offset(offset).Find(&tables).Error
	if e = response.CheckErr(err, "List Find err"); e != nil {
		return
	}
	genResp := make([]resp.GenTableResp, 0, len(tables))
	for _, table := range tables {
		genResp = append(genResp, genTableRespFromMakeadmin(table))
	}
	return response.PageResp{
		PageNo:   page.PageNo,
		PageSize: page.PageSize,
		Count:    count,
		Lists:    genResp,
	}, nil
}

// Detail 生成详情
func (genSrv generateService) Detail(id uint) (res resp.GenTableDetailResp, e error) {
	genTb, err := genSrv.findTable(id)
	if e = response.CheckErrDBNotRecord(err, "查询的数据不存在!"); e != nil {
		return
	}
	if e = response.CheckErr(err, "Detail Find err"); e != nil {
		return
	}
	columns, err := genSrv.listColumns(id)
	if e = response.CheckErr(err, "Detail Find err"); e != nil {
		return
	}
	return genTableDetailRespFromMakeadmin(genTb, columns), e
}

// ImportTable 导入表结构
func (genSrv generateService) ImportTable(tableNames []string) (e error) {
	var dbTbs []resp.DbTableResp
	err := generator.GenUtil.GetDbTablesQueryByNames(genSrv.db, tableNames).Find(&dbTbs).Error
	if e = response.CheckErr(err, "ImportTable Find tables err"); e != nil {
		return
	}
	var tables []legacygen.GenTable
	response.Copy(&tables, dbTbs)
	if len(tables) == 0 {
		e = response.AssertArgumentError.Make("表不存在!")
		return
	}
	err = genSrv.db.Transaction(func(tx *gorm.DB) error {
		for i := 0; i < len(tables); i++ {
			//生成表信息
			legacyTable := generator.GenUtil.InitTable(tables[i])
			genTable := makeadminTableFromLegacy(legacyTable)
			txErr := tx.Create(&genTable).Error
			if te := response.CheckErr(txErr, "ImportTable Create table err"); te != nil {
				return te
			}
			// 生成列信息
			var columns []legacygen.GenTableColumn
			txErr = generator.GenUtil.GetDbTableColumnsQueryByName(genSrv.db, tables[i].TableName).Find(&columns).Error
			if te := response.CheckErr(txErr, "ImportTable Find columns err"); te != nil {
				return te
			}
			for j := 0; j < len(columns); j++ {
				legacyColumn := generator.GenUtil.InitColumn(uint(genTable.ID), columns[j])
				column := makeadminColumnFromLegacy(legacyColumn)
				txErr = tx.Create(&column).Error
				if te := response.CheckErr(txErr, "ImportTable Create column err"); te != nil {
					return te
				}
			}
		}
		return nil
	})
	e = response.CheckErr(err, "ImportTable Transaction err")
	return
}

// SyncTable 同步表结构
func (genSrv generateService) SyncTable(id uint) (e error) {
	//旧数据
	genTable, err := genSrv.findTable(id)
	if e = response.CheckErrDBNotRecord(err, "生成数据不存在！"); e != nil {
		return
	}
	if e = response.CheckErr(err, "SyncTable First err"); e != nil {
		return
	}
	genTableCols, err := genSrv.listColumns(id)
	if e = response.CheckErr(err, "SyncTable Find err"); e != nil {
		return
	}
	if len(genTableCols) <= 0 {
		e = response.AssertArgumentError.Make("旧数据异常！")
		return
	}
	prevColMap := make(map[string]legacygen.GenTableColumn)
	for i := 0; i < len(genTableCols); i++ {
		prevColMap[genTableCols[i].ColumnName] = legacyColumnFromMakeadmin(genTableCols[i])
	}
	//新数据
	var columns []legacygen.GenTableColumn
	err = generator.GenUtil.GetDbTableColumnsQueryByName(genSrv.db, genTable.SourceTable).Find(&columns).Error
	if e = response.CheckErr(err, "SyncTable Find new err"); e != nil {
		return
	}
	if len(columns) <= 0 {
		e = response.AssertArgumentError.Make("同步结构失败,原表结构不存在！")
		return
	}
	//事务处理
	err = genSrv.db.Transaction(func(tx *gorm.DB) error {
		//处理新增和更新
		for i := 0; i < len(columns); i++ {
			col := generator.GenUtil.InitColumn(id, columns[i])
			if prevCol, ok := prevColMap[columns[i].ColumnName]; ok {
				//更新
				col.ID = prevCol.ID
				if col.IsList == 0 {
					col.DictType = prevCol.DictType
					col.QueryType = prevCol.QueryType
				}
				if (prevCol.IsRequired == 1 && prevCol.IsPk == 0 && prevCol.IsInsert == 1) || prevCol.IsEdit == 1 {
					col.HtmlType = prevCol.HtmlType
					col.IsRequired = prevCol.IsRequired
				}
				nextCol := makeadminColumnFromLegacy(col)
				txErr := tx.Save(&nextCol).Error
				if te := response.CheckErr(txErr, "SyncTable Save err"); te != nil {
					return te
				}
			} else {
				//新增
				nextCol := makeadminColumnFromLegacy(col)
				txErr := tx.Create(&nextCol).Error
				if te := response.CheckErr(txErr, "SyncTable Create err"); te != nil {
					return te
				}
			}
		}
		//处理删除
		colNames := make([]string, len(columns))
		for i := 0; i < len(columns); i++ {
			colNames[i] = columns[i].ColumnName
		}
		delColIds := make([]uint, 0)
		for _, prevCol := range prevColMap {
			if !util.ToolsUtil.Contains(colNames, prevCol.ColumnName) {
				delColIds = append(delColIds, prevCol.ID)
			}
		}
		if len(delColIds) > 0 {
			txErr := tx.Delete(&makeadmin.CodegenColumn{}, "id in ?", delColIds).Error
			if te := response.CheckErr(txErr, "SyncTable Delete err"); te != nil {
				return te
			}
		}
		return nil
	})
	e = response.CheckErr(err, "SyncTable Transaction err")
	return
}

// EditTable 编辑表结构
func (genSrv generateService) EditTable(editReq req.EditTableReq) (e error) {
	if editReq.GenTpl == generator.GenConstants.TplTree {
		if editReq.TreePrimary == "" {
			e = response.AssertArgumentError.Make("树主ID不能为空！")
			return
		}
		if editReq.TreeParent == "" {
			e = response.AssertArgumentError.Make("树父ID不能为空！")
			return
		}
	}
	genTable, err := genSrv.findTable(editReq.ID)
	if e = response.CheckErrDBNotRecord(err, "数据已丢失！"); e != nil {
		return
	}
	if e = response.CheckErr(err, "EditTable First err"); e != nil {
		return
	}
	applyEditReqToMakeadminTable(&genTable, editReq)
	err = genSrv.db.Transaction(func(tx *gorm.DB) error {
		txErr := tx.Save(&genTable).Error
		if te := response.CheckErr(txErr, "EditTable Save GenTable err"); te != nil {
			return te
		}
		for i := 0; i < len(editReq.Columns); i++ {
			var col makeadmin.CodegenColumn
			txErr = tx.Where("table_id = ? AND id = ?", genTable.ID, editReq.Columns[i].ID).Limit(1).First(&col).Error
			if te := response.CheckErrDBNotRecord(txErr, "字段数据已丢失！"); te != nil {
				return te
			}
			if te := response.CheckErr(txErr, "EditTable First GenTableColumn err"); te != nil {
				return te
			}
			applyEditReqToMakeadminColumn(&col, editReq.Columns[i])
			txErr = tx.Save(&col).Error
			if te := response.CheckErr(txErr, "EditTable Save GenTableColumn err"); te != nil {
				return te
			}
		}
		return nil
	})
	e = response.CheckErr(err, "EditTable Transaction err")
	return
}

// DelTable 删除表结构
func (genSrv generateService) DelTable(ids []uint) (e error) {
	if len(ids) == 0 {
		return nil
	}
	id64s := uintSliceToUint64(ids)
	err := genSrv.db.Transaction(func(tx *gorm.DB) error {
		txErr := tx.Model(&makeadmin.CodegenTable{}).
			Where("tenant_id = ? AND id in ? AND delete_time = ?", makeadmin.GlobalTenantID, id64s, 0).
			Update("delete_time", time.Now().Unix()).
			Error
		if te := response.CheckErr(txErr, "DelTable Delete GenTable err"); te != nil {
			return te
		}
		txErr = tx.Delete(&makeadmin.CodegenColumn{}, "table_id in ?", id64s).Error
		if te := response.CheckErr(txErr, "DelTable Delete GenTableColumn err"); te != nil {
			return te
		}
		return nil
	})
	e = response.CheckErr(err, "DelTable Transaction err")
	return
}

// getSubTableInfo 根据主表获取子表主键和列信息
func (genSrv generateService) getSubTableInfo(genTable legacygen.GenTable) (pkCol legacygen.GenTableColumn, cols []legacygen.GenTableColumn, e error) {
	if genTable.SubTableName == "" || genTable.SubTableFk == "" {
		return
	}
	var table makeadmin.CodegenTable
	err := genSrv.db.
		Where("tenant_id = ? AND table_name = ? AND delete_time = ?", makeadmin.GlobalTenantID, genTable.SubTableName, 0).
		Order("id desc").
		Limit(1).
		First(&table).
		Error
	if e = response.CheckErrDBNotRecord(err, "子表记录丢失！"); e != nil {
		return
	}
	if e = response.CheckErr(err, "getSubTableInfo First err"); e != nil {
		return
	}
	err = generator.GenUtil.GetDbTableColumnsQueryByName(genSrv.db, genTable.SubTableName).Find(&cols).Error
	if e = response.CheckErr(err, "getSubTableInfo Find err"); e != nil {
		return
	}
	pkCol = generator.GenUtil.InitColumn(uint(table.ID), generator.GenUtil.GetTablePriCol(cols))
	return
}

// renderCodeByTable 根据主表和模板文件渲染模板代码
func (genSrv generateService) renderCodeByTable(genTable legacygen.GenTable) (res map[string]string, e error) {
	var rawColumns []makeadmin.CodegenColumn
	err := genSrv.db.Where("table_id = ?", genTable.ID).Order("sort").Find(&rawColumns).Error
	if e = response.CheckErr(err, "renderCodeByTable Find err"); e != nil {
		return
	}
	columns := legacyColumnsFromMakeadmin(rawColumns)
	//获取子表信息
	pkCol, cols, err := genSrv.getSubTableInfo(genTable)
	if e = response.CheckErr(err, "renderCodeByTable getSubTableInfo err"); e != nil {
		return
	}
	//获取模板变量信息
	vars := generator.TemplateUtil.PrepareVars(genTable, columns, pkCol, cols)
	//生成模板内容
	res = make(map[string]string)
	for _, tplPath := range generator.TemplateUtil.GetTemplatePaths(genTable.GenTpl) {
		res[tplPath], err = generator.TemplateUtil.Render(tplPath, vars)
		if e = response.CheckErr(err, "renderCodeByTable Render err"); e != nil {
			return
		}
	}
	return
}

func renderCodeByLegacyTable(genTable legacygen.GenTable, columns []legacygen.GenTableColumn) (res map[string]string, e error) {
	vars := generator.TemplateUtil.PrepareVars(genTable, columns, legacygen.GenTableColumn{}, nil)
	res = make(map[string]string)
	for _, tplPath := range generator.TemplateUtil.GetTemplatePaths(genTable.GenTpl) {
		tplCode, err := generator.TemplateUtil.Render(tplPath, vars)
		if e = response.CheckErr(err, "renderCodeByLegacyTable Render err"); e != nil {
			return
		}
		res[strings.ReplaceAll(tplPath, ".tpl", "")] = tplCode
	}
	return
}

// PreviewCode 预览代码
func (genSrv generateService) PreviewCode(id uint) (res map[string]string, e error) {
	table, err := genSrv.findTable(id)
	if e = response.CheckErrDBNotRecord(err, "记录丢失！"); e != nil {
		return
	}
	if e = response.CheckErr(err, "PreviewCode First err"); e != nil {
		return
	}
	genTable := legacyTableFromMakeadmin(table)
	//获取模板内容
	tplCodeMap, err := genSrv.renderCodeByTable(genTable)
	if e = response.CheckErr(err, "PreviewCode renderCodeByTable err"); e != nil {
		return
	}
	res = make(map[string]string)
	for tplPath, tplCode := range tplCodeMap {
		res[strings.ReplaceAll(tplPath, ".tpl", "")] = tplCode
	}
	return
}

// GenCode 生成代码 (自定义路径)
func (genSrv generateService) GenCode(tableName string) (e error) {
	table, err := genSrv.findTableByName(tableName)
	if e = response.CheckErrDBNotRecord(err, "记录丢失！"); e != nil {
		return
	}
	if e = response.CheckErr(err, "GenCode First err"); e != nil {
		return
	}
	genTable := legacyTableFromMakeadmin(table)
	//获取模板内容
	tplCodeMap, err := genSrv.renderCodeByTable(genTable)
	if e = response.CheckErr(err, "GenCode renderCodeByTable err"); e != nil {
		return
	}
	//获取生成根路径
	basePath := generator.TemplateUtil.GetGenPath(genTable)
	//生成代码文件
	err = generator.TemplateUtil.GenCodeFiles(tplCodeMap, genTable.ModuleName, basePath)
	if e = response.CheckErr(err, "GenCode GenCodeFiles err"); e != nil {
		return
	}
	return
}

// genZipCode 生成代码 (压缩包下载)
func (genSrv generateService) genZipCode(zipWriter *zip.Writer, tableName string) (e error) {
	table, err := genSrv.findTableByName(tableName)
	if e = response.CheckErrDBNotRecord(err, "记录丢失！"); e != nil {
		return
	}
	if e = response.CheckErr(err, "genZipCode First err"); e != nil {
		return
	}
	genTable := legacyTableFromMakeadmin(table)
	//获取模板内容
	tplCodeMap, err := genSrv.renderCodeByTable(genTable)
	if e = response.CheckErr(err, "genZipCode renderCodeByTable err"); e != nil {
		return
	}
	//压缩文件
	err = generator.TemplateUtil.GenZip(zipWriter, tplCodeMap, genTable.ModuleName)
	if e = response.CheckErr(err, "genZipCode GenZip err"); e != nil {
		return
	}
	return
}

// DownloadCode 下载代码
func (genSrv generateService) DownloadCode(tableNames []string) ([]byte, error) {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)
	for _, tableName := range tableNames {
		err := genSrv.genZipCode(zipWriter, tableName)
		if err != nil {
			return nil, response.CheckErr(err, "DownloadCode genZipCode for %s err", tableName)
		}
	}
	err := zipWriter.Close()
	if err != nil {
		return nil, response.CheckErr(err, "DownloadCode zipWriter.Close err")
	}
	return buf.Bytes(), nil
}

func (genSrv generateService) findTable(id uint) (makeadmin.CodegenTable, error) {
	var table makeadmin.CodegenTable
	err := genSrv.db.
		Where("tenant_id = ? AND id = ? AND delete_time = ?", makeadmin.GlobalTenantID, id, 0).
		Limit(1).
		First(&table).
		Error
	return table, err
}

func (genSrv generateService) findTableByName(tableName string) (makeadmin.CodegenTable, error) {
	var table makeadmin.CodegenTable
	err := genSrv.db.
		Where("tenant_id = ? AND table_name = ? AND delete_time = ?", makeadmin.GlobalTenantID, strings.TrimSpace(tableName), 0).
		Order("id desc").
		Limit(1).
		First(&table).
		Error
	return table, err
}

func (genSrv generateService) listColumns(tableID uint) ([]makeadmin.CodegenColumn, error) {
	var columns []makeadmin.CodegenColumn
	err := genSrv.db.
		Where("table_id = ?", tableID).
		Order("sort").
		Find(&columns).
		Error
	return columns, err
}

func normalizePage(page request.PageReq) request.PageReq {
	if page.PageNo <= 0 {
		page.PageNo = 1
	}
	if page.PageSize <= 0 {
		page.PageSize = 20
	}
	if page.PageSize > 60 {
		page.PageSize = 60
	}
	return page
}

func genTableRespFromMakeadmin(table makeadmin.CodegenTable) resp.GenTableResp {
	return resp.GenTableResp{
		ID:           uint(table.ID),
		GenType:      legacyGenTypeFromMakeadmin(table.GenerateType),
		TableName:    table.SourceTable,
		TableComment: table.TableComment,
		CreateTime:   core.TsTime(table.CreateTime),
		UpdateTime:   core.TsTime(table.UpdateTime),
	}
}

func genTableDetailRespFromMakeadmin(table makeadmin.CodegenTable, columns []makeadmin.CodegenColumn) resp.GenTableDetailResp {
	legacyTable := legacyTableFromMakeadmin(table)
	columnResp := make([]resp.GenColumnResp, 0, len(columns))
	for _, column := range columns {
		columnResp = append(columnResp, genColumnRespFromMakeadmin(column))
	}
	return resp.GenTableDetailResp{
		Base: resp.GenTableBaseResp{
			ID:           uint(table.ID),
			TableName:    table.SourceTable,
			TableComment: table.TableComment,
			EntityName:   table.EntityName,
			AuthorName:   table.AuthorName,
			Remarks:      table.Remark,
			CreateTime:   core.TsTime(table.CreateTime),
			UpdateTime:   core.TsTime(table.UpdateTime),
		},
		Gen: resp.GenTableGenResp{
			GenTpl:       legacyTable.GenTpl,
			GenType:      legacyTable.GenType,
			GenPath:      legacyTable.GenPath,
			ModuleName:   legacyTable.ModuleName,
			FunctionName: legacyTable.FunctionName,
			TreePrimary:  legacyTable.TreePrimary,
			TreeParent:   legacyTable.TreeParent,
			TreeName:     legacyTable.TreeName,
			SubTableName: legacyTable.SubTableName,
			SubTableFk:   legacyTable.SubTableFk,
		},
		Column: columnResp,
	}
}

func genColumnRespFromMakeadmin(column makeadmin.CodegenColumn) resp.GenColumnResp {
	return resp.GenColumnResp{
		ID:            uint(column.ID),
		ColumnName:    column.ColumnName,
		ColumnComment: column.ColumnComment,
		ColumnLength:  column.ColumnLength,
		ColumnType:    column.ColumnType,
		JavaType:      column.GoType,
		JavaField:     column.GoField,
		IsRequired:    column.IsRequired,
		IsInsert:      column.IsInsert,
		IsEdit:        column.IsEdit,
		IsList:        column.IsList,
		IsQuery:       column.IsQuery,
		QueryType:     column.QueryType,
		HtmlType:      column.HTMLType,
		DictType:      column.DictType,
		CreateTime:    core.TsTime(column.CreateTime),
		UpdateTime:    core.TsTime(column.UpdateTime),
	}
}

func legacyTableFromMakeadmin(table makeadmin.CodegenTable) legacygen.GenTable {
	options := decodeTableOptions(table.Options)
	return legacygen.GenTable{
		ID:           uint(table.ID),
		TableName:    table.SourceTable,
		TableComment: table.TableComment,
		SubTableName: options.SubTableName,
		SubTableFk:   options.SubTableFk,
		AuthorName:   table.AuthorName,
		EntityName:   table.EntityName,
		ModuleName:   table.ModuleName,
		FunctionName: table.FunctionName,
		TreePrimary:  options.TreePrimary,
		TreeParent:   options.TreeParent,
		TreeName:     options.TreeName,
		GenTpl:       normalizeTemplateType(table.TemplateType),
		GenType:      legacyGenTypeFromMakeadmin(table.GenerateType),
		GenPath:      normalizeGenPath(table.GeneratePath),
		Remarks:      table.Remark,
		CreateTime:   table.CreateTime,
		UpdateTime:   table.UpdateTime,
	}
}

func makeadminTableFromLegacy(table legacygen.GenTable) makeadmin.CodegenTable {
	moduleName := strings.TrimSpace(table.ModuleName)
	if moduleName == "" {
		moduleName = generator.GenUtil.ToModuleName(table.TableName)
	}
	entityName := strings.TrimSpace(table.EntityName)
	if entityName == "" {
		entityName = generator.GenUtil.ToClassName(table.TableName)
	}
	functionName := strings.TrimSpace(table.FunctionName)
	if functionName == "" {
		functionName = strings.Replace(table.TableComment, "表", "", -1)
	}
	return makeadmin.CodegenTable{
		ID:           uint64(table.ID),
		TenantID:     makeadmin.GlobalTenantID,
		SourceTable:  table.TableName,
		TableComment: table.TableComment,
		ModuleName:   moduleName,
		PackageName:  config.GenConfig.PackageName,
		BusinessName: moduleName,
		EntityName:   entityName,
		FunctionName: functionName,
		AuthorName:   table.AuthorName,
		TemplateType: normalizeTemplateType(table.GenTpl),
		GenerateType: makeadminGenTypeFromLegacy(table.GenType),
		GeneratePath: normalizeGenPath(table.GenPath),
		Options: encodeTableOptions(codegenTableOptions{
			TreePrimary:  table.TreePrimary,
			TreeParent:   table.TreeParent,
			TreeName:     table.TreeName,
			SubTableName: table.SubTableName,
			SubTableFk:   table.SubTableFk,
		}),
		Remark:     table.Remarks,
		CreateTime: table.CreateTime,
		UpdateTime: table.UpdateTime,
	}
}

func legacyColumnFromMakeadmin(column makeadmin.CodegenColumn) legacygen.GenTableColumn {
	return legacygen.GenTableColumn{
		ID:            uint(column.ID),
		TableID:       uint(column.TableID),
		ColumnName:    column.ColumnName,
		ColumnComment: column.ColumnComment,
		ColumnLength:  column.ColumnLength,
		ColumnType:    column.ColumnType,
		JavaType:      column.GoType,
		JavaField:     column.GoField,
		IsPk:          column.IsPK,
		IsIncrement:   column.IsIncrement,
		IsRequired:    column.IsRequired,
		IsInsert:      column.IsInsert,
		IsEdit:        column.IsEdit,
		IsList:        column.IsList,
		IsQuery:       column.IsQuery,
		QueryType:     column.QueryType,
		HtmlType:      column.HTMLType,
		DictType:      column.DictType,
		Sort:          int(column.Sort),
		CreateTime:    column.CreateTime,
		UpdateTime:    column.UpdateTime,
	}
}

func legacyColumnsFromMakeadmin(columns []makeadmin.CodegenColumn) []legacygen.GenTableColumn {
	result := make([]legacygen.GenTableColumn, 0, len(columns))
	for _, column := range columns {
		result = append(result, legacyColumnFromMakeadmin(column))
	}
	return result
}

func makeadminColumnFromLegacy(column legacygen.GenTableColumn) makeadmin.CodegenColumn {
	goField := column.JavaField
	if goField == "" {
		goField = column.ColumnName
	}
	goType := column.JavaType
	if goType == "" {
		goType = generator.GoConstants.TypeString
	}
	return makeadmin.CodegenColumn{
		ID:            uint64(column.ID),
		TableID:       uint64(column.TableID),
		ColumnName:    column.ColumnName,
		ColumnComment: column.ColumnComment,
		ColumnType:    column.ColumnType,
		ColumnLength:  column.ColumnLength,
		GoType:        goType,
		GoField:       goField,
		JSONField:     util.StringUtil.ToSnakeCase(goField),
		IsPK:          column.IsPk,
		IsIncrement:   column.IsIncrement,
		IsRequired:    column.IsRequired,
		IsInsert:      column.IsInsert,
		IsEdit:        column.IsEdit,
		IsList:        column.IsList,
		IsQuery:       column.IsQuery,
		QueryType:     column.QueryType,
		HTMLType:      column.HtmlType,
		DictType:      column.DictType,
		Sort:          uint16(column.Sort),
		CreateTime:    column.CreateTime,
		UpdateTime:    column.UpdateTime,
	}
}

func applyEditReqToMakeadminTable(table *makeadmin.CodegenTable, editReq req.EditTableReq) {
	subTableName := strings.Replace(editReq.SubTableName, config.Config.DbTablePrefix, "", 1)
	table.SourceTable = editReq.TableName
	table.TableComment = editReq.TableComment
	table.ModuleName = editReq.ModuleName
	table.BusinessName = editReq.ModuleName
	table.PackageName = config.GenConfig.PackageName
	table.EntityName = editReq.EntityName
	table.FunctionName = editReq.FunctionName
	table.AuthorName = editReq.AuthorName
	table.TemplateType = normalizeTemplateType(editReq.GenTpl)
	table.GenerateType = makeadminGenTypeFromLegacy(editReq.GenType)
	table.GeneratePath = normalizeGenPath(editReq.GenPath)
	table.Options = encodeTableOptions(codegenTableOptions{
		TreePrimary:  editReq.TreePrimary,
		TreeParent:   editReq.TreeParent,
		TreeName:     editReq.TreeName,
		SubTableName: subTableName,
		SubTableFk:   editReq.SubTableFk,
	})
	table.Remark = editReq.Remarks
	table.UpdateTime = time.Now().Unix()
}

func applyEditReqToMakeadminColumn(column *makeadmin.CodegenColumn, editReq req.EditColumn) {
	column.ColumnComment = editReq.ColumnComment
	column.GoField = editReq.JavaField
	column.JSONField = util.StringUtil.ToSnakeCase(editReq.JavaField)
	column.IsRequired = editReq.IsRequired
	column.IsInsert = editReq.IsInsert
	column.IsEdit = editReq.IsEdit
	column.IsList = editReq.IsList
	column.IsQuery = editReq.IsQuery
	column.QueryType = editReq.QueryType
	column.HTMLType = editReq.HtmlType
	column.DictType = editReq.DictType
	column.UpdateTime = time.Now().Unix()
}

func decodeTableOptions(value string) codegenTableOptions {
	var options codegenTableOptions
	if strings.TrimSpace(value) == "" {
		return options
	}
	_ = json.Unmarshal([]byte(value), &options)
	return options
}

func encodeTableOptions(options codegenTableOptions) string {
	body, err := json.Marshal(options)
	if err != nil {
		return "{}"
	}
	return string(body)
}

func normalizeTemplateType(templateType string) string {
	if strings.TrimSpace(templateType) == makeadmin.CodegenTemplateTree {
		return makeadmin.CodegenTemplateTree
	}
	return makeadmin.CodegenTemplateCRUD
}

func normalizeGenPath(genPath string) string {
	if strings.TrimSpace(genPath) == "" {
		return defaultGenPath
	}
	return genPath
}

func makeadminGenTypeFromLegacy(genType int) string {
	if genType == legacyGenTypePath {
		return makeadminGenerateTypePath
	}
	return makeadminGenerateTypeZip
}

func legacyGenTypeFromMakeadmin(genType string) int {
	switch strings.ToLower(strings.TrimSpace(genType)) {
	case makeadminGenerateTypePath, "custom", "custom_path":
		return legacyGenTypePath
	default:
		return legacyGenTypeZip
	}
}

func uintSliceToUint64(ids []uint) []uint64 {
	result := make([]uint64, 0, len(ids))
	for _, id := range ids {
		result = append(result, uint64(id))
	}
	return result
}
