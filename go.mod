module github.com/paysuper/paysuper-management-api

require (
	github.com/ProtocolONE/authone-jwt-verifier-golang v0.0.0-20190327070329-4dd563b01681
	github.com/ProtocolONE/geoip-service v0.0.0-20190903084234-1d5ae6b96679
	github.com/ProtocolONE/go-core v1.10.3
	github.com/SebastiaanKlippert/go-wkhtmltopdf v1.4.1 // indirect
	github.com/aws/aws-sdk-go v1.23.16
	github.com/fatih/color v1.7.0
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8
	github.com/go-log/log v0.1.0
	github.com/golang/protobuf v1.3.2
	github.com/google/uuid v1.1.1
	github.com/google/wire v0.3.0
	github.com/karlseguin/expect v1.0.1 // indirect
	github.com/labstack/echo/v4 v4.1.6
	github.com/micro/go-micro v1.8.0
	github.com/micro/go-plugins v1.2.0
	github.com/paysuper/paysuper-aws-manager v0.0.1
	github.com/paysuper/paysuper-billing-server v0.0.0-20191003093339-de34a11f5870
	github.com/paysuper/paysuper-payment-link v0.0.0-20190903143854-b799a77c03ce
	github.com/paysuper/paysuper-recurring-repository v1.0.123
	github.com/paysuper/paysuper-reporter v0.0.0-20191003072342-610371fc9395
	github.com/paysuper/paysuper-tax-service v0.0.0-20190903084038-7849f394f122
	github.com/pkg/errors v0.8.1
	github.com/spf13/cobra v0.0.5
	github.com/stretchr/testify v1.4.0
	github.com/ttacon/libphonenumber v1.0.1
	github.com/wsxiaoys/terminal v0.0.0-20160513160801-0940f3fc43a0 // indirect
	go.uber.org/automaxprocs v1.2.0
	gopkg.in/go-playground/validator.v9 v9.29.1
	gopkg.in/karlseguin/expect.v1 v1.0.1 // indirect
)

replace (
	github.com/gogo/protobuf v0.0.0-20190410021324-65acae22fc9 => github.com/gogo/protobuf v1.2.2-0.20190723190241-65acae22fc9d
	github.com/hashicorp/consul => github.com/hashicorp/consul v1.5.1
)

go 1.12
