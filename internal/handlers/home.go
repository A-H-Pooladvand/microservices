package handlers

import (
	"github.com/labstack/echo/v4"
	"po/internal/app"
	"po/internal/vault"
)

func Home(c echo.Context) error {
	ctx := c.(*app.Context)

	client, err := vault.New()

	if err != nil {
		return ctx.R().BadRequest(err.Error())
	}

	secret, err := client.Get(c.Request().Context(), "database")

	if err != nil {
		return ctx.R().BadRequest(err.Error())
	}

	return ctx.R().Ok(secret.Data)
}
