package adapter

import (
	"context"
	"time"

	"gorm.io/gorm"

	"go-makeadmin/config"
	"go-makeadmin/core"
	"go-makeadmin/makeadmin/repository"
	makeadminsvc "go-makeadmin/makeadmin/service"
	"go-makeadmin/model/makeadmin"
	"go-makeadmin/util"
)

type IndexAdapter interface {
	Console(ctx context.Context) (map[string]interface{}, error)
	Config(ctx context.Context) (map[string]interface{}, error)
}

type indexAdapter struct {
	db *gorm.DB
}

func NewIndexAdapter(db *gorm.DB) IndexAdapter {
	return indexAdapter{db: db}
}

func (adapter indexAdapter) Console(ctx context.Context) (map[string]interface{}, error) {
	website, err := adapter.settingService().WebsiteDetail(ctx, makeadmin.GlobalTenantID)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	date := make([]string, 0, 15)
	for i := 14; i >= 0; i-- {
		date = append(date, now.AddDate(0, 0, -i).Format(core.DateFormat))
	}
	return map[string]interface{}{
		"version": map[string]interface{}{
			"name":    website.Name,
			"version": config.Config.Version,
			"website": "https://www.go-makeadmin.cn",
			"based":   "Go、Gin、Gorm、Vue3、Element Plus、MySQL、Redis",
			"links": map[string]string{
				"gitee":   "https://gitee.com/codepeep/go-makeadmin",
				"website": "https://www.go-makeadmin.cn",
			},
		},
		"today": map[string]interface{}{
			"time":        "2022-08-11 15:08:29",
			"todayVisits": 10,
			"totalVisits": 100,
			"todaySales":  30,
			"totalSales":  65,
			"todayOrder":  12,
			"totalOrder":  255,
			"todayUsers":  120,
			"totalUsers":  360,
		},
		"visitor": map[string]interface{}{
			"date": date,
			"list": []int{12, 13, 11, 5, 8, 22, 14, 9, 456, 62, 78, 12, 18, 22, 46},
		},
	}, nil
}

func (adapter indexAdapter) Config(ctx context.Context) (map[string]interface{}, error) {
	srv := adapter.settingService()
	website, err := srv.WebsiteDetail(ctx, makeadmin.GlobalTenantID)
	if err != nil {
		return nil, err
	}
	copyright, err := srv.CopyrightDetail(ctx, makeadmin.GlobalTenantID)
	if err != nil {
		return nil, err
	}
	copyrightItems := make([]map[string]string, 0, len(copyright))
	for _, item := range copyright {
		copyrightItems = append(copyrightItems, map[string]string{
			"name": item.Name,
			"link": item.Link,
		})
	}
	return map[string]interface{}{
		"webName":     website.Name,
		"webLogo":     util.UrlUtil.ToAbsoluteUrl(website.Logo),
		"webFavicon":  util.UrlUtil.ToAbsoluteUrl(website.Favicon),
		"webBackdrop": util.UrlUtil.ToAbsoluteUrl(website.Backdrop),
		"ossDomain":   config.Config.PublicUrl,
		"copyright":   copyrightItems,
	}, nil
}

func (adapter indexAdapter) settingService() makeadminsvc.SettingService {
	return makeadminsvc.NewSettingService(repository.NewSettingRepository(adapter.db))
}
