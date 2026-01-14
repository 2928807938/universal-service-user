# Universal User Service - 使用文档

## 📋 目录

- [项目简介](#项目简介)
- [快速开始](#快速开始)
- [一、SDK 使用文档](#一、sdk-使用文档)
  - [安装 SDK](#安装-sdk)
  - [快速开始](#快速开始-1)
  - [配置选项](#配置选项)
  - [API 方法](#api-方法)
  - [钩子系统](#钩子系统)
  - [完整示例](#完整示例)
- [二、API 使用文档](#二、api-使用文档)
  - [部署 API 服务](#部署-api-服务)
  - [API 接口列表](#api-接口列表)
  - [认证接口](#认证接口)
  - [用户管理接口](#用户管理接口)
  - [验证码接口](#验证码接口)
  - [错误码说明](#错误码说明)
- [三、多租户配置中心](#三、多租户配置中心)
  - [概述](#概述)
  - [应用注册](#应用注册)
  - [配置文件上传](#配置文件上传)
  - [使用租户ID](#使用租户id)
  - [环境隔离](#环境隔离)
  - [数据库表结构](#数据库表结构)
  - [最佳实践](#最佳实践)
- [四、前端直接对接文档](#四、前端直接对接文档)
  - [快速启动](#快速启动)
  - [前端对接示例](#前端对接示例)
  - [Token 管理](#token-管理)
- [五、配置说明](#五、配置说明)
  - [配置文件](#配置文件)
  - [环境变量](#环境变量)
  - [数据库配置](#数据库配置)
  - [Redis 配置](#redis-配置)
  - [邮件/短信配置](#邮件短信配置)

---

## 项目简介

**Universal User Service** 是一个通用的用户管理系统，提供完整的用户认证、授权和账户管理能力。支持三种使用方式：

### 使用方式对比

| 方式 | 说明 | 适用场景 | 优势 |
|------|------|----------|------|
| **SDK 模式** | 引入 Go 包，直接调用函数 | Go 语言项目、需要深度集成 | 高性能、类型安全、可定制 |
| **API 模式** | 启动独立服务，通过 HTTP 调用 | 微服务架构、跨语言调用 | 解耦、跨语言、易于扩展 |
| **一站式模式** | 直接部署，前端直接对接 | 快速启动、独立用户中心 | 开箱即用、��单快速 |

### 核心特性

- ✅ **多租户配置中心** - 支持多应用、多环境的统一配置管理
- ✅ **自动建表** - 首次启动自动创建数据库表结构
- ✅ **多数据库支持** - PostgreSQL / MySQL / SQLite
- ✅ **多种登录方式** - 用户名、邮箱、手机号、第三方登录
- ✅ **验证码系统** - 邮箱验证码、短信验证码
- ✅ **JWT 认证** - Access Token + Refresh Token 双令牌机制
- ✅ **登录防刷** - IP 限频、账号限频、设备限频
- ✅ **钩子系统** - 支持业务扩展（注册后、登录前等）
- ✅ **规则引擎** - 强大的数据验证规则引擎

### 技术栈

- **Web 框架**: CloudWeGo Hertz
- **ORM**: GORM
- **数据库**: PostgreSQL / MySQL / SQLite
- **缓存**: Redis
- **认证**: JWT

---

## 快速开始

### 5 分钟上手 Universal User Service

根据你的项目需求，选择适合的使用方式：

---

#### 方式一：SDK 模式（推荐 Go 项目）

**适用场景**：Go 语言项目，直接在代码中集成用户管理功能

```bash
# 1. 引入 SDK
go get github.com/your-org/universal-service-user/sdk
```

```go
// 2. 在你的代码中使用（main.go）
package main

import (
    "context"
    "fmt"
    "github.com/your-org/universal-service-user/sdk"
    "github.com/redis/go-redis/v9"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func main() {
    // 连接数据库
    db, _ := gorm.Open(postgres.Open("host=localhost user=postgres password=postgres dbname=mydb"), &gorm.Config{})
    rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

    // 创建 SDK 客户端
    client, _ := sdk.New(
        sdk.WithDatabase(db),
        sdk.WithRedis(rdb),
        sdk.WithAutoMigrate(true),
    )
    defer client.Close()

    // 用户注册
    user, _ := client.Register(context.Background(), &sdk.RegisterRequest{
        Email:    "user@example.com",
        Password: "password123",
        Code:     "123456",
    })
    fmt.Printf("用户注册成功！ID: %d\n", user.ID)

    // 用户登录
    resp, _ := client.Login(context.Background(), &sdk.LoginRequest{
        LoginType: "email",
        Email:     "user@example.com",
        Password:  "password123",
    })
    fmt.Printf("登录成功！Token: %s\n", resp.AccessToken)
}
```

```bash
# 3. 运行你的项目
go run main.go
```

📖 [查看 SDK 完整文档](#一、sdk-使用文档)

---

#### 方式二：API 模式

**适用场景**：微服务架构、跨语言调用、前后端分离

```bash
# 1. 启动 API 服务
git clone https://github.com/your-org/universal-service-user.git
cd universal-service-user
cp config.example.yaml config.yaml  # 修改数据库配置
go run cmd/api/main.go
```

```bash
# 2. 注册应用,获取租户ID
curl -X POST http://localhost:8080/api/v1/apps/register \
  -H "Content-Type: application/json" \
  -d '{"app_name":"我的应用","description":"应用描述"}'

# 响应: {"code":0,"message":"success","data":{"tenant_id":"550e8400-..."}}
```

```bash
# 3. 使用租户ID调用其他API
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -H "X-Tenant-Id: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-App-Environment: prod" \
  -d '{"email":"user@example.com","password":"password123","code":"123456"}'

# 响应
# {"code":0,"message":"success","data":{"id":1,"username":"user1",...}}
```

📖 [查看 API 完整文档](#二、api-使用文档)

---

#### 方式三：一站式模式

**适用场景**：快速启动、独立用户中心、前端项目

```bash
# 1. 启动服务
git clone https://github.com/your-org/universal-service-user.git
cd universal-service-user
cp config.example.yaml config.yaml  # 修改数据库配置
go run cmd/api/main.go
```

```bash
# 2. 注册应用
curl -X POST http://localhost:8080/api/v1/apps/register \
  -H "Content-Type: application/json" \
  -d '{"app_name":"前端应用"}'
# 记录返回的 tenant_id
```

```javascript
// 3. 前端直接调用 (记得添加 X-Tenant-Id 请求头)
const response = await fetch('http://localhost:8080/api/v1/auth/login', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'X-Tenant-Id': 'your-tenant-id',  // 替换为你的租户ID
    'X-App-Environment': 'prod'
  },
  body: JSON.stringify({
    login_type: 'email',
    email: 'user@example.com',
    password: 'password123'
  })
});

const { access_token, user_info } = (await response.json()).data;
localStorage.setItem('token', access_token);
```

📖 [查看前端对接完整文档](#四、前端直接对接文档)

---

### 需要准备什么？

- **数据库**: PostgreSQL / MySQL / SQLite（选一个）
- **Redis**: 用于缓存和会话管理
- **Go 1.18+**: 开发环境

---

### 下一步

- 📖 阅读详细文档了解完整功能
- ⚙️ 查看 [配置说明](#四、配置说明) 自定义配置

---

## 一、SDK 使用文档

### 安装 SDK

#### 引入 SDK

在你的 Go 项目中引入 SDK：

```bash
# 1. 初始化 Go Module（如果还没有）
go mod init your-project

# 2. 引入 SDK
go get github.com/your-org/universal-service-user/sdk
```

**依赖说明**：

SDK 会自动安装以下依赖：

- **GORM**: `gorm.io/gorm` - ORM 框架
- **数据库驱动**（根据使用情况选择）:
  - PostgreSQL: `gorm.io/driver/postgres`
  - MySQL: `gorm.io/driver/mysql`
  - SQLite: `gorm.io/driver/sqlite`
- **Redis**: `github.com/redis/go-redis/v9`
- **JWT**: `github.com/golang-jwt/jwt/v5`
- **验证码**: `github.com/mojocn/base64Captcha`

#### 本地引用

如果您想使用本地开发版本：

```bash
# 1. 将项目克隆到本地
git clone https://github.com/your-org/universal-service-user.git
cd universal-service-user

# 2. 在您的项目中引用本地路径
cd your-project
go mod edit -replace github.com/your-org/universal-service-user=/path/to/universal-service-user
go mod tidy
```

### 快速开始

#### 最简示例

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/redis/go-redis/v9"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "github.com/your-org/universal-service-user/sdk"
)

func main() {
    ctx := context.Background()

    // 1. 初始化数据库
    dsn := "host=localhost user=postgres password=postgres dbname=mydb port=5432"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    // 2. 初始化 Redis
    rdb := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })

    // 3. 创建 SDK 客户端
    client, err := sdk.New(
        sdk.WithDatabase(db),
        sdk.WithRedis(rdb),
        sdk.WithAutoMigrate(true),  // 自动建表
    )
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // 4. 使用 SDK
    user, err := client.Register(ctx, &sdk.RegisterRequest{
        Email:    "user@example.com",
        Password: "password123",
        Code:     "123456",
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("用户注册成功: %+v\n", user)
}
```
{
  "login_type": "email",        // 登录方式: username / email / phone
  "email": "user@example.com",  // 邮箱（login_type=email 时必填）
  "username": "john",           // 用户名（login_type=username 时必填）
  "phone": "13800138000",       // 手机号（login_type=phone 时必填）
  "password": "password123",    // 密码（username/email 登录时必填）
  "code": "123456"              // 验证码（phone 登录时必填）
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 7200,
    "user_info": {
      "id": 1,
      "username": "john",
      "email": "user@example.com",
      "status": 1,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  }
}
```

**登录方式示例**:

```bash
# 1. 用户名 + 密码登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "login_type": "username",
    "username": "john",
    "password": "password123"
  }'

# 2. 邮箱 + 密码登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "login_type": "email",
    "email": "user@example.com",
    "password": "password123"
  }'

# 3. 手机号 + 验证码登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "login_type": "phone",
    "phone": "13800138000",
    "code": "123456"
  }'
```

#### 2. 用户登出

**接口**: `POST /api/v1/auth/logout`

**请求头**:
```
Authorization: Bearer {access_token}
```

**请求参数**:

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."  // 可选
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": null
}
```

**示例**:

```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "refresh_token": "YOUR_REFRESH_TOKEN"
  }'
```

#### 3. 刷新令牌

**接口**: `POST /api/v1/auth/refresh`

**请求参数**:

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 7200
  }
}
```

**示例**:

```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "YOUR_REFRESH_TOKEN"
  }'
```

### 用户管理接口

#### 1. 用户注册

**接口**: `POST /api/v1/users/register`

**请求参数**:

```json
{
  "email": "user@example.com",   // 邮箱（与手机号二选一）
  "phone": "13800138000",        // 手机号（与邮箱二选一）
  "password": "password123",     // 密码（6-32位）
  "code": "123456",              // 验证码
  "username": "john"             // 用户名（可选，默认生成）
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "username": "john",
    "email": "user@example.com",
    "status": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

**注册流程**:

```bash
# 1. 先发送验证码
curl -X POST http://localhost:8080/api/v1/verification/code/send \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "target": "user@example.com",
    "scene": "register"
  }'

# 2. 使用验证码注册
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "code": "123456"
  }'
```

#### 2. 获取用户信息

**接口**: `GET /api/v1/users/:id`

**请求头**:
```
Authorization: Bearer {access_token}
```

**路径参数**:
- `id`: 用户 ID

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "username": "john",
    "email": "user@example.com",
    "status": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

**示例**:

```bash
curl -X GET http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

#### 3. 更新用户信息

**接口**: `PUT /api/v1/users/:id`

**请求头**:
```
Authorization: Bearer {access_token}
```

**路径参数**:
- `id`: 用户 ID

**请求参数**:

```json
{
  "username": "new_username",   // 可选
  "nickname": "John Doe",       // 可选
  "avatar": "https://...",      // 可选
  "status": 1                   // 可选（0:未激活 1:正常 2:禁用 3:锁定）
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "username": "new_username",
    "email": "user@example.com",
    "status": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

**示例**:

```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "username": "new_username"
  }'
```

#### 4. 修改密码

**接口**: `POST /api/v1/users/password/change`

**请求头**:
```
Authorization: Bearer {access_token}
```

**请求参数**:

```json
{
  "old_password": "old_password123",
  "new_password": "new_password123"
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": null
}
```

#### 5. 重置密码（忘记密码）

**接口**: `POST /api/v1/users/password/reset`

**请求参数**:

```json
{
  "email": "user@example.com",   // 邮箱（与手机号二选一）
  "phone": "13800138000",        // 手机号（与邮箱二选一）
  "code": "123456",              // 验证码
  "new_password": "new_password123"
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": null
}
```

**重置密码流程**:

```bash
# 1. 先发送验证码
curl -X POST http://localhost:8080/api/v1/verification/code/send \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "target": "user@example.com",
    "scene": "reset_password"
  }'

# 2. 使用验证码重置密码
curl -X POST http://localhost:8080/api/v1/users/password/reset \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "code": "123456",
    "new_password": "new_password123"
  }'
```

### 验证码接口

#### 1. 发送验证码

**接口**: `POST /api/v1/verification/code/send`

**请求参数**:

```json
{
  "type": "email",              // 类型: email / phone
  "target": "user@example.com", // 目标：邮箱地址或手机号
  "scene": "register"           // 场景: register / login / reset_password
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "验证码已发送"
  }
}
```

**场景说明**:

| 场景 | 说明 |
|------|------|
| `register` | 注册验证码 |
| `login` | 登录验证码 |
| `reset_password` | 重置密码验证码 |

**示例**:

```bash
# 发送邮箱验证码
curl -X POST http://localhost:8080/api/v1/verification/code/send \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "target": "user@example.com",
    "scene": "register"
  }'

# 发送短信验证码
curl -X POST http://localhost:8080/api/v1/verification/code/send \
  -H "Content-Type: application/json" \
  -d '{
    "type": "phone",
    "target": "13800138000",
    "scene": "login"
  }'
```

#### 2. 验证验证码

**接口**: `POST /api/v1/verification/code/verify`

**请求参数**:

```json
{
  "type": "email",              // 类型: email / phone
  "target": "user@example.com", // 目标：邮箱地址或手机号
  "scene": "register",          // 场景
  "code": "123456"              // 验证码
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "valid": true,
    "message": "验证码验证通过"
  }
}
```

**示例**:

```bash
curl -X POST http://localhost:8080/api/v1/verification/code/verify \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "target": "user@example.com",
    "scene": "register",
    "code": "123456"
  }'
```

### 错误码说明

#### 通用错误码

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 10001 | 参数错误 |
| 10002 | 数据库错误 |
| 10003 | Redis 错误 |
| 99999 | 系统内部错误 |

#### 用户模块错误码

| 错误码 | 说明 |
|--------|------|
| 10004 | 用户不存在 |
| 10005 | 用户已存在 |
| 10006 | 邮箱已被其他账号绑定 |
| 10007 | 手机号已被其他账号绑定 |
| 10008 | 密码错误 |
| 10009 | 旧密码错误 |
| 10010 | 用户状态异常（未激活/禁用/锁定） |

#### 验证码错误码

| 错误码 | 说明 |
|--------|------|
| 20001 | 验证码错误 |
| 20002 | 验证码已过期 |
| 20003 | 验证码发送频率过高 |
| 20004 | 验证码不存在 |

#### 认证错误码

| 错误码 | 说明 |
|--------|------|
| 30001 | Token 无效 |
| 30002 | Token 已过期 |
| 30003 | Refresh Token 无效 |
| 30004 | 登录失败次数过多，账号已锁定 |
| 30005 | IP 访问频率过高 |

**错误响应示例**:

```json
{
  "code": 10006,
  "message": "邮箱已被其他账号绑定",
  "data": null
}
```

---

## 三、多租户配置中心

### 概述

Universal User Service 现在支持多租户架构,允许一个服务实例同时为多个应用提供用户管理服务。每个应用拥有独立的配置空间和数据隔离。

### 应用注册

在使用用户服务之前,需要先注册应用以获取租户ID(TenantID):

```bash
# 注册应用
curl -X POST http://localhost:8080/api/v1/apps/register \
  -H "Content-Type: application/json" \
  -d '{
    "app_name": "我的应用",
    "description": "应用描述",
    "email": "admin@example.com"
  }'

# 响应
{
  "code": 0,
  "message": "success",
  "data": {
    "tenant_id": "550e8400-e29b-41d4-a716-446655440000"
  }
}
```

### 配置文件上传

应用注册成功后,需要使用 CLI 工具上传配置文件到数据库。配置文件包含了用户服务的各种配置项(如验证码设置、邮件配置等)。

#### 完整工作流程

```
┌────────────────────────────────────────────────────────────────┐
│                   第一步：注册应用                                │
└────────────────────────────────────────────────────────────────┘

curl -X POST http://localhost:8080/api/v1/apps/register \
  -d '{"app_name":"我的应用"}'

↓ 返回 tenant_id: "550e8400-..."

┌────────────────────────────────────────────────────────────────┐
│                   第二步：准备配置文件                            │
└────────────────────────────────────────────────────────────────┘

创建 config.yaml，包含：
  - app.tenant_id (从第一步获取)
  - app.environment (prod/dev/test)
  - email 配置
  - sms 配置
  - jwt 配置
  ...其他配置

┌────────────────────────────────────────────────────────────────┐
│                   第三步：上传配置到数据库                        │
└────────────────────────────────────────────────────────────────┘

./universal-service-cli -c config.yaml -u

↓ 连接到 Universal User Service 数据库

  ├─ 验证 tenant_id 是否存在于 apps 表
  ├─ 将 config.yaml 转换为 JSON
  ├─ 存储到 app_configs 表
  └─ 记录到 config_histories 表

┌────────────────────────────────────────────────────────────────┐
│                   第四步：启动 API 服务                          │
└────────────────────────────────────────────────────────────────┘

go run cmd/api/main.go

↓ 服务启动时：

  └─ 从数据库读取租户配置
      └─ 根据请求头的 X-Tenant-Id 和 X-App-Environment
          └─ 从 app_configs 表加载对应的配置

┌────────────────────────────────────────────────────────────────┐
│                   第五步：使用服务                               │
└────────────────────────────────────────────────────────────────┘

curl -X POST http://localhost:8080/api/v1/users/register \
  -H "X-Tenant-Id: 550e8400-..." \
  -H "X-App-Environment: prod" \
  ...

↓ 服务根据租户配置处理请求
```

#### 1. 编译 CLI 工具

```bash
# 进入项目目录
cd universal-service-user

# 编译 CLI 工具
go build -o universal-service-cli cmd/universal-service/main.go

# 或直接运行
go run cmd/universal-service/main.go -c config.yaml -u
```

#### 2. 准备配置文件

创建 YAML 格式的配置文件 `config.yaml`:

> **重要**: 本项目采用**完全多租户架构**，每个租户拥有独立的用户数据库。

```yaml
# ============================================================
# 租户应用配置文件
# ============================================================
# 说明：
# 1. app.tenant_id 和 app.environment 是必需的
# 2. user_database 是必需的（配置租户的用户数据库）
# 3. 其他配置根据需要添加
# 4. 您可以添加任何自定义配置，服务会忽略不需要的部分
# ============================================================

# 应用信息（必需）
app:
  tenant_id: "550e8400-e29b-41d4-a716-446655440000"  # 替换为你的 TenantID
  environment: "prod"                                # 环境: prod / dev / test / staging

# ============================================================
# 用户数据库配置（必需）
# ============================================================
# 说明：每个租户需要提供自己的数据库，用于存储用户数据
#      如果 auto_create_tables=true，服务会自动创建表
user_database:
  driver: "postgres"              # 数据库类型: postgres / mysql / sqlite
  host: "your-user-db-host.com"   # 数据库主机
  port: 5432                      # 数据库端口
  user: "your-db-user"            # 数据库用户
  password: "your-db-password"    # 数据库密码
  dbname: "your_app_users"        # 数据库名称
  sslmode: "require"              # PostgreSQL SSL模式
  auto_create_tables: true        # 自动创建表（推荐）
  max_open_conns: 100             # 最大连接数
  max_idle_conns: 10              # 最大空闲连接数
  conn_max_lifetime: 3600         # 连接最大生命周期（秒）

# 邮件配置（可选，如需发送验证码邮件）
email:
  enabled: true
  smtp:
    host: "smtp.example.com"
    port: 587
    username: "noreply@example.com"
    password: "your-password"
    from: "系统通知 <noreply@example.com>"
  templates:
    register:
      subject: "注册验证码"
      body: "您的注册验证码是：{code}，有效期5分钟"
    login:
      subject: "登录验证码"
      body: "您的登录验证码是：{code}，有效期5分钟"
    reset_password:
      subject: "重置密码验证码"
      body: "您的重置密码验证码是：{code}，有效期10分钟"

# 短信配置(可选)
sms:
  enabled: false
  provider: "tencent"         # tencent / aliyun
  tencent:
    secret_id: "your-secret-id"
    secret_key: "your-secret-key"
    app_id: "your-app-id"
    sign: "您的签名"
    templates:
      register: "123456"
      login: "123456"
      reset_password: "123456"

# JWT 配置（可选）
jwt:
  secret: "your-jwt-secret-key"
  access_token_expire: 7200       # 2小时(秒)
  refresh_token_expire: 604800    # 7天(秒)
  issuer: "user-service"

# 验证码配置（可选）
verification:
  code_length: 6                  # 验证码长度
  expire: 300                     # 过期时间(秒)
  rate_limit: 60                  # 发送间隔(秒)

# 功能开关（可选）
features:
  registration: true              # 是否开放注册
  email_login: true               # 邮箱登录
  phone_login: false              # 手机号登录
  username_login: true            # 用户名登录
  oauth_login: false              # 第三方登录

# 登录防刷配置（可选）
login_rate_limit:
  enabled: true
  account:
    max_failures: 5               # 最大失败次数
    lock_duration: 1800           # 锁定时长(秒)

# 您可以添加任何自定义配置
# 服务会忽略它不认识的配置项
custom_settings:
  theme_color: "#FF5722"
  max_users: 10000
```

#### 3. 上传配置文件

```bash
# 方式一: 使用编译后的工具
./universal-service-cli \
  -c config.yaml \
  -u

# 方式二: 直接运行
go run cmd/universal-service/main.go \
  -c config.yaml \
  -u
```

#### 4. 数据库连接配置

> **重要说明**: CLI 工具需要连接到 **Universal User Service 的中心数据库**。

本项目采用**完全多租户架构**，包含两种数据库：

### 1. 中心数据库（配置数据库）

CLI 工具和 API 服务启动时都连接到这个数据库，包含：
- `apps` - 存储应用信息（用于验证租户是否存在）
- `app_configs` - 存储租户配置（包含用户数据库连接信息、邮件配置等）
- `config_histories` - 存储配置变更历史

### 2. 租户用户数据库（用户数据）

每个租户拥有**独立的用户数据库**，包含：
- `users` - 用户主表
- `user_profiles` - 用户扩展信息表
- `user_login_logs` - 登录日志表
- `user_sessions` - 会话表
- `user_oauth_bindings` - OAuth绑定表

**工作原理**：
1. 租户配置文件中包含 `user_database` 配置（见上面的配置文件示例）
2. 服务启动时连接到中心数据库
3. 收到请求时，根据 `X-Tenant-Id` 从中心数据库读取租户配置
4. 根据配置中的 `user_database` 连接到租户的用户数据库
5. 如果 `auto_create_tables=true`，自动在用户数据库中创建表
6. 使用用户数据库处理请求

```
┌────────────────────────────────────────────────────────────┐
│                   完全多租户架构                             │
└────────────────────────────────────────────────────────────┘

中心数据库 (config.yaml 中的 database 配置)
├─ apps (租户应用)
├─ app_configs (租户配置，包含 user_database 连接信息)
└─ config_histories (配置历史)

租户A的用户数据库 (config.yaml 中 user_database 配置)
├─ users
├─ user_profiles
├─ user_login_logs
├─ user_sessions
└─ user_oauth_bindings

租户B的用户数据库 (各自的 user_database 配置)
├─ users
├─ user_profiles
└─ ...
```

可以通过以下方式配置**中心数据库**连接:

**方式一: 使用命令行参数**

```bash
./universal-service-cli \
  -c config.yaml \
  -u \
  --db-driver postgres \
  --db-host localhost \
  --db-port 5432 \
  --db-user postgres \
  --db-password your-password \
  --db-name universal_service_user
```

**方式二: 使用环境变量**

```bash
export CONFIG_DB_DRIVER=postgres
export CONFIG_DB_HOST=localhost
export CONFIG_DB_PORT=5432
export CONFIG_DB_USER=postgres
export CONFIG_DB_PASSWORD=your-password
export CONFIG_DB_NAME=universal_service_user

./universal-service-cli -c config.yaml -u
```

**方式三: 使用 DSN**

```bash
export CONFIG_DB_DSN="host=localhost user=postgres password=your-password dbname=universal_service_user port=5432 sslmode=disable"

./universal-service-cli -c config.yaml -u
```

#### 5. 参数说明

| 参数 | 短参数 | 说明 | 默认值 | 必填 |
|------|--------|------|--------|------|
| `--config` | `-c` | 配置文件路径 | - | ✅ |
| `--upload-only` | `-u` | 仅上传模式 | - | ✅ |
| `--db-driver` | - | 数据库类型 | postgres | ❌ |
| `--db-host` | - | 数据库主机 | localhost | ❌ |
| `--db-port` | - | 数据库端口 | 5432 | ❌ |
| `--db-user` | - | 数据库用户 | postgres | ❌ |
| `--db-password` | - | 数据库密码 | - | ❌ |
| `--db-name` | - | 数据库名称 | - | ❌ |
| `--db-dsn` | - | 数据库 DSN | - | ❌ |

#### 6. 上传成功示例

```bash
$ ./universal-service-cli -c config.yaml -u
upload completed: tenant_id=550e8400-e29b-41d4-a716-446655440000 environment=prod
```

#### 7. 配置管理说明

- **首次上传**: 创建新配置记录,版本号为 1
- **更新配置**: 自动递增版本号,记录历史
- **配置验证**: 上传前自动验证租户是否存在
- **配置历史**: 每次更新都会保存历史记录,支持回滚
- **JSON 转换**: YAML 配置会自动转换为 JSON 存储到数据库

#### 8. 用户数据库初始化

在上传配置文件之前，您需要准备租户的用户数据库。有两种方式：

##### 方式一：自动建表（推荐）

在配置文件中设置 `auto_create_tables: true`：

```yaml
user_database:
  driver: "postgres"
  host: "your-db-host.com"
  port: 5432
  user: "your-db-user"
  password: "your-db-password"
  dbname: "your_app_users"
  auto_create_tables: true  # 服务会自动创建表
```

**优点**：简单方便，无需手动执行SQL
**适用场景**：首次部署，数据库为空

##### 方式二：手动执行SQL脚本

如果您想手动控制建表过程，或查看表结构，可以使用我们提供的 SQL 脚本。

**步骤1：创建用户数据库**

```bash
# PostgreSQL
createdb -U postgres your_app_users

# MySQL
mysql -u root -p -e "CREATE DATABASE your_app_users CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# SQLite
touch your_app_users.db
```

**步骤2：执行建表脚本**

根据您选择的数据库类型，执行对应的SQL脚本：

```bash
# PostgreSQL
psql -U postgres -d your_app_users -f scripts/schema/postgresql.sql

# MySQL
mysql -u root -p your_app_users < scripts/schema/mysql.sql

# SQLite
sqlite3 your_app_users.db < scripts/schema/sqlite.sql
```

**步骤3：验证表创建**

```bash
# PostgreSQL
psql -U postgres -d your_app_users -c "\dt"

# MySQL
mysql -u root -p your_app_users -e "SHOW TABLES;"

# SQLite
sqlite3 your_app_users.db ".tables"
```

应该看到以下5张表：
- `users`
- `user_profiles`
- `user_login_logs`
- `user_sessions`
- `user_oauth_bindings`

**SQL 脚本位置**：
- PostgreSQL: `scripts/schema/postgresql.sql`
- MySQL: `scripts/schema/mysql.sql`
- SQLite: `scripts/schema/sqlite.sql`

详细说明请参考：[租户用户数据库建表脚本](../scripts/schema/README.md)

#### 9. 多环境配置

为不同环境创建不同的配置文件:

```bash
# 生产环境
./universal-service-cli -c config-prod.yaml -u

# 开发环境
./universal-service-cli -c config-dev.yaml -u

# 测试环境
./universal-service-cli -c config-test.yaml -u
```

配置文件中使用相同的 `tenant_id`,但不同的 `environment`:

```yaml
# config-prod.yaml
app:
  tenant_id: "550e8400-e29b-41d4-a716-446655440000"
  environment: "prod"

# config-dev.yaml
app:
  tenant_id: "550e8400-e29b-41d4-a716-446655440000"
  environment: "dev"
```

### 使用租户ID

获取 TenantID 后,在所有 API 请求中添加以下请求头:

```bash
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -H "X-Tenant-Id: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-App-Environment: prod" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "code": "123456"
  }'
```

### 环境隔离

通过 `X-App-Environment` 请求头实现环境隔离:

- `prod` - 生产环境(默认)
- `dev` - 开发环境
- `test` - 测试环境
- `staging` - 预发布环境

```bash
# 开发环境
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "X-Tenant-Id: your-tenant-id" \
  -H "X-App-Environment: dev" \
  ...

# 测试环境
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "X-Tenant-Id: your-tenant-id" \
  -H "X-App-Environment: test" \
  ...
```

### 数据库表结构

多租户配置中心使用以下数据库表:

#### apps 表

存储应用基本信息:

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int | 主键 |
| tenant_id | varchar(64) | 租户ID(唯一) |
| app_name | varchar(128) | 应用名称 |
| description | varchar(512) | 应用描述 |
| status | varchar(32) | 状态 |
| created_at | timestamp | 创建时间 |
| updated_at | timestamp | 更新时间 |

#### app_configs 表

存储应用配置数据(JSON格式):

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int | 主键 |
| tenant_id | varchar(64) | 租户ID |
| environment | varchar(32) | 环境 |
| config_data | json | 配置数据 |
| version | int | 版本号 |
| created_at | timestamp | 创建时间 |
| updated_at | timestamp | 更新时间 |

#### config_histories 表

存储配置变更历史:

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int | 主键 |
| tenant_id | varchar(64) | 租户ID |
| environment | varchar(32) | 环境 |
| config_data | json | 配置数据 |
| version | int | 版本号 |
| change_reason | varchar(255) | 变更原因 |
| created_at | timestamp | 创建时间 |

### 最佳实践

1. **应用隔离**: 不同应用使用不同的 TenantID
2. **环境管理**: 同一应用的不同环境使用相同的 TenantID,通过 Environment 区分
3. **配置安全**: 妥善保管 TenantID,避免泄露
4. **配置版本**: 配置变更会自动记录历史,支持回滚

---

## 二、SDK 使用文档

### 安装 SDK

#### 方式一：Go Module（推荐）

在您的 Go 项目中引入 SDK：

```bash
# 1. 初始化 Go Module（如果还没有）
go mod init your-project

# 2. 引入 SDK
go get github.com/your-org/universal-service-user/sdk
```

**依赖说明**：

SDK 会自动安装以下依赖：

- **GORM**: `gorm.io/gorm` - ORM 框架
- **数据库驱动**（根据使用情况选择）:
  - PostgreSQL: `gorm.io/driver/postgres`
  - MySQL: `gorm.io/driver/mysql`
  - SQLite: `gorm.io/driver/sqlite`
- **Redis**: `github.com/redis/go-redis/v9`
- **JWT**: `github.com/golang-jwt/jwt/v5`
- **验证码**: `github.com/mojocn/base64Captcha`

#### 方式二：本地引用

如果您想使用本地开发版本：

```bash
# 1. 将项目克隆到本地
git clone https://github.com/your-org/universal-service-user.git
cd universal-service-user

# 2. 在您的项目中引用本地路径
cd your-project
go mod edit -replace github.com/your-org/universal-service-user=/path/to/universal-service-user
go mod tidy
```

### 快速开始

#### 1. 最简示例

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/redis/go-redis/v9"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "github.com/your-org/universal-service-user/sdk"
)

func main() {
    ctx := context.Background()

    // 1. 初始化数据库
    dsn := "host=localhost user=postgres password=postgres dbname=mydb port=5432"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    // 2. 初始化 Redis
    rdb := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })

    // 3. 创建 SDK 客户端
    client, err := sdk.New(
        sdk.WithDatabase(db),
        sdk.WithRedis(rdb),
        sdk.WithAutoMigrate(true),  // 自动建表
    )
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // 4. 使用 SDK
    user, err := client.Register(ctx, &sdk.RegisterRequest{
        Email:    "user@example.com",
        Password: "password123",
        Code:     "123456",
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("用户注册成功: %+v\n", user)
}
```

### 配置选项

#### 基础配置

| 选项 | 类型 | 说明 | 必填 |
|------|------|------|------|
| `WithDatabase` | `*gorm.DB` | 数据库连接 | ✅ |
| `WithRedis` | `*redis.Client` | Redis 连接 | ✅ |
| `WithAutoMigrate` | `bool` | 是否自动建表 | ❌ |
| `WithConfig` | `*config.Config` | 完整配置对象 | ❌ |

#### JWT 配置

```go
sdk.WithJWT(
    "your-jwt-secret",    // 密钥
    2*time.Hour,          // Access Token 过期时间
    7*24*time.Hour,       // Refresh Token 过期时间
    "user-service",       // 签发者
)
```

#### 验证码配置

```go
sdk.WithVerification(
    6,              // 验证码长度
    5*time.Minute,  // 过期时间
    60*time.Second, // 发送间隔
)
```

#### 邮件配置

```go
sdk.WithEmailProvider(&email.SMTPConfig{
    Host:     "smtp.example.com",
    Port:     587,
    Username: "noreply@example.com",
    Password: "your-password",
    From:     "系统通知 <noreply@example.com>",
    Templates: map[string]string{
        "register": "验证码：{code}",
        "login": "验证码：{code}",
        "reset_password": "验证码：{code}",
    },
})
```

### API 方法

#### 1. 用户注册

```go
user, err := client.Register(ctx, &sdk.RegisterRequest{
    Email:    "user@example.com",
    Password: "password123",
    Code:     "123456",
})
```

#### 2. 用户登录

```go
// 邮箱 + 密码登录
resp, err := client.Login(ctx, &sdk.LoginRequest{
    LoginType: "email",
    Email:     "user@example.com",
    Password:  "password123",
})

// 手机号 + 验证码登录
resp, err := client.Login(ctx, &sdk.LoginRequest{
    LoginType: "phone",
    Phone:     "13800138000",
    Code:      "123456",
})
```

#### 3. 其他方法

```go
// 登出
client.Logout(ctx, accessToken, refreshToken)

// 刷新令牌
accessToken, refreshToken, expiresIn, err := client.RefreshToken(ctx, oldRefreshToken)

// 获取用户信息
user, err := client.GetUser(ctx, userID)

// 更新用户信息
username := "new_username"
user, err := client.UpdateUser(ctx, userID, &sdk.UpdateUserRequest{
    Username: &username,
})

// 重置密码
client.ResetPassword(ctx, &sdk.ResetPasswordRequest{
    Email:       "user@example.com",
    Code:        "123456",
    NewPassword: "newpassword123",
})

// 发送验证码
client.SendVerificationCode(ctx, &sdk.SendVerificationCodeRequest{
    Target: "user@example.com",
    Scene:  "register",
})

// 验证验证码
valid, err := client.VerifyCode(ctx, &sdk.VerifyCodeRequest{
    Target: "user@example.com",
    Scene:  "register",
    Code:   "123456",
})
```

### 钩子系统

SDK 提供了强大的钩子系统，允许您在关键操作前后执行自定义逻辑。

#### 可用的钩子

| 钩子 | 说明 | 执行时机 |
|------|------|----------|
| `BeforeCreate` | 创建用户前 | 用户注册前 |
| `AfterCreate` | 创建用户后 | 用户注册后 |
| `BeforeLogin` | 登录前 | 用户登录前 |
| `AfterLogin` | 登录后 | 用户登录后 |
| `BeforeLogout` | 登出前 | 用户登出前 |
| `AfterLogout` | 登出后 | 用户登出后 |
| `BeforeUpdate` | 更新用户前 | 更新用户信息前 |
| `AfterUpdate` | 更新用户后 | 更新用户信息后 |
| `BeforeResetPassword` | 重置密码前 | 重置密码前 |
| `AfterResetPassword` | 重置密码后 | 重置密码后 |

#### 注册钩子示例

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/your-org/universal-service-user/hook"
    "github.com/your-org/universal-service-user/sdk"
)

// 定义钩子处理函数
func afterCreateUserHook(ctx context.Context, metadata hook.Metadata) error {
    // 获取用户信息
    userID := metadata["user_id"]
    user := metadata["user"]

    fmt.Printf("用户创建成功！ID: %v, 用户信息: %+v\n", userID, user)

    // 发送欢迎邮件
    // go sendWelcomeEmail(user)

    // 添加到 CRM 系统
    // go addToCRM(user)

    return nil
}

func beforeLoginHook(ctx context.Context, metadata hook.Metadata) error {
    fmt.Printf("用户准备登录...\n")

    // 记录登录尝试
    // logLoginAttempt(ctx)

    // 安全检查
    // if isSuspiciousLogin(ctx) {
    //     return fmt.Errorf("可疑登录，��拦截")
    // }

    return nil
}

func afterLoginHook(ctx context.Context, metadata hook.Metadata) error {
    userID := metadata["user_id"]

    // 记录登录日志
    fmt.Printf("用户 %d 登录成功\n", userID)

    // 更新最后登录时间
    // updateLastLoginTime(userID)

    return nil
}

func main() {
    // ... 初始化 SDK

    // 注册钩子
    client.RegisterHook(hook.AfterCreate, afterCreateUserHook)
    client.RegisterHook(hook.BeforeLogin, beforeLoginHook)
    client.RegisterHook(hook.AfterLogin, afterLoginHook)

    // 正常使用 SDK，钩子会自动触发
    user, err := client.Register(ctx, &sdk.RegisterRequest{
        Email:    "user@example.com",
        Password: "password123",
        Code:     "123456",
    })
    // afterCreateUserHook 会自动执行

    resp, err := client.Login(ctx, &sdk.LoginRequest{
        LoginType: "email",
        Email:     "user@example.com",
        Password:  "password123",
    })
    // beforeLoginHook 和 afterLoginHook 会自动执行
}
```

#### 钩子最佳实践

**1. 异步处理耗时操作**

```go
func afterCreateUserHook(ctx context.Context, metadata hook.Metadata) error {
    user := metadata["user"]

    // 使用 goroutine 异步处理，避免阻塞主流程
    go func() {
        // 发送欢迎邮件
        sendWelcomeEmail(user)

        // 添加到推荐系统
        addToRecommendationSystem(user)
    }()

    return nil
}
```

**2. 错误处理**

```go
func afterCreateUserHook(ctx context.Context, metadata hook.Metadata) error {
    userID := metadata["user_id"].(int)

    // 如果是关键操作，返回错误会回滚整个流程
    if err := syncToExternalSystem(userID); err != nil {
        // 记录错误日志
        log.Printf("同步到外部系统失败: %v", err)

        // 根据业务需求决定是否返回错误
        // return err  // 返回错误会导致注册失败
        return nil  // 不返回错误，只记录日志
    }

    return nil
}
```

**3. 钩子链**

```go
// 可以注册多个同名钩子，它们会按注册顺序依次执行
client.RegisterHook(hook.AfterCreate, sendWelcomeEmail)
client.RegisterHook(hook.AfterCreate, addToCRM)
client.RegisterHook(hook.AfterCreate, initUserProfile)

// 执行顺序：
// 1. sendWelcomeEmail
// 2. addToCRM
// 3. initUserProfile
```

### 完整示例

#### 示例 1：Web 应用用户系统

```go
package main

import (
    "context"
    "log"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/redis/go-redis/v9"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "github.com/your-org/universal-service-user/hook"
    "github.com/your-org/universal-service-user/sdk"
)

var client *sdk.Client

func init() {
    // 初始化数据库
    db, err := gorm.Open(postgres.Open("host=localhost user=postgres password=postgres dbname=mydb"), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    // 初始化 Redis
    rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

    // 初始化 SDK
    client, err = sdk.New(
        sdk.WithDatabase(db),
        sdk.WithRedis(rdb),
        sdk.WithAutoMigrate(true),
        sdk.WithEmailProvider(emailConfig),
    )
    if err != nil {
        log.Fatal(err)
    }

    // 注册钩子
    client.RegisterHook(hook.AfterCreate, func(ctx context.Context, metadata hook.Metadata) error {
        log.Printf("新用户注册: %+v", metadata["user"])
        return nil
    })
}

// 注册接口
func registerHandler(c *gin.Context) {
    var req struct {
        Email    string `json:"email"`
        Password string `json:"password"`
        Code     string `json:"code"`
    }
    if err := c.BindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "参数错误"})
        return
    }

    user, err := client.Register(c.Request.Context(), &sdk.RegisterRequest{
        Email:    req.Email,
        Password: req.Password,
        Code:     req.Code,
    })

    if err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{"data": user})
}

// 登录接口
func loginHandler(c *gin.Context) {
    var req struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    if err := c.BindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "参数错误"})
        return
    }

    resp, err := client.Login(c.Request.Context(), &sdk.LoginRequest{
        LoginType: "email",
        Email:     req.Email,
        Password:  req.Password,
    })

    if err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{"data": resp})
}

// 获取用户信息
func getUserHandler(c *gin.Context) {
    userIDStr := c.Param("id")
    userID, _ := strconv.Atoi(userIDStr)

    user, err := client.GetUser(c.Request.Context(), userID)
    if err != nil {
        c.JSON(404, gin.H{"error": "用户不存在"})
        return
    }

    c.JSON(200, gin.H{"data": user})
}

func main() {
    r := gin.Default()

    r.POST("/register", registerHandler)
    r.POST("/login", loginHandler)
    r.GET("/users/:id", getUserHandler)

    r.Run(":8080")
}
```

#### 示例 2：CLI 工具

```go
package main

import (
    "bufio"
    "context"
    "fmt"
    "log"
    "os"
    "strings"

    "github.com/redis/go-redis/v9"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "github.com/your-org/universal-service-user/sdk"
)

func main() {
    // 初始化 SDK（使用 SQLite）
    db, _ := gorm.Open(sqlite.Open("user.db"), &gorm.Config{})
    rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

    client, err := sdk.New(
        sdk.WithDatabase(db),
        sdk.WithRedis(rdb),
        sdk.WithAutoMigrate(true),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    reader := bufio.NewReader(os.Stdin)

    for {
        fmt.Println("\n=== 用户管理系统 ===")
        fmt.Println("1. 注册")
        fmt.Println("2. 登录")
        fmt.Println("3. 退出")
        fmt.Print("请选择: ")

        var choice string
        fmt.Scanln(&choice)

        switch choice {
        case "1":
            fmt.Print("邮箱: ")
            email, _ := reader.ReadString('\n')
            email = strings.TrimSpace(email)

            fmt.Print("密码: ")
            password, _ := reader.ReadString('\n')
            password = strings.TrimSpace(password)

            fmt.Print("验证码: ")
            code, _ := reader.ReadString('\n')
            code = strings.TrimSpace(code)

            user, err := client.Register(context.Background(), &sdk.RegisterRequest{
                Email:    email,
                Password: password,
                Code:     code,
            })

            if err != nil {
                fmt.Printf("❌ 注册失败: %v\n", err)
            } else {
                fmt.Printf("✅ 注册成功！用户 ID: %d\n", user.ID)
            }

        case "2":
            fmt.Print("邮箱: ")
            email, _ := reader.ReadString('\n')
            email = strings.TrimSpace(email)

            fmt.Print("密码: ")
            password, _ := reader.ReadString('\n')
            password = strings.TrimSpace(password)

            resp, err := client.Login(context.Background(), &sdk.LoginRequest{
                LoginType: "email",
                Email:     email,
                Password:  password,
            })

            if err != nil {
                fmt.Printf("❌ 登录失败: %v\n", err)
            } else {
                fmt.Printf("✅ 登录成功！\n")
                fmt.Printf("用户名: %s\n", resp.UserInfo.Username)
                fmt.Printf("Token: %s\n", resp.AccessToken[:20]+"...")
            }

        case "3":
            fmt.Println("再见！")
            return
        }
    }
}
```

#### 示例 3：微服务集成

```go
package main

import (
    "context"
    "log"

    "github.com/redis/go-redis/v9"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "github.com/your-org/universal-service-user/sdk"
)

// 用户服务（用户中心微服务）
type UserService struct {
    client *sdk.Client
}

func NewUserService() *UserService {
    db, _ := gorm.Open(mysql.Open("user:password@tcp(localhost:3306)/user_db"), &gorm.Config{})
    rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

    client, _ := sdk.New(
        sdk.WithDatabase(db),
        sdk.WithRedis(rdb),
        sdk.WithAutoMigrate(true),
    )

    return &UserService{client: client}
}

// 注册用户
func (s *UserService) Register(ctx context.Context, email, password, code string) (int, error) {
    user, err := s.client.Register(ctx, &sdk.RegisterRequest{
        Email:    email,
        Password: password,
        Code:     code,
    })

    if err != nil {
        return 0, err
    }

    return user.ID, nil
}

// 验证用户
func (s *UserService) ValidateUser(ctx context.Context, email, password string) (int, error) {
    resp, err := s.client.Login(ctx, &sdk.LoginRequest{
        LoginType: "email",
        Email:     email,
        Password:  password,
    })

    if err != nil {
        return 0, err
    }

    return resp.UserInfo.ID, nil
}

// 获取用户信息
func (s *UserService) GetUserInfo(ctx context.Context, userID int) (*sdk.UserInfo, error) {
    return s.client.GetUser(ctx, userID)
}

func main() {
    userService := NewUserService()

    // 在微服务中使用
    userID, err := userService.Register(context.Background(), "user@example.com", "password123", "123456")
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("用户注册成功，ID: %d", userID)
}
```

---

## 三、前端直接对接文档

### 部署服务

#### 1. 启动 API 服务

按照「API 使用文档」中的步骤启动服务：

```bash
# 1. 配置数据库和 Redis
cp config.example.yaml config.yaml
vim config.yaml  # 修改配置

# 2. 运行服务
go run cmd/api/main.go
```

#### 2. 验证服务状态

```bash
# 检查健康状态
curl http://localhost:8080/health

# 预期响应
{"status":"ok"}
```

服务成功启动后，前端可以通过 `http://localhost:8080` 访问 API。

### 前端对接示例

#### React + Axios

##### 1. 创建 API 服务

```javascript
// src/services/api.js
import axios from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
});

// 请求拦截器
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// 响应拦截器：自动刷新 Token
api.interceptors.response.use(
  (response) => response.data,
  async (error) => {
    const originalRequest = error.config;

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        const refreshToken = localStorage.getItem('refresh_token');
        const response = await axios.post(`${API_BASE_URL}/api/v1/auth/refresh`, {
          refresh_token: refreshToken,
        });

        const { access_token, refresh_token } = response.data.data;
        localStorage.setItem('access_token', access_token);
        localStorage.setItem('refresh_token', refresh_token);

        originalRequest.headers.Authorization = `Bearer ${access_token}`;
        return api(originalRequest);
      } catch (refreshError) {
        localStorage.clear();
        window.location.href = '/login';
        return Promise.reject(refreshError);
      }
    }

    return Promise.reject(error);
  }
);

export default api;
```

##### 2. 登录示例

```javascript
// src/services/userService.js
import api from './api';

export const login = async (email, password) => {
  const response = await api.post('/api/v1/auth/login', {
    login_type: 'email',
    email,
    password,
  });

  const { access_token, refresh_token, user_info } = response.data;

  localStorage.setItem('access_token', access_token);
  localStorage.setItem('refresh_token', refresh_token);
  localStorage.setItem('user', JSON.stringify(user_info));

  return user_info;
};

export const register = async (email, password, code) => {
  return await api.post('/api/v1/users/register', {
    email,
    password,
    code,
  });
};

export const logout = async () => {
  const refreshToken = localStorage.getItem('refresh_token');
  await api.post('/api/v1/auth/logout', {
    refresh_token: refreshToken,
  });
  localStorage.clear();
};
```

##### 3. 登录页面

```javascript
// src/pages/Login.jsx
import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { login } from '../services/userService';

const Login = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await login(email, password);
      navigate('/dashboard');
    } catch (err) {
      setError(err.response?.data?.message || '登录失败');
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <input
        type="email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        placeholder="邮箱"
      />
      <input
        type="password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        placeholder="密码"
      />
      {error && <div className="error">{error}</div>}
      <button type="submit">登录</button>
    </form>
  );
};

export default Login;
```

#### Vue 3 + Axios

##### 1. 创建 API 服务

```javascript
// src/api/user.js
import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export default api;
```

##### 2. Pinia Store

```javascript
// src/stores/auth.js
import { defineStore } from 'pinia';
import api from '../api/user';

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: JSON.parse(localStorage.getItem('user')) || null,
  }),

  actions: {
    async login(email, password) {
      const response = await api.post('/api/v1/auth/login', {
        login_type: 'email',
        email,
        password,
      });

      const { access_token, refresh_token, user_info } = response.data;

      localStorage.setItem('access_token', access_token);
      localStorage.setItem('refresh_token', refresh_token);
      localStorage.setItem('user', JSON.stringify(user_info));

      this.user = user_info;
    },

    logout() {
      localStorage.clear();
      this.user = null;
    },
  },
});
```

### Token 管理

#### Token 存储策略

**LocalStorage（简单方便）**

```javascript
localStorage.setItem('access_token', access_token);
const token = localStorage.getItem('access_token');
```

**SessionStorage（更安全）**

```javascript
sessionStorage.setItem('access_token', access_token);
```

#### 自动刷新 Token

```javascript
// 在 axios 响应拦截器中处理
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      const response = await api.post('/api/v1/auth/refresh', {
        refresh_token: localStorage.getItem('refresh_token'),
      });

      const { access_token } = response.data;
      localStorage.setItem('access_token', access_token);

      originalRequest.headers.Authorization = `Bearer ${access_token}`;
      return api(originalRequest);
    }

    return Promise.reject(error);
  }
);
```

#### 路由守卫

**React Router**

```javascript
const ProtectedRoute = ({ children }) => {
  const token = localStorage.getItem('access_token');

  if (!token) {
    return <Navigate to="/login" replace />;
  }

  return children;
};

// 使用
<Route
  path="/dashboard"
  element={
    <ProtectedRoute>
      <Dashboard />
    </ProtectedRoute>
  }
/>
```

**Vue Router**

```javascript
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('access_token');

  if (to.meta.requiresAuth && !token) {
    next('/login');
  } else {
    next();
  }
});
```

---

## 四、配置说明

### 配置文件

项目使用 YAML 格式的配置文件 `config.yaml`。

#### 复制配置文件

```bash
cp config.example.yaml config.yaml
```

#### 配置文件结构

```yaml
# 服务器配置
server:
  port: "8080"              # 服务端口
  mode: "debug"             # 运行模式: debug / release

# 数据库配置
database:
  driver: "postgres"        # 数据库类型: postgres / mysql / sqlite
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "postgres"
  dbname: "universal_service_user"
  auto_migrate: true        # 自动建表

# Redis 配置
redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0

# JWT 配置
jwt:
  secret: "your-jwt-secret"
  expire: 7200              # Access Token 过期时间（秒）
  refresh_expire: 604800    # Refresh Token 过期时间（秒）

# 验证码配置
verification:
  code_length: 6            # 验证码长度
  expire: 300               # 过期时间（秒）
  rate_limit: 60            # 发送间隔（秒）

# 邮件配置
email:
  enabled: true
  smtp:
    host: "smtp.example.com"
    port: 587
    username: "noreply@example.com"
    password: "your-password"
    from: "系统通知 <noreply@example.com>"
  templates:
    register: "您的注册验证码是：{code}，有效期5分钟"
    reset_password: "您的重置密码验证码是：{code}，有效期5分钟"

# 短信配置（可选）
sms:
  enabled: false
  provider: "tencent"       # tencent / aliyun
```

### 环境变量

支持通过环境变量覆盖配置文件中的敏感信息。

#### 环境变量列表

| 环境变量 | 说明 | 对应配置项 |
|---------|------|-----------|
| `SERVER_PORT` | 服务端口 | `server.port` |
| `SERVER_MODE` | 运行模式 | `server.mode` |
| `DATABASE_DRIVER` | 数据库类型 | `database.driver` |
| `DATABASE_HOST` | 数据库主机 | `database.host` |
| `DATABASE_PORT` | 数据库端口 | `database.port` |
| `DATABASE_USER` | 数据库用户 | `database.user` |
| `DATABASE_PASSWORD` | 数据库密码 | `database.password` |
| `DATABASE_DBNAME` | 数据库名称 | `database.dbname` |
| `REDIS_HOST` | Redis 主机 | `redis.host` |
| `REDIS_PORT` | Redis 端口 | `redis.port` |
| `REDIS_PASSWORD` | Redis 密码 | `redis.password` |
| `JWT_SECRET` | JWT 密钥 | `jwt.secret` |
| `EMAIL_SMTP_PASSWORD` | 邮箱密码 | `email.smtp.password` |

#### 使用环境变量

```bash
# 设置环境变量
export JWT_SECRET="your-production-jwt-secret"
export DATABASE_PASSWORD="your-db-password"
export EMAIL_SMTP_PASSWORD="your-email-password"

# 运行服务
./user-service
```

或使用 `.env` 文件：

```bash
# .env
JWT_SECRET=your-production-jwt-secret
DATABASE_PASSWORD=your-db-password
EMAIL_SMTP_PASSWORD=your-email-password
```

### 数据库配置

#### PostgreSQL

```yaml
database:
  driver: "postgres"
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "postgres"
  dbname: "universal_service_user"
  max_open_conns: 100
  max_idle_conns: 10
  conn_max_lifetime: 3600
```

**安装和启动 PostgreSQL**：

```bash
# Ubuntu/Debian
sudo apt-get install postgresql
sudo systemctl start postgresql

# macOS
brew install postgresql
brew services start postgresql

# Windows
# 下载安装包: https://www.postgresql.org/download/windows/
```

#### MySQL

```yaml
database:
  driver: "mysql"
  host: "localhost"
  port: 3306
  user: "root"
  password: "password"
  dbname: "universal_service_user"
```

**安装和启动 MySQL**：

```bash
# Ubuntu/Debian
sudo apt-get install mysql-server
sudo systemctl start mysql

# macOS
brew install mysql
brew services start mysql

# Windows
# 下载安装包: https://dev.mysql.com/downloads/mysql/
```

#### SQLite（开发环境）

```yaml
database:
  driver: "sqlite"
  dbname: "./data/user.db"  # 文件路径
```

SQLite 无需安装，Go 会自动处理。

### Redis 配置

#### 本地 Redis

```yaml
redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0
  pool_size: 10
```

**安装和启动 Redis**：

```bash
# Ubuntu/Debian
sudo apt-get install redis-server
sudo systemctl start redis-server

# macOS
brew install redis
brew services start redis

# Windows
# 下载 Redis for Windows: https://github.com/microsoftarchive/redis/releases
# 或使用 WSL 运行 Linux 版本
```

#### Redis Cluster

```yaml
redis:
  addrs:
    - "redis-node1:6379"
    - "redis-node2:6379"
    - "redis-node3:6379"
  password: "your-password"
```

### 邮件/短信配置

#### 邮件配置

##### Gmail

```yaml
email:
  enabled: true
  smtp:
    host: "smtp.gmail.com"
    port: 587
    username: "your-email@gmail.com"
    password: "your-app-password"  # 使用应用专用密码
    from: "系统通知 <your-email@gmail.com>"
```

##### QQ 邮箱

```yaml
email:
  enabled: true
  smtp:
    host: "smtp.qq.com"
    port: 587
    username: "your-email@qq.com"
    password: "your-authorization-code"  # 使用授权码
    from: "系统通知 <your-email@qq.com>"
```

##### 163 邮箱

```yaml
email:
  enabled: true
  smtp:
    host: "smtp.163.com"
    port: 465
    username: "your-email@163.com"
    password: "your-authorization-code"
    from: "系统通知 <your-email@163.com>"
```

#### 短信配置

##### 腾讯云短信

```yaml
sms:
  enabled: true
  provider: "tencent"
  tencent:
    secret_id: "your-secret-id"
    secret_key: "your-secret-key"
    app_id: "your-app-id"
    sign: "您的签名"
    templates:
      register: "123456"        # 模板 ID
      reset_password: "123456"
      login: "123456"
```

##### 阿里云短信

```yaml
sms:
  enabled: true
  provider: "aliyun"
  aliyun:
    access_key_id: "your-access-key-id"
    access_key_secret: "your-access-key-secret"
    sign_name: "您的签名"
    templates:
      register: "SMS_123456789"
      reset_password: "SMS_987654321"
      login: "SMS_111111111"
```

---

## 五、部署指南

### 二进制部署

#### 编译

```bash
# Linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o user-service-linux-amd64 cmd/api/main.go

# macOS
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o user-service-darwin-amd64 cmd/api/main.go

# Windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o user-service-windows-amd64.exe cmd/api/main.go
```

#### 运行

```bash
# 1. 创建配置文件
cp config.example.yaml config.yaml

# 2. 编辑配置文件
vim config.yaml

# 3. 运行
./user-service-linux-amd64
```

#### 使用 Systemd 管理（Linux）

创建服务文件 `/etc/systemd/system/user-service.service`:

```ini
[Unit]
Description=Universal User Service
After=network.target postgresql.service redis.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/user-service
ExecStart=/opt/user-service/user-service
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

启动服务:

```bash
# 重载配置
sudo systemctl daemon-reload

# 启动服务
sudo systemctl start user-service

# 开机自启
sudo systemctl enable user-service

# 查看状态
sudo systemctl status user-service

# 查看日志
sudo journalctl -u user-service -f
```

### 源码部署

#### 克隆项目

```bash
git clone https://github.com/your-org/universal-service-user.git
cd universal-service-user
```

#### 安装依赖

```bash
go mod download
```

#### 配置

```bash
cp config.example.yaml config.yaml
vim config.yaml
```

#### 运行

```bash
go run cmd/api/main.go
```

### 生产环境建议

#### 1. 使用 Nginx 反向代理

**nginx.conf**:

```nginx
upstream user_service {
    server localhost:8080;
}

server {
    listen 80;
    server_name your-domain.com;

    # 重定向到 HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    # 日志
    access_log /var/log/nginx/user-service-access.log;
    error_log /var/log/nginx/user-service-error.log;

    # 代理设置
    location / {
        proxy_pass http://user_service;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # 超时设置
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # 静态文件缓存
    location ~* \.(jpg|jpeg|png|gif|ico|css|js)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
}
```

#### 2. 使用 PM2 管理（Node.js 生态）

```bash
# 安装 PM2
npm install -g pm2

# 启动服务
pm2 start user-service --name user-service

# 查看状态
pm2 status

# 查看日志
pm2 logs user-service

# 重启服务
pm2 restart user-service

# 开机自启
pm2 startup
pm2 save
```

#### 3. 监控和日志

**健康检查**:

```bash
# 添加健康检查
curl http://localhost:8080/health
```

**日志收集**:

```yaml
# config.yaml
logging:
  level: "info"                # 生产环境使用 info
  format: "json"               # 使用 JSON 格式便于解析
  enable_sensible_mask: true   # 脱敏敏感信息
```

**集成监控系统**:

- Prometheus + Grafana
- ELK Stack (Elasticsearch, Logstash, Kibana)
- 阿里云日志服务

#### 4. 安全加固

**修改默认密钥**:

```yaml
jwt:
  secret: "your-strong-random-secret-key-at-least-32-characters"
```

**启用 HTTPS**:

使用 Let's Encrypt 免费证书：

```bash
# 安装 certbot
sudo apt-get install certbot

# 获取证书
sudo certbot certonly --standalone -d your-domain.com

# 自动续期
sudo certbot renew --dry-run
```

**防火墙配置**:

```bash
# 只开放必要端口
sudo ufw allow 22/tcp   # SSH
sudo ufw allow 80/tcp   # HTTP
sudo ufw allow 443/tcp  # HTTPS
sudo ufw enable
```

#### 5. 性能优化

**数据库连接池**:

```yaml
database:
  max_open_conns: 100    # 最大打开连接数
  max_idle_conns: 10     # 最大空闲连接数
  conn_max_lifetime: 3600  # 连接最大生命周期（秒）
```

**Redis 连接池**:

```yaml
redis:
  pool_size: 10
```

**启用 Gzip**:

在 Nginx 中启用：

```nginx
gzip on;
gzip_vary on;
gzip_min_length 1024;
gzip_types text/plain text/css text/xml text/javascript application/json application/javascript;
```

#### 6. 备份策略

**数据库备份**:

```bash
#!/bin/bash
# backup.sh

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backups/postgres"

# PostgreSQL 备份
pg_dump -U postgres universal_service_user > $BACKUP_DIR/user_$DATE.sql

# MySQL 备份
# mysqldump -u root -p universal_service_user > $BACKUP_DIR/user_$DATE.sql

# 压缩备份
gzip $BACKUP_DIR/user_$DATE.sql

# 删除7天前的备份
find $BACKUP_DIR -name "*.sql.gz" -mtime +7 -delete
```

**定时备份**：

```bash
# 添加到 crontab
crontab -e

# 每天凌晨 2 点备份
0 2 * * * /opt/backup.sh
```

---

## 六、常见问题

### 1. 数据库连接失败

**错误**: `failed to connect to database`

**解决方案**:

- 检查数据库是否启动
- 检查配置文件中的连接信息
- 检查防火墙设置
- 检查数据库用户权限

### 2. Redis 连接失败

**错误**: `failed to connect to redis`

**解决方案**:

- 检查 Redis 是否启动
- 检查 Redis 密码配置
- 检查 Redis 地址和端口

### 3. 验证码发送失败

**错误**: `failed to send verification code`

**解决方案**:

- 检查邮件/短信服务配置
- 检查网络连接
- 检查服务商账户余额
- 查看详细错误日志

### 4. Token 验证失败

**错误**: `invalid token`

**解决方案**:

- 检查 JWT 密钥配置
- 检查 Token 是否过期
- 检查 Token 格式是否正确

### 5. CORS 跨域问题

**错误**: `CORS policy: No 'Access-Control-Allow-Origin' header`

**解决方案**:

- 检查服务端 CORS 配置
- 检查前端请求域名是否在允许列表中

---

## 七、更新日志

### v1.0.0 (2024-01-01)

- ✅ 用户注册、登录、登出
- ✅ 邮箱/手机号验证码
- ✅ JWT 认证
- ✅ 用户信息管理
- ✅ 密码重置
- ✅ SDK 模式支持
- ✅ 钩子系统
- ✅ 登录防刷

---

## 八、技术支持

- **GitHub Issues**: https://github.com/your-org/universal-service-user/issues
- **文档**: https://docs.your-domain.com
- **邮箱**: support@your-domain.com

---

**感谢使用 Universal User Service！**
