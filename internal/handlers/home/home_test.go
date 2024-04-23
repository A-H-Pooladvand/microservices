package home

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"po/internal/app"
	"po/internal/webserver/middlewares"
	"reflect"
	"testing"
)

func TestHomeHandler_Home_Success(t *testing.T) {
	// Set up mock dependencies
	// ... (if necessary)

	// Create a new HomeHandler instance
	handler := New()

	// Create a mock Echo context
	e := echo.New()
	e.Use(middlewares.Context)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	ctx := &app.Context{
		Context: c,
	}

	_ = handler.Home(ctx)
	
	// Assert the response status and body
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	expectedResponse := map[string]any{
		"ok": true,
		"data": map[string]any{
			"message": "Hello World",
			"ok":      true,
		},
	}

	var actualResponse map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &actualResponse); err != nil {
		t.Errorf("Failed to unmarshal response body: %v", err)
	}

	if !reflect.DeepEqual(expectedResponse, actualResponse) {
		t.Errorf("Unexpected response body: %v", actualResponse)
	}
}
