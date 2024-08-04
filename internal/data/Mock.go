package data

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
	var User *User
	return User, nil
}

func (m *MockUserStorage) GetAll() ([]User, error) {
	return nil, nil
}

func (m *MockUserStorage) Update(user *User) error {
	return nil
}

func (m *MockUserStorage) Delete(id int64) error {
	return nil
}
