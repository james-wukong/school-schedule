package postgres

import (
	"context"
	"errors"

	"github.com/james-wukong/school-schedule/internal/domain/timeslot"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type timeslotRepository struct {
	db  *gorm.DB
	log *zerolog.Logger
}

func NewTimeslotRepository(db *gorm.DB, log *zerolog.Logger) timeslot.Repository {
	return &timeslotRepository{db: db, log: log}
}

// Implement the schoolRepository interface methods here,
// using GORM to interact with the PostgreSQL database.
func (r *timeslotRepository) Create(ctx context.Context, entity *timeslot.Timeslots) error {
	// First, create the school in the database
	err := r.db.WithContext(ctx).Create(entity).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error creating timeslot in database")
		return err
	}
	return nil
}

// Implement GetByID method
func (r *timeslotRepository) GetByID(ctx context.Context, id int64) (*timeslot.Timeslots, error) {
	var record timeslot.Timeslots
	err := r.db.WithContext(ctx).
		Preload("Semester.School").
		First(&record, "id = ?", id).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, timeslot.ErrTimeslotNotFound // Return nil if not found
		}
		return nil, err
	}
	return &record, nil
}

// Implement GetBySchoolID method
func (r *timeslotRepository) GetBySemesterID(
	ctx context.Context,
	semesterID int64,
) ([]*timeslot.Timeslots, error) {
	var rows []*timeslot.Timeslots
	err := r.db.WithContext(ctx).
		Preload("Semester.School").
		Find(&rows, "semester_id = ?", semesterID).
		Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// Implement Update method
func (r *timeslotRepository) Update(ctx context.Context, entity *timeslot.Timeslots) error {
	// Update the school record in the database
	err := r.db.WithContext(ctx).Save(entity).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error updating timeslot in database")
	}

	return err
}

// Implement Delete method
func (r *timeslotRepository) Delete(ctx context.Context, id int64) error {
	// Delete the school record from the database
	err := r.db.WithContext(ctx).Delete(&timeslot.Timeslots{}, "id = ?", id).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error deleting a timeslot from database")
	}

	return err
}

// Implement List method
func (r *timeslotRepository) List(
	ctx context.Context,
	filter *timeslot.TimeslotFilterEntity,
) ([]*timeslot.Timeslots, error) {
	var rows []*timeslot.Timeslots
	query := r.db.WithContext(ctx).Model(&timeslot.Timeslots{})
	if filter != nil {
		if filter.SemesterID != nil {
			query = query.Where("semester_id = ", filter.SemesterID)
		}
		if filter.DayOfWeek != nil {
			query = query.Where("day_of_week = ?", filter.DayOfWeek)
		}
		if filter.StartTime != nil {
			query = query.Where("start_time = ?", filter.StartTime)
		}
	}

	err := query.Limit(filter.Limit).
		Offset((filter.Page - 1) * filter.Limit).
		Find(&rows).
		Error
	return rows, err
}
