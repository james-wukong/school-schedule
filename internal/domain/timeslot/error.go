package timeslot

import "errors"

var (
	ErrTimeslotNotFound      = errors.New("timeslot not found")
	ErrTimeslotAlreadyExists = errors.New("timeslot already exists")
)
