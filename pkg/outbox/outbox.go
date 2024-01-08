package outbox

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/KyKyPy3/clean/internal/infrastructure/config"
	kafkaClient "github.com/KyKyPy3/clean/pkg/kafka"
	"github.com/segmentio/kafka-go"
	"time"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/jmoiron/sqlx"

	"github.com/KyKyPy3/clean/pkg/latch"
	"github.com/KyKyPy3/clean/pkg/logger"
)

type Event interface {
}

type Options struct {
	Heartbeat time.Duration
}

type outbox struct {
	Id       int64
	Topic    string
	Payload  []byte
	Consumed bool
}

type Manager struct {
	cfg    *config.Config
	db     *sqlx.DB
	getter *trmsqlx.CtxGetter
	logger logger.Logger
}

func New(cfg *config.Config, db *sqlx.DB, getter *trmsqlx.CtxGetter, logger logger.Logger) Manager {
	return Manager{
		db:     db,
		getter: getter,
		logger: logger,
		cfg:    cfg,
	}
}

func (m Manager) Publish(ctx context.Context, event Event) error {
	e, err := json.Marshal(event)
	if err != nil {
		return err
	}

	stmt, err := m.getter.DefaultTrOrDB(ctx, m.db).PreparexContext(ctx, "INSERT INTO outbox (topic, payload) VALUES ($1, $2)")
	if err != nil {
		return err
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
			m.logger.Errorf("can't close create statement, err: %w", err)
		}
	}()

	_, err = stmt.ExecContext(ctx, "test_topic", e)
	if err != nil {
		return err
	}

	return nil
}

func (m Manager) Start(ctx context.Context, lock *latch.CountDownLatch, options Options) {
	lock.Add(1)
	kafkaProducer := kafkaClient.NewProducer(m.logger, m.cfg.Kafka.Brokers)

	go func() {
		defer lock.Done()
		ticker := time.NewTicker(options.Heartbeat)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				m.logger.Debugf("Consume outbox events")
				err := m.Consume(ctx, kafkaProducer)
				if err != nil {
					m.logger.Errorf("failed to send outbox messages to broker, err: %w", err)
				}
				// Consume events and publish to kafka
			case <-ctx.Done():
				kafkaProducer.Close()
				return
			}
		}
	}()
}

func (m Manager) Consume(ctx context.Context, producer kafkaClient.Producer) error {
	outboxes := []outbox{}

	tx, err := m.db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			rb := tx.Rollback()
			if rb != nil {
				err = rb
			}
		}
	}()

	locked, err := m.getLock(ctx, tx)
	if err != nil {
		return fmt.Errorf("failed get lock: %w", err)
	}
	if !locked {
		return nil
	}

	// TODO: Now each tim ewe process only 50 messages
	err = tx.SelectContext(
		ctx,
		&outboxes,
		`SELECT * FROM outbox WHERE consumed=FALSE ORDER BY id ASC LIMIT $1`,
		50,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		} else {
			return err
		}
	}

	if len(outboxes) == 0 {
		return nil
	}

	m.logger.Debugf("Events: %v", outboxes)

	ids := make([]int64, len(outboxes))
	for k, message := range outboxes {
		err := producer.PublishMessage(ctx, kafka.Message{
			Topic: message.Topic,
			Value: message.Payload,
			Time:  time.Now().UTC(),
		})
		if err != nil {
			continue
		}

		ids[k] = message.Id
	}
	query, args, err := sqlx.In("UPDATE outbox SET consumed=TRUE WHERE id IN(?)", ids)
	if err != nil {
		return fmt.Errorf("expanding ids to consume: %w", err)
	}
	query = tx.Rebind(query)
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("fail update consumed events state, err: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (m Manager) getLock(ctx context.Context, tx *sqlx.Tx) (bool, error) {
	ok := false
	err := tx.GetContext(ctx, &ok, "SELECT pg_try_advisory_xact_lock(123)")
	if err != nil {
		return false, err
	}

	return ok, nil
}
