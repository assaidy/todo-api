INSERT INTO todos (user_id, title, description, status, created_at)
VALUES (?, ?, ?, ?, ?)
RETURNING id;
