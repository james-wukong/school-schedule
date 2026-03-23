package schedule

import (
	"context"
	"time"
)

const (
	RedisUserPrefix = "schedule:"
	RedisUerTTL     = 2 * 24 * time.Hour
)

type RedisCache interface {
	GetByID(ctx context.Context, id int64) (*Schedules, error)
	GetBySchoolID(ctx context.Context, schoolID int64) ([]*Schedules, error)
	GetByRequirementID(ctx context.Context, requirementID int64) ([]*Schedules, error)
	GetByVersion(ctx context.Context, schoolID int64, version float64) ([]*Schedules, error)
	Set(ctx context.Context, schedule *Schedules) error
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, schedule *Schedules) error
}
