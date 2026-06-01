package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"go-makeadmin/model/makeadmin"
)

type DictTypeFilter struct {
	Name      string
	Code      string
	Status    int8
	StatusSet bool
}

type DictItemFilter struct {
	Name      string
	Value     string
	Status    int8
	StatusSet bool
}

type DictRepository interface {
	ListAllDictTypes(ctx context.Context) ([]makeadmin.DictType, error)
	ListDictTypes(ctx context.Context, filter DictTypeFilter, limit int, offset int) ([]makeadmin.DictType, int64, error)
	FindDictTypeByID(ctx context.Context, id uint64) (makeadmin.DictType, error)
	FindDictTypeByCode(ctx context.Context, code string) (makeadmin.DictType, error)
	CountDictTypesByName(ctx context.Context, name string, excludeID uint64) (int64, error)
	CountDictTypesByCode(ctx context.Context, code string, excludeID uint64) (int64, error)
	CreateDictType(ctx context.Context, dictType makeadmin.DictType) error
	UpdateDictType(ctx context.Context, dictType makeadmin.DictType) error
	DeleteDictTypes(ctx context.Context, ids []uint64) error
	ListAllDictItemsByTypeID(ctx context.Context, typeID uint64, filter DictItemFilter) ([]makeadmin.DictItem, error)
	ListDictItemsByTypeID(ctx context.Context, typeID uint64, filter DictItemFilter, limit int, offset int) ([]makeadmin.DictItem, int64, error)
	FindDictItemByID(ctx context.Context, id uint64) (makeadmin.DictItem, error)
	CountDictItemsByValue(ctx context.Context, typeID uint64, value string, excludeID uint64) (int64, error)
	CreateDictItem(ctx context.Context, item makeadmin.DictItem) error
	UpdateDictItem(ctx context.Context, item makeadmin.DictItem) error
	DeleteDictItems(ctx context.Context, ids []uint64) error
}

type dictRepository struct {
	db *gorm.DB
}

func NewDictRepository(db *gorm.DB) DictRepository {
	return &dictRepository{db: db}
}

func (repo dictRepository) ListAllDictTypes(ctx context.Context) ([]makeadmin.DictType, error) {
	var dictTypes []makeadmin.DictType
	err := repo.dictTypeQuery(ctx, DictTypeFilter{}).
		Order("sort DESC, id ASC").
		Find(&dictTypes).
		Error
	return dictTypes, err
}

func (repo dictRepository) ListDictTypes(ctx context.Context, filter DictTypeFilter, limit int, offset int) ([]makeadmin.DictType, int64, error) {
	query := repo.dictTypeQuery(ctx, filter)
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	var dictTypes []makeadmin.DictType
	err := query.
		Limit(limit).
		Offset(offset).
		Order("sort DESC, id ASC").
		Find(&dictTypes).
		Error
	return dictTypes, count, err
}

func (repo dictRepository) FindDictTypeByID(ctx context.Context, id uint64) (dictType makeadmin.DictType, err error) {
	err = repo.db.WithContext(ctx).
		Where("id = ? AND delete_time = ?", id, 0).
		Limit(1).
		First(&dictType).
		Error
	return
}

func (repo dictRepository) FindDictTypeByCode(ctx context.Context, code string) (dictType makeadmin.DictType, err error) {
	err = repo.db.WithContext(ctx).
		Where("code = ? AND delete_time = ?", code, 0).
		Limit(1).
		First(&dictType).
		Error
	return
}

func (repo dictRepository) CountDictTypesByName(ctx context.Context, name string, excludeID uint64) (int64, error) {
	query := repo.db.WithContext(ctx).Model(&makeadmin.DictType{}).Where("name = ? AND delete_time = ?", name, 0)
	if excludeID > 0 {
		query = query.Where("id <> ?", excludeID)
	}
	var count int64
	err := query.Count(&count).Error
	return count, err
}

func (repo dictRepository) CountDictTypesByCode(ctx context.Context, code string, excludeID uint64) (int64, error) {
	query := repo.db.WithContext(ctx).Model(&makeadmin.DictType{}).Where("code = ? AND delete_time = ?", code, 0)
	if excludeID > 0 {
		query = query.Where("id <> ?", excludeID)
	}
	var count int64
	err := query.Count(&count).Error
	return count, err
}

func (repo dictRepository) CreateDictType(ctx context.Context, dictType makeadmin.DictType) error {
	return repo.db.WithContext(ctx).Create(&dictType).Error
}

func (repo dictRepository) UpdateDictType(ctx context.Context, dictType makeadmin.DictType) error {
	return repo.db.WithContext(ctx).
		Model(&makeadmin.DictType{}).
		Where("id = ? AND delete_time = ?", dictType.ID, 0).
		Updates(map[string]interface{}{
			"code":        dictType.Code,
			"name":        dictType.Name,
			"remark":      dictType.Remark,
			"status":      dictType.Status,
			"update_time": time.Now().Unix(),
		}).
		Error
}

func (repo dictRepository) DeleteDictTypes(ctx context.Context, ids []uint64) error {
	now := time.Now().Unix()
	return repo.db.WithContext(ctx).
		Model(&makeadmin.DictType{}).
		Where("id IN ? AND delete_time = ?", ids, 0).
		Updates(map[string]interface{}{
			"delete_time": now,
			"update_time": now,
		}).
		Error
}

func (repo dictRepository) ListAllDictItemsByTypeID(ctx context.Context, typeID uint64, filter DictItemFilter) ([]makeadmin.DictItem, error) {
	var items []makeadmin.DictItem
	err := repo.dictItemQuery(ctx, typeID, filter).
		Order("sort DESC, id ASC").
		Find(&items).
		Error
	return items, err
}

func (repo dictRepository) ListDictItemsByTypeID(ctx context.Context, typeID uint64, filter DictItemFilter, limit int, offset int) ([]makeadmin.DictItem, int64, error) {
	query := repo.dictItemQuery(ctx, typeID, filter)
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	var items []makeadmin.DictItem
	err := query.
		Limit(limit).
		Offset(offset).
		Order("sort DESC, id ASC").
		Find(&items).
		Error
	return items, count, err
}

func (repo dictRepository) FindDictItemByID(ctx context.Context, id uint64) (item makeadmin.DictItem, err error) {
	err = repo.db.WithContext(ctx).
		Where("id = ? AND delete_time = ?", id, 0).
		Limit(1).
		First(&item).
		Error
	return
}

func (repo dictRepository) CountDictItemsByValue(ctx context.Context, typeID uint64, value string, excludeID uint64) (int64, error) {
	query := repo.db.WithContext(ctx).
		Model(&makeadmin.DictItem{}).
		Where("type_id = ? AND item_value = ? AND delete_time = ?", typeID, value, 0)
	if excludeID > 0 {
		query = query.Where("id <> ?", excludeID)
	}
	var count int64
	err := query.Count(&count).Error
	return count, err
}

func (repo dictRepository) CreateDictItem(ctx context.Context, item makeadmin.DictItem) error {
	return repo.db.WithContext(ctx).Create(&item).Error
}

func (repo dictRepository) UpdateDictItem(ctx context.Context, item makeadmin.DictItem) error {
	return repo.db.WithContext(ctx).
		Model(&makeadmin.DictItem{}).
		Where("id = ? AND delete_time = ?", item.ID, 0).
		Updates(map[string]interface{}{
			"type_id":     item.TypeID,
			"item_label":  item.ItemLabel,
			"item_value":  item.ItemValue,
			"remark":      item.Remark,
			"sort":        item.Sort,
			"status":      item.Status,
			"update_time": time.Now().Unix(),
		}).
		Error
}

func (repo dictRepository) DeleteDictItems(ctx context.Context, ids []uint64) error {
	now := time.Now().Unix()
	return repo.db.WithContext(ctx).
		Model(&makeadmin.DictItem{}).
		Where("id IN ? AND delete_time = ?", ids, 0).
		Updates(map[string]interface{}{
			"delete_time": now,
			"update_time": now,
		}).
		Error
}

func (repo dictRepository) dictTypeQuery(ctx context.Context, filter DictTypeFilter) *gorm.DB {
	query := repo.db.WithContext(ctx).Model(&makeadmin.DictType{}).Where("delete_time = ?", 0)
	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	if filter.Code != "" {
		query = query.Where("code LIKE ?", "%"+filter.Code+"%")
	}
	if filter.StatusSet {
		query = query.Where("status = ?", filter.Status)
	}
	return query
}

func (repo dictRepository) dictItemQuery(ctx context.Context, typeID uint64, filter DictItemFilter) *gorm.DB {
	query := repo.db.WithContext(ctx).Model(&makeadmin.DictItem{}).Where("type_id = ? AND delete_time = ?", typeID, 0)
	if filter.Name != "" {
		query = query.Where("item_label LIKE ?", "%"+filter.Name+"%")
	}
	if filter.Value != "" {
		query = query.Where("item_value LIKE ?", "%"+filter.Value+"%")
	}
	if filter.StatusSet {
		query = query.Where("status = ?", filter.Status)
	}
	return query
}
