package converter

import (
	"github.com/2928807938/universal-service-user/user/domain/entity"
	"github.com/2928807938/universal-service-user/user/domain/enum"
	"github.com/2928807938/universal-service-user/user/domain/valueobject"
	infraEntity "github.com/2928807938/universal-service-user/user/infrastructure/entity"

	basegorm "github.com/2928807938/universal-service-user/share/repository/gorm"
)

// UserConverter 用户转换器
type UserConverter struct{}

// NewUserConverter 创建用户转换器
func NewUserConverter() *UserConverter {
	return &UserConverter{}
}

// ToEntity 将 PO 转换为领域实体
func (c *UserConverter) ToEntity(po *infraEntity.UserPO) *entity.User {
	if po == nil {
		return nil
	}

	return &entity.User{
		ID:           po.ID,
		Username:     po.Username,
		Nickname:     po.Nickname,
		Avatar:       po.Avatar,
		Email:        valueobject.Email(po.Email), // 直接类型转换
		Phone:        valueobject.Phone(po.Phone), // 直接类型转换
		PasswordHash: po.PasswordHash,
		Status:       enum.UserStatus(po.Status),
		AuditFields: basegorm.AuditFields{
			CreatedAt: po.CreatedAt,
			UpdatedAt: po.UpdatedAt,
			Version:   po.Version,
		},
	}
}

// ToPO 将领域实体转换为 PO
func (c *UserConverter) ToPO(user *entity.User) *infraEntity.UserPO {
	if user == nil {
		return nil
	}

	po := &infraEntity.UserPO{
		ID:           user.ID,
		Username:     user.Username,
		Nickname:     user.Nickname,
		Avatar:       user.Avatar,
		Email:        user.Email.String(),
		Phone:        user.Phone.String(),
		PasswordHash: user.PasswordHash,
		Status:       int(user.Status),
	}
	po.CreatedAt = user.CreatedAt
	po.UpdatedAt = user.UpdatedAt
	po.Version = user.Version
	return po
}
