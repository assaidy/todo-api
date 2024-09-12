package handlers

import (
	"github.com/assaidy/todo-api/models"
)

type TodoHandler struct {
	Store models.Store
}

func NewTodoHandler(s models.Store) *TodoHandler {
	return &TodoHandler{
		Store: s,
	}
}
