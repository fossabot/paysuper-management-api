package dispatcher

import (
	"bytes"
	"fmt"
	"github.com/Nerufa/go-shared/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/paysuper/paysuper-management-api/internal/dispatcher/common"
	"io/ioutil"
	"net/http"
	"strconv"
)

// RecoverMiddleware
func (d *Dispatcher) RecoverMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					d.L().Critical("[PANIC RECOVER] %s", logger.Args(err.Error()), logger.Stack("stacktrace"))
					c.Error(err)
				}
			}()
			return next(c)
		}
	}
}

// GetUserDetailsMiddleware
func (d *Dispatcher) GetUserDetailsMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		auth := ctx.Request().Header.Get(echo.HeaderAuthorization)

		if auth == "" {
			return common.ErrorMessageAuthorizationHeaderNotFound
		}

		match := common.TokenRegex.FindStringSubmatch(auth)

		if len(match) < 1 {
			return common.ErrorMessageAuthorizationTokenNotFound
		}

		u, err := d.appSet.JwtVerifier.GetUserInfo(ctx.Request().Context(), match[1])

		if err != nil {
			return common.ErrorMessageAuthorizedUserNotFound
		}

		user := common.ExtractUserContext(ctx)
		user.Email = u.Email
		common.SetUserContext(ctx, user)

		return next(ctx)
	}
}

// LimitOffsetSortPreMiddleware
func (d *Dispatcher) LimitOffsetSortPreMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Limit
		if c.Request().Method != http.MethodGet {
			return next(c)
		}
		limit, err := strconv.Atoi(c.QueryParam(common.QueryParameterNameLimit))
		if err != nil {
			limit = int(d.globalCfg.LimitDefault)
		}
		// Offset
		offset, err := strconv.Atoi(c.QueryParam(common.QueryParameterNameOffset))
		if err != nil {
			offset = int(d.globalCfg.OffsetDefault)
		}
		// Sort
		qParams := c.QueryParams()
		sort := common.DefaultSort
		if s, ok := qParams[common.QueryParameterNameSort]; ok {
			sort = s
		}
		//
		common.SetCursorContext(c, &common.Cursor{
			Limit:  int32(limit),
			Offset: int32(offset),
			Sort:   sort,
		})
		return next(c)
	}
}

// RawBodyPreMiddleware
func (d *Dispatcher) RawBodyPreMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		buf, _ := ioutil.ReadAll(c.Request().Body)
		rdr := ioutil.NopCloser(bytes.NewBuffer(buf))
		c.Request().Body = rdr
		common.SetRawBodyContext(c, buf)
		return next(c)
	}
}

// BodyDumpMiddleware
func (d *Dispatcher) BodyDumpMiddleware() echo.MiddlewareFunc {
	return middleware.BodyDump(func(ctx echo.Context, reqBody, resBody []byte) {
		data := map[string]interface{}{
			"request_headers":  common.RequestResponseHeadersToString(ctx.Request().Header),
			"request_body":     string(reqBody),
			"response_headers": common.RequestResponseHeadersToString(ctx.Response().Header()),
			"response_body":    string(resBody),
		}
		d.L().Info(ctx.Path(), logger.WithFields(data))
	})
}
