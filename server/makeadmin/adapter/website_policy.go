package adapter

import (
	"context"

	"gorm.io/gorm"

	"go-makeadmin/admin/schemas/req"
	"go-makeadmin/makeadmin/repository"
	makeadminsvc "go-makeadmin/makeadmin/service"
	"go-makeadmin/model/makeadmin"
)

type CopyrightAdapter interface {
	Available(ctx context.Context) bool
	Detail(ctx context.Context) ([]map[string]interface{}, error)
	Save(ctx context.Context, items []req.SettingCopyrightItemReq) error
}

type ProtocolAdapter interface {
	Available(ctx context.Context) bool
	Detail(ctx context.Context) (map[string]interface{}, error)
	Save(ctx context.Context, protocolReq req.SettingProtocolReq) error
}

type copyrightAdapter struct {
	db *gorm.DB
}

type protocolAdapter struct {
	db *gorm.DB
}

func NewCopyrightAdapter(db *gorm.DB) CopyrightAdapter {
	return copyrightAdapter{db: db}
}

func NewProtocolAdapter(db *gorm.DB) ProtocolAdapter {
	return protocolAdapter{db: db}
}

func (adapter copyrightAdapter) Available(ctx context.Context) bool {
	return settingExists(ctx, adapter.db, "website", "copyright")
}

func (adapter copyrightAdapter) Detail(ctx context.Context) ([]map[string]interface{}, error) {
	items, err := adapter.settingService().CopyrightDetail(ctx, makeadmin.GlobalTenantID)
	if err != nil {
		return nil, err
	}
	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"name": item.Name,
			"link": item.Link,
		})
	}
	return result, nil
}

func (adapter copyrightAdapter) Save(ctx context.Context, items []req.SettingCopyrightItemReq) error {
	copyrightItems := make([]makeadminsvc.CopyrightItem, 0, len(items))
	for _, item := range items {
		copyrightItems = append(copyrightItems, makeadminsvc.CopyrightItem{
			Name: item.Name,
			Link: item.Link,
		})
	}
	return adapter.settingService().SaveCopyright(ctx, makeadmin.GlobalTenantID, copyrightItems)
}

func (adapter copyrightAdapter) settingService() makeadminsvc.SettingService {
	return makeadminsvc.NewSettingService(repository.NewSettingRepository(adapter.db))
}

func (adapter protocolAdapter) Available(ctx context.Context) bool {
	if adapter.db == nil || !adapter.db.Migrator().HasTable(&makeadmin.Setting{}) {
		return false
	}
	var count int64
	err := adapter.db.WithContext(ctx).
		Model(&makeadmin.Setting{}).
		Where("tenant_id = ? AND setting_group = ? AND setting_key IN ?", makeadmin.GlobalTenantID, "protocol", []string{"service", "privacy"}).
		Count(&count).
		Error
	return err == nil && count >= 2
}

func (adapter protocolAdapter) Detail(ctx context.Context) (map[string]interface{}, error) {
	setting, err := adapter.settingService().ProtocolDetail(ctx, makeadmin.GlobalTenantID)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"service": map[string]interface{}{
			"name":    setting.Service.Name,
			"content": setting.Service.Content,
		},
		"privacy": map[string]interface{}{
			"name":    setting.Privacy.Name,
			"content": setting.Privacy.Content,
		},
	}, nil
}

func (adapter protocolAdapter) Save(ctx context.Context, protocolReq req.SettingProtocolReq) error {
	return adapter.settingService().SaveProtocol(ctx, makeadmin.GlobalTenantID, makeadminsvc.ProtocolSetting{
		Service: makeadminsvc.ProtocolItem{
			Name:    protocolReq.Service.Name,
			Content: protocolReq.Service.Content,
		},
		Privacy: makeadminsvc.ProtocolItem{
			Name:    protocolReq.Privacy.Name,
			Content: protocolReq.Privacy.Content,
		},
	})
}

func (adapter protocolAdapter) settingService() makeadminsvc.SettingService {
	return makeadminsvc.NewSettingService(repository.NewSettingRepository(adapter.db))
}

func settingExists(ctx context.Context, db *gorm.DB, group string, key string) bool {
	if db == nil || !db.Migrator().HasTable(&makeadmin.Setting{}) {
		return false
	}
	var count int64
	err := db.WithContext(ctx).
		Model(&makeadmin.Setting{}).
		Where("tenant_id = ? AND setting_group = ? AND setting_key = ?", makeadmin.GlobalTenantID, group, key).
		Count(&count).
		Error
	return err == nil && count > 0
}
