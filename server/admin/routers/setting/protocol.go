package setting

import (
	"github.com/gin-gonic/gin"
	"go-makeadmin/admin/schemas/req"
	"go-makeadmin/core"
	"go-makeadmin/core/response"
	makeadminadapter "go-makeadmin/makeadmin/adapter"
	"go-makeadmin/middleware"
	"go-makeadmin/util"
)

var ProtocolGroup = core.Group("/setting", newProtocolHandler, regProtocol, middleware.TokenAuth())

func newProtocolHandler(makeadminProtocol makeadminadapter.ProtocolAdapter) *protocolHandler {
	return &protocolHandler{makeadminProtocol: makeadminProtocol}
}

func regProtocol(rg *gin.RouterGroup, group *core.GroupBase) error {
	return group.Reg(func(handle *protocolHandler) {
		rg.GET("/protocol/detail", handle.detail)
		rg.POST("/protocol/save", handle.save)
	})
}

type protocolHandler struct {
	makeadminProtocol makeadminadapter.ProtocolAdapter
}

// detail 获取政策信息
func (ph protocolHandler) detail(c *gin.Context) {
	res, err := ph.makeadminProtocol.Detail(c.Request.Context())
	response.CheckAndRespWithData(c, res, err)
}

// save 保存政策信息
func (ph protocolHandler) save(c *gin.Context) {
	var pReq req.SettingProtocolReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyJSON(c, &pReq)) {
		return
	}
	response.CheckAndResp(c, ph.makeadminProtocol.Save(c.Request.Context(), pReq))
}
