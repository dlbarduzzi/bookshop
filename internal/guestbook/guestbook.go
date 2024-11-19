package guestbook

import (
	"fmt"
	"log/slog"
	"sync"
)

type Guestbook struct {
	logger *slog.Logger
	wg     *sync.WaitGroup
}

func NewGuesbook(logger *slog.Logger) *Guestbook {
	return &Guestbook{
		wg:     &sync.WaitGroup{},
		logger: logger,
	}
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
