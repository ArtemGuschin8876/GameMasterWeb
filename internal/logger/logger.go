package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func NewLogger() zerolog.Logger {
	consoleWriter := zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.Out = os.Stdout
		w.TimeFormat = "3:04PM"
	})

	return zerolog.New(consoleWriter).With().Timestamp().Logger()
}
