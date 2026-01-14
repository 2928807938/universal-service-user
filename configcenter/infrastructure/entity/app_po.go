package entity

import "time"

// AppPO maps to apps table.
type AppPO struct {
	ID          int       `gorm:"primaryKey;autoIncrement"`
	TenantID    string    `gorm:"type:varchar(64);uniqueIndex"`
	AppName     string    `gorm:"type:varchar(128)"`
	Description string    `gorm:"type:varchar(512)"`
	Status      string    `gorm:"type:varchar(32);default:active;index"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (AppPO) TableName() string {
	return "apps"
}
