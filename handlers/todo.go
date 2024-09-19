package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/assaidy/todo-api/models"
	"github.com/assaidy/todo-api/repo"
	"github.com/assaidy/todo-api/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type TodoHandler struct {
	repo *repo.Repo
}

func NewTodoHandler(r *repo.Repo) *TodoHandler {
	return &TodoHandler{
		repo: r,
	}
}

func (h *TodoHandler) HandleCreateTodo(w http.ResponseWriter, r *http.Request) error {
	userId, ok := utils.GetUserIdFromContext(r.Context())
	if !ok {
		return utils.ForbiddenError()
	}

	// check if there's a user with that id
	if exists, err := h.repo.CheckUserIdExists(userId); err != nil {
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

	if err := h.repo.InsertTodo(&todo); err != nil {
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
	if exists, err := h.repo.CheckUserIdExists(userId); err != nil {
		return err
	} else if !exists {
		return utils.ForbiddenError()
	}

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	status := r.URL.Query().Get("status")

	// Default to page 1 and limit 10 if not specified
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	var todos []*models.Todo
	switch status {
	case "todo", "doing", "done":
		todos, err = h.repo.GetAllTodosByUserIdWithStatusFilter(userId, limit, offset, status)
		if err != nil {
			return err
		}
	default:
		todos, err = h.repo.GetAllTodosByUserId(userId, limit, offset)
		if err != nil {
			return err
		}
	}

	return utils.WriteJSON(w, http.StatusOK, map[string]any{
		"data":  todos,
		"page":  page,
		"limit": limit,
		"total": len(todos),
	})
}

func (h *TodoHandler) HandleDeleteAllTodosByUser(w http.ResponseWriter, r *http.Request) error {
	userId, ok := utils.GetUserIdFromContext(r.Context())
	if !ok {
		return utils.ForbiddenError()
	}

	// check if there's a user with that id
	if exists, err := h.repo.CheckUserIdExists(userId); err != nil {
		return err
	} else if !exists {
		return utils.ForbiddenError()
	}

	if err := h.repo.DeleteAllTodoByUserId(userId); err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *TodoHandler) HandleDeleteTodoById(w http.ResponseWriter, r *http.Request) error {
	userId, ok := utils.GetUserIdFromContext(r.Context())
	if !ok {
		return utils.ForbiddenError()
	}

	// check if there's a user with that id
	if exists, err := h.repo.CheckUserIdExists(userId); err != nil {
		return err
	} else if !exists {
		return utils.ForbiddenError()
	}

	todoId, _ := strconv.Atoi(mux.Vars(r)["id"])

	if err := h.repo.DeleteTodoByIdAndUserId(todoId, userId); err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *TodoHandler) HandleUpdateTodoById(w http.ResponseWriter, r *http.Request) error {
	userId, ok := utils.GetUserIdFromContext(r.Context())
	if !ok {
		return utils.ForbiddenError()
	}

	// check if there's a user with that id
	if exists, err := h.repo.CheckUserIdExists(userId); err != nil {
		return err
	} else if !exists {
		return utils.ForbiddenError()
	}

	todoId, _ := strconv.Atoi(mux.Vars(r)["id"])

	req := models.TodoCreateOrUpdateRequest{}
	if err := utils.ParseJSON(r, &req); err != nil {
		return err
	}

	todo := models.Todo{
		Id:          todoId,
		UserId:      userId,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
	}

	if err := h.repo.UpdateTodo(&todo); err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, &todo)
}
