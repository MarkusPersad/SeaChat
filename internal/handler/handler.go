package handler

import (
	"SeaChat/internal/database"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	db database.Service
}

type HandlerInterface interface {
	HealthHandler(c *fiber.Ctx) error
	InitDB(models ...any) error
	Account
	Friend
}

func New() HandlerInterface {
	return &Handler{
		db: database.New(),
	}
}

func (h *Handler) InitDB(models ...any) error {
	return h.db.InitDB(models...)
}

// healthHandler 获取健康状态
// @Summary 获取健康状态
// @Description 获取健康状态
// @Tags 服务管理
// @Accept  json
// @Produce  json
// @Success 200 {object} any "返回结果"
// @Failure 200 {object} any "返回结果"
// @Router /health [get]
func (h *Handler) HealthHandler(c *fiber.Ctx) error {
	return c.JSON(h.db.Health())
}
