package teacher

import (
	"context"
)

type Repository interface {
	// Create creates a new teacher in the repository
	// and returns the created teacher or an error if the operation fails.
	Create(ctx context.Context, teacher *Teachers) error

	// GetByID retrieves a teacher by their unique identifier.
	// It returns the teacher or an error if the teacher is not found.
	GetByID(ctx context.Context, id int64) (*Teachers, error)
	GetBySchoolID(ctx context.Context, teacherID int64) ([]*Teachers, error)
	GetByEmployeeID(ctx context.Context, employeeID int64) ([]*Teachers, error)

	// GetByFirstName retrieves a teacher by teacher code.
	// It returns the teacher or an error if the teacher is not found.
	GetByName(ctx context.Context, firstName, lastName string) ([]*Teachers, error)

	// Update updates an existing teacher's information in the repository.
	// It returns the updated teacher or an error if the operation fails.
	Update(ctx context.Context, teacher *Teachers) error

	// Delete removes a teacher from the repository by their unique identifier.
	// It returns an error if the operation fails.
	Delete(ctx context.Context, id int64) error

	// List retrieves all teachers from the repository.
	// It returns a slice of teachers or an error if the operation fails.
	List(ctx context.Context, filter *TeacherFilterEntity) ([]*Teachers, error)
}
