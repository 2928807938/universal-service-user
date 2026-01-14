package service

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/2928807938/universal-service-user/api/config-api/dto/request"
	"github.com/2928807938/universal-service-user/api/config-api/dto/vo"
	"github.com/2928807938/universal-service-user/configcenter/domain/entity"
	"github.com/2928807938/universal-service-user/configcenter/domain/repository"
	"github.com/2928807938/universal-service-user/share/errors"
)

// AppService handles app registration.
type AppService struct {
	appRepo repository.AppRepository
}

func NewAppService(appRepo repository.AppRepository) *AppService {
	return &AppService{appRepo: appRepo}
}

func (s *AppService) Register(ctx context.Context, req *request.RegisterAppRequest) (*vo.AppRegisterVo, error) {
	if req == nil {
		return nil, errors.ErrBadRequest("request is required")
	}
	appName := strings.TrimSpace(req.AppName)
	if appName == "" {
		return nil, errors.ErrBadRequest("app_name is required")
	}

	tenantID := uuid.New().String()
	exists, err := s.appRepo.ExistsByTenantID(ctx, tenantID)
	if err != nil {
		return nil, errors.ErrInternal("failed to check tenant_id", err)
	}
	for exists {
		tenantID = uuid.New().String()
		exists, err = s.appRepo.ExistsByTenantID(ctx, tenantID)
		if err != nil {
			return nil, errors.ErrInternal("failed to check tenant_id", err)
		}
	}

	app := &entity.App{
		TenantID:    tenantID,
		AppName:     appName,
		Description: strings.TrimSpace(req.Description),
		Status:      "active",
	}
	if err := s.appRepo.Create(ctx, app); err != nil {
		return nil, errors.ErrInternal("failed to create app", err)
	}

	return &vo.AppRegisterVo{TenantID: tenantID}, nil
}
