package bookshop

import (
	"database/sql"
	"fmt"
	"log/slog"
	"sync"

	"github.com/dlbarduzzi/bookshop/internal/bookshop/model"
)

type Bookshop struct {
	config *Config
	logger *slog.Logger
	models model.Models
	wg     *sync.WaitGroup
}

func NewBookshop(db *sql.DB, logger *slog.Logger, config *Config) (*Bookshop, error) {
	cfg, err := config.parse()
	if err != nil {
		return nil, err
	}

	return &Bookshop{
		config: cfg,
		logger: logger,
		models: model.NewModels(db),
		wg:     &sync.WaitGroup{},
	}, nil
}

func (b *Bookshop) Port() int {
	return b.config.Port
}

func (b *Bookshop) Background(fn func()) {
	b.wg.Add(1)

	go func() {
		defer b.wg.Done()

		defer func() {
			if err := recover(); err != nil {
				b.logger.Error(fmt.Sprintf("%v", err))
			}
		}()

		fn()
	}()
}

func (b *Bookshop) Shutdown() {
	b.wg.Wait()
}
