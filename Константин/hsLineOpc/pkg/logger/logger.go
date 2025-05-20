package logger

import (
	"log/slog"
	"os"
)

func SetupLogging(level slog.Level) {
	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}),
	)

	slog.SetDefault(log)
}
