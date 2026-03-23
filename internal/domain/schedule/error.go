package schedule

import "errors"

var (
	ErrScheduleNotFound      = errors.New("schedule not found")
	ErrScheduleAlreadyExists = errors.New("schedule already exists")
)
