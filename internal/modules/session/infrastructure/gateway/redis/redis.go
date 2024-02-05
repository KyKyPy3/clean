package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/modules/session/application/ports"
	"github.com/KyKyPy3/clean/internal/modules/session/domain/entity"
	"github.com/KyKyPy3/clean/pkg/logger"
)

const (
	basePrefix = "sessions:"
)

type sessionRedisStorage struct {
	basePrefix string
	db         *redis.Client
	logger     logger.Logger
	tracer     trace.Tracer
}

func NewSessionRedisStorage(db *redis.Client, logger logger.Logger) ports.SessionRedisStorage {
	return &sessionRedisStorage{basePrefix: basePrefix, db: db, logger: logger, tracer: otel.Tracer("")}
}

func (s *sessionRedisStorage) Get(ctx context.Context, tokenID common.UID) (entity.Token, error) {
	_, span := s.tracer.Start(ctx, "sessionRedisStorage.Get")
	defer span.End()

	tokenBytes, err := s.db.Get(ctx, s.createKey(tokenID.String())).Bytes()
	if err != nil {
		return entity.Token{}, err
	}

	token := DBToken{}
	if err = json.Unmarshal(tokenBytes, &token); err != nil {
		return entity.Token{}, err
	}

	return TokenFromDB(token)
}

func (s *sessionRedisStorage) Set(ctx context.Context, tokenID common.UID, token entity.Token) error {
	_, span := s.tracer.Start(ctx, "sessionRedisStorage.Set")
	defer span.End()

	t := TokenToDB(token)
	tokenBytes, err := json.Marshal(&t)
	if err != nil {
		return err
	}

	now := time.Now()
	if err := s.db.Set(ctx, s.createKey(tokenID.String()), tokenBytes, time.Unix(token.ExpiresIn(), 0).Sub(now)).Err(); err != nil {
		return err
	}

	return nil
}

func (s *sessionRedisStorage) Delete(ctx context.Context, tokenID common.UID) error {
	_, span := s.tracer.Start(ctx, "sessionRedisStorage.Delete")
	defer span.End()

	err := s.db.Del(ctx, s.createKey(tokenID.String())).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *sessionRedisStorage) createKey(tokenID string) string {
	return fmt.Sprintf("%s: %s", s.basePrefix, tokenID)
}
