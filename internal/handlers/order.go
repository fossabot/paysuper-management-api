package handlers

import (
	"fmt"
	"github.com/ProtocolONE/go-core/logger"
	"github.com/ProtocolONE/go-core/provider"
	"github.com/labstack/echo/v4"
	"github.com/paysuper/paysuper-billing-server/pkg"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/billing"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/grpc"
	"github.com/paysuper/paysuper-management-api/internal/dispatcher/common"
	paylinkServiceConst "github.com/paysuper/paysuper-payment-link/pkg"
	"github.com/paysuper/paysuper-payment-link/proto"
	"net/http"
	"time"
)

const (
	orderIdPath              = "/order/:id"
	paylinkIdPath            = "/paylink/:id"
	orderCreatePath          = "/order/create"
	orderPath                = "/order"
	paymentPath              = "/payment"
	orderRefundsPath         = "/order/:order_id/refunds"
	orderRefundsIdsPath      = "/order/:order_id/refunds/:refund_id"
	orderReplaceCodePath     = "/order/:order_id/replace_code"
	orderLanguagePath        = "/orders/:order_id/language"
	orderCustomerPath        = "/orders/:order_id/customer"
	orderBillingAddressPath  = "/orders/:order_id/billing_address"
	orderNotifySalesPath     = "/orders/:order_id/notify_sale"
	orderNotifyNewRegionPath = "/orders/:order_id/notify_new_region"
)

const (
	orderFormTemplateName  = "order.html"
	orderInlineFormUrlMask = "%s://%s/order/%s"
	errorTemplateName      = "error.html"
)

// swagger:model CreateOrderJsonProjectResponse
type CreateOrderJsonProjectResponse struct {
	// Order unique identifier.
	Id string `json:"id"`
	//description TODO
	PaymentFormUrl string `json:"payment_form_url"`
	//description TODO
	PaymentFormData *grpc.PaymentFormJsonData `json:"payment_form_data,omitempty"`
}

type OrderListRefundsBinder struct {
	dispatch common.HandlerSet
	provider.LMT
	cfg common.Config
}

func (b *OrderListRefundsBinder) Bind(i interface{}, ctx echo.Context) error {
	db := new(echo.DefaultBinder)
	err := db.Bind(i, ctx)

	if err != nil {
		return err
	}

	structure := i.(*grpc.ListRefundsRequest)
	structure.OrderId = ctx.Param(common.RequestParameterOrderId)

	if structure.Limit <= 0 {
		structure.Limit = b.cfg.LimitDefault
	}

	return nil
}

// NewOrderListRefundsBinder
func NewOrderListRefundsBinde(set common.HandlerSet, cfg common.Config) *OrderListRefundsBinder {
	return &OrderListRefundsBinder{
		dispatch: set,
		LMT:      &set.AwareSet,
		cfg:      cfg,
	}
}

type OrderRoute struct {
	dispatch common.HandlerSet
	cfg      common.Config
	provider.LMT
}

func NewOrderRoute(set common.HandlerSet, cfg *common.Config) *OrderRoute {
	set.AwareSet.Logger = set.AwareSet.Logger.WithFields(logger.Fields{"router": "OrderRoute"})
	return &OrderRoute{
		dispatch: set,
		LMT:      &set.AwareSet,
		cfg:      *cfg,
	}
}

func (h *OrderRoute) Route(groups *common.Groups) {

	// swagger:operation GET /order/{id} Order orderIdPathGetOrderForm
	// ---
	// summary: TODO.
	// description: TODO.
	groups.Common.GET(orderIdPath, h.getOrderForm)

	// swagger:operation GET /paylink/{id} Order paylinkIdPathGetOrderForPaylink
	// ---
	// summary: TODO.
	// description: TODO.
	groups.Common.GET(paylinkIdPath, h.getOrderForPaylink) // TODO: Need a test

	// swagger:operation GET /order/create Order orderCreatePathCreateFromFormData
	// ---
	// summary: TODO.
	// description: TODO.
	groups.Common.GET(orderCreatePath, h.createFromFormData) // TODO: Need a test

	// swagger:operation POST /order/create Order orderCreatePathCreateFromFormData
	// ---
	// summary: TODO.
	// description: TODO.
	groups.Common.POST(orderCreatePath, h.createFromFormData) // TODO: Need a test

	// swagger:operation POST /api/v1/order Order orderPathCreateJson
	// ---
	// summary: Create a payment order with parameters of products.
	// description: Create a payment order.
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: data
	//   in: body
	//   description: TODO
	//   required: true
	//   schema:
	//     $ref: '#/definitions/BillingServer_OrderCreateRequest'
	// responses:
	//   200:
	//     description: Object which contain data to render payment form.
	//     schema:
	//       $ref: '#/definitions/CreateOrderJsonProjectResponse'
	//   400:
	//     description: Object with error message.
	//     schema:
	//       $ref: '#/definitions/BillingSever_ResponseError'
	//   500:
	//     description: Object with error message.
	//     schema:
	//       $ref: '#/definitions/BillingSever_ResponseError'
	groups.AuthProject.POST(orderPath, h.createJson) // TODO: Need a test

	// swagger:operation POST /api/v1/payment Order paymentPathProcessCreatePayment
	// ---
	// summary: TODO.
	// description: TODO.
	groups.AuthProject.POST(paymentPath, h.processCreatePayment) // TODO: Need a test

	// swagger:operation GET /admin/api/v1/order Order orderPathListOrdersPublic
	// ---
	// summary: TODO.
	// description: TODO.
	groups.AuthUser.GET(orderPath, h.listOrdersPublic)

	// swagger:operation GET /admin/api/v1/order/{id} Order orderIdPathGetOrderPublic
	// ---
	// summary: TODO.
	// description: TODO.
	groups.AuthUser.GET(orderIdPath, h.getOrderPublic) // TODO: Need a test

	// swagger:operation GET /admin/api/v1/order/{order_id}/refunds Order orderRefundsPathListRefunds
	// ---
	// summary: TODO.
	// description: TODO.
	groups.AuthUser.GET(orderRefundsPath, h.listRefunds)

	// swagger:operation GET /admin/api/v1/order/{order_id}/refunds/{refund_id} Order orderRefundsIdsPathGetRefund
	// ---
	// summary: TODO.
	// description: TODO.
	groups.AuthUser.GET(orderRefundsIdsPath, h.getRefund)

	// swagger:operation POST /admin/api/v1/order/{order_id}/refunds Order orderRefundsPathCreateRefund
	// ---
	// summary: TODO.
	// description: TODO.
	groups.AuthUser.POST(orderRefundsPath, h.createRefund)

	// swagger:operation PUT /admin/api/v1/order/{order_id}/replace_code Order orderReplaceCodePathReplaceCode
	// ---
	// summary: TODO.
	// description: TODO.
	groups.AuthUser.PUT(orderReplaceCodePath, h.replaceCode)

	// swagger:operation PATCH /api/v1/orders/{order_id}/language Order orderLanguagePathChangeLanguage
	// ---
	// summary: TODO.
	// description: TODO.
	groups.AuthProject.PATCH(orderLanguagePath, h.changeLanguage)

	// swagger:operation PATCH /api/v1/orders/{order_id}/customer Order orderCustomerPathChangeCustomer
	// ---
	// summary: TODO.
	// description: TODO.
	groups.AuthProject.PATCH(orderCustomerPath, h.changeCustomer)

	// swagger:operation POST /api/v1/orders/{order_id}/billing_address Order orderBillingAddressPathProcessBillingAddress
	// ---
	// summary: TODO.
	// description: TODO.
	groups.AuthProject.POST(orderBillingAddressPath, h.processBillingAddress)

	// swagger:operation POST /api/v1/orders/{order_id}/notify_sale Order orderNotifySalesPathNotifySale
	// ---
	// summary: TODO.
	// description: TODO.
	groups.AuthProject.POST(orderNotifySalesPath, h.notifySale)

	// swagger:operation POST /api/v1/orders/{order_id}/notify_new_region Order orderNotifyNewRegionPathNotifyNewRegion
	// ---
	// summary: TODO.
	// description: TODO.
	groups.AuthProject.POST(orderNotifyNewRegionPath, h.notifyNewRegion)
}

// @Summary Create order with HTML form
// @Description Create a payment order use GET or POST HTML form
// @Tags Payment Order
// @Accept multipart/form-data
// @Accept application/x-www-form-urlencoded
// @Produce html
// @Param PP_PROJECT_ID query string true "Project unique identifier in Protocol One payment solution"
// @Param PP_AMOUNT query float64 true "Order amount"
// @Param PP_ACCOUNT query string true "User unique account in project"
// @Param PP_ORDER_ID query string false "Unique order identifier in project. This field not required, BUT we're recommend send this field always"
// @Param PP_PAYMENT_METHOD query string false "Payment method identifier in Protocol One payment solution"
// @Param PP_DESCRIPTION query string false "Order description. If this field not send in request, then we're create standard order description"
// @Param PP_CURRENCY query string false "Order currency by ISO 4217 (3 chars). If this field send, then we're process amount in this currency"
// @Param PP_REGION query string false "User (payer) region code by ISO 3166-1 (2 chars) for check project packages. If this field not send, then user region will be get from user ip"
// @Param PP_PAYER_EMAIL query string false "User (payer) email"
// @Param PP_PAYER_PHONE query string false "User (payer) phone"
// @Param PP_URL_VERIFY query string false "URL for payment data verification request to project. This field can be send if it allowed in project admin panel"
// @Param PP_URL_NOTIFY query string false "URL for payment notification request to project. This field can be send if it allowed in project admin panel"
// @Param PP_URL_SUCCESS query string false "URL for redirect user after successfully completed payment. This field can be send if it allowed in project admin panel"
// @Param PP_URL_FAIL query string false "URL for redirect user after failed payment. This field can be send if it allowed in project admin panel"
// @Param PP_SIGNATURE query string false "Signature of request to verify that the data has not been changed. This field not required, BUT we're recommend send this field always"
// @Param Other query string false "Any fields on the project side that do not match the names of the reserved fields"
// @Success 302 {string} html "Redirect user to form entering payment requisites"
// @Failure 400 {string} html "Redirect user to page with error description"
// @Failure 500 {string} html "Redirect user to page with error description"
// @Router /order/create [post]
func (h *OrderRoute) createFromFormData(ctx echo.Context) error {
	req := &billing.OrderCreateRequest{
		PayerIp: ctx.RealIP(),
		IsJson:  false,
	}

	if err := (&common.OrderFormBinder{}).Bind(req, ctx); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestDataInvalid)
	}

	req.IssuerUrl = ctx.Request().Header.Get(common.HeaderReferer)

	if err := h.dispatch.Validate.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	orderResponse, err := h.dispatch.Services.Billing.OrderCreateProcess(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "OrderCreateProcess", req)
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorUnknown)
	}

	if orderResponse.Status != http.StatusOK {
		return echo.NewHTTPError(int(orderResponse.Status), orderResponse.Message)
	}

	rUrl := "/order/" + orderResponse.Item.Id

	return ctx.Redirect(http.StatusFound, rUrl)
}

// Create order from json request.
// Order can be create:
// 1) By project host2host request with sending user (customer) information.
// 2) By payment form client request with sending prepare created user (customer) identification token.
// 3) By payment form client request without anything user identification information.
func (h *OrderRoute) createJson(ctx echo.Context) error {
	req := &billing.OrderCreateRequest{}
	err := (&common.OrderJsonBinder{}).Bind(req, ctx)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestParamsIncorrect)
	}

	err = h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	// If request contain user object then paysuper must check request signature
	if req.User != nil {
		httpErr := common.CheckProjectAuthRequestSignature(h.dispatch, ctx, req.ProjectId)

		if httpErr != nil {
			return httpErr
		}
	}

	ctxReq := ctx.Request().Context()
	req.IssuerUrl = ctx.Request().Header.Get(common.HeaderReferer)

	var (
		order         *billing.Order
		orderResponse *grpc.OrderCreateProcessResponse
	)

	// If request contain prepared order identifier than try to get order by this identifier
	if req.PspOrderUuid != "" {
		req1 := &grpc.IsOrderCanBePayingRequest{
			OrderId:   req.PspOrderUuid,
			ProjectId: req.ProjectId,
		}
		rsp1, err := h.dispatch.Services.Billing.IsOrderCanBePaying(ctxReq, req1)

		if err != nil {
			common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "IsOrderCanBePaying", req)
			return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
		}

		if rsp1.Status != pkg.ResponseStatusOk {
			return echo.NewHTTPError(int(rsp1.Status), rsp1.Message)
		}

		order = rsp1.Item
	} else {
		orderResponse, err = h.dispatch.Services.Billing.OrderCreateProcess(ctxReq, req)

		if err != nil {
			common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "OrderCreateProcess", req)
			return echo.NewHTTPError(http.StatusBadRequest, common.ErrorUnknown)
		}

		if orderResponse.Status != http.StatusOK {
			return echo.NewHTTPError(int(orderResponse.Status), orderResponse.Message)
		}

		order = orderResponse.Item
	}

	response := &CreateOrderJsonProjectResponse{
		Id:             order.Uuid,
		PaymentFormUrl: fmt.Sprintf(pkg.OrderInlineFormUrlMask, ctx.Request().URL.Scheme, ctx.Request().URL.Host, order.Uuid),
	}

	if h.cfg.ReturnPaymentForm {
		req2 := &grpc.PaymentFormJsonDataRequest{
			OrderId: order.Uuid,
			Scheme:  ctx.Request().URL.Scheme,
			Host:    ctx.Request().URL.Host,
			Ip:      ctx.RealIP(),
		}
		rsp2, err := h.dispatch.Services.Billing.PaymentFormJsonDataProcess(ctxReq, req2)

		if err != nil {
			common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "PaymentFormJsonDataProcess", req)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		if rsp2.Status != pkg.ResponseStatusOk {
			return echo.NewHTTPError(int(rsp2.Status), rsp2.Message)
		}

		response.PaymentFormData = rsp2.Item
	}

	return ctx.JSON(http.StatusOK, response)
}

func (h *OrderRoute) getOrderForm(ctx echo.Context) error {

	id := ctx.Param(common.RequestParameterId)

	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestParamsIncorrect)
	}

	cookie, err := ctx.Cookie(common.CustomerTokenCookiesName)

	req := &grpc.PaymentFormJsonDataRequest{
		OrderId: id,
		Scheme:  ctx.Request().URL.Scheme,
		Host:    ctx.Request().URL.Host,
		Locale:  ctx.Request().Header.Get(common.HeaderAcceptLanguage),
		Ip:      ctx.RealIP(),
	}

	if err == nil && cookie != nil && cookie.Value != "" {
		req.Cookie = cookie.Value
	}

	res, err := h.dispatch.Services.Billing.PaymentFormJsonDataProcess(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "PaymentFormJsonDataProcess", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}

	if res.Status != http.StatusOK {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	if res.Item.Cookie != "" {
		cookie := new(http.Cookie)
		cookie.Name = common.CustomerTokenCookiesName
		cookie.Value = res.Item.Cookie
		cookie.Expires = time.Now().Add(h.cfg.CustomerTokenCookiesLifetime)
		cookie.HttpOnly = true
		ctx.SetCookie(cookie)
	}

	return ctx.Render(
		http.StatusOK,
		orderFormTemplateName,
		map[string]interface{}{
			"Order":                   res,
			"WebSocketUrl":            h.cfg.WebsocketUrl,
			"PaymentFormJsLibraryUrl": h.cfg.PaymentFormJsLibraryUrl,
		},
	)
}

// Create order from payment link and redirect to order payment form
func (h *OrderRoute) getOrderForPaylink(ctx echo.Context) error {
	paylinkId := ctx.Param(common.RequestParameterId)

	req := &paylink.PaylinkRequest{
		Id: paylinkId,
	}

	err := h.dispatch.Validate.Struct(req)

	if err != nil {
		h.L().Error("Cannot validate request", logger.PairArgs("err", err.Error(), common.ErrorFieldRequest, req))

		return ctx.Render(http.StatusBadRequest, errorTemplateName, map[string]interface{}{})
	}

	ctxReq := ctx.Request().Context()
	pl, err := h.dispatch.Services.PayLink.GetPaylink(ctxReq, req)

	if err != nil {
		return ctx.Render(http.StatusNotFound, errorTemplateName, map[string]interface{}{})
	}

	oReq := &billing.OrderCreateRequest{
		ProjectId: pl.ProjectId,
		PayerIp:   ctx.RealIP(),
		Products:  pl.Products,
		PrivateMetadata: map[string]string{
			"PaylinkId": paylinkId,
		},
		IssuerUrl:  ctx.Request().Header.Get(common.HeaderReferer),
		IsEmbedded: false,
	}
	params := ctx.QueryParams()
	if v, ok := params[common.RequestParameterUtmSource]; ok {
		oReq.PrivateMetadata[common.RequestParameterUtmSource] = v[0]
	}
	if v, ok := params[common.RequestParameterUtmMedium]; ok {
		oReq.PrivateMetadata[common.RequestParameterUtmMedium] = v[0]
	}
	if v, ok := params[common.RequestParameterUtmCampaign]; ok {
		oReq.PrivateMetadata[common.RequestParameterUtmCampaign] = v[0]
	}

	orderResponse, err := h.dispatch.Services.Billing.OrderCreateProcess(ctxReq, oReq)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "OrderCreateProcess", req)
		return ctx.Render(http.StatusBadRequest, errorTemplateName, map[string]interface{}{})
	}

	if orderResponse.Status != http.StatusOK {
		return echo.NewHTTPError(int(orderResponse.Status), orderResponse.Message)
	}

	inlineFormRedirectUrl := fmt.Sprintf(orderInlineFormUrlMask, ctx.Request().URL.Scheme, ctx.Request().URL.Host, orderResponse.Item.Uuid)
	qs := ctx.QueryString()

	if qs != "" {
		inlineFormRedirectUrl += "?" + qs
	}

	go func() {
		_, err := h.dispatch.Services.PayLink.IncrPaylinkVisits(ctxReq, &paylink.PaylinkRequest{Id: paylinkId})

		if err != nil {
			common.LogSrvCallFailedGRPC(h.L(), err, paylinkServiceConst.ServiceName, "IncrPaylinkVisits", req)
		}
	}()

	return ctx.Redirect(http.StatusFound, inlineFormRedirectUrl)
}

// @Description Get order by id
// @Example curl -X GET -H 'Authorization: Bearer %access_token_here%' -H 'Content-Type: application/json' \
//  https://api.paysuper.online/admin/api/v1/order/%order_id_here%
func (h *OrderRoute) getOrderPublic(ctx echo.Context) error {
	req := &grpc.GetOrderRequest{
		Id: ctx.Param(common.RequestParameterId),
	}

	err := h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.GetOrderPublic(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "GetOrderPublic", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}

	if res.Status != pkg.ResponseStatusOk {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.JSON(http.StatusOK, res.Item)
}

// @Description Get orders list
// @Example curl -X GET -H 'Authorization: Bearer %access_token_here%' -H 'Content-Type: application/json' \
//  https://api.paysuper.online/admin/api/v1/order?project[]=%project_identifier_here%
func (h *OrderRoute) listOrdersPublic(ctx echo.Context) error {

	req := &grpc.ListOrdersRequest{}
	err := ctx.Bind(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestParamsIncorrect)
	}

	if req.Limit <= 0 {
		req.Limit = h.cfg.LimitDefault
	}

	if req.Offset <= 0 {
		req.Offset = h.cfg.OffsetDefault
	}

	err = h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.FindAllOrdersPublic(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "FindAllOrdersPublic", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}

	if res.Status != pkg.ResponseStatusOk {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.JSON(http.StatusOK, res.Item)
}

// Create payment by order
// route POST /api/v1/payment
func (h *OrderRoute) processCreatePayment(ctx echo.Context) error {
	data := make(map[string]string)
	err := (&common.PaymentCreateProcessBinder{}).Bind(data, ctx)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestDataInvalid)
	}

	req := &grpc.PaymentCreateRequest{
		Data:           data,
		AcceptLanguage: ctx.Request().Header.Get(common.HeaderAcceptLanguage),
		UserAgent:      ctx.Request().Header.Get(common.HeaderUserAgent),
		Ip:             ctx.RealIP(),
	}
	res, err := h.dispatch.Services.Billing.PaymentCreateProcess(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "PaymentCreateProcess", req)
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorUnknown)
	}

	if res.Status != pkg.ResponseStatusOk {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	body := map[string]interface{}{
		"redirect_url":  res.RedirectUrl,
		"need_redirect": res.NeedRedirect,
	}

	return ctx.JSON(http.StatusOK, body)
}

func (h *OrderRoute) getRefund(ctx echo.Context) error {
	req := &grpc.GetRefundRequest{
		OrderId:  ctx.Param(common.RequestParameterOrderId),
		RefundId: ctx.Param(common.RequestParameterRefundId),
	}

	err := h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.GetRefund(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "GetRefund", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}

	if res.Status != pkg.ResponseStatusOk {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.JSON(http.StatusOK, res.Item)
}

func (h *OrderRoute) listRefunds(ctx echo.Context) error {
	req := &grpc.ListRefundsRequest{}
	err := (&OrderListRefundsBinder{}).Bind(req, ctx)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestParamsIncorrect)
	}

	err = h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.ListRefunds(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "ListRefunds", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}

	return ctx.JSON(http.StatusOK, res)
}

func (h *OrderRoute) replaceCode(ctx echo.Context) error {
	req := &grpc.ChangeCodeInOrderRequest{}
	if err := ctx.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestParamsIncorrect)
	}

	req.OrderId = ctx.Param("order_id")
	if err := h.dispatch.Validate.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res := &grpc.ChangeCodeInOrderResponse{}

	res, err := h.dispatch.Services.Billing.ChangeCodeInOrder(ctx.Request().Context(), req)
	if err != nil {
		h.L().Error(common.InternalErrorTemplate, logger.WithFields(logger.Fields{"err": err.Error()}))
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}

	if res.Status != pkg.ResponseStatusOk {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.JSON(http.StatusOK, res.Order)
}

func (h *OrderRoute) createRefund(ctx echo.Context) error {
	authUser := common.ExtractUserContext(ctx)
	req := &grpc.CreateRefundRequest{}
	err := ctx.Bind(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestParamsIncorrect)
	}

	req.OrderId = ctx.Param(common.RequestParameterOrderId)
	err = h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	req.CreatorId = authUser.Id
	res, err := h.dispatch.Services.Billing.CreateRefund(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "CreateRefund", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}

	if res.Status != pkg.ResponseStatusOk {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.JSON(http.StatusCreated, res.Item)
}

func (h *OrderRoute) changeLanguage(ctx echo.Context) error {
	orderId := ctx.Param(common.RequestParameterOrderId)

	if orderId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorIncorrectOrderId)
	}

	req := &grpc.PaymentFormUserChangeLangRequest{}
	err := ctx.Bind(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestParamsIncorrect)
	}

	req.AcceptLanguage = ctx.Request().Header.Get(common.HeaderAcceptLanguage)
	req.UserAgent = ctx.Request().Header.Get(common.HeaderUserAgent)
	req.Ip = ctx.RealIP()
	req.OrderId = orderId
	err = h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.PaymentFormLanguageChanged(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "PaymentFormLanguageChanged", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}

	if res.Status != pkg.ResponseStatusOk {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.JSON(http.StatusOK, res.Item)
}

func (h *OrderRoute) changeCustomer(ctx echo.Context) error {
	orderId := ctx.Param(common.RequestParameterOrderId)

	if orderId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorIncorrectOrderId)
	}

	req := &grpc.PaymentFormUserChangePaymentAccountRequest{}
	err := ctx.Bind(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestParamsIncorrect)
	}

	req.AcceptLanguage = ctx.Request().Header.Get(common.HeaderAcceptLanguage)
	req.UserAgent = ctx.Request().Header.Get(common.HeaderUserAgent)
	req.Ip = ctx.RealIP()
	req.OrderId = orderId
	err = h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.PaymentFormPaymentAccountChanged(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "PaymentFormPaymentAccountChanged", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}

	if res.Status != pkg.ResponseStatusOk {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.JSON(http.StatusOK, res.Item)
}

func (h *OrderRoute) processBillingAddress(ctx echo.Context) error {
	orderId := ctx.Param(common.RequestParameterOrderId)

	if orderId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorIncorrectOrderId)
	}

	req := &grpc.ProcessBillingAddressRequest{}
	err := ctx.Bind(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestParamsIncorrect)
	}

	req.OrderId = orderId
	err = h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.ProcessBillingAddress(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "ProcessBillingAddress", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}

	if res.Status != pkg.ResponseStatusOk {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.JSON(http.StatusOK, res.Item)
}

/*
Switching sales notifications for order customer
POST /api/v1/orders/:order_id/notify_sale
@Param [email] string
@Param enableNotification string true|false
*/
func (h *OrderRoute) notifySale(ctx echo.Context) error {
	orderId := ctx.Param(common.RequestParameterOrderId)

	if orderId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorIncorrectOrderId)
	}

	req := &grpc.SetUserNotifyRequest{}
	err := ctx.Bind(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestParamsIncorrect)
	}

	req.OrderUuid = orderId
	err = h.dispatch.Validate.Struct(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	_, err = h.dispatch.Services.Billing.SetUserNotifySales(ctx.Request().Context(), req)
	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "SetUserNotifySales", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}

	return ctx.NoContent(http.StatusNoContent)
}

/*
Switching notifications customer about new regions allowed to make payments
POST /api/v1/orders/:order_uuid/notify_new_region
@Param [email] string
@Param enableNotification string true|false
*/
func (h *OrderRoute) notifyNewRegion(ctx echo.Context) error {
	orderUuid := ctx.Param(common.RequestParameterOrderId)

	if orderUuid == "" {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorIncorrectOrderId)
	}

	req := &grpc.SetUserNotifyRequest{}
	err := ctx.Bind(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestParamsIncorrect)
	}

	req.OrderUuid = orderUuid
	err = h.dispatch.Validate.Struct(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	_, err = h.dispatch.Services.Billing.SetUserNotifyNewRegion(ctx.Request().Context(), req)
	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "SetUserNotifyNewRegion", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}

	return ctx.NoContent(http.StatusNoContent)
}
