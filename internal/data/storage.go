package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Storage struct {
	User UserModel
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		User: UserModel{DB: db},
	}
}
