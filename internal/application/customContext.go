package application

import (
	"github.com/labstack/echo/v4"
)

type CustomContext struct {
	echo.Context
	App *Application
}
