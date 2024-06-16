package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func jsendSuccess(c echo.Context, data interface{}) error {
	res := response{
		Status: "success",
		Data:   data,
	}
	return c.JSON(http.StatusOK, res)
}

func jsendError(c echo.Context, message string) error {
	res := response{
		Status:  "error",
		Message: message,
	}
	return c.JSON(http.StatusOK, res)
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

func (app *application) renderHTML(c echo.Context, fileName string, s any) error {
	ts, ok := app.templates[fileName+".html"]
	if !ok {
		return c.String(http.StatusBadRequest, "template doesn't exist in cache")
	}

	err := ts.Execute(c.Response().Writer, s)
	if err != nil {
		return jsendError(c, "error execute template files")
	}

	return nil
}
