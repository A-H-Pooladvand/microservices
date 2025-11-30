package web

import (
	"github.com/a-h-pooladvand/microservices/internal/web/middlewares"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// RegisterMiddlewares registers all HTTP middlewares.
func RegisterMiddlewares(e *echo.Echo) {
	// Skip paths for observability
	skipPaths := map[string]bool{
		"/health":  true,
		"/metric":  true,
		"/metrics": true,
	}

	// Recovery middleware (must be first)
	e.Use(middlewares.Recover(middlewares.RecoverConfig{}))

	// CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"},
	}))

	// Request ID middleware
	e.Use(middleware.RequestID())

	// Tracing middleware (OpenTelemetry)
	e.Use(middlewares.Tracing(middlewares.TracingConfig{
		SkipPaths: skipPaths,
	}))

	// Metrics middleware
	e.Use(middlewares.Metrics(middlewares.MetricsConfig{
		Namespace: "http",
		SkipPaths: skipPaths,
	}))

	// Logging middleware
	e.Use(middlewares.Logging(middlewares.LoggingConfig{
		SkipPaths: skipPaths,
	}))

	// Custom context middleware
	e.Use(middlewares.Context)

	// Prometheus metrics endpoint handler
	e.Use(echoprometheus.NewMiddleware("app"))

	// Secure headers
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "DENY",
		HSTSMaxAge:            31536000,
		ContentSecurityPolicy: "default-src 'self'",
	}))

	// Body limit
	e.Use(middleware.BodyLimit("2M"))
}
