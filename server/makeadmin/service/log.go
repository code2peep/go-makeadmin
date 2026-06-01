package service

import (
	"context"

	"go-makeadmin/makeadmin/repository"
	"go-makeadmin/model/makeadmin"
)

type LoginLogFilter struct {
	Username  string
	Status    int
	StartTime int64
	EndTime   int64
}

type AuditLogFilter struct {
	Title     string
	Username  string
	IP        string
	Method    string
	Status    int
	Path      string
	StartTime int64
	EndTime   int64
}

type LogPageInput struct {
	PageNo   int
	PageSize int
}

type LoginLogPage struct {
	Items []makeadmin.LoginLog
	Count int64
}

type AuditLogPage struct {
	Items []repository.AuditLogRow
	Count int64
}

type LogService interface {
	ListLoginLogs(ctx context.Context, tenantID uint64, filter LoginLogFilter, page LogPageInput) (LoginLogPage, error)
	ListAuditLogs(ctx context.Context, tenantID uint64, filter AuditLogFilter, page LogPageInput) (AuditLogPage, error)
}

type logService struct {
	repo repository.LogRepository
}

func NewLogService(repo repository.LogRepository) LogService {
	return logService{repo: repo}
}

func (srv logService) ListLoginLogs(ctx context.Context, tenantID uint64, filter LoginLogFilter, page LogPageInput) (LoginLogPage, error) {
	items, count, err := srv.repo.ListLoginLogs(ctx, repository.LoginLogFilter{
		TenantID:  tenantID,
		Username:  filter.Username,
		Status:    filter.Status,
		StartTime: filter.StartTime,
		EndTime:   filter.EndTime,
	}, logPageLimit(page), logPageOffset(page))
	if err != nil {
		return LoginLogPage{}, err
	}
	return LoginLogPage{Items: items, Count: count}, nil
}

func (srv logService) ListAuditLogs(ctx context.Context, tenantID uint64, filter AuditLogFilter, page LogPageInput) (AuditLogPage, error) {
	items, count, err := srv.repo.ListAuditLogs(ctx, repository.AuditLogFilter{
		TenantID:  tenantID,
		Title:     filter.Title,
		Username:  filter.Username,
		IP:        filter.IP,
		Method:    filter.Method,
		Status:    filter.Status,
		Path:      filter.Path,
		StartTime: filter.StartTime,
		EndTime:   filter.EndTime,
	}, logPageLimit(page), logPageOffset(page))
	if err != nil {
		return AuditLogPage{}, err
	}
	return AuditLogPage{Items: items, Count: count}, nil
}

func logPageLimit(page LogPageInput) int {
	if page.PageSize <= 0 {
		return 20
	}
	return page.PageSize
}

func logPageOffset(page LogPageInput) int {
	pageNo := page.PageNo
	if pageNo <= 0 {
		pageNo = 1
	}
	return logPageLimit(page) * (pageNo - 1)
}
