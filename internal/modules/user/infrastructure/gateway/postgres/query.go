package postgres

import _ "embed"

var (
	//go:embed query/fetch.sql
	fetchSQL string

	//go:embed query/create.sql
	createSQL string

	//go:embed query/getByEmail.sql
	getByEmailSQL string

	//go:embed query/getByID.sql
	getByIDSQL string

	//go:embed query/delete.sql
	deleteSQL string
)
