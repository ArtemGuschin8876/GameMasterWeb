package data

import (
	"database/sql"
	"errors"
	"log"
)

type UserModel struct {
	DB *sql.DB
}

var (
	ErrDuplicateNickname = errors.New("nickname already exists")
	ErrDuplicateEmail    = errors.New("email already exists")
)

func (m UserModel) Add(user *Users) error {

	isUnique, err := m.isEmailUnique(user.Email)
	if err != nil {
		return err
	}

	if !isUnique {
		return ErrDuplicateEmail
	}

	query := `
		INSERT INTO users (name, nickname, email, city, about, image)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
		`

	args := []interface{}{
		user.Name,
		user.Nickname,
		user.Email,
		user.City,
		user.About,
		user.Image,
	}

	err = m.DB.QueryRow(query, args...).Scan(&user.ID)
	if err != nil {
		log.Println("Error executing query:", err)
	}

	return err
}

func (m UserModel) Get(id int64) (*Users, error) {

	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
	SELECT id, name, nickname, email, city, about, image
	FROM users
	WHERE id = $1`

	var user Users

	err := m.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Nickname,
		&user.Email,
		&user.City,
		&user.About,
		&user.Image,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (m UserModel) GetAll() ([]Users, error) {

	query := `
	SELECT id, name, nickname, email, city, about, image
	FROM users
	`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}

	var users []Users

	for rows.Next() {
		var user Users

		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Nickname,
			&user.Email,
			&user.City,
			&user.About,
			&user.Image,
		); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (m UserModel) Update(user *Users) error {

	query := `
	UPDATE users
	SET name = $1, nickname = $2, email = $3, city = $4, about = $5, image = $6
	WHERE id = $7`

	args := []interface{}{
		user.Name,
		user.Nickname,
		user.Email,
		user.City,
		user.About,
		user.Image,
		user.ID,
	}

	_, err := m.DB.Exec(query, args...)
	return err
}

func (m UserModel) Delete(id int64) error {

	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
	DELETE FROM users
	WHERE id = $1`

	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m UserModel) isEmailUnique(email string) (bool, error) {

	var id int

	query := `
	SELECT id FROM users WHERE email = $1`

	err := m.DB.QueryRow(query, email).Scan(&id)
	if err == sql.ErrNoRows {
		return true, nil
	}

	if err != nil {
		return false, err
	}

	return false, nil
}
