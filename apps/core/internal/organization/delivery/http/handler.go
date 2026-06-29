package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"lokiforce.com/apps/core/internal/organization/application"
	"lokiforce.com/apps/core/pkg/response"
)

type OrgHandler struct {
	usecase application.OrgUsecase
}

func NewOrgHandler(usecase application.OrgUsecase) *OrgHandler {
	return &OrgHandler{usecase: usecase}
}

type createOrgRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func (h *OrgHandler) Create(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req createOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	input := application.CreateOrgInput{
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     userID.(string),
	}

	out, err := h.usecase.CreateOrg(input)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Created(c, out)
}

func (h *OrgHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	out, err := h.usecase.GetOrgByID(id)
	if err != nil {
		response.Fail(c, http.StatusNotFound, err.Error())
		return
	}

	response.OK(c, out)
}

func (h *OrgHandler) ListMyOrgs(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	out, err := h.usecase.ListUserOrgs(userID.(string))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(c, out)
}
