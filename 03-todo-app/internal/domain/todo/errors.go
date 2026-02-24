package todo

import "errors"

var (
	// ErrNotFound is returned when todo is not found
	ErrNotFound = errors.New("todo not found")

	// ErrInvalidInput is returned when input validation fails
	ErrInvalidInput = errors.New("invalid input")

	// ErrUnauthorized is returned when user tries to access todo they don't own
	ErrUnauthorized = errors.New("unauthorized: todo does not belong to user")

	// ErrInvalidUserID is returned when user ID is invalid
	ErrInvalidUserID = errors.New("invalid user id")
)
