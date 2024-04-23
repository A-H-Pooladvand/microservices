//go:build wireinject
// +build wireinject

package cmd

import (
	"github.com/google/wire"
)

func InitializeGreeter() Greeter {
	wire.NewSet()
	wire.Build(NewMessage, NewGreeter)
	return Greeter{}
}

type Message string

type FileLogger struct {
	Message Message
}

func NewFileLogger(message Message) {
	wire.NewSet()
	return FileLogger{Message: message}
}

type DBLogger struct {
	Message Message
}

func NewDBLogger(message Message) {
	return DBLogger{Message: message}
}
