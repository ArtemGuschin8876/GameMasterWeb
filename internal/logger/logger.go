package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func NewLogger() zerolog.Logger {
	consoleWriter := zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.Out = os.Stdout
		w.TimeFormat = "2006-01-02 03:04:05.000PM"
	})

	return zerolog.New(consoleWriter).With().Timestamp().Caller().Logger()
}
