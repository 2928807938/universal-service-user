package converter

import (
	"universal-service-user/api/user-api/dto/vo"
	"universal-service-user/share/utils"
	"universal-service-user/user/domain/entity"
)

// UserConverter 用户转换器
type UserConverter struct{}

// NewUserConverter 创建用户转换器
func NewUserConverter() *UserConverter {
	return &UserConverter{}
}

// ToVo 将领域实体转换为视图对象（包含脱敏）
func (c *UserConverter) ToVo(user *entity.User) *vo.UserVo {
	if user == nil {
		return nil
	}

	// 邮箱脱敏
	maskedEmail := utils.MaskEmail(user.Email.String())

	return &vo.UserVo{
		ID:        user.ID,
		Username:  user.Username,
		Email:     maskedEmail, // 脱敏后的邮箱
		Status:    int(user.Status),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}