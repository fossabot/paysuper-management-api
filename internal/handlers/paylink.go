package handlers

import (
	"github.com/ProtocolONE/go-core/logger"
	"github.com/ProtocolONE/go-core/provider"
	"github.com/labstack/echo/v4"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/grpc"
	"github.com/paysuper/paysuper-management-api/internal/dispatcher/common"
	"github.com/paysuper/paysuper-payment-link/proto"
	"net/http"
)

const (
	paylinksProjectIdPath = "/paylinks/project/:project_id"
	paylinksIdPath        = "/paylinks/:id"
	paylinksStartPath     = "/paylinks/:id/stat"
	paylinksUrlPath       = "/paylinks/:id/url"
	paylinksPath          = "/paylinks"
)

type PayLinkRoute struct {
	dispatch common.HandlerSet
	cfg      common.Config
	provider.LMT
}

func NewPayLinkRoute(set common.HandlerSet, cfg *common.Config) *PayLinkRoute {
	set.AwareSet.Logger = set.AwareSet.Logger.WithFields(logger.Fields{"router": "PayLinkRoute"})
	return &PayLinkRoute{
		dispatch: set,
		LMT:      &set.AwareSet,
		cfg:      *cfg,
	}
}

func (h *PayLinkRoute) Route(groups *common.Groups) {
	groups.AuthUser.GET(paylinksProjectIdPath, h.getPaylinksList)
	groups.AuthUser.GET(paylinksIdPath, h.getPaylink)
	groups.AuthUser.GET(paylinksStartPath, h.getPaylinkStat)
	groups.AuthUser.GET(paylinksUrlPath, h.getPaylinkUrl)
	groups.AuthUser.DELETE(paylinksIdPath, h.deletePaylink)
	groups.AuthUser.POST(paylinksPath, h.createPaylink)
	groups.AuthUser.PUT(paylinksIdPath, h.updatePaylink)
}

// @Description Get list of paylinks for project, for authenticated merchant
// @Example GET /admin/api/v1/paylinks/project/21784001599a47e5a69ac28f7af2ec22?offset=0&limit=10
func (h *PayLinkRoute) getPaylinksList(ctx echo.Context) error {
	authUser := common.ExtractUserContext(ctx)
	req := &paylink.GetPaylinksRequest{}
	err := (&common.PaylinksListBinder{
		LimitDefault:  h.cfg.LimitDefault,
		OffsetDefault: h.cfg.OffsetDefault,
	}).Bind(req, ctx)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestParamsIncorrect)
	}

	merchant, err := h.dispatch.Services.Billing.GetMerchantBy(ctx.Request().Context(), &grpc.GetMerchantByRequest{UserId: authUser.Id})
	if err != nil || merchant.Item == nil {
		if err != nil {
			h.L().Error(common.InternalErrorTemplate, logger.WithFields(logger.Fields{"err": err.Error()}))
		}
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}

	req.MerchantId = merchant.Item.Id

	err = h.dispatch.Validate.Struct(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.PayLink.GetPaylinks(ctx.Request().Context(), req)
	if err != nil {
		h.L().Error(common.InternalErrorTemplate, logger.WithFields(logger.Fields{"err": err.Error()}))
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorInternal)
	}

	return ctx.JSON(http.StatusOK, res)
}

// @Description Get paylink, for authenticated merchant
// @Example GET /admin/api/v1/paylinks/21784001599a47e5a69ac28f7af2ec22
func (h *PayLinkRoute) getPaylink(ctx echo.Context) error {
	id := ctx.Param(common.RequestParameterId)

	req := &paylink.PaylinkRequest{
		Id: id,
	}
	err := h.dispatch.Validate.Struct(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.PayLink.GetPaylink(ctx.Request().Context(), req)
	if err != nil {
		h.L().Error(common.InternalErrorTemplate, logger.WithFields(logger.Fields{"err": err.Error()}))
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorInternal)
	}

	return ctx.JSON(http.StatusOK, res)
}

// @Description Get stat for paylink
// @Example GET /admin/api/v1/paylinks/21784001599a47e5a69ac28f7af2ec22/stat
func (h *PayLinkRoute) getPaylinkStat(ctx echo.Context) error {
	id := ctx.Param(common.RequestParameterId)

	req := &paylink.PaylinkRequest{
		Id: id,
	}
	err := h.dispatch.Validate.Struct(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.PayLink.GetPaylinkStat(ctx.Request().Context(), req)
	if err != nil {
		h.L().Error(common.InternalErrorTemplate, logger.WithFields(logger.Fields{"err": err.Error()}))
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorInternal)
	}

	return ctx.JSON(http.StatusOK, res)
}

// @Description paylink public url
// @Example GET /admin/api/v1/paylinks/21784001599a47e5a69ac28f7af2ec22/url?utm_source=3wefwe&utm_medium=njytrn&utm_campaign=bdfbh5
func (h *PayLinkRoute) getPaylinkUrl(ctx echo.Context) error {
	req := &paylink.GetPaylinkURLRequest{}
	err := (&common.PaylinksUrlBinder{}).Bind(req, ctx)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestParamsIncorrect)
	}

	err = h.dispatch.Validate.Struct(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.PayLink.GetPaylinkURL(ctx.Request().Context(), req)
	if err != nil {
		h.L().Error(common.InternalErrorTemplate, logger.WithFields(logger.Fields{"err": err.Error()}))
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorInternal)
	}

	return ctx.JSON(http.StatusOK, res)
}

// @Description Get paylink, for authenticated merchant
// @Example DELETE /admin/api/v1/paylinks/21784001599a47e5a69ac28f7af2ec22
func (h *PayLinkRoute) deletePaylink(ctx echo.Context) error {
	id := ctx.Param(common.RequestParameterId)

	req := &paylink.PaylinkRequest{
		Id: id,
	}
	err := h.dispatch.Validate.Struct(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	_, err = h.dispatch.Services.PayLink.DeletePaylink(ctx.Request().Context(), req)
	if err != nil {
		h.L().Error(common.InternalErrorTemplate, logger.WithFields(logger.Fields{"err": err.Error()}))
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorInternal)
	}

	return ctx.NoContent(http.StatusNoContent)
}

// @Description Create paylink, for authenticated merchant
// @Example POST /admin/api/v1/paylinks
func (h *PayLinkRoute) createPaylink(ctx echo.Context) error {
	return h.createOrUpdatePaylink(ctx, &common.PaylinksCreateBinder{})
}

// @Description Update paylink, for authenticated merchant
// @Example PUT /admin/api/v1/paylinks/21784001599a47e5a69ac28f7af2ec22
func (h *PayLinkRoute) updatePaylink(ctx echo.Context) error {
	return h.createOrUpdatePaylink(ctx, &common.PaylinksUpdateBinder{})
}

func (h *PayLinkRoute) createOrUpdatePaylink(ctx echo.Context, binder echo.Binder) error {
	authUser := common.ExtractUserContext(ctx)
	req := &paylink.CreatePaylinkRequest{}
	err := binder.Bind(req, ctx)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestParamsIncorrect)
	}

	merchant, err := h.dispatch.Services.Billing.GetMerchantBy(ctx.Request().Context(), &grpc.GetMerchantByRequest{UserId: authUser.Id})
	if err != nil || merchant.Item == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}

	req.MerchantId = merchant.Item.Id

	err = h.dispatch.Validate.Struct(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.PayLink.CreateOrUpdatePaylink(ctx.Request().Context(), req)
	if err != nil {
		h.L().Error(common.InternalErrorTemplate, logger.WithFields(logger.Fields{"err": err.Error()}))
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorInternal)
	}

	return ctx.JSON(http.StatusOK, res)
}
