package dispatcher

import (
	"context"
	"github.com/Nerufa/go-shared/invoker"
	"github.com/Nerufa/go-shared/logger"
	"github.com/Nerufa/go-shared/provider"
	jwtverifier "github.com/ProtocolONE/authone-jwt-verifier-golang"
	jwtMiddleware "github.com/ProtocolONE/authone-jwt-verifier-golang/middleware/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/paysuper/paysuper-management-api/internal/dispatcher/common"
	"html/template"
	"net/http"
)

// Dispatcher
type Dispatcher struct {
	ctx    context.Context
	cfg    Config
	appSet AppSet
	provider.LMT
	globalCfg *common.Config
}

// dispatch
func (d *Dispatcher) Dispatch(echoHttp *echo.Echo) error {

	t, e := template.New("").Funcs(common.FuncMap).ParseGlob(d.cfg.WorkDir + "/web/template/*.html")
	if e != nil {
		return e
	}
	echoHttp.Renderer = common.NewTemplate(t)

	// Called after routes
	echoHttp.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: logger.NewLevelWriter(d.L(), logger.LevelInfo),
		Format: `{"id":"${id}","remote_ip":"${remote_ip}",` +
			`"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
			`"status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}"` +
			`,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}`,
	}))                                 // 3
	echoHttp.Use(d.RecoverMiddleware()) // 2
	echoHttp.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowHeaders: []string{"authorization", "content-type"},
	}))                                 // 1
	// Called before routes
	echoHttp.Use(d.RawBodyPreMiddleware)         // 2
	echoHttp.Use(d.LimitOffsetSortPreMiddleware) // 1
	// init group routes
	grp := &common.Groups{
		AuthProject: echoHttp.Group(common.AuthProjectGroupPath),
		Access:      echoHttp.Group(common.AccessGroupPath),
		AuthUser:    echoHttp.Group(common.AuthUserGroupPath),
		WebHooks:    echoHttp.Group(common.WebHookGroupPath),
		Common:      echoHttp,
	}
	d.authProjectGroup(grp.AuthProject)
	d.accessGroup(grp.Access)
	d.authUserGroup(grp.AuthUser)
	d.webHookGroup(grp.WebHooks)
	// init routes
	for _, handler := range d.appSet.Handlers {
		handler.Route(grp)
	}
	return nil
}

func (d *Dispatcher) commonRoutes(echoHttp *echo.Echo) {
	echoHttp.Static("/", d.cfg.WorkDir+"/web/static")
	echoHttp.Static("/spec", d.cfg.WorkDir+"/api")
	echoHttp.GET("/docs", func(ctx echo.Context) error {
		return ctx.Render(http.StatusOK, "docs.html", map[string]interface{}{})
	})
}

func (d *Dispatcher) authProjectGroup(grp *echo.Group) {
	// Called after routes
	grp.Use(d.BodyDumpMiddleware()) // 1
}

func (d *Dispatcher) accessGroup(grp *echo.Group) {
	// Called after routes
	grp.Use(d.RecoverMiddleware()) // 1
}

func (d *Dispatcher) authUserGroup(grp *echo.Group) {
	// Called after routes
	if !d.globalCfg.DisableAuthMiddleware {
		grp.Use(
			common.ContextWrapperCallback(func(c echo.Context, next echo.HandlerFunc) error {
				handleFn := jwtMiddleware.AuthOneJwtCallableWithConfig(
					d.appSet.JwtVerifier,
					func(ui *jwtverifier.UserInfo) {
						user := common.ExtractUserContext(c)
						user.Id = ui.UserID
						user.Name = "System User"
						user.Merchants = make(map[string]bool)
						user.Roles = make(map[string]bool)
						common.SetUserContext(c, user)
					},
				)(next)
				return handleFn(c)
			}),
		) // 1
		// Called before routes
		grp.Use(d.GetUserDetailsMiddleware) // 1
	}
}

func (d *Dispatcher) webHookGroup(grp *echo.Group) {
	// Called after routes
	grp.Use(d.BodyDumpMiddleware()) // 1
}

// Config
type Config struct {
	Debug   bool `fallback:"shared.debug"`
	WorkDir string
	invoker *invoker.Invoker
}

// OnReload
func (c *Config) OnReload(callback func(ctx context.Context)) {
	c.invoker.OnReload(callback)
}

// Reload
func (c *Config) Reload(ctx context.Context) {
	c.invoker.Reload(ctx)
}

// AppSet
type AppSet struct {
	Handlers    common.Handlers
	Services    common.Services
	JwtVerifier *jwtverifier.JwtVerifier
}

// New
func New(ctx context.Context, set provider.AwareSet, appSet AppSet, cfg *Config, globalCfg *common.Config) *Dispatcher {
	set.Logger = set.Logger.WithFields(logger.Fields{"service": common.Prefix})
	return &Dispatcher{
		ctx:       ctx,
		cfg:       *cfg,
		appSet:    appSet,
		LMT:       &set,
		globalCfg: globalCfg,
	}
}
