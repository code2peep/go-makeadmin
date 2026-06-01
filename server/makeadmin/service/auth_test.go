package service

import (
	"context"
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"go-makeadmin/makeadmin/security"
	"go-makeadmin/model/makeadmin"
)

type fakeAuthRepository struct {
	admin          makeadmin.Admin
	profile        makeadmin.AdminProfile
	roleIDs        []uint64
	permissionCode []string
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

type fixedTokenGenerator struct {
	token string
}

func (generator fixedTokenGenerator) Generate() (string, error) {
	return generator.token, nil
}

type fakeSessionStore struct {
	token    string
	ttl      int
	identity Identity
	deleted  string
}

func (store *fakeSessionStore) Save(ctx context.Context, token string, identity Identity, ttlSeconds int) error {
	store.token = token
	store.identity = identity
	store.ttl = ttlSeconds
	return nil
}

func (store *fakeSessionStore) Delete(ctx context.Context, token string) error {
	store.deleted = token
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
	srv := NewAuthServiceWithDependencies(repo, hasher, fixedTokenGenerator{token: "fixed-token"}, store, 3600)

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
	if result.Token != "fixed-token" || result.ExpiresIn != 3600 {
		t.Fatalf("Login() result = %#v, want fixed token and ttl", result)
	}
	if store.token != "fixed-token" || store.identity.AdminID != 1 || store.ttl != 3600 {
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
	srv := NewAuthServiceWithDependencies(repo, hasher, fixedTokenGenerator{token: "fixed-token"}, store, 3600)

	_, err = srv.Login(context.Background(), LoginInput{
		TenantID: makeadmin.GlobalTenantID,
		Username: "admin",
		Password: "wrong-secret",
		IP:       "127.0.0.1",
	})
	if !errors.Is(err, ErrInvalidCredential) {
		t.Fatalf("Login() error = %v, want ErrInvalidCredential", err)
	}
	if store.token != "" || repo.lastLoginTime != 0 {
		t.Fatalf("failed login touched session/update: token=%q lastLogin=%d", store.token, repo.lastLoginTime)
	}
	if len(repo.loginLogs) != 1 || repo.loginLogs[0].Status != 0 || repo.loginLogs[0].AdminID != 1 {
		t.Fatalf("login logs = %#v, want one failure log for admin", repo.loginLogs)
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
