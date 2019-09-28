package handlers

import (
	"github.com/globalsign/mgo/bson"
	"github.com/paysuper/paysuper-management-api/internal/dispatcher/common"
	"github.com/paysuper/paysuper-management-api/internal/mock"
	"github.com/paysuper/paysuper-management-api/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/url"
	"testing"
)

type VatReportsTestSuite struct {
	suite.Suite
	router *VatReportsRoute
	caller *test.EchoReqResCaller
}

func Test_VatReports(t *testing.T) {
	suite.Run(t, new(VatReportsTestSuite))
}

func (suite *VatReportsTestSuite) SetupTest() {
	var e error
	settings := test.DefaultSettings()
	srv := common.Services{
		Billing: mock.NewBillingServerOkMock(),
	}
	suite.caller, e = test.SetUp(settings, srv, func(set *test.TestSet, mw test.Middleware) common.Handlers {
		suite.router = NewVatReportsRoute(set.HandlerSet, set.GlobalConfig)
		return common.Handlers{
			suite.router,
		}
	})
	if e != nil {
		panic(e)
	}
}

func (suite *VatReportsTestSuite) TearDownTest() {}

func (suite *VatReportsTestSuite) TestVatReports_getVatReportsDashboard() {
	res, err := suite.caller.Builder().
		Method(http.MethodGet).
		Path(common.AuthUserGroupPath + vatReportsPath).
		Init(test.ReqInitApplicationForm()).
		Exec(suite.T())

	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), http.StatusOK, res.Code)
		assert.NotEmpty(suite.T(), res.Body.String())
	}
}

func (suite *VatReportsTestSuite) TestVatReports_getVatReportsForCountry() {

	q := make(url.Values)
	q.Set("limit", "100")
	q.Set("offset", "200")

	res, err := suite.caller.Builder().
		Method(http.MethodGet).
		SetQueryParams(q).
		Params(":country", "ru").
		Path(common.AuthUserGroupPath + vatReportsCountryPath).
		Init(test.ReqInitApplicationForm()).
		Exec(suite.T())

	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), http.StatusOK, res.Code)
		assert.NotEmpty(suite.T(), res.Body.String())
	}
}

func (suite *VatReportsTestSuite) TestVatReports_getVatReportTransactions() {

	q := make(url.Values)
	q.Set("limit", "100")
	q.Set("offset", "200")

	res, err := suite.caller.Builder().
		Method(http.MethodGet).
		SetQueryParams(q).
		Params(":"+common.RequestParameterId, bson.NewObjectId().Hex()).
		Path(common.AuthUserGroupPath + vatReportsDetailsPath).
		Init(test.ReqInitApplicationForm()).
		Exec(suite.T())

	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), http.StatusOK, res.Code)
		assert.NotEmpty(suite.T(), res.Body.String())
	}
}

func (suite *VatReportsTestSuite) TestVatReports_updateVatReportStatus() {

	bodyJson := `{"status": "paid"}`

	res, err := suite.caller.Builder().
		Method(http.MethodPost).
		Params(":"+common.RequestParameterId, bson.NewObjectId().Hex()).
		Path(common.AuthUserGroupPath + vatReportsStatusPath).
		Init(test.ReqInitJSON()).
		BodyString(bodyJson).
		Exec(suite.T())

	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), http.StatusNoContent, res.Code)
	}
}
