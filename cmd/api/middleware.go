package main

import (
	"net/http"

	"gamemasterweb.net/internal/logger"
	"github.com/labstack/echo/v4"
)

var zeroLog = logger.NewLogger()

func RecoverPanic(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {
		defer func() {
			if err := recover(); err != nil {
				zeroLog.Info().Msgf("Panic occurred: %v", err)
				errResponse := map[string]interface{}{
					"error": "Internal Server Error",
				}
				c.JSON(http.StatusOK, errResponse)
			}

		}()

		return next(c)
	}
}
