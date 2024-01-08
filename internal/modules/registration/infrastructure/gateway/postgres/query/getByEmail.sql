SELECT id, email, verified
FROM registrations
WHERE email = $1
