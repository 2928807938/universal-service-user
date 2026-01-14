# 通用用户服务设计文档

## 一、项目概述

### 1.1 项目定位

通用用户服务（Universal User Service）是一个可复用的用户管理解决方案，提供完整的用户认证、授权和账户管理能力。

### 1.2 使用方式

| 方式 | 说明 | 适用场景 |
|------|------|----------|
| **API 模式** | 通过 HTTP/RPC 调用远程服务 | 微服务架构、前后端分离项目 |
| **SDK 模式** | 引入 Go 包，调用函数 | Go 语言项目、需要深度集成 |
| **一站式模式** | 直接部署，前端直接对接 | 快速启动、独立用户中心 |

### 1.3 核心原则

- **自动建表**：服务提供完整表结构，使用 API/SDK 时自动执行建表 SQL
- **数据库兼容**：支持多种数据库

| 数据库 | 支持状态 |
|--------|----------|
| PostgreSQL | ✅ 完全支持 |
| MySQL | ✅ 完全支持 |
| SQLite | ✅ 完全支持 |

- **可扩展钩子**：关键操作支持前置/后置钩子
- **多厂商适配**：短信、第三方登录等支持多厂商配置
- **DDD 架构**：保持领域驱动设计的清晰边界

---

## 二、功能模块划分

```
┌─────────────────────────────────────────────────────────────┐
│                     通用用户服务                              │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │  用户核心    │  │  认证模块    │  │   第三方登录模块     │  │
│  │  - 注册     │  │  - 登录      │  │   - 微信           │  │
│  │  - 信息管理  │  │  - Token    │  │   - 支付宝         │  │
│  │  - 状态管理  │  │  - 会话管理  │  │   - QQ  
│  │                                    - 自定义适配
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
│                                                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │  验证码模块  │  │  邮箱模块    │  │    短信模块         │  │
│  │  - 生成     │  │  - 发送邮件  │  │   - 腾讯云          │  │
│  │  - 校验     │  │  - 模板管理  │  │   - 阿里云          │  │
│  │  - 过期管理  │  │  - 限频控制  │  │   - 自定义适配      │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐    │
│  │                    钩子系统                          │    │
│  │   BeforeCreate / AfterCreate                       │    │
│  └─────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
```

---

## 三、模块详细设计

### 3.1 用户核心模块

#### 功能清单

| 功能 | 描述 |
|------|------|
| 用户注册 | 通过邮箱或手机号注册账号 |
| 用户登录 | 支持用户名+密码、邮箱+密码、手机号+验证码、微信、支付宝、QQ |
| 修改密码 | 已登录用户修改密码 |
| 忘记密码 | 通过邮箱或手机号找回密码 |
| 用户信息管理 | 查询、更新用户信息 |

#### 注册流程

```
┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐
│  输入信息  │────▶│ 发送验证码 │────▶│ 验证验证码 │────▶│ 唯一性校验  │────▶│  创建用户 │
└──────────┘     └──────────┘     └──────────┘     └──────────┘     └──────────┘
                      │                │                  │                 │
                      ▼                ▼                  ▼                 ▼
                 邮箱/短信          校验通过          检查邮箱/手机号     触发钩子
                                                      是否已被绑定
```

**校验规则：**
- 邮箱已被其他账号绑定 → 返回错误码 10006
- 手机号已被其他账号绑定 → 返回错误码 10007

#### 忘记密码流程

```
┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐
│ 输入账号  │────▶│ 发送验证码 │────▶│ 验证验证码 │────▶│  重置密码 │
└──────────┘     └──────────┘     └──────────┘     └──────────┘
    │                                                    │
    ▼                                                    ▼
 邮箱/手机号                                          触发钩子
```

#### 更换邮箱流程

```
┌────────────┐     ┌────────────┐     ┌────────────┐
│ 输入旧邮箱  │────▶│发送验证码到旧邮箱│────▶│ 验证旧邮箱  │
└────────────┘     └────────────┘     └────────────┘
                                               │
                                               ▼
                                        ┌────────────┐
                                        │  验证通过    │
                                        └────────────┘
                                               │
                                               ▼
┌────────────┐     ┌────────────┐     ┌────────────┐     ┌────────────┐
│ 输入新邮箱  │────▶│发送验证码到新邮箱│────▶│ 验证新邮箱  │────▶│ 唯一性校验  │
└────────────┘     └────────────┘     └────────────┘     └────────────┘
                                                            │
                                                            ▼
                                                     ┌────────────┐
                                                     │  更新邮箱   │
                                                     └────────────┘
                                                            │
                                                            ▼
                                                       触发钩子
```

**安全说明：**
- 需同时验证旧邮箱和新邮箱的所有权
- 两个验证码都验证通过后，执行唯一性校验
- **新邮箱已被其他账号绑定 → 返回错误码 10006，操作终止**
- 防止恶意修改他人邮箱

#### 更换手机号流程

```
┌────────────┐     ┌────────────┐     ┌────────────┐
│ 输入旧手机号 │────▶│发送短信到旧手机号│────▶│ 验证旧手机号 │
└────────────┘     └────────────┘     └────────────┘
                                               │
                                               ▼
                                        ┌────────────┐
                                        │  验证通过    │
                                        └────────────┘
                                               │
                                               ▼
┌────────────┐     ┌────────────┐     ┌────────────┐     ┌────────────┐
│ 输入新手机号 │────▶│发送短信到新手机号│────▶│ 验证新手机号 │────▶│ 唯一性校验  │
└────────────┘     └────────────┘     └────────────┘     └────────────┘
                                                            │
                                                            ▼
                                                     ┌────────────┐
                                                     │  更新手机号  │
                                                     └────────────┘
                                                            │
                                                            ▼
                                                       触发钩子
```

**安全说明：**
- 需同时验证旧手机号和新手机号的所有权
- 两个验证码都验证通过后，执行唯一性校验
- **新手机号已被其他账号绑定 → 返回错误码 10007，操作终止**
- 防止恶意修改他人手机号

#### 登录流程

支持六种登录方式：

```
                      ┌─────────────────────────────────────────────────────────────┐
                      │                          用户登录                              │
                      └─────────────────────────────────────────────────────────────┘
                                               │
     ┌─────────────────┬─────────────────┬─────┴─────┬─────────────────┬─────────────────┐
     │                 │                 │           │                 │                 │
     ▼                 ▼                 ▼           ▼                 ▼                 ▼
┌──────────┐    ┌──────────┐    ┌──────────┐  ┌──────────┐    ┌──────────┐    ┌──────────┐
│用户名+密码 │    │ 邮箱+密码 │    │手机号+验证码│  │  微信登录 │    │ 支付宝登录 │    │  QQ登录   │
│  登录    │    │   登录   │    │   登录    │  │(OAuth 2.0)│    │(OAuth 2.0)│    │(OAuth 2.0)│
└──────────┘    └──────────┘    └──────────┘  └────────────┘    └──────────┘    └──────────┘
     │                 │                 │           │                 │                 │
     └─────────────────┴─────────────────┴───────────┴─────────────────┴─────────────────┘
                                               ▼
                                      ┌──────────────────┐
                                      │    验证通过       │
                                      │   生成Token      │
                                      └──────────────────┘
                                               ▼
                                      ┌──────────────────┐
                                      │  AfterLogin钩子   │
                                      │  登录日志/设备记录│
                                      └──────────────────┘
```

**说明**：
- 基础登录方式（用户名、邮箱、手机号）和第三方登录（微信、支付宝、QQ）均支持
- 第三方登录走 OAuth 2.0 授权流程，回调后绑定或创建账号
- 所有登录方式均触发 `BeforeLogin` 和 `AfterLogin` 钩子

#### 值对象设计

**邮箱值对象（Email）**

| 方法 | 说明 |
|------|------|
| `NewEmail(email string)` | 创建邮箱，自动格式化并校验 |
| `String() string` | 获取邮箱字符串 |
| `Domain() string` | 获取邮箱域名 |
| `LocalPart() string` | 获取@前本地部分 |

**手机号值对象（Phone）**

| 方法 | 说明 |
|------|------|
| `NewPhone(phone string, region string)` | 创建手机号，校验格式 |
| `String() string` | 获取格式化后的手机号 |
| `IsValid() bool` | 校验手机号是否有效 |
| `Region() string` | 获取国家/地区代码 |

**手机号校验规则（中国大陆）**

```
格式：^1[3-9]\d{9}$
说明：以1开头，第二位为3-9，共11位数字
示例：13800138000、15912345678
```

---

### 3.2 验证码模块

#### 核心规则

| 规则 | 说明 |
|------|------|
| 有效期 | 5 分钟 |
| 发送间隔 | 60 秒 |
| 匹配校验 | 邮箱/手机号 + 验证码 + 场景 |
| 使用次数 | 在有效期内只能使用一次，验证成功后失效；超过有效期后不能再使用 |

#### 验证码场景枚举

| 场景 | 用途 | 发送方式 |
|------|------|----------|
| `REGISTER` | 注册 | 邮箱/短信 |
| `LOGIN` | 登录 | 邮箱/短信 |
| `RESET_PASSWORD` | 重置密码 | 邮箱/短信 |
| `CHANGE_EMAIL_OLD` | 更换邮箱：验证旧邮箱 | 邮箱 |
| `CHANGE_EMAIL_NEW` | 更换邮箱：验证新邮箱 | 邮箱 |
| `CHANGE_PHONE_OLD` | 更换手机号：验证旧手机号 | 短信 |
| `CHANGE_PHONE_NEW` | 更换手机号：验证新手机号 | 短信 |

#### 存储设计

**验证码仅存储在 Redis 中，不需要数据库表。**

Redis Key 设计：

```
verification_code:{type}:{target}:{scene}
```

示例：
- `verification_code:email:test@example.com:register`
- `verification_code:phone:13800138000:reset_password`

Redis 数据结构（Hash）：

| 字段 | 说明 |
|------|------|
| code | 验证码（加密存储） |
| expires_at | 过期时间戳 |
| verified | 是否已验证 |

利用 Redis 的 TTL 机制自动处理过期。

---

### 3.3 邮箱模块

#### 功能清单

| 功能 | 描述 |
|------|------|
| 发送验证码 | 发送验证码邮件 |
| 限频控制 | 60 秒内不可重复发送 |
| 模板管理 | 支持自定义邮件模板 |

#### 配置项

```yaml
email:
  smtp:
    host: "smtp.example.com"
    port: 587
    username: "noreply@example.com"
    password: "***"
    from: "系统通知 <noreply@example.com>"
  templates:
    register: "您的注册验证码是：{code}，有效期5分钟"
    reset_password: "您的重置密码验证码是：{code}，有效期5分钟"
  rate_limit:
    interval: 60  # 秒
```

---

### 3.4 短信模块

#### 设计思路

采用**适配器模式 + 工厂方法模式 + 策略模式**，定义统一的短信发送接口，各厂商实现各自的适配器。

#### 架构设计

```
┌─────────────────────────────────────────────────────────────┐
│                     SMS 管理层                              │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │ Provider    │  │  Registry   │  │      Factory        │  │
│  │  接口定义    │  │  注册中心    │  │     Provider工厂    │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
├────────────────────────────────────────────────────────────┤
│                     内置 Providers                           │
│  ┌──────────┐ ┌──────────┐                                 │
│  │ 腾讯云   │ │  阿里云   │                                 │
│  └──────────┘ └──────────┘                                 │
├────────────────────────────────────────────────────────────┤
│                    用户自定义扩展                            │
│  ┌───────────────────────────────────────────────────────┐    │
│  │    用户可实现 Provider 接口添加自定义短信平台          │    │
│  └─────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
```

#### 支持厂商

| 厂商 | 说明 |
|------|------|
| 腾讯云 | Tencent Cloud SMS |
| 阿里云 | Alibaba Cloud SMS |
| 自定义 | 用户可实现自定义适配器 |

#### Provider 接口定义

| 方法 | 说明 |
|------|------|
| GetName | 获取厂商唯一标识 |
| GetDisplayName | 获取厂商显示名称 |
| Send | 发送短信 |
| ValidateConfig | 验证配置是否完整有效 |

#### 目录结构

```
share/sms/
├── provider.go              # Provider 接口定义
├── registry.go              # Provider 注册中心
├── factory.go               # Provider 工厂
├── config.go                # 通用配置结构
└── providers/               # 内置 Provider 实现
    ├── tencent.go           # 腾讯云
    └── aliyun.go            # 阿里云
```

#### 配置模板（腾讯云示例）

```yaml
sms:
  provider: "tencent"  # tencent / aliyun / custom
  tencent:
    secret_id: "your-secret-id"
    secret_key: "your-secret-key"
    app_id: "your-app-id"
    sign: "your-sign"
    templates:
      register: "template-id-for-register"
      reset_password: "template-id-for-reset"
  rate_limit:
    interval: 60  # 秒
```

#### 配置模板（阿里云示例）

```yaml
sms:
  provider: "aliyun"
  aliyun:
    access_key_id: "your-access-key-id"
    access_key_secret: "your-access-key-secret"
    sign_name: "your-sign-name"
    templates:
      register: "SMS_123456789"
      reset_password: "SMS_987654321"
```

#### 适配器接口

```
SMSProvider (接口)
├── Send(phone, templateID, params) error
├── GetProviderName() string
└── Validate() error

实现类：
├── TencentSMSProvider
├── AliyunSMSProvider
└── CustomSMSProvider
```

---

### 3.5 第三方登录模块

#### 设计思路

采用**适配器模式**，定义统一的 OAuth Provider 接口，各平台实现各自的适配器。使用者按需配置需要的平台，也支持自定义实现。

#### 架构设计

```
┌─────────────────────────────────────────────────────────────┐
│                    OAuth 管理层                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │ Provider    │  │  Registry   │  │      Factory        │  │
│  │  接口定义    │  │  注册中心    │  │     Provider工厂    │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
├─────────────────────────────────────────────────────────────┤
│                     内置 Providers                           │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐                     │
│  │  微信    │ │  支付宝   │ │    QQ    │                     │
│  └──────────┘ └──────────┘ └──────────┘                     │
├─────────────────────────────────────────────────────────────┤
│                    用户自定义扩展                            │
│  ┌───────────────────────────────────────────────────────┐    │
│  │    用户可实现 Provider 接口添加自定义 OAuth 平台      │    │
│  └─────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
```

#### Provider 接口定义

| 方法 | 说明 |
|------|------|
| GetName | 获取平台唯一标识 |
| GetDisplayName | 获取平台显示名称 |
| GetAuthURL | 生成授权跳转 URL |
| ExchangeToken | 用授权码换取 Access Token |
| GetUserInfo | 获取第三方用户信息 |
| ValidateConfig | 验证配置是否完整有效 |

#### 目录结构

```
share/oauth/
├── provider.go              # Provider 接口定义
├── registry.go              # Provider 注册中心
├── factory.go               # Provider 工厂
├── config.go                # 通用配置结构
├── types.go                 # 用户信息、令牌等通用类型
└── providers/               # 内置 Provider 实现
    ├── wechat.go            # 微信
    ├── alipay.go            # 支付宝
    └── qq.go                # QQ
```

#### 登录流程

```
┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐
│ 发起授权  │────▶│ 跳转平台  │────▶│ 用户授权  │────▶│ 回调处理 │
└──────────┘     └──────────┘     └──────────┘     └──────────┘
                                                        │
                      ┌─────────────────────────────────┘
                      ▼
               ┌──────────────┐
               │ Provider获取  │
               │  用户信息    │
               └──────────────┘
                      │
        ┌─────────────┴─────────────┐
        ▼                           ▼
   ┌──────────┐              ┌──────────┐
   │ 已绑定账号 │              │ 未绑定账号 │
   └──────────┘              └──────────┘
        │                           │
        ▼                           ▼
   直接登录                   创建/绑定账号
```

#### 内置平台支持

| 平台 | 协议 | 状态 |
|------|------|------|
| 微信 | OAuth 2.0 | 内置 |
| 支付宝 | OAuth 2.0 | 内置 |
| QQ | OAuth 2.0 | 内置 |

#### 扩展支持

使用者可通过实现 Provider 接口添加其他平台，例如：
- GitHub
- Google
- 钉钉
- 企业微信
- 飞书
- 或其他自定义 OAuth 平台

#### 配置方式

**配置文件方式**（适用于常见平台）：

使用平台名称作为 Key，配置对应参数：

| 参数 | 说明 |
|------|------|
| enabled | 是否启用该平台 |
| app_id / client_id | 应用 ID |
| app_secret / client_secret | 应用密钥 |
| redirect_uri | 回调地址 |
| scopes | 授权范围 |

**代码注册方式**（适用于自定义平台）：

通过注册中心动态添加自定义 Provider

#### 用户信息统一格式

各 Provider 返回的用户信息统一转换为以下格式：

| 字段 | 说明 |
|------|------|
| provider | 平台标识 |
| open_id | 平台用户唯一标识 |
| union_id | 联合标识（可选，如微信） |
| nickname | 昵称 |
| avatar | 头像 URL |
| raw | 原始数据（用于扩展） |

#### 自定义 Provider 扩展

使用者可通过以下方式添加自定义平台：

1. 实现 Provider 接口
2. 通过注册中心注册
3. 在配置文件中引用

支持场景：
- 企业内部 SSO 系统
- 未开源的第三方平台
- 特殊定制的认证流程

---

### 3.6 钩子系统

#### 设计思路

在关键操作的前后提供扩展点，使用者可以注入自定义逻辑。

#### 钩子类型

| 钩子 | 触发时机 | 用途示例 |
|------|----------|----------|
| `BeforeCreate` | 创建用户前 | 参数校验、数据补充 |
| `AfterCreate` | 创建用户后 | 发送欢迎邮件、初始化数据 |
| `BeforeUpdate` | 更新用户前 | 权限校验、日志记录 |
| `AfterUpdate` | 更新用户后 | 同步外部系统 |
| `BeforeDelete` | 删除用户前 | 关联数据检查 |
| `AfterDelete` | 删除用户后 | 清理关联数据 |
| `BeforeLogin` | 登录前 | 账号状态检查、风控 |
| `AfterLogin` | 登录后 | 登录日志、设备记录 |
| `BeforeResetPassword` | 重置密码前 | 安全校验 |
| `AfterResetPassword` | 重置密码后 | 通知用户、踢出其他会话 |

#### 钩子上下文

钩子函数接收一个上下文对象，包含：

| 字段 | 说明 |
|------|------|
| `Context` | Go 标准上下文 |
| `User` | 用户实体（如果存在） |
| `Request` | 原始请求数据 |
| `Metadata` | 自定义元数据 |

#### 使用方式

```
SDK 模式：直接注册钩子函数
API 模式：通过 Webhook 回调或消息队列
```

---

## 四、数据模型设计

### 4.0 自动建表机制

#### 设计思路

服务提供预定义的表结构 SQL 文件，使用者在使用 SDK 或启动 API 服务时，系统会自动检测并创建缺失的表。

#### 实现方式

| 步骤 | 描述 |
|------|------|
| 1. 检测表存在 | 启动时检查必需的表是否存在 |
| 2. 自动建表 | 表不存在则执行对应的建表 SQL |
| 3. 版本管理 | 支持表结构版本号，自动迁移 |
| 4. 多数据库支持 | 根据数据库类型选择对应的 SQL 脚本 |

#### SQL 文件组织

```
share/migrations/
├── postgres/
│   ├── 001_users.up.sql
│   ├── 001_users.down.sql
│   └── 002_user_oauth_bindings.up.sql
├── mysql/
│   └── ...
└── sqlite/
    └── ...
```

#### SDK 使用示例

```go
// 初始化时自动建表
userService := usersdk.New(
    usersdk.WithDatabase(db),
    usersdk.WithAutoMigrate(true),  // 自动建表
)
```

#### API 服务使用示例

```go
// 服务启动时自动建表
server := NewServer(
    WithDatabase(db),
    WithAutoMigrate(true),
)
server.Run()
```

### 4.1 通用审计字段

所有业务表都应包含以下通用审计字段，用于跟踪数据的创建、修改和软删除。

#### 字段定义

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INT | 主键，自增 |
| created_at | TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 更新时间 |
| deleted_at | TIMESTAMP | 软删除时间（NULL 表示未删除） |
| version | INT | 版本号（乐观锁，默认为 1，每次更新递增） |

#### 字段说明

| 字段 | 说明 |
|------|------|
| `id` | 自增主键，性能更好 |
| `created_at` | 记录数据创建时间，创建后不可修改 |
| `updated_at` | 记录数据最后更新时间，每次更新自动刷新 |
| `deleted_at` | 软删除标记，NULL 表示有效，非 NULL 表示已删除 |
| `version` | 乐观锁版本号，更新时检查版本号防止并发冲突 |

#### 使用示例

**表结构定义时引用审计字段：**

```
users 表字段：
├── 业务字段
│   ├── username
│   ├── email
│   ├── phone
│   ├── password_hash
│   └── status
└── 审计字段（继承通用审计字段）
    ├── id
    ├── created_at
    ├── updated_at
    ├── deleted_at
    └── version
```

**更新操作（乐观锁）：**

```go
// 更新时检查版本号
UPDATE users
SET username = ?, status = ?, version = version + 1, updated_at = NOW()
WHERE id = ? AND version = ?
// 如果版本号不匹配，影响行数为 0，更新失败
```

### 4.2 用户表（users）

| 字段 | 类型 | 说明 |
|------|------|------|
| username | VARCHAR(64) | 用户名（唯一，可空） |
| nickname | VARCHAR(64) | 昵称（用户显示名称） |
| avatar | VARCHAR(512) | 头像URL |
| email | VARCHAR(255) | 邮箱（唯一，可空） |
| phone | VARCHAR(20) | 手机号（唯一，可空） |
| password_hash | VARCHAR(255) | 密码哈希（可空，第三方登录用户可能无密码） |
| status | INT | 状态（0:未激活 1:正常 2:禁用 3:锁定） |
| *id | INT | 继承审计字段（主键，自增） |
| *created_at | TIMESTAMP | 继承审计字段（创建时间） |
| *updated_at | TIMESTAMP | 继承审计字段（更新时间） |
| *deleted_at | TIMESTAMP | 继承审计字段（软删除时间） |
| *version | INT | 继承审计字段（版本号） |

#### 用户状态枚举

| 状态值 | 名称 | 说明 |
|--------|------|------|
| 0 | INACTIVE | 未激活 |
| 1 | ACTIVE | 正常 |
| 2 | DISABLED | 禁用（管理员手动禁用） |
| 3 | LOCKED | 锁定（多次登录失败自动锁定） |

#### 索引设计

| 索引 | 字段 | 类型 |
|------|------|------|
| PRIMARY | id | 主键索引 |
| idx_username | username | 唯一索引 |
| idx_email | email | 唯一索引 |
| idx_phone | phone | 唯一索引 |
| idx_status | status | 普通索引 |
| idx_deleted_at | deleted_at | 普通索引 |

### 4.3 第三方账号绑定表（user_oauth_bindings）

| 字段 | 类型 | 说明 |
|------|------|------|
| user_id | INT | 用户ID（外键） |
| provider | VARCHAR(32) | 平台（wechat/alipay/qq） |
| open_id | VARCHAR(128) | 平台用户唯一标识 |
| union_id | VARCHAR(128) | 平台联合标识（可空） |
| nickname | VARCHAR(64) | 平台昵称 |
| avatar | VARCHAR(512) | 平台头像 |
| access_token | VARCHAR(512) | 访问令牌（加密存储） |
| refresh_token | VARCHAR(512) | 刷新令牌（加密存储） |
| expires_at | TIMESTAMP | 令牌过期时间 |
| *id | INT | 继承审计字段（主键，自增） |
| *created_at | TIMESTAMP | 继承审计字段（创建时间） |
| *updated_at | TIMESTAMP | 继承审计字段（更新时间） |
| *deleted_at | TIMESTAMP | 继承审计字段（软删除时间，可选） |
| *version | INT | 继承审计字段（版本号） |

#### 索引设计

| 索引 | 字段 | 类型 |
|------|------|------|
| PRIMARY | id | 主键索引 |
| idx_user_provider | user_id + provider | 复合索引 |
| idx_open_id | open_id | 普通索引 |

### 4.4 登录日志表（user_login_logs）

| 字段 | 类型 | 说明 |
|------|------|------|
| user_id | INT | 用户ID（外键） |
| login_type | VARCHAR(32) | 登录方式（username/email/phone/wechat/alipay/qq） |
| login_at | TIMESTAMP | 登录时间 |
| logout_at | TIMESTAMP | 登出时间（NULL 表示未登出） |
| ip | VARCHAR(64) | 登录IP |
| device_type | VARCHAR(32) | 设备类型（ios/android/web/desktop） |
| device_name | VARCHAR(128) | 设备名称（如 "iPhone 15 Pro"） |
| browser | VARCHAR(64) | 浏览器（Chrome/Safari/Firefox） |
| os | VARCHAR(32) | 操作系统（iOS/Android/Windows/macOS） |
| status | INT | 状态（1:成功 2:失败） |
| fail_reason | VARCHAR(128) | 失败原因（成功时为空） |
| *id | INT | 继承审计字段（主键，自增） |
| *created_at | TIMESTAMP | 继承审计字段（创建时间） |
| *version | INT | 继承审计字段（版本号，日志表不需要 soft delete） |

#### 登录状态枚举

| 状态值 | 名称 | 说明 |
|--------|------|------|
| 1 | SUCCESS | 登录成功 |
| 2 | FAILED | 登录失败 |

#### 索引设计

| 索引 | 字段 | 类型 |
|------|------|------|
| PRIMARY | id | 主键索引 |
| idx_user_id | user_id | 普通索引 |
| idx_user_login_at | user_id + login_at | 复合索引 |
| idx_login_at | login_at | 普通索引 |
| idx_ip | ip | 普通索引 |

### 4.5 会话表（user_sessions）

| 字段 | 类型 | 说明 |
|------|------|------|
| user_id | INT | 用户ID（外键） |
| refresh_token_hash | VARCHAR(128) | Refresh Token 哈希值（用于精确匹配会话） |
| device_type | VARCHAR(32) | 设备类型（ios/android/web/desktop） |
| device_name | VARCHAR(128) | 设备名称 |
| ip | VARCHAR(64) | 登录IP |
| last_active_at | TIMESTAMP | 最后活跃时间 |
| *id | INT | 继承审计字段（主键，自增） |
| *created_at | TIMESTAMP | 继承审计字段（创建时间） |
| *updated_at | TIMESTAMP | 继承审计字段（更新时间） |
| *version | INT | 继承审计字段（版本号，会话表不需要 soft delete） |

#### 索引设计

| 索引 | 字段 | 类型 |
|------|------|------|
| PRIMARY | id | 主键索引 |
| idx_user_id | user_id | 普通索引 |
| idx_token_hash | refresh_token_hash | 唯一索引 |
| idx_user_active | user_id + last_active_at | 复合索引 |

#### 表说明

- 记录用户当前活跃的登录会话
- 支持多设备登录管理功能
- 用户可以查看当前登录设备列表
- 支持踢出指定设备
- Refresh Token 存储在 Redis 中，TTL 自动管理过期
- 清理过期会话时，通过查询 Redis 判断 Token 是否存在

### 4.6 用户扩展信息表（user_profiles）

| 字段 | 类型 | 说明 |
|------|------|------|
| user_id | INT | 用户ID（外键，关联 users.id） |
| bio | VARCHAR(500) | 个人简介/描述 |
| gender | INT | 性别（0:未知 1:男 2:女） |
| birthday | DATE | 生日 |
| location | VARCHAR(128) | 所在地（如"北京 海淀"） |
| company | VARCHAR(128) | 公司 |
| position | VARCHAR(64) | 职位 |
| *id | INT | 继承审计字段（主键，自增） |
| *created_at | TIMESTAMP | 继承审计字段（创建时间） |
| *updated_at | TIMESTAMP | 继承审计字段（更新时间） |
| *version | INT | 继承审计字段（版本号） |

#### 性别枚举

| 状态值 | 名称 | 说明 |
|--------|------|------|
| 0 | UNKNOWN | 未知/未设置 |
| 1 | MALE | 男 |
| 2 | FEMALE | 女 |

#### 索引设计

| 索引 | 字段 | 类型 |
|------|------|------|
| PRIMARY | id | 主键索引 |
| idx_user_id | user_id | 唯一索引（一对一关系） |

#### 表说明

- 与 users 表一对一关系
- 存储用户的扩展个人信息
- 查询用户基本信息时不需要关联此表
- 只在访问个人主页时才关联查询

---

## 五、项目结构设计

```
universal-service-user/
├── go.work
├── bom/                          # BOM 依赖管理
├── share/                        # 公共组件
│   ├── errors/                   # 错误定义
│   ├── types/                    # 通用类型
│   ├── config/                   # 配置加载
│   └── middleware/               # 中间件
│
├── user/                         # 用户聚合
│   ├── domain/                   # 领域层
│   │   ├── entity/               # 领域实体
│   │   ├── repository/           # 仓储接口
│   │   ├── service/              # 领域服务
│   │   ├── valueobject/          # 值对象
│   │   ├── event/                # 领域事件
│   │   └── enum/                 # 枚举
│   └── infrastructure/           # 基础设施层
│       ├── entity/               # PO 实体
│       ├── converter/            # 转换器
│       └── repository/           # 仓储实现
│
├── verification/                 # 验证码聚合（新增）
│   ├── domain/
│   │   ├── valueobject/          # 验证码值对象
│   │   ├── repository/           # 仓储接口
│   │   └── service/              # 验证码领域服务
│   └── infrastructure/
│       └── repository/           # Redis 实现
│
├── notification/                 # 通知聚合（新增）
│   ├── domain/
│   │   ├── service/              # 通知领域服务
│   │   └── provider/             # 提供者接口
│   └── infrastructure/
│       ├── email/                # 邮箱实现
│       └── sms/                  # 短信实现
│           ├── tencent/          # 腾讯云
│           ├── aliyun/           # 阿里云
│           └── provider.go       # 提供者工厂
│
├── oauth/                        # OAuth 聚合（新增）
│   ├── domain/
│   │   ├── entity/               # OAuth 绑定实体
│   │   ├── repository/           # 仓储接口
│   │   └── service/              # OAuth 领域服务
│   └── infrastructure/
│       ├── wechat/               # 微信实现
│       ├── alipay/               # 支付宝实现
│       ├── qq/                   # QQ 实现
│       └── provider.go           # 提供者工厂
│
├── auth/                         # 认证聚合（新增）
│   ├── domain/
│   │   └── service/              # 认证领域服务（登录、Token）
│   └── infrastructure/
│       └── jwt/                  # JWT 实现
│
├── hook/                         # 钩子系统（新增）
│   ├── types.go                  # 钩子类型定义
│   ├── registry.go               # 钩子注册中心
│   └── executor.go               # 钩子执行器
│
├── api/                          # API 层
│   └── user-api/
│       ├── dto/                  # 数据传输对象
│       ├── service/              # 应用服务
│       ├── http/                 # HTTP 处理器
│       └── rpc/                  # RPC 处理器（可选）
│
├── sdk/                          # SDK 封装（新增）
│   ├── client.go                 # SDK 客户端
│   ├── options.go                # 配置选项
│   └── hooks.go                  # 钩子注册
│
└── cmd/
    └── api/                      # 主程序入口
```

---

## 六、接口设计

### 6.1 用户相关

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/users/register | 用户注册 |
| POST | /api/v1/users/login | 用户登录 |
| POST | /api/v1/users/logout | 用户登出 |
| POST | /api/v1/users/password/reset | 重置密码（忘记密码） |
| PUT | /api/v1/users/password | 修改密码（已登录） |
| GET | /api/v1/users/me | 获取当前用户信息 |
| PUT | /api/v1/users/me | 更新当前用户信息 |

#### 6.1.1 用户注册

**请求参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| email | string | 是* | 邮箱（与手机号二选一） |
| phone | string | 是* | 手机号（与邮箱二选一） |
| password | string | 是 | 密码（6-32位） |
| code | string | 是 | 验证码 |
| username | string | 否 | 用户名（默认生成） |

#### 6.1.2 用户登录

支持三种登录方式：

**方式一：用户名 + 密码**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| login_type | string | 是 | 固定值 "username" |
| username | string | 是 | 用户名 |
| password | string | 是 | 密码 |

**方式二：邮箱 + 密码**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| login_type | string | 是 | 固定值 "email" |
| email | string | 是 | 邮箱 |
| password | string | 是 | 密码 |

**方式三：手机号 + 验证码**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| login_type | string | 是 | 固定值 "phone" |
| phone | string | 是 | 手机号 |
| code | string | 是 | 验证码 |

**响应参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| access_token | string | 访问令牌 |
| refresh_token | string | 刷新令牌 |
| expires_in | int | 过期时间（秒） |
| user_info | object | 用户信息 |

### 6.2 验证码相关

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/verification/email/send | 发送邮箱验证码 |
| POST | /api/v1/verification/sms/send | 发送短信验证码 |
| POST | /api/v1/verification/verify | 验证验证码 |

### 6.3 第三方登录

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/oauth/{provider}/authorize | 获取授权 URL |
| GET | /api/v1/oauth/{provider}/callback | 授权回调 |
| POST | /api/v1/oauth/{provider}/binding | 绑定第三方账号 |
| DELETE | /api/v1/oauth/{provider}/binding | 解绑第三方账号 |

---

## 七、SDK 使用示例

### 7.1 基础使用

```go
// 初始化 SDK
userService := usersdk.New(
    usersdk.WithDatabase(db),
    usersdk.WithRedis(redis),
    usersdk.WithEmailConfig(emailConfig),
    usersdk.WithSMSProvider("tencent", tencentConfig),
)

// 注册用户
user, err := userService.Register(ctx, &RegisterRequest{
    Email:    "test@example.com",
    Password: "password123",
    Code:     "123456",
})
```

### 7.2 注册钩子

```go
// 注册创建用户后的钩子
userService.OnAfterCreate(func(ctx *hook.Context) error {
    // 发送欢迎邮件
    sendWelcomeEmail(ctx.User.Email)

    // 初始化用户默认数据
    initUserDefaults(ctx.User.ID)

    return nil
})

// 注册登录前的钩子
userService.OnBeforeLogin(func(ctx *hook.Context) error {
    // 风控检查
    if isRisky(ctx.Request) {
        return errors.New("登录风险较高，请稍后重试")
    }
    return nil
})
```

---

## 八、配置文件模板

```yaml
# config.yaml
server:
  port: 8080
  mode: "release"  # debug / release
  read_timeout: 60
  write_timeout: 60

database:
  driver: "postgres"           # postgres / mysql / sqlite
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "postgres"
  dbname: "user_service"
  max_open_conns: 100
  max_idle_conns: 10
  conn_max_lifetime: 3600      # 秒

  # 自动建表配置
  auto_migrate: true           # 是否自动建表

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0
  pool_size: 10

jwt:
  secret: "your-jwt-secret"    # 建议使用环境变量
  expire: 7200                 # Access Token 过期时间（秒）
  refresh_expire: 604800       # Refresh Token 过期时间（秒，默认7天）
  issuer: "user-service"

verification:
  code_length: 6               # 验证码长度
  expire: 300                  # 验证码过期时间（秒，5分钟）
  rate_limit: 60               # 发送间隔（秒，60秒）

email:
  enabled: true
  smtp:
    host: "smtp.example.com"
    port: 587
    username: "noreply@example.com"
    password: "***"
    from: "系统通知 <noreply@example.com>"
  templates:
    register: "您的注册验证码是：{code}，有效期5分钟"
    reset_password: "您的重置密码验证码是：{code}，有效期5分钟"
    login: "您的登录验证码是：{code}，有效期5分钟"
  rate_limit:
    interval: 60               # 秒

sms:
  enabled: false
  provider: "tencent"          # tencent / aliyun / custom
  tencent:
    secret_id: ""
    secret_key: ""
    app_id: ""
    sign: ""
    templates:
      register: ""             # 腾讯云短信模板 ID
      reset_password: ""
      login: ""
  aliyun:
    access_key_id: ""
    access_key_secret: ""
    sign_name: ""
    templates:
      register: "SMS_123456789"
      reset_password: "SMS_987654321"
      login: "SMS_111111111"
  rate_limit:
    interval: 60               # 秒

oauth:
  wechat:
    enabled: false
    app_id: ""
    app_secret: ""
    redirect_uri: "https://your-domain.com/auth/wechat/callback"
  alipay:
    enabled: false
    app_id: ""
    private_key: ""
    alipay_public_key: ""
    redirect_uri: "https://your-domain.com/auth/alipay/callback"
  qq:
    enabled: false
    app_id: ""
    app_key: ""
    redirect_uri: "https://your-domain.com/auth/qq/callback"

# 登录防刷配置
login_rate_limit:
  enabled: true
  ip:
    max_attempts: 5            # 同一 IP 最大尝试次数
    window: 60                 # 时间窗口（秒）
    block_duration: 600        # 封禁时长（秒，10分钟）
  account:
    max_failures: 5            # 同一账号最大失败次数
    window: 600                # 时间窗口（秒，10分钟）
    lock_duration: 1800        # 锁定时长（秒，30分钟）
  device:
    max_failures: 10           # 同一设备最大失败次数
    window: 3600               # 时间窗口（秒，1小时）

# Webhook 配置（API 模式下钩子回调）
webhook:
  enabled: false
  url: "https://your-domain.com/webhook/hooks"
  secret: "your-webhook-secret"  # 签名验证密钥
  events:                      # 需要回调的事件
    - "user.created"
    - "user.login"
    - "user.password_reset"
  timeout: 5                   # 请求超时（秒）
  retry: 3                     # 失败重试次数

# 日志配置
logging:
  level: "info"                # debug / info / warn / error
  format: "json"               # json / text
  enable_sensible_mask: true   # 敏感信息脱敏
```

---

## 九、实施计划

### 第一阶段：核心功能

1. 验证码模块（邮箱 + Redis 存储）
2. 用户注册（邮箱验证码注册）
3. 用户登录（邮箱 + 密码）
4. 修改密码
5. 忘记密码（邮箱找回）

### 第二阶段：短信能力

1. 短信适配器接口设计
2. 腾讯云短信实现
3. 阿里云短信实现
4. 手机号注册/登录/找回密码

### 第三阶段：钩子系统

1. 钩子类型定义
2. 钩子注册中心
3. 钩子执行器
4. SDK 钩子注册封装

### 第四阶段：第三方登录

1. OAuth 基础架构
2. 微信登录实现
3. 支付宝登录实现
4. QQ 登录实现

### 第五阶段：SDK 封装

1. SDK 客户端设计
2. 配置选项封装
3. 文档和示例

---

## 十、安全考虑

### 10.1 密码安全

| 措施 | 说明 |
|------|------|
| 加密算法 | bcrypt，cost=12 |
| 密码强度 | 最少6位，建议8位以上，包含字母数字 |
| 密码错误响应 | 统一返回"用户名或密码错误"，不区分具体错误 |
| 密码修改 | 修改密码后需重新登录 |

### 10.2 验证码安全

| 措施 | 说明 |
|------|------|
| 有效期 | 5 分钟，超时自动失效 |
| 发送限频 | 同一目标 60 秒内只能发送一次 |
| 使用次数 | 在有效期内只能使用一次，验证成功后立即失效 |
| 格式 | 6 位数字 |

### 10.3 登录防刷机制

**限制规则：**

| 限制类型 | 规则 | 说明 |
|----------|------|------|
| IP 限制 | 同一 IP 1 分钟内最多尝试 5 次 | 超过后锁定该 IP 10 分钟 |
| 账号限制 | 同一账号 10 分钟内最多失败 5 次 | 超过后锁定账号 30 分钟 |
| 设备指纹 | 同一设备 1 小时内最多失败 10 次 | 防止换 IP 绕过 |

**Redis Key 设计：**

```
# IP 限制
login:limit:ip:{ip_address}          # 计数，TTL 60秒
login:block:ip:{ip_address}          # 封禁标记，TTL 600秒

# 账号限制
login:fail:account:{account}         # 失败次数，TTL 600秒
login:block:account:{account}        # 账号封禁，TTL 1800秒

# 设备限制
login:fail:device:{device_id}        # 失败次数，TTL 3600秒
```

**错误响应：**

| 场景 | HTTP 状态码 | 错误码 |
|------|-------------|--------|
| 用户名或密码错误 | 400 | 10001 |
| 验证码错误 | 400 | 10002 |
| 验证码已过期 | 400 | 10003 |
| IP 请求过于频繁 | 429 | 20001 |
| 账号已锁定 | 403 | 20002 |

### 10.4 JWT 与黑名单机制

#### Token 设计

**Access Token：**
- 有效期：2 小时
- 存储：仅存储在客户端（Header 或 Cookie）
- 内容：用户 ID、签发时间、过期时间

**Refresh Token：**
- 有效期：7 天
- 存储：Redis + 数据库（支持撤销）
- 内容：用户 ID、设备标识、签发时间

#### 黑名单机制

**场景：**
1. 用户主动登出
2. 修改密码后踢出所有设备
3. 管理员禁用用户

**Redis 黑名单 Key 设计：**

```
# Token 黑名单
token:blacklist:{jti}               # JTI 为 JWT 的唯一标识，TTL = token 剩余有效期

# 用户所有 Token 黑名单（批量踢出）
user:blacklist:{user_id}            # 用户黑名单时间戳（密码修改时间）
```

**验证流程：**

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  请求Token  │────▶│ 解析JWT     │────▶│ 检查黑名单  │────▶│ 验证通过    │
└─────────────┘     └─────────────┘     └─────────────┘     └─────────────┘
                           │                    │
                           ▼                    ▼
                      解析失败            在黑名单中
                           │                    │
                           ▼                    ▼
                      返回401            返回401
```

**Token 刷新机制：**

| 步骤 | 说明 |
|------|------|
| 1 | 客户端使用 refresh_token 请求刷新 |
| 2 | 服务端验证 refresh_token 有效性 |
| 3 | 检查用户状态（是否被禁用、密码是否修改） |
| 4 | 生成新的 access_token |
| 5 | 可选：同时刷新 refresh_token |

### 10.5 错误码设计

#### 错误码规范

```
{模块码(2位)}{具体错误(3位)}
```

| 模块码 | 模块 |
|--------|------|
| 10 | 用户相关 |
| 20 | 验证码相关 |
| 30 | 认证相关 |
| 40 | 第三方登录 |
| 99 | 系统错误 |

#### 用户相关错误码 (10xxx)

| 错误码 | 说明 | HTTP 状态码 | 触发场景 |
|--------|------|-------------|----------|
| 10001 | 用户名或密码错误 | 400 | 登录时用户名或密码不匹配 |
| 10002 | 用户不存在 | 404 | 查询不存在的用户 ID |
| 10003 | 用户已被禁用 | 403 | 登录或操作被禁用的用户 |
| 10004 | 用户已被锁定 | 403 | 登录失败次数过多导致账号锁定 |
| 10005 | 用户名已存在 | 400 | 注册/设置用户名时，用户名已被使用 |
| 10006 | 邮箱已存在 | 400 | **1) 注册时邮箱已被其他账号绑定<br>2) 更换邮箱时，新邮箱已被其他账号绑定<br>3) 添加邮箱时，该邮箱已被其他账号绑定** |
| 10007 | 手机号已存在 | 400 | **1) 注册时手机号已被其他账号绑定<br>2) 更换手机号时，新手机号已被其他账号绑定<br>3) 添加手机号时，该手机号已被其他账号绑定** |
| 10008 | 密码强度不够 | 400 | 设置或修改密码时，密码不符合安全要求 |

#### 验证码相关错误码 (20xxx)

| 错误码 | 说明 | HTTP 状态码 |
|--------|------|-------------|
| 20001 | 验证码错误 | 400 |
| 20002 | 验证码已过期 | 400 |
| 20003 | 验证码已使用 | 400 |
| 20004 | 验证码发送过于频繁 | 429 |
| 20005 | 验证码不存在 | 404 |

#### 认证相关错误码 (30xxx)

| 错误码 | 说明 | HTTP 状态码 |
|--------|------|-------------|
| 30001 | 未登录 | 401 |
| 30002 | Token 已过期 | 401 |
| 30003 | Token 无效 | 401 |
| 30004 | Token 已被撤销 | 401 |
| 30501 | IP 请求过于频繁 | 429 |
| 30502 | 登录失败次数过多，账号已锁定 | 403 |

#### 第三方登录错误码 (40xxx)

| 错误码 | 说明 | HTTP 状态码 |
|--------|------|-------------|
| 40001 | 不支持的登录平台 | 400 |
| 40002 | 第三方授权失败 | 400 |
| 40003 | 获取用户信息失败 | 400 |
| 40004 | 账号已绑定其他用户 | 400 |

#### 系统错误码 (99xxx)

| 错误码 | 说明 | HTTP 状态码 |
|--------|------|-------------|
| 99001 | 系统内部错误 | 500 |
| 99002 | 数据库错误 | 500 |
| 99003 | Redis 错误 | 500 |
| 99004 | 邮件发送失败 | 500 |
| 99005 | 短信发送失败 | 500 |

#### 错误响应格式

```json
{
  "code": 10001,
  "message": "用户名或密码错误",
  "details": {},
  "request_id": "req-123456789",
  "timestamp": 1704067200
}
```

### 10.6 敏感数据保护

| 数据类型 | 保护措施 |
|----------|----------|
| 密码 | bcrypt 加密 |
| 手机号 | 数据库加密存储 + 日志脱敏 |
| 邮箱 | 日志脱敏 |
| OAuth Token | AES 加密存储 |
| API 密钥 | 环境变量，不写入代码 |

### 10.7 日志脱敏规则

| 原始数据 | 脱敏后 |
|----------|--------|
| `13800138000` | `138****8000` |
| `test@example.com` | `t***@example.com` |
| `password123` | `******` |

---

## 十一、总结

本设计方案基于现有 DDD 架构进行扩展，新增验证码、通知、OAuth、认证四个聚合，并引入钩子系统提供扩展能力。通过适配器模式支持多厂商短信和第三方登录，通过 SDK 封装提供便捷的集成方式。