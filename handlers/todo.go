package handlers

import (
	"net/http"
	"time"

	"github.com/assaidy/todo-api/models"
	"github.com/assaidy/todo-api/utils"
	"github.com/go-playground/validator/v10"
)

type TodoHandler struct {
	Store models.Store
}

func NewTodoHandler(s models.Store) *TodoHandler {
	return &TodoHandler{
		Store: s,
	}
}

func (h *TodoHandler) HandleCreateTodo(w http.ResponseWriter, r *http.Request) error {
	userId, ok := utils.GetUserIdFromContext(r.Context())
	if !ok {
		return utils.ForbiddenError()
	}

	// check if there's a user with that id
	if exists, err := h.Store.CheckUserIdExists(userId); err != nil {
		return err
	} else if !exists {
		return utils.ForbiddenError()
	}

	req := models.TodoCreateOrUpdateRequest{}
	if err := utils.ParseJSON(r, &req); err != nil {
		return err
	}

	if err := utils.Validate.Struct(req); err != nil {
		errors := err.(validator.ValidationErrors)
		return utils.InvalidRequestData(errors.Error())
	}

	// if !slices.Contains(models.TodoStatus, req.Status) {
	//     return utils.InvalidRequestData("invalid todo status")
	// }

	todo := models.Todo{
		UserId:      userId,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		CreatedAt:   time.Now().UTC(),
	}

	if err := h.Store.InsertTodo(&todo); err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusCreated, &todo)
}

func (h *TodoHandler) HandleGetAllTodosByUser(w http.ResponseWriter, r *http.Request) error {
	userId, ok := utils.GetUserIdFromContext(r.Context())
	if !ok {
		return utils.ForbiddenError()
	}

	// check if there's a user with that id
	if exists, err := h.Store.CheckUserIdExists(userId); err != nil {
		return err
	} else if !exists {
		return utils.ForbiddenError()
	}

	todos, err := h.Store.GetAllTodosByUserId(userId)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, todos)
}

func (h *TodoHandler) HandleDeleteAllTodosByUser(w http.ResponseWriter, r *http.Request) error {
	userId, ok := utils.GetUserIdFromContext(r.Context())
	if !ok {
		return utils.ForbiddenError()
	}

	// check if there's a user with that id
	if exists, err := h.Store.CheckUserIdExists(userId); err != nil {
		return err
	} else if !exists {
		return utils.ForbiddenError()
	}

	if err := h.Store.DeleteAllTodoByUserId(userId); err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusNoContent, nil)
}
