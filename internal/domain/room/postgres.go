package room

import (
	"context"
)

type Repository interface {
	// Create creates a new room in the repository
	// and returns the created room or an error if the operation fails.
	Create(ctx context.Context, room *Rooms) error

	// GetByID retrieves a room by their unique identifier.
	// It returns the room or an error if the room is not found.
	GetByID(ctx context.Context, id int64) (*Rooms, error)
	GetBySchoolID(ctx context.Context, schoolID int64) ([]*Rooms, error)

	// Update updates an existing room's information in the repository.
	// It returns the updated room or an error if the operation fails.
	Update(ctx context.Context, room *Rooms) error

	// Delete removes a room from the repository by their unique identifier.
	// It returns an error if the operation fails.
	Delete(ctx context.Context, id int64) error
}
