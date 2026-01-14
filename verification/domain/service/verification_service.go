package service

import (
	"context"
	"time"

	"github.com/2928807938/universal-service-user/verification/domain/enum"
	verrors "github.com/2928807938/universal-service-user/verification/domain/errors"
	"github.com/2928807938/universal-service-user/verification/domain/repository"
	"github.com/2928807938/universal-service-user/verification/domain/valueobject"
)

// VerificationConfig 验证码配置
type VerificationConfig struct {
	CodeLength     int           // 验证码长度（默认 6 位）
	ExpireDuration time.Duration // 过期时间（默认 5 分钟）
	SendInterval   time.Duration // 发送间隔（默认 60 秒）
}

// DefaultConfig 默认配置
func DefaultConfig() *VerificationConfig {
	return &VerificationConfig{
		CodeLength:     6,
		ExpireDuration: 5 * time.Minute,
		SendInterval:   60 * time.Second,
	}
}

// VerificationService 验证码领域服务
type VerificationService struct {
	repo   repository.VerificationRepository
	config *VerificationConfig
}

// NewVerificationService 创建验证码服务
func NewVerificationService(repo repository.VerificationRepository, config *VerificationConfig) *VerificationService {
	if config == nil {
		config = DefaultConfig()
	}
	return &VerificationService{
		repo:   repo,
		config: config,
	}
}

// Generate 生成验证码
// target: 邮箱或手机号
// scene: 验证码场景
// 返回生成的验证码
func (s *VerificationService) Generate(ctx context.Context, target string, scene enum.VerificationScene) (*valueobject.VerificationCode, error) {
	// 1. 检查发送频率限制
	canSend, err := s.repo.CheckSendLimit(ctx, target, scene, s.config.SendInterval)
	if err != nil {
		return nil, verrors.WrapVerificationError(verrors.VerificationStoreFailed, "检查发送频率失败", err)
	}
	if !canSend {
		return nil, verrors.ErrCodeTooFrequent
	}

	// 2. 生成验证码
	code, err := valueobject.NewVerificationCode(target, scene, s.config.CodeLength, s.config.ExpireDuration)
	if err != nil {
		return nil, err
	}

	// 3. 保存验证码到 Redis
	if err := s.repo.Save(ctx, target, scene, code); err != nil {
		return nil, verrors.WrapVerificationError(verrors.VerificationStoreFailed, "保存验证码失败", err)
	}

	return code, nil
}

// Verify 验证验证码
// target: 邮箱或手机号
// scene: 验证码场景
// code: 用户输入的验证码
// 返回 error 表示验证失败
func (s *VerificationService) Verify(ctx context.Context, target string, scene enum.VerificationScene, code string) error {
	// 1. 从 Redis 获取验证码
	savedCode, err := s.repo.Get(ctx, target, scene)
	if err != nil {
		return err
	}
	if savedCode == nil {
		return verrors.ErrCodeNotFound
	}

	// 2. 验证验证码
	if err := savedCode.Verify(code); err != nil {
		return err
	}

	// 3. 标记为已验证
	if err := s.repo.MarkAsVerified(ctx, target, scene); err != nil {
		return verrors.WrapVerificationError(verrors.VerificationStoreFailed, "标记验证码为已验证失败", err)
	}

	return nil
}

// VerifyAndDelete 验证验证码并删除
// target: 邮箱或手机号
// scene: 验证码场景
// code: 用户输入的验证码
// 验证成功后会自动删除验证码
func (s *VerificationService) VerifyAndDelete(ctx context.Context, target string, scene enum.VerificationScene, code string) error {
	// 1. 验证验证码
	if err := s.Verify(ctx, target, scene, code); err != nil {
		return err
	}

	// 2. 删除验证码
	if err := s.repo.Delete(ctx, target, scene); err != nil {
		// 删除失败不影响验证结果，只记录错误
		// 因为 Redis 会自动过期
		return nil
	}

	return nil
}

// CheckExists 检查验证码是否存在
// target: 邮箱或手机号
// scene: 验证码场景
func (s *VerificationService) CheckExists(ctx context.Context, target string, scene enum.VerificationScene) (bool, error) {
	code, err := s.repo.Get(ctx, target, scene)
	if err != nil {
		return false, err
	}
	return code != nil && !code.IsExpired(), nil
}
