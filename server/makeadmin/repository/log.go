package repository

import (
	"context"

	"gorm.io/gorm"

	"go-makeadmin/model/makeadmin"
)

type LoginLogFilter struct {
	TenantID  uint64
	Username  string
	Status    int
	StartTime int64
	EndTime   int64
	DataScope DataScopeFilter
}

type AuditLogFilter struct {
	TenantID  uint64
	Title     string
	Username  string
	IP        string
	Method    string
	Status    int
	Path      string
	StartTime int64
	EndTime   int64
	DataScope DataScopeFilter
}

type AuditLogRow struct {
	makeadmin.AuditLog
	Username string
	Nickname string
}

type LogRepository interface {
	ListLoginLogs(ctx context.Context, filter LoginLogFilter, limit int, offset int) ([]makeadmin.LoginLog, int64, error)
	ListAuditLogs(ctx context.Context, filter AuditLogFilter, limit int, offset int) ([]AuditLogRow, int64, error)
}

type logRepository struct {
	db *gorm.DB
}

func NewLogRepository(db *gorm.DB) LogRepository {
	return logRepository{db: db}
}

func (repo logRepository) ListLoginLogs(ctx context.Context, filter LoginLogFilter, limit int, offset int) ([]makeadmin.LoginLog, int64, error) {
	query := repo.loginLogQuery(ctx, filter)
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	var logs []makeadmin.LoginLog
	err := query.Limit(limit).Offset(offset).Order("id DESC").Find(&logs).Error
	return logs, count, err
}

func (repo logRepository) ListAuditLogs(ctx context.Context, filter AuditLogFilter, limit int, offset int) ([]AuditLogRow, int64, error) {
	query := repo.auditLogQuery(ctx, filter)
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	var logs []AuditLogRow
	err := query.Limit(limit).Offset(offset).Order("log.id DESC").Find(&logs).Error
	return logs, count, err
}

func (repo logRepository) loginLogQuery(ctx context.Context, filter LoginLogFilter) *gorm.DB {
	query := repo.db.WithContext(ctx).Model(&makeadmin.LoginLog{}).Where("tenant_id = ?", filter.TenantID)
	if filter.Username != "" {
		query = query.Where("username LIKE ?", "%"+filter.Username+"%")
	}
	if filter.Status == 1 {
		query = query.Where("status = ?", 1)
	} else if filter.Status == 2 {
		query = query.Where("status = ?", 0)
	}
	if filter.StartTime > 0 {
		query = query.Where("create_time >= ?", filter.StartTime)
	}
	if filter.EndTime > 0 {
		query = query.Where("create_time <= ?", filter.EndTime)
	}
	query = applyDataScopeFilter(repo.db, query, filter.TenantID, "admin_id", filter.DataScope)
	return query
}

func (repo logRepository) auditLogQuery(ctx context.Context, filter AuditLogFilter) *gorm.DB {
	query := repo.db.WithContext(ctx).
		Table("ma_audit_log AS log").
		Select("log.*, admin.username, profile.nickname").
		Joins("LEFT JOIN ma_admin AS admin ON log.admin_id = admin.id").
		Joins("LEFT JOIN ma_admin_profile AS profile ON profile.admin_id = admin.id").
		Where("log.tenant_id = ?", filter.TenantID)
	if filter.Title != "" {
		query = query.Where("log.action LIKE ?", "%"+filter.Title+"%")
	}
	if filter.Username != "" {
		query = query.Where("admin.username LIKE ?", "%"+filter.Username+"%")
	}
	if filter.IP != "" {
		query = query.Where("log.ip LIKE ?", "%"+filter.IP+"%")
	}
	if filter.Method != "" {
		query = query.Where("log.method = ?", filter.Method)
	}
	if filter.Status == 1 {
		query = query.Where("log.status = ?", 1)
	} else if filter.Status == 2 {
		query = query.Where("log.status = ?", 0)
	}
	if filter.Path != "" {
		query = query.Where("log.path = ?", filter.Path)
	}
	if filter.StartTime > 0 {
		query = query.Where("log.create_time >= ?", filter.StartTime)
	}
	if filter.EndTime > 0 {
		query = query.Where("log.create_time <= ?", filter.EndTime)
	}
	query = applyDataScopeFilter(repo.db, query, filter.TenantID, "log.admin_id", filter.DataScope)
	return query
}
