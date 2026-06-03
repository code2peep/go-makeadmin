package common

import (
	"go-makeadmin/config"
	"go-makeadmin/core/response"
	"go-makeadmin/util"
	"gorm.io/gorm"
)

type IIndexService interface {
	Console() (res map[string]interface{}, e error)
	Config() (res map[string]interface{}, e error)
}

// NewIndexService 初始化
func NewIndexService(db *gorm.DB) IIndexService {
	return &indexService{db: db}
}

// indexService 主页服务实现类
type indexService struct {
	db *gorm.DB
}

// Console 控制台数据
func (iSrv indexService) Console() (res map[string]interface{}, e error) {
	// 版本信息
	name, err := util.ConfigUtil.GetVal(iSrv.db, "website", "name", "go-makeadmin")
	if e = response.CheckErr(err, "Console Get err"); e != nil {
		return
	}
	version := map[string]interface{}{
		"name":    name,
		"version": config.Config.Version,
		"website": "https://github.com/code2peep/go-makeadmin",
		"based":   "Go、Gin、Gorm、Vue3、Element Plus、MySQL、Redis",
		"links": map[string]string{
			"github":  "https://github.com/code2peep/go-makeadmin",
			"website": "https://github.com/code2peep/go-makeadmin",
		},
	}
	return map[string]interface{}{
		"version": version,
		"framework": map[string]interface{}{
			"stage":           "P4.10 P4 冻结验收",
			"database":        "go_makeadmin",
			"tables":          "ma_*",
			"auth":            "JWT + Redis session",
			"moduleLifecycle": "manifest + codegen + install/uninstall apply",
		},
		"milestones": []map[string]string{
			{"name": "P1 核心后台", "status": "已冻结", "summary": "登录、菜单、权限、设置、字典、文件、日志和代码生成器切到 ma_*。"},
			{"name": "P2 权限租户", "status": "已冻结", "summary": "JWT、Redis session、租户上下文、数据权限和模块生命周期命令完成。"},
			{"name": "P3 模块产品化", "status": "已冻结", "summary": "脚手架、codegen、manifest、安装卸载和 apply 结果闭环完成。"},
			{"name": "P4 可见后台", "status": "已冻结", "summary": "工作台、模块中心、核心页面入口、空态和失败态完成可见验收。"},
		},
		"validation": []map[string]string{
			{"name": "无库验证", "status": "通过", "scope": "runtime residue、Go test、type-check、build、npm audit"},
			{"name": "模块工具链", "status": "通过", "scope": "manifest、脚手架、codegen、安装卸载计划、写入门禁"},
			{"name": "模块中心", "status": "通过", "scope": "内嵌预览、apply 结果、状态清单"},
			{"name": "核心页面入口", "status": "通过", "scope": "菜单、角色、管理员、部门、网站信息、缓存、日志"},
			{"name": "P4 冻结验收", "status": "通过", "scope": "P4 完成面、人工测试入口、P5 进入条件已收敛"},
			{"name": "本地 API", "status": "可用", "scope": "http://127.0.0.1:18000/api"},
			{"name": "管理端", "status": "可用", "scope": "http://127.0.0.1:5173"},
		},
		"corePages": []map[string]string{
			{"name": "菜单权限", "route": "/menu", "status": "已验收", "scope": "菜单树、权限字符、路由显隐"},
			{"name": "角色管理", "route": "/role", "status": "已验收", "scope": "角色列表、授权入口、数据权限"},
			{"name": "管理员", "route": "/admin", "status": "已验收", "scope": "账号列表、组织岗位、启停"},
			{"name": "组织部门", "route": "/department", "status": "已验收", "scope": "部门树、负责人、状态"},
			{"name": "网站信息", "route": "/information", "status": "已验收", "scope": "站点名称、Logo、备案基础信息"},
			{"name": "系统缓存", "route": "/cache", "status": "已验收", "scope": "缓存清理、本地运行状态"},
			{"name": "系统日志", "route": "/journal", "status": "已验收", "scope": "管理员操作日志、登录日志"},
		},
	}, nil
}

// Config 公共配置
func (iSrv indexService) Config() (res map[string]interface{}, e error) {
	website, err := util.ConfigUtil.Get(iSrv.db, "website")
	if e = response.CheckErr(err, "Config Get err"); e != nil {
		return
	}
	copyrightStr, err := util.ConfigUtil.GetVal(iSrv.db, "website", "copyright", "")
	if e = response.CheckErr(err, "Config GetVal err"); e != nil {
		return
	}
	var copyright []map[string]string
	if copyrightStr != "" {
		err = util.ToolsUtil.JsonToObj(copyrightStr, &copyright)
		if e = response.CheckErr(err, "Config JsonToObj err"); e != nil {
			return
		}
	} else {
		copyright = []map[string]string{}
	}
	return map[string]interface{}{
		"webName":     website["name"],
		"webLogo":     util.UrlUtil.ToAbsoluteUrl(website["logo"]),
		"webFavicon":  util.UrlUtil.ToAbsoluteUrl(website["favicon"]),
		"webBackdrop": util.UrlUtil.ToAbsoluteUrl(website["backdrop"]),
		"ossDomain":   config.Config.PublicUrl,
		"copyright":   copyright,
	}, nil
}
