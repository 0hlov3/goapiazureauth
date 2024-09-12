package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitializeZapCustomLogger(logLevel string) *zap.Logger {
	newLogLevel, err := zap.ParseAtomicLevel(logLevel)
	if err != nil {
		fmt.Println("Invalid log level, defaulting to Info level")
		newLogLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	config := zap.Config{
		Encoding:    "json",
		Level:       newLogLevel,
		OutputPaths: []string{"stdout"},
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:     "level",
			TimeKey:      "time",
			CallerKey:    "file",
			MessageKey:   "msg",
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	logger, err := config.Build()
	if err != nil {
		fmt.Println("Error building logger")
		return nil
	}

	return logger
}
