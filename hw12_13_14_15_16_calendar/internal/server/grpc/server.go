package internalgrpc

import (
	syslog "log"
	"net"

	"github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/app"
	eventspb "github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/pb/events"
	"github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/server/grpc/interceptors"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

// Server представляет grpc сервер приложения.
type Server struct {
	server *grpc.Server
}

// NewServer конструктор для grpc сервера.
func NewServer(logger *syslog.Logger, app app.IApp) *Server {
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(interceptors.NewUnaryServerLoggingInterceptor(logger)),
	)
	impl := NewEventsServer(app)
	eventspb.RegisterEventsServer(server, impl)

	return &Server{server: server}
}

// Start запускает grpc сервер.
func (s *Server) Start(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.Wrapf(err, "[grpc::Start]: can't get listener for %q", addr)
	}

	log.Info().Msgf("calendar GRPC is running on %v", addr)

	if err = s.server.Serve(listener); err != nil {
		return errors.Wrap(err, "[grpc::Start]")
	}

	return nil
}

// Stop останавливает grpc сервер с поддержкой graceful shutdown.
func (s *Server) Stop() {
	s.server.GracefulStop()
	log.Info().Msg("calendar GRPC gracefully shutdown")
}
