package user

import (
	"github.com/a-h-pooladvand/microservices/internal/app"
	"github.com/a-h-pooladvand/microservices/internal/domain/service"
	"github.com/a-h-pooladvand/microservices/internal/dto"
	"github.com/a-h-pooladvand/microservices/internal/handler/http/user/request"
	"github.com/a-h-pooladvand/microservices/internal/log"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type HandlerParams struct {
	fx.In
	UserService service.User
}

type HttpHandler struct {
	userService service.User
}

func NewHttpHandler(p HandlerParams) HttpHandler {
	return HttpHandler{
		userService: p.UserService,
	}
}

// Create godoc
//
//	@Summary		Create a new user
//	@Description	Creates a new user in the system.
//	@tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		request.Request	true	"User data"
//	@Success		200		{object}	response.UserResponse
//	@Failure		400		{object}	response.Response
//	@Failure		422		{object}	response.Response
//	@Failure		404		{object}	response.Response
//	@Failure		500		{object}	response.Response
//	@Router			/api/v1/users [post]
func (h HttpHandler) Create(c echo.Context) error {
	ctx := app.GetContext(c)

	_, span := ctx.Start("Handler", trace.WithAttributes(
		attribute.String("component", "User"),
		attribute.String("operation", "Create"),
	))
	defer span.End()

	var r request.Request

	if err := ctx.Validate(&r); err != nil {
		span.RecordError(err)
		return err
	}

	span.AddEvent("inserting user into database")
	user, err := h.userService.Create(ctx.GetContext(), r.ToUserEntity())

	if err != nil {
		log.Error("failed to create user", zap.Error(err), zap.String("name", r.Name))
		span.RecordError(err)
		return ctx.R().ServerError()
	}
	span.AddEvent("user created successfully")

	return ctx.R().Ok(
		dto.ToUserResponse(*user),
	)
}
