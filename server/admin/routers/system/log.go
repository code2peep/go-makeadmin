package system

import (
	"github.com/gin-gonic/gin"
	"go-makeadmin/admin/schemas/req"
	"go-makeadmin/core"
	"go-makeadmin/core/request"
	"go-makeadmin/core/response"
	makeadminadapter "go-makeadmin/makeadmin/adapter"
	"go-makeadmin/middleware"
	"go-makeadmin/util"
)

var LogGroup = core.Group("/system", newLogHandler, regLog, middleware.TokenAuth())

func newLogHandler(makeadminLog makeadminadapter.LogAdapter) *logHandler {
	return &logHandler{makeadminLog: makeadminLog}
}

func regLog(rg *gin.RouterGroup, group *core.GroupBase) error {
	return group.Reg(func(handle *logHandler) {
		rg.GET("/log/operate", handle.operate)
		rg.GET("/log/login", handle.login)
	})
}

type logHandler struct {
	makeadminLog makeadminadapter.LogAdapter
}

// operate 操作日志
func (lh logHandler) operate(c *gin.Context) {
	var page request.PageReq
	var logReq req.SystemLogOperateReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyQuery(c, &page)) {
		return
	}
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyQuery(c, &logReq)) {
		return
	}
	res, err := lh.makeadminLog.Operate(c.Request.Context(), page, logReq)
	response.CheckAndRespWithData(c, res, err)
}

// login 登录日志
func (lh logHandler) login(c *gin.Context) {
	var page request.PageReq
	var logReq req.SystemLogLoginReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyQuery(c, &page)) {
		return
	}
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyQuery(c, &logReq)) {
		return
	}
	res, err := lh.makeadminLog.Login(c.Request.Context(), page, logReq)
	response.CheckAndRespWithData(c, res, err)
}
