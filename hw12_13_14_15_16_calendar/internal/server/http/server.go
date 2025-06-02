package internalhttp

import (
	"context"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// Server представляет http сервер приложения.
type Server struct {
	server *http.Server
}

// NewServer конструктор для http сервера.
func NewServer(addr string, handler http.Handler) *Server {
	server := http.Server{
		Addr:        addr,
		Handler:     handler,
		ReadTimeout: 10 * time.Second,
	}

	return &Server{server: &server}
}

// Start запускает http сервер.
func (s *Server) Start(ctx context.Context) error {
	log.Info().Msgf("calendar HTTP is running on %v", s.server.Addr)

	if err := s.server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return errors.Wrap(err, "[internalhttp::Start]: server closed")
		}
	}
	<-ctx.Done()
	return nil
}

// Stop останавливает http сервер с поддержкой graceful shutdown.
func (s *Server) Stop(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "[internalhttp::Stop]: graceful shutdown failed")
	}
	log.Info().Msg("calendar HTTP gracefully shutdown")
	return nil
}
