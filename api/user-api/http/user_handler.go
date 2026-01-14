package http

import (
	"context"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"universal-service-user/api/user-api/dto/request"
	"universal-service-user/api/user-api/service"
	"universal-service-user/share/errors"
	"universal-service-user/share/types"
)

// UserHandler 用户处理器(用户管理相关)
type UserHandler struct {
	userAppService *service.UserAppService
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userAppService *service.UserAppService) *UserHandler {
	return &UserHandler{
		userAppService: userAppService,
	}
}

// Register 用户注册
// @Summary 用户注册
// @Description 通过邮箱或手机号注册用户
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body request.RegisterRequest true "注册请求"
// @Success 200 {object} types.Response{data=vo.UserVo}
// @Router /api/v1/users/register [post]
func (h *UserHandler) Register(ctx context.Context, c *app.RequestContext) {
	var req request.RegisterRequest
	if err := c.BindAndValidate(&req); err != nil {
		errors.HandleError(ctx, c, err)
		return
	}

	resp, err := h.userAppService.Register(ctx, &req)
	if err != nil {
		errors.HandleError(ctx, c, err)
		return
	}

	c.JSON(consts.StatusOK, types.Success(resp))
}

// GetUser 获取用户
// @Summary 获取用户信息
// @Description 根据用户ID获取用户信息
// @Tags 用户
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} types.Response{data=vo.UserVo}
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(ctx context.Context, c *app.RequestContext) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errors.HandleError(ctx, c, err)
		return
	}

	user, err := h.userAppService.GetUser(ctx, id)
	if err != nil {
		errors.HandleError(ctx, c, err)
		return
	}

	c.JSON(consts.StatusOK, types.Success(user))
}

// UpdateUser 更新用户
// @Summary 更新用户信息
// @Description 更新用户信息
// @Tags 用户
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param request body request.UpdateUserRequest true "更新用户请求"
// @Success 200 {object} types.Response{data=vo.UserVo}
// @Router /api/v1/users/{id} [put]
func (h *UserHandler) UpdateUser(ctx context.Context, c *app.RequestContext) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errors.HandleError(ctx, c, err)
		return
	}

	var req request.UpdateUserRequest
	if err := c.BindAndValidate(&req); err != nil {
		errors.HandleError(ctx, c, err)
		return
	}

	user, err := h.userAppService.UpdateUser(ctx, id, &req)
	if err != nil {
		errors.HandleError(ctx, c, err)
		return
	}

	c.JSON(consts.StatusOK, types.Success(user))
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 已登录用户修改密码
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body request.ChangePasswordRequest true "修改密码请求"
// @Success 200 {object} types.Response
// @Router /api/v1/users/password/change [post]
func (h *UserHandler) ChangePassword(ctx context.Context, c *app.RequestContext) {
	var req request.ChangePasswordRequest
	if err := c.BindAndValidate(&req); err != nil {
		errors.HandleError(ctx, c, err)
		return
	}

	// TODO: 从 token 中获取当前用户ID
	// userID := c.GetInt("user_id")

	c.JSON(consts.StatusOK, types.Success(nil))
}

// ResetPassword 重置密码(忘记密码)
// @Summary 重置密码(忘记密码)
// @Description 通过邮箱或手机号验证码重置密码
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body request.ResetPasswordRequest true "重置密码请求"
// @Success 200 {object} types.Response
// @Router /api/v1/users/password/reset [post]
func (h *UserHandler) ResetPassword(ctx context.Context, c *app.RequestContext) {
	var req request.ResetPasswordRequest
	if err := c.BindAndValidate(&req); err != nil {
		errors.HandleError(ctx, c, err)
		return
	}

	err := h.userAppService.ResetPassword(ctx, &req)
	if err != nil {
		errors.HandleError(ctx, c, err)
		return
	}

	c.JSON(consts.StatusOK, types.Success(nil))
}
