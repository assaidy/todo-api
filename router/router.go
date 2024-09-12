package router

import (
	"net/http"

	"github.com/assaidy/todo-api/handlers"
	"github.com/assaidy/todo-api/models"
	"github.com/assaidy/todo-api/utils"
	"github.com/gorilla/mux"
)

func NewRouter(s models.Store) http.Handler {
	router := mux.NewRouter().StrictSlash(true)
	protected := router.PathPrefix("").Subrouter()
	protected.Use(utils.WithJWT)

	userH := handlers.NewUserHandler(s)
	todoH := handlers.NewTodoHandler(s)

	protected.HandleFunc("/register", utils.Make(userH.HandleRegisterUser)).Methods("POST")
	protected.HandleFunc("/login", utils.Make(userH.HandleLoginUser)).Methods("POST")
	protected.HandleFunc("/users/{id:[0-9+]}", utils.Make(userH.HandleGetUserById)).Methods("GET")

	protected.HandleFunc("/users/{id:[0-9+]}", utils.Make(userH.HandleDeleteUserById)).Methods("DELETE")
	protected.HandleFunc("/todos", utils.Make(todoH.HandleGetAllTodos)).Methods("GET")
	protected.HandleFunc("/todos", utils.Make(todoH.HandleDeleteAllTodos)).Methods("DELETE")
	protected.HandleFunc("/todos", utils.Make(todoH.HandleCreateTodo)).Methods("POST")
	// protected.HandleFunc("/todos/{id:[0-9+]}", utils.Make(todoH.HandleGetTodoById)).Methods("GET")
	protected.HandleFunc("/todos/{id:[0-9+]}", utils.Make(todoH.HandleUpdateTodoById)).Methods("PUT")
	protected.HandleFunc("/todos/{id:[0-9+]}", utils.Make(todoH.HandleDeleteTodoById)).Methods("DELETE")

	return router
}
