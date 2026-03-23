package postgres

import (
	"context"
	"errors"

	"github.com/james-wukong/school-schedule/internal/domain/class"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type classRepository struct {
	db  *gorm.DB
	log *zerolog.Logger
}

func NewClassRepository(db *gorm.DB, log *zerolog.Logger) class.Repository {
	return &classRepository{db: db, log: log}
}

// Implement the classRepository interface methods here,
// using GORM to interact with the PostgreSQL database.
func (r *classRepository) Create(ctx context.Context, entity *class.Classes) error {
	// First, create the school in the database
	err := r.db.WithContext(ctx).Create(entity).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error creating class in database")
		return err
	}
	return nil
}

// Implement GetByID method
func (r *classRepository) GetByID(ctx context.Context, id int64) (*class.Classes, error) {
	var record class.Classes
	err := r.db.WithContext(ctx).
		Preload("Semester.School").
		First(&record, "id = ?", id).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, class.ErrClassNotFound // Return nil if not found
		}
		return nil, err
	}
	return &record, nil
}

// Implement GetBySchoolID method
func (r *classRepository) GetBySemesterID(
	ctx context.Context,
	semesterID int64,
) ([]*class.Classes, error) {
	var rows []*class.Classes
	err := r.db.WithContext(ctx).
		Preload("Semester.School").
		Find(&rows, "semester_id = ", semesterID).
		Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// Implement Update method
func (r *classRepository) Update(ctx context.Context, entity *class.Classes) error {
	// Update the school record in the database
	err := r.db.WithContext(ctx).Save(entity).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error updating class in database")
	}

	return err
}

// Implement Delete method
func (r *classRepository) Delete(ctx context.Context, id int64) error {
	// Delete the school record from the database
	err := r.db.WithContext(ctx).Delete(&class.Classes{}, "id = ?", id).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error deleting a class from database")
	}

	return err
}
