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
	balancePath         = "/balance"
	balanceMerchantPath = "/balance/:merchant_id"
)

type BalanceRoute struct {
	dispatch common.HandlerSet
	cfg      common.Config
	provider.LMT
}

// NewBalanceRoute
func NewBalanceRoute(set common.HandlerSet, cfg *common.Config) *BalanceRoute {
	set.AwareSet.Logger = set.AwareSet.Logger.WithFields(logger.Fields{"router": "BalanceRoute"})
	return &BalanceRoute{
		dispatch: set,
		LMT:      &set.AwareSet,
		cfg:      *cfg,
	}
}

func (h *BalanceRoute) Route(groups *common.Groups) {
	groups.AuthUser.GET(balancePath, h.getMerchantBalance)
	groups.AuthUser.GET(balanceMerchantPath, h.getMerchantBalance)
}

// Get merchant balance
// GET /admin/api/v1/balance - for current merchant
// GET /admin/api/v1/balance/:merchant_id - for any merchant by it's id
func (h *BalanceRoute) getMerchantBalance(ctx echo.Context) error {
	authUser := common.ExtractUserContext(ctx)

	req := &grpc.GetMerchantBalanceRequest{}
	merchantId := ctx.Param(common.RequestParameterMerchantId)

	if merchantId == "" {
		mReq := &grpc.GetMerchantByRequest{UserId: authUser.Id}
		merchant, err := h.dispatch.Services.Billing.GetMerchantBy(ctx.Request().Context(), mReq)
		if err != nil {
			common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "GetMerchantBy", req)
			return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
		}
		if merchant.Status != http.StatusOK {
			return echo.NewHTTPError(int(merchant.Status), merchant.Message)
		}
		merchantId = merchant.Item.Id
	}

	req.MerchantId = merchantId

	err := h.dispatch.Validate.Struct(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.GetMerchantBalance(ctx.Request().Context(), req)
	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "GetMerchantBalance", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}
	if res.Status != http.StatusOK {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}
	return ctx.JSON(http.StatusOK, res.Item)
}
