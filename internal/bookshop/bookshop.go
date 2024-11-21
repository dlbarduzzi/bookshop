package bookshop

import (
	"fmt"
	"log/slog"
	"sync"
)

type Bookshop struct {
	logger *slog.Logger
	wg     *sync.WaitGroup
}

func NewBookshop(logger *slog.Logger) *Bookshop {
	return &Bookshop{
		wg:     &sync.WaitGroup{},
		logger: logger,
	}
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
