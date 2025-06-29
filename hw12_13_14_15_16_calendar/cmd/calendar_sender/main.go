package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/app/sender"
	"github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/app/sender/adapters"
	"github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/config"
	"github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/logger"
	memorystorage "github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	cfg, err := config.NewConfig()
	if err != nil {
		cancel()
		log.Fatal().Err(err).Msg("failed to get service config") //nolint:gocritic
	}

	if err = logger.ConfigureLogging(cfg.Logger); err != nil {
		cancel()
		log.Fatal().Err(err).Msg("failed to configure logging")
	}

	rmq, err := adapters.NewRabbitMQClient(cfg.MessageQueueConfig.URL, cfg.MessageQueueConfig.Queue)
	if err != nil {
		cancel()
		log.Fatal().Err(err).Msg("failed to create amqp client")
	}

	var repo sender.IRepository

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

	app := sender.NewApp(repo, rmq)

	go func() {
		<-ctx.Done()

		if err = app.Stop(); err != nil {
			log.Fatal().Err(err).Msg("failed to close amqp connection")
		}
	}()

	if err = app.Run(ctx); err != nil {
		cancel()
		log.Fatal().Err(err).Msg("failed to run sender")
	}
}
