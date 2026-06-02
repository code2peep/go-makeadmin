package adapter

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"go-makeadmin/admin/schemas/req"
	"go-makeadmin/admin/schemas/resp"
	"go-makeadmin/core"
	"go-makeadmin/core/request"
	"go-makeadmin/core/response"
	"go-makeadmin/makeadmin/repository"
	makeadminsvc "go-makeadmin/makeadmin/service"
	"go-makeadmin/model/makeadmin"
)

type LogAdapter interface {
	Available(ctx context.Context) bool
	Operate(ctx context.Context, page request.PageReq, logReq req.SystemLogOperateReq) (response.PageResp, error)
	Login(ctx context.Context, page request.PageReq, logReq req.SystemLogLoginReq) (response.PageResp, error)
}

type logAdapter struct {
	db *gorm.DB
}

func NewLogAdapter(db *gorm.DB) LogAdapter {
	return logAdapter{db: db}
}

func (adapter logAdapter) Available(ctx context.Context) bool {
	if adapter.db == nil {
		return false
	}
	return adapter.db.Migrator().HasTable(&makeadmin.LoginLog{}) &&
		adapter.db.Migrator().HasTable(&makeadmin.AuditLog{})
}

func (adapter logAdapter) Operate(ctx context.Context, page request.PageReq, logReq req.SystemLogOperateReq) (response.PageResp, error) {
	page = normalizePage(page)
	result, err := adapter.logService().ListAuditLogs(ctx, tenantIDFromContext(ctx), makeadminsvc.AuditLogFilter{
		Title:     logReq.Title,
		Username:  logReq.Username,
		IP:        logReq.Ip,
		Method:    logReq.Type,
		Status:    logReq.Status,
		Path:      logReq.Url,
		StartTime: logReq.StartTime.Unix(),
		EndTime:   logReq.EndTime.Unix(),
		DataScope: dataScopeFromContext(ctx),
	}, makeadminsvc.LogPageInput{PageNo: page.PageNo, PageSize: page.PageSize})
	if err != nil {
		return response.PageResp{}, err
	}
	return response.PageResp{
		PageNo:   page.PageNo,
		PageSize: page.PageSize,
		Count:    result.Count,
		Lists:    auditLogResponses(result.Items),
	}, nil
}

func (adapter logAdapter) Login(ctx context.Context, page request.PageReq, logReq req.SystemLogLoginReq) (response.PageResp, error) {
	page = normalizePage(page)
	result, err := adapter.logService().ListLoginLogs(ctx, tenantIDFromContext(ctx), makeadminsvc.LoginLogFilter{
		Username:  logReq.Username,
		Status:    logReq.Status,
		StartTime: logReq.StartTime.Unix(),
		EndTime:   logReq.EndTime.Unix(),
		DataScope: dataScopeFromContext(ctx),
	}, makeadminsvc.LogPageInput{PageNo: page.PageNo, PageSize: page.PageSize})
	if err != nil {
		return response.PageResp{}, err
	}
	return response.PageResp{
		PageNo:   page.PageNo,
		PageSize: page.PageSize,
		Count:    result.Count,
		Lists:    loginLogResponses(result.Items),
	}, nil
}

func (adapter logAdapter) logService() makeadminsvc.LogService {
	return makeadminsvc.NewLogService(repository.NewLogRepository(adapter.db))
}

func loginLogResponses(items []makeadmin.LoginLog) []resp.SystemLogLoginResp {
	result := make([]resp.SystemLogLoginResp, 0, len(items))
	for _, item := range items {
		result = append(result, resp.SystemLogLoginResp{
			ID:         uint(item.ID),
			Username:   item.Username,
			Ip:         item.IP,
			Os:         item.OS,
			Browser:    item.Browser,
			Status:     legacyLogStatus(item.Status),
			CreateTime: core.TsTime(item.CreateTime),
		})
	}
	return result
}

func auditLogResponses(items []repository.AuditLogRow) []resp.SystemLogOperateResp {
	result := make([]resp.SystemLogOperateResp, 0, len(items))
	for _, item := range items {
		result = append(result, resp.SystemLogOperateResp{
			ID:         uint(item.ID),
			Username:   item.Username,
			Nickname:   item.Nickname,
			Type:       item.Method,
			Title:      item.Action,
			Method:     item.Method,
			Ip:         item.IP,
			Url:        item.Path,
			Args:       item.RequestBody,
			Error:      item.Error,
			Status:     legacyLogStatus(item.Status),
			TaskTime:   fmt.Sprintf("%dms", item.DurationMS),
			StartTime:  core.TsTime(item.StartTime),
			EndTime:    core.TsTime(item.EndTime),
			CreateTime: core.TsTime(item.CreateTime),
		})
	}
	return result
}

func legacyLogStatus(status uint8) int {
	if status == 1 {
		return 1
	}
	return 2
}
