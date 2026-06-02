package repository

import (
	"context"

	"gorm.io/gorm"

	"go-makeadmin/model/makeadmin"
)

type AuthRepository interface {
	FindAdminByID(ctx context.Context, adminID uint64) (makeadmin.Admin, error)
	FindAdminByUsername(ctx context.Context, username string) (makeadmin.Admin, error)
	FindAdminProfileByAdminID(ctx context.Context, adminID uint64) (makeadmin.AdminProfile, error)
	FindTenantByID(ctx context.Context, tenantID uint64) (makeadmin.Tenant, error)
	FindTenantMember(ctx context.Context, tenantID uint64, adminID uint64) (makeadmin.TenantMember, error)
	ListTenantMembershipsByAdminID(ctx context.Context, adminID uint64) ([]TenantMembership, error)
	ListRoleIDsByAdminID(ctx context.Context, tenantID uint64, adminID uint64) ([]uint64, error)
	ListPermissionCodesByRoleIDs(ctx context.Context, tenantID uint64, roleIDs []uint64) ([]string, error)
	FindPrimaryAdminOrg(ctx context.Context, tenantID uint64, adminID uint64) (makeadmin.AdminOrg, error)
	ListDataScopesByRoleIDs(ctx context.Context, tenantID uint64, roleIDs []uint64) ([]makeadmin.DataScope, error)
	ListOrgUnits(ctx context.Context, tenantID uint64) ([]makeadmin.OrgUnit, error)
	ListVisibleRouteMenus(ctx context.Context) ([]makeadmin.Menu, error)
	ListMenuPermissionCodes(ctx context.Context) (map[uint64][]string, error)
	UpdateAdminLoginInfo(ctx context.Context, adminID uint64, ip string, loginTime int64) error
	CreateLoginLog(ctx context.Context, loginLog makeadmin.LoginLog) error
}

type TenantMembership struct {
	TenantID   uint64
	Code       string
	Name       string
	MemberType string
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (repo authRepository) FindAdminByID(ctx context.Context, adminID uint64) (admin makeadmin.Admin, err error) {
	err = repo.db.WithContext(ctx).
		Where("id = ? AND delete_time = ?", adminID, 0).
		Limit(1).
		First(&admin).
		Error
	return
}

func (repo authRepository) FindAdminByUsername(ctx context.Context, username string) (admin makeadmin.Admin, err error) {
	err = repo.db.WithContext(ctx).
		Where("username = ? AND delete_time = ?", username, 0).
		Limit(1).
		First(&admin).
		Error
	return
}

func (repo authRepository) FindAdminProfileByAdminID(ctx context.Context, adminID uint64) (profile makeadmin.AdminProfile, err error) {
	err = repo.db.WithContext(ctx).
		Where("admin_id = ?", adminID).
		Limit(1).
		First(&profile).
		Error
	return
}

func (repo authRepository) FindTenantByID(ctx context.Context, tenantID uint64) (tenant makeadmin.Tenant, err error) {
	err = repo.db.WithContext(ctx).
		Where("id = ? AND delete_time = ?", tenantID, 0).
		Limit(1).
		First(&tenant).
		Error
	return
}

func (repo authRepository) FindTenantMember(ctx context.Context, tenantID uint64, adminID uint64) (member makeadmin.TenantMember, err error) {
	err = repo.db.WithContext(ctx).
		Where("tenant_id = ? AND admin_id = ? AND status = ? AND delete_time = ?", tenantID, adminID, makeadmin.StatusEnabled, 0).
		Limit(1).
		First(&member).
		Error
	return
}

func (repo authRepository) ListTenantMembershipsByAdminID(ctx context.Context, adminID uint64) ([]TenantMembership, error) {
	var rows []TenantMembership
	err := repo.db.WithContext(ctx).
		Table("ma_tenant_member AS member").
		Select("tenant.id AS tenant_id, tenant.code, tenant.name, member.member_type").
		Joins("INNER JOIN ma_tenant AS tenant ON tenant.id = member.tenant_id").
		Where("member.admin_id = ? AND member.status = ? AND member.delete_time = ?", adminID, makeadmin.StatusEnabled, 0).
		Where("tenant.status = ? AND tenant.delete_time = ?", makeadmin.StatusEnabled, 0).
		Order("tenant.id ASC").
		Find(&rows).
		Error
	return rows, err
}

func (repo authRepository) ListRoleIDsByAdminID(ctx context.Context, tenantID uint64, adminID uint64) (roleIDs []uint64, err error) {
	err = repo.db.WithContext(ctx).
		Model(&makeadmin.AdminRole{}).
		Where("tenant_id = ? AND admin_id = ?", tenantID, adminID).
		Pluck("role_id", &roleIDs).
		Error
	return
}

func (repo authRepository) ListPermissionCodesByRoleIDs(ctx context.Context, tenantID uint64, roleIDs []uint64) (codes []string, err error) {
	if len(roleIDs) == 0 {
		return []string{}, nil
	}
	err = repo.db.WithContext(ctx).
		Table("ma_permission AS p").
		Joins("INNER JOIN ma_role_permission AS rp ON rp.permission_id = p.id").
		Where("rp.tenant_id = ? AND rp.role_id IN ? AND p.status = ?", tenantID, roleIDs, makeadmin.StatusEnabled).
		Order("p.sort DESC, p.id ASC").
		Distinct().
		Pluck("p.code", &codes).
		Error
	return
}

func (repo authRepository) FindPrimaryAdminOrg(ctx context.Context, tenantID uint64, adminID uint64) (adminOrg makeadmin.AdminOrg, err error) {
	err = repo.db.WithContext(ctx).
		Where("tenant_id = ? AND admin_id = ? AND is_primary = ? AND status = ? AND delete_time = ?", tenantID, adminID, 1, makeadmin.StatusEnabled, 0).
		Limit(1).
		First(&adminOrg).
		Error
	return
}

func (repo authRepository) ListDataScopesByRoleIDs(ctx context.Context, tenantID uint64, roleIDs []uint64) ([]makeadmin.DataScope, error) {
	if len(roleIDs) == 0 {
		return []makeadmin.DataScope{}, nil
	}
	var scopes []makeadmin.DataScope
	err := repo.db.WithContext(ctx).
		Table("ma_data_scope AS scope").
		Select("scope.*").
		Joins("INNER JOIN ma_role_data_scope AS rds ON rds.data_scope_id = scope.id").
		Where("rds.tenant_id = ? AND rds.role_id IN ?", tenantID, roleIDs).
		Where("scope.tenant_id = ? AND scope.status = ? AND scope.delete_time = ?", tenantID, makeadmin.StatusEnabled, 0).
		Order("scope.id ASC").
		Find(&scopes).
		Error
	return scopes, err
}

func (repo authRepository) ListOrgUnits(ctx context.Context, tenantID uint64) (orgs []makeadmin.OrgUnit, err error) {
	err = repo.db.WithContext(ctx).
		Where("tenant_id = ? AND status = ? AND delete_time = ?", tenantID, makeadmin.StatusEnabled, 0).
		Order("id ASC").
		Find(&orgs).
		Error
	return
}

func (repo authRepository) ListVisibleRouteMenus(ctx context.Context) (menus []makeadmin.Menu, err error) {
	err = repo.db.WithContext(ctx).
		Where("menu_type IN ? AND is_visible = ? AND status = ? AND delete_time = ?", []string{
			makeadmin.MenuTypeCatalog,
			makeadmin.MenuTypePage,
		}, 1, makeadmin.StatusEnabled, 0).
		Order("sort DESC, id ASC").
		Find(&menus).
		Error
	return
}

func (repo authRepository) ListMenuPermissionCodes(ctx context.Context) (map[uint64][]string, error) {
	type row struct {
		MenuID uint64
		Code   string
	}
	var rows []row
	err := repo.db.WithContext(ctx).
		Table("ma_menu_permission AS mp").
		Select("mp.menu_id, p.code").
		Joins("INNER JOIN ma_permission AS p ON p.id = mp.permission_id").
		Where("p.status = ?", makeadmin.StatusEnabled).
		Find(&rows).
		Error
	if err != nil {
		return nil, err
	}
	result := make(map[uint64][]string, len(rows))
	for _, item := range rows {
		result[item.MenuID] = append(result[item.MenuID], item.Code)
	}
	return result, nil
}

func (repo authRepository) UpdateAdminLoginInfo(ctx context.Context, adminID uint64, ip string, loginTime int64) error {
	return repo.db.WithContext(ctx).
		Model(&makeadmin.Admin{}).
		Where("id = ?", adminID).
		Updates(map[string]interface{}{
			"last_login_ip":   ip,
			"last_login_time": loginTime,
		}).
		Error
}

func (repo authRepository) CreateLoginLog(ctx context.Context, loginLog makeadmin.LoginLog) error {
	return repo.db.WithContext(ctx).Create(&loginLog).Error
}
