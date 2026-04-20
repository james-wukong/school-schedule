package postgres

import (
	"context"
	"errors"

	tt "github.com/james-wukong/school-schedule/internal/domain/teacher_timeslot"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type teacherTimeslotRepository struct {
	db  *gorm.DB
	log *zerolog.Logger
}

func NewTeacherTimeslotRepository(db *gorm.DB, log *zerolog.Logger) *teacherTimeslotRepository {
	return &teacherTimeslotRepository{db: db, log: log}
}

// Implement the teacherTimeslotRepository interface methods here,
// using GORM to interact with the PostgreSQL database.
func (r *teacherTimeslotRepository) Create(ctx context.Context, entity *tt.TeacherTimeslots) error {
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
func (r *teacherTimeslotRepository) GetByIDs(
	ctx context.Context,
	entity *tt.TeacherTimeslots,
) (*tt.TeacherTimeslots, error) {
	var record tt.TeacherTimeslots
	err := r.db.WithContext(ctx).
		Preload("Teacher").
		Preload("Timeslot").
		First(&record, "teacher_id = ? and timeslot_id = ?",
			entity.TeacherID, entity.TimeslotID,
		).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, tt.ErrTeacherTimeslotNotFound // Return nil if not found
		}
		return nil, err
	}
	return &record, nil
}

// Implement GetByTeacherID method
func (r *teacherTimeslotRepository) GetByTeacherID(
	ctx context.Context,
	teacherID int64,
) ([]*tt.TeacherTimeslots, error) {
	var rows []*tt.TeacherTimeslots
	err := r.db.WithContext(ctx).
		Preload("Teacher").
		Preload("Timeslot").
		Find(&rows, "teacher_id = ", teacherID).
		Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// Implement GetBySubjectID method
func (r *teacherTimeslotRepository) GetByTimeslotID(
	ctx context.Context,
	slotID int64,
) ([]*tt.TeacherTimeslots, error) {
	var rows []*tt.TeacherTimeslots
	err := r.db.WithContext(ctx).
		Preload("Teacher").
		Preload("Timeslot").
		Find(&rows, "timeslot_id = ", slotID).
		Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// Implement Update method
func (r *teacherTimeslotRepository) Update(
	ctx context.Context,
	entity *tt.TeacherTimeslots,
) error {
	// Update the school record in the database
	err := r.db.WithContext(ctx).Save(entity).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error updating teacher in database")
	}

	return err
}
