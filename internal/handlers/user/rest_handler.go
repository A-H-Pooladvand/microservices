package user

import (
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"po/api/proto/user/v1"
	"po/internal/app"
	_ "po/internal/entities"
	"po/pkg/trace"
)

type RestHandler struct {
	service *Service
	tracer  trace.Tracer
}

type RestHandlerParams struct {
	fx.In
	Service *Service
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
//	@Summary		Lists all users.
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
	ctx, _ := app.GetContext(c)

	tracer := h.tracer.FromContext(ctx.GetContext())

	spanContext, span := tracer.Start(
		ctx.GetContext(),
		"User RestHandler",
	)

	defer span.End()

	h.service.GetAllUsers(spanContext)

	conn, err := grpc.DialContext(
		spanContext,
		"127.0.0.1:8501",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := user.NewUserServiceClient(conn)

	r, err := client.Index(spanContext, &user.UserRequest{
		FirstName: "Amirhossein",
		LastName:  "Pooladvand",
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())

	return ctx.R().Ok(map[string]any{
		"data": "Hello World",
	})
}

func (h RestHandler) Show(c echo.Context) error {
	ctx, _ := app.GetContext(c)
	tracer := h.tracer.FromContext(ctx.GetContext())

	_, span := tracer.Start(
		c.Request().Context(),
		"User RestHandler",
	)

	defer span.End()

	return ctx.R().Ok("ok")
}
