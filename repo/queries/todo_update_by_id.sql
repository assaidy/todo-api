UPDATE todos
SET 
    title = $1,
    description = $2,
    status = $3
WHERE id = $4 AND user_id = $5;
