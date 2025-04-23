package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger() *zap.Logger {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"stdout", "./logs/app.log"}
	cfg.ErrorOutputPaths = []string{"stdout", "./logs/error.log"}
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := cfg.Build()

	if err != nil {
		panic(fmt.Sprintf("failed to initialize logger: %v", err))
	}

	zap.ReplaceGlobals(logger)

	return logger
}
