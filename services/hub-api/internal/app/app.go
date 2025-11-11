package app

import "context"

// Application wires the Hub domain with external services.
type Application struct {
	// TODO: add dependencies (repositories, caches, etc.).
}

// Start boots the application lifecycle. It returns when the application
// shuts down or when an error is encountered during initialization.
func (a *Application) Start(ctx context.Context) error {
	_ = ctx
	// TODO: boot HTTP server, background workers, etc.
	return nil
}
