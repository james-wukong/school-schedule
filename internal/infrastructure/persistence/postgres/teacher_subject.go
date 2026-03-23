package postgres

import (
	"context"
	"errors"

	ts "github.com/james-wukong/school-schedule/internal/domain/teacher_subject"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type teacherSubjectRepository struct {
	db  *gorm.DB
	log *zerolog.Logger
}

func NewTeacherSubjectRepository(db *gorm.DB, log *zerolog.Logger) ts.Repository {
	return &teacherSubjectRepository{db: db, log: log}
}

// Implement the teacherSubjectRepository interface methods here,
// using GORM to interact with the PostgreSQL database.
func (r *teacherSubjectRepository) Create(ctx context.Context, entity *ts.TeacherSubjects) error {
	// First, create the school in the database
	err := r.db.WithContext(ctx).
		Create(entity).
		Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error creating a teacher subject in database")
		return err
	}
	return nil
}

// Implement GetByID method
func (r *teacherSubjectRepository) GetByIDs(
	ctx context.Context,
	entity *ts.TeacherSubjects,
) (*ts.TeacherSubjects, error) {
	var record ts.TeacherSubjects
	err := r.db.WithContext(ctx).
		Preload("Teacher").
		Preload("Subject").
		First(&record, "teacher_id = ? and subject_id = ?",
			entity.TeacherID, entity.SubjectID,
		).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ts.ErrTeacherSubjectNotFound // Return nil if not found
		}
		return nil, err
	}
	return &record, nil
}

// Implement GetByTeacherID method
func (r *teacherSubjectRepository) GetByTeacherID(
	ctx context.Context,
	teacherID int64,
) ([]*ts.TeacherSubjects, error) {
	var rows []*ts.TeacherSubjects
	err := r.db.WithContext(ctx).
		Preload("Teacher").
		Preload("Subject").
		Find(&rows, "teacher_id = ", teacherID).
		Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// Implement GetBySubjectID method
func (r *teacherSubjectRepository) GetBySubjectID(
	ctx context.Context,
	subjectID int64,
) ([]*ts.TeacherSubjects, error) {
	var rows []*ts.TeacherSubjects
	err := r.db.WithContext(ctx).
		Preload("Teacher").
		Preload("Subject").
		Find(&rows, "subject_id = ", subjectID).
		Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// Implement Update method
func (r *teacherSubjectRepository) Update(
	ctx context.Context,
	entity *ts.TeacherSubjects,
) error {
	// Update the school record in the database
	err := r.db.WithContext(ctx).Save(entity).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error updating teacher in database")
	}

	return err
}
