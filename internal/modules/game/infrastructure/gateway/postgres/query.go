package postgres

import _ "embed"

var (
	//go:embed query/fetch.sql
	FetchSQL string

	//go:embed query/create.sql
	CreateSQL string
)
