package main

import (
	"context"
	"fmt"
	"os"

	"github.com/dlbarduzzi/bookshop/internal/bookshop"
	"github.com/dlbarduzzi/bookshop/internal/logging"
	"github.com/dlbarduzzi/bookshop/internal/registry"
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
	log.Info(fmt.Sprintf("Bookshop app running on port %d", bookshopConfig.Port))

	return nil
}

func getBookshopConfig(v *viper.Viper) *bookshop.Config {
	return &bookshop.Config{
		Port: v.GetInt("BS_APP_PORT"),
	}
}
