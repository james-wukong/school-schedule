package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/james-wukong/school-schedule/internal/app"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application, err := app.Bootstrap(ctx)
	if err != nil {
		log.Fatal(err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	app.Shutdown(ctx, application)
}
