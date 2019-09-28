package handlers

import (
	"github.com/globalsign/mgo/bson"
	"github.com/golang/protobuf/ptypes"
	billingMocks "github.com/paysuper/paysuper-billing-server/pkg/mocks"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/billing"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/grpc"
	"github.com/paysuper/paysuper-management-api/internal/dispatcher/common"
	"github.com/paysuper/paysuper-management-api/internal/test"
	"github.com/stretchr/testify/assert"
	mock2 "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/url"
	"testing"
)

var payoutMock = &billing.PayoutDocument{
	Id:                 bson.NewObjectId().Hex(),
	SourceId:           []string{bson.NewObjectId().Hex(), bson.NewObjectId().Hex()},
	Transaction:        "",
	Amount:             100500,
	Currency:           "EUR",
	Status:             "pending",
	Description:        "royalty for june-july 2019",
	Destination:        &billing.MerchantBanking{},
	CreatedAt:          ptypes.TimestampNow(),
	UpdatedAt:          ptypes.TimestampNow(),
	ArrivalDate:        ptypes.TimestampNow(),
	FailureCode:        "",
	FailureMessage:     "",
	FailureTransaction: "",
	MerchantId:         bson.NewObjectId().Hex(),
	SignatureData: &billing.PayoutDocumentSignatureData{
		DetailsUrl:          "http://localhost",
		FilesUrl:            "http://localhost",
		SignatureRequestId:  bson.NewObjectId().Hex(),
		MerchantSignatureId: bson.NewObjectId().Hex(),
		PsSignatureId:       bson.NewObjectId().Hex(),
	},
	HasMerchantSignature:  false,
	HasPspSignature:       false,
	SignedDocumentFileUrl: "",
}

type PayoutDocumentsTestSuite struct {
	suite.Suite
	router *PayoutDocumentsRoute
	caller *test.EchoReqResCaller
}

func Test_PayoutDocuments(t *testing.T) {
	suite.Run(t, new(PayoutDocumentsTestSuite))
}

func (suite *PayoutDocumentsTestSuite) SetupTest() {

	billingService := &billingMocks.BillingService{}

	billingService.On("GetPayoutDocuments", mock2.Anything, mock2.Anything).
		Return(&grpc.GetPayoutDocumentsResponse{
			Status: http.StatusOK,
			Data: &grpc.PayoutDocumentsPaginate{
				Count: 1,
				Items: []*billing.PayoutDocument{payoutMock},
			},
		}, nil)

	billingService.On("CreatePayoutDocument", mock2.Anything, mock2.Anything).
		Return(&grpc.PayoutDocumentResponse{
			Status: http.StatusOK,
			Item:   payoutMock,
		}, nil)

	billingService.On("UpdatePayoutDocument", mock2.Anything, mock2.Anything).
		Return(&grpc.PayoutDocumentResponse{
			Status: http.StatusOK,
			Item:   payoutMock,
		}, nil)

	var e error
	settings := test.DefaultSettings()
	srv := common.Services{
		Billing: billingService,
	}
	suite.caller, e = test.SetUp(settings, srv, func(set *test.TestSet, mw test.Middleware) common.Handlers {
		suite.router = NewPayoutDocumentsRoute(set.HandlerSet, set.GlobalConfig)
		return common.Handlers{
			suite.router,
		}
	})
	if e != nil {
		panic(e)
	}
}

func (suite *PayoutDocumentsTestSuite) TearDownTest() {}

func (suite *PayoutDocumentsTestSuite) TestPayoutDocuments_Ok_getPayoutDocumentsList() {
	q := make(url.Values)
	q.Set("status", "pending")
	q.Set("merchant_id", "5bdc39a95d1e1100019fb7df")
	q.Set("limit", "10")
	q.Set("signed", "true")
	q.Set("offset", "10")

	res, err := suite.caller.Builder().
		Method(http.MethodGet).
		SetQueryParams(q).
		Path(common.AuthUserGroupPath + payoutsPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), http.StatusOK, res.Code)
		assert.NotEmpty(suite.T(), res.Body.String())
	}
}

func (suite *BalanceTestSuite) TestPayoutDocuments_Fail_getPayoutDocumentsList_MerchantNotFound() {
	assert.Equal(suite.T(), common.TestStubImplementMe, "implement me!")
}

func (suite *BalanceTestSuite) TestPayoutDocuments_Fail_getPayoutDocumentsList_MerchantIdInvalid() {
	assert.Equal(suite.T(), common.TestStubImplementMe, "implement me!")
}

func (suite *BalanceTestSuite) TestPayoutDocuments_Fail_getPayoutDocumentsList_StatusInvalid() {
	assert.Equal(suite.T(), common.TestStubImplementMe, "implement me!")
}

func (suite *BalanceTestSuite) TestPayoutDocuments_Fail_getPayoutDocumentsList_LimitInvalid() {
	assert.Equal(suite.T(), common.TestStubImplementMe, "implement me!")
}

func (suite *BalanceTestSuite) TestPayoutDocuments_Fail_getPayoutDocumentsList_OffsetInvalid() {
	assert.Equal(suite.T(), common.TestStubImplementMe, "implement me!")
}

func (suite *BalanceTestSuite) TestPayoutDocuments_Fail_getPayoutDocumentsList_GrpcError() {
	assert.Equal(suite.T(), common.TestStubImplementMe, "implement me!")
}

func (suite *BalanceTestSuite) TestPayoutDocuments_Fail_getPayoutDocumentsList_ResponseStatusNotOk() {
	assert.Equal(suite.T(), common.TestStubImplementMe, "implement me!")
}

func (suite *PayoutDocumentsTestSuite) TestPayoutDocuments_Ok_getPayoutDocument() {

	res, err := suite.caller.Builder().
		Method(http.MethodGet).
		Params(":"+common.RequestParameterId, bson.NewObjectId().Hex()).
		Path(common.AuthUserGroupPath + payoutsPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), http.StatusOK, res.Code)
		assert.NotEmpty(suite.T(), res.Body.String())
	}
}

func (suite *BalanceTestSuite) TestPayoutDocuments_Fail_getPayoutDocument_NotFound() {
	assert.Equal(suite.T(), common.TestStubImplementMe, "implement me!")
}

func (suite *BalanceTestSuite) TestPayoutDocuments_Fail_getPayoutDocument_InvalidId() {
	assert.Equal(suite.T(), common.TestStubImplementMe, "implement me!")
}

func (suite *BalanceTestSuite) TestPayoutDocuments_Fail_getPayoutDocument_GrpcError() {
	assert.Equal(suite.T(), common.TestStubImplementMe, "implement me!")
}

func (suite *BalanceTestSuite) TestPayoutDocuments_Fail_getPayoutDocument_ResponseStatusNotOk() {
	assert.Equal(suite.T(), common.TestStubImplementMe, "implement me!")
}

func (suite *PayoutDocumentsTestSuite) TestPayoutDocuments_Ok_createPayoutDocument() {
	bodyJson := `{"description": "royalty for june-july 2019", "merchant_id": "5bdc39a95d1e1100019fb7df"}`

	res, err := suite.caller.Builder().
		Method(http.MethodPost).
		Path(common.AuthUserGroupPath + payoutsPath).
		Init(test.ReqInitJSON()).
		BodyString(bodyJson).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, res.Code)
	assert.NotEmpty(suite.T(), res.Body.String())
}

func (suite *BalanceTestSuite) TestPayoutDocuments_Fail_createPayoutDocument_MerchantNotFound() {
	assert.Equal(suite.T(), common.TestStubImplementMe, "implement me!")
}

func (suite *BalanceTestSuite) TestPayoutDocuments_Fail_createPayoutDocument_MerchantIdInvalid() {
	assert.Equal(suite.T(), common.TestStubImplementMe, "implement me!")
}

func (suite *BalanceTestSuite) TestPayoutDocuments_Fail_createPayoutDocument_ValidationError() {
	assert.Equal(suite.T(), common.TestStubImplementMe, "implement me!")
}

func (suite *BalanceTestSuite) TestPayoutDocuments_Fail_createPayoutDocument_GrpcError() {
	assert.Equal(suite.T(), common.TestStubImplementMe, "implement me!")
}

func (suite *BalanceTestSuite) TestPayoutDocuments_Fail_createPayoutDocument_ResponseStatusNotOk() {
	assert.Equal(suite.T(), common.TestStubImplementMe, "implement me!")
}

func (suite *PayoutDocumentsTestSuite) TestPayoutDocuments_Ok_updatePayoutDocument() {
	bodyJson := `{"status": "failed", "failure_code": "account_closed"}`

	res, err := suite.caller.Builder().
		Method(http.MethodPost).
		Params(":"+common.RequestParameterId, bson.NewObjectId().Hex()).
		Path(common.AuthUserGroupPath + payoutsIdPath).
		Init(test.ReqInitJSON()).
		BodyString(bodyJson).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, res.Code)
	assert.NotEmpty(suite.T(), res.Body.String())
}

func (suite *PayoutDocumentsTestSuite) TestPayoutDocuments_Ok_updatePayoutDocument_NotModified() {
	assert.Equal(suite.T(), common.TestStubImplementMe, "implement me!")
}

func (suite *BalanceTestSuite) TestPayoutDocuments_Fail_updatePayoutDocument_InvalidId() {
	assert.Equal(suite.T(), common.TestStubImplementMe, "implement me!")
}

func (suite *BalanceTestSuite) TestPayoutDocuments_Fail_updatePayoutDocument_GrpcError() {
	assert.Equal(suite.T(), common.TestStubImplementMe, "implement me!")
}

func (suite *BalanceTestSuite) TestPayoutDocuments_Fail_updatePayoutDocument_ResponseStatusNotOk() {
	assert.Equal(suite.T(), common.TestStubImplementMe, "implement me!")
}
