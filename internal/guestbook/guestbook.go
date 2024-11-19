package guestbook

import (
	"fmt"
	"log/slog"
	"sync"
	"time"
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

func (g *Guestbook) FakeHandler() {
	g.logger.Info("This is the start of FakeHandler...")
	g.Background(func() {
		fmt.Println("START FakeHandler background...")
		time.Sleep(time.Second * 3)
		fmt.Println("END FakeHandler background...")
	})
	g.logger.Info("This is the end of FakeHandler. Goodbye!")
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
