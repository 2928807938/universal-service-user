package entity

import (
	"encoding/json"
	"time"
)

// AppConfig stores tenant config data and version.
type AppConfig struct {
	ID          int
	TenantID    string
	Environment string
	ConfigData  json.RawMessage
	Version     int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
