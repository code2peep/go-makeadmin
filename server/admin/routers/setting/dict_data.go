package setting

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

var DictDataGroup = core.Group("/setting", newDictDataHandler, regDictData, middleware.TokenAuth())

func newDictDataHandler(makeadminDict makeadminadapter.DictAdapter) *dictDataHandler {
	return &dictDataHandler{makeadminDict: makeadminDict}
}

func regDictData(rg *gin.RouterGroup, group *core.GroupBase) error {
	return group.Reg(func(handle *dictDataHandler) {
		rg.GET("/dict/data/all", handle.all)
		rg.GET("/dict/data/list", handle.list)
		rg.GET("/dict/data/detail", handle.detail)
		rg.POST("/dict/data/add", handle.add)
		rg.POST("/dict/data/edit", handle.edit)
		rg.POST("/dict/data/del", handle.del)
	})
}

type dictDataHandler struct {
	makeadminDict makeadminadapter.DictAdapter
}

// all 字典数据所有
func (ddh dictDataHandler) all(c *gin.Context) {
	var allReq req.SettingDictDataListReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyQuery(c, &allReq)) {
		return
	}
	res, err := ddh.makeadminDict.DataAll(c.Request.Context(), allReq)
	response.CheckAndRespWithData(c, res, err)
}

// list 字典数据列表
func (ddh dictDataHandler) list(c *gin.Context) {
	var page request.PageReq
	var listReq req.SettingDictDataListReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyQuery(c, &page)) {
		return
	}
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyQuery(c, &listReq)) {
		return
	}
	res, err := ddh.makeadminDict.DataList(c.Request.Context(), page, listReq)
	response.CheckAndRespWithData(c, res, err)
}

// detail 字典数据详情
func (ddh dictDataHandler) detail(c *gin.Context) {
	var detailReq req.SettingDictDataDetailReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyQuery(c, &detailReq)) {
		return
	}
	res, err := ddh.makeadminDict.DataDetail(c.Request.Context(), detailReq.ID)
	response.CheckAndRespWithData(c, res, err)
}

// add 字典数据新增
func (ddh dictDataHandler) add(c *gin.Context) {
	var addReq req.SettingDictDataAddReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyJSON(c, &addReq)) {
		return
	}
	response.CheckAndResp(c, ddh.makeadminDict.DataAdd(c.Request.Context(), addReq))
}

// edit 字典数据编辑
func (ddh dictDataHandler) edit(c *gin.Context) {
	var editReq req.SettingDictDataEditReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyJSON(c, &editReq)) {
		return
	}
	response.CheckAndResp(c, ddh.makeadminDict.DataEdit(c.Request.Context(), editReq))
}

// del 字典数据删除
func (ddh dictDataHandler) del(c *gin.Context) {
	var delReq req.SettingDictDataDelReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyJSON(c, &delReq)) {
		return
	}
	response.CheckAndResp(c, ddh.makeadminDict.DataDel(c.Request.Context(), delReq))
}
