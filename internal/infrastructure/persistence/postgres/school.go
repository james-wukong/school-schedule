// Package postgres implements repositories using PostgreSQL as the database.
// It implements the interfaces defined in the domain,
// providing methods for creating, retrieving, updating, and deleting records
// in a PostgreSQL database using GORM as the ORM.
package postgres

import (
	"context"
	"errors"
	"strings"

	"github.com/james-wukong/school-schedule/internal/domain/school"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type schoolRepository struct {
	db  *gorm.DB
	log *zerolog.Logger
}

func NewSchoolRepository(db *gorm.DB, log *zerolog.Logger) *schoolRepository {
	return &schoolRepository{db: db, log: log}
}

// Implement the schoolRepository interface methods here,
// using GORM to interact with the PostgreSQL database.
func (r *schoolRepository) Create(ctx context.Context, entity *school.Schools) error {
	// First, create the school in the database
	err := r.db.WithContext(ctx).Create(entity).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error creating school in database")
		return err
	}
	return nil
}

// Implement GetByID method
func (r *schoolRepository) GetByID(ctx context.Context, id int64) (*school.Schools, error) {
	var record school.Schools
	err := r.db.WithContext(ctx).First(&record, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, school.ErrSchoolNotFound // Return nil if not found
		}
		return nil, err
	}
	return &record, nil
}

// Implement GetByCode method
func (r *schoolRepository) GetByCode(ctx context.Context, code string) (*school.Schools, error) {
	var record school.Schools
	err := r.db.WithContext(ctx).First(&record, "code = ?", code).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return domain error
			return nil, school.ErrSchoolNotFound
		}
		return nil, err
	}

	return &record, nil
}

// Implement Update method
func (r *schoolRepository) Update(ctx context.Context, entity *school.Schools) error {
	// Update the school record in the database
	err := r.db.WithContext(ctx).Save(entity).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error updating school in database")
	}

	return err
}

// Implement Delete method
func (r *schoolRepository) Delete(ctx context.Context, id int64) error {
	// Delete the school record from the database
	err := r.db.WithContext(ctx).Delete(&school.Schools{}, "id = ?", id).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error deleting a school from database")
	}

	return err
}

// Implement List method
func (r *schoolRepository) List(
	ctx context.Context,
	filter *school.SchoolFilterEntity,
) ([]*school.Schools, error) {
	var rows []*school.Schools
	query := r.db.WithContext(ctx).Model(&school.Schools{})
	if filter != nil {
		if filter.Email != nil && *filter.Email != "" {
			query = query.Where("email = ?", filter.Email)
		}
		if filter.IsActive != nil {
			query = query.Where("is_active = ?", filter.IsActive)
		}
		if filter.Name != nil && *filter.Name != "" {
			query = query.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(*filter.Name)+"%")
		}
		if filter.Code != nil && *filter.Code != "" {
			query = query.Where("LOWER(code) LIKE ?", "%"+strings.ToLower(*filter.Name)+"%")
		}
	}

	err := query.Limit(filter.Limit).
		Offset((filter.Page - 1) * filter.Limit).
		Find(&rows).
		Error
	return rows, err
}
