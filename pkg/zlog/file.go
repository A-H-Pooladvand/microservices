package zlog

import (
	"fmt"
	"go.elastic.co/ecszap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"time"
)

func FileWriter() zapcore.WriteSyncer {
	l := &lumberjack.Logger{
		Filename:   fmt.Sprintf("logs/%s.logger", time.Now().Format("2006-01-02")),
		MaxSize:    2,  // megabytes
		MaxAge:     30, // days
		MaxBackups: 3,
		LocalTime:  false,
		Compress:   false,
	}

	return zapcore.AddSync(l)
}

func FileEncoder() ecszap.EncoderConfig {
	return ecszap.NewDefaultEncoderConfig()
}
