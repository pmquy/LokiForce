package http

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"lokiforce.com/apps/core/internal/user/application"
	"lokiforce.com/apps/core/pkg/response"
)

type UserHandler struct {
	usecase application.UserUsecase
}

func NewUserHandler(usecase application.UserUsecase) *UserHandler {
	return &UserHandler{usecase: usecase}
}

type registerRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *UserHandler) register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Register request binding failed", "error", err)
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	input := application.RegisterUserInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	out, err := h.usecase.RegisterUser(c.Request.Context(), input)
	if err != nil {
		slog.Error("RegisterUser usecase execution failed", "error", err)
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	slog.Info("User registered successfully", "userID", out.UserID)
	response.Created(c, gin.H{
		"message": "user registered successfully",
		"user_id": out.UserID,
	})
}

func (h *UserHandler) login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Login request binding failed", "error", err)
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	input := application.LoginUserInput{
		Email:    req.Email,
		Password: req.Password,
	}

	out, err := h.usecase.LoginUser(c.Request.Context(), input)
	if err != nil {
		slog.Warn("LoginUser usecase execution failed", "error", err)
		response.Fail(c, http.StatusUnauthorized, err.Error())
		return
	}

	slog.Info("User logged in successfully")
	response.OK(c, gin.H{
		"token": out.Token,
	})
}

func (h *UserHandler) getProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		slog.Warn("GetProfile called without userID in context")
		response.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	out, err := h.usecase.GetUserByID(c.Request.Context(), userID.(string))
	if err != nil {
		slog.Error("GetUserByID usecase execution failed", "userID", userID, "error", err)
		response.Fail(c, http.StatusNotFound, err.Error())
		return
	}

	response.OK(c, out)
}
