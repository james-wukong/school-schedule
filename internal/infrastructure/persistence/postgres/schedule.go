package postgres

import (
	"context"
	"errors"
	"math/rand/v2"

	"github.com/james-wukong/school-schedule/internal/domain/schedule"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type scheduleRepository struct {
	db  *gorm.DB
	log *zerolog.Logger
}

func NewScheduleRepository(db *gorm.DB, log *zerolog.Logger) *scheduleRepository {
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

func (r *scheduleRepository) CreateInBatches(
	ctx context.Context, t []*schedule.Schedules,
) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// remove all rows with same semester id and version number
		if len(t) == 0 {
			return nil
		}
		if err := tx.Delete(&schedule.Schedules{},
			"semester_id = ? and version = ?",
			t[0].SemesterID,
			t[0].Version,
		).Error; err != nil {
			return err
		}
		// create schedules
		return tx.CreateInBatches(t, 100).Error
	})
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

func (r *scheduleRepository) CreateVersionNumber(
	ctx context.Context, semesterID int64,
) decimal.Decimal {
	var m schedule.Schedules
	err := r.db.WithContext(ctx).
		Select("id", "version").
		Where("semester_id = ?", semesterID).
		Order("version DESC").
		First(&m).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return decimal.NewFromFloat(1.00)
		}
	}
	randMin := 0.01
	randMax := 0.06

	res := randMin + rand.Float64()*(randMax-randMin)
	return m.Version.Add(decimal.NewFromFloat(res))
}
