package guestbook

import (
	"database/sql"
	"fmt"
	"log/slog"
	"sync"

	"github.com/dlbarduzzi/guestbook/internal/guestbook/model"
)

type Guestbook struct {
	config *Config
	logger *slog.Logger
	models model.Models
	wg     *sync.WaitGroup
}

func NewGuestbook(db *sql.DB, logger *slog.Logger, config *Config) (*Guestbook, error) {
	cfg, err := config.parse()
	if err != nil {
		return nil, err
	}

	return &Guestbook{
		config: cfg,
		logger: logger,
		models: model.NewModels(db),
		wg:     &sync.WaitGroup{},
	}, nil
}

func (g *Guestbook) Port() int {
	return g.config.Port
}

func (g *Guestbook) Background(fn func()) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()

		defer func() {
			if err := recover(); err != nil {
				g.logger.Error(fmt.Sprintf("%v", err))
			}
		}()

		fn()
	}()
}

func (g *Guestbook) Shutdown() {
	g.wg.Wait()
}
