package application

import (
	"github.com/labstack/echo/v4"
)

type AppContext struct {
	echo.Context
	App *Application
}
