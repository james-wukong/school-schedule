// Package roomtimeslot defines the room_subjects entity and related value objects.
// It represents how data looks in the database or business rules.
package roomtimeslot

import (
	"github.com/james-wukong/school-schedule/internal/domain/room"
	"github.com/james-wukong/school-schedule/internal/domain/timeslot"
)

// RoomTimeslots represents the room_subjects table in PostgreSQL.
// It uses GORM tags to handle identity columns and automatic timestamps.
type RoomTimeslots struct {
	// Primary Key composition
	RoomID     int64 `gorm:"primaryKey;column:room_id;not null" json:"room_id"`
	TimeslotID int64 `gorm:"primaryKey;column:timeslot_id;not null" json:"timeslot_id"`

	// Relationships (Belongs To)
	// These allow GORM to perform Preload("Room") or Preload("Timeslot")
	Room     *room.Rooms         `gorm:"foreignKey:RoomID;constraint:OnDelete:CASCADE" json:"room,omitempty"`
	Timeslot *timeslot.Timeslots `gorm:"foreignKey:TimeslotID;constraint:OnDelete:CASCADE" json:"timeslot,omitempty"`
}

func NewRoomTimeslots(roomID, slotID int64) *RoomTimeslots {
	return &RoomTimeslots{
		RoomID:     roomID,
		TimeslotID: slotID,
	}
}
