package main

import (
	"context"
	"os"

	"github.com/dlbarduzzi/guestbook/internal/guestbook"
	"github.com/dlbarduzzi/guestbook/internal/logging"
	"github.com/dlbarduzzi/guestbook/internal/registry"
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

	reg, err := registry.NewRegistry()
	if err != nil {
		return err
	}

	port := reg.GetInt("GB_PORT")

	srv := server.NewServer(port, logger)
	app := guestbook.NewGuesbook(logger)

	srv.RunBeforeShutdown(func() {
		app.Shutdown()
	})

	return srv.Start(ctx, app.Routes())
}
