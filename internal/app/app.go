package app

import (
	"github.com/a-h-pooladvand/microservices/internal/log"
	"github.com/labstack/echo/v4"
	"os"
)

func GetContext(c echo.Context) *Context {
	ctx, ok := c.(*Context)

	if !ok {
		log.Panic("unable to get context")
	}

	return ctx
}

func Env() string {
	return os.Getenv("APP_ENV")
}

func IsProd() bool {
	return Env() == ""
}

func IsDev() bool {
	return !IsProd()
}
