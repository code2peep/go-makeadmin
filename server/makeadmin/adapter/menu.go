package adapter

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"go-makeadmin/admin/schemas/req"
	"go-makeadmin/admin/schemas/resp"
	"go-makeadmin/core"
	"go-makeadmin/core/response"
	"go-makeadmin/makeadmin/repository"
	makeadminsvc "go-makeadmin/makeadmin/service"
	"go-makeadmin/model/makeadmin"
	"go-makeadmin/util"
)

type MenuAdapter interface {
	Available(ctx context.Context) bool
	List(ctx context.Context) ([]interface{}, error)
	Detail(ctx context.Context, id uint) (resp.SystemAuthMenuResp, error)
	Add(ctx context.Context, addReq req.SystemAuthMenuAddReq) error
	Edit(ctx context.Context, editReq req.SystemAuthMenuEditReq) error
	Del(ctx context.Context, id uint) error
}

type menuAdapter struct {
	db *gorm.DB
}

func NewMenuAdapter(db *gorm.DB) MenuAdapter {
	return menuAdapter{db: db}
}

func (adapter menuAdapter) Available(ctx context.Context) bool {
	if adapter.db == nil ||
		!adapter.db.Migrator().HasTable(&makeadmin.Menu{}) ||
		!adapter.db.Migrator().HasTable(&makeadmin.Permission{}) ||
		!adapter.db.Migrator().HasTable(&makeadmin.MenuPermission{}) {
		return false
	}
	var count int64
	err := adapter.db.WithContext(ctx).
		Model(&makeadmin.Menu{}).
		Where("delete_time = ?", 0).
		Count(&count).
		Error
	return err == nil && count > 0
}

func (adapter menuAdapter) List(ctx context.Context) ([]interface{}, error) {
	items, err := adapter.menuService().List(ctx)
	if err != nil {
		return nil, mapMenuError(err)
	}
	return util.ArrayUtil.ListToTree(
		util.ConvertUtil.StructsToMaps(menuResponses(items)), "id", "pid", "children"), nil
}

func (adapter menuAdapter) Detail(ctx context.Context, id uint) (resp.SystemAuthMenuResp, error) {
	item, err := adapter.menuService().Detail(ctx, uint64(id))
	if err != nil {
		return resp.SystemAuthMenuResp{}, mapMenuError(err)
	}
	return menuResponse(item), nil
}

func (adapter menuAdapter) Add(ctx context.Context, addReq req.SystemAuthMenuAddReq) error {
	return mapMenuError(adapter.menuService().Add(ctx, makeadminsvc.MenuInput{
		ParentID:  uint64(addReq.Pid),
		MenuType:  addReq.MenuType,
		MenuName:  addReq.MenuName,
		MenuIcon:  addReq.MenuIcon,
		MenuSort:  addReq.MenuSort,
		Perms:     addReq.Perms,
		Paths:     addReq.Paths,
		Component: addReq.Component,
		Selected:  addReq.Selected,
		Params:    addReq.Params,
		IsCache:   addReq.IsCache,
		IsShow:    addReq.IsShow,
		IsDisable: addReq.IsDisable,
	}))
}

func (adapter menuAdapter) Edit(ctx context.Context, editReq req.SystemAuthMenuEditReq) error {
	return mapMenuError(adapter.menuService().Edit(ctx, makeadminsvc.MenuInput{
		ID:        uint64(editReq.ID),
		ParentID:  uint64(editReq.Pid),
		MenuType:  editReq.MenuType,
		MenuName:  editReq.MenuName,
		MenuIcon:  editReq.MenuIcon,
		MenuSort:  editReq.MenuSort,
		Perms:     editReq.Perms,
		Paths:     editReq.Paths,
		Component: editReq.Component,
		Selected:  editReq.Selected,
		Params:    editReq.Params,
		IsCache:   editReq.IsCache,
		IsShow:    editReq.IsShow,
		IsDisable: editReq.IsDisable,
	}))
}

func (adapter menuAdapter) Del(ctx context.Context, id uint) error {
	return mapMenuError(adapter.menuService().Delete(ctx, uint64(id)))
}

func (adapter menuAdapter) menuService() makeadminsvc.MenuService {
	return makeadminsvc.NewMenuService(repository.NewMenuRepository(adapter.db))
}

func menuResponses(items []makeadminsvc.MenuItem) []resp.SystemAuthMenuResp {
	result := make([]resp.SystemAuthMenuResp, 0, len(items))
	for _, item := range items {
		result = append(result, menuResponse(item))
	}
	return result
}

func menuResponse(item makeadminsvc.MenuItem) resp.SystemAuthMenuResp {
	return resp.SystemAuthMenuResp{
		ID:         uint(item.ID),
		Pid:        uint(item.ParentID),
		MenuType:   item.MenuType,
		MenuName:   item.MenuName,
		MenuIcon:   item.MenuIcon,
		MenuSort:   item.MenuSort,
		Perms:      item.Perms,
		Paths:      item.Paths,
		Component:  item.Component,
		Selected:   item.Selected,
		Params:     item.Params,
		IsCache:    item.IsCache,
		IsShow:     item.IsShow,
		IsDisable:  item.IsDisable,
		CreateTime: core.TsTime(item.CreateTime),
		UpdateTime: core.TsTime(item.UpdateTime),
	}
}

func mapMenuError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, makeadminsvc.ErrMenuNotFound):
		return response.AssertArgumentError.Make("菜单已不存在!")
	case errors.Is(err, makeadminsvc.ErrParentMenuNotFound):
		return response.AssertArgumentError.Make("父级菜单已不存在!")
	case errors.Is(err, makeadminsvc.ErrMenuSelfParent):
		return response.AssertArgumentError.Make("父级菜单不能是自己!")
	case errors.Is(err, makeadminsvc.ErrMenuHasChildren):
		return response.AssertArgumentError.Make("请先删除子菜单再操作！")
	case errors.Is(err, makeadminsvc.ErrMenuPermsExists):
		return response.AssertArgumentError.Make("权限字符已存在!")
	default:
		return err
	}
}
