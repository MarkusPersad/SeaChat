package server

import (
	errormiddleware "SeaChat/internal/middleware/ErrorMiddleware"
	jwtmiddleware "SeaChat/internal/middleware/JwtMiddleware"
	recovermiddleware "SeaChat/internal/middleware/RecoverMiddleware"
	"SeaChat/pkg/constants"
	"SeaChat/pkg/entity"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/contrib/fiberzerolog"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	recoverer "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog/log"
)

func (s *FiberServer) RegisterFiberRoutes() {
	// Apply CORS middleware
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: false, // credentials require explicit origins
		MaxAge:           300,
	}))

	// Apply logging middleware
	s.App.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &log.Logger,
	}))

	// Apply rate limiting middleware
	s.App.Use(limiter.New(limiter.Config{
		Max: constants.LIMITER_TIMES,
		Expiration: constants.LIMITER_TIME*time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))

	//Apply monitor middleware
	s.App.Get("/metrics", monitor.New(monitor.Config{Title: fmt.Sprintf("%s Metrics Page",os.Getenv("APP_NAME"))}))

	// Apply recover middleware
	s.App.Use(recoverer.New(recoverer.Config{
		Next: nil,
		EnableStackTrace: true,
		StackTraceHandler: recovermiddleware.StackTraceHandler,
	}))

	// Apply swagger middleware
	s.App.Use(swagger.New(swagger.Config{
		BasePath: "/",
		FilePath: "./docs/swagger.json",
		Title: os.Getenv("APP_NAME"),
		CacheAge: 3600,
		Path: "docs",
	}))
	// Apply jwt middleware
	s.App.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: []byte(os.Getenv("JWT_SECRET")),
			JWTAlg: "HS256",
		},
		ContextKey: constants.JWT_CONTEXT_KEY,
		Claims: &entity.SeaClaim{},
		ErrorHandler: errormiddleware.JwtErrorHandler,
		Filter: jwtmiddleware.JwtFilter,

	}))
	// Apply compress middleware
	s.App.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	s.App.Get("/", s.HelloWorldHandler)

	s.App.Get("/health", s.HealthHandler)
	api := s.App.Group("/api")
	account := api.Group("/account")
	account.Get("/getcaptcha",s.GetCaptcha)
	account.Post("/register", s.Register)
	account.Post("/login", s.Login)
	account.Post("/getuserinfo",s.GetUserInfo)
	account.Post("/logout",s.Logout)
	friend := api.Group("/friend")
	friend.Post("/add",s.AddFriend)

}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}

	return c.JSON(resp)
}
