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
	GetByID(ctx context.Context, id int64) (*Requirements, error)
	GetBySchoolID(ctx context.Context, schoolID int64) ([]*Requirements, error)
	GetBySubjectID(ctx context.Context, subjectID int64) ([]*Requirements, error)
	GetByTeacherID(ctx context.Context, teacherID int64) ([]*Requirements, error)
	GetByClassID(ctx context.Context, classID int64) ([]*Requirements, error)
	GetByVersion(ctx context.Context, schoolID int64, version float64) ([]*Requirements, error)
	Set(ctx context.Context, requirement *Requirements) error
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, requirement *Requirements) error
}
