package vo

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string  `json:"access_token"`  // 访问令牌
	RefreshToken string  `json:"refresh_token"` // 刷新令牌
	ExpiresIn    int     `json:"expires_in"`    // 过期时间(秒)
	UserInfo     *UserVo `json:"user_info"`     // 用户信息
}

// RefreshTokenResponse 刷新令牌响应
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`  // 新的访问令牌
	RefreshToken string `json:"refresh_token"` // 新的刷新令牌(可选)
	ExpiresIn    int    `json:"expires_in"`    // 过期时间(秒)
}

// SendCodeResponse 发送验证码响应
type SendCodeResponse struct {
	Message string `json:"message"` // 提示信息
}

// VerifyCodeResponse 验证验证码响应
type VerifyCodeResponse struct {
	Valid   bool   `json:"valid"`   // 是否有效
	Message string `json:"message"` // 提示信息
}
