package semester

import (
	"context"
)

type Repository interface {
	// Create creates a new semester in the repository
	// and returns the created semester or an error if the operation fails.
	Create(ctx context.Context, semester *Semesters) error

	// GetByID retrieves a semester by their unique identifier.
	// It returns the semester or an error if the semester is not found.
	GetByID(ctx context.Context, id int64) (*Semesters, error)
	GetBySchoolID(ctx context.Context, semesterID int64) ([]*Semesters, error)

	// Update updates an existing semester's information in the repository.
	// It returns the updated semester or an error if the operation fails.
	Update(ctx context.Context, semester *Semesters) error

	// Delete removes a semester from the repository by their unique identifier.
	// It returns an error if the operation fails.
	Delete(ctx context.Context, id int64) error

	// List retrieves all semesters from the repository.
	// It returns a slice of semesters or an error if the operation fails.
	List(ctx context.Context, filter *SemesterFilterEntity) ([]*Semesters, error)
}
