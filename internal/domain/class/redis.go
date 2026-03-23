package class

import (
	"context"
	"time"
)

const (
	RedisUserPrefix = "class:"
	RedisUerTTL     = 2 * 24 * time.Hour
)

type RedisCache interface {
	GetByID(ctx context.Context, id int64) (*Classes, error)
	GetBySemesterID(ctx context.Context, semesterID int64) ([]*Classes, error)
	Set(ctx context.Context, class *Classes) error
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, class *Classes) error
}
