package semester

import (
	"context"
	"time"
)

const (
	RedisUserPrefix = "semester:"
	RedisUerTTL     = 2 * 24 * time.Hour
)

type RedisCache interface {
	GetByID(ctx context.Context, id int64) (*Semesters, error)
	GetBySchoolID(ctx context.Context, schoolID int64) ([]*Semesters, error)
	Set(ctx context.Context, semester *Semesters) error
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, semester *Semesters) error
}
