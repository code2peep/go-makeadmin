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

var CopyrightGroup = core.Group("/setting", newCopyrightHandler, regCopyright, middleware.TokenAuth())

func newCopyrightHandler(srv setting.ISettingCopyrightService, makeadminCopyright makeadminadapter.CopyrightAdapter) *copyrightHandler {
	return &copyrightHandler{srv: srv, makeadminCopyright: makeadminCopyright}
}

func regCopyright(rg *gin.RouterGroup, group *core.GroupBase) error {
	return group.Reg(func(handle *copyrightHandler) {
		rg.GET("/copyright/detail", handle.detail)
		rg.POST("/copyright/save", handle.save)
	})
}

type copyrightHandler struct {
	srv                setting.ISettingCopyrightService
	makeadminCopyright makeadminadapter.CopyrightAdapter
}

// detail 获取备案信息
func (ch copyrightHandler) detail(c *gin.Context) {
	if ch.makeadminCopyright.Available(c.Request.Context()) {
		res, err := ch.makeadminCopyright.Detail(c.Request.Context())
		response.CheckAndRespWithData(c, res, err)
		return
	}
	res, err := ch.srv.Detail()
	response.CheckAndRespWithData(c, res, err)
}

// save 保存备案信息
func (ch copyrightHandler) save(c *gin.Context) {
	var cReqs []req.SettingCopyrightItemReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyJSONArray(c, &cReqs)) {
		return
	}
	if ch.makeadminCopyright.Available(c.Request.Context()) {
		response.CheckAndResp(c, ch.makeadminCopyright.Save(c.Request.Context(), cReqs))
		return
	}
	response.CheckAndResp(c, ch.srv.Save(cReqs))
}
