package handlers

import (
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo/v4"
	"github.com/paysuper/paysuper-management-api/internal/dispatcher/common"
	"github.com/paysuper/paysuper-management-api/internal/mock"
	"github.com/paysuper/paysuper-management-api/internal/test"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type MerchantUsersTestSuite struct {
	suite.Suite
	router *MerchantUsersRoute
	caller *test.EchoReqResCaller
}

func Test_MerchantUsers(t *testing.T) {
	suite.Run(t, new(MerchantUsersTestSuite))
}

func (suite *MerchantUsersTestSuite) SetupTest() {
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
		suite.router = NewMerchantUsersRoute(set.HandlerSet, set.GlobalConfig)
		return common.Handlers{
			suite.router,
		}
	})
	if e != nil {
		panic(e)
	}
}

func (suite *MerchantUsersTestSuite) TearDownTest() {}

func (suite *MerchantUsersTestSuite) TestMerchantUsers_GetList_ValidationError() {
	shouldBe := require.New(suite.T())

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Params(":"+common.RequestParameterMerchantId, "").
		Path(common.AuthUserGroupPath + merchantUsers).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	shouldBe.Error(err)
	hErr, ok := err.(*echo.HTTPError)
	shouldBe.True(ok)
	shouldBe.Equal(400, hErr.Code)
	shouldBe.NotEmpty(hErr.Message)
}

func (suite *MerchantUsersTestSuite) TestMerchantUsers_GetList_Ok() {
	shouldBe := require.New(suite.T())

	res, err := suite.caller.Builder().
		Method(http.MethodGet).
		Params(":"+common.RequestParameterMerchantId, bson.NewObjectId().Hex()).
		Path(common.AuthUserGroupPath + merchantUsers).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	shouldBe.NoError(err)
	shouldBe.Equal(http.StatusOK, res.Code)
	shouldBe.NotEmpty(res.Body.String())
}
