package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	port              int
	logger            *slog.Logger
	runBeforeShutdown []func()
}

func NewServer(port int, logger *slog.Logger) *Server {
	return &Server{
		port:              port,
		logger:            logger,
		runBeforeShutdown: make([]func(), 0),
	}
}

func (s *Server) Start(ctx context.Context, handler http.Handler) error {
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
		s.logger.Info("server received shutdown signal", slog.String("signal", sig.String()))

		for _, fn := range s.runBeforeShutdown {
			fn()
		}

		ctx, cancel := context.WithTimeout(ctx, time.Second*30)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			shutdownErr <- err
		}

		s.logger.Info("server completing background tasks")

		shutdownErr <- nil
	}()

	s.logger.Info("server starting", slog.Int("port", s.port))

	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownErr
	if err != nil {
		return err
	}

	s.logger.Info("server stopped", slog.Int("port", s.port))

	return nil
}

func (s *Server) RunBeforeShutdown(fn func()) {
	s.runBeforeShutdown = append(s.runBeforeShutdown, fn)
}
