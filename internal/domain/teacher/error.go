package teacher

import "errors"

var (
	ErrTeacherNotFound      = errors.New("teacher not found")
	ErrTeacherAlreadyExists = errors.New("teacher already exists")
)
