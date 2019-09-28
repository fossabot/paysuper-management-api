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
	keysIdPath = "/keys/:key_id"
)

type KeyRoute struct {
	dispatch common.HandlerSet
	cfg      common.Config
	provider.LMT
}

func NewKeyRoute(set common.HandlerSet, cfg *common.Config) *KeyRoute {
	set.AwareSet.Logger = set.AwareSet.Logger.WithFields(logger.Fields{"router": "KeyRoute"})
	return &KeyRoute{
		dispatch: set,
		LMT:      &set.AwareSet,
		cfg:      *cfg,
	}
}

func (h *KeyRoute) Route(groups *common.Groups) {
	groups.AuthUser.GET(keysIdPath, h.getKeyInfo)
}

func (h *KeyRoute) getKeyInfo(ctx echo.Context) error {
	req := &grpc.KeyForOrderRequest{
		KeyId: ctx.Param("key_id"),
	}

	if err := h.dispatch.Validate.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.GetKeyByID(ctx.Request().Context(), req)
	if err != nil {
		h.L().Error(common.InternalErrorTemplate, logger.PairArgs("err", err.Error()))
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorInternal)
	}

	if res.Status != pkg.ResponseStatusOk {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.JSON(http.StatusOK, res.Key)
}
