SELECT
    id,
    title,
    description,
    status,
    created_at
FROM todos
WHERE user_id = $1
ORDER BY created_at DESC; -- newest first
