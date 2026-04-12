// Package main is the entry point for the orders API application.
// It initializes the application, sets up signal handling for graceful shutdown, and starts the application.
package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	app "github.com/james-wukong/school-schedule/internal/application"
	"github.com/james-wukong/school-schedule/internal/config"
	"github.com/james-wukong/school-schedule/internal/infrastructure/logger"

	// Import the generated docs package (the path may vary based on your project structure)
	_ "github.com/james-wukong/school-schedule/docs"
)

//	@title						Your API Title
//	@version					1.0
//	@description				link: http://127.0.0.1:8092/swagger/index.html
//	@host						localhost:8092
//	@BasePath					/api/v1
//	@schemes					http
//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/
//
// cmd: swag init --parseDependency --parseInternal -d cmd/api,internal/interface/http/handler
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	appLog := logger.New(logger.LogConfig{
		EnableConsole: true,
		FilePath:      "logs/app.log",
		MaxSize:       5,    // Rotate every 5MB
		MaxBackups:    10,   // Keep last 10 files
		Compress:      true, // Save disk space
	})

	// 1. Load Configuration (Viper)
	cfg := config.InitConfig()

	// 2. Initialize the App Container (Dependency Injection)
	// This sets up HTTPServer, DB, Redis, Loggers, routes, and Usecases
	application, err := app.Bootstrap(ctx, &appLog)
	if err != nil {
		appLog.Error().Err(err).Msg("Failed to bootstrap application")
	}

	// 3. Start Server in a separate Goroutine
	go func() {
		appLog.Info().Msgf("Server starting on port %d", cfg.App.Port)
		if err := application.HTTPServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLog.Error().Err(err).Msg("Failed to listen and serve")
		}
	}()

	// 5. Wait for Interrupt Signal (Graceful Shutdown)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Block here until a signal is received

	appLog.Info().Msg("Shutting down server...")

	// 6. Graceful Shutdown context (5 second timeout)
	if err := app.Shutdown(ctx, application, &appLog); err != nil {
		appLog.Error().Err(err).Msg("Server forced to shutdown")
	}

	appLog.Info().Msg("Server exiting")
}
