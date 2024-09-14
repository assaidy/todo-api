SELECT 1
FROM todo
WHERE id = $1 AND user_id = $2
LIMIT 1;
