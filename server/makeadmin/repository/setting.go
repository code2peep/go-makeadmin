package repository

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"go-makeadmin/model/makeadmin"
)

type SettingRepository interface {
	ListSettingsByGroup(ctx context.Context, tenantID uint64, group string) (map[string]string, error)
	SaveSetting(ctx context.Context, setting makeadmin.Setting) error
}

type settingRepository struct {
	db *gorm.DB
}

func NewSettingRepository(db *gorm.DB) SettingRepository {
	return settingRepository{db: db}
}

func (repo settingRepository) ListSettingsByGroup(ctx context.Context, tenantID uint64, group string) (map[string]string, error) {
	var settings []makeadmin.Setting
	err := repo.db.WithContext(ctx).
		Where("tenant_id = ? AND setting_group = ?", tenantID, group).
		Find(&settings).
		Error
	if err != nil {
		return nil, err
	}
	result := make(map[string]string, len(settings))
	for _, setting := range settings {
		result[setting.SettingKey] = setting.SettingValue
	}
	return result, nil
}

func (repo settingRepository) SaveSetting(ctx context.Context, setting makeadmin.Setting) error {
	return repo.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "tenant_id"},
				{Name: "setting_group"},
				{Name: "setting_key"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"setting_value",
				"value_type",
				"is_public",
				"remark",
				"update_time",
			}),
		}).
		Create(&setting).
		Error
}
