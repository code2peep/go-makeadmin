package system

import (
	"errors"

	"github.com/gin-gonic/gin"
	"go-makeadmin/admin/schemas/req"
	"go-makeadmin/admin/service/system"
	"go-makeadmin/core"
	"go-makeadmin/core/response"
	makeadminadapter "go-makeadmin/makeadmin/adapter"
	"go-makeadmin/middleware"
	"go-makeadmin/util"
)

var LoginGroup = core.Group("/system", newLoginHandler, regLogin, middleware.TokenAuth())

func newLoginHandler(srv system.ISystemLoginService, makeadminAdapter makeadminadapter.SystemAdapter) *loginHandler {
	return &loginHandler{srv: srv, makeadminAdapter: makeadminAdapter}
}

func regLogin(rg *gin.RouterGroup, group *core.GroupBase) error {
	return group.Reg(func(handle *loginHandler) {
		rg.POST("/login", handle.login)
		rg.POST("/logout", handle.logout)
	})
}

type loginHandler struct {
	srv              system.ISystemLoginService
	makeadminAdapter makeadminadapter.SystemAdapter
}

// login 登录系统
func (lh loginHandler) login(c *gin.Context) {
	var loginReq req.SystemLoginReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyJSON(c, &loginReq)) {
		return
	}
	if lh.makeadminAdapter.Available(c.Request.Context()) {
		res, err := lh.makeadminAdapter.Login(c, &loginReq)
		if !errors.Is(err, makeadminadapter.ErrUnavailable) {
			response.CheckAndRespWithData(c, res, err)
			return
		}
	}
	res, err := lh.srv.Login(c, &loginReq)
	response.CheckAndRespWithData(c, res, err)
}

// logout 登录退出
func (lh loginHandler) logout(c *gin.Context) {
	var logoutReq req.SystemLogoutReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyHeader(c, &logoutReq)) {
		return
	}
	if makeadminadapter.IsMakeAdminContext(c) {
		response.CheckAndResp(c, lh.makeadminAdapter.Logout(c.Request.Context(), logoutReq.Token))
		return
	}
	response.CheckAndResp(c, lh.srv.Logout(&logoutReq))
}
