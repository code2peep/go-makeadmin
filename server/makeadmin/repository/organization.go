package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"go-makeadmin/model/makeadmin"
)

type OrgUnitFilter struct {
	Name      string
	Status    uint8
	StatusSet bool
}

type PositionFilter struct {
	Code      string
	Name      string
	Status    uint8
	StatusSet bool
}

type OrgUnitRepository interface {
	ListOrgUnits(ctx context.Context, tenantID uint64, filter OrgUnitFilter) ([]makeadmin.OrgUnit, error)
	FindOrgUnitByID(ctx context.Context, tenantID uint64, id uint64) (makeadmin.OrgUnit, error)
	CountRootOrgUnits(ctx context.Context, tenantID uint64, excludeID uint64) (int64, error)
	CountChildOrgUnits(ctx context.Context, tenantID uint64, parentID uint64) (int64, error)
	CountActiveAdminsByOrgID(ctx context.Context, tenantID uint64, orgID uint64) (int64, error)
	CreateOrgUnit(ctx context.Context, org makeadmin.OrgUnit) error
	UpdateOrgUnit(ctx context.Context, org makeadmin.OrgUnit) error
	DeleteOrgUnit(ctx context.Context, tenantID uint64, id uint64) error
}

type PositionRepository interface {
	ListAllPositions(ctx context.Context, tenantID uint64) ([]makeadmin.Position, error)
	ListPositions(ctx context.Context, tenantID uint64, filter PositionFilter, limit int, offset int) ([]makeadmin.Position, int64, error)
	FindPositionByID(ctx context.Context, tenantID uint64, id uint64) (makeadmin.Position, error)
	CountPositionsByCode(ctx context.Context, tenantID uint64, code string, excludeID uint64) (int64, error)
	CountPositionsByName(ctx context.Context, tenantID uint64, name string, excludeID uint64) (int64, error)
	CountActiveAdminsByPositionID(ctx context.Context, tenantID uint64, positionID uint64) (int64, error)
	CreatePosition(ctx context.Context, position makeadmin.Position) error
	UpdatePosition(ctx context.Context, position makeadmin.Position) error
	DeletePosition(ctx context.Context, tenantID uint64, id uint64) error
}

type orgUnitRepository struct {
	db *gorm.DB
}

type positionRepository struct {
	db *gorm.DB
}

func NewOrgUnitRepository(db *gorm.DB) OrgUnitRepository {
	return orgUnitRepository{db: db}
}

func NewPositionRepository(db *gorm.DB) PositionRepository {
	return positionRepository{db: db}
}

func (repo orgUnitRepository) ListOrgUnits(ctx context.Context, tenantID uint64, filter OrgUnitFilter) ([]makeadmin.OrgUnit, error) {
	query := repo.db.WithContext(ctx).
		Where("tenant_id = ? AND delete_time = ?", tenantID, 0)
	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	if filter.StatusSet {
		query = query.Where("status = ?", filter.Status)
	}
	var orgs []makeadmin.OrgUnit
	err := query.Order("sort DESC, id DESC").Find(&orgs).Error
	return orgs, err
}

func (repo orgUnitRepository) FindOrgUnitByID(ctx context.Context, tenantID uint64, id uint64) (org makeadmin.OrgUnit, err error) {
	err = repo.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ? AND delete_time = ?", tenantID, id, 0).
		Limit(1).
		First(&org).
		Error
	return
}

func (repo orgUnitRepository) CountRootOrgUnits(ctx context.Context, tenantID uint64, excludeID uint64) (int64, error) {
	query := repo.db.WithContext(ctx).
		Model(&makeadmin.OrgUnit{}).
		Where("tenant_id = ? AND parent_id = ? AND delete_time = ?", tenantID, 0, 0)
	if excludeID > 0 {
		query = query.Where("id <> ?", excludeID)
	}
	var count int64
	err := query.Count(&count).Error
	return count, err
}

func (repo orgUnitRepository) CountChildOrgUnits(ctx context.Context, tenantID uint64, parentID uint64) (int64, error) {
	var count int64
	err := repo.db.WithContext(ctx).
		Model(&makeadmin.OrgUnit{}).
		Where("tenant_id = ? AND parent_id = ? AND delete_time = ?", tenantID, parentID, 0).
		Count(&count).
		Error
	return count, err
}

func (repo orgUnitRepository) CountActiveAdminsByOrgID(ctx context.Context, tenantID uint64, orgID uint64) (int64, error) {
	var count int64
	err := repo.db.WithContext(ctx).
		Model(&makeadmin.AdminOrg{}).
		Joins("INNER JOIN ma_admin AS admin ON admin.id = ma_admin_org.admin_id").
		Where("ma_admin_org.tenant_id = ? AND ma_admin_org.org_id = ? AND ma_admin_org.delete_time = ? AND admin.delete_time = ?", tenantID, orgID, 0, 0).
		Count(&count).
		Error
	return count, err
}

func (repo orgUnitRepository) CreateOrgUnit(ctx context.Context, org makeadmin.OrgUnit) error {
	return repo.db.WithContext(ctx).Create(&org).Error
}

func (repo orgUnitRepository) UpdateOrgUnit(ctx context.Context, org makeadmin.OrgUnit) error {
	return repo.db.WithContext(ctx).
		Model(&makeadmin.OrgUnit{}).
		Where("tenant_id = ? AND id = ? AND delete_time = ?", org.TenantID, org.ID, 0).
		Updates(map[string]interface{}{
			"parent_id":   org.ParentID,
			"name":        org.Name,
			"status":      org.Status,
			"sort":        org.Sort,
			"update_time": time.Now().Unix(),
		}).
		Error
}

func (repo orgUnitRepository) DeleteOrgUnit(ctx context.Context, tenantID uint64, id uint64) error {
	now := time.Now().Unix()
	return repo.db.WithContext(ctx).
		Model(&makeadmin.OrgUnit{}).
		Where("tenant_id = ? AND id = ? AND delete_time = ?", tenantID, id, 0).
		Updates(map[string]interface{}{
			"delete_time": now,
			"update_time": now,
		}).
		Error
}

func (repo positionRepository) ListAllPositions(ctx context.Context, tenantID uint64) ([]makeadmin.Position, error) {
	var positions []makeadmin.Position
	err := repo.db.WithContext(ctx).
		Where("tenant_id = ? AND delete_time = ?", tenantID, 0).
		Order("sort DESC, id DESC").
		Find(&positions).
		Error
	return positions, err
}

func (repo positionRepository) ListPositions(ctx context.Context, tenantID uint64, filter PositionFilter, limit int, offset int) ([]makeadmin.Position, int64, error) {
	query := repo.db.WithContext(ctx).
		Model(&makeadmin.Position{}).
		Where("tenant_id = ? AND delete_time = ?", tenantID, 0)
	if filter.Code != "" {
		query = query.Where("code LIKE ?", "%"+filter.Code+"%")
	}
	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	if filter.StatusSet {
		query = query.Where("status = ?", filter.Status)
	}
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	var positions []makeadmin.Position
	err := query.
		Limit(limit).
		Offset(offset).
		Order("sort DESC, id DESC").
		Find(&positions).
		Error
	return positions, count, err
}

func (repo positionRepository) FindPositionByID(ctx context.Context, tenantID uint64, id uint64) (position makeadmin.Position, err error) {
	err = repo.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ? AND delete_time = ?", tenantID, id, 0).
		Limit(1).
		First(&position).
		Error
	return
}

func (repo positionRepository) CountPositionsByCode(ctx context.Context, tenantID uint64, code string, excludeID uint64) (int64, error) {
	query := repo.db.WithContext(ctx).
		Model(&makeadmin.Position{}).
		Where("tenant_id = ? AND code = ? AND delete_time = ?", tenantID, code, 0)
	if excludeID > 0 {
		query = query.Where("id <> ?", excludeID)
	}
	var count int64
	err := query.Count(&count).Error
	return count, err
}

func (repo positionRepository) CountPositionsByName(ctx context.Context, tenantID uint64, name string, excludeID uint64) (int64, error) {
	query := repo.db.WithContext(ctx).
		Model(&makeadmin.Position{}).
		Where("tenant_id = ? AND name = ? AND delete_time = ?", tenantID, name, 0)
	if excludeID > 0 {
		query = query.Where("id <> ?", excludeID)
	}
	var count int64
	err := query.Count(&count).Error
	return count, err
}

func (repo positionRepository) CountActiveAdminsByPositionID(ctx context.Context, tenantID uint64, positionID uint64) (int64, error) {
	var count int64
	err := repo.db.WithContext(ctx).
		Model(&makeadmin.AdminOrg{}).
		Joins("INNER JOIN ma_admin AS admin ON admin.id = ma_admin_org.admin_id").
		Where("ma_admin_org.tenant_id = ? AND ma_admin_org.position_id = ? AND ma_admin_org.delete_time = ? AND admin.delete_time = ?", tenantID, positionID, 0, 0).
		Count(&count).
		Error
	return count, err
}

func (repo positionRepository) CreatePosition(ctx context.Context, position makeadmin.Position) error {
	return repo.db.WithContext(ctx).Create(&position).Error
}

func (repo positionRepository) UpdatePosition(ctx context.Context, position makeadmin.Position) error {
	return repo.db.WithContext(ctx).
		Model(&makeadmin.Position{}).
		Where("tenant_id = ? AND id = ? AND delete_time = ?", position.TenantID, position.ID, 0).
		Updates(map[string]interface{}{
			"code":        position.Code,
			"name":        position.Name,
			"remark":      position.Remark,
			"status":      position.Status,
			"sort":        position.Sort,
			"update_time": time.Now().Unix(),
		}).
		Error
}

func (repo positionRepository) DeletePosition(ctx context.Context, tenantID uint64, id uint64) error {
	now := time.Now().Unix()
	return repo.db.WithContext(ctx).
		Model(&makeadmin.Position{}).
		Where("tenant_id = ? AND id = ? AND delete_time = ?", tenantID, id, 0).
		Updates(map[string]interface{}{
			"delete_time": now,
			"update_time": now,
		}).
		Error
}
