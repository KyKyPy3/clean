package redis

import (
	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/modules/session/domain/entity"
)

// DBToken Database session representation
type DBToken struct {
	ID        string
	UserID    string
	ExpiresIn int64
}

// TokenFromDB Convert database token model to domain model
func TokenFromDB(dbToken DBToken) (entity.Token, error) {
	entityID, err := common.ParseUID(dbToken.ID)
	if err != nil {
		return entity.Token{}, err
	}

	userID, err := common.ParseUID(dbToken.UserID)
	if err != nil {
		return entity.Token{}, err
	}

	token := entity.Hydrate(entityID, userID, dbToken.ExpiresIn)

	return token, nil
}

// TokenToDB Convert domain token model to database model
func TokenToDB(session entity.Token) DBToken {
	return DBToken{
		ID:        session.ID().String(),
		UserID:    session.UserID().String(),
		ExpiresIn: session.ExpiresIn(),
	}
}
