package repo

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"time"

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

func (pg *PostgresDB) InsertUser(user *models.UserCreateOrUpdate) (*models.User, error) {
	query, err := sqlFiles.ReadFile("user_insert.sql")
	if err != nil {
		return nil, err
	}

	var (
		userId   int64
		joinedAt time.Time
	)

	err = pg.DB.QueryRow(string(query), user.Name, user.Email, user.Password).Scan(&userId, &joinedAt)
	if err != nil {
		return nil, err
	}

	return &models.User{
		Id:       userId,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		JoinedAt: joinedAt,
	}, nil
}

func (pg *PostgresDB) GetUserById(id int64) (*models.User, error) {
	query, err := sqlFiles.ReadFile("user_get_by_id.sql")
	if err != nil {
		return nil, err
	}

	user := &models.User{Id: id}

	// NOTE: validate email in the handler first
	// and check if email exists with pg.CheckEmailExists(email)
	err = pg.DB.QueryRow(string(query), id).Scan(&user.Name, &user.Email, &user.JoinedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NotFoundError(fmt.Sprintf("no user with id %d found", id))
		}
		return nil, err
	}

	return user, nil
}

func (pg *PostgresDB) UpdateUserById(id int64, user *models.UserCreateOrUpdate) error {
	query, err := sqlFiles.ReadFile("user_update_by_id.sql")
	if err != nil {
		return err
	}

	// NOTE: validate email in the handler first
	// and check if email exists with pg.CheckEmailExists(email)
	res, err := pg.DB.Exec(string(query), user.Name, user.Email, user.Password, id)
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

func (pg *PostgresDB) DeleteUserById(id int64) error {
	query, err := sqlFiles.ReadFile("user_delete_by_id.sql")
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
	query, err := sqlFiles.ReadFile("user_check_email_exists.sql")
	if err != nil {
		return false, err
	}

	var emailExists int

	err = pg.DB.QueryRow(string(query), email).Scan(&emailExists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
