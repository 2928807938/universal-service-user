package request

// LoginRequest 用户登录请求
type LoginRequest struct {
	LoginType string `json:"login_type" vd:"$=='username'||$=='email'||$=='phone'"` // 登录方式: username / email / phone
	Username  string `json:"username"`                                               // 用户名
	Email     string `json:"email"`                                                  // 邮箱
	Phone     string `json:"phone"`                                                  // 手机号
	Password  string `json:"password"`                                               // 密码
	Code      string `json:"code"`                                                   // 验证码(手机号登录时需要)
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" vd:"len($)>0"` // 刷新令牌
}

// LogoutRequest 登出请求
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"` // 刷新令牌(可选,如果提供则只登出该设备)
}
