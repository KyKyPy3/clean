package postgres_test

import (
	"context"
	"testing"

	"github.com/KyKyPy3/clean/internal/common"
	"github.com/KyKyPy3/clean/internal/user/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestFetch(t *testing.T) {
	users, err := repo.Fetch(context.Background(), 10)
	require.NoError(t, err)
	require.Len(t, users, 2)
}

func TestCreate(t *testing.T) {
	user := entity.User{
		Name:       "Bob",
		Surname:    "Smith",
		Middlename: "Joseph",
		Email:      "bob@email.com",
		Password:   "12345",
	}

	user, err := repo.Create(context.Background(), user)
	require.NoError(t, err)
	require.NotNil(t, user.ID)
}

func TestGetByID(t *testing.T) {
	user, err := repo.GetByID(context.Background(), uuid.MustParse("2b0c8791-2136-46b6-bc38-b33038ca2e80"))
	require.NoError(t, err)
	require.NotNil(t, user)

	user, err = repo.GetByID(context.Background(), uuid.MustParse("2b0c1111-2136-46b6-bc38-b33038ca2e80"))
	require.NotNil(t, err)
	require.ErrorIs(t, err, common.ErrNotFound)
	require.Equal(t, entity.User{}, user)
}

func TestGetByEmail(t *testing.T) {
	user, err := repo.GetByEmail(context.Background(), "alise@email.com")
	require.NoError(t, err)
	require.NotNil(t, user)

	user, err = repo.GetByEmail(context.Background(), "rob@email.com")
	require.NotNil(t, err)
	require.ErrorIs(t, err, common.ErrNotFound)
	require.Equal(t, entity.User{}, user)
}

func TestDelete(t *testing.T) {
	err := repo.Delete(context.Background(), uuid.MustParse("2b0c8791-2136-46b6-bc38-b33038ca2e80"))
	require.NoError(t, err)
}
