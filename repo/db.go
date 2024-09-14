package repo

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/assaidy/todo-api/models"
	"github.com/assaidy/todo-api/utils"
	_ "github.com/lib/pq"
)

//go:embed queries/*.sql
var sqlFiles embed.FS

type PostgresDB struct {
	DB *sql.DB
}

func NewPostgresDB(conn string) (*PostgresDB, error) {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresDB{DB: db}, nil
}

func (pg *PostgresDB) InsertUser(user *models.User) error {
	query, err := sqlFiles.ReadFile("queries/user_insert.sql")
	if err != nil {
		return err
	}

	err = pg.DB.QueryRow(string(query), user.Name, user.Email, user.Password).Scan(&user.Id)
	if err != nil {
		return err
	}

	return nil
}

func (pg *PostgresDB) GetUserById(id int64) (*models.User, error) {
	query, err := sqlFiles.ReadFile("queries/user_get_by_id.sql")
	if err != nil {
		return nil, err
	}

	user := &models.User{Id: id}

	err = pg.DB.QueryRow(string(query), id).Scan(&user.Name, &user.Email, &user.Password, &user.JoinedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NotFoundError(fmt.Sprintf("no user with id %d found", id))
		}
		return nil, err
	}

	return user, nil
}

func (pg *PostgresDB) GetUserByEmail(email string) (*models.User, error) {
	query, err := sqlFiles.ReadFile("queries/user_get_by_email.sql")
	if err != nil {
		return nil, err
	}

	user := &models.User{Email: email}

	err = pg.DB.QueryRow(string(query), email).Scan(&user.Id, &user.Name, &user.Password, &user.JoinedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NotFoundError(fmt.Sprintf("no user with email '%s' found", email))
		}
		return nil, err
	}

	return user, nil
}

func (pg *PostgresDB) UpdateUser(user *models.User) error {
	query, err := sqlFiles.ReadFile("queries/user_update_by_id.sql")
	if err != nil {
		return err
	}

	res, err := pg.DB.Exec(string(query), user.Name, user.Email, user.Password, user.Id)
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

func (pg *PostgresDB) DeleteUserById(id int64) error {
	query, err := sqlFiles.ReadFile("queries/user_delete_by_id.sql")
	if err != nil {
		return err
	}

	res, err := pg.DB.Exec(string(query), id)
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

func (pg *PostgresDB) CheckEmailExists(email string) (bool, error) {
	query, err := sqlFiles.ReadFile("queries/user_check_email_exists.sql")
	if err != nil {
		return false, err
	}

	err = pg.DB.QueryRow(string(query), email).Scan(new(int))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (pg *PostgresDB) InsertTodo(todo *models.Todo) error {
	query, err := sqlFiles.ReadFile("queries/todo_insert.sql")
	if err != nil {
		return err
	}

	// NOTE: validate title & status in the handler first
	err = pg.DB.QueryRow(string(query), todo.UserId, todo.Title, todo.Description, todo.Status, todo.CreatedAt).Scan(&todo.Id)
	if err != nil {
		return err
	}

	return nil
}

func (pg *PostgresDB) UpdateTodo(todo *models.Todo) error {
	query, err := sqlFiles.ReadFile("queries/todo_update_by_id.sql")
	if err != nil {
		return err
	}

	// NOTE: validate title & status in the handler first
	res, err := pg.DB.Exec(string(query), todo.Title, todo.Description, todo.Status, todo.Id, todo.UserId)
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

func (pg *PostgresDB) DeleteTodoByIdAndUserId(tid, uid int64) error {
	query, err := sqlFiles.ReadFile("queries/todo_delete_by_id.sql")
	if err != nil {
		return err
	}

	res, err := pg.DB.Exec(string(query), tid, uid)
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

// NOTE: result is sorted by the creation date (most recent first)
func (pg *PostgresDB) GetAllTodosByUserId(uid int64) ([]*models.Todo, error) {
	query, err := sqlFiles.ReadFile("queries/todo_get_all_by_user_id.sql")
	if err != nil {
		return nil, err
	}

	todos := []*models.Todo{}

	rows, err := pg.DB.Query(string(query), uid)
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

// TODO: get all todos in pages + filtering

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
