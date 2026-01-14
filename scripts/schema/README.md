# 租户用户数据库建表脚本

本目录包含 Universal User Service 租户用户数据库的建表脚本。

## 📋 目录结构

```
scripts/schema/
├── README.md        # 本文件
├── postgresql.sql   # PostgreSQL 建表脚本
├── mysql.sql        # MySQL 建表脚本
└── sqlite.sql       # SQLite 建表脚本
```

## 🗄️ 数据库表说明

建表脚本会创建以下5张表：

### 1. users - 用户主表
存储用户基本信息

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INT | 主键 |
| username | VARCHAR(64) | 用户名（唯一，可空） |
| nickname | VARCHAR(64) | 昵称 |
| avatar | VARCHAR(512) | 头像URL |
| email | VARCHAR(255) | 邮箱（唯一，可空） |
| phone | VARCHAR(20) | 手机号（唯一，可空） |
| password_hash | VARCHAR(255) | 密码哈希 |
| status | INT | 状态（0:未激活 1:正常 2:禁用 3:锁定） |
| created_at | TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 更新时间 |
| deleted_at | TIMESTAMP | 软删除时间 |
| version | INT | 版本号 |

### 2. user_profiles - 用户扩展信息表
存储用户扩展信息（一对一关系）

### 3. user_login_logs - 用户登录日志表
记录用户登录历史

### 4. user_sessions - 会话表
存储用户会话信息

### 5. user_oauth_bindings - OAuth绑定表
存储第三方账号绑定信息

## 🚀 使用方法

### PostgreSQL

```bash
# 方式一：使用 psql 命令
psql -U postgres -d your_tenant_database -f postgresql.sql

# 方式二：登录后执行
psql -U postgres -d your_tenant_database
\i postgresql.sql

# 方式三：从 stdin
cat postgresql.sql | psql -U postgres -d your_tenant_database
```

### MySQL

```bash
# 方式一：使用 mysql 命令
mysql -u root -p your_tenant_database < mysql.sql

# 方式二：登录后执行
mysql -u root -p your_tenant_database
source mysql.sql;

# 方式三：从 stdin
cat mysql.sql | mysql -u root -p your_tenant_database
```

### SQLite

```bash
# 方式一：使用 sqlite3 命令
sqlite3 your_tenant.db < sqlite.sql

# 方式二：登录后执行
sqlite3 your_tenant.db
.read sqlite.sql

# 方式三：从 stdin
cat sqlite.sql | sqlite3 your_tenant.db
```

## 📝 注意事项

### 1. 数据库创建

在执行建表脚本前，请确保已创建数据库：

**PostgreSQL:**
```bash
createdb -U postgres your_tenant_database
```

**MySQL:**
```sql
CREATE DATABASE your_tenant_database
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;
```

**SQLite:**
```bash
# SQLite 会自动创建数据库文件，无需预先创建
touch your_tenant.db
```

### 2. 字符集

- **PostgreSQL**: 默认使用 UTF-8 编码
- **MySQL**: 使用 `utf8mb4` 字符集和 `utf8mb4_unicode_ci` 排序规则
- **SQLite**: 默认支持 UTF-8

### 3. 索引

所有脚本已自动创建必要的索引：
- 唯一索引：username, email, phone
- 普通索引：status, deleted_at, user_id 等
- 复合索引：user_id + login_at 等

### 4. 外键约束

- 使用 `ON DELETE CASCADE`：当用户被删除时，自动删除关联记录
- PostgreSQL/MySQL/SQLite 都支持外键约束

### 5. 软删除

users 表和 user_oauth_bindings 表支持软删除：
- 使用 `deleted_at` 字段标记删除时间
- NULL 表示未删除
- 非NULL表示已删除

### 6. 时间戳

- PostgreSQL: 使用 `TIMESTAMP` 类型
- MySQL: 使用 `TIMESTAMP` 类型，自动更新
- SQLite: 使用 `DATETIME` 类型

## 🔧 自动建表（推荐）

如果您不想手动执行 SQL 脚本，可以在租户配置文件中启用自动建表：

```yaml
user_database:
  driver: "postgres"
  host: "your-host"
  port: 5432
  user: "your-user"
  password: "your-password"
  dbname: "your-database"
  auto_create_tables: true  # 自动创建表
```

服务会自动使用 GORM 的 AutoMigrate 功能创建表结构。

## ⚠️ 生产环境建议

### 1. 数据库用户权限

为租户数据库创建专用用户，只授予必要权限：

```sql
-- PostgreSQL
CREATE USER tenant_user WITH PASSWORD 'secure_password';
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO tenant_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO tenant_user;

-- MySQL
CREATE USER 'tenant_user'@'%' IDENTIFIED BY 'secure_password';
GRANT SELECT, INSERT, UPDATE, DELETE ON your_tenant_database.* TO 'tenant_user'@'%';
FLUSH PRIVILEGES;
```

### 2. 备份策略

定期备份租户数据库：

```bash
# PostgreSQL
pg_dump -U postgres your_tenant_database > backup_$(date +%Y%m%d).sql

# MySQL
mysqldump -u root -p your_tenant_database > backup_$(date +%Y%m%d).sql

# SQLite
cp your_tenant.db backup_$(date +%Y%m%d).db
```

### 3. 性能优化

- 根据查询需求添加额外索引
- 定期清理登录日志（如保留最近3个月）
- 定期清理软删除的数据
- 使用连接池（已内置）

## 📚 相关文档

- [完整使用文档](../../Universal User Service - 使用文档.md)
- [多租户配置中心](../../Universal User Service - 使用文档.md#三多租户配置中心)

## 🆘 常见问题

### Q: 表已存在怎么办？

A: 脚本使用了 `CREATE TABLE IF NOT EXISTS`，可以安全重复执行。如果表结构需要更新，请手动处理。

### Q: 如何查看表结构？

```bash
# PostgreSQL
\dt your_tenant_database

# MySQL
DESCRIBE users;

# SQLite
.schema users
```

### Q: 如何清空表数据？

```sql
-- 清空所有用户数据（慎用！）
TRUNCATE TABLE user_login_logs;
TRUNCATE TABLE user_sessions;
TRUNCATE TABLE user_oauth_bindings;
TRUNCATE TABLE user_profiles;
TRUNCATE TABLE users;
```

### Q: 自动建表失败怎么办？

A:
1. 检查数据库用户权限
2. 检查数据库连接配置
3. 查看服务日志获取详细错误信息
4. 手动执行对应的 SQL 脚本

## 📞 技术支持

如有问题，请查阅：
- GitHub Issues: https://github.com/your-org/universal-service-user/issues
- 完整文档: `../../Universal User Service - 使用文档.md`
