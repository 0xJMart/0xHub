package domain_test

import (
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/0xHub/hub-api/internal/domain"
)

func TestProjectValidateSuccess(t *testing.T) {
	projectID := uuid.New()
	tagID := uuid.New()
	linkID := uuid.New()
	mediaID := uuid.New()

	project := domain.Project{
		ID:      projectID,
		Title:   "Homelab Inventory",
		Slug:    "homelab-inventory",
		Status:  domain.ProjectStatusActive,
		Summary: "Track gear",
		Tags: []domain.Tag{
			{
				ID:          tagID,
				Name:        "infra",
				DisplayName: "Infrastructure",
			},
		},
		Links: []domain.ProjectLink{
			{
				ID:        linkID,
				ProjectID: projectID,
				LinkType:  "repo",
				URL:       "https://github.com/example/project",
			},
		},
		Media: []domain.MediaAsset{
			{
				ID:               mediaID,
				StoragePath:      "projects/homelab-inventory/hero.png",
				OriginalFilename: "hero.png",
				Status:           domain.MediaStatusAvailable,
			},
		},
	}

	if err := project.Validate().AsError(); err != nil {
		t.Fatalf("expected project to be valid, got %v", err)
	}
}

func TestProjectValidateCollectsErrors(t *testing.T) {
	project := domain.Project{
		Status: "invalid",
		Tags: []domain.Tag{
			{},
		},
		Links: []domain.ProjectLink{
			{},
		},
		Media: []domain.MediaAsset{
			{
				Status:    "bad",
				SizeBytes: -100,
			},
		},
	}

	errs := project.Validate()
	if !errs.HasErrors() {
		t.Fatalf("expected validation errors, got none")
	}
	if len(errs) < 4 {
		t.Fatalf("expected multiple errors, got %d (%v)", len(errs), errs)
	}
}

func TestTagValidate(t *testing.T) {
	tag := domain.Tag{
		ID:          uuid.New(),
		Name:        "automation",
		DisplayName: "Automation",
	}

	if tag.Validate().HasErrors() {
		t.Fatal("expected tag to be valid")
	}
}

func TestTagValidateErrors(t *testing.T) {
	tag := domain.Tag{}
	errs := tag.Validate()
	if !errs.HasErrors() {
		t.Fatal("expected tag validation to fail")
	}
}

func TestMediaAssetValidate(t *testing.T) {
	asset := domain.MediaAsset{
		ID:               uuid.New(),
		StoragePath:      "projects/demo/file.txt",
		OriginalFilename: "file.txt",
		Status:           domain.MediaStatusPending,
		SizeBytes:        123,
	}
	if asset.Validate().HasErrors() {
		t.Fatal("expected media asset to be valid")
	}
}

func TestProjectLinkValidateErrors(t *testing.T) {
	link := domain.ProjectLink{}
	errs := link.Validate()
	if !errs.HasErrors() {
		t.Fatal("expected link validation to fail")
	}
}

func TestValidationErrorsFormatting(t *testing.T) {
	errs := domain.ValidationErrors{
		{Field: "title", Message: "cannot be empty"},
		{Field: "slug", Message: "invalid"},
	}
	msg := errs.Error()
	if !strings.Contains(msg, "title") || !strings.Contains(msg, "slug") {
		t.Fatalf("unexpected error message: %s", msg)
	}
}
