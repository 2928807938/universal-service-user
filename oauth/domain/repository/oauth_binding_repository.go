package repository

import (
	"context"
	"github.com/2928807938/universal-service-user/oauth/domain/entity"
	"github.com/2928807938/universal-service-user/share/repository"
)

// OAuthBindingRepository OAuth 绑定仓储接口
type OAuthBindingRepository interface {
	repository.BaseRepository[entity.OAuthBinding, int]

	// FindByUserIDAndProvider 根据用户 ID 和平台查询绑定
	// ctx: 上下文
	// userID: 用户 ID
	// provider: 平台标识
	// 返回绑定实体和错误信息
	FindByUserIDAndProvider(ctx context.Context, userID int, provider string) (*entity.OAuthBinding, error)

	// FindByProviderAndOpenID 根据平台和 OpenID 查询绑定
	// ctx: 上下文
	// provider: 平台标识
	// openID: 平台用户唯一标识
	// 返回绑定实体和错误信息
	FindByProviderAndOpenID(ctx context.Context, provider, openID string) (*entity.OAuthBinding, error)

	// FindByProviderAndUnionID 根据平台和 UnionID 查询绑定
	// ctx: 上下文
	// provider: 平台标识
	// unionID: 平台联合标识
	// 返回绑定实体和错误信息
	FindByProviderAndUnionID(ctx context.Context, provider, unionID string) (*entity.OAuthBinding, error)

	// FindByUserID 根据用户 ID 查询所有绑定
	// ctx: 上下文
	// userID: 用户 ID
	// 返回绑定实体列表和错误信息
	FindByUserID(ctx context.Context, userID int) ([]*entity.OAuthBinding, error)

	// DeleteByUserIDAndProvider 根据用户 ID 和平台删除绑定
	// ctx: 上下文
	// userID: 用户 ID
	// provider: 平台标识
	// 返回错误信息
	DeleteByUserIDAndProvider(ctx context.Context, userID int, provider string) error

	// ExistsByProviderAndOpenID 判断绑定是否存在
	// ctx: 上下文
	// provider: 平台标识
	// openID: 平台用户唯一标识
	// 返回是否存在和错误信息
	ExistsByProviderAndOpenID(ctx context.Context, provider, openID string) (bool, error)
}
