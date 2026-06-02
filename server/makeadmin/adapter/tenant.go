package adapter

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"

	"go-makeadmin/core/response"
	"go-makeadmin/makeadmin/repository"
	makeadmintenant "go-makeadmin/makeadmin/tenant"
)

func tenantIDFromContext(ctx context.Context) uint64 {
	return makeadmintenant.IDFromContext(ctx)
}

func dataScopeFromContext(ctx context.Context) repository.DataScopeFilter {
	identity, ok := IdentityFromRequestContext(ctx)
	if !ok {
		return repository.DataScopeFilter{Enabled: true, NoAccess: true}
	}
	return identity.DataScope
}

func tenantContextFromGin(c *gin.Context) (makeadmintenant.Context, error) {
	tenantCtx, err := makeadmintenant.ResolveLogin(makeadmintenant.HeaderValue(c.Request.Header))
	if err != nil {
		return makeadmintenant.Context{}, err
	}
	c.Request = c.Request.WithContext(makeadmintenant.WithContext(c.Request.Context(), tenantCtx))
	return tenantCtx, nil
}

func mapTenantError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, makeadmintenant.ErrInvalidTenantHeader):
		return response.AssertArgumentError.Make("租户参数无效")
	case errors.Is(err, makeadmintenant.ErrTenantUnsupported):
		return response.NoPermission.Make("当前阶段不支持切换租户")
	default:
		return err
	}
}
