package models

type Todo struct {
	Id          int64  `json:"id"`
	UserId      int64  `json:"userId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdAt"`
}

type TodoCreateOrUpdateRequest struct {
    Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Status      string `json:"status" validate:"required"`
}
