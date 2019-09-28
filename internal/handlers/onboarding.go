package handlers

import (
	"bytes"
	"fmt"
	"github.com/Nerufa/go-shared/config"
	"github.com/Nerufa/go-shared/logger"
	"github.com/Nerufa/go-shared/provider"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo/v4"
	awsWrapper "github.com/paysuper/paysuper-aws-manager"
	"github.com/paysuper/paysuper-billing-server/pkg"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/billing"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/grpc"
	"github.com/paysuper/paysuper-management-api/internal/dispatcher/common"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

const (
	merchantsPath                      = "/merchants"
	merchantsIdPath                    = "/merchants/:id"
	merchantsUserPath                  = "/merchants/user"
	merchantsCompanyPath               = "/merchants/company"
	merchantsContactsPath              = "/merchants/contacts"
	merchantsBankingPath               = "/merchants/banking"
	merchantsIdCompanyPath             = "/merchants/:id/company"
	merchantsIdContactsPath            = "/merchants/:id/contacts"
	merchantsIdBankingPath             = "/merchants/:id/banking"
	merchantsIdStatusCompanyPath       = "/merchants/:id/status"
	merchantsIdChangeStatusCompanyPath = "/merchants/:id/change-status"
	merchantsNotificationsPath         = "/merchants/:merchant_id/notifications"
	merchantsIdAgreementPath           = "/merchants/:id/agreement"
	merchantsAgreementDocumentPath     = "/merchants/:id/agreement/document"
	merchantsAgreementSignaturePath    = "/merchants/:id/agreement/signature"
	merchantsNotificationsIdPath       = "/merchants/:merchant_id/notifications/:notification_id"
	merchantsNotificationsMarkReadPath = "/merchants/:merchant_id/notifications/:notification_id/mark-as-read"
	merchantsTariffsPath               = "/merchants/tariffs"
	merchantsIdTariffsPath             = "/merchants/:id/tariffs"
)

const (
	agreementFileMask      = "agreement_%s.pdf"
	agreementContentType   = "application/pdf"
	agreementExtension     = "pdf"
	agreementUrlMask       = "%s://%s/admin/api/v1/merchants/%s/agreement/document"
	agreementUploadMaxSize = 3145728
)

type OnboardingFileMetadata struct {
	Name        string `json:"name"`
	Extension   string `json:"extension"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
}

type OnboardingFileData struct {
	Url      string                  `json:"url"`
	Metadata *OnboardingFileMetadata `json:"metadata"`
}

type OnboardingRoute struct {
	dispatch   common.HandlerSet
	awsManager awsWrapper.AwsManagerInterface
	cfg        common.Config
	provider.LMT
}

func NewOnboardingRoute(set common.HandlerSet, initial config.Initial, awsManager awsWrapper.AwsManagerInterface, globalCfg *common.Config) *OnboardingRoute {
	set.AwareSet.Logger = set.AwareSet.Logger.WithFields(logger.Fields{"router": "OnboardingRoute"})
	return &OnboardingRoute{
		dispatch:   set,
		LMT:        &set.AwareSet,
		cfg:        *globalCfg,
		awsManager: awsManager,
	}
}

func (h *OnboardingRoute) Route(groups *common.Groups) {
	groups.AuthUser.GET(merchantsPath, h.listMerchants)
	groups.AuthUser.GET(merchantsIdPath, h.getMerchant)
	groups.AuthUser.GET(merchantsUserPath, h.getMerchantByUser)

	groups.AuthUser.PUT(merchantsCompanyPath, h.setMerchantCompany)
	groups.AuthUser.PUT(merchantsContactsPath, h.setMerchantContacts)
	groups.AuthUser.PUT(merchantsBankingPath, h.setMerchantBanking)
	groups.AuthUser.PUT(merchantsIdCompanyPath, h.setMerchantCompany)
	groups.AuthUser.PUT(merchantsIdContactsPath, h.setMerchantContacts)
	groups.AuthUser.PUT(merchantsIdBankingPath, h.setMerchantBanking)
	groups.AuthUser.GET(merchantsIdStatusCompanyPath, h.getMerchantStatus)

	groups.AuthUser.PUT(merchantsIdChangeStatusCompanyPath, h.changeMerchantStatus)
	groups.AuthUser.PATCH(merchantsIdPath, h.changeAgreement)

	groups.AuthUser.GET(merchantsIdAgreementPath, h.generateAgreement)
	groups.AuthUser.GET(merchantsAgreementDocumentPath, h.getAgreementDocument)
	groups.AuthUser.POST(merchantsAgreementDocumentPath, h.uploadAgreementDocument)
	groups.AuthUser.PUT(merchantsAgreementSignaturePath, h.createAgreementSignature)

	groups.AuthUser.POST(merchantsNotificationsPath, h.createNotification)
	groups.AuthUser.GET(merchantsNotificationsIdPath, h.getNotification)
	groups.AuthUser.GET(merchantsNotificationsPath, h.listNotifications)
	groups.AuthUser.PUT(merchantsNotificationsMarkReadPath, h.markAsReadNotification)

	groups.AuthUser.GET(merchantsTariffsPath, h.getTariffRates)
	groups.AuthUser.POST(merchantsIdTariffsPath, h.setTariffRates)
}

func (h *OnboardingRoute) getMerchant(ctx echo.Context) error {
	id := ctx.Param(common.RequestParameterId)

	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorIdIsEmpty)
	}

	req := &grpc.GetMerchantByRequest{MerchantId: id}
	res, err := h.dispatch.Services.Billing.GetMerchantBy(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "GetMerchantBy", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}

	if res.Status != pkg.ResponseStatusOk {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.JSON(http.StatusOK, res.Item)
}

func (h *OnboardingRoute) getMerchantByUser(ctx echo.Context) error {
	authUser := common.ExtractUserContext(ctx)
	if authUser.Id == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, common.ErrorMessageAccessDenied)
	}

	req := &grpc.GetMerchantByRequest{UserId: authUser.Id}
	res, err := h.dispatch.Services.Billing.GetMerchantBy(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "GetMerchantBy", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}

	if res.Status != pkg.ResponseStatusOk {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.JSON(http.StatusOK, res.Item)
}

// @Description List merchants by conditions
// @Example curl -X GET 'Authorization: Bearer %access_token_here%' \
//  'https://api.paysuper.online/admin/api/v1/merchants?received_date_from=1568332800'
func (h *OnboardingRoute) listMerchants(ctx echo.Context) error {
	req := &grpc.MerchantListingRequest{}
	err := (&common.OnboardingMerchantListingBinder{
		LimitDefault:  h.cfg.LimitDefault,
		OffsetDefault: h.cfg.OffsetDefault,
	}).Bind(req, ctx)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestParamsIncorrect)
	}

	err = h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.ListMerchants(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "ListMerchants", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}

	return ctx.JSON(http.StatusOK, res)
}

func (h *OnboardingRoute) changeMerchantStatus(ctx echo.Context) error {
	authUser := common.ExtractUserContext(ctx)
	req := &grpc.MerchantChangeStatusRequest{}
	err := (&common.OnboardingChangeMerchantStatusBinder{}).Bind(req, ctx)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestParamsIncorrect)
	}

	err = h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	req.UserId = authUser.Id
	res, err := h.dispatch.Services.Billing.ChangeMerchantStatus(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "ChangeMerchantStatus", req)
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorUnknown)
	}

	if res.Status != http.StatusOK {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.JSON(http.StatusOK, res.Item)
}

func (h *OnboardingRoute) createNotification(ctx echo.Context) error {
	authUser := common.ExtractUserContext(ctx)
	req := &grpc.NotificationRequest{}
	err := (&common.OnboardingCreateNotificationBinder{}).Bind(req, ctx)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestParamsIncorrect)
	}

	err = h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	req.UserId = authUser.Id
	res, err := h.dispatch.Services.Billing.CreateNotification(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "CreateNotification", req)
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorUnknown)
	}

	if res.Status != http.StatusOK {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.JSON(http.StatusCreated, res.Item)
}

func (h *OnboardingRoute) getNotification(ctx echo.Context) error {
	merchantId := ctx.Param(common.RequestParameterMerchantId)
	notificationId := ctx.Param(common.RequestParameterNotificationId)

	if merchantId == "" || bson.IsObjectIdHex(merchantId) == false {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorIncorrectMerchantId)
	}

	if notificationId == "" || bson.IsObjectIdHex(notificationId) == false {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorIncorrectNotificationId)
	}

	req := &grpc.GetNotificationRequest{
		MerchantId:     merchantId,
		NotificationId: notificationId,
	}
	res, err := h.dispatch.Services.Billing.GetNotification(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "GetNotification", req)
		return echo.NewHTTPError(http.StatusNotFound, common.ErrorNotificationNotFound)
	}

	return ctx.JSON(http.StatusOK, res)
}

func (h *OnboardingRoute) listNotifications(ctx echo.Context) error {
	req := &grpc.ListingNotificationRequest{}
	err := (&common.OnboardingNotificationsListBinder{
		LimitDefault:  h.cfg.LimitDefault,
		OffsetDefault: h.cfg.OffsetDefault,
	}).Bind(req, ctx)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestParamsIncorrect)
	}

	err = h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.ListNotifications(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "ListNotifications", req)
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorUnknown)
	}

	return ctx.JSON(http.StatusOK, res)
}

func (h *OnboardingRoute) markAsReadNotification(ctx echo.Context) error {
	merchantId := ctx.Param(common.RequestParameterMerchantId)
	notificationId := ctx.Param(common.RequestParameterNotificationId)

	if merchantId == "" || bson.IsObjectIdHex(merchantId) == false {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorIncorrectMerchantId)
	}

	if notificationId == "" || bson.IsObjectIdHex(notificationId) == false {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorIncorrectNotificationId)
	}

	req := &grpc.GetNotificationRequest{
		MerchantId:     merchantId,
		NotificationId: notificationId,
	}
	res, err := h.dispatch.Services.Billing.MarkNotificationAsRead(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "MarkNotificationAsRead", req)
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorUnknown)
	}

	return ctx.JSON(http.StatusOK, res)
}

func (h *OnboardingRoute) changeAgreement(ctx echo.Context) error {
	req := &grpc.ChangeMerchantDataRequest{}
	binder := common.NewChangeMerchantDataRequestBinder(h.dispatch, h.cfg)
	err := binder.Bind(req, ctx)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestParamsIncorrect)
	}

	err = h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.ChangeMerchantData(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "ChangeMerchantData", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}

	if res.Status != pkg.ResponseStatusOk {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.JSON(http.StatusOK, res.Item)
}

func (h *OnboardingRoute) generateAgreement(ctx echo.Context) error {
	merchantId := ctx.Param(common.RequestParameterId)

	if merchantId == "" || bson.IsObjectIdHex(merchantId) == false {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestParamsIncorrect)
	}

	ctxReq := ctx.Request().Context()
	req := &grpc.GetMerchantByRequest{MerchantId: merchantId}
	res, err := h.dispatch.Services.Billing.GetMerchantBy(ctxReq, req)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}

	if res.Status != pkg.ResponseStatusOk {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	if res.Item.S3AgreementName != "" {
		filePath := os.TempDir() + string(os.PathSeparator) + res.Item.S3AgreementName
		_, err = h.awsManager.Download(ctxReq, filePath, &awsWrapper.DownloadInput{FileName: res.Item.S3AgreementName})

		if err != nil {
			h.L().Error("AWS api call to download file failed", logger.PairArgs("err", err.Error(), "file_name", res.Item.S3AgreementName))

			return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
		}

		fData, err := h.getAgreementStructure(ctx, merchantId, agreementExtension, agreementContentType, filePath)

		if err != nil {
			h.L().Error("Get agreement structure failed", logger.PairArgs("err", err.Error(), "merchant_id", merchantId))

			return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorInternal)
		}

		return ctx.JSON(http.StatusOK, fData)
	}

	if res.Item.CanGenerateAgreement() == false {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorMessageAgreementCanNotBeGenerate)
	}

	buf := new(bytes.Buffer)
	data := map[string]interface{}{"Merchant": res.Item}
	err = ctx.Echo().Renderer.Render(buf, common.AgreementPageTemplateName, data, ctx)

	if err != nil {
		h.L().Error("Render agreement document failed", logger.PairArgs("err", err.Error(), "merchant_id", merchantId))

		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorInternal)
	}

	pdf, err := wkhtmltopdf.NewPDFGenerator()

	if err != nil {
		h.L().Error("New pdf generator instance create failed", logger.PairArgs("err", err.Error(), "merchant_id", merchantId))

		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorInternal)
	}

	pdf.AddPage(wkhtmltopdf.NewPageReader(strings.NewReader(buf.String())))
	err = pdf.Create()

	if err != nil {
		h.L().Error("Create pdf document of agreement failed", logger.PairArgs("err", err.Error(), "merchant_id", merchantId))

		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorInternal)
	}

	fileName := fmt.Sprintf(agreementFileMask, res.Item.Id)
	filePath := os.TempDir() + string(os.PathSeparator) + fileName
	err = pdf.WriteFile(filePath)

	if err != nil {
		h.L().Error("Write generated agreement failed", logger.PairArgs("err", err.Error(), "merchant_id", merchantId))

		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorInternal)
	}

	_, err = h.awsManager.Upload(ctxReq, &awsWrapper.UploadInput{Body: bytes.NewReader(pdf.Bytes()), FileName: fileName})

	if err != nil {
		h.L().Error("AWS api call to upload file failed", logger.PairArgs("err", err.Error(), "file_name", fileName))

		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorInternal)
	}

	req1 := &grpc.SetMerchantS3AgreementRequest{MerchantId: merchantId, S3AgreementName: fileName}
	_, err = h.dispatch.Services.Billing.SetMerchantS3Agreement(ctx.Request().Context(), req1)

	if err != nil {
		h.L().Error(pkg.ErrorGrpcServiceCallFailed, logger.PairArgs("err", err.Error(), common.ErrorFieldService, pkg.ServiceName, common.ErrorFieldMethod, "SetMerchantS3Agreement", common.ErrorFieldRequest, req1))
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}

	fData, err := h.getAgreementStructure(ctx, merchantId, agreementExtension, agreementContentType, filePath)

	if err != nil {
		h.L().Error("Get agreement structure failed", logger.PairArgs("err", err.Error(), "merchant_id", merchantId))

		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorInternal)
	}

	return ctx.JSON(http.StatusOK, fData)
}

func (h *OnboardingRoute) getAgreementDocument(ctx echo.Context) error {
	merchantId := ctx.Param(common.RequestParameterId)

	if merchantId == "" || bson.IsObjectIdHex(merchantId) == false {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestParamsIncorrect)
	}

	ctxReq := ctx.Request().Context()
	req := &grpc.GetMerchantByRequest{MerchantId: merchantId}
	res, err := h.dispatch.Services.Billing.GetMerchantBy(ctxReq, req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "GetMerchantBy", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}

	if res.Status != pkg.ResponseStatusOk {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	if res.Item.S3AgreementName == "" {
		return echo.NewHTTPError(http.StatusNotFound, common.ErrorMessageAgreementNotGenerated)
	}

	filePath := os.TempDir() + string(os.PathSeparator) + res.Item.S3AgreementName
	_, err = h.awsManager.Download(ctxReq, filePath, &awsWrapper.DownloadInput{FileName: res.Item.S3AgreementName})

	if err != nil {
		h.L().Error("AWS api call to download file failed", logger.PairArgs("err", err.Error(), "file_name", res.Item.S3AgreementName))

		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorAgreementFileNotExist)
	}

	return ctx.File(filePath)
}

func (h *OnboardingRoute) uploadAgreementDocument(ctx echo.Context) error {
	merchantId := ctx.Param(common.RequestParameterId)

	if merchantId == "" || bson.IsObjectIdHex(merchantId) == false {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestParamsIncorrect)
	}

	ctxReq := ctx.Request().Context()
	req := &grpc.GetMerchantByRequest{MerchantId: merchantId}
	res, err := h.dispatch.Services.Billing.GetMerchantBy(ctxReq, req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "GetMerchantBy", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}

	if res.Status != pkg.ResponseStatusOk {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	file, err := ctx.FormFile(common.RequestParameterFile)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorNotMultipartForm)
	}

	src, err := h.validateUpload(file)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	defer src.Close()

	fileName := fmt.Sprintf(agreementFileMask, res.Item.Id)
	filePath := os.TempDir() + string(os.PathSeparator) + fileName
	dst, err := os.Create(filePath)

	if err != nil {
		h.L().Error("Upload new version of agreement failed", logger.PairArgs("err", err.Error(), "merchant_id", merchantId))

		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorInternal)
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)

	if err != nil {
		h.L().Error("Upload new version of agreement failed", logger.PairArgs("err", err.Error(), "merchant_id", merchantId))

		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorInternal)
	}

	_, err = h.awsManager.Upload(ctxReq, &awsWrapper.UploadInput{Body: dst, FileName: fileName})

	if err != nil {
		h.L().Error("AWS api call to upload file failed", logger.PairArgs("err", err.Error(), "file_name", res.Item.S3AgreementName))

		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUploadFailed)
	}

	req1 := &grpc.SetMerchantS3AgreementRequest{MerchantId: merchantId, S3AgreementName: fileName}
	_, err = h.dispatch.Services.Billing.SetMerchantS3Agreement(ctx.Request().Context(), req1)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "SetMerchantS3Agreement", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}

	fData, err := h.getAgreementStructure(ctx, merchantId, agreementExtension, agreementContentType, filePath)

	if err != nil {
		h.L().Error("Get agreement structure failed", logger.PairArgs("err", err.Error(), "merchant_id", merchantId))

		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorInternal)
	}

	return ctx.JSON(http.StatusOK, fData)
}

func (h *OnboardingRoute) getAgreementStructure(
	ctx echo.Context,
	merchantId, ext, ct, fPath string,
) (interface{}, error) {
	file, err := os.Open(fPath)

	if err != nil {
		return nil, common.ErrorMessageAgreementNotFound
	}

	defer func() {
		if err := file.Close(); err != nil {
			return
		}
	}()

	fi, err := file.Stat()

	if err != nil {
		return nil, common.ErrorMessageAgreementNotFound
	}

	data := &OnboardingFileData{
		Url: fmt.Sprintf(agreementUrlMask, ctx.Request().URL.Scheme, ctx.Request().URL.Host, merchantId),
		Metadata: &OnboardingFileMetadata{
			Name:        fi.Name(),
			Extension:   ext,
			ContentType: ct,
			Size:        fi.Size(),
		},
	}

	return data, nil
}

func (h *OnboardingRoute) validateUpload(file *multipart.FileHeader) (multipart.File, error) {
	if file.Size > agreementUploadMaxSize {
		return nil, common.ErrorMessageAgreementUploadMaxSize
	}

	src, err := file.Open()

	if err != nil {
		h.L().Error("validate upload error", logger.PairArgs("err", err.Error()))
		return nil, common.ErrorUnknown
	}

	buffer := make([]byte, 512)
	_, err = src.Read(buffer)

	if err != nil {
		h.L().Error("validate upload error", logger.PairArgs("err", err.Error()))
		return nil, common.ErrorUnknown
	}

	_, err = src.Seek(0, 0)

	if err != nil {
		h.L().Error("validate upload error", logger.PairArgs("err", err.Error()))
		return nil, common.ErrorUnknown
	}

	ct := http.DetectContentType(buffer)

	if ct != agreementContentType {
		return nil, common.ErrorMessageAgreementContentType
	}

	return src, nil
}

// @Description Set company information in merchant onboarding process
// @Example curl -X PUT -H 'Authorization: Bearer %access_token_here%' -H 'Content-Type: application/json' \
//  -d '{"name": "Roga and Copita LLC", "alternative_name": "Apple Inc", "website": "http://localhost", "country": "RU",
//    	"state": "St.Petersburg", "zip": "190000", "city": "St.Petersburg", "address": "Nevskiy st. 1"}' \
//  https://api.paysuper.online/admin/api/v1/merchants/company
//
// @Example curl -X PUT -H 'Authorization: Bearer %access_token_here%' -H 'Content-Type: application/json' \
//  -d '{"name": "Roga and Copita LLC", "alternative_name": "Apple Inc", "website": "http://localhost", "country": "RU",
//    	"state": "St.Petersburg", "zip": "190000", "city": "St.Petersburg", "address": "Nevskiy st. 2"}' \
//  https://api.paysuper.online/admin/api/v1/merchants/5d4847f61986ee46ec581e26/company
func (h *OnboardingRoute) setMerchantCompany(ctx echo.Context) error {
	authUser := common.ExtractUserContext(ctx)
	in := &billing.MerchantCompanyInfo{}
	err := ctx.Bind(in)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestParamsIncorrect)
	}

	req := &grpc.OnboardingRequest{
		Company: in,
		Id:      ctx.Param(common.RequestParameterId),
		User: &billing.MerchantUser{
			Id:    authUser.Id,
			Email: authUser.Email,
		},
	}
	err = h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.ChangeMerchant(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "ChangeMerchant", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}

	if res.Status != pkg.ResponseStatusOk {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.JSON(http.StatusOK, res.Item)
}

// @Description Set company contact information in merchant onboarding process
// @Example curl -X PUT -H 'Authorization: Bearer %access_token_here%' -H 'Content-Type: application/json' \
//  -d '{"authorized": {"name": "Unit Test", "email": "test@unit.test", "phone": "1234567890", "position": "CEO"},
//    	"technical": {"name": "Unit Test", "email": "test@unit.test", "phone": "1234567890"}}' \
//  https://api.paysuper.online/admin/api/v1/merchants/contacts
//
// @Example curl -X PUT -H 'Authorization: Bearer %access_token_here%' -H 'Content-Type: application/json' \
//  -d '{"authorized": {"name": "Unit Test", "email": "test@unit.test", "phone": "1234567891", "position": "CEO"},
//    	"technical": {"name": "Unit Test", "email": "test@unit.test", "phone": "1234567890"}}' \
//  https://api.paysuper.online/admin/api/v1/merchants/5d4847f61986ee46ec581e26/contacts
func (h *OnboardingRoute) setMerchantContacts(ctx echo.Context) error {
	authUser := common.ExtractUserContext(ctx)
	in := &billing.MerchantContact{}
	err := ctx.Bind(in)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestParamsIncorrect)
	}

	req := &grpc.OnboardingRequest{
		Contacts: in,
		Id:       ctx.Param(common.RequestParameterId),
		User: &billing.MerchantUser{
			Id:    authUser.Id,
			Email: authUser.Email,
		},
	}
	err = h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.ChangeMerchant(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "ChangeMerchant", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}

	if res.Status != pkg.ResponseStatusOk {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.JSON(http.StatusOK, res.Item)
}

// @Description Set company banking information in merchant onboarding process
// @Example curl -X PUT -H 'Authorization: Bearer %access_token_here%' -H 'Content-Type: application/json' \
//  -d '{"name": "Bank Name-Spb.", "address": "St.Petersburg, Nevskiy st. 1",
//  	"account_number": "408000000001", "swift": "ALFARUMM", "correspondent_account": "408000000001"}' \
//  https://api.paysuper.online/admin/api/v1/merchants/banking
//
// @Example curl -X PUT -H 'Authorization: Bearer %access_token_here%' -H 'Content-Type: application/json' \
//  -d '{"name": "Bank Name-Spb.", "address": "St.Petersburg, Nevskiy st. 1",
//  	"account_number": "408000000001", "swift": "ALFARUMM", "correspondent_account": "408000000002"}' \
//  https://api.paysuper.online/admin/api/v1/merchants/5d4847f61986ee46ec581e26/banking
func (h *OnboardingRoute) setMerchantBanking(ctx echo.Context) error {
	authUser := common.ExtractUserContext(ctx)
	in := &billing.MerchantBanking{}
	err := ctx.Bind(in)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestParamsIncorrect)
	}

	req := &grpc.OnboardingRequest{
		Banking: in,
		Id:      ctx.Param(common.RequestParameterId),
		User: &billing.MerchantUser{
			Id:    authUser.Id,
			Email: authUser.Email,
		},
	}
	err = h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.ChangeMerchant(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "ChangeMerchant", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}

	if res.Status != pkg.ResponseStatusOk {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.JSON(http.StatusOK, res.Item)
}

// @Description Get merchant completion information
// @Example curl -X GET -H 'Authorization: Bearer %access_token_here%' -H 'Content-Type: application/json' \
//	-d '{"signer_type": 0}'
// https://api.paysuper.online/admin/api/v1/merchants/5d4847f61986ee46ec581e26/status
func (h *OnboardingRoute) getMerchantStatus(ctx echo.Context) error {
	req := &grpc.SetMerchantS3AgreementRequest{
		MerchantId: ctx.Param(common.RequestParameterId),
	}

	err := h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.GetMerchantOnboardingCompleteData(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "ChangeMerchant", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}

	if res.Status != pkg.ResponseStatusOk {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.JSON(http.StatusOK, res.Item)
}

// @Description get hellosign (https://www.hellosign.com) signature to sign license agreement
// @Example @Example curl -X PUT -H 'Authorization: Bearer %access_token_here%' -H 'Content-Type: application/json' \
// 		https://api.paysuper.online/admin/api/v1/merchants/ffffffffffffffffffffffff/agreement/signature
func (h *OnboardingRoute) createAgreementSignature(ctx echo.Context) error {
	req := &grpc.GetMerchantAgreementSignUrlRequest{}
	err := ctx.Bind(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestParamsIncorrect)
	}

	req.MerchantId = ctx.Param(common.RequestParameterId)
	err = h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	res, err := h.dispatch.Services.Billing.GetMerchantAgreementSignUrl(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "GetMerchantAgreementSignUrl", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}

	if res.Status != pkg.ResponseStatusOk {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.JSON(http.StatusOK, res.Item)
}

// @Description get list of merchants tariffs rates by conditions
// @Example @Example curl -X GET -H 'Authorization: Bearer %access_token_here%' -H 'Content-Type: application/json' \
// 		https://api.paysuper.online/admin/api/v1/merchants/tariffs?region=CIS&payout_currency=USD
func (h *OnboardingRoute) getTariffRates(ctx echo.Context) error {
	req := &grpc.GetMerchantTariffRatesRequest{}
	err := ctx.Bind(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestParamsIncorrect)
	}

	err = h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	req.Region = common.TariffRegions[req.Region]
	res, err := h.dispatch.Services.Billing.GetMerchantTariffRates(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "GetMerchantTariffRates", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}

	if res.Status != pkg.ResponseStatusOk {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.JSON(http.StatusOK, res.Item)
}

// @Description set tariff to merchant
// @Example @Example curl -X POST -H 'Authorization: Bearer %access_token_here%' -H 'Content-Type: application/json' \
//		-d '{"region": "CIS", "payout_currency": "USD", "amount_from": 0.75, "amount_to": 5}'
// 		https://api.paysuper.online/admin/api/v1/merchants/ffffffffffffffffffffffff/tariffs
func (h *OnboardingRoute) setTariffRates(ctx echo.Context) error {
	req := &grpc.SetMerchantTariffRatesRequest{}
	err := ctx.Bind(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorRequestParamsIncorrect)
	}

	req.MerchantId = ctx.Param(common.RequestParameterId)
	err = h.dispatch.Validate.Struct(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.GetValidationError(err))
	}

	req.Region = common.TariffRegions[req.Region]
	res, err := h.dispatch.Services.Billing.SetMerchantTariffRates(ctx.Request().Context(), req)

	if err != nil {
		common.LogSrvCallFailedGRPC(h.L(), err, pkg.ServiceName, "SetMerchantTariffRates", req)
		return echo.NewHTTPError(http.StatusInternalServerError, common.ErrorUnknown)
	}

	if res.Status != pkg.ResponseStatusOk {
		return echo.NewHTTPError(int(res.Status), res.Message)
	}

	return ctx.NoContent(http.StatusOK)
}
