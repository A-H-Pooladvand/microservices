package middlewares

import (
	"time"

	"github.com/a-h-pooladvand/microservices/internal/log"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// LoggingConfig holds configuration for the logging middleware.
type LoggingConfig struct {
	SkipPaths map[string]bool
}

// Logging returns a middleware that logs HTTP requests with trace context.
func Logging(cfg LoggingConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip logging for certain paths
			if cfg.SkipPaths != nil && cfg.SkipPaths[c.Path()] {
				return next(c)
			}

			startTime := time.Now()
			req := c.Request()
			ctx := req.Context()

			// Process request
			err := next(c)

			// Log request
			duration := time.Since(startTime)
			statusCode := c.Response().Status

			fields := []zap.Field{
				zap.String("method", req.Method),
				zap.String("path", req.URL.Path),
				zap.String("query", req.URL.RawQuery),
				zap.Int("status", statusCode),
				zap.Duration("duration", duration),
				zap.Int64("response_size", c.Response().Size),
				zap.String("client_ip", c.RealIP()),
				zap.String("user_agent", req.UserAgent()),
			}

			if err != nil {
				fields = append(fields, zap.Error(err))
				log.ErrorCtx(ctx, "HTTP request failed", fields...)
			} else if statusCode >= 500 {
				log.ErrorCtx(ctx, "HTTP server error", fields...)
			} else if statusCode >= 400 {
				log.WarnCtx(ctx, "HTTP client error", fields...)
			} else {
				log.InfoCtx(ctx, "HTTP request completed", fields...)
			}

			return err
		}
	}
}
