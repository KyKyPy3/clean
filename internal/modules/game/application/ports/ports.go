package ports

import (
	"context"

	"github.com/KyKyPy3/clean/internal/modules/game/domain/entity"
	"github.com/KyKyPy3/clean/pkg/mediator"
)

type Mediator interface {
	Publish(ctx context.Context, events ...mediator.Event) error
}

type TrManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) (err error)
}

type GamePgStorage interface {
	Fetch(ctx context.Context, limit, offset int64) ([]entity.Game, error)
	Create(ctx context.Context, registration entity.Game) error
}
