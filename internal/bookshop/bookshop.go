package bookshop

import (
	"context"
	"log/slog"
	"sync"

	"github.com/dlbarduzzi/bookshop/internal/logging"
)

type Bookshop struct {
	config *Config
	logger *slog.Logger
	wg     sync.WaitGroup
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

func (bs *Bookshop) WaitGroup() *sync.WaitGroup {
	return &bs.wg
}
