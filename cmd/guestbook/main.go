package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/dlbarduzzi/guestbook/internal/guestbook"
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

	srv := server.NewServer(8080, logger)
	app := guestbook.NewGuesbook(logger)

	app.FakeHandler()

	srv.RunBeforeShutdown(func() {
		fmt.Println("Start before shutdown...")
		time.Sleep(time.Second * 1)
		fmt.Println("End before shutdown!")
	})

	srv.RunBeforeShutdown(func() {
		app.Shutdown()
	})

	return srv.Start(ctx, nil)
}
