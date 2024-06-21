package main

import (
	"net/http"
	"os"

	"gamemasterweb.net/cmd/api/api_handlers"
	"gamemasterweb.net/internal/app"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type pathsSwagger struct {
	filePathSwagger   string
	pathStaticSwagger string
}

func routes(app *app.Application) *echo.Echo {

	e := echo.New()

	e.Use(RecoverPanic)

	pathSwagger := pathsSwagger{
		filePathSwagger:   os.Getenv("SWAGGER_FILE"),
		pathStaticSwagger: os.Getenv("STATIC_SWAGGER"),
	}

	secretKeySession := os.Getenv("SECRET_KEY_FOR_SESSION")
	if secretKeySession == "" {
		e.Logger.Fatal("SECRET_KEY_FOR_SESSION environment variable is required")
	}

	e.Static("/swagger/", pathSwagger.pathStaticSwagger)
	e.File("/docs/api/swagger.json", pathSwagger.filePathSwagger)

	e.Use(session.Middleware(sessions.NewCookieStore([]byte(secretKeySession))))

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusSeeOther, "/users")
	})

	e.GET("/users", func(c echo.Context) error {
		return api_handlers.ShowAllUsersHandler(c, app)
	})

	e.GET("/users/:id", func(c echo.Context) error {
		return api_handlers.ShowOneUserHandler(c, app)
	})

	e.GET("/users/new", func(c echo.Context) error {
		return api_handlers.AddUsersHandler(c, app)
	})

	e.GET("/users/edit/:id", func(c echo.Context) error {
		return api_handlers.EditUserFormHandler(c, app)
	})

	e.POST("/users", func(c echo.Context) error {
		return api_handlers.AddUsersHandler(c, app)
	})

	e.POST("/users/edit/:id", func(c echo.Context) error {
		return api_handlers.UpdateUsersHandler(c, app)
	})

	e.DELETE("/users/:id", func(c echo.Context) error {
		return api_handlers.DeleteUsersHandler(c, app)
	})

	checkRoutesPath(e, app)

	return e

}

func checkRoutesPath(e *echo.Echo, app *app.Application) {

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
