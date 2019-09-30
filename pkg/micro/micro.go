package micro

import (
	"context"
	"github.com/ProtocolONE/go-core/invoker"
	"github.com/ProtocolONE/go-core/logger"
	"github.com/ProtocolONE/go-core/provider"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	mlog "github.com/micro/go-micro/util/log"
	"github.com/micro/go-plugins/client/selector/static"
)

// Micro
type Micro struct {
	ctx context.Context
	cfg Config
	srv micro.Service
	provider.LMT
}

// Client
func (m *Micro) Client() client.Client {
	return m.srv.Client()
}

// ListenAndServe
func (m *Micro) ListenAndServe() (err error) {

	mlog.SetLogger(NewLoggerAdapter(m.L(), logger.LevelInfo))

	m.L().Info("start listen and serve micro service")

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-m.ctx.Done()
		m.L().Info("context cancelled, shutdown is raised")
		if e := m.srv.Server().Stop(); e != nil {
			m.L().Error("graceful shutdown error, %v", logger.Args(e))
		} else {
			cancel()
		}
	}()

	if err = m.srv.Server().Start(); err != nil {
		return
	}

	<-ctx.Done()

	m.L().Info("micro service stopped successfully")
	return nil
}

// Config
type Config struct {
	Debug    bool `fallback:"shared.debug"`
	Name     string
	Version  string `default:"latest"`
	Selector string
	Bind     string
	invoker  *invoker.Invoker
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
func New(ctx context.Context, set provider.AwareSet, cfg *Config) *Micro {
	set.Logger = set.Logger.WithFields(logger.Fields{"service": Prefix, "service_name": cfg.Name})
	options := []micro.Option{
		micro.Name(cfg.Name),
		micro.Version(cfg.Version),
	}
	if cfg.Selector == "static" {
		options = append(options, micro.Selector(static.NewSelector()))
	}
	return &Micro{
		ctx: ctx,
		cfg: *cfg,
		LMT: &set,
		srv: micro.NewService(options...),
	}
}
