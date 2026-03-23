package roomtimeslot

import (
	"context"
	"time"
)

const (
	RedisUserPrefix = "room_timeslot:"
	RedisUerTTL     = 2 * 24 * time.Hour
)

type RedisCache interface {
	GetByIDs(ctx context.Context, entity *RoomTimeslots) (*RoomTimeslots, error)
	GetByRoomID(ctx context.Context, roomID int64) ([]*RoomTimeslots, error)
	GetByTimeslotID(ctx context.Context, slotID int64) ([]*RoomTimeslots, error)
	Set(ctx context.Context, entity *RoomTimeslots) error
	Update(ctx context.Context, entity *RoomTimeslots) error
}
