package log

import (
	"log/slog"
	"os"
)

const (
	// LOCAL defines the local environment
	LOCAL = "local"

	// DEV defines the dev environment
	DEV = "dev"

	// PROD defines the prod environment
	PROD = "prod"
)

func New(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case LOCAL:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case DEV:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case PROD:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
