package http

import (
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
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	input := application.RegisterUserInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	out, err := h.usecase.RegisterUser(input)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Created(c, gin.H{
		"message": "user registered successfully",
		"user_id": out.UserID,
	})
}

func (h *UserHandler) login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	input := application.LoginUserInput{
		Email:    req.Email,
		Password: req.Password,
	}

	out, err := h.usecase.LoginUser(input)
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, err.Error())
		return
	}

	response.OK(c, gin.H{
		"token": out.Token,
	})
}

func (h *UserHandler) getProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	out, err := h.usecase.GetUserByID(userID.(string))
	if err != nil {
		response.Fail(c, http.StatusNotFound, err.Error())
		return
	}

	response.OK(c, out)
}
