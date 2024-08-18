package api_handlers

import (
	"log"

	"gamemasterweb.net/internal/application"
)

func NewUserForm(c application.AppContext) error {

	app := c.App

	err := app.RenderHTML(c, "addUser", nil)
	if err != nil {
		log.Println("error render html in newUserForm")
	}

	return nil
}
