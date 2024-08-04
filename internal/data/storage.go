package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Storage struct {
	User UserStorage
}

type UserStorage interface {
	Add(user *User) error
	Get(id int64) (*User, error)
	GetAll() ([]User, error)
	Update(user *User) error
	Delete(id int64) error
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		User: UserModel{DB: db},
	}
}
