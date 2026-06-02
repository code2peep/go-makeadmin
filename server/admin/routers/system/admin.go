package system

import (
	"github.com/gin-gonic/gin"
	"go-makeadmin/admin/schemas/req"
	"go-makeadmin/config"
	"go-makeadmin/core"
	"go-makeadmin/core/request"
	"go-makeadmin/core/response"
	makeadminadapter "go-makeadmin/makeadmin/adapter"
	"go-makeadmin/middleware"
	"go-makeadmin/util"
)

var AdminGroup = core.Group("/system", newAdminHandler, regAdmin, middleware.TokenAuth())

func newAdminHandler(makeadminAdapter makeadminadapter.SystemAdapter, makeadminAdmin makeadminadapter.AdminAdapter) *adminHandler {
	return &adminHandler{makeadminAdapter: makeadminAdapter, makeadminAdmin: makeadminAdmin}
}

func regAdmin(rg *gin.RouterGroup, group *core.GroupBase) error {
	return group.Reg(func(handle *adminHandler) {
		rg.GET("/admin/self", handle.self)
		rg.GET("/admin/list", handle.list)
		rg.GET("/admin/detail", handle.detail)
		rg.POST("/admin/add", middleware.RecordLog("管理员新增"), handle.add)
		rg.POST("/admin/edit", middleware.RecordLog("管理员编辑"), handle.edit)
		rg.POST("/admin/upInfo", middleware.RecordLog("管理员更新"), handle.upInfo)
		rg.POST("/admin/del", middleware.RecordLog("管理员删除"), handle.del)
		rg.POST("/admin/disable", middleware.RecordLog("管理员状态切换"), handle.disable)
	})
}

type adminHandler struct {
	makeadminAdapter makeadminadapter.SystemAdapter
	makeadminAdmin   makeadminadapter.AdminAdapter
}

// self 管理员信息
func (ah adminHandler) self(c *gin.Context) {
	adminId := config.AdminConfig.GetAdminId(c)
	res, err := ah.makeadminAdapter.Self(c.Request.Context(), uint64(adminId))
	response.CheckAndRespWithData(c, res, err)
}

// list 管理员列表
func (ah adminHandler) list(c *gin.Context) {
	var page request.PageReq
	var listReq req.SystemAuthAdminListReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyQuery(c, &page)) {
		return
	}
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyQuery(c, &listReq)) {
		return
	}
	res, err := ah.makeadminAdmin.List(c.Request.Context(), page, listReq)
	response.CheckAndRespWithData(c, res, err)
}

// detail 管理员详细
func (ah adminHandler) detail(c *gin.Context) {
	var detailReq req.SystemAuthAdminDetailReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyQuery(c, &detailReq)) {
		return
	}
	res, err := ah.makeadminAdmin.Detail(c.Request.Context(), detailReq.ID)
	response.CheckAndRespWithData(c, res, err)
}

// add 管理员新增
func (ah adminHandler) add(c *gin.Context) {
	var addReq req.SystemAuthAdminAddReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyJSON(c, &addReq)) {
		return
	}
	response.CheckAndResp(c, ah.makeadminAdmin.Add(c.Request.Context(), addReq))
}

// edit 管理员编辑
func (ah adminHandler) edit(c *gin.Context) {
	var editReq req.SystemAuthAdminEditReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyJSON(c, &editReq)) {
		return
	}
	response.CheckAndResp(c, ah.makeadminAdmin.Edit(c.Request.Context(), editReq))
}

// upInfo 管理员更新
func (ah adminHandler) upInfo(c *gin.Context) {
	var updateReq req.SystemAuthAdminUpdateReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyJSON(c, &updateReq)) {
		return
	}
	response.CheckAndResp(c, ah.makeadminAdmin.UpdateSelf(c, updateReq, config.AdminConfig.GetAdminId(c)))
}

// del 管理员删除
func (ah adminHandler) del(c *gin.Context) {
	var delReq req.SystemAuthAdminDelReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyJSON(c, &delReq)) {
		return
	}
	response.CheckAndResp(c, ah.makeadminAdmin.Del(c.Request.Context(), config.AdminConfig.GetAdminId(c), delReq.ID))
}

// disable 管理员状态切换
func (ah adminHandler) disable(c *gin.Context) {
	var disableReq req.SystemAuthAdminDisableReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyJSON(c, &disableReq)) {
		return
	}
	response.CheckAndResp(c, ah.makeadminAdmin.Disable(c.Request.Context(), config.AdminConfig.GetAdminId(c), disableReq.ID))
}
