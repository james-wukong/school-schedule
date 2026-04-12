package postgres

import (
	"context"

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

// Implement GetByCode method
func (r *requirementRepository) GetByVersion(
	ctx context.Context,
	schoolID, semesterID int64,
	version float64,
) ([]*req.Requirements, error) {
	var rows []*req.Requirements
	err := r.db.WithContext(ctx).
		Preload("Class").
		Preload("School").
		Preload("Semester").
		Preload("Subject").
		Preload("Teacher").
		Preload("Teacher.Subjects").
		Preload("Teacher.Timeslots").
		Find(&rows, "school_id = ? AND semester_id = ? AND version = ?",
			schoolID, semesterID, version).
		Error

	if err != nil {
		return nil, err
	}

	return rows, nil
}
