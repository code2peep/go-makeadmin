package repository

import (
	"context"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"go-makeadmin/model/makeadmin"
)

type MenuRepository interface {
	ListMenus(ctx context.Context) ([]makeadmin.Menu, error)
	FindMenuByID(ctx context.Context, id uint64) (makeadmin.Menu, error)
	CountChildMenus(ctx context.Context, parentID uint64) (int64, error)
	ListPermissionCodesByMenuID(ctx context.Context, menuID uint64) ([]string, error)
	FindPermissionByCode(ctx context.Context, code string) (makeadmin.Permission, error)
	CountPermissionCode(ctx context.Context, code string, excludeID uint64) (int64, error)
	CreateMenuWithPermission(ctx context.Context, menu makeadmin.Menu, permission *makeadmin.Permission) (uint64, error)
	UpdateMenuWithPermission(ctx context.Context, menu makeadmin.Menu, permission *makeadmin.Permission) error
	DeleteMenu(ctx context.Context, menuID uint64) error
}

type menuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) MenuRepository {
	return menuRepository{db: db}
}

func (repo menuRepository) ListMenus(ctx context.Context) ([]makeadmin.Menu, error) {
	var menus []makeadmin.Menu
	err := repo.db.WithContext(ctx).
		Where("delete_time = ?", 0).
		Order("sort DESC, id ASC").
		Find(&menus).
		Error
	return menus, err
}

func (repo menuRepository) FindMenuByID(ctx context.Context, id uint64) (menu makeadmin.Menu, err error) {
	err = repo.db.WithContext(ctx).
		Where("id = ? AND delete_time = ?", id, 0).
		Limit(1).
		First(&menu).
		Error
	return
}

func (repo menuRepository) CountChildMenus(ctx context.Context, parentID uint64) (int64, error) {
	var count int64
	err := repo.db.WithContext(ctx).
		Model(&makeadmin.Menu{}).
		Where("parent_id = ? AND delete_time = ?", parentID, 0).
		Count(&count).
		Error
	return count, err
}

func (repo menuRepository) ListPermissionCodesByMenuID(ctx context.Context, menuID uint64) ([]string, error) {
	var codes []string
	err := repo.db.WithContext(ctx).
		Table("ma_permission AS p").
		Joins("INNER JOIN ma_menu_permission AS mp ON mp.permission_id = p.id").
		Where("mp.menu_id = ?", menuID).
		Order("p.id ASC").
		Pluck("p.code", &codes).
		Error
	return codes, err
}

func (repo menuRepository) FindPermissionByCode(ctx context.Context, code string) (permission makeadmin.Permission, err error) {
	err = repo.db.WithContext(ctx).
		Where("code = ?", code).
		Limit(1).
		First(&permission).
		Error
	return
}

func (repo menuRepository) CountPermissionCode(ctx context.Context, code string, excludeID uint64) (int64, error) {
	query := repo.db.WithContext(ctx).
		Model(&makeadmin.Permission{}).
		Where("code = ?", code)
	if excludeID > 0 {
		query = query.Where("id <> ?", excludeID)
	}
	var count int64
	err := query.Count(&count).Error
	return count, err
}

func (repo menuRepository) CreateMenuWithPermission(ctx context.Context, menu makeadmin.Menu, permission *makeadmin.Permission) (uint64, error) {
	err := repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&menu).Error; err != nil {
			return err
		}
		permissionID, err := upsertPermission(ctx, tx, permission)
		if err != nil {
			return err
		}
		return replaceMenuPermission(ctx, tx, menu.ID, permissionID)
	})
	return menu.ID, err
}

func (repo menuRepository) UpdateMenuWithPermission(ctx context.Context, menu makeadmin.Menu, permission *makeadmin.Permission) error {
	now := time.Now().Unix()
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&makeadmin.Menu{}).
			Where("id = ? AND delete_time = ?", menu.ID, 0).
			Updates(map[string]interface{}{
				"parent_id":   menu.ParentID,
				"menu_type":   menu.MenuType,
				"name":        menu.Name,
				"icon":        menu.Icon,
				"route_path":  menu.RoutePath,
				"route_name":  menu.RouteName,
				"component":   menu.Component,
				"redirect":    menu.Redirect,
				"active_path": menu.ActivePath,
				"meta":        menu.Meta,
				"is_visible":  menu.IsVisible,
				"is_cache":    menu.IsCache,
				"status":      menu.Status,
				"sort":        menu.Sort,
				"update_time": now,
			}).Error; err != nil {
			return err
		}
		permissionID, err := upsertPermission(ctx, tx, permission)
		if err != nil {
			return err
		}
		return replaceMenuPermission(ctx, tx, menu.ID, permissionID)
	})
}

func (repo menuRepository) DeleteMenu(ctx context.Context, menuID uint64) error {
	now := time.Now().Unix()
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		permissionIDs, err := listPermissionIDsByMenuID(ctx, tx, menuID)
		if err != nil {
			return err
		}
		if err := tx.Model(&makeadmin.Menu{}).
			Where("id = ? AND delete_time = ?", menuID, 0).
			Updates(map[string]interface{}{
				"delete_time": now,
				"update_time": now,
			}).Error; err != nil {
			return err
		}
		if err := tx.Where("menu_id = ?", menuID).Delete(&makeadmin.MenuPermission{}).Error; err != nil {
			return err
		}
		return pruneOrphanRolePermissions(ctx, tx, permissionIDs)
	})
}

func upsertPermission(ctx context.Context, tx *gorm.DB, permission *makeadmin.Permission) (uint64, error) {
	if permission == nil || strings.TrimSpace(permission.Code) == "" {
		return 0, nil
	}
	now := time.Now().Unix()
	permission.Code = strings.TrimSpace(permission.Code)
	permission.CreateTime = now
	permission.UpdateTime = now
	if err := tx.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "code"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"name":        permission.Name,
				"module":      permission.Module,
				"resource":    permission.Resource,
				"action":      permission.Action,
				"status":      permission.Status,
				"sort":        permission.Sort,
				"update_time": now,
			}),
		}).
		Create(permission).
		Error; err != nil {
		return 0, err
	}
	var saved makeadmin.Permission
	if err := tx.WithContext(ctx).
		Where("code = ?", permission.Code).
		Limit(1).
		First(&saved).
		Error; err != nil {
		return 0, err
	}
	return saved.ID, nil
}

func replaceMenuPermission(ctx context.Context, tx *gorm.DB, menuID uint64, permissionID uint64) error {
	oldPermissionIDs, err := listPermissionIDsByMenuID(ctx, tx, menuID)
	if err != nil {
		return err
	}
	if err := tx.WithContext(ctx).
		Where("menu_id = ?", menuID).
		Delete(&makeadmin.MenuPermission{}).
		Error; err != nil {
		return err
	}
	if permissionID > 0 {
		if err := tx.WithContext(ctx).
			Create(&makeadmin.MenuPermission{
				MenuID:       menuID,
				PermissionID: permissionID,
				CreateTime:   time.Now().Unix(),
			}).Error; err != nil {
			return err
		}
		if err := migrateRolePermissions(ctx, tx, oldPermissionIDs, permissionID); err != nil {
			return err
		}
	}
	return pruneOrphanRolePermissions(ctx, tx, oldPermissionIDs)
}

func listPermissionIDsByMenuID(ctx context.Context, tx *gorm.DB, menuID uint64) ([]uint64, error) {
	var permissionIDs []uint64
	err := tx.WithContext(ctx).
		Model(&makeadmin.MenuPermission{}).
		Where("menu_id = ?", menuID).
		Pluck("permission_id", &permissionIDs).
		Error
	return permissionIDs, err
}

func migrateRolePermissions(ctx context.Context, tx *gorm.DB, oldPermissionIDs []uint64, newPermissionID uint64) error {
	if len(oldPermissionIDs) == 0 || newPermissionID == 0 {
		return nil
	}
	type roleRow struct {
		TenantID uint64
		RoleID   uint64
	}
	var rows []roleRow
	if err := tx.WithContext(ctx).
		Model(&makeadmin.RolePermission{}).
		Where("permission_id IN ?", oldPermissionIDs).
		Distinct().
		Find(&rows).
		Error; err != nil {
		return err
	}
	if len(rows) == 0 {
		return nil
	}
	now := time.Now().Unix()
	for _, row := range rows {
		if err := tx.WithContext(ctx).
			Clauses(clause.OnConflict{DoNothing: true}).
			Create(&makeadmin.RolePermission{
				TenantID:     row.TenantID,
				RoleID:       row.RoleID,
				PermissionID: newPermissionID,
				CreateTime:   now,
			}).Error; err != nil {
			return err
		}
	}
	return nil
}

func pruneOrphanRolePermissions(ctx context.Context, tx *gorm.DB, permissionIDs []uint64) error {
	for _, permissionID := range permissionIDs {
		var count int64
		if err := tx.WithContext(ctx).
			Model(&makeadmin.MenuPermission{}).
			Where("permission_id = ?", permissionID).
			Count(&count).
			Error; err != nil {
			return err
		}
		if count > 0 {
			continue
		}
		if err := tx.WithContext(ctx).
			Where("permission_id = ?", permissionID).
			Delete(&makeadmin.RolePermission{}).
			Error; err != nil {
			return err
		}
	}
	return nil
}
