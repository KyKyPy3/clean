package postgres_test

import (
	"context"
	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/modules/user/domain/value_object"
	"testing"

	"github.com/KyKyPy3/clean/internal/domain"
	"github.com/KyKyPy3/clean/internal/modules/user/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestFetch(t *testing.T) {
	users, err := repo.Fetch(context.Background(), 10, 0)
	require.NoError(t, err)
	require.Len(t, users, 2)
}

func TestCreate(t *testing.T) {
	fullName, _ := value_object.NewFullName("Bob", "Smith", "Joseph")
	email, _ := common.NewEmail("bob@email.com")
	user := entity.User{}
	user.SetID(common.NewUID())
	user.SetFullName(fullName)
	user.SetEmail(email)
	user.SetPassword("12345")

	user, err := repo.Create(context.Background(), user)
	require.NoError(t, err)
	require.NotNil(t, user.GetID())
}

func TestGetByID(t *testing.T) {
	user, err := repo.GetByID(context.Background(), common.NewWithSpecifiedID(uuid.MustParse("2b0c8791-2136-46b6-bc38-b33038ca2e80")))
	require.NoError(t, err)
	require.NotNil(t, user)

	user, err = repo.GetByID(context.Background(), common.NewWithSpecifiedID(uuid.MustParse("2b0c1111-2136-46b6-bc38-b33038ca2e80")))
	require.NotNil(t, err)
	require.ErrorIs(t, err, domain.ErrNotFound)
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
	require.ErrorIs(t, err, domain.ErrNotFound)
	require.Equal(t, entity.User{}, user)
}

func TestDelete(t *testing.T) {
	err := repo.Delete(context.Background(), common.NewWithSpecifiedID(uuid.MustParse("2b0c8791-2136-46b6-bc38-b33038ca2e80")))
	require.NoError(t, err)
}
