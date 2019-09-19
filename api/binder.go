package api

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/labstack/echo/v4"
	"github.com/paysuper/paysuper-billing-server/pkg"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/billing"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/grpc"
	"github.com/paysuper/paysuper-payment-link/proto"
	"go.uber.org/zap"
	"io/ioutil"
	"strconv"
)

const (
	binderErrorQueryParamsIsEmpty   = "required params not found"
	binderErrorPeriodIsRequire      = "period is required field"
	binderErrorFromIsRequire        = "date from is required field"
	binderErrorToIsRequire          = "date to is required field"
	binderErrorUnknownRevenuePeriod = "unknown revenue period"
)

type OrderFormBinder struct{}
type OrderJsonBinder struct{}
type OrderRevenueDynamicRequestBinder struct{}
type OrderAccountingPaymentRequestBinder struct{}
type PaymentCreateProcessBinder struct{}
type OnboardingMerchantListingBinder struct{}
type OnboardingChangeMerchantStatusBinder struct{}
type OnboardingNotificationsListBinder struct{}
type OnboardingGetPaymentMethodBinder struct{}
type OnboardingChangePaymentMethodBinder struct{}
type OnboardingCreateNotificationBinder struct{}
type PaylinksListBinder struct{}
type PaylinksUrlBinder struct{}
type PaylinksCreateBinder struct{}
type PaylinksUpdateBinder struct{}
type ProductsGetProductsListBinder struct{}
type ProductsCreateProductBinder struct{}
type ProductsUpdateProductBinder struct{}
type ChangeMerchantDataRequestBinder struct {
	*Api
}
type ChangeProjectRequestBinder struct {
	*Api
}
type ProtobufBinder struct{}

func (m *ProtobufBinder) Bind(i interface{}, ctx echo.Context) error {
	st, ok := i.(proto.Message)

	if !ok {
		return errors.New("some error")
	}

	err := jsonpb.Unmarshal(ctx.Request().Body, st)

	if err != nil {
		return err
	}

	return nil
}

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

func (cb *OnboardingMerchantListingBinder) Bind(i interface{}, ctx echo.Context) (err error) {
	db := new(echo.DefaultBinder)
	err = db.Bind(i, ctx)

	if err != nil {
		return err
	}

	params := ctx.QueryParams()
	structure := i.(*grpc.MerchantListingRequest)

	if structure.Limit <= 0 {
		structure.Limit = LimitDefault
	}

	if v, ok := params[requestParameterIsSigned]; ok {
		if v[0] == "0" || v[0] == "false" {
			structure.IsSigned = 2
		} else {
			if v[0] == "1" || v[0] == "true" {
				structure.IsSigned = 2
			} else {
				return errorRequestParamsIncorrect
			}
		}
	}

	return
}

func (cb *OnboardingNotificationsListBinder) Bind(i interface{}, ctx echo.Context) error {
	db := new(echo.DefaultBinder)
	err := db.Bind(i, ctx)

	if err != nil {
		return err
	}

	params := ctx.QueryParams()
	structure := i.(*grpc.ListingNotificationRequest)
	structure.MerchantId = ctx.Param(requestParameterMerchantId)

	if structure.Limit <= 0 {
		structure.Limit = LimitDefault
	}

	if v, ok := params[requestParameterIsSystem]; ok {
		if v[0] == "0" || v[0] == "false" {
			structure.IsSystem = 1
		} else {
			structure.IsSystem = 2
		}
	}

	return nil
}

func (cb *OnboardingGetPaymentMethodBinder) Bind(i interface{}, ctx echo.Context) error {
	merchantId := ctx.Param(requestParameterMerchantId)
	paymentMethodId := ctx.Param(requestParameterPaymentMethodId)

	if merchantId == "" || bson.IsObjectIdHex(merchantId) == false {
		return errorIncorrectMerchantId
	}

	if paymentMethodId == "" || bson.IsObjectIdHex(paymentMethodId) == false {
		return errorIncorrectPaymentMethodId
	}

	structure := i.(*grpc.GetMerchantPaymentMethodRequest)
	structure.MerchantId = merchantId
	structure.PaymentMethodId = paymentMethodId

	return nil
}

func (cb *OnboardingChangePaymentMethodBinder) Bind(i interface{}, ctx echo.Context) error {
	db := new(echo.DefaultBinder)
	err := db.Bind(i, ctx)

	if err != nil {
		return err
	}

	structure := i.(*grpc.MerchantPaymentMethodRequest)
	merchantId := ctx.Param(requestParameterMerchantId)
	methodId := ctx.Param(requestParameterPaymentMethodId)

	if merchantId == "" || bson.IsObjectIdHex(merchantId) == false {
		return errorIncorrectMerchantId
	}

	if methodId == "" || bson.IsObjectIdHex(methodId) == false ||
		structure.PaymentMethod.Id != methodId {
		return errorIncorrectPaymentMethodId
	}

	structure.MerchantId = merchantId

	return nil
}

func (b *OnboardingChangeMerchantStatusBinder) Bind(i interface{}, ctx echo.Context) error {
	db := new(echo.DefaultBinder)
	err := db.Bind(i, ctx)

	if err != nil {
		return err
	}

	merchantId := ctx.Param(requestParameterId)

	if merchantId == "" || bson.IsObjectIdHex(merchantId) == false {
		return errorIncorrectMerchantId
	}

	structure := i.(*grpc.MerchantChangeStatusRequest)
	structure.MerchantId = merchantId

	return nil
}

func (b *OnboardingCreateNotificationBinder) Bind(i interface{}, ctx echo.Context) error {
	db := new(echo.DefaultBinder)
	err := db.Bind(i, ctx)

	if err != nil {
		return err
	}

	merchantId := ctx.Param(requestParameterMerchantId)

	if merchantId == "" || bson.IsObjectIdHex(merchantId) == false {
		return errorIncorrectMerchantId
	}

	structure := i.(*grpc.NotificationRequest)
	structure.MerchantId = merchantId

	return nil
}

func (b *PaylinksListBinder) Bind(i interface{}, ctx echo.Context) error {
	params := ctx.QueryParams()
	structure := i.(*paylink.GetPaylinksRequest)

	structure.Limit = LimitDefault
	structure.Offset = OffsetDefault

	if v, ok := params[requestParameterLimit]; ok {
		i, err := strconv.ParseInt(v[0], 10, 32)
		if err != nil {
			return err
		}
		structure.Limit = uint32(i)
	}

	if v, ok := params[requestParameterOffset]; ok {
		i, err := strconv.ParseInt(v[0], 10, 32)
		if err != nil {
			return err
		}
		structure.Offset = uint32(i)
	}

	structure.ProjectId = ctx.Param(requestParameterProjectId)

	return nil
}

func (b *PaylinksUrlBinder) Bind(i interface{}, ctx echo.Context) error {
	params := ctx.QueryParams()
	structure := i.(*paylink.GetPaylinkURLRequest)

	id := ctx.Param(requestParameterId)
	structure.Id = id

	if v, ok := params[requestParameterUtmSource]; ok {
		structure.UtmSource = v[0]
	}

	if v, ok := params[requestParameterUtmMedium]; ok {
		structure.UtmMedium = v[0]
	}

	if v, ok := params[requestParameterUtmCampaign]; ok {
		structure.UtmCampaign = v[0]
	}

	return nil
}

func (b *PaylinksCreateBinder) Bind(i interface{}, ctx echo.Context) error {
	db := new(echo.DefaultBinder)
	err := db.Bind(i, ctx)
	if err != nil {
		return err
	}

	structure := i.(*paylink.CreatePaylinkRequest)
	structure.Id = ""

	return nil
}

func (b *PaylinksUpdateBinder) Bind(i interface{}, ctx echo.Context) error {
	id := ctx.Param(requestParameterId)
	if id == "" {
		return errorIncorrectPaylinkId
	}
	db := new(echo.DefaultBinder)
	err := db.Bind(i, ctx)
	if err != nil {
		return err
	}

	structure := i.(*paylink.CreatePaylinkRequest)
	structure.Id = id

	return nil
}

func (b *ProductsGetProductsListBinder) Bind(i interface{}, ctx echo.Context) error {
	limit := int32(LimitDefault)
	offset := int32(OffsetDefault)

	params := ctx.QueryParams()

	if v, ok := params[requestParameterLimit]; ok {
		i, err := strconv.ParseInt(v[0], 10, 32)
		if err != nil {
			return err
		}
		limit = int32(i)
	}

	if v, ok := params[requestParameterOffset]; ok {
		i, err := strconv.ParseInt(v[0], 10, 32)
		if err != nil {
			return err
		}
		offset = int32(i)
	}

	structure := i.(*grpc.ListProductsRequest)
	structure.Limit = limit
	structure.Offset = offset

	if v, ok := params[requestParameterName]; ok {
		if v[0] != "" {
			structure.Name = v[0]
		}
	}

	if v, ok := params[requestParameterSku]; ok {
		if v[0] != "" {
			structure.Sku = v[0]
		}
	}

	if v, ok := params[requestParameterProjectId]; ok {
		if v[0] != "" {
			structure.ProjectId = v[0]
		}
	}

	return nil
}

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

func (b *ProductsUpdateProductBinder) Bind(i interface{}, ctx echo.Context) error {
	id := ctx.Param(requestParameterId)
	if id == "" || bson.IsObjectIdHex(id) == false {
		return errorIncorrectProductId
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

func (b *ChangeMerchantDataRequestBinder) Bind(i interface{}, ctx echo.Context) error {
	req := make(map[string]interface{})

	db := new(echo.DefaultBinder)
	err := db.Bind(&req, ctx)

	if err != nil {
		return errorRequestParamsIncorrect
	}

	merchantId := ctx.Param(requestParameterId)

	if merchantId == "" || bson.IsObjectIdHex(merchantId) == false {
		return errorIncorrectMerchantId
	}

	mReq := &grpc.GetMerchantByRequest{MerchantId: merchantId}
	mRsp, err := b.billingService.GetMerchantBy(context.Background(), mReq)

	if err != nil {
		zap.S().Errorf(`Call billing server method "GetMerchantBy" failed`, "error", err.Error(), "request", mReq)
		return errorUnknown
	}

	if mRsp.Status != pkg.ResponseStatusOk {
		return mRsp.Message
	}

	structure := i.(*grpc.ChangeMerchantDataRequest)
	structure.MerchantId = merchantId
	structure.AgreementType = mRsp.Item.AgreementType
	structure.HasMerchantSignature = mRsp.Item.HasMerchantSignature
	structure.HasPspSignature = mRsp.Item.HasPspSignature
	structure.AgreementSentViaMail = mRsp.Item.AgreementSentViaMail
	structure.MailTrackingLink = mRsp.Item.MailTrackingLink

	if v, ok := req[requestParameterAgreementType]; ok {
		if tv, ok := v.(float64); !ok {
			return errorMessageAgreementTypeIncorrectType
		} else {
			structure.AgreementType = int32(tv)
		}
	}

	if v, ok := req[requestParameterHasMerchantSignature]; ok {
		if tv, ok := v.(bool); !ok {
			return errorMessageHasMerchantSignatureIncorrectType
		} else {
			structure.HasMerchantSignature = tv
		}
	}

	if v, ok := req[requestParameterHasPspSignature]; ok {
		if tv, ok := v.(bool); !ok {
			return errorMessageHasPspSignatureIncorrectType
		} else {
			structure.HasPspSignature = tv
		}
	}

	if v, ok := req[requestParameterAgreementSentViaMail]; ok {
		if tv, ok := v.(bool); !ok {
			return errorMessageAgreementSentViaMailIncorrectType
		} else {
			structure.AgreementSentViaMail = tv
		}
	}

	if v, ok := req[requestParameterMailTrackingLink]; ok {
		if tv, ok := v.(string); !ok {
			return errorMessageMailTrackingLinkIncorrectType
		} else {
			structure.MailTrackingLink = tv
		}
	}

	return nil
}

func (b *ChangeProjectRequestBinder) Bind(i interface{}, ctx echo.Context) error {
	req := make(map[string]interface{})

	db := new(echo.DefaultBinder)
	err := db.Bind(&req, ctx)

	if err != nil {
		return errorRequestParamsIncorrect
	}

	projectId := ctx.Param(requestParameterId)

	if projectId == "" || bson.IsObjectIdHex(projectId) == false {
		return errorIncorrectProjectId
	}

	pReq := &grpc.GetProjectRequest{ProjectId: projectId}
	pRsp, err := b.billingService.GetProject(context.Background(), pReq)

	if err != nil {
		zap.S().Errorf(`Call billing server method "GetProject" failed`, "error", err.Error(), "request", pReq)
		return errorUnknown
	}

	if pRsp.Status != pkg.ResponseStatusOk {
		return pRsp.Message
	}

	structure := i.(*billing.Project)
	structure.Id = projectId
	structure.MerchantId = pRsp.Item.MerchantId
	structure.Name = pRsp.Item.Name
	structure.Image = pRsp.Item.Image
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

	if v, ok := req[requestParameterName]; ok {
		tv, ok := v.(map[string]interface{})

		if !ok || len(tv) <= 0 {
			return errorMessageNameIncorrectType
		}

		for k, tvv := range tv {
			structure.Name[k] = tvv.(string)
		}
	}

	if v, ok := req[requestParameterImage]; ok {
		if tv, ok := v.(string); !ok {
			return errorMessageImageIncorrectType
		} else {
			structure.Image = tv
		}
	}

	if v, ok := req[requestParameterCallbackCurrency]; ok {
		if tv, ok := v.(string); !ok {
			return errorMessageCallbackCurrencyIncorrectType
		} else {
			structure.CallbackCurrency = tv
		}
	}

	if v, ok := req[requestParameterCallbackProtocol]; ok {
		if tv, ok := v.(string); !ok {
			return errorMessageCallbackProtocolIncorrectType
		} else {
			structure.CallbackProtocol = tv
		}
	}

	if v, ok := req[requestParameterCreateOrderAllowedUrls]; ok {
		tv, ok := v.([]interface{})

		if !ok {
			return errorMessageCreateOrderAllowedUrlsIncorrectType
		}

		structure.CreateOrderAllowedUrls = []string{}

		for _, tvv := range tv {
			structure.CreateOrderAllowedUrls = append(structure.CreateOrderAllowedUrls, tvv.(string))
		}
	}

	if v, ok := req[requestParameterAllowDynamicNotifyUrls]; ok {
		if tv, ok := v.(bool); !ok {
			return errorMessageAllowDynamicNotifyUrlsIncorrectType
		} else {
			structure.AllowDynamicNotifyUrls = tv
		}
	}

	if v, ok := req[requestParameterAllowDynamicRedirectUrls]; ok {
		if tv, ok := v.(bool); !ok {
			return errorMessageAllowDynamicRedirectUrlsIncorrectType
		} else {
			structure.AllowDynamicRedirectUrls = tv
		}
	}

	if v, ok := req[requestParameterLimitsCurrency]; ok {
		if tv, ok := v.(string); !ok {
			return errorMessageLimitsCurrencyIncorrectType
		} else {
			structure.LimitsCurrency = tv
		}
	}

	if v, ok := req[requestParameterMinPaymentAmount]; ok {
		if tv, ok := v.(float64); !ok {
			return errorMessageMinPaymentAmountIncorrectType
		} else {
			structure.MinPaymentAmount = tv
		}
	}

	if v, ok := req[requestParameterMaxPaymentAmount]; ok {
		if tv, ok := v.(float64); !ok {
			return errorMessageMaxPaymentAmountIncorrectType
		} else {
			structure.MaxPaymentAmount = tv
		}
	}

	if v, ok := req[requestParameterNotifyEmails]; ok {
		tv, ok := v.([]interface{})

		if !ok {
			return errorMessageNotifyEmailsIncorrectType
		}

		structure.NotifyEmails = []string{}

		for _, tvv := range tv {
			structure.NotifyEmails = append(structure.NotifyEmails, tvv.(string))
		}
	}

	if v, ok := req[requestParameterIsProductsCheckout]; ok {
		if tv, ok := v.(bool); !ok {
			return errorMessageIsProductsCheckoutIncorrectType
		} else {
			structure.IsProductsCheckout = tv
		}
	}

	if v, ok := req[requestParameterSecretKey]; ok {
		if tv, ok := v.(string); !ok {
			return errorMessageSecretKeyIncorrectType
		} else {
			structure.SecretKey = tv
		}
	}

	if v, ok := req[requestParameterSignatureRequired]; ok {
		if tv, ok := v.(bool); !ok {
			return errorMessageSignatureRequiredIncorrectType
		} else {
			structure.SignatureRequired = tv
		}
	}

	if v, ok := req[requestParameterSendNotifyEmail]; ok {
		if tv, ok := v.(bool); !ok {
			return errorMessageSendNotifyEmailIncorrectType
		} else {
			structure.SendNotifyEmail = tv
		}
	}

	if v, ok := req[requestParameterUrlCheckAccount]; ok {
		if tv, ok := v.(string); !ok {
			return errorMessageUrlCheckAccountIncorrectType
		} else {
			structure.UrlCheckAccount = tv
		}
	}

	if v, ok := req[requestParameterUrlProcessPayment]; ok {
		if tv, ok := v.(string); !ok {
			return errorMessageUrlProcessPaymentIncorrectType
		} else {
			structure.UrlProcessPayment = tv
		}
	}

	if v, ok := req[requestParameterUrlRedirectFail]; ok {
		if tv, ok := v.(string); !ok {
			return errorMessageUrlRedirectFailIncorrectType
		} else {
			structure.UrlRedirectFail = tv
		}
	}

	if v, ok := req[requestParameterUrlRedirectSuccess]; ok {
		if tv, ok := v.(string); !ok {
			return errorMessageUrlRedirectSuccessIncorrectType
		} else {
			structure.UrlRedirectSuccess = tv
		}
	}

	if v, ok := req[requestParameterStatus]; ok {
		if tv, ok := v.(float64); !ok {
			return errorMessageStatusIncorrectType
		} else {
			structure.Status = int32(tv)
		}
	}

	if v, ok := req[requestParameterUrlChargebackPayment]; ok {
		if tv, ok := v.(string); !ok {
			return errorMessageUrlChargebackPayment
		} else {
			structure.UrlChargebackPayment = tv
		}
	}

	if v, ok := req[requestParameterUrlCancelPayment]; ok {
		if tv, ok := v.(string); !ok {
			return errorMessageUrlCancelPayment
		} else {
			structure.UrlCancelPayment = tv
		}
	}

	if v, ok := req[requestParameterUrlFraudPayment]; ok {
		if tv, ok := v.(string); !ok {
			return errorMessageUrlFraudPayment
		} else {
			structure.UrlFraudPayment = tv
		}
	}

	if v, ok := req[requestParameterUrlRefundPayment]; ok {
		if tv, ok := v.(string); !ok {
			return errorMessageUrlRefundPayment
		} else {
			structure.UrlRefundPayment = tv
		}
	}

	return nil
}
