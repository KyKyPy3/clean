package postgres

import (
	"context"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/KyKyPy3/clean/internal/modules/game/application/ports"
	"github.com/KyKyPy3/clean/internal/modules/game/domain/entity"
	"github.com/KyKyPy3/clean/pkg/logger"
)

type gamePgStorage struct {
	db     *sqlx.DB
	logger logger.Logger
	tracer trace.Tracer
	getter *trmsqlx.CtxGetter
}

func NewGamePgStorage(db *sqlx.DB, getter *trmsqlx.CtxGetter, logger logger.Logger) ports.GamePgStorage {
	return &gamePgStorage{
		db:     db,
		logger: logger,
		getter: getter,
		tracer: otel.Tracer(""),
	}
}

// Fetch users with given limit.
// TODO: think about offset - use numeric or time offset?
func (g *gamePgStorage) Fetch(ctx context.Context, limit, offset int64) ([]entity.Game, error) {
	ctx, span := g.tracer.Start(ctx, "gamePgStorage.Fetch")
	defer span.End()

	stmt, err := g.getter.DefaultTrOrDB(ctx, g.db).PreparexContext(ctx, FetchSQL)
	if err != nil {
		return nil, errors.Wrap(err, "Fetch.PreparexContext")
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			g.logger.Errorf("[gamePgStorage.Fetch] can't close fetch statement, err: %v", err)
		}
	}()

	rows, err := stmt.QueryxContext(ctx, limit, offset)
	if err != nil || rows.Err() != nil {
		g.logger.Errorf("[gamePgStorage.Fetch] Can't fetch gmae with limit %d and offset %d, err: %v", limit, offset, err)
		return nil, errors.Wrap(err, "Fetch.QueryxContext")
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			g.logger.Errorf("[gamePgStorage.Fetch] Can't close fetched game rows, err: %v", errRow)
		}
	}()

	result := make([]entity.Game, 0)
	for rows.Next() {
		game := DBGame{}

		err = rows.StructScan(&game)
		if err != nil {
			g.logger.Errorf("[gamePgStorage.Fetch] Can't scan game data. err: %v", err)
			return nil, errors.Wrap(err, "Fetch.StructScan")
		}

		var gameEntity entity.Game
		gameEntity, err = GameFromDB(game)
		if err != nil {
			g.logger.Errorf("[gamePgStorage.Fetch] Can't convert game data to domain entity. err: %v", err)
			return nil, errors.Wrap(err, "Fetch.GameFromDB")
		}
		result = append(result, gameEntity)
	}

	return result, nil
}

// Create new game.
func (g *gamePgStorage) Create(ctx context.Context, d entity.Game) error {
	ctx, span := g.tracer.Start(ctx, "gamePgStorage.Create")
	defer span.End()

	stmt, err := g.getter.DefaultTrOrDB(ctx, g.db).PreparexContext(ctx, CreateSQL)
	if err != nil {
		return errors.Wrap(err, "Create.PreparexContext")
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			g.logger.Errorf("can't close create statement, err: %v", err)
		}
	}()

	game := GameToDB(d)
	if err = stmt.QueryRowxContext(
		ctx,
		game.ID,
		game.Name,
	).StructScan(&game); err != nil {
		return errors.Wrap(err, "Create.QueryRowxContext")
	}

	return nil
}
