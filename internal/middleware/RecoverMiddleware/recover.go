package recovermiddleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)


func StackTraceHandler(_ *fiber.Ctx,e any){
	log.Error().Msgf("Recovered error: %v", e)
}
