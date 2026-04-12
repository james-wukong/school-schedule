// Package app contains the application logic for the orders API, including initialization and shutdown procedures.
package application

import (
	"context"

	"github.com/rs/zerolog"
)

func Shutdown(ctx context.Context, app *App, log *zerolog.Logger) error {
	log.Info().Msg("Shutting down application...")

	// shutdown HTTP server with context timeout
	if err := app.HTTPServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Error shutting down HTTP server")
		return err
	}
	// Close database and Redis connections if they exist
	sqlDB, _ := app.Database.DB.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}

	// Close Redis connection if it exists
	if app.Redis != nil {
		_ = app.Redis.Close()
	}

	log.Info().Msg("Shutting down complete...")
	return nil
}
