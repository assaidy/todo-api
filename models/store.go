package models

type Store interface {
	InsertUser(*UserCreateOrUpdate) (*User, error)
	GetUserById(int64) (*User, error)
	UpdateUserById(int64, *UserCreateOrUpdate) error
	DeleteUserById(int64) error
	CheckEmailExists(string) (bool, error)
}
