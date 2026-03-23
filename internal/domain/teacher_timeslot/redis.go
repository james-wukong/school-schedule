package teachertimeslot

import (
	"context"
	"time"
)

const (
	RedisUserPrefix = "teacher_timeslot:"
	RedisUerTTL     = 2 * 24 * time.Hour
)

type RedisCache interface {
	GetByIDs(ctx context.Context, entity *TeacherTimeslots) (*TeacherTimeslots, error)
	GetByTeacherID(ctx context.Context, teacherID int64) ([]*TeacherTimeslots, error)
	GetByTimeslotID(ctx context.Context, slotID int64) ([]*TeacherTimeslots, error)
	Set(ctx context.Context, entity *TeacherTimeslots) error
	Update(ctx context.Context, entity *TeacherTimeslots) error
}
