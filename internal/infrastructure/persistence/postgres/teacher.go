package postgres

import (
	"context"
	"errors"
	"strings"

	"github.com/james-wukong/school-schedule/internal/domain/teacher"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type teacherRepository struct {
	db  *gorm.DB
	log *zerolog.Logger
}

func NewTeacherRepository(db *gorm.DB, log *zerolog.Logger) *teacherRepository {
	return &teacherRepository{db: db, log: log}
}

// Implement the schoolRepository interface methods here,
// using GORM to interact with the PostgreSQL database.
func (r *teacherRepository) Create(ctx context.Context, entity *teacher.Teachers) error {
	// First, create the school in the database
	err := r.db.WithContext(ctx).Create(entity).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error creating teacher in database")
		return err
	}
	return nil
}

// Implement GetByID method
func (r *teacherRepository) GetByID(ctx context.Context, id int64) (*teacher.Teachers, error) {
	var record teacher.Teachers
	err := r.db.WithContext(ctx).
		Preload("Subjects").
		Preload("School").
		Preload("Timeslots").
		First(&record, "id = ?", id).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, teacher.ErrTeacherNotFound // Return nil if not found
		}
		return nil, err
	}
	return &record, nil
}

// Implement GetBySchoolID method
func (r *teacherRepository) GetBySchoolID(
	ctx context.Context,
	schoolID int64,
) ([]*teacher.Teachers, error) {
	var rows []*teacher.Teachers
	err := r.db.WithContext(ctx).
		Preload("School").
		Preload("Subjects").
		Preload("Timeslots").
		Find(&rows, "school_id = ?", schoolID).
		Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// Implement GetBySchoolID method
func (r *teacherRepository) GetByEmployeeID(
	ctx context.Context,
	employeeID int64,
) ([]*teacher.Teachers, error) {
	var rows []*teacher.Teachers
	err := r.db.WithContext(ctx).
		Preload("School").
		Preload("Subjects").
		Preload("Timeslots").
		Find(&rows, "employee_id = ", employeeID).
		Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// Implement GetByCode method
func (r *teacherRepository) GetByName(
	ctx context.Context,
	firstName, lastName string,
) ([]*teacher.Teachers, error) {
	var rows []*teacher.Teachers
	err := r.db.WithContext(ctx).
		Preload("School").
		Preload("Subjects").
		Preload("Timeslots").
		Find(&rows, "LOWER(first_name) LIKE ? and LOWER(last_name) = ?",
			"%"+strings.ToLower(firstName)+"%",
			strings.ToLower(lastName)).
		Error

	if err != nil {
		return nil, err
	}

	return rows, nil
}

// Implement Update method
func (r *teacherRepository) Update(ctx context.Context, entity *teacher.Teachers) error {
	// Update the school record in the database
	err := r.db.WithContext(ctx).Save(entity).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error updating teacher in database")
	}

	return err
}

// Implement Delete method
func (r *teacherRepository) Delete(ctx context.Context, id int64) error {
	// Delete the school record from the database
	err := r.db.WithContext(ctx).Delete(&teacher.Teachers{}, "id = ?", id).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error deleting a teacher from database")
	}

	return err
}

// Implement List method
func (r *teacherRepository) List(
	ctx context.Context,
	filter *teacher.TeacherFilterEntity,
) ([]*teacher.Teachers, error) {
	var rows []*teacher.Teachers
	query := r.db.WithContext(ctx).Model(&teacher.Teachers{})
	if filter != nil {
		if filter.FirstName != nil && *filter.FirstName != "" {
			query = query.Where("LOWER(first_name) LIKE ?",
				"%"+strings.ToLower(*filter.FirstName)+"%",
			)
		}
		if filter.LastName != nil && *filter.LastName != "" {
			query = query.Where("LOWER(last_name) LIKE ?",
				"%"+strings.ToLower(*filter.LastName)+"%",
			)
		}
		if filter.SchoolID != nil {
			query = query.Where("school_id = ", filter.SchoolID)
		}
		if filter.EmployeeID != nil {
			query = query.Where("employee_id = ", filter.EmployeeID)
		}
		if filter.IsActive != nil {
			query = query.Where("is_active = ?", filter.IsActive)
		}
	}

	err := query.Limit(filter.Limit).
		Offset((filter.Page - 1) * filter.Limit).
		Find(&rows).
		Error
	return rows, err
}
