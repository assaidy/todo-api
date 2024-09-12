SELECT 
    name,
    email,
    password,
    joined_at
FROM users
WHERE id = ?;
