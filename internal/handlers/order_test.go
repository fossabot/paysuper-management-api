package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/paysuper/paysuper-billing-server/pkg"
	billMock "github.com/paysuper/paysuper-billing-server/pkg/mocks"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/billing"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/grpc"
	"github.com/paysuper/paysuper-management-api/internal/dispatcher/common"
	"github.com/paysuper/paysuper-management-api/internal/mock"
	"github.com/paysuper/paysuper-management-api/internal/test"
	"github.com/stretchr/testify/assert"
	mock2 "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/url"
	"testing"
	"time"
)

type OrderTestSuite struct {
	suite.Suite
	router *OrderRoute
	caller *test.EchoReqResCaller
}

func Test_Order(t *testing.T) {
	suite.Run(t, new(OrderTestSuite))
}

func (suite *OrderTestSuite) SetupTest() {
	user := &common.AuthUser{
		Id: "ffffffffffffffffffffffff",
	}

	var e error
	settings := test.DefaultSettings()
	srv := common.Services{
		Billing: mock.NewBillingServerOkMock(),
		PayLink: mock.NewPaymentLinkOkMock(),
	}
	suite.caller, e = test.SetUp(settings, srv, func(set *test.TestSet, mw test.Middleware) common.Handlers {
		mw.Pre(test.PreAuthUserMiddleware(user))
		suite.router = NewOrderRoute(set.HandlerSet, set.GlobalConfig)
		return common.Handlers{
			suite.router,
		}
	})
	if e != nil {
		panic(e)
	}
}

func (suite *OrderTestSuite) TearDownTest() {}

func (suite *OrderTestSuite) TestOrder_GetRefund_Ok() {

	res, err := suite.caller.Builder().
		Method(http.MethodGet).
		Params(":order_id", uuid.New().String()).
		Params(":refund_id", bson.NewObjectId().Hex()).
		Path(common.AuthUserGroupPath + orderRefundsIdsPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, res.Code)
	assert.NotEmpty(suite.T(), res.Body.String())

	refund := &billing.JsonRefund{}
	err = json.Unmarshal(res.Body.Bytes(), refund)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), refund.Id)
	assert.NotEmpty(suite.T(), refund.Currency)
	assert.Len(suite.T(), refund.Currency, 3)
}

func (suite *OrderTestSuite) TestOrder_GetRefund_RefundIdEmpty_Error() {

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Params(":order_id", uuid.New().String()).
		Path(common.AuthUserGroupPath + orderRefundsIdsPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Regexp(suite.T(), common.NewValidationError("RefundId"), httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_GetRefund_OrderIdEmpty_Error() {

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Params(":order_id", uuid.New().String()).
		Path(common.AuthUserGroupPath + orderRefundsIdsPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Regexp(suite.T(), common.NewValidationError("OrderId"), httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_GetRefund_BillingServerError() {

	suite.router.dispatch.Services.Billing = mock.NewBillingServerSystemErrorMock()

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Params(":order_id", uuid.New().String()).
		Params(":refund_id", bson.NewObjectId().Hex()).
		Path(common.AuthUserGroupPath + orderRefundsIdsPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusInternalServerError, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorUnknown, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_GetRefund_BillingServer_RefundNotFound_Error() {

	suite.router.dispatch.Services.Billing = mock.NewBillingServerErrorMock()

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Params(":order_id", uuid.New().String()).
		Params(":refund_id", bson.NewObjectId().Hex()).
		Path(common.AuthUserGroupPath + orderRefundsIdsPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusNotFound, httpErr.Code)
	assert.Equal(suite.T(), mock.SomeError, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_ListRefunds_Ok() {

	res, err := suite.caller.Builder().
		Method(http.MethodGet).
		Params(":order_id", uuid.New().String()).
		Path(common.AuthUserGroupPath + orderRefundsPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, res.Code)
	assert.NotEmpty(suite.T(), res.Body.String())
}

func (suite *OrderTestSuite) TestOrder_ListRefunds_BindError() {

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Path(common.AuthUserGroupPath + orderRefundsPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Regexp(suite.T(), common.NewValidationError("OrderId"), httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_ListRefunds_BillingServerError() {

	suite.router.dispatch.Services.Billing = mock.NewBillingServerSystemErrorMock()

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Params(":order_id", uuid.New().String()).
		Path(common.AuthUserGroupPath + orderRefundsPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusInternalServerError, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorUnknown, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_CreateRefund_Ok() {
	data := `{"amount": 10, "reason": "test"}`

	res, err := suite.caller.Builder().
		Method(http.MethodPost).
		Params(":order_id", uuid.New().String()).
		Path(common.AuthUserGroupPath + orderRefundsPath).
		Init(test.ReqInitJSON()).
		BodyString(data).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusCreated, res.Code)
	assert.NotEmpty(suite.T(), res.Body.String())
}

func (suite *OrderTestSuite) TestOrder_CreateRefund_BindError() {
	data := `{"amount": "qwerty", "reason": "test"}`

	_, err := suite.caller.Builder().
		Method(http.MethodPost).
		Params(":order_id", uuid.New().String()).
		Path(common.AuthUserGroupPath + orderRefundsPath).
		Init(test.ReqInitJSON()).
		BodyString(data).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorRequestParamsIncorrect, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_CreateRefund_ValidationError() {
	data := `{"amount": -10, "reason": "test"}`

	_, err := suite.caller.Builder().
		Method(http.MethodPost).
		Params(":order_id", uuid.New().String()).
		Path(common.AuthUserGroupPath + orderRefundsPath).
		Init(test.ReqInitJSON()).
		BodyString(data).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Regexp(suite.T(), common.NewValidationError("Amount"), httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_CreateRefund_BillingServerError() {
	data := `{"amount": 10, "reason": "test"}`

	suite.router.dispatch.Services.Billing = mock.NewBillingServerSystemErrorMock()

	_, err := suite.caller.Builder().
		Method(http.MethodPost).
		Params(":order_id", uuid.New().String()).
		Path(common.AuthUserGroupPath + orderRefundsPath).
		Init(test.ReqInitJSON()).
		BodyString(data).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusInternalServerError, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorUnknown, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_CreateRefund_BillingServer_CreateError() {
	data := `{"amount": 10, "reason": "test"}`

	suite.router.dispatch.Services.Billing = mock.NewBillingServerErrorMock()

	_, err := suite.caller.Builder().
		Method(http.MethodPost).
		Params(":order_id", uuid.New().String()).
		Path(common.AuthUserGroupPath + orderRefundsPath).
		Init(test.ReqInitJSON()).
		BodyString(data).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), mock.SomeError, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_ChangeLanguage_Ok() {
	body := `{"lang": "en"}`

	res, err := suite.caller.Builder().
		Method(http.MethodPatch).
		Params(":order_id", uuid.New().String()).
		Path(common.AuthProjectGroupPath + orderLanguagePath).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, res.Code)
	assert.NotEmpty(suite.T(), res.Body.String())
}

func (suite *OrderTestSuite) TestOrder_ChangeLanguage_OrderIdEmpty_Error() {
	body := `{"lang": "en"}`

	_, err := suite.caller.Builder().
		Method(http.MethodPatch).
		Params(":order_id", "").
		Path(common.AuthProjectGroupPath + orderLanguagePath).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorIncorrectOrderId, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_ChangeLanguage_BindError() {

	_, err := suite.caller.Builder().
		Method(http.MethodPatch).
		Params(":order_id", uuid.New().String()).
		Path(common.AuthProjectGroupPath + orderLanguagePath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorRequestParamsIncorrect, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_ChangeLanguage_ValidationError() {
	body := `{"lang": "en"}`

	_, err := suite.caller.Builder().
		Method(http.MethodPatch).
		Params(":order_id", "some_value").
		Path(common.AuthProjectGroupPath + orderLanguagePath).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Regexp(suite.T(), common.NewValidationError("OrderId"), httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_ChangeLanguage_BillingServerSystemError() {
	body := `{"lang": "en"}`

	suite.router.dispatch.Services.Billing = mock.NewBillingServerSystemErrorMock()

	_, err := suite.caller.Builder().
		Method(http.MethodPatch).
		Params(":order_id", uuid.New().String()).
		Path(common.AuthProjectGroupPath + orderLanguagePath).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusInternalServerError, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorUnknown, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_ChangeLanguage_BillingServerErrorResult_Error() {
	body := `{"lang": "en"}`

	suite.router.dispatch.Services.Billing = mock.NewBillingServerErrorMock()

	_, err := suite.caller.Builder().
		Method(http.MethodPatch).
		Params(":order_id", uuid.New().String()).
		Path(common.AuthProjectGroupPath + orderLanguagePath).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), mock.SomeError, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_ChangePaymentAccount_Ok() {
	body := `{"method_id": "000000000000000000000000", "account": "4000000000000002"}`

	res, err := suite.caller.Builder().
		Method(http.MethodPatch).
		Params(":order_id", uuid.New().String()).
		Path(common.AuthProjectGroupPath + orderCustomerPath).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, res.Code)
	assert.NotEmpty(suite.T(), res.Body.String())
}

func (suite *OrderTestSuite) TestOrder_ChangePaymentAccount_OrderIdEmpty_Error() {
	body := `{"method_id": "000000000000000000000000", "account": "4000000000000002"}`

	_, err := suite.caller.Builder().
		Method(http.MethodPatch).
		Params(":order_id", "").
		Path(common.AuthProjectGroupPath + orderCustomerPath).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorIncorrectOrderId, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_ChangePaymentAccount_BindError() {
	_, err := suite.caller.Builder().
		Method(http.MethodPatch).
		Params(":order_id", uuid.New().String()).
		Path(common.AuthProjectGroupPath + orderCustomerPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorRequestParamsIncorrect, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_ChangePaymentAccount_ValidationError() {
	body := `{"method_id": "some_value", "account": "4000000000000002"}`

	_, err := suite.caller.Builder().
		Method(http.MethodPatch).
		Params(":order_id", uuid.New().String()).
		Path(common.AuthProjectGroupPath + orderCustomerPath).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Regexp(suite.T(), common.NewValidationError("MethodId"), httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_ChangePaymentAccount_BillingServerSystemError() {
	body := `{"method_id": "000000000000000000000000", "account": "4000000000000002"}`

	suite.router.dispatch.Services.Billing = mock.NewBillingServerSystemErrorMock()

	_, err := suite.caller.Builder().
		Method(http.MethodPatch).
		Params(":order_id", uuid.New().String()).
		Path(common.AuthProjectGroupPath + orderCustomerPath).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusInternalServerError, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorUnknown, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_ChangePaymentAccount_BillingServerErrorResult_Error() {
	body := `{"method_id": "000000000000000000000000", "account": "4000000000000002"}`

	suite.router.dispatch.Services.Billing = mock.NewBillingServerErrorMock()

	_, err := suite.caller.Builder().
		Method(http.MethodPatch).
		Params(":order_id", uuid.New().String()).
		Path(common.AuthProjectGroupPath + orderCustomerPath).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), mock.SomeError, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_CalculateAmounts_Ok() {
	body := `{"country": "US", "zip": "98001"}`

	res, err := suite.caller.Builder().
		Method(http.MethodPost).
		Params(":order_id", uuid.New().String()).
		Path(common.AuthProjectGroupPath + orderBillingAddressPath).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, res.Code)
	assert.NotEmpty(suite.T(), res.Body.String())
}

func (suite *OrderTestSuite) TestOrder_CalculateAmounts_NoUSA_Ok() {
	body := `{"country": "RU"}`

	res, err := suite.caller.Builder().
		Method(http.MethodPost).
		Params(":order_id", uuid.New().String()).
		Path(common.AuthProjectGroupPath + orderBillingAddressPath).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, res.Code)
	assert.NotEmpty(suite.T(), res.Body.String())
}

func (suite *OrderTestSuite) TestOrder_CalculateAmounts_OrderIdEmpty_Error() {
	body := `{"country": "US", "zip": "98001"}`

	_, err := suite.caller.Builder().
		Method(http.MethodPost).
		Params(":order_id", "").
		Path(common.AuthProjectGroupPath + orderBillingAddressPath).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorIncorrectOrderId, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_CalculateAmounts_BindError() {

	_, err := suite.caller.Builder().
		Method(http.MethodPost).
		Params(":order_id", uuid.New().String()).
		Path(common.AuthProjectGroupPath + orderBillingAddressPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorRequestParamsIncorrect, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_CalculateAmounts_ValidationError() {
	body := `{"country": "some_value", "zip": "98001"}`

	_, err := suite.caller.Builder().
		Method(http.MethodPost).
		Params(":order_id", uuid.New().String()).
		Path(common.AuthProjectGroupPath + orderBillingAddressPath).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Regexp(suite.T(), common.NewValidationError("Country"), httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_CalculateAmounts_ValidationZipError() {
	body := `{"country": "US", "zip": "00"}`

	_, err := suite.caller.Builder().
		Method(http.MethodPost).
		Params(":order_id", uuid.New().String()).
		Path(common.AuthProjectGroupPath + orderBillingAddressPath).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)

	msg, ok := httpErr.Message.(*grpc.ResponseErrorMessage)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), common.ErrorMessageIncorrectZip, msg)
	assert.Regexp(suite.T(), "Zip", msg.Details)
}

func (suite *OrderTestSuite) TestOrder_CalculateAmounts_BillingServerSystemError() {
	body := `{"country": "US", "zip": "98001"}`

	suite.router.dispatch.Services.Billing = mock.NewBillingServerSystemErrorMock()

	_, err := suite.caller.Builder().
		Method(http.MethodPost).
		Params(":order_id", uuid.New().String()).
		Path(common.AuthProjectGroupPath + orderBillingAddressPath).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusInternalServerError, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorUnknown, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_CalculateAmounts_BillingServerErrorResult_Error() {
	body := `{"country": "US", "zip": "98001"}`

	suite.router.dispatch.Services.Billing = mock.NewBillingServerErrorMock()

	_, err := suite.caller.Builder().
		Method(http.MethodPost).
		Params(":order_id", uuid.New().String()).
		Path(common.AuthProjectGroupPath + orderBillingAddressPath).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), mock.SomeError, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_CreateJson_Ok() {
	order := &billing.OrderCreateRequest{
		ProjectId:     bson.NewObjectId().Hex(),
		PaymentMethod: "BANKCARD",
		Currency:      "RUB",
		Amount:        100,
		Description:   "unit test",
		OrderId:       bson.NewObjectId().Hex(),
	}

	b, err := json.Marshal(order)
	assert.NoError(suite.T(), err)

	res, err := suite.caller.Builder().
		Method(http.MethodPost).
		Path(common.AuthProjectGroupPath + orderPath).
		Init(test.ReqInitJSON()).
		BodyBytes(b).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, res.Code)
	assert.NotEmpty(suite.T(), res.Body.String())

	var response *CreateOrderJsonProjectResponse
	err = json.Unmarshal(res.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response.PaymentFormData)
}

func (suite *OrderTestSuite) TestOrder_CreateJson_WithUser_Ok() {
	order := &billing.OrderCreateRequest{
		ProjectId:     bson.NewObjectId().Hex(),
		PaymentMethod: "BANKCARD",
		Currency:      "RUB",
		Amount:        100,
		Description:   "unit test",
		OrderId:       bson.NewObjectId().Hex(),
		User: &billing.OrderUser{
			ExternalId:    bson.NewObjectId().Hex(),
			Ip:            "127.0.0.1",
			Locale:        "ru",
			Name:          "Unit Test",
			Email:         "test@unit.test",
			EmailVerified: true,
			Metadata:      map[string]string{"field1": "val1", "field2": "val2"},
		},
	}

	b, err := json.Marshal(order)
	assert.NoError(suite.T(), err)

	reqInit := func(request *http.Request, middleware test.Middleware) {
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		request.Header.Set(common.HeaderXApiSignatureHeader, "signature")
	}

	res, err := suite.caller.Builder().
		Method(http.MethodPost).
		Path(common.AuthProjectGroupPath + orderPath).
		Init(reqInit).
		BodyBytes(b).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, res.Code)
	assert.NotEmpty(suite.T(), res.Body.String())
}

func (suite *OrderTestSuite) TestOrder_CreateJson_BindError() {
	body := `{"project_id": "` + bson.NewObjectId().Hex() + `", "amount": "qwerty"}`

	_, err := suite.caller.Builder().
		Method(http.MethodPost).
		Path(common.AuthProjectGroupPath + orderPath).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorRequestParamsIncorrect, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_CreateJson_WithUser_EmptyRequestSignature_Error() {
	order := &billing.OrderCreateRequest{
		ProjectId:     bson.NewObjectId().Hex(),
		PaymentMethod: "BANKCARD",
		Currency:      "RUB",
		Amount:        100,
		Description:   "unit test",
		OrderId:       bson.NewObjectId().Hex(),
		User: &billing.OrderUser{
			ExternalId:    bson.NewObjectId().Hex(),
			Ip:            "127.0.0.1",
			Locale:        "ru",
			Name:          "Unit Test",
			Email:         "test@unit.test",
			EmailVerified: true,
			Metadata:      map[string]string{"field1": "val1", "field2": "val2"},
		},
	}

	b, err := json.Marshal(order)
	assert.NoError(suite.T(), err)

	_, err = suite.caller.Builder().
		Method(http.MethodPost).
		Path(common.AuthProjectGroupPath + orderPath).
		Init(test.ReqInitJSON()).
		BodyBytes(b).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorMessageSignatureHeaderIsEmpty, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_CreateJson_WithUser_BillingServerSystemError() {
	order := &billing.OrderCreateRequest{
		ProjectId:     bson.NewObjectId().Hex(),
		PaymentMethod: "BANKCARD",
		Currency:      "RUB",
		Amount:        100,
		Description:   "unit test",
		OrderId:       bson.NewObjectId().Hex(),
		User: &billing.OrderUser{
			ExternalId:    bson.NewObjectId().Hex(),
			Ip:            "127.0.0.1",
			Locale:        "ru",
			Name:          "Unit Test",
			Email:         "test@unit.test",
			EmailVerified: true,
			Metadata:      map[string]string{"field1": "val1", "field2": "val2"},
		},
	}

	b, err := json.Marshal(order)
	assert.NoError(suite.T(), err)

	suite.router.dispatch.Services.Billing = mock.NewBillingServerSystemErrorMock()

	reqInit := func(request *http.Request, middleware test.Middleware) {
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		request.Header.Set(common.HeaderXApiSignatureHeader, "signature")
	}

	_, err = suite.caller.Builder().
		Method(http.MethodPost).
		Path(common.AuthProjectGroupPath + orderPath).
		Init(reqInit).
		BodyBytes(b).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusInternalServerError, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorUnknown, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_CreateJson_WithUser_BillingServerResultFail_Error() {
	order := &billing.OrderCreateRequest{
		ProjectId:     bson.NewObjectId().Hex(),
		PaymentMethod: "BANKCARD",
		Currency:      "RUB",
		Amount:        100,
		Description:   "unit test",
		OrderId:       bson.NewObjectId().Hex(),
		User: &billing.OrderUser{
			ExternalId:    bson.NewObjectId().Hex(),
			Ip:            "127.0.0.1",
			Locale:        "ru",
			Name:          "Unit Test",
			Email:         "test@unit.test",
			EmailVerified: true,
			Metadata:      map[string]string{"field1": "val1", "field2": "val2"},
		},
	}

	b, err := json.Marshal(order)
	assert.NoError(suite.T(), err)

	suite.router.dispatch.Services.Billing = mock.NewBillingServerErrorMock()

	reqInit := func(request *http.Request, middleware test.Middleware) {
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		request.Header.Set(common.HeaderXApiSignatureHeader, "signature")
	}

	_, err = suite.caller.Builder().
		Method(http.MethodPost).
		Path(common.AuthProjectGroupPath + orderPath).
		Init(reqInit).
		BodyBytes(b).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), mock.SomeError, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_CreateJson_ValidationError() {
	body := `{"amount": -10}`

	_, err := suite.caller.Builder().
		Method(http.MethodPost).
		Path(common.AuthProjectGroupPath + orderPath).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Regexp(suite.T(), common.NewValidationError("Amount"), httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_CreateJson_OrderCreateError() {
	order := &billing.OrderCreateRequest{
		ProjectId:     mock.SomeMerchantId,
		PaymentMethod: "BANKCARD",
		Currency:      "RUB",
		Amount:        100,
		Description:   "unit test",
		OrderId:       bson.NewObjectId().Hex(),
	}

	b, err := json.Marshal(order)
	assert.NoError(suite.T(), err)

	suite.router.dispatch.Services.Billing = mock.NewBillingServerSystemErrorMock()

	_, err = suite.caller.Builder().
		Method(http.MethodPost).
		Path(common.AuthProjectGroupPath + orderPath).
		Init(test.ReqInitJSON()).
		BodyBytes(b).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), mock.SomeError, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_CreateJson_PaymentFormJsonDataProcessError() {
	order := &billing.OrderCreateRequest{
		ProjectId:     mock.SomeMerchantId1,
		PaymentMethod: "BANKCARD",
		Currency:      "RUB",
		Amount:        100,
		Description:   "unit test",
		OrderId:       bson.NewObjectId().Hex(),
	}

	b, err := json.Marshal(order)
	assert.NoError(suite.T(), err)

	suite.router.dispatch.Services.Billing = mock.NewBillingServerSystemErrorMock()

	_, err = suite.caller.Builder().
		Method(http.MethodPost).
		Path(common.AuthProjectGroupPath + orderPath).
		Init(test.ReqInitJSON()).
		BodyBytes(b).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), mock.SomeError, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_CreateJson_ProductionEnvironment_Ok() {
	order := &billing.OrderCreateRequest{
		ProjectId:     bson.NewObjectId().Hex(),
		PaymentMethod: "BANKCARD",
		Currency:      "RUB",
		Amount:        100,
		Description:   "unit test",
		OrderId:       bson.NewObjectId().Hex(),
	}

	b, err := json.Marshal(order)
	assert.NoError(suite.T(), err)

	suite.router.cfg.ReturnPaymentForm = false

	res, err := suite.caller.Builder().
		Method(http.MethodPost).
		Params(":"+common.RequestParameterId, uuid.New().String()).
		Path(common.AuthProjectGroupPath + orderPath).
		Init(test.ReqInitJSON()).
		BodyBytes(b).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, res.Code)
	assert.NotEmpty(suite.T(), res.Body.String())

	var response *CreateOrderJsonProjectResponse
	err = json.Unmarshal(res.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Nil(suite.T(), response.PaymentFormData)
}

func (suite *OrderTestSuite) TestOrder_GetOrderForm_Ok() {

	res, err := suite.caller.Builder().
		Method(http.MethodGet).
		Params(":"+common.RequestParameterId, uuid.New().String()).
		Path(orderIdPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, res.Code)
	assert.NotEmpty(suite.T(), res.Body.String())
	assert.Equal(suite.T(), echo.MIMETextHTMLCharsetUTF8, res.Header().Get(echo.HeaderContentType))

	cookies := res.Result().Cookies()
	assert.True(suite.T(), len(cookies) == 1)
	assert.Equal(suite.T(), common.CustomerTokenCookiesName, cookies[0].Name)
	assert.True(suite.T(), cookies[0].HttpOnly)
}

func (suite *OrderTestSuite) TestOrder_GetOrderForm_TokenCookieExist_Ok() {

	cookie := new(http.Cookie)
	cookie.Name = common.CustomerTokenCookiesName
	cookie.Value = bson.NewObjectId().Hex()
	cookie.Expires = time.Now().Add(suite.router.cfg.CustomerTokenCookiesLifetime)
	cookie.HttpOnly = true

	res, err := suite.caller.Builder().
		Method(http.MethodGet).
		AddCookie(cookie).
		Params(":"+common.RequestParameterId, mock.SomeMerchantId1).
		Path(orderIdPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, res.Code)
	assert.NotEmpty(suite.T(), res.Body.String())
	assert.Equal(suite.T(), echo.MIMETextHTMLCharsetUTF8, res.Header().Get(echo.HeaderContentType))

	cookies := res.Result().Cookies()

	assert.True(suite.T(), len(cookies) > 0)
	assert.Equal(suite.T(), cookie.Name, cookies[0].Name)
	assert.Equal(suite.T(), cookie.Value, cookies[0].Value)
}

func (suite *OrderTestSuite) TestOrder_GetOrderForm_ParameterIdNotFound_Error() {

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Params(":"+common.RequestParameterId, "").
		Path(orderIdPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusNotFound, httpErr.Code)
}

func (suite *OrderTestSuite) TestOrder_GetOrderForm_BillingServerSystemError() {

	suite.router.dispatch.Services.Billing = mock.NewBillingServerSystemErrorMock()

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Params(":"+common.RequestParameterId, bson.NewObjectId().Hex()).
		Path(orderIdPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusInternalServerError, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorUnknown, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_GetOrders_Ok() {

	bs := &billMock.BillingService{}
	bs.On("FindAllOrdersPublic", mock2.Anything, mock2.Anything, mock2.Anything).
		Return(
			&grpc.ListOrdersPublicResponse{
				Status: pkg.ResponseStatusOk,
				Item: &grpc.ListOrdersPublicResponseItem{
					Count: 1,
					Items: []*billing.OrderViewPublic{},
				},
			},
			nil,
		)
	suite.router.dispatch.Services.Billing = bs

	res, err := suite.caller.Builder().
		Method(http.MethodGet).
		Path(common.AuthUserGroupPath + orderPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, res.Code)
	assert.NotEmpty(suite.T(), res.Body.String())
}

func (suite *OrderTestSuite) TestOrder_GetOrders_BillingServerError() {

	bs := &billMock.BillingService{}
	bs.On("FindAllOrdersPublic", mock2.Anything, mock2.Anything, mock2.Anything).
		Return(nil, errors.New("some error"))
	suite.router.dispatch.Services.Billing = bs

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Path(common.AuthUserGroupPath + orderPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusInternalServerError, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorUnknown, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_GetOrders_BindError_Id() {
	q := make(url.Values)
	q.Set(common.RequestParameterId, bson.NewObjectId().Hex())
	suite.testGetOrdersBindError(q, fmt.Sprintf(common.ErrorMessageMask, "Id", "uuid"))
}

func (suite *OrderTestSuite) TestOrder_GetOrders_BindError_Project() {
	q := url.Values{common.RequestParameterProject: []string{"foo"}}
	suite.testGetOrdersBindError(q, fmt.Sprintf(common.ErrorMessageMask, "Project[0]", "hexadecimal"))
}

func (suite *OrderTestSuite) TestOrder_GetOrders_BindError_PaymentMethod() {
	q := url.Values{common.RequestParameterPaymentMethod: []string{"foo"}}
	suite.testGetOrdersBindError(q, fmt.Sprintf(common.ErrorMessageMask, "PaymentMethod[0]", "hexadecimal"))
}

func (suite *OrderTestSuite) TestOrder_GetOrders_BindError_Country() {
	q := url.Values{common.RequestParameterCountries: []string{"foo"}}
	suite.testGetOrdersBindError(q, fmt.Sprintf(common.ErrorMessageMask, "Country[0]", "len"))
}

func (suite *OrderTestSuite) testGetOrdersBindError(q url.Values, error string) {

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		SetQueryParams(q).
		Path(common.AuthUserGroupPath + orderPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), common.NewValidationError(error), httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_CreateJson_WithPreparedOrderId_Ok() {
	order := &billing.OrderCreateRequest{
		ProjectId:    bson.NewObjectId().Hex(),
		PspOrderUuid: uuid.New().String(),
	}

	b, err := json.Marshal(order)
	assert.NoError(suite.T(), err)

	res, err := suite.caller.Builder().
		Method(http.MethodPost).
		Path(common.AuthProjectGroupPath + orderPath).
		Init(test.ReqInitJSON()).
		BodyBytes(b).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, res.Code)
	assert.NotEmpty(suite.T(), res.Body.String())
}

func (suite *OrderTestSuite) TestOrder_CreateJson_WithPreparedOrderId_BillingServiceSystemError() {
	order := &billing.OrderCreateRequest{
		ProjectId:    bson.NewObjectId().Hex(),
		PspOrderUuid: uuid.New().String(),
	}

	b, err := json.Marshal(order)
	assert.NoError(suite.T(), err)

	suite.router.dispatch.Services.Billing = mock.NewBillingServerSystemErrorMock()

	_, err = suite.caller.Builder().
		Method(http.MethodPost).
		Path(common.AuthProjectGroupPath + orderPath).
		Init(test.ReqInitJSON()).
		BodyBytes(b).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusInternalServerError, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorUnknown, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_CreateJson_WithPreparedOrderId_BillingServiceResultError() {
	order := &billing.OrderCreateRequest{
		ProjectId:    bson.NewObjectId().Hex(),
		PspOrderUuid: uuid.New().String(),
	}

	b, err := json.Marshal(order)
	assert.NoError(suite.T(), err)

	suite.router.dispatch.Services.Billing = mock.NewBillingServerErrorMock()

	_, err = suite.caller.Builder().
		Method(http.MethodPost).
		Path(common.AuthProjectGroupPath + orderPath).
		Init(test.ReqInitJSON()).
		BodyBytes(b).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), mock.SomeError, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_ChangeOrderCode_Ok() {
	shouldBe := require.New(suite.T())

	changeOrderRequest := &grpc.ChangeCodeInOrderRequest{
		KeyProductId: bson.NewObjectId().Hex(),
	}
	b, err := json.Marshal(changeOrderRequest)
	assert.NoError(suite.T(), err)

	billingService := &billMock.BillingService{}
	billingService.On("ChangeCodeInOrder", mock2.Anything, mock2.Anything).Return(&grpc.ChangeCodeInOrderResponse{
		Status: pkg.ResponseStatusOk,
		Order:  &billing.Order{},
	}, nil)

	suite.router.dispatch.Services.Billing = billingService

	res, err := suite.caller.Builder().
		Method(http.MethodPut).
		Params(":order_id", bson.NewObjectId().Hex()).
		Path(common.AuthUserGroupPath + orderReplaceCodePath).
		Init(test.ReqInitJSON()).
		BodyBytes(b).
		Exec(suite.T())

	shouldBe.NoError(err)
	shouldBe.Nil(err)
	shouldBe.EqualValues(http.StatusOK, res.Code)
	shouldBe.NotEmpty(res.Body.String())
}

func (suite *OrderTestSuite) TestOrder_ChangeOrderCode_ServiceError() {
	shouldBe := require.New(suite.T())

	changeOrderRequest := &grpc.ChangeCodeInOrderRequest{
		KeyProductId: bson.NewObjectId().Hex(),
	}
	b, err := json.Marshal(changeOrderRequest)
	assert.NoError(suite.T(), err)

	billingService := &billMock.BillingService{}
	billingService.On("ChangeCodeInOrder", mock2.Anything, mock2.Anything).Return(nil, errors.New("some error"))

	suite.router.dispatch.Services.Billing = billingService

	_, err = suite.caller.Builder().
		Method(http.MethodPut).
		Params(":order_id", bson.NewObjectId().Hex()).
		Path(common.AuthUserGroupPath + orderReplaceCodePath).
		Init(test.ReqInitJSON()).
		BodyBytes(b).
		Exec(suite.T())

	shouldBe.NotNil(err)
	httpErr, ok := err.(*echo.HTTPError)
	shouldBe.True(ok)
	shouldBe.EqualValues(http.StatusInternalServerError, httpErr.Code)
}

func (suite *OrderTestSuite) TestOrder_ChangeOrderCode_ErrorInService() {
	shouldBe := require.New(suite.T())

	changeOrderRequest := &grpc.ChangeCodeInOrderRequest{
		KeyProductId: bson.NewObjectId().Hex(),
	}
	b, err := json.Marshal(changeOrderRequest)
	assert.NoError(suite.T(), err)

	billingService := &billMock.BillingService{}
	billingService.On("ChangeCodeInOrder", mock2.Anything, mock2.Anything).Return(&grpc.ChangeCodeInOrderResponse{
		Status: 400,
	}, nil)

	suite.router.dispatch.Services.Billing = billingService

	_, err = suite.caller.Builder().
		Method(http.MethodPut).
		Params(":order_id", bson.NewObjectId().Hex()).
		Path(common.AuthUserGroupPath + orderReplaceCodePath).
		Init(test.ReqInitJSON()).
		BodyBytes(b).
		Exec(suite.T())

	shouldBe.NotNil(err)
	httpErr, ok := err.(*echo.HTTPError)
	shouldBe.True(ok)
	shouldBe.EqualValues(http.StatusBadRequest, httpErr.Code)
}

func (suite *OrderTestSuite) TestOrder_ChangeOrderCode_ValidationError() {
	shouldBe := require.New(suite.T())

	// Missing key product id
	changeOrderRequest := &grpc.ChangeCodeInOrderRequest{
		KeyProductId: "",
	}
	b, err := json.Marshal(changeOrderRequest)
	assert.NoError(suite.T(), err)

	_, err = suite.caller.Builder().
		Method(http.MethodPut).
		Params(":order_id", bson.NewObjectId().Hex()).
		Path(common.AuthUserGroupPath + orderReplaceCodePath).
		Init(test.ReqInitJSON()).
		BodyBytes(b).
		Exec(suite.T())

	shouldBe.NotNil(err)
	httpErr, ok := err.(*echo.HTTPError)
	shouldBe.True(ok)
	shouldBe.EqualValues(http.StatusBadRequest, httpErr.Code)

	// Wrong order id
	changeOrderRequest = &grpc.ChangeCodeInOrderRequest{
		KeyProductId: bson.NewObjectId().Hex(),
	}
	b, err = json.Marshal(changeOrderRequest)
	assert.NoError(suite.T(), err)
}
