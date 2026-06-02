package tenant

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"go-makeadmin/model/makeadmin"
)

const (
	HeaderTenantID = "X-Tenant-ID"

	SourceDefault = "default"
	SourceJWT     = "jwt"
	SourceHeader  = "header"
)

var (
	ErrInvalidTenantHeader = errors.New("makeadmin invalid tenant header")
	ErrTenantMismatch      = errors.New("makeadmin tenant mismatch")
	ErrTenantUnsupported   = errors.New("makeadmin tenant unsupported")
)

type Context struct {
	TenantID uint64
	Source   string
}

type contextKey struct{}

func DefaultContext() Context {
	return Context{TenantID: makeadmin.GlobalTenantID, Source: SourceDefault}
}

func WithContext(ctx context.Context, tenantCtx Context) context.Context {
	if tenantCtx.Source == "" {
		tenantCtx.Source = SourceDefault
	}
	return context.WithValue(ctx, contextKey{}, tenantCtx)
}

func FromContext(ctx context.Context) (Context, bool) {
	if ctx == nil {
		return DefaultContext(), false
	}
	tenantCtx, ok := ctx.Value(contextKey{}).(Context)
	if !ok {
		return DefaultContext(), false
	}
	if tenantCtx.Source == "" {
		tenantCtx.Source = SourceDefault
	}
	return tenantCtx, true
}

func IDFromContext(ctx context.Context) uint64 {
	tenantCtx, _ := FromContext(ctx)
	return tenantCtx.TenantID
}

func HeaderValue(header http.Header) string {
	if header == nil {
		return ""
	}
	return header.Get(HeaderTenantID)
}

func ResolveLogin(headerValue string) (Context, error) {
	tenantID, provided, err := parseHeaderTenantID(headerValue)
	if err != nil {
		return Context{}, err
	}
	if !provided {
		return DefaultContext(), nil
	}
	if tenantID != makeadmin.GlobalTenantID {
		return Context{}, ErrTenantUnsupported
	}
	return Context{TenantID: tenantID, Source: SourceHeader}, nil
}

func ResolveAuthenticated(claimTenantID uint64, headerValue string) (Context, error) {
	source := SourceJWT
	tenantID := claimTenantID
	if tenantID == 0 {
		tenantID = makeadmin.GlobalTenantID
	}

	headerTenantID, provided, err := parseHeaderTenantID(headerValue)
	if err != nil {
		return Context{}, err
	}
	if provided {
		if headerTenantID != tenantID {
			return Context{}, ErrTenantMismatch
		}
		source = SourceHeader
	}
	return Context{TenantID: tenantID, Source: source}, nil
}

func parseHeaderTenantID(value string) (uint64, bool, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, false, nil
	}
	tenantID, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, true, ErrInvalidTenantHeader
	}
	return tenantID, true, nil
}
