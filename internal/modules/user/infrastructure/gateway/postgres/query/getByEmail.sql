SELECT id, name, surname, middlename, email, password, created_at, updated_at
FROM users
WHERE email = $1
