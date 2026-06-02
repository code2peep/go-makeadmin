package common

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

var AlbumGroup = core.Group("/common", newAlbumHandler, regAlbum, middleware.TokenAuth())

func newAlbumHandler(makeadminFile makeadminadapter.FileAdapter) *albumHandler {
	return &albumHandler{makeadminFile: makeadminFile}
}

func regAlbum(rg *gin.RouterGroup, group *core.GroupBase) error {
	return group.Reg(func(handle *albumHandler) {
		rg.GET("/album/albumList", handle.albumList)
		rg.POST("/album/albumRename", middleware.RecordLog("相册文件重命名"), handle.albumRename)
		rg.POST("/album/albumMove", middleware.RecordLog("相册文件移动"), handle.albumMove)
		rg.POST("/album/albumDel", middleware.RecordLog("相册文件删除"), handle.albumDel)
		rg.GET("/album/cateList", handle.cateList)
		rg.POST("/album/cateAdd", middleware.RecordLog("相册分类新增"), handle.cateAdd)
		rg.POST("/album/cateRename", middleware.RecordLog("相册分类重命名"), handle.cateRename)
		rg.POST("/album/cateDel", middleware.RecordLog("相册分类删除"), handle.cateDel)
	})
}

type albumHandler struct {
	makeadminFile makeadminadapter.FileAdapter
}

// albumList 相册文件列表
func (ah albumHandler) albumList(c *gin.Context) {
	var page request.PageReq
	var listReq req.CommonAlbumListReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyQuery(c, &page)) {
		return
	}
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyQuery(c, &listReq)) {
		return
	}
	res, err := ah.makeadminFile.AlbumList(c.Request.Context(), page, listReq)
	response.CheckAndRespWithData(c, res, err)
}

// albumRename 相册文件重命名
func (ah albumHandler) albumRename(c *gin.Context) {
	var rnReq req.CommonAlbumRenameReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyJSON(c, &rnReq)) {
		return
	}
	response.CheckAndResp(c, ah.makeadminFile.AlbumRename(c.Request.Context(), rnReq.ID, rnReq.Name))
}

// albumMove 相册文件移动
func (ah albumHandler) albumMove(c *gin.Context) {
	var mvReq req.CommonAlbumMoveReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyJSON(c, &mvReq)) {
		return
	}
	response.CheckAndResp(c, ah.makeadminFile.AlbumMove(c.Request.Context(), mvReq.Ids, mvReq.Cid))
}

// albumDel 相册文件删除
func (ah albumHandler) albumDel(c *gin.Context) {
	var delReq req.CommonAlbumDelReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyJSON(c, &delReq)) {
		return
	}
	response.CheckAndResp(c, ah.makeadminFile.AlbumDel(c.Request.Context(), delReq.Ids))
}

// cateList 类目列表
func (ah albumHandler) cateList(c *gin.Context) {
	var listReq req.CommonCateListReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyQuery(c, &listReq)) {
		return
	}
	res, err := ah.makeadminFile.CateList(c.Request.Context(), listReq)
	response.CheckAndRespWithData(c, res, err)
}

// cateAdd 类目新增
func (ah albumHandler) cateAdd(c *gin.Context) {
	var addReq req.CommonCateAddReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyJSON(c, &addReq)) {
		return
	}
	response.CheckAndResp(c, ah.makeadminFile.CateAdd(c.Request.Context(), addReq))
}

// cateRename 类目命名
func (ah albumHandler) cateRename(c *gin.Context) {
	var rnReq req.CommonCateRenameReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyJSON(c, &rnReq)) {
		return
	}
	response.CheckAndResp(c, ah.makeadminFile.CateRename(c.Request.Context(), rnReq.ID, rnReq.Name))
}

// cateDel 类目删除
func (ah albumHandler) cateDel(c *gin.Context) {
	var delReq req.CommonCateDelReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyJSON(c, &delReq)) {
		return
	}
	response.CheckAndResp(c, ah.makeadminFile.CateDel(c.Request.Context(), delReq.ID))
}
