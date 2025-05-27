package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

// Server представляет основной сервер приложения.
type Server struct {
	server *http.Server
}

// Application интерфейс приложения.
type Application interface { // TODO
}

// NewServer конструктор для основного сервера.
func NewServer(_ Application, cfg config.ServerConfig, router *chi.Mux) *Server {
	server := http.Server{
		Addr:        fmt.Sprintf("%v:%v", cfg.Host, cfg.Port),
		Handler:     router,
		ReadTimeout: 10 * time.Second,
	}

	return &Server{server: &server}
}

// Start запускает основной сервер.
func (s *Server) Start(ctx context.Context) error {
	if err := s.server.ListenAndServe(); err != nil {
		return errors.Wrap(err, "[internalhttp::Start]: server closed")
	}
	<-ctx.Done()
	return nil
}

// Stop останавливает основной сервер с поддержкой graceful shutdown.
func (s *Server) Stop(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "[internalhttp::Stop]: graceful shutdown failed")
	}
	return nil
}
