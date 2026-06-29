package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"lokiforce.com/apps/core/internal/project/application"
	"lokiforce.com/apps/core/pkg/response"
)

type ProjectHandler struct {
	usecase application.ProjectUsecase
}

func NewProjectHandler(usecase application.ProjectUsecase) *ProjectHandler {
	return &ProjectHandler{usecase: usecase}
}

type createProjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	OrgID       string `json:"org_id" binding:"required"`
}

func (h *ProjectHandler) Create(c *gin.Context) {
	var req createProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	input := application.CreateProjectInput{
		Name:        req.Name,
		Description: req.Description,
		OrgID:       req.OrgID,
	}

	out, err := h.usecase.CreateProject(input)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Created(c, out)
}

func (h *ProjectHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	out, err := h.usecase.GetProjectByID(id)
	if err != nil {
		response.Fail(c, http.StatusNotFound, err.Error())
		return
	}

	response.OK(c, out)
}

func (h *ProjectHandler) ListByOrg(c *gin.Context) {
	orgID := c.Query("org_id")
	if orgID == "" {
		response.Fail(c, http.StatusBadRequest, "org_id query parameter is required")
		return
	}

	out, err := h.usecase.ListOrgProjects(orgID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(c, out)
}
