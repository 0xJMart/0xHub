package store

import (
	"sync"

	"0xhub/backend/internal/models"
)

// Store is an in-memory store for projects
type Store struct {
	mu       sync.RWMutex
	projects map[string]*models.Project
}

// NewStore creates a new in-memory store
func NewStore() *Store {
	return &Store{
		projects: make(map[string]*models.Project),
	}
}

// GetAll returns all projects
func (s *Store) GetAll() []*models.Project {
	s.mu.RLock()
	defer s.mu.RUnlock()

	projects := make([]*models.Project, 0, len(s.projects))
	for _, p := range s.projects {
		projects = append(projects, p)
	}
	return projects
}

// GetByID returns a project by ID
func (s *Store) GetByID(id string) (*models.Project, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	project, exists := s.projects[id]
	return project, exists
}

// Create creates a new project
func (s *Store) Create(project *models.Project) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.projects[project.ID] = project
}

// Update updates an existing project
func (s *Store) Update(project *models.Project) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.projects[project.ID]; !exists {
		return false
	}
	s.projects[project.ID] = project
	return true
}

// Delete deletes a project by ID
func (s *Store) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.projects[id]; !exists {
		return false
	}
	delete(s.projects, id)
	return true
}

