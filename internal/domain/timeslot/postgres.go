package timeslot

import (
	"context"
)

type Repository interface {
	// Create creates a new timeslot in the repository
	// and returns the created timeslot or an error if the operation fails.
	Create(ctx context.Context, slot *Timeslots) error

	// GetByID retrieves a timeslot by their unique identifier.
	// It returns the timeslot or an error if the timeslot is not found.
	GetByID(ctx context.Context, id int64) (*Timeslots, error)
	GetBySemesterID(ctx context.Context, semesterID int64) ([]*Timeslots, error)

	// Update updates an existing timeslot's information in the repository.
	// It returns the updated timeslot or an error if the operation fails.
	Update(ctx context.Context, timeslot *Timeslots) error

	// Delete removes a timeslot from the repository by their unique identifier.
	// It returns an error if the operation fails.
	Delete(ctx context.Context, id int64) error

	// List retrieves all timeslots from the repository.
	// It returns a slice of timeslots or an error if the operation fails.
	List(ctx context.Context, filter *TimeslotFilterEntity) ([]*Timeslots, error)
}
