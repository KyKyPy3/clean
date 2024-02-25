package dto

import "github.com/KyKyPy3/clean/internal/modules/game/domain/entity"

type GameDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updateAt"`
}

type CreateGameDTO struct {
	Name string `json:"name" validate:"required"`
}

type FetchGamesDTO struct {
	Limit  int64 `query:"limit" validate:"gte=0,lte=1000"`
	Offset int64 `query:"offset" validate:"gte=0,lte=1000"`
}

// GameToResponse - Convert domain game model to response model.
func GameToResponse(user entity.Game) GameDTO {
	return GameDTO{
		ID:        user.ID().String(),
		Name:      user.Name(),
		CreatedAt: user.CreatedAt().String(),
		UpdatedAt: user.UpdatedAt().String(),
	}
}
