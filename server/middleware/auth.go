package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"
	"go-makeadmin/config"
	"go-makeadmin/core"
	"go-makeadmin/core/response"
	makeadminadapter "go-makeadmin/makeadmin/adapter"
	"go-makeadmin/makeadmin/repository"
	makeadminsvc "go-makeadmin/makeadmin/service"
	"go-makeadmin/model/makeadmin"
	"go-makeadmin/util"
	"strconv"
	"strings"
	"sync"
)

var (
	authServiceOnce sync.Once
	makeadminAuth   makeadminsvc.AuthService
)

func initAuthServices() {
	authServiceOnce.Do(func() {
		db := core.GetDB()
		makeadminAuth = makeadminsvc.NewAuthServiceWithPasswordHasher(repository.NewAuthRepository(db), nil)
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
	tokenKey := makeadminsvc.SessionTokenKeyPrefix + token
	existCnt := util.RedisUtil.Exists(tokenKey)
	if existCnt < 0 {
		response.Fail(c, response.SystemError)
		c.Abort()
		return true
	}
	if existCnt == 0 {
		return false
	}

	initAuthServices()
	uidStr := util.RedisUtil.Get(tokenKey)
	uid64, err := strconv.ParseUint(uidStr, 10, 64)
	if err != nil {
		core.Logger.Errorf("MakeAdminTokenAuth ParseUint uidStr err: err=[%+v]", err)
		response.Fail(c, response.TokenInvalid)
		c.Abort()
		return true
	}

	identity, err := makeadminAuth.BuildIdentityByAdminID(c.Request.Context(), makeadmin.GlobalTenantID, uid64)
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

	if util.RedisUtil.TTL(tokenKey) < 1800 {
		util.RedisUtil.Expire(tokenKey, makeadminsvc.DefaultSessionTTLSeconds)
	}

	roleID := "0"
	if len(identity.RoleIDs) > 0 {
		roleID = strconv.FormatUint(identity.RoleIDs[0], 10)
	}
	c.Set(config.AdminConfig.ReqAdminIdKey, uint(identity.AdminID))
	c.Set(config.AdminConfig.ReqRoleIdKey, roleID)
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
