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

var PostGroup = core.Group("/system", newPostHandler, regPost, middleware.TokenAuth())

func newPostHandler(makeadminPosition makeadminadapter.PositionAdapter) *postHandler {
	return &postHandler{makeadminPosition: makeadminPosition}
}

func regPost(rg *gin.RouterGroup, group *core.GroupBase) error {
	return group.Reg(func(handle *postHandler) {
		rg.GET("/post/all", handle.all)
		rg.GET("/post/list", handle.list)
		rg.GET("/post/detail", handle.detail)
		rg.POST("/post/add", handle.add)
		rg.POST("/post/edit", handle.edit)
		rg.POST("/post/del", handle.del)
	})
}

type postHandler struct {
	makeadminPosition makeadminadapter.PositionAdapter
}

// all 岗位所有
func (ph postHandler) all(c *gin.Context) {
	res, err := ph.makeadminPosition.All(c.Request.Context())
	response.CheckAndRespWithData(c, res, err)
}

// list 岗位列表
func (ph postHandler) list(c *gin.Context) {
	var page request.PageReq
	var listReq req.SystemAuthPostListReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyQuery(c, &page)) {
		return
	}
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyQuery(c, &listReq)) {
		return
	}
	res, err := ph.makeadminPosition.List(c.Request.Context(), page, listReq)
	response.CheckAndRespWithData(c, res, err)
}

// detail 岗位详情
func (ph postHandler) detail(c *gin.Context) {
	var detailReq req.SystemAuthPostDetailReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyQuery(c, &detailReq)) {
		return
	}
	res, err := ph.makeadminPosition.Detail(c.Request.Context(), detailReq.ID)
	response.CheckAndRespWithData(c, res, err)
}

// add 岗位新增
func (ph postHandler) add(c *gin.Context) {
	var addReq req.SystemAuthPostAddReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyBody(c, &addReq)) {
		return
	}
	response.CheckAndResp(c, ph.makeadminPosition.Add(c.Request.Context(), addReq))
}

// edit 岗位编辑
func (ph postHandler) edit(c *gin.Context) {
	var editReq req.SystemAuthPostEditReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyBody(c, &editReq)) {
		return
	}
	response.CheckAndResp(c, ph.makeadminPosition.Edit(c.Request.Context(), editReq))
}

// del 岗位删除
func (ph postHandler) del(c *gin.Context) {
	var delReq req.SystemAuthPostDelReq
	if response.IsFailWithResp(c, util.VerifyUtil.VerifyBody(c, &delReq)) {
		return
	}
	response.CheckAndResp(c, ph.makeadminPosition.Del(c.Request.Context(), delReq.ID))
}
