package handlers

import (
	"github.com/ProtocolONE/go-core/logger"
	"github.com/ProtocolONE/go-core/provider"
	"github.com/labstack/echo/v4"
	"github.com/paysuper/paysuper-billing-server/pkg"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/grpc"
	"github.com/paysuper/paysuper-management-api/internal/dispatcher/common"
	"net/http"
)

const (
	payoutsPath   = "/payout_documents"
	payoutsIdPath = "/payout_documents/:id"
)

type PayoutDocumentsRoute struct {
	dispatch common.HandlerSet
	cfg      common.Config
	provider.LMT
}

// NewPayoutDocumentsRoute
func NewPayoutDocumentsRoute(set common.HandlerSet, cfg *common.Config) *PayoutDocumentsRoute {
	set.AwareSet.Logger = set.AwareSet.Logger.WithFields(logger.Fields{"router": "PayoutDocumentsRoute"})
	return &PayoutDocumentsRoute{
		dispatch: set,
		LMT:      &set.AwareSet,
		cfg:      *cfg,
	}
}

func (h *PayoutDocumentsRoute) Route(groups *common.Groups) {
	groups.AuthUser.GET(payoutsPath, h.getPayoutDocumentsList)
	groups.AuthUser.GET(payoutsIdPath, h.getPayoutDocument)
	groups.AuthUser.POST(payoutsPath, h.createPayoutDocument)
	groups.AuthUser.POST(payoutsIdPath, h.updatePayoutDocument)
}

// Get payout documents list with filters and pagination
// GET /admin/api/v1/payout_documents?payout_document_id=5ced34d689fce60bf4440829
// GET /admin/api/v1/payout_documents?status=pending&merchant_id=5bdc39a95d1e1100019fb7df&limit=10&offset=0
// GET /admin/api/v1/payout_documents?status=pending&limit=10&offset=0
func (h *PayoutDocumentsRoute) getPayoutDocumentsList(ctx echo.Context) error {
	req := &grpc.GetPayoutDocumentsRequest{}
	err := ctx.Bind(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestParamsIncorrect)
	}

	err = h.dispatch.Validate.Struct(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.GetPayoutDocuments(ctx.Request().Context(), req)
	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "GetPayoutDocuments", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}
	if res.Status != http.StatusOK {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}
	return ctx.JSON(http.StatusOK, res.Data)
}

// Get payout document
// GET /admin/api/v1/payout_documents/5ced34d689fce60bf4440829
func (h *PayoutDocumentsRoute) getPayoutDocument(ctx echo.Context) error {
	req := &grpc.GetPayoutDocumentsRequest{}
	req.PayoutDocumentId = ctx.Param(common.RequestParameterId)

	err := h.dispatch.Validate.Struct(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.GetPayoutDocuments(ctx.Request().Context(), req)
	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "GetPayoutDocuments", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}
	if res.Status != http.StatusOK {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}
	if len(res.Data.Items) == 0 {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	return ctx.JSON(http.StatusOK, res.Data.Items[0])
}

// Create payout document
// POST /admin/api/v1/payout_documents
//
// @Example curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" \
//      -H "Authorization: Bearer %access_token_here%" \
//      -d '{ "description": "royalty for june-july 2019", "merchant_id": "5bdc39a95d1e1100019fb7df" }' \
//      https://api.paysuper.online/admin/api/v1/payout_documents
func (h *PayoutDocumentsRoute) createPayoutDocument(ctx echo.Context) error {
	req := &grpc.CreatePayoutDocumentRequest{}
	err := ctx.Bind(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestParamsIncorrect)
	}

	req.Ip = ctx.RealIP()

	err = h.dispatch.Validate.Struct(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.CreatePayoutDocument(ctx.Request().Context(), req)
	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "CreatePayoutDocument", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}
	if res.Status != http.StatusOK {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}
	return ctx.JSON(http.StatusOK, res.Item)
}

// Update payout document by admin
// POST /admin/api/v1/payout_documents/5ced34d689fce60bf4440829
//
// @Example curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" \
//      -H "Authorization: Bearer %access_token_here%" \
//      -d '{"status": "failed", "failure_code": "account_closed"}' \
//      https://api.paysuper.online/admin/api/v1/payout_documents/5ced34d689fce60bf4440829
func (h *PayoutDocumentsRoute) updatePayoutDocument(ctx echo.Context) error {

	req := &grpc.UpdatePayoutDocumentRequest{}
	err := ctx.Bind(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestParamsIncorrect)
	}
	req.PayoutDocumentId = ctx.Param(common.RequestParameterId)
	req.Ip = ctx.RealIP()

	err = h.dispatch.Validate.Struct(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.UpdatePayoutDocument(ctx.Request().Context(), req)
	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "UpdatePayoutDocument", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}
	if res.Status != http.StatusOK {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}
	return ctx.JSON(http.StatusOK, res.Item)
}
