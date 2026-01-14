package repository

import (
	"context"
	"github.com/2928807938/universal-service-user/configcenter/domain/entity"
)

// AppConfigRepository manages tenant configs.
type AppConfigRepository interface {
	FindByTenantAndEnv(ctx context.Context, tenantID, environment string) (*entity.AppConfig, error)
	Create(ctx context.Context, cfg *entity.AppConfig) error
	Update(ctx context.Context, cfg *entity.AppConfig) error
}
