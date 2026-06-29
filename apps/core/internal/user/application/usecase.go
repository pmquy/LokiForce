package application

import "context"

type RegisterUserInput struct {
	Username string
	Email    string
	Password string
}

type RegisterUserOutput struct {
	UserID string
}

type LoginUserInput struct {
	Email    string
	Password string
}

type LoginUserOutput struct {
	Token string
}

type UserProfileOutput struct {
	ID       string
	Username string
	Email    string
	Role     string
}

type UserUsecase interface {
	RegisterUser(ctx context.Context, input RegisterUserInput) (RegisterUserOutput, error)
	LoginUser(ctx context.Context, input LoginUserInput) (LoginUserOutput, error)
	GetUserByID(ctx context.Context, id string) (UserProfileOutput, error)
}
