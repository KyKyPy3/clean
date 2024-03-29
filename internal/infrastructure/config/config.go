package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Certs    CertsConfig
	Jwt      JwtConfig
	Logger   LoggerConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	Kafka    KafkaConfig
}

type ServerConfig struct {
	Host         string
	Port         string
	SSL          bool
	Version      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	Name         string
}

type CertsConfig struct {
	PrivateKey string
	PublicKey  string
}

type JwtConfig struct {
	AccessTokenMaxAge  time.Duration
	RefreshTokenMaxAge time.Duration
}

type LoggerConfig struct {
	Mode     string
	Level    string
	Encoding string
}

type PostgresConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	DBName       string
	SSLMode      bool
	MaxOpenConn  int
	ConnLifetime time.Duration
	MaxIdleTime  time.Duration
}

type RedisConfig struct {
	Host        string
	Port        string
	Password    string
	DB          int
	MinIdleConn int
	PoolSize    int
	PoolTimeout time.Duration
}

type KafkaConfig struct {
	Brokers []string
	GroupID string
}

func NewConfig(path string) (*Config, error) {
	v := viper.New()

	v.SetConfigFile(path)
	v.AddConfigPath(".")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			return nil, errors.New("config file not found")
		}

		return nil, err
	}

	config := &Config{}
	err := v.Unmarshal(config)
	if err != nil {
		return nil, fmt.Errorf("unable to decode config file: %w", err)
	}

	return config, nil
}
