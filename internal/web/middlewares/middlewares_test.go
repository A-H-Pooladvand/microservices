package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/a-h-pooladvand/microservices/internal/web/middlewares"
	"github.com/labstack/echo/v4"
)

func TestTracingMiddleware(t *testing.T) {
	e := echo.New()

	// Create handler
	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}

	// Apply middleware
	middleware := middlewares.Tracing(middlewares.TracingConfig{})
	h := middleware(handler)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/test")

	// Execute
	err := h(c)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestTracingMiddleware_SkipPaths(t *testing.T) {
	e := echo.New()

	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}

	middleware := middlewares.Tracing(middlewares.TracingConfig{
		SkipPaths: map[string]bool{
			"/health": true,
		},
	})
	h := middleware(handler)

	// Create request for skipped path
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/health")

	err := h(c)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestMetricsMiddleware(t *testing.T) {
	e := echo.New()

	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}

	middleware := middlewares.Metrics(middlewares.MetricsConfig{})
	h := middleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/test")

	err := h(c)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestLoggingMiddleware(t *testing.T) {
	e := echo.New()

	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}

	middleware := middlewares.Logging(middlewares.LoggingConfig{})
	h := middleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/test")

	err := h(c)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestRecoverMiddleware(t *testing.T) {
	e := echo.New()

	// Handler that panics
	handler := func(c echo.Context) error {
		panic("test panic")
	}

	middleware := middlewares.Recover(middlewares.RecoverConfig{})
	h := middleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/test")

	// Should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("middleware should have recovered from panic, but got: %v", r)
		}
	}()

	_ = h(c) // Error is expected from the panic recovery
}

func TestContextMiddleware(t *testing.T) {
	e := echo.New()

	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}

	middleware := middlewares.Context
	h := middleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h(c)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}
