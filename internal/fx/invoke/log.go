package invoke

import (
	"context"
	"errors"
	"fmt"
	"github.com/a-h-pooladvand/microservices/config"
	"go.elastic.co/apm/module/apmzap/v2"
	"go.elastic.co/ecszap"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"syscall"
	"time"
)

func Log(lc fx.Lifecycle, config *config.Config) {
	config.Attach(bootLogger)
	bootLogger(config)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			err := zap.L().Sync()

			if errors.Is(err, syscall.EINVAL) {
				return nil
			}

			return err
		},
	})
}

func bootLogger(config *config.Config) {
	// Disabled logstash since we currently want to collect logs via filebeat
	/*ls := zapcore.AddSync(
		l.Connection(),
	)*/

	var cores []zapcore.Core

	if config.App.Debuggable() {
		//Stdout Writer
		std := zapcore.NewCore(
			StdoutEncoder(),
			zapcore.AddSync(os.Stdout),
			zap.NewAtomicLevelAt(zapcore.DebugLevel),
		)

		cores = append(cores, std)
	}

	file := ecszap.NewCore(
		ecszap.NewDefaultEncoderConfig(),
		FileWriter(),
		zap.NewAtomicLevelAt(fileLevel(config)),
	)

	cores = append(cores, file)

	core := zapcore.NewTee(
		cores...,

	// Logstash Writer
	//ecszap.NewCore(ecszap.NewDefaultEncoderConfig(), ls, zap.NewAtomicLevelAt(zap.InfoLevel)),
	)

	logger := zap.New(
		core,
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.WrapCore((&apmzap.Core{}).WrapCore),
	)

	zap.ReplaceGlobals(logger)
}

func StdoutEncoder() zapcore.Encoder {
	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zapcore.NewConsoleEncoder(developmentCfg)
}

func FileWriter() zapcore.WriteSyncer {
	l := &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/logs/%s.log", "", time.Now().Format("2006-01-02")),
		MaxSize:    2,  // megabytes
		MaxAge:     30, // days
		MaxBackups: 3,
		LocalTime:  false,
		Compress:   false,
	}

	return zapcore.AddSync(l)
}

func fileLevel(config *config.Config) zapcore.Level {
	if config.App.Debuggable() {
		return zapcore.DebugLevel
	}

	return zapcore.InfoLevel
}
