-- ============================================================================
-- Universal User Service - 租户用户数据库建表脚本 (SQLite)
-- ============================================================================
--
-- 说明：
--   此脚本用于创建租户的用户数据库表结构
--   每个租户应该有自己独立的数据库文件
--
-- 使用方法：
--   sqlite3 tenant.db < sqlite.sql
--   或
--   sqlite3 tenant.db
--   sqlite> .read sqlite.sql
--
-- ============================================================================

-- ----------------------------------------------------------------------------
-- 1. 用户主表 (users)
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE,
    nickname TEXT,
    avatar TEXT,
    email TEXT UNIQUE,
    phone TEXT UNIQUE,
    password_hash TEXT,
    status INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at DATETIME NOT NULL DEFAULT (datetime('now')),
    deleted_at DATETIME,
    version INTEGER NOT NULL DEFAULT 1
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username) WHERE username IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email) WHERE email IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone) WHERE phone IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at) WHERE deleted_at IS NOT NULL;

-- ----------------------------------------------------------------------------
-- 2. 用户扩展信息表 (user_profiles)
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS user_profiles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL UNIQUE,
    bio TEXT,
    gender INTEGER NOT NULL DEFAULT 0,
    birthday DATE,
    location TEXT,
    company TEXT,
    position TEXT,
    created_at DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at DATETIME NOT NULL DEFAULT (datetime('now')),
    version INTEGER NOT NULL DEFAULT 1,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- ----------------------------------------------------------------------------
-- 3. 用户登录日志表 (user_login_logs)
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS user_login_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    login_type TEXT NOT NULL,
    login_at DATETIME NOT NULL,
    logout_at DATETIME,
    ip TEXT,
    device_type TEXT,
    device_name TEXT,
    browser TEXT,
    os TEXT,
    status INTEGER NOT NULL,
    fail_reason TEXT,
    created_at DATETIME NOT NULL DEFAULT (datetime('now')),
    version INTEGER NOT NULL DEFAULT 1,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_login_logs_user ON user_login_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_login_logs_user_login_at ON user_login_logs(user_id, login_at);
CREATE INDEX IF NOT EXISTS idx_login_logs_ip ON user_login_logs(ip);

-- ----------------------------------------------------------------------------
-- 4. 会话表 (user_sessions)
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS user_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    refresh_token_hash TEXT NOT NULL UNIQUE,
    device_type TEXT,
    device_name TEXT,
    ip TEXT,
    last_active_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at DATETIME NOT NULL DEFAULT (datetime('now')),
    version INTEGER NOT NULL DEFAULT 1,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_sessions_user ON user_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_user_active ON user_sessions(user_id, last_active_at);

-- ----------------------------------------------------------------------------
-- 5. OAuth 绑定表 (user_oauth_bindings)
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS user_oauth_bindings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    provider TEXT NOT NULL,
    open_id TEXT NOT NULL,
    union_id TEXT,
    nickname TEXT,
    avatar TEXT,
    access_token TEXT,
    refresh_token TEXT,
    expires_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at DATETIME NOT NULL DEFAULT (datetime('now')),
    deleted_at DATETIME,
    version INTEGER NOT NULL DEFAULT 1,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_oauth_user_provider ON user_oauth_bindings(user_id, provider);
CREATE INDEX IF NOT EXISTS idx_oauth_open_id ON user_oauth_bindings(open_id);
CREATE INDEX IF NOT EXISTS idx_oauth_union_id ON user_oauth_bindings(union_id) WHERE union_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_oauth_deleted_at ON user_oauth_bindings(deleted_at) WHERE deleted_at IS NOT NULL;

-- ============================================================================
-- 说明
-- ============================================================================
--
-- 用户状态枚举 (users.status):
--   0: INACTIVE  - 未激活
--   1: ACTIVE    - 正常
--   2: DISABLED  - 禁用（管理员手动禁用）
--   3: LOCKED    - 锁定（多次登录失败自动锁定）
--
-- 性别枚举 (user_profiles.gender):
--   0: 未知
--   1: 男
--   2: 女
--
-- 登录状态枚举 (user_login_logs.status):
--   1: 成功
--   2: 失败
--
-- 登录方式枚举 (user_login_logs.login_type):
--   username - 用户名+密码
--   email    - 邮箱+密码
--   phone    - 手机号+验证码
--   wechat   - 微信登录
--   alipay   - 支付宝登录
--   qq       - QQ登录
--
-- OAuth 平台 (user_oauth_bindings.provider):
--   wechat - 微信
--   alipay - 支付宝
--   qq     - QQ
--
-- ============================================================================
