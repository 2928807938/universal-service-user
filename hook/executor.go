package hook

import (
	"fmt"
	"log"
	"time"
)

// Executor 钩子执行器
// 负责执行注册在注册中心的钩子函数
type Executor struct {
	// registry 钩子注册中心
	registry *Registry

	// enableLog 是否启用日志
	enableLog bool

	// logger 日志记录器
	logger Logger
}

// Logger 日志接口
// 允许使用者自定义日志实现
type Logger interface {
	Info(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
}

// defaultLogger 默认日志实现
type defaultLogger struct{}

func (l *defaultLogger) Info(msg string, fields ...interface{}) {
	log.Printf("[INFO] %s %v", msg, fields)
}

func (l *defaultLogger) Error(msg string, fields ...interface{}) {
	log.Printf("[ERROR] %s %v", msg, fields)
}

func (l *defaultLogger) Debug(msg string, fields ...interface{}) {
	log.Printf("[DEBUG] %s %v", msg, fields)
}

// NewExecutor 创建新的钩子执行器
// registry: 钩子注册中心，如果为 nil 则使用全局注册中心
func NewExecutor(registry *Registry) *Executor {
	if registry == nil {
		registry = defaultRegistry
	}
	return &Executor{
		registry:  registry,
		enableLog: false,
		logger:    &defaultLogger{},
	}
}

// WithLogger 设置日志记录器
func (e *Executor) WithLogger(logger Logger) *Executor {
	e.logger = logger
	e.enableLog = true
	return e
}

// EnableLog 启用日志
func (e *Executor) EnableLog(enable bool) *Executor {
	e.enableLog = enable
	return e
}

// Execute 执行指定类型的所有钩子
// ctx: 钩子上下文
// 返回错误（如果有任何钩子执行失败）
//
// 执行顺序：按注册顺序依次执行
// 错误处理：如果某个钩子返回错误，立即停止执行并返回该错误
func (e *Executor) Execute(ctx *Context) error {
	hookType := ctx.HookType

	// 获取该类型的所有钩子
	hooks := e.registry.Get(hookType)

	// 如果没有注册钩子，直接返回
	if len(hooks) == 0 {
		if e.enableLog {
			e.logger.Debug("no hooks registered", "hookType", hookType)
		}
		return nil
	}

	if e.enableLog {
		e.logger.Info(
			"executing hooks",
			"hookType", hookType,
			"count", len(hooks),
		)
	}

	// 依次执行所有钩子
	for i, hookFunc := range hooks {
		startTime := time.Now()

		if e.enableLog {
			e.logger.Debug(
				"executing hook",
				"hookType", hookType,
				"index", i+1,
				"total", len(hooks),
			)
		}

		// 执行钩子
		err := hookFunc(ctx)

		// 记录执行时间
		duration := time.Since(startTime)

		if err != nil {
			if e.enableLog {
				e.logger.Error(
					"hook execution failed",
					"hookType", hookType,
					"index", i+1,
					"duration", duration,
					"error", err,
				)
			}
			return fmt.Errorf("hook execution failed at index %d: %w", i+1, err)
		}

		if e.enableLog {
			e.logger.Debug(
				"hook executed successfully",
				"hookType", hookType,
				"index", i+1,
				"duration", duration,
			)
		}
	}

	if e.enableLog {
		e.logger.Info(
			"all hooks executed successfully",
			"hookType", hookType,
			"count", len(hooks),
		)
	}

	return nil
}

// ExecuteAsync 异步执行钩子
// ctx: 钩子上下文
// callback: 执行完成后的回调函数（可选）
//
// 注意：异步执行时，钩子中的错误会通过回调函数返回，不会阻塞主流程
func (e *Executor) ExecuteAsync(ctx *Context, callback func(error)) {
	go func() {
		err := e.Execute(ctx)
		if callback != nil {
			callback(err)
		}
	}()
}

// MustExecute 执行钩子，如果失败则 panic
// ctx: 钩子上下文
//
// 仅在确定钩子必须成功执行的场景下使用
func (e *Executor) MustExecute(ctx *Context) {
	if err := e.Execute(ctx); err != nil {
		panic(fmt.Sprintf("hook execution failed: %v", err))
	}
}

// 全局钩子执行器
var defaultExecutor = NewExecutor(nil)

// Execute 使用全局执行器执行钩子
func Execute(ctx *Context) error {
	return defaultExecutor.Execute(ctx)
}

// ExecuteAsync 使用全局执行器异步执行钩子
func ExecuteAsync(ctx *Context, callback func(error)) {
	defaultExecutor.ExecuteAsync(ctx, callback)
}

// MustExecute 使用全局执行器执行钩子，失败则 panic
func MustExecute(ctx *Context) {
	defaultExecutor.MustExecute(ctx)
}

// SetLogger 为全局执行器设置日志记录器
func SetLogger(logger Logger) {
	defaultExecutor.WithLogger(logger)
}

// EnableGlobalLog 启用全局执行器的日志
func EnableGlobalLog(enable bool) {
	defaultExecutor.EnableLog(enable)
}

// GetDefaultExecutor 获取全局钩子执行器
// 用于在需要独立执行器的场景下获取全局实例
func GetDefaultExecutor() *Executor {
	return defaultExecutor
}
