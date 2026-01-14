# Rules - Go 规则引擎

一个可扩展的 Go 语言规则引擎，支持内置验证规则和用户自定义规则，提供简洁的链式调用 API。

## 特性

- **链式调用 API** - 简洁直观的 Fluent API
- **丰富的内置规则** - 覆盖常见验证场景
- **自定义规则支持** - 轻松扩展业务规则
- **两种验证模式** - Fast-fail 和全量收集
- **并发安全** - 规则对象不可变，无锁设计

## 快速开始

```go
import "universal-service-user/rules"

// 单字段验证
result := rules.ForField("Email").
    Required().
    Email().
    Validate(user)

if !result.IsValid() {
    fmt.Println(result.Error())
}
```

## 验证模式

### Fast-fail 模式

遇到第一个错误立即返回，适合 API 参数验证：

```go
result := rules.ForField("Username").Required().Length(3, 20).
    ForField("Email").Required().Email().
    Validate(user)

if !result.IsValid() {
    return result.Error() // 返回第一个错误
}
```

### 全量收集模式

收集所有错误后返回，适合表单验证：

```go
result := rules.ForField("Username").Required().Length(3, 20).
    ForField("Email").Required().Email().
    ForField("Age").Required().Range(18, 120).
    ValidateAll(user)

if !result.IsValid() {
    for _, err := range result.Errors() {
        fmt.Printf("字段 %s: %s\n", err.Field, err.Message)
    }
}
```

## 内置规则

### 基础规则

| 方法 | 说明 |
|------|------|
| `Required()` | 必填验证 |
| `NotNull()` | 非空验证 |
| `NotEmpty()` | 非空字符串/数组验证 |

### 字符串规则

| 方法 | 说明 |
|------|------|
| `Email()` | 邮箱格式 |
| `URL()` | URL 格式 |
| `Length(min, max)` | 长度范围 |
| `MinLength(min)` | 最小长度 |
| `MaxLength(max)` | 最大长度 |
| `Pattern(regex)` | 正则匹配 |
| `Contains(substr)` | 包含子串 |
| `StartsWith(prefix)` | 前缀匹配 |
| `EndsWith(suffix)` | 后缀匹配 |
| `Alpha()` | 只包含字母 |
| `Alphanumeric()` | 只包含字母和数字 |
| `Numeric()` | 只包含数字字符 |
| `UUID()` | UUID 格式 |
| `IP()` / `IPv4()` / `IPv6()` | IP 地址 |
| `Phone()` | 手机号（中国大陆） |
| `JSON()` | JSON 格式 |
| `Lowercase()` | 小写字母 |
| `Uppercase()` | 大写字母 |

### 数值规则

| 方法 | 说明 |
|------|------|
| `Min(value)` | 最小值 |
| `Max(value)` | 最大值 |
| `Range(min, max)` | 数值范围 |
| `Positive()` | 正数 |
| `Negative()` | 负数 |
| `NonNegative()` | 非负数 |

### 集合规则

| 方法 | 说明 |
|------|------|
| `In(values...)` | 值在集合中 |
| `NotIn(values...)` | 值不在集合中 |
| `ArrayLength(min, max)` | 数组长度范围 |
| `MinArrayLength(min)` | 数组最小长度 |
| `MaxArrayLength(max)` | 数组最大长度 |
| `Unique()` | 数组元素唯一 |

### 时间规则

| 方法 | 说明 |
|------|------|
| `Before(time)` | 早于指定时间 |
| `BeforeNow()` | 早于当前时间 |
| `After(time)` | 晚于指定时间 |
| `AfterNow()` | 晚于当前时间 |
| `BetweenTime(start, end)` | 时间范围 |
| `DateFormat(format)` | 日期格式 |

### 比较规则

| 方法 | 说明 |
|------|------|
| `Equals(value)` | 等于指定值 |
| `NotEquals(value)` | 不等于指定值 |
| `EqualsField(field)` | 等于另一字段（如密码确认） |
| `NotEqualsField(field)` | 不等于另一字段 |
| `GreaterThanField(field)` | 大于另一字段 |
| `LessThanField(field)` | 小于另一字段 |
| `RequiredIf(field, value)` | 条件必填 |
| `RequiredUnless(field, value)` | 条件非必填 |

## 自定义规则

### 函数方式

```go
// 注册自定义规则
rules.RegisterFunc("custom_id", "Custom ID validation",
    func(ctx rules.RuleContext) rules.RuleResult {
        id, ok := ctx.Value().(string)
        if !ok || !strings.HasPrefix(id, "ID-") {
            return rules.Failure(ctx.FieldName(), "custom_id", "ID must start with 'ID-'")
        }
        return rules.Success()
    })

// 使用自定义规则
customRule, _ := rules.GetRule("custom_id")
result := rules.ForField("UserId").Rule(customRule).Validate(data)
```

### 结构体方式

```go
// 定义规则结构体
type CustomIDRule struct {
    prefix string
}

func (r *CustomIDRule) Name() string        { return "custom_id" }
func (r *CustomIDRule) Description() string { return "Custom ID validation" }

func (r *CustomIDRule) Validate(ctx rules.RuleContext) rules.RuleResult {
    id, ok := ctx.Value().(string)
    if !ok || !strings.HasPrefix(id, r.prefix) {
        return rules.Failure(ctx.FieldName(), r.Name(),
            fmt.Sprintf("ID must start with '%s'", r.prefix))
    }
    return rules.Success()
}

// 注册并使用
rules.Register(&CustomIDRule{prefix: "ID-"})
```

## 验证结果

```go
result := rules.ForField("Email").Required().Email().ValidateAll(user)

// 检查是否有效
if result.IsValid() {
    // 验证通过
}

// 获取第一个错误
err := result.Error()

// 获取所有错误
errors := result.Errors()

// 获取错误映射 (field -> messages)
errorMap := result.ErrorMap()

// 获取每个字段的第一个错误
firstErrorMap := result.FirstErrorMap()

// 获取错误字符串
errorString := result.ErrorString()
```

## 使用场景示例

### 用户注册

```go
result := rules.ForField("Username").Required().Alphanumeric().Length(3, 20).
    ForField("Email").Required().Email().
    ForField("Phone").Required().Phone().
    ForField("Password").Required().MinLength(8).
    ForField("ConfirmPassword").Required().EqualsField("Password").
    ForField("Age").Required().Range(18, 120).
    ForField("Birthday").BeforeNow().
    ValidateAll(req)
```

### API 参数验证

```go
result := rules.ForField("ID").Required().UUID().
    ForField("Status").Required().In("active", "inactive", "pending").
    ForField("Page").Required().Positive().
    ForField("PageSize").Required().Range(1, 100).
    Validate(req)
```

### 条件验证

```go
result := rules.ForField("PaymentType").Required().In("credit", "debit", "cash").
    ForField("CardNumber").RequiredIf("PaymentType", "credit").Length(16, 16).
    ForField("ExpiryDate").RequiredIf("PaymentType", "credit").
    ValidateAll(req)
```

## 目录结构

```
rules/
├── go.mod
├── rules.go              # 统一入口
├── README.md             # 文档
├── core/                 # 核心层
│   ├── rule.go           # Rule 接口
│   ├── context.go        # RuleContext
│   ├── result.go         # ValidationResult
│   └── errors.go         # 错误类型
├── predefined/           # 预定义规则
│   ├── base.go           # 基础规则
│   ├── string.go         # 字符串规则
│   ├── numeric.go        # 数值规则
│   ├── collection.go     # 集合规则
│   ├── time.go           # 时间规则
│   └── compare.go        # 比较规则
├── definition/           # 自定义规则
│   ├── custom.go         # CustomRule
│   └── registry.go       # 规则注册中心
└── api/                  # API 层
    ├── validator.go      # Fluent API
    └── helper.go         # 辅助函数
```

## 设计原则

- **接口最小化** - 核心接口保持简洁稳定
- **不可变对象** - 规则对象不可变，并发安全
- **链式组合** - 通过链式调用组合多个规则
- **扩展灵活** - 支持自定义规则和动态注册
