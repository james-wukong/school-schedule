package teachertimeslot

import "errors"

var (
	ErrTeacherTimeslotNotFound      = errors.New("teacher timeslot not found")
	ErrTeacherTimeslotAlreadyExists = errors.New("teacher timeslot already exists")
)
