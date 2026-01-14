package http

import (
	"context"

	"github.com/2928807938/universal-service-user/api/user-api/dto/request"
	"github.com/2928807938/universal-service-user/api/user-api/service"
	"github.com/2928807938/universal-service-user/share/errors"
	"github.com/2928807938/universal-service-user/share/types"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// AuthHandler 认证处理器(只负责登录相关)
type AuthHandler struct {
	authAppService *service.AuthAppService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authAppService *service.AuthAppService) *AuthHandler {
	return &AuthHandler{
		authAppService: authAppService,
	}
}

// Login 用户登录
// @Summary 用户登录
// @Description 支持用户名+密码、邮箱+密码、手机号+验证码登录
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body request.LoginRequest true "登录请求"
// @Success 200 {object} types.Response{data=vo.LoginResponse}
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(ctx context.Context, c *app.RequestContext) {
	var req request.LoginRequest
	if err := c.BindAndValidate(&req); err != nil {
		errors.HandleError(ctx, c, err)
		return
	}

	resp, err := h.authAppService.Login(ctx, &req)
	if err != nil {
		errors.HandleError(ctx, c, err)
		return
	}

	c.JSON(consts.StatusOK, types.Success(resp))
}

// Logout 用户登出
// @Summary 用户登出
// @Description 登出当前或指定设备
// @Tags 认证
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {access_token}"
// @Param request body request.LogoutRequest true "登出请求"
// @Success 200 {object} types.Response
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(ctx context.Context, c *app.RequestContext) {
	var req request.LogoutRequest
	if err := c.BindAndValidate(&req); err != nil {
		errors.HandleError(ctx, c, err)
		return
	}

	// 从 Authorization header 中获取 accessToken
	authHeader := string(c.GetHeader("Authorization"))
	accessToken := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		accessToken = authHeader[7:]
	}

	// 调用应用服务登出
	if err := h.authAppService.Logout(ctx, accessToken, req.RefreshToken); err != nil {
		errors.HandleError(ctx, c, err)
		return
	}

	c.JSON(consts.StatusOK, types.Success(nil))
}

// RefreshToken 刷新令牌
// @Summary 刷新令牌
// @Description 使用 refresh_token 获取新的 access_token
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body request.RefreshTokenRequest true "刷新令牌请求"
// @Success 200 {object} types.Response{data=vo.RefreshTokenResponse}
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(ctx context.Context, c *app.RequestContext) {
	var req request.RefreshTokenRequest
	if err := c.BindAndValidate(&req); err != nil {
		errors.HandleError(ctx, c, err)
		return
	}

	// 调用应用服务刷新令牌
	resp, err := h.authAppService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		errors.HandleError(ctx, c, err)
		return
	}

	c.JSON(consts.StatusOK, types.Success(resp))
}
