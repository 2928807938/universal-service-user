package http

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/2928807938/universal-service-user/api/config-api/dto/request"
	"github.com/2928807938/universal-service-user/api/config-api/service"
	"github.com/2928807938/universal-service-user/share/errors"
	"github.com/2928807938/universal-service-user/share/types"
)

// AppHandler handles app management.
type AppHandler struct {
	appService *service.AppService
}

func NewAppHandler(appService *service.AppService) *AppHandler {
	return &AppHandler{appService: appService}
}

// Register registers a tenant app.
func (h *AppHandler) Register(ctx context.Context, c *app.RequestContext) {
	var req request.RegisterAppRequest
	if err := c.BindAndValidate(&req); err != nil {
		errors.HandleError(ctx, c, err)
		return
	}

	resp, err := h.appService.Register(ctx, &req)
	if err != nil {
		errors.HandleError(ctx, c, err)
		return
	}

	c.JSON(consts.StatusOK, types.Success(resp))
}
