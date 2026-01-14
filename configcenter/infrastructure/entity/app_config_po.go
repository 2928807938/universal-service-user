package entity

import (
	"time"

	"gorm.io/datatypes"
)

// AppConfigPO maps to app_configs table.
type AppConfigPO struct {
	ID          int            `gorm:"primaryKey;autoIncrement"`
	TenantID    string         `gorm:"type:varchar(64);index:idx_app_config_tenant_env,unique"`
	Environment string         `gorm:"type:varchar(32);index:idx_app_config_tenant_env,unique"`
	ConfigData  datatypes.JSON `gorm:"type:json"`
	Version     int            `gorm:"type:int;default:1"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
}

func (AppConfigPO) TableName() string {
	return "app_configs"
}
