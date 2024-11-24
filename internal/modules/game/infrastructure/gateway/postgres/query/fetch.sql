SELECT
    id,
    name,
    owner_id,
    created_at,
    updated_at
FROM games
ORDER BY created_at
LIMIT $1 OFFSET $2