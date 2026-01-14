package entity

import (
	"encoding/json"
	"time"
)

// ConfigHistory records config changes.
type ConfigHistory struct {
	ID           int
	TenantID     string
	OldConfig    json.RawMessage
	NewConfig    json.RawMessage
	Version      int
	ChangedBy    string
	ChangeReason string
	CreatedAt    time.Time
}
