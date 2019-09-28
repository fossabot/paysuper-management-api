package handlers

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo/v4"
	awsWrapper "github.com/paysuper/paysuper-aws-manager"
	awsWrapperMocks "github.com/paysuper/paysuper-aws-manager/pkg/mocks"
	"github.com/paysuper/paysuper-management-api/internal/dispatcher/common"
	"github.com/paysuper/paysuper-management-api/internal/mock"
	"github.com/paysuper/paysuper-management-api/internal/test"
	reporterMocks "github.com/paysuper/paysuper-reporter/pkg/mocks"
	reporterProto "github.com/paysuper/paysuper-reporter/pkg/proto"
	"github.com/stretchr/testify/assert"
	mock2 "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"os"
	"testing"
)

type ReportFileTestSuite struct {
	suite.Suite
	router *ReportFileRoute
	caller *test.EchoReqResCaller
}

func Test_ReportFile(t *testing.T) {
	suite.Run(t, new(ReportFileTestSuite))
}

func (suite *ReportFileTestSuite) SetupTest() {
	var e error
	settings := test.DefaultSettings()
	srv := common.Services{
		Billing: mock.NewBillingServerOkMock(),
	}
	suite.caller, e = test.SetUp(settings, srv, func(set *test.TestSet, mw test.Middleware) common.Handlers {

		downloadMockResultFn := func(
			ctx context.Context,
			filePath string,
			in *awsWrapper.DownloadInput,
			opts ...func(*s3manager.Downloader),
		) int64 {
			_, err := os.Stat(filePath)

			if err == nil {
				return 0
			}

			if !os.IsNotExist(err) {
				return 0
			}

			src, err := os.Open(set.Initial.WorkDir + "/test/test_pdf.pdf")
			if err != nil {
				return 0
			}
			defer src.Close()

			dst, err := os.Create(filePath)
			if err != nil {
				return 0
			}
			defer dst.Close()

			nBytes, err := io.Copy(dst, src)

			if err != nil {
				return 0
			}

			return nBytes
		}

		awsManagerMock := &awsWrapperMocks.AwsManagerInterface{}
		awsManagerMock.On("Upload", mock2.Anything, mock2.Anything, mock2.Anything).Return(&s3manager.UploadOutput{}, nil)
		awsManagerMock.On("Download", mock2.Anything, mock2.Anything, mock2.Anything, mock2.Anything).
			Return(downloadMockResultFn, nil)

		suite.router = NewReportFileRoute(set.HandlerSet, awsManagerMock, set.GlobalConfig)
		return common.Handlers{
			suite.router,
		}
	})
	if e != nil {
		panic(e)
	}
}

func (suite *ReportFileTestSuite) TestReportFile_create_Error_CreateFile() {
	data := `{"period_from": 1, "period_to": 2}`

	reporterService := &reporterMocks.ReporterService{}
	reporterService.
		On("CreateFile", mock2.Anything, mock2.Anything).
		Return(nil, errors.New("error"))
	suite.router.dispatch.Services.Reporter = reporterService

	_, err := suite.caller.Builder().
		Method(http.MethodPost).
		Path(common.AccessGroupPath + reportFilePath).
		Init(test.ReqInitJSON()).
		BodyString(data).
		Exec(suite.T())

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusInternalServerError, httpErr.Code)
	assert.Regexp(suite.T(), common.ErrorMessageCreateReportFile.Message, httpErr.Message)
}

func (suite *ReportFileTestSuite) TestReportFile_create_Ok() {
	data := `{"period_from": 1, "period_to": 2}`

	reporterService := &reporterMocks.ReporterService{}
	reporterService.
		On("CreateFile", mock2.Anything, mock2.Anything).
		Return(&reporterProto.CreateFileResponse{FileId: bson.NewObjectId().Hex()}, nil)
	suite.router.dispatch.Services.Reporter = reporterService

	_, err := suite.caller.Builder().
		Method(http.MethodPost).
		Path(common.AccessGroupPath + reportFilePath).
		Init(test.ReqInitJSON()).
		BodyString(data).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
}

func (suite *ReportFileTestSuite) TestReportFile_download_Error_EmptyId() {

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Path(common.AccessGroupPath + reportFileDownloadPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Regexp(suite.T(), common.ErrorRequestParamsIncorrect.Message, httpErr.Message)
}

func (suite *ReportFileTestSuite) TestReportFile_download_Error_ValidationFileEmpty() {

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Params(":"+common.RequestParameterFile, " ").
		Path(common.AccessGroupPath + reportFileDownloadPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Regexp(suite.T(), common.ErrorRequestParamsIncorrect, httpErr.Message)
}

func (suite *ReportFileTestSuite) TestReportFile_download_Error_ValidationFileIncorrect() {

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Params(":"+common.RequestParameterFile, "test").
		Path(common.AccessGroupPath + reportFileDownloadPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Regexp(suite.T(), common.ErrorRequestParamsIncorrect, httpErr.Message)
}

func (suite *ReportFileTestSuite) TestReportFile_download_Ok() {

	_, err := suite.caller.Builder().
		Method(http.MethodGet).
		Params(":"+common.RequestParameterFile, "string.csv").
		Path(common.AccessGroupPath + reportFileDownloadPath).
		Init(test.ReqInitJSON()).
		Exec(suite.T())

	assert.NoError(suite.T(), err)
}
