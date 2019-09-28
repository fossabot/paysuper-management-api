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
	tokenPath = "/tokens"
)

type TokenRoute struct {
	dispatch common.HandlerSet
	cfg      common.Config
	provider.LMT
}

func NewTokenRoute(set common.HandlerSet, cfg *common.Config) *TokenRoute {
	set.AwareSet.Logger = set.AwareSet.Logger.WithFields(logger.Fields{"router": "TokenRoute"})
	return &TokenRoute{
		dispatch: set,
		LMT:      &set.AwareSet,
		cfg:      *cfg,
	}
}

func (h *TokenRoute) Route(groups *common.Groups) {
	groups.AuthProject.POST(tokenPath, h.createToken)
}

func (h *TokenRoute) createToken(ctx echo.Context) error {
	req := &grpc.TokenRequest{}
	err := ctx.Bind(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestParamsIncorrect)
	}

	err = h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	err = common.CheckProjectAuthRequestSignature(h.dispatch, ctx, req.Settings.ProjectId)

	if err != nil {
		return err
	}

	res, err := h.dispatch.Services.Billing.CreateToken(ctx.Request().Context(), req)

	if err != nil {
		h.L().Error(common.InternalErrorTemplate, logger.WithFields(logger.Fields{"err": err.Error()}))
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}

	if res.Status != pkg.ResponseStatusOk {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.JSON(http.StatusOK, map[string]string{"token": res.Token})
}
