package entity

import (
	"time"

	"gorm.io/datatypes"
)

// ConfigHistoryPO maps to config_history table.
type ConfigHistoryPO struct {
	ID           int            `gorm:"primaryKey;autoIncrement"`
	TenantID     string         `gorm:"type:varchar(64);index"`
	OldConfig    datatypes.JSON `gorm:"type:json"`
	NewConfig    datatypes.JSON `gorm:"type:json"`
	Version      int            `gorm:"type:int"`
	ChangedBy    string         `gorm:"type:varchar(128)"`
	ChangeReason string         `gorm:"type:varchar(256)"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
}

func (ConfigHistoryPO) TableName() string {
	return "config_history"
}
