package bookshop

import (
	"context"
	"log/slog"

	"github.com/dlbarduzzi/bookshop/internal/logging"
)

const version = "unknown"

type Bookshop struct {
	config *Config
	logger *slog.Logger
}

func NewBookshop(ctx context.Context, cfg *Config) (*Bookshop, error) {
	log := logging.LoggerFromContext(ctx)

	cfg, err := cfg.parseConfig()
	if err != nil {
		return nil, err
	}

	return &Bookshop{
		config: cfg,
		logger: log,
	}, nil
}

func (bs *Bookshop) Port() int {
	return bs.config.Port
}
