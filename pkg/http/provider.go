package http

import (
	"context"
	"github.com/ProtocolONE/go-core/config"
	"github.com/ProtocolONE/go-core/invoker"
	"github.com/ProtocolONE/go-core/provider"
	"github.com/google/wire"
)

// Cfg
func Cfg(cfg config.Configurator) (*Config, func(), error) {
	c := &Config{
		invoker: invoker.NewInvoker(),
	}
	e := cfg.UnmarshalKeyOnReload(UnmarshalKey, c)
	return c, func() {}, e
}

// CfgTest
func CfgTest() (*Config, func(), error) {
	c := &Config{
		invoker: invoker.NewInvoker(),
	}
	return c, func() {}, nil
}

// Provider
func Provider(ctx context.Context, set provider.AwareSet, dispatcher Dispatcher, cfg *Config) (*HTTP, func(), error) {
	http := New(ctx, set, dispatcher, cfg)
	return http, func() {}, nil
}

var (
	WireSet = wire.NewSet(
		Provider,
		Cfg,
	)
	WireTestSet = wire.NewSet(
		Provider,
		CfgTest,
	)
)
