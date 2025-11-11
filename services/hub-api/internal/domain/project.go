package domain

import (
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ProjectStatus represents the lifecycle state of a project.
type ProjectStatus string

const (
	ProjectStatusDraft      ProjectStatus = "draft"
	ProjectStatusActive     ProjectStatus = "active"
	ProjectStatusArchived   ProjectStatus = "archived"
	ProjectStatusDeprecated ProjectStatus = "deprecated"
)

// Project captures the main content surfaced by the Hub.
type Project struct {
	ID          uuid.UUID
	CategoryID  *uuid.UUID
	Category    *Category
	Title       string
	Slug        string
	Summary     string
	Description string
	Status      ProjectStatus
	HeroMediaID *uuid.UUID
	HeroMedia   *MediaAsset
	Tags        []Tag
	Links       []ProjectLink
	Media       []MediaAsset
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ArchivedAt  *time.Time
}

// Validate ensures the project meets homelab constraints.
func (p Project) Validate() ValidationErrors {
	var errs ValidationErrors

	if p.ID == uuid.Nil {
		appendError(&errs, "id", "must be a valid UUID")
	}

	title := strings.TrimSpace(p.Title)
	if title == "" {
		appendError(&errs, "title", "cannot be empty")
	} else if len(title) > maxTitleLength {
		appendError(&errs, "title", "must be 160 characters or fewer")
	}

	validateSlug("slug", p.Slug, &errs)

	if summary := strings.TrimSpace(p.Summary); len(summary) > maxSummaryLength {
		appendError(&errs, "summary", "must be 512 characters or fewer")
	}

	switch p.Status {
	case ProjectStatusDraft, ProjectStatusActive, ProjectStatusArchived, ProjectStatusDeprecated:
	default:
		appendError(&errs, "status", "must be draft, active, archived, or deprecated")
	}

	for i, tag := range p.Tags {
		if tagErrs := tag.Validate(); tagErrs.HasErrors() {
			errs = append(errs, prefixErrors("tags["+strconv.Itoa(i)+"]", tagErrs)...)
		}
	}

	for i, link := range p.Links {
		if linkErrs := link.Validate(); linkErrs.HasErrors() {
			errs = append(errs, prefixErrors("links["+strconv.Itoa(i)+"]", linkErrs)...)
		}
	}

	for i, media := range p.Media {
		if mediaErrs := media.Validate(); mediaErrs.HasErrors() {
			errs = append(errs, prefixErrors("media["+strconv.Itoa(i)+"]", mediaErrs)...)
		}
	}

	if p.Category != nil {
		if categoryErrs := p.Category.Validate(); categoryErrs.HasErrors() {
			errs = append(errs, prefixErrors("category", categoryErrs)...)
		}
	}

	if p.HeroMedia != nil {
		if mediaErrs := p.HeroMedia.Validate(); mediaErrs.HasErrors() {
			errs = append(errs, prefixErrors("heroMedia", mediaErrs)...)
		}
	}

	return errs
}
