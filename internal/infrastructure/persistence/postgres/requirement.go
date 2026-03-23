package postgres

import (
	"context"
	"errors"

	req "github.com/james-wukong/school-schedule/internal/domain/requirement"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type requirementRepository struct {
	db  *gorm.DB
	log *zerolog.Logger
}

func NewRequirementRepository(db *gorm.DB, log *zerolog.Logger) req.Repository {
	return &requirementRepository{db: db, log: log}
}

// Implement the schoolRepository interface methods here,
// using GORM to interact with the PostgreSQL database.
func (r *requirementRepository) Create(ctx context.Context, entity *req.Requirements) error {
	// First, create the school in the database
	err := r.db.WithContext(ctx).Create(entity).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error creating requirement in database")
		return err
	}
	return nil
}

// Implement GetByID method
func (r *requirementRepository) GetByID(ctx context.Context, id int64) (*req.Requirements, error) {
	var record req.Requirements
	err := r.db.WithContext(ctx).
		Preload("Class").
		Preload("School").
		Preload("Subject").
		Preload("Teacher").
		Where("id = ?", id).
		First(&record).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, req.ErrRequirementNotFound // Return nil if not found
		}
		return nil, err
	}
	return &record, nil
}

// Implement GetBySchoolID method
func (r *requirementRepository) GetBySchoolID(
	ctx context.Context,
	schoolID int64,
) ([]*req.Requirements, error) {
	var rows []*req.Requirements
	err := r.db.WithContext(ctx).
		Preload("Class").
		Preload("School").
		Preload("Subject").
		Preload("Teacher").
		Find(&rows, "school_id = ", schoolID).
		Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// Implement GetBySubjectID method
func (r *requirementRepository) GetBySubjectID(
	ctx context.Context,
	subjectID int64,
) ([]*req.Requirements, error) {
	var rows []*req.Requirements
	err := r.db.WithContext(ctx).
		Preload("Class").
		Preload("School").
		Preload("Subject").
		Preload("Teacher").
		Find(&rows, "subject_id = ", subjectID).
		Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// Implement GetBySubjectID method
func (r *requirementRepository) GetByTeacherID(
	ctx context.Context,
	teacherID int64,
) ([]*req.Requirements, error) {
	var rows []*req.Requirements
	err := r.db.WithContext(ctx).
		Preload("Class").
		Preload("School").
		Preload("Subject").
		Preload("Teacher").
		Find(&rows, "teacher_id = ", teacherID).
		Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// Implement GetByClassID method
func (r *requirementRepository) GetByClassID(
	ctx context.Context,
	classID int64,
) ([]*req.Requirements, error) {
	var rows []*req.Requirements
	err := r.db.WithContext(ctx).
		Preload("Class").
		Preload("School").
		Preload("Subject").
		Preload("Teacher").
		Find(&rows, "class_id = ", classID).
		Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// Implement GetByCode method
func (r *requirementRepository) GetByVersion(
	ctx context.Context,
	schoolID int64,
	version float64,
) ([]*req.Requirements, error) {
	var rows []*req.Requirements
	err := r.db.WithContext(ctx).
		Preload("Class").
		Preload("School").
		Preload("Subject").
		Preload("Teacher").
		Find(&rows, "school_id = ? AND version = ?",
			schoolID, version).
		Error

	if err != nil {
		return nil, err
	}

	return rows, nil
}

// Implement Update method
func (r *requirementRepository) Update(ctx context.Context, entity *req.Requirements) error {
	// Update the school record in the database
	err := r.db.WithContext(ctx).Save(entity).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error updating requirement in database")
	}

	return err
}

// Implement Delete method
func (r *requirementRepository) Delete(ctx context.Context, id int64) error {
	// Delete the school record from the database
	err := r.db.WithContext(ctx).Delete(&req.Requirements{}, "id = ?", id).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error deleting a requirement from database")
	}

	return err
}

// Implement List method
func (r *requirementRepository) List(
	ctx context.Context,
	filter *req.RequirementFilterEntity,
) ([]*req.Requirements, error) {
	var rows []*req.Requirements
	query := r.db.WithContext(ctx).Model(&req.Requirements{})
	if filter != nil {
		if filter.SchoolID != nil {
			query = query.Where("school_id = ?", filter.SchoolID)
		}
		if filter.SubjectID != nil {
			query = query.Where("subject_id = ?", filter.SubjectID)
		}
		if filter.TeacherID != nil {
			query = query.Where("teacher_id = ?", filter.TeacherID)
		}
		if filter.ClassID != nil {
			query = query.Where("class_id = ?", filter.ClassID)
		}
		if filter.Version != nil {
			query = query.Where("version = ?", filter.Version)
		}
	}

	err := query.Limit(filter.Limit).
		Offset((filter.Page - 1) * filter.Limit).
		Find(&rows).
		Error
	return rows, err
}
