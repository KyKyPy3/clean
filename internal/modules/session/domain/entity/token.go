package entity

import (
	"fmt"

	"github.com/KyKyPy3/clean/internal/domain/common"
)

type Token struct {
	id        common.UID
	userID    common.UID
	expiresIn int64
}

func NewToken(userID common.UID, expiresIn int64) Token {
	return Token{
		id:        common.NewUID(),
		userID:    userID,
		expiresIn: expiresIn,
	}
}

func Hydrate(tokenID, userID common.UID, expiresIn int64) Token {
	return Token{
		id:        tokenID,
		userID:    userID,
		expiresIn: expiresIn,
	}
}

func (t *Token) ID() common.UID {
	return t.id
}

func (t *Token) UserID() common.UID {
	return t.userID
}

func (t *Token) ExpiresIn() int64 {
	return t.expiresIn
}

func (t *Token) IsEmpty() bool {
	return *t == Token{}
}

func (t *Token) String() string {
	return fmt.Sprintf(
		"Token{ID: %s, UserID: %s, ExpiresIn: %d}",
		t.ID(),
		t.UserID(),
		t.ExpiresIn(),
	)
}
