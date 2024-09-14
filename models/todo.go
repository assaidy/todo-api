package models

import "time"

type Todo struct {
	Id          int64     `json:"id"`
	UserId      int64     `json:"userId"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
}

type TodoCreateOrUpdateRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Status      string `json:"status" validate:"required,oneof=todo doing done"`
}

// var TodoStatus = []string{"todo", "doing", "done"}
