package service

import (
	"context"
	"universal-service-user/api/user-api/converter"
	"universal-service-user/api/user-api/dto/request"
	"universal-service-user/api/user-api/dto/vo"
	"universal-service-user/user/domain/enum"
	"universal-service-user/user/domain/repository"
	domainService "universal-service-user/user/domain/service"
	"universal-service-user/user/domain/valueobject"
	verificationEnum "universal-service-user/verification/domain/enum"
	verificationDomain "universal-service-user/verification/domain/service"
)

// UserAppService 用户应用服务
type UserAppService struct {
	userRepo            repository.UserRepository
	userDomainService   *domainService.UserDomainService
	verificationService *verificationDomain.VerificationService
	converter           *converter.UserConverter
}

// NewUserAppService 创建用户应用服务
func NewUserAppService(
	userRepo repository.UserRepository,
	verificationService *verificationDomain.VerificationService,
	loginAttemptService *domainService.LoginAttemptService,
) *UserAppService {
	return &UserAppService{
		userRepo:            userRepo,
		userDomainService:   domainService.NewUserDomainService(userRepo, loginAttemptService),
		verificationService: verificationService,
		converter:           converter.NewUserConverter(),
	}
}

// Register 用户注册
func (s *UserAppService) Register(ctx context.Context, req *request.RegisterRequest) (*vo.UserVo, error) {
	// 1. 确定验证目标（优先使用邮箱）
	target := req.Email
	if target == "" {
		target = req.Phone
	}

	// 2. 验证验证码（错误信息已在领域层定义）
	if err := s.verificationService.VerifyAndDelete(ctx, target, verificationEnum.SceneRegister, req.Code); err != nil {
		return nil, err
	}

	// 3. 创建密码哈希
	password, err := valueobject.NewPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// 4. 调用领域服务注册用户（领域服务会验证"至少提供一个"）
	user, err := s.userDomainService.RegisterUser(ctx, req.Email, req.Phone, password.Hash(), req.Username)
	if err != nil {
		return nil, err
	}

	// 5. 返回用户信息
	return s.converter.ToVo(user), nil
}

// GetUser 获取用户
func (s *UserAppService) GetUser(ctx context.Context, id int) (*vo.UserVo, error) {
	user, err := s.userDomainService.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.converter.ToVo(user), nil
}

// UpdateUser 更新用户
func (s *UserAppService) UpdateUser(ctx context.Context, id int, req *request.UpdateUserRequest) (*vo.UserVo, error) {
	username := ""
	if req.Username != nil {
		username = *req.Username
	}
	var status *enum.UserStatus
	if req.Status != nil {
		s := enum.UserStatus(*req.Status)
		status = &s
	}

	// 调用领域服务更新用户
	user, err := s.userDomainService.UpdateUser(ctx, id, username, status)
	if err != nil {
		return nil, err
	}

	// 保存更新
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return s.converter.ToVo(user), nil
}

// ResetPassword 重置密码(忘记密码)
func (s *UserAppService) ResetPassword(ctx context.Context, req *request.ResetPasswordRequest) error {
	// 1. 确定验证目标（优先使用邮箱）
	target := req.Email
	if target == "" {
		target = req.Phone
	}

	// 2. 验证验证码（错误信息已在领域层定义）
	if err := s.verificationService.VerifyAndDelete(ctx, target, verificationEnum.SceneResetPassword, req.Code); err != nil {
		return err
	}

	// 3. 创建新密码哈希
	newPassword, err := valueobject.NewPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// 4. 调用领域服务重置密码（领域服务会查找用户，用户不存在会返回错误）
	if err := s.userDomainService.ResetPassword(ctx, target, newPassword.Hash()); err != nil {
		return err
	}

	return nil
}
