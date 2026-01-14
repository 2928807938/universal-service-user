package entity

import "time"

// App represents a tenant application.
type App struct {
	ID          int
	TenantID    string
	AppName     string
	Description string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
