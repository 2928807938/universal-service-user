package enum

// VerificationScene 验证码场景枚举
type VerificationScene string

const (
	// SceneRegister 注册场景
	SceneRegister VerificationScene = "REGISTER"

	// SceneLogin 登录场景
	SceneLogin VerificationScene = "LOGIN"

	// SceneResetPassword 重置密码场景
	SceneResetPassword VerificationScene = "RESET_PASSWORD"

	// SceneChangeEmailOld 更换邮箱：验证旧邮箱
	SceneChangeEmailOld VerificationScene = "CHANGE_EMAIL_OLD"

	// SceneChangeEmailNew 更换邮箱：验证新邮箱
	SceneChangeEmailNew VerificationScene = "CHANGE_EMAIL_NEW"

	// SceneChangePhoneOld 更换手机号：验证旧手机号
	SceneChangePhoneOld VerificationScene = "CHANGE_PHONE_OLD"

	// SceneChangePhoneNew 更换手机号：验证新手机号
	SceneChangePhoneNew VerificationScene = "CHANGE_PHONE_NEW"
)

// String 返回场景字符串
func (s VerificationScene) String() string {
	return string(s)
}

// IsValid 验证场景是否有效
func (s VerificationScene) IsValid() bool {
	switch s {
	case SceneRegister,
		SceneLogin,
		SceneResetPassword,
		SceneChangeEmailOld,
		SceneChangeEmailNew,
		SceneChangePhoneOld,
		SceneChangePhoneNew:
		return true
	default:
		return false
	}
}

// Description 返回场景描述
func (s VerificationScene) Description() string {
	switch s {
	case SceneRegister:
		return "注册"
	case SceneLogin:
		return "登录"
	case SceneResetPassword:
		return "重置密码"
	case SceneChangeEmailOld:
		return "更换邮箱：验证旧邮箱"
	case SceneChangeEmailNew:
		return "更换邮箱：验证新邮箱"
	case SceneChangePhoneOld:
		return "更换手机号：验证旧手机号"
	case SceneChangePhoneNew:
		return "更换手机号：验证新手机号"
	default:
		return "未知场景"
	}
}
