package user

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"po/internal/app"
	"po/internal/handlers/user/dto"
	request2 "po/internal/handlers/user/request"
	"po/internal/transformer/user_transformer"
	"po/pkg/trace"
)

type RestHandler struct {
	service Service
	tracer  trace.Tracer
}

type RestHandlerParams struct {
	fx.In
	Service Service
	Tracer  trace.Tracer
}

func NewRestHandler(params RestHandlerParams) RestHandler {
	return RestHandler{
		service: params.Service,
		tracer:  params.Tracer,
	}
}

// Index godoc
//
//	@Summary		List all users
//	@Description	Retrieves a list of all users in the system.
//	@tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Account ID"
//	@Success		200	{object}	entities.UserResponse
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/api/v1/users [get]
func (h RestHandler) Index(c echo.Context) error {
	ctx := app.GetContext(c)

	tracer := h.tracer.FromContext(ctx.GetContext())

	spanContext, span := tracer.Start(
		ctx.GetContext(),
		"User RestHandler",
	)

	defer span.End()

	users, err := h.service.GetAllUsers(spanContext, dto.GetAllUsers{
		Filter: ctx.Filter(),
	})

	if err != nil {
		return ctx.R().SetMessage(err.Error()).NotFound()
	}

	return ctx.R().Ok(user_transformer.All(users))
}

func (h RestHandler) Show(c echo.Context) error {
	ctx := app.GetContext(c)
	tracer := h.tracer.FromContext(ctx.GetContext())

	_, span := tracer.Start(
		c.Request().Context(),
		"User RestHandler",
	)

	defer span.End()

	return ctx.R().Ok("ok")
}

func (h RestHandler) Create(c echo.Context) error {
	ctx := app.GetContext(c)

	tracer := h.tracer.FromContext(ctx.GetContext())

	_, span := tracer.Start(
		c.Request().Context(),
		"User RestHandler",
	)

	defer span.End()

	var request request2.Request

	if err := ctx.Validate(&request); err != nil {
		return err
	}

	return ctx.R().Ok("ok")
}
