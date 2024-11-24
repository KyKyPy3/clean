package postgres

import _ "embed"

var (
	//go:embed query/fetch.sql
	FetchSQL string

	//go:embed query/create.sql
	CreateSQL string

	//go:embed query/addUser.sql
	AddUserSQL string

	//go:embed query/getByID.sql
	GetByIDSQL string
)
