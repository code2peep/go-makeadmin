package middleware

import (
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go-makeadmin/config"
	"go-makeadmin/core"
	"go-makeadmin/core/response"
	makeadminadapter "go-makeadmin/makeadmin/adapter"
	"go-makeadmin/makeadmin/repository"
	makeadminsvc "go-makeadmin/makeadmin/service"
	makeadmintenant "go-makeadmin/makeadmin/tenant"
	"go-makeadmin/util"
)

var (
	authServiceOnce      sync.Once
	makeadminAuth        makeadminsvc.AuthService
	makeadminTokenCodec  makeadminsvc.TokenCodec
	makeadminSessionRepo makeadminsvc.SessionStore
)

func initAuthServices() {
	authServiceOnce.Do(func() {
		db := core.GetDB()
		makeadminTokenCodec = makeadminsvc.NewJWTTokenCodec(config.Config.Secret)
		makeadminSessionRepo = makeadminsvc.NewRedisSessionStore(core.GetRedis(), config.Config.RedisPrefix)
		makeadminAuth = makeadminsvc.NewAuthServiceWithDependencies(
			repository.NewAuthRepository(db),
			nil,
			makeadminTokenCodec,
			makeadminsvc.UnavailableSessionStore{},
			makeadminsvc.DefaultSessionTTLSeconds,
		)
	})
}

// TokenAuth Token认证中间件
func TokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 路由转权限
		auths := strings.ReplaceAll(strings.Replace(c.Request.URL.Path, "/api/", "", 1), "/", ":")

		// 免登录接口
		if util.ToolsUtil.Contains(config.AdminConfig.NotLoginUri, auths) {
			c.Next()
			return
		}

		// Token是否为空
		token := c.Request.Header.Get("token")
		if token == "" {
			response.Fail(c, response.TokenEmpty)
			c.Abort()
			return
		}
		if handleMakeAdminToken(c, auths, token) {
			return
		}
		response.Fail(c, response.TokenInvalid)
		c.Abort()
	}
}

func handleMakeAdminToken(c *gin.Context, auths string, token string) bool {
	initAuthServices()

	claims, err := makeadminTokenCodec.Parse(token)
	if err != nil {
		core.Logger.Errorf("MakeAdminTokenAuth Parse JWT err: err=[%+v]", err)
		response.Fail(c, response.TokenInvalid)
		c.Abort()
		return true
	}

	adminID, err := makeadminSessionRepo.FindAdminID(c.Request.Context(), claims.SessionID)
	if err != nil {
		core.Logger.Errorf("MakeAdminTokenAuth FindAdminID err: err=[%+v]", err)
		if errors.Is(err, makeadminsvc.ErrSessionStore) {
			response.Fail(c, response.SystemError)
		} else {
			response.Fail(c, response.TokenInvalid)
		}
		c.Abort()
		return true
	}
	if adminID != claims.AdminID {
		core.Logger.Errorf("MakeAdminTokenAuth admin mismatch: claim=[%d] state=[%d]", claims.AdminID, adminID)
		response.Fail(c, response.TokenInvalid)
		c.Abort()
		return true
	}

	tenantCtx, err := makeadmintenant.ResolveAuthenticated(claims.TenantID, makeadmintenant.HeaderValue(c.Request.Header))
	if err != nil {
		core.Logger.Errorf("MakeAdminTokenAuth resolve tenant err: err=[%+v]", err)
		response.Fail(c, response.NoPermission)
		c.Abort()
		return true
	}
	c.Request = c.Request.WithContext(makeadmintenant.WithContext(c.Request.Context(), tenantCtx))

	identity, err := makeadminAuth.BuildIdentityByAdminID(c.Request.Context(), tenantCtx.TenantID, adminID)
	if err != nil {
		core.Logger.Errorf("MakeAdminTokenAuth BuildIdentityByAdminID err: err=[%+v]", err)
		if errors.Is(err, makeadminsvc.ErrAdminDisabled) {
			response.Fail(c, response.LoginDisableError)
		} else {
			response.Fail(c, response.TokenInvalid)
		}
		c.Abort()
		return true
	}

	if remainingTTL := int(claims.ExpiresAt - time.Now().Unix()); remainingTTL > 0 {
		if err := makeadminSessionRepo.Refresh(c.Request.Context(), claims.SessionID, remainingTTL); err != nil {
			core.Logger.Errorf("MakeAdminTokenAuth Refresh session err: err=[%+v]", err)
		}
	}

	roleID := "0"
	if len(identity.RoleIDs) > 0 {
		roleID = strconv.FormatUint(identity.RoleIDs[0], 10)
	}
	c.Set(config.AdminConfig.ReqAdminIdKey, uint(identity.AdminID))
	c.Set(config.AdminConfig.ReqRoleIdKey, roleID)
	c.Set(config.AdminConfig.ReqTenantIdKey, identity.TenantID)
	c.Set(config.AdminConfig.ReqUsernameKey, identity.Username)
	c.Set(config.AdminConfig.ReqNicknameKey, identity.Nickname)
	makeadminadapter.MarkMakeAdminContext(c, identity)

	if util.ToolsUtil.Contains(config.AdminConfig.NotAuthUri, auths) || identity.IsSuper || hasMakeAdminPermission(identity.Permissions, auths) {
		c.Next()
		return true
	}

	response.Fail(c, response.NoPermission)
	c.Abort()
	return true
}

func hasMakeAdminPermission(permissions []string, auths string) bool {
	aliases := []string{auths}
	if strings.HasPrefix(auths, "common:upload:") {
		aliases = append(aliases, strings.TrimPrefix(auths, "common:"))
	}
	for _, permission := range permissions {
		if permission == "*" {
			return true
		}
		for _, alias := range aliases {
			if permission == alias {
				return true
			}
		}
	}
	return false
}
