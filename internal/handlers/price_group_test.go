package handlers

import (
	"errors"
	"github.com/labstack/echo/v4"
	billingMocks "github.com/paysuper/paysuper-billing-server/pkg/mocks"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/billing"
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

type PriceGroupTestSuite struct {
	suite.Suite
	router *PriceGroup
	caller *test.EchoReqResCaller
}

func Test_PriceGroup(t *testing.T) {
	suite.Run(t, new(PriceGroupTestSuite))
}

func (suite *PriceGroupTestSuite) SetupTest() {
	var e error
	settings := test.DefaultSettings()
	srv := common.Services{
		Billing: mock.NewBillingServerOkMock(),
	}
	suite.caller, e = test.SetUp(settings, srv, func(set *test.TestSet, mw test.Middleware) common.Handlers {
		suite.router = NewPriceGroup(set.HandlerSet, set.GlobalConfig)
		return common.Handlers{
			suite.router,
		}
	})
	if e != nil {
		panic(e)
	}
}

func (suite *PriceGroupTestSuite) TearDownTest() {}

func (suite *PriceGroupTestSuite) TestPaymentMethod_getPriceGroupByCountry_BindError_RequiredCountry() {
	data := `{"variable": "test"}`

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Path(common.AccessGroupPath + priceGroupCountryPath).
		Init(test.ReqInitJSON()).
		BodyString(data).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)

	msg, ok := httpErr.Message.(*grpc.ResponseErrorMessage)
	assert.True(suite.T(), ok)
	assert.Regexp(suite.T(), "field validation for 'Country' failed on the 'required' tag", msg.Details)
}

func (suite *PriceGroupTestSuite) TestPaymentMethod_getPriceGroupByCountry_Error_BillingServer() {
	data := `{"country": "RU"}`

	billingService := &billingMocks.BillingService{}
	billingService.On("GetPriceGroupByCountry", mock2.Anything, mock2.Anything).Return(nil, errors.New("error"))
	suite.router.dispatch.Services.Billing = billingService

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Path(common.AccessGroupPath + priceGroupCountryPath).
		Init(test.ReqInitJSON()).
		BodyString(data).
		Exec(suite.T())

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusInternalServerError, httpErr.Code)
	assert.Regexp(suite.T(), common.ErrorMessagePriceGroupByCountry.Message, httpErr.Message)
}

func (suite *PriceGroupTestSuite) TestPaymentMethod_getPriceGroupByCountry_Ok() {
	data := `{"country": "RU"}`

	billingService := &billingMocks.BillingService{}
	billingService.On("GetPriceGroupByCountry", mock2.Anything, mock2.Anything).Return(&billing.PriceGroup{}, nil)
	suite.router.dispatch.Services.Billing = billingService

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Path(common.AccessGroupPath + priceGroupCountryPath).
		Init(test.ReqInitJSON()).
		BodyString(data).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
}

func (suite *PriceGroupTestSuite) TestPaymentMethod_getCurrencyList_Error_BillingServer() {

	billingService := &billingMocks.BillingService{}
	billingService.On("GetPriceGroupCurrencies", mock2.Anything, mock2.Anything).Return(nil, errors.New("error"))
	suite.router.dispatch.Services.Billing = billingService

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Path(common.AccessGroupPath + priceGroupCurrenciesPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusInternalServerError, httpErr.Code)
	assert.Regexp(suite.T(), common.ErrorMessagePriceGroupCurrencyList.Message, httpErr.Message)
}

func (suite *PriceGroupTestSuite) TestPaymentMethod_getCurrencyList_Ok() {

	billingService := &billingMocks.BillingService{}
	billingService.On("GetPriceGroupCurrencies", mock2.Anything, mock2.Anything).Return(&grpc.PriceGroupCurrenciesResponse{}, nil)
	suite.router.dispatch.Services.Billing = billingService

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Path(common.AccessGroupPath + priceGroupCurrenciesPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
}

func (suite *PriceGroupTestSuite) TestPaymentMethod_getCurrencyByRegion_BindError_RequiredRegion() {
	data := `{"variable": "test"}`

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Path(common.AccessGroupPath + priceGroupRegionPath).
		Init(test.ReqInitJSON()).
		BodyString(data).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)

	msg, ok := httpErr.Message.(*grpc.ResponseErrorMessage)
	assert.True(suite.T(), ok)
	assert.Regexp(suite.T(), "field validation for 'Region' failed on the 'required' tag", msg.Details)
}

func (suite *PriceGroupTestSuite) TestPaymentMethod_getCurrencyByRegion_Error_BillingServer() {
	data := `{"region": "RUB"}`

	billingService := &billingMocks.BillingService{}
	billingService.On("GetPriceGroupCurrencyByRegion", mock2.Anything, mock2.Anything).Return(nil, errors.New("error"))
	suite.router.dispatch.Services.Billing = billingService

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Path(common.AccessGroupPath + priceGroupRegionPath).
		Init(test.ReqInitJSON()).
		BodyString(data).
		Exec(suite.T())

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusInternalServerError, httpErr.Code)
	assert.Regexp(suite.T(), common.ErrorMessagePriceGroupCurrencyByRegion.Message, httpErr.Message)
}

func (suite *PriceGroupTestSuite) TestPaymentMethod_getCurrencyByRegion_Ok() {
	data := `{"region": "RU"}`

	billingService := &billingMocks.BillingService{}
	billingService.On("GetPriceGroupCurrencyByRegion", mock2.Anything, mock2.Anything).Return(&grpc.PriceGroupCurrenciesResponse{}, nil)
	suite.router.dispatch.Services.Billing = billingService

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Path(common.AccessGroupPath + priceGroupRegionPath).
		Init(test.ReqInitJSON()).
		BodyString(data).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
}

func (suite *PriceGroupTestSuite) TestPaymentMethod_getRecommendedPrice_BindError_RequiredAmount() {
	data := `{"variable": "test"}`

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Path(common.AccessGroupPath + priceGroupRecommendedPath).
		Init(test.ReqInitJSON()).
		BodyString(data).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)

	msg, ok := httpErr.Message.(*grpc.ResponseErrorMessage)
	assert.True(suite.T(), ok)
	assert.Regexp(suite.T(), "field validation for 'Amount' failed on the 'required' tag", msg.Details)
}

func (suite *PriceGroupTestSuite) TestPaymentMethod_getRecommendedPrice_Error_BillingServer() {
	data := `{"amount": 1}`

	billingService := &billingMocks.BillingService{}
	billingService.On("GetPriceGroupRecommendedPrice", mock2.Anything, mock2.Anything).Return(nil, errors.New("error"))
	suite.router.dispatch.Services.Billing = billingService

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Path(common.AccessGroupPath + priceGroupRecommendedPath).
		Init(test.ReqInitJSON()).
		BodyString(data).
		Exec(suite.T())

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusInternalServerError, httpErr.Code)
	assert.Regexp(suite.T(), common.ErrorMessagePriceGroupRecommendedList.Message, httpErr.Message)
}

func (suite *PriceGroupTestSuite) TestPaymentMethod_getRecommendedPrice_Ok() {
	data := `{"amount": 1}`

	billingService := &billingMocks.BillingService{}
	billingService.On("GetPriceGroupRecommendedPrice", mock2.Anything, mock2.Anything).Return(&grpc.PriceGroupRecommendedPriceResponse{}, nil)
	suite.router.dispatch.Services.Billing = billingService

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Path(common.AccessGroupPath + priceGroupRecommendedPath).
		Init(test.ReqInitJSON()).
		BodyString(data).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
}
