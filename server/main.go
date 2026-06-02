package main

import (
	"github.com/gin-gonic/gin"
	adminRouters "go-makeadmin/admin/routers"
	admin "go-makeadmin/admin/service"
	"go-makeadmin/config"
	"go-makeadmin/core"
	"go-makeadmin/core/response"
	genRouters "go-makeadmin/generator/routers"
	gen "go-makeadmin/generator/service"
	"go-makeadmin/middleware"
	moduleRouters "go-makeadmin/modules/routers"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

// initDI 初始化DI
func initDI() {
	regFunctions := admin.InitFunctions
	regFunctions = append(regFunctions, gen.InitFunctions...)
	regFunctions = append(regFunctions, core.GetDB)
	for i := 0; i < len(regFunctions); i++ {
		if err := core.ProvideForDI(regFunctions[i]); err != nil {
			log.Fatalln(err)
		}
	}
}

// initRouter 初始化router
func initRouter() *gin.Engine {
	// 初始化gin
	gin.SetMode(config.Config.GinMode)
	router := gin.New()
	// 设置静态路径
	router.Static(config.Config.PublicPrefix, config.Config.UploadDirectory)
	router.Static(config.Config.StaticPath, config.Config.StaticDirectory)
	// 设置中间件
	router.Use(gin.Logger(), middleware.Cors(), middleware.ErrorRecover())
	// 演示模式
	if config.Config.DisallowModify {
		router.Use(middleware.ShowMode())
	}
	// 特殊异常处理
	router.NoMethod(response.NoMethod)
	router.NoRoute(response.NoRoute)
	// 注册路由
	group := router.Group("/api")
	//core.RegisterGroup(group, routers.CommonGroup, middleware.TokenAuth())
	//core.RegisterGroup(group, routers.MonitorGroup, middleware.TokenAuth())
	//core.RegisterGroup(group, routers.SettingGroup, middleware.TokenAuth())
	//core.RegisterGroup(group, routers.SystemGroup, middleware.TokenAuth())

	routers := adminRouters.InitRouters[:]
	routers = append(routers, genRouters.InitRouters...)
	routers = append(routers, moduleRouters.InitRouters()...)
	for i := 0; i < len(routers); i++ {
		core.RegisterGroup(group, routers[i])
	}
	return router
}

// initServer 初始化server
func initServer(router *gin.Engine) *http.Server {
	return &http.Server{
		Addr:           net.JoinHostPort(config.Config.ServerHost, strconv.Itoa(config.Config.ServerPort)),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

func main() {
	// 刷新日志缓冲
	defer core.Logger.Sync()
	// 程序结束前关闭数据库连接
	defer core.CloseDB()
	defer core.CloseRedis()
	// 初始化DI
	initDI()
	// 初始化router
	router := initRouter()
	// 初始化server
	s := initServer(router)
	// 运行服务
	log.Fatalln(s.ListenAndServe().Error())
}
