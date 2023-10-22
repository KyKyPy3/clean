package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
	Logger   LoggerConfig
	Redis    RedisConfig
}

type LoggerConfig struct {
	Mode     string
	Level    string
	Encoding string
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

type PostgresConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	DbName       string
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

func NewConfig(path string) (*Config, error) {
	v := viper.New()

	v.SetConfigFile(path)
	v.AddConfigPath(".")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("Config file not found")
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
