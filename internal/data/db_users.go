package data

import (
	"database/sql"
	"log"
)

type UserModel struct {
	DB *sql.DB
}

func (m UserModel) Add(user *Users) error {

	query := `
		INSERT INTO users (name, nickname, email, city, about, image)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	args := []interface{}{
		user.Name,
		user.Nickname,
		user.Email,
		user.City,
		user.About,
		user.Image,
	}

	err := m.DB.QueryRow(query, args...).Scan(&user.ID)
	if err != nil {
		log.Println("Error executing query:", err)
	}

	return err
}

func (m UserModel) Get() error {
	return nil
}

func (m UserModel) Update() error {
	return nil
}

func (m UserModel) Delete() error {
	return nil
}
