package domain_test

import (
	"testing"

	"lokiforce.com/apps/core/internal/user/domain"
)

func TestNewUser_ValidInput(t *testing.T) {

	user, err := domain.NewUser("uuid-123", "vinh", "vinh@gmail.com", "password123")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if user.Role != domain.MemberRole {
		t.Errorf("Expected default role to be member, got %v", user.Role)
	}
}

func TestNewUser_InvalidPassword(t *testing.T) {

	_, err := domain.NewUser("uuid-123", "vinh", "vinh@gmail.com", "short")
	if err != domain.ErrInvalidPassword {
		t.Errorf("Expected ErrInvalidPassword, got %v", err)
	}
}

func TestNewUser_InvalidEmail(t *testing.T) {

	_, err := domain.NewUser("uuid-123", "vinh", "", "password123")
	if err != domain.ErrInvalidEmail {
		t.Errorf("Expected ErrInvalidEmail, got %v", err)
	}
}
