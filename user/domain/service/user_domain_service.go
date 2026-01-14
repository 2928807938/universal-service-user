package service

import (
	"context"
	"fmt"
	"time"

	"universal-service-user/rules"
	basegorm "universal-service-user/share/repository/gorm"
	"universal-service-user/user/domain/entity"
	"universal-service-user/user/domain/enum"
	"universal-service-user/user/domain/errors"
	"universal-service-user/user/domain/repository"
	"universal-service-user/user/domain/valueobject"
)

// UserDomainService 用户领域服务
type UserDomainService struct {
	userRepo            repository.UserRepository
	loginAttemptService *LoginAttemptService
}

// NewUserDomainService 创建用户领域服务
func NewUserDomainService(userRepo repository.UserRepository, loginAttemptService *LoginAttemptService) *UserDomainService {
	return &UserDomainService{
		userRepo:            userRepo,
		loginAttemptService: loginAttemptService,
	}
}

// createUserRequest 用于验证的请求结构
type createUserRequest struct {
	Username     string
	Email        string
	PasswordHash string
}

// GetUser 获取用户（包含业务规则校验）
func (s *UserDomainService) GetUser(ctx context.Context, id int) (*entity.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.ErrUserNotFound
	}
	return user, nil
}

// ========== 登录相关 ==========

// LoginByUsername 通过用户名和密码登录
func (s *UserDomainService) LoginByUsername(ctx context.Context, username string, password string) (*entity.User, error) {
	// 验证参数
	result := rules.ForField("Username").Required().Length(1, 50).
		ForField("Password").Required().
		Validate(map[string]any{
			"Username": username,
			"Password": password,
		})

	if !result.IsValid() {
		return nil, errors.NewValidationError(result.Error().Error())
	}

	// 检查登录失败次数（防撞库）
	if s.loginAttemptService != nil {
		// 先检查是否已被锁定
		isLocked, ttl, err := s.loginAttemptService.IsLocked(ctx, username)
		if err != nil {
			return nil, fmt.Errorf("检查锁定状态失败: %w", err)
		}
		if isLocked {
			return nil, errors.NewAccountTempLockedError(ttl)
		}
	}

	// 查找用户
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		// 记录登录失败（用户不存在也算失败，防止枚举用户名）
		if s.loginAttemptService != nil {
			if recordErr := s.loginAttemptService.CheckAndRecordFailure(ctx, username); recordErr != nil {
				return nil, recordErr
			}
		}
		return nil, errors.ErrUserInvalidPassword // 统一返回密码错误，不暴露用户是否存在
	}

	// 检查用户状态
	if err := s.checkUserStatus(user); err != nil {
		return nil, err
	}

	// 验证密码
	pwd, err := valueobject.NewPasswordFromHash(user.PasswordHash)
	if err != nil {
		return nil, err
	}

	if !pwd.Verify(password) {
		// 记录登录失败
		if s.loginAttemptService != nil {
			if recordErr := s.loginAttemptService.CheckAndRecordFailure(ctx, username); recordErr != nil {
				return nil, recordErr
			}
		}
		return nil, errors.ErrUserInvalidPassword
	}

	// 登录成功，清除失败记录
	if s.loginAttemptService != nil {
		if err := s.loginAttemptService.ClearFailures(ctx, username); err != nil {
			// 记录日志但不影响登录
		}
	}

	return user, nil
}

// LoginByEmail 通过邮箱和密码登录
func (s *UserDomainService) LoginByEmail(ctx context.Context, emailStr string, password string) (*entity.User, error) {
	// 验证参数
	result := rules.ForField("Email").Required().Email().
		ForField("Password").Required().
		Validate(map[string]any{
			"Email":    emailStr,
			"Password": password,
		})

	if !result.IsValid() {
		return nil, errors.NewValidationError(result.Error().Error())
	}

	// 检查登录失败次数（防撞库）
	if s.loginAttemptService != nil {
		// 先检查是否已被锁定
		isLocked, ttl, err := s.loginAttemptService.IsLocked(ctx, emailStr)
		if err != nil {
			return nil, fmt.Errorf("检查锁定状态失败: %w", err)
		}
		if isLocked {
			return nil, errors.NewAccountTempLockedError(ttl)
		}
	}

	// 查找用户
	user, err := s.userRepo.FindByEmail(ctx, emailStr)
	if err != nil {
		return nil, err
	}
	if user == nil {
		// 记录登录失败（��户不存在也算失败，防止枚举邮箱）
		if s.loginAttemptService != nil {
			if recordErr := s.loginAttemptService.CheckAndRecordFailure(ctx, emailStr); recordErr != nil {
				return nil, recordErr
			}
		}
		return nil, errors.ErrUserInvalidPassword // 统一返回密码错误
	}

	// 检查用户状态
	if err := s.checkUserStatus(user); err != nil {
		return nil, err
	}

	// 验证密码
	pwd, err := valueobject.NewPasswordFromHash(user.PasswordHash)
	if err != nil {
		return nil, err
	}

	if !pwd.Verify(password) {
		// 记录登录失败
		if s.loginAttemptService != nil {
			if recordErr := s.loginAttemptService.CheckAndRecordFailure(ctx, emailStr); recordErr != nil {
				return nil, recordErr
			}
		}
		return nil, errors.ErrUserInvalidPassword
	}

	// 登录成功，清除失败记录
	if s.loginAttemptService != nil {
		if err := s.loginAttemptService.ClearFailures(ctx, emailStr); err != nil {
			// 记录日志但不影响登录
		}
	}

	return user, nil
}

// LoginByPhone 通过手机号和验证码登录（验证码验证由外层处理）
func (s *UserDomainService) LoginByPhone(ctx context.Context, phone string) (*entity.User, error) {
	// 验证参数
	result := rules.ForField("Phone").Required().Phone().
		Validate(map[string]any{
			"Phone": phone,
		})

	if !result.IsValid() {
		return nil, errors.NewValidationError(result.Error().Error())
	}

	// 查找用户
	user, err := s.userRepo.FindByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.ErrUserNotFound
	}

	// 检查用户状态
	if err := s.checkUserStatus(user); err != nil {
		return nil, err
	}

	return user, nil
}

// checkUserStatus 检查用户状态是否允许登录
func (s *UserDomainService) checkUserStatus(user *entity.User) error {
	if user.IsDisabled() {
		return errors.ErrUserDisabled
	}

	if user.IsLocked() {
		return errors.ErrUserLocked
	}

	if !user.CanLogin() {
		return errors.ErrUserStatusAbnormal
	}

	return nil
}

// ========== 密码管理 ==========

// ChangePassword 修改密码（已登录用户）
func (s *UserDomainService) ChangePassword(ctx context.Context, userID int, oldPassword string, newPassword string) error {
	// 验证参数
	result := rules.ForField("OldPassword").Required().
		ForField("NewPassword").Required().Length(6, 32).
		Validate(map[string]any{
			"OldPassword": oldPassword,
			"NewPassword": newPassword,
		})

	if !result.IsValid() {
		return errors.NewValidationError(result.Error().Error())
	}

	// 查找用户
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.ErrUserNotFound
	}

	// 验证旧密码
	pwd, err := valueobject.NewPasswordFromHash(user.PasswordHash)
	if err != nil {
		return err
	}

	if !pwd.Verify(oldPassword) {
		return errors.ErrUserInvalidPassword
	}

	// 生成新密码哈希
	newPwd, err := valueobject.NewPassword(newPassword)
	if err != nil {
		return err
	}

	// 更新密码
	user.ChangePassword(newPwd.Hash())

	return s.userRepo.Update(ctx, user)
}

// ResetPassword 重置密码（忘记密码场景，验证码验证由外层处理）
// identifier 可以是邮箱或手机号
func (s *UserDomainService) ResetPassword(ctx context.Context, identifier string, newPasswordHash string) error {
	// 查找用户（支持邮箱或手机号）
	user, err := s.FindUserByIdentifier(ctx, identifier, identifier)
	if err != nil {
		return err
	}

	// 更新密码
	user.ChangePassword(newPasswordHash)

	return s.userRepo.Update(ctx, user)
}

// ResetPasswordWithValidation 重置密码（包含密码验证）
func (s *UserDomainService) ResetPasswordWithValidation(ctx context.Context, identifier string, newPassword string) error {
	// 验证新密码
	result := rules.ForField("NewPassword").Required().Length(6, 32).
		Validate(map[string]any{
			"NewPassword": newPassword,
		})

	if !result.IsValid() {
		return errors.NewValidationError(result.Error().Error())
	}

	// 生成新密码哈希
	newPwd, err := valueobject.NewPassword(newPassword)
	if err != nil {
		return err
	}

	return s.ResetPassword(ctx, identifier, newPwd.Hash())
}

// ========== 邮箱/手机号管理 ==========

// UpdateEmail 更换邮箱（验证码验证由外层处理）
func (s *UserDomainService) UpdateEmail(ctx context.Context, userID int, newEmailStr string) error {
	// 验证邮箱格式
	result := rules.ForField("Email").Required().Email().
		Validate(map[string]any{
			"Email": newEmailStr,
		})

	if !result.IsValid() {
		return errors.NewValidationError(result.Error().Error())
	}

	// 创建邮箱值对象
	newEmail, err := valueobject.NewEmail(newEmailStr)
	if err != nil {
		return errors.ErrUserInvalidEmail
	}

	// 检查新邮箱是否已被其他用户使用
	exists, err := s.userRepo.ExistsByEmail(ctx, newEmailStr)
	if err != nil {
		return err
	}
	if exists {
		return errors.ErrUserEmailAlreadyExists
	}

	// 查找用户
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.ErrUserNotFound
	}

	// 更新邮箱
	user.UpdateEmail(newEmail)

	return s.userRepo.Update(ctx, user)
}

// UpdatePhone 更换手机号（验证码验证由外层处理）
func (s *UserDomainService) UpdatePhone(ctx context.Context, userID int, newPhoneStr string) error {
	// 创建手机号值对象（自动验证格式）
	newPhone, err := valueobject.NewPhone(newPhoneStr)
	if err != nil {
		return err
	}

	// 检查新手机号是否已被其他用户使用
	exists, err := s.userRepo.ExistsByPhone(ctx, newPhoneStr)
	if err != nil {
		return err
	}
	if exists {
		return errors.ErrUserPhoneAlreadyExists
	}

	// 查找用户
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.ErrUserNotFound
	}

	// 更新手机号
	user.UpdatePhone(newPhone)

	return s.userRepo.Update(ctx, user)
}

// ========== 用户信息更新 ==========

// UpdateUser 更新用户信息（通用方法）
func (s *UserDomainService) UpdateUser(ctx context.Context, userID int, username string, status *enum.UserStatus) (*entity.User, error) {
	// 查找用户
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.ErrUserNotFound
	}

	// 更新用户名
	if username != "" && username != user.Username {
		// 检查新用户名是否已被使用
		exists, err := s.userRepo.ExistsByUsername(ctx, username)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.ErrUserUsernameExists
		}
		user.Username = username
		user.Touch()
	}

	// 更新状态
	if status != nil && *status != user.Status {
		user.Status = *status
		user.Touch()
	}

	return user, nil
}

// UpdateUserProfile 更新用户个人信息
func (s *UserDomainService) UpdateUserProfile(ctx context.Context, userID int, nickname string, avatar string) error {
	// 查找用户
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.ErrUserNotFound
	}

	// 更新昵称
	if nickname != "" {
		result := rules.ForField("Nickname").Length(1, 64).
			Validate(map[string]any{
				"Nickname": nickname,
			})
		if !result.IsValid() {
			return errors.NewValidationError(result.Error().Error())
		}
		user.UpdateNickname(nickname)
	}

	// 更新头像
	if avatar != "" {
		result := rules.ForField("Avatar").URL().MaxLength(512).
			Validate(map[string]any{
				"Avatar": avatar,
			})
		if !result.IsValid() {
			return errors.NewValidationError(result.Error().Error())
		}
		user.UpdateAvatar(avatar)
	}

	return s.userRepo.Update(ctx, user)
}

// ========== 用户状态管理 ==========

// LockUser 锁定用户
func (s *UserDomainService) LockUser(ctx context.Context, userID int) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.ErrUserNotFound
	}

	user.Lock()
	return s.userRepo.Update(ctx, user)
}

// UnlockUser 解锁用户
func (s *UserDomainService) UnlockUser(ctx context.Context, userID int) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.ErrUserNotFound
	}

	user.Unlock()
	return s.userRepo.Update(ctx, user)
}

// ActivateUser 激活用户
func (s *UserDomainService) ActivateUser(ctx context.Context, userID int) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.ErrUserNotFound
	}

	user.Activate()
	return s.userRepo.Update(ctx, user)
}

// DisableUser 禁用用户
func (s *UserDomainService) DisableUser(ctx context.Context, userID int) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.ErrUserNotFound
	}

	user.Disable()
	return s.userRepo.Update(ctx, user)
}

// ========== 辅助方法 ==========

// FindUserByIdentifier 通过邮箱或手机号查找用户
func (s *UserDomainService) FindUserByIdentifier(ctx context.Context, email string, phone string) (*entity.User, error) {
	var user *entity.User
	var err error

	// 优先使用邮箱查找
	if email != "" {
		if _, emailErr := valueobject.NewEmail(email); emailErr == nil {
			user, err = s.userRepo.FindByEmail(ctx, email)
			if err != nil {
				return nil, err
			}
			if user != nil {
				return user, nil
			}
		}
	}

	// 使用手机号查找
	if phone != "" {
		if _, phoneErr := valueobject.NewPhone(phone); phoneErr == nil {
			user, err = s.userRepo.FindByPhone(ctx, phone)
			if err != nil {
				return nil, err
			}
			if user != nil {
				return user, nil
			}
		}
	}

	return nil, errors.ErrUserNotFound
}

// CheckUserExists 检查邮箱或手机号是否已存在
func (s *UserDomainService) CheckUserExists(ctx context.Context, email string, phone string) error {
	// 检查邮箱是否存在
	if email != "" {
		exists, err := s.userRepo.ExistsByEmail(ctx, email)
		if err != nil {
			return fmt.Errorf("检查邮箱是否存在失败: %w", err)
		}
		if exists {
			return errors.ErrUserEmailAlreadyExists
		}
	}

	// 检查手机号是否存在
	if phone != "" {
		exists, err := s.userRepo.ExistsByPhone(ctx, phone)
		if err != nil {
			return fmt.Errorf("检查手机号是否存在失败: %w", err)
		}
		if exists {
			return errors.ErrUserPhoneAlreadyExists
		}
	}

	return nil
}

// GenerateUsername 生成用户名
func (s *UserDomainService) GenerateUsername(email string, phone string) string {
	timestamp := time.Now().Unix()

	if email != "" {
		// 从邮箱提取用户名
		emailObj, err := valueobject.NewEmail(email)
		if err == nil {
			return fmt.Sprintf("%s_%d", emailObj.LocalPart(), timestamp)
		}
	}

	if phone != "" {
		// 使用手机号后4位
		if len(phone) >= 4 {
			return fmt.Sprintf("user_%s", phone[len(phone)-4:])
		}
	}

	return fmt.Sprintf("user_%d", timestamp)
}

// ========== 用户注册（支持多种方式） ==========

// RegisterUser 用户注册（统一注册流程）
func (s *UserDomainService) RegisterUser(ctx context.Context, email string, phone string, passwordHash string, username string) (*entity.User, error) {
	// 1. 验证至少提供邮箱或手机号
	if email == "" && phone == "" {
		return nil, errors.NewValidationError("邮箱和手机号至少提供一个")
	}

	// 2. 检查邮箱或手机号是否已存在
	if err := s.CheckUserExists(ctx, email, phone); err != nil {
		return nil, err
	}

	// 3. 生成用户名（如果未提供）
	if username == "" {
		username = s.GenerateUsername(email, phone)
	} else {
		// 检查用户名是否已存在
		exists, err := s.userRepo.ExistsByUsername(ctx, username)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.ErrUserUsernameExists
		}
	}

	// 4. 创建邮箱值对象（如果提供）
	var emailObj valueobject.Email
	if email != "" {
		var err error
		emailObj, err = valueobject.NewEmail(email)
		if err != nil {
			return nil, errors.ErrUserInvalidEmail
		}
	}

	// 5. 创建手机号值对象（如果提供）
	var phoneObj valueobject.Phone
	if phone != "" {
		var err error
		phoneObj, err = valueobject.NewPhone(phone)
		if err != nil {
			return nil, err
		}
	}

	// 6. 创建用户实体
	now := time.Now()
	user := &entity.User{
		ID:           0,
		Username:     username,
		Email:        emailObj,
		Phone:        phoneObj,
		PasswordHash: passwordHash,
		Status:       enum.UserStatusActive, // 默认激活（验证码已验证）
		AuditFields: basegorm.AuditFields{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	// 7. 保存用户
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("保存用户失败: %w", err)
	}

	return user, nil
}

// RegisterByEmail 通过邮箱注册
func (s *UserDomainService) RegisterByEmail(ctx context.Context, emailStr string, password string, nickname string) (*entity.User, error) {
	// 验证参数
	result := rules.ForField("Email").Required().Email().
		ForField("Password").Required().Length(6, 32).
		Validate(map[string]any{
			"Email":    emailStr,
			"Password": password,
		})

	if !result.IsValid() {
		return nil, errors.NewValidationError(result.Error().Error())
	}

	// 创建邮箱值对象
	email, err := valueobject.NewEmail(emailStr)
	if err != nil {
		return nil, errors.ErrUserInvalidEmail
	}

	// 检查邮箱是否已存在
	exists, err := s.userRepo.ExistsByEmail(ctx, emailStr)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.ErrUserEmailAlreadyExists
	}

	// 生成密码哈希
	pwd, err := valueobject.NewPassword(password)
	if err != nil {
		return nil, err
	}

	// 生成默认用户名（可选）
	username := fmt.Sprintf("user_%d", time.Now().UnixNano()%1000000)

	// 设置默认昵称
	if nickname == "" {
		nickname = "新用户"
	}

	// 创建用户实体
	now := time.Now()
	user := &entity.User{
		Username:     username,
		Nickname:     nickname,
		Email:        email,
		PasswordHash: pwd.Hash(),
		Status:       enum.UserStatusInactive, // 默认未激活，需要验证邮箱后激活
		AuditFields: basegorm.AuditFields{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	return user, nil
}

// RegisterByPhone 通过手机号注册
func (s *UserDomainService) RegisterByPhone(ctx context.Context, phoneStr string, password string, nickname string) (*entity.User, error) {
	// 验证密码
	result := rules.ForField("Password").Required().Length(6, 32).
		Validate(map[string]any{
			"Password": password,
		})

	if !result.IsValid() {
		return nil, errors.NewValidationError(result.Error().Error())
	}

	// 创建手机号值对象（自动验证格式）
	phone, err := valueobject.NewPhone(phoneStr)
	if err != nil {
		return nil, err
	}

	// 检查手机号是否已存在
	exists, err := s.userRepo.ExistsByPhone(ctx, phoneStr)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.ErrUserPhoneAlreadyExists
	}

	// 生成密码哈希
	pwd, err := valueobject.NewPassword(password)
	if err != nil {
		return nil, err
	}

	// 生成默认用户名
	username := fmt.Sprintf("user_%d", time.Now().UnixNano()%1000000)

	// 设置默认昵称
	if nickname == "" {
		nickname = "新用户"
	}

	// 创建用户实体
	now := time.Now()
	user := &entity.User{
		Username:     username,
		Nickname:     nickname,
		Phone:        phone,
		PasswordHash: pwd.Hash(),
		Status:       enum.UserStatusActive, // 手机号注册直接激活
		AuditFields: basegorm.AuditFields{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	return user, nil
}
