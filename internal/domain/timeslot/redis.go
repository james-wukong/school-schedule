package timeslot

import (
	"context"
	"time"
)

const (
	RedisUserPrefix = "timeslot:"
	RedisUerTTL     = 2 * 24 * time.Hour
)

type RedisCache interface {
	GetByID(ctx context.Context, id int64) (*Timeslots, error)
	GetBySemesterID(ctx context.Context, semesterID int64) ([]*Timeslots, error)
	Set(ctx context.Context, timeslot *Timeslots) error
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, timeslot *Timeslots) error
}
