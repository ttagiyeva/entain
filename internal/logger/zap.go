package logger

import (
	"context"

	"github.com/ttagiyeva/entain/internal/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewZapLogger returns a new instance of zap logger
func NewZapLogger(lc fx.Lifecycle, conf *config.Config) *zap.SugaredLogger {
	cfg := zap.Config{
		Encoding:    conf.Logger.Encoding,
		OutputPaths: []string{"stdout"},
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:    "level",
			MessageKey:  "message",
			FunctionKey: "function",
			TimeKey:     "time",
			EncodeLevel: zapcore.LowercaseLevelEncoder,
			EncodeTime:  zapcore.ISO8601TimeEncoder,
		},
		ErrorOutputPaths: []string{"stderr"},
	}

	cfg.Level.UnmarshalText([]byte(conf.Logger.Level))

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {

			return logger.Sync()
		},
	})

	return logger.Sugar()
}
