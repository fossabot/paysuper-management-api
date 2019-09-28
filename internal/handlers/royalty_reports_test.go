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

type RoyaltyReportsTestSuite struct {
	suite.Suite
	router *RoyaltyReportsRoute
	caller *test.EchoReqResCaller
}

func Test_RoyaltyReports(t *testing.T) {
	suite.Run(t, new(RoyaltyReportsTestSuite))
}

func (suite *RoyaltyReportsTestSuite) SetupTest() {
	var e error
	settings := test.DefaultSettings()
	srv := common.Services{
		Billing: mock.NewBillingServerOkMock(),
	}
	suite.caller, e = test.SetUp(settings, srv, func(set *test.TestSet, mw test.Middleware) common.Handlers {
		suite.router = NewRoyaltyReportsRoute(set.HandlerSet, set.GlobalConfig)
		return common.Handlers{
			suite.router,
		}
	})
	if e != nil {
		panic(e)
	}
}

func (suite *RoyaltyReportsTestSuite) TearDownTest() {}

func (suite *RoyaltyReportsTestSuite) TestRoyaltyReports_getRoyaltyReportsList() {

	res, err := suite.caller.Builder().
		SetQueryParam("id", "5ced34d689fce60bf4440829").
		SetQueryParam("merchant_id", "5ced34d689fce60bf444082b").
		Method(http.MethodGet).
		Path(common.AuthUserGroupPath + royaltyReportsPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), http.StatusOK, res.Code)
		assert.NotEmpty(suite.T(), res.Body.String())
	}
}

func (suite *RoyaltyReportsTestSuite) TestRoyaltyReports_getRoyaltyReport() {

	res, err := suite.caller.Builder().
		Params(":"+common.RequestParameterId, bson.NewObjectId().Hex()).
		Method(http.MethodGet).
		Path(common.AuthUserGroupPath + royaltyReportsIdPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), http.StatusOK, res.Code)
		assert.NotEmpty(suite.T(), res.Body.String())
	}
}

func (suite *RoyaltyReportsTestSuite) TestRoyaltyReports_listRoyaltyReportOrders() {

	res, err := suite.caller.Builder().
		SetQueryParam("limit", "100").
		SetQueryParam("offset", "200").
		Params(":"+common.RequestParameterId, bson.NewObjectId().Hex()).
		Method(http.MethodGet).
		Path(common.AuthUserGroupPath + royaltyReportsTransactionsPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), http.StatusOK, res.Code)
		assert.NotEmpty(suite.T(), res.Body.String())
	}
}

func (suite *RoyaltyReportsTestSuite) TestRoyaltyReports_MerchantReviewRoyaltyReport() {

	res, err := suite.caller.Builder().
		Params(":"+common.RequestParameterId, bson.NewObjectId().Hex()).
		Method(http.MethodPost).
		Path(common.AuthUserGroupPath + royaltyReportsAcceptPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), http.StatusNoContent, res.Code)
	}
}

func (suite *RoyaltyReportsTestSuite) TestRoyaltyReports_merchantDeclineRoyaltyReport() {

	res, err := suite.caller.Builder().
		Params(":"+common.RequestParameterId, bson.NewObjectId().Hex()).
		Method(http.MethodPost).
		Path(common.AuthUserGroupPath + royaltyReportsDeclinePath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), http.StatusNoContent, res.Code)
	}
}

func (suite *RoyaltyReportsTestSuite) TestRoyaltyReports_changeRoyaltyReport() {

	bodyJson := `{"status": "accepted", "correction": {"amount": 100500, "reason": "just for fun :)"}, "payout_id": "5bdc39a95d1e1100019fb7df"}`

	res, err := suite.caller.Builder().
		Params(":"+common.RequestParameterId, bson.NewObjectId().Hex()).
		Method(http.MethodPost).
		Path(common.AuthUserGroupPath + royaltyReportsChangePath).
		Init(test.ReqInitJSON()).
		BodyString(bodyJson).
		Exec(suite.T())

	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), http.StatusNoContent, res.Code)
	}
}
