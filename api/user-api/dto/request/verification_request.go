package request

// SendCodeRequest 发送验证码请求
type SendCodeRequest struct {
	Type   string `json:"type" vd:"$=='email'||$=='phone'"` // 类型: email / phone
	Target string `json:"target" vd:"len($)>0"`             // 目标: 邮箱地址或手机号
	Scene  string `json:"scene" vd:"len($)>0"`              // 场景: register / login / reset_password 等
}

// VerifyCodeRequest 验证验证码请求
type VerifyCodeRequest struct {
	Type   string `json:"type" vd:"$=='email'||$=='phone'"` // 类型: email / phone
	Target string `json:"target" vd:"len($)>0"`             // 目标: 邮箱地址或手机号
	Scene  string `json:"scene" vd:"len($)>0"`              // 场景
	Code   string `json:"code" vd:"len($)==6"`              // 验证码
}

// ResetPasswordRequest 重置密码请求(忘记密码场景)
type ResetPasswordRequest struct {
	Email       string `json:"email"`                                     // 邮箱(与手机号二选一)
	Phone       string `json:"phone"`                                     // 手机号(与邮箱二选一)
	Code        string `json:"code" vd:"len($)==6"`                       // 验证码
	NewPassword string `json:"new_password" vd:"len($)>=6 && len($)<=32"` // 新密码
}
