package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
	"task-api/pkg/config"
)

func CreateLogger(cfg config.Logger) (*zap.Logger, error) {
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	var encoder zapcore.Encoder
	if strings.ToLower(cfg.AppEnv) == "production" {
		encoderCfg.EncodeLevel = zapcore.LowercaseLevelEncoder
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	}
	var cores []zapcore.Core
	if cfg.Output == "stdout" || cfg.Output == "both" {
		consoleCore := zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), level)
		cores = append(cores, consoleCore)
	}
	if cfg.Output == "file" || cfg.Output == "both" {
		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.FilePath,
			MaxSize:    cfg.MaxSizeMB,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAgeDays,
			Compress:   cfg.Compress,
		})
		fileEncoder := zapcore.NewJSONEncoder(encoderCfg)
		fileCore := zapcore.NewCore(fileEncoder, fileWriter, level)
		cores = append(cores, fileCore)
	}
	combinedCore := zapcore.NewTee(cores...)
	logger := zap.New(combinedCore, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))
	return logger, nil
}
