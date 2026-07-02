package logger

import (
	"log/slog"
	"os"
)

func New() *slog.Logger {

	handler := slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		},
	)

	return slog.New(handler)
}
