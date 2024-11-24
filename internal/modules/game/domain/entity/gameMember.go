package entity

import (
	"github.com/KyKyPy3/clean/internal/domain/common"
)

// GameMember struct.
type GameMember struct {
	id     common.UID
	name   string
	userID common.UID
	gameID common.UID
}

func NewGroupMember(name string, userID, gameID common.UID) (GameMember, error) {
	return GameMember{
		id:     common.NewUID(),
		name:   name,
		userID: userID,
		gameID: gameID,
	}, nil
}
