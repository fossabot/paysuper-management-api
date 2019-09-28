package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/micro/go-micro/client"
	"github.com/paysuper/paysuper-management-api/internal/dispatcher/common"
	"github.com/paysuper/paysuper-management-api/internal/mock"
	"github.com/paysuper/paysuper-management-api/internal/test"
	"github.com/paysuper/paysuper-tax-service/proto"
	"github.com/stretchr/testify/assert"
	mock2 "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/url"
	"strconv"
	"testing"
)

type TaxServiceMock struct {
	mock2.Mock
}

type TaxesTestSuite struct {
	suite.Suite
	router *TaxesRoute
	caller *test.EchoReqResCaller
}

func Test_Taxes(t *testing.T) {
	suite.Run(t, new(TaxesTestSuite))
}

func (suite *TaxesTestSuite) SetupTest() {
	var e error
	settings := test.DefaultSettings()
	srv := common.Services{
		Billing: mock.NewBillingServerOkMock(),
		Tax:     createNewTaxServiceMock(),
	}
	suite.caller, e = test.SetUp(settings, srv, func(set *test.TestSet, mw test.Middleware) common.Handlers {
		suite.router = NewTaxesRoute(set.HandlerSet, set.GlobalConfig)
		return common.Handlers{
			suite.router,
		}
	})
	if e != nil {
		panic(e)
	}
}

func (suite *TaxesTestSuite) Test_GetRates() {
	rates := getRates(suite.T(), suite, "RU", "", "", "", 0, 0)
	assert.Len(suite.T(), rates, int(suite.router.cfg.LimitDefault))
	assert.Equal(suite.T(), "RU", rates[0].Country)

	rates = getRates(suite.T(), suite, "", "City", "", "", 0, 0)
	assert.Len(suite.T(), rates, int(suite.router.cfg.LimitDefault))
	assert.Equal(suite.T(), "City", rates[0].City)

	rates = getRates(suite.T(), suite, "", "", "00001", "", 0, 0)
	assert.Len(suite.T(), rates, int(suite.router.cfg.LimitDefault))
	assert.Equal(suite.T(), "00001", rates[0].Zip)

	rates = getRates(suite.T(), suite, "", "", "", "NY", 0, 0)
	assert.Len(suite.T(), rates, int(suite.router.cfg.LimitDefault))
	assert.Equal(suite.T(), "NY", rates[0].State)

	rates = getRates(suite.T(), suite, "", "", "", "NY", 1, 0)
	assert.Len(suite.T(), rates, 1)
	assert.Equal(suite.T(), "NY", rates[0].State)

	rates = getRates(suite.T(), suite, "", "", "", "NY", 1, 2)
	assert.Len(suite.T(), rates, 1)
	assert.Equal(suite.T(), "NY", rates[0].State)
	assert.EqualValues(suite.T(), 2, rates[0].Id)
}

func getRates(t *testing.T, suite *TaxesTestSuite, country, city, zip, state string, limit, offset int) []*tax_service.TaxRate {
	t.Helper()

	q := make(url.Values)
	q.Set("city", city)
	q.Set("country", country)
	q.Set("zip", zip)
	q.Set("state", state)
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	if offset > 0 {
		q.Set("offset", strconv.Itoa(offset))
	}

	res, err := suite.caller.Builder().
		Method(http.MethodGet).
		Path(common.AuthUserGroupPath + taxesPath).
		SetQueryParams(q).
		Exec(suite.T())

	assert.NoError(t, err)

	var response []*tax_service.TaxRate
	err = json.Unmarshal(res.Body.Bytes(), &response)

	return response
}

func (suite *TaxesTestSuite) Test_GetRatesError() {
	q := make(url.Values)
	q.Set("city", "fail")

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Path(common.AuthUserGroupPath + taxesPath).
		SetQueryParams(q).
		Exec(suite.T())

	if assert.Error(suite.T(), err) {
		assert.Equal(suite.T(), http.StatusInternalServerError, err.(*echo.HTTPError).Code)
	}
}

func (suite *TaxesTestSuite) Test_CreateTax() {
	obj := &tax_service.TaxRate{
		Zip:  "00001",
		City: "City",
	}

	b, _ := json.Marshal(obj)
	res, err := suite.caller.Builder().
		Method(http.MethodPost).
		Path(common.AuthUserGroupPath + taxesPath).
		Init(test.ReqInitJSON()).
		BodyBytes(b).
		Exec(suite.T())

	if assert.NoError(suite.T(), err) {
		obj1 := &tax_service.TaxRate{}
		err := json.Unmarshal(res.Body.Bytes(), obj1)
		assert.NoError(suite.T(), err)

		assert.Equal(suite.T(), http.StatusOK, res.Code)
		assert.Equal(suite.T(), obj.Zip, obj1.Zip)
		assert.Equal(suite.T(), obj.City, obj1.City)
	}
}

func (suite *TaxesTestSuite) Test_CreateTaxWithError() {
	testCreateTaxWithError(suite.T(), suite, nil, http.StatusBadRequest)
	testCreateTaxWithError(suite.T(), suite, &tax_service.TaxRate{Id: 1}, http.StatusInternalServerError)
}

func testCreateTaxWithError(t *testing.T, suite *TaxesTestSuite, obj *tax_service.TaxRate, code int) {
	t.Helper()

	var body []byte
	if obj != nil {
		body, _ = json.Marshal(obj)
	}

	_, err := suite.caller.Builder().
		Method(http.MethodPost).
		Params(":"+common.RequestParameterId, "1").
		Path(common.AuthUserGroupPath + taxesPath).
		Init(test.ReqInitJSON()).
		BodyBytes(body).
		Exec(suite.T())

	if assert.Error(t, err) {
		assert.Equal(t, code, err.(*echo.HTTPError).Code)
	}
}

func (suite *TaxesTestSuite) Test_DeleteTax() {
	res, err := suite.caller.Builder().
		Method(http.MethodDelete).
		Params(":"+common.RequestParameterId, "1").
		Path(common.AuthUserGroupPath + taxesIDPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), http.StatusOK, res.Code)
	}
}

func (suite *TaxesTestSuite) Test_DeleteTaxWithInvalidId() {
	testDeleteTaxWithError(suite.T(), suite, " ", http.StatusBadRequest)
	testDeleteTaxWithError(suite.T(), suite, "string", http.StatusBadRequest)
	testDeleteTaxWithError(suite.T(), suite, "0.1", http.StatusBadRequest)
	testDeleteTaxWithError(suite.T(), suite, "0", http.StatusInternalServerError)
}

func testDeleteTaxWithError(t *testing.T, suite *TaxesTestSuite, id string, code int) {
	t.Helper()
	_, err := suite.caller.Builder().
		Method(http.MethodDelete).
		Params(":"+common.RequestParameterId, id).
		Path(common.AuthUserGroupPath + taxesIDPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	if assert.Error(t, err) {
		assert.Equal(t, code, err.(*echo.HTTPError).Code)
	}
}

func createNewTaxServiceMock() tax_service.TaxService {
	return &TaxServiceMock{}
}

func (ts *TaxServiceMock) GetRate(ctx context.Context, in *tax_service.GetRateRequest, opts ...client.CallOption) (*tax_service.GetRateResponse, error) {
	panic("this method is not implemented in mock")
}

func (ts *TaxServiceMock) GetRates(ctx context.Context, in *tax_service.GetRatesRequest, opts ...client.CallOption) (*tax_service.GetRatesResponse, error) {
	if in.City == "fail" {
		return nil, errors.New("Invalid request")
	}
	res := &tax_service.GetRatesResponse{}
	for i := 0; i < int(in.Limit); i++ {
		res.Rates = append(
			res.Rates,
			&tax_service.TaxRate{
				Id:      uint32(i) + uint32(in.Offset),
				City:    in.City,
				Zip:     in.Zip,
				Country: in.Country,
				Rate:    0.1,
				State:   in.State,
			},
		)
	}

	return res, nil
}

func (ts *TaxServiceMock) CreateOrUpdate(ctx context.Context, in *tax_service.TaxRate, opts ...client.CallOption) (*tax_service.TaxRate, error) {
	if in.Id == 1 {
		return nil, errors.New("Invalid request")
	}

	return in, nil
}

func (ts *TaxServiceMock) DeleteRateById(ctx context.Context, in *tax_service.DeleteRateRequest, opts ...client.CallOption) (*tax_service.DeleteRateResponse, error) {
	if in.Id == 1 {
		return &tax_service.DeleteRateResponse{}, nil
	}

	return nil, errors.New("Invalid request")
}
