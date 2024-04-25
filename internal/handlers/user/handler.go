package user

import (
	"github.com/labstack/echo/v4"
	"po/internal/app"
)

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h Handler) Index(c echo.Context) error {
	ctx, _ := app.GetContext(c)

	return ctx.R().Ok(map[string]any{
		"data": "Hello World",
	})
}
