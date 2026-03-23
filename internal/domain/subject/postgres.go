package subject

import (
	"context"
)

type Repository interface {
	// Create creates a new subject in the repository
	// and returns the created subject or an error if the operation fails.
	Create(ctx context.Context, subject *Subjects) error

	// GetByID retrieves a subject by their unique identifier.
	// It returns the subject or an error if the subject is not found.
	GetByID(ctx context.Context, id int64) (*Subjects, error)
	GetBySchoolID(ctx context.Context, subjectID int64) ([]*Subjects, error)

	// GetByCode retrieves a subject by subject code.
	// It returns the subject or an error if the subject is not found.
	GetByCode(ctx context.Context, code string) (*Subjects, error)

	// Update updates an existing subject's information in the repository.
	// It returns the updated subject or an error if the operation fails.
	Update(ctx context.Context, subject *Subjects) error

	// Delete removes a subject from the repository by their unique identifier.
	// It returns an error if the operation fails.
	Delete(ctx context.Context, id int64) error

	// List retrieves all subjects from the repository.
	// It returns a slice of subjects or an error if the operation fails.
	List(ctx context.Context, filter *SubjectFilterEntity) ([]*Subjects, error)
}
