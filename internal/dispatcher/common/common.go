package common

import (
	"github.com/Nerufa/go-shared/provider"
	"github.com/ProtocolONE/geoip-service/pkg/proto"
	"github.com/labstack/echo/v4"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/grpc"
	paylink "github.com/paysuper/paysuper-payment-link/proto"
	"github.com/paysuper/paysuper-recurring-repository/pkg/proto/repository"
	reporterProto "github.com/paysuper/paysuper-reporter/pkg/proto"
	tax_service "github.com/paysuper/paysuper-tax-service/proto"
	"gopkg.in/go-playground/validator.v9"
)

const (
	Prefix                   = "internal.dispatcher"
	UnmarshalKey             = "dispatcher"
	UnmarshalGlobalConfigKey = "dispatcher.global"
	AuthProjectGroupPath     = "/api/v1"
	AccessGroupPath          = "/api/v1/s"
	AuthUserGroupPath        = "/admin/api/v1"
	WebHookGroupPath         = "/webhook"
)

// Cursor
type Cursor struct {
	Limit, Offset int32
	Sort          []string
}

// ExtractUserContext
func ExtractUserContext(ctx echo.Context) *AuthUser {
	if user, ok := ctx.Get("user").(*AuthUser); ok {
		return user
	}
	return &AuthUser{}
}

// ExtractRawBodyContext
func ExtractRawBodyContext(ctx echo.Context) []byte {
	if rawBody, ok := ctx.Get("rawBody").([]byte); ok {
		return rawBody
	}
	return nil
}

// ExtractCursorContext
func ExtractCursorContext(ctx echo.Context) *Cursor {
	if cursor, ok := ctx.Get("cursor").(*Cursor); ok {
		return cursor
	}
	return &Cursor{}
}

// SetUserContext
func SetUserContext(ctx echo.Context, user *AuthUser) {
	ctx.Set("user", user)
}

// SetRawBodyContext
func SetRawBodyContext(ctx echo.Context, rawBody []byte) {
	ctx.Set("rawBody", rawBody)
}

// SetCursorContext
func SetCursorContext(ctx echo.Context, cursor *Cursor) {
	ctx.Set("cursor", cursor)
}

// Groups
type Groups struct {
	AuthProject *echo.Group
	Access      *echo.Group
	AuthUser    *echo.Group
	WebHooks    *echo.Group
	Common      *echo.Echo
}

// Handler
type Handler interface {
	Route(groups *Groups)
}

// Validate
type Validator interface {
	Use(validator *validator.Validate)
}

// Services
type Services struct {
	Repository repository.RepositoryService
	Geo        proto.GeoIpService
	Billing    grpc.BillingService
	Tax        tax_service.TaxService
	PayLink    paylink.PaylinkService
	Reporter   reporterProto.ReporterService
}

// Handlers
type Handlers []Handler

// HandlerSet
type HandlerSet struct {
	Services Services
	Validate *validator.Validate
	AwareSet provider.AwareSet
}

// AuthUser
type AuthUser struct {
	Id        string
	Name      string
	Email     string
	Roles     map[string]bool
	Merchants map[string]bool
}
