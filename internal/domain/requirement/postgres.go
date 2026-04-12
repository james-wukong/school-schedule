package requirement

import (
	"context"
)

type Repository interface {
	// GetByCode retrieves a requirement by requirement code.
	// It returns the requirement or an error if the requirement is not found.
	GetByVersion(ctx context.Context,
		schoolID, semesterID int64,
		version float64,
	) ([]*Requirements, error)
}
