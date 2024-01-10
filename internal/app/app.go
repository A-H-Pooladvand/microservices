package app

import (
	"po/pkg/logstash"
	"po/pkg/postgres"
	"sync"
)

var (
	app  *App
	once sync.Once
)

type App struct {
	Logstash *logstash.Client
	Postgres *postgres.Client
}

func (a *App) GracefulShutdown() {
	a.Postgres.Shutdown()
	a.Logstash.Shutdown()
}

func NewSingleton() *App {
	once.Do(func() {
		app = &App{}
	})

	return app
}

func Get() *App {
	return NewSingleton()
}
