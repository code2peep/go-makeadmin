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

var DictTypeGroup = core.Group("/setting", newDictTypeHandler, regDictType, middleware.TokenAuth())

func newDictTypeHandler(makeadminDict makeadminadapter.DictAdapter) *dictTypeHandler {
	return &dictTypeHandler{makeadminDict: makeadminDict}
}

func regDictType(rg *gin.RouterGroup, group *core.GroupBase) error {
	return group.Reg(func(handle *dictTypeHandler) {
		rg.GET("/dict/type/all", handle.all)
		rg.GET("/dict/type/list", handle.list)
		rg.GET("/dict/type/detail", handle.detail)
		rg.POST("/dict/type/add", handle.add)
		rg.POST("/dict/type/edit", handle.edit)
		rg.POST("/dict/type/del", handle.del)
	})
}

type dictTypeHandler struct {
	makeadminDict makeadminadapter.DictAdapter
}

// all 字典类型所有
func (dth dictTypeHandler) all(c *gin.Context) {
	res, err := dth.makeadminDict.TypeAll(c.Request.Context())
	response.CheckAndRespWithData(c, res, err)
}

// list 字典类型列表
func (dth dictTypeHandler) list(c *gin.Context) {
	var page request.PageReq
	var listReq req.SettingDictTypeListReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyQuery(c, &page)) {
		return
	}
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyQuery(c, &listReq)) {
		return
	}
	res, err := dth.makeadminDict.TypeList(c.Request.Context(), page, listReq)
	response.CheckAndRespWithData(c, res, err)
}

// detail 字典类型详情
func (dth dictTypeHandler) detail(c *gin.Context) {
	var detailReq req.SettingDictTypeDetailReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyQuery(c, &detailReq)) {
		return
	}
	res, err := dth.makeadminDict.TypeDetail(c.Request.Context(), detailReq.ID)
	response.CheckAndRespWithData(c, res, err)
}

// add 字典类型新增
func (dth dictTypeHandler) add(c *gin.Context) {
	var addReq req.SettingDictTypeAddReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyJSON(c, &addReq)) {
		return
	}
	response.CheckAndResp(c, dth.makeadminDict.TypeAdd(c.Request.Context(), addReq))
}

// edit 字典类型编辑
func (dth dictTypeHandler) edit(c *gin.Context) {
	var editReq req.SettingDictTypeEditReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyJSON(c, &editReq)) {
		return
	}
	response.CheckAndResp(c, dth.makeadminDict.TypeEdit(c.Request.Context(), editReq))
}

// del 字典类型删除
func (dth dictTypeHandler) del(c *gin.Context) {
	var delReq req.SettingDictTypeDelReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyJSON(c, &delReq)) {
		return
	}
	response.CheckAndResp(c, dth.makeadminDict.TypeDel(c.Request.Context(), delReq))
}
