package service

import (
	"context"
	"testing"

	"go-makeadmin/makeadmin/repository"
	"go-makeadmin/model/makeadmin"
)

type fakeLogRepository struct {
	loginFilter repository.LoginLogFilter
	auditFilter repository.AuditLogFilter
	limit       int
	offset      int
	loginLogs   []makeadmin.LoginLog
	auditLogs   []repository.AuditLogRow
}

func (repo *fakeLogRepository) ListLoginLogs(ctx context.Context, filter repository.LoginLogFilter, limit int, offset int) ([]makeadmin.LoginLog, int64, error) {
	repo.loginFilter = filter
	repo.limit = limit
	repo.offset = offset
	return repo.loginLogs, int64(len(repo.loginLogs)), nil
}

func (repo *fakeLogRepository) ListAuditLogs(ctx context.Context, filter repository.AuditLogFilter, limit int, offset int) ([]repository.AuditLogRow, int64, error) {
	repo.auditFilter = filter
	repo.limit = limit
	repo.offset = offset
	return repo.auditLogs, int64(len(repo.auditLogs)), nil
}

func TestListLoginLogsBuildsRepositoryFilter(t *testing.T) {
	repo := &fakeLogRepository{
		loginLogs: []makeadmin.LoginLog{{ID: 1, Username: "admin", Status: 1}},
	}
	srv := NewLogService(repo)

	result, err := srv.ListLoginLogs(context.Background(), makeadmin.GlobalTenantID, LoginLogFilter{
		Username:  "admin",
		Status:    1,
		StartTime: 100,
		EndTime:   200,
	}, LogPageInput{PageNo: 2, PageSize: 10})
	if err != nil {
		t.Fatalf("ListLoginLogs() error = %v", err)
	}
	if result.Count != 1 || len(result.Items) != 1 {
		t.Fatalf("ListLoginLogs() = %#v", result)
	}
	if repo.loginFilter.TenantID != makeadmin.GlobalTenantID ||
		repo.loginFilter.Username != "admin" ||
		repo.loginFilter.Status != 1 ||
		repo.loginFilter.StartTime != 100 ||
		repo.loginFilter.EndTime != 200 {
		t.Fatalf("ListLoginLogs() filter = %#v", repo.loginFilter)
	}
	if repo.limit != 10 || repo.offset != 10 {
		t.Fatalf("ListLoginLogs() page limit=%d offset=%d", repo.limit, repo.offset)
	}
}

func TestListAuditLogsBuildsRepositoryFilter(t *testing.T) {
	repo := &fakeLogRepository{
		auditLogs: []repository.AuditLogRow{{
			AuditLog: makeadmin.AuditLog{ID: 2, Action: "上传图片", Method: "POST", Status: 1},
			Username: "admin",
			Nickname: "admin",
		}},
	}
	srv := NewLogService(repo)

	result, err := srv.ListAuditLogs(context.Background(), makeadmin.GlobalTenantID, AuditLogFilter{
		Title:     "上传",
		Username:  "admin",
		IP:        "127.0.0.1",
		Method:    "POST",
		Status:    2,
		Path:      "/api/common/upload/image",
		StartTime: 100,
		EndTime:   200,
	}, LogPageInput{PageNo: 1, PageSize: 20})
	if err != nil {
		t.Fatalf("ListAuditLogs() error = %v", err)
	}
	if result.Count != 1 || len(result.Items) != 1 {
		t.Fatalf("ListAuditLogs() = %#v", result)
	}
	if repo.auditFilter.TenantID != makeadmin.GlobalTenantID ||
		repo.auditFilter.Title != "上传" ||
		repo.auditFilter.Username != "admin" ||
		repo.auditFilter.IP != "127.0.0.1" ||
		repo.auditFilter.Method != "POST" ||
		repo.auditFilter.Status != 2 ||
		repo.auditFilter.Path != "/api/common/upload/image" ||
		repo.auditFilter.StartTime != 100 ||
		repo.auditFilter.EndTime != 200 {
		t.Fatalf("ListAuditLogs() filter = %#v", repo.auditFilter)
	}
	if repo.limit != 20 || repo.offset != 0 {
		t.Fatalf("ListAuditLogs() page limit=%d offset=%d", repo.limit, repo.offset)
	}
}

func TestLogPageDefaults(t *testing.T) {
	repo := &fakeLogRepository{}
	srv := NewLogService(repo)

	_, err := srv.ListLoginLogs(context.Background(), makeadmin.GlobalTenantID, LoginLogFilter{}, LogPageInput{})
	if err != nil {
		t.Fatalf("ListLoginLogs() error = %v", err)
	}
	if repo.limit != 20 || repo.offset != 0 {
		t.Fatalf("ListLoginLogs() default page limit=%d offset=%d", repo.limit, repo.offset)
	}
}
