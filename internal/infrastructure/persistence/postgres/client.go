// Package postgres provides a Postgres client for database operations.
// It initializes the client with the provided configuration and ensures a successful connection to the Postgres server. If the connection fails, it logs the error using the internal logger and returns the error for further handling.
package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/james-wukong/school-schedule/internal/config"
	"github.com/james-wukong/school-schedule/internal/infrastructure/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Initialize logger with console output only
var conLog = logger.New(logger.LogConfig{
	EnableConsole: true,
	FilePath:      "",
	// MaxSize:       5,    // Rotate every 5MB
	// MaxBackups:    10,   // Keep last 10 files
	// Compress:      false, // Save disk space
})

func dsnFromConfig(cfg config.PostgresConfig) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DB,
		cfg.SSL,
	)
}

func NewPool(ctx context.Context, cfg config.PostgresConfig) (*pgxpool.Pool, error) {
	dsn := dsnFromConfig(cfg)

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	poolCfg.MaxConns = 20
	poolCfg.MinConns = 2
	poolCfg.MaxConnLifetime = time.Hour

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		conLog.Error().Err(err).Msg("Failed to connect to Postgres")
		return nil, err
	}

	return pool, nil
}

func NewGormDB(ctx context.Context, cfg config.PostgresConfig) (*gorm.DB, error) {
	dsn := dsnFromConfig(cfg)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		conLog.Error().Err(err).Msg("GORM Failed to connect to Postgres")
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		conLog.Error().Err(err).Msg("GORM Failed to return sql.DB")
		return nil, err
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
