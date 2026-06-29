package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"lokiforce.com/apps/core/pkg/middleware"
)

type mockTokenValidator struct {
	validateFunc func(token string) (string, error)
}

func (m *mockTokenValidator) ValidateToken(token string) (string, error) {
	return m.validateFunc(token)
}

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		authHeader     string
		validateFunc   func(token string) (string, error)
		expectedStatus int
		expectedUserID string
	}{
		{
			name:       "Success Auth",
			authHeader: "Bearer valid_token",
			validateFunc: func(token string) (string, error) {
				if token == "valid_token" {
					return "user-123", nil
				}
				return "", errors.New("invalid")
			},
			expectedStatus: http.StatusOK,
			expectedUserID: "user-123",
		},
		{
			name:           "Missing Auth Header",
			authHeader:     "",
			validateFunc:   nil,
			expectedStatus: http.StatusUnauthorized,
			expectedUserID: "",
		},
		{
			name:           "Invalid Format Auth Header",
			authHeader:     "Basic credentials",
			validateFunc:   nil,
			expectedStatus: http.StatusUnauthorized,
			expectedUserID: "",
		},
		{
			name:       "Failed Validation",
			authHeader: "Bearer invalid_token",
			validateFunc: func(token string) (string, error) {
				return "", errors.New("token expired")
			},
			expectedStatus: http.StatusUnauthorized,
			expectedUserID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			validator := &mockTokenValidator{validateFunc: tt.validateFunc}
			r.Use(middleware.AuthMiddleware(validator))
			r.GET("/test", func(c *gin.Context) {
				userID, exists := c.Get("userID")
				if exists {
					c.String(http.StatusOK, userID.(string))
				} else {
					c.String(http.StatusInternalServerError, "no user ID")
				}
			})

			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				if w.Body.String() != tt.expectedUserID {
					t.Errorf("Expected body to be %s, got %s", tt.expectedUserID, w.Body.String())
				}
			}
		})
	}
}
