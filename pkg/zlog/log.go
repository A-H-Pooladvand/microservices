package zlog

import (
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"po/pkg/logstash"
)

func Boot() {
	ls := zapcore.AddSync(
		logstash.Get().Conn,
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
	)

	zap.ReplaceGlobals(logger)
}
