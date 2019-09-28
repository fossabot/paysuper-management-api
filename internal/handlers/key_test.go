package handlers

import (
	"errors"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo/v4"
	billMock "github.com/paysuper/paysuper-billing-server/pkg/mocks"
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

type KeyTestSuite struct {
	suite.Suite
	router *KeyRoute
	caller *test.EchoReqResCaller
}

func Test_GetKeyById(t *testing.T) {
	suite.Run(t, new(KeyTestSuite))
}

func (suite *KeyTestSuite) SetupTest() {
	var e error
	settings := test.DefaultSettings()
	srv := common.Services{
		Billing: mock.NewBillingServerOkMock(),
	}
	suite.caller, e = test.SetUp(settings, srv, func(set *test.TestSet, mw test.Middleware) common.Handlers {
		suite.router = NewKeyRoute(set.HandlerSet, set.GlobalConfig)
		return common.Handlers{
			suite.router,
		}
	})
	if e != nil {
		panic(e)
	}
}

func (suite *KeyTestSuite) TearDownTest() {}

func (suite *KeyTestSuite) TestGetKeyById_Ok() {

	billingService := &billMock.BillingService{}
	billingService.On("GetKeyByID", mock2.Anything, mock2.Anything).Return(&grpc.GetKeyForOrderRequestResponse{
		Status: 200,
		Key: &billing.Key{
			Id:   bson.NewObjectId().Hex(),
			Code: "XXXX-YYYY-ZZZZ",
		},
	}, nil)

	suite.router.dispatch.Services.Billing = billingService

	res, err := suite.caller.Builder().
		Method(http.MethodGet).
		Params(":key_id", bson.NewObjectId().Hex()).
		Path(common.AuthUserGroupPath + keysIdPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, res.Code)
	assert.NotEmpty(suite.T(), res.Body.String())
}

func (suite *KeyTestSuite) TestGetKeyById_InternalError() {

	billingService := &billMock.BillingService{}
	billingService.On("GetKeyByID", mock2.Anything, mock2.Anything).Return(nil, errors.New("some error"))

	suite.router.dispatch.Services.Billing = billingService

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Params(":key_id", bson.NewObjectId().Hex()).
		Path(common.AuthUserGroupPath + keysIdPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.NotNil(suite.T(), err)
	httpErr := err.(*echo.HTTPError)
	assert.EqualValues(suite.T(), 500, httpErr.Code)
}

func (suite *KeyTestSuite) TestGetKeyById_ServiceError() {

	billingService := &billMock.BillingService{}
	billingService.On("GetKeyByID", mock2.Anything, mock2.Anything).Return(&grpc.GetKeyForOrderRequestResponse{
		Status: 404,
	}, nil)

	suite.router.dispatch.Services.Billing = billingService

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Params(":key_id", bson.NewObjectId().Hex()).
		Path(common.AuthUserGroupPath + keysIdPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.NotNil(suite.T(), err)
	httpErr := err.(*echo.HTTPError)
	assert.EqualValues(suite.T(), 404, httpErr.Code)
}

func (suite *KeyTestSuite) TestGetKeyById_ValidationError() {

	billingService := &billMock.BillingService{}
	billingService.On("GetKeyByID", mock2.Anything, mock2.Anything).Return(&grpc.GetKeyForOrderRequestResponse{
		Status: 200,
		Key: &billing.Key{
			Id:   bson.NewObjectId().Hex(),
			Code: "XXXX-YYYY-ZZZZ",
		},
	}, nil)

	suite.router.dispatch.Services.Billing = billingService

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Params(":key_id", " ").
		Path(common.AuthUserGroupPath + keysIdPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.NotNil(suite.T(), err)
	httpErr := err.(*echo.HTTPError)
	assert.EqualValues(suite.T(), 400, httpErr.Code)
}
