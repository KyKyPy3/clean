SELECT
    id,
    name,
    surname,
    middlename,
    email,
    created_at,
    updated_at
FROM users
ORDER BY created_at
LIMIT $1 OFFSET $2
