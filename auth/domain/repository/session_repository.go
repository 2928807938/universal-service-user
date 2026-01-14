package repository

import (
	"context"

	"universal-service-user/auth/domain/entity"
)

// SessionRepository 会话仓储接口
type SessionRepository interface {
	// Create 创建会话
	Create(ctx context.Context, session *entity.Session) error

	// GetByID 根据 ID 获取会话
	GetByID(ctx context.Context, id int) (*entity.Session, error)

	// GetByRefreshTokenHash 根据 Refresh Token 哈希获取会话
	GetByRefreshTokenHash(ctx context.Context, tokenHash string) (*entity.Session, error)

	// FindByUserID 查询用户的所有会话
	FindByUserID(ctx context.Context, userID int) ([]*entity.Session, error)

	// Update 更新会话
	Update(ctx context.Context, session *entity.Session) error

	// Delete 删除会话
	Delete(ctx context.Context, id int) error

	// DeleteByUserID 删除用户的所有会话
	DeleteByUserID(ctx context.Context, userID int) error

	// DeleteByRefreshTokenHash 根据 Refresh Token 哈希删除会话
	DeleteByRefreshTokenHash(ctx context.Context, tokenHash string) error

	// DeleteExpiredSessions 删除过期会话
	DeleteExpiredSessions(ctx context.Context, expiredBefore int64) error
}
