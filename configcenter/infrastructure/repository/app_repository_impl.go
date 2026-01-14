package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/2928807938/universal-service-user/configcenter/domain/entity"
	domainRepo "github.com/2928807938/universal-service-user/configcenter/domain/repository"
	"github.com/2928807938/universal-service-user/configcenter/infrastructure/converter"
	infraEntity "github.com/2928807938/universal-service-user/configcenter/infrastructure/entity"
)

// AppRepositoryImpl implements AppRepository with GORM.
type AppRepositoryImpl struct {
	db        *gorm.DB
	converter *converter.AppConverter
}

func NewAppRepositoryImpl(db *gorm.DB) domainRepo.AppRepository {
	return &AppRepositoryImpl{
		db:        db,
		converter: converter.NewAppConverter(),
	}
}

func (r *AppRepositoryImpl) Create(ctx context.Context, app *entity.App) error {
	po := r.converter.ToPO(app)
	return r.db.WithContext(ctx).Create(po).Error
}

func (r *AppRepositoryImpl) FindByTenantID(ctx context.Context, tenantID string) (*entity.App, error) {
	if tenantID == "" {
		return nil, nil
	}
	var po infraEntity.AppPO
	err := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).First(&po).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return r.converter.ToEntity(&po), nil
}

func (r *AppRepositoryImpl) ExistsByTenantID(ctx context.Context, tenantID string) (bool, error) {
	if tenantID == "" {
		return false, nil
	}
	var count int64
	if err := r.db.WithContext(ctx).Model(&infraEntity.AppPO{}).
		Where("tenant_id = ?", tenantID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
