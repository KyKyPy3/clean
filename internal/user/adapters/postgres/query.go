package postgres

const (
	fetchSQL = `
		SELECT id, name, surname, middlename, email, created_at, updated_at
		FROM users
		ORDER BY created_at LIMIT $1
	`
	createSQL = `
		INSERT INTO users (name, surname, middlename, email)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, surname, middlename, email, created_at, updated_at
	`
	getByEmailSQL = `
		SELECT id, name, surname, middlename, email, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	getByIDSQL = `
		SELECT id, name, surname, middlename, email, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	deleteSQL = `
		DELETE FROM users
		WHERE id = $1
	`
)
