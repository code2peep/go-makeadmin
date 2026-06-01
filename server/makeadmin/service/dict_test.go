package service

import (
	"context"
	"errors"
	"testing"

	"gorm.io/gorm"

	"go-makeadmin/makeadmin/repository"
	"go-makeadmin/model/makeadmin"
)

type fakeDictRepository struct {
	types       []makeadmin.DictType
	items       []makeadmin.DictItem
	createdType makeadmin.DictType
	updatedType makeadmin.DictType
	createdItem makeadmin.DictItem
	updatedItem makeadmin.DictItem
}

func (repo *fakeDictRepository) ListAllDictTypes(ctx context.Context) ([]makeadmin.DictType, error) {
	return repo.types, nil
}

func (repo *fakeDictRepository) ListDictTypes(ctx context.Context, filter repository.DictTypeFilter, limit int, offset int) ([]makeadmin.DictType, int64, error) {
	return repo.types, int64(len(repo.types)), nil
}

func (repo *fakeDictRepository) FindDictTypeByID(ctx context.Context, id uint64) (makeadmin.DictType, error) {
	for _, dictType := range repo.types {
		if dictType.ID == id && dictType.DeleteTime == 0 {
			return dictType, nil
		}
	}
	return makeadmin.DictType{}, gorm.ErrRecordNotFound
}

func (repo *fakeDictRepository) FindDictTypeByCode(ctx context.Context, code string) (makeadmin.DictType, error) {
	for _, dictType := range repo.types {
		if dictType.Code == code && dictType.DeleteTime == 0 {
			return dictType, nil
		}
	}
	return makeadmin.DictType{}, gorm.ErrRecordNotFound
}

func (repo *fakeDictRepository) CountDictTypesByName(ctx context.Context, name string, excludeID uint64) (int64, error) {
	var count int64
	for _, dictType := range repo.types {
		if dictType.Name == name && dictType.ID != excludeID && dictType.DeleteTime == 0 {
			count++
		}
	}
	return count, nil
}

func (repo *fakeDictRepository) CountDictTypesByCode(ctx context.Context, code string, excludeID uint64) (int64, error) {
	var count int64
	for _, dictType := range repo.types {
		if dictType.Code == code && dictType.ID != excludeID && dictType.DeleteTime == 0 {
			count++
		}
	}
	return count, nil
}

func (repo *fakeDictRepository) CreateDictType(ctx context.Context, dictType makeadmin.DictType) error {
	repo.createdType = dictType
	return nil
}

func (repo *fakeDictRepository) UpdateDictType(ctx context.Context, dictType makeadmin.DictType) error {
	repo.updatedType = dictType
	return nil
}

func (repo *fakeDictRepository) DeleteDictTypes(ctx context.Context, ids []uint64) error {
	return nil
}

func (repo *fakeDictRepository) ListAllDictItemsByTypeID(ctx context.Context, typeID uint64, filter repository.DictItemFilter) ([]makeadmin.DictItem, error) {
	return repo.itemsByType(typeID), nil
}

func (repo *fakeDictRepository) ListDictItemsByTypeID(ctx context.Context, typeID uint64, filter repository.DictItemFilter, limit int, offset int) ([]makeadmin.DictItem, int64, error) {
	items := repo.itemsByType(typeID)
	return items, int64(len(items)), nil
}

func (repo *fakeDictRepository) FindDictItemByID(ctx context.Context, id uint64) (makeadmin.DictItem, error) {
	for _, item := range repo.items {
		if item.ID == id && item.DeleteTime == 0 {
			return item, nil
		}
	}
	return makeadmin.DictItem{}, gorm.ErrRecordNotFound
}

func (repo *fakeDictRepository) CountDictItemsByValue(ctx context.Context, typeID uint64, value string, excludeID uint64) (int64, error) {
	var count int64
	for _, item := range repo.items {
		if item.TypeID == typeID && item.ItemValue == value && item.ID != excludeID && item.DeleteTime == 0 {
			count++
		}
	}
	return count, nil
}

func (repo *fakeDictRepository) CreateDictItem(ctx context.Context, item makeadmin.DictItem) error {
	repo.createdItem = item
	return nil
}

func (repo *fakeDictRepository) UpdateDictItem(ctx context.Context, item makeadmin.DictItem) error {
	repo.updatedItem = item
	return nil
}

func (repo *fakeDictRepository) DeleteDictItems(ctx context.Context, ids []uint64) error {
	return nil
}

func (repo *fakeDictRepository) itemsByType(typeID uint64) []makeadmin.DictItem {
	result := make([]makeadmin.DictItem, 0)
	for _, item := range repo.items {
		if item.TypeID == typeID && item.DeleteTime == 0 {
			result = append(result, item)
		}
	}
	return result
}

func TestDictTypeAddRejectsDuplicateCode(t *testing.T) {
	srv := NewDictService(&fakeDictRepository{
		types: []makeadmin.DictType{{ID: 1, Code: "status", Name: "状态"}},
	})

	err := srv.AddType(context.Background(), DictTypeInput{
		Name:   "启停状态",
		Code:   "status",
		Status: 1,
	})
	if !errors.Is(err, ErrDictTypeCodeExists) {
		t.Fatalf("AddType() error = %v, want ErrDictTypeCodeExists", err)
	}
}

func TestDictTypeEditUpdatesExisting(t *testing.T) {
	repo := &fakeDictRepository{
		types: []makeadmin.DictType{{ID: 1, Code: "menu_type", Name: "菜单类型", Status: 1}},
	}
	srv := NewDictService(repo)

	err := srv.EditType(context.Background(), DictTypeInput{
		ID:     1,
		Name:   "菜单分类",
		Code:   "menu_kind",
		Remark: "updated",
		Status: 0,
	})
	if err != nil {
		t.Fatalf("EditType() error = %v", err)
	}
	if repo.updatedType.ID != 1 || repo.updatedType.Code != "menu_kind" || repo.updatedType.Status != 0 {
		t.Fatalf("EditType() updated = %#v", repo.updatedType)
	}
}

func TestDictItemsResolveTypeCode(t *testing.T) {
	repo := &fakeDictRepository{
		types: []makeadmin.DictType{{ID: 3, Code: "storage_type", Name: "存储类型"}},
		items: []makeadmin.DictItem{
			{ID: 6, TypeID: 3, ItemLabel: "本地", ItemValue: "local"},
			{ID: 7, TypeID: 3, ItemLabel: "七牛云", ItemValue: "qiniu"},
		},
	}
	srv := NewDictService(repo)

	items, err := srv.AllItems(context.Background(), "storage_type", repository.DictItemFilter{})
	if err != nil {
		t.Fatalf("AllItems() error = %v", err)
	}
	if len(items) != 2 || items[0].Value != "local" {
		t.Fatalf("AllItems() = %#v", items)
	}
}

func TestDictItemAddRejectsDuplicateValue(t *testing.T) {
	repo := &fakeDictRepository{
		types: []makeadmin.DictType{{ID: 3, Code: "storage_type", Name: "存储类型"}},
		items: []makeadmin.DictItem{{ID: 6, TypeID: 3, ItemLabel: "本地", ItemValue: "local"}},
	}
	srv := NewDictService(repo)

	err := srv.AddItem(context.Background(), DictItemInput{
		TypeID: 3,
		Name:   "本地文件",
		Value:  "local",
		Status: 1,
	})
	if !errors.Is(err, ErrDictItemValueExists) {
		t.Fatalf("AddItem() error = %v, want ErrDictItemValueExists", err)
	}
}

func TestDictItemEditUpdatesExisting(t *testing.T) {
	repo := &fakeDictRepository{
		types: []makeadmin.DictType{{ID: 3, Code: "storage_type", Name: "存储类型"}},
		items: []makeadmin.DictItem{{ID: 7, TypeID: 3, ItemLabel: "七牛云", ItemValue: "qiniu", Status: 1}},
	}
	srv := NewDictService(repo)

	err := srv.EditItem(context.Background(), DictItemInput{
		ID:     7,
		TypeID: 3,
		Name:   "七牛云 OSS",
		Value:  "qiniu",
		Remark: "updated",
		Sort:   990,
		Status: 0,
	})
	if err != nil {
		t.Fatalf("EditItem() error = %v", err)
	}
	if repo.updatedItem.ID != 7 || repo.updatedItem.ItemLabel != "七牛云 OSS" || repo.updatedItem.Status != 0 {
		t.Fatalf("EditItem() updated = %#v", repo.updatedItem)
	}
}
