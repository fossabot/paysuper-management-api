package common

import "time"

type Auth1 struct {
	Issuer       string `envconfig:"AUTH1_ISSUER" default:"https://dev-auth1.tst.protocol.one"`
	ClientId     string `envconfig:"AUTH1_CLIENTID" required:"true"`
	ClientSecret string `envconfig:"AUTH1_CLIENTSECRET" required:"true"`
	RedirectUrl  string `envconfig:"AUTH1_REDIRECTURL" required:"true"`
}

type Config struct {
	Auth1

	PaymentFormJsLibraryUrl string `envconfig:"PAYMENT_FORM_JS_LIBRARY_URL" required:"true"`
	WebsocketUrl            string `envconfig:"WEBSOCKET_URL" default:"wss://cf.tst.protocol.one/connection/websocket"`

	AwsAccessKeyIdAgreement     string `envconfig:"AWS_ACCESS_KEY_ID_AGREEMENT" required:"true"`
	AwsSecretAccessKeyAgreement string `envconfig:"AWS_SECRET_ACCESS_KEY_AGREEMENT" required:"true"`
	AwsRegionAgreement          string `envconfig:"AWS_REGION_AGREEMENT" default:"eu-west-1"`
	AwsBucketAgreement          string `envconfig:"AWS_BUCKET_AGREEMENT" required:"true"`

	AwsAccessKeyIdReporter     string `envconfig:"AWS_ACCESS_KEY_ID_REPORTER" required:"true"`
	AwsSecretAccessKeyReporter string `envconfig:"AWS_SECRET_ACCESS_KEY_REPORTER" required:"true"`
	AwsRegionReporter          string `envconfig:"AWS_REGION_REPORTER" default:"eu-west-1"`
	AwsBucketReporter          string `envconfig:"AWS_BUCKET_REPORTER" required:"true"`

	LimitDefault                 int32 `default:"100"`
	OffsetDefault                int32 `default:"0"`
	LimitMax                     int32 `default:"1000"`
	ReturnPaymentForm            bool  `envconfig:"DEBUG_RETURN_PAYMENT_FORM"`
	DisableAuthMiddleware        bool
	CustomerTokenCookiesLifetime time.Duration // CustomerTokenCookiesLifetime = 2592000
}
