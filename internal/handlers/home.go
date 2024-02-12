package handlers

import (
	"github.com/labstack/echo/v4"
	"go.elastic.co/apm/v2"
	"po/internal/app"
)

func Home(c echo.Context) error {
	ctx := c.(*app.Context)

	tx := apm.DefaultTracer().StartTransaction("hello", "request")
	defer tx.End()

	tx.Context.SetCustom("key", "value")

	span := tx.StartSpan("work", "custom", nil)
	defer span.End()

	return ctx.R().Ok("Hello World")
}
