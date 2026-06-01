package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"go-makeadmin/makeadmin/repository"
	"go-makeadmin/model/makeadmin"
)

const FileTypeOther = "file"

var (
	ErrFileNotFound            = errors.New("makeadmin file not found")
	ErrFileCategoryNotFound    = errors.New("makeadmin file category not found")
	ErrFileCategoryInUse       = errors.New("makeadmin file category in use")
	ErrFileCategoryHasChildren = errors.New("makeadmin file category has children")
	ErrUnsupportedFileType     = errors.New("makeadmin unsupported file type")
)

type FileFilter struct {
	CategoryID uint64
	FileType   string
	Name       string
}

type FilePageInput struct {
	PageNo   int
	PageSize int
}

type FileItem struct {
	ID         uint64
	CategoryID uint64
	AdminID    uint64
	FileType   string
	Name       string
	URI        string
	Ext        string
	Size       int64
	CreateTime int64
	UpdateTime int64
}

type FileCategoryItem struct {
	ID         uint64
	ParentID   uint64
	Name       string
	FileType   string
	CreateTime int64
	UpdateTime int64
}

type FilePage struct {
	Items []FileItem
	Count int64
}

type FileInput struct {
	TenantID   uint64
	CategoryID uint64
	AdminID    uint64
	FileType   string
	Name       string
	URI        string
	URL        string
	Ext        string
	Size       int64
}

type FileCategoryInput struct {
	TenantID uint64
	ParentID uint64
	FileType string
	Name     string
}

type FileService interface {
	ListFiles(ctx context.Context, tenantID uint64, filter FileFilter, page FilePageInput) (FilePage, error)
	RenameFile(ctx context.Context, tenantID uint64, id uint64, name string) error
	MoveFiles(ctx context.Context, tenantID uint64, ids []uint64, categoryID uint64) error
	AddFile(ctx context.Context, input FileInput) (FileItem, error)
	DeleteFiles(ctx context.Context, tenantID uint64, ids []uint64) error
	ListCategories(ctx context.Context, tenantID uint64, fileType string, name string) ([]FileCategoryItem, error)
	AddCategory(ctx context.Context, input FileCategoryInput) error
	RenameCategory(ctx context.Context, tenantID uint64, id uint64, name string) error
	DeleteCategory(ctx context.Context, tenantID uint64, id uint64) error
}

type fileService struct {
	repo repository.FileRepository
}

func NewFileService(repo repository.FileRepository) FileService {
	return fileService{repo: repo}
}

func (srv fileService) ListFiles(ctx context.Context, tenantID uint64, filter FileFilter, page FilePageInput) (FilePage, error) {
	if filter.FileType != "" && !isSupportedFileType(filter.FileType) {
		return FilePage{}, ErrUnsupportedFileType
	}
	files, count, err := srv.repo.ListFiles(ctx, repository.FileFilter{
		TenantID:   tenantID,
		CategoryID: filter.CategoryID,
		FileType:   filter.FileType,
		Name:       filter.Name,
	}, pageLimitForFile(page), pageOffsetForFile(page))
	if err != nil {
		return FilePage{}, err
	}
	return FilePage{Items: fileItemsFromModels(files), Count: count}, nil
}

func (srv fileService) RenameFile(ctx context.Context, tenantID uint64, id uint64, name string) error {
	files, err := srv.repo.FindFilesByIDs(ctx, tenantID, []uint64{id})
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return ErrFileNotFound
	}
	return srv.repo.UpdateFileName(ctx, tenantID, id, name)
}

func (srv fileService) MoveFiles(ctx context.Context, tenantID uint64, ids []uint64, categoryID uint64) error {
	files, err := srv.repo.FindFilesByIDs(ctx, tenantID, ids)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return ErrFileNotFound
	}
	if categoryID > 0 {
		if _, err := srv.repo.FindFileCategoryByID(ctx, tenantID, categoryID); err != nil {
			return mapFileRecordError(err, ErrFileCategoryNotFound)
		}
	}
	return srv.repo.MoveFiles(ctx, tenantID, ids, categoryID)
}

func (srv fileService) AddFile(ctx context.Context, input FileInput) (FileItem, error) {
	if !isSupportedFileType(input.FileType) {
		return FileItem{}, ErrUnsupportedFileType
	}
	if input.CategoryID > 0 {
		if _, err := srv.repo.FindFileCategoryByID(ctx, input.TenantID, input.CategoryID); err != nil {
			return FileItem{}, mapFileRecordError(err, ErrFileCategoryNotFound)
		}
	}
	file, err := srv.repo.CreateFile(ctx, makeadmin.File{
		TenantID:      input.TenantID,
		CategoryID:    input.CategoryID,
		OwnerAdminID:  input.AdminID,
		FileType:      input.FileType,
		StorageDriver: StorageAliasLocal,
		OriginalName:  input.Name,
		FileName:      input.Name,
		URI:           strings.TrimPrefix(input.URI, "/"),
		URL:           input.URL,
		Ext:           input.Ext,
		Size:          input.Size,
		Status:        makeadmin.StatusEnabled,
	})
	if err != nil {
		return FileItem{}, err
	}
	return fileItemFromModel(file), nil
}

func (srv fileService) DeleteFiles(ctx context.Context, tenantID uint64, ids []uint64) error {
	files, err := srv.repo.FindFilesByIDs(ctx, tenantID, ids)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return ErrFileNotFound
	}
	return srv.repo.DeleteFiles(ctx, tenantID, ids)
}

func (srv fileService) ListCategories(ctx context.Context, tenantID uint64, fileType string, name string) ([]FileCategoryItem, error) {
	if fileType != "" && !isSupportedFileType(fileType) {
		return nil, ErrUnsupportedFileType
	}
	categories, err := srv.repo.ListFileCategories(ctx, repository.FileCategoryFilter{
		TenantID: tenantID,
		FileType: fileType,
		Name:     name,
	})
	if err != nil {
		return nil, err
	}
	return fileCategoryItemsFromModels(categories), nil
}

func (srv fileService) AddCategory(ctx context.Context, input FileCategoryInput) error {
	if !isSupportedFileType(input.FileType) {
		return ErrUnsupportedFileType
	}
	if input.ParentID > 0 {
		if _, err := srv.repo.FindFileCategoryByID(ctx, input.TenantID, input.ParentID); err != nil {
			return mapFileRecordError(err, ErrFileCategoryNotFound)
		}
	}
	return srv.repo.CreateFileCategory(ctx, makeadmin.FileCategory{
		TenantID: input.TenantID,
		ParentID: input.ParentID,
		Code:     newCategoryCode(input.FileType),
		Name:     input.Name,
		FileType: input.FileType,
		Status:   makeadmin.StatusEnabled,
	})
}

func (srv fileService) RenameCategory(ctx context.Context, tenantID uint64, id uint64, name string) error {
	if _, err := srv.repo.FindFileCategoryByID(ctx, tenantID, id); err != nil {
		return mapFileRecordError(err, ErrFileCategoryNotFound)
	}
	return srv.repo.UpdateFileCategoryName(ctx, tenantID, id, name)
}

func (srv fileService) DeleteCategory(ctx context.Context, tenantID uint64, id uint64) error {
	if _, err := srv.repo.FindFileCategoryByID(ctx, tenantID, id); err != nil {
		return mapFileRecordError(err, ErrFileCategoryNotFound)
	}
	childCount, err := srv.repo.CountChildCategories(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if childCount > 0 {
		return ErrFileCategoryHasChildren
	}
	fileCount, err := srv.repo.CountFilesByCategoryID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if fileCount > 0 {
		return ErrFileCategoryInUse
	}
	return srv.repo.DeleteFileCategory(ctx, tenantID, id)
}

func fileItemsFromModels(models []makeadmin.File) []FileItem {
	result := make([]FileItem, 0, len(models))
	for _, model := range models {
		result = append(result, fileItemFromModel(model))
	}
	return result
}

func fileItemFromModel(model makeadmin.File) FileItem {
	name := model.OriginalName
	if name == "" {
		name = model.FileName
	}
	return FileItem{
		ID:         model.ID,
		CategoryID: model.CategoryID,
		AdminID:    model.OwnerAdminID,
		FileType:   model.FileType,
		Name:       name,
		URI:        model.URI,
		Ext:        model.Ext,
		Size:       model.Size,
		CreateTime: model.CreateTime,
		UpdateTime: model.UpdateTime,
	}
}

func fileCategoryItemsFromModels(models []makeadmin.FileCategory) []FileCategoryItem {
	result := make([]FileCategoryItem, 0, len(models))
	for _, model := range models {
		result = append(result, FileCategoryItem{
			ID:         model.ID,
			ParentID:   model.ParentID,
			Name:       model.Name,
			FileType:   model.FileType,
			CreateTime: model.CreateTime,
			UpdateTime: model.UpdateTime,
		})
	}
	return result
}

func mapFileRecordError(err error, notFound error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return notFound
	}
	return err
}

func pageLimitForFile(page FilePageInput) int {
	if page.PageSize <= 0 {
		return 20
	}
	return page.PageSize
}

func pageOffsetForFile(page FilePageInput) int {
	pageNo := page.PageNo
	if pageNo <= 0 {
		pageNo = 1
	}
	return pageLimitForFile(page) * (pageNo - 1)
}

func isSupportedFileType(fileType string) bool {
	return fileType == makeadmin.FileTypeImage || fileType == makeadmin.FileTypeVideo || fileType == FileTypeOther
}

func newCategoryCode(fileType string) string {
	return fmt.Sprintf("%s_%d", fileType, time.Now().UnixNano())
}
