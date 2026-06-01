package service

import (
	"context"
	"errors"
	"testing"

	"gorm.io/gorm"

	"go-makeadmin/makeadmin/repository"
	"go-makeadmin/makeadmin/security"
	"go-makeadmin/model/makeadmin"
)

type fakeAdminRepository struct {
	admins             []makeadmin.Admin
	profiles           []makeadmin.AdminProfile
	roles              []makeadmin.Role
	orgs               []makeadmin.OrgUnit
	positions          []makeadmin.Position
	adminRoles         map[uint64][]uint64
	adminOrgs          map[uint64]makeadmin.AdminOrg
	createdAdmin       makeadmin.Admin
	createdProfile     makeadmin.AdminProfile
	createdRoleID      uint64
	createdAdminOrg    makeadmin.AdminOrg
	updatedAdmin       makeadmin.Admin
	updatedProfile     makeadmin.AdminProfile
	updatedRoleID      uint64
	updatedAdminOrg    makeadmin.AdminOrg
	deletedAdminID     uint64
	toggledAdminID     uint64
	selfUpdatedAdmin   makeadmin.Admin
	selfUpdatedProfile makeadmin.AdminProfile
}

func (repo *fakeAdminRepository) ListAdmins(ctx context.Context, tenantID uint64, filter repository.AdminFilter, limit int, offset int) ([]makeadmin.Admin, int64, error) {
	return repo.admins, int64(len(repo.admins)), nil
}

func (repo *fakeAdminRepository) FindAdminByID(ctx context.Context, adminID uint64) (makeadmin.Admin, error) {
	for _, admin := range repo.admins {
		if admin.ID == adminID && admin.DeleteTime == 0 {
			return admin, nil
		}
	}
	return makeadmin.Admin{}, gorm.ErrRecordNotFound
}

func (repo *fakeAdminRepository) FindAdminByUsername(ctx context.Context, username string) (makeadmin.Admin, error) {
	for _, admin := range repo.admins {
		if admin.Username == username && admin.DeleteTime == 0 {
			return admin, nil
		}
	}
	return makeadmin.Admin{}, gorm.ErrRecordNotFound
}

func (repo *fakeAdminRepository) FindAdminProfileByAdminID(ctx context.Context, adminID uint64) (makeadmin.AdminProfile, error) {
	for _, profile := range repo.profiles {
		if profile.AdminID == adminID {
			return profile, nil
		}
	}
	return makeadmin.AdminProfile{}, gorm.ErrRecordNotFound
}

func (repo *fakeAdminRepository) CountAdminsByUsername(ctx context.Context, username string, excludeID uint64) (int64, error) {
	var count int64
	for _, admin := range repo.admins {
		if admin.Username == username && admin.ID != excludeID && admin.DeleteTime == 0 {
			count++
		}
	}
	return count, nil
}

func (repo *fakeAdminRepository) CountProfilesByNickname(ctx context.Context, nickname string, excludeAdminID uint64) (int64, error) {
	var count int64
	for _, profile := range repo.profiles {
		if profile.Nickname == nickname && profile.AdminID != excludeAdminID {
			count++
		}
	}
	return count, nil
}

func (repo *fakeAdminRepository) ListRoleIDsByAdminID(ctx context.Context, tenantID uint64, adminID uint64) ([]uint64, error) {
	if repo.adminRoles == nil {
		return []uint64{}, nil
	}
	return repo.adminRoles[adminID], nil
}

func (repo *fakeAdminRepository) ListRolesByIDs(ctx context.Context, tenantID uint64, roleIDs []uint64) ([]makeadmin.Role, error) {
	result := make([]makeadmin.Role, 0, len(roleIDs))
	for _, role := range repo.roles {
		for _, roleID := range roleIDs {
			if role.ID == roleID && role.TenantID == tenantID && role.DeleteTime == 0 {
				result = append(result, role)
			}
		}
	}
	return result, nil
}

func (repo *fakeAdminRepository) FindRoleByID(ctx context.Context, tenantID uint64, roleID uint64) (makeadmin.Role, error) {
	for _, role := range repo.roles {
		if role.ID == roleID && role.TenantID == tenantID && role.DeleteTime == 0 {
			return role, nil
		}
	}
	return makeadmin.Role{}, gorm.ErrRecordNotFound
}

func (repo *fakeAdminRepository) FindPrimaryAdminOrg(ctx context.Context, tenantID uint64, adminID uint64) (makeadmin.AdminOrg, error) {
	if repo.adminOrgs == nil {
		return makeadmin.AdminOrg{}, gorm.ErrRecordNotFound
	}
	adminOrg, ok := repo.adminOrgs[adminID]
	if !ok || adminOrg.DeleteTime > 0 {
		return makeadmin.AdminOrg{}, gorm.ErrRecordNotFound
	}
	return adminOrg, nil
}

func (repo *fakeAdminRepository) FindOrgUnitByID(ctx context.Context, tenantID uint64, orgID uint64) (makeadmin.OrgUnit, error) {
	for _, org := range repo.orgs {
		if org.ID == orgID && org.TenantID == tenantID && org.DeleteTime == 0 {
			return org, nil
		}
	}
	return makeadmin.OrgUnit{}, gorm.ErrRecordNotFound
}

func (repo *fakeAdminRepository) FindPositionByID(ctx context.Context, tenantID uint64, positionID uint64) (makeadmin.Position, error) {
	for _, position := range repo.positions {
		if position.ID == positionID && position.TenantID == tenantID && position.DeleteTime == 0 {
			return position, nil
		}
	}
	return makeadmin.Position{}, gorm.ErrRecordNotFound
}

func (repo *fakeAdminRepository) CreateAdminWithRelations(ctx context.Context, admin makeadmin.Admin, profile makeadmin.AdminProfile, roleID uint64, adminOrg makeadmin.AdminOrg) (uint64, error) {
	admin.ID = 9
	profile.AdminID = admin.ID
	adminOrg.AdminID = admin.ID
	repo.createdAdmin = admin
	repo.createdProfile = profile
	repo.createdRoleID = roleID
	repo.createdAdminOrg = adminOrg
	return admin.ID, nil
}

func (repo *fakeAdminRepository) UpdateAdminWithRelations(ctx context.Context, admin makeadmin.Admin, profile makeadmin.AdminProfile, roleID uint64, adminOrg makeadmin.AdminOrg, updatePassword bool) error {
	repo.updatedAdmin = admin
	repo.updatedProfile = profile
	repo.updatedRoleID = roleID
	repo.updatedAdminOrg = adminOrg
	return nil
}

func (repo *fakeAdminRepository) UpdateAdminSelf(ctx context.Context, admin makeadmin.Admin, profile makeadmin.AdminProfile, updatePassword bool) error {
	repo.selfUpdatedAdmin = admin
	repo.selfUpdatedProfile = profile
	return nil
}

func (repo *fakeAdminRepository) SoftDeleteAdmin(ctx context.Context, tenantID uint64, adminID uint64) error {
	repo.deletedAdminID = adminID
	return nil
}

func (repo *fakeAdminRepository) ToggleAdminStatus(ctx context.Context, admin makeadmin.Admin) error {
	repo.toggledAdminID = admin.ID
	return nil
}

type fakePasswordHasher struct {
	matched bool
}

func (hasher fakePasswordHasher) Hash(plain string) (security.PasswordDigest, error) {
	if len(plain) < 8 {
		return security.PasswordDigest{}, security.ErrPasswordTooShort
	}
	return security.PasswordDigest{Hash: "hash:" + plain}, nil
}

func (hasher fakePasswordHasher) Verify(plain string, digest security.PasswordDigest) (bool, error) {
	return hasher.matched, nil
}

func (hasher fakePasswordHasher) NeedsUpgrade(digest security.PasswordDigest) bool {
	return false
}

func TestAdminAddCreatesAdminRelations(t *testing.T) {
	repo := newAdminRepoFixture()
	srv := NewAdminServiceWithPasswordHasher(repo, fakePasswordHasher{matched: true})

	err := srv.Add(context.Background(), AdminInput{
		TenantID:   makeadmin.GlobalTenantID,
		OrgID:      1,
		PositionID: 1,
		Username:   " operator ",
		Nickname:   " 运营 ",
		Password:   "password123",
		Avatar:     "/api/static/avatar.png",
		RoleID:     2,
	})
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}
	if repo.createdAdmin.Username != "operator" ||
		repo.createdProfile.Nickname != "运营" ||
		repo.createdRoleID != 2 ||
		repo.createdAdminOrg.OrgID != 1 ||
		repo.createdAdmin.PasswordHash != "hash:password123" {
		t.Fatalf("Add() created admin=%#v profile=%#v role=%d org=%#v", repo.createdAdmin, repo.createdProfile, repo.createdRoleID, repo.createdAdminOrg)
	}
}

func TestAdminAddRejectsDuplicateUsername(t *testing.T) {
	repo := newAdminRepoFixture()
	repo.admins = append(repo.admins, makeadmin.Admin{ID: 8, Username: "operator"})
	srv := NewAdminServiceWithPasswordHasher(repo, fakePasswordHasher{matched: true})

	err := srv.Add(context.Background(), AdminInput{
		TenantID:   makeadmin.GlobalTenantID,
		OrgID:      1,
		PositionID: 1,
		Username:   "operator",
		Nickname:   "运营",
		Password:   "password123",
		RoleID:     2,
	})
	if !errors.Is(err, ErrAdminUsernameExists) {
		t.Fatalf("Add() error = %v, want ErrAdminUsernameExists", err)
	}
}

func TestAdminDetailMapsRelations(t *testing.T) {
	repo := newAdminRepoFixture()
	repo.admins = append(repo.admins, makeadmin.Admin{ID: 8, Username: "operator", Status: makeadmin.StatusDisabled})
	repo.profiles = append(repo.profiles, makeadmin.AdminProfile{AdminID: 8, Nickname: "运营", Avatar: "/avatar.png"})
	repo.adminRoles = map[uint64][]uint64{8: {2}}
	repo.adminOrgs = map[uint64]makeadmin.AdminOrg{8: {TenantID: makeadmin.GlobalTenantID, AdminID: 8, OrgID: 1, PositionID: 1, IsPrimary: 1}}
	srv := NewAdminServiceWithPasswordHasher(repo, fakePasswordHasher{matched: true})

	item, err := srv.Detail(context.Background(), makeadmin.GlobalTenantID, 8)
	if err != nil {
		t.Fatalf("Detail() error = %v", err)
	}
	if item.RoleID != 2 || item.RoleLabel != "运营角色" || item.OrgName != "总部" || item.IsDisable != 1 {
		t.Fatalf("Detail() = %#v", item)
	}
}

func TestAdminDeleteProtectsSystemAndSelf(t *testing.T) {
	repo := newAdminRepoFixture()
	repo.admins = []makeadmin.Admin{
		{ID: 1, Username: "admin", IsSuper: 1},
		{ID: 2, Username: "operator"},
	}
	srv := NewAdminServiceWithPasswordHasher(repo, fakePasswordHasher{matched: true})

	err := srv.Delete(context.Background(), makeadmin.GlobalTenantID, 2, 1)
	if !errors.Is(err, ErrSystemAdminProtected) {
		t.Fatalf("Delete() system error = %v", err)
	}
	err = srv.Delete(context.Background(), makeadmin.GlobalTenantID, 2, 2)
	if !errors.Is(err, ErrAdminSelfProtected) {
		t.Fatalf("Delete() self error = %v", err)
	}
}

func TestAdminDisableProtectsSelf(t *testing.T) {
	repo := newAdminRepoFixture()
	repo.admins = []makeadmin.Admin{{ID: 2, Username: "operator"}}
	srv := NewAdminServiceWithPasswordHasher(repo, fakePasswordHasher{matched: true})

	err := srv.Disable(context.Background(), 2, 2)
	if !errors.Is(err, ErrAdminSelfProtected) {
		t.Fatalf("Disable() error = %v, want ErrAdminSelfProtected", err)
	}
}

func TestAdminUpdateSelfVerifiesCurrentPassword(t *testing.T) {
	repo := newAdminRepoFixture()
	repo.admins = []makeadmin.Admin{{ID: 2, Username: "operator", PasswordHash: "old"}}
	srv := NewAdminServiceWithPasswordHasher(repo, fakePasswordHasher{matched: false})

	err := srv.UpdateSelf(context.Background(), AdminSelfInput{
		ID:           2,
		Nickname:     "新昵称",
		Password:     "password123",
		CurrPassword: "wrong",
	})
	if !errors.Is(err, ErrAdminPasswordInvalid) {
		t.Fatalf("UpdateSelf() error = %v, want ErrAdminPasswordInvalid", err)
	}

	srv = NewAdminServiceWithPasswordHasher(repo, fakePasswordHasher{matched: true})
	err = srv.UpdateSelf(context.Background(), AdminSelfInput{
		ID:           2,
		Nickname:     "新昵称",
		Password:     "password123",
		CurrPassword: "current123",
	})
	if err != nil {
		t.Fatalf("UpdateSelf() error = %v", err)
	}
	if repo.selfUpdatedAdmin.PasswordHash != "hash:password123" || repo.selfUpdatedProfile.Nickname != "新昵称" {
		t.Fatalf("UpdateSelf() admin=%#v profile=%#v", repo.selfUpdatedAdmin, repo.selfUpdatedProfile)
	}
}

func newAdminRepoFixture() *fakeAdminRepository {
	return &fakeAdminRepository{
		roles: []makeadmin.Role{
			{ID: 2, TenantID: makeadmin.GlobalTenantID, Name: "运营角色", Status: makeadmin.StatusEnabled},
		},
		orgs: []makeadmin.OrgUnit{
			{ID: 1, TenantID: makeadmin.GlobalTenantID, Name: "总部", Status: makeadmin.StatusEnabled},
		},
		positions: []makeadmin.Position{
			{ID: 1, TenantID: makeadmin.GlobalTenantID, Name: "管理员", Status: makeadmin.StatusEnabled},
		},
	}
}
