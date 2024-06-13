package main

import (
	"context"
	"os"

	"github.com/spf13/viper"

	"github.com/dlbarduzzi/bookshop/internal/bookshop"
	"github.com/dlbarduzzi/bookshop/internal/logging"
	"github.com/dlbarduzzi/bookshop/internal/registry"
	"github.com/dlbarduzzi/bookshop/internal/server"
)

func main() {
	log := logging.NewLoggerFromEnv()

	ctx := context.Background()
	ctx = logging.LoggerWithContext(ctx, log)

	if err := start(ctx); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func start(ctx context.Context) error {
	log := logging.LoggerFromContext(ctx)

	reg, err := registry.NewRegistry()
	if err != nil {
		return err
	}

	log.Info("database connection established")

	bookshopConfig := setBookshopConfig(reg)

	bs, err := bookshop.NewBookshop(ctx, bookshopConfig)
	if err != nil {
		return err
	}

	srv, err := server.NewServer(bs.Port())
	if err != nil {
		return err
	}

	return srv.Start(ctx, bs.Routes())
}

func setBookshopConfig(v *viper.Viper) *bookshop.Config {
	return &bookshop.Config{
		Port: v.GetInt("BOOKSHOP_APP_PORT"),
	}
}
