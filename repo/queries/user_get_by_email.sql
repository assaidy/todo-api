SELECT 
    id,
    name,
    password,
    joined_at
FROM users
WHERE email = $1;
