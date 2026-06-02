package adapter

import (
	"context"

	"gorm.io/gorm"

	"go-makeadmin/admin/schemas/req"
	"go-makeadmin/makeadmin/repository"
	makeadminsvc "go-makeadmin/makeadmin/service"
	"go-makeadmin/model/makeadmin"
	"go-makeadmin/util"
)

type WebsiteAdapter interface {
	Available(ctx context.Context) bool
	Detail(ctx context.Context) (map[string]string, error)
	Save(ctx context.Context, websiteReq req.SettingWebsiteReq) error
}

type websiteAdapter struct {
	db *gorm.DB
}

func NewWebsiteAdapter(db *gorm.DB) WebsiteAdapter {
	return websiteAdapter{db: db}
}

func (adapter websiteAdapter) Available(ctx context.Context) bool {
	if adapter.db == nil || !adapter.db.Migrator().HasTable(&makeadmin.Setting{}) {
		return false
	}
	var count int64
	err := adapter.db.WithContext(ctx).
		Model(&makeadmin.Setting{}).
		Where("tenant_id = ? AND setting_group = ? AND setting_key IN ?", tenantIDFromContext(ctx), "website", []string{
			"name",
			"logo",
			"favicon",
			"backdrop",
		}).
		Count(&count).
		Error
	return err == nil && count >= 4
}

func (adapter websiteAdapter) Detail(ctx context.Context) (map[string]string, error) {
	setting, err := adapter.settingService().WebsiteDetail(ctx, tenantIDFromContext(ctx))
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"name":     setting.Name,
		"logo":     util.UrlUtil.ToAbsoluteUrl(setting.Logo),
		"favicon":  util.UrlUtil.ToAbsoluteUrl(setting.Favicon),
		"backdrop": util.UrlUtil.ToAbsoluteUrl(setting.Backdrop),
	}, nil
}

func (adapter websiteAdapter) Save(ctx context.Context, websiteReq req.SettingWebsiteReq) error {
	return adapter.settingService().SaveWebsite(ctx, tenantIDFromContext(ctx), makeadminsvc.WebsiteSetting{
		Name:     websiteReq.Name,
		Logo:     util.UrlUtil.ToRelativeUrl(websiteReq.Logo),
		Favicon:  util.UrlUtil.ToRelativeUrl(websiteReq.Favicon),
		Backdrop: util.UrlUtil.ToRelativeUrl(websiteReq.Backdrop),
	})
}

func (adapter websiteAdapter) settingService() makeadminsvc.SettingService {
	return makeadminsvc.NewSettingService(repository.NewSettingRepository(adapter.db))
}
