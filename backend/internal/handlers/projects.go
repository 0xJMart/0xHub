package handlers

import (
	"net/http"

	"0xhub/backend/internal/models"
	"0xhub/backend/internal/store"

	"github.com/gin-gonic/gin"
)

// ProjectsHandler handles project-related HTTP requests
type ProjectsHandler struct {
	store *store.Store
}

// NewProjectsHandler creates a new projects handler
func NewProjectsHandler(store *store.Store) *ProjectsHandler {
	return &ProjectsHandler{
		store: store,
	}
}

// GetProjects returns all projects
func (h *ProjectsHandler) GetProjects(c *gin.Context) {
	projects := h.store.GetAll()
	c.JSON(http.StatusOK, gin.H{
		"projects": projects,
	})
}

// GetProject returns a single project by ID
func (h *ProjectsHandler) GetProject(c *gin.Context) {
	id := c.Param("id")
	project, exists := h.store.GetByID(id)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "project not found",
		})
		return
	}
	c.JSON(http.StatusOK, project)
}

// CreateProject creates a new project
func (h *ProjectsHandler) CreateProject(c *gin.Context) {
	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if project.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "id is required",
		})
		return
	}

	h.store.Create(&project)
	c.JSON(http.StatusCreated, project)
}

// UpdateProject updates an existing project
func (h *ProjectsHandler) UpdateProject(c *gin.Context) {
	id := c.Param("id")
	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	project.ID = id
	if !h.store.Update(&project) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "project not found",
		})
		return
	}

	c.JSON(http.StatusOK, project)
}

// DeleteProject deletes a project
func (h *ProjectsHandler) DeleteProject(c *gin.Context) {
	id := c.Param("id")
	if !h.store.Delete(id) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "project not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "project deleted",
	})
}
