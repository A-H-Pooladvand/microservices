package routes

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Message string

func NewMessage() Message {
	return Message("Hello World")
}

func NewGreeter(m Message) Greeter {
	return Greeter{Message: m}
}

type Greeter struct {
	Message Message // <- adding a Message field
}

func NewEvent(g Greeter) Event {
	return Event{Greeter: g}
}

func (e Event) Start() {
	msg := e.Greeter.Greet()
	fmt.Println(msg)
}

type Event struct {
	Greeter Greeter // <- adding a Greeter field
}

func (g Greeter) Greet() Message {
	return g.Message
}

func Register(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {

		return c.JSON(http.StatusOK, "OK")
	})
}
