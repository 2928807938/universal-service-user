package service

import (
	"context"
	"fmt"
	"github.com/2928807938/universal-service-user/configcenter/domain/entity"
	"github.com/2928807938/universal-service-user/configcenter/domain/repository"
	"github.com/2928807938/universal-service-user/share/errors"
	"sync"
)

// CachedConfig represents in-memory tenant config.
type CachedConfig struct {
	Config *entity.AppConfig
}

// ConfigCacheService loads configs and caches them per tenant+env.
type ConfigCacheService struct {
	repo  repository.AppConfigRepository
	cache sync.Map
}

// NewConfigCacheService creates a config cache service.
func NewConfigCacheService(repo repository.AppConfigRepository) *ConfigCacheService {
	return &ConfigCacheService{
		repo: repo,
	}
}

// GetConfig returns cached config or loads from DB.
func (s *ConfigCacheService) GetConfig(ctx context.Context, tenantID, environment string) (*CachedConfig, error) {
	if tenantID == "" || environment == "" {
		return nil, errors.ErrBadRequest("tenant_id and environment are required")
	}
	key := buildCacheKey(tenantID, environment)
	if cached, ok := s.cache.Load(key); ok {
		return cached.(*CachedConfig), nil
	}

	cfg, err := s.repo.FindByTenantAndEnv(ctx, tenantID, environment)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, errors.ErrNotFound("config not found")
	}

	cached := &CachedConfig{Config: cfg}
	s.cache.Store(key, cached)
	return cached, nil
}

// Invalidate removes cached config for a tenant/environment.
func (s *ConfigCacheService) Invalidate(tenantID, environment string) {
	if tenantID == "" || environment == "" {
		return
	}
	s.cache.Delete(buildCacheKey(tenantID, environment))
}

func buildCacheKey(tenantID, environment string) string {
	return fmt.Sprintf("%s:%s", tenantID, environment)
}
