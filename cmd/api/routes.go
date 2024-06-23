package main

import (
	"net/http"
	"os"

	"gamemasterweb.net/cmd/api/api_handlers"
	"gamemasterweb.net/internal/application"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type pathsSwagger struct {
	filePathSwagger   string
	pathStaticSwagger string
}

func routes(app *application.Application) *echo.Echo {

	e := echo.New()

	e.Use(RecoverPanic)

	pathSwagger := pathsSwagger{
		filePathSwagger:   os.Getenv("SWAGGER_FILE"),
		pathStaticSwagger: os.Getenv("STATIC_SWAGGER"),
	}

	secretKeySession := os.Getenv("SECRET_KEY")
	if secretKeySession == "" {
		e.Logger.Fatal("SECRET_KEY environment variable is required")
	}

	e.Static("/swagger/", pathSwagger.pathStaticSwagger)
	e.File("/docs/api/swagger.json", pathSwagger.filePathSwagger)

	e.Use(session.Middleware(sessions.NewCookieStore([]byte(secretKeySession))))

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &application.AppContext{Context: c, App: app}
			return next(cc)
		}
	})

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusSeeOther, "/users")
	})

	e.GET("/users", api_handlers.ShowListUsers)
	e.GET("/users/:id", api_handlers.ShowUser)
	e.GET("/users/new", api_handlers.NewUserForm)
	e.GET("/users/edit/:id", api_handlers.EditUserForm)

	e.POST("/users", api_handlers.CreateUser)
	e.POST("/users/:id", api_handlers.UpdateUser)

	e.DELETE("/users/:id", api_handlers.DeleteUser)

	checkRoutesPath(e, app)

	return e

}

func checkRoutesPath(e *echo.Echo, app *application.Application) {

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
			err := app.RenderHTML(c, "404", nil)
			if err != nil {
				msg = "err rendering 404 page"
				c.String(http.StatusNotFound, msg)
			}

		} else {
			res := app.Response
			res.Status = "error"
			res.Message = msg
			c.JSON(code, res)
		}
	}
}
