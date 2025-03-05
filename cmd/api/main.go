package main

import (
	"SeaChat/internal/server"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog/log"
)



func gracefulShutdown(fiberServer *server.FiberServer, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Logger.Info().Msg("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := fiberServer.ShutdownWithContext(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Logger.Info().Msg("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}


// @version 1.0
// @description This is an API for SeaChat Application

// @contact.name MarkusPersad
// @contact.email msp060308@gmail.co

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /api
func main() {

	server := server.New()

	server.RegisterFiberRoutes()

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	go func() {
		port, _ := strconv.Atoi(os.Getenv("PORT"))
		err := server.Listen(fmt.Sprintf(":%d", port),fiber.ListenConfig{
			EnablePrefork: os.Getenv("PREFORK") == "true",
		})
		if err != nil {
			log.Logger.Panic().Err(err).Msgf("http server error: %s", err)
		}
	}()

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(server, done)

	// Wait for the graceful shutdown to complete
	<-done
	log.Logger.Info().Msg("Graceful shutdown complete.")
}
