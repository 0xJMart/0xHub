package store

import (
	"0xhub/backend/internal/models"
	"testing"
)

func TestNewStore(t *testing.T) {
	store := NewStore()
	if store == nil {
		t.Fatal("NewStore() returned nil")
	}
	if store.projects == nil {
		t.Fatal("Store projects map is nil")
	}
	if len(store.projects) != 0 {
		t.Fatal("New store should be empty")
	}
}

func TestStore_Create(t *testing.T) {
	store := NewStore()
	project := &models.Project{
		ID:          "test-1",
		Name:        "Test Project",
		Description: "A test project",
		URL:         "https://test.com",
	}

	store.Create(project)

	if len(store.projects) != 1 {
		t.Fatalf("Expected 1 project, got %d", len(store.projects))
	}

	retrieved, exists := store.GetByID("test-1")
	if !exists {
		t.Fatal("Project should exist after creation")
	}
	if retrieved.ID != project.ID {
		t.Fatalf("Expected ID %s, got %s", project.ID, retrieved.ID)
	}
	if retrieved.Name != project.Name {
		t.Fatalf("Expected Name %s, got %s", project.Name, retrieved.Name)
	}
}

func TestStore_GetByID(t *testing.T) {
	store := NewStore()
	project := &models.Project{
		ID:          "test-1",
		Name:        "Test Project",
		Description: "A test project",
		URL:         "https://test.com",
	}

	store.Create(project)

	// Test existing project
	retrieved, exists := store.GetByID("test-1")
	if !exists {
		t.Fatal("Project should exist")
	}
	if retrieved.ID != project.ID {
		t.Fatalf("Expected ID %s, got %s", project.ID, retrieved.ID)
	}

	// Test non-existent project
	_, exists = store.GetByID("non-existent")
	if exists {
		t.Fatal("Non-existent project should not exist")
	}
}

func TestStore_GetAll(t *testing.T) {
	store := NewStore()

	// Test empty store
	projects := store.GetAll()
	if len(projects) != 0 {
		t.Fatalf("Expected 0 projects, got %d", len(projects))
	}

	// Add multiple projects
	project1 := &models.Project{ID: "test-1", Name: "Project 1", URL: "https://test1.com"}
	project2 := &models.Project{ID: "test-2", Name: "Project 2", URL: "https://test2.com"}
	project3 := &models.Project{ID: "test-3", Name: "Project 3", URL: "https://test3.com"}

	store.Create(project1)
	store.Create(project2)
	store.Create(project3)

	projects = store.GetAll()
	if len(projects) != 3 {
		t.Fatalf("Expected 3 projects, got %d", len(projects))
	}

	// Verify all projects are present
	ids := make(map[string]bool)
	for _, p := range projects {
		ids[p.ID] = true
	}
	if !ids["test-1"] || !ids["test-2"] || !ids["test-3"] {
		t.Fatal("Not all projects are present in GetAll()")
	}
}

func TestStore_Update(t *testing.T) {
	store := NewStore()
	project := &models.Project{
		ID:          "test-1",
		Name:        "Test Project",
		Description: "A test project",
		URL:         "https://test.com",
	}

	store.Create(project)

	// Update existing project
	updated := &models.Project{
		ID:          "test-1",
		Name:        "Updated Project",
		Description: "An updated test project",
		URL:         "https://updated.com",
	}

	success := store.Update(updated)
	if !success {
		t.Fatal("Update should succeed for existing project")
	}

	retrieved, exists := store.GetByID("test-1")
	if !exists {
		t.Fatal("Project should still exist after update")
	}
	if retrieved.Name != "Updated Project" {
		t.Fatalf("Expected Name 'Updated Project', got %s", retrieved.Name)
	}

	// Try to update non-existent project
	nonExistent := &models.Project{ID: "non-existent", Name: "Non Existent"}
	success = store.Update(nonExistent)
	if success {
		t.Fatal("Update should fail for non-existent project")
	}
}

func TestStore_Delete(t *testing.T) {
	store := NewStore()
	project := &models.Project{
		ID:          "test-1",
		Name:        "Test Project",
		Description: "A test project",
		URL:         "https://test.com",
	}

	store.Create(project)

	// Delete existing project
	success := store.Delete("test-1")
	if !success {
		t.Fatal("Delete should succeed for existing project")
	}

	_, exists := store.GetByID("test-1")
	if exists {
		t.Fatal("Project should not exist after deletion")
	}

	if len(store.projects) != 0 {
		t.Fatalf("Expected 0 projects after deletion, got %d", len(store.projects))
	}

	// Try to delete non-existent project
	success = store.Delete("non-existent")
	if success {
		t.Fatal("Delete should fail for non-existent project")
	}
}

func TestStore_ConcurrentAccess(t *testing.T) {
	store := NewStore()
	done := make(chan bool)

	// Concurrent writes
	go func() {
		for i := 0; i < 100; i++ {
			project := &models.Project{
				ID:   "concurrent-1",
				Name: "Concurrent Project",
				URL:  "https://concurrent.com",
			}
			store.Create(project)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 100; i++ {
			store.GetByID("concurrent-1")
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 100; i++ {
			store.GetAll()
		}
		done <- true
	}()

	// Wait for all goroutines to complete
	<-done
	<-done
	<-done

	// Verify final state
	_, exists := store.GetByID("concurrent-1")
	if !exists {
		t.Fatal("Project should exist after concurrent operations")
	}
}
