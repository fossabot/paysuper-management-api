package handlers

import (
	"encoding/json"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo/v4"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/billing"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/grpc"
	"github.com/paysuper/paysuper-management-api/internal/dispatcher/common"
	"github.com/paysuper/paysuper-management-api/internal/mock"
	"github.com/paysuper/paysuper-management-api/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type TokenTestSuite struct {
	suite.Suite
	router *TokenRoute
	caller *test.EchoReqResCaller
}

func Test_Customer(t *testing.T) {
	suite.Run(t, new(TokenTestSuite))
}

func (suite *TokenTestSuite) SetupTest() {
	var e error
	settings := test.DefaultSettings()
	srv := common.Services{
		Billing: mock.NewBillingServerOkMock(),
	}
	suite.caller, e = test.SetUp(settings, srv, func(set *test.TestSet, mw test.Middleware) common.Handlers {
		suite.router = NewTokenRoute(set.HandlerSet, set.GlobalConfig)
		return common.Handlers{
			suite.router,
		}
	})
	if e != nil {
		panic(e)
	}
}

func (suite *TokenTestSuite) TearDownTest() {}

func (suite *TokenTestSuite) TestToken_CreateToken_Ok() {
	body := &grpc.TokenRequest{
		User: &billing.TokenUser{
			Id: bson.NewObjectId().Hex(),
			Email: &billing.TokenUserEmailValue{
				Value: "test@unit.test",
			},
			Phone: &billing.TokenUserPhoneValue{
				Value: "1234567890",
			},
			Name: &billing.TokenUserValue{
				Value: "Unit Test",
			},
			Ip: &billing.TokenUserIpValue{
				Value: "127.0.0.1",
			},
			Locale: &billing.TokenUserLocaleValue{
				Value: "ru",
			},
			Address: &billing.OrderBillingAddress{
				Country:    "RU",
				City:       "St.Petersburg",
				PostalCode: "190000",
				State:      "SPE",
			},
		},
		Settings: &billing.TokenSettings{
			ProjectId:   bson.NewObjectId().Hex(),
			Currency:    "RUB",
			Amount:      100,
			Description: "test payment",
		},
	}

	b, err := json.Marshal(body)
	assert.NoError(suite.T(), err)

	reqInit := func(request *http.Request, middleware test.Middleware) {
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		request.Header.Set(common.HeaderXApiSignatureHeader, "signature")
	}

	res, err := suite.caller.Builder().
		Method(http.MethodPost).
		Path(common.AuthProjectGroupPath + tokenPath).
		Init(reqInit).
		BodyBytes(b).
		Exec(suite.T())

	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), http.StatusOK, res.Code)
	assert.NotEmpty(suite.T(), res.Body.String())
}

func (suite *TokenTestSuite) TestToken_CreateToken_BindError() {
	body := `{"user": "qwerty", "metadata": "qwerty"}`

	_, err := suite.caller.Builder().
		Method(http.MethodPost).
		Params(":"+common.RequestParameterId, "1").
		Path(common.AuthProjectGroupPath + tokenPath).
		Init(test.ReqInitJSON()).
		BodyString(body).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorRequestParamsIncorrect, httpErr.Message)
}

func (suite *TokenTestSuite) TestToken_CreateToken_ValidationError() {
	body := &grpc.TokenRequest{
		User: &billing.TokenUser{
			Id: bson.NewObjectId().Hex(),
			Email: &billing.TokenUserEmailValue{
				Value: "test@unit.test",
			},
			Phone: &billing.TokenUserPhoneValue{
				Value: "1234567890",
			},
			Name: &billing.TokenUserValue{
				Value: "Unit Test",
			},
			Ip: &billing.TokenUserIpValue{
				Value: "127.0.0.1",
			},
			Locale: &billing.TokenUserLocaleValue{
				Value: "ru",
			},
			Address: &billing.OrderBillingAddress{
				Country:    "RU",
				City:       "St.Petersburg",
				PostalCode: "190000",
				State:      "SPE",
			},
		},
		Settings: &billing.TokenSettings{
			ProjectId:   bson.NewObjectId().Hex(),
			Currency:    "RUB",
			Amount:      -100,
			Description: "test payment",
		},
	}

	b, err := json.Marshal(body)
	assert.NoError(suite.T(), err)

	_, err = suite.caller.Builder().
		Method(http.MethodPost).
		Params(":"+common.RequestParameterId, "1").
		Path(common.AuthProjectGroupPath + tokenPath).
		Init(test.ReqInitJSON()).
		BodyBytes(b).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Regexp(suite.T(), common.NewValidationError("Amount"), httpErr.Message)
}

func (suite *TokenTestSuite) TestToken_CreateToken_CheckProjectRequestSignature_System_Error() {
	body := &grpc.TokenRequest{
		User: &billing.TokenUser{
			Id: bson.NewObjectId().Hex(),
			Email: &billing.TokenUserEmailValue{
				Value: "test@unit.test",
			},
			Phone: &billing.TokenUserPhoneValue{
				Value: "1234567890",
			},
			Name: &billing.TokenUserValue{
				Value: "Unit Test",
			},
			Ip: &billing.TokenUserIpValue{
				Value: "127.0.0.1",
			},
			Locale: &billing.TokenUserLocaleValue{
				Value: "ru",
			},
			Address: &billing.OrderBillingAddress{
				Country:    "RU",
				City:       "St.Petersburg",
				PostalCode: "190000",
				State:      "SPE",
			},
		},
		Settings: &billing.TokenSettings{
			ProjectId:   bson.NewObjectId().Hex(),
			Currency:    "RUB",
			Amount:      100,
			Description: "test payment",
		},
	}

	b, err := json.Marshal(body)
	assert.NoError(suite.T(), err)

	suite.router.dispatch.Services.Billing = mock.NewBillingServerSystemErrorMock()

	reqInit := func(request *http.Request, middleware test.Middleware) {
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		request.Header.Set(common.HeaderXApiSignatureHeader, "signature")
	}

	_, err = suite.caller.Builder().
		Method(http.MethodPost).
		Path(common.AuthProjectGroupPath + tokenPath).
		Init(reqInit).
		BodyBytes(b).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusInternalServerError, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorUnknown, httpErr.Message)
}

func (suite *TokenTestSuite) TestToken_CreateToken_CheckProjectRequestSignature_ResultError() {
	body := &grpc.TokenRequest{
		User: &billing.TokenUser{
			Id: bson.NewObjectId().Hex(),
			Email: &billing.TokenUserEmailValue{
				Value: "test@unit.test",
			},
			Phone: &billing.TokenUserPhoneValue{
				Value: "1234567890",
			},
			Name: &billing.TokenUserValue{
				Value: "Unit Test",
			},
			Ip: &billing.TokenUserIpValue{
				Value: "127.0.0.1",
			},
			Locale: &billing.TokenUserLocaleValue{
				Value: "ru",
			},
			Address: &billing.OrderBillingAddress{
				Country:    "RU",
				City:       "St.Petersburg",
				PostalCode: "190000",
				State:      "SPE",
			},
		},
		Settings: &billing.TokenSettings{
			ProjectId:   bson.NewObjectId().Hex(),
			Currency:    "RUB",
			Amount:      100,
			Description: "test payment",
		},
	}

	b, err := json.Marshal(body)
	assert.NoError(suite.T(), err)

	suite.router.dispatch.Services.Billing = mock.NewBillingServerErrorMock()

	reqInit := func(request *http.Request, middleware test.Middleware) {
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		request.Header.Set(common.HeaderXApiSignatureHeader, "signature")
	}

	_, err = suite.caller.Builder().
		Method(http.MethodPost).
		Path(common.AuthProjectGroupPath + tokenPath).
		Init(reqInit).
		BodyBytes(b).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), mock.SomeError, httpErr.Message)
}

func (suite *TokenTestSuite) TestToken_CreateToken_ChangeCustomer_System_Error() {
	body := &grpc.TokenRequest{
		User: &billing.TokenUser{
			Id: bson.NewObjectId().Hex(),
			Email: &billing.TokenUserEmailValue{
				Value: "test@unit.test",
			},
			Phone: &billing.TokenUserPhoneValue{
				Value: "1234567890",
			},
			Name: &billing.TokenUserValue{
				Value: "Unit Test",
			},
			Ip: &billing.TokenUserIpValue{
				Value: "127.0.0.1",
			},
			Locale: &billing.TokenUserLocaleValue{
				Value: "ru",
			},
			Address: &billing.OrderBillingAddress{
				Country:    "RU",
				City:       "St.Petersburg",
				PostalCode: "190000",
				State:      "SPE",
			},
		},
		Settings: &billing.TokenSettings{
			ProjectId:   bson.NewObjectId().Hex(),
			Currency:    "RUB",
			Amount:      100,
			Description: "test payment",
		},
	}

	b, err := json.Marshal(body)
	assert.NoError(suite.T(), err)

	suite.router.dispatch.Services.Billing = mock.NewBillingServerSystemErrorMock()

	reqInit := func(request *http.Request, middleware test.Middleware) {
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		request.Header.Set(common.HeaderXApiSignatureHeader, "signature")
	}

	_, err = suite.caller.Builder().
		Method(http.MethodPost).
		Path(common.AuthProjectGroupPath + tokenPath).
		Init(reqInit).
		BodyBytes(b).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusInternalServerError, httpErr.Code)
	assert.Equal(suite.T(), common.ErrorUnknown, httpErr.Message)
}

func (suite *TokenTestSuite) TestToken_CreateToken_ChangeCustomer_ResultError() {
	body := &grpc.TokenRequest{
		User: &billing.TokenUser{
			Id: bson.NewObjectId().Hex(),
			Email: &billing.TokenUserEmailValue{
				Value: "test@unit.test",
			},
			Phone: &billing.TokenUserPhoneValue{
				Value: "1234567890",
			},
			Name: &billing.TokenUserValue{
				Value: "Unit Test",
			},
			Ip: &billing.TokenUserIpValue{
				Value: "127.0.0.1",
			},
			Locale: &billing.TokenUserLocaleValue{
				Value: "ru",
			},
			Address: &billing.OrderBillingAddress{
				Country:    "RU",
				City:       "St.Petersburg",
				PostalCode: "190000",
				State:      "SPE",
			},
		},
		Settings: &billing.TokenSettings{
			ProjectId:   bson.NewObjectId().Hex(),
			Currency:    "RUB",
			Amount:      100,
			Description: "test payment",
		},
	}

	b, err := json.Marshal(body)
	assert.NoError(suite.T(), err)

	suite.router.dispatch.Services.Billing = mock.NewBillingServerErrorMock()

	reqInit := func(request *http.Request, middleware test.Middleware) {
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		request.Header.Set(common.HeaderXApiSignatureHeader, "signature")
	}

	_, err = suite.caller.Builder().
		Method(http.MethodPost).
		Path(common.AuthProjectGroupPath + tokenPath).
		Init(reqInit).
		BodyBytes(b).
		Exec(suite.T())

	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), mock.SomeError, httpErr.Message)
}
