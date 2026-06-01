package service

import (
	"go-makeadmin/admin/service/common"
	"go-makeadmin/admin/service/setting"
	"go-makeadmin/admin/service/system"
	makeadminadapter "go-makeadmin/makeadmin/adapter"
)

var InitFunctions = []interface{}{
	// common
	common.NewAlbumService,
	common.NewIndexService,
	common.NewUploadService,
	// setting
	setting.NewSettingCopyrightService,
	setting.NewSettingDictDataService,
	setting.NewSettingDictTypeService,
	setting.NewSettingProtocolService,
	setting.NewSettingStorageService,
	setting.NewSettingWebsiteService,
	// system
	system.NewSystemAuthAdminService,
	system.NewSystemAuthDeptService,
	system.NewSystemAuthMenuService,
	system.NewSystemAuthPermService,
	system.NewSystemAuthPostService,
	system.NewSystemAuthRoleService,
	system.NewSystemLoginService,
	system.NewSystemLogsServer,
	makeadminadapter.NewCopyrightAdapter,
	makeadminadapter.NewDictAdapter,
	makeadminadapter.NewProtocolAdapter,
	makeadminadapter.NewStorageAdapter,
	makeadminadapter.NewSystemAdapter,
	makeadminadapter.NewWebsiteAdapter,
}
