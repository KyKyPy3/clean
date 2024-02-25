SELECT
    id,
    name,
    created_at,
    updated_at
FROM games
ORDER BY created_at
LIMIT $1 OFFSET $2