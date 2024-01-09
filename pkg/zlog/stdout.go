package zlog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func StdoutEncoder() zapcore.Encoder {
	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zapcore.NewConsoleEncoder(developmentCfg)
}

func StdoutWriter() zapcore.WriteSyncer {
	return zapcore.AddSync(os.Stdout)
}
