package room

import (
	"context"
	"time"
)

const (
	RedisUserPrefix = "room:"
	RedisUerTTL     = 2 * 24 * time.Hour
)

type RedisCache interface {
	GetByID(ctx context.Context, id int64) (*Rooms, error)
	GetBySchoolID(ctx context.Context, schoolID int64) ([]*Rooms, error)
	Set(ctx context.Context, room *Rooms) error
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, room *Rooms) error
}
