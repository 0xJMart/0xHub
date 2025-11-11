package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// ProjectLink represents an external reference associated with a project.
type ProjectLink struct {
	ID        uuid.UUID
	ProjectID uuid.UUID
	LinkType  string
	URL       string
	Label     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Validate checks the link metadata for completeness.
func (l ProjectLink) Validate() ValidationErrors {
	var errs ValidationErrors

	if l.ID == uuid.Nil {
		appendError(&errs, "id", "must be a valid UUID")
	}

	if l.ProjectID == uuid.Nil {
		appendError(&errs, "projectId", "must be a valid UUID")
	}

	if strings.TrimSpace(l.LinkType) == "" {
		appendError(&errs, "linkType", "cannot be empty")
	}

	validateURL("url", l.URL, &errs)

	return errs
}
