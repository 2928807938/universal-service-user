# Go 规则引擎架构设计文档

## 1. 概述

本文档描述了一个可扩展的 Go 语言规则引擎的架构设计。该规则引擎支持内置验证规则和用户自定义规则，采用分层架构设计，确保核心抽象的稳定性和外部扩展的灵活性。

## 2. 架构分层

整体架构分为五个层次，从底层到上层依次为：

```
┌─────────────────────────────────────────┐
│            API Layer (应用层)            │
├─────────────────────────────────────────┤
│     Implementation Layer (实现层)        │
├──────────────┬──────────────────────────┤
│ Predefined   │     Definition Layer      │
│   Layer      │     (自定义规则层)        │
│ (预定义规则)  │                           │
├──────────────┴──────────────────────────┤
│          Core Layer (核心层)            │
└─���───────────────────────────────────────┘
```

### 2.1 Core Layer（核心层）

**职责**：定义规则引擎的核心概念和契约，提供所有层共用的基础抽象。

**主要组件**：

- **Rule 接口**
  - 所有规则必须实现的核心接口
  - 定义验证方法：`Validate(ctx RuleContext) RuleResult`
  - 描述信息方法：`Name() string`, `Description() string`

- **RuleContext（验证上下文）**
  - 封装验证所需的输入数据
  - 包含字段名、字段值、元数据等
  - 提���上下文传递能力（支持跨规则数据共享）

- **RuleResult（验证结果）**
  - 统一的结果模型
  - 包含成功/失败状态
  - 错误信息和错误详情
  - 支持结果链（多个规则的组合结果）

- **RuleEngine 接口**
  - 规则引擎的抽象定义
  - 规则注册与管理
  - 规则执行编排

- **错误模型**
  - ValidationError：验证失败错误
  - RuleConfigError：规则配置错误
  - RuleExecutionError：规则执行错误

**设计原则**：
- 接口最小化，只定义必需的方法
- 不依赖任何具体实现
- 提供扩展点，支持用户自定义行为

### 2.2 Predefined Layer（预定义规则层）

**职责**：提供开箱即用的常用验证规则，覆盖大部分常见场景。这是底层实现，不对用户直接暴露，通过 Fluent API 间接使用。

**规则分类**：

1. **基础规则**
   - `RequiredRule`：必填验证
   - `NotNullRule`：非空验证
   - `NotEmptyRule`：非空字符串验证

2. **字符串规则**
   - `EmailRule`：邮箱格式验证
   - `LengthRule`：字符串长度范围
   - `PatternRule`：正则表达式匹配
   - `URLRule`：URL 格式验证
   - `ContainsRule`：包含子串验证

3. **数值规则**
   - `MinRule`：最小值验证
   - `MaxRule`：最大值验证
   - `RangeRule`：数值范围验证
   - `PositiveRule`：正数验证

4. **时间规则**
   - `BeforeRule`：早于指定时间
   - `AfterRule`：晚于指定时间
   - `BetweenRule`：时间范围验证

5. **集合规则**
   - `ArrayLengthRule`：数组/切片长度
   - `InRule`：值在指定集合中

**特点**：
- 不可变对象，线程安全
- 通过 Fluent API 的方法创建和配置
- 错误消息支持国际化
- 完全实现 Core 层的 Rule 接口

**注意**：用户不直接使用这些规则类，而是通过 Fluent API 间接调用。

### 2.3 Definition Layer（自定义规则层）

**职责**：支持用户动态定义和配置规则，实现运行时扩展。

**主要组件**：

1. **CustomRule（自定义规则）**
   - 通用的自定义规则实现
   - 接受验证函数作为参数
   - 支持闭包捕获上下文

2. **Rule Definition（规则定义）**
   - 支持多种定义格式：
     - 结构体定义
     - JSON/YAML 配置
     - 数据库存储
   - 规则元数据（名称、类型、参数）

3. **Rule Parser（规则解析器）**
   - 解析规则定义
   - 验证规则配置的正确性
   - 生成可执行的 Rule 实例

4. **Rule Registry（规则注册中心）**
   - 管理用户自定义规则类型
   - 支持规则的注册和查找
   - 防止重复注册

5. **Rule Compiler（规则编译器）**
   - 将规则定义预编译为可执行对象
   - 优化执行性能
   - 编译期错误检查

**使用场景**：
- 业务特定的验证逻辑
- 从配置文件动态加载规则
- A/B 测试不同的验证策略
- 运行时修改验证规则

### 2.4 Implementation Layer（实现层）

**职责**：实现 Core 层定义的接口，提供具体的执行逻辑。

**主要组件**：

1. **RuleEngineImpl（规则引擎实现）**
   - 实现 RuleEngine 接口
   - 管理自定义规则注册中心
   - 提供规则执行调度

2. **RuleChainExecutor（规则链执行器）**
   - 按链式调用顺序执行规则
   - 支持 **fast-fail 模式**（`Validate()`）：遇到第一个错误立即返回
   - 支持 **全量执行模式**（`ValidateAll()`）：收集所有字段的所有错误
   - 简单顺序执行，无并发

3. **ValidationResult（验证结果）**
   - 包含验证是否通过的状态
   - 单错误模式：返回第一个错误
   - 多错误模式：返回所有错误的列表
   - 提供友好的错误访问接口

4. **RuleConfigValidator（规则配置验证器）**
   - 在启动时验证自定义规则配置
   - 检查规则定义的合法性
   - 提前发现配置错误，避免运行时失败

**设计原则**：
- 简单顺序执行，易于理解和调试
- 错误快速返回或全量收集，满足不同场景
- 并发安全设计（规则对象不可变）
- 配置验证前置，减少运行时错误

### 2.5 API Layer（应用层）

**职责**：对外提供简洁易用的高层 API，屏蔽内部复杂性。用户只需要使用两种 API：**Fluent API** 和 **Custom Rule API**。

#### 2.5.1 Fluent API（链式验证 API）

**适用场景**：所有验证场景的统一接口，这是用户的主要使用方式。

提供简洁的链式调用 API，支持单个字段验证、多字段验证、规则组合：

```go
// 单个字段简单验证（fast-fail）
result := validator.
    ForField("email").
    Required().
    Email().
    Validate(user)

if !result.IsValid() {
    return result.Error()  // 返回第一个错误
}

// 单个字段复杂验证（多个规则组合）
result := validator.
    ForField("password").
    Required().
    Length(8, 32).
    Pattern(`^[a-zA-Z0-9]+$`).
    Validate(user)

// 多个字段验证（fast-fail）
result := validator.
    ForField("username").Required().Length(3, 20).
    ForField("email").Required().Email().
    ForField("age").Required().Range(18, 120).
    Validate(user)
// 按顺序验证，遇到第一个错误立即返回

// 多个字段验证（全量收集所有错误）
result := validator.
    ForField("username").Required().Length(3, 20).
    ForField("email").Required().Email().
    ForField("age").Required().Range(18, 120).
    ValidateAll(user)
// 收集所有字段的所有错误后返回

if !result.IsValid() {
    for _, err := range result.Errors() {
        fmt.Printf("字段 %s: %s\n", err.Field, err.Message)
    }
}

// 使用自定义规则
result := validator.
    ForField("customId").
    UseRule("custom_id_rule").
    Validate(data)
```

**核心方法**：

**字段指定**：
- `ForField(fieldName string) *Validator`：指定要验证的字段，可多次调用验证多个字段

**内置规则**（内部使用 Predefined Layer 的规则实现）：
- `Required() *Validator`：字段必填
- `Email() *Validator`：邮箱格式验证
- `Length(min, max int) *Validator`：字符串长度范围
- `Pattern(regex string) *Validator`：正则表达式验证
- `Min(value int) *Validator`：最小值
- `Max(value int) *Validator`：最大值
- `Range(min, max int) *Validator`：数值范围
- `URL() *Validator`：URL 格式验证
- `In(values ...interface{}) *Validator`：值在指定集合中

**自定义规则**：
- `UseRule(ruleName string) *Validator`：使用已注册的自定义规则

**执行验证**：
- `Validate(data interface{}) ValidationResult`：执行验证（fast-fail 模式）
- `ValidateAll(data interface{}) ValidationResult`：执行验证（全量收集错误模式）

**执行模式说明**：
- **Fast-fail 模式**（`Validate()`）：按链式调用顺序执行规则，遇到第一个错误立即返回，适合需要快速失败的场景
- **全量收集模式**（`ValidateAll()`）：执行所有字段的所有规则，收集所有错误后返回，适合表单验证等需要一次性展示所有错误的场景

#### 2.5.2 Custom Rule API（自定义规则 API）

**适用场景**：注册和使用业务特定的验证规则

```go
// 方式 1：注册自定义规则（函数方式）
validator.RegisterRule("custom_id", func(ctx RuleContext) RuleResult {
    value := ctx.Value()
    if isValidCustomID(value) {
        return RuleResult{Valid: true}
    }
    return RuleResult{
        Valid:   false,
        Message: "invalid custom id format",
    }
})

// 方式 2：注册自定义规则（结构体方式）
type CustomIDRule struct {
    prefix string
}

func (r *CustomIDRule) Validate(ctx RuleContext) RuleResult {
    id, ok := ctx.Value().(string)
    if !ok || !strings.HasPrefix(id, r.prefix) {
        return RuleResult{
            Valid:   false,
            Message: fmt.Sprintf("ID must start with '%s'", r.prefix),
        }
    }
    return RuleResult{Valid: true}
}

func (r *CustomIDRule) Name() string { return "custom_id" }
func (r *CustomIDRule) Description() string { return "Custom ID validation" }

// 注册
validator.RegisterRule("custom_id", &CustomIDRule{prefix: "ID-"})

// 使用自定义规则（在 Fluent API 中）
result := validator.
    ForField("userId").
    UseRule("custom_id").
    Validate(user)

// 配置验证（启动时调用，提前发现配置错误）
if err := validator.ValidateConfig(); err != nil {
    log.Fatalf("规则配置错误: %v", err)
}
```

**核心方法**：
- `RegisterRule(name string, rule Rule) error`：注册自定义规则（支持函数或实现 Rule 接口的结构体）
- `ValidateConfig() error`：验证所有已注册规则的配置是否正确（建议在启动时调用）

## 3. 核心设计考虑

### 3.1 扩展性设计

**接口隔离**
- Core 层接口保持最小化和稳定
- 不依赖具体实现，依赖倒置原则
- 新增规则类型不需要修改核心代码

**插件化机制**
- 规则注册中心支持动态加载
- 用户可实现 Rule 接口创建自定义规则
- 支持第三方规则库集成

**链式调用**
- Fluent API 通过链式调用组合多个规则
- 按链式顺序执行，逻辑清晰易懂
- 避免深层继承层次

### 3.2 灵活性设计

**多种定义方式**
- 链式方式：通过 Fluent API 直接使用内置规则
- 函数方式：注册自定义验证函数
- 结构体方式：实现 Rule 接口的自定义规则
- 配置方式：JSON/YAML 定义规则（可选）

**动态配置**
- 支持运行时注册自定义规则
- 支持从配置文件加载规则（可选）
- 配置验证前置，启动时发现错误

**执行模式选择**
- Fast-fail 模式：快速失败，适合 API 参数验证
- 全量收集模式：收集所有错误，适合表单验证

### 3.3 性能优化

**简单高效**
- 顺序执行，无并发开销
- 规则对象不可变，无锁设计
- Fast-fail 模式快速返回

**配置验证**
- 规则配置在启动时验证
- 提前发现配置错误
- 避免运行时解析和验证开销

### 3.4 错误处理

**分层错误模型**
- 错误分类清晰（验证错误、配置错误）
- 错误上下文完整（字段名、字段值、错误消息）
- ValidationResult 统一的结果模型

**错误信息**
- 支持国际化（i18n）
- 可自定义错误消息
- 提供详细的字段和错误原因

**错误收集**
- Fast-fail 模式：返回第一个错误
- 全量收集模式：收集所有字段的所有错误
- 结构化错误输出，易于前端展示

### 3.5 并发安全

**不可变对象**
- 预定义规则是不可变的
- 自定义规则建议设计为不可变
- 无锁设计，天然支持并发

**并发友好**
- Validator 实例可被多个 goroutine 安全使用
- 规则注册建议在启动时完成
- 验证执行无副作用，线程安全

## 4. 使用场景

### 4.1 简单验证（Fast-fail 模式）
使用链式 API，快速验证单个或多个字段：

```go
// 单个字段验证
result := validator.ForField("email").Required().Email().Validate(user)
if !result.IsValid() {
    return result.Error()  // 返回第一个错误
}

// 多规则组合
result := validator.
    ForField("password").
    Required().
    Length(8, 32).
    Pattern(`^[a-zA-Z0-9]+$`).
    Validate(user)

// 多字段验证（遇到第一个错误立即返回）
result := validator.
    ForField("username").Required().Length(3, 20).
    ForField("email").Required().Email().
    ForField("age").Required().Range(18, 120).
    Validate(user)
```

**典型应用**：
- API 参数验证
- 单个字段验证
- 需要快速失败的场景

### 4.2 表单验证（全量收集模式）
使用 `ValidateAll()` 收集所有错误，一次性展示：

```go
// 用户注册表单验证
result := validator.
    ForField("username").Required().Length(3, 20).
    ForField("email").Required().Email().
    ForField("password").Required().Length(8, 32).Pattern(`^[a-zA-Z0-9]+$`).
    ForField("age").Required().Range(18, 120).
    ValidateAll(user)  // 收集所有字段的所有错误

if !result.IsValid() {
    // 返回所有错误给前端
    errors := make(map[string]string)
    for _, err := range result.Errors() {
        errors[err.Field] = err.Message
    }
    return c.JSON(http.StatusBadRequest, errors)
}
```

**典型应用**：
- 用户注册表单验证
- 复杂业务规则验证
- 需要一次性展示所有错误的场景

### 4.3 自定义规则
使用 Custom Rule API 实现业务特定逻辑：

```go
// 注册自定义规则（函数方式）
validator.RegisterRule("custom_id", func(ctx RuleContext) RuleResult {
    id, ok := ctx.Value().(string)
    if !ok || !strings.HasPrefix(id, "ID-") || len(id) != 10 {
        return RuleResult{
            Valid:   false,
            Message: "ID must start with 'ID-' and be 10 characters long",
        }
    }
    return RuleResult{Valid: true}
})

// 使用自定义规则
result := validator.
    ForField("userId").
    UseRule("custom_id").
    Validate(user)
```

**典型应用**：
- 业务特定的验证逻辑
- 复杂的跨字段验证
- 需要调用外部服务的验证（如检查用户名是否已存在）

### 4.4 动态规则配置（可选）
从配置文件加载规则（可选特性）：

```go
// 从配置文件加载规则定义
func LoadCustomRulesFromConfig(validator *Validator) error {
    rules := LoadRulesFromYAML("validation_rules.yaml")

    for name, ruleDef := range rules {
        if err := validator.RegisterRule(name, ruleDef); err != nil {
            return err
        }
    }

    // 验证配置是否正确
    return validator.ValidateConfig()
}

// 启动时加载
if err := LoadCustomRulesFromConfig(validator); err != nil {
    log.Fatalf("加载规则配置失败: %v", err)
}

// 使用动态加载的规则
result := validator.
    ForField("field1").UseRule("rule_from_config").
    Validate(data)
```

**典型应用**：
- 配置文件驱动的验证
- 不同环境使用不同验证规则
- 运行时动态调整验证策略

## 5. 技术选型建议

### 5.1 配置格式（可选）
如果需要从配置文件加载规则：
- **YAML**：推荐，简洁、支持注释、适合复杂配置
- **JSON**：通用、易读、与前端交互友好
- **TOML**：Go 友好、类型明确

### 5.2 反射库
- **reflect**：Go 标准库，用于从结构体提取字段值
- 考虑使用字段标签（struct tags）定义验证规则（可选扩展）

### 5.3 可选增强
- **规则测试**：规则单元测试辅助工具
- **规则监控**：规则执行性能监控（可选）
- **国际化支持**：错误消息多语言
- **Struct Tags**：通过标签定义验证规则（如 `validate:"required,email"`）

## 6. 开发计划

### Phase 1: Core + Predefined（核心基础）
**目标**：建立核心抽象和基础规则

- 实现 Core 层接口和模型
  - Rule 接口
  - RuleContext 和 RuleResult
  - ValidationResult 错误模型
- 实现 Predefined 规则（底层实现）
  - 基础规则（Required, NotNull, NotEmpty）
  - 字符串规则（Email, Length, Pattern, URL）
  - 数值规则（Min, Max, Range）

### Phase 2: Fluent API（用户接口）
**目标**：提供简洁的链式调用接口

- 实现 Fluent API
  - ForField()
  - 内置规则方法（Required, Email, Length, Pattern 等）
  - Validate()（fast-fail 模式）
  - ValidateAll()（全量收集模式）
- 支持多字段验证
- 错误结果友好展示

### Phase 3: Custom Rule（自定义扩展）
**目标**：支持自定义规则

- Custom Rule API
  - RegisterRule()
  - UseRule()
  - 函数式自定义规则
  - 结构体式自定义规则
- 规则注册中心
- 规则配置验证器（ValidateConfig）

### Phase 4: 高级特性（可选）
**目标**：生产就绪和可选增强

- 配置文件加载（可选）
  - JSON/YAML 配置解析
  - 规则定义到 Rule 的转换
- 国际化支持（可选）
  - 错误消息多语言
  - 自定义错误消息
- Struct Tags 支持（可选）
  - 通过标签定义验证规则
- 规则监控（可选）
  - 执行时间统计
  - 错误率监控

## 7. 总结

该规则引擎设计具有以下特点：

### 核心优势
- ✅ **架构清晰**：Core → Predefined/Definition → Implementation → API，职责分明
- ✅ **API 极简**：用户只需掌握两种 API
  - **Fluent API**：链式调用，覆盖 90% 的验证场景
  - **Custom Rule API**：灵活的自定义规则支持
- ✅ **使用简单**：一行代码完成验证，学习成本低
- ✅ **执行清晰**：顺序执行，按链式调用顺序，易于理解和调试
- ✅ **灵活执行**：支持 fast-fail 和全量收集两种模式

### 设计特点
- ✅ **扩展性强**：支持自定义规则和动态配置（可选）
- ✅ **并发安全**：规则对象不可变，无锁设计，天然支持并发
- ✅ **错误友好**：结构化错误输出，易于前端展示
- ✅ **配置验证**：启动时验证规则配置，避免运行时错误

### 适用场景
- ✅ **API 参数验证**：fast-fail 模式，快速失败
- ✅ **表单验证**：全量收集模式，一次性展示所有错误
- ✅ **业务规则验证**：自定义规则，灵活扩展
- ✅ **配置驱动验证**：动态加载规则（可选）

这个架构在简洁性和灵活性之间取得了良好平衡，既能满足常见的验证需求，又能通过自定义规则支持复杂场景，是一个**生产级别**的规则引擎设计。