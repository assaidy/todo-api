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

	router.HandleFunc("/register", utils.Make(userH.HandleRegisterUser)).Methods("POST")
	router.HandleFunc("/login", utils.Make(userH.HandleLoginUser)).Methods("POST")

	protected.HandleFunc("/users/{id:[0-9+]}", utils.Make(userH.HandleDeleteUserById)).Methods("DELETE")
	protected.HandleFunc("/users/{id:[0-9+]}", utils.Make(userH.HandleUpdateUserById)).Methods("PUT")
    protected.HandleFunc("/todos", utils.Make(todoH.HandleCreateTodo)).Methods("POST")
	protected.HandleFunc("/todos", utils.Make(todoH.HandleGetAllTodosByUser)).Methods("GET")
	protected.HandleFunc("/todos", utils.Make(todoH.HandleDeleteAllTodosByUser)).Methods("DELETE")
	// // protected.HandleFunc("/todos/{id:[0-9+]}", utils.Make(todoH.HandleGetTodoById)).Methods("GET")
	// protected.HandleFunc("/todos/{id:[0-9+]}", utils.Make(todoH.HandleUpdateTodoById)).Methods("PUT")
	// protected.HandleFunc("/todos/{id:[0-9+]}", utils.Make(todoH.HandleDeleteTodoById)).Methods("DELETE")

	return router
}
// TODO: apply pagging to 'get_all_todos'
// TODO: apply filtering to 'get_all_todos'
