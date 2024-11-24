INSERT INTO games (id, name, owner_id)
VALUES ($1, $2, $3)
RETURNING id