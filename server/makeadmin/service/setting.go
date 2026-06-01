package service

import (
	"context"
	"encoding/json"
	"errors"

	"go-makeadmin/makeadmin/repository"
	"go-makeadmin/model/makeadmin"
)

const (
	StorageAliasLocal  = "local"
	StorageAliasQiniu  = "qiniu"
	StorageAliasAliyun = "aliyun"
	StorageAliasQcloud = "qcloud"
)

var ErrUnsupportedStorage = errors.New("makeadmin unsupported storage driver")

type WebsiteSetting struct {
	Name     string
	Logo     string
	Favicon  string
	Backdrop string
}

type CopyrightItem struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

type ProtocolItem struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type ProtocolSetting struct {
	Service ProtocolItem `json:"service"`
	Privacy ProtocolItem `json:"privacy"`
}

type StorageSetting struct {
	Alias     string `json:"alias"`
	Name      string `json:"name"`
	Describe  string `json:"describe"`
	Status    int    `json:"status"`
	Bucket    string `json:"bucket"`
	SecretKey string `json:"secretKey"`
	AccessKey string `json:"accessKey"`
	Domain    string `json:"domain"`
	Region    string `json:"region"`
}

type SettingService interface {
	WebsiteDetail(ctx context.Context, tenantID uint64) (WebsiteSetting, error)
	SaveWebsite(ctx context.Context, tenantID uint64, setting WebsiteSetting) error
	CopyrightDetail(ctx context.Context, tenantID uint64) ([]CopyrightItem, error)
	SaveCopyright(ctx context.Context, tenantID uint64, items []CopyrightItem) error
	ProtocolDetail(ctx context.Context, tenantID uint64) (ProtocolSetting, error)
	SaveProtocol(ctx context.Context, tenantID uint64, setting ProtocolSetting) error
	StorageList(ctx context.Context, tenantID uint64) ([]StorageSetting, error)
	StorageDetail(ctx context.Context, tenantID uint64, alias string) (StorageSetting, error)
	SaveStorage(ctx context.Context, tenantID uint64, setting StorageSetting) error
	ChangeStorage(ctx context.Context, tenantID uint64, alias string, status int) error
}

type settingService struct {
	repo repository.SettingRepository
}

func NewSettingService(repo repository.SettingRepository) SettingService {
	return settingService{repo: repo}
}

func (srv settingService) WebsiteDetail(ctx context.Context, tenantID uint64) (WebsiteSetting, error) {
	settings, err := srv.repo.ListSettingsByGroup(ctx, tenantID, "website")
	if err != nil {
		return WebsiteSetting{}, err
	}
	return WebsiteSetting{
		Name:     settings["name"],
		Logo:     settings["logo"],
		Favicon:  settings["favicon"],
		Backdrop: settings["backdrop"],
	}, nil
}

func (srv settingService) SaveWebsite(ctx context.Context, tenantID uint64, setting WebsiteSetting) error {
	items := []makeadmin.Setting{
		newWebsiteSetting(tenantID, "name", setting.Name, "Website name."),
		newWebsiteSetting(tenantID, "logo", setting.Logo, "Website logo."),
		newWebsiteSetting(tenantID, "favicon", setting.Favicon, "Website favicon."),
		newWebsiteSetting(tenantID, "backdrop", setting.Backdrop, "Login backdrop."),
	}
	for _, item := range items {
		if err := srv.repo.SaveSetting(ctx, item); err != nil {
			return err
		}
	}
	return nil
}

func (srv settingService) CopyrightDetail(ctx context.Context, tenantID uint64) ([]CopyrightItem, error) {
	settings, err := srv.repo.ListSettingsByGroup(ctx, tenantID, "website")
	if err != nil {
		return nil, err
	}
	return decodeCopyright(settings["copyright"])
}

func (srv settingService) SaveCopyright(ctx context.Context, tenantID uint64, items []CopyrightItem) error {
	value, err := json.Marshal(items)
	if err != nil {
		return err
	}
	return srv.repo.SaveSetting(ctx, makeadmin.Setting{
		TenantID:     tenantID,
		SettingGroup: "website",
		SettingKey:   "copyright",
		SettingValue: string(value),
		ValueType:    "json",
		IsPublic:     1,
		Remark:       "Copyright links.",
	})
}

func (srv settingService) ProtocolDetail(ctx context.Context, tenantID uint64) (ProtocolSetting, error) {
	settings, err := srv.repo.ListSettingsByGroup(ctx, tenantID, "protocol")
	if err != nil {
		return ProtocolSetting{}, err
	}
	service, err := decodeProtocolItem(settings["service"])
	if err != nil {
		return ProtocolSetting{}, err
	}
	privacy, err := decodeProtocolItem(settings["privacy"])
	if err != nil {
		return ProtocolSetting{}, err
	}
	return ProtocolSetting{Service: service, Privacy: privacy}, nil
}

func (srv settingService) SaveProtocol(ctx context.Context, tenantID uint64, setting ProtocolSetting) error {
	items := []struct {
		key    string
		value  ProtocolItem
		remark string
	}{
		{key: "service", value: setting.Service, remark: "Service protocol."},
		{key: "privacy", value: setting.Privacy, remark: "Privacy protocol."},
	}
	for _, item := range items {
		value, err := json.Marshal(item.value)
		if err != nil {
			return err
		}
		if err := srv.repo.SaveSetting(ctx, makeadmin.Setting{
			TenantID:     tenantID,
			SettingGroup: "protocol",
			SettingKey:   item.key,
			SettingValue: string(value),
			ValueType:    "json",
			IsPublic:     1,
			Remark:       item.remark,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (srv settingService) StorageList(ctx context.Context, tenantID uint64) ([]StorageSetting, error) {
	settings, err := srv.repo.ListSettingsByGroup(ctx, tenantID, "storage")
	if err != nil {
		return nil, err
	}
	result := make([]StorageSetting, 0, len(storageMetas()))
	for _, meta := range storageMetas() {
		setting, err := decodeStorageSetting(meta, settings[meta.Alias], settings["default"])
		if err != nil {
			return nil, err
		}
		result = append(result, setting)
	}
	return result, nil
}

func (srv settingService) StorageDetail(ctx context.Context, tenantID uint64, alias string) (StorageSetting, error) {
	meta, ok := storageMetaByAlias(alias)
	if !ok {
		return StorageSetting{}, ErrUnsupportedStorage
	}
	settings, err := srv.repo.ListSettingsByGroup(ctx, tenantID, "storage")
	if err != nil {
		return StorageSetting{}, err
	}
	return decodeStorageSetting(meta, settings[alias], settings["default"])
}

func (srv settingService) SaveStorage(ctx context.Context, tenantID uint64, setting StorageSetting) error {
	meta, ok := storageMetaByAlias(setting.Alias)
	if !ok {
		return ErrUnsupportedStorage
	}
	settings, err := srv.repo.ListSettingsByGroup(ctx, tenantID, "storage")
	if err != nil {
		return err
	}
	setting.Name = meta.Name
	setting.Describe = meta.Describe
	if err := srv.saveStorageDriver(ctx, tenantID, setting); err != nil {
		return err
	}
	return srv.saveStorageDefaultIfNeeded(ctx, tenantID, setting.Alias, setting.Status, settings["default"])
}

func (srv settingService) ChangeStorage(ctx context.Context, tenantID uint64, alias string, status int) error {
	if _, ok := storageMetaByAlias(alias); !ok {
		return ErrUnsupportedStorage
	}
	settings, err := srv.repo.ListSettingsByGroup(ctx, tenantID, "storage")
	if err != nil {
		return err
	}
	return srv.saveStorageDefaultIfNeeded(ctx, tenantID, alias, status, settings["default"])
}

func newWebsiteSetting(tenantID uint64, key string, value string, remark string) makeadmin.Setting {
	return makeadmin.Setting{
		TenantID:     tenantID,
		SettingGroup: "website",
		SettingKey:   key,
		SettingValue: value,
		ValueType:    "string",
		IsPublic:     1,
		Remark:       remark,
	}
}

func decodeCopyright(value string) ([]CopyrightItem, error) {
	if value == "" {
		return []CopyrightItem{}, nil
	}
	var items []CopyrightItem
	if err := json.Unmarshal([]byte(value), &items); err != nil {
		return nil, err
	}
	return items, nil
}

func decodeProtocolItem(value string) (ProtocolItem, error) {
	if value == "" {
		return ProtocolItem{}, nil
	}
	var item ProtocolItem
	if err := json.Unmarshal([]byte(value), &item); err != nil {
		return ProtocolItem{}, err
	}
	return item, nil
}

func (srv settingService) saveStorageDriver(ctx context.Context, tenantID uint64, setting StorageSetting) error {
	value, err := json.Marshal(map[string]string{
		"name":      setting.Name,
		"bucket":    setting.Bucket,
		"secretKey": setting.SecretKey,
		"accessKey": setting.AccessKey,
		"domain":    setting.Domain,
		"region":    setting.Region,
	})
	if err != nil {
		return err
	}
	return srv.repo.SaveSetting(ctx, makeadmin.Setting{
		TenantID:     tenantID,
		SettingGroup: "storage",
		SettingKey:   setting.Alias,
		SettingValue: string(value),
		ValueType:    "json",
		IsPublic:     0,
		Remark:       setting.Name + ".",
	})
}

func (srv settingService) saveStorageDefaultIfNeeded(ctx context.Context, tenantID uint64, alias string, status int, currentDefault string) error {
	nextDefault := currentDefault
	if status == 1 {
		nextDefault = alias
	}
	if status == 0 && currentDefault == alias {
		nextDefault = ""
	}
	if nextDefault == currentDefault {
		return nil
	}
	return srv.repo.SaveSetting(ctx, makeadmin.Setting{
		TenantID:     tenantID,
		SettingGroup: "storage",
		SettingKey:   "default",
		SettingValue: nextDefault,
		ValueType:    "string",
		IsPublic:     0,
		Remark:       "Default storage driver.",
	})
}

type storageMeta struct {
	Alias    string
	Name     string
	Describe string
}

func storageMetas() []storageMeta {
	return []storageMeta{
		{Alias: StorageAliasLocal, Name: "本地存储", Describe: "存储在本地服务器"},
		{Alias: StorageAliasQiniu, Name: "七牛云存储", Describe: "存储在七牛云对象存储"},
		{Alias: StorageAliasAliyun, Name: "阿里云存储", Describe: "存储在阿里云 OSS"},
		{Alias: StorageAliasQcloud, Name: "腾讯云存储", Describe: "存储在腾讯云 COS"},
	}
}

func storageMetaByAlias(alias string) (storageMeta, bool) {
	for _, meta := range storageMetas() {
		if meta.Alias == alias {
			return meta, true
		}
	}
	return storageMeta{}, false
}

func decodeStorageSetting(meta storageMeta, value string, defaultAlias string) (StorageSetting, error) {
	setting := StorageSetting{
		Alias:    meta.Alias,
		Name:     meta.Name,
		Describe: meta.Describe,
	}
	if defaultAlias == meta.Alias {
		setting.Status = 1
	}
	if value == "" {
		return setting, nil
	}
	var payload map[string]string
	if err := json.Unmarshal([]byte(value), &payload); err != nil {
		return StorageSetting{}, err
	}
	if payload["name"] != "" {
		setting.Name = payload["name"]
	}
	setting.Bucket = payload["bucket"]
	setting.SecretKey = payload["secretKey"]
	setting.AccessKey = payload["accessKey"]
	setting.Domain = payload["domain"]
	setting.Region = payload["region"]
	return setting, nil
}
