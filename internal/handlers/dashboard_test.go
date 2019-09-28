package handlers

import (
	"errors"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo/v4"
	"github.com/paysuper/paysuper-billing-server/pkg"
	billingMocks "github.com/paysuper/paysuper-billing-server/pkg/mocks"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/grpc"
	"github.com/paysuper/paysuper-management-api/internal/dispatcher/common"
	"github.com/paysuper/paysuper-management-api/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/url"
	"testing"
)

type DashboardTestSuite struct {
	suite.Suite
	router *DashboardRoute
	caller *test.EchoReqResCaller
}

func Test_Dashboard(t *testing.T) {
	suite.Run(t, new(DashboardTestSuite))
}

func (suite *DashboardTestSuite) SetupTest() {
	bs := &billingMocks.BillingService{}
	bs.On("GetDashboardMainReport", mock.Anything, mock.Anything, mock.Anything).
		Return(&grpc.GetDashboardMainResponse{Status: pkg.ResponseStatusOk, Item: &grpc.DashboardMainReport{}}, nil)
	bs.On("GetDashboardRevenueDynamicsReport", mock.Anything, mock.Anything, mock.Anything).
		Return(
			&grpc.GetDashboardRevenueDynamicsReportResponse{
				Status: pkg.ResponseStatusOk,
				Item:   &grpc.DashboardRevenueDynamicReport{},
			},
			nil,
		)
	bs.On("GetDashboardBaseReport", mock.Anything, mock.Anything, mock.Anything).
		Return(
			&grpc.GetDashboardBaseReportResponse{
				Status: pkg.ResponseStatusOk,
				Item:   &grpc.DashboardBaseReports{},
			},
			nil,
		)

	var e error
	settings := test.DefaultSettings()
	srv := common.Services{
		Billing: bs,
	}
	suite.caller, e = test.SetUp(settings, srv, func(set *test.TestSet, mw test.Middleware) common.Handlers {
		suite.router = NewDashboardRoute(set.HandlerSet, set.GlobalConfig)
		return common.Handlers{
			suite.router,
		}
	})
	if e != nil {
		panic(e)
	}
}

func (suite *DashboardTestSuite) TearDownTest() {}

func (suite *DashboardTestSuite) TestDashboard_GetMainReports_Ok() {
	q := make(url.Values)
	q.Set("period", "current_month")
	_, err := suite.caller.Builder().
		Path(common.AuthUserGroupPath+dashboardMainPath).
		SetQueryParams(q).
		Params(":"+common.RequestParameterId, bson.NewObjectId().Hex()).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
}

func (suite *DashboardTestSuite) TestDashboard_GetMainReports_ValidationError() {
	q := make(url.Values)
	q.Set("period", "123")
	_, err := suite.caller.Builder().
		Path(common.AuthUserGroupPath+dashboardMainPath).
		SetQueryParams(q).
		Params(":"+common.RequestParameterId, bson.NewObjectId().Hex()).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorIncorrectPeriod, httpErr.Message)
}

func (suite *DashboardTestSuite) TestDashboard_GetMainReports_BillingServerSystemError() {
	q := make(url.Values)
	q.Set("period", "current_month")

	bs := &billingMocks.BillingService{}
	bs.On("GetDashboardMainReport", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, errors.New("some error"))
	suite.router.dispatch.Services.Billing = bs

	_, err := suite.caller.Builder().
		Path(common.AuthUserGroupPath+dashboardMainPath).
		SetQueryParams(q).
		Params(":"+common.RequestParameterId, bson.NewObjectId().Hex()).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusInternalServerError, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorUnknown, httpErr.Message)
}

func (suite *DashboardTestSuite) TestDashboard_GetMainReports_BillingServerReturnError() {
	q := make(url.Values)
	q.Set("period", "current_month")

	bs := &billingMocks.BillingService{}
	bs.On("GetDashboardMainReport", mock.Anything, mock.Anything, mock.Anything).
		Return(
			&grpc.GetDashboardMainResponse{
				Status:  pkg.ResponseStatusBadData,
				Message: &grpc.ResponseErrorMessage{Message: "some error"},
			},
			nil,
		)
	suite.router.dispatch.Services.Billing = bs

	_, err := suite.caller.Builder().
		Path(common.AuthUserGroupPath+dashboardMainPath).
		SetQueryParams(q).
		Params(":"+common.RequestParameterId, bson.NewObjectId().Hex()).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)

	msg, ok := httpErr.Message.(*grpc.ResponseErrorMessage)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), "some error", msg.Message)
}

func (suite *DashboardTestSuite) TestDashboard_GetRevenueDynamicsReport_Ok() {
	q := make(url.Values)
	q.Set("period", "current_month")
	_, err := suite.caller.Builder().
		Path(common.AuthUserGroupPath+dashboardRevenueDynamicsPath).
		SetQueryParams(q).
		Params(":"+common.RequestParameterId, bson.NewObjectId().Hex()).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
}

func (suite *DashboardTestSuite) TestDashboard_GetRevenueDynamicsReport_ValidationError() {
	q := make(url.Values)
	q.Set("period", "123")

	_, err := suite.caller.Builder().
		Path(common.AuthUserGroupPath+dashboardRevenueDynamicsPath).
		SetQueryParams(q).
		Params(":"+common.RequestParameterId, bson.NewObjectId().Hex()).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorIncorrectPeriod, httpErr.Message)
}

func (suite *DashboardTestSuite) TestDashboard_GetRevenueDynamicsReport_BillingServerSystemError() {
	q := make(url.Values)
	q.Set("period", "current_month")

	bs := &billingMocks.BillingService{}
	bs.On("GetDashboardRevenueDynamicsReport", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, errors.New("some error"))
	suite.router.dispatch.Services.Billing = bs

	_, err := suite.caller.Builder().
		Path(common.AuthUserGroupPath+dashboardRevenueDynamicsPath).
		SetQueryParams(q).
		Params(":"+common.RequestParameterId, bson.NewObjectId().Hex()).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusInternalServerError, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorUnknown, httpErr.Message)
}

func (suite *DashboardTestSuite) TestDashboard_GetRevenueDynamicsReport_BillingServerReturnError() {
	q := make(url.Values)
	q.Set("period", "current_month")

	bs := &billingMocks.BillingService{}
	bs.On("GetDashboardRevenueDynamicsReport", mock.Anything, mock.Anything, mock.Anything).
		Return(
			&grpc.GetDashboardRevenueDynamicsReportResponse{
				Status:  pkg.ResponseStatusBadData,
				Message: &grpc.ResponseErrorMessage{Message: "some error"},
			},
			nil,
		)
	suite.router.dispatch.Services.Billing = bs

	_, err := suite.caller.Builder().
		Path(common.AuthUserGroupPath+dashboardRevenueDynamicsPath).
		SetQueryParams(q).
		Params(":"+common.RequestParameterId, bson.NewObjectId().Hex()).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)

	msg, ok := httpErr.Message.(*grpc.ResponseErrorMessage)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), "some error", msg.Message)
}

func (suite *DashboardTestSuite) TestDashboard_GetBaseReports_Ok() {
	q := make(url.Values)
	q.Set("period", "current_month")

	_, err := suite.caller.Builder().
		Path(common.AuthUserGroupPath+dashboardBasePath).
		SetQueryParams(q).
		Params(":"+common.RequestParameterId, bson.NewObjectId().Hex()).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
}

func (suite *DashboardTestSuite) TestDashboard_GetBaseReports_ValidationError() {
	q := make(url.Values)
	q.Set("period", "123")

	_, err := suite.caller.Builder().
		Path(common.AuthUserGroupPath+dashboardBasePath).
		SetQueryParams(q).
		Params(":"+common.RequestParameterId, bson.NewObjectId().Hex()).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorIncorrectPeriod, httpErr.Message)
}

func (suite *DashboardTestSuite) TestDashboard_GetBaseReports_BillingServerSystemError() {
	q := make(url.Values)
	q.Set("period", "current_month")

	bs := &billingMocks.BillingService{}
	bs.On("GetDashboardBaseReport", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, errors.New("some error"))
	suite.router.dispatch.Services.Billing = bs

	_, err := suite.caller.Builder().
		Path(common.AuthUserGroupPath+dashboardBasePath).
		SetQueryParams(q).
		Params(":"+common.RequestParameterId, bson.NewObjectId().Hex()).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusInternalServerError, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorUnknown, httpErr.Message)
}

func (suite *DashboardTestSuite) TestDashboard_GetBaseReports_BillingServerReturnError() {
	q := make(url.Values)
	q.Set("period", "current_month")

	bs := &billingMocks.BillingService{}
	bs.On("GetDashboardBaseReport", mock.Anything, mock.Anything, mock.Anything).
		Return(
			&grpc.GetDashboardBaseReportResponse{
				Status:  pkg.ResponseStatusBadData,
				Message: &grpc.ResponseErrorMessage{Message: "some error"},
			},
			nil,
		)
	suite.router.dispatch.Services.Billing = bs

	_, err := suite.caller.Builder().
		Path(common.AuthUserGroupPath+dashboardBasePath).
		SetQueryParams(q).
		Params(":"+common.RequestParameterId, bson.NewObjectId().Hex()).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)

	msg, ok := httpErr.Message.(*grpc.ResponseErrorMessage)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), "some error", msg.Message)
}
