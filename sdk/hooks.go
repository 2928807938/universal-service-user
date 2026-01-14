package sdk

import (
	"github.com/2928807938/universal-service-user/hook"
)

// =============================================================================
// 钩子注册方法
// =============================================================================

// OnBeforeCreate 注册用户创建前钩子
// 钩子函数接收 hook.Context，可以访问：
//   - ctx.Request: RegisterRequest
//   - ctx.Metadata: 自定义元数据
//
// 钩子返回错误将阻止用户创建
// 示例:
//
//	client.OnBeforeCreate(func(ctx *hook.Context) error {
//	    req := ctx.Request.(*sdk.RegisterRequest)
//	    log.Printf("准备创建用户: %s", req.Email)
//	    // 可以在这里进行额外的验证
//	    return nil
//	})
func (c *Client) OnBeforeCreate(hookFunc hook.HookFunc) {
	c.hookRegistry.Register(hook.BeforeCreate, hookFunc)
}

// OnAfterCreate 注册用户创建后钩子
// 钩子函数接收 hook.Context，可以访问：
//   - ctx.Request: RegisterRequest
//   - ctx.Metadata["user_id"]: 新创建的用户 ID
//   - ctx.Metadata["user"]: 新创建的用户信息 (*UserInfo)
//
// 钩子返回错误不影响用户创建结果，但会记录日志
// 示例:
//
//	client.OnAfterCreate(func(ctx *hook.Context) error {
//	    userID := ctx.Metadata["user_id"].(int)
//	    user := ctx.Metadata["user"].(*sdk.UserInfo)
//	    log.Printf("用户创建成功: ID=%d, Email=%s", userID, user.Email)
//	    // 发送欢迎邮件
//	    sendWelcomeEmail(user.Email)
//	    // 初始化用户默认数据
//	    initUserDefaults(userID)
//	    return nil
//	})
func (c *Client) OnAfterCreate(hookFunc hook.HookFunc) {
	c.hookRegistry.Register(hook.AfterCreate, hookFunc)
}

// OnBeforeUpdate 注册用户更新前钩子
// 钩子函数接收 hook.Context，可以访问：
//   - ctx.Request: UpdateUserRequest
//   - ctx.Metadata["user_id"]: 要更新的用户 ID
//
// 钩子返回错误将阻止用户更新
// 示例:
//
//	client.OnBeforeUpdate(func(ctx *hook.Context) error {
//	    userID := ctx.Metadata["user_id"].(int)
//	    req := ctx.Request.(*sdk.UpdateUserRequest)
//	    log.Printf("准备更新用户: ID=%d", userID)
//	    // 可以在这里进行权限检查
//	    return nil
//	})
func (c *Client) OnBeforeUpdate(hookFunc hook.HookFunc) {
	c.hookRegistry.Register(hook.BeforeUpdate, hookFunc)
}

// OnAfterUpdate 注册用户更新后钩子
// 钩子函数接收 hook.Context，可以访问：
//   - ctx.Request: UpdateUserRequest
//   - ctx.Metadata["user_id"]: 更新的用户 ID
//   - ctx.Metadata["user"]: 更新后的用户信息 (*UserInfo)
//
// 钩子返回错误不影响用户更新结果，但会记录日志
// 示例:
//
//	client.OnAfterUpdate(func(ctx *hook.Context) error {
//	    userID := ctx.Metadata["user_id"].(int)
//	    user := ctx.Metadata["user"].(*sdk.UserInfo)
//	    log.Printf("用户更新成功: ID=%d", userID)
//	    // 同步到外部系统
//	    syncToExternalSystem(user)
//	    return nil
//	})
func (c *Client) OnAfterUpdate(hookFunc hook.HookFunc) {
	c.hookRegistry.Register(hook.AfterUpdate, hookFunc)
}

// OnBeforeDelete 注册用户删除前钩子
// 钩子函数接收 hook.Context，可以访问：
//   - ctx.Metadata["user_id"]: 要删除的用户 ID
//
// 钩子返回错误将阻止用户删除
// 示例:
//
//	client.OnBeforeDelete(func(ctx *hook.Context) error {
//	    userID := ctx.Metadata["user_id"].(int)
//	    log.Printf("准备删除用户: ID=%d", userID)
//	    // 检查用户是否有关联数据
//	    if hasRelatedData(userID) {
//	        return errors.New("用户有关联数据，无法删除")
//	    }
//	    return nil
//	})
func (c *Client) OnBeforeDelete(hookFunc hook.HookFunc) {
	c.hookRegistry.Register(hook.BeforeDelete, hookFunc)
}

// OnAfterDelete 注册用户删除后钩子
// 钩子函数接收 hook.Context，可以访问：
//   - ctx.Metadata["user_id"]: 删除的用户 ID
//
// 钩子返回错误不影响用户删除结果，但会记录日志
// 示例:
//
//	client.OnAfterDelete(func(ctx *hook.Context) error {
//	    userID := ctx.Metadata["user_id"].(int)
//	    log.Printf("用户删除成功: ID=%d", userID)
//	    // 清理关联数据
//	    cleanupUserData(userID)
//	    return nil
//	})
func (c *Client) OnAfterDelete(hookFunc hook.HookFunc) {
	c.hookRegistry.Register(hook.AfterDelete, hookFunc)
}

// OnBeforeLogin 注册用户登录前钩子
// 钩子函数接收 hook.Context，可以访问：
//   - ctx.Request: LoginRequest
//
// 钩子返回错误将阻止用户登录
// 示例:
//
//	client.OnBeforeLogin(func(ctx *hook.Context) error {
//	    req := ctx.Request.(*sdk.LoginRequest)
//	    log.Printf("准备登录: %s", req.Email)
//	    // 风控检查
//	    if isRisky(req) {
//	        return errors.New("登录风险较高，请稍后重试")
//	    }
//	    return nil
//	})
func (c *Client) OnBeforeLogin(hookFunc hook.HookFunc) {
	c.hookRegistry.Register(hook.BeforeLogin, hookFunc)
}

// OnAfterLogin 注册用户登录后钩子
// 钩子函数接收 hook.Context，可以访问：
//   - ctx.Request: LoginRequest
//   - ctx.Metadata["user_id"]: 登录的用户 ID
//   - ctx.Metadata["user"]: 登录的用户信息 (*UserInfo)
//
// 钩子返回错误不影响登录结果，但会记录日志
// 示例:
//
//	client.OnAfterLogin(func(ctx *hook.Context) error {
//	    userID := ctx.Metadata["user_id"].(int)
//	    user := ctx.Metadata["user"].(*sdk.UserInfo)
//	    req := ctx.Request.(*sdk.LoginRequest)
//	    log.Printf("用户登录成功: ID=%d, Email=%s", userID, user.Email)
//	    // 记录登录日志
//	    logLoginActivity(userID, req.LoginType)
//	    // 更新最后登录时间
//	    updateLastLoginTime(userID)
//	    return nil
//	})
func (c *Client) OnAfterLogin(hookFunc hook.HookFunc) {
	c.hookRegistry.Register(hook.AfterLogin, hookFunc)
}

// OnBeforeLogout 注册用户登出前钩子
// 钩子函数接收 hook.Context，可以访问：
//   - ctx.Metadata["access_token"]: 访问令牌
//   - ctx.Metadata["refresh_token"]: 刷新令牌
//
// 钩子返回错误将阻止用户登出
// 示例:
//
//	client.OnBeforeLogout(func(ctx *hook.Context) error {
//	    log.Println("准备登出")
//	    return nil
//	})
func (c *Client) OnBeforeLogout(hookFunc hook.HookFunc) {
	c.hookRegistry.Register(hook.BeforeLogout, hookFunc)
}

// OnAfterLogout 注册用户登出后钩子
// 钩子函数接收 hook.Context
// 钩子返回错误不影响登出结果，但会记录日志
// 示例:
//
//	client.OnAfterLogout(func(ctx *hook.Context) error {
//	    log.Println("用户登出成功")
//	    // 清理相关缓存
//	    clearUserCache()
//	    return nil
//	})
func (c *Client) OnAfterLogout(hookFunc hook.HookFunc) {
	c.hookRegistry.Register(hook.AfterLogout, hookFunc)
}

// OnBeforeResetPassword 注册重置密码前钩子
// 钩子函数接收 hook.Context，可以访问：
//   - ctx.Request: ResetPasswordRequest
//
// 钩子返回错误将阻止密码重置
// 示例:
//
//	client.OnBeforeResetPassword(func(ctx *hook.Context) error {
//	    req := ctx.Request.(*sdk.ResetPasswordRequest)
//	    log.Printf("准备重置密码: %s", req.Email)
//	    // 安全检查
//	    return nil
//	})
func (c *Client) OnBeforeResetPassword(hookFunc hook.HookFunc) {
	c.hookRegistry.Register(hook.BeforeResetPassword, hookFunc)
}

// OnAfterResetPassword 注册重置密码后钩子
// 钩子函数接收 hook.Context，可以访问：
//   - ctx.Request: ResetPasswordRequest
//
// 钩子返回错误不影响密码重置结果，但会记录日志
// 示例:
//
//	client.OnAfterResetPassword(func(ctx *hook.Context) error {
//	    req := ctx.Request.(*sdk.ResetPasswordRequest)
//	    log.Printf("密码重置成功: %s", req.Email)
//	    // 通知用户密码已修改
//	    notifyPasswordChanged(req.Email)
//	    // 踢出其他会话
//	    kickoutOtherSessions(req.Email)
//	    return nil
//	})
func (c *Client) OnAfterResetPassword(hookFunc hook.HookFunc) {
	c.hookRegistry.Register(hook.AfterResetPassword, hookFunc)
}

// OnBeforeChangeEmail 注册更换邮箱前钩子
// 钩子函数接收 hook.Context，可以访问：
//   - ctx.Metadata["user_id"]: 用户 ID
//   - ctx.Metadata["old_email"]: 旧邮箱
//   - ctx.Metadata["new_email"]: 新邮箱
//
// 钩子返回错误将阻止邮箱更换
// 示例:
//
//	client.OnBeforeChangeEmail(func(ctx *hook.Context) error {
//	    userID := ctx.Metadata["user_id"].(int)
//	    oldEmail := ctx.Metadata["old_email"].(string)
//	    newEmail := ctx.Metadata["new_email"].(string)
//	    log.Printf("准备更换邮箱: 用户ID=%d, %s -> %s", userID, oldEmail, newEmail)
//	    return nil
//	})
func (c *Client) OnBeforeChangeEmail(hookFunc hook.HookFunc) {
	c.hookRegistry.Register(hook.BeforeChangeEmail, hookFunc)
}

// OnAfterChangeEmail 注册更换邮箱后钩子
// 钩子函数接收 hook.Context，可以访问：
//   - ctx.Metadata["user_id"]: 用户 ID
//   - ctx.Metadata["old_email"]: 旧邮箱
//   - ctx.Metadata["new_email"]: 新邮箱
//
// 钩子返回错误不影响邮箱更换结果，但会记录日志
// 示例:
//
//	client.OnAfterChangeEmail(func(ctx *hook.Context) error {
//	    userID := ctx.Metadata["user_id"].(int)
//	    newEmail := ctx.Metadata["new_email"].(string)
//	    log.Printf("邮箱更换成功: 用户ID=%d, 新邮箱=%s", userID, newEmail)
//	    // 发送通知邮件到新旧邮箱
//	    notifyEmailChanged(userID, newEmail)
//	    return nil
//	})
func (c *Client) OnAfterChangeEmail(hookFunc hook.HookFunc) {
	c.hookRegistry.Register(hook.AfterChangeEmail, hookFunc)
}

// OnBeforeChangePhone 注册更换手机号前钩子
// 钩子函数接收 hook.Context，可以访问：
//   - ctx.Metadata["user_id"]: 用户 ID
//   - ctx.Metadata["old_phone"]: 旧手机号
//   - ctx.Metadata["new_phone"]: 新手机号
//
// 钩子返回错误将阻止手机号更换
// 示例:
//
//	client.OnBeforeChangePhone(func(ctx *hook.Context) error {
//	    userID := ctx.Metadata["user_id"].(int)
//	    oldPhone := ctx.Metadata["old_phone"].(string)
//	    newPhone := ctx.Metadata["new_phone"].(string)
//	    log.Printf("准备更换手机号: 用户ID=%d, %s -> %s", userID, oldPhone, newPhone)
//	    return nil
//	})
func (c *Client) OnBeforeChangePhone(hookFunc hook.HookFunc) {
	c.hookRegistry.Register(hook.BeforeChangePhone, hookFunc)
}

// OnAfterChangePhone 注册更换手机号后钩子
// 钩子函数接收 hook.Context，可以访问：
//   - ctx.Metadata["user_id"]: 用户 ID
//   - ctx.Metadata["old_phone"]: 旧手机号
//   - ctx.Metadata["new_phone"]: 新手机号
//
// 钩子返回错误不影响手机号更换结果，但会记录日志
// 示例:
//
//	client.OnAfterChangePhone(func(ctx *hook.Context) error {
//	    userID := ctx.Metadata["user_id"].(int)
//	    newPhone := ctx.Metadata["new_phone"].(string)
//	    log.Printf("手机号更换成功: 用户ID=%d, 新手机号=%s", userID, newPhone)
//	    // 发送通知短信到新旧手机号
//	    notifyPhoneChanged(userID, newPhone)
//	    return nil
//	})
func (c *Client) OnAfterChangePhone(hookFunc hook.HookFunc) {
	c.hookRegistry.Register(hook.AfterChangePhone, hookFunc)
}

// =============================================================================
// 钩子管理方法
// =============================================================================

// ClearHooks 清空指定类型的所有钩子
// hookType: 钩子类型
// 示例:
//
//	client.ClearHooks(hook.BeforeCreate)
func (c *Client) ClearHooks(hookType hook.HookType) {
	c.hookRegistry.Clear(hookType)
}

// ClearAllHooks 清空所有钩子
// 示例:
//
//	client.ClearAllHooks()
func (c *Client) ClearAllHooks() {
	c.hookRegistry.ClearAll()
}

// HasHooks 检查是否注册了指定类型的钩子
// hookType: 钩子类型
// 返回: true 表示已注册该类型的钩子，false 表示未注册
// 示例:
//
//	if client.HasHooks(hook.BeforeCreate) {
//	    log.Println("已注册 BeforeCreate 钩子")
//	}
func (c *Client) HasHooks(hookType hook.HookType) bool {
	return c.hookRegistry.Has(hookType)
}

// CountHooks 获取指定类型已注册的钩子数量
// hookType: 钩子类型
// 返回: 该类型已注册的钩子数量
// 示例:
//
//	count := client.CountHooks(hook.BeforeCreate)
//	log.Printf("BeforeCreate 钩子数量: %d", count)
func (c *Client) CountHooks(hookType hook.HookType) int {
	return c.hookRegistry.Count(hookType)
}

// GetAllHookTypes 获取所有已注册钩子的类型
// 返回: 所有已注册钩子的类型列表
// 示例:
//
//	types := client.GetAllHookTypes()
//	for _, t := range types {
//	    log.Printf("已注册钩子类型: %s", t)
//	}
func (c *Client) GetAllHookTypes() []hook.HookType {
	return c.hookRegistry.GetAllTypes()
}

// GetHookRegistry 获取钩子注册中心
// 返回钩子注册中心实例，供高级用户直接操作
// 示例:
//
//	registry := client.GetHookRegistry()
//	hooks := registry.Get(hook.BeforeCreate)
func (c *Client) GetHookRegistry() *hook.Registry {
	return c.hookRegistry
}

// GetHookExecutor 获取钩子执行器
// 返回钩子执行器实例，供高级用户直接操作
// 示例:
//
//	executor := client.GetHookExecutor()
//	executor.EnableLog(true)
func (c *Client) GetHookExecutor() *hook.Executor {
	return c.hookExecutor
}
