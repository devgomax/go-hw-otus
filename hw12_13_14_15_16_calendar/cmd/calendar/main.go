package main

import (
	"context"
	syslog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/app"
	"github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/config"
	"github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/server/http"
	"github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
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

	var repo storage.Repository

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

	calendar := app.New(repo)

	file, err := os.OpenFile("server.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o664)
	if err != nil {
		cancel()
		repo.Close()
		log.Fatal().Err(err).Msg("Не удалось открыть файл для логов")
	}
	defer file.Close()

	servLogger := syslog.New(file, "", 0)

	router := chi.NewRouter()
	router.Use(internalhttp.NewLoggingMiddleware(servLogger))
	router.Use(middleware.Recoverer)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if _, inErr := w.Write([]byte("hello-world")); inErr != nil {
			log.Error().Err(inErr).Msg("failed to write response")
		}
	})

	server := internalhttp.NewServer(calendar, cfg.ServerConfig, router)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			log.Error().Err(err).Msg("failed to stop http server")
		}
	}()

	log.Info().Msg("calendar is running...")

	if err = server.Start(ctx); err != nil {
		cancel()
		repo.Close()
		file.Close()
		log.Fatal().Err(err).Msg("failed to start http server")
	}
}
