SELECT
    id,
    title,
    description,
    status,
    created_at
FROM todos
WHERE user_id = ?
ORDER BY created_at DESC; -- newest first

