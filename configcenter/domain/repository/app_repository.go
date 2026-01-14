package repository

import (
	"context"
	"universal-service-user/configcenter/domain/entity"
)

// AppRepository manages tenant apps.
type AppRepository interface {
	Create(ctx context.Context, app *entity.App) error
	FindByTenantID(ctx context.Context, tenantID string) (*entity.App, error)
	ExistsByTenantID(ctx context.Context, tenantID string) (bool, error)
}
