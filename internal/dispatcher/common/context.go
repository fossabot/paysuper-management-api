package common

import (
	"github.com/labstack/echo/v4"
)

// ContextWrapperCallback
func ContextWrapperCallback(fn func(ctx echo.Context, next echo.HandlerFunc) error) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			return fn(c, next)
		}
	}
}
