package postgres

import (
	"context"
	"github.com/KyKyPy3/clean/internal/modules/registration/application/ports"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/domain/core"
	"github.com/KyKyPy3/clean/internal/modules/registration/domain/entity"
	"github.com/KyKyPy3/clean/pkg/logger"
)

type registrationPgStorage struct {
	db     *sqlx.DB
	logger logger.Logger
	tracer trace.Tracer
	getter *trmsqlx.CtxGetter
}

func NewRegistrationPgStorage(db *sqlx.DB, getter *trmsqlx.CtxGetter, logger logger.Logger) ports.RegistrationPgStorage {
	return &registrationPgStorage{
		db:     db,
		logger: logger,
		getter: getter,
		tracer: otel.Tracer(""),
	}
}

// Create new registration
func (r *registrationPgStorage) Create(ctx context.Context, d entity.Registration) error {
	ctx, span := r.tracer.Start(ctx, "registrationPgStorage.Create")
	defer span.End()

	stmt, err := r.getter.DefaultTrOrDB(ctx, r.db).PreparexContext(ctx, createSQL)
	if err != nil {
		return errors.Wrap(err, "[registrationPgStorage.Create] PreparexContext")
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
			r.logger.Errorf("[registrationPgStorage.Create] can't close create statement, err: %w", err)
		}
	}()

	registration := RegistrationToDB(d)
	r.logger.Debugf("[registrationPgStorage.Create] Save registration %v", registration)

	if _, err := stmt.ExecContext(
		ctx,
		registration.ID,
		registration.Email,
	); err != nil {
		return errors.Wrap(err, "[registrationPgStorage.Create] QueryRowxContext")
	}

	return nil
}

func (r *registrationPgStorage) Update(ctx context.Context, d entity.Registration) error {
	ctx, span := r.tracer.Start(ctx, "registrationPgStorage.Update")
	defer span.End()

	stmt, err := r.getter.DefaultTrOrDB(ctx, r.db).PreparexContext(ctx, updateSQL)
	if err != nil {
		return errors.Wrap(err, "[registrationPgStorage.Update] PreparexContext")
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
			r.logger.Errorf("[registrationPgStorage.Update] can't close create statement, err: %w", err)
		}
	}()

	registration := RegistrationToDB(d)
	r.logger.Debugf("[registrationPgStorage.Update] Save registration %v", registration)

	if _, err := stmt.ExecContext(
		ctx,
		registration.ID,
		true,
	); err != nil {
		return errors.Wrap(err, "[registrationPgStorage.Update] QueryRowxContext")
	}

	return nil
}

func (r *registrationPgStorage) GetByID(ctx context.Context, id common.UID) (entity.Registration, error) {
	ctx, span := r.tracer.Start(ctx, "registrationPgStorage.GetByID")
	defer span.End()

	stmt, err := r.getter.DefaultTrOrDB(ctx, r.db).PreparexContext(ctx, getByIDSQL)
	if err != nil {
		return entity.Registration{}, errors.Wrap(err, "[registrationPgStorage.GetByID] PreparexContext")
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
			r.logger.Errorf("[registrationPgStorage.GetByID] can't close statement, err: %w", err)
		}
	}()

	rows, err := stmt.QueryxContext(ctx, id.String())
	if err != nil {
		r.logger.Errorf("[registrationPgStorage.GetByID] Can't fetch registration by id, err: %w", err)
		return entity.Registration{}, errors.Wrap(err, "[registrationPgStorage.GetByID] QueryxContext")
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			r.logger.Errorf("[registrationPgStorage.GetByID] Can't close fetched registration rows, err: %w", errRow)
		}
	}()

	result := make([]entity.Registration, 0)
	for rows.Next() {
		registration := DBRegistration{}

		err = rows.StructScan(&registration)
		if err != nil {
			r.logger.Errorf("[registrationPgStorage.GetByID] Can't scan registration data. err: %w", err)
			return entity.Registration{}, errors.Wrap(err, "[registrationPgStorage.GetByID] StructScan")
		}

		registrationEntity, err := RegistrationFromDB(registration)
		if err != nil {
			r.logger.Errorf("[registrationPgStorage.GetByID] Can't convert registration data to domain entity. err: %w", err)
			return entity.Registration{}, errors.Wrap(err, "[registrationPgStorage.GetByID] RegistrationFromDB")
		}

		result = append(result, registrationEntity)
	}

	if len(result) == 0 {
		return entity.Registration{}, core.ErrNotFound
	}

	return result[0], nil
}
