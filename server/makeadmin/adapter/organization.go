package adapter

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"go-makeadmin/admin/schemas/req"
	"go-makeadmin/admin/schemas/resp"
	"go-makeadmin/core"
	"go-makeadmin/core/request"
	"go-makeadmin/core/response"
	"go-makeadmin/makeadmin/repository"
	makeadminsvc "go-makeadmin/makeadmin/service"
	"go-makeadmin/model/makeadmin"
	"go-makeadmin/util"
)

type OrgUnitAdapter interface {
	Available(ctx context.Context) bool
	All(ctx context.Context) ([]resp.SystemAuthDeptResp, error)
	List(ctx context.Context, listReq req.SystemAuthDeptListReq) ([]interface{}, error)
	Detail(ctx context.Context, id uint) (resp.SystemAuthDeptResp, error)
	Add(ctx context.Context, addReq req.SystemAuthDeptAddReq) error
	Edit(ctx context.Context, editReq req.SystemAuthDeptEditReq) error
	Del(ctx context.Context, id uint) error
}

type PositionAdapter interface {
	Available(ctx context.Context) bool
	All(ctx context.Context) ([]resp.SystemAuthPostResp, error)
	List(ctx context.Context, page request.PageReq, listReq req.SystemAuthPostListReq) (response.PageResp, error)
	Detail(ctx context.Context, id uint) (resp.SystemAuthPostResp, error)
	Add(ctx context.Context, addReq req.SystemAuthPostAddReq) error
	Edit(ctx context.Context, editReq req.SystemAuthPostEditReq) error
	Del(ctx context.Context, id uint) error
}

type orgUnitAdapter struct {
	db *gorm.DB
}

type positionAdapter struct {
	db *gorm.DB
}

func NewOrgUnitAdapter(db *gorm.DB) OrgUnitAdapter {
	return orgUnitAdapter{db: db}
}

func NewPositionAdapter(db *gorm.DB) PositionAdapter {
	return positionAdapter{db: db}
}

func (adapter orgUnitAdapter) Available(ctx context.Context) bool {
	if adapter.db == nil ||
		!adapter.db.Migrator().HasTable(&makeadmin.OrgUnit{}) ||
		!adapter.db.Migrator().HasTable(&makeadmin.AdminOrg{}) {
		return false
	}
	var count int64
	err := adapter.db.WithContext(ctx).
		Model(&makeadmin.OrgUnit{}).
		Where("tenant_id = ? AND delete_time = ?", makeadmin.GlobalTenantID, 0).
		Count(&count).
		Error
	return err == nil && count > 0
}

func (adapter orgUnitAdapter) All(ctx context.Context) ([]resp.SystemAuthDeptResp, error) {
	items, err := adapter.orgService().List(ctx, makeadmin.GlobalTenantID, repository.OrgUnitFilter{})
	if err != nil {
		return nil, mapOrgUnitError(err)
	}
	result := make([]resp.SystemAuthDeptResp, 0, len(items))
	for _, item := range items {
		if item.ParentID == 0 {
			continue
		}
		result = append(result, deptResponse(item))
	}
	return result, nil
}

func (adapter orgUnitAdapter) List(ctx context.Context, listReq req.SystemAuthDeptListReq) ([]interface{}, error) {
	filter := repository.OrgUnitFilter{
		Name:      listReq.Name,
		Status:    statusFromListStop(listReq.IsStop),
		StatusSet: listReq.IsStop >= 0,
	}
	items, err := adapter.orgService().List(ctx, makeadmin.GlobalTenantID, filter)
	if err != nil {
		return nil, mapOrgUnitError(err)
	}
	return util.ArrayUtil.ListToTree(
		util.ConvertUtil.StructsToMaps(deptResponses(items)), "id", "pid", "children"), nil
}

func (adapter orgUnitAdapter) Detail(ctx context.Context, id uint) (resp.SystemAuthDeptResp, error) {
	item, err := adapter.orgService().Detail(ctx, makeadmin.GlobalTenantID, uint64(id))
	if err != nil {
		return resp.SystemAuthDeptResp{}, mapOrgUnitError(err)
	}
	return deptResponse(item), nil
}

func (adapter orgUnitAdapter) Add(ctx context.Context, addReq req.SystemAuthDeptAddReq) error {
	return mapOrgUnitError(adapter.orgService().Add(ctx, makeadminsvc.OrgUnitInput{
		TenantID: makeadmin.GlobalTenantID,
		ParentID: uint64(addReq.Pid),
		Name:     addReq.Name,
		IsStop:   addReq.IsStop,
		Sort:     addReq.Sort,
	}))
}

func (adapter orgUnitAdapter) Edit(ctx context.Context, editReq req.SystemAuthDeptEditReq) error {
	return mapOrgUnitError(adapter.orgService().Edit(ctx, makeadminsvc.OrgUnitInput{
		ID:       uint64(editReq.ID),
		TenantID: makeadmin.GlobalTenantID,
		ParentID: uint64(editReq.Pid),
		Name:     editReq.Name,
		IsStop:   editReq.IsStop,
		Sort:     editReq.Sort,
	}))
}

func (adapter orgUnitAdapter) Del(ctx context.Context, id uint) error {
	return mapOrgUnitError(adapter.orgService().Delete(ctx, makeadmin.GlobalTenantID, uint64(id)))
}

func (adapter orgUnitAdapter) orgService() makeadminsvc.OrgUnitService {
	return makeadminsvc.NewOrgUnitService(repository.NewOrgUnitRepository(adapter.db))
}

func (adapter positionAdapter) Available(ctx context.Context) bool {
	if adapter.db == nil ||
		!adapter.db.Migrator().HasTable(&makeadmin.Position{}) ||
		!adapter.db.Migrator().HasTable(&makeadmin.AdminOrg{}) {
		return false
	}
	var count int64
	err := adapter.db.WithContext(ctx).
		Model(&makeadmin.Position{}).
		Where("tenant_id = ? AND delete_time = ?", makeadmin.GlobalTenantID, 0).
		Count(&count).
		Error
	return err == nil && count > 0
}

func (adapter positionAdapter) All(ctx context.Context) ([]resp.SystemAuthPostResp, error) {
	items, err := adapter.positionService().ListAll(ctx, makeadmin.GlobalTenantID)
	if err != nil {
		return nil, mapPositionError(err)
	}
	return postResponses(items), nil
}

func (adapter positionAdapter) List(ctx context.Context, page request.PageReq, listReq req.SystemAuthPostListReq) (response.PageResp, error) {
	page = normalizePage(page)
	result, err := adapter.positionService().List(ctx, makeadmin.GlobalTenantID, repository.PositionFilter{
		Code:      listReq.Code,
		Name:      listReq.Name,
		Status:    statusFromListStop(listReq.IsStop),
		StatusSet: listReq.IsStop >= 0,
	}, page.PageNo, page.PageSize)
	if err != nil {
		return response.PageResp{}, mapPositionError(err)
	}
	return response.PageResp{
		PageNo:   page.PageNo,
		PageSize: page.PageSize,
		Count:    result.Count,
		Lists:    postResponses(result.Items),
	}, nil
}

func (adapter positionAdapter) Detail(ctx context.Context, id uint) (resp.SystemAuthPostResp, error) {
	item, err := adapter.positionService().Detail(ctx, makeadmin.GlobalTenantID, uint64(id))
	if err != nil {
		return resp.SystemAuthPostResp{}, mapPositionError(err)
	}
	return postResponse(item), nil
}

func (adapter positionAdapter) Add(ctx context.Context, addReq req.SystemAuthPostAddReq) error {
	return mapPositionError(adapter.positionService().Add(ctx, makeadminsvc.PositionInput{
		TenantID: makeadmin.GlobalTenantID,
		Code:     addReq.Code,
		Name:     addReq.Name,
		Remark:   addReq.Remarks,
		IsStop:   addReq.IsStop,
		Sort:     addReq.Sort,
	}))
}

func (adapter positionAdapter) Edit(ctx context.Context, editReq req.SystemAuthPostEditReq) error {
	return mapPositionError(adapter.positionService().Edit(ctx, makeadminsvc.PositionInput{
		ID:       uint64(editReq.ID),
		TenantID: makeadmin.GlobalTenantID,
		Code:     editReq.Code,
		Name:     editReq.Name,
		Remark:   editReq.Remarks,
		IsStop:   editReq.IsStop,
		Sort:     editReq.Sort,
	}))
}

func (adapter positionAdapter) Del(ctx context.Context, id uint) error {
	return mapPositionError(adapter.positionService().Delete(ctx, makeadmin.GlobalTenantID, uint64(id)))
}

func (adapter positionAdapter) positionService() makeadminsvc.PositionService {
	return makeadminsvc.NewPositionService(repository.NewPositionRepository(adapter.db))
}

func deptResponses(items []makeadminsvc.OrgUnitItem) []resp.SystemAuthDeptResp {
	result := make([]resp.SystemAuthDeptResp, 0, len(items))
	for _, item := range items {
		result = append(result, deptResponse(item))
	}
	return result
}

func deptResponse(item makeadminsvc.OrgUnitItem) resp.SystemAuthDeptResp {
	return resp.SystemAuthDeptResp{
		ID:         uint(item.ID),
		Pid:        uint(item.ParentID),
		Name:       item.Name,
		Sort:       item.Sort,
		IsStop:     item.IsStop,
		CreateTime: core.TsTime(item.CreateTime),
		UpdateTime: core.TsTime(item.UpdateTime),
	}
}

func postResponses(items []makeadminsvc.PositionItem) []resp.SystemAuthPostResp {
	result := make([]resp.SystemAuthPostResp, 0, len(items))
	for _, item := range items {
		result = append(result, postResponse(item))
	}
	return result
}

func postResponse(item makeadminsvc.PositionItem) resp.SystemAuthPostResp {
	return resp.SystemAuthPostResp{
		ID:         uint(item.ID),
		Code:       item.Code,
		Name:       item.Name,
		Remarks:    item.Remark,
		Sort:       item.Sort,
		IsStop:     item.IsStop,
		CreateTime: core.TsTime(item.CreateTime),
		UpdateTime: core.TsTime(item.UpdateTime),
	}
}

func statusFromListStop(isStop int8) uint8 {
	if isStop == 1 {
		return makeadmin.StatusDisabled
	}
	return makeadmin.StatusEnabled
}

func mapOrgUnitError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, makeadminsvc.ErrOrgUnitNotFound):
		return response.AssertArgumentError.Make("部门不存在!")
	case errors.Is(err, makeadminsvc.ErrParentOrgUnitNotFound):
		return response.AssertArgumentError.Make("上级部门不存在!")
	case errors.Is(err, makeadminsvc.ErrRootOrgUnitExists):
		return response.AssertArgumentError.Make("顶级部门只允许有一个!")
	case errors.Is(err, makeadminsvc.ErrRootOrgUnitParentLocked):
		return response.AssertArgumentError.Make("顶级部门不能修改上级!")
	case errors.Is(err, makeadminsvc.ErrRootOrgUnitDeleteLocked):
		return response.AssertArgumentError.Make("顶级部门不能删除!")
	case errors.Is(err, makeadminsvc.ErrOrgUnitSelfParent):
		return response.AssertArgumentError.Make("上级部门不能是自己!")
	case errors.Is(err, makeadminsvc.ErrOrgUnitHasChildren):
		return response.AssertArgumentError.Make("请先删除子级部门!")
	case errors.Is(err, makeadminsvc.ErrOrgUnitInUse):
		return response.AssertArgumentError.Make("该部门已被管理员使用,请先移除!")
	default:
		return err
	}
}

func mapPositionError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, makeadminsvc.ErrPositionNotFound):
		return response.AssertArgumentError.Make("岗位不存在!")
	case errors.Is(err, makeadminsvc.ErrPositionCodeExists), errors.Is(err, makeadminsvc.ErrPositionNameExists):
		return response.AssertArgumentError.Make("该岗位已存在!")
	case errors.Is(err, makeadminsvc.ErrPositionInUse):
		return response.AssertArgumentError.Make("该岗位已被管理员使用,请先移除!")
	default:
		return err
	}
}
