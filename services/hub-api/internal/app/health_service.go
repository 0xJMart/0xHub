package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/0xHub/hub-api/internal/ports"
)

// HealthService exposes health and readiness checks.
type HealthService struct {
	repo   ports.HealthRepository
	logger *slog.Logger
}

// NewHealthService constructs a HealthService.
func NewHealthService(repo ports.HealthRepository, logger *slog.Logger) *HealthService {
	if logger == nil {
		logger = slog.Default()
	}
	return &HealthService{
		repo:   repo,
		logger: logger,
	}
}

// Liveness returns nil when the service is alive.
func (s *HealthService) Liveness(_ context.Context) error {
	return nil
}

// Readiness checks dependent resources.
func (s *HealthService) Readiness(ctx context.Context) error {
	if s.repo == nil {
		return fmt.Errorf("health service: repository is not configured")
	}
	if err := s.repo.Ping(ctx); err != nil {
		s.logger.ErrorContext(ctx, "readiness probe failed", slog.String("error", err.Error()))
		return err
	}
	return nil
}
