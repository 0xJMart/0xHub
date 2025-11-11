package app

import (
	"context"
	"io"
	"log/slog"

	"github.com/google/uuid"

	"github.com/0xHub/hub-api/internal/domain"
	"github.com/0xHub/hub-api/internal/ports"
)

type projectRepoStub struct {
	listFunc      func(ctx context.Context, filter ports.ProjectFilter) ([]domain.Project, int, error)
	getBySlugFunc func(ctx context.Context, slug string) (*domain.Project, error)
	getByIDFunc   func(ctx context.Context, id uuid.UUID) (*domain.Project, error)
	createFunc    func(ctx context.Context, project domain.Project) error
	updateFunc    func(ctx context.Context, project domain.Project) error
	deleteFunc    func(ctx context.Context, id uuid.UUID) error
}

func (s *projectRepoStub) ListProjects(ctx context.Context, filter ports.ProjectFilter) ([]domain.Project, int, error) {
	if s.listFunc != nil {
		return s.listFunc(ctx, filter)
	}
	return nil, 0, nil
}

func (s *projectRepoStub) GetProjectBySlug(ctx context.Context, slug string) (*domain.Project, error) {
	if s.getBySlugFunc != nil {
		return s.getBySlugFunc(ctx, slug)
	}
	return nil, domain.ErrNotFound("project")
}

func (s *projectRepoStub) GetProjectByID(ctx context.Context, id uuid.UUID) (*domain.Project, error) {
	if s.getByIDFunc != nil {
		return s.getByIDFunc(ctx, id)
	}
	return nil, domain.ErrNotFound("project")
}

func (s *projectRepoStub) CreateProject(ctx context.Context, project domain.Project) error {
	if s.createFunc != nil {
		return s.createFunc(ctx, project)
	}
	return nil
}

func (s *projectRepoStub) UpdateProject(ctx context.Context, project domain.Project) error {
	if s.updateFunc != nil {
		return s.updateFunc(ctx, project)
	}
	return nil
}

func (s *projectRepoStub) DeleteProject(ctx context.Context, id uuid.UUID) error {
	if s.deleteFunc != nil {
		return s.deleteFunc(ctx, id)
	}
	return nil
}

type tagRepoStub struct {
	listFunc func(ctx context.Context) ([]domain.Tag, error)
}

func (s *tagRepoStub) ListTags(ctx context.Context) ([]domain.Tag, error) {
	if s.listFunc != nil {
		return s.listFunc(ctx)
	}
	return nil, nil
}

type healthRepoStub struct {
	pingFunc func(ctx context.Context) error
}

func (s *healthRepoStub) Ping(ctx context.Context) error {
	if s.pingFunc != nil {
		return s.pingFunc(ctx)
	}
	return nil
}

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))
}
