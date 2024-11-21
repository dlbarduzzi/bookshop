package main

import (
	"context"
	"os"

	"github.com/spf13/viper"

	"github.com/dlbarduzzi/bookshop/internal/bookshop"
	"github.com/dlbarduzzi/bookshop/internal/database"
	"github.com/dlbarduzzi/bookshop/internal/logging"
	"github.com/dlbarduzzi/bookshop/internal/registry"
	"github.com/dlbarduzzi/bookshop/internal/server"
)

func main() {
	logger := logging.NewLoggerFromEnv().With("app", "bookshop")

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

	port := reg.GetInt("APP_PORT")
	dbConfig := setDatabaseConfig(reg)

	db, err := database.NewDatabase(dbConfig)
	if err != nil {
		return err
	}

	defer db.Close()
	logger.Info("database connection established")

	srv := server.NewServer(port, logger)
	app := bookshop.NewBookshop(logger)

	srv.RunBeforeShutdown(func() {
		app.Shutdown()
	})

	return srv.Start(ctx, app.Routes())
}

func setDatabaseConfig(v *viper.Viper) *database.Config {
	return &database.Config{
		ConnectionURL:   v.GetString("DB_CONNECTION_URL"),
		MaxOpenConns:    v.GetInt("DB_MAX_OPEN_CONNS"),
		MaxIdleConns:    v.GetInt("DB_MAX_IDLE_CONNS"),
		ConnMaxIdleTime: v.GetDuration("DB_CONN_MAX_IDLE_TIME"),
	}
}
