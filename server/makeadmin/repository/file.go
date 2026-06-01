package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"go-makeadmin/model/makeadmin"
)

type FileCategoryFilter struct {
	TenantID uint64
	FileType string
	Name     string
}

type FileFilter struct {
	TenantID   uint64
	CategoryID uint64
	FileType   string
	Name       string
}

type FileRepository interface {
	ListFileCategories(ctx context.Context, filter FileCategoryFilter) ([]makeadmin.FileCategory, error)
	FindFileCategoryByID(ctx context.Context, tenantID uint64, id uint64) (makeadmin.FileCategory, error)
	CountFilesByCategoryID(ctx context.Context, tenantID uint64, categoryID uint64) (int64, error)
	CountChildCategories(ctx context.Context, tenantID uint64, parentID uint64) (int64, error)
	CreateFileCategory(ctx context.Context, category makeadmin.FileCategory) error
	UpdateFileCategoryName(ctx context.Context, tenantID uint64, id uint64, name string) error
	DeleteFileCategory(ctx context.Context, tenantID uint64, id uint64) error
	ListFiles(ctx context.Context, filter FileFilter, limit int, offset int) ([]makeadmin.File, int64, error)
	FindFilesByIDs(ctx context.Context, tenantID uint64, ids []uint64) ([]makeadmin.File, error)
	CreateFile(ctx context.Context, file makeadmin.File) (makeadmin.File, error)
	UpdateFileName(ctx context.Context, tenantID uint64, id uint64, name string) error
	MoveFiles(ctx context.Context, tenantID uint64, ids []uint64, categoryID uint64) error
	DeleteFiles(ctx context.Context, tenantID uint64, ids []uint64) error
}

type fileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) FileRepository {
	return &fileRepository{db: db}
}

func (repo fileRepository) ListFileCategories(ctx context.Context, filter FileCategoryFilter) ([]makeadmin.FileCategory, error) {
	query := repo.db.WithContext(ctx).
		Where("tenant_id = ? AND delete_time = ?", filter.TenantID, 0)
	if filter.FileType != "" {
		query = query.Where("file_type = ?", filter.FileType)
	}
	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	var categories []makeadmin.FileCategory
	err := query.Order("sort DESC, id ASC").Find(&categories).Error
	return categories, err
}

func (repo fileRepository) FindFileCategoryByID(ctx context.Context, tenantID uint64, id uint64) (category makeadmin.FileCategory, err error) {
	err = repo.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ? AND delete_time = ?", tenantID, id, 0).
		Limit(1).
		First(&category).
		Error
	return
}

func (repo fileRepository) CountFilesByCategoryID(ctx context.Context, tenantID uint64, categoryID uint64) (int64, error) {
	var count int64
	err := repo.db.WithContext(ctx).
		Model(&makeadmin.File{}).
		Where("tenant_id = ? AND category_id = ? AND delete_time = ?", tenantID, categoryID, 0).
		Count(&count).
		Error
	return count, err
}

func (repo fileRepository) CountChildCategories(ctx context.Context, tenantID uint64, parentID uint64) (int64, error) {
	var count int64
	err := repo.db.WithContext(ctx).
		Model(&makeadmin.FileCategory{}).
		Where("tenant_id = ? AND parent_id = ? AND delete_time = ?", tenantID, parentID, 0).
		Count(&count).
		Error
	return count, err
}

func (repo fileRepository) CreateFileCategory(ctx context.Context, category makeadmin.FileCategory) error {
	return repo.db.WithContext(ctx).Create(&category).Error
}

func (repo fileRepository) UpdateFileCategoryName(ctx context.Context, tenantID uint64, id uint64, name string) error {
	return repo.db.WithContext(ctx).
		Model(&makeadmin.FileCategory{}).
		Where("tenant_id = ? AND id = ? AND delete_time = ?", tenantID, id, 0).
		Updates(map[string]interface{}{
			"name":        name,
			"update_time": time.Now().Unix(),
		}).
		Error
}

func (repo fileRepository) DeleteFileCategory(ctx context.Context, tenantID uint64, id uint64) error {
	now := time.Now().Unix()
	return repo.db.WithContext(ctx).
		Model(&makeadmin.FileCategory{}).
		Where("tenant_id = ? AND id = ? AND delete_time = ?", tenantID, id, 0).
		Updates(map[string]interface{}{
			"delete_time": now,
			"update_time": now,
		}).
		Error
}

func (repo fileRepository) ListFiles(ctx context.Context, filter FileFilter, limit int, offset int) ([]makeadmin.File, int64, error) {
	query := repo.fileQuery(ctx, filter)
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	var files []makeadmin.File
	err := query.
		Limit(limit).
		Offset(offset).
		Order("id DESC").
		Find(&files).
		Error
	return files, count, err
}

func (repo fileRepository) FindFilesByIDs(ctx context.Context, tenantID uint64, ids []uint64) ([]makeadmin.File, error) {
	var files []makeadmin.File
	err := repo.db.WithContext(ctx).
		Where("tenant_id = ? AND id IN ? AND delete_time = ?", tenantID, ids, 0).
		Find(&files).
		Error
	return files, err
}

func (repo fileRepository) CreateFile(ctx context.Context, file makeadmin.File) (makeadmin.File, error) {
	err := repo.db.WithContext(ctx).Create(&file).Error
	return file, err
}

func (repo fileRepository) UpdateFileName(ctx context.Context, tenantID uint64, id uint64, name string) error {
	return repo.db.WithContext(ctx).
		Model(&makeadmin.File{}).
		Where("tenant_id = ? AND id = ? AND delete_time = ?", tenantID, id, 0).
		Updates(map[string]interface{}{
			"original_name": name,
			"update_time":   time.Now().Unix(),
		}).
		Error
}

func (repo fileRepository) MoveFiles(ctx context.Context, tenantID uint64, ids []uint64, categoryID uint64) error {
	return repo.db.WithContext(ctx).
		Model(&makeadmin.File{}).
		Where("tenant_id = ? AND id IN ? AND delete_time = ?", tenantID, ids, 0).
		Updates(map[string]interface{}{
			"category_id": categoryID,
			"update_time": time.Now().Unix(),
		}).
		Error
}

func (repo fileRepository) DeleteFiles(ctx context.Context, tenantID uint64, ids []uint64) error {
	now := time.Now().Unix()
	return repo.db.WithContext(ctx).
		Model(&makeadmin.File{}).
		Where("tenant_id = ? AND id IN ? AND delete_time = ?", tenantID, ids, 0).
		Updates(map[string]interface{}{
			"delete_time": now,
			"update_time": now,
		}).
		Error
}

func (repo fileRepository) fileQuery(ctx context.Context, filter FileFilter) *gorm.DB {
	query := repo.db.WithContext(ctx).
		Model(&makeadmin.File{}).
		Where("tenant_id = ? AND delete_time = ?", filter.TenantID, 0)
	if filter.CategoryID > 0 {
		query = query.Where("category_id = ?", filter.CategoryID)
	}
	if filter.FileType != "" {
		query = query.Where("file_type = ?", filter.FileType)
	}
	if filter.Name != "" {
		query = query.Where("original_name LIKE ?", "%"+filter.Name+"%")
	}
	return query
}
