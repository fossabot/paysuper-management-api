package http

import "github.com/labstack/echo/v4"

const (
	Prefix           = "internal.http"
	UnmarshalKey     = "http"
	UnmarshalKeyBind = "http.bind"
)

// Dispatcher
type Dispatcher interface {
	Dispatch(http *echo.Echo) error
}
