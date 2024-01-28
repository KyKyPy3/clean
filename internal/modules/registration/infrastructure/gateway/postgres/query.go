package postgres

import _ "embed"

var (
	//go:embed query/create.sql
	createSQL string

	//go:embed query/getByID.sql
	getByIDSQL string

	//go:embed query/update.sql
	updateSQL string
)
