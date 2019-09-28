package handlers

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo/v4"
	"github.com/paysuper/paysuper-billing-server/pkg"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/billing"
	"github.com/paysuper/paysuper-management-api/internal/dispatcher/common"
	"github.com/paysuper/paysuper-management-api/internal/mock"
	"github.com/paysuper/paysuper-management-api/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"strings"
	"testing"
	"time"
)

type CardPayTestSuite struct {
	suite.Suite
	router *CardPayWebHook
	caller *test.EchoReqResCaller
}

func Test_CardPayTestSuite(t *testing.T) {
	suite.Run(t, new(CardPayTestSuite))
}

func (suite *CardPayTestSuite) SetupTest() {
	var e error
	settings := test.DefaultSettings()
	srv := common.Services{
		Billing: mock.NewBillingServerOkMock(),
	}
	suite.caller, e = test.SetUp(settings, srv, func(set *test.TestSet, mw test.Middleware) common.Handlers {
		suite.router = NewCardPayWebHook(set.HandlerSet, set.GlobalConfig)
		return common.Handlers{
			suite.router,
		}
	})
	if e != nil {
		panic(e)
	}
}

func (suite *CardPayTestSuite) TearDownTest() {}

func (suite *CardPayTestSuite) TestCardPay_RefundCallback_Ok() {

	refundReq := &billing.CardPayRefundCallback{
		MerchantOrder: &billing.CardPayMerchantOrder{
			Id: bson.NewObjectId().Hex(),
		},
		PaymentMethod: "BANKCARD",
		PaymentData: &billing.CardPayRefundCallbackPaymentData{
			Id:              bson.NewObjectId().Hex(),
			RemainingAmount: 0,
		},
		RefundData: &billing.CardPayRefundCallbackRefundData{
			Amount:   100,
			Created:  time.Now().Format("2006-01-02T15:04:05Z"),
			Id:       bson.NewObjectId().Hex(),
			Currency: "RUB",
			Status:   pkg.CardPayPaymentResponseStatusCompleted,
			AuthCode: bson.NewObjectId().Hex(),
			Is_3D:    true,
			Rrn:      bson.NewObjectId().Hex(),
		},
		CallbackTime: time.Now().Format("2006-01-02T15:04:05Z"),
		Customer: &billing.CardPayCustomer{
			Email: "test@unut.test",
			Id:    "test@unut.test",
		},
	}

	b, err := json.Marshal(refundReq)
	assert.NoError(suite.T(), err)

	hash := sha512.New()
	hash.Write([]byte(string(b) + "secret_key"))

	path := common.WebHookGroupPath + cardPayWebHookRefundNotifyPath
	res, err := suite.caller.Request(http.MethodPost, path, bytes.NewReader(b), func(request *http.Request, middleware test.Middleware) {
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		request.Header.Set(common.CardPayPaymentResponseHeaderSignature, hex.EncodeToString(hash.Sum(nil)))
	})
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), http.StatusOK, res.Code)
	assert.Empty(suite.T(), res.Body.String())
}

func (suite *CardPayTestSuite) TestCardPay_RefundCallback_BindError() {

	refundReq := `{"payment_method": 11111}`
	hash := sha512.New()
	hash.Write([]byte(refundReq + "secret_key"))

	path := common.WebHookGroupPath + cardPayWebHookRefundNotifyPath
	_, err := suite.caller.Request(http.MethodPost, path, strings.NewReader(refundReq), func(request *http.Request, middleware test.Middleware) {
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		request.Header.Set(common.CardPayPaymentResponseHeaderSignature, hex.EncodeToString(hash.Sum(nil)))
	})
	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)

	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorRequestParamsIncorrect, httpErr.Message)
}

func (suite *CardPayTestSuite) TestCardPay_RefundCallback_ValidationError() {
	refundReq := &billing.CardPayRefundCallback{
		MerchantOrder: &billing.CardPayMerchantOrder{
			Id: bson.NewObjectId().Hex(),
		},
		PaymentMethod: "BANKCARD",
		RefundData: &billing.CardPayRefundCallbackRefundData{
			Amount:   100,
			Created:  time.Now().Format("2006-01-02T15:04:05Z"),
			Id:       bson.NewObjectId().Hex(),
			Currency: "RUB",
			Status:   pkg.CardPayPaymentResponseStatusCompleted,
			AuthCode: bson.NewObjectId().Hex(),
			Is_3D:    true,
			Rrn:      bson.NewObjectId().Hex(),
		},
		CallbackTime: time.Now().Format("2006-01-02T15:04:05Z"),
		Customer: &billing.CardPayCustomer{
			Email: "test@unut.test",
			Id:    "test@unut.test",
		},
	}

	b, err := json.Marshal(refundReq)
	assert.NoError(suite.T(), err)

	hash := sha512.New()
	hash.Write([]byte(string(b) + "secret_key"))

	path := common.WebHookGroupPath + cardPayWebHookRefundNotifyPath
	_, err = suite.caller.Request(http.MethodPost, path, bytes.NewReader(b), func(request *http.Request, middleware test.Middleware) {
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		request.Header.Set(common.CardPayPaymentResponseHeaderSignature, hex.EncodeToString(hash.Sum(nil)))
	})
	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Regexp(suite.T(), common.NewValidationError("PaymentData"), httpErr.Message)
}

func (suite *CardPayTestSuite) TestCardPay_RefundCallback_BillingServerSystemError() {

	refundReq := &billing.CardPayRefundCallback{
		MerchantOrder: &billing.CardPayMerchantOrder{
			Id: bson.NewObjectId().Hex(),
		},
		PaymentMethod: "BANKCARD",
		PaymentData: &billing.CardPayRefundCallbackPaymentData{
			Id:              bson.NewObjectId().Hex(),
			RemainingAmount: 0,
		},
		RefundData: &billing.CardPayRefundCallbackRefundData{
			Amount:   100,
			Created:  time.Now().Format("2006-01-02T15:04:05Z"),
			Id:       bson.NewObjectId().Hex(),
			Currency: "RUB",
			Status:   pkg.CardPayPaymentResponseStatusCompleted,
			AuthCode: bson.NewObjectId().Hex(),
			Is_3D:    true,
			Rrn:      bson.NewObjectId().Hex(),
		},
		CallbackTime: time.Now().Format("2006-01-02T15:04:05Z"),
		Customer: &billing.CardPayCustomer{
			Email: "test@unut.test",
			Id:    "test@unut.test",
		},
	}

	b, err := json.Marshal(refundReq)
	assert.NoError(suite.T(), err)

	hash := sha512.New()
	hash.Write([]byte(string(b) + "secret_key"))

	suite.router.dispatch.Services.Billing = mock.NewBillingServerSystemErrorMock()

	path := common.WebHookGroupPath + cardPayWebHookRefundNotifyPath
	_, err = suite.caller.Request(http.MethodPost, path, bytes.NewReader(b), func(request *http.Request, middleware test.Middleware) {
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		request.Header.Set(common.CardPayPaymentResponseHeaderSignature, hex.EncodeToString(hash.Sum(nil)))
	})
	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusInternalServerError, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorUnknown, httpErr.Message)
}

func (suite *CardPayTestSuite) TestCardPay_RefundCallback_BillingServer_Error() {
	refundReq := &billing.CardPayRefundCallback{
		MerchantOrder: &billing.CardPayMerchantOrder{
			Id: bson.NewObjectId().Hex(),
		},
		PaymentMethod: "BANKCARD",
		PaymentData: &billing.CardPayRefundCallbackPaymentData{
			Id:              bson.NewObjectId().Hex(),
			RemainingAmount: 0,
		},
		RefundData: &billing.CardPayRefundCallbackRefundData{
			Amount:   100,
			Created:  time.Now().Format("2006-01-02T15:04:05Z"),
			Id:       bson.NewObjectId().Hex(),
			Currency: "RUB",
			Status:   pkg.CardPayPaymentResponseStatusCompleted,
			AuthCode: bson.NewObjectId().Hex(),
			Is_3D:    true,
			Rrn:      bson.NewObjectId().Hex(),
		},
		CallbackTime: time.Now().Format("2006-01-02T15:04:05Z"),
		Customer: &billing.CardPayCustomer{
			Email: "test@unut.test",
			Id:    "test@unut.test",
		},
	}

	b, err := json.Marshal(refundReq)
	assert.NoError(suite.T(), err)

	hash := sha512.New()
	hash.Write([]byte(string(b) + "secret_key"))

	suite.router.dispatch.Services.Billing = mock.NewBillingServerErrorMock()

	path := common.WebHookGroupPath + cardPayWebHookRefundNotifyPath
	_, err = suite.caller.Request(http.MethodPost, path, bytes.NewReader(b), func(request *http.Request, middleware test.Middleware) {
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		request.Header.Set(common.CardPayPaymentResponseHeaderSignature, hex.EncodeToString(hash.Sum(nil)))
	})
	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusNotFound, httpErr.Code)
	assert.Equal(suite.T(), mock.SomeError.Message, httpErr.Message)
}

func (suite *CardPayTestSuite) TestCardPay_RefundCallback_BillingServerTemporary_Ok() {
	refundReq := &billing.CardPayRefundCallback{
		MerchantOrder: &billing.CardPayMerchantOrder{
			Id: bson.NewObjectId().Hex(),
		},
		PaymentMethod: "BANKCARD",
		PaymentData: &billing.CardPayRefundCallbackPaymentData{
			Id:              bson.NewObjectId().Hex(),
			RemainingAmount: 0,
		},
		RefundData: &billing.CardPayRefundCallbackRefundData{
			Amount:   100,
			Created:  time.Now().Format("2006-01-02T15:04:05Z"),
			Id:       bson.NewObjectId().Hex(),
			Currency: "RUB",
			Status:   pkg.CardPayPaymentResponseStatusCompleted,
			AuthCode: bson.NewObjectId().Hex(),
			Is_3D:    true,
			Rrn:      bson.NewObjectId().Hex(),
		},
		CallbackTime: time.Now().Format("2006-01-02T15:04:05Z"),
		Customer: &billing.CardPayCustomer{
			Email: "test@unut.test",
			Id:    "test@unut.test",
		},
	}

	b, err := json.Marshal(refundReq)
	assert.NoError(suite.T(), err)

	hash := sha512.New()
	hash.Write([]byte(string(b) + "secret_key"))

	suite.router.dispatch.Services.Billing = mock.NewBillingServerOkTemporaryMock()

	path := common.WebHookGroupPath + cardPayWebHookRefundNotifyPath
	res, err := suite.caller.Request(http.MethodPost, path, bytes.NewReader(b), func(request *http.Request, middleware test.Middleware) {
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		request.Header.Set(common.CardPayPaymentResponseHeaderSignature, hex.EncodeToString(hash.Sum(nil)))
	})

	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), res.Body.String())

	var body map[string]string
	err = json.Unmarshal(res.Body.Bytes(), &body)
	assert.NoError(suite.T(), err)

	v, ok := body["message"]
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), mock.SomeError.Message, v)
}
