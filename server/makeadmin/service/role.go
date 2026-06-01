package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"go-makeadmin/makeadmin/repository"
	"go-makeadmin/model/makeadmin"
)

var (
	ErrRoleNotFound        = errors.New("makeadmin role not found")
	ErrRoleNameExists      = errors.New("makeadmin role name exists")
	ErrRoleInUse           = errors.New("makeadmin role in use")
	ErrSystemRoleProtected = errors.New("makeadmin system role protected")
)

type RoleInput struct {
	ID        uint64
	TenantID  uint64
	Name      string
	Remark    string
	Sort      int
	IsDisable uint8
	MenuIDs   []uint64
}

type RoleItem struct {
	ID         uint64
	Name       string
	Remark     string
	MenuIDs    []uint64
	Member     int64
	Sort       uint16
	IsDisable  uint8
	CreateTime int64
	UpdateTime int64
}

type RolePage struct {
	Items []RoleItem
	Count int64
}

type RoleService interface {
	ListAll(ctx context.Context, tenantID uint64) ([]RoleItem, error)
	List(ctx context.Context, tenantID uint64, pageNo int, pageSize int) (RolePage, error)
	Detail(ctx context.Context, tenantID uint64, id uint64) (RoleItem, error)
	Add(ctx context.Context, input RoleInput) error
	Edit(ctx context.Context, input RoleInput) error
	Delete(ctx context.Context, tenantID uint64, id uint64) error
}

type roleService struct {
	repo repository.RoleRepository
}

func NewRoleService(repo repository.RoleRepository) RoleService {
	return roleService{repo: repo}
}

func (srv roleService) ListAll(ctx context.Context, tenantID uint64) ([]RoleItem, error) {
	roles, err := srv.repo.ListAllRoles(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	return srv.roleItems(ctx, tenantID, roles, false)
}

func (srv roleService) List(ctx context.Context, tenantID uint64, pageNo int, pageSize int) (RolePage, error) {
	roles, count, err := srv.repo.ListRoles(ctx, tenantID, rolePageLimit(pageSize), rolePageOffset(pageNo, pageSize))
	if err != nil {
		return RolePage{}, err
	}
	items, err := srv.roleItems(ctx, tenantID, roles, true)
	if err != nil {
		return RolePage{}, err
	}
	return RolePage{Items: items, Count: count}, nil
}

func (srv roleService) Detail(ctx context.Context, tenantID uint64, id uint64) (RoleItem, error) {
	role, err := srv.repo.FindRoleByID(ctx, tenantID, id)
	if err != nil {
		return RoleItem{}, mapRoleRecordError(err, ErrRoleNotFound)
	}
	return srv.roleItem(ctx, tenantID, role, true)
}

func (srv roleService) Add(ctx context.Context, input RoleInput) error {
	name := strings.TrimSpace(input.Name)
	if err := srv.ensureNameUnique(ctx, input.TenantID, name, 0); err != nil {
		return err
	}
	_, err := srv.repo.CreateRoleWithMenuIDs(ctx, makeadmin.Role{
		TenantID: input.TenantID,
		Code:     newRoleCode(),
		Name:     name,
		Remark:   input.Remark,
		Status:   statusFromDisable(input.IsDisable),
		Sort:     roleSortFromInt(input.Sort),
	}, input.MenuIDs)
	return err
}

func (srv roleService) Edit(ctx context.Context, input RoleInput) error {
	role, err := srv.repo.FindRoleByID(ctx, input.TenantID, input.ID)
	if err != nil {
		return mapRoleRecordError(err, ErrRoleNotFound)
	}
	name := strings.TrimSpace(input.Name)
	if err := srv.ensureNameUnique(ctx, input.TenantID, name, input.ID); err != nil {
		return err
	}
	role.Name = name
	role.Remark = input.Remark
	role.Status = statusFromDisable(input.IsDisable)
	role.Sort = roleSortFromInt(input.Sort)
	return srv.repo.UpdateRoleWithMenuIDs(ctx, role, input.MenuIDs)
}

func (srv roleService) Delete(ctx context.Context, tenantID uint64, id uint64) error {
	role, err := srv.repo.FindRoleByID(ctx, tenantID, id)
	if err != nil {
		return mapRoleRecordError(err, ErrRoleNotFound)
	}
	if role.IsSystem == 1 {
		return ErrSystemRoleProtected
	}
	count, err := srv.repo.CountAdminsByRoleID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrRoleInUse
	}
	return srv.repo.DeleteRole(ctx, tenantID, id)
}

func (srv roleService) ensureNameUnique(ctx context.Context, tenantID uint64, name string, excludeID uint64) error {
	count, err := srv.repo.CountRolesByName(ctx, tenantID, name, excludeID)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrRoleNameExists
	}
	return nil
}

func (srv roleService) roleItems(ctx context.Context, tenantID uint64, roles []makeadmin.Role, withMembers bool) ([]RoleItem, error) {
	items := make([]RoleItem, 0, len(roles))
	for _, role := range roles {
		item, err := srv.roleItem(ctx, tenantID, role, withMembers)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (srv roleService) roleItem(ctx context.Context, tenantID uint64, role makeadmin.Role, withMembers bool) (RoleItem, error) {
	menuIDs, err := srv.repo.ListMenuIDsByRoleID(ctx, tenantID, role.ID)
	if err != nil {
		return RoleItem{}, err
	}
	member := int64(0)
	if withMembers {
		member, err = srv.repo.CountAdminsByRoleID(ctx, tenantID, role.ID)
		if err != nil {
			return RoleItem{}, err
		}
	}
	return RoleItem{
		ID:         role.ID,
		Name:       role.Name,
		Remark:     role.Remark,
		MenuIDs:    menuIDs,
		Member:     member,
		Sort:       role.Sort,
		IsDisable:  disableFromStatus(role.Status),
		CreateTime: role.CreateTime,
		UpdateTime: role.UpdateTime,
	}, nil
}

func ParseRoleMenuIDs(raw string) []uint64 {
	if strings.TrimSpace(raw) == "" {
		return []uint64{}
	}
	parts := strings.Split(raw, ",")
	ids := make([]uint64, 0, len(parts))
	seen := map[uint64]struct{}{}
	for _, part := range parts {
		id, err := strconv.ParseUint(strings.TrimSpace(part), 10, 64)
		if err != nil || id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	return ids
}

func mapRoleRecordError(err error, notFound error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return notFound
	}
	return err
}

func statusFromDisable(isDisable uint8) uint8 {
	if isDisable == 1 {
		return makeadmin.StatusDisabled
	}
	return makeadmin.StatusEnabled
}

func disableFromStatus(status uint8) uint8 {
	if status == makeadmin.StatusDisabled {
		return 1
	}
	return 0
}

func roleSortFromInt(value int) uint16 {
	if value <= 0 {
		return 0
	}
	if value > 65535 {
		return 65535
	}
	return uint16(value)
}

func rolePageLimit(pageSize int) int {
	if pageSize <= 0 {
		return 20
	}
	return pageSize
}

func rolePageOffset(pageNo int, pageSize int) int {
	if pageNo <= 0 {
		pageNo = 1
	}
	return rolePageLimit(pageSize) * (pageNo - 1)
}

func newRoleCode() string {
	return fmt.Sprintf("role_%d", time.Now().UnixNano())
}
