package class

import "errors"

var (
	ErrClassNotFound      = errors.New("class not found")
	ErrClassAlreadyExists = errors.New("class already exists")
)
