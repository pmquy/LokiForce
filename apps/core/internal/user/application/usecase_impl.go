package application

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"lokiforce.com/apps/core/internal/user/domain"
	"lokiforce.com/apps/core/pkg/mail"
	"lokiforce.com/apps/core/pkg/mq"
)

func NewUserUsecase(
	repo domain.UserRepository,
	tokenService TokenService,
	mailService mail.MailService,
	msgQueue mq.MessageQueue,
) UserUsecase {
	return &userUsecaseImpl{
		repository:   repo,
		tokenService: tokenService,
		mailService:  mailService,
		mq:           msgQueue,
	}
}

type userUsecaseImpl struct {
	repository   domain.UserRepository
	tokenService TokenService
	mailService  mail.MailService
	mq           mq.MessageQueue
}

func (u *userUsecaseImpl) RegisterUser(ctx context.Context, input RegisterUserInput) (RegisterUserOutput, error) {
	id := uuid.NewString()
	user, err := domain.NewUser(id, input.Username, input.Email, input.Password)
	if err != nil {
		return RegisterUserOutput{}, err
	}

	err = u.repository.CreateUser(ctx, user)
	if err != nil {
		return RegisterUserOutput{}, err
	}

	emailSubject := "Welcome to LokiForce!"
	emailBody := fmt.Sprintf("Hi %s,\n\nThank you for registering at LokiForce portal!", user.Username)
	go func() {
		mailCtx := context.Background()
		if err := u.mailService.SendEmail(mailCtx, string(user.Email), emailSubject, emailBody); err != nil {
			slog.Error("Failed to send welcome email", "error", err, "email", string(user.Email))
		}
	}()

	eventPayload := mq.UserRegisteredEvent{
		UserID:   user.ID,
		Username: user.Username,
		Email:    string(user.Email),
	}
	go func() {
		mqCtx := context.Background()
		if err := u.mq.Publish(mqCtx, "user.registered", eventPayload); err != nil {
			slog.Error("Failed to publish user.registered event", "error", err, "userID", user.ID)
		}
	}()

	return RegisterUserOutput{
		UserID: user.ID,
	}, nil
}

func (u *userUsecaseImpl) LoginUser(ctx context.Context, input LoginUserInput) (LoginUserOutput, error) {
	user, err := u.repository.GetUserByEmail(ctx, input.Email)
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

func (u *userUsecaseImpl) GetUserByID(ctx context.Context, id string) (UserProfileOutput, error) {
	user, err := u.repository.GetUserByID(ctx, id)
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
