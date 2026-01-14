package converter

import (
	basegorm "universal-service-user/share/repository/gorm"
	"universal-service-user/oauth/domain/entity"
	infraEntity "universal-service-user/oauth/infrastructure/entity"
)

// OAuthBindingConverter OAuth 绑定转换器
type OAuthBindingConverter struct{}

// NewOAuthBindingConverter 创建 OAuth 绑定转换器
func NewOAuthBindingConverter() *OAuthBindingConverter {
	return &OAuthBindingConverter{}
}

// ToEntity 将 PO 转换为领域实体
func (c *OAuthBindingConverter) ToEntity(po *infraEntity.OAuthBindingPO) *entity.OAuthBinding {
	if po == nil {
		return nil
	}

	return &entity.OAuthBinding{
		ID:           po.ID,
		UserID:       po.UserID,
		Provider:     po.Provider,
		OpenID:       po.OpenID,
		UnionID:      po.UnionID,
		Nickname:     po.Nickname,
		Avatar:       po.Avatar,
		AccessToken:  po.AccessToken,
		RefreshToken: po.RefreshToken,
		ExpiresAt:    po.ExpiresAt,
		AuditFields: basegorm.AuditFields{
			CreatedAt: po.CreatedAt,
			UpdatedAt: po.UpdatedAt,
			Version:   po.Version,
		},
	}
}

// ToPO 将领域实体转换为 PO
func (c *OAuthBindingConverter) ToPO(binding *entity.OAuthBinding) *infraEntity.OAuthBindingPO {
	if binding == nil {
		return nil
	}

	po := &infraEntity.OAuthBindingPO{
		ID:           binding.ID,
		UserID:       binding.UserID,
		Provider:     binding.Provider,
		OpenID:       binding.OpenID,
		UnionID:      binding.UnionID,
		Nickname:     binding.Nickname,
		Avatar:       binding.Avatar,
		AccessToken:  binding.AccessToken,
		RefreshToken: binding.RefreshToken,
		ExpiresAt:    binding.ExpiresAt,
	}
	po.CreatedAt = binding.CreatedAt
	po.UpdatedAt = binding.UpdatedAt
	po.Version = binding.Version
	return po
}

// ToEntities 将 PO 列表转换为领域实体列表
func (c *OAuthBindingConverter) ToEntities(pos []*infraEntity.OAuthBindingPO) []*entity.OAuthBinding {
	if pos == nil {
		return nil
	}

	entities := make([]*entity.OAuthBinding, 0, len(pos))
	for _, po := range pos {
		entities = append(entities, c.ToEntity(po))
	}
	return entities
}

// ToPOs 将领域实体列表转换为 PO 列表
func (c *OAuthBindingConverter) ToPOs(entities []*entity.OAuthBinding) []*infraEntity.OAuthBindingPO {
	if entities == nil {
		return nil
	}

	pos := make([]*infraEntity.OAuthBindingPO, 0, len(entities))
	for _, entity := range entities {
		pos = append(pos, c.ToPO(entity))
	}
	return pos
}
