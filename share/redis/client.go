package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/2928807938/universal-service-user/share/config"

	"github.com/redis/go-redis/v9"
)

// Client Redis 客户端封装
type Client struct {
	*redis.Client
	config *config.RedisConfig
}

// NewClient 创建 Redis 客户端
func NewClient(cfg *config.RedisConfig) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.GetRedisAddr(),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis 连接失败: %w", err)
	}

	return &Client{
		Client: client,
		config: cfg,
	}, nil
}

// IncrAndGet 增加计数器并返回当前值
func (c *Client) IncrAndGet(ctx context.Context, key string) (int64, error) {
	result, err := c.Client.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("增加计数器失败: %w", err)
	}
	return result, nil
}

// SetEX 设置键值��指定过期时间（秒）
func (c *Client) SetEX(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if err := c.Client.Set(ctx, key, value, expiration).Err(); err != nil {
		return fmt.Errorf("设置键值失败: %w", err)
	}
	return nil
}

// Get 获取键值
func (c *Client) Get(ctx context.Context, key string) *redis.StringCmd {
	return c.Client.Get(ctx, key)
}

// GetTTL 获取键的剩余过期时间
func (c *Client) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := c.Client.TTL(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("获取 TTL 失败: %w", err)
	}
	return ttl, nil
}

// Del 删除键
func (c *Client) Del(ctx context.Context, keys ...string) error {
	if err := c.Client.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("删除键失败: %w", err)
	}
	return nil
}

// Exists 检查键是否存在
func (c *Client) Exists(ctx context.Context, keys ...string) (bool, error) {
	result, err := c.Client.Exists(ctx, keys...).Result()
	if err != nil {
		return false, fmt.Errorf("检查键存在性失败: %w", err)
	}
	return result > 0, nil
}

// Close 关闭连接
func (c *Client) Close() error {
	return c.Client.Close()
}
