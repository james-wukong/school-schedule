package teachertimeslot

import (
	"context"
)

type Repository interface {
	// Create creates a new teacher timeslot in the repository
	// and returns the created teacher timeslot or an error if the operation fails.
	Create(ctx context.Context, entity *TeacherTimeslots) error

	// GetByIDs retrieves a teacher timeslot by their unique identifier.
	// It returns the teacher timeslot or an error if the teacher timeslot is not found.
	GetByIDs(ctx context.Context, entity *TeacherTimeslots) (*TeacherTimeslots, error)
	GetByTeacherID(ctx context.Context, teacherID int64) ([]*TeacherTimeslots, error)
	GetByTimeslotID(ctx context.Context, slotID int64) ([]*TeacherTimeslots, error)

	// Update updates an existing teacher timeslot's information in the repository.
	// It returns the updated teacher timeslot or an error if the operation fails.
	Update(ctx context.Context, entity *TeacherTimeslots) error
}
