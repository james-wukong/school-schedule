package postgres

import (
	"context"
	"errors"

	"github.com/james-wukong/school-schedule/internal/domain/semester"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type semesterRepository struct {
	db  *gorm.DB
	log *zerolog.Logger
}

func NewSemesterRepository(db *gorm.DB, log *zerolog.Logger) *semesterRepository {
	return &semesterRepository{db: db, log: log}
}

// Implement the schoolRepository interface methods here,
// using GORM to interact with the PostgreSQL database.
func (r *semesterRepository) Create(ctx context.Context, entity *semester.Semesters) error {
	// First, create the school in the database
	err := r.db.WithContext(ctx).Create(entity).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error creating semester in database")
		return err
	}
	return nil
}

// Implement GetByID method
func (r *semesterRepository) GetByID(ctx context.Context, id int64) (*semester.Semesters, error) {
	var record semester.Semesters
	err := r.db.WithContext(ctx).
		Preload("School").
		First(&record, "id = ?", id).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, semester.ErrSemesterNotFound // Return nil if not found
		}
		return nil, err
	}
	return &record, nil
}

// Implement GetBySchoolID method
func (r *semesterRepository) GetBySchoolID(
	ctx context.Context,
	schoolID int64,
) ([]*semester.Semesters, error) {
	var rows []*semester.Semesters
	err := r.db.WithContext(ctx).
		Preload("School").
		Find(&rows, "school_id = ?", schoolID).
		Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// Implement Update method
func (r *semesterRepository) Update(ctx context.Context, entity *semester.Semesters) error {
	// Update the school record in the database
	err := r.db.WithContext(ctx).Save(entity).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error updating semester in database")
	}

	return err
}

// Implement Delete method
func (r *semesterRepository) Delete(ctx context.Context, id int64) error {
	// Delete the school record from the database
	err := r.db.WithContext(ctx).Delete(&semester.Semesters{}, "id = ?", id).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error deleting a semester from database")
	}

	return err
}

// Implement List method
func (r *semesterRepository) List(
	ctx context.Context,
	filter *semester.SemesterFilterEntity,
) ([]*semester.Semesters, error) {
	var rows []*semester.Semesters
	query := r.db.WithContext(ctx).Model(&semester.Semesters{})
	if filter != nil {
		if filter.SchoolID != nil {
			query = query.Where("school_id = ", filter.SchoolID)
		}
		if filter.Year != nil {
			query = query.Where("year = ?", filter.Year)
		}
		if filter.Semester != nil {
			query = query.Where("semester = ?", filter.Semester)
		}
	}

	err := query.Limit(filter.Limit).
		Offset((filter.Page - 1) * filter.Limit).
		Find(&rows).
		Error
	return rows, err
}
