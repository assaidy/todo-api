package models

type Store interface {
	InsertUser(*User) error
	GetUserById(int64) (*User, error)
	GetUserByEmail(string) (*User, error)
	UpdateUser(*User) error
	DeleteUserById(int64) error
	CheckEmailExists(string) (bool, error)

	InsertTodo(*Todo) error
	UpdateTodo(*Todo) error
	DeleteTodoByIdAndUserId(int64, int64) error
	GetAllTodosByUserId(int64) ([]*Todo, error)
	// CheckUserOwnsTodo(int64, int64) (bool, error)
}
