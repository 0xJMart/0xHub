package ports

import (
	"context"

	"github.com/0xHub/hub-api/internal/domain"
)

// TagRepository exposes read operations for tag aggregates.
type TagRepository interface {
	ListTags(ctx context.Context) ([]domain.Tag, error)
}
