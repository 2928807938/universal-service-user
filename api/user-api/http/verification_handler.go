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

// VerificationHandler 验证码处理器
type VerificationHandler struct {
	verificationAppService *service.VerificationAppService
}

// NewVerificationHandler 创建验证码处理器
func NewVerificationHandler(
	verificationAppService *service.VerificationAppService,
) *VerificationHandler {
	return &VerificationHandler{
		verificationAppService: verificationAppService,
	}
}

// SendCode 发送验证码
// @Summary 发送验证码
// @Description 发送邮箱或短信验证码
// @Tags 验证码
// @Accept json
// @Produce json
// @Param request body request.SendCodeRequest true "发送验证码请求"
// @Success 200 {object} types.Response{data=vo.SendCodeResponse}
// @Router /api/v1/verification/code/send [post]
func (h *VerificationHandler) SendCode(ctx context.Context, c *app.RequestContext) {
	var req request.SendCodeRequest
	if err := c.BindAndValidate(&req); err != nil {
		errors.HandleError(ctx, c, err)
		return
	}

	resp, err := h.verificationAppService.SendCode(ctx, &req)
	if err != nil {
		errors.HandleError(ctx, c, err)
		return
	}

	c.JSON(consts.StatusOK, types.Success(resp))
}

// VerifyCode 验证验证码
// @Summary 验证验证码
// @Description 验证邮箱或短信验证码是否正确
// @Tags 验证码
// @Accept json
// @Produce json
// @Param request body request.VerifyCodeRequest true "验证验证码请求"
// @Success 200 {object} types.Response{data=vo.VerifyCodeResponse}
// @Router /api/v1/verification/code/verify [post]
func (h *VerificationHandler) VerifyCode(ctx context.Context, c *app.RequestContext) {
	var req request.VerifyCodeRequest
	if err := c.BindAndValidate(&req); err != nil {
		errors.HandleError(ctx, c, err)
		return
	}

	resp, err := h.verificationAppService.VerifyCode(ctx, &req)
	if err != nil {
		errors.HandleError(ctx, c, err)
		return
	}

	c.JSON(consts.StatusOK, types.Success(resp))
}
