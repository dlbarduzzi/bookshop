package main

import (
	"context"
	"os"

	"github.com/dlbarduzzi/bookshop/internal/logging"
)

func main() {
	log := logging.NewLoggerFromEnv()

	ctx := context.Background()
	ctx = logging.LoggerWithContext(ctx, log)

	if err := start(ctx); err != nil {
		log.Error(err.Error())
		os.Exit(2)
	}
}

func start(ctx context.Context) error {
	log := logging.LoggerFromContext(ctx)
	log.Info("Welcome to my bookshop!")
	return nil
}
