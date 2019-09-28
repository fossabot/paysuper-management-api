package common

import (
	"github.com/Nerufa/go-shared/logger"
	"github.com/paysuper/paysuper-billing-server/pkg"
	"regexp"
)

type Dictionary struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

const (
	// LimitDefault  = 100
	// LimitMax      = 1000
	// OffsetDefault = 0

	RequestParameterId                       = "id"
	RequestParameterName                     = "name"
	RequestParameterSku                      = "sku"
	RequestParameterIsSigned                 = "is_signed"
	RequestParameterMerchantId               = "merchant_id"
	RequestParameterProject                  = "project[]"
	RequestParameterPaymentMethod            = "payment_method[]"
	RequestParameterCountry                  = "country"
	RequestParameterCountries                = "country[]"
	RequestParameterProjectId                = "project_id"
	RequestParameterPaymentMethodId          = "method_id"
	RequestParameterOrderId                  = "order_id"
	RequestParameterRefundId                 = "refund_id"
	RequestParameterNotificationId           = "notification_id"
	RequestParameterUserId                   = "user"
	RequestParameterLimit                    = "limit"
	RequestParameterOffset                   = "offset"
	RequestParameterFile                     = "file"
	RequestParameterUtmSource                = "utm_source"
	RequestParameterUtmMedium                = "utm_medium"
	RequestParameterUtmCampaign              = "utm_campaign"
	RequestParameterIsSystem                 = "is_system"
	RequestParameterAgreementType            = "agreement_type"
	RequestParameterHasMerchantSignature     = "has_merchant_signature"
	RequestParameterHasPspSignature          = "has_psp_signature"
	RequestParameterAgreementSentViaMail     = "agreement_sent_via_mail"
	RequestParameterMailTrackingLink         = "mail_tracking_link"
	RequestParameterImage                    = "image"
	RequestParameterCallbackCurrency         = "callback_currency"
	RequestParameterCallbackProtocol         = "callback_protocol"
	RequestParameterCreateOrderAllowedUrls   = "create_order_allowed_urls"
	RequestParameterAllowDynamicNotifyUrls   = "allow_dynamic_notify_urls"
	RequestParameterAllowDynamicRedirectUrls = "allow_dynamic_redirect_urls"
	RequestParameterLimitsCurrency           = "limits_currency"
	RequestParameterMinPaymentAmount         = "min_payment_amount"
	RequestParameterMaxPaymentAmount         = "max_payment_amount"
	RequestParameterNotifyEmails             = "notify_emails"
	RequestParameterIsProductsCheckout       = "is_products_checkout"
	RequestParameterSecretKey                = "secret_key"
	RequestParameterSignatureRequired        = "signature_required"
	RequestParameterSendNotifyEmail          = "send_notify_email"
	RequestParameterUrlCheckAccount          = "url_check_account"
	RequestParameterUrlProcessPayment        = "url_process_payment"
	RequestParameterUrlRedirectFail          = "url_redirect_fail"
	RequestParameterUrlRedirectSuccess       = "url_redirect_success"
	RequestParameterUrlChargebackPayment     = "url_chargeback_payment"
	RequestParameterUrlCancelPayment         = "url_cancel_payment"
	RequestParameterUrlFraudPayment          = "url_fraud_payment"
	RequestParameterUrlRefundPayment         = "url_refund_payment"
	RequestParameterStatus                   = "status"
	RequestAuthorizationTokenRegex           = "Bearer ([A-z0-9_.-]{10,})"
	RequestParameterZipUsa                   = "zip_usa"
	RequestParameterRateId                   = "rate_id"

	UserProfileFieldNumberOfEmployees = "NumberOfEmployees"
	UserProfileFieldAnnualIncome      = "AnnualIncome"
	UserProfileFieldCompanyName       = "CompanyName"
	UserProfileFieldPosition          = "Position"
	UserProfileFieldFirstName         = "FirstName"
	UserProfileFieldLastName          = "LastName"
	UserProfileFieldWebsite           = "Website"
	UserProfileFieldKindOfActivity    = "KindOfActivity"
	UserProfileFieldReview            = "Review"
	UserProfileFieldPageId            = "PageId"

	OrderFieldProjectId     = "PP_PROJECT_ID"
	OrderFieldSignature     = "PP_SIGNATURE"
	OrderFieldAmount        = "PP_AMOUNT"
	OrderFieldCurrency      = "PP_CURRENCY"
	OrderFieldAccount       = "PP_ACCOUNT"
	OrderFieldOrderId       = "PP_ORDER_ID"
	OrderFieldPaymentMethod = "PP_PAYMENT_METHOD"
	OrderFieldUrlVerify     = "PP_URL_VERIFY"
	OrderFieldUrlNotify     = "PP_URL_NOTIFY"
	OrderFieldUrlSuccess    = "PP_URL_SUCCESS"
	OrderFieldUrlFail       = "PP_URL_FAIL"
	OrderFieldPayerEmail    = "PP_PAYER_EMAIL"
	OrderFieldPayerPhone    = "PP_PAYER_PHONE"
	OrderFieldDescription   = "PP_DESCRIPTION"
	OrderFieldRegion        = "PP_REGION"

	QueryParameterNameLimit  = "limit"
	QueryParameterNameOffset = "offset"
	QueryParameterNameSort   = "sort[]"

	ErrorMessageMask = "field validation for '%s' failed on the '%s' tag"

	HeaderAcceptLanguage      = "Accept-Language"
	HeaderUserAgent           = "User-Agent"
	HeaderXApiSignatureHeader = "X-API-SIGNATURE"
	HeaderReferer             = "referer"

	// EnvironmentProduction        = "prod"
	CustomerTokenCookiesName = "_ps_ctkn"
	// CustomerTokenCookiesLifetime = 2592000

	CardPayPaymentResponseHeaderSignature = "Signature"

	AgreementPageTemplateName = "agreement.html"

	UserProfilePositionCEO               = "CEO"
	UserProfilePositionCTO               = "CTO"
	UserProfilePositionCMO               = "CMO"
	UserProfilePositionCFO               = "CFO"
	UserProfilePositionProjectManagement = "Project Management"
	UserProfilePositionGenericManagement = "Generic Management"
	UserProfilePositionSoftwareDeveloper = "Software Developer"
	UserProfilePositionMarketing         = "Marketing"
	UserProfilePositionSupport           = "Support"

	ErrorFieldService = "service"
	ErrorFieldMethod  = "method"
	ErrorFieldRequest = "request"

	InternalErrorTemplate = "internal error"
)

var (
	DefaultSort = []string{"_id"}

	OrderReservedWords = map[string]bool{
		OrderFieldProjectId:     true,
		OrderFieldSignature:     true,
		OrderFieldAmount:        true,
		OrderFieldCurrency:      true,
		OrderFieldAccount:       true,
		OrderFieldOrderId:       true,
		OrderFieldDescription:   true,
		OrderFieldPaymentMethod: true,
		OrderFieldUrlVerify:     true,
		OrderFieldUrlNotify:     true,
		OrderFieldUrlSuccess:    true,
		OrderFieldUrlFail:       true,
		OrderFieldPayerEmail:    true,
		OrderFieldPayerPhone:    true,
		OrderFieldRegion:        true,
	}

	ZipRegexp = map[string]*regexp.Regexp{
		"AF": regexp.MustCompile("^\\d{4}$"),
		"AX": regexp.MustCompile("^\\d{5}$"),
		"AL": regexp.MustCompile("^\\d{4}$"),
		"DZ": regexp.MustCompile("^\\d{5}$"),
		"AS": regexp.MustCompile("^\\d{5}(-{1}\\d{4,6})$"),
		"AD": regexp.MustCompile("^[Aa][Dd]\\d{3}$"),
		"AI": regexp.MustCompile("^[Aa][I][-][2][6][4][0]$"),
		"AR": regexp.MustCompile("^\\d{4}|[A-Za-z]\\d{4}[a-zA-Z]{3}$"),
		"AM": regexp.MustCompile("^\\d{4}$"),
		"AC": regexp.MustCompile("^[Aa][Ss][Cc][Nn]\\s{0,1}[1][Zz][Zz]$"),
		"AU": regexp.MustCompile("^\\d{4}$"),
		"AT": regexp.MustCompile("^\\d{4}$"),
		"AZ": regexp.MustCompile("^[Aa][Zz]\\d{4}$"),
		"BH": regexp.MustCompile("^\\d{3,4}$"),
		"BD": regexp.MustCompile("^\\d{4}$"),
		"BB": regexp.MustCompile("^[Aa][Zz]\\d{5}$"),
		"BY": regexp.MustCompile("^\\d{6}$"),
		"BE": regexp.MustCompile("^\\d{4}$"),
		"BM": regexp.MustCompile("^[A-Za-z]{2}\\s([A-Za-z]{2}|\\d{2})$"),
		"BT": regexp.MustCompile("^\\d{5}$"),
		"BO": regexp.MustCompile("^\\d{4}$"),
		"BA": regexp.MustCompile("^\\d{5}$"),
		"BR": regexp.MustCompile("^\\d{5}-\\d{3}$"),
		"IO": regexp.MustCompile("^[Bb]{2}[Nn][Dd]\\s{0,1}[1][Zz]{2}$"),
		"VG": regexp.MustCompile("^[Vv][Gg]\\d{4}$"),
		"BN": regexp.MustCompile("^[A-Za-z]{2}\\d{4}$"),
		"BG": regexp.MustCompile("^\\d{4}$"),
		"KH": regexp.MustCompile("^\\d{5}$"),
		"CA": regexp.MustCompile("^[A-Za-z]\\d[A-Za-z][ -]?\\d[A-Za-z]\\d$"),
		"CV": regexp.MustCompile("^\\d{4}$"),
		"KY": regexp.MustCompile("^[Kk][Yy]\\d[-\\s]{0,1}\\d{4}$"),
		"TD": regexp.MustCompile("^\\d{5}$"),
		"CL": regexp.MustCompile("^\\d{7}\\s\\(\\d{3}-\\d{4}\\)$"),
		"CN": regexp.MustCompile("^\\d{6}$"),
		"CX": regexp.MustCompile("^\\d{4}$"),
		"CC": regexp.MustCompile("^\\d{4}$"),
		"CO": regexp.MustCompile("^\\d{6}$"),
		"CD": regexp.MustCompile("^[Cc][Dd]$"),
		"CR": regexp.MustCompile("^\\d{4,5}$"),
		"HR": regexp.MustCompile("^\\d{5}$"),
		"CU": regexp.MustCompile("^\\d{5}$"),
		"CY": regexp.MustCompile("^\\d{4}$"),
		"CZ": regexp.MustCompile("^\\d{5}\\s\\(\\d{3}\\s\\d{2}\\)$"),
		"DK": regexp.MustCompile("^\\d{4}$"),
		"DO": regexp.MustCompile("^\\d{5}$"),
		"EC": regexp.MustCompile("^\\d{6}$"),
		"SV": regexp.MustCompile("^1101$"),
		"EG": regexp.MustCompile("^\\d{5}$"),
		"EE": regexp.MustCompile("^\\d{5}$"),
		"ET": regexp.MustCompile("^\\d{4}$"),
		"FK": regexp.MustCompile("^[Ff][Ii][Qq]{2}\\s{0,1}[1][Zz]{2}$"),
		"FO": regexp.MustCompile("^\\d{3}$"),
		"FI": regexp.MustCompile("^\\d{5}$"),
		"FR": regexp.MustCompile("^\\d{5}$"),
		"GF": regexp.MustCompile("^973\\d{2}$"),
		"PF": regexp.MustCompile("^987\\d{2}$"),
		"GA": regexp.MustCompile("^\\d{2}\\s[a-zA-Z-_ ]\\s\\d{2}$"),
		"GE": regexp.MustCompile("^\\d{4}$"),
		"DE": regexp.MustCompile("^\\d{2,5}$"),
		"GI": regexp.MustCompile("^[Gg][Xx][1]{2}\\s{0,1}[1][Aa]{2}$"),
		"GR": regexp.MustCompile("^\\d{3}\\s{0,1}\\d{2}$"),
		"GL": regexp.MustCompile("^\\d{4}$"),
		"GP": regexp.MustCompile("^971\\d{2}$"),
		"GU": regexp.MustCompile("^\\d{5}$"),
		"GT": regexp.MustCompile("^\\d{5}$"),
		"GG": regexp.MustCompile("^[A-Za-z]{2}\\d\\s{0,1}\\d[A-Za-z]{2}$"),
		"GW": regexp.MustCompile("^\\d{4}$"),
		"HT": regexp.MustCompile("^\\d{4}$"),
		"HM": regexp.MustCompile("^\\d{4}$"),
		"HN": regexp.MustCompile("^\\d{5}$"),
		"HU": regexp.MustCompile("^\\d{4}$"),
		"IS": regexp.MustCompile("^\\d{3}$"),
		"IN": regexp.MustCompile("^\\d{6}$"),
		"ID": regexp.MustCompile("^\\d{5}$"),
		"IR": regexp.MustCompile("^\\d{5}-\\d{5}$"),
		"IQ": regexp.MustCompile("^\\d{5}$"),
		"IM": regexp.MustCompile("^[Ii[Mm]\\d{1,2}\\s\\d\\[A-Z]{2}$"),
		"IL": regexp.MustCompile("^\\b\\d{5}(\\d{2})?$"),
		"IT": regexp.MustCompile("^\\d{5}$"),
		"JM": regexp.MustCompile("^\\d{2}$"),
		"JP": regexp.MustCompile("^\\d{7}\\s\\(\\d{3}-\\d{4}\\)$"),
		"JE": regexp.MustCompile("^[Jj][Ee]\\d\\s{0,1}\\d[A-Za-z]{2}$"),
		"JO": regexp.MustCompile("^\\d{5}$"),
		"KZ": regexp.MustCompile("^\\d{6}$"),
		"KE": regexp.MustCompile("^\\d{5}$"),
		"KR": regexp.MustCompile("^\\d{6}\\s\\(\\d{3}-\\d{3}\\)$"),
		"XK": regexp.MustCompile("^\\d{5}$"),
		"KW": regexp.MustCompile("^\\d{5}$"),
		"KG": regexp.MustCompile("^\\d{6}$"),
		"LV": regexp.MustCompile("^[Ll][Vv][- ]{0,1}\\d{4}$"),
		"LA": regexp.MustCompile("^\\d{5}$"),
		"LB": regexp.MustCompile("^\\d{4}\\s{0,1}\\d{4}$"),
		"LS": regexp.MustCompile("^\\d{3}$"),
		"LR": regexp.MustCompile("^\\d{4}$"),
		"LY": regexp.MustCompile("^\\d{5}$"),
		"LI": regexp.MustCompile("^\\d{4}$"),
		"LT": regexp.MustCompile("^[Ll][Tt][- ]{0,1}\\d{5}$"),
		"LU": regexp.MustCompile("^\\d{4}$"),
		"MK": regexp.MustCompile("^\\d{4}$"),
		"MG": regexp.MustCompile("^\\d{3}$"),
		"MV": regexp.MustCompile("^\\d{4,5}$"),
		"MY": regexp.MustCompile("^\\d{5}$"),
		"MT": regexp.MustCompile("^[A-Za-z]{3}\\s{0,1}\\d{4}$"),
		"MH": regexp.MustCompile("^\\d{5}$"),
		"MQ": regexp.MustCompile("^972\\d{2}$"),
		"YT": regexp.MustCompile("^976\\d{2}$"),
		"MX": regexp.MustCompile("^\\d{5}$"),
		"FM": regexp.MustCompile("^\\d{5}$"),
		"MD": regexp.MustCompile("^[Mm][Dd][- ]{0,1}\\d{4}$"),
		"MC": regexp.MustCompile("^980\\d{2}$"),
		"MN": regexp.MustCompile("^\\d{5}$"),
		"ME": regexp.MustCompile("^\\d{5}$"),
		"MS": regexp.MustCompile("^[Mm][Ss][Rr]\\s{0,1}\\d{4}$"),
		"MA": regexp.MustCompile("^\\d{5}$"),
		"MZ": regexp.MustCompile("^\\d{4}$"),
		"MM": regexp.MustCompile("^\\d{5}$"),
		"NA": regexp.MustCompile("^\\d{5}$"),
		"NP": regexp.MustCompile("^\\d{5}$"),
		"NL": regexp.MustCompile("^\\d{4}\\s{0,1}[A-Za-z]{2}$"),
		"NC": regexp.MustCompile("^988\\d{2}$"),
		"NZ": regexp.MustCompile("^\\d{4}$"),
		"NI": regexp.MustCompile("^\\d{5}$"),
		"NE": regexp.MustCompile("^\\d{4}$"),
		"NG": regexp.MustCompile("^\\d{6}$"),
		"NF": regexp.MustCompile("^\\d{4}$"),
		"MP": regexp.MustCompile("^\\d{5}$"),
		"NO": regexp.MustCompile("^\\d{4}$"),
		"OM": regexp.MustCompile("^\\d{3}$"),
		"PK": regexp.MustCompile("^\\d{5}$"),
		"PW": regexp.MustCompile("^\\d{5}$"),
		"PA": regexp.MustCompile("^\\d{6}$"),
		"PG": regexp.MustCompile("^\\d{3}$"),
		"PY": regexp.MustCompile("^\\d{4}$"),
		"PE": regexp.MustCompile("^\\d{5}$"),
		"PH": regexp.MustCompile("^\\d{4}$"),
		"PN": regexp.MustCompile("^[Pp][Cc][Rr][Nn]\\s{0,1}[1][Zz]{2}$"),
		"PL": regexp.MustCompile("^\\d{2}[- ]{0,1}\\d{3}$"),
		"PT": regexp.MustCompile("^\\d{4}$"),
		"PR": regexp.MustCompile("^\\d{5}$"),
		"RE": regexp.MustCompile("^974\\d{2}$"),
		"RO": regexp.MustCompile("^\\d{6}$"),
		"RU": regexp.MustCompile("^\\d{6}$"),
		"BL": regexp.MustCompile("^97133$"),
		"SH": regexp.MustCompile("^[Ss][Tt][Hh][Ll]\\s{0,1}[1][Zz]{2}$|^[Tt][Dd][Cc][Uu]\\s{0,1}[1][Zz]{2}$"),
		"MF": regexp.MustCompile("^97150$"),
		"PM": regexp.MustCompile("^97500$"),
		"VC": regexp.MustCompile("^[Vv][Cc]\\d{4}$"),
		"SM": regexp.MustCompile("^4789\\d$"),
		"SA": regexp.MustCompile("^\\d{5}(-{1}\\d{4})?$"),
		"SN": regexp.MustCompile("^\\d{5}$"),
		"RS": regexp.MustCompile("^\\d{5}$"),
		"SG": regexp.MustCompile("^\\d{6}$"),
		"SK": regexp.MustCompile("^\\d{5}\\s\\(\\d{3}\\s\\d{2}\\)$"),
		"SI": regexp.MustCompile("^([Ss][Ii][- ]{0,1}){0,1}\\d{4}$"),
		"ZA": regexp.MustCompile("^\\d{4}$"),
		"GS": regexp.MustCompile("^[Ss][Ii][Qq]{2}\\s{0,1}[1][Zz]{2}$"),
		"ES": regexp.MustCompile("^\\d{5}$"),
		"LK": regexp.MustCompile("^\\d{5}$"),
		"SD": regexp.MustCompile("^\\d{5}$"),
		"SZ": regexp.MustCompile("^[A-Za-z]\\d{3}$"),
		"SE": regexp.MustCompile("^\\d{3}\\s*\\d{2}$"),
		"CH": regexp.MustCompile("^\\d{4}$"),
		"SJ": regexp.MustCompile("^\\d{4}$"),
		"TW": regexp.MustCompile("^\\d{5}$"),
		"TJ": regexp.MustCompile("^\\d{6}$"),
		"TH": regexp.MustCompile("^\\d{5}$"),
		"TT": regexp.MustCompile("^\\d{6}$"),
		"TN": regexp.MustCompile("^\\d{4}$"),
		"TR": regexp.MustCompile("^\\d{5}$"),
		"TM": regexp.MustCompile("^\\d{6}$"),
		"TC": regexp.MustCompile("^[Tt][Kk][Cc][Aa]\\s{0,1}[1][Zz]{2}$"),
		"UA": regexp.MustCompile("^\\d{5}$"),
		"GB": regexp.MustCompile("^[A-Z]{1,2}[0-9R][0-9A-Z]?\\s*[0-9][A-Z-[CIKMOV]]{2}"),
		"US": regexp.MustCompile("^\\b\\d{5}\\b(?:[- ]{1}\\d{4})?$"),
		"UY": regexp.MustCompile("^\\d{5}$"),
		"VI": regexp.MustCompile("^\\d{5}$"),
		"UZ": regexp.MustCompile("^\\d{3} \\d{3}$"),
		"VA": regexp.MustCompile("^120$"),
		"VE": regexp.MustCompile("^\\d{4}(\\s[a-zA-Z]{1})?$"),
		"VN": regexp.MustCompile("^\\d{6}$"),
		"WF": regexp.MustCompile("^986\\d{2}$"),
		"ZM": regexp.MustCompile("^\\d{5}$"),
	}

	TariffRegions = map[string]string{
		"cis":             "CIS",
		"russia":          "Russia",
		"west_asia":       "West Asia",
		"eu":              "EU",
		"north_america":   "North America",
		"central_america": "Central America",
		"south_america":   "South America",
		"uk":              "United Kingdom",
		"worldwide":       "Worldwide",
		"south_pacific":   "South Pacific",
	}

	TestStubImplementMe = "implement me!"

	TokenRegex = regexp.MustCompile(RequestAuthorizationTokenRegex)
)

func LogSrvCallFailedGRPC(log logger.Logger, err error, name, method string, req interface{}) {
	log.Error(pkg.ErrorGrpcServiceCallFailed,
		logger.PairArgs(
			ErrorFieldService, name,
			ErrorFieldMethod, method,
		),
		logger.WithPrettyFields(logger.Fields{"err": err, ErrorFieldRequest: req}),
	)
}
