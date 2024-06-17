package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

type pathsSwagger struct {
	filePathSwagger   string
	pathStaticSwagger string
}

type response struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func (app *application) routes() *echo.Echo {
	e := echo.New()

	LoadEnv()
	pathSwagger := pathsSwagger{
		filePathSwagger:   os.Getenv("SWAGGER_FILE"),
		pathStaticSwagger: os.Getenv("STATIC_SWAGGER"),
	}

	e.Static("/swagger/", pathSwagger.pathStaticSwagger)
	e.File("/docs/api/swagger.json", pathSwagger.filePathSwagger)

	e.GET("/users", app.showAllUsersHandler)
	e.GET("/users/:id", app.showOneUserHandler)
	e.GET("/users/successfully", app.successfullCreatedUserHandler)
	e.GET("/users/new", app.showUserForm)

	e.POST("/users", app.addUsersHandler)

	e.PUT("/users/:id", app.updateUsersHandler)
	e.DELETE("/users/:id", app.deleteUsersHandler)

	app.checkRoutesPath(e)

	return e

}

func (app *application) checkRoutesPath(e *echo.Echo) {

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		msg := "the server was unable to process your request"

		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			if he.Message != nil {
				msg = he.Message.(string)
			}
		}

		if code == http.StatusNotFound {
			err := app.renderHTML(c, "404", nil)
			if err != nil {
				msg = "err rendering 404 page"
				c.String(http.StatusNotFound, msg)
			}

		} else {
			c.JSON(code, response{
				Status:  "error",
				Message: msg,
			})
		}
	}

}
