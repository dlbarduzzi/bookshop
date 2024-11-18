package main

import (
	"context"
	"os"

	"github.com/dlbarduzzi/guestbook/internal/logging"
	"github.com/dlbarduzzi/guestbook/internal/server"
)

func main() {
	logger := logging.NewLoggerFromEnv().With("app", "guestbook")

	ctx := context.Background()
	ctx = logging.LoggerWithContext(ctx, logger)

	if err := start(ctx); err != nil {
		logger.Error(err.Error())
		os.Exit(2)
	}
}

func start(ctx context.Context) error {
	logger := logging.LoggerFromContext(ctx)
	logger.Info("app running...")

	srv := server.NewServer(8000)

	return srv.Start(ctx, nil)
}
