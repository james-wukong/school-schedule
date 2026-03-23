package requirement

import "errors"

var (
	ErrRequirementNotFound      = errors.New("requirement not found")
	ErrRequirementAlreadyExists = errors.New("requirement already exists")
)
