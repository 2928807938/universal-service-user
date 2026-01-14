-- ============================================================================
-- Universal User Service - 租户用户数据库建表脚本 (MySQL)
-- ============================================================================
--
-- 说明：
--   此脚本用于创建租户的用户数据库表结构
--   每个租户应该有自己独立的数据库
--
-- 使用方法：
--   mysql -u root -p tenant_database < mysql.sql
--
-- ============================================================================

-- ----------------------------------------------------------------------------
-- 1. 用户主表 (users)
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(64) UNIQUE,
    nickname VARCHAR(64),
    avatar VARCHAR(512),
    email VARCHAR(255) UNIQUE,
    phone VARCHAR(20) UNIQUE,
    password_hash VARCHAR(255),
    status INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    version INT NOT NULL DEFAULT 1,
    INDEX idx_users_username (username),
    INDEX idx_users_email (email),
    INDEX idx_users_phone (phone),
    INDEX idx_users_status (status),
    INDEX idx_users_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- 2. 用户扩展信息表 (user_profiles)
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS user_profiles (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL UNIQUE,
    bio VARCHAR(500),
    gender INT NOT NULL DEFAULT 0,
    birthday DATE NULL,
    location VARCHAR(128),
    company VARCHAR(128),
    position VARCHAR(64),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    version INT NOT NULL DEFAULT 1,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- 3. 用户登录日志表 (user_login_logs)
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS user_login_logs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    login_type VARCHAR(32) NOT NULL,
    login_at TIMESTAMP NOT NULL,
    logout_at TIMESTAMP NULL,
    ip VARCHAR(64),
    device_type VARCHAR(32),
    device_name VARCHAR(128),
    browser VARCHAR(64),
    os VARCHAR(32),
    status INT NOT NULL,
    fail_reason VARCHAR(128),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    version INT NOT NULL DEFAULT 1,
    INDEX idx_login_logs_user (user_id),
    INDEX idx_login_logs_user_login_at (user_id, login_at),
    INDEX idx_login_logs_ip (ip),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- 4. 会话表 (user_sessions)
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS user_sessions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    refresh_token_hash VARCHAR(128) NOT NULL UNIQUE,
    device_type VARCHAR(32),
    device_name VARCHAR(128),
    ip VARCHAR(64),
    last_active_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    version INT NOT NULL DEFAULT 1,
    INDEX idx_sessions_user (user_id),
    INDEX idx_sessions_user_active (user_id, last_active_at),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- 5. OAuth 绑定表 (user_oauth_bindings)
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS user_oauth_bindings (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    provider VARCHAR(32) NOT NULL,
    open_id VARCHAR(128) NOT NULL,
    union_id VARCHAR(128),
    nickname VARCHAR(64),
    avatar VARCHAR(512),
    access_token VARCHAR(512),
    refresh_token VARCHAR(512),
    expires_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    version INT NOT NULL DEFAULT 1,
    INDEX idx_oauth_user_provider (user_id, provider),
    INDEX idx_oauth_open_id (open_id),
    INDEX idx_oauth_union_id (union_id),
    INDEX idx_oauth_deleted_at (deleted_at),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
