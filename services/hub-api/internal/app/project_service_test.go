package app

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/0xHub/hub-api/internal/domain"
	"github.com/0xHub/hub-api/internal/ports"
)

func TestProjectServiceListClampsPagination(t *testing.T) {
	ctx := context.Background()
	var captured ports.ProjectFilter
	repo := &projectRepoStub{
		listFunc: func(ctx context.Context, filter ports.ProjectFilter) ([]domain.Project, int, error) {
			captured = filter
			return []domain.Project{}, 0, nil
		},
	}
	service := NewProjectService(repo, newTestLogger())

	_, _, err := service.List(ctx, ports.ProjectFilter{
		Limit:  500,
		Offset: -10,
		Tag:    " infra ",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if captured.Limit != maxListLimit {
		t.Fatalf("expected limit to clamp to %d, got %d", maxListLimit, captured.Limit)
	}
	if captured.Offset != 0 {
		t.Fatalf("expected offset to clamp to 0, got %d", captured.Offset)
	}
	if captured.Tag != "infra" {
		t.Fatalf("expected tag to be trimmed, got %q", captured.Tag)
	}
}

func TestProjectServiceCreateValidatesProject(t *testing.T) {
	ctx := context.Background()
	repo := &projectRepoStub{}
	service := NewProjectService(repo, newTestLogger())

	err := service.Create(ctx, domain.Project{})
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	var valErrs domain.ValidationErrors
	if !errors.As(err, &valErrs) {
		t.Fatalf("expected validation errors, got %v", err)
	}
}

func TestProjectServiceGetBySlugRequiresValue(t *testing.T) {
	ctx := context.Background()
	repo := &projectRepoStub{}
	service := NewProjectService(repo, newTestLogger())

	if _, err := service.GetBySlug(ctx, ""); err == nil {
		t.Fatal("expected error when slug is empty")
	}
}

func TestProjectServiceDeleteRejectsNilID(t *testing.T) {
	ctx := context.Background()
	repo := &projectRepoStub{}
	service := NewProjectService(repo, newTestLogger())

	err := service.Delete(ctx, uuid.Nil)
	if err == nil {
		t.Fatal("expected error for nil UUID")
	}
	if !domain.IsNotFound(err) {
		t.Fatalf("expected not found error, got %v", err)
	}
}
