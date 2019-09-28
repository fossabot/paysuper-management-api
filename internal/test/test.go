package test

import (
	"bytes"
	"context"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo/v4"
	"github.com/paysuper/paysuper-management-api/internal/dispatcher/common"
	httpEcho "github.com/paysuper/paysuper-management-api/pkg/http"
	"github.com/stretchr/testify/assert"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)


var (
	HexId = bson.NewObjectId().Hex()
)

// EchoReqResCaller
type EchoReqResCaller struct {
	dispatcher      httpEcho.Dispatcher
	middlewareSetUp *MiddlewareTestUp
}

// NewTestRequest
func NewTestRequest(dispatcher httpEcho.Dispatcher, mw *MiddlewareTestUp) *EchoReqResCaller {
	return &EchoReqResCaller{
		dispatcher:      dispatcher,
		middlewareSetUp: mw,
	}
}

type Middleware interface {
	Post(middleware ...echo.MiddlewareFunc)
	Pre(middleware ...echo.MiddlewareFunc)
}

type MiddlewareTestUp struct {
	preMiddleware []echo.MiddlewareFunc
	useMiddleware []echo.MiddlewareFunc
}

type QueryBuilder struct {
	c            *EchoReqResCaller
	method, path string
	init         []func(request *http.Request, middleware Middleware)
	params       []string
	query        url.Values
	body         io.Reader
	cookie       []*http.Cookie
}

// Params
func (s *QueryBuilder) Params(params ...string) *QueryBuilder {
	s.params = append(s.params, params...)
	return s
}

// SetCookie
func (s *QueryBuilder) AddCookie(cookie *http.Cookie) *QueryBuilder {
	s.cookie = append(s.cookie, cookie)
	return s
}

// SetQueryParam
func (s *QueryBuilder) SetQueryParam(key, value string) *QueryBuilder {
	s.query.Set(key, value)
	return s
}

// SetQueryParams
func (s *QueryBuilder) SetQueryParams(values url.Values) *QueryBuilder {
	for key, value := range values {
		if len(value) != 0 {
			s.query.Set(key, value[0])
		}
	}
	return s
}

// Init
func (s *QueryBuilder) Init(init func(*http.Request, Middleware)) *QueryBuilder {
	s.init = append(s.init, init)
	return s
}

// BodyString
func (s *QueryBuilder) BodyString(body string) *QueryBuilder {
	s.body = strings.NewReader(body)
	return s
}

// BodyBytes
func (s *QueryBuilder) BodyBytes(body []byte) *QueryBuilder {
	s.body = bytes.NewReader(body)
	return s
}

// Method
func (s *QueryBuilder) Method(method string) *QueryBuilder {
	s.method = method
	return s
}

// Path
func (s *QueryBuilder) Path(path string) *QueryBuilder {
	s.path = path
	return s
}

// Body
func (s *QueryBuilder) Body(body io.Reader) *QueryBuilder {
	s.body = body
	return s
}

// Exec
func (s *QueryBuilder) Exec(t *testing.T) (*httptest.ResponseRecorder, error) {

	r := strings.NewReplacer(s.params...)
	s.path = r.Replace(s.path)

	u, err := url.Parse(s.path)
	assert.NoError(t, err)

	u.RawQuery = s.query.Encode()

	return s.c.Request(s.method, u.String(), s.body, func(request *http.Request, middleware Middleware) {
		for _, cookie := range s.cookie {
			request.AddCookie(cookie)
		}
		for _, init := range s.init {
			init(request, middleware)
		}
	})
}

// ExecFileUpload
func (s *QueryBuilder) ExecFileUpload(t *testing.T, params map[string]string, paramName, path string) (*httptest.ResponseRecorder, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := file.Close(); err != nil {
			return
		}
	}()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))

	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	err = writer.Close()

	if err != nil {
		return nil, err
	}

	s.Init(func(request *http.Request, middleware Middleware) {
		request.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	})

	s.Method(http.MethodPost)
	s.Body(body)

	return s.Exec(t)
}

// NewQueryBuilder
func NewQueryBuilder(c *EchoReqResCaller) *QueryBuilder {
	return &QueryBuilder{c: c, query: make(url.Values)}
}

// Post
func (m *MiddlewareTestUp) Post(middleware ...echo.MiddlewareFunc) {
	m.useMiddleware = middleware
}

// Pre
func (m *MiddlewareTestUp) Pre(middleware ...echo.MiddlewareFunc) {
	m.preMiddleware = middleware
}

// ListUse
func (m *MiddlewareTestUp) ListUse() []echo.MiddlewareFunc {
	return m.useMiddleware
}

// ListPre
func (m *MiddlewareTestUp) ListPre() []echo.MiddlewareFunc {
	return m.preMiddleware
}

type MiddlewareHandlerFunc func(echo.Context) (echo.Context, error)

// Request
func (c *EchoReqResCaller) Request(method, target string, body io.Reader, init func(*http.Request, Middleware)) (resRec *httptest.ResponseRecorder, err error) {
	e := echo.New()
	req := httptest.NewRequest(method, target, body)
	if init == nil {
		panic("request init function should be present")
	}
	init(req, c.middlewareSetUp)
	e.Pre(c.middlewareSetUp.ListPre()...)
	e.Use(c.middlewareSetUp.ListUse()...)
	resRec = httptest.NewRecorder()
	if err = c.dispatcher.Dispatch(e); err != nil {
		return
	}
	//
	e.HTTPErrorHandler = func(e error, context echo.Context) {
		err = e
	}
	//
	e.ServeHTTP(resRec, req)
	return
}

// Send
func (c *EchoReqResCaller) Builder() *QueryBuilder {
	return NewQueryBuilder(c)
}

// PreAuthUserMiddleware
func PreAuthUserMiddleware(user *common.AuthUser) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			common.SetUserContext(c, user)
			return next(c)
		}
	}
}

// DefaultSettings
func DefaultSettings() map[string]interface{} {
	return map[string]interface{}{
		"dispatcher": map[string]interface{}{
			"global": map[string]interface{}{
				"returnPaymentForm":            true,
				"disableAuthMiddleware":        true,
				"paymentFormJsLibraryUrl":      "unknown",
				"clientId":                     "unknown",
				"clientSecret":                 "unknown",
				"redirectUrl":                  "unknown",
				"awsAccessKeyIdAgreement":      "unknown",
				"awsSecretAccessKeyAgreement":  "unknown",
				"awsBucketAgreement":           "unknown",
				"awsAccessKeyIdReporter":       "unknown",
				"awsSecretAccessKeyReporter":   "unknown",
				"awsBucketReporter":            "unknown",
				"awsRegionAgreement":           "eu-west-1",
				"awsRegionReporter":            "eu-west-1",
				"awsBucketReporterr":           "eu-west-1",
				"customerTokenCookiesLifetime": "2592000s",
			},
		},
	}
}

// SetUp
func SetUp(settings map[string]interface{}, services common.Services, setUp func(*TestSet, Middleware) common.Handlers) (*EchoReqResCaller, error) {
	middlewareSetUp := &MiddlewareTestUp{}
	testSet, _, e := BuildTestSet(
		context.Background(),
		settings,
		services,
		nil,
	)
	d, _, e := BuildDispatcher(
		context.Background(),
		settings,
		services,
		setUp(testSet, middlewareSetUp),
		nil,
	)
	if e != nil {
		return nil, e
	}
	return NewTestRequest(d, middlewareSetUp), e
}

// ReqInitJSON
func ReqInitJSON() func(request *http.Request, middleware Middleware) {
	return func(request *http.Request, middleware Middleware) {
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}
}

// ReqInitMultipartForm
func ReqInitMultipartForm() func(request *http.Request, middleware Middleware) {
	return func(request *http.Request, middleware Middleware) {
		request.Header.Set(echo.HeaderContentType, echo.MIMEMultipartForm)
	}
}

// ReqInitApplicationForm
func ReqInitApplicationForm() func(request *http.Request, middleware Middleware) {
	return func(request *http.Request, middleware Middleware) {
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	}
}

// ReqInitXML
func ReqInitXML() func(request *http.Request, middleware Middleware) {
	return func(request *http.Request, middleware Middleware) {
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationXML)
	}
}
