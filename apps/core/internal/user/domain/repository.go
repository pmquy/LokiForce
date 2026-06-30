package domain

import "context"

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error

	GetUserByID(ctx context.Context, id string) (*User, error)

	GetUserByUsername(ctx context.Context, username string) (*User, error)

	GetUserByEmail(ctx context.Context, email string) (*User, error)

	UpdateUser(ctx context.Context, user *User) error

	DeleteUser(ctx context.Context, id string) error
}
