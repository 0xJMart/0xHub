package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// Category groups projects under a common theme.
type Category struct {
	ID          uuid.UUID
	Name        string
	Slug        string
	Description string
	SortOrder   int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Validate ensures the category contains the required metadata.
func (c Category) Validate() ValidationErrors {
	var errs ValidationErrors

	if c.ID == uuid.Nil {
		appendError(&errs, "id", "must be a valid UUID")
	}

	if strings.TrimSpace(c.Name) == "" {
		appendError(&errs, "name", "cannot be empty")
	}

	validateSlug("slug", c.Slug, &errs)

	return errs
}
