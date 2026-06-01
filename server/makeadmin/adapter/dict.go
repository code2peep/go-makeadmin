package adapter

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"go-makeadmin/admin/schemas/req"
	adminresp "go-makeadmin/admin/schemas/resp"
	"go-makeadmin/core"
	"go-makeadmin/core/request"
	"go-makeadmin/core/response"
	"go-makeadmin/makeadmin/repository"
	makeadminsvc "go-makeadmin/makeadmin/service"
	"go-makeadmin/model/makeadmin"
)

type DictAdapter interface {
	Available(ctx context.Context) bool
	TypeAll(ctx context.Context) ([]adminresp.SettingDictTypeResp, error)
	TypeList(ctx context.Context, page request.PageReq, listReq req.SettingDictTypeListReq) (response.PageResp, error)
	TypeDetail(ctx context.Context, id uint) (adminresp.SettingDictTypeResp, error)
	TypeAdd(ctx context.Context, addReq req.SettingDictTypeAddReq) error
	TypeEdit(ctx context.Context, editReq req.SettingDictTypeEditReq) error
	TypeDel(ctx context.Context, delReq req.SettingDictTypeDelReq) error
	DataAll(ctx context.Context, allReq req.SettingDictDataListReq) ([]adminresp.SettingDictDataResp, error)
	DataList(ctx context.Context, page request.PageReq, listReq req.SettingDictDataListReq) (response.PageResp, error)
	DataDetail(ctx context.Context, id uint) (adminresp.SettingDictDataResp, error)
	DataAdd(ctx context.Context, addReq req.SettingDictDataAddReq) error
	DataEdit(ctx context.Context, editReq req.SettingDictDataEditReq) error
	DataDel(ctx context.Context, delReq req.SettingDictDataDelReq) error
}

type dictAdapter struct {
	db *gorm.DB
}

func NewDictAdapter(db *gorm.DB) DictAdapter {
	return dictAdapter{db: db}
}

func (adapter dictAdapter) Available(ctx context.Context) bool {
	if adapter.db == nil ||
		!adapter.db.Migrator().HasTable(&makeadmin.DictType{}) ||
		!adapter.db.Migrator().HasTable(&makeadmin.DictItem{}) {
		return false
	}
	var count int64
	err := adapter.db.WithContext(ctx).
		Model(&makeadmin.DictType{}).
		Where("delete_time = ?", 0).
		Count(&count).
		Error
	return err == nil && count > 0
}

func (adapter dictAdapter) TypeAll(ctx context.Context) ([]adminresp.SettingDictTypeResp, error) {
	items, err := adapter.dictService().AllTypes(ctx)
	if err != nil {
		return nil, mapDictError(err)
	}
	return dictTypeResponses(items), nil
}

func (adapter dictAdapter) TypeList(ctx context.Context, page request.PageReq, listReq req.SettingDictTypeListReq) (response.PageResp, error) {
	page = normalizePage(page)
	result, err := adapter.dictService().ListTypes(ctx, repository.DictTypeFilter{
		Name:      listReq.DictName,
		Code:      listReq.DictType,
		Status:    listReq.DictStatus,
		StatusSet: listReq.DictStatus >= 0,
	}, page.PageSize, page.PageSize*(page.PageNo-1))
	if err != nil {
		return response.PageResp{}, mapDictError(err)
	}
	return response.PageResp{
		PageNo:   page.PageNo,
		PageSize: page.PageSize,
		Count:    result.Count,
		Lists:    dictTypeResponses(result.Items),
	}, nil
}

func (adapter dictAdapter) TypeDetail(ctx context.Context, id uint) (adminresp.SettingDictTypeResp, error) {
	item, err := adapter.dictService().DetailType(ctx, uint64(id))
	if err != nil {
		return adminresp.SettingDictTypeResp{}, mapDictError(err)
	}
	return dictTypeResponse(item), nil
}

func (adapter dictAdapter) TypeAdd(ctx context.Context, addReq req.SettingDictTypeAddReq) error {
	return mapDictError(adapter.dictService().AddType(ctx, makeadminsvc.DictTypeInput{
		Name:   addReq.DictName,
		Code:   addReq.DictType,
		Remark: addReq.DictRemark,
		Status: addReq.DictStatus,
	}))
}

func (adapter dictAdapter) TypeEdit(ctx context.Context, editReq req.SettingDictTypeEditReq) error {
	return mapDictError(adapter.dictService().EditType(ctx, makeadminsvc.DictTypeInput{
		ID:     uint64(editReq.ID),
		Name:   editReq.DictName,
		Code:   editReq.DictType,
		Remark: editReq.DictRemark,
		Status: editReq.DictStatus,
	}))
}

func (adapter dictAdapter) TypeDel(ctx context.Context, delReq req.SettingDictTypeDelReq) error {
	return mapDictError(adapter.dictService().DeleteTypes(ctx, uintSliceToUint64(delReq.Ids)))
}

func (adapter dictAdapter) DataAll(ctx context.Context, allReq req.SettingDictDataListReq) ([]adminresp.SettingDictDataResp, error) {
	items, err := adapter.dictService().AllItems(ctx, allReq.DictType, repository.DictItemFilter{
		Name:      allReq.Name,
		Value:     allReq.Value,
		Status:    allReq.Status,
		StatusSet: allReq.Status >= 0,
	})
	if err != nil {
		return nil, mapDictError(err)
	}
	return dictDataResponses(items), nil
}

func (adapter dictAdapter) DataList(ctx context.Context, page request.PageReq, listReq req.SettingDictDataListReq) (response.PageResp, error) {
	page = normalizePage(page)
	result, err := adapter.dictService().ListItems(ctx, listReq.DictType, repository.DictItemFilter{
		Name:      listReq.Name,
		Value:     listReq.Value,
		Status:    listReq.Status,
		StatusSet: listReq.Status >= 0,
	}, page.PageSize, page.PageSize*(page.PageNo-1))
	if err != nil {
		return response.PageResp{}, mapDictError(err)
	}
	return response.PageResp{
		PageNo:   page.PageNo,
		PageSize: page.PageSize,
		Count:    result.Count,
		Lists:    dictDataResponses(result.Items),
	}, nil
}

func (adapter dictAdapter) DataDetail(ctx context.Context, id uint) (adminresp.SettingDictDataResp, error) {
	item, err := adapter.dictService().DetailItem(ctx, uint64(id))
	if err != nil {
		return adminresp.SettingDictDataResp{}, mapDictError(err)
	}
	return dictDataResponse(item), nil
}

func (adapter dictAdapter) DataAdd(ctx context.Context, addReq req.SettingDictDataAddReq) error {
	return mapDictError(adapter.dictService().AddItem(ctx, makeadminsvc.DictItemInput{
		TypeID: uint64(addReq.TypeId),
		Name:   addReq.Name,
		Value:  addReq.Value,
		Remark: addReq.Remark,
		Sort:   addReq.Sort,
		Status: addReq.Status,
	}))
}

func (adapter dictAdapter) DataEdit(ctx context.Context, editReq req.SettingDictDataEditReq) error {
	return mapDictError(adapter.dictService().EditItem(ctx, makeadminsvc.DictItemInput{
		ID:     uint64(editReq.ID),
		TypeID: uint64(editReq.TypeId),
		Name:   editReq.Name,
		Value:  editReq.Value,
		Remark: editReq.Remark,
		Sort:   editReq.Sort,
		Status: editReq.Status,
	}))
}

func (adapter dictAdapter) DataDel(ctx context.Context, delReq req.SettingDictDataDelReq) error {
	return mapDictError(adapter.dictService().DeleteItems(ctx, uintSliceToUint64(delReq.Ids)))
}

func (adapter dictAdapter) dictService() makeadminsvc.DictService {
	return makeadminsvc.NewDictService(repository.NewDictRepository(adapter.db))
}

func dictTypeResponses(items []makeadminsvc.DictType) []adminresp.SettingDictTypeResp {
	result := make([]adminresp.SettingDictTypeResp, 0, len(items))
	for _, item := range items {
		result = append(result, dictTypeResponse(item))
	}
	return result
}

func dictTypeResponse(item makeadminsvc.DictType) adminresp.SettingDictTypeResp {
	return adminresp.SettingDictTypeResp{
		ID:         uint(item.ID),
		DictName:   item.Name,
		DictType:   item.Code,
		DictRemark: item.Remark,
		DictStatus: item.Status,
		CreateTime: core.TsTime(item.CreateTime),
		UpdateTime: core.TsTime(item.UpdateTime),
	}
}

func dictDataResponses(items []makeadminsvc.DictItem) []adminresp.SettingDictDataResp {
	result := make([]adminresp.SettingDictDataResp, 0, len(items))
	for _, item := range items {
		result = append(result, dictDataResponse(item))
	}
	return result
}

func dictDataResponse(item makeadminsvc.DictItem) adminresp.SettingDictDataResp {
	return adminresp.SettingDictDataResp{
		ID:         uint(item.ID),
		TypeId:     uint(item.TypeID),
		Name:       item.Name,
		Value:      item.Value,
		Remark:     item.Remark,
		Sort:       item.Sort,
		Status:     item.Status,
		CreateTime: core.TsTime(item.CreateTime),
		UpdateTime: core.TsTime(item.UpdateTime),
	}
}

func mapDictError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, makeadminsvc.ErrDictTypeNotFound):
		return response.AssertArgumentError.Make("字典类型不存在！")
	case errors.Is(err, makeadminsvc.ErrDictTypeNameExists):
		return response.AssertArgumentError.Make("字典名称已存在！")
	case errors.Is(err, makeadminsvc.ErrDictTypeCodeExists):
		return response.AssertArgumentError.Make("字典类型已存在！")
	case errors.Is(err, makeadminsvc.ErrDictItemNotFound):
		return response.AssertArgumentError.Make("字典数据不存在！")
	case errors.Is(err, makeadminsvc.ErrDictItemValueExists):
		return response.AssertArgumentError.Make("字典数据值已存在！")
	default:
		return err
	}
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

func uintSliceToUint64(ids []uint) []uint64 {
	result := make([]uint64, 0, len(ids))
	for _, id := range ids {
		result = append(result, uint64(id))
	}
	return result
}
