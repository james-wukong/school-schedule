package subject

import (
	"context"
	"time"
)

const (
	RedisUserPrefix = "subject:"
	RedisUerTTL     = 2 * 24 * time.Hour
)

type RedisCache interface {
	GetByID(ctx context.Context, id int64) (*Subjects, error)
	GetBySchoolID(ctx context.Context, schoolID int64) ([]*Subjects, error)
	GetByCode(ctx context.Context, code string) (*Subjects, error)
	Set(ctx context.Context, subject *Subjects) error
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, subject *Subjects) error
}
