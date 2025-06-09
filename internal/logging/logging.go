package logging

import (
	"log/slog"
	"os"
)

func Setup() *slog.Logger {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	return logger
}
