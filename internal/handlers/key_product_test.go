package handlers

import (
	"encoding/json"
	"errors"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo/v4"
	billMock "github.com/paysuper/paysuper-billing-server/pkg/mocks"
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

type KeyProductTestSuite struct {
	suite.Suite
	router *KeyProductRoute
	caller *test.EchoReqResCaller
}

func Test_keyProduct(t *testing.T) {
	suite.Run(t, new(KeyProductTestSuite))
}

func (suite *KeyProductTestSuite) SetupTest() {
	var e error
	settings := test.DefaultSettings()
	srv := common.Services{
		Billing: mock.NewBillingServerOkMock(),
		Geo:     mock.NewGeoIpServiceTestOk(),
	}
	suite.caller, e = test.SetUp(settings, srv, func(set *test.TestSet, mw test.Middleware) common.Handlers {
		suite.router = NewKeyProductRoute(set.HandlerSet, set.GlobalConfig)
		return common.Handlers{
			suite.router,
		}
	})
	if e != nil {
		panic(e)
	}
}

func (suite *KeyProductTestSuite) TearDownTest() {}

func (suite *KeyProductTestSuite) TestProject_RemovePlatform_Ok() {

	res, err := suite.caller.Builder().
		Method(http.MethodDelete).
		Params(":key_product_id", bson.NewObjectId().Hex()).
		Params(":platform_id", "steam").
		Path(common.AuthUserGroupPath + keyProductsPlatformIdPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, res.Code)
	assert.NotEmpty(suite.T(), res.Body.String())
}

func (suite *KeyProductTestSuite) TestProject_PublishKeyProduct_Ok() {

	res, err := suite.caller.Builder().
		Method(http.MethodPost).
		Params(":key_product_id", bson.NewObjectId().Hex()).
		Path(common.AuthUserGroupPath + keyProductsPublishPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, res.Code)
	assert.NotEmpty(suite.T(), res.Body.String())
}

func (suite *KeyProductTestSuite) TestProject_GetPlatformList_Error() {

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		SetQueryParam("limit", "qwe").
		Path(common.AuthUserGroupPath + platformsPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.Error(suite.T(), err)
}

func (suite *KeyProductTestSuite) TestProject_GetPlatformList_Ok() {

	res, err := suite.caller.Builder().
		Method(http.MethodGet).
		SetQueryParam("limit", "10").
		Path(common.AuthUserGroupPath + platformsPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, res.Code)
	assert.NotEmpty(suite.T(), res.Body.String())

	res, err = suite.caller.Builder().
		Method(http.MethodGet).
		SetQueryParam("limit", "").
		Path(common.AuthUserGroupPath + platformsPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, res.Code)
	assert.NotEmpty(suite.T(), res.Body.String())

	res, err = suite.caller.Builder().
		Method(http.MethodGet).
		SetQueryParam("limit", "300000").
		Path(common.AuthUserGroupPath + platformsPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, res.Code)
	assert.NotEmpty(suite.T(), res.Body.String())
}

func (suite *KeyProductTestSuite) TestProject_GetListKeyProduct_Ok() {

	res, err := suite.caller.Builder().
		Method(http.MethodGet).
		Path(common.AuthUserGroupPath + keyProductsPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, res.Code)
	assert.NotEmpty(suite.T(), res.Body.String())
}

func (suite *KeyProductTestSuite) TestProject_GetKeyProduct_ValidationError() {

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Params(":key_product_id", " ").
		Path(common.AuthUserGroupPath + keyProductsIdPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.Error(suite.T(), err)
	e, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), 400, e.Code)
	assert.NotEmpty(suite.T(), e.Message)
}

func (suite *KeyProductTestSuite) TestProject_GetKeyProduct_Ok() {

	res, err := suite.caller.Builder().
		Method(http.MethodGet).
		Params(":key_product_id", bson.NewObjectId().Hex()).
		Path(common.AuthUserGroupPath + keyProductsIdPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, res.Code)
	assert.NotEmpty(suite.T(), res.Body.String())
}

func (suite *KeyProductTestSuite) TestProject_CreateKeyProduct_Ok() {
	body := &grpc.CreateOrUpdateKeyProductRequest{
		MerchantId:      bson.NewObjectId().Hex(),
		Name:            map[string]string{"en": "A", "ru": "А"},
		Description:     map[string]string{"en": "A", "ru": "А"},
		DefaultCurrency: "RUB",
		ProjectId:       bson.NewObjectId().Hex(),
		Sku:             "some_sku",
		Object:          "key_product",
	}

	b, err := json.Marshal(&body)
	assert.NoError(suite.T(), err)

	res, err := suite.caller.Builder().
		Method(http.MethodPost).
		Params(":key_product_id", body.Id).
		Path(common.AuthUserGroupPath + keyProductsPath).
		Init(test.ReqInitJSON()).
		BodyBytes(b).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusCreated, res.Code)
	assert.NotEmpty(suite.T(), res.Body.String())
}

func (suite *KeyProductTestSuite) TestProject_ChangeKeyProduct_ValidationError() {
	body := &grpc.CreateOrUpdateKeyProductRequest{
		Id:              bson.NewObjectId().Hex(),
		MerchantId:      bson.NewObjectId().Hex(),
		DefaultCurrency: "RUB",
		ProjectId:       bson.NewObjectId().Hex(),
		Sku:             "some_sku",
		Object:          "key_product",
	}

	b, err := json.Marshal(&body)
	assert.NoError(suite.T(), err)

	_, err = suite.caller.Builder().
		Method(http.MethodPut).
		Params(":key_product_id", body.Id).
		Path(common.AuthUserGroupPath + keyProductsIdPath).
		Init(test.ReqInitJSON()).
		BodyBytes(b).
		Exec(suite.T())

	assert.Error(suite.T(), err)
	e, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), 400, e.Code)
	assert.NotEmpty(suite.T(), e.Message)
}

func (suite *KeyProductTestSuite) TestProject_ChangeKeyProduct_Ok() {
	body := &grpc.CreateOrUpdateKeyProductRequest{
		Id:              bson.NewObjectId().Hex(),
		MerchantId:      bson.NewObjectId().Hex(),
		Name:            map[string]string{"en": "A", "ru": "А"},
		Description:     map[string]string{"en": "A", "ru": "А"},
		DefaultCurrency: "RUB",
		ProjectId:       bson.NewObjectId().Hex(),
		Sku:             "some_sku",
		Object:          "key_product",
	}

	b, err := json.Marshal(&body)
	assert.NoError(suite.T(), err)

	res, err := suite.caller.Builder().
		Method(http.MethodPut).
		Params(":key_product_id", body.Id).
		Path(common.AuthUserGroupPath + keyProductsIdPath).
		Init(test.ReqInitJSON()).
		BodyBytes(b).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, res.Code)
	assert.NotEmpty(suite.T(), res.Body.String())
}

func (suite *KeyProductTestSuite) TestProject_CreateKeyProduct_ValidationError() {
	body := &grpc.CreateOrUpdateKeyProductRequest{
		MerchantId:      bson.NewObjectId().Hex(),
		Description:     map[string]string{"en": "A", "ru": "А"},
		DefaultCurrency: "RUB",
		ProjectId:       bson.NewObjectId().Hex(),
		Sku:             "some_sku",
		Object:          "key_product",
	}

	b, err := json.Marshal(&body)
	assert.NoError(suite.T(), err)

	_, err = suite.caller.Builder().
		Method(http.MethodPost).
		Path(common.AuthUserGroupPath + keyProductsPath).
		Init(test.ReqInitJSON()).
		BodyBytes(b).
		Exec(suite.T())

	assert.Error(suite.T(), err)
	e, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), 400, e.Code)
	assert.NotEmpty(suite.T(), e.Message)
}

func (suite *KeyProductTestSuite) TestProject_getKeyProduct_ValidationError() {

	billingService := &billMock.BillingService{}
	billingService.On("GetKeyProduct", mock2.Anything, mock2.Anything).Return(&grpc.KeyProductResponse{}, nil)
	suite.router.dispatch.Services.Billing = billingService

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Path(common.AuthProjectGroupPath + keyProductsIdPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.Error(suite.T(), err)
	err2, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, err2.Code)
	assert.NotEmpty(suite.T(), err2.Message)
}

func (suite *KeyProductTestSuite) TestProject_getKeyProduct_BillingServer() {

	billingService := &billMock.BillingService{}
	billingService.On("GetKeyProductInfo", mock2.Anything, mock2.Anything).Return(nil, errors.New("error"))
	suite.router.dispatch.Services.Billing = billingService

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Params(":key_product_id", bson.NewObjectId().Hex()).
		Path(common.AuthProjectGroupPath + keyProductsIdPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.Error(suite.T(), err)
	err2, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusInternalServerError, err2.Code)
	assert.NotEmpty(suite.T(), err2.Message)
}

func (suite *KeyProductTestSuite) TestProject_getKeyProductWithCurrency_Ok() {

	billingService := &billMock.BillingService{}
	billingService.On("GetKeyProductInfo", mock2.Anything, mock2.Anything).Return(&grpc.GetKeyProductInfoResponse{
		Status: 200,
		KeyProduct: &grpc.KeyProductInfo{
			LongDescription: "Description",
			Name:            "Name",
			Platforms: []*grpc.PlatformPriceInfo{
				{Name: "Steam", Id: "steam", Price: &grpc.ProductPriceInfo{Currency: "USD", Amount: 10, Region: "USD"}},
			},
		},
	}, nil)
	suite.router.dispatch.Services.Billing = billingService

	res, err := suite.caller.Builder().
		Method(http.MethodGet).
		SetQueryParam("currency", "USD").
		Params(":key_product_id", bson.NewObjectId().Hex()).
		Path(common.AuthProjectGroupPath + keyProductsIdPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, res.Code)
	assert.NotEmpty(suite.T(), res.Body.String())
}

func (suite *KeyProductTestSuite) TestProject_getKeyProductWithCountry_Ok() {

	billingService := &billMock.BillingService{}
	billingService.On("GetKeyProductInfo", mock2.Anything, mock2.Anything).Return(&grpc.GetKeyProductInfoResponse{
		Status: 200,
		KeyProduct: &grpc.KeyProductInfo{
			LongDescription: "Description",
			Name:            "Name",
			Platforms: []*grpc.PlatformPriceInfo{
				{Name: "Steam", Id: "steam", Price: &grpc.ProductPriceInfo{Currency: "USD", Amount: 10, Region: "USD"}},
			},
		},
	}, nil)
	suite.router.dispatch.Services.Billing = billingService

	res, err := suite.caller.Builder().
		Method(http.MethodGet).
		SetQueryParam("country", "RUS").
		Params(":key_product_id", bson.NewObjectId().Hex()).
		Path(common.AuthProjectGroupPath + keyProductsIdPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, res.Code)
	assert.NotEmpty(suite.T(), res.Body.String())
}

func (suite *KeyProductTestSuite) TestProject_getKeyProduct_Ok() {

	billingService := &billMock.BillingService{}
	billingService.On("GetKeyProductInfo", mock2.Anything, mock2.Anything).Return(&grpc.GetKeyProductInfoResponse{
		Status: 200,
		KeyProduct: &grpc.KeyProductInfo{
			LongDescription: "Description",
			Name:            "Name",
			Platforms: []*grpc.PlatformPriceInfo{
				{Name: "Steam", Id: "steam", Price: &grpc.ProductPriceInfo{Currency: "USD", Amount: 10, Region: "USD"}},
			},
		},
	}, nil)
	suite.router.dispatch.Services.Billing = billingService

	res, err := suite.caller.Builder().
		Method(http.MethodGet).
		Params(":key_product_id", bson.NewObjectId().Hex()).
		Path(common.AuthProjectGroupPath + keyProductsIdPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, res.Code)
	assert.NotEmpty(suite.T(), res.Body.String())
}
