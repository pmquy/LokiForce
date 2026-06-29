package http

import (
	"log/slog"
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
		slog.Warn("Create org called without userID in context")
		response.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req createOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Create org request binding failed", "error", err)
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	input := application.CreateOrgInput{
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     userID.(string),
	}

	out, err := h.usecase.CreateOrg(c.Request.Context(), input)
	if err != nil {
		slog.Error("CreateOrg usecase failed", "error", err)
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	slog.Info("Organization created successfully", "orgID", out.OrgID, "ownerID", userID)
	response.Created(c, out)
}

func (h *OrgHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	out, err := h.usecase.GetOrgByID(c.Request.Context(), id)
	if err != nil {
		slog.Warn("GetOrgByID failed", "orgID", id, "error", err)
		response.Fail(c, http.StatusNotFound, err.Error())
		return
	}

	response.OK(c, out)
}

func (h *OrgHandler) ListMyOrgs(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		slog.Warn("List my orgs called without userID in context")
		response.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	out, err := h.usecase.ListUserOrgs(c.Request.Context(), userID.(string))
	if err != nil {
		slog.Error("ListUserOrgs failed", "ownerID", userID, "error", err)
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(c, out)
}
