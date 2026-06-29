package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"lokiforce.com/apps/core/internal/user/application"
)

type mockUserUsecase struct {
	registerFunc func(ctx context.Context, input application.RegisterUserInput) (application.RegisterUserOutput, error)
	loginFunc    func(ctx context.Context, input application.LoginUserInput) (application.LoginUserOutput, error)
	getByIDFunc  func(ctx context.Context, id string) (application.UserProfileOutput, error)
}

func (m *mockUserUsecase) RegisterUser(ctx context.Context, input application.RegisterUserInput) (application.RegisterUserOutput, error) {
	return m.registerFunc(ctx, input)
}

func (m *mockUserUsecase) LoginUser(ctx context.Context, input application.LoginUserInput) (application.LoginUserOutput, error) {
	return m.loginFunc(ctx, input)
}

func (m *mockUserUsecase) GetUserByID(ctx context.Context, id string) (application.UserProfileOutput, error) {
	return m.getByIDFunc(ctx, id)
}

func TestUserHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success Registration", func(t *testing.T) {
		mockUsecase := &mockUserUsecase{
			registerFunc: func(ctx context.Context, input application.RegisterUserInput) (application.RegisterUserOutput, error) {
				return application.RegisterUserOutput{UserID: "user-123"}, nil
			},
		}

		handler := NewUserHandler(mockUsecase)
		r := gin.New()
		r.POST("/register", handler.register)

		body := map[string]string{
			"username": "vinh",
			"email":    "vinh@gmail.com",
			"password": "password123",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status Created (201), got %d", w.Code)
		}
	})

	t.Run("Invalid Payload", func(t *testing.T) {
		handler := NewUserHandler(&mockUserUsecase{})
		r := gin.New()
		r.POST("/register", handler.register)

		body := map[string]string{
			"username": "",
			"email":    "invalid-email",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status BadRequest (400), got %d", w.Code)
		}
	})
}

func TestUserHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success Login", func(t *testing.T) {
		mockUsecase := &mockUserUsecase{
			loginFunc: func(ctx context.Context, input application.LoginUserInput) (application.LoginUserOutput, error) {
				return application.LoginUserOutput{Token: "jwt_token"}, nil
			},
		}

		handler := NewUserHandler(mockUsecase)
		r := gin.New()
		r.POST("/login", handler.login)

		body := map[string]string{
			"email":    "vinh@gmail.com",
			"password": "password123",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status OK (200), got %d", w.Code)
		}
	})

	t.Run("Unauthorized Login", func(t *testing.T) {
		mockUsecase := &mockUserUsecase{
			loginFunc: func(ctx context.Context, input application.LoginUserInput) (application.LoginUserOutput, error) {
				return application.LoginUserOutput{}, errors.New("invalid credentials")
			},
		}

		handler := NewUserHandler(mockUsecase)
		r := gin.New()
		r.POST("/login", handler.login)

		body := map[string]string{
			"email":    "vinh@gmail.com",
			"password": "wrongpassword",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status Unauthorized (401), got %d", w.Code)
		}
	})
}

func TestUserHandler_GetProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success Get Profile", func(t *testing.T) {
		mockUsecase := &mockUserUsecase{
			getByIDFunc: func(ctx context.Context, id string) (application.UserProfileOutput, error) {
				return application.UserProfileOutput{
					ID:       "user-123",
					Username: "vinh",
					Email:    "vinh@gmail.com",
					Role:     "member",
				}, nil
			},
		}

		handler := NewUserHandler(mockUsecase)
		r := gin.New()
		r.GET("/profile", func(c *gin.Context) {
			c.Set("userID", "user-123")
		}, handler.getProfile)

		req, _ := http.NewRequest(http.MethodGet, "/profile", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status OK (200), got %d", w.Code)
		}
	})
}
