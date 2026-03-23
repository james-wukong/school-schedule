package semester

import "errors"

var (
	ErrSemesterNotFound      = errors.New("semester not found")
	ErrSemesterAlreadyExists = errors.New("semester already exists")
)
