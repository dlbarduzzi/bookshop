package main

import (
	"context"
	"os"

	"github.com/spf13/viper"

	"github.com/dlbarduzzi/bookshop/internal/bookshop"
	"github.com/dlbarduzzi/bookshop/internal/database"
	"github.com/dlbarduzzi/bookshop/internal/logging"
	"github.com/dlbarduzzi/bookshop/internal/mailer"
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

	databaseConfig := setDatabaseConfig(reg)
	bookshopConfig := setBookshopConfig(reg)

	db, err := database.NewDatabase(databaseConfig)
	if err != nil {
		return err
	}

	defer db.Close()
	log.Info("database connection established")

	mailerConfig := setMailerConfig(reg)
	mailer := mailer.NewMailer(&mailerConfig)

	bs, err := bookshop.NewBookshop(ctx, bookshopConfig, db, mailer)
	if err != nil {
		return err
	}

	srv, err := server.NewServer(bs.Port())
	if err != nil {
		return err
	}

	srv.WaitGroup = bs.WaitGroup()
	return srv.Start(ctx, bs.Routes())
}

func setBookshopConfig(v *viper.Viper) *bookshop.Config {
	return &bookshop.Config{
		Port: v.GetInt("BOOKSHOP_APP_PORT"),
	}
}

func setDatabaseConfig(v *viper.Viper) *database.Config {
	return &database.Config{
		ConnectionURL:   v.GetString("DB_CONNECTION_URL"),
		MaxOpenConns:    v.GetInt("DB_MAX_OPEN_CONNS"),
		MaxIdleConns:    v.GetInt("DB_MAX_IDLE_CONNS"),
		ConnMaxIdleTime: v.GetDuration("DB_CONN_MAX_IDLE_TIME"),
	}
}

func setMailerConfig(v *viper.Viper) mailer.Config {
	return mailer.Config{
		Host:     v.GetString("MAILER_HOST"),
		Port:     v.GetInt("MAILER_PORT"),
		Username: v.GetString("MAILER_USERNAME"),
		Password: v.GetString("MAILER_PASSWORD"),
		Sender:   v.GetString("MAILER_SENDER"),
	}
}
