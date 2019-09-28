package handlers

import (
	"github.com/Nerufa/go-shared/logger"
	"github.com/Nerufa/go-shared/provider"
	"github.com/labstack/echo/v4"
	"github.com/paysuper/paysuper-billing-server/pkg"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/billing"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/grpc"
	"github.com/paysuper/paysuper-management-api/internal/dispatcher/common"
	"net/http"
)

const (
	cardPayWebHookPaymentNotifyPath = "/cardpay/payment"
	cardPayWebHookRefundNotifyPath  = "/cardpay/refund"
)

type CardPayWebHook struct {
	dispatch common.HandlerSet
	cfg      common.Config
	provider.LMT
}

func NewCardPayWebHook(set common.HandlerSet, cfg *common.Config) *CardPayWebHook {
	set.AwareSet.Logger = set.AwareSet.Logger.WithFields(logger.Fields{"router": "CardPayWebHook"})
	return &CardPayWebHook{
		dispatch: set,
		LMT:      &set.AwareSet,
		cfg:      *cfg,
	}
}

func (h *CardPayWebHook) Route(groups *common.Groups) {
	groups.WebHooks.POST(cardPayWebHookPaymentNotifyPath, h.paymentCallback)
	groups.WebHooks.POST(cardPayWebHookRefundNotifyPath, h.refundCallback)
}

func (h *CardPayWebHook) paymentCallback(ctx echo.Context) error {

	st := &billing.CardPayPaymentCallback{}

	if err := ctx.Bind(st); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestDataInvalid)
	}

	if err := h.dispatch.Validate.Struct(st); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	req := &grpc.PaymentNotifyRequest{
		OrderId:   st.MerchantOrder.Id,
		Request:   common.ExtractRawBodyContext(ctx),
		Signature: ctx.Request().Header.Get(common.CardPayPaymentResponseHeaderSignature),
	}

	res, err := h.dispatch.Services.Billing.PaymentCallbackProcess(ctx.Request().Context(), req)

	if err != nil {
		h.L().Error(common.InternalErrorTemplate, logger.WithFields(logger.Fields{"err": err.Error()}))
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorUnknown)
	}

	var httpStatus int
	var message = map[string]string{"message": res.Error}

	switch res.Status {
	case pkg.StatusErrorValidation:
		httpStatus = http.StatusBadRequest
		break
	case pkg.StatusErrorSystem:
		httpStatus = http.StatusInternalServerError
		break
	case pkg.StatusTemporary:
		httpStatus = http.StatusGone
		break
	default:
		httpStatus = http.StatusOK
		message["message"] = "Payment successfully complete"
	}

	return ctx.JSON(httpStatus, message)
}

func (h *CardPayWebHook) refundCallback(ctx echo.Context) error {

	st := &billing.CardPayRefundCallback{}
	err := ctx.Bind(st)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestParamsIncorrect)
	}

	err = h.dispatch.Validate.Struct(st)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	req := &grpc.CallbackRequest{
		Handler:   pkg.PaymentSystemHandlerCardPay,
		Body:      common.ExtractRawBodyContext(ctx),
		Signature: ctx.Request().Header.Get(common.CardPayPaymentResponseHeaderSignature),
	}

	res, err := h.dispatch.Services.Billing.ProcessRefundCallback(ctx.Request().Context(), req)

	if err != nil {
		h.L().Error(common.InternalErrorTemplate, logger.WithFields(logger.Fields{"err": err.Error()}))
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}

	if res.Status != pkg.ResponseStatusOk {
		return echo.NewHTTPError(int(res.Status), res.Error)
	}

	if res.Error != "" {
		return ctx.JSON(http.StatusOK, map[string]string{"message": res.Error})
	}

	return ctx.NoContent(http.StatusOK)
}
