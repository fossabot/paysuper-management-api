package handlers

import (
	"github.com/Nerufa/go-shared/logger"
	"github.com/Nerufa/go-shared/provider"
	"github.com/labstack/echo/v4"
	"github.com/paysuper/paysuper-management-api/internal/dispatcher/common"
	"github.com/paysuper/paysuper-tax-service/proto"
	"net/http"
	"strconv"
)

const (
	taxesPath   = "/taxes"
	taxesIDPath = "/taxes/:id"
)

type TaxesRoute struct {
	dispatch common.HandlerSet
	cfg      common.Config
	provider.LMT
}

func NewTaxesRoute(set common.HandlerSet, cfg *common.Config) *TaxesRoute {
	set.AwareSet.Logger = set.AwareSet.Logger.WithFields(logger.Fields{"router": "TaxesRoute"})
	return &TaxesRoute{
		dispatch: set,
		LMT:      &set.AwareSet,
		cfg:      *cfg,
	}
}

func (h *TaxesRoute) Route(groups *common.Groups) {
	groups.AuthUser.GET(taxesPath, h.getTaxes)
	groups.AuthUser.POST(taxesPath, h.setTax)
	groups.AuthUser.DELETE(taxesIDPath, h.deleteTax)
}

func (h *TaxesRoute) getTaxes(ctx echo.Context) error {
	req := h.bindGetTaxes(ctx)
	res, err := h.dispatch.Services.Tax.GetRates(ctx.Request().Context(), req)

	if err != nil {
		h.L().Error(common.InternalErrorTemplate, logger.WithFields(logger.Fields{"err": err.Error()}))
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorInternal)
	}

	return ctx.JSON(http.StatusOK, res.Rates)
}

func (h *TaxesRoute) bindGetTaxes(ctx echo.Context) *tax_service.GetRatesRequest {
	structure := &tax_service.GetRatesRequest{}

	params := ctx.QueryParams()

	if v, ok := params["country"]; ok {
		structure.Country = string(v[0])
	}

	if v, ok := params["city"]; ok {
		structure.City = string(v[0])
	}

	if v, ok := params["state"]; ok {
		structure.State = string(v[0])
	}

	if v, ok := params["zip"]; ok {
		structure.Zip = string(v[0])
	}

	if v, ok := params[common.RequestParameterLimit]; ok {
		if i, err := strconv.ParseInt(v[0], 10, 32); err == nil {
			structure.Limit = int32(i)
		}
	} else {
		structure.Limit = h.cfg.LimitDefault
	}

	if v, ok := params[common.RequestParameterOffset]; ok {
		if i, err := strconv.ParseInt(v[0], 10, 32); err == nil {
			structure.Offset = int32(i)
		}
	} else {
		structure.Offset = h.cfg.OffsetDefault
	}

	return structure
}

func (h *TaxesRoute) setTax(ctx echo.Context) error {
	req := &tax_service.TaxRate{}
	err := ctx.Bind(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.NewValidationError(err.Error()))
	}

	res, err := h.dispatch.Services.Tax.CreateOrUpdate(ctx.Request().Context(), req)
	if err != nil {
		h.L().Error(common.InternalErrorTemplate, logger.WithFields(logger.Fields{"err": err.Error()}))
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorInternal)
	}

	return ctx.JSON(http.StatusOK, res)
}

func (h *TaxesRoute) deleteTax(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestParamsIncorrect)
	}

	value, err := strconv.Atoi(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestParamsIncorrect)
	}

	res, err := h.dispatch.Services.Tax.DeleteRateById(ctx.Request().Context(), &tax_service.DeleteRateRequest{Id: uint32(value)})
	if err != nil {
		h.L().Error(common.InternalErrorTemplate, logger.WithFields(logger.Fields{"err": err.Error()}))
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorInternal)
	}

	return ctx.JSON(http.StatusOK, res)
}
