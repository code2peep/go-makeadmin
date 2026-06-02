package common

import (
	"github.com/gin-gonic/gin"
	"go-makeadmin/core"
	"go-makeadmin/core/response"
	makeadminadapter "go-makeadmin/makeadmin/adapter"
	"go-makeadmin/middleware"
)

var IndexGroup = core.Group("/common", newIndexHandler, regIndex, middleware.TokenAuth())

func newIndexHandler(makeadminIndex makeadminadapter.IndexAdapter) *indexHandler {
	return &indexHandler{makeadminIndex: makeadminIndex}
}

func regIndex(rg *gin.RouterGroup, group *core.GroupBase) error {
	return group.Reg(func(handle *indexHandler) {
		rg.GET("/index/console", handle.console)
		rg.GET("/index/config", handle.config)
	})
}

type indexHandler struct {
	makeadminIndex makeadminadapter.IndexAdapter
}

// console 控制台
func (ih indexHandler) console(c *gin.Context) {
	res, err := ih.makeadminIndex.Console(c.Request.Context())
	response.CheckAndRespWithData(c, res, err)
}

// config 公共配置
func (ih indexHandler) config(c *gin.Context) {
	res, err := ih.makeadminIndex.Config(c.Request.Context())
	response.CheckAndRespWithData(c, res, err)
}
