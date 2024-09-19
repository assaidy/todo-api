package router

import (
	"net/http"

	"github.com/assaidy/todo-api/handlers"
	"github.com/assaidy/todo-api/repo"
	"github.com/assaidy/todo-api/utils"
	"github.com/gorilla/mux"
)

func NewRouter(r *repo.Repo) http.Handler {
	router := mux.NewRouter().StrictSlash(true)
	protected := router.PathPrefix("").Subrouter()
	protected.Use(utils.WithJWT)

	userH := handlers.NewUserHandler(r)
	todoH := handlers.NewTodoHandler(r)

	router.HandleFunc("/register", utils.Make(userH.HandleRegisterUser)).Methods("POST")
	router.HandleFunc("/login",    utils.Make(userH.HandleLoginUser)).Methods("POST")

	protected.HandleFunc("/users/{id:[0-9+]}", utils.Make(userH.HandleDeleteUserById)).Methods("DELETE")
	protected.HandleFunc("/users/{id:[0-9+]}", utils.Make(userH.HandleUpdateUserById)).Methods("PUT")
	protected.HandleFunc("/todos",             utils.Make(todoH.HandleCreateTodo)).Methods("POST")
	protected.HandleFunc("/todos",             utils.Make(todoH.HandleGetAllTodosByUser)).Methods("GET")
	protected.HandleFunc("/todos",             utils.Make(todoH.HandleDeleteAllTodosByUser)).Methods("DELETE")
	protected.HandleFunc("/todos/{id:[0-9+]}", utils.Make(todoH.HandleDeleteTodoById)).Methods("DELETE")
	protected.HandleFunc("/todos/{id:[0-9+]}", utils.Make(todoH.HandleUpdateTodoById)).Methods("PUT")
	// protected.HandleFunc("/todos/{id:[0-9+]}", utils.Make(todoH.HandleGetTodoById)).Methods("GET")

	return router
}

// TODO: apply filtering to 'get_all_todos'
// TODO: make a proper logger middleware and remove utils.Make(),or just use 'fiber'
