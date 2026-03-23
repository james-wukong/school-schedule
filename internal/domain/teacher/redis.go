package teacher

import (
	"context"
	"time"
)

const (
	RedisUserPrefix = "teacher:"
	RedisUerTTL     = 2 * 24 * time.Hour
)

type RedisCache interface {
	GetByID(ctx context.Context, id int64) (*Teachers, error)
	GetBySchoolID(ctx context.Context, schoolID int64) ([]*Teachers, error)
	GetByEmployeeID(ctx context.Context, employeeID int64) ([]*Teachers, error)
	GetByName(ctx context.Context, firstName, lastName string) ([]*Teachers, error)
	Set(ctx context.Context, teacher *Teachers) error
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, teacher *Teachers) error
}
