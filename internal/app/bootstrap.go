package app

import (
	"context"
	"log"

	"github.com/james-wukong/online-orders/internal/config"
	"github.com/james-wukong/online-orders/internal/infrastructure/postgres"
	"github.com/james-wukong/online-orders/internal/infrastructure/redis"
)

type App struct {
	DB    *postgres.DBWrapper
	Redis *redis.Client
}

type DBWrapper struct {
	Pool *pgxpool.Pool
}

func Bootstrap(ctx context.Context) (*App, error) {
	cfg := config.Load()

	pgPool, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		return nil, err
	}

	redisClient, err := redis.New(cfg.Redis)
	if err != nil {
		return nil, err
	}

	log.Println("Postgres & Redis connected")

	return &App{
		DB: &DBWrapper{Pool: pgPool},
		Redis: redisClient,
	}, nil
}
