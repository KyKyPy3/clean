INSERT INTO games (id, name)
VALUES ($1, $2)
RETURNING id