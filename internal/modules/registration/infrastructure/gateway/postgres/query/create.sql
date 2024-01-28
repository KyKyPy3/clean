INSERT INTO registrations (id, email)
VALUES ($1, $2)
RETURNING id