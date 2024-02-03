package infrastructure

import (
	"context"
	"time"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"

	"github.com/KyKyPy3/clean/internal/infrastructure/config"
	"github.com/KyKyPy3/clean/internal/infrastructure/queue"
	"github.com/KyKyPy3/clean/internal/modules/registration"
	email_gateway "github.com/KyKyPy3/clean/internal/modules/registration/infrastructure/gateway/email"
	queueGateway "github.com/KyKyPy3/clean/internal/modules/registration/infrastructure/gateway/queue"
	"github.com/KyKyPy3/clean/internal/modules/user"
	"github.com/KyKyPy3/clean/pkg/email"
	kafkaClient "github.com/KyKyPy3/clean/pkg/kafka"
	"github.com/KyKyPy3/clean/pkg/latch"
	"github.com/KyKyPy3/clean/pkg/logger"
	"github.com/KyKyPy3/clean/pkg/mediator"
	"github.com/KyKyPy3/clean/pkg/metric"
	"github.com/KyKyPy3/clean/pkg/outbox"
	"github.com/KyKyPy3/clean/pkg/postgres"
	redisClient "github.com/KyKyPy3/clean/pkg/redis"
	"github.com/KyKyPy3/clean/pkg/tracing"
)

type App struct {
	cfg         *config.Config
	pgClient    *sqlx.DB
	redisClient *redis.Client
	kafkaClient *kafka.Conn
	web         *Web
	consumer    *queue.Consumer
	lock        *latch.CountDownLatch
	logger      logger.Logger
}

func NewApp(
	ctx context.Context,
	cfg *config.Config,
	logger logger.Logger,
	lock *latch.CountDownLatch,
) *App {
	// Init pg client
	pgClient, err := postgres.New(ctx, postgres.Config{
		Host:         cfg.Postgres.Host,
		Port:         cfg.Postgres.Port,
		User:         cfg.Postgres.User,
		Password:     cfg.Postgres.Password,
		DbName:       cfg.Postgres.DbName,
		SSLMode:      cfg.Postgres.SSLMode,
		MaxOpenConn:  cfg.Postgres.MaxOpenConn,
		ConnLifetime: cfg.Postgres.ConnLifetime,
		MaxIdleTime:  cfg.Postgres.MaxIdleTime,
	})
	if err != nil {
		logger.Fatalf("Can't init Postgres database connection: %s", err)
	} else {
		logger.Infof("Postgres connected")
	}

	// Init redis client
	rdClient, err := redisClient.New(ctx, redisClient.Config{
		Host:        cfg.Redis.Host,
		Port:        cfg.Redis.Port,
		Password:    cfg.Redis.Password,
		DB:          cfg.Redis.DB,
		MinIdleConn: cfg.Redis.MinIdleConn,
		PoolSize:    cfg.Redis.PoolSize,
		PoolTimeout: cfg.Redis.PoolTimeout,
	})
	if err != nil {
		logger.Fatalf("Can't init Redis database connection: %s", err)
	} else {
		logger.Info("Redis connected")
	}

	// Init kafka client
	kfClient, err := kafkaClient.New(ctx, kafkaClient.Config{
		Brokers: cfg.Kafka.Brokers,
		GroupID: cfg.Kafka.GroupID,
	})
	if err != nil {
		logger.Fatalf("Can't init Kafka connection: %s", err)
	} else {
		brokers, err := kfClient.Brokers()
		if err != nil {
			logger.Fatalf("Can't get kafka brokers: %s", err)
		}
		logger.Info("Kafka connected to brokers %+v", brokers)
	}

	web := NewWeb(cfg, logger, lock)
	consumer := queue.NewConsumer(cfg, lock, logger)

	return &App{
		cfg:         cfg,
		logger:      logger,
		lock:        lock,
		pgClient:    pgClient,
		kafkaClient: kfClient,
		redisClient: rdClient,
		web:         web,
		consumer:    consumer,
	}
}

func (a *App) Shutdown() {
	_ = a.pgClient.Close()
	_ = a.redisClient.Close()
	_ = a.kafkaClient.Close()
	_ = a.consumer.Stop()
	_ = a.web.Shutdown()
}

func (a *App) Run(ctx context.Context) error {
	// Set trace provider
	traceShutdown, err := tracing.New(ctx, a.cfg.Server.Name)
	if err != nil {
		a.logger.Fatalf("Can't init trace: %s", err)
	} else {
		a.logger.Info("Tracing initialized")
	}
	a.lock.Add(1)
	go func() {
		<-ctx.Done()
		if err := traceShutdown(ctx); err != nil {
			a.logger.Fatalf("Can't shutdown trace client: %s", err)
		}
		a.lock.Done()
	}()

	// Set metric provider
	metricShutdown, err := metric.New(ctx, a.cfg.Server.Name)
	if err != nil {
		a.logger.Fatalf("Can't init metrics: %s", err)
	} else {
		a.logger.Info("Metrics initialized")
	}
	a.lock.Add(1)
	go func() {
		<-ctx.Done()
		if err := metricShutdown(ctx); err != nil {
			a.logger.Fatalf("Can't shutdown metrics client: %s", err)
		}
		a.lock.Done()
	}()

	kafkaProducer := kafkaClient.NewProducer(a.logger, a.cfg.Kafka.Brokers)
	a.lock.Add(1)
	go func() {
		<-ctx.Done()
		kafkaProducer.Close()
		a.lock.Done()
	}()

	// Collect metrics
	//meter := otel.GetMeterProvider().Meter("")
	//testMetric, _ := meter.Int64Counter(
	//	"test_metric",
	//	otelmetric.WithDescription("Number of jobs received by the compute node"),
	//)
	//testMetric.Add(ctx, 1)

	err = a.web.Start()
	if err != nil {
		return err
	}

	a.connectHandlers(ctx, kafkaProducer)

	err = a.consumer.Start(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) connectHandlers(ctx context.Context, producer kafkaClient.Producer) {
	// Health endpoint
	NewHealthHandlers(a.web.mountPoint(), a.logger, a)

	// Init core systems
	pubsub := mediator.New(a.logger)
	trManager := manager.Must(trmsqlx.NewDefaultFactory(a.pgClient))
	queue := queueGateway.NewQueue(producer)
	outboxMngr := outbox.New(a.cfg, a.pgClient, queue, trmsqlx.DefaultCtxGetter, a.logger)
	outboxMngr.Start(ctx, a.lock, outbox.Options{Heartbeat: time.Second * 15})
	apiMountPoint := a.web.mountPoint().Group("/api/v1")
	emailClient := email.New(a.logger)
	emailGateway := email_gateway.New(emailClient, a.logger)

	////////////////////////////////
	// Init user layout
	////////////////////////////////
	user.InitUserHandlers(a.pgClient, apiMountPoint, pubsub, trManager, a.logger)

	////////////////////////////////
	// Init registration layout
	////////////////////////////////
	registration.InitUserHandlers(a.pgClient, apiMountPoint, pubsub, a.consumer, trManager, emailGateway, outboxMngr, a.logger)
}
