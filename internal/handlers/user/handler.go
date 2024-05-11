package user

import (
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"po/internal/app"
	"po/internal/entities"
)

type Handler struct {
	Service *Service
	Tracer  trace.Tracer
}

type HandlerParams struct {
	fx.In
	Service *Service
	Tracer  trace.Tracer
}

func NewHandler(params HandlerParams) *Handler {
	return &Handler{
		Service: params.Service,
		Tracer:  params.Tracer,
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
func (h Handler) Index(c echo.Context) error {
	ctx, _ := app.GetContext(c)
	_ = entities.UserResponse{}

	//// Attributes represent additional key-value descriptors that can be bound
	//// to a metric observer or recorder.
	//commonAttrs := []attribute.KeyValue{
	//	attribute.String("attrA", "chocolate"),
	//	attribute.String("attrB", "raspberry"),
	//	attribute.String("attrC", "vanilla"),
	//}
	//
	//tracer := otel.Tracer("test-tracer")
	//
	//// work begins
	//spanContext, span := tracer.Start(
	//	c.Request().Context(),
	//	"CollectorExporter-Example",
	//	trace.WithAttributes(commonAttrs...),
	//)
	//defer span.End()
	//for i := 0; i < 10; i++ {
	//	_, iSpan := tracer.Start(spanContext, fmt.Sprintf("Sample-%d", i))
	//	log.Printf("Doing really hard work (%d / 10)\n", i+1)
	//
	//	<-time.After(time.Second)
	//	iSpan.End()
	//}
	//
	//log.Printf("Done!")

	return ctx.R().Ok(map[string]any{
		"data": "Hello World",
	})
}
