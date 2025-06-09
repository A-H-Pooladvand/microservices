package routes

import (
	"github.com/a-h-pooladvand/microservices/internal/handler"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// RegisterWebRoutes registers web routes.
func RegisterWebRoutes(e *echo.Echo, w *handler.Http) {
	e.GET("/metric", echoprometheus.NewHandler())
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("health", func(c echo.Context) error {
		return nil
	})

	v1 := e.Group("api/v1/")

	user := v1.Group("users")
	{
		user.POST("", w.User.Create)
	}

	//for _, route := range e.Routes() {
	//	fmt.Printf("method: %s Path: %s\n", route.Method, route.Path)
	//}
}
