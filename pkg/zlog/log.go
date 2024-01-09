package zlog

import (
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Boot() {
	core := zapcore.NewTee(
		zapcore.NewCore(StdoutEncoder(), StdoutWriter(), zap.NewAtomicLevelAt(zap.DebugLevel)),
		ecszap.NewCore(FileEncoder(), FileWriter(), zap.NewAtomicLevelAt(zap.InfoLevel)),
		ecszap.NewCore(LogstashEncoder(), LogstashWriter(), zap.NewAtomicLevelAt(zap.InfoLevel)),
	)

	logger := zap.New(
		core,
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.AddCaller(),
	)

	zap.ReplaceGlobals(logger)
}
