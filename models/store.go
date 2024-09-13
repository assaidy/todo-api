package models

type Store interface {
	InsertUser(*User) error
	GetUserById(int64) (*User, error)
	UpdateUser(*User) error
	DeleteUserById(int64) error
	CheckEmailExists(string) (bool, error)

	InsertTodo(*Todo) error
	UpdateTodo(*Todo) error
	DeleteTodoById(int64) error
	GetAllTodosByUserId(int64) ([]*Todo, error)
}
