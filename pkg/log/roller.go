package log

import (
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"time"
)

func FileRoller() io.Writer {
	return &lumberjack.Logger{
		Filename:   fmt.Sprintf("logs/%s.log", time.Now().Format("2006-01-02")),
		MaxSize:    2,  // megabytes
		MaxAge:     30, // days
		MaxBackups: 3,
		LocalTime:  false,
		Compress:   false,
	}
}
