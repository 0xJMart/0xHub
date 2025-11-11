package ports

import "context"

// HealthRepository exposes simple readiness probes for infrastructure
// dependencies.
type HealthRepository interface {
	Ping(ctx context.Context) error
}
