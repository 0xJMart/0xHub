package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/0xHub/hub-api/internal/adapters/postgres"
)

// Application wires the Hub domain with external services.
type Application struct {
	cfg       Config
	logger    *slog.Logger
	validator *validator.Validate

	projectSvc *ProjectService
	tagSvc     *TagService
	healthSvc  *HealthService

	store  *postgres.Store
	server *http.Server
}

// NewApplication constructs an application instance.
func NewApplication(cfg Config, logger *slog.Logger) *Application {
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return &Application{
		cfg:    cfg,
		logger: logger,
	}
}

// Start boots the HTTP server and blocks until shutdown.
func (a *Application) Start(ctx context.Context) error {
	if err := a.initialise(ctx); err != nil {
		return err
	}
	defer a.cleanup(context.Background())

	errCh := make(chan error, 1)
	go func() {
		a.logger.Info("hub-api listening", slog.String("addr", a.cfg.HTTPAddr))
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), a.cfg.ShutdownTimeout)
		defer cancel()
		if err := a.server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("http shutdown: %w", err)
		}
		<-errCh
	case err := <-errCh:
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *Application) initialise(ctx context.Context) error {
	if a.cfg.DatabaseURL == "" {
		return fmt.Errorf("application: database URL must be provided")
	}

	store, err := postgres.New(ctx, postgres.Config{
		URL:             a.cfg.DatabaseURL,
		MaxConns:        a.cfg.DatabaseMaxConns,
		MinConns:        a.cfg.DatabaseMinConns,
		MaxConnLifetime: a.cfg.DatabaseMaxConnLifetime,
	})
	if err != nil {
		return err
	}
	a.store = store

	a.validator = validator.New()
	a.projectSvc = NewProjectService(store, a.logger)
	a.tagSvc = NewTagService(store, a.logger)
	a.healthSvc = NewHealthService(store, a.logger)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery(), requestLogger(a.logger))
	a.registerRoutes(router)

	a.server = &http.Server{
		Addr:         a.cfg.HTTPAddr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return nil
}

func (a *Application) cleanup(ctx context.Context) {
	_ = ctx
	if a.store != nil {
		a.store.Close()
		a.store = nil
	}
}
