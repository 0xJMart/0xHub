package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/0xHub/hub-api/internal/domain"
	"github.com/0xHub/hub-api/internal/ports"
)

// TagService coordinates tag read operations.
type TagService struct {
	repo   ports.TagRepository
	logger *slog.Logger
}

// NewTagService constructs a TagService instance.
func NewTagService(repo ports.TagRepository, logger *slog.Logger) *TagService {
	if logger == nil {
		logger = slog.Default()
	}
	return &TagService{
		repo:   repo,
		logger: logger,
	}
}

// List returns all known tags.
func (s *TagService) List(ctx context.Context) ([]domain.Tag, error) {
	if s.repo == nil {
		return nil, fmt.Errorf("tag service: repository is not configured")
	}
	tags, err := s.repo.ListTags(ctx)
	if err != nil {
		s.logger.ErrorContext(ctx, "list tags failed", slog.String("error", err.Error()))
		return nil, err
	}
	return tags, nil
}
