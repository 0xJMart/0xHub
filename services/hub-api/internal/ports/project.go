package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/0xHub/hub-api/internal/domain"
)

// ProjectFilter declares pagination and filter options available when querying
// for projects.
type ProjectFilter struct {
	Tag      string
	Category string
	Search   string
	Limit    int
	Offset   int
}

// ProjectRepository exposes persistence operations over Project aggregates.
type ProjectRepository interface {
	ListProjects(ctx context.Context, filter ProjectFilter) ([]domain.Project, int, error)
	GetProjectBySlug(ctx context.Context, slug string) (*domain.Project, error)
	GetProjectByID(ctx context.Context, id uuid.UUID) (*domain.Project, error)
	CreateProject(ctx context.Context, project domain.Project) error
	UpdateProject(ctx context.Context, project domain.Project) error
	DeleteProject(ctx context.Context, id uuid.UUID) error
}
