package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"gamemasterweb.net/internal/data"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func (app *application) readIDParam(c echo.Context) (int64, error) {

	IDParam := c.Param("id")

	id, err := strconv.ParseInt(IDParam, 10, 64)
	if err != nil || id < 1 {
		log.Println(id)
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

func (app *application) writeJSON(c echo.Context) error {

	id, err := app.readIDParam(c)
	log.Printf(
		"err is %v", err,
	)

	if err != nil {
		return c.JSON(200, response{
			Status: "error",
			//data:    nil,
			Message: "problema in Oleg",
		})
		//return c.String(http.StatusNotFound, "the requested resource could not be found")
	}

	user := data.Users{
		ID:       id,
		Name:     "Oleg",
		Nickname: "Parlis",
		Email:    "OlegSuka@gmail.com",
		City:     "Saratov",
		About:    "I am pidor",
	}

	return c.JSON(http.StatusOK, user)
}
