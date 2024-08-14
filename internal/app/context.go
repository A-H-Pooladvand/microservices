package app

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"po/internal/Filter"
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

func (c *Context) Validate(v any) error {
	if err := c.Bind(v); err != nil {
		_ = c.R().BadRequest(err.Error())

		return errors.New("invalid request")
	}

	if err := c.Context.Validate(v); err != nil {
		_ = c.R().UnprocessableEntity(err.Error())

		return errors.New("invalid request")
	}

	return nil
}

func (c *Context) Filter() *Filter.Filter {
	filters := c.FormValue("filters")

	if filters == "" {
		return nil
	}

	filter := new(Filter.Filter)

	err := json.Unmarshal(
		[]byte(filters),
		filter,
	)

	if err != nil {
		zap.L().Panic("failed to unmarshal filters", zap.Error(err))

		return nil
	}

	return filter
}
