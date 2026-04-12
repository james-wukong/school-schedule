package postgres

import (
	"context"
	"errors"
	"strings"

	"github.com/james-wukong/school-schedule/internal/domain/subject"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type subjectRepository struct {
	db  *gorm.DB
	log *zerolog.Logger
}

func NewSubjectRepository(db *gorm.DB, log *zerolog.Logger) subject.Repository {
	return &subjectRepository{db: db, log: log}
}

// Implement the schoolRepository interface methods here,
// using GORM to interact with the PostgreSQL database.
func (r *subjectRepository) Create(ctx context.Context, entity *subject.Subjects) error {
	// First, create the school in the database
	err := r.db.WithContext(ctx).Create(entity).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error creating subject in database")
		return err
	}
	return nil
}

// Implement GetByID method
func (r *subjectRepository) GetByID(ctx context.Context, id int64) (*subject.Subjects, error) {
	var record subject.Subjects
	err := r.db.WithContext(ctx).
		Preload("School").
		First(&record, "id = ?", id).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, subject.ErrSubjectNotFound // Return nil if not found
		}
		return nil, err
	}
	return &record, nil
}

// Implement GetBySchoolID method
func (r *subjectRepository) GetBySchoolID(
	ctx context.Context,
	schoolID int64,
) ([]*subject.Subjects, error) {
	var rows []*subject.Subjects
	err := r.db.WithContext(ctx).
		Preload("School").
		Find(&rows, "school_id = ?", schoolID).
		Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// Implement GetByCode method
func (r *subjectRepository) GetByCode(ctx context.Context, code string) (*subject.Subjects, error) {
	var record subject.Subjects
	err := r.db.WithContext(ctx).
		Preload("School").
		First(&record, "code = ?", code).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return domain error
			return nil, subject.ErrSubjectNotFound
		}
		return nil, err
	}

	return &record, nil
}

// Implement Update method
func (r *subjectRepository) Update(ctx context.Context, entity *subject.Subjects) error {
	// Update the school record in the database
	err := r.db.WithContext(ctx).Save(entity).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error updating subject in database")
	}

	return err
}

// Implement Delete method
func (r *subjectRepository) Delete(ctx context.Context, id int64) error {
	// Delete the school record from the database
	err := r.db.WithContext(ctx).Delete(&subject.Subjects{}, "id = ?", id).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error deleting a subject from database")
	}

	return err
}

// Implement List method
func (r *subjectRepository) List(
	ctx context.Context,
	filter *subject.SubjectFilterEntity,
) ([]*subject.Subjects, error) {
	var rows []*subject.Subjects
	query := r.db.WithContext(ctx).Model(&subject.Subjects{})
	if filter != nil {
		if filter.Name != nil && *filter.Name != "" {
			query = query.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(*filter.Name)+"%")
		}
		if filter.Code != nil && *filter.Code != "" {
			query = query.Where("LOWER(code) LIKE ?", "%"+strings.ToLower(*filter.Code)+"%")
		}
		if filter.SchoolID != nil {
			query = query.Where("school_id = ", filter.SchoolID)
		}
		if filter.IsHeavy != nil {
			query = query.Where("is_heavy = ?", filter.IsHeavy)
		}
	}

	err := query.Limit(filter.Limit).
		Offset((filter.Page - 1) * filter.Limit).
		Find(&rows).
		Error
	return rows, err
}
