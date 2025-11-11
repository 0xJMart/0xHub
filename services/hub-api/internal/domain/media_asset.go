package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// MediaStatus captures the lifecycle state of a media asset.
type MediaStatus string

const (
	MediaStatusPending   MediaStatus = "pending"
	MediaStatusAvailable MediaStatus = "available"
	MediaStatusErrored   MediaStatus = "errored"
)

// MediaAsset represents a stored file related to a project.
type MediaAsset struct {
	ID               uuid.UUID
	ProjectID        *uuid.UUID
	StoragePath      string
	OriginalFilename string
	MimeType         string
	SizeBytes        int64
	Description      string
	Status           MediaStatus
	CreatedAt        time.Time
	UpdatedAt        time.Time
	ExpiresAt        *time.Time
}

// Validate ensures the asset references and metadata are valid.
func (m MediaAsset) Validate() ValidationErrors {
	var errs ValidationErrors

	if m.ID == uuid.Nil {
		appendError(&errs, "id", "must be a valid UUID")
	}

	if strings.TrimSpace(m.StoragePath) == "" {
		appendError(&errs, "storagePath", "cannot be empty")
	}

	if strings.TrimSpace(m.OriginalFilename) == "" {
		appendError(&errs, "originalFilename", "cannot be empty")
	}

	switch m.Status {
	case MediaStatusPending, MediaStatusAvailable, MediaStatusErrored:
	default:
		appendError(&errs, "status", "must be pending, available, or errored")
	}

	if m.SizeBytes < 0 {
		appendError(&errs, "sizeBytes", "cannot be negative")
	}

	return errs
}
