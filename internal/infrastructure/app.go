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
	"github.com/KyKyPy3/clean/internal/modules/registration/application"
	"github.com/KyKyPy3/clean/internal/modules/registration/domain/event"
	registrationHandlers "github.com/KyKyPy3/clean/internal/modules/registration/infrastructure/controller/http/v1"
	registrationPg "github.com/KyKyPy3/clean/internal/modules/registration/infrastructure/gateway/postgres"
	"github.com/KyKyPy3/clean/internal/modules/registration/infrastructure/gateway/queue"
	userUsecase "github.com/KyKyPy3/clean/internal/modules/user/application/usecase"
	userService "github.com/KyKyPy3/clean/internal/modules/user/domain/service"
	userHandlers "github.com/KyKyPy3/clean/internal/modules/user/infrastructure/controller/http/v1"
	userPg "github.com/KyKyPy3/clean/internal/modules/user/infrastructure/gateway/postgres"
	userRedis "github.com/KyKyPy3/clean/internal/modules/user/infrastructure/gateway/redis"
	kafkaClient "github.com/KyKyPy3/clean/pkg/kafka"
	"github.com/KyKyPy3/clean/pkg/latch"
	"github.com/KyKyPy3/clean/pkg/logger"
	"github.com/KyKyPy3/clean/pkg/mediator"
	"github.com/KyKyPy3/clean/pkg/metric"
	"github.com/KyKyPy3/clean/pkg/outbox"
	postgresClient "github.com/KyKyPy3/clean/pkg/postgres"
	redisClient "github.com/KyKyPy3/clean/pkg/redis"
	"github.com/KyKyPy3/clean/pkg/tracing"
)

const (
	queueTopic = "registration"
)

type App struct {
	cfg         *config.Config
	pgClient    *sqlx.DB
	redisClient *redis.Client
	kafkaClient *kafka.Conn
	web         *Web
	consumer    *Consumer
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
	pgClient, err := postgresClient.New(ctx, postgresClient.Config{
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
	consumer := NewConsumer(cfg, lock, logger)

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

	err = a.consumer.Start(ctx)
	if err != nil {
		return err
	}

	a.connectHandlers(ctx, kafkaProducer)

	return nil
}

func (a *App) connectHandlers(ctx context.Context, producer kafkaClient.Producer) {
	// Health endpoint
	NewHealthHandlers(a.web.mountPoint(), a.logger, a)

	// Init core systems
	pubsub := mediator.New(a.logger)
	trManager := manager.Must(trmsqlx.NewDefaultFactory(a.pgClient))
	queue := queue.NewQueue(producer)
	outboxMngr := outbox.New(a.cfg, a.pgClient, queue, trmsqlx.DefaultCtxGetter, a.logger)
	outboxMngr.Start(ctx, a.lock, outbox.Options{Heartbeat: time.Second * 15})

	// Init user layers
	apiMountPoint := a.web.mountPoint().Group("/api/v1")

	userPgStorage := userPg.NewUserPgStorage(a.pgClient, trmsqlx.DefaultCtxGetter, a.logger)
	userRedisStorage := userRedis.NewUserRedisStorage(a.redisClient, a.logger)
	userSrv := userService.NewUserService(userPgStorage, userRedisStorage, a.logger)
	userUsecase := userUsecase.NewUserUsecase(userSrv, trManager, a.logger)

	userHandlers.NewUserHandlers(apiMountPoint, userUsecase, a.logger)

	// Init registration layout
	regPgStorage := registrationPg.NewRegistrationPgStorage(a.pgClient, trmsqlx.DefaultCtxGetter, a.logger)
	regUniqPolicy := application.NewUniquenessPolicy(userPgStorage, a.logger)
	regApplication := application.NewApplication(regPgStorage, regUniqPolicy, trManager, pubsub, a.logger)
	pubsub.Subscribe(event.RegistrationCreated, func(ctx context.Context, e mediator.Event) error {
		a.logger.Debugf("Receive domain event %v", e)

		err := outboxMngr.Publish(ctx, queueTopic, e)
		if err != nil {
			return err
		}

		return nil
	})

	pubsub.Subscribe(event.RegistrationVerified, func(ctx context.Context, e mediator.Event) error {
		a.logger.Debugf("Receive domain event %v", e)

		return nil
	})

	registrationHandlers.NewRegistrationHandlers(apiMountPoint, regApplication, a.logger)
}
