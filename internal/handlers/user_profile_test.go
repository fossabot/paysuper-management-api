package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/grpc"
	"github.com/paysuper/paysuper-management-api/internal/dispatcher/common"
	"github.com/paysuper/paysuper-management-api/internal/mock"
	"github.com/paysuper/paysuper-management-api/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type UserProfileTestSuite struct {
	suite.Suite
	router *UserProfileRoute
	caller *test.EchoReqResCaller
}

func Test_UserProfile(t *testing.T) {
	suite.Run(t, new(UserProfileTestSuite))
}

func (suite *UserProfileTestSuite) SetupTest() {
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
		suite.router = NewUserProfileRoute(set.HandlerSet, set.GlobalConfig)
		return common.Handlers{
			suite.router,
		}
	})
	if e != nil {
		panic(e)
	}
}

func (suite *UserProfileTestSuite) TearDownTest() {}

func (suite *UserProfileTestSuite) TestUserProfile_GetUserProfile_Ok() {
	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Params(":"+common.RequestParameterId, "qwerty").
		Path(common.AuthUserGroupPath + userProfilePath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
}

func (suite *UserProfileTestSuite) TestUserProfile_GetUserProfile_ValidationError() {

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Params(":"+common.RequestParameterId, "qwerty").
		Path(common.AuthUserGroupPath + userProfilePathId).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.Error(suite.T(), err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)

	msg, ok := httpErr.Message.(*grpc.ResponseErrorMessage)
	assert.True(suite.T(), ok)
	assert.Regexp(suite.T(), "ProfileId", msg.Details)
}

func (suite *UserProfileTestSuite) TestUserProfile_GetUserProfile_BillingServerSystemError() {
	suite.router.dispatch.Services.Billing = mock.NewBillingServerSystemErrorMock()

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Path(common.AuthUserGroupPath + userProfilePath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusInternalServerError, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorUnknown, httpErr.Message)
}

func (suite *UserProfileTestSuite) TestUserProfile_GetUserProfile_BillingServerReturnError() {
	suite.router.dispatch.Services.Billing = mock.NewBillingServerErrorMock()

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Path(common.AuthUserGroupPath + userProfilePath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), mock.SomeError, httpErr.Message)
}

func (suite *UserProfileTestSuite) TestUserProfile_SetUserProfile_Ok() {
	body := `{
		"personal": {"first_name": "unit test", "last_name": "test-unit", "position": "Software Developer"},
		"company": {
			"company_name": "Unit Test.-444", 
			"website": "http://localhost",
			"annual_income": {"from": 0, "to": 1000},
			"number_of_employees": {"from": 1, "to": 10},
			"kind_of_activity": "other"
		}
	}`

	_, err := suite.caller.Builder().
		Method(http.MethodPatch).
		Path(common.AuthUserGroupPath + userProfilePath).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
}

func (suite *UserProfileTestSuite) TestUserProfile_SetUserProfile_BindError() {
	body := `{"help": {"product_promotion_and_development": "unit test"}}`

	_, err := suite.caller.Builder().
		Method(http.MethodPatch).
		Path(common.AuthUserGroupPath + userProfilePath).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorRequestParamsIncorrect, httpErr.Message)
}

func (suite *UserProfileTestSuite) TestUserProfile_SetUserProfile_ValidationError() {
	body := `{"help": {"product_promotion_and_development": false}}`

	reqInit := func(request *http.Request, middleware test.Middleware) {
		middleware.Pre(test.PreAuthUserMiddleware(&common.AuthUser{}))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}

	_, err := suite.caller.Builder().
		Method(http.MethodPatch).
		Path(common.AuthUserGroupPath + userProfilePath).
		Init(reqInit).
		BodyString(body).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)

	err1, ok := httpErr.Message.(*grpc.ResponseErrorMessage)
	assert.True(suite.T(), ok)

	assert.Regexp(suite.T(), "UserId", err1.Details)
}

func (suite *UserProfileTestSuite) TestUserProfile_SetUserProfile_ValidationUserNameError() {
	body := `{"personal": {"first_name": "unit test♂", "last_name": "test-unit", "position": "Software Developer"}}`

	_, err := suite.caller.Builder().
		Method(http.MethodPatch).
		Path(common.AuthUserGroupPath + userProfilePath).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)

	err1, ok := httpErr.Message.(*grpc.ResponseErrorMessage)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), common.ErrorMessageIncorrectFirstName, err1)
	assert.Regexp(suite.T(), "Name", err1.Details)
}

func (suite *UserProfileTestSuite) TestUserProfile_SetUserProfile_ValidationUserPositionError() {
	body := `{"personal": {"first_name": "unit test", "last_name": "test-unit", "position": "qwerty"}}`

	_, err := suite.caller.Builder().
		Method(http.MethodPatch).
		Path(common.AuthUserGroupPath + userProfilePath).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)

	err1, ok := httpErr.Message.(*grpc.ResponseErrorMessage)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), common.ErrorMessageIncorrectPosition, err1)
	assert.Regexp(suite.T(), "Position", err1.Details)
}

func (suite *UserProfileTestSuite) TestUserProfile_SetUserProfile_ValidationAnnualIncomeError() {
	body := `{
		"company": {
			"company_name": "Unit Test", 
			"website": "http://localhost", 
			"annual_income": {"from": 78898, "to": 9998},
			"number_of_employees": {"from": 1, "to": 10}
		}
	}`

	_, err := suite.caller.Builder().
		Method(http.MethodPatch).
		Path(common.AuthUserGroupPath + userProfilePath).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)

	err1, ok := httpErr.Message.(*grpc.ResponseErrorMessage)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), common.ErrorMessageIncorrectAnnualIncome, err1)
	assert.Regexp(suite.T(), "AnnualIncome", err1.Details)
}

func (suite *UserProfileTestSuite) TestUserProfile_SetUserProfile_ValidationNumberOfEmployeesError() {
	body := `{
		"company": {
			"company_name": "Unit Test", 
			"website": "http://localhost", 
			"annual_income": {"from": 0, "to": 1000},
			"number_of_employees": {"from": 23872, "to": 129}
		}
	}`

	_, err := suite.caller.Builder().
		Method(http.MethodPatch).
		Path(common.AuthUserGroupPath + userProfilePath).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)

	err1, ok := httpErr.Message.(*grpc.ResponseErrorMessage)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), common.ErrorMessageIncorrectNumberOfEmployees, err1)
	assert.Regexp(suite.T(), "NumberOfEmployees", err1.Details)
}

func (suite *UserProfileTestSuite) TestUserProfile_SetUserProfile_ValidationCompanyNameError() {
	body := `{
		"company": {
			"company_name": "Unit Test♂", 
			"website": "http://localhost", 
			"annual_income": {"from": 0, "to": 1000},
			"number_of_employees": {"from": 1, "to": 10}
		}
	}`

	_, err := suite.caller.Builder().
		Method(http.MethodPatch).
		Path(common.AuthUserGroupPath + userProfilePath).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)

	err1, ok := httpErr.Message.(*grpc.ResponseErrorMessage)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), common.ErrorMessageIncorrectCompanyName, err1)
	assert.Regexp(suite.T(), "CompanyName", err1.Details)
}

func (suite *UserProfileTestSuite) TestUserProfile_SetUserProfile_BillingServerSystemError() {
	body := `{"help": {"product_promotion_and_development": false}}`
	suite.router.dispatch.Services.Billing = mock.NewBillingServerSystemErrorMock()

	_, err := suite.caller.Builder().
		Method(http.MethodPatch).
		Path(common.AuthUserGroupPath + userProfilePath).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusInternalServerError, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorUnknown, httpErr.Message)
}

func (suite *UserProfileTestSuite) TestUserProfile_SetUserProfile_BillingServerReturnError() {
	body := `{"help": {"product_promotion_and_development": false}}`
	suite.router.dispatch.Services.Billing = mock.NewBillingServerErrorMock()

	_, err := suite.caller.Builder().
		Method(http.MethodPatch).
		Path(common.AuthUserGroupPath + userProfilePath).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), mock.SomeError, httpErr.Message)
}

func (suite *UserProfileTestSuite) TestUserProfile_ConfirmEmail_Ok() {
	body := `{"token": "123456789"}`

	_, err := suite.caller.Builder().
		Method(http.MethodPut).
		Path(common.AuthProjectGroupPath + userProfileConfirmEmailPath).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
}

func (suite *UserProfileTestSuite) TestUserProfile_ConfirmEmail_EmptyToken_Error() {
	_, err := suite.caller.Builder().
		Method(http.MethodPut).
		Path(common.AuthProjectGroupPath + userProfileConfirmEmailPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorRequestParamsIncorrect, httpErr.Message)
}

func (suite *UserProfileTestSuite) TestUserProfile_ConfirmEmail_BillingServerSystemError() {
	body := `{"token": "123456789"}`
	suite.router.dispatch.Services.Billing = mock.NewBillingServerSystemErrorMock()

	_, err := suite.caller.Builder().
		Method(http.MethodPut).
		Path(common.AuthProjectGroupPath + userProfileConfirmEmailPath).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusInternalServerError, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorUnknown, httpErr.Message)
}

func (suite *UserProfileTestSuite) TestUserProfile_ConfirmEmail_BillingServerReturnError() {
	body := `{"token": "123456789"}`
	suite.router.dispatch.Services.Billing = mock.NewBillingServerErrorMock()

	_, err := suite.caller.Builder().
		Method(http.MethodPut).
		Path(common.AuthProjectGroupPath + userProfileConfirmEmailPath).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), mock.SomeError, httpErr.Message)
}

func (suite *UserProfileTestSuite) TestUserProfile_CreatePageReview_Ok() {
	body := `{"review": "some review text", "page_id": "primary_onboarding"}`
	_, err := suite.caller.Builder().
		Method(http.MethodPost).
		Path(common.AuthUserGroupPath + userProfilePathFeedback).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
}

func (suite *UserProfileTestSuite) TestUserProfile_CreatePageReview_Unauthorized_Error() {

	reqInit := func(request *http.Request, middleware test.Middleware) {
		middleware.Pre(test.PreAuthUserMiddleware(&common.AuthUser{}))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}

	_, err := suite.caller.Builder().
		Method(http.MethodPost).
		Path(common.AuthUserGroupPath + userProfilePathFeedback).
		Init(reqInit).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusUnauthorized, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorMessageAccessDenied, httpErr.Message)
}

func (suite *UserProfileTestSuite) TestUserProfile_CreatePageReview_BindError() {

	body := `{"review": "some review text", "page_id": "merchant_onboarding"}`

	_, err := suite.caller.Builder().
		Method(http.MethodPost).
		Path(common.AuthUserGroupPath + userProfilePathFeedback).
		Init(test.ReqInitXML()).
		BodyString(body).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorRequestParamsIncorrect, httpErr.Message)
}

func (suite *UserProfileTestSuite) TestUserProfile_CreatePageReview_ValidatePageIdError() {

	body := `{"review": "some review text", "page_id": "unknown_page"}`

	_, err := suite.caller.Builder().
		Method(http.MethodPost).
		Path(common.AuthUserGroupPath + userProfilePathFeedback).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorMessageIncorrectPageId, httpErr.Message)
}

func (suite *UserProfileTestSuite) TestUserProfile_CreatePageReview_ValidateReviewError() {

	body := `{"review": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "page_id": "primary_onboarding"}`

	_, err := suite.caller.Builder().
		Method(http.MethodPost).
		Path(common.AuthUserGroupPath + userProfilePathFeedback).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorMessageIncorrectReview, httpErr.Message)
}

func (suite *UserProfileTestSuite) TestUserProfile_CreatePageReview_BillingServerSystemError() {

	suite.router.dispatch.Services.Billing = mock.NewBillingServerSystemErrorMock()
	body := `{"review": "some review text", "page_id": "primary_onboarding"}`

	_, err := suite.caller.Builder().
		Method(http.MethodPost).
		Path(common.AuthUserGroupPath + userProfilePathFeedback).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusInternalServerError, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorUnknown, httpErr.Message)
}

func (suite *UserProfileTestSuite) TestUserProfile_CreatePageReview_BillingServerResultError() {

	suite.router.dispatch.Services.Billing = mock.NewBillingServerErrorMock()
	body := `{"review": "some review text", "page_id": "primary_onboarding"}`

	_, err := suite.caller.Builder().
		Method(http.MethodPost).
		Path(common.AuthUserGroupPath + userProfilePathFeedback).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), mock.SomeError, httpErr.Message)
}
