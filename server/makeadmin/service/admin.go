package service

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"

	"go-makeadmin/makeadmin/repository"
	"go-makeadmin/makeadmin/security"
	"go-makeadmin/model/makeadmin"
)

var (
	ErrAdminNotFound         = errors.New("makeadmin admin not found")
	ErrAdminUsernameExists   = errors.New("makeadmin admin username exists")
	ErrAdminNicknameExists   = errors.New("makeadmin admin nickname exists")
	ErrAdminRoleNotFound     = errors.New("makeadmin admin role not found")
	ErrAdminRoleDisabled     = errors.New("makeadmin admin role disabled")
	ErrAdminOrgNotFound      = errors.New("makeadmin admin org not found")
	ErrAdminOrgDisabled      = errors.New("makeadmin admin org disabled")
	ErrAdminPositionNotFound = errors.New("makeadmin admin position not found")
	ErrAdminPositionDisabled = errors.New("makeadmin admin position disabled")
	ErrAdminPasswordInvalid  = errors.New("makeadmin admin password invalid")
	ErrSystemAdminProtected  = errors.New("makeadmin system admin protected")
	ErrAdminSelfProtected    = errors.New("makeadmin admin self protected")
)

type AdminInput struct {
	ID           uint64
	TenantID     uint64
	OrgID        uint64
	PositionID   uint64
	Username     string
	Nickname     string
	Password     string
	Avatar       string
	RoleID       uint64
	IsDisable    uint8
	IsMultipoint uint8
}

type AdminSelfInput struct {
	ID           uint64
	Nickname     string
	Avatar       string
	Password     string
	CurrPassword string
}

type AdminItem struct {
	ID            uint64
	Username      string
	Nickname      string
	Avatar        string
	RoleID        uint64
	RoleLabel     string
	OrgID         uint64
	PositionID    uint64
	OrgName       string
	IsMultipoint  uint8
	IsDisable     uint8
	LastLoginIP   string
	LastLoginTime int64
	CreateTime    int64
	UpdateTime    int64
	IsSuper       bool
}

type AdminPage struct {
	Items []AdminItem
	Count int64
}

type AdminService interface {
	List(ctx context.Context, tenantID uint64, filter repository.AdminFilter, pageNo int, pageSize int) (AdminPage, error)
	Detail(ctx context.Context, tenantID uint64, id uint64) (AdminItem, error)
	Add(ctx context.Context, input AdminInput) error
	Edit(ctx context.Context, input AdminInput) error
	UpdateSelf(ctx context.Context, input AdminSelfInput) error
	Delete(ctx context.Context, tenantID uint64, currentAdminID uint64, id uint64) error
	Disable(ctx context.Context, currentAdminID uint64, id uint64) error
}

type adminService struct {
	repo   repository.AdminRepository
	hasher security.PasswordHasher
}

func NewAdminService(repo repository.AdminRepository) AdminService {
	return NewAdminServiceWithPasswordHasher(repo, security.NewBcryptPasswordHasher(0))
}

func NewAdminServiceWithPasswordHasher(repo repository.AdminRepository, hasher security.PasswordHasher) AdminService {
	if hasher == nil {
		hasher = security.NewBcryptPasswordHasher(0)
	}
	return adminService{repo: repo, hasher: hasher}
}

func (srv adminService) List(ctx context.Context, tenantID uint64, filter repository.AdminFilter, pageNo int, pageSize int) (AdminPage, error) {
	admins, count, err := srv.repo.ListAdmins(ctx, tenantID, filter, adminPageLimit(pageSize), adminPageOffset(pageNo, pageSize))
	if err != nil {
		return AdminPage{}, err
	}
	items := make([]AdminItem, 0, len(admins))
	for _, admin := range admins {
		item, err := srv.adminItem(ctx, tenantID, admin)
		if err != nil {
			return AdminPage{}, err
		}
		if admin.IsSuper == 1 {
			item.RoleLabel = "系统管理员"
		}
		items = append(items, item)
	}
	return AdminPage{Items: items, Count: count}, nil
}

func (srv adminService) Detail(ctx context.Context, tenantID uint64, id uint64) (AdminItem, error) {
	admin, err := srv.repo.FindAdminByID(ctx, id)
	if err != nil {
		return AdminItem{}, mapAdminRecordError(err, ErrAdminNotFound)
	}
	item, err := srv.adminItem(ctx, tenantID, admin)
	if err != nil {
		return AdminItem{}, err
	}
	if admin.IsSuper == 1 {
		item.RoleID = 0
		item.RoleLabel = "0"
	}
	return item, nil
}

func (srv adminService) Add(ctx context.Context, input AdminInput) error {
	username := strings.TrimSpace(input.Username)
	nickname := strings.TrimSpace(input.Nickname)
	if err := srv.ensureAdminUnique(ctx, username, nickname, 0); err != nil {
		return err
	}
	if err := srv.validateAdminRelations(ctx, input.TenantID, input.RoleID, input.OrgID, input.PositionID); err != nil {
		return err
	}
	digest, err := srv.hasher.Hash(strings.TrimSpace(input.Password))
	if err != nil {
		return mapAdminPasswordError(err)
	}
	return srv.createAdmin(ctx, input, username, nickname, digest)
}

func (srv adminService) Edit(ctx context.Context, input AdminInput) error {
	current, err := srv.repo.FindAdminByID(ctx, input.ID)
	if err != nil {
		return mapAdminRecordError(err, ErrAdminNotFound)
	}
	username := strings.TrimSpace(input.Username)
	if current.IsSuper == 1 {
		username = current.Username
		input.RoleID = 0
		input.IsDisable = 0
	}
	nickname := strings.TrimSpace(input.Nickname)
	if err := srv.ensureAdminUnique(ctx, username, nickname, input.ID); err != nil {
		return err
	}
	if current.IsSuper == 0 {
		if err := srv.validateAdminRelations(ctx, input.TenantID, input.RoleID, input.OrgID, input.PositionID); err != nil {
			return err
		}
	} else if err := srv.validateOrgPosition(ctx, input.TenantID, input.OrgID, input.PositionID); err != nil {
		return err
	}
	digest := security.PasswordDigest{Hash: current.PasswordHash, Salt: current.PasswordSalt}
	updatePassword := strings.TrimSpace(input.Password) != ""
	if updatePassword {
		var err error
		digest, err = srv.hasher.Hash(strings.TrimSpace(input.Password))
		if err != nil {
			return mapAdminPasswordError(err)
		}
	}
	return srv.repo.UpdateAdminWithRelations(ctx, makeadmin.Admin{
		ID:           input.ID,
		Username:     username,
		PasswordHash: digest.Hash,
		PasswordSalt: digest.Salt,
		IsSuper:      current.IsSuper,
		Status:       statusFromDisable(input.IsDisable),
	}, makeadmin.AdminProfile{
		AdminID:  input.ID,
		Nickname: nickname,
		Avatar:   input.Avatar,
	}, input.RoleID, makeadmin.AdminOrg{
		TenantID:   input.TenantID,
		AdminID:    input.ID,
		OrgID:      input.OrgID,
		PositionID: input.PositionID,
	}, updatePassword)
}

func (srv adminService) UpdateSelf(ctx context.Context, input AdminSelfInput) error {
	current, err := srv.repo.FindAdminByID(ctx, input.ID)
	if err != nil {
		return mapAdminRecordError(err, ErrAdminNotFound)
	}
	matched, err := srv.hasher.Verify(strings.TrimSpace(input.CurrPassword), security.PasswordDigest{
		Hash: current.PasswordHash,
		Salt: current.PasswordSalt,
	})
	if err != nil {
		return mapAdminPasswordError(err)
	}
	if !matched {
		return ErrAdminPasswordInvalid
	}
	digest := security.PasswordDigest{Hash: current.PasswordHash, Salt: current.PasswordSalt}
	updatePassword := strings.TrimSpace(input.Password) != ""
	if updatePassword {
		digest, err = srv.hasher.Hash(strings.TrimSpace(input.Password))
		if err != nil {
			return mapAdminPasswordError(err)
		}
	}
	return srv.repo.UpdateAdminSelf(ctx, makeadmin.Admin{
		ID:           input.ID,
		PasswordHash: digest.Hash,
		PasswordSalt: digest.Salt,
	}, makeadmin.AdminProfile{
		AdminID:  input.ID,
		Nickname: strings.TrimSpace(input.Nickname),
		Avatar:   input.Avatar,
	}, updatePassword)
}

func (srv adminService) Delete(ctx context.Context, tenantID uint64, currentAdminID uint64, id uint64) error {
	admin, err := srv.repo.FindAdminByID(ctx, id)
	if err != nil {
		return mapAdminRecordError(err, ErrAdminNotFound)
	}
	if admin.IsSuper == 1 || id == 1 {
		return ErrSystemAdminProtected
	}
	if id == currentAdminID {
		return ErrAdminSelfProtected
	}
	return srv.repo.SoftDeleteAdmin(ctx, tenantID, id)
}

func (srv adminService) Disable(ctx context.Context, currentAdminID uint64, id uint64) error {
	admin, err := srv.repo.FindAdminByID(ctx, id)
	if err != nil {
		return mapAdminRecordError(err, ErrAdminNotFound)
	}
	if id == currentAdminID {
		return ErrAdminSelfProtected
	}
	if admin.IsSuper == 1 || id == 1 {
		return ErrSystemAdminProtected
	}
	return srv.repo.ToggleAdminStatus(ctx, admin)
}

func (srv adminService) createAdmin(ctx context.Context, input AdminInput, username string, nickname string, digest security.PasswordDigest) error {
	avatar := input.Avatar
	if avatar == "" {
		avatar = "/api/static/backend_avatar.png"
	}
	_, err := srv.repo.CreateAdminWithRelations(ctx, makeadmin.Admin{
		Username:     username,
		PasswordHash: digest.Hash,
		PasswordSalt: digest.Salt,
		Status:       statusFromDisable(input.IsDisable),
	}, makeadmin.AdminProfile{
		Nickname: nickname,
		Avatar:   avatar,
	}, input.RoleID, makeadmin.AdminOrg{
		TenantID:   input.TenantID,
		OrgID:      input.OrgID,
		PositionID: input.PositionID,
		IsPrimary:  1,
		Status:     makeadmin.StatusEnabled,
	})
	return err
}

func (srv adminService) adminItem(ctx context.Context, tenantID uint64, admin makeadmin.Admin) (AdminItem, error) {
	profile, err := srv.repo.FindAdminProfileByAdminID(ctx, admin.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return AdminItem{}, err
	}
	roleIDs, err := srv.repo.ListRoleIDsByAdminID(ctx, tenantID, admin.ID)
	if err != nil {
		return AdminItem{}, err
	}
	roles, err := srv.repo.ListRolesByIDs(ctx, tenantID, roleIDs)
	if err != nil {
		return AdminItem{}, err
	}
	adminOrg, err := srv.repo.FindPrimaryAdminOrg(ctx, tenantID, admin.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return AdminItem{}, err
	}
	orgName := ""
	if adminOrg.OrgID > 0 {
		org, err := srv.repo.FindOrgUnitByID(ctx, tenantID, adminOrg.OrgID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return AdminItem{}, err
		}
		orgName = org.Name
	}
	roleID := uint64(0)
	if len(roleIDs) > 0 {
		roleID = roleIDs[0]
	}
	return AdminItem{
		ID:            admin.ID,
		Username:      admin.Username,
		Nickname:      profile.Nickname,
		Avatar:        profile.Avatar,
		RoleID:        roleID,
		RoleLabel:     roleNames(roles),
		OrgID:         adminOrg.OrgID,
		PositionID:    adminOrg.PositionID,
		OrgName:       orgName,
		IsMultipoint:  1,
		IsDisable:     disableFromStatus(admin.Status),
		LastLoginIP:   admin.LastLoginIP,
		LastLoginTime: admin.LastLoginTime,
		CreateTime:    admin.CreateTime,
		UpdateTime:    admin.UpdateTime,
		IsSuper:       admin.IsSuper == 1,
	}, nil
}

func (srv adminService) ensureAdminUnique(ctx context.Context, username string, nickname string, excludeID uint64) error {
	count, err := srv.repo.CountAdminsByUsername(ctx, username, excludeID)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrAdminUsernameExists
	}
	count, err = srv.repo.CountProfilesByNickname(ctx, nickname, excludeID)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrAdminNicknameExists
	}
	return nil
}

func (srv adminService) validateAdminRelations(ctx context.Context, tenantID uint64, roleID uint64, orgID uint64, positionID uint64) error {
	role, err := srv.repo.FindRoleByID(ctx, tenantID, roleID)
	if err != nil {
		return mapAdminRecordError(err, ErrAdminRoleNotFound)
	}
	if role.Status == makeadmin.StatusDisabled {
		return ErrAdminRoleDisabled
	}
	return srv.validateOrgPosition(ctx, tenantID, orgID, positionID)
}

func (srv adminService) validateOrgPosition(ctx context.Context, tenantID uint64, orgID uint64, positionID uint64) error {
	org, err := srv.repo.FindOrgUnitByID(ctx, tenantID, orgID)
	if err != nil {
		return mapAdminRecordError(err, ErrAdminOrgNotFound)
	}
	if org.Status == makeadmin.StatusDisabled {
		return ErrAdminOrgDisabled
	}
	position, err := srv.repo.FindPositionByID(ctx, tenantID, positionID)
	if err != nil {
		return mapAdminRecordError(err, ErrAdminPositionNotFound)
	}
	if position.Status == makeadmin.StatusDisabled {
		return ErrAdminPositionDisabled
	}
	return nil
}

func roleNames(roles []makeadmin.Role) string {
	names := make([]string, 0, len(roles))
	for _, role := range roles {
		names = append(names, role.Name)
	}
	return strings.Join(names, ",")
}

func mapAdminRecordError(err error, notFound error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return notFound
	}
	return err
}

func mapAdminPasswordError(err error) error {
	if errors.Is(err, security.ErrPasswordEmpty) ||
		errors.Is(err, security.ErrPasswordTooShort) ||
		errors.Is(err, security.ErrPasswordTooLong) ||
		errors.Is(err, security.ErrPasswordPlaceholder) ||
		errors.Is(err, security.ErrPasswordUnsupported) {
		return ErrAdminPasswordInvalid
	}
	return err
}

func adminPageLimit(pageSize int) int {
	if pageSize <= 0 {
		return 20
	}
	return pageSize
}

func adminPageOffset(pageNo int, pageSize int) int {
	if pageNo <= 0 {
		pageNo = 1
	}
	return adminPageLimit(pageSize) * (pageNo - 1)
}
