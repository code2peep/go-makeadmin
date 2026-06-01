package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"go-makeadmin/model/makeadmin"
)

type RoleRepository interface {
	ListAllRoles(ctx context.Context, tenantID uint64) ([]makeadmin.Role, error)
	ListRoles(ctx context.Context, tenantID uint64, limit int, offset int) ([]makeadmin.Role, int64, error)
	FindRoleByID(ctx context.Context, tenantID uint64, id uint64) (makeadmin.Role, error)
	CountRolesByName(ctx context.Context, tenantID uint64, name string, excludeID uint64) (int64, error)
	CountAdminsByRoleID(ctx context.Context, tenantID uint64, roleID uint64) (int64, error)
	ListMenuIDsByRoleID(ctx context.Context, tenantID uint64, roleID uint64) ([]uint64, error)
	CreateRoleWithMenuIDs(ctx context.Context, role makeadmin.Role, menuIDs []uint64) (makeadmin.Role, error)
	UpdateRoleWithMenuIDs(ctx context.Context, role makeadmin.Role, menuIDs []uint64) error
	DeleteRole(ctx context.Context, tenantID uint64, roleID uint64) error
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return roleRepository{db: db}
}

func (repo roleRepository) ListAllRoles(ctx context.Context, tenantID uint64) ([]makeadmin.Role, error) {
	var roles []makeadmin.Role
	err := repo.db.WithContext(ctx).
		Where("tenant_id = ? AND delete_time = ?", tenantID, 0).
		Order("sort DESC, id DESC").
		Find(&roles).
		Error
	return roles, err
}

func (repo roleRepository) ListRoles(ctx context.Context, tenantID uint64, limit int, offset int) ([]makeadmin.Role, int64, error) {
	query := repo.db.WithContext(ctx).
		Model(&makeadmin.Role{}).
		Where("tenant_id = ? AND delete_time = ?", tenantID, 0)
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	var roles []makeadmin.Role
	err := query.
		Limit(limit).
		Offset(offset).
		Order("sort DESC, id DESC").
		Find(&roles).
		Error
	return roles, count, err
}

func (repo roleRepository) FindRoleByID(ctx context.Context, tenantID uint64, id uint64) (role makeadmin.Role, err error) {
	err = repo.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ? AND delete_time = ?", tenantID, id, 0).
		Limit(1).
		First(&role).
		Error
	return
}

func (repo roleRepository) CountRolesByName(ctx context.Context, tenantID uint64, name string, excludeID uint64) (int64, error) {
	query := repo.db.WithContext(ctx).
		Model(&makeadmin.Role{}).
		Where("tenant_id = ? AND name = ? AND delete_time = ?", tenantID, name, 0)
	if excludeID > 0 {
		query = query.Where("id <> ?", excludeID)
	}
	var count int64
	err := query.Count(&count).Error
	return count, err
}

func (repo roleRepository) CountAdminsByRoleID(ctx context.Context, tenantID uint64, roleID uint64) (int64, error) {
	var count int64
	err := repo.db.WithContext(ctx).
		Model(&makeadmin.AdminRole{}).
		Joins("INNER JOIN ma_admin AS admin ON admin.id = ma_admin_role.admin_id").
		Where("ma_admin_role.tenant_id = ? AND ma_admin_role.role_id = ? AND admin.delete_time = ?", tenantID, roleID, 0).
		Count(&count).
		Error
	return count, err
}

func (repo roleRepository) ListMenuIDsByRoleID(ctx context.Context, tenantID uint64, roleID uint64) ([]uint64, error) {
	var menuIDs []uint64
	err := repo.db.WithContext(ctx).
		Table("ma_menu_permission AS mp").
		Joins("INNER JOIN ma_role_permission AS rp ON rp.permission_id = mp.permission_id").
		Where("rp.tenant_id = ? AND rp.role_id = ?", tenantID, roleID).
		Distinct().
		Order("mp.menu_id ASC").
		Pluck("mp.menu_id", &menuIDs).
		Error
	return menuIDs, err
}

func (repo roleRepository) CreateRoleWithMenuIDs(ctx context.Context, role makeadmin.Role, menuIDs []uint64) (makeadmin.Role, error) {
	err := repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&role).Error; err != nil {
			return err
		}
		return replaceRolePermissionsByMenuIDs(ctx, tx, role.TenantID, role.ID, menuIDs)
	})
	return role, err
}

func (repo roleRepository) UpdateRoleWithMenuIDs(ctx context.Context, role makeadmin.Role, menuIDs []uint64) error {
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&makeadmin.Role{}).
			Where("tenant_id = ? AND id = ? AND delete_time = ?", role.TenantID, role.ID, 0).
			Updates(map[string]interface{}{
				"name":        role.Name,
				"remark":      role.Remark,
				"status":      role.Status,
				"sort":        role.Sort,
				"update_time": time.Now().Unix(),
			}).Error; err != nil {
			return err
		}
		return replaceRolePermissionsByMenuIDs(ctx, tx, role.TenantID, role.ID, menuIDs)
	})
}

func (repo roleRepository) DeleteRole(ctx context.Context, tenantID uint64, roleID uint64) error {
	now := time.Now().Unix()
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&makeadmin.Role{}).
			Where("tenant_id = ? AND id = ? AND delete_time = ?", tenantID, roleID, 0).
			Updates(map[string]interface{}{
				"delete_time": now,
				"update_time": now,
			}).Error; err != nil {
			return err
		}
		return tx.Where("tenant_id = ? AND role_id = ?", tenantID, roleID).Delete(&makeadmin.RolePermission{}).Error
	})
}

func replaceRolePermissionsByMenuIDs(ctx context.Context, tx *gorm.DB, tenantID uint64, roleID uint64, menuIDs []uint64) error {
	if err := tx.WithContext(ctx).
		Where("tenant_id = ? AND role_id = ?", tenantID, roleID).
		Delete(&makeadmin.RolePermission{}).
		Error; err != nil {
		return err
	}
	if len(menuIDs) == 0 {
		return nil
	}
	var permissionIDs []uint64
	if err := tx.WithContext(ctx).
		Model(&makeadmin.MenuPermission{}).
		Where("menu_id IN ?", menuIDs).
		Distinct().
		Pluck("permission_id", &permissionIDs).
		Error; err != nil {
		return err
	}
	if len(permissionIDs) == 0 {
		return nil
	}
	rolePermissions := make([]makeadmin.RolePermission, 0, len(permissionIDs))
	for _, permissionID := range permissionIDs {
		rolePermissions = append(rolePermissions, makeadmin.RolePermission{
			TenantID:     tenantID,
			RoleID:       roleID,
			PermissionID: permissionID,
		})
	}
	return tx.WithContext(ctx).Create(&rolePermissions).Error
}
