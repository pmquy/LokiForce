package application

import (
	"github.com/google/uuid"
	"lokiforce.com/apps/core/internal/user/domain"
)

func NewUserUsecase(repo domain.UserRepository, tokenService TokenService) UserUsecase {
	return &userUsecaseImpl{

		repository:   repo,
		tokenService: tokenService,
	}
}

type userUsecaseImpl struct {
	repository   domain.UserRepository
	tokenService TokenService
}

func (u *userUsecaseImpl) RegisterUser(input RegisterUserInput) (RegisterUserOutput, error) {
	id := uuid.NewString()
	user, err := domain.NewUser(id, input.Username, input.Email, input.Password)

	if err != nil {
		return RegisterUserOutput{}, err
	}

	err = u.repository.CreateUser(user)

	if err != nil {
		return RegisterUserOutput{}, err
	}

	return RegisterUserOutput{
		UserID: user.ID,
	}, nil
}

func (u *userUsecaseImpl) LoginUser(input LoginUserInput) (LoginUserOutput, error) {
	user, err := u.repository.GetUserByEmail(input.Email)
	if err != nil {
		return LoginUserOutput{}, domain.ErrInvalidCredentials
	}

	if !user.VerifyPassword(input.Password) {
		return LoginUserOutput{}, domain.ErrInvalidCredentials
	}

	token, err := u.tokenService.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		return LoginUserOutput{}, err
	}

	return LoginUserOutput{
		Token: token,
	}, nil
}

func (u *userUsecaseImpl) GetUserByID(id string) (UserProfileOutput, error) {
	user, err := u.repository.GetUserByID(id)
	if err != nil {
		return UserProfileOutput{}, err
	}

	return UserProfileOutput{
		ID:       user.ID,
		Username: user.Username,
		Email:    string(user.Email),
		Role:     string(user.Role),
	}, nil
}
