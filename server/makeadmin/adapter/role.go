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
)

type RoleAdapter interface {
	Available(ctx context.Context) bool
	All(ctx context.Context) ([]resp.SystemAuthRoleSimpleResp, error)
	List(ctx context.Context, page request.PageReq) (response.PageResp, error)
	Detail(ctx context.Context, id uint) (resp.SystemAuthRoleResp, error)
	Add(ctx context.Context, addReq req.SystemAuthRoleAddReq) error
	Edit(ctx context.Context, editReq req.SystemAuthRoleEditReq) error
	Del(ctx context.Context, id uint) error
}

type roleAdapter struct {
	db *gorm.DB
}

func NewRoleAdapter(db *gorm.DB) RoleAdapter {
	return roleAdapter{db: db}
}

func (adapter roleAdapter) Available(ctx context.Context) bool {
	if adapter.db == nil ||
		!adapter.db.Migrator().HasTable(&makeadmin.Role{}) ||
		!adapter.db.Migrator().HasTable(&makeadmin.RolePermission{}) ||
		!adapter.db.Migrator().HasTable(&makeadmin.MenuPermission{}) {
		return false
	}
	var count int64
	err := adapter.db.WithContext(ctx).
		Model(&makeadmin.Role{}).
		Where("tenant_id = ? AND delete_time = ?", tenantIDFromContext(ctx), 0).
		Count(&count).
		Error
	return err == nil && count > 0
}

func (adapter roleAdapter) All(ctx context.Context) ([]resp.SystemAuthRoleSimpleResp, error) {
	items, err := adapter.roleService().ListAll(ctx, tenantIDFromContext(ctx))
	if err != nil {
		return nil, mapRoleError(err)
	}
	result := make([]resp.SystemAuthRoleSimpleResp, 0, len(items))
	for _, item := range items {
		result = append(result, resp.SystemAuthRoleSimpleResp{
			ID:         uint(item.ID),
			Name:       item.Name,
			CreateTime: core.TsTime(item.CreateTime),
			UpdateTime: core.TsTime(item.UpdateTime),
		})
	}
	return result, nil
}

func (adapter roleAdapter) List(ctx context.Context, page request.PageReq) (response.PageResp, error) {
	page = normalizePage(page)
	result, err := adapter.roleService().List(ctx, tenantIDFromContext(ctx), page.PageNo, page.PageSize)
	if err != nil {
		return response.PageResp{}, mapRoleError(err)
	}
	return response.PageResp{
		PageNo:   page.PageNo,
		PageSize: page.PageSize,
		Count:    result.Count,
		Lists:    roleResponses(result.Items),
	}, nil
}

func (adapter roleAdapter) Detail(ctx context.Context, id uint) (resp.SystemAuthRoleResp, error) {
	item, err := adapter.roleService().Detail(ctx, tenantIDFromContext(ctx), uint64(id))
	if err != nil {
		return resp.SystemAuthRoleResp{}, mapRoleError(err)
	}
	return roleResponse(item), nil
}

func (adapter roleAdapter) Add(ctx context.Context, addReq req.SystemAuthRoleAddReq) error {
	tenantID := tenantIDFromContext(ctx)
	return mapRoleError(adapter.roleService().Add(ctx, makeadminsvc.RoleInput{
		TenantID:  tenantID,
		Name:      addReq.Name,
		Remark:    addReq.Remark,
		Sort:      addReq.Sort,
		IsDisable: addReq.IsDisable,
		MenuIDs:   makeadminsvc.ParseRoleMenuIDs(addReq.MenuIds),
	}))
}

func (adapter roleAdapter) Edit(ctx context.Context, editReq req.SystemAuthRoleEditReq) error {
	tenantID := tenantIDFromContext(ctx)
	return mapRoleError(adapter.roleService().Edit(ctx, makeadminsvc.RoleInput{
		ID:        uint64(editReq.ID),
		TenantID:  tenantID,
		Name:      editReq.Name,
		Remark:    editReq.Remark,
		Sort:      editReq.Sort,
		IsDisable: editReq.IsDisable,
		MenuIDs:   makeadminsvc.ParseRoleMenuIDs(editReq.MenuIds),
	}))
}

func (adapter roleAdapter) Del(ctx context.Context, id uint) error {
	return mapRoleError(adapter.roleService().Delete(ctx, tenantIDFromContext(ctx), uint64(id)))
}

func (adapter roleAdapter) roleService() makeadminsvc.RoleService {
	return makeadminsvc.NewRoleService(repository.NewRoleRepository(adapter.db))
}

func roleResponses(items []makeadminsvc.RoleItem) []resp.SystemAuthRoleResp {
	result := make([]resp.SystemAuthRoleResp, 0, len(items))
	for _, item := range items {
		result = append(result, roleResponse(item))
	}
	return result
}

func roleResponse(item makeadminsvc.RoleItem) resp.SystemAuthRoleResp {
	return resp.SystemAuthRoleResp{
		ID:         uint(item.ID),
		Name:       item.Name,
		Remark:     item.Remark,
		Menus:      roleMenuIDs(item.MenuIDs),
		Member:     item.Member,
		Sort:       item.Sort,
		IsDisable:  item.IsDisable,
		CreateTime: core.TsTime(item.CreateTime),
		UpdateTime: core.TsTime(item.UpdateTime),
	}
}

func roleMenuIDs(ids []uint64) []uint {
	result := make([]uint, 0, len(ids))
	for _, id := range ids {
		result = append(result, uint(id))
	}
	return result
}

func mapRoleError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, makeadminsvc.ErrRoleNotFound):
		return response.AssertArgumentError.Make("角色已不存在!")
	case errors.Is(err, makeadminsvc.ErrRoleNameExists):
		return response.AssertArgumentError.Make("角色名称已存在!")
	case errors.Is(err, makeadminsvc.ErrRoleInUse):
		return response.AssertArgumentError.Make("角色已被管理员使用,请先移除!")
	case errors.Is(err, makeadminsvc.ErrSystemRoleProtected):
		return response.AssertArgumentError.Make("系统角色不能删除!")
	default:
		return err
	}
}
