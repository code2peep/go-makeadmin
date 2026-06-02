package adapter

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"go-makeadmin/config"
	"go-makeadmin/makeadmin/repository"
	makeadminsvc "go-makeadmin/makeadmin/service"
	makeadmintenant "go-makeadmin/makeadmin/tenant"
)

func TestMarkMakeAdminContextStoresTenantContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodGet, "/api/system/admin/self", nil)
	tenantCtx := makeadmintenant.Context{TenantID: 7, Source: makeadmintenant.SourceJWT}
	req = req.WithContext(makeadmintenant.WithContext(req.Context(), tenantCtx))
	c.Request = req

	identity := makeadminsvc.Identity{
		AdminID:  3,
		TenantID: 7,
		Username: "admin",
		DataScope: repository.DataScopeFilter{
			Enabled: true,
			Self:    true,
			AdminID: 3,
		},
	}
	MarkMakeAdminContext(c, identity)

	if got := config.AdminConfig.GetTenantId(c); got != 7 {
		t.Fatalf("GetTenantId() = %d, want 7", got)
	}
	gotTenantCtx, ok := makeadmintenant.FromContext(c.Request.Context())
	if !ok || gotTenantCtx.TenantID != 7 || gotTenantCtx.Source != makeadmintenant.SourceJWT {
		t.Fatalf("FromContext() = %#v, %v", gotTenantCtx, ok)
	}
	gotIdentity, ok := IdentityFromContext(c)
	if !ok || gotIdentity.AdminID != 3 || gotIdentity.TenantID != 7 {
		t.Fatalf("IdentityFromContext() = %#v, %v", gotIdentity, ok)
	}
	gotRequestIdentity, ok := IdentityFromRequestContext(c.Request.Context())
	if !ok || gotRequestIdentity.AdminID != 3 || gotRequestIdentity.TenantID != 7 || !gotRequestIdentity.DataScope.Self {
		t.Fatalf("IdentityFromRequestContext() = %#v, %v", gotRequestIdentity, ok)
	}
}

func TestDataScopeFromContextFailsClosedWithoutIdentity(t *testing.T) {
	scope := dataScopeFromContext(context.Background())

	if !scope.Enabled || !scope.NoAccess {
		t.Fatalf("dataScopeFromContext() = %#v, want no access", scope)
	}
}

func TestTenantContextFromGinRejectsUnsupportedLoginTenant(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/system/login", nil)
	c.Request.Header.Set(makeadmintenant.HeaderTenantID, "2")

	if _, err := tenantContextFromGin(c); !errors.Is(err, makeadmintenant.ErrTenantUnsupported) {
		t.Fatalf("tenantContextFromGin() error = %v, want ErrTenantUnsupported", err)
	}
}
