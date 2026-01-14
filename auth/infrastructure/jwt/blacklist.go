package jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Blacklist Token 黑名单接口
type Blacklist interface {
	// Add 添加 Token 到黑名单
	Add(ctx context.Context, jti string, ttl time.Duration) error

	// IsBlacklisted 检查 Token 是否在黑名单中
	IsBlacklisted(ctx context.Context, jti string) (bool, error)

	// AddUserBlacklist 添加用户黑名单时间戳（用于批量撤销）
	AddUserBlacklist(ctx context.Context, userID int, revokedAt time.Time) error

	// IsUserBlacklisted 检查用户是否被列入黑名单
	IsUserBlacklisted(ctx context.Context, userID int, tokenIssuedAt time.Time) (bool, error)
}

// RedisBlacklist Redis 实现的 Token 黑名单
type RedisBlacklist struct {
	client *redis.Client
}

// NewRedisBlacklist 创建 Redis 黑名单
func NewRedisBlacklist(client *redis.Client) *RedisBlacklist {
	return &RedisBlacklist{
		client: client,
	}
}

// Add 添加 Token 到黑名单
func (r *RedisBlacklist) Add(ctx context.Context, jti string, ttl time.Duration) error {
	key := r.tokenBlacklistKey(jti)
	return r.client.Set(ctx, key, "1", ttl).Err()
}

// IsBlacklisted 检查 Token 是否在黑名单中
func (r *RedisBlacklist) IsBlacklisted(ctx context.Context, jti string) (bool, error) {
	key := r.tokenBlacklistKey(jti)
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

// AddUserBlacklist 添加用户黑名单时间戳
func (r *RedisBlacklist) AddUserBlacklist(ctx context.Context, userID int, revokedAt time.Time) error {
	key := r.userBlacklistKey(userID)
	// 用户黑名单永久有效（或设置一个较长的过期时间，如 30 天）
	return r.client.Set(ctx, key, revokedAt.Unix(), 30*24*time.Hour).Err()
}

// IsUserBlacklisted 检查用户是否被列入黑名单
func (r *RedisBlacklist) IsUserBlacklisted(ctx context.Context, userID int, tokenIssuedAt time.Time) (bool, error) {
	key := r.userBlacklistKey(userID)

	revokedAtStr, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil // 用户未被列入黑名单
	}
	if err != nil {
		return false, err
	}

	// 将存储的时间戳字符串转换为 int64
	var revokedAt int64
	_, err = fmt.Sscanf(revokedAtStr, "%d", &revokedAt)
	if err != nil {
		return false, err
	}

	// 如果 Token 签发时间早于用户黑名单时间，则该 Token 已被撤销
	return tokenIssuedAt.Unix() < revokedAt, nil
}

// tokenBlacklistKey 生成 Token 黑名单的 Redis Key
func (r *RedisBlacklist) tokenBlacklistKey(jti string) string {
	return fmt.Sprintf("token:blacklist:%s", jti)
}

// userBlacklistKey 生成用户黑名单的 Redis Key
func (r *RedisBlacklist) userBlacklistKey(userID int) string {
	return fmt.Sprintf("user:blacklist:%d", userID)
}
