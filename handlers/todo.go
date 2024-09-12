package handlers

import "database/sql"

type TodoHandler struct {
	DB *sql.DB
}

func NewTodoHandler(db *sql.DB) *TodoHandler {
	return &TodoHandler{
		DB: db,
	}
}
