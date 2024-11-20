package main

import (
	"context"
	"os"

	"github.com/spf13/viper"

	"github.com/dlbarduzzi/guestbook/internal/database"
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

	port := reg.GetInt("GB_APP_PORT")
	dbConfig := setDatabaseConfig(reg)

	db, err := database.NewDatabase(dbConfig)
	if err != nil {
		return err
	}

	defer db.Close()
	logger.Info("database connection established")

	srv := server.NewServer(port, logger)
	app := guestbook.NewGuesbook(logger)

	srv.RunBeforeShutdown(func() {
		app.Shutdown()
	})

	return srv.Start(ctx, app.Routes())
}

func setDatabaseConfig(v *viper.Viper) *database.Config {
	return &database.Config{
		ConnectionURL:   v.GetString("GB_DATABASE_CONNECTION_URL"),
		MaxOpenConns:    v.GetInt("GB_DATABASE_MAX_OPEN_CONNS"),
		MaxIdleConns:    v.GetInt("GB_DATABASE_MAX_IDLE_CONNS"),
		ConnMaxIdleTime: v.GetDuration("GB_DATABASE_CONN_MAX_IDLE_TIME"),
	}
}
