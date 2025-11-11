package app

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/0xHub/hub-api/internal/domain"
	"github.com/0xHub/hub-api/internal/ports"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestHandleListProjectsReturnsSummary(t *testing.T) {
	projectID := uuid.New()
	tagID := uuid.New()
	categoryID := uuid.New()
	heroID := uuid.New()

	repo := &projectRepoStub{
		listFunc: func(_ context.Context, _ ports.ProjectFilter) ([]domain.Project, int, error) {
			return []domain.Project{
				{
					ID:      projectID,
					Title:   "Homelab Inventory",
					Slug:    "homelab-inventory",
					Status:  domain.ProjectStatusActive,
					Summary: "Track gear",
					Category: &domain.Category{
						ID:          categoryID,
						Name:        "Automation",
						Slug:        "automation",
						Description: "Automation projects",
						SortOrder:   5,
					},
					Tags: []domain.Tag{
						{
							ID:          tagID,
							Name:        "infra",
							DisplayName: "Infrastructure",
						},
					},
					HeroMedia: &domain.MediaAsset{
						ID:               heroID,
						StoragePath:      "projects/homelab-inventory/hero.png",
						OriginalFilename: "hero.png",
						Status:           domain.MediaStatusAvailable,
					},
				},
			}, 1, nil
		},
	}

	app := &Application{
		logger:     newTestLogger(),
		projectSvc: NewProjectService(repo, newTestLogger()),
		tagSvc:     NewTagService(&tagRepoStub{}, newTestLogger()),
		healthSvc:  NewHealthService(&healthRepoStub{}, newTestLogger()),
	}

	router := gin.New()
	router.Use(gin.Recovery())
	app.registerRoutes(router)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/projects", nil)
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var body projectListResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(body.Items) != 1 {
		t.Fatalf("expected 1 project, got %d", len(body.Items))
	}

	item := body.Items[0]
	if item.ID != projectID.String() {
		t.Fatalf("expected project ID %s, got %s", projectID, item.ID)
	}
	if item.Category == nil || item.Category.ID != categoryID.String() {
		t.Fatalf("expected category to be populated")
	}
	if len(item.Tags) != 1 || item.Tags[0].ID != tagID.String() {
		t.Fatalf("expected tag to be included")
	}
	if item.HeroMedia == nil || item.HeroMedia.ID != heroID.String() {
		t.Fatalf("expected hero media to be included")
	}
}

func TestHandleListProjectsValidationError(t *testing.T) {
	app := &Application{
		logger:     newTestLogger(),
		projectSvc: NewProjectService(&projectRepoStub{}, newTestLogger()),
		tagSvc:     NewTagService(&tagRepoStub{}, newTestLogger()),
		healthSvc:  NewHealthService(&healthRepoStub{}, newTestLogger()),
	}

	router := gin.New()
	router.Use(gin.Recovery())
	app.registerRoutes(router)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/projects?limit=abc", nil)
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestHandleReadyzUnavailable(t *testing.T) {
	app := &Application{
		logger:     newTestLogger(),
		projectSvc: NewProjectService(&projectRepoStub{}, newTestLogger()),
		tagSvc:     NewTagService(&tagRepoStub{}, newTestLogger()),
		healthSvc: NewHealthService(&healthRepoStub{
			pingFunc: func(ctx context.Context) error {
				return errors.New("database down")
			},
		}, newTestLogger()),
	}

	router := gin.New()
	router.Use(gin.Recovery())
	app.registerRoutes(router)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status 503, got %d", rec.Code)
	}
}

func TestHandleGetProjectNotFound(t *testing.T) {
	app := &Application{
		logger: newTestLogger(),
		projectSvc: NewProjectService(&projectRepoStub{
			getBySlugFunc: func(ctx context.Context, slug string) (*domain.Project, error) {
				return nil, domain.ErrNotFound("project")
			},
		}, newTestLogger()),
		tagSvc:    NewTagService(&tagRepoStub{}, newTestLogger()),
		healthSvc: NewHealthService(&healthRepoStub{}, newTestLogger()),
	}

	router := gin.New()
	router.Use(gin.Recovery())
	app.registerRoutes(router)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/projects/unknown", nil)
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
}
