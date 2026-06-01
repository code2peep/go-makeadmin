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

var (
	ErrOrgUnitNotFound         = errors.New("makeadmin org unit not found")
	ErrParentOrgUnitNotFound   = errors.New("makeadmin parent org unit not found")
	ErrRootOrgUnitExists       = errors.New("makeadmin root org unit exists")
	ErrRootOrgUnitParentLocked = errors.New("makeadmin root org unit parent locked")
	ErrRootOrgUnitDeleteLocked = errors.New("makeadmin root org unit delete locked")
	ErrOrgUnitSelfParent       = errors.New("makeadmin org unit self parent")
	ErrOrgUnitHasChildren      = errors.New("makeadmin org unit has children")
	ErrOrgUnitInUse            = errors.New("makeadmin org unit in use")
	ErrPositionNotFound        = errors.New("makeadmin position not found")
	ErrPositionCodeExists      = errors.New("makeadmin position code exists")
	ErrPositionNameExists      = errors.New("makeadmin position name exists")
	ErrPositionInUse           = errors.New("makeadmin position in use")
)

type OrgUnitInput struct {
	ID       uint64
	TenantID uint64
	ParentID uint64
	Name     string
	IsStop   uint8
	Sort     int
}

type OrgUnitItem struct {
	ID         uint64
	ParentID   uint64
	Name       string
	Sort       uint16
	IsStop     uint8
	CreateTime int64
	UpdateTime int64
}

type PositionInput struct {
	ID       uint64
	TenantID uint64
	Code     string
	Name     string
	Remark   string
	IsStop   uint8
	Sort     int
}

type PositionItem struct {
	ID         uint64
	Code       string
	Name       string
	Remark     string
	Sort       uint16
	IsStop     uint8
	CreateTime int64
	UpdateTime int64
}

type PositionPage struct {
	Items []PositionItem
	Count int64
}

type OrgUnitService interface {
	List(ctx context.Context, tenantID uint64, filter repository.OrgUnitFilter) ([]OrgUnitItem, error)
	Detail(ctx context.Context, tenantID uint64, id uint64) (OrgUnitItem, error)
	Add(ctx context.Context, input OrgUnitInput) error
	Edit(ctx context.Context, input OrgUnitInput) error
	Delete(ctx context.Context, tenantID uint64, id uint64) error
}

type PositionService interface {
	ListAll(ctx context.Context, tenantID uint64) ([]PositionItem, error)
	List(ctx context.Context, tenantID uint64, filter repository.PositionFilter, pageNo int, pageSize int) (PositionPage, error)
	Detail(ctx context.Context, tenantID uint64, id uint64) (PositionItem, error)
	Add(ctx context.Context, input PositionInput) error
	Edit(ctx context.Context, input PositionInput) error
	Delete(ctx context.Context, tenantID uint64, id uint64) error
}

type orgUnitService struct {
	repo repository.OrgUnitRepository
}

type positionService struct {
	repo repository.PositionRepository
}

func NewOrgUnitService(repo repository.OrgUnitRepository) OrgUnitService {
	return orgUnitService{repo: repo}
}

func NewPositionService(repo repository.PositionRepository) PositionService {
	return positionService{repo: repo}
}

func (srv orgUnitService) List(ctx context.Context, tenantID uint64, filter repository.OrgUnitFilter) ([]OrgUnitItem, error) {
	orgs, err := srv.repo.ListOrgUnits(ctx, tenantID, filter)
	if err != nil {
		return nil, err
	}
	return orgUnitItemsFromModels(orgs), nil
}

func (srv orgUnitService) Detail(ctx context.Context, tenantID uint64, id uint64) (OrgUnitItem, error) {
	org, err := srv.repo.FindOrgUnitByID(ctx, tenantID, id)
	if err != nil {
		return OrgUnitItem{}, mapOrganizationRecordError(err, ErrOrgUnitNotFound)
	}
	return orgUnitItemFromModel(org), nil
}

func (srv orgUnitService) Add(ctx context.Context, input OrgUnitInput) error {
	name := strings.TrimSpace(input.Name)
	if input.ParentID == 0 {
		count, err := srv.repo.CountRootOrgUnits(ctx, input.TenantID, 0)
		if err != nil {
			return err
		}
		if count > 0 {
			return ErrRootOrgUnitExists
		}
	} else if _, err := srv.repo.FindOrgUnitByID(ctx, input.TenantID, input.ParentID); err != nil {
		return mapOrganizationRecordError(err, ErrParentOrgUnitNotFound)
	}
	return srv.repo.CreateOrgUnit(ctx, makeadmin.OrgUnit{
		TenantID: input.TenantID,
		ParentID: input.ParentID,
		Code:     newOrgUnitCode(),
		Name:     name,
		Status:   statusFromStop(input.IsStop),
		Sort:     normalizeSort(input.Sort),
	})
}

func (srv orgUnitService) Edit(ctx context.Context, input OrgUnitInput) error {
	current, err := srv.repo.FindOrgUnitByID(ctx, input.TenantID, input.ID)
	if err != nil {
		return mapOrganizationRecordError(err, ErrOrgUnitNotFound)
	}
	if current.ParentID == 0 && input.ParentID > 0 {
		return ErrRootOrgUnitParentLocked
	}
	if input.ID == input.ParentID {
		return ErrOrgUnitSelfParent
	}
	if input.ParentID > 0 {
		if _, err := srv.repo.FindOrgUnitByID(ctx, input.TenantID, input.ParentID); err != nil {
			return mapOrganizationRecordError(err, ErrParentOrgUnitNotFound)
		}
	} else {
		count, err := srv.repo.CountRootOrgUnits(ctx, input.TenantID, input.ID)
		if err != nil {
			return err
		}
		if count > 0 {
			return ErrRootOrgUnitExists
		}
	}
	current.ParentID = input.ParentID
	current.Name = strings.TrimSpace(input.Name)
	current.Status = statusFromStop(input.IsStop)
	current.Sort = normalizeSort(input.Sort)
	return srv.repo.UpdateOrgUnit(ctx, current)
}

func (srv orgUnitService) Delete(ctx context.Context, tenantID uint64, id uint64) error {
	current, err := srv.repo.FindOrgUnitByID(ctx, tenantID, id)
	if err != nil {
		return mapOrganizationRecordError(err, ErrOrgUnitNotFound)
	}
	if current.ParentID == 0 {
		return ErrRootOrgUnitDeleteLocked
	}
	count, err := srv.repo.CountChildOrgUnits(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrOrgUnitHasChildren
	}
	count, err = srv.repo.CountActiveAdminsByOrgID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrOrgUnitInUse
	}
	return srv.repo.DeleteOrgUnit(ctx, tenantID, id)
}

func (srv positionService) ListAll(ctx context.Context, tenantID uint64) ([]PositionItem, error) {
	positions, err := srv.repo.ListAllPositions(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	return positionItemsFromModels(positions), nil
}

func (srv positionService) List(ctx context.Context, tenantID uint64, filter repository.PositionFilter, pageNo int, pageSize int) (PositionPage, error) {
	positions, count, err := srv.repo.ListPositions(ctx, tenantID, filter, orgPageLimit(pageSize), orgPageOffset(pageNo, pageSize))
	if err != nil {
		return PositionPage{}, err
	}
	return PositionPage{Items: positionItemsFromModels(positions), Count: count}, nil
}

func (srv positionService) Detail(ctx context.Context, tenantID uint64, id uint64) (PositionItem, error) {
	position, err := srv.repo.FindPositionByID(ctx, tenantID, id)
	if err != nil {
		return PositionItem{}, mapOrganizationRecordError(err, ErrPositionNotFound)
	}
	return positionItemFromModel(position), nil
}

func (srv positionService) Add(ctx context.Context, input PositionInput) error {
	code := strings.TrimSpace(input.Code)
	if code == "" {
		code = newPositionCode()
	}
	name := strings.TrimSpace(input.Name)
	if err := srv.ensurePositionUnique(ctx, input.TenantID, code, name, 0); err != nil {
		return err
	}
	return srv.repo.CreatePosition(ctx, makeadmin.Position{
		TenantID: input.TenantID,
		Code:     code,
		Name:     name,
		Remark:   input.Remark,
		Status:   statusFromStop(input.IsStop),
		Sort:     normalizeSort(input.Sort),
	})
}

func (srv positionService) Edit(ctx context.Context, input PositionInput) error {
	current, err := srv.repo.FindPositionByID(ctx, input.TenantID, input.ID)
	if err != nil {
		return mapOrganizationRecordError(err, ErrPositionNotFound)
	}
	code := strings.TrimSpace(input.Code)
	if code == "" {
		code = current.Code
	}
	name := strings.TrimSpace(input.Name)
	if err := srv.ensurePositionUnique(ctx, input.TenantID, code, name, input.ID); err != nil {
		return err
	}
	current.Code = code
	current.Name = name
	current.Remark = input.Remark
	current.Status = statusFromStop(input.IsStop)
	current.Sort = normalizeSort(input.Sort)
	return srv.repo.UpdatePosition(ctx, current)
}

func (srv positionService) Delete(ctx context.Context, tenantID uint64, id uint64) error {
	if _, err := srv.repo.FindPositionByID(ctx, tenantID, id); err != nil {
		return mapOrganizationRecordError(err, ErrPositionNotFound)
	}
	count, err := srv.repo.CountActiveAdminsByPositionID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrPositionInUse
	}
	return srv.repo.DeletePosition(ctx, tenantID, id)
}

func (srv positionService) ensurePositionUnique(ctx context.Context, tenantID uint64, code string, name string, excludeID uint64) error {
	count, err := srv.repo.CountPositionsByName(ctx, tenantID, name, excludeID)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrPositionNameExists
	}
	count, err = srv.repo.CountPositionsByCode(ctx, tenantID, code, excludeID)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrPositionCodeExists
	}
	return nil
}

func orgUnitItemsFromModels(models []makeadmin.OrgUnit) []OrgUnitItem {
	result := make([]OrgUnitItem, 0, len(models))
	for _, item := range models {
		result = append(result, orgUnitItemFromModel(item))
	}
	return result
}

func orgUnitItemFromModel(model makeadmin.OrgUnit) OrgUnitItem {
	return OrgUnitItem{
		ID:         model.ID,
		ParentID:   model.ParentID,
		Name:       model.Name,
		Sort:       model.Sort,
		IsStop:     stopFromStatus(model.Status),
		CreateTime: model.CreateTime,
		UpdateTime: model.UpdateTime,
	}
}

func positionItemsFromModels(models []makeadmin.Position) []PositionItem {
	result := make([]PositionItem, 0, len(models))
	for _, item := range models {
		result = append(result, positionItemFromModel(item))
	}
	return result
}

func positionItemFromModel(model makeadmin.Position) PositionItem {
	return PositionItem{
		ID:         model.ID,
		Code:       model.Code,
		Name:       model.Name,
		Remark:     model.Remark,
		Sort:       model.Sort,
		IsStop:     stopFromStatus(model.Status),
		CreateTime: model.CreateTime,
		UpdateTime: model.UpdateTime,
	}
}

func mapOrganizationRecordError(err error, notFound error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return notFound
	}
	return err
}

func statusFromStop(isStop uint8) uint8 {
	if isStop == 1 {
		return makeadmin.StatusDisabled
	}
	return makeadmin.StatusEnabled
}

func stopFromStatus(status uint8) uint8 {
	if status == makeadmin.StatusDisabled {
		return 1
	}
	return 0
}

func orgPageLimit(pageSize int) int {
	if pageSize <= 0 {
		return 20
	}
	return pageSize
}

func orgPageOffset(pageNo int, pageSize int) int {
	if pageNo <= 0 {
		pageNo = 1
	}
	return orgPageLimit(pageSize) * (pageNo - 1)
}

func newOrgUnitCode() string {
	return fmt.Sprintf("org_%d", time.Now().UnixNano())
}

func newPositionCode() string {
	return fmt.Sprintf("position_%d", time.Now().UnixNano())
}
