package main

import (
	"context"
	"os"

	"github.com/spf13/viper"

	"github.com/dlbarduzzi/bookshop/internal/bookshop"
	"github.com/dlbarduzzi/bookshop/internal/logging"
	"github.com/dlbarduzzi/bookshop/internal/registry"
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

	bookshopConfig := setBookshopConfig(reg)

	bs, err := bookshop.NewBookshop(ctx, bookshopConfig)
	if err != nil {
		return err
	}

	log.Info("application running", "port", bs.Port())

	return nil
}

func setBookshopConfig(v *viper.Viper) *bookshop.Config {
	return &bookshop.Config{
		Port: v.GetInt("BOOKSHOP_APP_PORT"),
	}
}
