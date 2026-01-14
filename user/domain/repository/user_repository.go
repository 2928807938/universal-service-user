package repository

import (
	"context"

	baseRepo "universal-service-user/share/repository"
	"universal-service-user/user/domain/entity"
)

// UserRepository 用户仓储接口，继承可查询仓储
type UserRepository interface {
	// 继承可查询仓储（包含 CRUD、分页、条件查询等）
	baseRepo.QueryableRepository[entity.User, int]

	// FindByEmail 根据邮箱查找用户
	FindByEmail(ctx context.Context, email string) (*entity.User, error)

	// FindByUsername 根据用户名查找用户
	FindByUsername(ctx context.Context, username string) (*entity.User, error)

	// FindByPhone 根据手机号查找用户
	FindByPhone(ctx context.Context, phone string) (*entity.User, error)

	// ExistsByEmail 检查邮箱是否存在
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// ExistsByUsername 检查用户名是否存在
	ExistsByUsername(ctx context.Context, username string) (bool, error)

	// ExistsByPhone 检查手机号是否存在
	ExistsByPhone(ctx context.Context, phone string) (bool, error)
}
