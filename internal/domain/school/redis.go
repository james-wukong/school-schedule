package school

import (
	"context"
	"time"
)

const (
	RedisUserPrefix = "school:"
	RedisUerTTL     = 2 * 24 * time.Hour
)

type RedisCache interface {
	GetByID(ctx context.Context, id int64) (*Schools, error)
	GetByCode(ctx context.Context, code string) (*Schools, error)
	Set(ctx context.Context, school *Schools) error
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, school *Schools) error
	// MapEmailToID(ctx context.Context, email string, id int64) error
	// DeleteEmailToID(ctx context.Context, email string) error
}
