INSERT INTO users (name, email, password) 
VALUES (?, ?, ?) 
RETURNING id, joined_at;
