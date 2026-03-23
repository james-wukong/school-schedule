package postgres

import (
	"context"
	"errors"

	tt "github.com/james-wukong/school-schedule/internal/domain/room_timeslot"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type roomTimeslotRepository struct {
	db  *gorm.DB
	log *zerolog.Logger
}

func NewRoomTimeslotRepository(db *gorm.DB, log *zerolog.Logger) tt.Repository {
	return &roomTimeslotRepository{db: db, log: log}
}

// Implement the roomTimeslotRepository interface methods here,
// using GORM to interact with the PostgreSQL database.
func (r *roomTimeslotRepository) Create(ctx context.Context, entity *tt.RoomTimeslots) error {
	// First, create the school in the database
	err := r.db.WithContext(ctx).
		Create(entity).
		Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error creating a room subject in database")
		return err
	}
	return nil
}

// Implement GetByID method
func (r *roomTimeslotRepository) GetByIDs(
	ctx context.Context,
	entity *tt.RoomTimeslots,
) (*tt.RoomTimeslots, error) {
	var record tt.RoomTimeslots
	err := r.db.WithContext(ctx).
		Preload("Room").
		Preload("Timeslot").
		First(&record, "room_id = ? and timeslot_id = ?",
			entity.RoomID, entity.TimeslotID,
		).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, tt.ErrRoomTimeslotNotFound // Return nil if not found
		}
		return nil, err
	}
	return &record, nil
}

// Implement GetByRoomID method
func (r *roomTimeslotRepository) GetByRoomID(
	ctx context.Context,
	roomID int64,
) ([]*tt.RoomTimeslots, error) {
	var rows []*tt.RoomTimeslots
	err := r.db.WithContext(ctx).
		Preload("Room").
		Preload("Timeslot").
		Find(&rows, "room_id = ", roomID).
		Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// Implement GetBySubjectID method
func (r *roomTimeslotRepository) GetByTimeslotID(
	ctx context.Context,
	slotID int64,
) ([]*tt.RoomTimeslots, error) {
	var rows []*tt.RoomTimeslots
	err := r.db.WithContext(ctx).
		Preload("Room").
		Preload("Timeslot").
		Find(&rows, "timeslot_id = ", slotID).
		Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// Implement Update method
func (r *roomTimeslotRepository) Update(
	ctx context.Context,
	entity *tt.RoomTimeslots,
) error {
	// Update the school record in the database
	err := r.db.WithContext(ctx).Save(entity).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error updating room in database")
	}

	return err
}
