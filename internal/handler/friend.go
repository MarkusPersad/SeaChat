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
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Friend interface {
	AddFriend(ctx *fiber.Ctx) error
	AcceptFriend(ctx *fiber.Ctx) error
}

// AddFriend 添加好友(发送请求)
// @Summary 添加好友(发送请求)
// @Description 添加好友(发送请求)
// @Tags Friend
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param userInfo body request.UserInfo true "用户信息"
// @Success 200 {object} response.Response 
// @Failure 200 {object} response.Response 
// @Router /friend/add [post]
func(h *Handler) AddFriend(ctx *fiber.Ctx) error {
	tokenString,userId,err := utils.TokenCheck(ctx,h.db,false)
	if err != nil {
		return err
	}
	var userInfo request.UserInfo
	if err := ctx.BodyParser(&userInfo); err != nil {
		return exception.ErrBadRequest
	}
	if err := utils.Validate(userInfo); err != nil {
		return err
	}
	err = h.db.Transaction(ctx.UserContext(),func(ctx context.Context) error {
	var user model.User
	if err := h.db.GetDB(ctx).Model(&model.User{}).
	Where("user_id = ?",userInfo.Info).
	Or("email = ?",userInfo.Info).
	Or("user_name = ?",userInfo.Info).
	First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return exception.ErrUserNotFound
		}
		return err
	}
	friend := model.Friend{
		FGID: utils.GetSortedID("",userId,user.UserID),
		UserID: user.UserID,
		FriendID: userId,
		Status: constants.FRIEND_STATUS_WAITING,
	}
	if err := h.db.GetDB(ctx).Create(&friend).Error; err != nil {
		return err
	}
		return nil
	})
	if err != nil {
		log.Logger.Error().Err(err).Msgf("add friend failed:%v",err)
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(response.Success("发送好友请求成功",nil,tokenString))
}
// AcceptFriend 接受好友(同意请求)
// @Summary 接受好友(同意请求)
// @Description 接受好友(同意请求)
// @Tags Friend
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param userInfo body request.UserInfo true "用户信息"
// @Success 200 {object} response.Response 
// @Failure 200 {object} response.Response 
// @Router /friend/accept [post]
func(h *Handler) AcceptFriend(ctx *fiber.Ctx) error {
	tokenString,userId,err := utils.TokenCheck(ctx,h.db,false)
	if err != nil {
		return err
	}
	var userInfo request.UserInfo
	if err := ctx.BodyParser(&userInfo); err != nil {
		return exception.ErrBadRequest
	}
	if err := utils.Validate(userInfo); err != nil {
		return err
	}
	err = h.db.Transaction(ctx.UserContext(),func(ctx context.Context) error {
		if err := h.db.GetDB(ctx).Model(&model.Friend{}).Where("fg_id = ?",utils.GetSortedID("",userId,userInfo.Info)).
		Updates(&model.Friend{Status: constants.FRIEND_STATUS_FRIENDLY}).Error; err != nil {
			log.Logger.Error().Err(err).Msgf("accept friend failed:%v",err)
			return err
		}
		return nil
	})
	if err != nil {
		log.Logger.Error().Err(err).Msgf("accept friend failed:%v",err)
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(response.Success("接受好友请求成功",nil,tokenString))
}