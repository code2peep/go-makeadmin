package service

import (
	"context"
	"errors"
	"testing"

	"gorm.io/gorm"

	"go-makeadmin/makeadmin/repository"
	"go-makeadmin/model/makeadmin"
)

type fakeFileRepository struct {
	categories []makeadmin.FileCategory
	files      []makeadmin.File

	fileCount  int64
	childCount int64

	createdCategory makeadmin.FileCategory
	createdFile     makeadmin.File
	movedCategoryID uint64
	deletedCategory uint64
}

func (repo *fakeFileRepository) ListFileCategories(ctx context.Context, filter repository.FileCategoryFilter) ([]makeadmin.FileCategory, error) {
	return repo.categories, nil
}

func (repo *fakeFileRepository) FindFileCategoryByID(ctx context.Context, tenantID uint64, id uint64) (makeadmin.FileCategory, error) {
	for _, category := range repo.categories {
		if category.ID == id && category.TenantID == tenantID && category.DeleteTime == 0 {
			return category, nil
		}
	}
	return makeadmin.FileCategory{}, gorm.ErrRecordNotFound
}

func (repo *fakeFileRepository) CountFilesByCategoryID(ctx context.Context, tenantID uint64, categoryID uint64) (int64, error) {
	return repo.fileCount, nil
}

func (repo *fakeFileRepository) CountChildCategories(ctx context.Context, tenantID uint64, parentID uint64) (int64, error) {
	return repo.childCount, nil
}

func (repo *fakeFileRepository) CreateFileCategory(ctx context.Context, category makeadmin.FileCategory) error {
	repo.createdCategory = category
	return nil
}

func (repo *fakeFileRepository) UpdateFileCategoryName(ctx context.Context, tenantID uint64, id uint64, name string) error {
	return nil
}

func (repo *fakeFileRepository) DeleteFileCategory(ctx context.Context, tenantID uint64, id uint64) error {
	repo.deletedCategory = id
	return nil
}

func (repo *fakeFileRepository) ListFiles(ctx context.Context, filter repository.FileFilter, limit int, offset int) ([]makeadmin.File, int64, error) {
	return repo.files, int64(len(repo.files)), nil
}

func (repo *fakeFileRepository) FindFilesByIDs(ctx context.Context, tenantID uint64, ids []uint64) ([]makeadmin.File, error) {
	idSet := make(map[uint64]struct{}, len(ids))
	for _, id := range ids {
		idSet[id] = struct{}{}
	}
	result := make([]makeadmin.File, 0)
	for _, file := range repo.files {
		if _, ok := idSet[file.ID]; ok && file.TenantID == tenantID && file.DeleteTime == 0 {
			result = append(result, file)
		}
	}
	return result, nil
}

func (repo *fakeFileRepository) CreateFile(ctx context.Context, file makeadmin.File) (makeadmin.File, error) {
	if file.ID == 0 {
		file.ID = 9
	}
	repo.createdFile = file
	return file, nil
}

func (repo *fakeFileRepository) UpdateFileName(ctx context.Context, tenantID uint64, id uint64, name string) error {
	return nil
}

func (repo *fakeFileRepository) MoveFiles(ctx context.Context, tenantID uint64, ids []uint64, categoryID uint64) error {
	repo.movedCategoryID = categoryID
	return nil
}

func (repo *fakeFileRepository) DeleteFiles(ctx context.Context, tenantID uint64, ids []uint64) error {
	return nil
}

func TestAddFileCreatesMakeadminFile(t *testing.T) {
	repo := &fakeFileRepository{
		categories: []makeadmin.FileCategory{{ID: 1, TenantID: makeadmin.GlobalTenantID, FileType: makeadmin.FileTypeImage}},
	}
	srv := NewFileService(repo)

	item, err := srv.AddFile(context.Background(), FileInput{
		TenantID:   makeadmin.GlobalTenantID,
		CategoryID: 1,
		AdminID:    2,
		FileType:   makeadmin.FileTypeImage,
		Name:       "avatar.png",
		URI:        "image/20260601/avatar.png",
		Ext:        "png",
		Size:       1024,
	})
	if err != nil {
		t.Fatalf("AddFile() error = %v", err)
	}
	if item.ID != 9 || repo.createdFile.StorageDriver != StorageAliasLocal || repo.createdFile.OwnerAdminID != 2 {
		t.Fatalf("AddFile() item=%#v created=%#v", item, repo.createdFile)
	}
}

func TestAddFileRejectsMissingCategory(t *testing.T) {
	srv := NewFileService(&fakeFileRepository{})

	_, err := srv.AddFile(context.Background(), FileInput{
		TenantID:   makeadmin.GlobalTenantID,
		CategoryID: 404,
		FileType:   makeadmin.FileTypeImage,
		Name:       "lost.png",
		URI:        "image/lost.png",
	})
	if !errors.Is(err, ErrFileCategoryNotFound) {
		t.Fatalf("AddFile() error = %v, want ErrFileCategoryNotFound", err)
	}
}

func TestMoveFilesValidatesTargetCategory(t *testing.T) {
	repo := &fakeFileRepository{
		files: []makeadmin.File{{ID: 1, TenantID: makeadmin.GlobalTenantID}},
		categories: []makeadmin.FileCategory{
			{ID: 2, TenantID: makeadmin.GlobalTenantID, FileType: makeadmin.FileTypeImage},
		},
	}
	srv := NewFileService(repo)

	err := srv.MoveFiles(context.Background(), makeadmin.GlobalTenantID, []uint64{1}, 2)
	if err != nil {
		t.Fatalf("MoveFiles() error = %v", err)
	}
	if repo.movedCategoryID != 2 {
		t.Fatalf("MoveFiles() category = %d, want 2", repo.movedCategoryID)
	}
}

func TestDeleteCategoryRejectsChildrenAndFiles(t *testing.T) {
	repo := &fakeFileRepository{
		categories: []makeadmin.FileCategory{{ID: 1, TenantID: makeadmin.GlobalTenantID}},
		childCount: 1,
	}
	srv := NewFileService(repo)

	err := srv.DeleteCategory(context.Background(), makeadmin.GlobalTenantID, 1)
	if !errors.Is(err, ErrFileCategoryHasChildren) {
		t.Fatalf("DeleteCategory() child error = %v", err)
	}

	repo.childCount = 0
	repo.fileCount = 1
	err = srv.DeleteCategory(context.Background(), makeadmin.GlobalTenantID, 1)
	if !errors.Is(err, ErrFileCategoryInUse) {
		t.Fatalf("DeleteCategory() file error = %v", err)
	}
}

func TestAddCategoryCreatesGeneratedCode(t *testing.T) {
	repo := &fakeFileRepository{}
	srv := NewFileService(repo)

	err := srv.AddCategory(context.Background(), FileCategoryInput{
		TenantID: makeadmin.GlobalTenantID,
		FileType: makeadmin.FileTypeImage,
		Name:     "素材",
	})
	if err != nil {
		t.Fatalf("AddCategory() error = %v", err)
	}
	if repo.createdCategory.Code == "" || repo.createdCategory.Status != makeadmin.StatusEnabled {
		t.Fatalf("AddCategory() created = %#v", repo.createdCategory)
	}
}
