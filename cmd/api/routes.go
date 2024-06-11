package main

import (
	"log"
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

	e.GET("/users/:id", app.showOneUserHandler)
	e.POST("/users", app.addUsersHandler)
	e.PUT("/users/:id", app.updateUsersHandler)
	e.DELETE("/users/:id", app.deleteUsersHandler)

	err := checkRoutesPath(e)
	if err != nil {
		log.Fatal("error path verification")
	}

	return e
}

func checkRoutesPath(e *echo.Echo) error {

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
			msg = "the requested resource could not be found"
			c.JSON(http.StatusNotFound, response{
				Status:  "fail",
				Message: msg,
			})

		} else {
			c.JSON(code, response{
				Status:  "error",
				Message: msg,
			})
		}
	}

	return nil
}
