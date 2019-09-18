package api

import (
	"context"
	"encoding/json"
	"fmt"
	jwtverifier "github.com/ProtocolONE/authone-jwt-verifier-golang"
	"gopkg.in/go-playground/validator.v9"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	jwtMiddleware "github.com/ProtocolONE/authone-jwt-verifier-golang/middleware/echo"
	"github.com/ProtocolONE/geoip-service/pkg"
	"github.com/ProtocolONE/geoip-service/pkg/proto"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/micro/go-micro"
	"github.com/micro/go-plugins/client/selector/static"
	casbinMiddleware "github.com/paysuper/echo-casbin-middleware"
	"github.com/paysuper/paysuper-billing-server/pkg"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/billing"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/grpc"
	"github.com/paysuper/paysuper-management-api/config"
	"github.com/paysuper/paysuper-management-api/utils"
	paylinkServiceConst "github.com/paysuper/paysuper-payment-link/pkg"
	"github.com/paysuper/paysuper-payment-link/proto"
	"github.com/paysuper/paysuper-recurring-repository/pkg/constant"
	"github.com/paysuper/paysuper-recurring-repository/pkg/proto/repository"
	reporterPkg "github.com/paysuper/paysuper-reporter/pkg"
	reporterProto "github.com/paysuper/paysuper-reporter/pkg/proto"
	taxServiceConst "github.com/paysuper/paysuper-tax-service/pkg"
	"github.com/paysuper/paysuper-tax-service/proto"
	"github.com/ttacon/libphonenumber"
	"go.uber.org/zap"
)

var funcMap = template.FuncMap{
	"Marshal": func(v interface{}) template.JS {
		a, _ := json.Marshal(v)
		return template.JS(a)
	},
}

type Template struct {
	templates *template.Template
}

type ServerInitParams struct {
	Config     *config.Config
	Logger     *zap.SugaredLogger
	HttpScheme string
	Auth1      *config.Auth1
}

type Merchant struct {
	Identifier string
}

type GetParams struct {
	limit  int32
	offset int32
	sort   []string
}

type Order struct {
	PayerPhone *libphonenumber.PhoneNumber
}

type AuthUser struct {
	Id        string
	Name      string
	Email     string
	Roles     map[string]bool
	Merchants map[string]bool
}

type Api struct {
	Http     *echo.Echo
	config   *config.Config
	logger   *zap.SugaredLogger
	validate *validator.Validate

	accessRouteGroup    *echo.Group
	webhookRouteGroup   *echo.Group
	apiAuthProjectGroup *echo.Group
	jwtVerifier         *jwtverifier.JwtVerifier

	authUserRouteGroup *echo.Group
	authUser           *AuthUser

	httpScheme string

	service        micro.Service
	serviceContext context.Context
	serviceCancel  context.CancelFunc

	repository      repository.RepositoryService
	geoService      proto.GeoIpService
	billingService  grpc.BillingService
	taxService      tax_service.TaxService
	paylinkService  paylink.PaylinkService
	reporterService reporterProto.ReporterService

	rawBody      string
	reqSignature string

	Merchant
	GetParams
	Order
}

func NewServer(p *ServerInitParams) (*Api, error) {
	api := &Api{
		Http:       echo.New(),
		logger:     p.Logger,
		validate:   validator.New(),
		httpScheme: p.HttpScheme,
		config:     p.Config,
	}
	api.InitService()

	jwtVerifierSettings := jwtverifier.Config{
		ClientID:     p.Auth1.ClientId,
		ClientSecret: p.Auth1.ClientSecret,
		Scopes:       []string{"openid", "offline"},
		RedirectURL:  p.Auth1.RedirectUrl,
		Issuer:       p.Auth1.Issuer,
	}
	api.jwtVerifier = jwtverifier.NewJwtVerifier(jwtVerifierSettings)

	renderer := &Template{
		templates: template.Must(template.New("").Funcs(funcMap).ParseGlob("web/template/*.html")),
	}
	api.Http.Renderer = renderer

	api.Http.Static("/", "web/static")
	api.Http.Static("/spec", "spec")

	err := api.registerValidators()

	if err != nil {
		return nil, err
	}

	api.accessRouteGroup = api.Http.Group("/api/v1/s")

	api.accessRouteGroup.Use(
		jwtMiddleware.AuthOneJwtCallableWithConfig(
			api.jwtVerifier,
			func(ui *jwtverifier.UserInfo) {
				api.Merchant.Identifier = ui.UserID
			},
		),
	)
	api.accessRouteGroup.Use(middleware.Logger())
	api.accessRouteGroup.Use(middleware.Recover())

	api.authUserRouteGroup = api.Http.Group(apiAuthUserGroupPath)
	api.authUserRouteGroup.Use(
		jwtMiddleware.AuthOneJwtCallableWithConfig(
			api.jwtVerifier,
			func(ui *jwtverifier.UserInfo) {
				api.authUser = &AuthUser{
					Id:        ui.UserID,
					Name:      "System User",
					Merchants: make(map[string]bool),
					Roles:     make(map[string]bool),
				}
			},
		),
	)
	api.authUserRouteGroup.Use(api.getUserDetailsMiddleware)
	api.authUserRouteGroup.Use(casbinMiddleware.Middleware(api.service.Client(), casbinMiddleware.EnforceModePermissive))
	api.webhookRouteGroup = api.Http.Group(apiWebHookGroupPath)
	api.webhookRouteGroup.Use(middleware.BodyDump(func(ctx echo.Context, reqBody, resBody []byte) {
		data := []interface{}{
			"request_headers", utils.RequestResponseHeadersToString(ctx.Request().Header),
			"request_body", string(reqBody),
			"response_headers", utils.RequestResponseHeadersToString(ctx.Response().Header()),
			"response_body", string(resBody),
		}

		api.logger.Infow(ctx.Path(), data...)
	}))
	api.webhookRouteGroup.Use(api.RawBodyMiddleware)

	api.apiAuthProjectGroup = api.Http.Group(apiAuthProjectGroupPath)
	api.apiAuthProjectGroup.Use(middleware.BodyDump(func(ctx echo.Context, reqBody, resBody []byte) {
		data := []interface{}{
			"request_headers", utils.RequestResponseHeadersToString(ctx.Request().Header),
			"request_body", string(reqBody),
			"response_headers", utils.RequestResponseHeadersToString(ctx.Response().Header()),
			"response_body", string(resBody),
		}

		api.logger.Infow(ctx.Path(), data...)
	}))
	api.apiAuthProjectGroup.Use(api.RawBodyMiddleware)
	api.Http.Use(api.RawBodyMiddleware)

	api.Http.Use(api.LimitOffsetSortMiddleware)
	api.Http.Use(middleware.Logger())
	api.Http.Use(middleware.Recover())
	api.Http.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowHeaders: []string{"authorization", "content-type"},
	}))

	api.
		InitCountryRoutes().
		InitProductRoutes().
		InitProjectRoutes().
		InitOrderV1Routes().
		InitPaylinkRoutes().
		InitCardPayWebHookRoutes().
		InitPaymentCostRoutes().
		initTaxesRoutes().
		initTokenRoutes().
		initZipCodeRoutes().
		initPaymentMethodRoutes().
		initPriceGroupRoutes().
		initUserProfileRoutes().
		initVatReportsRoutes().
		initRoyaltyReportsRoutes().
		initKeyProductRoutes().
		initDashboardRoutes().
		initReportFileRoute()

	_, err = api.initOnboardingRoutes()

	if err != nil {
		return nil, err
	}

	api.Http.GET("/docs", func(ctx echo.Context) error {
		return ctx.Render(http.StatusOK, "docs.html", map[string]interface{}{})
	})

	return api, nil
}

func (api *Api) Start() error {
	go func() {
		if err := api.service.Run(); err != nil {
			return
		}
	}()

	go func() {
		if err := api.Http.Start(":3001"); err != nil {
			api.Http.Logger.Info("shutting down the server")
		}
	}()

	return nil
}

func (api *Api) InitService() {
	api.serviceContext, api.serviceCancel = context.WithCancel(context.Background())

	options := []micro.Option{
		micro.Name("p1payapi"),
		micro.Version(constant.PayOneMicroserviceVersion),
	}

	log.Println("Initialize micro service")

	if os.Getenv("MICRO_SELECTOR") == "static" {
		log.Println("Use micro selector `static`")
		options = append(options, micro.Selector(static.NewSelector()))
	}

	api.service = micro.NewService(options...)
	api.service.Init()

	api.repository = repository.NewRepositoryService(constant.PayOneRepositoryServiceName, api.service.Client())
	api.geoService = proto.NewGeoIpService(geoip.ServiceName, api.service.Client())
	api.billingService = grpc.NewBillingService(pkg.ServiceName, api.service.Client())
	api.taxService = tax_service.NewTaxService(taxServiceConst.ServiceName, api.service.Client())
	api.paylinkService = paylink.NewPaylinkService(paylinkServiceConst.ServiceName, api.service.Client())
	api.reporterService = reporterProto.NewReporterService(reporterPkg.ServiceName, api.service.Client())
}

func (api *Api) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := api.Http.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Http server exiting")

	api.serviceCancel()
	log.Println("Micro server exiting")
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func (api *Api) getUserDetailsMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		auth := ctx.Request().Header.Get(echo.HeaderAuthorization)

		if auth == "" {
			return errorMessageAuthorizationHeaderNotFound
		}

		r := regexp.MustCompile(requestAuthorizationTokenRegex)
		match := r.FindStringSubmatch(auth)

		if len(match) < 1 {
			return errorMessageAuthorizationTokenNotFound
		}

		u, err := api.jwtVerifier.GetUserInfo(ctx.Request().Context(), match[1])

		if err != nil {
			return errorMessageAuthorizedUserNotFound
		}

		api.authUser.Email = u.Email

		return next(ctx)
	}
}

func (api *Api) getValidationError(err error) (rspErr *grpc.ResponseErrorMessage) {
	vErr := err.(validator.ValidationErrors)[0]
	val, ok := validationErrors[vErr.Field()]

	if ok {
		rspErr = val
	} else {
		val, ok := validationNamespaceErrors[vErr.StructNamespace()]

		if ok {
			rspErr = val
		} else {
			if vErr.Tag() == requestParameterZipUsa {
				rspErr = errorMessageIncorrectZip
			} else {
				rspErr = errorValidationFailed
			}
		}
	}

	rspErr.Details = fmt.Sprintf(errorMessageMask, vErr.Field(), vErr.Tag())

	return rspErr
}

func (api *Api) onboardingBeforeHandler(st interface{}, ctx echo.Context) *echo.HTTPError {
	err := ctx.Bind(st)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errorRequestParamsIncorrect)
	}

	err = api.validate.Struct(st)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.getValidationError(err))
	}

	return nil
}

func (api *Api) logError(msg string, data []interface{}) {
	zap.S().Errorw(fmt.Sprintf("[PAYSUPER_MANAGEMENT_API] %s", msg), data...)
}

func (api *Api) isProductionEnvironment() bool {
	return api.config.Environment == EnvironmentProduction
}

func (api *Api) checkProjectAuthRequestSignature(ctx echo.Context, projectId string) error {
	signature := ctx.Request().Header.Get(HeaderXApiSignatureHeader)

	if signature == "" {
		return echo.NewHTTPError(http.StatusBadRequest, errorMessageSignatureHeaderIsEmpty)
	}

	req := &grpc.CheckProjectRequestSignatureRequest{Body: api.rawBody, ProjectId: projectId, Signature: signature}
	rsp, err := api.billingService.CheckProjectRequestSignature(ctx.Request().Context(), req)

	if err != nil {
		zap.S().Errorf("internal error", "err", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, errorUnknown)
	}

	if rsp.Status != pkg.ResponseStatusOk {
		return echo.NewHTTPError(int(rsp.Status), rsp.Message)
	}

	return nil
}

func (api *Api) registerValidators() error {
	if err := api.validate.RegisterValidation("phone", api.PhoneValidator); err != nil {
		return err
	}

	if err := api.validate.RegisterValidation("uuid", api.UuidValidator); err != nil {
		return err
	}

	if err := api.validate.RegisterValidation("zip_usa", api.ZipUsaValidator); err != nil {
		return err
	}

	if err := api.validate.RegisterValidation("name", api.NameValidator); err != nil {
		return err
	}

	if err := api.validate.RegisterValidation("position", api.PositionValidator); err != nil {
		return err
	}

	api.validate.RegisterStructValidation(api.CompanyValidator, grpc.UserProfileCompany{})
	api.validate.RegisterStructValidation(api.MerchantCompanyValidator, billing.MerchantCompanyInfo{})
	api.validate.RegisterStructValidation(api.MerchantTariffRatesValidator, grpc.GetMerchantTariffRatesRequest{})

	if err := api.validate.RegisterValidation("company_name", api.CompanyNameValidator); err != nil {
		return err
	}

	if err := api.validate.RegisterValidation("swift", api.SwiftValidator); err != nil {
		return err
	}

	if err := api.validate.RegisterValidation("city", api.CityValidator); err != nil {
		return err
	}

	if err := api.validate.RegisterValidation("world_region", api.WorldRegionValidator); err != nil {
		return err
	}

	return nil
}
