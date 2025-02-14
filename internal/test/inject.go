// +build wireinject

package test

import (
	"context"
	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/ProtocolONE/go-core/v2/pkg/metric"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/ProtocolONE/go-core/v2/pkg/tracing"
	"github.com/google/wire"
	"github.com/paysuper/paysuper-management-api/internal/dispatcher"
	"github.com/paysuper/paysuper-management-api/internal/dispatcher/common"
	"github.com/paysuper/paysuper-management-api/internal/validators"
	"gopkg.in/go-playground/validator.v9"
	"os"
)

type TestSet struct {
	AwareSet     provider.AwareSet
	Configurator config.Configurator
	GlobalConfig *common.Config
	HandlerSet   common.HandlerSet
	Initial      config.Initial
}

// ProviderTestSet
func ProviderTestSet(initial config.Initial, awareSet provider.AwareSet, srv common.Services, configurator config.Configurator, globalConfig *common.Config, validate *validator.Validate) (*TestSet, func(), error) {
	t := &TestSet{
		AwareSet:     awareSet,
		Configurator: configurator,
		GlobalConfig: globalConfig,
		HandlerSet: common.HandlerSet{
			AwareSet: awareSet,
			Validate: validate,
			Services: srv,
		},
		Initial: initial,
	}
	return t, func() {}, nil
}

// ProviderTestInitial
func ProviderTestInitial() config.Initial {
	wd := os.Getenv("WD")
	if len(wd) == 0 {
		wd = "./../../"
	}
	return config.Initial{WorkDir: wd}
}

// BuildTestSet
func BuildTestSet(ctx context.Context, settings config.Settings, srv common.Services, observer invoker.Observer) (*TestSet, func(), error) {
	panic(
		wire.Build(
			ProviderTestInitial,
			ProviderTestSet,
			config.WireTestSet, // Configurator | Dependencies: config.Settings
			logger.WireTestSet, // Logger
			metric.WireTestSet, // Scope
			tracing.WireTestSet,
			wire.Struct(new(provider.AwareSet), "*"),
			validators.WireSet,
			dispatcher.ProviderGlobalCfg,
			dispatcher.ProviderValidators,
		),
	)
}

// BuildDispatcher
func BuildDispatcher(ctx context.Context, settings config.Settings, srv common.Services, handlers common.Handlers, observer invoker.Observer) (*dispatcher.Dispatcher, func(), error) {
	panic(
		wire.Build(
			ProviderTestInitial,
			config.WireTestSet, // Configurator | Dependencies: config.Settings
			logger.WireTestSet, // Logger
			metric.WireTestSet, // Scope
			tracing.WireTestSet,
			wire.Struct(new(provider.AwareSet), "*"),
			dispatcher.WireTestSet, // Dispatcher | Dependencies: AwareSet, ValidatorSet, Services, Handlers, Configurator
		),
	)
}
