package errs

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrNotExists          = errors.New("resource does not exist")
	ErrAlreadyExists      = errors.New("resource already exists")
)
