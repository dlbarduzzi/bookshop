package bookshop

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/dlbarduzzi/bookshop/internal/bookshop/model"
	"github.com/dlbarduzzi/bookshop/internal/logging"
)

const version = "unknown"

type Bookshop struct {
	config *Config
	logger *slog.Logger
	models model.Models
}

func NewBookshop(ctx context.Context, cfg *Config, db *sql.DB) (*Bookshop, error) {
	log := logging.LoggerFromContext(ctx)

	cfg, err := cfg.parseConfig()
	if err != nil {
		return nil, err
	}

	return &Bookshop{
		config: cfg,
		logger: log,
		models: model.NewModels(db),
	}, nil
}

func (bs *Bookshop) Port() int {
	return bs.config.Port
}
