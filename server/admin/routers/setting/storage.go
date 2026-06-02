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

var StorageGroup = core.Group("/setting", newStorageHandler, regStorage, middleware.TokenAuth())

func newStorageHandler(makeadminStorage makeadminadapter.StorageAdapter) *storageHandler {
	return &storageHandler{makeadminStorage: makeadminStorage}
}

func regStorage(rg *gin.RouterGroup, group *core.GroupBase) error {
	return group.Reg(func(handle *storageHandler) {
		rg.GET("/storage/list", handle.list)
		rg.GET("/storage/detail", handle.detail)
		rg.POST("/storage/edit", handle.edit)
		rg.POST("/storage/change", handle.change)
	})
}

type storageHandler struct {
	makeadminStorage makeadminadapter.StorageAdapter
}

// list 存储列表
func (sh storageHandler) list(c *gin.Context) {
	res, err := sh.makeadminStorage.List(c.Request.Context())
	response.CheckAndRespWithData(c, res, err)
}

// detail 存储详情
func (sh storageHandler) detail(c *gin.Context) {
	var detailReq req.SettingStorageDetailReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyQuery(c, &detailReq)) {
		return
	}
	res, err := sh.makeadminStorage.Detail(c.Request.Context(), detailReq.Alias)
	response.CheckAndRespWithData(c, res, err)
}

// edit 存储编辑
func (sh storageHandler) edit(c *gin.Context) {
	var editReq req.SettingStorageEditReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyBody(c, &editReq)) {
		return
	}
	response.CheckAndResp(c, sh.makeadminStorage.Edit(c.Request.Context(), editReq))
}

// change 存储切换
func (sh storageHandler) change(c *gin.Context) {
	var changeReq req.SettingStorageChangeReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyBody(c, &changeReq)) {
		return
	}
	response.CheckAndResp(c, sh.makeadminStorage.Change(c.Request.Context(), changeReq.Alias, changeReq.Status))
}
