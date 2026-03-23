package requirement

import (
	"context"
)

type Repository interface {
	// Create creates a new requirement in the repository
	// and returns the created requirement or an error if the operation fails.
	Create(ctx context.Context, requirement *Requirements) error

	// GetByID retrieves a requirement by their unique identifier.
	// It returns the requirement or an error if the requirement is not found.
	GetByID(ctx context.Context, id int64) (*Requirements, error)
	GetBySchoolID(ctx context.Context, requirementID int64) ([]*Requirements, error)
	GetBySubjectID(ctx context.Context, subjectID int64) ([]*Requirements, error)
	GetByTeacherID(ctx context.Context, teacherID int64) ([]*Requirements, error)
	GetByClassID(ctx context.Context, classID int64) ([]*Requirements, error)

	// GetByCode retrieves a requirement by requirement code.
	// It returns the requirement or an error if the requirement is not found.
	GetByVersion(ctx context.Context, requirementID int64, version float64) ([]*Requirements, error)

	// Update updates an existing requirement's information in the repository.
	// It returns the updated requirement or an error if the operation fails.
	Update(ctx context.Context, requirement *Requirements) error

	// Delete removes a requirement from the repository by their unique identifier.
	// It returns an error if the operation fails.
	Delete(ctx context.Context, id int64) error

	// List retrieves all requirements from the repository.
	// It returns a slice of requirements or an error if the operation fails.
	List(ctx context.Context, filter *RequirementFilterEntity) ([]*Requirements, error)
}
