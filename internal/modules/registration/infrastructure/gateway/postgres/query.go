package postgres

import _ "embed"

var (
	//go:embed query/create.sql
	createSQL string

	//go:embed query/getByEmail.sql
	getByEmailSQL string
)
