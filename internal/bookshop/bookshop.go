package bookshop

import (
	"context"
	"database/sql"
	"log/slog"
	"sync"

	"github.com/dlbarduzzi/bookshop/internal/bookshop/model"
	"github.com/dlbarduzzi/bookshop/internal/logging"
	"github.com/dlbarduzzi/bookshop/internal/mailer"
)

const version = "unknown"

type Bookshop struct {
	config *Config
	logger *slog.Logger
	models model.Models
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

func NewBookshop(ctx context.Context, cfg *Config, db *sql.DB, mailer mailer.Mailer) (*Bookshop, error) {
	log := logging.LoggerFromContext(ctx)

	cfg, err := cfg.parseConfig()
	if err != nil {
		return nil, err
	}

	return &Bookshop{
		config: cfg,
		logger: log,
		models: model.NewModels(db),
		mailer: mailer,
	}, nil
}

func (bs *Bookshop) Port() int {
	return bs.config.Port
}

func (bs *Bookshop) WaitGroup() *sync.WaitGroup {
	return &bs.wg
}
