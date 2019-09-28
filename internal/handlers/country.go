package handlers

import (
	"github.com/Nerufa/go-shared/logger"
	"github.com/Nerufa/go-shared/provider"
	"github.com/labstack/echo/v4"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/billing"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/grpc"
	"github.com/paysuper/paysuper-management-api/internal/dispatcher/common"
	"net/http"
)

type CountryApiV1 struct {
	dispatch common.HandlerSet
	cfg      common.Config
	provider.LMT
}

func NewCountryApiV1(set common.HandlerSet, cfg *common.Config) *CountryApiV1 {
	set.AwareSet.Logger = set.AwareSet.Logger.WithFields(logger.Fields{"router": "CountryApiV1"})
	return &CountryApiV1{
		dispatch: set,
		LMT:      &set.AwareSet,
		cfg:      *cfg,
	}
}

func (h *CountryApiV1) Route(groups *common.Groups) {
	groups.AuthProject.GET("/country", h.get)
	groups.AuthProject.GET("/country/:code", h.getById)
}

// Get full list of currencies
// GET /api/v1/country
func (h *CountryApiV1) get(ctx echo.Context) error {

	res, err := h.dispatch.Services.Billing.GetCountriesList(ctx.Request().Context(), &grpc.EmptyRequest{})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError /*ErrorCountriesListError*/, err)
	}

	return ctx.JSON(http.StatusOK, res)
}

// Get country by ISO 3166-1 alpha 2 country code
// GET /api/v1/country/{code}
func (h *CountryApiV1) getById(ctx echo.Context) error {
	code := ctx.Param("code")

	if len(code) != 2 {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorIncorrectCountryIdentifier)
	}

	req := &billing.GetCountryRequest{
		IsoCode: code,
	}
	err := h.dispatch.Validate.Struct(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.GetCountry(ctx.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, common.ErrorCountryNotFound)
	}

	return ctx.JSON(http.StatusOK, res)
}
