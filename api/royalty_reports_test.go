package api

import (
	"github.com/labstack/echo/v4"
	"github.com/paysuper/paysuper-management-api/internal/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type RoyaltyReportsTestSuite struct {
	suite.Suite
	router *royaltyReportsRoute
	api    *Api
}

func Test_RoyaltyReports(t *testing.T) {
	suite.Run(t, new(RoyaltyReportsTestSuite))
}

func (suite *RoyaltyReportsTestSuite) SetupTest() {
	suite.api = &Api{
		Http:           echo.New(),
		validate:       validator.New(),
		billingService: mock.NewBillingServerOkMock(),
		authUser: &AuthUser{
			Id: "ffffffffffffffffffffffff",
		},
	}

	suite.api.authUserRouteGroup = suite.api.Http.Group(apiAuthUserGroupPath)
	suite.router = &royaltyReportsRoute{Api: suite.api}
}

func (suite *RoyaltyReportsTestSuite) TearDownTest() {}

func (suite *RoyaltyReportsTestSuite) TestRoyaltyReports_getRoyaltyReportsList() {
	e := echo.New()
	q := make(url.Values)
	q.Set("id", "5ced34d689fce60bf4440829")
	q.Set("merchant_id", "5ced34d689fce60bf444082b")
	req := httptest.NewRequest(http.MethodGet, "/royalty_reports?"+q.Encode(), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := e.NewContext(req, rsp)

	ctx.SetPath("/royalty_reports?" + q.Encode())

	err := suite.router.getRoyaltyReportsList(ctx)

	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), http.StatusOK, rsp.Code)
		assert.NotEmpty(suite.T(), rsp.Body.String())
	}
}

func (suite *RoyaltyReportsTestSuite) TestRoyaltyReports_listRoyaltyReportOrders() {
	e := echo.New()
	q := make(url.Values)
	q.Set("limit", "100")
	q.Set("offset", "200")
	req := httptest.NewRequest(http.MethodGet, "/royalty_reports/details/5ced34d689fce60bf4440829?"+q.Encode(), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := e.NewContext(req, rsp)

	ctx.SetPath("/royalty_reports/details/:" + requestParameterId + "?" + q.Encode())
	ctx.SetParamNames(requestParameterId)
	ctx.SetParamValues("5ced34d689fce60bf4440829")

	err := suite.router.listRoyaltyReportOrders(ctx)

	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), http.StatusOK, rsp.Code)
		assert.NotEmpty(suite.T(), rsp.Body.String())
	}
}
