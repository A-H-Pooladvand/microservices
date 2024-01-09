package zlog

import (
	"go.elastic.co/ecszap"
	"go.uber.org/zap/zapcore"
	"net"
)

func LogstashWriter() zapcore.WriteSyncer {
	conn, err := net.Dial("tcp", "127.0.0.1:5000") // Replace "localhost" with your Logstash server's hostname or IP address
	if err != nil {
		panic(err)
	}

	//defer conn.Close()

	return zapcore.AddSync(conn)
}

func LogstashEncoder() ecszap.EncoderConfig {
	return ecszap.NewDefaultEncoderConfig()
}
