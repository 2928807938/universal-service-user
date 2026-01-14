-- ============================================================================
-- Universal User Service - 租户用户数据库建表脚本 (PostgreSQL)
-- ============================================================================
--
-- 说明：
--   此脚本用于创建租户的用户数据库表结构
--   每个租户应该有自己独立的数据库
--
-- 使用方法：
--   psql -U postgres -d tenant_database -f postgresql.sql
--
-- ============================================================================

-- ----------------------------------------------------------------------------
-- 1. 用户主表 (users)
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(64) UNIQUE,
    nickname VARCHAR(64),
    avatar VARCHAR(512),
    email VARCHAR(255) UNIQUE,
    phone VARCHAR(20) UNIQUE,
    password_hash VARCHAR(255),
    status INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    version INT NOT NULL DEFAULT 1
);

-- 索引
CREATE UNIQUE INDEX idx_users_username ON users(username) WHERE username IS NOT NULL;
CREATE UNIQUE INDEX idx_users_email ON users(email) WHERE email IS NOT NULL;
CREATE UNIQUE INDEX idx_users_phone ON users(phone) WHERE phone IS NOT NULL;
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_deleted_at ON users(deleted_at) WHERE deleted_at IS NOT NULL;

-- 注释
COMMENT ON TABLE users IS '用户主表';
COMMENT ON COLUMN users.username IS '用户名（唯一，可空）';
COMMENT ON COLUMN users.nickname IS '昵称';
COMMENT ON COLUMN users.avatar IS '头像URL';
COMMENT ON COLUMN users.email IS '邮箱（唯一，可空）';
COMMENT ON COLUMN users.phone IS '手机号（唯一，可空）';
COMMENT ON COLUMN users.password_hash IS '密码哈希（第三方登录用户可空）';
COMMENT ON COLUMN users.status IS '状态（0:未激活 1:正常 2:禁用 3:锁定）';

-- ----------------------------------------------------------------------------
-- 2. 用户扩展信息表 (user_profiles)
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS user_profiles (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL UNIQUE,
    bio VARCHAR(500),
    gender INT NOT NULL DEFAULT 0,
    birthday DATE,
    location VARCHAR(128),
    company VARCHAR(128),
    position VARCHAR(64),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    version INT NOT NULL DEFAULT 1,
    CONSTRAINT fk_user_profiles_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 注释
COMMENT ON TABLE user_profiles IS '用户扩展信息表';
COMMENT ON COLUMN user_profiles.gender IS '性别（0:未知 1:男 2:女）';

-- ----------------------------------------------------------------------------
-- 3. 用户登录日志表 (user_login_logs)
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS user_login_logs (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    login_type VARCHAR(32) NOT NULL,
    login_at TIMESTAMP NOT NULL,
    logout_at TIMESTAMP,
    ip VARCHAR(64),
    device_type VARCHAR(32),
    device_name VARCHAR(128),
    browser VARCHAR(64),
    os VARCHAR(32),
    status INT NOT NULL,
    fail_reason VARCHAR(128),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    version INT NOT NULL DEFAULT 1,
    CONSTRAINT fk_user_login_logs_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 索引
CREATE INDEX idx_login_logs_user ON user_login_logs(user_id);
CREATE INDEX idx_login_logs_user_login_at ON user_login_logs(user_id, login_at);
CREATE INDEX idx_login_logs_ip ON user_login_logs(ip);

-- 注释
COMMENT ON TABLE user_login_logs IS '用户登录日志表';
COMMENT ON COLUMN user_login_logs.login_type IS '登录方式（username/email/phone/wechat/alipay/qq）';
COMMENT ON COLUMN user_login_logs.status IS '状态（1:成功 2:失败）';

-- ----------------------------------------------------------------------------
-- 4. 会话表 (user_sessions)
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS user_sessions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    refresh_token_hash VARCHAR(128) NOT NULL UNIQUE,
    device_type VARCHAR(32),
    device_name VARCHAR(128),
    ip VARCHAR(64),
    last_active_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    version INT NOT NULL DEFAULT 1,
    CONSTRAINT fk_user_sessions_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 索引
CREATE INDEX idx_sessions_user ON user_sessions(user_id);
CREATE INDEX idx_sessions_user_active ON user_sessions(user_id, last_active_at);

-- 注释
COMMENT ON TABLE user_sessions IS '用户会话表';

-- ----------------------------------------------------------------------------
-- 5. OAuth 绑定表 (user_oauth_bindings)
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS user_oauth_bindings (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    provider VARCHAR(32) NOT NULL,
    open_id VARCHAR(128) NOT NULL,
    union_id VARCHAR(128),
    nickname VARCHAR(64),
    avatar VARCHAR(512),
    access_token VARCHAR(512),
    refresh_token VARCHAR(512),
    expires_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    version INT NOT NULL DEFAULT 1,
    CONSTRAINT fk_user_oauth_bindings_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 索引
CREATE INDEX idx_oauth_user_provider ON user_oauth_bindings(user_id, provider);
CREATE INDEX idx_oauth_open_id ON user_oauth_bindings(open_id);
CREATE INDEX idx_oauth_union_id ON user_oauth_bindings(union_id) WHERE union_id IS NOT NULL;
CREATE INDEX idx_oauth_deleted_at ON user_oauth_bindings(deleted_at) WHERE deleted_at IS NOT NULL;

-- 注释
COMMENT ON TABLE user_oauth_bindings IS 'OAuth第三方账号绑定表';
COMMENT ON COLUMN user_oauth_bindings.provider IS '平台（wechat/alipay/qq）';

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