package app

import (
	"context"
	"github.com/labstack/echo/v4"
	"po/internal/response"
)

type Context struct {
	echo.Context
}

func (c *Context) GetContext() context.Context {
	return c.Context.Request().Context()
}

func (c *Context) R() *response.Response {
	return response.New(c.Context)
}
