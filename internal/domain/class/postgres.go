package class

import (
	"context"
)

type Repository interface {
	// Create creates a new class in the repository
	// and returns the created class or an error if the operation fails.
	Create(ctx context.Context, class *Classes) error

	// GetByID retrieves a class by their unique identifier.
	// It returns the class or an error if the class is not found.
	GetByID(ctx context.Context, id int64) (*Classes, error)
	GetBySemesterID(ctx context.Context, semesterID int64) ([]*Classes, error)

	// Update updates an existing class's information in the repository.
	// It returns the updated class or an error if the operation fails.
	Update(ctx context.Context, class *Classes) error

	// Delete removes a class from the repository by their unique identifier.
	// It returns an error if the operation fails.
	Delete(ctx context.Context, id int64) error
}
