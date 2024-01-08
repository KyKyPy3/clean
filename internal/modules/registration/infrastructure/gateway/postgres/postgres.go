package postgres

import (
	"context"
	"github.com/KyKyPy3/clean/internal/domain"
	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/modules/registration/domain/service"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/KyKyPy3/clean/internal/modules/registration/domain/entity"
	"github.com/KyKyPy3/clean/pkg/logger"
)

type registrationPgStorage struct {
	db     *sqlx.DB
	logger logger.Logger
	tracer trace.Tracer
	getter *trmsqlx.CtxGetter
}

func NewRegistrationPgStorage(db *sqlx.DB, getter *trmsqlx.CtxGetter, logger logger.Logger) service.RegistrationPgStorage {
	return &registrationPgStorage{
		db:     db,
		logger: logger,
		getter: getter,
		tracer: otel.Tracer(""),
	}
}

// Create new user
func (r *registrationPgStorage) Create(ctx context.Context, d entity.Registration) error {
	ctx, span := r.tracer.Start(ctx, "registrationPgStorage.Create")
	defer span.End()

	stmt, err := r.getter.DefaultTrOrDB(ctx, r.db).PreparexContext(ctx, createSQL)
	if err != nil {
		return errors.Wrap(err, "Create.PreparexContext")
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
			r.logger.Errorf("can't close create statement, err: %w", err)
		}
	}()

	registration := RegistrationToDB(d)
	r.logger.Debugf("Save registration %v", registration)

	if _, err := stmt.ExecContext(
		ctx,
		registration.ID,
		registration.Email,
	); err != nil {
		return errors.Wrap(err, "Create.QueryRowxContext")
	}

	return nil
}

// GetByEmail Get registration by email
func (r *registrationPgStorage) GetByEmail(ctx context.Context, email common.Email) (entity.Registration, error) {
	ctx, span := r.tracer.Start(ctx, "registrationPgStorage.GetByEmail")
	defer span.End()

	stmt, err := r.getter.DefaultTrOrDB(ctx, r.db).PreparexContext(ctx, getByEmailSQL)
	if err != nil {
		return entity.Registration{}, errors.Wrap(err, "GetByEmail.PreparexContext")
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
			r.logger.Errorf("can't close getByEmail statement, err: %w", err)
		}
	}()

	rows, err := stmt.QueryxContext(ctx, email.String())
	if err != nil {
		r.logger.Errorf("Can't fetch registration by email, err: %w", err)
		return entity.Registration{}, errors.Wrap(err, "GetByEmail.QueryxContext")
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			r.logger.Errorf("Can't close fetched registration rows, err: %w", errRow)
		}
	}()

	result := make([]entity.Registration, 0)
	for rows.Next() {
		registration := DbRegistration{}

		err = rows.StructScan(&registration)
		if err != nil {
			r.logger.Errorf("Can't scan registration data. err: %w", err)
			return entity.Registration{}, errors.Wrap(err, "GetByEmail.StructScan")
		}

		registrationEntity, err := RegistrationFromDB(registration)
		if err != nil {
			r.logger.Errorf("Can't convert registration data to domain entity. err: %w", err)
			return entity.Registration{}, errors.Wrap(err, "GetByEmail.RegistrationFromDB")
		}

		result = append(result, registrationEntity)
	}

	if len(result) == 0 {
		return entity.Registration{}, domain.ErrNotFound
	}

	return result[0], nil
}
