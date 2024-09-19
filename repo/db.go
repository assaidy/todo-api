package repo

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/assaidy/todo-api/models"
	"github.com/assaidy/todo-api/utils"
	_ "github.com/lib/pq"
)

type Repo struct {
	DB *sql.DB
}

func New(conn string) (*Repo, error) {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Repo{DB: db}, nil
}

func (r *Repo) InsertUser(user *models.User) error {
    err := r.DB.QueryRow(QOInsertUser, user.Name, user.Email, user.Password).Scan(&user.Id)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) GetUserById(id int64) (*models.User, error) {
	user := &models.User{Id: id}

    err := r.DB.QueryRow(QMGetUserById, id).Scan(&user.Name, &user.Email, &user.Password, &user.JoinedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NotFoundError(fmt.Sprintf("no user with id %d found", id))
		}
		return nil, err
	}

	return user, nil
}

func (r *Repo) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{Email: email}

    err := r.DB.QueryRow(QMGetUserByEmail, email).Scan(&user.Id, &user.Name, &user.Password, &user.JoinedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NotFoundError(fmt.Sprintf("no user with email '%s' found", email))
		}
		return nil, err
	}

	return user, nil
}

func (r *Repo) UpdateUser(user *models.User) error {
	res, err := r.DB.Exec(QEUpdateUser, user.Name, user.Email, user.Password, user.Id)
	if err != nil {
		return err
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows == 0 {
		return utils.NotFoundError(fmt.Sprintf("no user with id %d found", user.Id))
	}

	return nil
}

func (r *Repo) DeleteUserById(id int64) error {
	res, err := r.DB.Exec(QEDeleteUser, id)
	if err != nil {
		return err
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows == 0 {
		return utils.NotFoundError(fmt.Sprintf("no user with id %d found", id))
	}

	return nil
}

func (r *Repo) CheckEmailExists(email string) (bool, error) {
    err := r.DB.QueryRow(QOCheckEmailExists, email).Scan(new(int))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *Repo) CheckUserIdExists(id int64) (bool, error) {
    err := r.DB.QueryRow(QOCheckUserIdExists, id).Scan(new(int))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *Repo) InsertTodo(todo *models.Todo) error {
    err := r.DB.QueryRow(QOInsertTodo, todo.UserId, todo.Title, todo.Description, todo.Status, todo.CreatedAt).Scan(&todo.Id)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) UpdateTodo(todo *models.Todo) error {
	res, err := r.DB.Exec(QEUpdateTodo, todo.Title, todo.Description, todo.Status, todo.Id, todo.UserId)
	if err != nil {
		return err
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows == 0 {
		return utils.NotFoundError(fmt.Sprintf("no todo with id %d found", todo.Id))
	}

	return nil
}

func (r *Repo) DeleteTodoByIdAndUserId(tid, uid int64) error {
	res, err := r.DB.Exec(QEDeleteTodo, tid, uid)
	if err != nil {
		return err
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows == 0 {
		return utils.NotFoundError(fmt.Sprintf("no todo with id %d found for user with id %d", tid, uid))
	}

	return nil
}

func (r *Repo) DeleteAllTodoByUserId(uid int64) error {
    _, err := r.DB.Exec(QEDeleteAllTodosByUser, uid)
	if err != nil {
		return err
	}

	// affectedRows, err := res.RowsAffected()
	// if err != nil {
	// 	return err
	// }
	// if affectedRows == 0 {
	// 	return utils.NotFoundError(fmt.Sprintf("no todos found for user with id %d", uid))
	// }

	return nil
}

// NOTE: result is sorted by the creation date (most recent first)
func (r *Repo) GetAllTodosByUserId(uid int64) ([]*models.Todo, error) {
	todos := []*models.Todo{}

	rows, err := r.DB.Query(QMGetAllTodosByUser, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		t := models.Todo{UserId: uid}
		if err := rows.Scan(&t.Id, &t.Title, &t.Description, &t.Status, &t.CreatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

// func (pg *PostgresDB) CheckUserOwnsTodo(tid, uid int64) (bool, error) {
// 	query, err := sqlFiles.ReadFile("queries/todo_check_owner.sql")
// 	if err != nil {
// 		return false, err
// 	}

// 	err = pg.DB.QueryRow(string(query), tid, uid).Scan(new(int))
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return false, nil
// 		}
// 		return false, err
// 	}

// 	return true, nil
// }
