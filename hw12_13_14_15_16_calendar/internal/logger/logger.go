package logger

import (
	"cmp"
	"os"
	"time"

	"github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// ConfigureLogging настраивает основной логгер приложения.
func ConfigureLogging(cfg config.LoggerConfig) error {
	lvl, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		return errors.Wrap(err, "[logger::ConfigureLogging]")
	}

	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	log.Logger = log.Output(output)

	zerolog.SetGlobalLevel(lvl)
	zerolog.DisableSampling(cfg.DisableSampling)
	zerolog.TimestampFieldName = cmp.Or(cfg.TimestampFieldName, zerolog.TimestampFieldName)
	zerolog.LevelFieldName = cmp.Or(cfg.LevelFieldName, zerolog.LevelFieldName)
	zerolog.MessageFieldName = cmp.Or(cfg.MessageFieldName, zerolog.MessageFieldName)
	zerolog.ErrorFieldName = cmp.Or(cfg.ErrorFieldName, zerolog.ErrorFieldName)
	zerolog.TimeFieldFormat = cmp.Or(cfg.TimeFieldFormat, zerolog.TimeFieldFormat)

	return nil
}
