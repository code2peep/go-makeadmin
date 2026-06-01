package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"go-makeadmin/model/makeadmin"
)

type AdminFilter struct {
	Username string
	Nickname string
	RoleID   uint64
	RoleSet  bool
}

type AdminRepository interface {
	ListAdmins(ctx context.Context, tenantID uint64, filter AdminFilter, limit int, offset int) ([]makeadmin.Admin, int64, error)
	FindAdminByID(ctx context.Context, adminID uint64) (makeadmin.Admin, error)
	FindAdminByUsername(ctx context.Context, username string) (makeadmin.Admin, error)
	FindAdminProfileByAdminID(ctx context.Context, adminID uint64) (makeadmin.AdminProfile, error)
	CountAdminsByUsername(ctx context.Context, username string, excludeID uint64) (int64, error)
	CountProfilesByNickname(ctx context.Context, nickname string, excludeAdminID uint64) (int64, error)
	ListRoleIDsByAdminID(ctx context.Context, tenantID uint64, adminID uint64) ([]uint64, error)
	ListRolesByIDs(ctx context.Context, tenantID uint64, roleIDs []uint64) ([]makeadmin.Role, error)
	FindRoleByID(ctx context.Context, tenantID uint64, roleID uint64) (makeadmin.Role, error)
	FindPrimaryAdminOrg(ctx context.Context, tenantID uint64, adminID uint64) (makeadmin.AdminOrg, error)
	FindOrgUnitByID(ctx context.Context, tenantID uint64, orgID uint64) (makeadmin.OrgUnit, error)
	FindPositionByID(ctx context.Context, tenantID uint64, positionID uint64) (makeadmin.Position, error)
	CreateAdminWithRelations(ctx context.Context, admin makeadmin.Admin, profile makeadmin.AdminProfile, roleID uint64, adminOrg makeadmin.AdminOrg) (uint64, error)
	UpdateAdminWithRelations(ctx context.Context, admin makeadmin.Admin, profile makeadmin.AdminProfile, roleID uint64, adminOrg makeadmin.AdminOrg, updatePassword bool) error
	UpdateAdminSelf(ctx context.Context, admin makeadmin.Admin, profile makeadmin.AdminProfile, updatePassword bool) error
	SoftDeleteAdmin(ctx context.Context, tenantID uint64, adminID uint64) error
	ToggleAdminStatus(ctx context.Context, admin makeadmin.Admin) error
}

type adminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return adminRepository{db: db}
}

func (repo adminRepository) ListAdmins(ctx context.Context, tenantID uint64, filter AdminFilter, limit int, offset int) ([]makeadmin.Admin, int64, error) {
	query := repo.db.WithContext(ctx).
		Model(&makeadmin.Admin{}).
		Where("ma_admin.delete_time = ?", 0)
	if filter.Username != "" {
		query = query.Where("ma_admin.username LIKE ?", "%"+filter.Username+"%")
	}
	if filter.Nickname != "" {
		query = query.
			Joins("INNER JOIN ma_admin_profile AS profile ON profile.admin_id = ma_admin.id").
			Where("profile.nickname LIKE ?", "%"+filter.Nickname+"%")
	}
	if filter.RoleSet {
		query = query.Where("ma_admin.id IN (?)",
			repo.db.Model(&makeadmin.AdminRole{}).
				Select("admin_id").
				Where("tenant_id = ? AND role_id = ?", tenantID, filter.RoleID))
	}
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	var admins []makeadmin.Admin
	err := query.
		Limit(limit).
		Offset(offset).
		Order("ma_admin.id DESC").
		Find(&admins).
		Error
	return admins, count, err
}

func (repo adminRepository) FindAdminByID(ctx context.Context, adminID uint64) (admin makeadmin.Admin, err error) {
	err = repo.db.WithContext(ctx).
		Where("id = ? AND delete_time = ?", adminID, 0).
		Limit(1).
		First(&admin).
		Error
	return
}

func (repo adminRepository) FindAdminByUsername(ctx context.Context, username string) (admin makeadmin.Admin, err error) {
	err = repo.db.WithContext(ctx).
		Where("username = ? AND delete_time = ?", username, 0).
		Limit(1).
		First(&admin).
		Error
	return
}

func (repo adminRepository) FindAdminProfileByAdminID(ctx context.Context, adminID uint64) (profile makeadmin.AdminProfile, err error) {
	err = repo.db.WithContext(ctx).
		Where("admin_id = ?", adminID).
		Limit(1).
		First(&profile).
		Error
	return
}

func (repo adminRepository) CountAdminsByUsername(ctx context.Context, username string, excludeID uint64) (int64, error) {
	query := repo.db.WithContext(ctx).
		Model(&makeadmin.Admin{}).
		Where("username = ? AND delete_time = ?", username, 0)
	if excludeID > 0 {
		query = query.Where("id <> ?", excludeID)
	}
	var count int64
	err := query.Count(&count).Error
	return count, err
}

func (repo adminRepository) CountProfilesByNickname(ctx context.Context, nickname string, excludeAdminID uint64) (int64, error) {
	query := repo.db.WithContext(ctx).
		Model(&makeadmin.AdminProfile{}).
		Where("nickname = ?", nickname)
	if excludeAdminID > 0 {
		query = query.Where("admin_id <> ?", excludeAdminID)
	}
	var count int64
	err := query.Count(&count).Error
	return count, err
}

func (repo adminRepository) ListRoleIDsByAdminID(ctx context.Context, tenantID uint64, adminID uint64) (roleIDs []uint64, err error) {
	err = repo.db.WithContext(ctx).
		Model(&makeadmin.AdminRole{}).
		Where("tenant_id = ? AND admin_id = ?", tenantID, adminID).
		Order("role_id ASC").
		Pluck("role_id", &roleIDs).
		Error
	return
}

func (repo adminRepository) ListRolesByIDs(ctx context.Context, tenantID uint64, roleIDs []uint64) ([]makeadmin.Role, error) {
	if len(roleIDs) == 0 {
		return []makeadmin.Role{}, nil
	}
	var roles []makeadmin.Role
	err := repo.db.WithContext(ctx).
		Where("tenant_id = ? AND id IN ? AND delete_time = ?", tenantID, roleIDs, 0).
		Order("sort DESC, id ASC").
		Find(&roles).
		Error
	return roles, err
}

func (repo adminRepository) FindRoleByID(ctx context.Context, tenantID uint64, roleID uint64) (role makeadmin.Role, err error) {
	err = repo.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ? AND delete_time = ?", tenantID, roleID, 0).
		Limit(1).
		First(&role).
		Error
	return
}

func (repo adminRepository) FindPrimaryAdminOrg(ctx context.Context, tenantID uint64, adminID uint64) (adminOrg makeadmin.AdminOrg, err error) {
	err = repo.db.WithContext(ctx).
		Where("tenant_id = ? AND admin_id = ? AND is_primary = ? AND delete_time = ?", tenantID, adminID, 1, 0).
		Limit(1).
		First(&adminOrg).
		Error
	return
}

func (repo adminRepository) FindOrgUnitByID(ctx context.Context, tenantID uint64, orgID uint64) (org makeadmin.OrgUnit, err error) {
	err = repo.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ? AND delete_time = ?", tenantID, orgID, 0).
		Limit(1).
		First(&org).
		Error
	return
}

func (repo adminRepository) FindPositionByID(ctx context.Context, tenantID uint64, positionID uint64) (position makeadmin.Position, err error) {
	err = repo.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ? AND delete_time = ?", tenantID, positionID, 0).
		Limit(1).
		First(&position).
		Error
	return
}

func (repo adminRepository) CreateAdminWithRelations(ctx context.Context, admin makeadmin.Admin, profile makeadmin.AdminProfile, roleID uint64, adminOrg makeadmin.AdminOrg) (uint64, error) {
	err := repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&admin).Error; err != nil {
			return err
		}
		profile.AdminID = admin.ID
		if err := tx.Create(&profile).Error; err != nil {
			return err
		}
		if roleID > 0 {
			if err := tx.Create(&makeadmin.AdminRole{
				TenantID:   adminOrg.TenantID,
				AdminID:    admin.ID,
				RoleID:     roleID,
				CreateTime: time.Now().Unix(),
			}).Error; err != nil {
				return err
			}
		}
		if adminOrg.OrgID > 0 || adminOrg.PositionID > 0 {
			adminOrg.AdminID = admin.ID
			return tx.Create(&adminOrg).Error
		}
		return nil
	})
	return admin.ID, err
}

func (repo adminRepository) UpdateAdminWithRelations(ctx context.Context, admin makeadmin.Admin, profile makeadmin.AdminProfile, roleID uint64, adminOrg makeadmin.AdminOrg, updatePassword bool) error {
	now := time.Now().Unix()
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		adminUpdates := map[string]interface{}{
			"username":    admin.Username,
			"status":      admin.Status,
			"update_time": now,
		}
		if updatePassword {
			adminUpdates["password_hash"] = admin.PasswordHash
			adminUpdates["password_salt"] = admin.PasswordSalt
		}
		if admin.IsSuper == 1 {
			delete(adminUpdates, "username")
			adminUpdates["status"] = makeadmin.StatusEnabled
		}
		if err := tx.Model(&makeadmin.Admin{}).
			Where("id = ? AND delete_time = ?", admin.ID, 0).
			Updates(adminUpdates).
			Error; err != nil {
			return err
		}
		if err := upsertAdminProfile(tx, profile, now); err != nil {
			return err
		}
		if admin.IsSuper == 0 {
			if err := replaceAdminRole(tx, adminOrg.TenantID, admin.ID, roleID, now); err != nil {
				return err
			}
		}
		return replacePrimaryAdminOrg(tx, adminOrg, now)
	})
}

func (repo adminRepository) UpdateAdminSelf(ctx context.Context, admin makeadmin.Admin, profile makeadmin.AdminProfile, updatePassword bool) error {
	now := time.Now().Unix()
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		adminUpdates := map[string]interface{}{"update_time": now}
		if updatePassword {
			adminUpdates["password_hash"] = admin.PasswordHash
			adminUpdates["password_salt"] = admin.PasswordSalt
		}
		if len(adminUpdates) > 1 {
			if err := tx.Model(&makeadmin.Admin{}).
				Where("id = ? AND delete_time = ?", admin.ID, 0).
				Updates(adminUpdates).
				Error; err != nil {
				return err
			}
		}
		return upsertAdminProfile(tx, profile, now)
	})
}

func (repo adminRepository) SoftDeleteAdmin(ctx context.Context, tenantID uint64, adminID uint64) error {
	now := time.Now().Unix()
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&makeadmin.Admin{}).
			Where("id = ? AND delete_time = ?", adminID, 0).
			Updates(map[string]interface{}{
				"delete_time": now,
				"update_time": now,
			}).Error; err != nil {
			return err
		}
		if err := tx.Where("tenant_id = ? AND admin_id = ?", tenantID, adminID).
			Delete(&makeadmin.AdminRole{}).
			Error; err != nil {
			return err
		}
		return tx.Model(&makeadmin.AdminOrg{}).
			Where("tenant_id = ? AND admin_id = ? AND delete_time = ?", tenantID, adminID, 0).
			Updates(map[string]interface{}{
				"delete_time": now,
				"update_time": now,
			}).Error
	})
}

func (repo adminRepository) ToggleAdminStatus(ctx context.Context, admin makeadmin.Admin) error {
	now := time.Now().Unix()
	status := makeadmin.StatusDisabled
	if admin.Status == makeadmin.StatusDisabled {
		status = makeadmin.StatusEnabled
	}
	return repo.db.WithContext(ctx).
		Model(&makeadmin.Admin{}).
		Where("id = ? AND delete_time = ?", admin.ID, 0).
		Updates(map[string]interface{}{
			"status":      status,
			"update_time": now,
		}).
		Error
}

func upsertAdminProfile(tx *gorm.DB, profile makeadmin.AdminProfile, now int64) error {
	profile.UpdateTime = now
	if profile.CreateTime == 0 {
		profile.CreateTime = now
	}
	return tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "admin_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"nickname":    profile.Nickname,
			"avatar":      profile.Avatar,
			"email":       profile.Email,
			"mobile":      profile.Mobile,
			"remark":      profile.Remark,
			"update_time": now,
		}),
	}).Create(&profile).Error
}

func replaceAdminRole(tx *gorm.DB, tenantID uint64, adminID uint64, roleID uint64, now int64) error {
	if err := tx.Where("tenant_id = ? AND admin_id = ?", tenantID, adminID).
		Delete(&makeadmin.AdminRole{}).
		Error; err != nil {
		return err
	}
	if roleID == 0 {
		return nil
	}
	return tx.Create(&makeadmin.AdminRole{
		TenantID:   tenantID,
		AdminID:    adminID,
		RoleID:     roleID,
		CreateTime: now,
	}).Error
}

func replacePrimaryAdminOrg(tx *gorm.DB, adminOrg makeadmin.AdminOrg, now int64) error {
	if err := tx.Model(&makeadmin.AdminOrg{}).
		Where("tenant_id = ? AND admin_id = ? AND is_primary = ? AND delete_time = ?", adminOrg.TenantID, adminOrg.AdminID, 1, 0).
		Updates(map[string]interface{}{
			"delete_time": now,
			"update_time": now,
		}).Error; err != nil {
		return err
	}
	if adminOrg.OrgID == 0 && adminOrg.PositionID == 0 {
		return nil
	}
	adminOrg.IsPrimary = 1
	adminOrg.Status = makeadmin.StatusEnabled
	adminOrg.CreateTime = now
	adminOrg.UpdateTime = now
	return tx.Create(&adminOrg).Error
}
