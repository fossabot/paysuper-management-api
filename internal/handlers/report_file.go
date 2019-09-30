package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/ProtocolONE/go-core/logger"
	"github.com/ProtocolONE/go-core/provider"
	"github.com/labstack/echo/v4"
	awsWrapper "github.com/paysuper/paysuper-aws-manager"
	"github.com/paysuper/paysuper-management-api/internal/dispatcher/common"
	reporterPkg "github.com/paysuper/paysuper-reporter/pkg"
	reporterProto "github.com/paysuper/paysuper-reporter/pkg/proto"
	"net/http"
	"os"
	"strings"
)

const (
	reportFilePath         = "/report_file"
	reportFileDownloadPath = "/report_file/download/:file"
)

type reportFileRequest struct {
	Id         string                 `json:"id" form:"id" bson:"_id"`
	MerchantId string                 `json:"merchant_id" form:"merchant_id" bson:"merchant_id"`
	FileType   string                 `json:"file_type" form:"file_type" bson:"file_type"`
	ReportType string                 `json:"report_type" form:"report_type" bson:"report_type"`
	Template   string                 `json:"template" form:"template" bson:"template"`
	Params     map[string]interface{} `json:"params" form:"params" bson:"params"`
}

type ReportFileRoute struct {
	dispatch   common.HandlerSet
	awsManager awsWrapper.AwsManagerInterface
	cfg        common.Config
	provider.LMT
}

func NewReportFileRoute(set common.HandlerSet, awsManager awsWrapper.AwsManagerInterface, cfg *common.Config) *ReportFileRoute {
	set.AwareSet.Logger = set.AwareSet.Logger.WithFields(logger.Fields{"router": "ReportFileRoute"})
	return &ReportFileRoute{
		dispatch:   set,
		LMT:        &set.AwareSet,
		cfg:        *cfg,
		awsManager: awsManager,
	}
}

func (h *ReportFileRoute) Route(groups *common.Groups) {
	groups.Access.POST(reportFilePath, h.create)
	groups.Access.GET(reportFileDownloadPath, h.download)
}

// Send a request to create a report for download.
// POST /api/v1/s/report_file
//
// @Example curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" \
//      -H "Authorization: Bearer %access_token_here%" \
//      -d '{"file_type": "pdf", "period_from": 1566727410, "period_to": "1566736763"}' \
//      https://api.paysuper.online/api/v1/s/report_file
//
func (h *ReportFileRoute) create(ctx echo.Context) error {
	authUser := common.ExtractUserContext(ctx)
	data := &reportFileRequest{}
	if err := ctx.Bind(data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestDataInvalid)
	}

	params, err := json.Marshal(data.Params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorRequestDataInvalid)
	}

	err = h.dispatch.Validate.Struct(data)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	req := &reporterProto.ReportFile{
		UserId:           authUser.Id,
		MerchantId:       data.MerchantId,
		ReportType:       data.ReportType,
		FileType:         data.FileType,
		Template:         data.Template,
		Params:           params,
		SendNotification: true,
	}

	res, err := h.dispatch.Services.Reporter.CreateFile(ctx.Request().Context(), req)
	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, reporterPkg.ServiceName, "CreateFile", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorMessageCreateReportFile)
	}

	return ctx.JSON(http.StatusOK, res)
}

// Send a request to create a report for download.
// GET /api/v1/s/report_file/download/5ced34d689fce60bf4440829.csv
//
// @Example curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" \
//      -H "Authorization: Bearer %access_token_here%" \
//      https://api.paysuper.online/api/v1/s/report_file/download/5ced34d689fce60bf4440829.csv
//
func (h *ReportFileRoute) download(ctx echo.Context) error {
	authUser := common.ExtractUserContext(ctx)
	file := ctx.Param("file")
	if file == "" {
		h.L().Error("unable to find the file")
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestParamsIncorrect)
	}

	params := strings.Split(file, ".")

	if len(params) != 2 {
		h.L().Error("incorrect of file string")
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestParamsIncorrect)
	}

	fileName := fmt.Sprintf(reporterPkg.FileMask, authUser.Id, params[0], params[1])
	filePath := os.TempDir() + string(os.PathSeparator) + fileName
	_, err := h.awsManager.Download(ctx.Request().Context(), filePath, &awsWrapper.DownloadInput{FileName: fileName})

	if err != nil {
		h.L().Error("unable to find the file id")
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorMessageDownloadReportFile)
	}

	defer os.Remove(filePath)
	return ctx.File(filePath)
}
