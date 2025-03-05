package recoverery

import (
	"runtime/debug"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

func StackTraceHandler(_ fiber.Ctx, e any){
	log.Logger.Panic().Msgf("panic: %v\n%s\n", e, debug.Stack())
}
