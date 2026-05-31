package routers

import (
	"go-makeadmin/admin/routers/common"
	"go-makeadmin/admin/routers/monitor"
	"go-makeadmin/admin/routers/setting"
	"go-makeadmin/admin/routers/system"
	"go-makeadmin/core"
)

var InitRouters = []*core.GroupBase{
	// common
	common.AlbumGroup,
	common.IndexGroup,
	common.UploadGroup,
	// monitor
	monitor.MonitorGroup,
	// setting
	setting.CopyrightGroup,
	setting.DictDataGroup,
	setting.DictTypeGroup,
	setting.ProtocolGroup,
	setting.StorageGroup,
	setting.WebsiteGroup,
	// system
	system.AdminGroup,
	system.DeptGroup,
	system.LogGroup,
	system.LoginGroup,
	system.MenuGroup,
	system.PostGroup,
	system.RoleGroup,
}
