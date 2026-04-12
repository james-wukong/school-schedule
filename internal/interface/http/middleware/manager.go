package middleware

import (
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Manager struct {
	log   *zerolog.Logger
	db    *gorm.DB
	redis *redis.Client
}

func NewManager(l *zerolog.Logger, db *gorm.DB, redis *redis.Client) *Manager {
	return &Manager{log: l, db: db, redis: redis}
}
