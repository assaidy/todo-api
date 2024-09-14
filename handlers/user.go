package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/assaidy/todo-api/models"
	"github.com/assaidy/todo-api/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	Store models.Store
}

func NewUserHandler(s models.Store) *UserHandler {
	return &UserHandler{Store: s}
}

func (h *UserHandler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) error {
	req := models.UserCreateOrUpdateRequest{}
	if err := utils.ParseJSON(r, &req); err != nil {
		return err
	}

	if err := utils.Validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return utils.InvalidRequestData(validationErrors.Error())
	}

	exists, err := h.Store.CheckEmailExists(req.Email)
	if err != nil {
		return err
	}
	if exists {
		return utils.AlreadyExistsError(fmt.Sprintf("email '%s' is already taken", req.Email))
	}

	encryptedPassword, err := utils.Encrypt(req.Password)
	if err != nil {
		return err
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: encryptedPassword,
		JoinedAt: time.Now().UTC(),
	}

	if err := h.Store.InsertUser(&user); err != nil {
		return err
	}

	token, err := utils.CreateToken(user.Id)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusCreated, map[string]any{
		"token": token,
		"user":  user,
	})
}

func (h *UserHandler) HandleLoginUser(w http.ResponseWriter, r *http.Request) error {
	req := models.UserLoginRequest{}
	if err := utils.ParseJSON(r, &req); err != nil {
		return err
	}

	if err := utils.Validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return utils.InvalidRequestData(validationErrors.Error())
	}

	user, err := h.Store.GetUserByEmail(req.Email)
	if err != nil {
		return err
	}

	dencryptedPassword, err := utils.Decrypt(user.Password)
	if err != nil {
		return err
	}

	if req.Password != dencryptedPassword {
		return utils.NotFoundError("invalid password")
	}

	token, err := utils.CreateToken(user.Id)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, map[string]any{
		"token": token,
		"user":  user,
	})
}

func (h *UserHandler) HandleUpdateUserById(w http.ResponseWriter, r *http.Request) error {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	user, err := h.Store.GetUserById(int64(id))
	if err != nil {
		return err
	}

	userId, ok := utils.GetUserIdFromContext(r.Context())
	if !ok || user.Id != userId {
		return utils.ForbiddenError()
	}

	req := models.UserCreateOrUpdateRequest{}
	if err := utils.ParseJSON(r, &req); err != nil {
		return err
	}

	if err := utils.Validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return utils.InvalidRequestData(validationErrors.Error())
	}

	if req.Email != user.Email {
		exists, err := h.Store.CheckEmailExists(req.Email)
		if err != nil {
			return err
		}
		if exists {
			return utils.AlreadyExistsError(fmt.Sprintf("email '%s' is already taken", req.Email))
		}
	}

	user.Name = req.Name
	user.Email = req.Email
	user.Password = req.Password

	if err := h.Store.UpdateUser(user); err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, &user)
}

func (h *UserHandler) HandleDeleteUserById(w http.ResponseWriter, r *http.Request) error {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	user, err := h.Store.GetUserById(int64(id))
	if err != nil {
		return err
	}

	userId, ok := utils.GetUserIdFromContext(r.Context())
	if !ok || user.Id != userId {
		return utils.ForbiddenError()
	}

	if err := h.Store.DeleteUserById(user.Id); err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusNoContent, nil)
}
