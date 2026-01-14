package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/2928807938/universal-service-user/verification/domain/enum"
	verrors "github.com/2928807938/universal-service-user/verification/domain/errors"
	"github.com/2928807938/universal-service-user/verification/domain/repository"
	"github.com/2928807938/universal-service-user/verification/domain/valueobject"
)

// RedisVerificationRepository Redis 验证码仓储实现
type RedisVerificationRepository struct {
	client *redis.Client
}

// NewRedisVerificationRepository 创建 Redis 验证码仓储
func NewRedisVerificationRepository(client *redis.Client) repository.VerificationRepository {
	return &RedisVerificationRepository{
		client: client,
	}
}

// buildKey 构建 Redis Key
// 格式：verification_code:{type}:{target}:{scene}
// type: email 或 phone
func (r *RedisVerificationRepository) buildKey(target string, scene enum.VerificationScene) string {
	// 判断是邮箱还是手机号
	targetType := "email"
	if !strings.Contains(target, "@") {
		targetType = "phone"
	}
	return fmt.Sprintf("verification_code:%s:%s:%s", targetType, target, scene.String())
}

// buildSendLimitKey 构建发送频率限制 Key
// 格式：verification_send_limit:{type}:{target}:{scene}
func (r *RedisVerificationRepository) buildSendLimitKey(target string, scene enum.VerificationScene) string {
	targetType := "email"
	if !strings.Contains(target, "@") {
		targetType = "phone"
	}
	return fmt.Sprintf("verification_send_limit:%s:%s:%s", targetType, target, scene.String())
}

// Save 保存验证码到 Redis
func (r *RedisVerificationRepository) Save(ctx context.Context, target string, scene enum.VerificationScene, code *valueobject.VerificationCode) error {
	key := r.buildKey(target, scene)

	// 序列化验证码对象
	data, err := json.Marshal(code)
	if err != nil {
		return verrors.WrapVerificationError(verrors.VerificationStoreFailed, "序列化验证码失败", err)
	}

	// 保存到 Redis，设置 TTL
	ttl := time.Until(code.ExpiresAt)
	if ttl <= 0 {
		return verrors.ErrCodeExpired
	}

	if err := r.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return verrors.WrapVerificationError(verrors.VerificationStoreFailed, "保存验证码到Redis失败", err)
	}

	// 设置发送频率限制标记
	sendLimitKey := r.buildSendLimitKey(target, scene)
	if err := r.client.Set(ctx, sendLimitKey, "1", ttl).Err(); err != nil {
		// 发送限制设置失败不影响验证码保存
		return nil
	}

	return nil
}

// Get 获取验证码
func (r *RedisVerificationRepository) Get(ctx context.Context, target string, scene enum.VerificationScene) (*valueobject.VerificationCode, error) {
	key := r.buildKey(target, scene)

	// 从 Redis 获取数据
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 验证码不存在
		}
		return nil, verrors.WrapVerificationError(verrors.VerificationStoreFailed, "从Redis获取验证码失败", err)
	}

	// 反序列化
	var code valueobject.VerificationCode
	if err := json.Unmarshal([]byte(data), &code); err != nil {
		return nil, verrors.WrapVerificationError(verrors.VerificationStoreFailed, "反序列化验证码失败", err)
	}

	return &code, nil
}

// Delete 删除验证码
func (r *RedisVerificationRepository) Delete(ctx context.Context, target string, scene enum.VerificationScene) error {
	key := r.buildKey(target, scene)

	if err := r.client.Del(ctx, key).Err(); err != nil {
		return verrors.WrapVerificationError(verrors.VerificationStoreFailed, "删除验证码失败", err)
	}

	return nil
}

// CheckSendLimit 检查发送频率限制
func (r *RedisVerificationRepository) CheckSendLimit(ctx context.Context, target string, scene enum.VerificationScene, interval time.Duration) (bool, error) {
	sendLimitKey := r.buildSendLimitKey(target, scene)

	// 检查是否存在限制标记
	exists, err := r.client.Exists(ctx, sendLimitKey).Result()
	if err != nil {
		return false, verrors.WrapVerificationError(verrors.VerificationStoreFailed, "检查发送频率限制失败", err)
	}

	// 如果存在限制标记，说明发送过于频繁
	if exists > 0 {
		return false, nil
	}

	return true, nil
}

// MarkAsVerified 标记验证码为已验证
func (r *RedisVerificationRepository) MarkAsVerified(ctx context.Context, target string, scene enum.VerificationScene) error {
	key := r.buildKey(target, scene)

	// 获取当前验证码
	code, err := r.Get(ctx, target, scene)
	if err != nil {
		return err
	}
	if code == nil {
		return verrors.ErrCodeNotFound
	}

	// 标记为已验证
	code.MarkAsVerified()

	// 重新保存
	data, err := json.Marshal(code)
	if err != nil {
		return verrors.WrapVerificationError(verrors.VerificationStoreFailed, "序列化验证码失败", err)
	}

	// 更新 Redis，保持原有的 TTL
	ttl := time.Until(code.ExpiresAt)
	if ttl <= 0 {
		return verrors.ErrCodeExpired
	}

	if err := r.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return verrors.WrapVerificationError(verrors.VerificationStoreFailed, "更新验证码状态失败", err)
	}

	return nil
}
