package handlers

import (
	"github.com/ProtocolONE/go-core/logger"
	"github.com/ProtocolONE/go-core/provider"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo/v4"
	"github.com/paysuper/paysuper-billing-server/pkg"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/grpc"
	"github.com/paysuper/paysuper-management-api/internal/dispatcher/common"
	"net/http"
)

const (
	productsPath         = "/products"
	productsMerchantPath = "/products"
	productsIdPath       = "/products/:id"
	productsPricesPath   = "/products/:id/prices"
)

type ProductRoute struct {
	dispatch common.HandlerSet
	cfg      common.Config
	provider.LMT
}

func NewProductRoute(set common.HandlerSet, cfg *common.Config) *ProductRoute {
	set.AwareSet.Logger = set.AwareSet.Logger.WithFields(logger.Fields{"router": "ProductRoute"})
	return &ProductRoute{
		dispatch: set,
		LMT:      &set.AwareSet,
		cfg:      *cfg,
	}
}

func (h *ProductRoute) Route(groups *common.Groups) {
	groups.AuthUser.GET(productsPath, h.getProductsList)
	groups.AuthUser.GET(productsMerchantPath, h.getProductsList)
	groups.AuthUser.POST(productsPath, h.createProduct)
	groups.AuthUser.GET(productsIdPath, h.getProduct)
	groups.AuthUser.PUT(productsIdPath, h.updateProduct)
	groups.AuthUser.DELETE(productsIdPath, h.deleteProduct)
	groups.AuthUser.GET(productsPricesPath, h.getProductPrices)    // TODO: Need test
	groups.AuthUser.PUT(productsPricesPath, h.updateProductPrices) // TODO: Need test
}

// @Description Get list of products for authenticated merchant
// @Example GET /admin/api/v1/products?name=car&sku=ru_0&project_id=5bdc39a95d1e1100019fb7df&offset=0&limit=10
func (h *ProductRoute) getProductsList(ctx echo.Context) error {
	authUser := common.ExtractUserContext(ctx)
	req := &grpc.ListProductsRequest{}
	err := (&common.ProductsGetProductsListBinder{
		LimitDefault:  h.cfg.LimitDefault,
		OffsetDefault: h.cfg.OffsetDefault,
	}).Bind(req, ctx)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestParamsIncorrect)
	}

	reqCtx := ctx.Request().Context()
	merchantId := ctx.Param(common.RequestParameterId)

	if merchantId == "" {
		merchant, err := h.dispatch.Services.Billing.GetMerchantBy(reqCtx, &grpc.GetMerchantByRequest{UserId: authUser.Id})

		if err != nil {
			common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "GetMerchantBy", req)
			return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
		}

		if merchant.Status != pkg.ResponseStatusOk {
			return echo.NewHTTPError(int(merchant.Status), merchant.Message)
		}

		merchantId = merchant.Item.Id
	}

	req.MerchantId = merchantId
	err = h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.ListProducts(reqCtx, req)
	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "ListProducts", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorInternal)
	}

	return ctx.JSON(http.StatusOK, res)
}

// @Description Get product for authenticated merchant
// @Example GET /admin/api/v1/products/5c99288068add43f74be9c1d
func (h *ProductRoute) getProduct(ctx echo.Context) error {
	authUser := common.ExtractUserContext(ctx)

	id := ctx.Param(common.RequestParameterId)
	if id == "" || bson.IsObjectIdHex(id) == false {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorIncorrectProductId)
	}

	merchant, err := h.dispatch.Services.Billing.GetMerchantBy(ctx.Request().Context(), &grpc.GetMerchantByRequest{UserId: authUser.Id})
	if err != nil || merchant.Item == nil {
		if err != nil {
			h.L().Error(common.InternalErrorTemplate, logger.WithFields(logger.Fields{"err": err.Error()}))
		}
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}

	req := &grpc.RequestProduct{
		Id:         id,
		MerchantId: merchant.Item.Id,
	}

	err = h.dispatch.Validate.Struct(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.GetProduct(ctx.Request().Context(), req)
	if err != nil {
		h.L().Error(common.InternalErrorTemplate, logger.WithFields(logger.Fields{"err": err.Error()}))
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorInternal)
	}

	return ctx.JSON(http.StatusOK, res)
}

// @Description Delete product for authenticated merchant
// @Example DELETE /admin/api/v1/products/5c99288068add43f74be9c1d
func (h *ProductRoute) deleteProduct(ctx echo.Context) error {
	authUser := common.ExtractUserContext(ctx)
	id := ctx.Param(common.RequestParameterId)
	if id == "" || bson.IsObjectIdHex(id) == false {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorIncorrectProductId)
	}

	merchant, err := h.dispatch.Services.Billing.GetMerchantBy(ctx.Request().Context(), &grpc.GetMerchantByRequest{UserId: authUser.Id})
	if err != nil || merchant.Item == nil {
		if err != nil {
			h.L().Error(common.InternalErrorTemplate, logger.WithFields(logger.Fields{"err": err.Error()}))
		}
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}

	req := &grpc.RequestProduct{
		Id:         id,
		MerchantId: merchant.Item.Id,
	}

	err = h.dispatch.Validate.Struct(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	_, err = h.dispatch.Services.Billing.DeleteProduct(ctx.Request().Context(), req)
	if err != nil {
		h.L().Error(common.InternalErrorTemplate, logger.WithFields(logger.Fields{"err": err.Error()}))
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorInternal)
	}

	return ctx.NoContent(http.StatusNoContent)
}

// @Description Create new product for authenticated merchant
// @Example curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" \
//      -H "Authorization: Bearer %access_token_here%" \
//      -d '{"object": "product", "type": "simple_product", "sku": "ru_0_doom_2", "name": {"en": "Doom II"},
//          "default_currency": "USD", "enabled": true, "prices": [{"amount": 12.93, "currency": "USD"}],
//          "description": {"en": "Doom II description"}, "long_description": {}, "project_id": "5bdc39a95d1e1100019fb7df"}' \
//      https://api.paysuper.online/admin/api/v1/products
func (h *ProductRoute) createProduct(ctx echo.Context) error {
	return h.createOrUpdateProduct(ctx, &common.ProductsCreateProductBinder{})
}

// @Description Update existing product for authenticated merchant
// @Example curl -X PUT -H "Accept: application/json" -H "Content-Type: application/json" \
//      -H "Authorization: Bearer %access_token_here%" \
//      -d '{"object": "product", "type": "simple_product", "sku": "ru_0_doom_4", "name": {"en": "Doom IV"},
//          "default_currency": "USD", "enabled": true, "prices": [{"amount": 146.00, "currency": "USD"}],
//          "description": {"en": "Doom IV description"}, "long_description": {}, "project_id": "5bdc39a95d1e1100019fb7df"}' \
//      https://api.paysuper.online/admin/api/v1/products/5c99288068add43f74be9c1d
func (h *ProductRoute) updateProduct(ctx echo.Context) error {
	return h.createOrUpdateProduct(ctx, &common.ProductsUpdateProductBinder{})
}

func (h *ProductRoute) createOrUpdateProduct(ctx echo.Context, binder echo.Binder) error {
	authUser := common.ExtractUserContext(ctx)
	req := &grpc.Product{}
	err := binder.Bind(req, ctx)

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

	res, err := h.dispatch.Services.Billing.CreateOrUpdateProduct(ctx.Request().Context(), req)
	if err != nil {
		h.L().Error(common.InternalErrorTemplate, logger.WithFields(logger.Fields{"err": err.Error()}))
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorInternal)
	}

	return ctx.JSON(http.StatusOK, res)
}

func (h *ProductRoute) getProductPrices(ctx echo.Context) error {
	id := ctx.Param(common.RequestParameterId)
	if id == "" || bson.IsObjectIdHex(id) == false {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorIncorrectProductId)
	}

	req := &grpc.RequestProduct{
		Id: id,
	}

	if err := h.dispatch.Validate.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.GetProductPrices(ctx.Request().Context(), req)
	if err != nil {
		h.L().Error(common.InternalErrorTemplate, logger.WithFields(logger.Fields{"err": err.Error()}))
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorMessageGetProductPrice)
	}

	return ctx.JSON(http.StatusOK, res)
}

func (h *ProductRoute) updateProductPrices(ctx echo.Context) error {
	id := ctx.Param(common.RequestParameterId)
	if id == "" || bson.IsObjectIdHex(id) == false {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorIncorrectProductId)
	}

	req := &grpc.UpdateProductPricesRequest{}

	if err := h.dispatch.Validate.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.UpdateProductPrices(ctx.Request().Context(), req)
	if err != nil {
		h.L().Error(common.InternalErrorTemplate, logger.WithFields(logger.Fields{"err": err.Error()}))
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorMessageUpdateProductPrice)
	}

	return ctx.JSON(http.StatusOK, res)
}
