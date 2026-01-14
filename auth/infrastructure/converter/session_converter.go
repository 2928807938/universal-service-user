package converter

import (
	basegorm "universal-service-user/share/repository/gorm"
	"universal-service-user/auth/domain/entity"
	infraEntity "universal-service-user/auth/infrastructure/entity"
)

// SessionConverter 会话转换器
type SessionConverter struct{}

// NewSessionConverter 创建会话转换器
func NewSessionConverter() *SessionConverter {
	return &SessionConverter{}
}

// ToEntity 将 PO 转换为领域实体
func (c *SessionConverter) ToEntity(po *infraEntity.SessionPO) *entity.Session {
	if po == nil {
		return nil
	}

	return &entity.Session{
		ID:               po.ID,
		UserID:           po.UserID,
		RefreshTokenHash: po.RefreshTokenHash,
		DeviceType:       po.DeviceType,
		DeviceName:       po.DeviceName,
		IP:               po.IP,
		LastActiveAt:     po.LastActiveAt,
		AuditFields: basegorm.AuditFields{
			CreatedAt: po.CreatedAt,
			UpdatedAt: po.UpdatedAt,
			Version:   po.Version,
		},
	}
}

// ToPO 将领域实体转换为 PO
func (c *SessionConverter) ToPO(session *entity.Session) *infraEntity.SessionPO {
	if session == nil {
		return nil
	}

	po := &infraEntity.SessionPO{
		ID:               session.ID,
		UserID:           session.UserID,
		RefreshTokenHash: session.RefreshTokenHash,
		DeviceType:       session.DeviceType,
		DeviceName:       session.DeviceName,
		IP:               session.IP,
		LastActiveAt:     session.LastActiveAt,
	}
	po.CreatedAt = session.CreatedAt
	po.UpdatedAt = session.UpdatedAt
	po.Version = session.Version
	return po
}

// ToEntities 批量转换为领域实体
func (c *SessionConverter) ToEntities(pos []*infraEntity.SessionPO) []*entity.Session {
	if pos == nil {
		return nil
	}

	entities := make([]*entity.Session, 0, len(pos))
	for _, po := range pos {
		entities = append(entities, c.ToEntity(po))
	}
	return entities
}
