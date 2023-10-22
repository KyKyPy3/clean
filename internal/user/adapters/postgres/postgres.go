package postgres

import (
	"context"
	"fmt"
	"github.com/KyKyPy3/clean/internal/common"
	"github.com/KyKyPy3/clean/internal/user/domain/entity"
	"github.com/KyKyPy3/clean/internal/user/domain/service"
	"github.com/KyKyPy3/clean/pkg/logger"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
)

type userPgStorage struct {
	db     *sqlx.DB
	logger logger.Logger
}

func NewUserPgStorage(db *sqlx.DB, logger logger.Logger) service.UserPgStorage {
	return &userPgStorage{db: db, logger: logger}
}

// Fetch users with given limit
// TODO: think about offset - use numeric or time offset?
func (u *userPgStorage) Fetch(ctx context.Context, limit int64) ([]entity.User, error) {
	ctx, span := otel.Tracer("").Start(ctx, "userPgStorage.Fetch")
	defer span.End()

	stmt, err := u.db.PreparexContext(ctx, fetchSQL)
	if err != nil {
		return []entity.User{}, errors.Wrap(err, "Fetch.PreparexContext")
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
			u.logger.Errorf("can't close fetch statement, err: %w", err)
		}
	}()

	rows, err := stmt.QueryxContext(ctx, limit)
	if err != nil {
		u.logger.Errorf("Can't fetch user with limit %d, err: %w", limit, err)
		return []entity.User{}, errors.Wrap(err, "Fetch.QueryxContext")
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			u.logger.Errorf("Can't close fetched user rows, err: %w", errRow)
		}
	}()

	result := make([]entity.User, 0)
	for rows.Next() {
		user := DbUser{}

		err = rows.StructScan(&user)

		result = append(result, UserFromDB(user))
	}

	return result, nil
}

// Create new user
func (u *userPgStorage) Create(ctx context.Context, d entity.User) (entity.User, error) {
	ctx, span := otel.Tracer("").Start(ctx, "userPgStorage.Create")
	defer span.End()

	stmt, err := u.db.PreparexContext(ctx, createSQL)
	if err != nil {
		return entity.User{}, errors.Wrap(err, "Create.PreparexContext")
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
			u.logger.Errorf("can't close create statement, err: %w", err)
		}
	}()

	user := UserToDB(d)
	if err := stmt.QueryRowxContext(
		ctx,
		user.Name,
		user.Surname,
		user.Middlename,
		user.Email,
	).StructScan(user); err != nil {
		return entity.User{}, errors.Wrap(err, "Create.QueryRowxContext")
	}

	return UserFromDB(user), nil
}

// Get user by id
func (u *userPgStorage) GetByID(ctx context.Context, id common.ID) (entity.User, error) {
	ctx, span := otel.Tracer("").Start(ctx, "userPgStorage.GetByID")
	defer span.End()

	stmt, err := u.db.PreparexContext(ctx, getByIDSQL)
	if err != nil {
		return entity.User{}, errors.Wrap(err, "GetByID.PreparexContext")
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
			u.logger.Errorf("can't close getById statement, err: %w", err)
		}
	}()

	rows, err := stmt.QueryxContext(ctx, id)
	if err != nil {
		u.logger.Errorf("Can't fetch user by id, err: %w", err)
		return entity.User{}, errors.Wrap(err, "GetByID.QueryxContext")
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			u.logger.Errorf("Can't close fetched user rows, err: %w", errRow)
		}
	}()

	result := make([]entity.User, 0)
	for rows.Next() {
		user := DbUser{}

		err = rows.StructScan(&user)

		result = append(result, UserFromDB(user))
	}

	if len(result) == 0 {
		return entity.User{}, common.ErrNotFound
	}

	return result[0], nil
}

// Get user by email
func (u *userPgStorage) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	ctx, span := otel.Tracer("").Start(ctx, "userPgStorage.GetByEmail")
	defer span.End()

	stmt, err := u.db.PreparexContext(ctx, getByEmailSQL)
	if err != nil {
		return entity.User{}, errors.Wrap(err, "GetByEmail.PreparexContext")
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
			u.logger.Errorf("can't close getByEmail statement, err: %w", err)
		}
	}()

	rows, err := stmt.QueryxContext(ctx, email)
	if err != nil {
		u.logger.Errorf("Can't fetch user by email, err: %w", err)
		return entity.User{}, errors.Wrap(err, "GetByEmail.QueryxContext")
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			u.logger.Errorf("Can't close fetched user rows, err: %w", errRow)
		}
	}()

	result := make([]entity.User, 0)
	for rows.Next() {
		user := DbUser{}

		err = rows.StructScan(&user)

		result = append(result, UserFromDB(user))
	}

	if len(result) == 0 {
		return entity.User{}, common.ErrNotFound
	}

	return result[0], nil
}

// Delete user by provided id
func (u *userPgStorage) Delete(ctx context.Context, id common.ID) error {
	ctx, span := otel.Tracer("").Start(ctx, "userPgStorage.Delete")
	defer span.End()

	stmt, err := u.db.PreparexContext(ctx, deleteSQL)
	if err != nil {
		return errors.Wrap(err, "Delete.PreparexContext")
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
			u.logger.Errorf("can't close delete statement, err: %w", err)
		}
	}()

	res, err := stmt.ExecContext(ctx, id)
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
