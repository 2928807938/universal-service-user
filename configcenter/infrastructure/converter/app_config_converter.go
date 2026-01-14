package converter

import (
	"encoding/json"

	"gorm.io/datatypes"

	"universal-service-user/configcenter/domain/entity"
	infraEntity "universal-service-user/configcenter/infrastructure/entity"
)

// AppConfigConverter converts between domain and persistence models.
type AppConfigConverter struct{}

func NewAppConfigConverter() *AppConfigConverter {
	return &AppConfigConverter{}
}

func (c *AppConfigConverter) ToPO(cfg *entity.AppConfig) *infraEntity.AppConfigPO {
	if cfg == nil {
		return nil
	}
	return &infraEntity.AppConfigPO{
		ID:          cfg.ID,
		TenantID:    cfg.TenantID,
		Environment: cfg.Environment,
		ConfigData:  datatypes.JSON(cfg.ConfigData),
		Version:     cfg.Version,
		CreatedAt:   cfg.CreatedAt,
		UpdatedAt:   cfg.UpdatedAt,
	}
}

func (c *AppConfigConverter) ToEntity(po *infraEntity.AppConfigPO) *entity.AppConfig {
	if po == nil {
		return nil
	}
	return &entity.AppConfig{
		ID:          po.ID,
		TenantID:    po.TenantID,
		Environment: po.Environment,
		ConfigData:  json.RawMessage(po.ConfigData),
		Version:     po.Version,
		CreatedAt:   po.CreatedAt,
		UpdatedAt:   po.UpdatedAt,
	}
}
