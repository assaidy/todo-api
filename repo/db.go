package repo

import (
	"database/sql"
	"embed"
	_ "github.com/lib/pq"
)

//go:embed queries/*.sql
var sqlFiles embed.FS

func NewPostgresDB(conn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// handle different database tasks
