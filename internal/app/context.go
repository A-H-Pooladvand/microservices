package app

import (
	"context"
	"github.com/labstack/echo/v4"
	"po/internal/response"
)

// Context represents the context of the request
type Context struct {
	echo.Context
}

// NewContext returns a new context instance
func NewContext(ctx echo.Context) *Context {
	return &Context{Context: ctx}
}

// GetContext returns the context of the request
func (c *Context) GetContext() context.Context {
	return c.Context.Request().Context()
}

// R returns a new response instance
func (c *Context) R() *response.Response {
	return response.New(c.Context)
}
