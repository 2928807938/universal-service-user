package service

import (
	"context"

	"universal-service-user/api/user-api/dto/request"
	"universal-service-user/api/user-api/dto/vo"
	notificationDomain "universal-service-user/notification/domain/service"
	verificationEnum "universal-service-user/verification/domain/enum"
	verificationErrors "universal-service-user/verification/domain/errors"
	verificationDomain "universal-service-user/verification/domain/service"
)

// VerificationAppService 验证码应用服务
type VerificationAppService struct {
	verificationService *verificationDomain.VerificationService
	notificationService *notificationDomain.NotificationService
}

// NewVerificationAppService 创建验证码应用服务
func NewVerificationAppService(
	verificationService *verificationDomain.VerificationService,
	notificationService *notificationDomain.NotificationService,
) *VerificationAppService {
	return &VerificationAppService{
		verificationService: verificationService,
		notificationService: notificationService,
	}
}

// SendCode 发送验证码
func (s *VerificationAppService) SendCode(ctx context.Context, req *request.SendCodeRequest) (*vo.SendCodeResponse, error) {
	// 1. 转换场景
	scene, err := s.parseScene(req.Scene)
	if err != nil {
		return nil, err
	}

	// 2. 生成验证码（错误信息已在领域层定义）
	code, err := s.verificationService.Generate(ctx, req.Target, scene)
	if err != nil {
		return nil, err
	}

	// 3. 发送验证码
	switch req.Type {
	case "email":
		err = s.notificationService.SendEmailCode(ctx, req.Target, code.Code, scene)
	case "phone":
		err = s.notificationService.SendSMSCode(ctx, req.Target, code.Code, scene)
	default:
		return nil, verificationErrors.NewVerificationError(verificationErrors.VerificationTargetInvalid, "不支持的发送类型: "+req.Type)
	}

	if err != nil {
		return nil, verificationErrors.WrapVerificationError(verificationErrors.VerificationGenerateFailed, "发送验证码失败", err)
	}

	return &vo.SendCodeResponse{
		Message: "验证码已发送,有效期5分钟",
	}, nil
}

// VerifyCode 验证验证码
func (s *VerificationAppService) VerifyCode(ctx context.Context, req *request.VerifyCodeRequest) (*vo.VerifyCodeResponse, error) {
	scene, err := s.parseScene(req.Scene)
	if err != nil {
		return nil, err
	}

	err = s.verificationService.Verify(ctx, req.Target, scene, req.Code)
	if err != nil {
		return &vo.VerifyCodeResponse{
			Valid:   false,
			Message: err.Error(),
		}, nil
	}

	return &vo.VerifyCodeResponse{
		Valid:   true,
		Message: "验证码正确",
	}, nil
}

// parseScene 解析场景
func (s *VerificationAppService) parseScene(sceneStr string) (verificationEnum.VerificationScene, error) {
	scene := verificationEnum.VerificationScene(sceneStr)
	if !scene.IsValid() {
		return "", verificationErrors.NewVerificationError(verificationErrors.VerificationSceneInvalid, "无效的验证码场景: "+sceneStr)
	}
	return scene, nil
}
