package data

import (
	"fmt"
	"strconv"
)

type MockUserStorage struct {
	Users map[string]*User
	Err   error
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

func (m *MockUserStorage) Get(id int64) (*User, error) {
	idStr := strconv.FormatInt(id, 10)
	user, exists := m.Users[idStr]
	if !exists {
		return nil, fmt.Errorf("user with id %d not found", id)
	}
	return user, nil
}

func (m *MockUserStorage) GetAll() ([]User, error) {
	var users []User

	for _, user := range m.Users {
		users = append(users, *user)
	}

	return users, nil
}

func (m *MockUserStorage) Update(user *User) error {
	var foundUser *User
	var oldNickname string

	for nickname, storedUser := range m.Users {
		if storedUser.ID == user.ID {
			foundUser = storedUser
			oldNickname = nickname
			break
		}
	}

	if foundUser == nil {
		return fmt.Errorf("user not found")
	}

	if oldNickname != user.Nickname {
		delete(m.Users, oldNickname)
	}

	m.Users[user.Nickname] = user

	return nil
}

func (m *MockUserStorage) Delete(id int64) error {
	for nickname, user := range m.Users {
		if user.ID == id {
			delete(m.Users, nickname)
			return nil
		}
	}
	return fmt.Errorf("user not found")
}
