package handlers

import (
	"github.com/globalsign/mgo/bson"
	"github.com/paysuper/paysuper-management-api/internal/dispatcher/common"
	"github.com/paysuper/paysuper-management-api/internal/mock"
	"github.com/paysuper/paysuper-management-api/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type BalanceTestSuite struct {
	suite.Suite
	router *BalanceRoute
	caller *test.EchoReqResCaller
}

func Test_Balance(t *testing.T) {
	suite.Run(t, new(BalanceTestSuite))
}

func (suite *BalanceTestSuite) SetupTest() {
	user := &common.AuthUser{
		Id:    "ffffffffffffffffffffffff",
		Email: "test@unit.test",
	}
	var e error
	settings := test.DefaultSettings()
	srv := common.Services{
		Billing: mock.NewBillingServerOkMock(),
	}
	suite.caller, e = test.SetUp(settings, srv, func(set *test.TestSet, mw test.Middleware) common.Handlers {
		mw.Pre(test.PreAuthUserMiddleware(user))
		suite.router = NewBalanceRoute(set.HandlerSet, set.GlobalConfig)
		return common.Handlers{
			suite.router,
		}
	})
	if e != nil {
		panic(e)
	}
}

func (suite *BalanceTestSuite) TearDownTest() {}

func (suite *BalanceTestSuite) TestBalance_Ok_getBalanceForCurrentMerchant() {

	res, err := suite.caller.Builder().
		Method(http.MethodGet).
		Path(common.AuthUserGroupPath + balancePath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), http.StatusOK, res.Code)
		assert.NotEmpty(suite.T(), res.Body.String())
	}
}

func (suite *BalanceTestSuite) TestBalance_Ok_getBalanceForOtherMerchant() {

	res, err := suite.caller.Builder().
		Method(http.MethodGet).
		Params(":"+common.RequestParameterMerchantId, bson.NewObjectId().Hex()).
		Path(common.AuthUserGroupPath + balanceMerchantPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), http.StatusOK, res.Code)
		assert.NotEmpty(suite.T(), res.Body.String())
	}
}

func (suite *BalanceTestSuite) TestBalance_Fail_getBalanceForMerchantNotFound() {
	assert.Equal(suite.T(), common.TestStubImplementMe, "implement me!")
}

func (suite *BalanceTestSuite) TestBalance_Fail_getBalanceForMerchantIdInvalid() {
	assert.Equal(suite.T(), common.TestStubImplementMe, "implement me!")
}

func (suite *BalanceTestSuite) TestBalance_Fail_GrpcError() {
	assert.Equal(suite.T(), common.TestStubImplementMe, "implement me!")
}

func (suite *BalanceTestSuite) TestBalance_Fail_ResponseStatusNotOk() {
	assert.Equal(suite.T(), common.TestStubImplementMe, "implement me!")
}
