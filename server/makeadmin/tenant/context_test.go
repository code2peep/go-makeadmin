package tenant

import (
	"context"
	"errors"
	"testing"
)

func TestContextDefaultsToGlobalTenant(t *testing.T) {
	if got := IDFromContext(context.Background()); got != 0 {
		t.Fatalf("IDFromContext() = %d, want 0", got)
	}
	tenantCtx, ok := FromContext(WithContext(context.Background(), Context{TenantID: 9, Source: SourceJWT}))
	if !ok || tenantCtx.TenantID != 9 || tenantCtx.Source != SourceJWT {
		t.Fatalf("FromContext() = %#v, %v", tenantCtx, ok)
	}
}

func TestResolveLoginOnlyAllowsGlobalTenantInP2(t *testing.T) {
	tenantCtx, err := ResolveLogin("")
	if err != nil || tenantCtx.TenantID != 0 || tenantCtx.Source != SourceDefault {
		t.Fatalf("ResolveLogin(empty) = %#v, %v", tenantCtx, err)
	}
	tenantCtx, err = ResolveLogin("0")
	if err != nil || tenantCtx.TenantID != 0 || tenantCtx.Source != SourceHeader {
		t.Fatalf("ResolveLogin(0) = %#v, %v", tenantCtx, err)
	}
	if _, err = ResolveLogin("2"); !errors.Is(err, ErrTenantUnsupported) {
		t.Fatalf("ResolveLogin(2) error = %v, want ErrTenantUnsupported", err)
	}
}

func TestResolveAuthenticatedRejectsTenantMismatch(t *testing.T) {
	tenantCtx, err := ResolveAuthenticated(3, "")
	if err != nil || tenantCtx.TenantID != 3 || tenantCtx.Source != SourceJWT {
		t.Fatalf("ResolveAuthenticated(jwt) = %#v, %v", tenantCtx, err)
	}
	tenantCtx, err = ResolveAuthenticated(3, "3")
	if err != nil || tenantCtx.TenantID != 3 || tenantCtx.Source != SourceHeader {
		t.Fatalf("ResolveAuthenticated(header) = %#v, %v", tenantCtx, err)
	}
	if _, err = ResolveAuthenticated(3, "4"); !errors.Is(err, ErrTenantMismatch) {
		t.Fatalf("ResolveAuthenticated(mismatch) error = %v, want ErrTenantMismatch", err)
	}
	if _, err = ResolveAuthenticated(3, "abc"); !errors.Is(err, ErrInvalidTenantHeader) {
		t.Fatalf("ResolveAuthenticated(invalid) error = %v, want ErrInvalidTenantHeader", err)
	}
}
