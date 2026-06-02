package system

import (
	"github.com/gin-gonic/gin"
	"go-makeadmin/admin/schemas/req"
	"go-makeadmin/core"
	"go-makeadmin/core/response"
	makeadminadapter "go-makeadmin/makeadmin/adapter"
	"go-makeadmin/middleware"
	"go-makeadmin/util"
)

var LoginGroup = core.Group("/system", newLoginHandler, regLogin, middleware.TokenAuth())

func newLoginHandler(makeadminAdapter makeadminadapter.SystemAdapter) *loginHandler {
	return &loginHandler{makeadminAdapter: makeadminAdapter}
}

func regLogin(rg *gin.RouterGroup, group *core.GroupBase) error {
	return group.Reg(func(handle *loginHandler) {
		rg.POST("/login", handle.login)
		rg.POST("/logout", handle.logout)
	})
}

type loginHandler struct {
	makeadminAdapter makeadminadapter.SystemAdapter
}

// login 登录系统
func (lh loginHandler) login(c *gin.Context) {
	var loginReq req.SystemLoginReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyJSON(c, &loginReq)) {
		return
	}
	res, err := lh.makeadminAdapter.Login(c, &loginReq)
	response.CheckAndRespWithData(c, res, err)
}

// logout 登录退出
func (lh loginHandler) logout(c *gin.Context) {
	var logoutReq req.SystemLogoutReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyHeader(c, &logoutReq)) {
		return
	}
	response.CheckAndResp(c, lh.makeadminAdapter.Logout(c.Request.Context(), logoutReq.Token))
}
