package main

import (
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

	e.GET("/users/:id", app.showUsersHandler)
	
	return e
}
