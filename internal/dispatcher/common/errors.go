package common

import (
	"github.com/paysuper/paysuper-billing-server/pkg/proto/grpc"
)

// NewManagementApiResponseError
func NewManagementApiResponseError(code, msg string, details ...string) *grpc.ResponseErrorMessage {
	var det string
	if len(details) > 0 && details[0] != "" {
		det = details[0]
	} else {
		det = ""
	}
	return &grpc.ResponseErrorMessage{Code: code, Message: msg, Details: det}
}

// NewValidationError
func NewValidationError(details string) *grpc.ResponseErrorMessage {
	return NewManagementApiResponseError(ErrorValidationFailed.Code, ErrorValidationFailed.Message, details)
}

const (
	ErrorNamespaceMerchantCompanyInfoName                 = "OnboardingRequest.Company.Name"
	ErrorNamespaceMerchantCompanyInfoAlternativeName      = "OnboardingRequest.Company.AlternativeName"
	ErrorNamespaceMerchantCompanyInfoWebsite              = "OnboardingRequest.Company.Website"
	ErrorNamespaceMerchantCompanyInfoCountry              = "OnboardingRequest.Company.Country"
	ErrorNamespaceMerchantCompanyInfoState                = "OnboardingRequest.Company.State"
	ErrorNamespaceMerchantCompanyInfoZip                  = "OnboardingRequest.Company.Zip"
	ErrorNamespaceMerchantCompanyInfoCity                 = "OnboardingRequest.Company.City"
	ErrorNamespaceMerchantCompanyInfoAddress              = "OnboardingRequest.Company.Address"
	ErrorNamespaceMerchantContactAuthorized               = "OnboardingRequest.Contacts.Authorized"
	ErrorNamespaceMerchantContactTechnical                = "OnboardingRequest.Contacts.Technical"
	ErrorNamespaceMerchantContactAuthorizedName           = "OnboardingRequest.Contacts.Authorized.Name"
	ErrorNamespaceMerchantContactAuthorizedEmail          = "OnboardingRequest.Contacts.Authorized.Email"
	ErrorNamespaceMerchantContactAuthorizedPhone          = "OnboardingRequest.Contacts.Authorized.Phone"
	ErrorNamespaceMerchantContactAuthorizedPosition       = "OnboardingRequest.Contacts.Authorized.Position"
	ErrorNamespaceMerchantContactTechnicalName            = "OnboardingRequest.Contacts.Technical.Name"
	ErrorNamespaceMerchantContactTechnicalEmail           = "OnboardingRequest.Contacts.Technical.Email"
	ErrorNamespaceMerchantContactTechnicalPhone           = "OnboardingRequest.Contacts.Technical.Phone"
	ErrorNamespaceMerchantBankingCurrency                 = "OnboardingRequest.Banking.Currency"
	ErrorNamespaceMerchantBankingName                     = "OnboardingRequest.Banking.Name"
	ErrorNamespaceMerchantBankingAddress                  = "OnboardingRequest.Banking.Address"
	ErrorNamespaceMerchantBankingAccountNumber            = "OnboardingRequest.Banking.AccountNumber"
	ErrorNamespaceMerchantBankingSwift                    = "OnboardingRequest.Banking.Swift"
	ErrorNamespaceMerchantBankingCorrespondentAccount     = "OnboardingRequest.Banking.CorrespondentAccount"
	ErrorNamespaceGetDashboardMainRequestPeriod           = "GetDashboardMainRequest.Period"
	ErrorNamespaceGetDashboardMainRequestMerchantId       = "GetDashboardMainRequest.MerchantId"
	ErrorNamespaceGetDashboardBaseReportRequestPeriod     = "GetDashboardBaseReportRequest.Period"
	ErrorNamespaceGetDashboardBaseReportRequestMerchantId = "GetDashboardBaseReportRequest.MerchantId"
)

var (
	ErrorUnknown                                      = NewManagementApiResponseError("ma000001", "unknown error. try request later")
	ErrorValidationFailed                             = NewManagementApiResponseError("ma000002", "validation failed")
	ErrorInternal                                     = NewManagementApiResponseError("ma000003", InternalErrorTemplate)
	ErrorMessageAccessDenied                          = NewManagementApiResponseError("ma000004", "access denied")
	ErrorIdIsEmpty                                    = NewManagementApiResponseError("ma000005", "identifier can't be empty")
	ErrorIncorrectMerchantId                          = NewManagementApiResponseError("ma000006", "incorrect merchant identifier")
	ErrorIncorrectNotificationId                      = NewManagementApiResponseError("ma000007", "incorrect notification identifier")
	ErrorIncorrectOrderId                             = NewManagementApiResponseError("ma000008", "incorrect order identifier")
	ErrorIncorrectProductId                           = NewManagementApiResponseError("ma000009", "incorrect product identifier")
	ErrorIncorrectCountryIdentifier                   = NewManagementApiResponseError("ma000010", "incorrect country identifier")
	ErrorIncorrectCurrencyIdentifier                  = NewManagementApiResponseError("ma000011", "incorrect currency identifier")
	ErrorMessageOrdersNotFound                        = NewManagementApiResponseError("ma000012", "orders not found")
	ErrorCountryNotFound                              = NewManagementApiResponseError("ma000013", "country not found")
	ErrorCurrencyNotFound                             = NewManagementApiResponseError("ma000014", "currency not found")
	ErrorNotificationNotFound                         = NewManagementApiResponseError("ma000015", "notification not found")
	ErrorMessageAgreementCanNotBeGenerate             = NewManagementApiResponseError("ma000020", "agreement can't be generated for not checked merchant data")
	ErrorMessageAgreementNotGenerated                 = NewManagementApiResponseError("ma000021", "agreement for merchant not generated early")
	ErrorMessageSignatureHeaderIsEmpty                = NewManagementApiResponseError("ma000022", "header with request signature can't be empty")
	ErrorRequestParamsIncorrect                       = NewManagementApiResponseError("ma000023", "incorrect request parameters")
	ErrorEmailFieldIncorrect                          = NewManagementApiResponseError("ma000024", "incorrect email")
	ErrorRequestDataInvalid                           = NewManagementApiResponseError("ma000026", "request data invalid")
	ErrorCountriesListError                           = NewManagementApiResponseError("ma000027", "countries list error")
	ErrorAgreementFileNotExist                        = NewManagementApiResponseError("ma000028", "file for the specified key does not exist")
	ErrorNotMultipartForm                             = NewManagementApiResponseError("ma000029", "no multipart boundary param in Content-Type")
	ErrorUploadFailed                                 = NewManagementApiResponseError("ma000030", "upload failed")
	ErrorIncorrectProjectId                           = NewManagementApiResponseError("ma000031", "incorrect project identifier")
	ErrorIncorrectPaymentMethodId                     = NewManagementApiResponseError("ma000032", "incorrect payment method identifier")
	ErrorIncorrectPaylinkId                           = NewManagementApiResponseError("ma000033", "incorrect paylink identifier")
	ErrorMessageAuthorizationHeaderNotFound           = NewManagementApiResponseError("ma000034", "authorization header not found")
	ErrorMessageAuthorizationTokenNotFound            = NewManagementApiResponseError("ma000035", "authorization token not found")
	ErrorMessageAuthorizedUserNotFound                = NewManagementApiResponseError("ma000036", "information about authorized user not found")
	ErrorMessageStatusIncorrectType                   = NewManagementApiResponseError("ma000037", "status parameter has incorrect type")
	ErrorMessageAgreementNotFound                     = NewManagementApiResponseError("ma000038", "agreement for merchant not found")
	ErrorMessageAgreementUploadMaxSize                = NewManagementApiResponseError("ma000039", "agreement document max upload size exceeded")
	ErrorMessageAgreementContentType                  = NewManagementApiResponseError("ma000040", "agreement document type must be a pdf")
	ErrorMessageAgreementTypeIncorrectType            = NewManagementApiResponseError("ma000041", "agreement type parameter have incorrect type")
	ErrorMessageHasMerchantSignatureIncorrectType     = NewManagementApiResponseError("ma000042", "merchant signature parameter has incorrect type")
	ErrorMessageHasPspSignatureIncorrectType          = NewManagementApiResponseError("ma000043", "paysuper signature parameter has incorrect type")
	ErrorMessageAgreementSentViaMailIncorrectType     = NewManagementApiResponseError("ma000044", "agreement sent via email parameter has incorrect type")
	ErrorMessageMailTrackingLinkIncorrectType         = NewManagementApiResponseError("ma000045", "mail tracking link parameter has incorrect type")
	ErrorMessageNameIncorrectType                     = NewManagementApiResponseError("ma000046", "name parameter has incorrect type")
	ErrorMessageImageIncorrectType                    = NewManagementApiResponseError("ma000047", "image parameter has incorrect type")
	ErrorMessageCallbackCurrencyIncorrectType         = NewManagementApiResponseError("ma000048", "callback currency parameter has incorrect type")
	ErrorMessageCallbackProtocolIncorrectType         = NewManagementApiResponseError("ma000049", "callback protocol parameter has incorrect type")
	ErrorMessageCreateOrderAllowedUrlsIncorrectType   = NewManagementApiResponseError("ma000050", "create order allowed urls parameter has incorrect type")
	ErrorMessageAllowDynamicNotifyUrlsIncorrectType   = NewManagementApiResponseError("ma000051", "allow dynamic notify urls parameter has incorrect type")
	ErrorMessageAllowDynamicRedirectUrlsIncorrectType = NewManagementApiResponseError("ma000052", "allow dynamic redirect urls parameter has incorrect type")
	ErrorMessageLimitsCurrencyIncorrectType           = NewManagementApiResponseError("ma000053", "limits currency parameter has incorrect type")
	ErrorMessageMinPaymentAmountIncorrectType         = NewManagementApiResponseError("ma000054", "min payment amount parameter has incorrect type")
	ErrorMessageMaxPaymentAmountIncorrectType         = NewManagementApiResponseError("ma000055", "max payment amount parameter has incorrect type")
	ErrorMessageNotifyEmailsIncorrectType             = NewManagementApiResponseError("ma000056", "notify emails parameter has incorrect type")
	ErrorMessageIsProductsCheckoutIncorrectType       = NewManagementApiResponseError("ma000057", "is products checkout parameter has incorrect type")
	ErrorMessageSecretKeyIncorrectType                = NewManagementApiResponseError("ma000058", "secret key parameter has incorrect type")
	ErrorMessageSignatureRequiredIncorrectType        = NewManagementApiResponseError("ma000059", "signature required parameter has incorrect type")
	ErrorMessageSendNotifyEmailIncorrectType          = NewManagementApiResponseError("ma000060", "send notify email parameter has incorrect type")
	ErrorMessageUrlCheckAccountIncorrectType          = NewManagementApiResponseError("ma000061", "url check account parameter has incorrect type")
	ErrorMessageUrlProcessPaymentIncorrectType        = NewManagementApiResponseError("ma000062", "url process payment parameter has incorrect type")
	ErrorMessageUrlRedirectFailIncorrectType          = NewManagementApiResponseError("ma000063", "url redirect fail parameter has incorrect type")
	ErrorMessageUrlRedirectSuccessIncorrectType       = NewManagementApiResponseError("ma000064", "url redirect success parameter has incorrect type")
	ErrorMessageUrlChargebackPayment                  = NewManagementApiResponseError("ma000065", "url chargeback payment parameter has incorrect type")
	ErrorMessageUrlCancelPayment                      = NewManagementApiResponseError("ma000066", "url cancel payment parameter has incorrect type")
	ErrorMessageUrlFraudPayment                       = NewManagementApiResponseError("ma000067", "url fraud payment parameter has incorrect type")
	ErrorMessageUrlRefundPayment                      = NewManagementApiResponseError("ma000068", "url refund payment parameter has incorrect type")
	ErrorMessagePriceGroupByCountry                   = NewManagementApiResponseError("ma000069", "unable to get price group by country")
	ErrorMessagePriceGroupCurrencyList                = NewManagementApiResponseError("ma000070", "unable to get price group currencies")
	ErrorMessagePriceGroupCurrencyByRegion            = NewManagementApiResponseError("ma000071", "unable to get price group currency by region")
	ErrorMessagePriceGroupRecommendedList             = NewManagementApiResponseError("ma000072", "unable to get price group recommended prices")
	ErrorMessageGetProductPrice                       = NewManagementApiResponseError("ma000072", "unable to get price of product")
	ErrorMessageUpdateProductPrice                    = NewManagementApiResponseError("ma000072", "unable to update price of product")
	ErrorMessageIncorrectZip                          = NewManagementApiResponseError("ma000073", "incorrect zip code")
	ErrorMessageIncorrectNumberOfEmployees            = NewManagementApiResponseError("ma000074", "incorrect number of employees value")
	ErrorMessageIncorrectAnnualIncome                 = NewManagementApiResponseError("ma000075", "incorrect annual income value")
	ErrorMessageIncorrectCompanyName                  = NewManagementApiResponseError("ma000076", "incorrect company name")
	ErrorMessageIncorrectPosition                     = NewManagementApiResponseError("ma000077", "incorrect position")
	ErrorMessageIncorrectFirstName                    = NewManagementApiResponseError("ma000078", "incorrect first name")
	ErrorMessageIncorrectLastName                     = NewManagementApiResponseError("ma000079", "incorrect last name")
	ErrorMessageIncorrectWebsite                      = NewManagementApiResponseError("ma000080", "incorrect website")
	ErrorMessageIncorrectKindOfActivity               = NewManagementApiResponseError("ma000081", "incorrect kind of activity")
	ErrorMessageIncorrectReview                       = NewManagementApiResponseError("ma000082", "review must be text with length lower than or equal 500 characters")
	ErrorMessageIncorrectPageId                       = NewManagementApiResponseError("ma000083", "review page identifier must be one of next values: primary_onboarding, merchant_onboarding")
	ErrorMessageKeyProductIdInvalid                   = NewManagementApiResponseError("ma000082", "key product id is invalid")
	ErrorMessagePlatformIdInvalid                     = NewManagementApiResponseError("ma000083", "platform id is invalid")

	ErrorMessageIncorrectAlternativeName          = NewManagementApiResponseError("ma000084", "incorrect brand")
	ErrorMessageIncorrectState                    = NewManagementApiResponseError("ma000085", "incorrect state")
	ErrorMessageIncorrectCity                     = NewManagementApiResponseError("ma000086", "incorrect city")
	ErrorMessageIncorrectAddress                  = NewManagementApiResponseError("ma000087", "incorrect address")
	ErrorMessageRequiredContactAuthorized         = NewManagementApiResponseError("ma000088", "company authorized contact information is required")
	ErrorMessageRequiredContactTechnical          = NewManagementApiResponseError("ma000089", "company technical contact information is required")
	ErrorMessageIncorrectName                     = NewManagementApiResponseError("ma000090", "incorrect name")
	ErrorMessageIncorrectPhone                    = NewManagementApiResponseError("ma000091", "incorrect phone")
	ErrorMessageIncorrectBankName                 = NewManagementApiResponseError("ma000092", "incorrect bank name")
	ErrorMessageIncorrectBankAddress              = NewManagementApiResponseError("ma000093", "incorrect bank address")
	ErrorMessageIncorrectBankAccountNumber        = NewManagementApiResponseError("ma000094", "incorrect bank accounting number")
	ErrorMessageIncorrectBankSwift                = NewManagementApiResponseError("ma000095", "incorrect bank swift code")
	ErrorMessageIncorrectBankCorrespondentAccount = NewManagementApiResponseError("ma000096", "incorrect bank correspondent account")
	ErrorMessageFileNotFound                      = NewManagementApiResponseError("ma000097", "file with key was not specified")
	ErrorMessageCantReadFile                      = NewManagementApiResponseError("ma000098", "file can not be read")
	ErrorIncorrectPeriod                          = NewManagementApiResponseError("ma000099", "incorrect period")
	ErrorMessageMerchantNotFound                  = NewManagementApiResponseError("ma000100", "merchant not found")
	ErrorMessageCreateReportFile                  = NewManagementApiResponseError("ma000101", "unable to create report file")
	ErrorMessageDownloadReportFile                = NewManagementApiResponseError("ma000102", "unable to download report file")

	ValidationErrors = map[string]*grpc.ResponseErrorMessage{
		UserProfileFieldNumberOfEmployees: ErrorMessageIncorrectNumberOfEmployees,
		UserProfileFieldAnnualIncome:      ErrorMessageIncorrectAnnualIncome,
		UserProfileFieldCompanyName:       ErrorMessageIncorrectCompanyName,
		UserProfileFieldPosition:          ErrorMessageIncorrectPosition,
		UserProfileFieldFirstName:         ErrorMessageIncorrectFirstName,
		UserProfileFieldLastName:          ErrorMessageIncorrectLastName,
		UserProfileFieldWebsite:           ErrorMessageIncorrectWebsite,
		UserProfileFieldKindOfActivity:    ErrorMessageIncorrectKindOfActivity,
		UserProfileFieldReview:            ErrorMessageIncorrectReview,
		UserProfileFieldPageId:            ErrorMessageIncorrectPageId,
	}

	ValidationNamespaceErrors = map[string]*grpc.ResponseErrorMessage{
		ErrorNamespaceMerchantCompanyInfoName:                 ErrorMessageIncorrectCompanyName,
		ErrorNamespaceMerchantCompanyInfoAlternativeName:      ErrorMessageIncorrectAlternativeName,
		ErrorNamespaceMerchantCompanyInfoWebsite:              ErrorMessageIncorrectWebsite,
		ErrorNamespaceMerchantCompanyInfoCountry:              ErrorIncorrectCountryIdentifier,
		ErrorNamespaceMerchantCompanyInfoState:                ErrorMessageIncorrectState,
		ErrorNamespaceMerchantCompanyInfoZip:                  ErrorMessageIncorrectZip,
		ErrorNamespaceMerchantCompanyInfoCity:                 ErrorMessageIncorrectCity,
		ErrorNamespaceMerchantCompanyInfoAddress:              ErrorMessageIncorrectAddress,
		ErrorNamespaceMerchantContactAuthorized:               ErrorMessageRequiredContactAuthorized,
		ErrorNamespaceMerchantContactTechnical:                ErrorMessageRequiredContactTechnical,
		ErrorNamespaceMerchantContactAuthorizedName:           ErrorMessageIncorrectName,
		ErrorNamespaceMerchantContactAuthorizedEmail:          ErrorEmailFieldIncorrect,
		ErrorNamespaceMerchantContactAuthorizedPhone:          ErrorMessageIncorrectPhone,
		ErrorNamespaceMerchantContactAuthorizedPosition:       ErrorMessageIncorrectPosition,
		ErrorNamespaceMerchantContactTechnicalName:            ErrorMessageIncorrectName,
		ErrorNamespaceMerchantContactTechnicalEmail:           ErrorEmailFieldIncorrect,
		ErrorNamespaceMerchantContactTechnicalPhone:           ErrorMessageIncorrectPhone,
		ErrorNamespaceMerchantBankingCurrency:                 ErrorIncorrectCurrencyIdentifier,
		ErrorNamespaceMerchantBankingName:                     ErrorMessageIncorrectBankName,
		ErrorNamespaceMerchantBankingAddress:                  ErrorMessageIncorrectBankAddress,
		ErrorNamespaceMerchantBankingAccountNumber:            ErrorMessageIncorrectBankAccountNumber,
		ErrorNamespaceMerchantBankingSwift:                    ErrorMessageIncorrectBankSwift,
		ErrorNamespaceMerchantBankingCorrespondentAccount:     ErrorMessageIncorrectBankCorrespondentAccount,
		ErrorNamespaceGetDashboardMainRequestMerchantId:       ErrorIncorrectMerchantId,
		ErrorNamespaceGetDashboardMainRequestPeriod:           ErrorIncorrectPeriod,
		ErrorNamespaceGetDashboardBaseReportRequestPeriod:     ErrorIncorrectPeriod,
		ErrorNamespaceGetDashboardBaseReportRequestMerchantId: ErrorIncorrectMerchantId,
	}
)
