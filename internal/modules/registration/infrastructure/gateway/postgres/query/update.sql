UPDATE registrations
SET verified = $2
WHERE id = $1
RETURNING id