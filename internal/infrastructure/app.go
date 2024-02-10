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
	reg_postgres "github.com/KyKyPy3/clean/internal/modules/registration/infrastructure/gateway/postgres"
	queueGateway "github.com/KyKyPy3/clean/internal/modules/registration/infrastructure/gateway/queue"
	"github.com/KyKyPy3/clean/internal/modules/session"
	"github.com/KyKyPy3/clean/internal/modules/session/infrastructure/controller/middleware"
	session_redis "github.com/KyKyPy3/clean/internal/modules/session/infrastructure/gateway/redis"
	"github.com/KyKyPy3/clean/internal/modules/user"
	user_postgres "github.com/KyKyPy3/clean/internal/modules/user/infrastructure/gateway/postgres"
	"github.com/KyKyPy3/clean/pkg/email"
	"github.com/KyKyPy3/clean/pkg/jwt"
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

const heartbeatInterval = time.Second * 15

type App struct {
	cfg         *config.Config
	pgClient    *sqlx.DB
	redisClient *redis.Client
	kafkaClient *kafka.Conn
	web         *Web
	jwt         *jwt.JWT
	consumer    *queue.Consumer
	producer    kafkaClient.Producer
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
		DBName:       cfg.Postgres.DBName,
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
		var brokers []kafka.Broker
		brokers, err = kfClient.Brokers()
		if err != nil {
			logger.Fatalf("Can't get kafka brokers: %s", err)
		}
		logger.Info("Kafka connected to brokers %+v", brokers)
	}

	// Init JWT manager
	jwtManager, err := jwt.NewJWT(cfg.Certs.PrivateKey, cfg.Certs.PublicKey)
	if err != nil {
		logger.Fatalf("Can't parse certs: %s", err)
	}

	web := NewWeb(cfg, logger, lock)
	kafkaProducer := kafkaClient.NewProducer(logger, cfg.Kafka.Brokers)
	consumer := queue.NewConsumer(cfg, lock, logger)

	return &App{
		cfg:         cfg,
		logger:      logger,
		lock:        lock,
		pgClient:    pgClient,
		kafkaClient: kfClient,
		jwt:         jwtManager,
		redisClient: rdClient,
		producer:    kafkaProducer,
		web:         web,
		consumer:    consumer,
	}
}

func (a *App) Shutdown() {
	_ = a.producer.Close()
	_ = a.web.Shutdown()
	_ = a.pgClient.Close()
	_ = a.redisClient.Close()
	_ = a.kafkaClient.Close()
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
		if err = traceShutdown(ctx); err != nil {
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
		if err = metricShutdown(ctx); err != nil {
			a.logger.Fatalf("Can't shutdown metrics client: %s", err)
		}
		a.lock.Done()
	}()

	// Collect metrics
	// meter := otel.GetMeterProvider().Meter("")
	// testMetric, _ := meter.Int64Counter(
	//	"test_metric",
	//	otelmetric.WithDescription("Number of jobs received by the compute node"),
	// )
	// testMetric.Add(ctx, 1)

	err = a.web.Start()
	if err != nil {
		return err
	}

	a.connectHandlers(ctx)

	err = a.consumer.Start(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) connectHandlers(ctx context.Context) {
	mountPoint := a.web.mountPoint()

	// Health endpoint
	NewHealthHandlers(mountPoint, a.logger, a.pgClient, a.redisClient, a.kafkaClient)

	// Init core systems
	pubsub := mediator.New(a.logger)
	trManager := manager.Must(trmsqlx.NewDefaultFactory(a.pgClient))
	queue := queueGateway.NewQueue(a.producer)
	outboxMngr := outbox.New(a.cfg, a.pgClient, queue, trmsqlx.DefaultCtxGetter, a.logger)
	outboxMngr.Start(ctx, a.lock, outbox.Options{Heartbeat: heartbeatInterval})
	emailClient := email.New(a.logger)
	emailGateway := email_gateway.New(emailClient, a.logger)

	userPgStorage := user_postgres.NewUserPgStorage(a.pgClient, trmsqlx.DefaultCtxGetter, a.logger)
	regPgStorage := reg_postgres.NewRegistrationPgStorage(a.pgClient, trmsqlx.DefaultCtxGetter, a.logger)
	sessionStorage := session_redis.NewSessionRedisStorage(a.redisClient, a.logger)

	authMiddleware := middleware.NewAuthMiddleware(a.jwt, sessionStorage, a.logger)
	publicMountPoint := mountPoint.Group("/api/v1")
	privateMountPoint := mountPoint.Group("/api/v1", authMiddleware.Process)

	////////////////////////////////
	// Init user layout
	////////////////////////////////
	user.InitHandlers(
		ctx,
		userPgStorage,
		privateMountPoint,
		pubsub,
		trManager,
		a.logger,
	)

	////////////////////////////////
	// Init session layout
	////////////////////////////////
	session.InitHandlers(
		userPgStorage,
		sessionStorage,
		publicMountPoint,
		privateMountPoint,
		a.cfg,
		a.jwt,
		a.logger,
	)

	////////////////////////////////
	// Init registration layout
	////////////////////////////////
	registration.InitHandlers(
		ctx,
		userPgStorage,
		regPgStorage,
		publicMountPoint,
		pubsub,
		a.consumer,
		trManager,
		emailGateway,
		outboxMngr,
		a.logger,
	)
}
