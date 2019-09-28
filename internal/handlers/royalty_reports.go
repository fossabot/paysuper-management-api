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
	royaltyReportsPath             = "/royalty_reports"
	royaltyReportsIdPath           = "/royalty_reports/:id"
	royaltyReportsTransactionsPath = "/royalty_reports/:id/transactions"
	royaltyReportsAcceptPath       = "/royalty_reports/:id/accept"
	royaltyReportsDeclinePath      = "/royalty_reports/:id/decline"
	royaltyReportsChangePath       = "/royalty_reports/:id/change"
)

type RoyaltyReportsRoute struct {
	dispatch common.HandlerSet
	cfg      common.Config
	provider.LMT
}

func NewRoyaltyReportsRoute(set common.HandlerSet, cfg *common.Config) *RoyaltyReportsRoute {
	set.AwareSet.Logger = set.AwareSet.Logger.WithFields(logger.Fields{"router": "RoyaltyReportsRoute"})
	return &RoyaltyReportsRoute{
		dispatch: set,
		LMT:      &set.AwareSet,
		cfg:      *cfg,
	}
}

func (h *RoyaltyReportsRoute) Route(groups *common.Groups) {
	groups.AuthUser.GET(royaltyReportsPath, h.getRoyaltyReportsList)
	groups.AuthUser.GET(royaltyReportsTransactionsPath, h.getRoyaltyReport)
	groups.AuthUser.GET(royaltyReportsIdPath, h.listRoyaltyReportOrders)
	groups.AuthUser.POST(royaltyReportsAcceptPath, h.MerchantReviewRoyaltyReport)
	groups.AuthUser.POST(royaltyReportsDeclinePath, h.merchantDeclineRoyaltyReport)
	groups.AuthUser.POST(royaltyReportsChangePath, h.changeRoyaltyReport)
}

// Get royalty reports list by params (by merchant, for period) with pagination
// GET /admin/api/v1/royalty_reports
func (h *RoyaltyReportsRoute) getRoyaltyReportsList(ctx echo.Context) error {
	req := &grpc.ListRoyaltyReportsRequest{}
	err := ctx.Bind(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.NewValidationError(err.Error()))
	}

	err = h.dispatch.Validate.Struct(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.ListRoyaltyReports(ctx.Request().Context(), req)
	if err != nil {

		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "ListRoyaltyReports", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)

	}
	if res.Status != http.StatusOK {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}
	return ctx.JSON(http.StatusOK, res.Data)
}

// Get royalty reports list by id
// GET /admin/api/v1/royalty_reports
func (h *RoyaltyReportsRoute) getRoyaltyReport(ctx echo.Context) error {
	req := &grpc.GetRoyaltyReportRequest{
		ReportId: ctx.Param(common.RequestParameterId),
	}

	err := h.dispatch.Validate.Struct(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.GetRoyaltyReport(ctx.Request().Context(), req)
	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "GetRoyaltyReport", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}
	if res.Status != http.StatusOK {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}
	return ctx.JSON(http.StatusOK, res.Item)
}

// Get transactions for royalty report
// GET /admin/api/v1/royalty_reports/5ced34d689fce60bf4440829/transactions
func (h *RoyaltyReportsRoute) listRoyaltyReportOrders(ctx echo.Context) error {
	req := &grpc.ListRoyaltyReportOrdersRequest{}
	err := ctx.Bind(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.NewValidationError(err.Error()))
	}

	req.ReportId = ctx.Param(common.RequestParameterId)

	err = h.dispatch.Validate.Struct(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.ListRoyaltyReportOrders(ctx.Request().Context(), req)
	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "ListRoyaltyReportOrders", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}
	if res.Status != http.StatusOK {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}
	return ctx.JSON(http.StatusOK, res.Data)
}

// Accept royalty report by merchant
// POST /admin/api/v1/royalty_reports/5ced34d689fce60bf4440829/accept
func (h *RoyaltyReportsRoute) MerchantReviewRoyaltyReport(ctx echo.Context) error {

	req := &grpc.MerchantReviewRoyaltyReportRequest{
		IsAccepted: true,
		Ip:         ctx.RealIP(),
		ReportId:   ctx.Param(common.RequestParameterId),
	}

	res, err := h.dispatch.Services.Billing.MerchantReviewRoyaltyReport(ctx.Request().Context(), req)
	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "MerchantReviewRoyaltyReport", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}
	if res.Status != http.StatusOK {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}
	return ctx.NoContent(http.StatusNoContent)
}

// Decline royalty report by merchant and start a dispute
// POST /admin/api/v1/royalty_reports/5ced34d689fce60bf4440829/decline
func (h *RoyaltyReportsRoute) merchantDeclineRoyaltyReport(ctx echo.Context) error {

	req := &grpc.MerchantReviewRoyaltyReportRequest{
		IsAccepted: false,
		Ip:         ctx.RealIP(),
		ReportId:   ctx.Param(common.RequestParameterId),
	}

	res, err := h.dispatch.Services.Billing.MerchantReviewRoyaltyReport(ctx.Request().Context(), req)
	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "MerchantReviewRoyaltyReport", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}
	if res.Status != http.StatusOK {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}
	return ctx.NoContent(http.StatusNoContent)
}

// Change royalty report by admin
// POST /admin/api/v1/royalty_reports/5ced34d689fce60bf4440829/change
//
// @Example curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" \
//      -H "Authorization: Bearer %access_token_here%" \
//      -d '{"status": "Accepted", "correction": {"amount": 100500, "reason": "just for fun :)"}, payout_id: "5bdc39a95d1e1100019fb7df"}' \
//      https://api.paysuper.online/admin/api/v1/royalty_reports/5ced34d689fce60bf4440829/change
func (h *RoyaltyReportsRoute) changeRoyaltyReport(ctx echo.Context) error {
	req := &grpc.ChangeRoyaltyReportRequest{}
	err := ctx.Bind(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.NewValidationError(err.Error()))
	}

	req.ReportId = ctx.Param(common.RequestParameterId)
	req.Ip = ctx.RealIP()

	err = h.dispatch.Validate.Struct(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.ChangeRoyaltyReport(ctx.Request().Context(), req)
	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "ChangeRoyaltyReport", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}
	if res.Status != http.StatusOK {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}
	return ctx.NoContent(http.StatusNoContent)
}
