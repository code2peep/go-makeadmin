package system

import (
	"github.com/gin-gonic/gin"
	"go-makeadmin/admin/schemas/req"
	"go-makeadmin/admin/service/system"
	"go-makeadmin/core"
	"go-makeadmin/core/response"
	makeadminadapter "go-makeadmin/makeadmin/adapter"
	"go-makeadmin/middleware"
	"go-makeadmin/util"
)

var DeptGroup = core.Group("/system", newDeptHandler, regDept, middleware.TokenAuth())

func newDeptHandler(srv system.ISystemAuthDeptService, makeadminOrg makeadminadapter.OrgUnitAdapter) *deptHandler {
	return &deptHandler{srv: srv, makeadminOrg: makeadminOrg}
}

func regDept(rg *gin.RouterGroup, group *core.GroupBase) error {
	return group.Reg(func(handle *deptHandler) {
		rg.GET("/dept/all", handle.all)
		rg.GET("/dept/list", handle.list)
		rg.GET("/dept/detail", handle.detail)
		rg.POST("/dept/add", handle.add)
		rg.POST("/dept/edit", handle.edit)
		rg.POST("/dept/del", handle.del)
	})
}

type deptHandler struct {
	srv          system.ISystemAuthDeptService
	makeadminOrg makeadminadapter.OrgUnitAdapter
}

// all 部门所有
func (dh deptHandler) all(c *gin.Context) {
	if dh.makeadminOrg.Available(c.Request.Context()) {
		res, err := dh.makeadminOrg.All(c.Request.Context())
		response.CheckAndRespWithData(c, res, err)
		return
	}
	res, err := dh.srv.All()
	response.CheckAndRespWithData(c, res, err)
}

// list 部门列表
func (dh deptHandler) list(c *gin.Context) {
	var listReq req.SystemAuthDeptListReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyQuery(c, &listReq)) {
		return
	}
	if dh.makeadminOrg.Available(c.Request.Context()) {
		res, err := dh.makeadminOrg.List(c.Request.Context(), listReq)
		response.CheckAndRespWithData(c, res, err)
		return
	}
	res, err := dh.srv.List(listReq)
	response.CheckAndRespWithData(c, res, err)
}

// detail 部门详情
func (dh deptHandler) detail(c *gin.Context) {
	var detailReq req.SystemAuthDeptDetailReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyQuery(c, &detailReq)) {
		return
	}
	if dh.makeadminOrg.Available(c.Request.Context()) {
		res, err := dh.makeadminOrg.Detail(c.Request.Context(), detailReq.ID)
		response.CheckAndRespWithData(c, res, err)
		return
	}
	res, err := dh.srv.Detail(detailReq.ID)
	response.CheckAndRespWithData(c, res, err)
}

// add 部门新增
func (dh deptHandler) add(c *gin.Context) {
	var addReq req.SystemAuthDeptAddReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyBody(c, &addReq)) {
		return
	}
	if dh.makeadminOrg.Available(c.Request.Context()) {
		response.CheckAndResp(c, dh.makeadminOrg.Add(c.Request.Context(), addReq))
		return
	}
	response.CheckAndResp(c, dh.srv.Add(addReq))
}

// edit 部门编辑
func (dh deptHandler) edit(c *gin.Context) {
	var editReq req.SystemAuthDeptEditReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyBody(c, &editReq)) {
		return
	}
	if dh.makeadminOrg.Available(c.Request.Context()) {
		response.CheckAndResp(c, dh.makeadminOrg.Edit(c.Request.Context(), editReq))
		return
	}
	response.CheckAndResp(c, dh.srv.Edit(editReq))
}

// del 部门删除
func (dh deptHandler) del(c *gin.Context) {
	var delReq req.SystemAuthDeptDelReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyBody(c, &delReq)) {
		return
	}
	if dh.makeadminOrg.Available(c.Request.Context()) {
		response.CheckAndResp(c, dh.makeadminOrg.Del(c.Request.Context(), delReq.ID))
		return
	}
	response.CheckAndResp(c, dh.srv.Del(delReq.ID))
}
