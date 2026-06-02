package service

import (
	"context"
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"go-makeadmin/makeadmin/security"
	"go-makeadmin/model/makeadmin"
)

type fakeAuthRepository struct {
	admin          makeadmin.Admin
	profile        makeadmin.AdminProfile
	roleIDs        []uint64
	permissionCode []string
	dataScopes     []makeadmin.DataScope
	primaryOrg     makeadmin.AdminOrg
	orgs           []makeadmin.OrgUnit
	menus          []makeadmin.Menu
	menuPerms      map[uint64][]string
	loginLogs      []makeadmin.LoginLog
	lastLoginIP    string
	lastLoginTime  int64
}

func (repo *fakeAuthRepository) FindAdminByUsername(context.Context, string) (makeadmin.Admin, error) {
	return repo.admin, nil
}

func (repo *fakeAuthRepository) FindAdminByID(context.Context, uint64) (makeadmin.Admin, error) {
	return repo.admin, nil
}

func (repo *fakeAuthRepository) FindAdminProfileByAdminID(context.Context, uint64) (makeadmin.AdminProfile, error) {
	return repo.profile, nil
}

func (repo *fakeAuthRepository) ListRoleIDsByAdminID(context.Context, uint64, uint64) ([]uint64, error) {
	return repo.roleIDs, nil
}

func (repo *fakeAuthRepository) ListPermissionCodesByRoleIDs(context.Context, uint64, []uint64) ([]string, error) {
	return repo.permissionCode, nil
}

func (repo *fakeAuthRepository) FindPrimaryAdminOrg(context.Context, uint64, uint64) (makeadmin.AdminOrg, error) {
	if repo.primaryOrg.OrgID == 0 {
		return makeadmin.AdminOrg{}, gorm.ErrRecordNotFound
	}
	return repo.primaryOrg, nil
}

func (repo *fakeAuthRepository) ListDataScopesByRoleIDs(context.Context, uint64, []uint64) ([]makeadmin.DataScope, error) {
	return repo.dataScopes, nil
}

func (repo *fakeAuthRepository) ListOrgUnits(context.Context, uint64) ([]makeadmin.OrgUnit, error) {
	return repo.orgs, nil
}

func (repo *fakeAuthRepository) ListVisibleRouteMenus(context.Context) ([]makeadmin.Menu, error) {
	return repo.menus, nil
}

func (repo *fakeAuthRepository) ListMenuPermissionCodes(context.Context) (map[uint64][]string, error) {
	return repo.menuPerms, nil
}

func (repo *fakeAuthRepository) UpdateAdminLoginInfo(ctx context.Context, adminID uint64, ip string, loginTime int64) error {
	repo.lastLoginIP = ip
	repo.lastLoginTime = loginTime
	return nil
}

func (repo *fakeAuthRepository) CreateLoginLog(ctx context.Context, loginLog makeadmin.LoginLog) error {
	repo.loginLogs = append(repo.loginLogs, loginLog)
	return nil
}

type fixedTokenCodec struct {
	token     string
	sessionID string
}

func (codec fixedTokenCodec) Issue(identity Identity, ttlSeconds int) (SessionToken, error) {
	return SessionToken{AccessToken: codec.token, SessionID: codec.sessionID}, nil
}

func (codec fixedTokenCodec) Parse(token string) (TokenClaims, error) {
	if token != codec.token {
		return TokenClaims{}, ErrTokenInvalid
	}
	return TokenClaims{SessionID: codec.sessionID, AdminID: 1, TenantID: makeadmin.GlobalTenantID}, nil
}

type fakeSessionStore struct {
	sessionID string
	ttl       int
	identity  Identity
	deleted   string
	refreshed string
}

func (store *fakeSessionStore) Save(ctx context.Context, sessionID string, identity Identity, ttlSeconds int) error {
	store.sessionID = sessionID
	store.identity = identity
	store.ttl = ttlSeconds
	return nil
}

func (store *fakeSessionStore) Delete(ctx context.Context, sessionID string) error {
	store.deleted = sessionID
	return nil
}

func (store *fakeSessionStore) FindAdminID(ctx context.Context, sessionID string) (uint64, error) {
	if sessionID != store.sessionID {
		return 0, ErrTokenInvalid
	}
	return store.identity.AdminID, nil
}

func (store *fakeSessionStore) Refresh(ctx context.Context, sessionID string, ttlSeconds int) error {
	store.refreshed = sessionID
	store.ttl = ttlSeconds
	return nil
}

func TestBuildIdentityByUsernameSuperAdmin(t *testing.T) {
	srv := NewAuthService(&fakeAuthRepository{
		admin: makeadmin.Admin{
			ID:       1,
			Username: "admin",
			IsSuper:  1,
			Status:   makeadmin.StatusEnabled,
		},
		profile: makeadmin.AdminProfile{
			AdminID:  1,
			Nickname: "Admin",
			Avatar:   "/api/static/backend_avatar.png",
		},
	})

	identity, err := srv.BuildIdentityByUsername(context.Background(), makeadmin.GlobalTenantID, "admin")
	if err != nil {
		t.Fatalf("BuildIdentityByUsername() error = %v", err)
	}
	if !identity.IsSuper {
		t.Fatal("BuildIdentityByUsername() expected super admin identity")
	}
	if len(identity.Permissions) != 1 || identity.Permissions[0] != "*" {
		t.Fatalf("BuildIdentityByUsername() permissions = %#v, want wildcard", identity.Permissions)
	}
	if !identity.DataScope.Enabled || !identity.DataScope.All {
		t.Fatalf("BuildIdentityByUsername() data scope = %#v, want all", identity.DataScope)
	}
}

func TestBuildIdentityByUsernameResolvesOrgTreeDataScope(t *testing.T) {
	srv := NewAuthService(&fakeAuthRepository{
		admin: makeadmin.Admin{
			ID:       7,
			Username: "operator",
			Status:   makeadmin.StatusEnabled,
		},
		roleIDs:        []uint64{2},
		permissionCode: []string{"system:admin:list"},
		dataScopes: []makeadmin.DataScope{{
			ID:        1,
			TenantID:  makeadmin.GlobalTenantID,
			ScopeType: makeadmin.ScopeTypeOrgTree,
			Status:    makeadmin.StatusEnabled,
		}},
		primaryOrg: makeadmin.AdminOrg{TenantID: makeadmin.GlobalTenantID, AdminID: 7, OrgID: 10, Status: makeadmin.StatusEnabled},
		orgs: []makeadmin.OrgUnit{
			{ID: 10, TenantID: makeadmin.GlobalTenantID, ParentID: 0, Status: makeadmin.StatusEnabled},
			{ID: 11, TenantID: makeadmin.GlobalTenantID, ParentID: 10, Status: makeadmin.StatusEnabled},
			{ID: 12, TenantID: makeadmin.GlobalTenantID, ParentID: 11, Status: makeadmin.StatusEnabled},
		},
	})

	identity, err := srv.BuildIdentityByUsername(context.Background(), makeadmin.GlobalTenantID, "operator")
	if err != nil {
		t.Fatalf("BuildIdentityByUsername() error = %v", err)
	}
	if identity.DataScope.All || identity.DataScope.Self || identity.DataScope.NoAccess {
		t.Fatalf("BuildIdentityByUsername() data scope flags = %#v", identity.DataScope)
	}
	wantOrgIDs := []uint64{10, 11, 12}
	if len(identity.DataScope.OrgIDs) != len(wantOrgIDs) {
		t.Fatalf("BuildIdentityByUsername() org ids = %#v, want %#v", identity.DataScope.OrgIDs, wantOrgIDs)
	}
	for i, id := range wantOrgIDs {
		if identity.DataScope.OrgIDs[i] != id {
			t.Fatalf("BuildIdentityByUsername() org ids = %#v, want %#v", identity.DataScope.OrgIDs, wantOrgIDs)
		}
	}
}

func TestAuthenticateByUsernameVerifiesBcryptPassword(t *testing.T) {
	hasher := security.NewBcryptPasswordHasher(bcrypt.MinCost)
	digest, err := hasher.Hash("makeadmin-secret")
	if err != nil {
		t.Fatalf("Hash() error = %v", err)
	}
	srv := NewAuthServiceWithPasswordHasher(&fakeAuthRepository{
		admin: makeadmin.Admin{
			ID:           1,
			Username:     "admin",
			PasswordHash: digest.Hash,
			PasswordSalt: digest.Salt,
			IsSuper:      1,
			Status:       makeadmin.StatusEnabled,
		},
	}, hasher)

	identity, err := srv.AuthenticateByUsername(context.Background(), makeadmin.GlobalTenantID, "admin", "makeadmin-secret")
	if err != nil {
		t.Fatalf("AuthenticateByUsername() error = %v", err)
	}
	if identity.AdminID != 1 || !identity.IsSuper {
		t.Fatalf("AuthenticateByUsername() identity = %#v, want super admin", identity)
	}

	_, err = srv.AuthenticateByUsername(context.Background(), makeadmin.GlobalTenantID, "admin", "wrong-secret")
	if !errors.Is(err, ErrInvalidCredential) {
		t.Fatalf("AuthenticateByUsername() wrong password error = %v, want ErrInvalidCredential", err)
	}
}

func TestLoginWritesSessionAndAudit(t *testing.T) {
	hasher := security.NewBcryptPasswordHasher(bcrypt.MinCost)
	digest, err := hasher.Hash("makeadmin-secret")
	if err != nil {
		t.Fatalf("Hash() error = %v", err)
	}
	repo := &fakeAuthRepository{
		admin: makeadmin.Admin{
			ID:           1,
			Username:     "admin",
			PasswordHash: digest.Hash,
			IsSuper:      1,
			Status:       makeadmin.StatusEnabled,
		},
	}
	store := &fakeSessionStore{}
	srv := NewAuthServiceWithDependencies(repo, hasher, fixedTokenCodec{token: "fixed-jwt", sessionID: "fixed-session"}, store, 3600)

	result, err := srv.Login(context.Background(), LoginInput{
		TenantID: makeadmin.GlobalTenantID,
		Username: "admin",
		Password: "makeadmin-secret",
		IP:       "127.0.0.1",
		OS:       "macOS",
		Browser:  "Chrome",
	})
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if result.Token != "fixed-jwt" || result.ExpiresIn != 3600 {
		t.Fatalf("Login() result = %#v, want fixed JWT and ttl", result)
	}
	if store.sessionID != "fixed-session" || store.identity.AdminID != 1 || store.ttl != 3600 {
		t.Fatalf("session store = %#v, want saved identity", store)
	}
	if repo.lastLoginIP != "127.0.0.1" || repo.lastLoginTime == 0 {
		t.Fatalf("last login ip/time = %q/%d, want updated", repo.lastLoginIP, repo.lastLoginTime)
	}
	if len(repo.loginLogs) != 1 || repo.loginLogs[0].Status != 1 {
		t.Fatalf("login logs = %#v, want one success log", repo.loginLogs)
	}
}

func TestLoginRecordsFailedPassword(t *testing.T) {
	hasher := security.NewBcryptPasswordHasher(bcrypt.MinCost)
	digest, err := hasher.Hash("makeadmin-secret")
	if err != nil {
		t.Fatalf("Hash() error = %v", err)
	}
	repo := &fakeAuthRepository{
		admin: makeadmin.Admin{
			ID:           1,
			Username:     "admin",
			PasswordHash: digest.Hash,
			IsSuper:      1,
			Status:       makeadmin.StatusEnabled,
		},
	}
	store := &fakeSessionStore{}
	srv := NewAuthServiceWithDependencies(repo, hasher, fixedTokenCodec{token: "fixed-jwt", sessionID: "fixed-session"}, store, 3600)

	_, err = srv.Login(context.Background(), LoginInput{
		TenantID: makeadmin.GlobalTenantID,
		Username: "admin",
		Password: "wrong-secret",
		IP:       "127.0.0.1",
	})
	if !errors.Is(err, ErrInvalidCredential) {
		t.Fatalf("Login() error = %v, want ErrInvalidCredential", err)
	}
	if store.sessionID != "" || repo.lastLoginTime != 0 {
		t.Fatalf("failed login touched session/update: sessionID=%q lastLogin=%d", store.sessionID, repo.lastLoginTime)
	}
	if len(repo.loginLogs) != 1 || repo.loginLogs[0].Status != 0 || repo.loginLogs[0].AdminID != 1 {
		t.Fatalf("login logs = %#v, want one failure log for admin", repo.loginLogs)
	}
}

func TestLogoutDeletesSessionState(t *testing.T) {
	store := &fakeSessionStore{}
	srv := NewAuthServiceWithDependencies(
		&fakeAuthRepository{},
		nil,
		fixedTokenCodec{token: "fixed-jwt", sessionID: "fixed-session"},
		store,
		3600,
	)

	if err := srv.Logout(context.Background(), "fixed-jwt"); err != nil {
		t.Fatalf("Logout() error = %v", err)
	}
	if store.deleted != "fixed-session" {
		t.Fatalf("deleted session = %q, want fixed-session", store.deleted)
	}
}

func TestListRouteMenusIncludesParentCatalog(t *testing.T) {
	srv := NewAuthService(&fakeAuthRepository{
		menus: []makeadmin.Menu{
			{ID: 100, MenuType: makeadmin.MenuTypeCatalog, Name: "权限管理", RoutePath: "/permission", Sort: 100},
			{ID: 101, ParentID: 100, MenuType: makeadmin.MenuTypePage, Name: "管理员", RoutePath: "/permission/admin", Sort: 10},
			{ID: 200, MenuType: makeadmin.MenuTypePage, Name: "未授权", RoutePath: "/hidden", Sort: 1},
		},
		menuPerms: map[uint64][]string{
			101: {"system:admin:list"},
			200: {"hidden:view"},
		},
	})
	identity := Identity{
		TenantID:    makeadmin.GlobalTenantID,
		RoleIDs:     []uint64{2},
		Permissions: []string{"system:admin:list"},
	}

	menus, err := srv.ListRouteMenus(context.Background(), identity)
	if err != nil {
		t.Fatalf("ListRouteMenus() error = %v", err)
	}
	if len(menus) != 2 {
		t.Fatalf("ListRouteMenus() len = %d, want 2", len(menus))
	}
	if menus[0].ID != 100 || menus[1].ID != 101 {
		t.Fatalf("ListRouteMenus() ids = [%d, %d], want [100, 101]", menus[0].ID, menus[1].ID)
	}
}
