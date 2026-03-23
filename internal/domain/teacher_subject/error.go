package teachersubject

import "errors"

var (
	ErrTeacherSubjectNotFound      = errors.New("teacher subject not found")
	ErrTeacherSubjectAlreadyExists = errors.New("teacher subject already exists")
)
