package postgres

import (
	"context"
	"errors"

	"github.com/james-wukong/school-schedule/internal/domain/room"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type roomRepository struct {
	db  *gorm.DB
	log *zerolog.Logger
}

func NewRoomRepository(db *gorm.DB, log *zerolog.Logger) *roomRepository {
	return &roomRepository{db: db, log: log}
}

// Implement the roomRepository interface methods here,
// using GORM to interact with the PostgreSQL database.
func (r *roomRepository) Create(ctx context.Context, entity *room.Rooms) error {
	// First, create the school in the database
	err := r.db.WithContext(ctx).Create(entity).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error creating room in database")
		return err
	}
	return nil
}

// Implement GetByID method
func (r *roomRepository) GetByID(ctx context.Context, id int64) (*room.Rooms, error) {
	var record room.Rooms
	err := r.db.WithContext(ctx).
		Preload("School").
		First(&record, "id = ?", id).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, room.ErrRoomNotFound // Return nil if not found
		}
		return nil, err
	}
	return &record, nil
}

// Implement GetBySchoolID method
func (r *roomRepository) GetBySchoolID(
	ctx context.Context,
	schoolID int64,
) ([]*room.Rooms, error) {
	var rows []*room.Rooms
	err := r.db.WithContext(ctx).
		Preload("School").
		Find(&rows, "school_id = ?", schoolID).
		Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// Implement Update method
func (r *roomRepository) Update(ctx context.Context, entity *room.Rooms) error {
	// Update the school record in the database
	err := r.db.WithContext(ctx).Save(entity).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error updating room in database")
	}

	return err
}

// Implement Delete method
func (r *roomRepository) Delete(ctx context.Context, id int64) error {
	// Delete the school record from the database
	err := r.db.WithContext(ctx).Delete(&room.Rooms{}, "id = ?", id).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error deleting a room from database")
	}

	return err
}
