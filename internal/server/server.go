package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/dlbarduzzi/guestbook/internal/logging"
)

type Server struct {
	port int
	wg   sync.WaitGroup
}

func NewServer(port int) *Server {
	return &Server{
		port: port,
	}
}

func (s *Server) Start(ctx context.Context, handler http.Handler) error {
	logger := logging.LoggerFromContext(ctx)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      handler,
		IdleTimeout:  time.Second * 60,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
	}

	shutdownErr := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		sig := <-quit
		logger.Info("server received shutdown signal", slog.String("signal", sig.String()))

		ctx, cancel := context.WithTimeout(ctx, time.Second*30)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			shutdownErr <- err
		}

		logger.Info("server completing background tasks")

		s.wg.Wait()
		shutdownErr <- nil
	}()

	logger.Info("server starting", slog.Int("port", s.port))

	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownErr
	if err != nil {
		return err
	}

	logger.Info("server stopped", slog.Int("port", s.port))

	return nil
}

func (s *Server) Background(ctx context.Context, fn func()) {
	s.wg.Add(1)
	logger := logging.LoggerFromContext(ctx)

	go func() {
		defer s.wg.Done()

		defer func() {
			if err := recover(); err != nil {
				logger.Error(fmt.Sprintf("%v", err))
			}
		}()
	}()

	fn()
}
