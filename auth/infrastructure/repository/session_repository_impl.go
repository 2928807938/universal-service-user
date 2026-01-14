package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/2928807938/universal-service-user/auth/domain/entity"
	domainRepo "github.com/2928807938/universal-service-user/auth/domain/repository"
	"github.com/2928807938/universal-service-user/auth/infrastructure/converter"
	infraEntity "github.com/2928807938/universal-service-user/auth/infrastructure/entity"
	"github.com/2928807938/universal-service-user/share/repository"
	basegorm "github.com/2928807938/universal-service-user/share/repository/gorm"
)

// SessionRepositoryImpl 会话仓储实现
type SessionRepositoryImpl struct {
	db        *gorm.DB
	repo      *basegorm.QueryableGormRepository[infraEntity.SessionPO, int]
	converter *converter.SessionConverter
}

// NewSessionRepositoryImpl 创建会话仓储实现
func NewSessionRepositoryImpl(db *gorm.DB) domainRepo.SessionRepository {
	return &SessionRepositoryImpl{
		db:        db,
		repo:      basegorm.NewQueryableGormRepository[infraEntity.SessionPO, int](db),
		converter: converter.NewSessionConverter(),
	}
}

// Create 创建会话
func (r *SessionRepositoryImpl) Create(ctx context.Context, session *entity.Session) error {
	po := r.converter.ToPO(session)
	return r.repo.Create(ctx, po)
}

// GetByID 根据 ID 获取会话
func (r *SessionRepositoryImpl) GetByID(ctx context.Context, id int) (*entity.Session, error) {
	po, err := r.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if po == nil {
		return nil, nil
	}
	return r.converter.ToEntity(po), nil
}

// GetByRefreshTokenHash 根据 Refresh Token 哈希获取会话
func (r *SessionRepositoryImpl) GetByRefreshTokenHash(ctx context.Context, tokenHash string) (*entity.Session, error) {
	poList, err := r.repo.Where(ctx, repository.Eq("refresh_token_hash", tokenHash))
	if err != nil {
		return nil, err
	}
	if len(poList) == 0 {
		return nil, nil
	}
	return r.converter.ToEntity(poList[0]), nil
}

// FindByUserID 查询用户的所有会话
func (r *SessionRepositoryImpl) FindByUserID(ctx context.Context, userID int) ([]*entity.Session, error) {
	poList, err := r.repo.Where(ctx, repository.Eq("user_id", userID))
	if err != nil {
		return nil, err
	}
	return r.converter.ToEntities(poList), nil
}

// Update 更新会话
func (r *SessionRepositoryImpl) Update(ctx context.Context, session *entity.Session) error {
	po := r.converter.ToPO(session)
	return r.repo.Update(ctx, po)
}

// Delete 删除会话
func (r *SessionRepositoryImpl) Delete(ctx context.Context, id int) error {
	return r.repo.Delete(ctx, id)
}

// DeleteByUserID 删除用户的所有会话
func (r *SessionRepositoryImpl) DeleteByUserID(ctx context.Context, userID int) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&infraEntity.SessionPO{}).Error
}

// DeleteByRefreshTokenHash 根据 Refresh Token 哈希删除会话
func (r *SessionRepositoryImpl) DeleteByRefreshTokenHash(ctx context.Context, tokenHash string) error {
	return r.db.WithContext(ctx).
		Where("refresh_token_hash = ?", tokenHash).
		Delete(&infraEntity.SessionPO{}).Error
}

// DeleteExpiredSessions 删除过期会话
func (r *SessionRepositoryImpl) DeleteExpiredSessions(ctx context.Context, expiredBefore int64) error {
	return r.db.WithContext(ctx).
		Where("last_active_at < ?", expiredBefore).
		Delete(&infraEntity.SessionPO{}).Error
}
