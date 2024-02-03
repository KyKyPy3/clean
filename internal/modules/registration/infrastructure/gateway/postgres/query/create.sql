INSERT INTO registrations (id, email, password)
VALUES ($1, $2, $3)
RETURNING id