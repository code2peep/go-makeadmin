package service

import (
	"context"
	"testing"

	"go-makeadmin/model/makeadmin"
)

type fakeSettingRepository struct {
	values map[string]string
	saved  []makeadmin.Setting
}

func (repo *fakeSettingRepository) ListSettingsByGroup(ctx context.Context, tenantID uint64, group string) (map[string]string, error) {
	return repo.values, nil
}

func (repo *fakeSettingRepository) SaveSetting(ctx context.Context, setting makeadmin.Setting) error {
	repo.saved = append(repo.saved, setting)
	return nil
}

func TestWebsiteDetail(t *testing.T) {
	srv := NewSettingService(&fakeSettingRepository{
		values: map[string]string{
			"name":     "go-makeadmin",
			"logo":     "/api/static/backend_logo.png",
			"favicon":  "/api/static/backend_favicon.ico",
			"backdrop": "/api/static/backend_backdrop.png",
		},
	})

	setting, err := srv.WebsiteDetail(context.Background(), makeadmin.GlobalTenantID)
	if err != nil {
		t.Fatalf("WebsiteDetail() error = %v", err)
	}
	if setting.Name != "go-makeadmin" || setting.Logo == "" || setting.Favicon == "" || setting.Backdrop == "" {
		t.Fatalf("WebsiteDetail() = %#v, want complete website setting", setting)
	}
}

func TestSaveWebsite(t *testing.T) {
	repo := &fakeSettingRepository{}
	srv := NewSettingService(repo)

	err := srv.SaveWebsite(context.Background(), makeadmin.GlobalTenantID, WebsiteSetting{
		Name:     "go-makeadmin",
		Logo:     "/api/static/backend_logo.png",
		Favicon:  "/api/static/backend_favicon.ico",
		Backdrop: "/api/static/backend_backdrop.png",
	})
	if err != nil {
		t.Fatalf("SaveWebsite() error = %v", err)
	}
	if len(repo.saved) != 4 {
		t.Fatalf("SaveWebsite() saved %d settings, want 4", len(repo.saved))
	}
	if repo.saved[0].SettingGroup != "website" || repo.saved[0].SettingKey != "name" {
		t.Fatalf("SaveWebsite() first saved item = %#v", repo.saved[0])
	}
}

func TestCopyrightDetailAndSave(t *testing.T) {
	repo := &fakeSettingRepository{
		values: map[string]string{
			"copyright": `[{"name":"go-makeadmin","link":"https://example.com"}]`,
		},
	}
	srv := NewSettingService(repo)

	items, err := srv.CopyrightDetail(context.Background(), makeadmin.GlobalTenantID)
	if err != nil {
		t.Fatalf("CopyrightDetail() error = %v", err)
	}
	if len(items) != 1 || items[0].Name != "go-makeadmin" {
		t.Fatalf("CopyrightDetail() = %#v", items)
	}

	err = srv.SaveCopyright(context.Background(), makeadmin.GlobalTenantID, items)
	if err != nil {
		t.Fatalf("SaveCopyright() error = %v", err)
	}
	if len(repo.saved) != 1 || repo.saved[0].SettingKey != "copyright" || repo.saved[0].ValueType != "json" {
		t.Fatalf("SaveCopyright() saved = %#v", repo.saved)
	}
}

func TestProtocolDetailAndSave(t *testing.T) {
	repo := &fakeSettingRepository{
		values: map[string]string{
			"service": `{"name":"服务协议","content":"service"}`,
			"privacy": `{"name":"隐私协议","content":"privacy"}`,
		},
	}
	srv := NewSettingService(repo)

	setting, err := srv.ProtocolDetail(context.Background(), makeadmin.GlobalTenantID)
	if err != nil {
		t.Fatalf("ProtocolDetail() error = %v", err)
	}
	if setting.Service.Name != "服务协议" || setting.Privacy.Content != "privacy" {
		t.Fatalf("ProtocolDetail() = %#v", setting)
	}

	err = srv.SaveProtocol(context.Background(), makeadmin.GlobalTenantID, setting)
	if err != nil {
		t.Fatalf("SaveProtocol() error = %v", err)
	}
	if len(repo.saved) != 2 || repo.saved[0].SettingKey != "service" || repo.saved[1].SettingKey != "privacy" {
		t.Fatalf("SaveProtocol() saved = %#v", repo.saved)
	}
}

func TestStorageListAndDetail(t *testing.T) {
	repo := &fakeSettingRepository{
		values: map[string]string{
			"default": "local",
			"local":   `{"name":"本地存储"}`,
			"qiniu":   `{"name":"七牛云存储","bucket":"","secretKey":"","accessKey":"","domain":""}`,
			"aliyun":  `{"name":"阿里云OSS","bucket":"","secretKey":"","accessKey":"","domain":""}`,
			"qcloud":  `{"name":"腾讯云OSS","bucket":"","secretKey":"","accessKey":"","domain":"","region":""}`,
		},
	}
	srv := NewSettingService(repo)

	list, err := srv.StorageList(context.Background(), makeadmin.GlobalTenantID)
	if err != nil {
		t.Fatalf("StorageList() error = %v", err)
	}
	if len(list) != 4 || list[0].Alias != StorageAliasLocal || list[0].Status != 1 {
		t.Fatalf("StorageList() = %#v", list)
	}

	qcloud, err := srv.StorageDetail(context.Background(), makeadmin.GlobalTenantID, StorageAliasQcloud)
	if err != nil {
		t.Fatalf("StorageDetail() error = %v", err)
	}
	if qcloud.Alias != StorageAliasQcloud || qcloud.Region != "" {
		t.Fatalf("StorageDetail() qcloud = %#v", qcloud)
	}
}

func TestSaveStorageUpdatesDriverAndDefault(t *testing.T) {
	repo := &fakeSettingRepository{
		values: map[string]string{
			"default": "local",
			"local":   `{"name":"本地存储"}`,
		},
	}
	srv := NewSettingService(repo)

	err := srv.SaveStorage(context.Background(), makeadmin.GlobalTenantID, StorageSetting{
		Alias:     StorageAliasQiniu,
		Status:    1,
		Bucket:    "bucket",
		AccessKey: "ak",
		SecretKey: "sk",
		Domain:    "https://static.example.com",
	})
	if err != nil {
		t.Fatalf("SaveStorage() error = %v", err)
	}
	if len(repo.saved) != 2 {
		t.Fatalf("SaveStorage() saved %d settings, want 2", len(repo.saved))
	}
	if repo.saved[0].SettingKey != StorageAliasQiniu || repo.saved[0].ValueType != "json" {
		t.Fatalf("SaveStorage() driver saved = %#v", repo.saved[0])
	}
	if repo.saved[1].SettingKey != "default" || repo.saved[1].SettingValue != StorageAliasQiniu {
		t.Fatalf("SaveStorage() default saved = %#v", repo.saved[1])
	}
}

func TestChangeStorageOnlyClearsCurrentDefault(t *testing.T) {
	repo := &fakeSettingRepository{
		values: map[string]string{
			"default": StorageAliasQcloud,
		},
	}
	srv := NewSettingService(repo)

	err := srv.ChangeStorage(context.Background(), makeadmin.GlobalTenantID, StorageAliasQiniu, 0)
	if err != nil {
		t.Fatalf("ChangeStorage() error = %v", err)
	}
	if len(repo.saved) != 0 {
		t.Fatalf("ChangeStorage() saved = %#v, want no default change", repo.saved)
	}

	err = srv.ChangeStorage(context.Background(), makeadmin.GlobalTenantID, StorageAliasQcloud, 0)
	if err != nil {
		t.Fatalf("ChangeStorage() current default error = %v", err)
	}
	if len(repo.saved) != 1 || repo.saved[0].SettingKey != "default" || repo.saved[0].SettingValue != "" {
		t.Fatalf("ChangeStorage() current default saved = %#v", repo.saved)
	}
}
