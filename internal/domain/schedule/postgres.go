package schedule

import (
	"context"
)

type Repository interface {
	// Create creates a new schedule in the repository
	// and returns the created schedule or an error if the operation fails.
	Create(ctx context.Context, schedule *Schedules) error

	// GetByID retrieves a schedule by their unique identifier.
	// It returns the schedule or an error if the schedule is not found.
	GetByID(ctx context.Context, id int64) (*Schedules, error)
	GetBySchoolID(ctx context.Context, schoolID int64) ([]*Schedules, error)
	GetByRequirementID(ctx context.Context, requirementID int64) ([]*Schedules, error)

	// GetByCode retrieves a schedule by schedule code.
	// It returns the schedule or an error if the schedule is not found.
	GetByVersion(ctx context.Context, schoolID int64, version float64) ([]*Schedules, error)

	// Update updates an existing schedule's information in the repository.
	// It returns the updated schedule or an error if the operation fails.
	Update(ctx context.Context, schedule *Schedules) error

	// Delete removes a schedule from the repository by their unique identifier.
	// It returns an error if the operation fails.
	Delete(ctx context.Context, id int64) error

	// List retrieves all schedules from the repository.
	// It returns a slice of schedules or an error if the operation fails.
	List(ctx context.Context, filter *ScheduleFilterEntity) ([]*Schedules, error)
}
