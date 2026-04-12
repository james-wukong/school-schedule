package requirement

import (
	"context"
	"time"
)

const (
	RedisUserPrefix = "requirement:"
	RedisUerTTL     = 2 * 24 * time.Hour
)

type RedisCache interface {
	GetByVersion(ctx context.Context,
		schoolID, semesterID int64,
		version float64,
	) ([]*Requirements, error)
}
