UPDATE users
SET name = $2,
    surname = $3,
    middlename = $4,
    email = $5
WHERE id = $1
RETURNING id