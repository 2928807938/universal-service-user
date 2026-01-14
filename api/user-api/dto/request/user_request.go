package request

// RegisterRequest 用户注册请求
type RegisterRequest struct {
	Email    string `json:"email" vd:"email($)"`                    // 邮箱(与手机号二选一)
	Phone    string `json:"phone"`                                   // 手机号(与邮箱二选一)
	Password string `json:"password" vd:"len($)>=6 && len($)<=32"`  // 密码(6-32位)
	Code     string `json:"code" vd:"len($)==6"`                     // 验证码
	Username string `json:"username"`                                // 用户名(可选,默认生成)
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" vd:"len($)>2 && len($)<51"`
	Email    string `json:"email" vd:"email($)"`
	Password string `json:"password" vd:"len($)>5 && len($)<51"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Username *string `json:"username,omitempty" vd:"len($)>2 && len($)<51"`
	Nickname *string `json:"nickname,omitempty" vd:"len($)>0 && len($)<=64"`
	Avatar   *string `json:"avatar,omitempty" vd:"len($)<=512"`
	Status   *int    `json:"status,omitempty" vd:"$>=0 && $<=3"`
}

// ChangePasswordRequest 修改密码请求(已登录用户)
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" vd:"len($)>=6 && len($)<=32"` // 旧密码
	NewPassword string `json:"new_password" vd:"len($)>=6 && len($)<=32"` // 新密码
}

// ListUsersRequest 用户列表请求
type ListUsersRequest struct {
	Page     int `query:"page"`
	PageSize int `query:"page_size"`
}

// SetDefaults 设置默认值
func (r *ListUsersRequest) SetDefaults() {
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.PageSize <= 0 {
		r.PageSize = 10
	}
}
