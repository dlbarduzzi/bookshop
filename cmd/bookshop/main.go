package main

import (
	"context"
	"os"

	"github.com/dlbarduzzi/bookshop/internal/bookshop"
	"github.com/dlbarduzzi/bookshop/internal/logging"
	"github.com/dlbarduzzi/bookshop/internal/registry"
	"github.com/dlbarduzzi/bookshop/internal/server"
	"github.com/spf13/viper"
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

	reg, err := registry.NewRegistry()
	if err != nil {
		return err
	}

	bookshopConfig := getBookshopConfig(reg)
	log.Info("database connection established")

	srv, err := server.NewServer(bookshopConfig.Port)
	if err != nil {
		return err
	}

	return srv.Start(ctx, nil)
}

func getBookshopConfig(v *viper.Viper) *bookshop.Config {
	return &bookshop.Config{
		Port: v.GetInt("BS_APP_PORT"),
	}
}
