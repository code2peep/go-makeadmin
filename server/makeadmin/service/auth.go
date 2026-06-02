package service

import (
	"context"
	"errors"
	"sort"
	"time"

	"gorm.io/gorm"

	"go-makeadmin/makeadmin/repository"
	"go-makeadmin/makeadmin/security"
	"go-makeadmin/model/makeadmin"
)

var (
	ErrAdminDisabled     = errors.New("makeadmin admin is disabled")
	ErrAdminDeleted      = errors.New("makeadmin admin is deleted")
	ErrInvalidCredential = errors.New("makeadmin invalid credential")
	ErrTenantNotFound    = errors.New("makeadmin tenant not found")
	ErrTenantDisabled    = errors.New("makeadmin tenant is disabled")
	ErrTenantForbidden   = errors.New("makeadmin tenant access forbidden")
)

type Identity struct {
	AdminID     uint64
	Username    string
	Nickname    string
	Avatar      string
	IsSuper     bool
	Status      uint8
	TenantID    uint64
	RoleIDs     []uint64
	Permissions []string
	DataScope   repository.DataScopeFilter
}

type RouteMenu struct {
	ID         uint64
	ParentID   uint64
	MenuType   string
	Name       string
	Perms      string
	Icon       string
	RoutePath  string
	RouteName  string
	Component  string
	Redirect   string
	ActivePath string
	Meta       string
	IsCache    uint8
	Sort       uint16
}

type LoginInput struct {
	TenantID uint64
	Username string
	Password string
	IP       string
	OS       string
	Browser  string
}

type LoginResult struct {
	Token     string
	ExpiresIn int
	Identity  Identity
}

type TenantItem struct {
	ID         uint64
	Code       string
	Name       string
	MemberType string
	IsCurrent  bool
}

type AuthService interface {
	Login(ctx context.Context, input LoginInput) (LoginResult, error)
	Logout(ctx context.Context, token string) error
	AuthenticateByUsername(ctx context.Context, tenantID uint64, username string, plainPassword string) (Identity, error)
	BuildIdentityByAdminID(ctx context.Context, tenantID uint64, adminID uint64) (Identity, error)
	BuildIdentityByUsername(ctx context.Context, tenantID uint64, username string) (Identity, error)
	ListTenants(ctx context.Context, adminID uint64, currentTenantID uint64) ([]TenantItem, error)
	SwitchTenant(ctx context.Context, adminID uint64, tenantID uint64) (LoginResult, error)
	ListRouteMenus(ctx context.Context, identity Identity) ([]RouteMenu, error)
}

type authService struct {
	repo           repository.AuthRepository
	passwordHasher security.PasswordHasher
	tokenCodec     TokenCodec
	sessionStore   SessionStore
	sessionTTL     int
}

func NewAuthService(repo repository.AuthRepository) AuthService {
	return NewAuthServiceWithPasswordHasher(repo, security.NewBcryptPasswordHasher(0))
}

func NewAuthServiceWithPasswordHasher(repo repository.AuthRepository, passwordHasher security.PasswordHasher) AuthService {
	return NewAuthServiceWithDependencies(
		repo,
		passwordHasher,
		NewJWTTokenCodec("go-makeadmin-dev-secret"),
		UnavailableSessionStore{},
		DefaultSessionTTLSeconds,
	)
}

func NewAuthServiceWithDependencies(
	repo repository.AuthRepository,
	passwordHasher security.PasswordHasher,
	tokenCodec TokenCodec,
	sessionStore SessionStore,
	sessionTTL int,
) AuthService {
	if passwordHasher == nil {
		passwordHasher = security.NewBcryptPasswordHasher(0)
	}
	if tokenCodec == nil {
		tokenCodec = NewJWTTokenCodec("go-makeadmin-dev-secret")
	}
	if sessionStore == nil {
		sessionStore = UnavailableSessionStore{}
	}
	if sessionTTL <= 0 {
		sessionTTL = DefaultSessionTTLSeconds
	}
	return &authService{
		repo:           repo,
		passwordHasher: passwordHasher,
		tokenCodec:     tokenCodec,
		sessionStore:   sessionStore,
		sessionTTL:     sessionTTL,
	}
}

func (srv authService) Login(ctx context.Context, input LoginInput) (LoginResult, error) {
	admin, err := srv.repo.FindAdminByUsername(ctx, input.Username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		_ = srv.recordLoginLog(ctx, loginLogInput{
			TenantID: input.TenantID,
			Username: input.Username,
			IP:       input.IP,
			OS:       input.OS,
			Browser:  input.Browser,
			Message:  ErrInvalidCredential.Error(),
		})
		return LoginResult{}, ErrInvalidCredential
	}
	if err != nil {
		return LoginResult{}, err
	}

	matched, err := srv.passwordHasher.Verify(input.Password, security.PasswordDigest{
		Hash: admin.PasswordHash,
		Salt: admin.PasswordSalt,
	})
	if err != nil {
		_ = srv.recordLoginLog(ctx, loginLogInputFromAdmin(input, admin, err.Error()))
		return LoginResult{}, err
	}
	if !matched {
		_ = srv.recordLoginLog(ctx, loginLogInputFromAdmin(input, admin, ErrInvalidCredential.Error()))
		return LoginResult{}, ErrInvalidCredential
	}

	identity, err := srv.buildIdentityByAdmin(ctx, input.TenantID, admin)
	if err != nil {
		_ = srv.recordLoginLog(ctx, loginLogInputFromAdmin(input, admin, err.Error()))
		return LoginResult{}, err
	}

	sessionToken, err := srv.issueSession(ctx, identity)
	if err != nil {
		_ = srv.recordLoginLog(ctx, loginLogInputFromAdmin(input, admin, err.Error()))
		return LoginResult{}, err
	}

	loginTime := time.Now().Unix()
	if err := srv.repo.UpdateAdminLoginInfo(ctx, admin.ID, input.IP, loginTime); err != nil {
		_ = srv.recordLoginLog(ctx, loginLogInputFromAdmin(input, admin, err.Error()))
		return LoginResult{}, err
	}

	if err := srv.recordLoginLog(ctx, loginLogInputFromAdmin(input, admin, "")); err != nil {
		return LoginResult{}, err
	}
	return LoginResult{
		Token:     sessionToken.AccessToken,
		ExpiresIn: srv.sessionTTL,
		Identity:  identity,
	}, nil
}

func (srv authService) issueSession(ctx context.Context, identity Identity) (SessionToken, error) {
	sessionToken, err := srv.tokenCodec.Issue(identity, srv.sessionTTL)
	if err != nil {
		return SessionToken{}, err
	}
	if err := srv.sessionStore.Save(ctx, sessionToken.SessionID, identity, srv.sessionTTL); err != nil {
		return SessionToken{}, err
	}
	return sessionToken, nil
}

func (srv authService) Logout(ctx context.Context, token string) error {
	if token == "" {
		return nil
	}
	claims, err := srv.tokenCodec.Parse(token)
	if err != nil {
		return err
	}
	return srv.sessionStore.Delete(ctx, claims.SessionID)
}

func (srv authService) AuthenticateByUsername(ctx context.Context, tenantID uint64, username string, plainPassword string) (Identity, error) {
	admin, err := srv.repo.FindAdminByUsername(ctx, username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return Identity{}, ErrInvalidCredential
	}
	if err != nil {
		return Identity{}, err
	}

	matched, err := srv.passwordHasher.Verify(plainPassword, security.PasswordDigest{
		Hash: admin.PasswordHash,
		Salt: admin.PasswordSalt,
	})
	if err != nil {
		return Identity{}, err
	}
	if !matched {
		return Identity{}, ErrInvalidCredential
	}

	return srv.buildIdentityByAdmin(ctx, tenantID, admin)
}

func (srv authService) BuildIdentityByUsername(ctx context.Context, tenantID uint64, username string) (Identity, error) {
	admin, err := srv.repo.FindAdminByUsername(ctx, username)
	if err != nil {
		return Identity{}, err
	}
	return srv.buildIdentityByAdmin(ctx, tenantID, admin)
}

func (srv authService) BuildIdentityByAdminID(ctx context.Context, tenantID uint64, adminID uint64) (Identity, error) {
	admin, err := srv.repo.FindAdminByID(ctx, adminID)
	if err != nil {
		return Identity{}, err
	}
	return srv.buildIdentityByAdmin(ctx, tenantID, admin)
}

func (srv authService) buildIdentityByAdmin(ctx context.Context, tenantID uint64, admin makeadmin.Admin) (Identity, error) {
	if admin.DeleteTime != 0 {
		return Identity{}, ErrAdminDeleted
	}
	if admin.Status != makeadmin.StatusEnabled {
		return Identity{}, ErrAdminDisabled
	}
	if err := srv.ensureTenantAccess(ctx, admin.ID, tenantID); err != nil {
		return Identity{}, err
	}

	profile, err := srv.repo.FindAdminProfileByAdminID(ctx, admin.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return Identity{}, err
	}

	identity := Identity{
		AdminID:  admin.ID,
		Username: admin.Username,
		Nickname: profile.Nickname,
		Avatar:   profile.Avatar,
		IsSuper:  admin.IsSuper == 1,
		Status:   admin.Status,
		TenantID: tenantID,
	}
	if identity.Nickname == "" {
		identity.Nickname = admin.Username
	}

	if identity.IsSuper {
		identity.Permissions = []string{"*"}
		identity.DataScope = repository.DataScopeFilter{
			Enabled: true,
			All:     true,
			AdminID: admin.ID,
		}
		return identity, nil
	}

	roleIDs, err := srv.repo.ListRoleIDsByAdminID(ctx, tenantID, admin.ID)
	if err != nil {
		return Identity{}, err
	}
	identity.RoleIDs = roleIDs

	permissions, err := srv.repo.ListPermissionCodesByRoleIDs(ctx, tenantID, roleIDs)
	if err != nil {
		return Identity{}, err
	}
	identity.Permissions = permissions
	dataScope, err := srv.resolveDataScope(ctx, tenantID, admin.ID, roleIDs)
	if err != nil {
		return Identity{}, err
	}
	identity.DataScope = dataScope
	return identity, nil
}

func (srv authService) ListTenants(ctx context.Context, adminID uint64, currentTenantID uint64) ([]TenantItem, error) {
	rows, err := srv.repo.ListTenantMembershipsByAdminID(ctx, adminID)
	if err != nil {
		return nil, err
	}
	items := []TenantItem{{
		ID:         makeadmin.GlobalTenantID,
		Code:       "default",
		Name:       "默认租户",
		MemberType: "default",
		IsCurrent:  currentTenantID == makeadmin.GlobalTenantID,
	}}
	for _, row := range rows {
		if row.TenantID == makeadmin.GlobalTenantID {
			continue
		}
		items = append(items, TenantItem{
			ID:         row.TenantID,
			Code:       row.Code,
			Name:       row.Name,
			MemberType: row.MemberType,
			IsCurrent:  row.TenantID == currentTenantID,
		})
	}
	return items, nil
}

func (srv authService) SwitchTenant(ctx context.Context, adminID uint64, tenantID uint64) (LoginResult, error) {
	identity, err := srv.BuildIdentityByAdminID(ctx, tenantID, adminID)
	if err != nil {
		return LoginResult{}, err
	}
	sessionToken, err := srv.issueSession(ctx, identity)
	if err != nil {
		return LoginResult{}, err
	}
	return LoginResult{
		Token:     sessionToken.AccessToken,
		ExpiresIn: srv.sessionTTL,
		Identity:  identity,
	}, nil
}

func (srv authService) ensureTenantAccess(ctx context.Context, adminID uint64, tenantID uint64) error {
	if tenantID == makeadmin.GlobalTenantID {
		return nil
	}
	tenant, err := srv.repo.FindTenantByID(ctx, tenantID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrTenantNotFound
	}
	if err != nil {
		return err
	}
	if tenant.Status != makeadmin.StatusEnabled {
		return ErrTenantDisabled
	}
	_, err = srv.repo.FindTenantMember(ctx, tenantID, adminID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrTenantForbidden
	}
	return err
}

func (srv authService) ListRouteMenus(ctx context.Context, identity Identity) ([]RouteMenu, error) {
	menus, err := srv.repo.ListVisibleRouteMenus(ctx)
	if err != nil {
		return nil, err
	}
	menuPermissions, err := srv.repo.ListMenuPermissionCodes(ctx)
	if err != nil {
		return nil, err
	}
	if identity.IsSuper {
		return toRouteMenus(menus, menuPermissions), nil
	}
	allowedPermissions := toSet(identity.Permissions)
	allowedMenuIDs := make(map[uint64]struct{})
	menuByID := make(map[uint64]makeadmin.Menu, len(menus))

	for _, menu := range menus {
		menuByID[menu.ID] = menu
		for _, code := range menuPermissions[menu.ID] {
			if _, ok := allowedPermissions[code]; ok {
				allowedMenuIDs[menu.ID] = struct{}{}
				break
			}
		}
	}

	for menuID := range allowedMenuIDs {
		for parentID := menuByID[menuID].ParentID; parentID != 0; parentID = menuByID[parentID].ParentID {
			if _, ok := menuByID[parentID]; !ok {
				break
			}
			allowedMenuIDs[parentID] = struct{}{}
		}
	}

	filtered := make([]makeadmin.Menu, 0, len(allowedMenuIDs))
	for _, menu := range menus {
		if _, ok := allowedMenuIDs[menu.ID]; ok {
			filtered = append(filtered, menu)
		}
	}
	sortMenus(filtered)
	return toRouteMenus(filtered, menuPermissions), nil
}

func toRouteMenus(menus []makeadmin.Menu, menuPermissions map[uint64][]string) []RouteMenu {
	sortMenus(menus)
	result := make([]RouteMenu, 0, len(menus))
	for _, menu := range menus {
		perms := ""
		if codes := menuPermissions[menu.ID]; len(codes) > 0 {
			perms = codes[0]
		}
		result = append(result, RouteMenu{
			ID:         menu.ID,
			ParentID:   menu.ParentID,
			MenuType:   menu.MenuType,
			Name:       menu.Name,
			Perms:      perms,
			Icon:       menu.Icon,
			RoutePath:  menu.RoutePath,
			RouteName:  menu.RouteName,
			Component:  menu.Component,
			Redirect:   menu.Redirect,
			ActivePath: menu.ActivePath,
			Meta:       menu.Meta,
			IsCache:    menu.IsCache,
			Sort:       menu.Sort,
		})
	}
	return result
}

func sortMenus(menus []makeadmin.Menu) {
	sort.SliceStable(menus, func(i, j int) bool {
		if menus[i].Sort == menus[j].Sort {
			return menus[i].ID < menus[j].ID
		}
		return menus[i].Sort > menus[j].Sort
	})
}

func toSet(values []string) map[string]struct{} {
	result := make(map[string]struct{}, len(values))
	for _, value := range values {
		result[value] = struct{}{}
	}
	return result
}

type loginLogInput struct {
	TenantID uint64
	AdminID  uint64
	Username string
	IP       string
	OS       string
	Browser  string
	Message  string
}

func loginLogInputFromAdmin(input LoginInput, admin makeadmin.Admin, message string) loginLogInput {
	return loginLogInput{
		TenantID: input.TenantID,
		AdminID:  admin.ID,
		Username: admin.Username,
		IP:       input.IP,
		OS:       input.OS,
		Browser:  input.Browser,
		Message:  message,
	}
}

func (srv authService) recordLoginLog(ctx context.Context, input loginLogInput) error {
	status := uint8(1)
	if input.Message != "" {
		status = 0
	}
	return srv.repo.CreateLoginLog(ctx, makeadmin.LoginLog{
		TenantID: input.TenantID,
		AdminID:  input.AdminID,
		Username: input.Username,
		IP:       input.IP,
		OS:       input.OS,
		Browser:  input.Browser,
		Status:   status,
		Message:  input.Message,
	})
}
