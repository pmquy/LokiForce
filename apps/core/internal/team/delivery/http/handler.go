package http

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"lokiforce.com/apps/core/internal/team/application"
	"lokiforce.com/apps/core/pkg/response"
)

type TeamHandler struct {
	usecase application.TeamUsecase
}

func NewTeamHandler(usecase application.TeamUsecase) *TeamHandler {
	return &TeamHandler{usecase: usecase}
}

type createTeamRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	OrgID       string `json:"org_id" binding:"required"`
}

func (h *TeamHandler) Create(c *gin.Context) {
	var req createTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Create team request binding failed", "error", err)
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	input := application.CreateTeamInput{
		Name:        req.Name,
		Description: req.Description,
		OrgID:       req.OrgID,
	}

	out, err := h.usecase.CreateTeam(c.Request.Context(), input)
	if err != nil {
		slog.Error("CreateTeam usecase failed", "error", err)
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	slog.Info("Team created successfully", "teamID", out.TeamID, "orgID", req.OrgID)
	response.Created(c, out)
}

func (h *TeamHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	out, err := h.usecase.GetTeamByID(c.Request.Context(), id)
	if err != nil {
		slog.Warn("GetTeamByID failed", "teamID", id, "error", err)
		response.Fail(c, http.StatusNotFound, err.Error())
		return
	}

	response.OK(c, out)
}

func (h *TeamHandler) ListByOrg(c *gin.Context) {
	orgID := c.Query("org_id")
	if orgID == "" {
		response.Fail(c, http.StatusBadRequest, "org_id query parameter is required")
		return
	}

	out, err := h.usecase.ListOrgTeams(c.Request.Context(), orgID)
	if err != nil {
		slog.Error("ListOrgTeams failed", "orgID", orgID, "error", err)
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(c, out)
}
