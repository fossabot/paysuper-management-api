package handlers

import (
	"github.com/Nerufa/go-shared/logger"
	"github.com/Nerufa/go-shared/provider"
	"github.com/labstack/echo/v4"
	"github.com/paysuper/paysuper-billing-server/pkg"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/grpc"
	"github.com/paysuper/paysuper-management-api/internal/dispatcher/common"
	"net/http"
)

const (
	dashboardMainPath            = "/merchants/:id/dashboard/main"
	dashboardRevenueDynamicsPath = "/merchants/:id/dashboard/revenue_dynamics"
	dashboardBasePath            = "/merchants/:id/dashboard/base"
)

type DashboardRoute struct {
	dispatch common.HandlerSet
	cfg      common.Config
	provider.LMT
}

func NewDashboardRoute(set common.HandlerSet, cfg *common.Config) *DashboardRoute {
	set.AwareSet.Logger = set.AwareSet.Logger.WithFields(logger.Fields{"router": "DashboardRoute"})
	return &DashboardRoute{
		dispatch: set,
		LMT:      &set.AwareSet,
		cfg:      *cfg,
	}
}

func (h *DashboardRoute) Route(groups *common.Groups) {
	groups.AuthUser.GET(dashboardMainPath, h.getMainReports)
	groups.AuthUser.GET(dashboardRevenueDynamicsPath, h.getRevenueDynamicsReport)
	groups.AuthUser.GET(dashboardBasePath, h.getBaseReports)
}

// @Description get main reports data for dashboard
// @Example curl -X GET -H 'Authorization: Bearer %access_token_here%' \
//  https://api.paysuper.online/admin/api/v1/merchants/ffffffffffffffffffffffff/dashboard/main?period=previous_month
func (h *DashboardRoute) getMainReports(ctx echo.Context) error {
	req := &grpc.GetDashboardMainRequest{}
	err := ctx.Bind(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestParamsIncorrect)
	}

	req.MerchantId = ctx.Param(common.RequestParameterId)
	err = h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.GetDashboardMainReport(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "GetDashboardMainReport", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}

	if res.Status != pkg.ResponseStatusOk {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.JSON(http.StatusOK, res.Item)
}

// @Description get revenue dynamics report data for dashboard
// @Example curl -X GET -H 'Authorization: Bearer %access_token_here%' \
//  'http://127.0.0.1:3001/admin/api/v1/merchants/ffffffffffffffffffffffff/dashboard/revenue_dynamics?period=previous_month'
func (h *DashboardRoute) getRevenueDynamicsReport(ctx echo.Context) error {
	req := &grpc.GetDashboardMainRequest{}
	err := ctx.Bind(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestParamsIncorrect)
	}

	req.MerchantId = ctx.Param(common.RequestParameterId)
	err = h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.GetDashboardRevenueDynamicsReport(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "GetDashboardMainReport", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}

	if res.Status != pkg.ResponseStatusOk {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.JSON(http.StatusOK, res.Item)
}

// @Description get base reports data for dashboard
// @Example curl -X GET -H 'Authorization: Bearer %access_token_here%' \
//  https://api.paysuper.online/admin/api/v1/merchants/ffffffffffffffffffffffff/dashboard/base?period=previous_month
func (h *DashboardRoute) getBaseReports(ctx echo.Context) error {
	req := &grpc.GetDashboardBaseReportRequest{}
	err := ctx.Bind(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestParamsIncorrect)
	}

	req.MerchantId = ctx.Param(common.RequestParameterId)
	err = h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.GetDashboardBaseReport(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "GetDashboardMainReport", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}

	if res.Status != pkg.ResponseStatusOk {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.JSON(http.StatusOK, res.Item)
}
