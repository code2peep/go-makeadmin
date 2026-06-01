package service

import (
	"context"
	"errors"
	"testing"

	"gorm.io/gorm"

	"go-makeadmin/model/makeadmin"
)

type fakeMenuRepository struct {
	menus             []makeadmin.Menu
	permissions       []makeadmin.Permission
	menuPermissions   map[uint64][]string
	childCount        int64
	createdMenu       makeadmin.Menu
	createdPermission *makeadmin.Permission
	updatedMenu       makeadmin.Menu
	updatedPermission *makeadmin.Permission
	deletedMenuID     uint64
}

func (repo *fakeMenuRepository) ListMenus(ctx context.Context) ([]makeadmin.Menu, error) {
	return repo.menus, nil
}

func (repo *fakeMenuRepository) FindMenuByID(ctx context.Context, id uint64) (makeadmin.Menu, error) {
	for _, menu := range repo.menus {
		if menu.ID == id && menu.DeleteTime == 0 {
			return menu, nil
		}
	}
	return makeadmin.Menu{}, gorm.ErrRecordNotFound
}

func (repo *fakeMenuRepository) CountChildMenus(ctx context.Context, parentID uint64) (int64, error) {
	return repo.childCount, nil
}

func (repo *fakeMenuRepository) ListPermissionCodesByMenuID(ctx context.Context, menuID uint64) ([]string, error) {
	if repo.menuPermissions == nil {
		return []string{}, nil
	}
	return repo.menuPermissions[menuID], nil
}

func (repo *fakeMenuRepository) FindPermissionByCode(ctx context.Context, code string) (makeadmin.Permission, error) {
	for _, permission := range repo.permissions {
		if permission.Code == code {
			return permission, nil
		}
	}
	return makeadmin.Permission{}, gorm.ErrRecordNotFound
}

func (repo *fakeMenuRepository) CountPermissionCode(ctx context.Context, code string, excludeID uint64) (int64, error) {
	var count int64
	for _, permission := range repo.permissions {
		if permission.Code == code && permission.ID != excludeID {
			count++
		}
	}
	return count, nil
}

func (repo *fakeMenuRepository) CreateMenuWithPermission(ctx context.Context, menu makeadmin.Menu, permission *makeadmin.Permission) (uint64, error) {
	repo.createdMenu = menu
	repo.createdPermission = permission
	return 9, nil
}

func (repo *fakeMenuRepository) UpdateMenuWithPermission(ctx context.Context, menu makeadmin.Menu, permission *makeadmin.Permission) error {
	repo.updatedMenu = menu
	repo.updatedPermission = permission
	return nil
}

func (repo *fakeMenuRepository) DeleteMenu(ctx context.Context, menuID uint64) error {
	repo.deletedMenuID = menuID
	return nil
}

func TestMenuAddCreatesPageAndPermission(t *testing.T) {
	repo := &fakeMenuRepository{
		menus: []makeadmin.Menu{{ID: 100, MenuType: makeadmin.MenuTypeCatalog}},
	}
	srv := NewMenuService(repo)

	err := srv.Add(context.Background(), MenuInput{
		ParentID:  100,
		MenuType:  "C",
		MenuName:  "烟测页面",
		MenuIcon:  "el-icon-Document",
		MenuSort:  10,
		Perms:     "system:smoke:list",
		Paths:     "smoke",
		Component: "permission/menu/index",
		Selected:  "permission/menu",
		Params:    "id=1",
		IsCache:   1,
		IsShow:    1,
		IsDisable: 0,
	})
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}
	if repo.createdMenu.MenuType != makeadmin.MenuTypePage ||
		repo.createdMenu.RoutePath != "/smoke" ||
		repo.createdMenu.Meta != `{"params":"id=1"}` ||
		repo.createdPermission.Code != "system:smoke:list" ||
		repo.createdPermission.Module != "system" ||
		repo.createdPermission.Resource != "smoke" ||
		repo.createdPermission.Action != "list" {
		t.Fatalf("Add() menu=%#v permission=%#v", repo.createdMenu, repo.createdPermission)
	}
}

func TestMenuAddRejectsDuplicatePermission(t *testing.T) {
	repo := &fakeMenuRepository{
		permissions: []makeadmin.Permission{{ID: 1, Code: "system:menu:list"}},
	}
	srv := NewMenuService(repo)

	err := srv.Add(context.Background(), MenuInput{
		MenuType: "A",
		MenuName: "菜单列表",
		Perms:    "system:menu:list",
	})
	if !errors.Is(err, ErrMenuPermsExists) {
		t.Fatalf("Add() error = %v, want ErrMenuPermsExists", err)
	}
}

func TestMenuEditAllowsCurrentPermission(t *testing.T) {
	repo := &fakeMenuRepository{
		menus: []makeadmin.Menu{{ID: 120, MenuType: makeadmin.MenuTypePage}},
		permissions: []makeadmin.Permission{
			{ID: 17, Code: "system:menu:list"},
		},
		menuPermissions: map[uint64][]string{120: []string{"system:menu:list"}},
	}
	srv := NewMenuService(repo)

	err := srv.Edit(context.Background(), MenuInput{
		ID:       120,
		MenuType: "C",
		MenuName: "菜单管理",
		Perms:    "system:menu:list",
		Paths:    "permission/menu",
		IsShow:   1,
	})
	if err != nil {
		t.Fatalf("Edit() error = %v", err)
	}
	if repo.updatedPermission == nil || repo.updatedPermission.ID != 17 {
		t.Fatalf("Edit() permission = %#v", repo.updatedPermission)
	}
}

func TestMenuDeleteRejectsChildren(t *testing.T) {
	repo := &fakeMenuRepository{
		menus:      []makeadmin.Menu{{ID: 100}},
		childCount: 1,
	}
	srv := NewMenuService(repo)

	err := srv.Delete(context.Background(), 100)
	if !errors.Is(err, ErrMenuHasChildren) {
		t.Fatalf("Delete() error = %v, want ErrMenuHasChildren", err)
	}
}

func TestMenuDetailMapsLegacyFields(t *testing.T) {
	repo := &fakeMenuRepository{
		menus: []makeadmin.Menu{{
			ID:         120,
			ParentID:   100,
			MenuType:   makeadmin.MenuTypePage,
			Name:       "菜单管理",
			RoutePath:  "/permission/menu",
			ActivePath: "/permission",
			Meta:       `{"params":"id=1"}`,
			IsVisible:  1,
			IsCache:    1,
			Status:     makeadmin.StatusDisabled,
			Sort:       10,
		}},
		menuPermissions: map[uint64][]string{120: []string{"system:menu:list"}},
	}
	srv := NewMenuService(repo)

	item, err := srv.Detail(context.Background(), 120)
	if err != nil {
		t.Fatalf("Detail() error = %v", err)
	}
	if item.MenuType != "C" || item.Paths != "permission/menu" || item.Selected != "permission" || item.Params != "id=1" || item.IsDisable != 1 {
		t.Fatalf("Detail() = %#v", item)
	}
}
