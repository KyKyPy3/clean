INSERT INTO users (id, name, surname, middlename, email)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, name, surname, middlename, email, created_at, updated_at
