package gorm

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DatabaseStrategy 数据库连接策略接口
type DatabaseStrategy interface {
	// GetDialector 获取数据库方言器
	GetDialector(config *DatabaseConfig) (gorm.Dialector, error)
	// ValidateConfig 验证配置是否有效
	ValidateConfig(config *DatabaseConfig) error
	// GetDefaultPort 获取默认端口
	GetDefaultPort() int
}

// MySQLStrategy MySQL 数据库策略
type MySQLStrategy struct{}

func (s *MySQLStrategy) GetDialector(config *DatabaseConfig) (gorm.Dialector, error) {
	charset := config.Charset
	if charset == "" {
		charset = "utf8mb4"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
		charset,
	)

	return mysql.Open(dsn), nil
}

func (s *MySQLStrategy) ValidateConfig(config *DatabaseConfig) error {
	if config.Host == "" {
		return fmt.Errorf("MySQL host is required")
	}
	if config.Database == "" {
		return fmt.Errorf("MySQL database name is required")
	}
	return nil
}

func (s *MySQLStrategy) GetDefaultPort() int {
	return 3306
}

// PostgreSQLStrategy PostgreSQL 数据库策略
type PostgreSQLStrategy struct{}

func (s *PostgreSQLStrategy) GetDialector(config *DatabaseConfig) (gorm.Dialector, error) {
	sslMode := config.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.Username,
		config.Password,
		config.Database,
		sslMode,
	)

	return postgres.Open(dsn), nil
}

func (s *PostgreSQLStrategy) ValidateConfig(config *DatabaseConfig) error {
	if config.Host == "" {
		return fmt.Errorf("PostgreSQL host is required")
	}
	if config.Database == "" {
		return fmt.Errorf("PostgreSQL database name is required")
	}
	return nil
}

func (s *PostgreSQLStrategy) GetDefaultPort() int {
	return 5432
}

// SQLiteStrategy SQLite 数据库策略
type SQLiteStrategy struct{}

func (s *SQLiteStrategy) GetDialector(config *DatabaseConfig) (gorm.Dialector, error) {
	if config.Database == "" {
		return nil, fmt.Errorf("SQLite database file path is required")
	}
	return sqlite.Open(config.Database), nil
}

func (s *SQLiteStrategy) ValidateConfig(config *DatabaseConfig) error {
	if config.Database == "" {
		return fmt.Errorf("SQLite database file path is required")
	}
	return nil
}

func (s *SQLiteStrategy) GetDefaultPort() int {
	return 0 // SQLite 不需要端口
}

// StrategyRegistry 策略注册器
type StrategyRegistry struct {
	strategies map[DatabaseType]DatabaseStrategy
}

// NewStrategyRegistry 创建策略注册器
func NewStrategyRegistry() *StrategyRegistry {
	registry := &StrategyRegistry{
		strategies: make(map[DatabaseType]DatabaseStrategy),
	}

	// 注册默认策略
	registry.Register(MySQL, &MySQLStrategy{})
	registry.Register(PostgreSQL, &PostgreSQLStrategy{})
	registry.Register(SQLite, &SQLiteStrategy{})

	return registry
}

// Register 注册策略
func (r *StrategyRegistry) Register(dbType DatabaseType, strategy DatabaseStrategy) {
	r.strategies[dbType] = strategy
}

// Get 获取策略
func (r *StrategyRegistry) Get(dbType DatabaseType) (DatabaseStrategy, error) {
	strategy, exists := r.strategies[dbType]
	if !exists {
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
	return strategy, nil
}

// GetSupportedTypes 获取支持的数据库类型列表
func (r *StrategyRegistry) GetSupportedTypes() []DatabaseType {
	types := make([]DatabaseType, 0, len(r.strategies))
	for dbType := range r.strategies {
		types = append(types, dbType)
	}
	return types
}
