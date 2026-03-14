// Package redis provides a Redis client for caching and other operations.
// It initializes the client with the provided configuration and ensures a successful connection to the Redis server. If the connection fails, it logs the error using the internal logger and returns the error for further handling.
package redis

import (
	"context"
	"fmt"

	"github.com/james-wukong/school-schedule/internal/config"
	"github.com/james-wukong/school-schedule/internal/infrastructure/logger"
	"github.com/redis/go-redis/v9"
)

func New(cfg config.RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		// Initialize logger with console output only
		conLog := logger.New(logger.LogConfig{
			EnableConsole: true,
			FilePath:      "",
			// MaxSize:       5,    // Rotate every 5MB
			// MaxBackups:    10,   // Keep last 10 files
			// Compress:      false, // Save disk space
		})
		conLog.Error().Err(err).Msg("Failed to connect to Redis")
		return nil, err
	}

	return rdb, nil
}
