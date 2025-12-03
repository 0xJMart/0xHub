package backend

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func TestNewClient(t *testing.T) {
	client := NewClient("http://localhost:8080")
	assert.NotNil(t, client)
	assert.Equal(t, "http://localhost:8080", client.baseURL)
	assert.NotNil(t, client.httpClient)
	assert.Equal(t, 10*time.Second, client.httpClient.Timeout)
}

func TestClient_HealthCheck_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/health", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	client := NewClient(server.URL)
	err := client.HealthCheck()
	assert.NoError(t, err)
}

func TestClient_HealthCheck_Failure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	err := client.HealthCheck()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "backend API error")
}

func TestClient_CreateProject_Success(t *testing.T) {
	var receivedProject Project
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/projects", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		err := json.NewDecoder(r.Body).Decode(&receivedProject)
		require.NoError(t, err)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(receivedProject)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	project := &Project{
		ID:          "test-1",
		Name:        "Test Project",
		Description: "A test project",
		URL:         "https://test.com",
		Category:    "testing",
		Status:      "active",
	}

	err := client.CreateProject(project)
	assert.NoError(t, err)
	assert.Equal(t, "test-1", receivedProject.ID)
	assert.Equal(t, "Test Project", receivedProject.Name)
}

func TestClient_CreateProject_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid project"})
	}))
	defer server.Close()

	client := NewClient(server.URL)
	project := &Project{
		ID:   "test-1",
		Name: "Test Project",
		URL:  "https://test.com",
	}

	err := client.CreateProject(project)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "backend API error")
}

func TestClient_UpdateProject_Success(t *testing.T) {
	var receivedProject Project
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/projects/test-1", r.URL.Path)
		assert.Equal(t, http.MethodPut, r.Method)

		err := json.NewDecoder(r.Body).Decode(&receivedProject)
		require.NoError(t, err)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(receivedProject)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	project := &Project{
		ID:          "test-1",
		Name:        "Updated Project",
		Description: "An updated project",
		URL:         "https://updated.com",
	}

	err := client.UpdateProject("test-1", project)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Project", receivedProject.Name)
}

func TestClient_UpdateProject_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "project not found"})
	}))
	defer server.Close()

	client := NewClient(server.URL)
	project := &Project{
		ID:   "test-1",
		Name: "Updated Project",
		URL:  "https://updated.com",
	}

	err := client.UpdateProject("test-1", project)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "backend API error")
}

func TestClient_GetProject_Success(t *testing.T) {
	expectedProject := Project{
		ID:          "test-1",
		Name:        "Test Project",
		Description: "A test project",
		URL:         "https://test.com",
		Category:    "testing",
		Status:      "active",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/projects/test-1", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expectedProject)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	project, err := client.GetProject("test-1")
	assert.NoError(t, err)
	assert.NotNil(t, project)
	assert.Equal(t, expectedProject.ID, project.ID)
	assert.Equal(t, expectedProject.Name, project.Name)
	assert.Equal(t, expectedProject.Description, project.Description)
}

func TestClient_GetProject_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "project not found"})
	}))
	defer server.Close()

	client := NewClient(server.URL)
	project, err := client.GetProject("non-existent")
	assert.Error(t, err)
	assert.Nil(t, project)
}

func TestClient_DeleteProject_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/projects/test-1", r.URL.Path)
		assert.Equal(t, http.MethodDelete, r.Method)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "deleted"})
	}))
	defer server.Close()

	client := NewClient(server.URL)
	err := client.DeleteProject("test-1")
	assert.NoError(t, err)
}

func TestClient_DeleteProject_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "project not found"})
	}))
	defer server.Close()

	client := NewClient(server.URL)
	err := client.DeleteProject("non-existent")
	assert.Error(t, err)
}

func TestClient_Timeout(t *testing.T) {
	// Create a server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(15 * time.Second) // Longer than client timeout
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	project := &Project{
		ID:   "test-1",
		Name: "Test Project",
		URL:  "https://test.com",
	}

	err := client.CreateProject(project)
	assert.Error(t, err)
	// Check for timeout-related error messages
	assert.True(t, 
		contains(err.Error(), "timeout") || 
		contains(err.Error(), "deadline exceeded") ||
		contains(err.Error(), "context deadline"),
		"Error should contain timeout-related message, got: %s", err.Error())
}

