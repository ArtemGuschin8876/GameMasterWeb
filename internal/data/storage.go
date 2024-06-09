package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Storage struct {
	Users UserModel
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Users: UserModel{DB: db},
	}
}
