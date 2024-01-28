package outbox

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/jmoiron/sqlx"

	"github.com/KyKyPy3/clean/internal/infrastructure/config"
	"github.com/KyKyPy3/clean/pkg/latch"
	"github.com/KyKyPy3/clean/pkg/logger"
)

type Publisher interface {
	Publish(ctx context.Context, event Message) error
}

type Message struct {
	Kind    string
	Topic   string
	Payload []byte
}

type Event interface {
	Kind() string
}

type Options struct {
	Heartbeat time.Duration
}

type outbox struct {
	Id       int64
	Topic    string
	Kind     string
	Payload  []byte
	Consumed bool
}

type Manager struct {
	cfg       *config.Config
	db        *sqlx.DB
	getter    *trmsqlx.CtxGetter
	logger    logger.Logger
	publisher Publisher
}

func New(cfg *config.Config, db *sqlx.DB, publisher Publisher, getter *trmsqlx.CtxGetter, logger logger.Logger) Manager {
	return Manager{
		db:        db,
		getter:    getter,
		logger:    logger,
		cfg:       cfg,
		publisher: publisher,
	}
}

func (m Manager) Publish(ctx context.Context, topic string, event Event) error {
	e, err := json.Marshal(event)
	if err != nil {
		return err
	}

	stmt, err := m.getter.DefaultTrOrDB(ctx, m.db).PreparexContext(ctx, "INSERT INTO outbox (topic, kind, payload) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
			m.logger.Errorf("can't close create statement, err: %w", err)
		}
	}()

	_, err = stmt.ExecContext(ctx, topic, event.Kind(), e)
	if err != nil {
		return err
	}

	return nil
}

func (m Manager) Start(ctx context.Context, lock *latch.CountDownLatch, options Options) {
	lock.Add(1)

	go func() {
		defer lock.Done()
		ticker := time.NewTicker(options.Heartbeat)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				m.logger.Debugf("Consume outbox events")
				err := m.Consume(ctx)
				if err != nil {
					m.logger.Errorf("failed to send outbox messages to broker, err: %w", err)
				}
				// Consume events and publish to queue
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (m Manager) Consume(ctx context.Context) error {
	var outboxes []outbox

	tx, err := m.db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if panicErr := recover(); panicErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				err = fmt.Errorf("panic Error + rollback Error: %w", panicErr.(error))
			} else {
				err = fmt.Errorf("panic Error: %w", panicErr.(error))
			}
		}

		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				err = fmt.Errorf("rollback Error: %w", rollbackErr)
			}
		} else {
			err = tx.Commit()
		}
	}()

	locked, err := m.getLock(ctx, tx)
	if err != nil {
		return fmt.Errorf("failed get lock: %w", err)
	}
	if !locked {
		m.logger.Warnf("can't aquire lock, err: %v", err)
		return nil
	}

	// TODO: Now each time we process only 50 messages
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
		err := m.publisher.Publish(ctx, Message{
			Topic:   message.Topic,
			Kind:    message.Kind,
			Payload: message.Payload,
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
