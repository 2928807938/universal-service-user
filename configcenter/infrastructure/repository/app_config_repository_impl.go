package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"universal-service-user/configcenter/domain/entity"
	domainRepo "universal-service-user/configcenter/domain/repository"
	"universal-service-user/configcenter/infrastructure/converter"
	infraEntity "universal-service-user/configcenter/infrastructure/entity"
)

// AppConfigRepositoryImpl implements AppConfigRepository with GORM.
type AppConfigRepositoryImpl struct {
	db        *gorm.DB
	converter *converter.AppConfigConverter
}

func NewAppConfigRepositoryImpl(db *gorm.DB) domainRepo.AppConfigRepository {
	return &AppConfigRepositoryImpl{
		db:        db,
		converter: converter.NewAppConfigConverter(),
	}
}

func (r *AppConfigRepositoryImpl) FindByTenantAndEnv(ctx context.Context, tenantID, environment string) (*entity.AppConfig, error) {
	if tenantID == "" || environment == "" {
		return nil, nil
	}
	var po infraEntity.AppConfigPO
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND environment = ?", tenantID, environment).
		First(&po).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return r.converter.ToEntity(&po), nil
}

func (r *AppConfigRepositoryImpl) Create(ctx context.Context, cfg *entity.AppConfig) error {
	po := r.converter.ToPO(cfg)
	return r.db.WithContext(ctx).Create(po).Error
}

func (r *AppConfigRepositoryImpl) Update(ctx context.Context, cfg *entity.AppConfig) error {
	po := r.converter.ToPO(cfg)
	return r.db.WithContext(ctx).Save(po).Error
}
