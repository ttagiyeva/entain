package logger

import (
	"log/slog"
	"os"

	"github.com/ttagiyeva/entain/internal/config"
)

var logLevels = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

// NewLogger returns a new instance of slog logger.
func NewLogger(conf *config.Config) *slog.Logger {
	level, ok := logLevels[conf.Logger.Level]
	if !ok {
		level = slog.LevelError
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: level}))
}
