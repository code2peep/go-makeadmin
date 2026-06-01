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

var WebsiteGroup = core.Group("/setting", newWebsiteHandler, regWebsite, middleware.TokenAuth())

func newWebsiteHandler(srv setting.ISettingWebsiteService, makeadminWebsite makeadminadapter.WebsiteAdapter) *websiteHandler {
	return &websiteHandler{srv: srv, makeadminWebsite: makeadminWebsite}
}

func regWebsite(rg *gin.RouterGroup, group *core.GroupBase) error {
	return group.Reg(func(handle *websiteHandler) {
		rg.GET("/website/detail", handle.detail)
		rg.POST("/website/save", handle.save)
	})
}

type websiteHandler struct {
	srv              setting.ISettingWebsiteService
	makeadminWebsite makeadminadapter.WebsiteAdapter
}

// detail 获取网站信息
func (wh websiteHandler) detail(c *gin.Context) {
	if wh.makeadminWebsite.Available(c.Request.Context()) {
		res, err := wh.makeadminWebsite.Detail(c.Request.Context())
		response.CheckAndRespWithData(c, res, err)
		return
	}
	res, err := wh.srv.Detail()
	response.CheckAndRespWithData(c, res, err)
}

// save 保存网站信息
func (wh websiteHandler) save(c *gin.Context) {
	var wsReq req.SettingWebsiteReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyJSON(c, &wsReq)) {
		return
	}
	if wh.makeadminWebsite.Available(c.Request.Context()) {
		response.CheckAndResp(c, wh.makeadminWebsite.Save(c.Request.Context(), wsReq))
		return
	}
	response.CheckAndResp(c, wh.srv.Save(wsReq))
}
