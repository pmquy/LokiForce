package jwt_test

import (
	"testing"

	"lokiforce.com/apps/core/internal/user/infrastructure/jwt"
)

func TestTokenService_GenerateAndValidateToken(t *testing.T) {
	secret := "test_secret_key"
	service := jwt.NewJWTService(secret)

	userID := "user-123"
	role := "admin"

	token, err := service.GenerateToken(userID, role)
	if err != nil {
		t.Fatalf("Expected no error generating token, got %v", err)
	}

	if token == "" {
		t.Fatalf("Expected a non-empty token string")
	}

	parsedUserID, err := service.ValidateToken(token)
	if err != nil {
		t.Fatalf("Expected no error validating token, got %v", err)
	}

	if parsedUserID != userID {
		t.Errorf("Expected parsed user ID to be %s, got %s", userID, parsedUserID)
	}
}

func TestTokenService_InvalidToken(t *testing.T) {
	service := jwt.NewJWTService("secret")
	_, err := service.ValidateToken("invalid.token.here")
	if err == nil {
		t.Errorf("Expected error validating invalid token, got nil")
	}
}
