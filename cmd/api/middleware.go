package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func RecoverPanic(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic occurred: %v", err)
				errResponse := map[string]interface{}{
					"error": "Internal Server Error",
				}
				c.JSON(http.StatusOK, errResponse)
			}

		}()

		return next(c)
	}
}
