package common

import (
	"bytes"
	"context"
	"fmt"
	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo/v4"
	"github.com/paysuper/paysuper-billing-server/pkg"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/billing"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/grpc"
	"io/ioutil"
	"strconv"
)

type OrderFormBinder struct{}
type OrderJsonBinder struct{}
type OrderRevenueDynamicRequestBinder struct{}
type OrderAccountingPaymentRequestBinder struct{}
type PaymentCreateProcessBinder struct{}
type OnboardingMerchantListingBinder struct {
	LimitDefault, OffsetDefault int32
}
type OnboardingChangeMerchantStatusBinder struct{}
type OnboardingNotificationsListBinder struct {
	LimitDefault, OffsetDefault int32
}
type OnboardingGetPaymentMethodBinder struct{}
type OnboardingChangePaymentMethodBinder struct{}
type OnboardingCreateNotificationBinder struct{}
type ProductsGetProductsListBinder struct {
	LimitDefault, OffsetDefault int32
}
type ProductsCreateProductBinder struct{}
type ProductsUpdateProductBinder struct{}

// ChangeMerchantDataRequestBinder
type ChangeMerchantDataRequestBinder struct {
	dispatch HandlerSet
	provider.LMT
	cfg Config
}

// NewChangeMerchantDataRequestBinder
func NewChangeMerchantDataRequestBinder(set HandlerSet, cfg Config) *ChangeMerchantDataRequestBinder {
	return &ChangeMerchantDataRequestBinder{
		dispatch: set,
		LMT:      &set.AwareSet,
		cfg:      cfg,
	}
}

// ChangeProjectRequestBinder
type ChangeProjectRequestBinder struct {
	dispatch HandlerSet
	provider.LMT
	cfg Config
}

// NewChangeProjectRequestBinder
func NewChangeProjectRequestBinder(set HandlerSet, cfg Config) *ChangeProjectRequestBinder {
	return &ChangeProjectRequestBinder{
		dispatch: set,
		LMT:      &set.AwareSet,
		cfg:      cfg,
	}
}

// Bind
func (cb *OrderFormBinder) Bind(i interface{}, ctx echo.Context) (err error) {
	db := new(echo.DefaultBinder)

	if err = db.Bind(i, ctx); err != nil {
		return err
	}

	params, err := ctx.FormParams()
	addParams := make(map[string]string)
	rawParams := make(map[string]string)

	if err != nil {
		return err
	}

	o := i.(*billing.OrderCreateRequest)

	for key, value := range params {
		if _, ok := OrderReservedWords[key]; !ok {
			addParams[key] = value[0]
		}

		rawParams[key] = value[0]
	}

	o.Other = addParams
	o.RawParams = rawParams

	return
}

// Bind
func (cb *OrderJsonBinder) Bind(i interface{}, ctx echo.Context) (err error) {
	var buf []byte

	if ctx.Request().Body != nil {
		buf, err = ioutil.ReadAll(ctx.Request().Body)
		rdr := ioutil.NopCloser(bytes.NewBuffer(buf))

		if err != nil {
			return err
		}

		ctx.Request().Body = rdr
	}

	db := new(echo.DefaultBinder)
	if err = db.Bind(i, ctx); err != nil {
		return err
	}

	structure := i.(*billing.OrderCreateRequest)
	structure.RawBody = string(buf)

	return
}

// Bind
func (cb *PaymentCreateProcessBinder) Bind(i interface{}, ctx echo.Context) (err error) {
	db := new(echo.DefaultBinder)
	untypedData := make(map[string]interface{})

	if err = db.Bind(&untypedData, ctx); err != nil {
		return
	}

	data := i.(map[string]string)

	for k, v := range untypedData {
		switch sv := v.(type) {
		case bool:
			data[k] = "0"

			if sv == true {
				data[k] = "1"
			}
			break
		default:
			data[k] = fmt.Sprintf("%v", sv)
		}
	}

	return
}

// Bind
func (cb *OnboardingMerchantListingBinder) Bind(i interface{}, ctx echo.Context) (err error) {
	db := new(echo.DefaultBinder)
	err = db.Bind(i, ctx)

	if err != nil {
		return err
	}

	params := ctx.QueryParams()
	structure := i.(*grpc.MerchantListingRequest)

	if structure.Limit <= 0 {
		structure.Limit = cb.LimitDefault
	}

	if v, ok := params[RequestParameterIsSigned]; ok {
		if v[0] == "0" || v[0] == "false" {
			structure.IsSigned = 2
		} else {
			if v[0] == "1" || v[0] == "true" {
				structure.IsSigned = 2
			} else {
				return ErrorRequestParamsIncorrect
			}
		}
	}

	return
}

// Bind
func (cb *OnboardingNotificationsListBinder) Bind(i interface{}, ctx echo.Context) error {
	db := new(echo.DefaultBinder)
	err := db.Bind(i, ctx)

	if err != nil {
		return err
	}

	params := ctx.QueryParams()
	structure := i.(*grpc.ListingNotificationRequest)
	structure.MerchantId = ctx.Param(RequestParameterMerchantId)

	if structure.Limit <= 0 {
		structure.Limit = cb.LimitDefault
	}

	if v, ok := params[RequestParameterIsSystem]; ok {
		if v[0] == "0" || v[0] == "false" {
			structure.IsSystem = 1
		} else {
			structure.IsSystem = 2
		}
	}

	return nil
}

// Bind
func (cb *OnboardingGetPaymentMethodBinder) Bind(i interface{}, ctx echo.Context) error {
	merchantId := ctx.Param(RequestParameterMerchantId)
	paymentMethodId := ctx.Param(RequestParameterPaymentMethodId)

	if merchantId == "" || bson.IsObjectIdHex(merchantId) == false {
		return ErrorIncorrectMerchantId
	}

	if paymentMethodId == "" || bson.IsObjectIdHex(paymentMethodId) == false {
		return ErrorIncorrectPaymentMethodId
	}

	structure := i.(*grpc.GetMerchantPaymentMethodRequest)
	structure.MerchantId = merchantId
	structure.PaymentMethodId = paymentMethodId

	return nil
}

// Bind
func (cb *OnboardingChangePaymentMethodBinder) Bind(i interface{}, ctx echo.Context) error {
	db := new(echo.DefaultBinder)
	err := db.Bind(i, ctx)

	if err != nil {
		return err
	}

	structure := i.(*grpc.MerchantPaymentMethodRequest)
	merchantId := ctx.Param(RequestParameterMerchantId)
	methodId := ctx.Param(RequestParameterPaymentMethodId)

	if merchantId == "" || bson.IsObjectIdHex(merchantId) == false {
		return ErrorIncorrectMerchantId
	}

	if methodId == "" || bson.IsObjectIdHex(methodId) == false ||
		structure.PaymentMethod.Id != methodId {
		return ErrorIncorrectPaymentMethodId
	}

	structure.MerchantId = merchantId

	return nil
}

// Bind
func (b *OnboardingChangeMerchantStatusBinder) Bind(i interface{}, ctx echo.Context) error {
	db := new(echo.DefaultBinder)
	err := db.Bind(i, ctx)

	if err != nil {
		return err
	}

	merchantId := ctx.Param(RequestParameterId)

	if merchantId == "" || bson.IsObjectIdHex(merchantId) == false {
		return ErrorIncorrectMerchantId
	}

	structure := i.(*grpc.MerchantChangeStatusRequest)
	structure.MerchantId = merchantId

	return nil
}

// Bind
func (b *OnboardingCreateNotificationBinder) Bind(i interface{}, ctx echo.Context) error {
	db := new(echo.DefaultBinder)
	err := db.Bind(i, ctx)

	if err != nil {
		return err
	}

	merchantId := ctx.Param(RequestParameterMerchantId)

	if merchantId == "" || bson.IsObjectIdHex(merchantId) == false {
		return ErrorIncorrectMerchantId
	}

	structure := i.(*grpc.NotificationRequest)
	structure.MerchantId = merchantId

	return nil
}

// Bind
func (b *ProductsGetProductsListBinder) Bind(i interface{}, ctx echo.Context) error {
	limit := int32(b.LimitDefault)
	offset := int32(b.OffsetDefault)

	params := ctx.QueryParams()

	if v, ok := params[RequestParameterLimit]; ok {
		i, err := strconv.ParseInt(v[0], 10, 32)
		if err != nil {
			return err
		}
		limit = int32(i)
	}

	if v, ok := params[RequestParameterOffset]; ok {
		i, err := strconv.ParseInt(v[0], 10, 32)
		if err != nil {
			return err
		}
		offset = int32(i)
	}

	structure := i.(*grpc.ListProductsRequest)
	structure.Limit = limit
	structure.Offset = offset

	if v, ok := params[RequestParameterName]; ok {
		if v[0] != "" {
			structure.Name = v[0]
		}
	}

	if v, ok := params[RequestParameterSku]; ok {
		if v[0] != "" {
			structure.Sku = v[0]
		}
	}

	if v, ok := params[RequestParameterProjectId]; ok {
		if v[0] != "" {
			structure.ProjectId = v[0]
		}
	}

	return nil
}

// Bind
func (b *ProductsCreateProductBinder) Bind(i interface{}, ctx echo.Context) error {
	db := new(echo.DefaultBinder)
	err := db.Bind(i, ctx)
	if err != nil {
		return err
	}

	structure := i.(*grpc.Product)
	structure.Id = ""

	return nil
}

// Bind
func (b *ProductsUpdateProductBinder) Bind(i interface{}, ctx echo.Context) error {
	id := ctx.Param(RequestParameterId)
	if id == "" || bson.IsObjectIdHex(id) == false {
		return ErrorIncorrectProductId
	}
	db := new(echo.DefaultBinder)
	err := db.Bind(i, ctx)
	if err != nil {
		return err
	}

	structure := i.(*grpc.Product)
	structure.Id = id

	return nil
}

// Bind
func (b *ChangeMerchantDataRequestBinder) Bind(i interface{}, ctx echo.Context) error {
	req := make(map[string]interface{})

	db := new(echo.DefaultBinder)
	err := db.Bind(&req, ctx)

	if err != nil {
		return ErrorRequestParamsIncorrect
	}

	merchantId := ctx.Param(RequestParameterId)

	if merchantId == "" || bson.IsObjectIdHex(merchantId) == false {
		return ErrorIncorrectMerchantId
	}

	mReq := &grpc.GetMerchantByRequest{MerchantId: merchantId}
	mRsp, err := b.dispatch.Services.Billing.GetMerchantBy(context.Background(), mReq)

	if err != nil {
		b.L().Error(`Call billing server method "GetMerchantBy" failed`, logger.Args("error", err.Error(), "request", mReq))
		return ErrorUnknown
	}

	if mRsp.Status != pkg.ResponseStatusOk {
		return mRsp.Message
	}

	structure := i.(*grpc.ChangeMerchantDataRequest)
	structure.MerchantId = merchantId
	structure.HasMerchantSignature = mRsp.Item.HasMerchantSignature
	structure.HasPspSignature = mRsp.Item.HasPspSignature

	if v, ok := req[RequestParameterHasMerchantSignature]; ok {
		if tv, ok := v.(bool); !ok {
			return ErrorMessageHasMerchantSignatureIncorrectType
		} else {
			structure.HasMerchantSignature = tv
		}
	}

	if v, ok := req[RequestParameterHasPspSignature]; ok {
		if tv, ok := v.(bool); !ok {
			return ErrorMessageHasPspSignatureIncorrectType
		} else {
			structure.HasPspSignature = tv
		}
	}

	return nil
}

// Bind
func (b *ChangeProjectRequestBinder) Bind(i interface{}, ctx echo.Context) error {
	req := make(map[string]interface{})

	// Read the content
	var bodyBytes []byte
	if ctx.Request().Body != nil {
		bodyBytes, _ = ioutil.ReadAll(ctx.Request().Body)
	}

	// Restore the io.ReadCloser to its original state
	ctx.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	projectReq := &billing.Project{}
	if err := ctx.Bind(projectReq); err != nil {
		return ErrorRequestParamsIncorrect
	}

	// Restore the io.ReadCloser to its original state
	ctx.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	db := new(echo.DefaultBinder)
	err := db.Bind(&req, ctx)

	if err != nil {
		return ErrorRequestParamsIncorrect
	}

	projectId := ctx.Param(RequestParameterId)

	if projectId == "" || bson.IsObjectIdHex(projectId) == false {
		return ErrorIncorrectProjectId
	}

	pReq := &grpc.GetProjectRequest{ProjectId: projectId}
	pRsp, err := b.dispatch.Services.Billing.GetProject(context.Background(), pReq)

	if err != nil {
		b.L().Error(`Call billing server method "GetProject" failed`, logger.Args("error", err.Error(), "request", pReq))
		return ErrorUnknown
	}

	if pRsp.Status != pkg.ResponseStatusOk {
		return pRsp.Message
	}

	structure := i.(*billing.Project)
	structure.Id = projectId
	structure.MerchantId = pRsp.Item.MerchantId
	structure.Name = pRsp.Item.Name
	structure.CallbackCurrency = pRsp.Item.CallbackCurrency
	structure.CallbackProtocol = pRsp.Item.CallbackProtocol
	structure.CreateOrderAllowedUrls = pRsp.Item.CreateOrderAllowedUrls
	structure.AllowDynamicNotifyUrls = pRsp.Item.AllowDynamicNotifyUrls
	structure.AllowDynamicRedirectUrls = pRsp.Item.AllowDynamicRedirectUrls
	structure.LimitsCurrency = pRsp.Item.LimitsCurrency
	structure.MinPaymentAmount = pRsp.Item.MinPaymentAmount
	structure.MaxPaymentAmount = pRsp.Item.MaxPaymentAmount
	structure.NotifyEmails = pRsp.Item.NotifyEmails
	structure.IsProductsCheckout = pRsp.Item.IsProductsCheckout
	structure.SecretKey = pRsp.Item.SecretKey
	structure.SignatureRequired = pRsp.Item.SignatureRequired
	structure.SendNotifyEmail = pRsp.Item.SendNotifyEmail
	structure.UrlCheckAccount = pRsp.Item.UrlCheckAccount
	structure.UrlProcessPayment = pRsp.Item.UrlProcessPayment
	structure.UrlRedirectFail = pRsp.Item.UrlRedirectFail
	structure.UrlRedirectSuccess = pRsp.Item.UrlRedirectSuccess
	structure.Status = pRsp.Item.Status
	structure.ShortDescription = pRsp.Item.ShortDescription
	structure.Cover = pRsp.Item.Cover
	structure.FullDescription = pRsp.Item.FullDescription
	structure.Localizations = pRsp.Item.Localizations
	structure.Currencies = pRsp.Item.Currencies
	structure.VirtualCurrency = pRsp.Item.VirtualCurrency

	if v, ok := req[RequestParameterName]; ok {
		tv, ok := v.(map[string]interface{})

		if !ok || len(tv) <= 0 {
			return ErrorMessageNameIncorrectType
		}

		for k, tvv := range tv {
			structure.Name[k] = tvv.(string)
		}
	}

	if v, ok := req[RequestParameterCallbackCurrency]; ok {
		if tv, ok := v.(string); !ok {
			return ErrorMessageCallbackCurrencyIncorrectType
		} else {
			structure.CallbackCurrency = tv
		}
	}

	if v, ok := req[RequestParameterCallbackProtocol]; ok {
		if tv, ok := v.(string); !ok {
			return ErrorMessageCallbackProtocolIncorrectType
		} else {
			structure.CallbackProtocol = tv
		}
	}

	if v, ok := req[RequestParameterCreateOrderAllowedUrls]; ok {
		tv, ok := v.([]interface{})

		if !ok {
			return ErrorMessageCreateOrderAllowedUrlsIncorrectType
		}

		structure.CreateOrderAllowedUrls = []string{}

		for _, tvv := range tv {
			structure.CreateOrderAllowedUrls = append(structure.CreateOrderAllowedUrls, tvv.(string))
		}
	}

	if v, ok := req[RequestParameterAllowDynamicNotifyUrls]; ok {
		if tv, ok := v.(bool); !ok {
			return ErrorMessageAllowDynamicNotifyUrlsIncorrectType
		} else {
			structure.AllowDynamicNotifyUrls = tv
		}
	}

	if v, ok := req[RequestParameterAllowDynamicRedirectUrls]; ok {
		if tv, ok := v.(bool); !ok {
			return ErrorMessageAllowDynamicRedirectUrlsIncorrectType
		} else {
			structure.AllowDynamicRedirectUrls = tv
		}
	}

	if v, ok := req[RequestParameterLimitsCurrency]; ok {
		if tv, ok := v.(string); !ok {
			return ErrorMessageLimitsCurrencyIncorrectType
		} else {
			structure.LimitsCurrency = tv
		}
	}

	if v, ok := req[RequestParameterMinPaymentAmount]; ok {
		if tv, ok := v.(float64); !ok {
			return ErrorMessageMinPaymentAmountIncorrectType
		} else {
			structure.MinPaymentAmount = tv
		}
	}

	if v, ok := req[RequestParameterMaxPaymentAmount]; ok {
		if tv, ok := v.(float64); !ok {
			return ErrorMessageMaxPaymentAmountIncorrectType
		} else {
			structure.MaxPaymentAmount = tv
		}
	}

	if v, ok := req[RequestParameterNotifyEmails]; ok {
		tv, ok := v.([]interface{})

		if !ok {
			return ErrorMessageNotifyEmailsIncorrectType
		}

		structure.NotifyEmails = []string{}

		for _, tvv := range tv {
			structure.NotifyEmails = append(structure.NotifyEmails, tvv.(string))
		}
	}

	if v, ok := req[RequestParameterIsProductsCheckout]; ok {
		if tv, ok := v.(bool); !ok {
			return ErrorMessageIsProductsCheckoutIncorrectType
		} else {
			structure.IsProductsCheckout = tv
		}
	}

	if v, ok := req[RequestParameterSecretKey]; ok {
		if tv, ok := v.(string); !ok {
			return ErrorMessageSecretKeyIncorrectType
		} else {
			structure.SecretKey = tv
		}
	}

	if v, ok := req[RequestParameterSignatureRequired]; ok {
		if tv, ok := v.(bool); !ok {
			return ErrorMessageSignatureRequiredIncorrectType
		} else {
			structure.SignatureRequired = tv
		}
	}

	if v, ok := req[RequestParameterSendNotifyEmail]; ok {
		if tv, ok := v.(bool); !ok {
			return ErrorMessageSendNotifyEmailIncorrectType
		} else {
			structure.SendNotifyEmail = tv
		}
	}

	if v, ok := req[RequestParameterUrlCheckAccount]; ok {
		if tv, ok := v.(string); !ok {
			return ErrorMessageUrlCheckAccountIncorrectType
		} else {
			structure.UrlCheckAccount = tv
		}
	}

	if v, ok := req[RequestParameterUrlProcessPayment]; ok {
		if tv, ok := v.(string); !ok {
			return ErrorMessageUrlProcessPaymentIncorrectType
		} else {
			structure.UrlProcessPayment = tv
		}
	}

	if v, ok := req[RequestParameterUrlRedirectFail]; ok {
		if tv, ok := v.(string); !ok {
			return ErrorMessageUrlRedirectFailIncorrectType
		} else {
			structure.UrlRedirectFail = tv
		}
	}

	if v, ok := req[RequestParameterUrlRedirectSuccess]; ok {
		if tv, ok := v.(string); !ok {
			return ErrorMessageUrlRedirectSuccessIncorrectType
		} else {
			structure.UrlRedirectSuccess = tv
		}
	}

	if v, ok := req[RequestParameterStatus]; ok {
		if tv, ok := v.(float64); !ok {
			return ErrorMessageStatusIncorrectType
		} else {
			structure.Status = int32(tv)
		}
	}

	if v, ok := req[RequestParameterUrlChargebackPayment]; ok {
		if tv, ok := v.(string); !ok {
			return ErrorMessageUrlChargebackPayment
		} else {
			structure.UrlChargebackPayment = tv
		}
	}

	if v, ok := req[RequestParameterUrlCancelPayment]; ok {
		if tv, ok := v.(string); !ok {
			return ErrorMessageUrlCancelPayment
		} else {
			structure.UrlCancelPayment = tv
		}
	}

	if v, ok := req[RequestParameterUrlFraudPayment]; ok {
		if tv, ok := v.(string); !ok {
			return ErrorMessageUrlFraudPayment
		} else {
			structure.UrlFraudPayment = tv
		}
	}

	if v, ok := req[RequestParameterUrlRefundPayment]; ok {
		if tv, ok := v.(string); !ok {
			return ErrorMessageUrlRefundPayment
		} else {
			structure.UrlRefundPayment = tv
		}
	}

	if _, ok := req[RequestParameterFullDescription]; ok {
		structure.FullDescription = projectReq.FullDescription
	}

	if _, ok := req[RequestParameterShortDescription]; ok {
		structure.ShortDescription = projectReq.ShortDescription
	}

	if _, ok := req[RequestParameterCover]; ok {
		structure.Cover = projectReq.Cover
	}

	if _, ok := req[RequestParameterLocalizations]; ok {
		structure.Localizations = projectReq.Localizations
	}

	if _, ok := req[RequestParameterCurrencies]; ok {
		structure.Currencies = projectReq.Currencies
	}

	if _, ok := req[RequestParameterVirtualCurrency]; ok {
		structure.VirtualCurrency = projectReq.VirtualCurrency
	}

	return nil
}
