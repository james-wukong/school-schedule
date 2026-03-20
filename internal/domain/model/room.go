package model

import (
	"maps"
)

type RoomID int
type RoomType string

const (
	Regular RoomType = "Regular"
	Lab     RoomType = "Lab"
	Gym     RoomType = "Gym"
)

type Room struct {
	// AvailableTimes: map of day -> available times
	ID             RoomID
	TenantID       int
	Name           string
	Capacity       int
	AvailableTimes map[DayOfWeek][]string
	Type           RoomType
}

func NewRoom(entity Room) *Room {
	r := &entity
	if len(entity.AvailableTimes) > 0 {
		maps.Copy(r.AvailableTimes, entity.AvailableTimes)
	}
	if r.Capacity == 0 {
		r.Capacity = 45
	}
	if r.Type == "" {
		r.Type = Regular
	}
	return r
}

func (c Room) CanFit(size int) bool {
	return c.Capacity >= size
}
