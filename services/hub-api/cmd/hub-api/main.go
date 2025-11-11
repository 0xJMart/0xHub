package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/0xHub/hub-api/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil {
		log.Fatalf("hub-api failed: %v", err)
	}
}

func run(ctx context.Context) error {
	cfg := app.LoadConfig()
	application := app.NewApplication(cfg, nil)
	return application.Start(ctx)
}
