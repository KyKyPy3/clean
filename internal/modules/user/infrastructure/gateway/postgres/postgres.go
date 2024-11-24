package postgres

import (
	"context"
	"fmt"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/domain/core"
	"github.com/KyKyPy3/clean/internal/modules/user/application/ports"
	"github.com/KyKyPy3/clean/internal/modules/user/domain/entity"
	"github.com/KyKyPy3/clean/pkg/logger"
)

type userPgStorage struct {
	db     *sqlx.DB
	logger logger.Logger
	tracer trace.Tracer
	getter *trmsqlx.CtxGetter
}

func NewUserPgStorage(db *sqlx.DB, getter *trmsqlx.CtxGetter, logger logger.Logger) ports.UserPgStorage {
	return &userPgStorage{
		db:     db,
		logger: logger,
		getter: getter,
		tracer: otel.Tracer(""),
	}
}

// Fetch users with given limit.
// TODO: think about offset - use numeric or time offset?
func (u *userPgStorage) Fetch(ctx context.Context, limit, offset int64) ([]entity.User, error) {
	ctx, span := u.tracer.Start(ctx, "userPgStorage.Fetch")
	defer span.End()

	stmt, err := u.getter.DefaultTrOrDB(ctx, u.db).PreparexContext(ctx, FetchSQL)
	if err != nil {
		return nil, errors.Wrap(err, "Fetch.PreparexContext")
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			u.logger.Errorf("[userPgStorage.Fetch] can't close fetch statement, err: %w", err)
		}
	}()

	rows, err := stmt.QueryxContext(ctx, limit, offset)
	if err != nil || rows.Err() != nil {
		u.logger.Errorf("[userPgStorage.Fetch] Can't fetch user with limit %d and offset %d, err: %w", limit, offset, err)
		return nil, errors.Wrap(err, "Fetch.QueryxContext")
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			u.logger.Errorf("[userPgStorage.Fetch] Can't close fetched user rows, err: %w", errRow)
		}
	}()

	result := make([]entity.User, 0)
	for rows.Next() {
		user := DBUser{}

		err = rows.StructScan(&user)
		if err != nil {
			u.logger.Errorf("[userPgStorage.Fetch] Can't scan user data. err: %w", err)
			return nil, errors.Wrap(err, "Fetch.StructScan")
		}

		var userEntity entity.User
		userEntity, err = UserFromDB(user)
		if err != nil {
			u.logger.Errorf("[userPgStorage.Fetch] Can't convert user data to domain entity. err: %w", err)
			return nil, errors.Wrap(err, "Fetch.UserFromDB")
		}
		result = append(result, userEntity)
	}

	return result, nil
}

// Create new user.
func (u *userPgStorage) Create(ctx context.Context, d entity.User) error {
	ctx, span := u.tracer.Start(ctx, "userPgStorage.Create")
	defer span.End()

	stmt, err := u.getter.DefaultTrOrDB(ctx, u.db).PreparexContext(ctx, CreateSQL)
	if err != nil {
		return errors.Wrap(err, "Create.PreparexContext")
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			u.logger.Errorf("can't close create statement, err: %w", err)
		}
	}()

	user := UserToDB(d)
	if err = stmt.QueryRowxContext(
		ctx,
		user.ID,
		user.Name,
		user.Surname,
		user.Middlename,
		user.Email,
		user.Password,
	).StructScan(&user); err != nil {
		return errors.Wrap(err, "Create.QueryRowxContext")
	}

	return nil
}

func (u *userPgStorage) Update(ctx context.Context, d entity.User) error {
	ctx, span := u.tracer.Start(ctx, "userPgStorage.Update")
	defer span.End()

	stmt, err := u.getter.DefaultTrOrDB(ctx, u.db).PreparexContext(ctx, UpdateSQL)
	if err != nil {
		return errors.Wrap(err, "Update.PreparexContext")
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			u.logger.Errorf("can't close create statement, err: %w", err)
		}
	}()

	user := UserToDB(d)
	if err = stmt.QueryRowxContext(
		ctx,
		user.ID,
		user.Name,
		user.Surname,
		user.Middlename,
		user.Email,
	).StructScan(&user); err != nil {
		return errors.Wrap(err, "Update.QueryRowxContext")
	}

	return nil
}

// GetByID Get user by id.
func (u *userPgStorage) GetByID(ctx context.Context, id common.UID) (entity.User, error) {
	ctx, span := u.tracer.Start(ctx, "userPgStorage.GetByID")
	defer span.End()

	return u.getByParam(ctx, GetByIDSQL, id.GetID())
}

// GetByEmail Get user by email.
func (u *userPgStorage) GetByEmail(ctx context.Context, email common.Email) (entity.User, error) {
	ctx, span := u.tracer.Start(ctx, "userPgStorage.GetByEmail")
	defer span.End()

	return u.getByParam(ctx, GetByEmailSQL, email.String())
}

func (u *userPgStorage) getByParam(ctx context.Context, sqlQuery string, param any) (entity.User, error) {
	stmt, err := u.getter.DefaultTrOrDB(ctx, u.db).PreparexContext(ctx, sqlQuery)
	if err != nil {
		return entity.User{}, errors.Wrap(err, "getByParam.PreparexContext")
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			u.logger.Errorf("can't close getByParam statement, err: %v", err)
		}
	}()

	rows, err := stmt.QueryxContext(ctx, param)
	if err != nil || rows.Err() != nil {
		u.logger.Errorf("Can't fetch user by param %v, err: %v", param, err)
		return entity.User{}, errors.Wrap(err, "getByParam.QueryxContext")
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			u.logger.Errorf("Can't close fetched user rows, err: %v", errRow)
		}
	}()

	result := make([]entity.User, 0)
	for rows.Next() {
		user := DBUser{}

		err = rows.StructScan(&user)
		if err != nil {
			u.logger.Errorf("Can't scan user data. err: %v", err)
			return entity.User{}, errors.Wrap(err, "getByParam.StructScan")
		}

		var userEntity entity.User
		userEntity, err = UserFromDB(user)
		if err != nil {
			u.logger.Errorf("Can't convert user data to domain entity. err: %v", err)
			return entity.User{}, errors.Wrap(err, "getByParam.UserFromDB")
		}

		result = append(result, userEntity)
	}

	if len(result) == 0 {
		return entity.User{}, core.ErrNotFound
	}

	return result[0], nil
}

// Delete user by provided id.
func (u *userPgStorage) Delete(ctx context.Context, id common.UID) error {
	ctx, span := u.tracer.Start(ctx, "userPgStorage.Delete")
	defer span.End()

	stmt, err := u.getter.DefaultTrOrDB(ctx, u.db).PreparexContext(ctx, DeleteSQL)
	if err != nil {
		return errors.Wrap(err, "Delete.PreparexContext")
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			u.logger.Errorf("can't close delete statement, err: %w", err)
		}
	}()

	res, err := stmt.ExecContext(ctx, id.GetID())
	if err != nil {
		return errors.Wrap(err, "Delete.ExecContext")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "Delete.RowsAffected")
	}

	if rowsAffected != 1 {
		return fmt.Errorf("delete operation affected %d row", rowsAffected)
	}

	return nil
}
