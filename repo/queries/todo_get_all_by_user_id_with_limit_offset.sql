-- if limit is negative (-1) it will ignore the limitation 
-- and return all possible rows
SELECT
    id,
    title,
    description,
    created_at,
    status
FROM todos
WHERE user_id = $1
ORDER BY created_at DESC -- newest first
LIMIT $2
OFFSET $3;
