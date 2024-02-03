INSERT INTO users (id, name, surname, middlename, email, password)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id
