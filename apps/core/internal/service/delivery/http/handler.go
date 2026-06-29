package http

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"lokiforce.com/apps/core/internal/service/application"
	"lokiforce.com/apps/core/pkg/response"
)

type ServiceHandler struct {
	usecase application.ServiceUsecase
}

func NewServiceHandler(usecase application.ServiceUsecase) *ServiceHandler {
	return &ServiceHandler{usecase: usecase}
}

type createServiceRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	ProjectID   string `json:"project_id" binding:"required"`
	TemplateID  string `json:"template_id" binding:"required"`
}

type updateServiceRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func (h *ServiceHandler) Create(c *gin.Context) {
	var req createServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Create service binding failed", "error", err)
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	input := application.CreateServiceInput{
		Name:        req.Name,
		Description: req.Description,
		ProjectID:   req.ProjectID,
		TemplateID:  req.TemplateID,
	}

	out, err := h.usecase.CreateService(c.Request.Context(), input)
	if err != nil {
		slog.Error("CreateService failed", "error", err)
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	slog.Info("Service created and scaffolded successfully", "serviceID", out.ServiceID, "projectID", req.ProjectID, "repoURL", out.RepositoryURL)
	response.Created(c, out)
}

func (h *ServiceHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	out, err := h.usecase.GetServiceByID(c.Request.Context(), id)
	if err != nil {
		slog.Warn("GetServiceByID failed", "serviceID", id, "error", err)
		response.Fail(c, http.StatusNotFound, err.Error())
		return
	}

	response.OK(c, out)
}

func (h *ServiceHandler) ListByProject(c *gin.Context) {
	projectID := c.Query("project_id")
	if projectID == "" {
		response.Fail(c, http.StatusBadRequest, "project_id query parameter is required")
		return
	}

	out, err := h.usecase.ListProjectServices(c.Request.Context(), projectID)
	if err != nil {
		slog.Error("ListProjectServices failed", "projectID", projectID, "error", err)
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(c, out)
}

func (h *ServiceHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req updateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Update service binding failed", "error", err)
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	input := application.UpdateServiceInput{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.usecase.UpdateService(c.Request.Context(), input); err != nil {
		slog.Error("UpdateService failed", "serviceID", id, "error", err)
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	slog.Info("Service updated successfully", "serviceID", id)
	response.OK(c, gin.H{"message": "service updated successfully"})
}

func (h *ServiceHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.usecase.DeleteService(c.Request.Context(), id); err != nil {
		slog.Error("DeleteService failed", "serviceID", id, "error", err)
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	slog.Info("Service deleted successfully", "serviceID", id)
	response.OK(c, gin.H{"message": "service deleted successfully"})
}

func (h *ServiceHandler) ListTemplates(c *gin.Context) {
	out, err := h.usecase.ListTemplates(c.Request.Context())
	if err != nil {
		slog.Error("ListTemplates failed", "error", err)
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.OK(c, out)
}
