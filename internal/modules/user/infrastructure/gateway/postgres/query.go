package postgres

import _ "embed"

var (
	//go:embed query/fetch.sql
	FetchSQL string

	//go:embed query/create.sql
	CreateSQL string

	//go:embed query/update.sql
	UpdateSQL string

	//go:embed query/getByEmail.sql
	GetByEmailSQL string

	//go:embed query/getByID.sql
	GetByIDSQL string

	//go:embed query/delete.sql
	DeleteSQL string
)
