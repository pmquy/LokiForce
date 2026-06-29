package domain

import "context"

type UserRepository interface {
	// CreateUser creates a new user in the repository.
	CreateUser(ctx context.Context, user *User) error

	// GetUserByID retrieves a user by their ID.
	GetUserByID(ctx context.Context, id string) (*User, error)

	// GetUserByUsername retrieves a user by their username.
	GetUserByUsername(ctx context.Context, username string) (*User, error)

	GetUserByEmail(ctx context.Context, email string) (*User, error)

	// UpdateUser updates an existing user's information.
	UpdateUser(ctx context.Context, user *User) error

	// DeleteUser removes a user from the repository by their ID.
	DeleteUser(ctx context.Context, id string) error
}

