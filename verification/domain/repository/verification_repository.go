package repository

import (
	"context"
	"time"

	"github.com/2928807938/universal-service-user/verification/domain/enum"
	"github.com/2928807938/universal-service-user/verification/domain/valueobject"
)

// VerificationRepository 验证码仓储接口
type VerificationRepository interface {
	// Save 保存验证码到存储（Redis）
	// target: 邮箱或手机号
	// scene: 验证码场景
	// code: 验证码值对象
	Save(ctx context.Context, target string, scene enum.VerificationScene, code *valueobject.VerificationCode) error

	// Get 获取验证码
	// target: 邮箱或手机号
	// scene: 验证码场景
	Get(ctx context.Context, target string, scene enum.VerificationScene) (*valueobject.VerificationCode, error)

	// Delete 删除验证码（验证成功后调用）
	// target: 邮箱或手机号
	// scene: 验证码场景
	Delete(ctx context.Context, target string, scene enum.VerificationScene) error

	// CheckSendLimit 检查发送频率限制
	// target: 邮箱或手机号
	// scene: 验证码场景
	// interval: 发送间隔（秒）
	// 返回 true 表示可以发送，false 表示发送过于频繁
	CheckSendLimit(ctx context.Context, target string, scene enum.VerificationScene, interval time.Duration) (bool, error)

	// MarkAsVerified 标记验证码为已验证
	// target: 邮箱或手机号
	// scene: 验证码场景
	MarkAsVerified(ctx context.Context, target string, scene enum.VerificationScene) error
}
