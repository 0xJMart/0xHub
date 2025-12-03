package handlers

import (
	"0xhub/backend/internal/models"
	"0xhub/backend/internal/store"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	testStore := store.NewStore()
	handler := NewProjectsHandler(testStore)

	router := gin.New()
	api := router.Group("/api")
	{
		api.GET("/projects", handler.GetProjects)
		api.GET("/projects/:id", handler.GetProject)
		api.POST("/projects", handler.CreateProject)
		api.PUT("/projects/:id", handler.UpdateProject)
		api.DELETE("/projects/:id", handler.DeleteProject)
	}

	return router
}

func TestGetProjects_Empty(t *testing.T) {
	router := setupRouter()
	req, _ := http.NewRequest("GET", "/api/projects", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	projects, ok := response["projects"].([]interface{})
	require.True(t, ok)
	assert.Equal(t, 0, len(projects))
}

func TestGetProjects_WithData(t *testing.T) {
	testStore := store.NewStore()
	testStore.Create(&models.Project{
		ID:          "test-1",
		Name:        "Test Project",
		Description: "A test project",
		URL:         "https://test.com",
	})

	// Recreate handler with populated store
	handler := NewProjectsHandler(testStore)
	gin.SetMode(gin.TestMode)
	testRouter := gin.New()
	api := testRouter.Group("/api")
	api.GET("/projects", handler.GetProjects)

	req, _ := http.NewRequest("GET", "/api/projects", nil)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	projects, ok := response["projects"].([]interface{})
	require.True(t, ok)
	assert.Equal(t, 1, len(projects))
}

func TestGetProject_Exists(t *testing.T) {
	testStore := store.NewStore()
	testStore.Create(&models.Project{
		ID:          "test-1",
		Name:        "Test Project",
		Description: "A test project",
		URL:         "https://test.com",
	})

	handler := NewProjectsHandler(testStore)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api := router.Group("/api")
	api.GET("/projects/:id", handler.GetProject)

	req, _ := http.NewRequest("GET", "/api/projects/test-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var project models.Project
	err := json.Unmarshal(w.Body.Bytes(), &project)
	require.NoError(t, err)
	assert.Equal(t, "test-1", project.ID)
	assert.Equal(t, "Test Project", project.Name)
}

func TestGetProject_NotFound(t *testing.T) {
	router := setupRouter()
	req, _ := http.NewRequest("GET", "/api/projects/non-existent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "project not found", response["error"])
}

func TestCreateProject_Success(t *testing.T) {
	router := setupRouter()
	project := models.Project{
		ID:          "new-project",
		Name:        "New Project",
		Description: "A new project",
		URL:         "https://new.com",
	}

	jsonData, _ := json.Marshal(project)
	req, _ := http.NewRequest("POST", "/api/projects", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var created models.Project
	err := json.Unmarshal(w.Body.Bytes(), &created)
	require.NoError(t, err)
	assert.Equal(t, project.ID, created.ID)
	assert.Equal(t, project.Name, created.Name)
}

func TestCreateProject_MissingID(t *testing.T) {
	router := setupRouter()
	project := models.Project{
		Name:        "New Project",
		Description: "A new project",
		URL:         "https://new.com",
	}

	jsonData, _ := json.Marshal(project)
	req, _ := http.NewRequest("POST", "/api/projects", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"].(string), "id is required")
}

func TestCreateProject_InvalidJSON(t *testing.T) {
	router := setupRouter()
	req, _ := http.NewRequest("POST", "/api/projects", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateProject_Success(t *testing.T) {
	testStore := store.NewStore()
	testStore.Create(&models.Project{
		ID:          "test-1",
		Name:        "Test Project",
		Description: "A test project",
		URL:         "https://test.com",
	})

	handler := NewProjectsHandler(testStore)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api := router.Group("/api")
	api.PUT("/projects/:id", handler.UpdateProject)

	updated := models.Project{
		Name:        "Updated Project",
		Description: "An updated project",
		URL:         "https://updated.com",
	}

	jsonData, _ := json.Marshal(updated)
	req, _ := http.NewRequest("PUT", "/api/projects/test-1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result models.Project
	err := json.Unmarshal(w.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, "test-1", result.ID)
	assert.Equal(t, "Updated Project", result.Name)
}

func TestUpdateProject_NotFound(t *testing.T) {
	router := setupRouter()
	updated := models.Project{
		Name: "Updated Project",
		URL:  "https://updated.com",
	}

	jsonData, _ := json.Marshal(updated)
	req, _ := http.NewRequest("PUT", "/api/projects/non-existent", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "project not found", response["error"])
}

func TestDeleteProject_Success(t *testing.T) {
	testStore := store.NewStore()
	testStore.Create(&models.Project{
		ID:          "test-1",
		Name:        "Test Project",
		Description: "A test project",
		URL:         "https://test.com",
	})

	handler := NewProjectsHandler(testStore)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api := router.Group("/api")
	api.DELETE("/projects/:id", handler.DeleteProject)

	req, _ := http.NewRequest("DELETE", "/api/projects/test-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify project is deleted
	_, exists := testStore.GetByID("test-1")
	assert.False(t, exists, "Project should be deleted")
}

func TestDeleteProject_NotFound(t *testing.T) {
	router := setupRouter()
	req, _ := http.NewRequest("DELETE", "/api/projects/non-existent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "project not found", response["error"])
}

