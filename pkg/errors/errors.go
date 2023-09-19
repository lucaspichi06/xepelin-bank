package custom_errors

import "errors"

var (
	// handler errors
	ErrNotFound    = errors.New("registry not found")
	ErrEmptyFields = errors.New("fields can't be empty")
	ErrInvalidJSON = errors.New("invalid json")
	ErrInvalidID   = errors.New("invalid param id")

	// account errors
	ErrAccountExist = errors.New("there is already an account with this name")
)
