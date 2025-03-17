package handler

import (
	"SeaChat/internal/model"
	"SeaChat/pkg/common/request"
	"SeaChat/pkg/common/response"
	"SeaChat/pkg/constants"
	"SeaChat/pkg/exception"
	"SeaChat/pkg/utils"
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
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

// Register 用户注册
// @Summary 用户注册
// @Description 用户注册
// @Tags Account
// @Accept json
// @Produce json
// @Param register body request.Register true "用户注册信息"
// @Success 200 {object} response.Response
// @Failure 200 {object} response.Response
// @Router /api/account/register [post]
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
	if err := h.db.GetDB(ctx.UserContext()).Model(&model.User{}).Select("id").Where("user_name = ?", register.UserName).Or("email = ?", register.Email).First(&user).Error; err == nil {
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

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录
// @Tags Account
// @Accept json
// @Produce json
// @Param login body request.Login true "用户登录"
// @Success 200 {object} response.Response
// @Failure 200 {object} response.Response
// @Router /api/account/login [post]
func(h *Handler) Login(ctx *fiber.Ctx) error {
	var login request.Login
	if err := ctx.BodyParser(&login);err != nil {
		return exception.ErrBadRequest
	}
	if err := utils.Validate(&login);err != nil {
		log.Logger.Error().Err(err).Msgf("字段校验失败：%v", err)
		return err
	}
	if !utils.VerifyCaptcha(h.db,login.CheckCodeKey,login.CheckCode){
		return exception.ErrCaptchaInvalid
	}
	var user model.User
	var tokenString string
	err := h.db.Transaction(ctx.UserContext(),func(ctx context.Context) error{
		if err := h.db.GetDB(ctx).Model(&model.User{}).Where("email = ?",login.Email).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return exception.ErrUserNotFound
			}
			log.Error().Err(err).Msgf("用户登录失败：%v", err)
			return err
		}
		if err := utils.CompareHashPassword(user.Password, login.Password); err != nil {
			log.Logger.Error().Err(err).Msgf("密码校验失败：%v", err)
			return exception.ErrPasswordInvalid
		}
		if err := h.db.GetDB(ctx).Model(&model.User{}).Where("user_id",user.UserID).Update("status",constants.USER_ONLINE).Error; err != nil {
			log.Error().Err(err).Msgf("用户登录失败：%v", err)
			return err
		}
		if token,err := utils.GenerateTokenString(user.UserID); err != nil {
			log.Error().Err(err).Msgf("token生成失败：%v", err)
			return err
		} else {
			if err = h.db.SetAndTime(ctx,constants.JWT_CONTEXT_KEY+":"+user.UserID,token,constants.TOKEN_EXPIRATION*60); err != nil {
				log.Error().Err(err).Msgf("token存储失败：%v", err)
				return err
			}
			tokenString = token
		}
		return nil
	})
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(response.Success("登录成功", tokenString))
}
