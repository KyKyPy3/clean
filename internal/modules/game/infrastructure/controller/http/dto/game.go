package dto

type GameShortDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	User      string `json:"user"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updateAt"`
}

type GameDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	User      string `json:"user"`
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

type FetchGameDTO struct {
	ID string `json:"id"`
}
