package postgres_test

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KyKyPy3/clean/internal/modules/user/infrastructure/gateway/postgres"
	"github.com/KyKyPy3/clean/pkg/logger"
)

func TestFetch(t *testing.T) {
	// Create logger
	// TODO: add discard logger here
	log := logger.NewLogger(logger.Config{
		Mode: "test",
	})
	log.Init()

	t.Run("Successfully retrieve users", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer func() {
			_ = mockDB.Close()
		}()
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

		repo := postgres.NewUserPgStorage(sqlxDB, trmsqlx.DefaultCtxGetter, log)

		var limit int64 = 10
		var offset int64

		userID := uuid.New().String()
		name := "Ivan"
		surname := "Ivanov"
		middlename := "Ivanovich"
		email := "ivan@email.com"
		createdAt := time.Now()
		updatedAt := time.Now()

		rows := mock.
			NewRows([]string{
				"id",
				"name",
				"surname",
				"middlename",
				"email",
				"created_at",
				"updated_at",
			}).
			AddRow(
				userID,
				name,
				surname,
				middlename,
				email,
				createdAt,
				updatedAt,
			)

		mock.ExpectPrepare(regexp.QuoteMeta(postgres.FetchSQL))
		mock.ExpectQuery(regexp.QuoteMeta(postgres.FetchSQL)).
			WithArgs(limit, offset).
			WillReturnRows(rows)

		users, err := repo.Fetch(context.TODO(), limit, offset)
		require.NoError(t, err)

		assert.Len(t, users, 1)

		assert.Equal(t, users[0].ID().String(), userID)
		assert.Equal(t, users[0].FullName().FirstName(), name)
		assert.Equal(t, users[0].FullName().LastName(), surname)
		assert.Equal(t, users[0].FullName().MiddleName(), middlename)
		assert.Equal(t, users[0].Email().String(), email)
		assert.WithinDuration(t, users[0].CreatedAt(), createdAt, 10*time.Millisecond)
		assert.WithinDuration(t, users[0].UpdatedAt(), updatedAt, 10*time.Millisecond)

		// ensure that all expectations are met in the mock
		errExpectations := mock.ExpectationsWereMet()
		assert.NoError(t, errExpectations)
	})

	t.Run("Successfully retrieve empty result", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer func() {
			_ = mockDB.Close()
		}()
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

		repo := postgres.NewUserPgStorage(sqlxDB, trmsqlx.DefaultCtxGetter, log)

		var limit int64 = 10
		var offset int64

		rows := mock.
			NewRows([]string{
				"id",
				"name",
				"surname",
				"middlename",
				"email",
				"created_at",
				"updated_at",
			})

		mock.ExpectPrepare(regexp.QuoteMeta(postgres.FetchSQL))
		mock.ExpectQuery(regexp.QuoteMeta(postgres.FetchSQL)).
			WithArgs(limit, offset).
			WillReturnRows(rows)

		users, err := repo.Fetch(context.TODO(), limit, offset)
		require.NoError(t, err)

		assert.Empty(t, users)

		// ensure that all expectations are met in the mock
		errExpectations := mock.ExpectationsWereMet()
		assert.NoError(t, errExpectations)
	})

	t.Run("Failed retrieve users", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer func() {
			_ = mockDB.Close()
		}()
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

		repo := postgres.NewUserPgStorage(sqlxDB, trmsqlx.DefaultCtxGetter, log)

		var limit int64 = 10
		var offset int64

		mock.ExpectPrepare(regexp.QuoteMeta(postgres.FetchSQL))
		mock.ExpectQuery(regexp.QuoteMeta(postgres.FetchSQL)).
			WithArgs(limit, offset).
			WillReturnError(sql.ErrNoRows)

		users, err := repo.Fetch(context.TODO(), limit, offset)
		require.ErrorIs(t, err, sql.ErrNoRows)
		assert.Empty(t, users)

		// ensure that all expectations are met in the mock
		errExpectations := mock.ExpectationsWereMet()
		assert.NoError(t, errExpectations)
	})

	t.Run("Failed prepare statement", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer func() {
			_ = mockDB.Close()
		}()
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

		repo := postgres.NewUserPgStorage(sqlxDB, trmsqlx.DefaultCtxGetter, log)

		mock.ExpectPrepare(regexp.QuoteMeta(postgres.FetchSQL)).
			WillReturnError(sql.ErrConnDone)

		users, err := repo.Fetch(context.TODO(), int64(10), int64(0))
		require.ErrorIs(t, err, sql.ErrConnDone)
		assert.Empty(t, users)

		// ensure that all expectations are met in the mock
		errExpectations := mock.ExpectationsWereMet()
		assert.NoError(t, errExpectations)
	})

	t.Run("Failed scan user", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer func() {
			_ = mockDB.Close()
		}()
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

		repo := postgres.NewUserPgStorage(sqlxDB, trmsqlx.DefaultCtxGetter, log)

		var limit int64 = 10
		var offset int64

		rows := mock.
			NewRows([]string{
				"id",
				"name",
				"surname",
				"middlename",
				"email",
				"created_at",
				"updated_at",
			}).
			AddRow(
				nil,
				nil,
				nil,
				nil,
				nil,
				nil,
				nil,
			)

		mock.ExpectPrepare(regexp.QuoteMeta(postgres.FetchSQL))
		mock.ExpectQuery(regexp.QuoteMeta(postgres.FetchSQL)).
			WithArgs(limit, offset).
			WillReturnRows(rows)

		_, err = repo.Fetch(context.TODO(), limit, offset)
		require.Error(t, err)

		// ensure that all expectations are met in the mock
		errExpectations := mock.ExpectationsWereMet()
		assert.NoError(t, errExpectations)
	})
}
