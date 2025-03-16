package handler

import (
	"SeaChat/pkg/common/response"
	"SeaChat/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// GetCaptcha 获取验证码
// @Summary 获取验证码
// @Description 获取验证码
// @Tags Account
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 200 {object} response.Response
// @Router /api/account/getcaptcha [get]
func(h *Handler) GetCaptcha(ctx *fiber.Ctx) error {
	if captchaData,err := utils.GenerateCaptcha(h.db); err != nil {
		return err
	} else {
		return ctx.Status(fiber.StatusOK).JSON(response.Success(fiber.StatusOK,  "验证码获取成功", captchaData))
	}
}
