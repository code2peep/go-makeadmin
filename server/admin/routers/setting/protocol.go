package setting

import (
	"github.com/gin-gonic/gin"
	"go-makeadmin/admin/schemas/req"
	"go-makeadmin/admin/service/setting"
	"go-makeadmin/core"
	"go-makeadmin/core/response"
	makeadminadapter "go-makeadmin/makeadmin/adapter"
	"go-makeadmin/middleware"
	"go-makeadmin/util"
)

var ProtocolGroup = core.Group("/setting", newProtocolHandler, regProtocol, middleware.TokenAuth())

func newProtocolHandler(srv setting.ISettingProtocolService, makeadminProtocol makeadminadapter.ProtocolAdapter) *protocolHandler {
	return &protocolHandler{srv: srv, makeadminProtocol: makeadminProtocol}
}

func regProtocol(rg *gin.RouterGroup, group *core.GroupBase) error {
	return group.Reg(func(handle *protocolHandler) {
		rg.GET("/protocol/detail", handle.detail)
		rg.POST("/protocol/save", handle.save)
	})
}

type protocolHandler struct {
	srv               setting.ISettingProtocolService
	makeadminProtocol makeadminadapter.ProtocolAdapter
}

// detail 获取政策信息
func (ph protocolHandler) detail(c *gin.Context) {
	if ph.makeadminProtocol.Available(c.Request.Context()) {
		res, err := ph.makeadminProtocol.Detail(c.Request.Context())
		response.CheckAndRespWithData(c, res, err)
		return
	}
	res, err := ph.srv.Detail()
	response.CheckAndRespWithData(c, res, err)
}

// save 保存政策信息
func (ph protocolHandler) save(c *gin.Context) {
	var pReq req.SettingProtocolReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyJSON(c, &pReq)) {
		return
	}
	if ph.makeadminProtocol.Available(c.Request.Context()) {
		response.CheckAndResp(c, ph.makeadminProtocol.Save(c.Request.Context(), pReq))
		return
	}
	response.CheckAndResp(c, ph.srv.Save(pReq))
}
