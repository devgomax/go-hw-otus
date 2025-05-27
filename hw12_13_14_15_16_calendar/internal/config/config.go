package config

import (
	"github.com/pkg/errors"  //nolint:depguard
	"github.com/spf13/pflag" //nolint:depguard
	"github.com/spf13/viper" //nolint:depguard
)

func init() {
	pflag.String(ConfigPath, "/etc/calendar/config.yaml", "Path to configuration file")
	pflag.Parse()
}

// Key строковый алиас для ключей конфигурации.
type Key = string

const (
	// ConfigPath ключ конфигурации, указывающий путь до конфиг-файла.
	ConfigPath Key = "config"
)

// DBType строковый алиас для поддерживаемых типов БД.
type DBType = string

// Поддерживаемые типы БД.
const (
	DBTypeSQL      DBType = "sql"
	DBTypeInMemory DBType = "in-memory"
)

// DBConfig модель конфига для БД.
type DBConfig struct {
	DBType DBType `mapstructure:"db_type"`
}

// LoggerConfig модель конфига для логгера.
type LoggerConfig struct {
	Level              string `mapstructure:"level"`
	DisableSampling    bool   `mapstructure:"disable_sampling"`
	TimestampFieldName string `mapstructure:"timestamp_field_name"`
	LevelFieldName     string `mapstructure:"level_field_name"`
	MessageFieldName   string `mapstructure:"message_field_name"`
	ErrorFieldName     string `mapstructure:"error_field_name"`
	TimeFieldFormat    string `mapstructure:"time_field_format"`
}

// ServerConfig модель конфига для сервера.
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

// Config модель основного конфига приложения.
type Config struct {
	Logger       LoggerConfig `mapstructure:"logger"`
	Database     DBConfig     `mapstructure:"database"`
	ServerConfig ServerConfig `mapstructure:"server"`
}

// NewConfig конструктор для основного конфига приложения.
func NewConfig() (*Config, error) {
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return nil, errors.Wrap(err, "[main::NewConfig]: failed to bind flag set to config")
	}

	var c Config

	viper.SetConfigFile(viper.GetString(ConfigPath))
	viper.AutomaticEnv()

	err := viper.Unmarshal(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
