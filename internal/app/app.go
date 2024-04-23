package app

import (
	"context"
	"github.com/fatih/color"
	"github.com/labstack/echo/v4"
	"os"
	"po/pkg/logstash"
	"po/pkg/postgres"
	"strings"
	"sync"
)

var (
	app  *App
	once sync.Once
)

type App struct {
	Ctx      *context.Context
	Logstash *logstash.Client
	Postgres *postgres.Client
}

func (a *App) GracefulShutdown() {
	a.Postgres.Shutdown()
	a.Logstash.Shutdown()
}

func (a *App) Context() context.Context {
	return *a.Ctx
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

func Production() bool {
	env := strings.ToLower(os.Getenv("APP_ENV"))

	return env == "prod" || env == "production"
}

func Local() bool {
	return !Production()
}

func LocalMessage() {
	if Production() {
		return
	}

	warning := color.New(color.FgHiYellow).Add(color.Bold)
	_, _ = warning.Println("----------------------------------------------------------------------------------------------------")
	_, _ = warning.Println("| ⚠️ Warning: Application is running in LOCAL environment. ⚠️                                      |")
	_, _ = warning.Println("| ⚠️ If this is unintended, please switch to the PRODUCTION environment for accurate results. ⚠️   |")
	_, _ = warning.Println("----------------------------------------------------------------------------------------------------")
}

func GetContext(c echo.Context) (*Context, bool) {
	ctx, ok := c.(*Context)

	return ctx, ok
}
