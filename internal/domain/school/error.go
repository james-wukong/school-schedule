package school

import "errors"

var (
	ErrSchoolNotFound           = errors.New("school not found")
	ErrInvalidAuth              = errors.New("invalid authorization")
	ErrSchoolEmailAlreadyExists = errors.New("school email already exists")
	ErrInvalidEmailFormat       = errors.New("invalid email format")
	ErrSchoolAlreadyExists      = errors.New("school already exists")
)
