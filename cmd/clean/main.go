package main

import (
	"context"
	"github.com/KyKyPy3/clean/config"
	"github.com/KyKyPy3/clean/internal/app"
	"github.com/KyKyPy3/clean/pkg/kafka"
	"github.com/KyKyPy3/clean/pkg/logger"
	"github.com/KyKyPy3/clean/pkg/metric"
	"github.com/KyKyPy3/clean/pkg/postgres"
	"github.com/KyKyPy3/clean/pkg/redis"
	"github.com/KyKyPy3/clean/pkg/tracing"
	"github.com/KyKyPy3/clean/pkg/utils"
	"log"
)

const DefaultConfigFile = "config.yml"

func main() {
	ctx := context.Background()
	log.Println("Starting clean microservice")

	// Try to get config file name from environment
	configFile := utils.GetEnvVar("config", DefaultConfigFile)

	cfg, err := config.NewConfig(configFile)
	if err != nil {
		log.Fatalf("Loading config: %v", err)
	}

	// Init logger
	appLogger := logger.NewLogger(&cfg.Logger)
	appLogger.Init()
	appLogger.Infof("Version: %s, LogLevel: %s, SSL: %v", cfg.Server.Version, cfg.Logger.Level, cfg.Server.SSL)

	// Init pg client
	pgClient, err := postgres.New(ctx, &cfg.Postgres)
	if err != nil {
		appLogger.Fatalf("Can't init Postgres database connection: %s", err)
	} else {
		appLogger.Infof("Postgres connected")
	}
	defer func() { _ = pgClient.Close() }()

	// Init redis client
	redisClient, err := redis.New(ctx, &cfg.Redis)
	if err != nil {
		appLogger.Fatalf("Can't init Redis database connection: %s", err)
	} else {
		appLogger.Info("Redis connected")
	}
	defer func() { _ = redisClient.Close() }()

	// Init kafka client
	kafkaClient, err := kafka.New(ctx, &cfg.Kafka)
	if err != nil {
		appLogger.Fatalf("Can't init Kafka connection: %s", err)
	} else {
		brokers, err := kafkaClient.Brokers()
		if err != nil {
			appLogger.Fatalf("Can't get kafka brokers: %s", err)
		}
		appLogger.Info("Kafka connected to brokers %+v", brokers)
	}
	defer func() { _ = kafkaClient.Close() }()

	// Set trace provider
	traceShutdown, err := tracing.New(ctx, cfg.Server.Name)
	if err != nil {
		appLogger.Fatalf("Can't init trace: %s", err)
	} else {
		appLogger.Info("Tracing initialized")
	}
	defer func() {
		if err := traceShutdown(ctx); err != nil {
			appLogger.Fatalf("Can't shutdown trace client: %s", err)
		}
	}()

	// Set metric provider
	metricShutdown, err := metric.New(ctx, cfg.Server.Name)
	if err != nil {
		appLogger.Fatalf("Can't init metrics: %s", err)
	} else {
		appLogger.Info("Metrics initialized")
	}
	defer func() {
		if err := metricShutdown(ctx); err != nil {
			appLogger.Fatalf("Can't shutdown metrics client: %s", err)
		}
	}()

	// Run our service
	srv := app.New(cfg, appLogger, pgClient, redisClient, kafkaClient)
	if err = srv.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
