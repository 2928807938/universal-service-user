package converter

import (
	"universal-service-user/configcenter/domain/entity"
	infraEntity "universal-service-user/configcenter/infrastructure/entity"
)

// AppConverter converts between domain and persistence models.
type AppConverter struct{}

func NewAppConverter() *AppConverter {
	return &AppConverter{}
}

func (c *AppConverter) ToPO(app *entity.App) *infraEntity.AppPO {
	if app == nil {
		return nil
	}
	return &infraEntity.AppPO{
		ID:          app.ID,
		TenantID:    app.TenantID,
		AppName:     app.AppName,
		Description: app.Description,
		Status:      app.Status,
		CreatedAt:   app.CreatedAt,
		UpdatedAt:   app.UpdatedAt,
	}
}

func (c *AppConverter) ToEntity(po *infraEntity.AppPO) *entity.App {
	if po == nil {
		return nil
	}
	return &entity.App{
		ID:          po.ID,
		TenantID:    po.TenantID,
		AppName:     po.AppName,
		Description: po.Description,
		Status:      po.Status,
		CreatedAt:   po.CreatedAt,
		UpdatedAt:   po.UpdatedAt,
	}
}
