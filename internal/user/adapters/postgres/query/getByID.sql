SELECT id, name, surname, middlename, email, created_at, updated_at
FROM users
WHERE id = $1
