package converter

import (
	"encoding/json"

	"gorm.io/datatypes"

	"universal-service-user/configcenter/domain/entity"
	infraEntity "universal-service-user/configcenter/infrastructure/entity"
)

// ConfigHistoryConverter converts between domain and persistence models.
type ConfigHistoryConverter struct{}

func NewConfigHistoryConverter() *ConfigHistoryConverter {
	return &ConfigHistoryConverter{}
}

func (c *ConfigHistoryConverter) ToPO(history *entity.ConfigHistory) *infraEntity.ConfigHistoryPO {
	if history == nil {
		return nil
	}
	return &infraEntity.ConfigHistoryPO{
		ID:           history.ID,
		TenantID:     history.TenantID,
		OldConfig:    datatypes.JSON(history.OldConfig),
		NewConfig:    datatypes.JSON(history.NewConfig),
		Version:      history.Version,
		ChangedBy:    history.ChangedBy,
		ChangeReason: history.ChangeReason,
		CreatedAt:    history.CreatedAt,
	}
}

func (c *ConfigHistoryConverter) ToEntity(po *infraEntity.ConfigHistoryPO) *entity.ConfigHistory {
	if po == nil {
		return nil
	}
	return &entity.ConfigHistory{
		ID:           po.ID,
		TenantID:     po.TenantID,
		OldConfig:    json.RawMessage(po.OldConfig),
		NewConfig:    json.RawMessage(po.NewConfig),
		Version:      po.Version,
		ChangedBy:    po.ChangedBy,
		ChangeReason: po.ChangeReason,
		CreatedAt:    po.CreatedAt,
	}
}
