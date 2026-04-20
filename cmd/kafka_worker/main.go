package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/james-wukong/school-schedule/internal/config"
	"github.com/james-wukong/school-schedule/internal/infrastructure/logger"
	infraPostgre "github.com/james-wukong/school-schedule/internal/infrastructure/persistence/postgres"
	kafkautils "github.com/james-wukong/school-schedule/internal/utils"
)

func main() {
	// ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	workerLog := logger.New(logger.LogConfig{
		EnableConsole: true,
		FilePath:      "/app/logs/worker.log",
		MaxSize:       5,    // Rotate every 5MB
		MaxBackups:    10,   // Keep last 10 files
		Compress:      true, // Save disk space
	})

	// 1. Load Configurations and create *gorm.DB
	cfg := config.InitConfig()
	// Initialize Postgres connection pool
	db, err := infraPostgre.NewGormDB(ctx, cfg.Database.Postgres)
	if err != nil {
		workerLog.Error().Err(err).Msg("Failed to return *gorm.DB")
	}

	// 2. Start the consumer loop
	// This function never returns; it blocks while waiting for work
	kafkautils.StartScheduleConsumer(ctx, db, cfg, &workerLog)
}
