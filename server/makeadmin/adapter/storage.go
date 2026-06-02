package adapter

import (
	"context"

	"gorm.io/gorm"

	"go-makeadmin/admin/schemas/req"
	"go-makeadmin/makeadmin/repository"
	makeadminsvc "go-makeadmin/makeadmin/service"
	"go-makeadmin/model/makeadmin"
)

type StorageAdapter interface {
	Available(ctx context.Context) bool
	List(ctx context.Context) ([]map[string]interface{}, error)
	Detail(ctx context.Context, alias string) (map[string]interface{}, error)
	Edit(ctx context.Context, editReq req.SettingStorageEditReq) error
	Change(ctx context.Context, alias string, status int) error
}

type storageAdapter struct {
	db *gorm.DB
}

func NewStorageAdapter(db *gorm.DB) StorageAdapter {
	return storageAdapter{db: db}
}

func (adapter storageAdapter) Available(ctx context.Context) bool {
	if adapter.db == nil || !adapter.db.Migrator().HasTable(&makeadmin.Setting{}) {
		return false
	}
	var count int64
	err := adapter.db.WithContext(ctx).
		Model(&makeadmin.Setting{}).
		Where("tenant_id = ? AND setting_group = ? AND setting_key IN ?", tenantIDFromContext(ctx), "storage", []string{
			"default",
			makeadminsvc.StorageAliasLocal,
			makeadminsvc.StorageAliasQiniu,
			makeadminsvc.StorageAliasAliyun,
			makeadminsvc.StorageAliasQcloud,
		}).
		Count(&count).
		Error
	return err == nil && count >= 5
}

func (adapter storageAdapter) List(ctx context.Context) ([]map[string]interface{}, error) {
	settings, err := adapter.settingService().StorageList(ctx, tenantIDFromContext(ctx))
	if err != nil {
		return nil, err
	}
	result := make([]map[string]interface{}, 0, len(settings))
	for _, setting := range settings {
		result = append(result, map[string]interface{}{
			"name":     setting.Name,
			"alias":    setting.Alias,
			"describe": setting.Describe,
			"status":   setting.Status,
		})
	}
	return result, nil
}

func (adapter storageAdapter) Detail(ctx context.Context, alias string) (map[string]interface{}, error) {
	setting, err := adapter.settingService().StorageDetail(ctx, tenantIDFromContext(ctx), alias)
	if err != nil {
		return nil, err
	}
	return storageDetailMap(setting), nil
}

func (adapter storageAdapter) Edit(ctx context.Context, editReq req.SettingStorageEditReq) error {
	return adapter.settingService().SaveStorage(ctx, tenantIDFromContext(ctx), makeadminsvc.StorageSetting{
		Alias:     editReq.Alias,
		Status:    editReq.Status,
		Bucket:    editReq.Bucket,
		SecretKey: editReq.SecretKey,
		AccessKey: editReq.AccessKey,
		Domain:    editReq.Domain,
		Region:    editReq.Region,
	})
}

func (adapter storageAdapter) Change(ctx context.Context, alias string, status int) error {
	return adapter.settingService().ChangeStorage(ctx, tenantIDFromContext(ctx), alias, status)
}

func (adapter storageAdapter) settingService() makeadminsvc.SettingService {
	return makeadminsvc.NewSettingService(repository.NewSettingRepository(adapter.db))
}

func storageDetailMap(setting makeadminsvc.StorageSetting) map[string]interface{} {
	return map[string]interface{}{
		"name":      setting.Name,
		"alias":     setting.Alias,
		"status":    setting.Status,
		"bucket":    setting.Bucket,
		"secretKey": setting.SecretKey,
		"accessKey": setting.AccessKey,
		"domain":    setting.Domain,
		"region":    setting.Region,
	}
}
