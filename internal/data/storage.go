package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Storage struct {
	User     UserModel
	UserMock UserStorage
}

type UserStorage interface {
	Add(user *User) error
}

type MockUserStorage struct {
	Users map[string]*User
	Err   error
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		User: UserModel{DB: db},
	}
}

func (m *MockUserStorage) Add(user *User) error {
	if m.Err != nil {
		return m.Err
	}

	if _, exists := m.Users[user.Email]; exists {
		return ErrDuplicateEmail
	}

	m.Users[user.Email] = user
	return nil
}
