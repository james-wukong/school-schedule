package app

import (
	"context"
	"log"
)

func Shutdown(ctx context.Context, app *App) {
	log.Println("Shutting down...")

	if app.DB != nil {
		app.DB.Pool.Close()
	}

	if app.Redis != nil {
		_ = app.Redis.Close()
	}

	log.Println("Shutdown complete")
}
