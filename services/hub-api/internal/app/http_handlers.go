package app

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/0xHub/hub-api/internal/domain"
	"github.com/0xHub/hub-api/internal/ports"
)

func (a *Application) registerRoutes(router *gin.Engine) {
	router.GET("/healthz", a.handleHealthz)
	router.GET("/readyz", a.handleReadyz)

	router.GET("/projects", a.handleListProjects)
	router.GET("/projects/:slug", a.handleGetProjectBySlug)

	router.GET("/tags", a.handleListTags)
}

func (a *Application) handleHealthz(c *gin.Context) {
	c.JSON(http.StatusOK, healthResponse{
		Status:    "ok",
		Timestamp: time.Now().UTC(),
	})
}

func (a *Application) handleReadyz(c *gin.Context) {
	if err := a.healthSvc.Readiness(c.Request.Context()); err != nil {
		a.respondError(c, http.StatusServiceUnavailable, "dependencies unavailable", err.Error())
		return
	}
	c.JSON(http.StatusOK, healthResponse{
		Status:    "ready",
		Timestamp: time.Now().UTC(),
	})
}

type listProjectsQuery struct {
	Tag      string `form:"tag"`
	Category string `form:"category"`
	Search   string `form:"search"`
	Limit    int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset   int    `form:"offset" binding:"omitempty,min=0"`
}

func (a *Application) handleListProjects(c *gin.Context) {
	var query listProjectsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		a.respondError(c, http.StatusBadRequest, "invalid query parameters", err.Error())
		return
	}

	filter := ports.ProjectFilter{
		Tag:      query.Tag,
		Category: query.Category,
		Search:   query.Search,
		Limit:    query.Limit,
		Offset:   query.Offset,
	}

	projects, total, err := a.projectSvc.List(c.Request.Context(), filter)
	if err != nil {
		a.respondServiceError(c, err)
		return
	}

	items := make([]projectSummaryResponse, 0, len(projects))
	for _, project := range projects {
		items = append(items, toProjectSummary(project))
	}

	c.JSON(http.StatusOK, projectListResponse{
		Items: items,
		Total: total,
	})
}

func (a *Application) handleGetProjectBySlug(c *gin.Context) {
	slug := c.Param("slug")
	project, err := a.projectSvc.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		a.respondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, toProjectDetail(*project))
}

func (a *Application) handleListTags(c *gin.Context) {
	tags, err := a.tagSvc.List(c.Request.Context())
	if err != nil {
		a.respondServiceError(c, err)
		return
	}

	items := make([]tagResponse, 0, len(tags))
	for _, tag := range tags {
		items = append(items, toTagResponse(tag))
	}

	c.JSON(http.StatusOK, items)
}

func (a *Application) respondServiceError(c *gin.Context, err error) {
	switch e := err.(type) {
	case domain.ValidationErrors:
		a.respondValidationErrors(c, e)
	default:
		if domain.IsNotFound(err) {
			a.respondError(c, http.StatusNotFound, "resource not found", err.Error())
			return
		}
		a.respondError(c, http.StatusInternalServerError, "internal error", err.Error())
	}
}

type healthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

type errorResponse struct {
	Error         string `json:"error"`
	Detail        string `json:"detail,omitempty"`
	CorrelationID string `json:"correlationId,omitempty"`
}

func (a *Application) respondError(c *gin.Context, status int, message, detail string) {
	correlationID := uuid.NewString()
	if status >= http.StatusInternalServerError {
		a.logger.Error("request failed",
			slog.String("correlationId", correlationID),
			slog.String("detail", detail),
			slog.Int("status", status),
			slog.String("path", c.FullPath()),
		)
	}
	c.JSON(status, errorResponse{
		Error:         message,
		Detail:        detail,
		CorrelationID: correlationID,
	})
}

func (a *Application) respondValidationErrors(c *gin.Context, errs domain.ValidationErrors) {
	a.respondError(c, http.StatusBadRequest, "validation failed", errs.Error())
}

type projectListResponse struct {
	Items []projectSummaryResponse `json:"items"`
	Total int                      `json:"total"`
}

type projectSummaryResponse struct {
	ID        string            `json:"id"`
	Title     string            `json:"title"`
	Slug      string            `json:"slug"`
	Summary   string            `json:"summary,omitempty"`
	Status    string            `json:"status"`
	Category  *categoryResponse `json:"category,omitempty"`
	Tags      []tagResponse     `json:"tags"`
	HeroMedia *mediaResponse    `json:"heroMedia,omitempty"`
}

type projectDetailResponse struct {
	projectSummaryResponse
	Description string          `json:"description,omitempty"`
	Links       []linkResponse  `json:"links"`
	Media       []mediaResponse `json:"media"`
}

type categoryResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description,omitempty"`
	SortOrder   int    `json:"sortOrder,omitempty"`
}

type tagResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Color       string `json:"color,omitempty"`
}

type linkResponse struct {
	ID    string `json:"id"`
	Type  string `json:"linkType"`
	URL   string `json:"url"`
	Label string `json:"label,omitempty"`
}

type mediaResponse struct {
	ID               string     `json:"id"`
	ProjectID        *string    `json:"projectId,omitempty"`
	StoragePath      string     `json:"storagePath"`
	OriginalFilename string     `json:"originalFilename"`
	MimeType         string     `json:"mimeType,omitempty"`
	SizeBytes        int64      `json:"sizeBytes,omitempty"`
	Description      string     `json:"description,omitempty"`
	Status           string     `json:"status"`
	CreatedAt        *time.Time `json:"createdAt,omitempty"`
	UpdatedAt        *time.Time `json:"updatedAt,omitempty"`
	ExpiresAt        *time.Time `json:"expiresAt,omitempty"`
}

func toProjectSummary(project domain.Project) projectSummaryResponse {
	resp := projectSummaryResponse{
		ID:      project.ID.String(),
		Title:   project.Title,
		Slug:    project.Slug,
		Summary: project.Summary,
		Status:  string(project.Status),
		Tags:    make([]tagResponse, 0, len(project.Tags)),
	}

	if project.Category != nil {
		resp.Category = &categoryResponse{
			ID:          project.Category.ID.String(),
			Name:        project.Category.Name,
			Slug:        project.Category.Slug,
			Description: project.Category.Description,
			SortOrder:   project.Category.SortOrder,
		}
	}

	for _, tag := range project.Tags {
		resp.Tags = append(resp.Tags, toTagResponse(tag))
	}

	if project.HeroMedia != nil {
		hero := toMediaResponse(*project.HeroMedia)
		resp.HeroMedia = &hero
	}

	return resp
}

func toProjectDetail(project domain.Project) projectDetailResponse {
	detail := projectDetailResponse{
		projectSummaryResponse: toProjectSummary(project),
		Description:            project.Description,
		Links:                  make([]linkResponse, 0, len(project.Links)),
		Media:                  make([]mediaResponse, 0, len(project.Media)),
	}

	for _, link := range project.Links {
		detail.Links = append(detail.Links, linkResponse{
			ID:    link.ID.String(),
			Type:  link.LinkType,
			URL:   link.URL,
			Label: link.Label,
		})
	}

	for _, media := range project.Media {
		detail.Media = append(detail.Media, toMediaResponse(media))
	}

	return detail
}

func toTagResponse(tag domain.Tag) tagResponse {
	return tagResponse{
		ID:          tag.ID.String(),
		Name:        tag.Name,
		DisplayName: tag.DisplayName,
		Color:       tag.Color,
	}
}

func toMediaResponse(media domain.MediaAsset) mediaResponse {
	resp := mediaResponse{
		ID:               media.ID.String(),
		StoragePath:      media.StoragePath,
		OriginalFilename: media.OriginalFilename,
		MimeType:         media.MimeType,
		SizeBytes:        media.SizeBytes,
		Description:      media.Description,
		Status:           string(media.Status),
	}

	if media.ProjectID != nil {
		projectID := media.ProjectID.String()
		resp.ProjectID = &projectID
	}
	if !media.CreatedAt.IsZero() {
		resp.CreatedAt = &media.CreatedAt
	}
	if !media.UpdatedAt.IsZero() {
		resp.UpdatedAt = &media.UpdatedAt
	}
	if media.ExpiresAt != nil {
		resp.ExpiresAt = media.ExpiresAt
	}

	return resp
}
