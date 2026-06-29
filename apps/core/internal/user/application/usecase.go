package application

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
	RegisterUser(input RegisterUserInput) (RegisterUserOutput, error)
	LoginUser(input LoginUserInput) (LoginUserOutput, error)
	GetUserByID(id string) (UserProfileOutput, error)
}
