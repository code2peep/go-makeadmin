package adapter

import (
	"context"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"go-makeadmin/admin/schemas/req"
	"go-makeadmin/admin/schemas/resp"
	"go-makeadmin/core"
	"go-makeadmin/core/request"
	"go-makeadmin/core/response"
	"go-makeadmin/makeadmin/repository"
	makeadminsvc "go-makeadmin/makeadmin/service"
	"go-makeadmin/model/makeadmin"
	"go-makeadmin/util"
)

type AdminAdapter interface {
	Available(ctx context.Context) bool
	List(ctx context.Context, page request.PageReq, listReq req.SystemAuthAdminListReq) (response.PageResp, error)
	Detail(ctx context.Context, id uint) (resp.SystemAuthAdminResp, error)
	Add(ctx context.Context, addReq req.SystemAuthAdminAddReq) error
	Edit(ctx context.Context, editReq req.SystemAuthAdminEditReq) error
	UpdateSelf(c *gin.Context, updateReq req.SystemAuthAdminUpdateReq, adminID uint) error
	Del(ctx context.Context, currentAdminID uint, id uint) error
	Disable(ctx context.Context, currentAdminID uint, id uint) error
}

type adminAdapter struct {
	db *gorm.DB
}

func NewAdminAdapter(db *gorm.DB) AdminAdapter {
	return adminAdapter{db: db}
}

func (adapter adminAdapter) Available(ctx context.Context) bool {
	if adapter.db == nil ||
		!adapter.db.Migrator().HasTable(&makeadmin.Admin{}) ||
		!adapter.db.Migrator().HasTable(&makeadmin.AdminProfile{}) ||
		!adapter.db.Migrator().HasTable(&makeadmin.AdminRole{}) ||
		!adapter.db.Migrator().HasTable(&makeadmin.AdminOrg{}) {
		return false
	}
	var count int64
	err := adapter.db.WithContext(ctx).
		Model(&makeadmin.Admin{}).
		Where("delete_time = ? AND password_hash <> ''", 0).
		Count(&count).
		Error
	return err == nil && count > 0
}

func (adapter adminAdapter) List(ctx context.Context, page request.PageReq, listReq req.SystemAuthAdminListReq) (response.PageResp, error) {
	page = normalizePage(page)
	tenantID := tenantIDFromContext(ctx)
	result, err := adapter.adminService().List(ctx, tenantID, repository.AdminFilter{
		Username: listReq.Username,
		Nickname: listReq.Nickname,
		RoleID:   uint64(listReq.Role),
		RoleSet:  listReq.Role >= 0,
	}, page.PageNo, page.PageSize)
	if err != nil {
		return response.PageResp{}, mapAdminError(err)
	}
	return response.PageResp{
		PageNo:   page.PageNo,
		PageSize: page.PageSize,
		Count:    result.Count,
		Lists:    adminResponses(result.Items, false),
	}, nil
}

func (adapter adminAdapter) Detail(ctx context.Context, id uint) (resp.SystemAuthAdminResp, error) {
	item, err := adapter.adminService().Detail(ctx, tenantIDFromContext(ctx), uint64(id))
	if err != nil {
		return resp.SystemAuthAdminResp{}, mapAdminError(err)
	}
	return adminResponse(item, true), nil
}

func (adapter adminAdapter) Add(ctx context.Context, addReq req.SystemAuthAdminAddReq) error {
	tenantID := tenantIDFromContext(ctx)
	return mapAdminError(adapter.adminService().Add(ctx, makeadminsvc.AdminInput{
		TenantID:     tenantID,
		OrgID:        uint64(addReq.DeptId),
		PositionID:   uint64(addReq.PostId),
		Username:     addReq.Username,
		Nickname:     addReq.Nickname,
		Password:     addReq.Password,
		Avatar:       util.UrlUtil.ToRelativeUrl(addReq.Avatar),
		RoleID:       uint64(addReq.Role),
		IsDisable:    addReq.IsDisable,
		IsMultipoint: addReq.IsMultipoint,
	}))
}

func (adapter adminAdapter) Edit(ctx context.Context, editReq req.SystemAuthAdminEditReq) error {
	tenantID := tenantIDFromContext(ctx)
	return mapAdminError(adapter.adminService().Edit(ctx, makeadminsvc.AdminInput{
		ID:           uint64(editReq.ID),
		TenantID:     tenantID,
		OrgID:        uint64(editReq.DeptId),
		PositionID:   uint64(editReq.PostId),
		Username:     editReq.Username,
		Nickname:     editReq.Nickname,
		Password:     editReq.Password,
		Avatar:       util.UrlUtil.ToRelativeUrl(editReq.Avatar),
		RoleID:       uint64(editReq.Role),
		IsDisable:    editReq.IsDisable,
		IsMultipoint: editReq.IsMultipoint,
	}))
}

func (adapter adminAdapter) UpdateSelf(c *gin.Context, updateReq req.SystemAuthAdminUpdateReq, adminID uint) error {
	return mapAdminError(adapter.adminService().UpdateSelf(c.Request.Context(), makeadminsvc.AdminSelfInput{
		ID:           uint64(adminID),
		Nickname:     updateReq.Nickname,
		Avatar:       util.UrlUtil.ToRelativeUrl(updateReq.Avatar),
		Password:     updateReq.Password,
		CurrPassword: updateReq.CurrPassword,
	}))
}

func (adapter adminAdapter) Del(ctx context.Context, currentAdminID uint, id uint) error {
	return mapAdminError(adapter.adminService().Delete(ctx, tenantIDFromContext(ctx), uint64(currentAdminID), uint64(id)))
}

func (adapter adminAdapter) Disable(ctx context.Context, currentAdminID uint, id uint) error {
	return mapAdminError(adapter.adminService().Disable(ctx, uint64(currentAdminID), uint64(id)))
}

func (adapter adminAdapter) adminService() makeadminsvc.AdminService {
	return makeadminsvc.NewAdminService(repository.NewAdminRepository(adapter.db))
}

func adminResponses(items []makeadminsvc.AdminItem, detail bool) []resp.SystemAuthAdminResp {
	result := make([]resp.SystemAuthAdminResp, 0, len(items))
	for _, item := range items {
		result = append(result, adminResponse(item, detail))
	}
	return result
}

func adminResponse(item makeadminsvc.AdminItem, detail bool) resp.SystemAuthAdminResp {
	role := item.RoleLabel
	if detail {
		role = strconv.FormatUint(item.RoleID, 10)
	}
	return resp.SystemAuthAdminResp{
		ID:            uint(item.ID),
		Username:      item.Username,
		Nickname:      item.Nickname,
		Avatar:        util.UrlUtil.ToAbsoluteUrl(item.Avatar),
		Role:          role,
		DeptId:        uint(item.OrgID),
		PostId:        uint(item.PositionID),
		Dept:          item.OrgName,
		IsMultipoint:  item.IsMultipoint,
		IsDisable:     item.IsDisable,
		LastLoginIp:   item.LastLoginIP,
		LastLoginTime: core.TsTime(item.LastLoginTime),
		CreateTime:    core.TsTime(item.CreateTime),
		UpdateTime:    core.TsTime(item.UpdateTime),
	}
}

func mapAdminError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, makeadminsvc.ErrAdminNotFound):
		return response.AssertArgumentError.Make("账号已不存在!")
	case errors.Is(err, makeadminsvc.ErrAdminUsernameExists):
		return response.AssertArgumentError.Make("账号已存在换一个吧！")
	case errors.Is(err, makeadminsvc.ErrAdminNicknameExists):
		return response.AssertArgumentError.Make("昵称已存在换一个吧！")
	case errors.Is(err, makeadminsvc.ErrAdminRoleNotFound):
		return response.AssertArgumentError.Make("角色已不存在!")
	case errors.Is(err, makeadminsvc.ErrAdminRoleDisabled):
		return response.AssertArgumentError.Make("当前角色已被禁用!")
	case errors.Is(err, makeadminsvc.ErrAdminOrgNotFound):
		return response.AssertArgumentError.Make("部门已不存在!")
	case errors.Is(err, makeadminsvc.ErrAdminOrgDisabled):
		return response.AssertArgumentError.Make("当前部门已被停用!")
	case errors.Is(err, makeadminsvc.ErrAdminPositionNotFound):
		return response.AssertArgumentError.Make("岗位已不存在!")
	case errors.Is(err, makeadminsvc.ErrAdminPositionDisabled):
		return response.AssertArgumentError.Make("当前岗位已被停用!")
	case errors.Is(err, makeadminsvc.ErrAdminPasswordInvalid):
		return response.Failed.Make("密码必须在8~72位，且当前密码必须正确")
	case errors.Is(err, makeadminsvc.ErrSystemAdminProtected):
		return response.AssertArgumentError.Make("系统管理员不允许操作!")
	case errors.Is(err, makeadminsvc.ErrAdminSelfProtected):
		return response.AssertArgumentError.Make("不能操作自己!")
	default:
		return err
	}
}
