// +build wireinject

package daemon

import (
	"context"
	"github.com/ProtocolONE/go-core/config"
	"github.com/ProtocolONE/go-core/invoker"
	"github.com/ProtocolONE/go-core/provider"
	"github.com/google/wire"
	"github.com/paysuper/paysuper-management-api/internal/dispatcher"
	"github.com/paysuper/paysuper-management-api/internal/handlers"
	"github.com/paysuper/paysuper-management-api/internal/validators"
	"github.com/paysuper/paysuper-management-api/pkg/http"
	"github.com/paysuper/paysuper-management-api/pkg/micro"
)

// BuildHTTP
func BuildHTTP(ctx context.Context, initial config.Initial, observer invoker.Observer) (*http.HTTP, func(), error) {
	panic(
		wire.Build(
			provider.Set,
			wire.Bind(new(http.Dispatcher), new(*dispatcher.Dispatcher)),
			wire.Struct(new(provider.AwareSet), "*"),
			micro.WireSet,
			http.WireSet,
			validators.WireSet,
			dispatcher.WireSet,
			handlers.ProviderHandlers,
		),
	)
}

// BuildMicro
func BuildMicro(ctx context.Context, initial config.Initial, observer invoker.Observer) (*micro.Micro, func(), error) {
	panic(
		wire.Build(
			micro.WireSet,
			provider.Set,
			wire.Struct(new(provider.AwareSet), "*")),
	)
}
