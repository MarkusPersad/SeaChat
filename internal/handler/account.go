package handler

import (
	"SeaChat/internal/model"
	"SeaChat/pkg/common/request"
	"SeaChat/pkg/common/response"
	"SeaChat/pkg/exception"
	"SeaChat/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
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
		return ctx.Status(fiber.StatusOK).JSON(response.Success("验证码获取成功", captchaData))
	}
}

func(h *Handler)Register(ctx *fiber.Ctx) error {
	var register request.Register
	if err := ctx.BodyParser(&register); err != nil {
		return exception.ErrBadRequest
	}
	if err := utils.Validate(&register); err != nil {
		log.Logger.Error().Err(err).Msgf("字段校验失败：%v", err)
		return err
	}
	if !utils.VerifyCaptcha(h.db, register.CheckCodeKey, register.CheckCode){
		return exception.ErrCaptchaInvalid
	}
	var user model.User
	if err := h.db.GetDB(ctx.UserContext()).Model(&model.User{}).Where("user_name = ?", register.UserName).Or("email = ?", register.Email).First(&user).Error; err == nil {
		return exception.ErrUserAlreadyExists
	}
	user.UserID = uuid.New().String()
	user.Email = register.Email
	user.UserName = register.UserName
	if password,err := utils.GeneratePassword(register.Password); err != nil {
		log.Error().Err(err).Msgf("密码生成失败：%v", err)
		return err
	} else {
		user.Password = password
		if err := h.db.GetDB(ctx.UserContext()).Create(&user).Error; err != nil {
			log.Error().Err(err).Msgf("用户注册失败：%v", err)
			return err
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(response.Success("注册成功", nil))
}
