package http

import (
	"context"
	"github.com/Nerufa/go-shared/invoker"
	"github.com/Nerufa/go-shared/logger"
	"github.com/Nerufa/go-shared/provider"
	"github.com/labstack/echo/v4"
	"net/http"
)

// HTTP
type HTTP struct {
	ctx        context.Context
	cfg        Config
	dispatcher Dispatcher
	provider.LMT
}

// ListenAndServe
func (h *HTTP) ListenAndServe() (err error) {

	server := echo.New()
	server.HideBanner = true
	server.HidePort = true
	server.Debug = h.cfg.Debug

	h.dispatcher.Dispatch(server)

	h.L().Info("start listen and serve http at %v", logger.Args(h.cfg.Bind))

	go func() {
		<-h.ctx.Done()
		h.L().Info("context cancelled, shutdown is raised")
		if e := server.Shutdown(context.Background()); e != nil {
			h.L().Error("graceful shutdown error, %v", logger.Args(e))
		}
	}()

	if err = server.Start(h.cfg.Bind); err != nil {
		if err == http.ErrServerClosed {
			err = nil
			return
		}
	}

	h.L().Info("http server stopped successfully")
	return nil
}

// Config
type Config struct {
	Debug   bool   `fallback:"shared.debug"`
	Bind    string `required:"true"`
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

// New
func New(ctx context.Context, set provider.AwareSet, dispatcher Dispatcher, cfg *Config) *HTTP {
	set.Logger = set.Logger.WithFields(logger.Fields{"service": Prefix})
	return &HTTP{
		ctx:        ctx,
		cfg:        *cfg,
		dispatcher: dispatcher,
		LMT:        &set,
	}
}
