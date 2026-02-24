package user

import "errors"

var (
	// ErrNotFound is returned when user is not found
	ErrNotFound = errors.New("user not found")

	// ErrInvalidInput is returned when input validation fails
	ErrInvalidInput = errors.New("invalid input")

	// ErrWeakPassword is returned when password is too weak
	ErrWeakPassword = errors.New("password must be at least 8 characters")

	// ErrDuplicateEmail is returned when email already exists
	ErrDuplicateEmail = errors.New("email already exists")

	// ErrInvalidPassword is returned when password verification fails
	ErrInvalidPassword = errors.New("invalid email or password")

	// ErrUserInactive is returned when user is not active
	ErrUserInactive = errors.New("user is inactive")
)
