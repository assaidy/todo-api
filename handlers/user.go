package handlers

import (
	"github.com/assaidy/todo-api/models"
)

type UserHandler struct {
	Store models.Store
}

func NewUserHandler(s models.Store) *UserHandler {
	return &UserHandler{Store: s}
}
