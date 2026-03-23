package school

import (
	"context"
)

type Repository interface {
	// Create creates a new school in the repository
	// and returns the created school or an error if the operation fails.
	Create(ctx context.Context, school *Schools) error

	// GetByID retrieves a school by their unique identifier.
	// It returns the school or an error if the school is not found.
	GetByID(ctx context.Context, id int64) (*Schools, error)

	// GetByCode retrieves a school by school code.
	// It returns the school or an error if the school is not found.
	GetByCode(ctx context.Context, code string) (*Schools, error)

	// Update updates an existing school's information in the repository.
	// It returns the updated school or an error if the operation fails.
	Update(ctx context.Context, school *Schools) error

	// Delete removes a school from the repository by their unique identifier.
	// It returns an error if the operation fails.
	Delete(ctx context.Context, id int64) error

	// List retrieves all schools from the repository.
	// It returns a slice of schools or an error if the operation fails.
	List(ctx context.Context, filter *SchoolFilterEntity) ([]*Schools, error)
}
