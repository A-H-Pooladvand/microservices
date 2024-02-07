package response

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type Response struct {
	Ctx        echo.Context `json:"-"`
	StatusCode int          `json:"-"`
	OK         bool         `json:"ok"`
	Message    string       `json:"message,omitempty"`
	Data       any          `json:"data,omitempty"`
	Errors     any          `json:"errors,omitempty"`
}

func New(ctx echo.Context) *Response {
	return &Response{Ctx: ctx}
}

func (r *Response) Ok(data any) error {
	return r.Success(http.StatusOK, data)
}

func (r *Response) Forbidden() error {
	return r.Ctx.JSON(http.StatusForbidden, nil)
}

func (r *Response) BadRequest(msg string) error {
	return r.Ctx.JSON(http.StatusBadRequest, map[string]string{
		"message": msg,
	})
}

func (r *Response) NotFound() error {
	return r.Error(http.StatusNotFound, "")
}

func (r *Response) ServerError() error {
	return r.Error(http.StatusInternalServerError, "")
}

func (r *Response) Conflict() error {
	return r.Error(http.StatusConflict, "")
}

func (r *Response) json() error {
	return r.Ctx.JSON(r.StatusCode, r)
}

func (r *Response) setStatusCode(code int) {
	r.StatusCode = code
}

func (r *Response) setOk(v bool) {
	r.OK = v
}

func (r *Response) setData(i any) {
	r.Data = i
}

func (r *Response) SetMessage(message string) *Response {
	r.Message = message

	return r
}

func (r *Response) Error(statusCode int, v any) error {
	r.setOk(false)
	r.setStatusCode(statusCode)

	r.Errors = v

	return r.json()
}

func (r *Response) Success(statusCode int, v any) error {
	r.setOk(true)
	r.setData(v)
	r.setStatusCode(statusCode)

	return r.json()
}

func (r *Response) UnprocessableEntity(v any) error {
	return r.Error(http.StatusUnprocessableEntity, v)
}
