package setting

import (
	"github.com/gin-gonic/gin"
	"go-makeadmin/admin/schemas/req"
	"go-makeadmin/admin/service/setting"
	"go-makeadmin/core"
	"go-makeadmin/core/response"
	"go-makeadmin/middleware"
	"go-makeadmin/util"
)

var ProtocolGroup = core.Group("/setting", newProtocolHandler, regProtocol, middleware.TokenAuth())

func newProtocolHandler(srv setting.ISettingProtocolService) *protocolHandler {
	return &protocolHandler{srv: srv}
}

func regProtocol(rg *gin.RouterGroup, group *core.GroupBase) error {
	return group.Reg(func(handle *protocolHandler) {
		rg.GET("/protocol/detail", handle.detail)
		rg.POST("/protocol/save", handle.save)
	})
}

type protocolHandler struct {
	srv setting.ISettingProtocolService
}

//detail 获取政策信息
func (ph protocolHandler) detail(c *gin.Context) {
	res, err := ph.srv.Detail()
	response.CheckAndRespWithData(c, res, err)
}

//save 保存政策信息
func (ph protocolHandler) save(c *gin.Context) {
	var pReq req.SettingProtocolReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyJSON(c, &pReq)) {
		return
	}
	response.CheckAndResp(c, ph.srv.Save(pReq))
}
