package roomtimeslot

import (
	"context"
)

type Repository interface {
	// Create creates a new room timeslot in the repository
	// and returns the created room timeslot or an error if the operation fails.
	Create(ctx context.Context, entity *RoomTimeslots) error

	// GetByIDs retrieves a room timeslot by their unique identifier.
	// It returns the room timeslot or an error if the room timeslot is not found.
	GetByIDs(ctx context.Context, entity *RoomTimeslots) (*RoomTimeslots, error)
	GetByRoomID(ctx context.Context, roomID int64) ([]*RoomTimeslots, error)
	GetByTimeslotID(ctx context.Context, slotID int64) ([]*RoomTimeslots, error)

	// Update updates an existing room timeslot's information in the repository.
	// It returns the updated room timeslot or an error if the operation fails.
	Update(ctx context.Context, entity *RoomTimeslots) error
}
