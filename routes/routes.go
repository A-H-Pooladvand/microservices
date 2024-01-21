package routes

import (
	"github.com/labstack/echo/v4"
	"go.elastic.co/apm/v2"
	"net/http"
	"os"
)

func Register(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		os.Setenv("ELASTIC_APM_SERVER_URL", "https://127.0.0.1:8200")
		os.Setenv("ELASTIC_APM_SECRET_TOKEN", "supersecrettoken")
		os.Setenv("ELASTIC_APM_SERVER_CA_CERT_FILE", "ca.crt")
		os.Setenv("ELASTIC_APM_LOG_LEVE", "debug")
		os.Setenv("ELASTIC_APM_LOG_FILE", "log.log")
		os.Setenv("ELASTIC_APM_ENVIRONMENT", "staging")

		tx := apm.DefaultTracer().StartTransaction("GET /api/v1", "request")
		defer tx.End()
		tx.Result = "HTTP 2xx"
		tx.Context.SetLabel("region", "us-east-1")
		return c.JSON(http.StatusOK, "OK")
	})
}
