package middlewares

import (
	"runtime/debug"

	"github.com/a-h-pooladvand/microservices/internal/log"
	"github.com/a-h-pooladvand/microservices/pkg/observability"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// RecoverConfig holds configuration for the recover middleware.
type RecoverConfig struct {
	StackAll           bool
	PrintStack         bool
	DisableStackAll    bool
	DisablePrintStack  bool
}

// Recover returns a middleware that recovers from panics.
func Recover(cfg RecoverConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					ctx := c.Request().Context()
					span := observability.SpanFromContext(ctx)

					var err error
					switch e := r.(type) {
					case error:
						err = e
					case string:
						err = echo.NewHTTPError(500, e)
					default:
						err = echo.NewHTTPError(500, "internal server error")
					}

					// Record error in span
					if span != nil {
						observability.SetSpanError(span, err)
					}

					// Get stack trace
					stack := debug.Stack()

					// Log the panic
					log.ErrorCtx(ctx, "panic recovered",
						zap.Any("panic", r),
						zap.ByteString("stack", stack),
					)

					// Return error
					c.Error(err)
				}
			}()

			return next(c)
		}
	}
}
