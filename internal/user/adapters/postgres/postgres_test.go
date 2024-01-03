package postgres

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/KyKyPy3/clean/config"
	"github.com/KyKyPy3/clean/pkg/logger"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestFetch(t *testing.T) {
	// Create logger
	// TODO: add discard logger here
	loggerCfg := &config.LoggerConfig{Mode: "test"}
	logger := logger.NewLogger(loggerCfg)
	logger.Init()

	t.Run("Successfully retrieve users", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

		repo := NewUserPgStorage(sqlxDB, logger)

		var limit int64 = 10
		var offset int64 = 0

		userID := uuid.New()
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

		mock.ExpectPrepare(regexp.QuoteMeta(fetchSQL))
		mock.ExpectQuery(regexp.QuoteMeta(fetchSQL)).
			WithArgs(limit, offset).
			WillReturnRows(rows)

		users, err := repo.Fetch(context.TODO(), limit, offset)
		assert.NoError(t, err)

		assert.Len(t, users, 1)

		assert.Equal(t, users[0].ID(), userID)
		assert.Equal(t, users[0].FirstName(), name)
		assert.Equal(t, users[0].LastName(), surname)
		assert.Equal(t, users[0].MiddleName(), middlename)
		assert.Equal(t, users[0].Email(), email)
		assert.WithinDuration(t, users[0].CreatedAt(), createdAt, 10*time.Millisecond)
		assert.WithinDuration(t, users[0].UpdatedAt(), updatedAt, 10*time.Millisecond)

		// ensure that all expectations are met in the mock
		errExpectations := mock.ExpectationsWereMet()
		assert.Nil(t, errExpectations)
	})

	t.Run("Successfully retrieve empty result", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

		repo := NewUserPgStorage(sqlxDB, logger)

		var limit int64 = 10
		var offset int64 = 0

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

		mock.ExpectPrepare(regexp.QuoteMeta(fetchSQL))
		mock.ExpectQuery(regexp.QuoteMeta(fetchSQL)).
			WithArgs(limit, offset).
			WillReturnRows(rows)

		users, err := repo.Fetch(context.TODO(), limit, offset)
		assert.NoError(t, err)

		assert.Len(t, users, 0)

		// ensure that all expectations are met in the mock
		errExpectations := mock.ExpectationsWereMet()
		assert.Nil(t, errExpectations)
	})

	t.Run("Failed retrieve users", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

		repo := NewUserPgStorage(sqlxDB, logger)

		var limit int64 = 10
		var offset int64 = 0

		mock.ExpectPrepare(regexp.QuoteMeta(fetchSQL))
		mock.ExpectQuery(regexp.QuoteMeta(fetchSQL)).
			WithArgs(limit, offset).
			WillReturnError(sql.ErrNoRows)

		users, err := repo.Fetch(context.TODO(), limit, offset)
		assert.ErrorIs(t, err, sql.ErrNoRows)
		assert.Len(t, users, 0)

		// ensure that all expectations are met in the mock
		errExpectations := mock.ExpectationsWereMet()
		assert.Nil(t, errExpectations)
	})

	t.Run("Failed prepare statement", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

		repo := NewUserPgStorage(sqlxDB, logger)

		mock.ExpectPrepare(regexp.QuoteMeta(fetchSQL)).
			WillReturnError(sql.ErrConnDone)

		users, err := repo.Fetch(context.TODO(), int64(10), int64(0))
		assert.ErrorIs(t, err, sql.ErrConnDone)
		assert.Len(t, users, 0)

		// ensure that all expectations are met in the mock
		errExpectations := mock.ExpectationsWereMet()
		assert.Nil(t, errExpectations)
	})

	t.Run("Failed scan user", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

		repo := NewUserPgStorage(sqlxDB, logger)

		var limit int64 = 10
		var offset int64 = 0

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

		mock.ExpectPrepare(regexp.QuoteMeta(fetchSQL))
		mock.ExpectQuery(regexp.QuoteMeta(fetchSQL)).
			WithArgs(limit, offset).
			WillReturnRows(rows)

		_, err = repo.Fetch(context.TODO(), limit, offset)
		assert.Error(t, err)

		// ensure that all expectations are met in the mock
		errExpectations := mock.ExpectationsWereMet()
		assert.Nil(t, errExpectations)
	})
}
