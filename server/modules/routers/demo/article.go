package demo

import (
	"github.com/gin-gonic/gin"

	"go-makeadmin/core"
	"go-makeadmin/core/response"
	"go-makeadmin/middleware"
)

var ArticleGroup = core.Group("/", newArticleHandler, regArticle, middleware.TokenAuth())

var InitRouters = []*core.GroupBase{
	ArticleGroup,
}

func newArticleHandler() *articleHandler {
	return &articleHandler{}
}

func regArticle(rg *gin.RouterGroup, group *core.GroupBase) error {
	return group.Reg(func(handle *articleHandler) {
		rg.GET("/article/list", handle.list)
		rg.GET("/article/detail", handle.detail)
		rg.POST("/article/add", handle.readonly)
		rg.POST("/article/edit", handle.readonly)
		rg.POST("/article/del", handle.readonly)
	})
}

type articleHandler struct{}

func (ah articleHandler) list(c *gin.Context) {
	response.OkWithData(c, response.PageResp{
		Count:    0,
		PageNo:   1,
		PageSize: 20,
		Lists:    []interface{}{},
	})
}

func (ah articleHandler) detail(c *gin.Context) {
	response.OkWithData(c, gin.H{
		"module":            "article",
		"runtimeRegistered": true,
	})
}

func (ah articleHandler) readonly(c *gin.Context) {
	response.FailWithMsg(c, response.Failed, "demo module is read-only")
}
