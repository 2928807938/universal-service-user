package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/2928807938/universal-service-user/configcenter/domain/entity"
	domainRepo "github.com/2928807938/universal-service-user/configcenter/domain/repository"
	"github.com/2928807938/universal-service-user/configcenter/infrastructure/converter"
)

// ConfigHistoryRepositoryImpl implements ConfigHistoryRepository with GORM.
type ConfigHistoryRepositoryImpl struct {
	db        *gorm.DB
	converter *converter.ConfigHistoryConverter
}

func NewConfigHistoryRepositoryImpl(db *gorm.DB) domainRepo.ConfigHistoryRepository {
	return &ConfigHistoryRepositoryImpl{
		db:        db,
		converter: converter.NewConfigHistoryConverter(),
	}
}

func (r *ConfigHistoryRepositoryImpl) Create(ctx context.Context, history *entity.ConfigHistory) error {
	po := r.converter.ToPO(history)
	return r.db.WithContext(ctx).Create(po).Error
}
