package application_test

import (
	"context"
	"testing"

	"lokiforce.com/apps/core/internal/user/application"
	"lokiforce.com/apps/core/internal/user/domain"
	"lokiforce.com/apps/core/internal/user/domain/mocks"
	"lokiforce.com/apps/core/internal/user/infrastructure/jwt"
)

func TestRegisterUser_ValidInput(t *testing.T) {
	repo := mocks.NewMockUserRepository()
	tokenService := jwt.NewJWTService("test_secret")
	usecase := application.NewUserUsecase(repo, tokenService)

	input := application.RegisterUserInput{
		Username: "vinh",
		Email:    "vinh@gmail.com",
		Password: "password123",
	}

	output, err := usecase.RegisterUser(context.Background(), input)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if output.UserID == "" {
		t.Errorf("Expected a valid user ID, got empty string")
	}
}

func TestRegisterUser_InvalidPassword(t *testing.T) {
	repo := mocks.NewMockUserRepository()
	tokenService := jwt.NewJWTService("test_secret")
	usecase := application.NewUserUsecase(repo, tokenService)

	input := application.RegisterUserInput{
		Username: "vinh",
		Email:    "vinh@gmail.com",
		Password: "short",
	}
	_, err := usecase.RegisterUser(context.Background(), input)
	if err != domain.ErrInvalidPassword {
		t.Errorf("Expected ErrInvalidPassword, got %v", err)
	}
}

func TestRegisterUser_InvalidEmail(t *testing.T) {
	repo := mocks.NewMockUserRepository()
	tokenService := jwt.NewJWTService("test_secret")
	usecase := application.NewUserUsecase(repo, tokenService)

	input := application.RegisterUserInput{
		Username: "vinh",
		Email:    "",
		Password: "password123",
	}

	_, err := usecase.RegisterUser(context.Background(), input)
	if err != domain.ErrInvalidEmail {
		t.Errorf("Expected ErrInvalidEmail, got %v", err)
	}
}
