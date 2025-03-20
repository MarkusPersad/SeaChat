package handler

import (
	"SeaChat/internal/model"
	"SeaChat/pkg/common/request"
	"SeaChat/pkg/common/response"
	"SeaChat/pkg/constants"
	"SeaChat/pkg/entity"
	"SeaChat/pkg/exception"
	"SeaChat/pkg/utils"
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Account interface {
	GetCaptcha(ctx *fiber.Ctx) error
	Register(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
	GetUserInfo(ctx *fiber.Ctx) error
	Logout(ctx *fiber.Ctx) error
}

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
		if user.Status != constants.USER_OFFLINE {
			return exception.ErrUserStatusInvalid
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

// GetUserInfo 获取用户信息
// @Summary 获取用户信息
// @Description 获取用户信息
// @Tags Account
// @Accept json
// @Produce json
// @Param userinfo body request.UserInfo true "用户信息"
// @Success 200 {object} response.Response
// @Failure 200 {object} response.Response
// @Router /api/account/getuserinfo [post]
func(h *Handler)GetUserInfo(ctx *fiber.Ctx) error {
	tokenString,err := utils.TokenCheck(ctx,h.db,false)
	if err != nil {
		return err
	}
	var userInfo request.UserInfo
	if err := ctx.BodyParser(&userInfo); err != nil {
		return exception.ErrBadRequest
	}
	if err := utils.Validate(&userInfo); err != nil {
		log.Logger.Error().Err(err).Msgf("字段校验失败：%v", err)
		return err
	}
	var user model.User
	if err := h.db.GetDB(ctx.UserContext()).Model(&model.User{}).Where("user_id = ?",userInfo.Info).
	Or("email = ?",userInfo.Info).
	Or("user_name = ?",userInfo.Info).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return exception.ErrUserNotFound
		} else {
			return err
		}
	}
	userDetails := entity.UserDetails{
		UserID: user.UserID,
		UserName: user.UserName,
		Email: user.Email,
		Status: user.Status,
		Avatar: user.Avatar,
	}
	return ctx.Status(fiber.StatusOK).JSON(response.Success("获取用户信息成功",userDetails,tokenString))
}

// Logout 用户登出
// @Summary 用户登出
// @Description 用户登出
// @Tags Account
// @Accept json
// @Produce json
// @Param userinfo body request.UserInfo true "用户信息"
// @Success 200 {object} response.Response
// @Failure 200 {object} response.Response
// @Router /api/account/logout [post]
func(h *Handler)Logout(ctx *fiber.Ctx) error {
	_,err := utils.TokenCheck(ctx,h.db,false)
	if err != nil {
		return err
	}
	claims := ctx.Locals(constants.JWT_CONTEXT_KEY).(*jwt.Token).Claims.(*entity.SeaClaim)
	var userInfo request.UserInfo
	if err := ctx.BodyParser(&userInfo); err != nil {
		return exception.ErrBadRequest
	}
	if err := utils.Validate(&userInfo); err != nil {
		log.Logger.Error().Err(err).Msgf("字段校验失败：%v", err)
		return err
	}
	err = h.db.Transaction(ctx.UserContext(),func(ctx context.Context) error {
		var user model.User
		if err := h.db.GetDB(ctx).Model(&model.User{}).Where("user_id = ?",userInfo.Info).
		First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound{
				return exception.ErrUserNotFound
			}
			return err
		}
		if user.UserID != claims.UserID {
			return exception.ErrPermissionDenied
		}
		user.Status = constants.USER_OFFLINE
		if err := h.db.GetDB(ctx).Model(&model.User{}).Updates(&user).Error; err != nil {
			log.Error().Err(err).Msgf("用户登出失败：%v", err)
			return err
		}
		if err := h.db.DelValue(ctx,constants.JWT_CONTEXT_KEY+":"+user.UserID);err != nil {
			log.Error().Err(err).Msgf("token删除失败：%v", err)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(response.Success("登出成功",nil))
}