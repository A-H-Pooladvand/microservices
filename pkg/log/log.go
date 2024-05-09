package log

import (
	"context"
	"errors"
	"go.elastic.co/apm/module/apmzap/v2"
	"go.elastic.co/ecszap"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"po/pkg/logstash"
	"syscall"
)

func Invoke(lc fx.Lifecycle, l *logstash.Client) {
	ls := zapcore.AddSync(
		l.Connection(),
	)

	core := zapcore.NewTee(
		// Stdout Writer
		zapcore.NewCore(StdoutEncoder(), zapcore.AddSync(os.Stdout), zap.NewAtomicLevelAt(zap.DebugLevel)),
		// File Writer
		ecszap.NewCore(ecszap.NewDefaultEncoderConfig(), FileWriter(), zap.NewAtomicLevelAt(zap.InfoLevel)),
		// Logstash Writer
		ecszap.NewCore(ecszap.NewDefaultEncoderConfig(), ls, zap.NewAtomicLevelAt(zap.InfoLevel)),
	)

	logger := zap.New(
		core,
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.WrapCore((&apmzap.Core{}).WrapCore),
	)

	zap.ReplaceGlobals(logger)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			err := zap.L().Sync()

			if !errors.Is(err, syscall.EINVAL) {
				return err
			}

			return nil
		},
	})
}

func Error(msg string, fields ...zap.Field) {
	zap.L().Error(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	zap.L().Info(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	zap.L().Panic(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	zap.L().Warn(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	zap.L().Fatal(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	zap.L().Debug(msg, fields...)
}
