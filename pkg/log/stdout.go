package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func StdoutEncoder() zapcore.Encoder {
	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zapcore.NewConsoleEncoder(developmentCfg)
}
