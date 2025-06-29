package config

import (
	"net"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

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

// GetAddr возвращает строку вида "host:port".
func (sc *ServerConfig) GetAddr() string {
	return net.JoinHostPort(sc.Host, sc.Port)
}

// MessageQueueConfig модель конфига для очереди сообщений.
type MessageQueueConfig struct {
	URL   string `mapstructure:"url"`
	Queue string `mapstructure:"queue"`
}

// SchedulerConfig модель конфига для сервиса-планировщика.
type SchedulerConfig struct {
	DBReadInterval time.Duration `mapstructure:"db_read_interval"`
}

// Config модель основного конфига приложения.
type Config struct {
	Logger             LoggerConfig       `mapstructure:"logger"`
	Database           DBConfig           `mapstructure:"database"`
	GRPCConfig         ServerConfig       `mapstructure:"grpc"`
	HTTPConfig         ServerConfig       `mapstructure:"http"`
	MessageQueueConfig MessageQueueConfig `mapstructure:"amqp"`
	SchedulerConfig    SchedulerConfig    `mapstructure:"scheduler"`
}

// NewConfig конструктор для основного конфига приложения.
func NewConfig() (*Config, error) {
	configPath := pflag.String(ConfigPath, "/etc/calendar/config.toml", "Path to configuration file")
	pflag.Parse()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return nil, errors.Wrap(err, "[main::NewConfig]: failed to bind flag set to config")
	}

	var c Config

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()
	viper.SetConfigFile(*configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "[main::NewConfig]: failed to discover and read config file")
	}

	err := viper.Unmarshal(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
