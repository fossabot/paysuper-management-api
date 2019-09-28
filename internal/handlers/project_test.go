package handlers

import (
	"encoding/json"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo/v4"
	"github.com/paysuper/paysuper-billing-server/pkg"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/billing"
	"github.com/paysuper/paysuper-management-api/internal/dispatcher/common"
	"github.com/paysuper/paysuper-management-api/internal/mock"
	"github.com/paysuper/paysuper-management-api/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type ProjectTestSuite struct {
	suite.Suite
	router *ProjectRoute
	caller *test.EchoReqResCaller
}

func Test_Project(t *testing.T) {
	suite.Run(t, new(ProjectTestSuite))
}

func (suite *ProjectTestSuite) SetupTest() {
	var e error
	settings := test.DefaultSettings()
	srv := common.Services{
		Billing: mock.NewBillingServerOkMock(),
	}
	suite.caller, e = test.SetUp(settings, srv, func(set *test.TestSet, mw test.Middleware) common.Handlers {
		suite.router = NewProjectRoute(set.HandlerSet, set.GlobalConfig)
		return common.Handlers{
			suite.router,
		}
	})
	if e != nil {
		panic(e)
	}
}

func (suite *ProjectTestSuite) TearDownTest() {}

func (suite *ProjectTestSuite) TestProject_CreateProject_Ok() {
	body := &billing.Project{
		MerchantId:         bson.NewObjectId().Hex(),
		Name:               map[string]string{"en": "A", "ru": "А"},
		CallbackCurrency:   "RUB",
		CallbackProtocol:   pkg.ProjectCallbackProtocolEmpty,
		LimitsCurrency:     "RUB",
		MinPaymentAmount:   0,
		MaxPaymentAmount:   15000,
		IsProductsCheckout: false,
	}

	b, err := json.Marshal(&body)
	assert.NoError(suite.T(), err)

	res, err := suite.caller.Builder().
		Method(http.MethodPost).
		Path(common.AuthUserGroupPath + projectsPath).
		Init(test.ReqInitJSON()).
		BodyBytes(b).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusCreated, res.Code)
	assert.NotEmpty(suite.T(), res.Body.String())
}

func (suite *ProjectTestSuite) TestProject_CreateProject_BindError() {
	body := `{"name": "qwerty"}`

	_, err := suite.caller.Builder().
		Method(http.MethodPost).
		Path(common.AuthUserGroupPath + projectsPath).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	//assert.Equal(suite.T(), ErrorRequestParamsIncorrect, httpErr.Message)
}

func (suite *ProjectTestSuite) TestProject_CreateProject_ValidationError() {
	body := &billing.Project{
		MerchantId:         bson.NewObjectId().Hex(),
		Name:               map[string]string{"en": "A", "ru": "А"},
		CallbackCurrency:   "RUB",
		CallbackProtocol:   pkg.ProjectCallbackProtocolEmpty,
		LimitsCurrency:     "RUB",
		MinPaymentAmount:   -100,
		MaxPaymentAmount:   15000,
		IsProductsCheckout: false,
	}

	b, err := json.Marshal(&body)
	assert.NoError(suite.T(), err)

	_, err = suite.caller.Builder().
		Method(http.MethodPost).
		Path(common.AuthUserGroupPath + projectsPath).
		Init(test.ReqInitJSON()).
		BodyBytes(b).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Regexp(suite.T(), common.NewValidationError("MinPaymentAmount"), httpErr.Message)
}

func (suite *ProjectTestSuite) TestProject_CreateProject_BillingServerError() {
	body := &billing.Project{
		MerchantId:         bson.NewObjectId().Hex(),
		Name:               map[string]string{"en": "A", "ru": "А"},
		CallbackCurrency:   "RUB",
		CallbackProtocol:   pkg.ProjectCallbackProtocolEmpty,
		LimitsCurrency:     "RUB",
		MinPaymentAmount:   100,
		MaxPaymentAmount:   15000,
		IsProductsCheckout: false,
	}

	b, err := json.Marshal(&body)
	assert.NoError(suite.T(), err)

	suite.router.dispatch.Services.Billing = mock.NewBillingServerSystemErrorMock()

	_, err = suite.caller.Builder().
		Method(http.MethodPost).
		Path(common.AuthUserGroupPath + projectsPath).
		Init(test.ReqInitJSON()).
		BodyBytes(b).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusInternalServerError, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorUnknown, httpErr.Message)
}

func (suite *ProjectTestSuite) TestProject_CreateProject_BillingServerResultError() {
	body := &billing.Project{
		MerchantId:         bson.NewObjectId().Hex(),
		Name:               map[string]string{"en": "A", "ru": "А"},
		CallbackCurrency:   "RUB",
		CallbackProtocol:   pkg.ProjectCallbackProtocolEmpty,
		LimitsCurrency:     "RUB",
		MinPaymentAmount:   100,
		MaxPaymentAmount:   15000,
		IsProductsCheckout: false,
	}

	b, err := json.Marshal(&body)
	assert.NoError(suite.T(), err)

	suite.router.dispatch.Services.Billing = mock.NewBillingServerErrorMock()

	_, err = suite.caller.Builder().
		Method(http.MethodPost).
		Path(common.AuthUserGroupPath + projectsPath).
		Init(test.ReqInitJSON()).
		BodyBytes(b).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), mock.SomeError, httpErr.Message)
}

func (suite *ProjectTestSuite) TestProject_UpdateProject_Ok() {
	body := `{"min_payment_amount": 10}`

	res, err := suite.caller.Builder().
		Method(http.MethodPatch).
		Params(":"+common.RequestParameterId, bson.NewObjectId().Hex()).
		Path(common.AuthUserGroupPath + projectsIdPath).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, res.Code)
	assert.NotEmpty(suite.T(), res.Body.String())
}

func (suite *ProjectTestSuite) TestProject_UpdateProject_BindError() {
	body := `{"name": "qwerty"}`

	_, err := suite.caller.Builder().
		Method(http.MethodPatch).
		Params(":"+common.RequestParameterId, bson.NewObjectId().Hex()).
		Path(common.AuthUserGroupPath + projectsIdPath).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	//assert.Equal(suite.T(), ErrorRequestParamsIncorrect, httpErr.Message)
}

func (suite *ProjectTestSuite) TestProject_UpdateProject_ValidationError() {
	body := `{"min_payment_amount": -10}`

	_, err := suite.caller.Builder().
		Method(http.MethodPatch).
		Params(":"+common.RequestParameterId, bson.NewObjectId().Hex()).
		Path(common.AuthUserGroupPath + projectsIdPath).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Regexp(suite.T(), common.NewValidationError("MinPaymentAmount"), httpErr.Message)
}

func (suite *ProjectTestSuite) TestProject_UpdateProject_BillingServerError() {
	body := `{"min_payment_amount": 10}`

	suite.router.dispatch.Services.Billing = mock.NewBillingServerSystemErrorMock()

	_, err := suite.caller.Builder().
		Method(http.MethodPatch).
		Params(":"+common.RequestParameterId, mock.SomeMerchantId).
		Path(common.AuthUserGroupPath + projectsIdPath).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusInternalServerError, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorUnknown, httpErr.Message)
}

func (suite *ProjectTestSuite) TestProject_UpdateProject_BillingServerResultError() {
	body := `{"min_payment_amount": 10}`

	suite.router.dispatch.Services.Billing = mock.NewBillingServerErrorMock()

	_, err := suite.caller.Builder().
		Method(http.MethodPatch).
		Params(":"+common.RequestParameterId, mock.SomeMerchantId).
		Path(common.AuthUserGroupPath + projectsIdPath).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), mock.SomeError, httpErr.Message)
}

func (suite *ProjectTestSuite) TestProject_GetProject_Ok() {

	res, err := suite.caller.Builder().
		Method(http.MethodGet).
		Params(":"+common.RequestParameterId, bson.NewObjectId().Hex()).
		Path(common.AuthUserGroupPath + projectsIdPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, res.Code)
	assert.NotEmpty(suite.T(), res.Body.String())
}

func (suite *ProjectTestSuite) TestProject_GetProject_ValidationError() {

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Path(common.AuthUserGroupPath + projectsIdPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Regexp(suite.T(), common.NewValidationError("ProjectId"), httpErr.Message)
}

func (suite *ProjectTestSuite) TestProject_GetProject_BillingServerError() {

	suite.router.dispatch.Services.Billing = mock.NewBillingServerSystemErrorMock()

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Params(":"+common.RequestParameterId, bson.NewObjectId().Hex()).
		Path(common.AuthUserGroupPath + projectsIdPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusInternalServerError, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorUnknown, httpErr.Message)
}

func (suite *ProjectTestSuite) TestProject_GetProject_BillingServerResultError() {

	suite.router.dispatch.Services.Billing = mock.NewBillingServerErrorMock()

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Params(":"+common.RequestParameterId, bson.NewObjectId().Hex()).
		Path(common.AuthUserGroupPath + projectsIdPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), mock.SomeError, httpErr.Message)
}

func (suite *ProjectTestSuite) TestProject_ListProjects_Ok() {

	res, err := suite.caller.Builder().
		Method(http.MethodGet).
		SetQueryParam(common.RequestParameterLimit, "-100").
		Path(common.AuthUserGroupPath + projectsPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, res.Code)
	assert.NotEmpty(suite.T(), res.Body.String())
}

func (suite *ProjectTestSuite) TestProject_ListProjects_BindError() {

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		SetQueryParam(common.RequestParameterLimit, "qwerty").
		Path(common.AuthUserGroupPath + projectsPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorRequestParamsIncorrect, httpErr.Message)
}

func (suite *ProjectTestSuite) TestProject_ListProjects_ValidationError() {

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		SetQueryParam(common.RequestParameterOffset, "-100").
		Path(common.AuthUserGroupPath + projectsPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Regexp(suite.T(), common.NewValidationError("Offset"), httpErr.Message)
}

func (suite *ProjectTestSuite) TestProject_ListProjects_BillingServerError() {

	suite.router.dispatch.Services.Billing = mock.NewBillingServerSystemErrorMock()

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		SetQueryParam(common.RequestParameterLimit, "100").
		Path(common.AuthUserGroupPath + projectsPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusInternalServerError, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorUnknown, httpErr.Message)
}

func (suite *ProjectTestSuite) TestProject_DeleteProject_Ok() {

	res, err := suite.caller.Builder().
		Method(http.MethodDelete).
		Params(":"+common.RequestParameterId, bson.NewObjectId().Hex()).
		Path(common.AuthUserGroupPath + projectsIdPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, res.Code)
	assert.NotEmpty(suite.T(), res.Body.String())
}

func (suite *ProjectTestSuite) TestProject_DeleteProject_ValidateError() {

	_, err := suite.caller.Builder().
		Method(http.MethodDelete).
		Path(common.AuthUserGroupPath + projectsIdPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Regexp(suite.T(), common.NewValidationError("ProjectId"), httpErr.Message)
}

func (suite *ProjectTestSuite) TestProject_DeleteProject_BillingServerError() {

	suite.router.dispatch.Services.Billing = mock.NewBillingServerSystemErrorMock()

	_, err := suite.caller.Builder().
		Method(http.MethodDelete).
		Params(":"+common.RequestParameterId, bson.NewObjectId().Hex()).
		Path(common.AuthUserGroupPath + projectsIdPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusInternalServerError, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorUnknown, httpErr.Message)
}

func (suite *ProjectTestSuite) TestProject_DeleteProject_BillingServerResultError() {

	suite.router.dispatch.Services.Billing = mock.NewBillingServerErrorMock()

	_, err := suite.caller.Builder().
		Method(http.MethodDelete).
		Params(":"+common.RequestParameterId, bson.NewObjectId().Hex()).
		Path(common.AuthUserGroupPath + projectsIdPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), mock.SomeError, httpErr.Message)
}
