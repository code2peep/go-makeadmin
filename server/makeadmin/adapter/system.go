package adapter

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"go-makeadmin/admin/schemas/req"
	"go-makeadmin/admin/schemas/resp"
	"go-makeadmin/config"
	"go-makeadmin/core"
	"go-makeadmin/core/response"
	"go-makeadmin/makeadmin/repository"
	"go-makeadmin/makeadmin/security"
	makeadminsvc "go-makeadmin/makeadmin/service"
	makeadmintenant "go-makeadmin/makeadmin/tenant"
	"go-makeadmin/model/makeadmin"
	"go-makeadmin/util"
)

const (
	ContextAuthSourceKey = "makeadmin_auth_source"
	ContextIdentityKey   = "makeadmin_identity"
	ContextTenantKey     = "makeadmin_tenant"
	ContextAuthSource    = "makeadmin"
)

var ErrUnavailable = errors.New("makeadmin adapter is unavailable")

type requestIdentityContextKey struct{}

type SystemAdapter interface {
	Available(ctx context.Context) bool
	Login(c *gin.Context, loginReq *req.SystemLoginReq) (resp.SystemLoginResp, error)
	Logout(ctx context.Context, token string) error
	Self(ctx context.Context, adminID uint64) (resp.SystemAuthAdminSelfResp, error)
	MenuRoute(ctx context.Context, adminID uint64) ([]interface{}, error)
}

type systemAdapter struct {
	db *gorm.DB
}

func NewSystemAdapter(db *gorm.DB) SystemAdapter {
	return systemAdapter{db: db}
}

func MarkMakeAdminContext(c *gin.Context, identity makeadminsvc.Identity) {
	tenantCtx, ok := makeadmintenant.FromContext(c.Request.Context())
	if !ok {
		tenantCtx = makeadmintenant.Context{TenantID: identity.TenantID, Source: makeadmintenant.SourceJWT}
		c.Request = c.Request.WithContext(makeadmintenant.WithContext(c.Request.Context(), tenantCtx))
	}
	c.Request = c.Request.WithContext(WithIdentityRequestContext(c.Request.Context(), identity))
	c.Set(ContextAuthSourceKey, ContextAuthSource)
	c.Set(ContextIdentityKey, identity)
	c.Set(ContextTenantKey, tenantCtx)
	c.Set(config.AdminConfig.ReqTenantIdKey, tenantCtx.TenantID)
}

func WithIdentityRequestContext(ctx context.Context, identity makeadminsvc.Identity) context.Context {
	return context.WithValue(ctx, requestIdentityContextKey{}, identity)
}

func IsMakeAdminContext(c *gin.Context) bool {
	source, ok := c.Get(ContextAuthSourceKey)
	return ok && source == ContextAuthSource
}

func IdentityFromContext(c *gin.Context) (makeadminsvc.Identity, bool) {
	value, ok := c.Get(ContextIdentityKey)
	if !ok {
		return makeadminsvc.Identity{}, false
	}
	identity, ok := value.(makeadminsvc.Identity)
	return identity, ok
}

func IdentityFromRequestContext(ctx context.Context) (makeadminsvc.Identity, bool) {
	if ctx == nil {
		return makeadminsvc.Identity{}, false
	}
	identity, ok := ctx.Value(requestIdentityContextKey{}).(makeadminsvc.Identity)
	return identity, ok
}

func (adapter systemAdapter) Available(ctx context.Context) bool {
	if adapter.db == nil || !adapter.db.Migrator().HasTable(&makeadmin.Admin{}) {
		return false
	}
	var count int64
	err := adapter.db.WithContext(ctx).
		Model(&makeadmin.Admin{}).
		Where("delete_time = ? AND password_hash <> ''", 0).
		Where("password_hash NOT LIKE ? AND password_hash NOT LIKE ?", "%INSTALL_TIME%", "%REPLACE_ME%").
		Count(&count).
		Error
	return err == nil && count > 0
}

func (adapter systemAdapter) Login(c *gin.Context, loginReq *req.SystemLoginReq) (resp.SystemLoginResp, error) {
	if !adapter.Available(c.Request.Context()) {
		return resp.SystemLoginResp{}, ErrUnavailable
	}
	tenantCtx, err := tenantContextFromGin(c)
	if err != nil {
		return resp.SystemLoginResp{}, mapTenantError(err)
	}
	ua := core.UAParser.Parse(c.GetHeader("user-agent"))
	result, err := adapter.authService(true).Login(c.Request.Context(), makeadminsvc.LoginInput{
		TenantID: tenantCtx.TenantID,
		Username: loginReq.Username,
		Password: loginReq.Password,
		IP:       c.ClientIP(),
		OS:       ua.Os.Family,
		Browser:  ua.UserAgent.Family,
	})
	if err != nil {
		return resp.SystemLoginResp{}, mapAuthError(err)
	}
	return resp.SystemLoginResp{Token: result.Token}, nil
}

func (adapter systemAdapter) Logout(ctx context.Context, token string) error {
	return adapter.authService(true).Logout(ctx, token)
}

func (adapter systemAdapter) Self(ctx context.Context, adminID uint64) (resp.SystemAuthAdminSelfResp, error) {
	identity, err := adapter.authService(false).BuildIdentityByAdminID(ctx, tenantIDFromContext(ctx), adminID)
	if err != nil {
		return resp.SystemAuthAdminSelfResp{}, mapAuthError(err)
	}
	user := resp.SystemAuthAdminSelfOneResp{
		ID:        uint(identity.AdminID),
		Username:  identity.Username,
		Nickname:  identity.Nickname,
		Avatar:    util.UrlUtil.ToAbsoluteUrl(identity.Avatar),
		Role:      roleLabel(identity),
		IsDisable: disabledFlag(identity),
	}
	return resp.SystemAuthAdminSelfResp{User: user, Permissions: identity.Permissions}, nil
}

func (adapter systemAdapter) MenuRoute(ctx context.Context, adminID uint64) ([]interface{}, error) {
	auth := adapter.authService(false)
	identity, err := auth.BuildIdentityByAdminID(ctx, tenantIDFromContext(ctx), adminID)
	if err != nil {
		return nil, mapAuthError(err)
	}
	menus, err := auth.ListRouteMenus(ctx, identity)
	if err != nil {
		return nil, err
	}
	return util.ArrayUtil.ListToTree(routeMenuMaps(menus), "id", "pid", "children"), nil
}

func (adapter systemAdapter) authService(withSession bool) makeadminsvc.AuthService {
	repo := repository.NewAuthRepository(adapter.db)
	sessionStore := makeadminsvc.SessionStore(makeadminsvc.UnavailableSessionStore{})
	if withSession {
		sessionStore = makeadminsvc.NewRedisSessionStore(core.GetRedis(), config.Config.RedisPrefix)
	}
	return makeadminsvc.NewAuthServiceWithDependencies(
		repo,
		nil,
		makeadminsvc.NewJWTTokenCodec(config.Config.Secret),
		sessionStore,
		makeadminsvc.DefaultSessionTTLSeconds,
	)
}

func mapAuthError(err error) error {
	switch {
	case errors.Is(err, ErrUnavailable), errors.Is(err, security.ErrPasswordPlaceholder):
		return ErrUnavailable
	case errors.Is(err, makeadminsvc.ErrInvalidCredential), errors.Is(err, gorm.ErrRecordNotFound):
		return response.LoginAccountError
	case errors.Is(err, makeadminsvc.ErrAdminDisabled):
		return response.LoginDisableError
	case errors.Is(err, makeadminsvc.ErrAdminDeleted):
		return response.LoginAccountError
	default:
		return err
	}
}

func roleLabel(identity makeadminsvc.Identity) string {
	if identity.IsSuper {
		return "超级管理员"
	}
	if len(identity.RoleIDs) == 0 {
		return ""
	}
	return strconv.FormatUint(identity.RoleIDs[0], 10)
}

func disabledFlag(identity makeadminsvc.Identity) uint8 {
	if identity.Status == makeadmin.StatusEnabled {
		return 0
	}
	return 1
}

func routeMenuMaps(menus []makeadminsvc.RouteMenu) []map[string]interface{} {
	menuByID := make(map[uint64]makeadminsvc.RouteMenu, len(menus))
	for _, menu := range menus {
		menuByID[menu.ID] = menu
	}
	result := make([]map[string]interface{}, 0, len(menus))
	for _, menu := range menus {
		result = append(result, map[string]interface{}{
			"id":         uint(menu.ID),
			"pid":        uint(menu.ParentID),
			"menuType":   toLegacyMenuType(menu.MenuType),
			"menuName":   menu.Name,
			"menuIcon":   menu.Icon,
			"menuSort":   menu.Sort,
			"perms":      menu.Perms,
			"paths":      toLegacyRoutePath(menu, menuByID),
			"component":  menu.Component,
			"selected":   strings.Trim(menu.ActivePath, "/"),
			"params":     makeadminsvc.ParamsFromMenuMeta(menu.Meta),
			"isCache":    menu.IsCache,
			"isShow":     uint8(1),
			"isDisable":  uint8(0),
			"createTime": core.TsTime(0),
			"updateTime": core.TsTime(0),
		})
	}
	return result
}

func toLegacyMenuType(menuType string) string {
	switch menuType {
	case makeadmin.MenuTypeCatalog:
		return "M"
	case makeadmin.MenuTypePage:
		return "C"
	default:
		return "A"
	}
}

func toLegacyRoutePath(menu makeadminsvc.RouteMenu, menuByID map[uint64]makeadminsvc.RouteMenu) string {
	current := strings.Trim(menu.RoutePath, "/")
	parent, ok := menuByID[menu.ParentID]
	if !ok {
		return current
	}
	parentPath := strings.Trim(parent.RoutePath, "/")
	if parentPath != "" && strings.HasPrefix(current, parentPath+"/") {
		return strings.TrimPrefix(current, parentPath+"/")
	}
	return current
}
