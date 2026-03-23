package subject

import "errors"

var (
	ErrSubjectNotFound      = errors.New("subject not found")
	ErrSubjectAlreadyExists = errors.New("subject already exists")
)
