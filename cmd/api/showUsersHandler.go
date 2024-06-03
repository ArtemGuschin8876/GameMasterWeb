package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (app *application) showUsersHandler(c echo.Context) error {
	return c.String(http.StatusOK, "HIIIII!")
}
