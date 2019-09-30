package handlers

import (
	"github.com/ProtocolONE/go-core/logger"
	"github.com/ProtocolONE/go-core/provider"
	"github.com/labstack/echo/v4"
	"github.com/paysuper/paysuper-billing-server/pkg"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/billing"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/grpc"
	"github.com/paysuper/paysuper-management-api/internal/dispatcher/common"
	"net/http"
)

type PaymentCostRoute struct {
	dispatch common.HandlerSet
	cfg      common.Config
	provider.LMT
}

func NewPaymentCostRoute(set common.HandlerSet, cfg *common.Config) *PaymentCostRoute {
	set.AwareSet.Logger = set.AwareSet.Logger.WithFields(logger.Fields{"router": "PaymentCostRoute"})
	return &PaymentCostRoute{
		dispatch: set,
		LMT:      &set.AwareSet,
		cfg:      *cfg,
	}
}

const (
	paymentCostsChannelSystemPath        = "/payment_costs/channel/system"
	paymentCostsChannelSystemAllPath     = "/payment_costs/channel/system/all"
	paymentCostsChannelMerchantPath      = "/payment_costs/channel/merchant/:id"
	paymentCostsChannelMerchantAllPath   = "/payment_costs/channel/merchant/:id/all"
	paymentCostsChannelSystemIdPath      = "/payment_costs/channel/system/:id"
	paymentCostsChannelMerchantIdsPath   = "/payment_costs/channel/merchant/:merchant_id/:rate_id"
	paymentCostsMoneyBackAllPath         = "/payment_costs/money_back/system/all"
	paymentCostsMoneyBackMerchantPath    = "/payment_costs/money_back/merchant/:id"
	paymentCostsMoneyBackMerchantAllPath = "/payment_costs/money_back/merchant/:id/all"
	paymentCostsMoneyBackSystemPath      = "/payment_costs/money_back/system"
	paymentCostsMoneyBackSystemIdPath    = "/payment_costs/money_back/system/:id"
)

func (h *PaymentCostRoute) Route(groups *common.Groups) {
	groups.AuthUser.GET(paymentCostsChannelSystemAllPath, h.getAllPaymentChannelCostSystem)
	groups.AuthUser.GET(paymentCostsChannelMerchantAllPath, h.getAllPaymentChannelCostMerchant) //надо править
	groups.AuthUser.GET(paymentCostsMoneyBackAllPath, h.getAllMoneyBackCostSystem)
	groups.AuthUser.GET(paymentCostsMoneyBackMerchantAllPath, h.getAllMoneyBackCostMerchant) //надо править

	groups.AuthUser.GET(paymentCostsChannelSystemPath, h.getPaymentChannelCostSystem)
	groups.AuthUser.GET(paymentCostsChannelMerchantPath, h.getPaymentChannelCostMerchant)
	groups.AuthUser.GET(paymentCostsMoneyBackSystemPath, h.getMoneyBackCostSystem)
	groups.AuthUser.GET(paymentCostsMoneyBackMerchantPath, h.getMoneyBackCostMerchant)

	groups.AuthUser.DELETE(paymentCostsChannelSystemIdPath, h.deletePaymentChannelCostSystem)
	groups.AuthUser.DELETE(paymentCostsChannelMerchantPath, h.deletePaymentChannelCostMerchant)
	groups.AuthUser.DELETE(paymentCostsMoneyBackSystemIdPath, h.deleteMoneyBackCostSystem)
	groups.AuthUser.DELETE(paymentCostsMoneyBackMerchantPath, h.deleteMoneyBackCostMerchant)

	groups.AuthUser.POST(paymentCostsChannelSystemPath, h.setPaymentChannelCostSystem)
	groups.AuthUser.POST(paymentCostsChannelMerchantPath, h.setPaymentChannelCostMerchant)
	groups.AuthUser.POST(paymentCostsMoneyBackSystemPath, h.setMoneyBackCostSystem)
	groups.AuthUser.POST(paymentCostsMoneyBackMerchantPath, h.setMoneyBackCostMerchant)

	groups.AuthUser.PUT(paymentCostsChannelSystemIdPath, h.setPaymentChannelCostSystem)
	groups.AuthUser.PUT(paymentCostsChannelMerchantIdsPath, h.setPaymentChannelCostMerchant)
	groups.AuthUser.PUT(paymentCostsMoneyBackSystemIdPath, h.setMoneyBackCostSystem)
	groups.AuthUser.PUT(paymentCostsChannelMerchantIdsPath, h.setMoneyBackCostMerchant)
}

// @Description Get system costs for payments operations
// @Example curl -X GET -H "Authorization: Bearer %access_token_here%"  \
// 		https://api.paysuper.online/admin/api/v1/payment_costs/channel/system?name=VISA&region=CIS&country=AZ
func (h *PaymentCostRoute) getPaymentChannelCostSystem(ctx echo.Context) error {
	req := &billing.PaymentChannelCostSystemRequest{}
	err := ctx.Bind(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestDataInvalid)
	}

	err = h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.GetPaymentChannelCostSystem(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "GetPaymentChannelCostSystem", req)

		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorInternal)
	}

	if res.Status != http.StatusOK {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.JSON(http.StatusOK, res.Item)
}

// @Description Get merchant costs for payment operations
// @Example curl -X GET -H "Authorization: Bearer %access_token_here%"  \
// 		https://api.paysuper.online/admin/api/v1/payment_costs/channel/merchant/ffffffffffffffffffffffff?name=VISA&region=CIS&country=AZ&payout_currency=USD&amount=100
func (h *PaymentCostRoute) getPaymentChannelCostMerchant(ctx echo.Context) error {
	req := &billing.PaymentChannelCostMerchantRequest{}
	err := ctx.Bind(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestDataInvalid)
	}

	req.MerchantId = ctx.Param(common.RequestParameterId)
	err = h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.GetPaymentChannelCostMerchant(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "GetPaymentChannelCostMerchant", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorInternal)
	}

	if res.Status != http.StatusOK {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.JSON(http.StatusOK, res.Item)
}

// @Description Get system costs for money back operations
// @Example curl -X GET -H "Authorization: Bearer %access_token_here%"  \
// 		https://api.paysuper.online/admin/api/v1/payment_costs/money_back/system?name=VISA&region=CIS&country=AZ&payout_currency=USD&days=10&undo_reason=chargeback&payment_stage=1
func (h *PaymentCostRoute) getMoneyBackCostSystem(ctx echo.Context) error {
	req := &billing.MoneyBackCostSystemRequest{}
	err := ctx.Bind(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestDataInvalid)
	}

	err = h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.GetMoneyBackCostSystem(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "GetMoneyBackCostSystem", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorInternal)
	}

	if res.Status != http.StatusOK {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.JSON(http.StatusOK, res.Item)
}

// @Description Get merchant costs for money back operations
// @Example curl -X GET -H "Authorization: Bearer %access_token_here%"  \
// 		https://api.paysuper.online/admin/api/v1/payment_costs/money_back/merchant/ffffffffffffffffffffffff??name=VISA&region=CIS&country=AZ&payout_currency=USD&days=10&undo_reason=chargeback&payment_stage=1
func (h *PaymentCostRoute) getMoneyBackCostMerchant(ctx echo.Context) error {
	req := &billing.MoneyBackCostMerchantRequest{}
	err := ctx.Bind(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestDataInvalid)
	}

	req.MerchantId = ctx.Param(common.RequestParameterId)
	err = h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.GetMoneyBackCostMerchant(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "GetMoneyBackCostMerchant", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorInternal)
	}

	if res.Status != http.StatusOK {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.JSON(http.StatusOK, res.Item)
}

// @Description Delete system costs for payment operations
// @Example curl -X DELETE -H "Authorization: Bearer %access_token_here%" -H "Content-Type: application/json" \
// 		https://api.paysuper.online/admin/api/v1/payment_costs/channel/system/ffffffffffffffffffffffff
func (h *PaymentCostRoute) deletePaymentChannelCostSystem(ctx echo.Context) error {
	req := &billing.PaymentCostDeleteRequest{Id: ctx.Param(common.RequestParameterId)}
	err := h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.DeletePaymentChannelCostSystem(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "DeletePaymentChannelCostSystem", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorInternal)
	}

	if res.Status != http.StatusOK {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.NoContent(http.StatusNoContent)
}

// @Description Delete merchant costs for payment operations
// @Example curl -X DELETE -H "Authorization: Bearer %access_token_here%" -H "Content-Type: application/json" \
// 		https://api.paysuper.online/admin/api/v1/payment_costs/channel/merchant/ffffffffffffffffffffffff
func (h *PaymentCostRoute) deletePaymentChannelCostMerchant(ctx echo.Context) error {
	req := &billing.PaymentCostDeleteRequest{Id: ctx.Param(common.RequestParameterId)}
	err := h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.DeletePaymentChannelCostMerchant(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "DeletePaymentChannelCostMerchant", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorInternal)
	}

	if res.Status != http.StatusOK {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.NoContent(http.StatusNoContent)
}

// @Description Delete system costs for money back operations
// @Example curl -X DELETE -H "Authorization: Bearer %access_token_here%" -H "Content-Type: application/json" \
// 		https://api.paysuper.online/admin/api/v1/payment_costs/money_back/system/ffffffffffffffffffffffff
func (h *PaymentCostRoute) deleteMoneyBackCostSystem(ctx echo.Context) error {
	req := &billing.PaymentCostDeleteRequest{Id: ctx.Param(common.RequestParameterId)}
	err := h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.DeleteMoneyBackCostSystem(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "DeleteMoneyBackCostSystem", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorInternal)
	}

	if res.Status != http.StatusOK {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.NoContent(http.StatusNoContent)
}

// @Description Delete merchant costs for money back operations
// @Example curl -X DELETE -H "Authorization: Bearer %access_token_here%" -H "Content-Type: application/json" \
// 		https://api.paysuper.online/admin/api/v1/payment_costs/money_back/merchant/ffffffffffffffffffffffff
func (h *PaymentCostRoute) deleteMoneyBackCostMerchant(ctx echo.Context) error {
	req := &billing.PaymentCostDeleteRequest{Id: ctx.Param(common.RequestParameterId)}
	err := h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.DeleteMoneyBackCostMerchant(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "DeleteMoneyBackCostMerchant", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorInternal)
	}

	if res.Status != http.StatusOK {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.NoContent(http.StatusNoContent)
}

// @Description Create and update system costs for payments operations
// @Example curl -X POST -H "Authorization: Bearer %access_token_here%" -H "Content-Type: application/json" \
//      -d '{"name": "VISA", "region": "CIS", "country": "AZ", "percent": 0.01, "fix_amount": 2.34,
//      "fix_amount_currency": "USD"}' \
//      https://api.paysuper.online/admin/api/v1/payment_costs/channel/system
//
// @Example curl -X PUT -H "Authorization: Bearer %access_token_here%" -H "Content-Type: application/json" \
//      -d '{"name": "VISA", "region": "CIS", "country": "AZ", "percent": 0.01, "fix_amount": 2.34,
//      "fix_amount_currency": "USD"}' \
//      https://api.paysuper.online/admin/api/v1/payment_costs/channel/system/ffffffffffffffffffffffff
func (h *PaymentCostRoute) setPaymentChannelCostSystem(ctx echo.Context) error {
	req := &billing.PaymentChannelCostSystem{}
	err := ctx.Bind(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestDataInvalid)
	}

	if pcId := ctx.Param(common.RequestParameterId); pcId != "" {
		req.Id = pcId
	}

	err = h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.SetPaymentChannelCostSystem(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "SetPaymentChannelCostSystem", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorInternal)
	}

	if res.Status != http.StatusOK {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.JSON(http.StatusOK, res.Item)
}

// @Description Create and update merchant costs for payments operations
//  @Example curl -X POST -H "Authorization: Bearer %access_token_here%" -H "Content-Type: application/json" \
//      -d '{"name": "VISA", "region": "CIS", "country": "AZ", "min_amount": 0.75, "method_percent": 0.01,
// 			"method_fix_amount": 2.34, "ps_percent": 0.05, "ps_fixed_fee": 2, "ps_fixed_fee_currency": "EUR",
// 			"payout_currency": "USD"}' \
//      https://api.paysuper.online/admin/api/v1/payment_costs/channel/merchant/ffffffffffffffffffffffff
//
// @Example curl -X PUT -H "Authorization: Bearer %access_token_here%" -H "Content-Type: application/json" \
//      -d '{"name": "VISA", "region": "CIS", "country": "AZ", "min_amount": 0.75, "method_percent": 0.01,
//      "method_fix_amount": 2.34, "ps_percent": 0.05, "ps_fixed_fee": 2, "ps_fixed_fee_currency": "EUR",
//      "payout_currency": "USD"}' \
//      https://api.paysuper.online/admin/api/v1/payment_costs/channel/merchant/ffffffffffffffffffffffff/aaaaaaaaaaaaaaaaaaaaaaaa
func (h *PaymentCostRoute) setPaymentChannelCostMerchant(ctx echo.Context) error {
	req := &billing.PaymentChannelCostMerchant{}
	err := ctx.Bind(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestDataInvalid)
	}

	req.MerchantId = ctx.Param(common.RequestParameterId)

	if ctx.Request().Method == http.MethodPut {
		req.Id = ctx.Param(common.RequestParameterRateId)
	}

	err = h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.SetPaymentChannelCostMerchant(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "SetPaymentChannelCostMerchant", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorInternal)
	}

	if res.Status != http.StatusOK {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.JSON(http.StatusOK, res.Item)
}

// @Description Create and update system costs for money back operations
// @Example curl -X POST -H 'Authorization: Bearer %access_token_here%' -H "Content-Type: application/json" \
//		-d '{"name": "VISA", "region": "CIS", "country": "AZ", "percent": 0.01, "fix_amount": 2.34,
//		"payout_currency": "USD", "undo_reason": "chargeback", "days_from": 0, "payment_stage": 1}' \
// 		https://api.paysuper.online/admin/api/v1/payment_costs/money_back/system
//
// @Example curl -X PUT -H 'Authorization: Bearer %access_token_here%' -H "Content-Type: application/json" \
//		-d '{"name": "VISA", "region": "CIS", "country": "AZ", "percent": 0.01, "fix_amount": 2.34,
//		"payout_currency": "USD", "undo_reason": "chargeback", "days_from": 0, "payment_stage": 1}' \
// 		https://api.paysuper.online/admin/api/v1/payment_costs/money_back/system/ffffffffffffffffffffffff
func (h *PaymentCostRoute) setMoneyBackCostSystem(ctx echo.Context) error {
	req := &billing.MoneyBackCostSystem{}
	err := ctx.Bind(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestDataInvalid)
	}

	if pcId := ctx.Param(common.RequestParameterId); pcId != "" {
		req.Id = pcId
	}

	err = h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.SetMoneyBackCostSystem(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "SetMoneyBackCostSystem", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorInternal)
	}

	if res.Status != http.StatusOK {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.JSON(http.StatusOK, res.Item)
}

// @Description Create and update merchant costs for money back operations
// @Example curl -X POST -H 'Authorization: Bearer %access_token_here%' -H "Content-Type: application/json" \
//		-d '{"name": "VISA", "region": "CIS", "country": "AZ", "percent": 0.01, "fix_amount": 2.34,
//		"payout_currency": "USD", "undo_reason": "chargeback", "days_from": 0, "payment_stage": 1,
//		"is_paid_by_merchant": true}' \
// 		https://api.paysuper.online/admin/api/v1/payment_costs/money_back/merchant/ffffffffffffffffffffffff
//
// @Example curl -X PUT -H 'Authorization: Bearer %access_token_here%' -H "Content-Type: application/json" \
//		-d '{"name": "VISA", "region": "CIS", "country": "AZ", "percent": 0.01, "fix_amount": 2.34,
//		"payout_currency": "USD", "undo_reason": "chargeback", "days_from": 0, "payment_stage": 1,
//		"is_paid_by_merchant": true}' \
// 		https://api.paysuper.online/admin/api/v1/payment_costs/money_back/merchant/ffffffffffffffffffffffff/aaaaaaaaaaaaaaaaaaaaaaaa
func (h *PaymentCostRoute) setMoneyBackCostMerchant(ctx echo.Context) error {
	req := &billing.MoneyBackCostMerchant{}
	err := ctx.Bind(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestDataInvalid)
	}

	req.MerchantId = ctx.Param(common.RequestParameterId)

	if ctx.Request().Method == http.MethodPut {
		req.Id = ctx.Param(common.RequestParameterRateId)
	}

	err = h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.SetMoneyBackCostMerchant(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "SetMoneyBackCostMerchant", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorInternal)
	}

	if res.Status != http.StatusOK {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.JSON(http.StatusOK, res.Item)
}

// @Description Get all system costs for payments
// @Example @Example curl -X GET -H 'Authorization: Bearer %access_token_here%' \
// 		https://api.paysuper.online/admin/api/v1/payment_costs/channel/system/all
func (h *PaymentCostRoute) getAllPaymentChannelCostSystem(ctx echo.Context) error {
	res, err := h.dispatch.Services.Billing.GetAllPaymentChannelCostSystem(ctx.Request().Context(), &grpc.EmptyRequest{})

	if err != nil {
		h.L().Error(pkg.ErrorGrpcServiceCallFailed, logger.PairArgs("err", err.Error(), common.ErrorFieldService, pkg.ServiceName, common.ErrorFieldMethod, "GetAllPaymentChannelCostSystem"))
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorInternal)
	}

	if res.Status != http.StatusOK {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.JSON(http.StatusOK, res.Item)
}

// @Description Get all merchant costs for payments operations
// @Example @Example curl -X GET -H 'Authorization: Bearer %access_token_here%' \
// 		https://api.paysuper.online/admin/api/v1/payment_costs/channel/merchant/ffffffffffffffffffffffff/all
func (h *PaymentCostRoute) getAllPaymentChannelCostMerchant(ctx echo.Context) error {
	req := &billing.PaymentChannelCostMerchantListRequest{MerchantId: ctx.Param(common.RequestParameterId)}
	err := h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.GetAllPaymentChannelCostMerchant(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "GetAllPaymentChannelCostMerchant", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorInternal)
	}

	if res.Status != http.StatusOK {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.JSON(http.StatusOK, res.Item)
}

// @Description Get all system costs for money back operations
// @Example @Example curl -X GET -H 'Authorization: Bearer %access_token_here%' \
// 		https://api.paysuper.online/admin/api/v1/payment_costs/money_back/system/all
func (h *PaymentCostRoute) getAllMoneyBackCostSystem(ctx echo.Context) error {
	res, err := h.dispatch.Services.Billing.GetAllMoneyBackCostSystem(ctx.Request().Context(), &grpc.EmptyRequest{})

	if err != nil {
		h.L().Error(pkg.ErrorGrpcServiceCallFailed, logger.PairArgs("err", err.Error(), common.ErrorFieldService, pkg.ServiceName, common.ErrorFieldMethod, "GetAllMoneyBackCostSystem"))
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorInternal)
	}

	if res.Status != http.StatusOK {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.JSON(http.StatusOK, res.Item)
}

// @Description Get all merchant costs for money back operations
// @Example @Example curl -X GET -H 'Authorization: Bearer %access_token_here%' \
// 		https://api.paysuper.online/admin/api/v1/payment_costs/money_back/merchant/ffffffffffffffffffffffff/all
func (h *PaymentCostRoute) getAllMoneyBackCostMerchant(ctx echo.Context) error {
	req := &billing.MoneyBackCostMerchantListRequest{MerchantId: ctx.Param(common.RequestParameterId)}
	err := h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.GetAllMoneyBackCostMerchant(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "GetAllMoneyBackCostMerchant", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorInternal)
	}

	if res.Status != http.StatusOK {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.JSON(http.StatusOK, res.Item)
}
