package handlers

import (
	"errors"
	"github.com/labstack/echo/v4"
	billingMocks "github.com/paysuper/paysuper-billing-server/pkg/mocks"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/grpc"
	"github.com/paysuper/paysuper-management-api/internal/dispatcher/common"
	"github.com/paysuper/paysuper-management-api/internal/mock"
	"github.com/paysuper/paysuper-management-api/internal/test"
	"github.com/stretchr/testify/assert"
	mock2 "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type PaymentMethodTestSuite struct {
	suite.Suite
	router *PaymentMethodApiV1
	caller *test.EchoReqResCaller
}

func Test_PaymentMethod(t *testing.T) {
	suite.Run(t, new(PaymentMethodTestSuite))
}

func (suite *PaymentMethodTestSuite) SetupTest() {
	var e error
	settings := test.DefaultSettings()
	srv := common.Services{
		Billing: mock.NewBillingServerOkMock(),
	}
	suite.caller, e = test.SetUp(settings, srv, func(set *test.TestSet, mw test.Middleware) common.Handlers {
		suite.router = NewPaymentMethodApiV1(set.HandlerSet, set.GlobalConfig)
		return common.Handlers{
			suite.router,
		}
	})
	if e != nil {
		panic(e)
	}
}

func (suite *PaymentMethodTestSuite) TearDownTest() {}

func (suite *PaymentMethodTestSuite) TestPaymentMethod_create_BindError_RequiredName() {
	data := `{"payment_system_id": "payment_system_id"}`

	res, err := suite.caller.Builder().
		Method(http.MethodPost).
		Path(common.AccessGroupPath + paymentMethodPath).
		Init(test.ReqInitJSON()).
		BodyString(data).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusOK, res.Code)

	msg, ok := httpErr.Message.(*grpc.ResponseErrorMessage)
	assert.True(suite.T(), ok)
	assert.Regexp(suite.T(), "Name", msg.Details)
}

func (suite *PaymentMethodTestSuite) TestPaymentMethod_create_BindError_RequiredPaymentSystemId() {
	data := `{"name": "name"}`

	_, err := suite.caller.Builder().
		Method(http.MethodPost).
		Path(common.AccessGroupPath + paymentMethodPath).
		Init(test.ReqInitJSON()).
		BodyString(data).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)

	msg, ok := httpErr.Message.(*grpc.ResponseErrorMessage)
	assert.True(suite.T(), ok)
	assert.Regexp(suite.T(), "PaymentSystemId", msg.Details)
}

func (suite *PaymentMethodTestSuite) TestPaymentMethod_create_BindError_Name() {
	data := `{"name": "!", "payment_system_id": "payment_system_id"}`

	_, err := suite.caller.Builder().
		Method(http.MethodPost).
		Path(common.AccessGroupPath + paymentMethodPath).
		Init(test.ReqInitJSON()).
		BodyString(data).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)

	msg, ok := httpErr.Message.(*grpc.ResponseErrorMessage)
	assert.True(suite.T(), ok)
	assert.Regexp(suite.T(), "field validation for 'Name' failed on the 'alphanum' tag", msg.Details)
}

func (suite *PaymentMethodTestSuite) TestPaymentMethod_create_BindError_PaymentSystemId() {
	data := `{"name": "name", "payment_system_id": "1"}`

	_, err := suite.caller.Builder().
		Method(http.MethodPost).
		Path(common.AccessGroupPath + paymentMethodPath).
		Init(test.ReqInitJSON()).
		BodyString(data).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)

	msg, ok := httpErr.Message.(*grpc.ResponseErrorMessage)
	assert.True(suite.T(), ok)
	assert.Regexp(suite.T(), "field validation for 'PaymentSystemId' failed on the 'len' tag", msg.Details)
}

func (suite *PaymentMethodTestSuite) TestPaymentMethod_Error_BillingServer() {
	data := `{"name": "name", "payment_system_id": "507f1f77bcf86cd799439011"}`

	billingService := &billingMocks.BillingService{}
	billingService.On("CreateOrUpdatePaymentMethod", mock2.Anything, mock2.Anything).Return(nil, errors.New("error"))
	suite.router.dispatch.Services.Billing = billingService

	_, err := suite.caller.Builder().
		Method(http.MethodPost).
		Path(common.AccessGroupPath + paymentMethodPath).
		Init(test.ReqInitJSON()).
		BodyString(data).
		Exec(suite.T())

	assert.Error(suite.T(), err)
}

func (suite *PaymentMethodTestSuite) TestPaymentMethod_create_Ok() {
	data := `{"name": "name", "payment_system_id": "507f1f77bcf86cd799439011"}`

	billingService := &billingMocks.BillingService{}
	billingService.On("CreateOrUpdatePaymentMethod", mock2.Anything, mock2.Anything).Return(&grpc.ChangePaymentMethodResponse{}, nil)
	suite.router.dispatch.Services.Billing = billingService

	_, err := suite.caller.Builder().
		Method(http.MethodPost).
		Path(common.AccessGroupPath + paymentMethodPath).
		Init(test.ReqInitJSON()).
		BodyString(data).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
}

func (suite *PaymentMethodTestSuite) TestPaymentMethod_update_Ok() {
	data := `{"name": "name", "payment_system_id": "507f1f77bcf86cd799439011"}`

	billingService := &billingMocks.BillingService{}
	billingService.On("CreateOrUpdatePaymentMethod", mock2.Anything, mock2.Anything).Return(&grpc.ChangePaymentMethodResponse{}, nil)
	suite.router.dispatch.Services.Billing = billingService

	_, err := suite.caller.Builder().
		Method(http.MethodPut).
		Path(common.AccessGroupPath + paymentMethodIdPath).
		Init(test.ReqInitJSON()).
		BodyString(data).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
}

func (suite *PaymentMethodTestSuite) TestPaymentMethod_getProductionSettings_BindError_RequiredPaymentMethodId() {
	data := `{"currency_a3": "rub"}`

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Params(":"+common.RequestParameterId, "").
		Path(common.AccessGroupPath + paymentMethodProductionPath).
		Init(test.ReqInitJSON()).
		BodyString(data).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)

	msg, ok := httpErr.Message.(*grpc.ResponseErrorMessage)
	assert.True(suite.T(), ok)
	assert.Regexp(suite.T(), "field validation for 'PaymentMethodId' failed on the 'required' tag", msg.Details)
}

func (suite *PaymentMethodTestSuite) TestPaymentMethod_getProductionSettings_Error_BillingServer() {
	data := `{"currency_a3": "rub", "payment_method_id": "507f1f77bcf86cd799439011"}`

	billingService := &billingMocks.BillingService{}
	billingService.On("GetPaymentMethodProductionSettings", mock2.Anything, mock2.Anything).Return(nil, errors.New("error"))
	suite.router.dispatch.Services.Billing = billingService

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Path(common.AccessGroupPath + paymentMethodProductionPath).
		Init(test.ReqInitJSON()).
		BodyString(data).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Regexp(suite.T(), common.ErrorUnknown.Message, httpErr.Message)
}

func (suite *PaymentMethodTestSuite) TestPaymentMethod_getProductionSettings_Ok() {
	data := `{"currency_a3": "rub", "payment_method_id": "507f1f77bcf86cd799439011"}`

	billingService := &billingMocks.BillingService{}
	billingService.On("GetPaymentMethodProductionSettings", mock2.Anything, mock2.Anything).Return(&grpc.GetPaymentMethodSettingsResponse{}, nil)
	suite.router.dispatch.Services.Billing = billingService

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Path(common.AccessGroupPath + paymentMethodProductionPath).
		Init(test.ReqInitJSON()).
		BodyString(data).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
}

func (suite *PaymentMethodTestSuite) TestPaymentMethod_getTestSettings_BindError_RequiredPaymentMethodId() {
	data := `{"currency_a3": "rub"}`

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Params(":"+common.RequestParameterId, "").
		Path(common.AccessGroupPath + paymentMethodTestPath).
		Init(test.ReqInitJSON()).
		BodyString(data).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)

	msg, ok := httpErr.Message.(*grpc.ResponseErrorMessage)
	assert.True(suite.T(), ok)
	assert.Regexp(suite.T(), "field validation for 'PaymentMethodId' failed on the 'required' tag", msg.Details)
}

func (suite *PaymentMethodTestSuite) TestPaymentMethod_getTestSettings_Error_BillingServer() {
	data := `{"currency_a3": "rub", "payment_method_id": "507f1f77bcf86cd799439011"}`

	billingService := &billingMocks.BillingService{}
	billingService.On("GetPaymentMethodTestSettings", mock2.Anything, mock2.Anything).Return(nil, errors.New("error"))
	suite.router.dispatch.Services.Billing = billingService

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Params(":"+common.RequestParameterId, "").
		Path(common.AccessGroupPath + paymentMethodTestPath).
		Init(test.ReqInitJSON()).
		BodyString(data).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Regexp(suite.T(), common.ErrorUnknown.Message, httpErr.Message)
}

func (suite *PaymentMethodTestSuite) TestPaymentMethod_getTestSettings_Ok() {
	data := `{"currency_a3": "rub", "payment_method_id": "507f1f77bcf86cd799439011"}`

	billingService := &billingMocks.BillingService{}
	billingService.On("GetPaymentMethodTestSettings", mock2.Anything, mock2.Anything).Return(&grpc.GetPaymentMethodSettingsResponse{}, nil)
	suite.router.dispatch.Services.Billing = billingService

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Path(common.AccessGroupPath + paymentMethodTestPath).
		Init(test.ReqInitJSON()).
		BodyString(data).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
}

func (suite *PaymentMethodTestSuite) TestPaymentMethod_createProductionSettings_BindError_RequiredParams() {
	data := `{"payment_method_id": "507f1f77bcf86cd799439011"}`

	_, err := suite.caller.Builder().
		Method(http.MethodPost).
		Params(":"+common.RequestParameterId, "").
		Path(common.AccessGroupPath + paymentMethodProductionPath).
		Init(test.ReqInitJSON()).
		BodyString(data).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)

	msg, ok := httpErr.Message.(*grpc.ResponseErrorMessage)
	assert.True(suite.T(), ok)
	assert.Regexp(suite.T(), "field validation for 'Params' failed on the 'required' tag", msg.Details)
}

func (suite *PaymentMethodTestSuite) TestPaymentMethod_createProductionSettings_BindError_RequiredPaymentMethodId() {
	data := `{"params": {}}`

	_, err := suite.caller.Builder().
		Method(http.MethodPost).
		Params(":"+common.RequestParameterId, "").
		Path(common.AccessGroupPath + paymentMethodProductionPath).
		Init(test.ReqInitJSON()).
		BodyString(data).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)

	msg, ok := httpErr.Message.(*grpc.ResponseErrorMessage)
	assert.True(suite.T(), ok)
	assert.Regexp(suite.T(), "field validation for 'PaymentMethodId' failed on the 'required' tag", msg.Details)
}

func (suite *PaymentMethodTestSuite) TestPaymentMethod_updateProductionSettings_BindError_RequiredParams() {
	data := `{"payment_method_id": "507f1f77bcf86cd799439011"}`

	_, err := suite.caller.Builder().
		Method(http.MethodPut).
		Params(":"+common.RequestParameterId, "").
		Path(common.AccessGroupPath + paymentMethodProductionPath).
		Init(test.ReqInitJSON()).
		BodyString(data).
		Exec(suite.T())
	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)

	msg, ok := httpErr.Message.(*grpc.ResponseErrorMessage)
	assert.True(suite.T(), ok)
	assert.Regexp(suite.T(), "field validation for 'Params' failed on the 'required' tag", msg.Details)
}

func (suite *PaymentMethodTestSuite) TestPaymentMethod_updateProductionSettings_BindError_RequiredPaymentMethodId() {
	data := `{"params": {}}`

	_, err := suite.caller.Builder().
		Method(http.MethodPut).
		Params(":"+common.RequestParameterId, "").
		Path(common.AccessGroupPath + paymentMethodProductionPath).
		Init(test.ReqInitJSON()).
		BodyString(data).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)

	msg, ok := httpErr.Message.(*grpc.ResponseErrorMessage)
	assert.True(suite.T(), ok)
	assert.Regexp(suite.T(), "field validation for 'PaymentMethodId' failed on the 'required' tag", msg.Details)
}

func (suite *PaymentMethodTestSuite) TestPaymentMethod_createTestSettings_BindError_RequiredParams() {
	data := `{"payment_method_id": "507f1f77bcf86cd799439011"}`

	_, err := suite.caller.Builder().
		Method(http.MethodPost).
		Params(":"+common.RequestParameterId, "").
		Path(common.AccessGroupPath + paymentMethodProductionPath).
		Init(test.ReqInitJSON()).
		BodyString(data).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)

	msg, ok := httpErr.Message.(*grpc.ResponseErrorMessage)
	assert.True(suite.T(), ok)
	assert.Regexp(suite.T(), "field validation for 'Params' failed on the 'required' tag", msg.Details)
}

func (suite *PaymentMethodTestSuite) TestPaymentMethod_createTestSettings_BindError_RequiredPaymentMethodId() {
	data := `{"params": {}}`

	_, err := suite.caller.Builder().
		Method(http.MethodPost).
		Params(":"+common.RequestParameterId, "").
		Path(common.AccessGroupPath + paymentMethodProductionPath).
		Init(test.ReqInitJSON()).
		BodyString(data).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)

	msg, ok := httpErr.Message.(*grpc.ResponseErrorMessage)
	assert.True(suite.T(), ok)
	assert.Regexp(suite.T(), "field validation for 'PaymentMethodId' failed on the 'required' tag", msg.Details)
}

func (suite *PaymentMethodTestSuite) TestPaymentMethod_updateTestSettings_BindError_RequiredParams() {
	data := `{"payment_method_id": "507f1f77bcf86cd799439011"}`

	_, err := suite.caller.Builder().
		Method(http.MethodPut).
		Params(":"+common.RequestParameterId, "").
		Path(common.AccessGroupPath + paymentMethodProductionPath).
		Init(test.ReqInitJSON()).
		BodyString(data).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)

	msg, ok := httpErr.Message.(*grpc.ResponseErrorMessage)
	assert.True(suite.T(), ok)
	assert.Regexp(suite.T(), "field validation for 'Params' failed on the 'required' tag", msg.Details)
}

func (suite *PaymentMethodTestSuite) TestPaymentMethod_updateTestSettings_BindError_RequiredPaymentMethodId() {
	data := `{"params": {}}`

	_, err := suite.caller.Builder().
		Method(http.MethodPut).
		Params(":"+common.RequestParameterId, "").
		Path(common.AccessGroupPath + paymentMethodProductionPath).
		Init(test.ReqInitJSON()).
		BodyString(data).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)

	msg, ok := httpErr.Message.(*grpc.ResponseErrorMessage)
	assert.True(suite.T(), ok)
	assert.Regexp(suite.T(), "field validation for 'PaymentMethodId' failed on the 'required' tag", msg.Details)
}
