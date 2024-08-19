package application

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"path"
	"path/filepath"
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

type DB struct{ DSN string }

type Config struct{ Port int }

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

func LoadEnv() error {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	return nil
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

func (app *Application) JsonError(c echo.Context, data interface{}) error {
	res := Response{
		Status: "error",
		Data:   data,
	}
	return c.JSON(http.StatusOK, res)
}

func (app *Application) Respond(c echo.Context, status int, jsonResponse interface{}, htmlTmpl string, tmplData data.TemplateData) error {
	if c.Request().Header.Get("Accept") == "application/json" {
		return app.JsendSuccess(c, jsonResponse)
	}
	return c.Render(status, htmlTmpl, tmplData)
}

func (app *Application) WithAppContext(handler func(AppContext) error) func(echo.Context) error {
	return func(c echo.Context) error {

		appCtx := AppContext{
			Context: c,
			App:     app,
		}

		return handler(appCtx)
	}
}

// func NewTemplateCache() (map[string]*template.Template, error) {
// 	cache := map[string]*template.Template{}

// 	pages := []string{
// 		"./static/ui/html/tableAllUsers.html",
// 		"./static/ui/html/table.html",
// 		"./static/ui/html/addUser.html",
// 		"./static/ui/html/404.html",
// 		"./static/ui/html/updateUserForms.html",
// 	}

// 	for _, page := range pages {
// 		log.Printf("Attempting to open template file: %s", page)
// 		ts, err := template.ParseFiles(page)
// 		if err != nil {
// 			log.Printf("Error loading template %s: %v", page, err)
// 			return nil, err
// 		}
// 		cache[filepath.Base(page)] = ts
// 	}
// 	return cache, nil
// }

func ReadTemplateFromRootPath(projectRootPath string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pathToTemplates := path.Join(projectRootPath, "static/ui/html/*.html")
	pages, _ := filepath.Glob(pathToTemplates)
	log.Printf("found %v templates to load", len(pages))

	for _, page := range pages {
		ts, err := template.ParseFiles(page)
		if err != nil {
			log.Printf("Error loading template %s: %v", page, err)
			return nil, err
		}
		cache[filepath.Base(page)] = ts
	}
	return cache, nil
}

func ReadTemplates() (map[string]*template.Template, error) {
	return ReadTemplateFromRootPath(".")
}
