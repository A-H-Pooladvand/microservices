package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var Sugar *zap.SugaredLogger

func Boot() {
	stdout := zapcore.AddSync(os.Stdout)
	file := zapcore.AddSync(FileRoller())

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, zap.NewAtomicLevelAt(zap.DebugLevel)),
		zapcore.NewCore(fileEncoder, file, zap.NewAtomicLevelAt(zap.InfoLevel)),
	)

	logger := zap.New(
		core,
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.AddCaller(),
	)

	Sugar = logger.Sugar()

	zap.ReplaceGlobals(logger)
}

func Error(args ...any) {
	Sugar.Error(args...)
}

func Info(args ...any) {
	Sugar.Info(args...)
}
