package domain

type UserRepository interface {
	// CreateUser creates a new user in the repository.
	CreateUser(user *User) error

	// GetUserByID retrieves a user by their ID.
	GetUserByID(id string) (*User, error)

	// GetUserByUsername retrieves a user by their username.
	GetUserByUsername(username string) (*User, error)

	GetUserByEmail(email string) (*User, error)

	// UpdateUser updates an existing user's information.
	UpdateUser(user *User) error

	// DeleteUser removes a user from the repository by their ID.
	DeleteUser(id string) error
}
