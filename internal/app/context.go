package app

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/a-h-pooladvand/microservices/internal/filter"
	metric2 "github.com/a-h-pooladvand/microservices/internal/metric"
	"github.com/a-h-pooladvand/microservices/internal/response"
	"github.com/a-h-pooladvand/microservices/internal/trace"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/metric"
	tracer "go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// Context represents the context of the request
type Context struct {
	echo.Context
	tracerContext context.Context
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

		return errors.New("failed to bind request data")
	}

	if err := c.Context.Validate(v); err != nil {
		_ = c.R().UnprocessableEntity(err.Error())

		return errors.New("failed to validate request data")
	}

	return nil
}

func (c *Context) Filter() *filter.Filter {
	filters := c.FormValue("filters")

	if filters == "" {
		return nil
	}

	f := new(filter.Filter)

	err := json.Unmarshal(
		[]byte(filters),
		f,
	)

	if err != nil {
		zap.L().Panic("failed to unmarshal filters", zap.Error(err))

		return nil
	}

	return f
}

func (c *Context) Start(spanName string, opts ...tracer.SpanStartOption) (context.Context, tracer.Span) {
	if c.tracerContext == nil {
		tc, span := trace.Tracer().Start(c.Request().Context(), spanName, opts...)

		c.SetContext(tc)

		return tc, span
	}

	tc, span := trace.
		FromContext(c.tracerContext).
		Start(c.tracerContext, spanName, opts...)

	c.SetContext(tc)

	return tc, span
}

func (c *Context) Span() tracer.Span {
	return tracer.SpanFromContext(c.tracerContext)
}

func (c *Context) TracerContext() context.Context {
	return c.tracerContext
}

func (c *Context) Meter() metric.Meter {
	return metric2.Meter()
}

func (c *Context) SetContext(ctx context.Context) {
	c.tracerContext = ctx
}
