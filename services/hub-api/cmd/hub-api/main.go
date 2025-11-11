package main

import (
	"context"
	"log"

	"github.com/0xHub/hub-api/internal/app"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatalf("hub-api failed: %v", err)
	}
}

func run(ctx context.Context) error {
	application := &app.Application{}
	return application.Start(ctx)
}
