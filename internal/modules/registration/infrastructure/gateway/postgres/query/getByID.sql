SELECT id, email, password, verified
FROM registrations
WHERE id = $1