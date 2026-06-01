package service

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"go-makeadmin/makeadmin/repository"
	"go-makeadmin/model/makeadmin"
)

var (
	ErrDictTypeNotFound    = errors.New("makeadmin dict type not found")
	ErrDictTypeNameExists  = errors.New("makeadmin dict type name exists")
	ErrDictTypeCodeExists  = errors.New("makeadmin dict type code exists")
	ErrDictItemNotFound    = errors.New("makeadmin dict item not found")
	ErrDictItemValueExists = errors.New("makeadmin dict item value exists")
)

type DictType struct {
	ID         uint64
	Name       string
	Code       string
	Remark     string
	Status     uint8
	CreateTime int64
	UpdateTime int64
}

type DictItem struct {
	ID         uint64
	TypeID     uint64
	Name       string
	Value      string
	Remark     string
	Sort       uint16
	Status     uint8
	CreateTime int64
	UpdateTime int64
}

type DictTypeInput struct {
	ID     uint64
	Name   string
	Code   string
	Remark string
	Status int8
}

type DictItemInput struct {
	ID     uint64
	TypeID uint64
	Name   string
	Value  string
	Remark string
	Sort   int
	Status int8
}

type DictTypePage struct {
	Items []DictType
	Count int64
}

type DictItemPage struct {
	Items []DictItem
	Count int64
}

type DictService interface {
	AllTypes(ctx context.Context) ([]DictType, error)
	ListTypes(ctx context.Context, filter repository.DictTypeFilter, limit int, offset int) (DictTypePage, error)
	DetailType(ctx context.Context, id uint64) (DictType, error)
	AddType(ctx context.Context, input DictTypeInput) error
	EditType(ctx context.Context, input DictTypeInput) error
	DeleteTypes(ctx context.Context, ids []uint64) error
	AllItems(ctx context.Context, typeCode string, filter repository.DictItemFilter) ([]DictItem, error)
	ListItems(ctx context.Context, typeCode string, filter repository.DictItemFilter, limit int, offset int) (DictItemPage, error)
	DetailItem(ctx context.Context, id uint64) (DictItem, error)
	AddItem(ctx context.Context, input DictItemInput) error
	EditItem(ctx context.Context, input DictItemInput) error
	DeleteItems(ctx context.Context, ids []uint64) error
}

type dictService struct {
	repo repository.DictRepository
}

func NewDictService(repo repository.DictRepository) DictService {
	return dictService{repo: repo}
}

func (srv dictService) AllTypes(ctx context.Context) ([]DictType, error) {
	types, err := srv.repo.ListAllDictTypes(ctx)
	if err != nil {
		return nil, err
	}
	return dictTypesFromModels(types), nil
}

func (srv dictService) ListTypes(ctx context.Context, filter repository.DictTypeFilter, limit int, offset int) (DictTypePage, error) {
	types, count, err := srv.repo.ListDictTypes(ctx, filter, limit, offset)
	if err != nil {
		return DictTypePage{}, err
	}
	return DictTypePage{Items: dictTypesFromModels(types), Count: count}, nil
}

func (srv dictService) DetailType(ctx context.Context, id uint64) (DictType, error) {
	dictType, err := srv.repo.FindDictTypeByID(ctx, id)
	if err != nil {
		return DictType{}, mapDictRecordError(err, ErrDictTypeNotFound)
	}
	return dictTypeFromModel(dictType), nil
}

func (srv dictService) AddType(ctx context.Context, input DictTypeInput) error {
	if err := srv.ensureTypeUnique(ctx, input.Name, input.Code, 0); err != nil {
		return err
	}
	return srv.repo.CreateDictType(ctx, makeadmin.DictType{
		Name:   input.Name,
		Code:   input.Code,
		Remark: input.Remark,
		Status: normalizeStatus(input.Status, 1),
	})
}

func (srv dictService) EditType(ctx context.Context, input DictTypeInput) error {
	current, err := srv.repo.FindDictTypeByID(ctx, input.ID)
	if err != nil {
		return mapDictRecordError(err, ErrDictTypeNotFound)
	}
	if err := srv.ensureTypeUnique(ctx, input.Name, input.Code, input.ID); err != nil {
		return err
	}
	current.Name = input.Name
	current.Code = input.Code
	current.Remark = input.Remark
	current.Status = normalizeStatus(input.Status, current.Status)
	return srv.repo.UpdateDictType(ctx, current)
}

func (srv dictService) DeleteTypes(ctx context.Context, ids []uint64) error {
	return srv.repo.DeleteDictTypes(ctx, ids)
}

func (srv dictService) AllItems(ctx context.Context, typeCode string, filter repository.DictItemFilter) ([]DictItem, error) {
	dictType, err := srv.repo.FindDictTypeByCode(ctx, typeCode)
	if err != nil {
		return nil, mapDictRecordError(err, ErrDictTypeNotFound)
	}
	items, err := srv.repo.ListAllDictItemsByTypeID(ctx, dictType.ID, filter)
	if err != nil {
		return nil, err
	}
	return dictItemsFromModels(items), nil
}

func (srv dictService) ListItems(ctx context.Context, typeCode string, filter repository.DictItemFilter, limit int, offset int) (DictItemPage, error) {
	dictType, err := srv.repo.FindDictTypeByCode(ctx, typeCode)
	if err != nil {
		return DictItemPage{}, mapDictRecordError(err, ErrDictTypeNotFound)
	}
	items, count, err := srv.repo.ListDictItemsByTypeID(ctx, dictType.ID, filter, limit, offset)
	if err != nil {
		return DictItemPage{}, err
	}
	return DictItemPage{Items: dictItemsFromModels(items), Count: count}, nil
}

func (srv dictService) DetailItem(ctx context.Context, id uint64) (DictItem, error) {
	item, err := srv.repo.FindDictItemByID(ctx, id)
	if err != nil {
		return DictItem{}, mapDictRecordError(err, ErrDictItemNotFound)
	}
	return dictItemFromModel(item), nil
}

func (srv dictService) AddItem(ctx context.Context, input DictItemInput) error {
	if _, err := srv.repo.FindDictTypeByID(ctx, input.TypeID); err != nil {
		return mapDictRecordError(err, ErrDictTypeNotFound)
	}
	if err := srv.ensureItemUnique(ctx, input.TypeID, input.Value, 0); err != nil {
		return err
	}
	return srv.repo.CreateDictItem(ctx, makeadmin.DictItem{
		TypeID:    input.TypeID,
		ItemLabel: input.Name,
		ItemValue: input.Value,
		Remark:    input.Remark,
		Sort:      normalizeSort(input.Sort),
		Status:    normalizeStatus(input.Status, 1),
	})
}

func (srv dictService) EditItem(ctx context.Context, input DictItemInput) error {
	current, err := srv.repo.FindDictItemByID(ctx, input.ID)
	if err != nil {
		return mapDictRecordError(err, ErrDictItemNotFound)
	}
	if _, err := srv.repo.FindDictTypeByID(ctx, input.TypeID); err != nil {
		return mapDictRecordError(err, ErrDictTypeNotFound)
	}
	if err := srv.ensureItemUnique(ctx, input.TypeID, input.Value, input.ID); err != nil {
		return err
	}
	current.TypeID = input.TypeID
	current.ItemLabel = input.Name
	current.ItemValue = input.Value
	current.Remark = input.Remark
	current.Sort = normalizeSort(input.Sort)
	current.Status = normalizeStatus(input.Status, current.Status)
	return srv.repo.UpdateDictItem(ctx, current)
}

func (srv dictService) DeleteItems(ctx context.Context, ids []uint64) error {
	return srv.repo.DeleteDictItems(ctx, ids)
}

func (srv dictService) ensureTypeUnique(ctx context.Context, name string, code string, excludeID uint64) error {
	count, err := srv.repo.CountDictTypesByName(ctx, name, excludeID)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrDictTypeNameExists
	}
	count, err = srv.repo.CountDictTypesByCode(ctx, code, excludeID)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrDictTypeCodeExists
	}
	return nil
}

func (srv dictService) ensureItemUnique(ctx context.Context, typeID uint64, value string, excludeID uint64) error {
	count, err := srv.repo.CountDictItemsByValue(ctx, typeID, value, excludeID)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrDictItemValueExists
	}
	return nil
}

func dictTypesFromModels(models []makeadmin.DictType) []DictType {
	result := make([]DictType, 0, len(models))
	for _, item := range models {
		result = append(result, dictTypeFromModel(item))
	}
	return result
}

func dictTypeFromModel(model makeadmin.DictType) DictType {
	return DictType{
		ID:         model.ID,
		Name:       model.Name,
		Code:       model.Code,
		Remark:     model.Remark,
		Status:     model.Status,
		CreateTime: model.CreateTime,
		UpdateTime: model.UpdateTime,
	}
}

func dictItemsFromModels(models []makeadmin.DictItem) []DictItem {
	result := make([]DictItem, 0, len(models))
	for _, item := range models {
		result = append(result, dictItemFromModel(item))
	}
	return result
}

func dictItemFromModel(model makeadmin.DictItem) DictItem {
	return DictItem{
		ID:         model.ID,
		TypeID:     model.TypeID,
		Name:       model.ItemLabel,
		Value:      model.ItemValue,
		Remark:     model.Remark,
		Sort:       model.Sort,
		Status:     model.Status,
		CreateTime: model.CreateTime,
		UpdateTime: model.UpdateTime,
	}
}

func mapDictRecordError(err error, notFound error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return notFound
	}
	return err
}

func normalizeStatus(status int8, fallback uint8) uint8 {
	if status == 0 || status == 1 {
		return uint8(status)
	}
	return fallback
}

func normalizeSort(sort int) uint16 {
	if sort < 0 {
		return 0
	}
	if sort > 65535 {
		return 65535
	}
	return uint16(sort)
}
