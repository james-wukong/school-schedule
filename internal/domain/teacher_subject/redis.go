package teachersubject

import (
	"context"
	"time"
)

const (
	RedisUserPrefix = "teacher_subject:"
	RedisUerTTL     = 2 * 24 * time.Hour
)

type RedisCache interface {
	GetByIDs(ctx context.Context, entity *TeacherSubjects) (*TeacherSubjects, error)
	GetByTeacherID(ctx context.Context, teacherID int64) ([]*TeacherSubjects, error)
	GetBySubjectID(ctx context.Context, subjectID int64) ([]*TeacherSubjects, error)
	Set(ctx context.Context, entity *TeacherSubjects) error
	Update(ctx context.Context, entity *TeacherSubjects) error
}
