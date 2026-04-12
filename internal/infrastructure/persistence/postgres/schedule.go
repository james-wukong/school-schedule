package postgres

import (
	"context"
	"errors"

	"github.com/james-wukong/school-schedule/internal/domain/schedule"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type scheduleRepository struct {
	db  *gorm.DB
	log *zerolog.Logger
}

func NewScheduleRepository(db *gorm.DB, log *zerolog.Logger) schedule.Repository {
	return &scheduleRepository{db: db, log: log}
}

// Implement the schoolRepository interface methods here,
// using GORM to interact with the PostgreSQL database.
func (r *scheduleRepository) Create(ctx context.Context, entity *schedule.Schedules) error {
	// First, create the school in the database
	err := r.db.WithContext(ctx).Create(entity).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error creating schedule in database")
		return err
	}
	return nil
}

// Implement GetByID method
func (r *scheduleRepository) GetByID(ctx context.Context, id int64) (*schedule.Schedules, error) {
	var record schedule.Schedules
	err := r.db.WithContext(ctx).
		Preload("Requirement").
		Preload("Room").
		Preload("School").
		Preload("Timeslot").
		Where("id = ?", id).
		First(&record).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schedule.ErrScheduleNotFound // Return nil if not found
		}
		return nil, err
	}
	return &record, nil
}

// Implement GetBySchoolID method
func (r *scheduleRepository) GetBySchoolID(
	ctx context.Context,
	schoolID int64,
) ([]*schedule.Schedules, error) {
	var rows []*schedule.Schedules
	err := r.db.WithContext(ctx).
		Preload("Requirement").
		Preload("Room").
		Preload("School").
		Preload("Timeslot").
		Find(&rows, "school_id = ?", schoolID).
		Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// Implement GetByRequirementID method
func (r *scheduleRepository) GetByRequirementID(
	ctx context.Context,
	requirementID int64,
) ([]*schedule.Schedules, error) {
	var rows []*schedule.Schedules
	err := r.db.WithContext(ctx).
		Preload("Requirement").
		Preload("Room").
		Preload("School").
		Preload("Timeslot").
		Find(&rows, "requirement = ", requirementID).
		Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// Implement GetByCode method
func (r *scheduleRepository) GetByVersion(
	ctx context.Context,
	schoolID int64,
	version float64,
) ([]*schedule.Schedules, error) {
	var rows []*schedule.Schedules
	err := r.db.WithContext(ctx).
		Preload("Requirement").
		Preload("Room").
		Preload("School").
		Preload("Timeslot").
		Find(&rows, "school_id = ? AND version = ?",
			schoolID, version).
		Error

	if err != nil {
		return nil, err
	}

	return rows, nil
}

// Implement Update method
func (r *scheduleRepository) Update(ctx context.Context, entity *schedule.Schedules) error {
	// Update the school record in the database
	err := r.db.WithContext(ctx).Save(entity).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error updating schedule in database")
	}

	return err
}

// Implement Delete method
func (r *scheduleRepository) Delete(ctx context.Context, id int64) error {
	// Delete the school record from the database
	err := r.db.WithContext(ctx).Delete(&schedule.Schedules{}, "id = ?", id).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error deleting a schedule from database")
	}

	return err
}

// Implement List method
func (r *scheduleRepository) List(
	ctx context.Context,
	filter *schedule.ScheduleFilterEntity,
) ([]*schedule.Schedules, error) {
	var rows []*schedule.Schedules
	query := r.db.WithContext(ctx).Model(&schedule.Schedules{})
	if filter != nil {
		if filter.SchoolID != nil {
			query = query.Where("school_id = ?", filter.SchoolID)
		}
		if filter.RequirementID != nil {
			query = query.Where("requirement_id = ?", filter.RequirementID)
		}
		if filter.RoomID != nil {
			query = query.Where("room_id = ?", filter.RoomID)
		}
		if filter.TimeslotID != nil {
			query = query.Where("timeslot_id = ?", filter.TimeslotID)
		}
		if filter.Version != nil {
			query = query.Where("version = ?", filter.Version)
		}
		if filter.Status != nil {
			query = query.Where("status = ?", filter.Status)
		}
	}

	err := query.Limit(filter.Limit).
		Offset((filter.Page - 1) * filter.Limit).
		Find(&rows).
		Error
	return rows, err
}
