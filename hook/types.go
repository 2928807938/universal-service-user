package hook

import (
	"context"
	"universal-service-user/user/domain/entity"
)

// HookType 钩子类型
type HookType string

const (
	// 用户创建钩子
	BeforeCreate HookType = "before_create"
	AfterCreate  HookType = "after_create"

	// 用户更新钩子
	BeforeUpdate HookType = "before_update"
	AfterUpdate  HookType = "after_update"

	// 用户删除钩子
	BeforeDelete HookType = "before_delete"
	AfterDelete  HookType = "after_delete"

	// 用户登录钩子
	BeforeLogin HookType = "before_login"
	AfterLogin  HookType = "after_login"

	// 密码重置钩子
	BeforeResetPassword HookType = "before_reset_password"
	AfterResetPassword  HookType = "after_reset_password"

	// 邮箱更换钩子
	BeforeChangeEmail HookType = "before_change_email"
	AfterChangeEmail  HookType = "after_change_email"

	// 手机号更换钩子
	BeforeChangePhone HookType = "before_change_phone"
	AfterChangePhone  HookType = "after_change_phone"

	// 登出钩子
	BeforeLogout HookType = "before_logout"
	AfterLogout  HookType = "after_logout"
)

// Context 钩子上下文
// 钩子函数接收一个上下文对象，包含执行钩子所需的所有信息
type Context struct {
	// Ctx Go 标准上下文
	Ctx context.Context

	// User 用户实体（如果存在）
	// 在 BeforeCreate 时可能为 nil
	// 在 AfterCreate 及其他钩子中应该存在
	User *entity.User

	// Request 原始请求数据
	// 包含触发此钩子的原始请求信息
	// 类型为 interface{} 以支持不同的请求类型
	Request interface{}

	// Metadata 自定义元数据
	// 可用于在钩子之间传递数据
	// 或存储钩子执行过程中的临时数据
	Metadata map[string]interface{}

	// HookType 当前钩子类型
	// 用于标识当前执行的是哪个钩子
	HookType HookType
}

// HookFunc 钩子函数类型
// 接收一个钩子上下文，返回错误（如果有）
// 如果钩子函数返回错误，将中断当前操作的执行
type HookFunc func(ctx *Context) error

// NewContext 创建钩子上下文
func NewContext(ctx context.Context, hookType HookType) *Context {
	return &Context{
		Ctx:      ctx,
		HookType: hookType,
		Metadata: make(map[string]interface{}),
	}
}

// WithUser 设置用户实体
func (c *Context) WithUser(user *entity.User) *Context {
	c.User = user
	return c
}

// WithRequest 设置原始请求数据
func (c *Context) WithRequest(request interface{}) *Context {
	c.Request = request
	return c
}

// SetMetadata 设置元数据
func (c *Context) SetMetadata(key string, value interface{}) {
	c.Metadata[key] = value
}

// GetMetadata 获取元数据
func (c *Context) GetMetadata(key string) (interface{}, bool) {
	value, ok := c.Metadata[key]
	return value, ok
}

// GetMetadataString 获取字符串类型的元数据
func (c *Context) GetMetadataString(key string) (string, bool) {
	value, ok := c.Metadata[key]
	if !ok {
		return "", false
	}
	str, ok := value.(string)
	return str, ok
}

// GetMetadataInt 获取整数类型的元数据
func (c *Context) GetMetadataInt(key string) (int, bool) {
	value, ok := c.Metadata[key]
	if !ok {
		return 0, false
	}
	i, ok := value.(int)
	return i, ok
}

// GetMetadataBool 获取布尔类型的元数据
func (c *Context) GetMetadataBool(key string) (bool, bool) {
	value, ok := c.Metadata[key]
	if !ok {
		return false, false
	}
	b, ok := value.(bool)
	return b, ok
}
