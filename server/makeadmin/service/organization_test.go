package service

import (
	"context"
	"errors"
	"testing"

	"gorm.io/gorm"

	"go-makeadmin/makeadmin/repository"
	"go-makeadmin/model/makeadmin"
)

type fakeOrgUnitRepository struct {
	orgs             []makeadmin.OrgUnit
	childCount       int64
	orgAdminCount    int64
	createdOrg       makeadmin.OrgUnit
	updatedOrg       makeadmin.OrgUnit
	deletedOrgUnitID uint64
}

func (repo *fakeOrgUnitRepository) ListOrgUnits(ctx context.Context, tenantID uint64, filter repository.OrgUnitFilter) ([]makeadmin.OrgUnit, error) {
	return repo.orgs, nil
}

func (repo *fakeOrgUnitRepository) FindOrgUnitByID(ctx context.Context, tenantID uint64, id uint64) (makeadmin.OrgUnit, error) {
	for _, org := range repo.orgs {
		if org.ID == id && org.TenantID == tenantID && org.DeleteTime == 0 {
			return org, nil
		}
	}
	return makeadmin.OrgUnit{}, gorm.ErrRecordNotFound
}

func (repo *fakeOrgUnitRepository) CountRootOrgUnits(ctx context.Context, tenantID uint64, excludeID uint64) (int64, error) {
	var count int64
	for _, org := range repo.orgs {
		if org.TenantID == tenantID && org.ParentID == 0 && org.ID != excludeID && org.DeleteTime == 0 {
			count++
		}
	}
	return count, nil
}

func (repo *fakeOrgUnitRepository) CountChildOrgUnits(ctx context.Context, tenantID uint64, parentID uint64) (int64, error) {
	return repo.childCount, nil
}

func (repo *fakeOrgUnitRepository) CountActiveAdminsByOrgID(ctx context.Context, tenantID uint64, orgID uint64) (int64, error) {
	return repo.orgAdminCount, nil
}

func (repo *fakeOrgUnitRepository) CreateOrgUnit(ctx context.Context, org makeadmin.OrgUnit) error {
	repo.createdOrg = org
	return nil
}

func (repo *fakeOrgUnitRepository) UpdateOrgUnit(ctx context.Context, org makeadmin.OrgUnit) error {
	repo.updatedOrg = org
	return nil
}

func (repo *fakeOrgUnitRepository) DeleteOrgUnit(ctx context.Context, tenantID uint64, id uint64) error {
	repo.deletedOrgUnitID = id
	return nil
}

type fakePositionRepository struct {
	positions         []makeadmin.Position
	positionAdminCnt  int64
	createdPosition   makeadmin.Position
	updatedPosition   makeadmin.Position
	deletedPositionID uint64
}

func (repo *fakePositionRepository) ListAllPositions(ctx context.Context, tenantID uint64) ([]makeadmin.Position, error) {
	return repo.positions, nil
}

func (repo *fakePositionRepository) ListPositions(ctx context.Context, tenantID uint64, filter repository.PositionFilter, limit int, offset int) ([]makeadmin.Position, int64, error) {
	return repo.positions, int64(len(repo.positions)), nil
}

func (repo *fakePositionRepository) FindPositionByID(ctx context.Context, tenantID uint64, id uint64) (makeadmin.Position, error) {
	for _, position := range repo.positions {
		if position.ID == id && position.TenantID == tenantID && position.DeleteTime == 0 {
			return position, nil
		}
	}
	return makeadmin.Position{}, gorm.ErrRecordNotFound
}

func (repo *fakePositionRepository) CountPositionsByCode(ctx context.Context, tenantID uint64, code string, excludeID uint64) (int64, error) {
	var count int64
	for _, position := range repo.positions {
		if position.TenantID == tenantID && position.Code == code && position.ID != excludeID && position.DeleteTime == 0 {
			count++
		}
	}
	return count, nil
}

func (repo *fakePositionRepository) CountPositionsByName(ctx context.Context, tenantID uint64, name string, excludeID uint64) (int64, error) {
	var count int64
	for _, position := range repo.positions {
		if position.TenantID == tenantID && position.Name == name && position.ID != excludeID && position.DeleteTime == 0 {
			count++
		}
	}
	return count, nil
}

func (repo *fakePositionRepository) CountActiveAdminsByPositionID(ctx context.Context, tenantID uint64, positionID uint64) (int64, error) {
	return repo.positionAdminCnt, nil
}

func (repo *fakePositionRepository) CreatePosition(ctx context.Context, position makeadmin.Position) error {
	repo.createdPosition = position
	return nil
}

func (repo *fakePositionRepository) UpdatePosition(ctx context.Context, position makeadmin.Position) error {
	repo.updatedPosition = position
	return nil
}

func (repo *fakePositionRepository) DeletePosition(ctx context.Context, tenantID uint64, id uint64) error {
	repo.deletedPositionID = id
	return nil
}

func TestOrgUnitAddRejectsDuplicateRoot(t *testing.T) {
	srv := NewOrgUnitService(&fakeOrgUnitRepository{
		orgs: []makeadmin.OrgUnit{{ID: 1, TenantID: makeadmin.GlobalTenantID, ParentID: 0}},
	})

	err := srv.Add(context.Background(), OrgUnitInput{TenantID: makeadmin.GlobalTenantID, ParentID: 0, Name: "总部"})
	if !errors.Is(err, ErrRootOrgUnitExists) {
		t.Fatalf("Add() error = %v, want ErrRootOrgUnitExists", err)
	}
}

func TestOrgUnitAddCreatesChild(t *testing.T) {
	repo := &fakeOrgUnitRepository{
		orgs: []makeadmin.OrgUnit{{ID: 1, TenantID: makeadmin.GlobalTenantID, ParentID: 0}},
	}
	srv := NewOrgUnitService(repo)

	err := srv.Add(context.Background(), OrgUnitInput{
		TenantID: makeadmin.GlobalTenantID,
		ParentID: 1,
		Name:     " 研发部 ",
		IsStop:   1,
		Sort:     12,
	})
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}
	if repo.createdOrg.ParentID != 1 || repo.createdOrg.Name != "研发部" || repo.createdOrg.Status != makeadmin.StatusDisabled {
		t.Fatalf("Add() created = %#v", repo.createdOrg)
	}
}

func TestOrgUnitEditRejectsInvalidParent(t *testing.T) {
	srv := NewOrgUnitService(&fakeOrgUnitRepository{
		orgs: []makeadmin.OrgUnit{{ID: 2, TenantID: makeadmin.GlobalTenantID, ParentID: 1}},
	})

	err := srv.Edit(context.Background(), OrgUnitInput{ID: 2, TenantID: makeadmin.GlobalTenantID, ParentID: 2, Name: "研发部"})
	if !errors.Is(err, ErrOrgUnitSelfParent) {
		t.Fatalf("Edit() error = %v, want ErrOrgUnitSelfParent", err)
	}
}

func TestOrgUnitDeleteRejectsRootChildrenAndUse(t *testing.T) {
	repo := &fakeOrgUnitRepository{
		orgs: []makeadmin.OrgUnit{{ID: 1, TenantID: makeadmin.GlobalTenantID, ParentID: 0}},
	}
	srv := NewOrgUnitService(repo)

	err := srv.Delete(context.Background(), makeadmin.GlobalTenantID, 1)
	if !errors.Is(err, ErrRootOrgUnitDeleteLocked) {
		t.Fatalf("Delete() root error = %v", err)
	}

	repo.orgs = []makeadmin.OrgUnit{{ID: 2, TenantID: makeadmin.GlobalTenantID, ParentID: 1}}
	repo.childCount = 1
	err = srv.Delete(context.Background(), makeadmin.GlobalTenantID, 2)
	if !errors.Is(err, ErrOrgUnitHasChildren) {
		t.Fatalf("Delete() children error = %v", err)
	}

	repo.childCount = 0
	repo.orgAdminCount = 1
	err = srv.Delete(context.Background(), makeadmin.GlobalTenantID, 2)
	if !errors.Is(err, ErrOrgUnitInUse) {
		t.Fatalf("Delete() in-use error = %v", err)
	}
}

func TestPositionAddGeneratesCodeAndRejectsDuplicateName(t *testing.T) {
	repo := &fakePositionRepository{}
	srv := NewPositionService(repo)

	err := srv.Add(context.Background(), PositionInput{
		TenantID: makeadmin.GlobalTenantID,
		Name:     " 管理员 ",
		Remark:   "admin",
		Sort:     10,
	})
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}
	if repo.createdPosition.Code == "" || repo.createdPosition.Name != "管理员" {
		t.Fatalf("Add() created = %#v", repo.createdPosition)
	}

	repo.positions = []makeadmin.Position{{ID: 1, TenantID: makeadmin.GlobalTenantID, Name: "管理员", Code: "admin"}}
	err = srv.Add(context.Background(), PositionInput{TenantID: makeadmin.GlobalTenantID, Code: "ops", Name: "管理员"})
	if !errors.Is(err, ErrPositionNameExists) {
		t.Fatalf("Add() duplicate error = %v", err)
	}
}

func TestPositionEditKeepsCodeWhenBlank(t *testing.T) {
	repo := &fakePositionRepository{
		positions: []makeadmin.Position{{ID: 1, TenantID: makeadmin.GlobalTenantID, Code: "admin", Name: "管理员"}},
	}
	srv := NewPositionService(repo)

	err := srv.Edit(context.Background(), PositionInput{
		ID:       1,
		TenantID: makeadmin.GlobalTenantID,
		Name:     "系统管理员",
	})
	if err != nil {
		t.Fatalf("Edit() error = %v", err)
	}
	if repo.updatedPosition.Code != "admin" || repo.updatedPosition.Name != "系统管理员" {
		t.Fatalf("Edit() updated = %#v", repo.updatedPosition)
	}
}

func TestPositionDeleteRejectsInUse(t *testing.T) {
	repo := &fakePositionRepository{
		positions:        []makeadmin.Position{{ID: 1, TenantID: makeadmin.GlobalTenantID, Code: "admin", Name: "管理员"}},
		positionAdminCnt: 1,
	}
	srv := NewPositionService(repo)

	err := srv.Delete(context.Background(), makeadmin.GlobalTenantID, 1)
	if !errors.Is(err, ErrPositionInUse) {
		t.Fatalf("Delete() error = %v, want ErrPositionInUse", err)
	}
}
