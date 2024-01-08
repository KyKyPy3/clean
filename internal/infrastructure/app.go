package infrastructure

import (
	"context"
	v1 "github.com/KyKyPy3/clean/internal/modules/registration/infrastructure/controller/http/v1"
	"github.com/KyKyPy3/clean/pkg/outbox"
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"time"

	"github.com/KyKyPy3/clean/internal/infrastructure/config"
	regUsecase "github.com/KyKyPy3/clean/internal/modules/registration/application/usecase"
	regService "github.com/KyKyPy3/clean/internal/modules/registration/domain/service"
	registrationPg "github.com/KyKyPy3/clean/internal/modules/registration/infrastructure/gateway/postgres"
	userUsecase "github.com/KyKyPy3/clean/internal/modules/user/application/usecase"
	userService "github.com/KyKyPy3/clean/internal/modules/user/domain/service"
	userHandlers "github.com/KyKyPy3/clean/internal/modules/user/infrastructure/controller/http/v1"
	userPg "github.com/KyKyPy3/clean/internal/modules/user/infrastructure/gateway/postgres"
	userRedis "github.com/KyKyPy3/clean/internal/modules/user/infrastructure/gateway/redis"
	kafkaClient "github.com/KyKyPy3/clean/pkg/kafka"
	"github.com/KyKyPy3/clean/pkg/latch"
	"github.com/KyKyPy3/clean/pkg/logger"
	"github.com/KyKyPy3/clean/pkg/metric"
	postgresClient "github.com/KyKyPy3/clean/pkg/postgres"
	redisClient "github.com/KyKyPy3/clean/pkg/redis"
	"github.com/KyKyPy3/clean/pkg/tracing"
)

type App struct {
	cfg         *config.Config
	pgClient    *sqlx.DB
	redisClient *redis.Client
	kafkaClient *kafka.Conn
	web         *Web
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
	pgClient, err := postgresClient.New(ctx, &cfg.Postgres)
	if err != nil {
		logger.Fatalf("Can't init Postgres database connection: %s", err)
	} else {
		logger.Infof("Postgres connected")
	}

	// Init redis client
	rdClient, err := redisClient.New(ctx, &cfg.Redis)
	if err != nil {
		logger.Fatalf("Can't init Redis database connection: %s", err)
	} else {
		logger.Info("Redis connected")
	}

	// Init kafka client
	kfClient, err := kafkaClient.New(ctx, &cfg.Kafka)
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

	return &App{
		cfg:         cfg,
		logger:      logger,
		lock:        lock,
		pgClient:    pgClient,
		kafkaClient: kfClient,
		redisClient: rdClient,
		web:         web,
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

	a.connectHandlers(ctx)

	return nil
}

func (a *App) connectHandlers(ctx context.Context) {
	// Health endpoint
	NewHealthHandlers(a.web.mountPoint(), a.logger, a)

	// Init user layers
	trManager := manager.Must(trmsqlx.NewDefaultFactory(a.pgClient))
	apiMountPoint := a.web.mountPoint().Group("/api/v1")
	outboxMngr := outbox.New(a.cfg, a.pgClient, trmsqlx.DefaultCtxGetter, a.logger)
	outboxMngr.Start(ctx, a.lock, outbox.Options{Heartbeat: time.Second * 5})

	userPgStorage := userPg.NewUserPgStorage(a.pgClient, trmsqlx.DefaultCtxGetter, a.logger)
	userRedisStorage := userRedis.NewUserRedisStorage(a.redisClient, a.logger)
	userSrv := userService.NewUserService(userPgStorage, userRedisStorage, a.logger)
	userUsecase := userUsecase.NewUserUsecase(userSrv, trManager, a.logger)

	userHandlers.NewUserHandlers(apiMountPoint, userUsecase, a.logger)

	// Init registration layout
	regPgStorage := registrationPg.NewRegistrationPgStorage(a.pgClient, trmsqlx.DefaultCtxGetter, a.logger)
	regUniqPolicy := regService.NewUniquenessPolicy(regPgStorage, a.logger)
	regSrv := regService.NewRegistrationService(regPgStorage, regUniqPolicy, a.logger)
	registrationUsecase := regUsecase.NewRegistrationUsecase(regSrv, trManager, outboxMngr, a.logger)

	v1.NewRegistrationHandlers(apiMountPoint, registrationUsecase, a.logger)
}
