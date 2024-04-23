package home

import (
	"github.com/labstack/echo/v4"
	"po/internal/app"
)

type Handler struct{}

func New() Handler {
	return Handler{}
}

func (h Handler) Home(c echo.Context) error {
	ctx, ok := app.GetContext(c)

	if !ok {
		return ctx.R().SetMessage("unexpected error").ServerError()
	}

	//x := client.Get(context.Background(), "Amirhossein")

	return ctx.R().Ok(map[string]any{
		"ok":      ok,
		"message": "Hello World",
	})
}
