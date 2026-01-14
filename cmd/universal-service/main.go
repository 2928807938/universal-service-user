package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	domainEntity "universal-service-user/configcenter/domain/entity"
	domainRepo "universal-service-user/configcenter/domain/repository"
	infraRepo "universal-service-user/configcenter/infrastructure/repository"
)

type cliConfig struct {
	ConfigPath string
	UploadOnly bool

	DBDriver  string
	DBHost    string
	DBPort    int
	DBUser    string
	DBPass    string
	DBName    string
	DBSSLMode string
	DBDSN     string
}

func main() {
	cfg := parseFlags()
	if !cfg.UploadOnly {
		log.Fatal("upload-only mode is required (-u or --upload-only)")
	}
	if cfg.ConfigPath == "" {
		log.Fatal("config file path is required (-c or --config)")
	}

	tenantID, environment, configJSON, err := loadTenantConfig(cfg.ConfigPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := openDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	appRepo := infraRepo.NewAppRepositoryImpl(db)
	if err := ensureTenantExists(context.Background(), appRepo, tenantID); err != nil {
		log.Fatalf("tenant validation failed: %v", err)
	}

	if err := uploadConfig(context.Background(), db, tenantID, environment, configJSON); err != nil {
		log.Fatalf("upload failed: %v", err)
	}

	log.Printf("upload completed: tenant_id=%s environment=%s", tenantID, environment)
}

func parseFlags() *cliConfig {
	cfg := &cliConfig{}
	flag.StringVar(&cfg.ConfigPath, "config", "", "path to config.yaml")
	flag.StringVar(&cfg.ConfigPath, "c", "", "path to config.yaml")
	flag.BoolVar(&cfg.UploadOnly, "upload-only", false, "upload config then exit")
	flag.BoolVar(&cfg.UploadOnly, "u", false, "upload config then exit")

	flag.StringVar(&cfg.DBDriver, "db-driver", envOrDefault("CONFIG_DB_DRIVER", "postgres"), "database driver")
	flag.StringVar(&cfg.DBHost, "db-host", envOrDefault("CONFIG_DB_HOST", "localhost"), "database host")
	flag.IntVar(&cfg.DBPort, "db-port", envIntOrDefault("CONFIG_DB_PORT", 5432), "database port")
	flag.StringVar(&cfg.DBUser, "db-user", envOrDefault("CONFIG_DB_USER", "postgres"), "database user")
	flag.StringVar(&cfg.DBPass, "db-password", envOrDefault("CONFIG_DB_PASSWORD", ""), "database password")
	flag.StringVar(&cfg.DBName, "db-name", envOrDefault("CONFIG_DB_NAME", ""), "database name")
	flag.StringVar(&cfg.DBSSLMode, "db-sslmode", envOrDefault("CONFIG_DB_SSLMODE", "disable"), "postgres sslmode")
	flag.StringVar(&cfg.DBDSN, "db-dsn", envOrDefault("CONFIG_DB_DSN", ""), "database DSN")

	flag.Parse()
	return cfg
}

func loadTenantConfig(path string) (string, string, json.RawMessage, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", "", nil, err
	}

	var cfg map[string]interface{}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return "", "", nil, err
	}

	appSection, ok := normalizeMap(cfg["app"])
	if !ok {
		return "", "", nil, fmt.Errorf("missing app section")
	}

	tenantID, _ := appSection["tenant_id"].(string)
	environment, _ := appSection["environment"].(string)
	tenantID = strings.TrimSpace(tenantID)
	environment = strings.TrimSpace(environment)
	if tenantID == "" || environment == "" {
		return "", "", nil, fmt.Errorf("app.tenant_id and app.environment are required")
	}

	normalized := normalizeValue(cfg)
	jsonData, err := json.Marshal(normalized)
	if err != nil {
		return "", "", nil, err
	}

	return tenantID, environment, jsonData, nil
}

func normalizeMap(value interface{}) (map[string]interface{}, bool) {
	switch v := value.(type) {
	case map[string]interface{}:
		return v, true
	case map[interface{}]interface{}:
		out := make(map[string]interface{}, len(v))
		for key, val := range v {
			keyStr, ok := key.(string)
			if !ok {
				continue
			}
			out[keyStr] = val
		}
		return out, true
	default:
		return nil, false
	}
}

func normalizeValue(value interface{}) interface{} {
	switch v := value.(type) {
	case map[string]interface{}:
		out := make(map[string]interface{}, len(v))
		for key, val := range v {
			out[key] = normalizeValue(val)
		}
		return out
	case map[interface{}]interface{}:
		out := make(map[string]interface{}, len(v))
		for key, val := range v {
			keyStr, ok := key.(string)
			if !ok {
				continue
			}
			out[keyStr] = normalizeValue(val)
		}
		return out
	case []interface{}:
		out := make([]interface{}, len(v))
		for i, item := range v {
			out[i] = normalizeValue(item)
		}
		return out
	default:
		return value
	}
}

func openDB(cfg *cliConfig) (*gorm.DB, error) {
	if cfg.DBDSN != "" {
		return gorm.Open(selectDialector(cfg.DBDriver, cfg.DBDSN), &gorm.Config{})
	}
	dsn, err := buildDSN(cfg)
	if err != nil {
		return nil, err
	}
	return gorm.Open(selectDialector(cfg.DBDriver, dsn), &gorm.Config{})
}

func selectDialector(driver, dsn string) gorm.Dialector {
	switch driver {
	case "mysql":
		return mysql.Open(dsn)
	case "sqlite":
		return sqlite.Open(dsn)
	default:
		return postgres.Open(dsn)
	}
}

func buildDSN(cfg *cliConfig) (string, error) {
	switch cfg.DBDriver {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName), nil
	case "sqlite":
		if cfg.DBName == "" {
			return "", fmt.Errorf("db-name is required for sqlite")
		}
		return cfg.DBName, nil
	default:
		if cfg.DBName == "" {
			return "", fmt.Errorf("db-name is required")
		}
		return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
			cfg.DBHost, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBPort, cfg.DBSSLMode), nil
	}
}

func ensureTenantExists(ctx context.Context, repo domainRepo.AppRepository, tenantID string) error {
	exists, err := repo.ExistsByTenantID(ctx, tenantID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("tenant_id not found: %s", tenantID)
	}
	return nil
}

func uploadConfig(ctx context.Context, db *gorm.DB, tenantID, environment string, configJSON json.RawMessage) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCfgRepo := infraRepo.NewAppConfigRepositoryImpl(tx)
		txHistoryRepo := infraRepo.NewConfigHistoryRepositoryImpl(tx)

		existing, err := txCfgRepo.FindByTenantAndEnv(ctx, tenantID, environment)
		if err != nil {
			return err
		}

		if existing == nil {
			newCfg := &domainEntity.AppConfig{
				TenantID:    tenantID,
				Environment: environment,
				ConfigData:  configJSON,
				Version:     1,
			}
			if err := txCfgRepo.Create(ctx, newCfg); err != nil {
				return err
			}

			history := &domainEntity.ConfigHistory{
				TenantID:     tenantID,
				OldConfig:    nil,
				NewConfig:    configJSON,
				Version:      1,
				ChangedBy:    "cli",
				ChangeReason: "upload",
			}
			return txHistoryRepo.Create(ctx, history)
		}

		newVersion := existing.Version + 1
		oldConfig := existing.ConfigData
		existing.ConfigData = configJSON
		existing.Version = newVersion

		if err := txCfgRepo.Update(ctx, existing); err != nil {
			return err
		}

		history := &domainEntity.ConfigHistory{
			TenantID:     tenantID,
			OldConfig:    oldConfig,
			NewConfig:    configJSON,
			Version:      newVersion,
			ChangedBy:    "cli",
			ChangeReason: "upload",
		}
		return txHistoryRepo.Create(ctx, history)
	})
}

func envOrDefault(key, def string) string {
	val := strings.TrimSpace(os.Getenv(key))
	if val == "" {
		return def
	}
	return val
}

func envIntOrDefault(key string, def int) int {
	val := strings.TrimSpace(os.Getenv(key))
	if val == "" {
		return def
	}
	var parsed int
	_, err := fmt.Sscanf(val, "%d", &parsed)
	if err != nil {
		return def
	}
	return parsed
}
