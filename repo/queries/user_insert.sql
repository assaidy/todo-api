INSERT INTO users (name, email, password) 
VALUES (?, ?, ?) 
RETURNING id;
