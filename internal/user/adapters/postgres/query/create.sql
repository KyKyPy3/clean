INSERT INTO users (name, surname, middlename, email)
VALUES ($1, $2, $3, $4)
RETURNING id, name, surname, middlename, email, created_at, updated_at
