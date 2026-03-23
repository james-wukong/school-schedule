package teachersubject

import (
	"context"
)

type Repository interface {
	// Create creates a new teacher subject in the repository
	// and returns the created teacher subject or an error if the operation fails.
	Create(ctx context.Context, entity *TeacherSubjects) error

	// GetByIDs retrieves a teacher subject by their unique identifier.
	// It returns the teacher subject or an error if the teacher subject is not found.
	GetByIDs(ctx context.Context, entity *TeacherSubjects) (*TeacherSubjects, error)
	GetByTeacherID(ctx context.Context, teacherID int64) ([]*TeacherSubjects, error)
	GetBySubjectID(ctx context.Context, subjectID int64) ([]*TeacherSubjects, error)

	// Update updates an existing teacher subject's information in the repository.
	// It returns the updated teacher subject or an error if the operation fails.
	Update(ctx context.Context, entity *TeacherSubjects) error
}
