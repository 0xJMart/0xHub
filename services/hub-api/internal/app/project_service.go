package app

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"

	"github.com/0xHub/hub-api/internal/domain"
	"github.com/0xHub/hub-api/internal/ports"
)

const (
	defaultListLimit = 20
	maxListLimit     = 100
)

// ProjectService coordinates domain operations around projects.
type ProjectService struct {
	repo   ports.ProjectRepository
	logger *slog.Logger
}

// NewProjectService constructs a ProjectService instance.
func NewProjectService(repo ports.ProjectRepository, logger *slog.Logger) *ProjectService {
	if logger == nil {
		logger = slog.Default()
	}
	return &ProjectService{
		repo:   repo,
		logger: logger,
	}
}

// List returns paginated projects matching the provided filters.
func (s *ProjectService) List(ctx context.Context, filter ports.ProjectFilter) ([]domain.Project, int, error) {
	if s.repo == nil {
		return nil, 0, fmt.Errorf("project service: repository is not configured")
	}

	filter.Tag = strings.TrimSpace(filter.Tag)
	filter.Category = strings.TrimSpace(filter.Category)
	filter.Search = strings.TrimSpace(filter.Search)

	if filter.Limit <= 0 {
		filter.Limit = defaultListLimit
	}
	if filter.Limit > maxListLimit {
		filter.Limit = maxListLimit
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	projects, total, err := s.repo.ListProjects(ctx, filter)
	if err != nil {
		s.logger.ErrorContext(ctx, "list projects failed", slog.String("error", err.Error()))
		return nil, 0, err
	}

	return projects, total, nil
}

// GetBySlug fetches a project aggregate by slug.
func (s *ProjectService) GetBySlug(ctx context.Context, slug string) (*domain.Project, error) {
	if s.repo == nil {
		return nil, fmt.Errorf("project service: repository is not configured")
	}
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return nil, fmt.Errorf("project service: slug cannot be empty")
	}

	project, err := s.repo.GetProjectBySlug(ctx, slug)
	if err != nil {
		if domain.IsNotFound(err) {
			return nil, err
		}
		s.logger.ErrorContext(ctx, "get project failed", slog.String("slug", slug), slog.String("error", err.Error()))
		return nil, err
	}

	return project, nil
}

// Create persists a new project aggregate after validation.
func (s *ProjectService) Create(ctx context.Context, project domain.Project) error {
	if s.repo == nil {
		return fmt.Errorf("project service: repository is not configured")
	}
	if project.ID == uuid.Nil {
		return domain.ValidationErrors{{Field: "id", Message: "must be a valid UUID"}}
	}
	if errs := project.Validate(); errs.HasErrors() {
		return errs
	}
	if err := s.repo.CreateProject(ctx, project); err != nil {
		s.logger.ErrorContext(ctx, "create project failed", slog.String("projectId", project.ID.String()), slog.String("error", err.Error()))
		return err
	}
	return nil
}

// Update mutates an existing project aggregate.
func (s *ProjectService) Update(ctx context.Context, project domain.Project) error {
	if s.repo == nil {
		return fmt.Errorf("project service: repository is not configured")
	}
	if project.ID == uuid.Nil {
		return domain.ErrNotFound("project")
	}
	if errs := project.Validate(); errs.HasErrors() {
		return errs
	}
	if err := s.repo.UpdateProject(ctx, project); err != nil {
		if domain.IsNotFound(err) {
			return err
		}
		s.logger.ErrorContext(ctx, "update project failed", slog.String("projectId", project.ID.String()), slog.String("error", err.Error()))
		return err
	}
	return nil
}

// Delete removes a project by identifier.
func (s *ProjectService) Delete(ctx context.Context, id uuid.UUID) error {
	if s.repo == nil {
		return fmt.Errorf("project service: repository is not configured")
	}
	if id == uuid.Nil {
		return domain.ErrNotFound("project")
	}
	if err := s.repo.DeleteProject(ctx, id); err != nil {
		if domain.IsNotFound(err) {
			return err
		}
		s.logger.ErrorContext(ctx, "delete project failed", slog.String("projectId", id.String()), slog.String("error", err.Error()))
		return err
	}
	return nil
}
