package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// Tag represents a free-form label that can be applied to projects.
type Tag struct {
	ID          uuid.UUID
	Name        string
	DisplayName string
	Color       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Validate ensures the tag has the required identifiers.
func (t Tag) Validate() ValidationErrors {
	var errs ValidationErrors

	if t.ID == uuid.Nil {
		appendError(&errs, "id", "must be a valid UUID")
	}

	name := strings.TrimSpace(t.Name)
	if name == "" {
		appendError(&errs, "name", "cannot be empty")
	} else if len(name) > maxTagLength {
		appendError(&errs, "name", "must be 64 characters or fewer")
	}

	display := strings.TrimSpace(t.DisplayName)
	if display == "" {
		appendError(&errs, "displayName", "cannot be empty")
	}

	return errs
}
