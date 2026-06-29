package domain

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")

	ErrInvalidCredentials = errors.New("invalid credentials")

	ErrInvalidPassword = errors.New("invalid password: must be at least 8 characters long")

	ErrInvalidEmail = errors.New("invalid email: must not be empty")
)
