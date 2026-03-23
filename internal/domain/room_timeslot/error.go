package roomtimeslot

import "errors"

var (
	ErrRoomTimeslotNotFound      = errors.New("room timeslot not found")
	ErrRoomTimeslotAlreadyExists = errors.New("room timeslot already exists")
)
