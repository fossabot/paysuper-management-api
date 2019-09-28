package handlers

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/grpc"
	"github.com/paysuper/paysuper-management-api/internal/dispatcher/common"
	"github.com/paysuper/paysuper-management-api/internal/mock"
	"github.com/paysuper/paysuper-management-api/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/url"
	"testing"
)

type ZipCodeTestSuite struct {
	suite.Suite
	router *ZipCodeRoute
	caller *test.EchoReqResCaller
}

func Test_ZipCode(t *testing.T) {
	suite.Run(t, new(ZipCodeTestSuite))
}

func (suite *ZipCodeTestSuite) SetupTest() {
	var e error
	settings := test.DefaultSettings()
	srv := common.Services{
		Billing: mock.NewBillingServerOkMock(),
	}
	suite.caller, e = test.SetUp(settings, srv, func(set *test.TestSet, mw test.Middleware) common.Handlers {
		suite.router = NewZipCodeRoute(set.HandlerSet, set.GlobalConfig)
		return common.Handlers{
			suite.router,
		}
	})
	if e != nil {
		panic(e)
	}
}

func (suite *ZipCodeTestSuite) TearDownTest() {}

func (suite *ZipCodeTestSuite) TestCheckZip_Ok() {
	q := make(url.Values)
	q.Set("country", "US")
	q.Set("zip", "98")

	res, err := suite.caller.Builder().
		Path(common.AuthProjectGroupPath + zipCodePath).
		SetQueryParams(q).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, res.Code)
	assert.NotEmpty(suite.T(), res.Body.String())

	data := &grpc.FindByZipCodeResponse{}
	err = json.Unmarshal(res.Body.Bytes(), data)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int32(1), data.Count)
	assert.Len(suite.T(), data.Items, 1)
}

func (suite *ZipCodeTestSuite) TestCheckZip_BindError() {
	q := make(url.Values)
	q.Set("country", "US")
	q.Set("zip", "98")
	q.Set("limit", "qwerty")

	_, err := suite.caller.Builder().
		Path(common.AuthProjectGroupPath + zipCodePath).
		SetQueryParams(q).
		Exec(suite.T())
	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorRequestParamsIncorrect, httpErr.Message)
}

func (suite *ZipCodeTestSuite) TestCheckZip_ValidateError() {
	q := make(url.Values)
	q.Set("zip", "98")

	_, err := suite.caller.Builder().
		Path(common.AuthProjectGroupPath + zipCodePath).
		SetQueryParams(q).
		Exec(suite.T())
	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Regexp(suite.T(), common.NewValidationError("Country"), httpErr.Message)
}

func (suite *ZipCodeTestSuite) TestCheckZip_BillingServerError() {
	q := make(url.Values)
	q.Set("country", "US")
	q.Set("zip", "98")

	suite.router.dispatch.Services.Billing = mock.NewBillingServerErrorMock()
	_, err := suite.caller.Builder().
		Path(common.AuthProjectGroupPath + zipCodePath).
		SetQueryParams(q).
		Exec(suite.T())
	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusInternalServerError, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorUnknown, httpErr.Message)
}
