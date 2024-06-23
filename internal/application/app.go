package application

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"gamemasterweb.net/internal/data"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

type Response struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

type Config struct {
	Port int
	DB   struct {
		DSN string
	}
}

type Application struct {
	Logger    *log.Logger
	Config    Config
	Storage   data.Storage
	Templates map[string]*template.Template
	Response  Response
}

func (app *Application) ReadIDParam(c echo.Context) (int64, error) {

	IDParam := c.Param("id")
	if IDParam == "" {
		return 0, errors.New("missing id parameter")
	}
	id, err := strconv.ParseInt(IDParam, 10, 64)
	if err != nil || id < 1 {
		log.Println(id)
		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}

func (app *Application) RenderHTML(c echo.Context, fileName string, s any) error {
	ts, ok := app.Templates[fileName+".html"]
	if !ok {
		app.Logger.Println("template doesn't exist in cache:", fileName)
		return c.String(http.StatusBadRequest, "template doesn't exist in cache")
	}

	err := ts.Execute(c.Response().Writer, s)
	if err != nil {
		app.Logger.Printf("error executing template %s: %v", fileName, err)
		return app.JsendError(c, "error execute template files")
	}

	return nil
}

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func (app *Application) JsendSuccess(c echo.Context, data interface{}) error {
	res := Response{
		Status: "success",
		Data:   data,
	}
	return c.JSON(http.StatusOK, res)
}

func (app *Application) JsendError(c echo.Context, message string) error {
	res := Response{
		Status:  "error",
		Message: message,
	}
	return c.JSON(http.StatusOK, res)
}
