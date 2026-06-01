package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"gorm.io/gorm"

	"go-makeadmin/makeadmin/repository"
	"go-makeadmin/model/makeadmin"
)

var (
	ErrMenuNotFound       = errors.New("makeadmin menu not found")
	ErrParentMenuNotFound = errors.New("makeadmin parent menu not found")
	ErrMenuSelfParent     = errors.New("makeadmin menu self parent")
	ErrMenuHasChildren    = errors.New("makeadmin menu has children")
	ErrMenuPermsExists    = errors.New("makeadmin menu permission code exists")
)

type MenuInput struct {
	ID        uint64
	ParentID  uint64
	MenuType  string
	MenuName  string
	MenuIcon  string
	MenuSort  int
	Perms     string
	Paths     string
	Component string
	Selected  string
	Params    string
	IsCache   uint8
	IsShow    uint8
	IsDisable uint8
}

type MenuItem struct {
	ID         uint64
	ParentID   uint64
	MenuType   string
	MenuName   string
	MenuIcon   string
	MenuSort   uint16
	Perms      string
	Paths      string
	Component  string
	Selected   string
	Params     string
	IsCache    uint8
	IsShow     uint8
	IsDisable  uint8
	CreateTime int64
	UpdateTime int64
}

type MenuService interface {
	List(ctx context.Context) ([]MenuItem, error)
	Detail(ctx context.Context, id uint64) (MenuItem, error)
	Add(ctx context.Context, input MenuInput) error
	Edit(ctx context.Context, input MenuInput) error
	Delete(ctx context.Context, id uint64) error
}

type menuService struct {
	repo repository.MenuRepository
}

func NewMenuService(repo repository.MenuRepository) MenuService {
	return menuService{repo: repo}
}

func (srv menuService) List(ctx context.Context) ([]MenuItem, error) {
	menus, err := srv.repo.ListMenus(ctx)
	if err != nil {
		return nil, err
	}
	return srv.menuItems(ctx, menus)
}

func (srv menuService) Detail(ctx context.Context, id uint64) (MenuItem, error) {
	menu, err := srv.repo.FindMenuByID(ctx, id)
	if err != nil {
		return MenuItem{}, mapMenuRecordError(err, ErrMenuNotFound)
	}
	return srv.menuItem(ctx, menu)
}

func (srv menuService) Add(ctx context.Context, input MenuInput) error {
	if err := srv.validateParent(ctx, input.ParentID, 0); err != nil {
		return err
	}
	permission, err := srv.permissionFromInput(ctx, input, 0)
	if err != nil {
		return err
	}
	_, err = srv.repo.CreateMenuWithPermission(ctx, menuFromInput(input), permission)
	return err
}

func (srv menuService) Edit(ctx context.Context, input MenuInput) error {
	current, err := srv.repo.FindMenuByID(ctx, input.ID)
	if err != nil {
		return mapMenuRecordError(err, ErrMenuNotFound)
	}
	if err := srv.validateParent(ctx, input.ParentID, input.ID); err != nil {
		return err
	}
	currentPermissionID, err := srv.currentPermissionID(ctx, input.ID, input.Perms)
	if err != nil {
		return err
	}
	permission, err := srv.permissionFromInput(ctx, input, currentPermissionID)
	if err != nil {
		return err
	}
	menu := menuFromInput(input)
	menu.ID = current.ID
	return srv.repo.UpdateMenuWithPermission(ctx, menu, permission)
}

func (srv menuService) Delete(ctx context.Context, id uint64) error {
	if _, err := srv.repo.FindMenuByID(ctx, id); err != nil {
		return mapMenuRecordError(err, ErrMenuNotFound)
	}
	count, err := srv.repo.CountChildMenus(ctx, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrMenuHasChildren
	}
	return srv.repo.DeleteMenu(ctx, id)
}

func (srv menuService) validateParent(ctx context.Context, parentID uint64, selfID uint64) error {
	if parentID == 0 {
		return nil
	}
	if parentID == selfID {
		return ErrMenuSelfParent
	}
	if _, err := srv.repo.FindMenuByID(ctx, parentID); err != nil {
		return mapMenuRecordError(err, ErrParentMenuNotFound)
	}
	return nil
}

func (srv menuService) permissionFromInput(ctx context.Context, input MenuInput, currentPermissionID uint64) (*makeadmin.Permission, error) {
	code := strings.TrimSpace(input.Perms)
	if code == "" || toMakeAdminMenuType(input.MenuType) == makeadmin.MenuTypeCatalog {
		return nil, nil
	}
	count, err := srv.repo.CountPermissionCode(ctx, code, currentPermissionID)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrMenuPermsExists
	}
	module, resource, action := permissionParts(code)
	return &makeadmin.Permission{
		ID:       currentPermissionID,
		Code:     code,
		Name:     strings.TrimSpace(input.MenuName),
		Module:   module,
		Resource: resource,
		Action:   action,
		Status:   statusFromDisable(input.IsDisable),
		Sort:     normalizeSort(input.MenuSort),
	}, nil
}

func (srv menuService) menuItems(ctx context.Context, menus []makeadmin.Menu) ([]MenuItem, error) {
	items := make([]MenuItem, 0, len(menus))
	for _, menu := range menus {
		item, err := srv.menuItem(ctx, menu)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (srv menuService) menuItem(ctx context.Context, menu makeadmin.Menu) (MenuItem, error) {
	codes, err := srv.repo.ListPermissionCodesByMenuID(ctx, menu.ID)
	if err != nil {
		return MenuItem{}, err
	}
	perms := ""
	if len(codes) > 0 {
		perms = codes[0]
	}
	return MenuItem{
		ID:         menu.ID,
		ParentID:   menu.ParentID,
		MenuType:   legacyMenuTypeFromMakeAdmin(menu.MenuType),
		MenuName:   menu.Name,
		MenuIcon:   menu.Icon,
		MenuSort:   menu.Sort,
		Perms:      perms,
		Paths:      strings.Trim(menu.RoutePath, "/"),
		Component:  menu.Component,
		Selected:   strings.Trim(menu.ActivePath, "/"),
		Params:     ParamsFromMenuMeta(menu.Meta),
		IsCache:    menu.IsCache,
		IsShow:     menu.IsVisible,
		IsDisable:  disableFromStatus(menu.Status),
		CreateTime: menu.CreateTime,
		UpdateTime: menu.UpdateTime,
	}, nil
}

func menuFromInput(input MenuInput) makeadmin.Menu {
	return makeadmin.Menu{
		ParentID:   input.ParentID,
		MenuType:   toMakeAdminMenuType(input.MenuType),
		Name:       strings.TrimSpace(input.MenuName),
		Icon:       input.MenuIcon,
		RoutePath:  normalizeRoutePath(input.Paths),
		RouteName:  routeNameFromPath(input.Paths),
		Component:  input.Component,
		ActivePath: normalizeRoutePath(input.Selected),
		Meta:       metaFromParams(input.Params),
		IsVisible:  input.IsShow,
		IsCache:    input.IsCache,
		Status:     statusFromDisable(input.IsDisable),
		Sort:       normalizeSort(input.MenuSort),
	}
}

func (srv menuService) currentPermissionID(ctx context.Context, menuID uint64, code string) (uint64, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return 0, nil
	}
	codes, err := srv.repo.ListPermissionCodesByMenuID(ctx, menuID)
	if err != nil {
		return 0, err
	}
	if len(codes) == 0 || codes[0] != code {
		return 0, nil
	}
	permission, err := srv.repo.FindPermissionByCode(ctx, code)
	if err != nil {
		return 0, mapMenuRecordError(err, ErrMenuNotFound)
	}
	return permission.ID, nil
}

func toMakeAdminMenuType(menuType string) string {
	switch menuType {
	case "M":
		return makeadmin.MenuTypeCatalog
	case "C":
		return makeadmin.MenuTypePage
	default:
		return makeadmin.MenuTypeAction
	}
}

func legacyMenuTypeFromMakeAdmin(menuType string) string {
	switch menuType {
	case makeadmin.MenuTypeCatalog:
		return "M"
	case makeadmin.MenuTypePage:
		return "C"
	default:
		return "A"
	}
}

func permissionParts(code string) (string, string, string) {
	parts := strings.Split(code, ":")
	module := ""
	resource := ""
	action := ""
	if len(parts) > 0 {
		module = parts[0]
	}
	if len(parts) > 1 {
		resource = parts[1]
	}
	if len(parts) > 2 {
		action = strings.Join(parts[2:], ":")
	}
	return module, resource, action
}

func normalizeRoutePath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" || strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}
	return "/" + strings.Trim(path, "/")
}

func routeNameFromPath(path string) string {
	path = strings.Trim(normalizeRoutePath(path), "/")
	return strings.ReplaceAll(path, "/", ".")
}

func metaFromParams(params string) string {
	params = strings.TrimSpace(params)
	if params == "" {
		return "{}"
	}
	raw, err := json.Marshal(map[string]string{"params": params})
	if err != nil {
		return "{}"
	}
	return string(raw)
}

func ParamsFromMenuMeta(meta string) string {
	var payload struct {
		Params string `json:"params"`
	}
	if err := json.Unmarshal([]byte(meta), &payload); err != nil {
		return ""
	}
	return payload.Params
}

func mapMenuRecordError(err error, notFound error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return notFound
	}
	return err
}
