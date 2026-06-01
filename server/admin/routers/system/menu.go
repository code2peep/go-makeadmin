package system

import (
	"github.com/gin-gonic/gin"
	"go-makeadmin/admin/schemas/req"
	"go-makeadmin/admin/service/system"
	"go-makeadmin/config"
	"go-makeadmin/core"
	"go-makeadmin/core/response"
	makeadminadapter "go-makeadmin/makeadmin/adapter"
	"go-makeadmin/middleware"
	"go-makeadmin/util"
)

var MenuGroup = core.Group("/system", newMenuHandler, regMenu, middleware.TokenAuth())

func newMenuHandler(srv system.ISystemAuthMenuService, makeadminAdapter makeadminadapter.SystemAdapter, makeadminMenu makeadminadapter.MenuAdapter) *menuHandler {
	return &menuHandler{srv: srv, makeadminAdapter: makeadminAdapter, makeadminMenu: makeadminMenu}
}

func regMenu(rg *gin.RouterGroup, group *core.GroupBase) error {
	return group.Reg(func(handle *menuHandler) {
		rg.GET("/menu/route", handle.route)
		rg.GET("/menu/list", handle.list)
		rg.GET("/menu/detail", handle.detail)
		rg.POST("/menu/add", handle.add)
		rg.POST("/menu/edit", handle.edit)
		rg.POST("/menu/del", handle.del)
	})
}

type menuHandler struct {
	srv              system.ISystemAuthMenuService
	makeadminAdapter makeadminadapter.SystemAdapter
	makeadminMenu    makeadminadapter.MenuAdapter
}

// route 菜单路由
func (mh menuHandler) route(c *gin.Context) {
	adminId := config.AdminConfig.GetAdminId(c)
	if makeadminadapter.IsMakeAdminContext(c) {
		res, err := mh.makeadminAdapter.MenuRoute(c.Request.Context(), uint64(adminId))
		response.CheckAndRespWithData(c, res, err)
		return
	}
	res, err := mh.srv.SelectMenuByRoleId(c, adminId)
	response.CheckAndRespWithData(c, res, err)
}

// list 菜单列表
func (mh menuHandler) list(c *gin.Context) {
	if mh.makeadminMenu.Available(c.Request.Context()) {
		res, err := mh.makeadminMenu.List(c.Request.Context())
		response.CheckAndRespWithData(c, res, err)
		return
	}
	res, err := mh.srv.List()
	response.CheckAndRespWithData(c, res, err)
}

// detail 菜单详情
func (mh menuHandler) detail(c *gin.Context) {
	var detailReq req.SystemAuthMenuDetailReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyQuery(c, &detailReq)) {
		return
	}
	if mh.makeadminMenu.Available(c.Request.Context()) {
		res, err := mh.makeadminMenu.Detail(c.Request.Context(), detailReq.ID)
		response.CheckAndRespWithData(c, res, err)
		return
	}
	res, err := mh.srv.Detail(detailReq.ID)
	response.CheckAndRespWithData(c, res, err)
}

// add 新增菜单
func (mh menuHandler) add(c *gin.Context) {
	var addReq req.SystemAuthMenuAddReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyJSON(c, &addReq)) {
		return
	}
	if mh.makeadminMenu.Available(c.Request.Context()) {
		response.CheckAndResp(c, mh.makeadminMenu.Add(c.Request.Context(), addReq))
		return
	}
	response.CheckAndResp(c, mh.srv.Add(addReq))
}

// edit 编辑菜单
func (mh menuHandler) edit(c *gin.Context) {
	var editReq req.SystemAuthMenuEditReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyJSON(c, &editReq)) {
		return
	}
	if mh.makeadminMenu.Available(c.Request.Context()) {
		response.CheckAndResp(c, mh.makeadminMenu.Edit(c.Request.Context(), editReq))
		return
	}
	response.CheckAndResp(c, mh.srv.Edit(editReq))
}

// del 删除菜单
func (mh menuHandler) del(c *gin.Context) {
	var delReq req.SystemAuthMenuDelReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyJSON(c, &delReq)) {
		return
	}
	if mh.makeadminMenu.Available(c.Request.Context()) {
		response.CheckAndResp(c, mh.makeadminMenu.Del(c.Request.Context(), delReq.ID))
		return
	}
	response.CheckAndResp(c, mh.srv.Del(delReq.ID))
}
