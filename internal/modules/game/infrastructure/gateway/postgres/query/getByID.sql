SELECT id,
       name,
       owner_id,
       created_at,
       updated_at
FROM games
WHERE id = $1