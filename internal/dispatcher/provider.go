package dispatcher

import (
	"context"
	"github.com/Nerufa/go-shared/config"
	"github.com/Nerufa/go-shared/invoker"
	"github.com/Nerufa/go-shared/provider"
	jwtverifier "github.com/ProtocolONE/authone-jwt-verifier-golang"
	geoip "github.com/ProtocolONE/geoip-service/pkg"
	"github.com/ProtocolONE/geoip-service/pkg/proto"
	"github.com/google/wire"
	"github.com/paysuper/paysuper-billing-server/pkg"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/billing"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/grpc"
	"github.com/paysuper/paysuper-management-api/internal/dispatcher/common"
	"github.com/paysuper/paysuper-management-api/internal/validators"
	"github.com/paysuper/paysuper-management-api/pkg/micro"
	paylinkServiceConst "github.com/paysuper/paysuper-payment-link/pkg"
	paylink "github.com/paysuper/paysuper-payment-link/proto"
	"github.com/paysuper/paysuper-recurring-repository/pkg/constant"
	"github.com/paysuper/paysuper-recurring-repository/pkg/proto/repository"
	reporterPkg "github.com/paysuper/paysuper-reporter/pkg"
	reporterProto "github.com/paysuper/paysuper-reporter/pkg/proto"
	taxServiceConst "github.com/paysuper/paysuper-tax-service/pkg"
	tax_service "github.com/paysuper/paysuper-tax-service/proto"
	"gopkg.in/go-playground/validator.v9"
)

// ProviderCfg
func ProviderCfg(cfg config.Configurator) (*Config, func(), error) {
	c := &Config{
		WorkDir: cfg.WorkDir(),
		invoker: invoker.NewInvoker(),
	}
	e := cfg.UnmarshalKeyOnReload(common.UnmarshalKey, c)
	return c, func() {}, e
}

// ProviderGlobalCfg
func ProviderGlobalCfg(cfg config.Configurator) (*common.Config, func(), error) {
	c := &common.Config{
	}
	e := cfg.UnmarshalKey(common.UnmarshalGlobalConfigKey, c)
	return c, func() {}, e
}

// ProviderJwtVerifier
func ProviderJwtVerifier(cfg *common.Config) *jwtverifier.JwtVerifier {
	return jwtverifier.NewJwtVerifier(jwtverifier.Config{
		ClientID:     cfg.Auth1.ClientId,
		ClientSecret: cfg.Auth1.ClientSecret,
		Scopes:       []string{"openid", "offline"},
		RedirectURL:  cfg.Auth1.RedirectUrl,
		Issuer:       cfg.Auth1.Issuer,
	})
}

// ProviderServices
func ProviderServices(srv *micro.Micro) common.Services {
	return common.Services{
		Repository: repository.NewRepositoryService(constant.PayOneRepositoryServiceName, srv.Client()),
		Geo:        proto.NewGeoIpService(geoip.ServiceName, srv.Client()),
		Billing:    grpc.NewBillingService(pkg.ServiceName, srv.Client()),
		Tax:        tax_service.NewTaxService(taxServiceConst.ServiceName, srv.Client()),
		PayLink:    paylink.NewPaylinkService(paylinkServiceConst.ServiceName, srv.Client()),
		Reporter:   reporterProto.NewReporterService(reporterPkg.ServiceName, srv.Client()),
	}
}

// ProviderValidators
func ProviderValidators(v *validators.ValidatorSet) (validate *validator.Validate, _ func(), err error) {
	validate = validator.New()
	if err = validate.RegisterValidation("phone", v.PhoneValidator); err != nil {
		return
	}
	if err = validate.RegisterValidation("uuid", v.UuidValidator); err != nil {
		return
	}
	if err = validate.RegisterValidation("zip_usa", v.ZipUsaValidator); err != nil {
		return
	}
	if err = validate.RegisterValidation("name", v.NameValidator); err != nil {
		return
	}
	if err = validate.RegisterValidation("position", v.PositionValidator); err != nil {
		return
	}
	validate.RegisterStructValidation(v.CompanyValidator, grpc.UserProfileCompany{})
	validate.RegisterStructValidation(v.MerchantCompanyValidator, billing.MerchantCompanyInfo{})
	validate.RegisterStructValidation(v.MerchantTariffRatesValidator, grpc.GetMerchantTariffRatesRequest{})
	if err = validate.RegisterValidation("company_name", v.CompanyNameValidator); err != nil {
		return
	}
	if err = validate.RegisterValidation("swift", v.SwiftValidator); err != nil {
		return
	}
	if err = validate.RegisterValidation("city", v.CityValidator); err != nil {
		return
	}
	if err = validate.RegisterValidation("world_region", v.WorldRegionValidator); err != nil {
		return
	}
	return validate, func() {}, nil
}

// ProviderDispatcher
func ProviderDispatcher(ctx context.Context, set provider.AwareSet, appSet AppSet, cfg *Config, globalCfg *common.Config) (*Dispatcher, func(), error) {
	d := New(ctx, set, appSet, cfg, globalCfg)
	return d, func() {}, nil
}

var (
	// Dependencies: go-shared/provider.AwareSet, internal/*validators.ValidatorSet, pkg/micro.Micro, go-shared/config.Configurator, ProviderHandlers
	WireSet = wire.NewSet(
		ProviderDispatcher,
		ProviderServices,
		ProviderJwtVerifier,
		ProviderValidators,
		ProviderCfg,
		ProviderGlobalCfg,
		wire.Struct(new(AppSet), "*"),
	)
	// Dependencies: go-shared/provider.AwareSet, internal/*validators.ValidatorSet, common.Services, common.Handlers, go-shared/config.Configurator
	WireTestSet = wire.NewSet(
		ProviderDispatcher,
		ProviderJwtVerifier,
		ProviderValidators,
		ProviderCfg,
		ProviderGlobalCfg,
		wire.Struct(new(AppSet), "*"),
	)
)
