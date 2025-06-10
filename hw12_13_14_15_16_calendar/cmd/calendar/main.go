package main

import (
	"context"
	syslog "log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/app/calendar"
	"github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/config"
	"github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/logger"
	eventspb "github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/pb/events"
	internalgrpc "github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/server/http"
	"github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/server/http/middleware"
	"github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	if pflag.Arg(0) == "version" {
		printVersion()
		return
	}

	cfg, err := config.NewConfig()
	if err != nil {
		cancel()
		log.Fatal().Err(err).Msg("failed to get service config") //nolint:gocritic
	}

	if err = logger.ConfigureLogging(cfg.Logger); err != nil {
		cancel()
		log.Fatal().Err(err).Msg("failed to configure logging")
	}

	var repo storage.IRepository

	switch cfg.Database.DBType {
	case config.DBTypeSQL:
		repo = sqlstorage.New()
	case config.DBTypeInMemory:
		repo = memorystorage.New()
	default:
		cancel()
		log.Fatal().Msg("invalid config value for db_type")
	}

	if err = repo.Connect(ctx, ""); err != nil { // Empty DNS for PG to process environment variables
		cancel()
		log.Fatal().Err(err).Msg("failed to connect to DB")
	}
	defer repo.Close()

	calendarApp := calendar.New(repo)

	file, err := os.OpenFile("server.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o664)
	if err != nil {
		cancel()
		repo.Close()
		log.Fatal().Err(err).Msg("Не удалось открыть файл для логов")
	}
	defer file.Close()

	servLogger := syslog.New(file, "", 0)

	serverGRPC := internalgrpc.NewServer(servLogger, calendarApp)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err = serverGRPC.Start(cfg.GRPCConfig.GetAddr()); err != nil {
			log.Fatal().Err(err).Msg("failed to run GRPC server")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		<-ctx.Done()

		serverGRPC.Stop()
	}()

	handler := chi.NewRouter()
	handler.Use(middleware.NewLoggingMiddleware(servLogger))
	handler.Use(chimiddleware.Recoverer)

	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err = eventspb.RegisterEventsHandlerFromEndpoint(ctx, gwmux, cfg.GRPCConfig.GetAddr(), opts); err != nil {
		cancel()
		repo.Close()
		file.Close()
		log.Fatal().Err(err).Msgf("failed to dial to %q", cfg.GRPCConfig.GetAddr())
	}

	handler.Handle("/*", gwmux)

	server := internalhttp.NewServer(cfg.HTTPConfig.GetAddr(), handler)

	wg.Add(1)
	go func() {
		defer wg.Done()

		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			log.Error().Err(err).Msg("failed to stop HTTP server")
		}
	}()

	if err = server.Start(ctx); err != nil {
		cancel()
		repo.Close()
		file.Close()
		log.Fatal().Err(err).Msg("failed to run HTTP server")
	}

	wg.Wait()
}
