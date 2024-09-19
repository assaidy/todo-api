package repo

// QM -> returns many
// QO -> returns one
// QE -> excute without return

// user ops
const (
	QOInsertUser = `
    INSERT INTO users (name, email, password) 
    VALUES ($1, $2, $3) 
    RETURNING id;`

	QMGetUserById = `
    SELECT 
        name,
        email,
        password,
        joined_at
    FROM users
    WHERE id = $1;`

	QMGetUserByEmail = `
    SELECT 
        id,
        name,
        password,
        joined_at
    FROM users
    WHERE email = $1;`

	QEUpdateUser = `
    UPDATE users 
    SET 
        name = $1,
        email = $2,
        password = $3 
    WHERE id = $4;`

	QEDeleteUser = `
    DELETE FROM users 
    WHERE id = $1;`

	QOCheckUserIdExists = `
    SELECT 1
    FROM users
    WHERE id = $1
    LIMIT 1;`

	QOCheckEmailExists = `
    SELECT 1 
    FROM users 
    WHERE email = $1 
    LIMIT 1;`
)

// todo ops
const (
	QOInsertTodo = `
    INSERT INTO todos (user_id, title, description, status, created_at)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING id;`

	QMGetAllTodosByUser = `
    SELECT
        id,
        title,
        description,
        status,
        created_at
    FROM todos
    WHERE user_id = $1
    ORDER BY created_at DESC; -- newest first`

	QMGetAllTodosByUserWithLimit = `
    -- if limit is negative (-1) it will ignore the limitation 
    -- and return all possible rows
    SELECT
        id,
        title,
        description,
        status,
        created_at
    FROM todos
    WHERE user_id = $1
    ORDER BY created_at DESC -- newest first
    LIMIT $2
    OFFSET $3;`

	QMGetAllTodosByUserWithStatusFilter = `
    -- if limit is negative (-1) it will ignore the limitation 
    -- and return all possible rows
    SELECT
        id,
        title,
        description,
        status,
        created_at
    FROM todos
    WHERE user_id = $1 AND status = $2
    ORDER BY created_at DESC -- newest first
    LIMIT $3
    OFFSET $4;`

	QEUpdateTodo = `
    UPDATE todos
    SET 
        title = $1,
        description = $2,
        status = $3
    WHERE id = $4 AND user_id = $5;`

	QEDeleteTodo = `
    DELETE FROM todos
    WHERE id = $1 AND user_id = $2;`

	QEDeleteAllTodosByUser = `
    DELETE FROM todos
    WHERE user_id = $1;`

	QOCheckUserOwnTodo = `
    SELECT 1
    FROM todo
    WHERE id = $1 AND user_id = $2
    LIMIT 1;`
)
