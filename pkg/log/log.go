package log

import (
	"github.com/fatih/color"
	"go.elastic.co/apm/module/apmzap/v2"
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net"
	"os"
	"po/pkg/logstash"
)

func Boot() {
	client, err := logstash.NewSingleton()

	var conn net.Conn

	// Since application should not go down for a log failure we won't panic
	if err != nil {
		conn1, conn2 := net.Pipe()
		defer conn1.Close()
		defer conn2.Close()

		conn = conn1
		color.Red("unable to connect to logstash: %v", err.Error())
	} else {
		conn = client.Conn
	}

	ls := zapcore.AddSync(
		conn,
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
		zap.WrapCore((&apmzap.Core{}).WrapCore),
	)

	zap.ReplaceGlobals(logger)
}
