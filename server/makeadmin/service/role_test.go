package service

import (
	"context"
	"errors"
	"testing"

	"gorm.io/gorm"

	"go-makeadmin/model/makeadmin"
)

type fakeRoleRepository struct {
	roles         []makeadmin.Role
	menuIDs       []uint64
	memberCount   int64
	createdRole   makeadmin.Role
	createdMenus  []uint64
	updatedRole   makeadmin.Role
	updatedMenus  []uint64
	deletedRoleID uint64
}

func (repo *fakeRoleRepository) ListAllRoles(ctx context.Context, tenantID uint64) ([]makeadmin.Role, error) {
	return repo.roles, nil
}

func (repo *fakeRoleRepository) ListRoles(ctx context.Context, tenantID uint64, limit int, offset int) ([]makeadmin.Role, int64, error) {
	return repo.roles, int64(len(repo.roles)), nil
}

func (repo *fakeRoleRepository) FindRoleByID(ctx context.Context, tenantID uint64, id uint64) (makeadmin.Role, error) {
	for _, role := range repo.roles {
		if role.ID == id && role.TenantID == tenantID && role.DeleteTime == 0 {
			return role, nil
		}
	}
	return makeadmin.Role{}, gorm.ErrRecordNotFound
}

func (repo *fakeRoleRepository) CountRolesByName(ctx context.Context, tenantID uint64, name string, excludeID uint64) (int64, error) {
	var count int64
	for _, role := range repo.roles {
		if role.TenantID == tenantID && role.Name == name && role.ID != excludeID && role.DeleteTime == 0 {
			count++
		}
	}
	return count, nil
}

func (repo *fakeRoleRepository) CountAdminsByRoleID(ctx context.Context, tenantID uint64, roleID uint64) (int64, error) {
	return repo.memberCount, nil
}

func (repo *fakeRoleRepository) ListMenuIDsByRoleID(ctx context.Context, tenantID uint64, roleID uint64) ([]uint64, error) {
	return repo.menuIDs, nil
}

func (repo *fakeRoleRepository) CreateRoleWithMenuIDs(ctx context.Context, role makeadmin.Role, menuIDs []uint64) (makeadmin.Role, error) {
	role.ID = 9
	repo.createdRole = role
	repo.createdMenus = menuIDs
	return role, nil
}

func (repo *fakeRoleRepository) UpdateRoleWithMenuIDs(ctx context.Context, role makeadmin.Role, menuIDs []uint64) error {
	repo.updatedRole = role
	repo.updatedMenus = menuIDs
	return nil
}

func (repo *fakeRoleRepository) DeleteRole(ctx context.Context, tenantID uint64, roleID uint64) error {
	repo.deletedRoleID = roleID
	return nil
}

func TestRoleAddCreatesRoleAndPermissions(t *testing.T) {
	repo := &fakeRoleRepository{}
	srv := NewRoleService(repo)

	err := srv.Add(context.Background(), RoleInput{
		TenantID:  makeadmin.GlobalTenantID,
		Name:      " 运营 ",
		Remark:    "ops",
		Sort:      12,
		IsDisable: 0,
		MenuIDs:   []uint64{1, 2},
	})
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}
	if repo.createdRole.Name != "运营" || repo.createdRole.Status != makeadmin.StatusEnabled || len(repo.createdMenus) != 2 {
		t.Fatalf("Add() role=%#v menus=%#v", repo.createdRole, repo.createdMenus)
	}
}

func TestRoleAddRejectsDuplicateName(t *testing.T) {
	srv := NewRoleService(&fakeRoleRepository{
		roles: []makeadmin.Role{{ID: 1, TenantID: makeadmin.GlobalTenantID, Name: "运营"}},
	})

	err := srv.Add(context.Background(), RoleInput{TenantID: makeadmin.GlobalTenantID, Name: "运营"})
	if !errors.Is(err, ErrRoleNameExists) {
		t.Fatalf("Add() error = %v, want ErrRoleNameExists", err)
	}
}

func TestRoleDetailIncludesMenusAndMembers(t *testing.T) {
	srv := NewRoleService(&fakeRoleRepository{
		roles:       []makeadmin.Role{{ID: 1, TenantID: makeadmin.GlobalTenantID, Name: "运营", Status: makeadmin.StatusDisabled}},
		menuIDs:     []uint64{10, 20},
		memberCount: 3,
	})

	item, err := srv.Detail(context.Background(), makeadmin.GlobalTenantID, 1)
	if err != nil {
		t.Fatalf("Detail() error = %v", err)
	}
	if item.IsDisable != 1 || item.Member != 3 || len(item.MenuIDs) != 2 {
		t.Fatalf("Detail() = %#v", item)
	}
}

func TestRoleDeleteRejectsProtectedAndUsedRole(t *testing.T) {
	repo := &fakeRoleRepository{
		roles: []makeadmin.Role{{ID: 1, TenantID: makeadmin.GlobalTenantID, IsSystem: 1}},
	}
	srv := NewRoleService(repo)

	err := srv.Delete(context.Background(), makeadmin.GlobalTenantID, 1)
	if !errors.Is(err, ErrSystemRoleProtected) {
		t.Fatalf("Delete() protected error = %v", err)
	}

	repo.roles = []makeadmin.Role{{ID: 2, TenantID: makeadmin.GlobalTenantID}}
	repo.memberCount = 1
	err = srv.Delete(context.Background(), makeadmin.GlobalTenantID, 2)
	if !errors.Is(err, ErrRoleInUse) {
		t.Fatalf("Delete() in-use error = %v", err)
	}
}

func TestParseRoleMenuIDs(t *testing.T) {
	ids := ParseRoleMenuIDs("1, 2,2,abc,0,3")
	if len(ids) != 3 || ids[0] != 1 || ids[1] != 2 || ids[2] != 3 {
		t.Fatalf("ParseRoleMenuIDs() = %#v", ids)
	}
}
