package postgres_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/domain/core"
	"github.com/KyKyPy3/clean/internal/modules/user/domain/entity"
	"github.com/KyKyPy3/clean/internal/modules/user/domain/value_object"
)

func TestFetch(t *testing.T) {
	users, err := repo.Fetch(context.Background(), 10, 0)
	require.NoError(t, err)
	require.Len(t, users, 2)
}

func TestCreate(t *testing.T) {
	fullName, _ := value_object.NewFullName("Bob", "Smith", "Joseph")
	email, _ := common.NewEmail("bob@email.com")

	user, _ := entity.NewUser(fullName, email, "12345", policy)

	err := repo.Create(context.Background(), user)
	require.NoError(t, err)
	require.NotNil(t, user.ID())
}

func TestGetByID(t *testing.T) {
	user, err := repo.GetByID(context.Background(), common.NewWithSpecifiedID(uuid.MustParse("2b0c8791-2136-46b6-bc38-b33038ca2e80")))
	require.NoError(t, err)
	require.NotNil(t, user)

	user, err = repo.GetByID(context.Background(), common.NewWithSpecifiedID(uuid.MustParse("2b0c1111-2136-46b6-bc38-b33038ca2e80")))
	require.NotNil(t, err)
	require.ErrorIs(t, err, core.ErrNotFound)
	require.Equal(t, entity.User{}, user)
}

func TestGetByEmail(t *testing.T) {
	aliceEmail, _ := common.NewEmail("alise@email.com")
	user, err := repo.GetByEmail(context.Background(), aliceEmail)
	require.NoError(t, err)
	require.NotNil(t, user)

	bobEmail, _ := common.NewEmail("rob@email.com")
	user, err = repo.GetByEmail(context.Background(), bobEmail)
	require.NotNil(t, err)
	require.ErrorIs(t, err, core.ErrNotFound)
	require.Equal(t, entity.User{}, user)
}

func TestDelete(t *testing.T) {
	err := repo.Delete(context.Background(), common.NewWithSpecifiedID(uuid.MustParse("2b0c8791-2136-46b6-bc38-b33038ca2e80")))
	require.NoError(t, err)
}
