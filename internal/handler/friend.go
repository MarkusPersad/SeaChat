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
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Friend interface {
	AddFriend(ctx *fiber.Ctx) error
}

// AddFriend 添加好友(发送请求)
// @Summary 添加好友(发送请求)
// @Description 添加好友(发送请求)
// @Tags Friend
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param userInfo body request.UserInfo true "用户信息"
// @Success 200 {object} response.Response true
// @Failure 200 {object} response.Response true
// @Router /friend/add [post]
func(h *Handler) AddFriend(ctx *fiber.Ctx) error {
	tokenString,err := utils.TokenCheck(ctx,h.db,false)
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
	userId := ctx.Locals(constants.JWT_CONTEXT_KEY).(*jwt.Token).Claims.(*entity.SeaClaim).UserID
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
		UserID: user.UserID,
		FriendID: userId,
		Status: constants.FRIEND_STATUS_WAITING,
	}
	if err := h.db.GetDB(ctx).Create(&friend).Error; err != nil {
		return err
	}
	friend.FriendID = user.UserID
	friend.UserID = userId
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