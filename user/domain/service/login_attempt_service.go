package service

import (
	"context"
	"fmt"
	"time"

	redisClient "universal-service-user/share/redis"
	userErrors "universal-service-user/user/domain/errors"

	"github.com/redis/go-redis/v9"
)

// LoginAttemptService 登录尝试服务（用于防撞库）
type LoginAttemptService struct {
	redisClient  *redisClient.Client
	maxFailures  int           // 最大失败次数
	lockDuration time.Duration // 锁定时长
}

// NewLoginAttemptService 创建登录尝试服务
func NewLoginAttemptService(redisClient *redisClient.Client, maxFailures int, lockDuration time.Duration) *LoginAttemptService {
	return &LoginAttemptService{
		redisClient:  redisClient,
		maxFailures:  maxFailures,
		lockDuration: lockDuration,
	}
}

// Redis 键名常量
const (
	// 登录失败次数键: login:failed:{username} 或 login:failed:{email}
	loginFailedKeyPrefix = "login:failed:"
	// 账号锁定键: login:locked:{username} 或 login:locked:{email}
	loginLockedKeyPrefix = "login:locked:"
)

// CheckAndRecordFailure 检查并记录登录失败
// 返回: error - 如果是锁定则包含剩余锁定时长
func (s *LoginAttemptService) CheckAndRecordFailure(ctx context.Context, identifier string) error {
	// 1. 先检查账号是否已被锁定
	lockedKey := loginLockedKeyPrefix + identifier
	isLocked, err := s.redisClient.Exists(ctx, lockedKey)
	if err != nil {
		return fmt.Errorf("检查锁定状态失败: %w", err)
	}

	if isLocked {
		// 账号已被锁定，获取剩余锁定时间
		ttl, _ := s.redisClient.GetTTL(ctx, lockedKey)
		return userErrors.NewAccountTempLockedError(ttl)
	}

	// 2. 记录本次登录失败
	failedKey := loginFailedKeyPrefix + identifier
	count, err := s.redisClient.IncrAndGet(ctx, failedKey)
	if err != nil {
		return fmt.Errorf("记录失败次数失败: %w", err)
	}

	// 3. 如果是第一次失败，设置过期时间为锁定时长的2倍（防止一直累计）
	if count == 1 {
		// 设置过期时间为锁定时长的2倍
		if err := s.redisClient.SetEX(ctx, failedKey, count, s.lockDuration*2); err != nil {
			return fmt.Errorf("设置失败次数过期时间失败: %w", err)
		}
	}

	// 4. 检查是否达到最大失败次数
	if count >= int64(s.maxFailures) {
		// 锁定账号
		if err := s.redisClient.SetEX(ctx, lockedKey, time.Now().Unix(), s.lockDuration); err != nil {
			return fmt.Errorf("锁定账号失败: %w", err)
		}

		// 清除失败次数记录（可选）
		if err := s.redisClient.Del(ctx, failedKey); err != nil {
			// 记录日志但不影响主流程
		}

		return userErrors.NewAccountTempLockedError(s.lockDuration)
	}

	// 未达到最大失败次数，返回剩余尝试次数
	remainingAttempts := s.maxFailures - int(count)
	return userErrors.NewLoginAttemptError(remainingAttempts)
}

// ClearFailures 清除登录失败记录（登录成功时调用）
func (s *LoginAttemptService) ClearFailures(ctx context.Context, identifier string) error {
	failedKey := loginFailedKeyPrefix + identifier
	if err := s.redisClient.Del(ctx, failedKey); err != nil {
		return fmt.Errorf("清除失败记录失败: %w", err)
	}
	return nil
}

// IsLocked 检查账号是否被锁定
func (s *LoginAttemptService) IsLocked(ctx context.Context, identifier string) (bool, time.Duration, error) {
	lockedKey := loginLockedKeyPrefix + identifier
	isLocked, err := s.redisClient.Exists(ctx, lockedKey)
	if err != nil {
		return false, 0, fmt.Errorf("检查锁定状态失败: %w", err)
	}

	if !isLocked {
		return false, 0, nil
	}

	// 获取剩余锁定时间
	ttl, err := s.redisClient.GetTTL(ctx, lockedKey)
	if err != nil {
		return true, 0, nil // 获取 TTL 失败，但账号确实被锁定
	}

	return true, ttl, nil
}

// GetRemainingAttempts 获得剩余尝试次数
func (s *LoginAttemptService) GetRemainingAttempts(ctx context.Context, identifier string) (int, error) {
	failedKey := loginFailedKeyPrefix + identifier
	countStr, err := s.redisClient.Get(ctx, failedKey).Result()

	// 键不存在，说明没有失败记录
	if err == redis.Nil {
		return s.maxFailures, nil
	}
	if err != nil {
		// 其他错误，返回默认值
		return s.maxFailures, nil
	}
	if countStr == "" {
		return s.maxFailures, nil
	}

	// Redis 返回的是字符串，需要转换
	var count int64
	if _, err := fmt.Sscanf(countStr, "%d", &count); err != nil {
		return s.maxFailures, nil // 解析失败，返回默认值
	}

	remaining := s.maxFailures - int(count)
	if remaining < 0 {
		remaining = 0
	}

	return remaining, nil
}

// Unlock 手动解锁账号
func (s *LoginAttemptService) Unlock(ctx context.Context, identifier string) error {
	lockedKey := loginLockedKeyPrefix + identifier
	failedKey := loginFailedKeyPrefix + identifier

	// 删除锁定标记
	if err := s.redisClient.Del(ctx, lockedKey); err != nil {
		return fmt.Errorf("解锁账号失败: %w", err)
	}

	// 同时清除失败记录
	if err := s.redisClient.Del(ctx, failedKey); err != nil {
		// 记录日志但不影响主流程
	}

	return nil
}
