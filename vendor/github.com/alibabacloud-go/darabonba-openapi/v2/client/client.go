// This file is auto-generated, don't edit it. Thanks.
// Description:
//
// This is for OpenApi SDK
package client

import (
	"io"

	spi "github.com/alibabacloud-go/alibabacloud-gateway-spi/client"
	openapiutil "github.com/alibabacloud-go/openapi-util/service"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	xml "github.com/alibabacloud-go/tea-xml/service"
	"github.com/alibabacloud-go/tea/tea"
	credential "github.com/aliyun/credentials-go/credentials"
)

type GlobalParameters struct {
	Headers map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	Queries map[string]*string `json:"queries,omitempty" xml:"queries,omitempty"`
}

func (s GlobalParameters) String() string {
	return tea.Prettify(s)
}

func (s GlobalParameters) GoString() string {
	return s.String()
}

func (s *GlobalParameters) SetHeaders(v map[string]*string) *GlobalParameters {
	s.Headers = v
	return s
}

func (s *GlobalParameters) SetQueries(v map[string]*string) *GlobalParameters {
	s.Queries = v
	return s
}

// Description:
//
// Model for initing client
type Config struct {
	// accesskey id
	AccessKeyId *string `json:"accessKeyId,omitempty" xml:"accessKeyId,omitempty"`
	// accesskey secret
	AccessKeySecret *string `json:"accessKeySecret,omitempty" xml:"accessKeySecret,omitempty"`
	// security token
	//
	// example:
	//
	// a.txt
	SecurityToken *string `json:"securityToken,omitempty" xml:"securityToken,omitempty"`
	// bearer token
	//
	// example:
	//
	// the-bearer-token
	BearerToken *string `json:"bearerToken,omitempty" xml:"bearerToken,omitempty"`
	// http protocol
	//
	// example:
	//
	// http
	Protocol *string `json:"protocol,omitempty" xml:"protocol,omitempty"`
	// http method
	//
	// example:
	//
	// GET
	Method *string `json:"method,omitempty" xml:"method,omitempty"`
	// region id
	//
	// example:
	//
	// cn-hangzhou
	RegionId *string `json:"regionId,omitempty" xml:"regionId,omitempty"`
	// read timeout
	//
	// example:
	//
	// 10
	ReadTimeout *int `json:"readTimeout,omitempty" xml:"readTimeout,omitempty"`
	// connect timeout
	//
	// example:
	//
	// 10
	ConnectTimeout *int `json:"connectTimeout,omitempty" xml:"connectTimeout,omitempty"`
	// http proxy
	//
	// example:
	//
	// http://localhost
	HttpProxy *string `json:"httpProxy,omitempty" xml:"httpProxy,omitempty"`
	// https proxy
	//
	// example:
	//
	// https://localhost
	HttpsProxy *string `json:"httpsProxy,omitempty" xml:"httpsProxy,omitempty"`
	// credential
	Credential credential.Credential `json:"credential,omitempty" xml:"credential,omitempty"`
	// endpoint
	//
	// example:
	//
	// cs.aliyuncs.com
	Endpoint *string `json:"endpoint,omitempty" xml:"endpoint,omitempty"`
	// proxy white list
	//
	// example:
	//
	// http://localhost
	NoProxy *string `json:"noProxy,omitempty" xml:"noProxy,omitempty"`
	// max idle conns
	//
	// example:
	//
	// 3
	MaxIdleConns *int `json:"maxIdleConns,omitempty" xml:"maxIdleConns,omitempty"`
	// network for endpoint
	//
	// example:
	//
	// public
	Network *string `json:"network,omitempty" xml:"network,omitempty"`
	// user agent
	//
	// example:
	//
	// Alibabacloud/1
	UserAgent *string `json:"userAgent,omitempty" xml:"userAgent,omitempty"`
	// suffix for endpoint
	//
	// example:
	//
	// aliyun
	Suffix *string `json:"suffix,omitempty" xml:"suffix,omitempty"`
	// socks5 proxy
	Socks5Proxy *string `json:"socks5Proxy,omitempty" xml:"socks5Proxy,omitempty"`
	// socks5 network
	//
	// example:
	//
	// TCP
	Socks5NetWork *string `json:"socks5NetWork,omitempty" xml:"socks5NetWork,omitempty"`
	// endpoint type
	//
	// example:
	//
	// internal
	EndpointType *string `json:"endpointType,omitempty" xml:"endpointType,omitempty"`
	// OpenPlatform endpoint
	//
	// example:
	//
	// openplatform.aliyuncs.com
	OpenPlatformEndpoint *string `json:"openPlatformEndpoint,omitempty" xml:"openPlatformEndpoint,omitempty"`
	// Deprecated
	//
	// credential type
	//
	// example:
	//
	// access_key
	Type *string `json:"type,omitempty" xml:"type,omitempty"`
	// Signature Version
	//
	// example:
	//
	// v1
	SignatureVersion *string `json:"signatureVersion,omitempty" xml:"signatureVersion,omitempty"`
	// Signature Algorithm
	//
	// example:
	//
	// ACS3-HMAC-SHA256
	SignatureAlgorithm *string `json:"signatureAlgorithm,omitempty" xml:"signatureAlgorithm,omitempty"`
	// Global Parameters
	GlobalParameters *GlobalParameters `json:"globalParameters,omitempty" xml:"globalParameters,omitempty"`
	// privite key for client certificate
	//
	// example:
	//
	// MIIEvQ
	Key *string `json:"key,omitempty" xml:"key,omitempty"`
	// client certificate
	//
	// example:
	//
	// -----BEGIN CERTIFICATE-----
	//
	// xxx-----END CERTIFICATE-----
	Cert *string `json:"cert,omitempty" xml:"cert,omitempty"`
	// server certificate
	//
	// example:
	//
	// -----BEGIN CERTIFICATE-----
	//
	// xxx-----END CERTIFICATE-----
	Ca *string `json:"ca,omitempty" xml:"ca,omitempty"`
	// disable HTTP/2
	//
	// example:
	//
	// false
	DisableHttp2 *bool `json:"disableHttp2,omitempty" xml:"disableHttp2,omitempty"`
}

func (s Config) String() string {
	return tea.Prettify(s)
}

func (s Config) GoString() string {
	return s.String()
}

func (s *Config) SetAccessKeyId(v string) *Config {
	s.AccessKeyId = &v
	return s
}

func (s *Config) SetAccessKeySecret(v string) *Config {
	s.AccessKeySecret = &v
	return s
}

func (s *Config) SetSecurityToken(v string) *Config {
	s.SecurityToken = &v
	return s
}

func (s *Config) SetBearerToken(v string) *Config {
	s.BearerToken = &v
	return s
}

func (s *Config) SetProtocol(v string) *Config {
	s.Protocol = &v
	return s
}

func (s *Config) SetMethod(v string) *Config {
	s.Method = &v
	return s
}

func (s *Config) SetRegionId(v string) *Config {
	s.RegionId = &v
	return s
}

func (s *Config) SetReadTimeout(v int) *Config {
	s.ReadTimeout = &v
	return s
}

func (s *Config) SetConnectTimeout(v int) *Config {
	s.ConnectTimeout = &v
	return s
}

func (s *Config) SetHttpProxy(v string) *Config {
	s.HttpProxy = &v
	return s
}

func (s *Config) SetHttpsProxy(v string) *Config {
	s.HttpsProxy = &v
	return s
}

func (s *Config) SetCredential(v credential.Credential) *Config {
	s.Credential = v
	return s
}

func (s *Config) SetEndpoint(v string) *Config {
	s.Endpoint = &v
	return s
}

func (s *Config) SetNoProxy(v string) *Config {
	s.NoProxy = &v
	return s
}

func (s *Config) SetMaxIdleConns(v int) *Config {
	s.MaxIdleConns = &v
	return s
}

func (s *Config) SetNetwork(v string) *Config {
	s.Network = &v
	return s
}

func (s *Config) SetUserAgent(v string) *Config {
	s.UserAgent = &v
	return s
}

func (s *Config) SetSuffix(v string) *Config {
	s.Suffix = &v
	return s
}

func (s *Config) SetSocks5Proxy(v string) *Config {
	s.Socks5Proxy = &v
	return s
}

func (s *Config) SetSocks5NetWork(v string) *Config {
	s.Socks5NetWork = &v
	return s
}

func (s *Config) SetEndpointType(v string) *Config {
	s.EndpointType = &v
	return s
}

func (s *Config) SetOpenPlatformEndpoint(v string) *Config {
	s.OpenPlatformEndpoint = &v
	return s
}

func (s *Config) SetType(v string) *Config {
	s.Type = &v
	return s
}

func (s *Config) SetSignatureVersion(v string) *Config {
	s.SignatureVersion = &v
	return s
}

func (s *Config) SetSignatureAlgorithm(v string) *Config {
	s.SignatureAlgorithm = &v
	return s
}

func (s *Config) SetGlobalParameters(v *GlobalParameters) *Config {
	s.GlobalParameters = v
	return s
}

func (s *Config) SetKey(v string) *Config {
	s.Key = &v
	return s
}

func (s *Config) SetCert(v string) *Config {
	s.Cert = &v
	return s
}

func (s *Config) SetCa(v string) *Config {
	s.Ca = &v
	return s
}

func (s *Config) SetDisableHttp2(v bool) *Config {
	s.DisableHttp2 = &v
	return s
}

type OpenApiRequest struct {
	Headers          map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	Query            map[string]*string `json:"query,omitempty" xml:"query,omitempty"`
	Body             interface{}        `json:"body,omitempty" xml:"body,omitempty"`
	Stream           io.Reader          `json:"stream,omitempty" xml:"stream,omitempty"`
	HostMap          map[string]*string `json:"hostMap,omitempty" xml:"hostMap,omitempty"`
	EndpointOverride *string            `json:"endpointOverride,omitempty" xml:"endpointOverride,omitempty"`
}

func (s OpenApiRequest) String() string {
	return tea.Prettify(s)
}

func (s OpenApiRequest) GoString() string {
	return s.String()
}

func (s *OpenApiRequest) SetHeaders(v map[string]*string) *OpenApiRequest {
	s.Headers = v
	return s
}

func (s *OpenApiRequest) SetQuery(v map[string]*string) *OpenApiRequest {
	s.Query = v
	return s
}

func (s *OpenApiRequest) SetBody(v interface{}) *OpenApiRequest {
	s.Body = v
	return s
}

func (s *OpenApiRequest) SetStream(v io.Reader) *OpenApiRequest {
	s.Stream = v
	return s
}

func (s *OpenApiRequest) SetHostMap(v map[string]*string) *OpenApiRequest {
	s.HostMap = v
	return s
}

func (s *OpenApiRequest) SetEndpointOverride(v string) *OpenApiRequest {
	s.EndpointOverride = &v
	return s
}

type Params struct {
	Action      *string `json:"action,omitempty" xml:"action,omitempty" require:"true"`
	Version     *string `json:"version,omitempty" xml:"version,omitempty" require:"true"`
	Protocol    *string `json:"protocol,omitempty" xml:"protocol,omitempty" require:"true"`
	Pathname    *string `json:"pathname,omitempty" xml:"pathname,omitempty" require:"true"`
	Method      *string `json:"method,omitempty" xml:"method,omitempty" require:"true"`
	AuthType    *string `json:"authType,omitempty" xml:"authType,omitempty" require:"true"`
	BodyType    *string `json:"bodyType,omitempty" xml:"bodyType,omitempty" require:"true"`
	ReqBodyType *string `json:"reqBodyType,omitempty" xml:"reqBodyType,omitempty" require:"true"`
	Style       *string `json:"style,omitempty" xml:"style,omitempty"`
}

func (s Params) String() string {
	return tea.Prettify(s)
}

func (s Params) GoString() string {
	return s.String()
}

func (s *Params) SetAction(v string) *Params {
	s.Action = &v
	return s
}

func (s *Params) SetVersion(v string) *Params {
	s.Version = &v
	return s
}

func (s *Params) SetProtocol(v string) *Params {
	s.Protocol = &v
	return s
}

func (s *Params) SetPathname(v string) *Params {
	s.Pathname = &v
	return s
}

func (s *Params) SetMethod(v string) *Params {
	s.Method = &v
	return s
}

func (s *Params) SetAuthType(v string) *Params {
	s.AuthType = &v
	return s
}

func (s *Params) SetBodyType(v string) *Params {
	s.BodyType = &v
	return s
}

func (s *Params) SetReqBodyType(v string) *Params {
	s.ReqBodyType = &v
	return s
}

func (s *Params) SetStyle(v string) *Params {
	s.Style = &v
	return s
}

type Client struct {
	Endpoint             *string
	RegionId             *string
	Protocol             *string
	Method               *string
	UserAgent            *string
	EndpointRule         *string
	EndpointMap          map[string]*string
	Suffix               *string
	ReadTimeout          *int
	ConnectTimeout       *int
	HttpProxy            *string
	HttpsProxy           *string
	Socks5Proxy          *string
	Socks5NetWork        *string
	NoProxy              *string
	Network              *string
	ProductId            *string
	MaxIdleConns         *int
	EndpointType         *string
	OpenPlatformEndpoint *string
	Credential           credential.Credential
	SignatureVersion     *string
	SignatureAlgorithm   *string
	Headers              map[string]*string
	Spi                  spi.ClientInterface
	GlobalParameters     *GlobalParameters
	Key                  *string
	Cert                 *string
	Ca                   *string
	DisableHttp2         *bool
}

// Description:
//
// # Init client with Config
//
// @param config - config contains the necessary information to create a client
func NewClient(config *Config) (*Client, error) {
	client := new(Client)
	err := client.Init(config)
	return client, err
}

func (client *Client) Init(config *Config) (_err error) {
	if tea.BoolValue(util.IsUnset(config)) {
		_err = tea.NewSDKError(map[string]interface{}{
			"code":    "ParameterMissing",
			"message": "'config' can not be unset",
		})
		return _err
	}

	if !tea.BoolValue(util.Empty(config.AccessKeyId)) && !tea.BoolValue(util.Empty(config.AccessKeySecret)) {
		if !tea.BoolValue(util.Empty(config.SecurityToken)) {
			config.Type = tea.String("sts")
		} else {
			config.Type = tea.String("access_key")
		}

		credentialConfig := &credential.Config{
			AccessKeyId:     config.AccessKeyId,
			Type:            config.Type,
			AccessKeySecret: config.AccessKeySecret,
		}
		credentialConfig.SecurityToken = config.SecurityToken
		client.Credential, _err = credential.NewCredential(credentialConfig)
		if _err != nil {
			return _err
		}

	} else if !tea.BoolValue(util.Empty(config.BearerToken)) {
		cc := &credential.Config{
			Type:        tea.String("bearer"),
			BearerToken: config.BearerToken,
		}
		client.Credential, _err = credential.NewCredential(cc)
		if _err != nil {
			return _err
		}

	} else if !tea.BoolValue(util.IsUnset(config.Credential)) {
		client.Credential = config.Credential
	}

	client.Endpoint = config.Endpoint
	client.EndpointType = config.EndpointType
	client.Network = config.Network
	client.Suffix = config.Suffix
	client.Protocol = config.Protocol
	client.Method = config.Method
	client.RegionId = config.RegionId
	client.UserAgent = config.UserAgent
	client.ReadTimeout = config.ReadTimeout
	client.ConnectTimeout = config.ConnectTimeout
	client.HttpProxy = config.HttpProxy
	client.HttpsProxy = config.HttpsProxy
	client.NoProxy = config.NoProxy
	client.Socks5Proxy = config.Socks5Proxy
	client.Socks5NetWork = config.Socks5NetWork
	client.MaxIdleConns = config.MaxIdleConns
	client.SignatureVersion = config.SignatureVersion
	client.SignatureAlgorithm = config.SignatureAlgorithm
	client.GlobalParameters = config.GlobalParameters
	client.Key = config.Key
	client.Cert = config.Cert
	client.Ca = config.Ca
	client.DisableHttp2 = config.DisableHttp2
	return nil
}

// Description:
//
// # Encapsulate the request and invoke the network
//
// @param action - api name
//
// @param version - product version
//
// @param protocol - http or https
//
// @param method - e.g. GET
//
// @param authType - authorization type e.g. AK
//
// @param bodyType - response body type e.g. String
//
// @param request - object of OpenApiRequest
//
// @param runtime - which controls some details of call api, such as retry times
//
// @return the response
func (client *Client) DoRPCRequest(action *string, version *string, protocol *string, method *string, authType *string, bodyType *string, request *OpenApiRequest, runtime *util.RuntimeOptions) (_result map[string]interface{}, _err error) {
	_err = tea.Validate(request)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Validate(runtime)
	if _err != nil {
		return _result, _err
	}
	_runtime := map[string]interface{}{
		"timeouted":      "retry",
		"key":            tea.StringValue(util.DefaultString(runtime.Key, client.Key)),
		"cert":           tea.StringValue(util.DefaultString(runtime.Cert, client.Cert)),
		"ca":             tea.StringValue(util.DefaultString(runtime.Ca, client.Ca)),
		"readTimeout":    tea.IntValue(util.DefaultNumber(runtime.ReadTimeout, client.ReadTimeout)),
		"connectTimeout": tea.IntValue(util.DefaultNumber(runtime.ConnectTimeout, client.ConnectTimeout)),
		"httpProxy":      tea.StringValue(util.DefaultString(runtime.HttpProxy, client.HttpProxy)),
		"httpsProxy":     tea.StringValue(util.DefaultString(runtime.HttpsProxy, client.HttpsProxy)),
		"noProxy":        tea.StringValue(util.DefaultString(runtime.NoProxy, client.NoProxy)),
		"socks5Proxy":    tea.StringValue(util.DefaultString(runtime.Socks5Proxy, client.Socks5Proxy)),
		"socks5NetWork":  tea.StringValue(util.DefaultString(runtime.Socks5NetWork, client.Socks5NetWork)),
		"maxIdleConns":   tea.IntValue(util.DefaultNumber(runtime.MaxIdleConns, client.MaxIdleConns)),
		"retry": map[string]interface{}{
			"retryable":   tea.BoolValue(runtime.Autoretry),
			"maxAttempts": tea.IntValue(util.DefaultNumber(runtime.MaxAttempts, tea.Int(3))),
		},
		"backoff": map[string]interface{}{
			"policy": tea.StringValue(util.DefaultString(runtime.BackoffPolicy, tea.String("no"))),
			"period": tea.IntValue(util.DefaultNumber(runtime.BackoffPeriod, tea.Int(1))),
		},
		"ignoreSSL": tea.BoolValue(runtime.IgnoreSSL),
	}

	_resp := make(map[string]interface{})
	for _retryTimes := 0; tea.BoolValue(tea.AllowRetry(_runtime["retry"], tea.Int(_retryTimes))); _retryTimes++ {
		if _retryTimes > 0 {
			_backoffTime := tea.GetBackoffTime(_runtime["backoff"], tea.Int(_retryTimes))
			if tea.IntValue(_backoffTime) > 0 {
				tea.Sleep(_backoffTime)
			}
		}

		_resp, _err = func() (map[string]interface{}, error) {
			request_ := tea.NewRequest()
			request_.Protocol = util.DefaultString(client.Protocol, protocol)
			request_.Method = method
			request_.Pathname = tea.String("/")
			globalQueries := make(map[string]*string)
			globalHeaders := make(map[string]*string)
			if !tea.BoolValue(util.IsUnset(client.GlobalParameters)) {
				globalParams := client.GlobalParameters
				if !tea.BoolValue(util.IsUnset(globalParams.Queries)) {
					globalQueries = globalParams.Queries
				}

				if !tea.BoolValue(util.IsUnset(globalParams.Headers)) {
					globalHeaders = globalParams.Headers
				}

			}

			extendsHeaders := make(map[string]*string)
			extendsQueries := make(map[string]*string)
			if !tea.BoolValue(util.IsUnset(runtime.ExtendsParameters)) {
				extendsParameters := runtime.ExtendsParameters
				if !tea.BoolValue(util.IsUnset(extendsParameters.Headers)) {
					extendsHeaders = extendsParameters.Headers
				}

				if !tea.BoolValue(util.IsUnset(extendsParameters.Queries)) {
					extendsQueries = extendsParameters.Queries
				}

			}

			request_.Query = tea.Merge(map[string]*string{
				"Action":         action,
				"Format":         tea.String("json"),
				"Version":        version,
				"Timestamp":      openapiutil.GetTimestamp(),
				"SignatureNonce": util.GetNonce(),
			}, globalQueries,
				extendsQueries,
				request.Query)
			headers, _err := client.GetRpcHeaders()
			if _err != nil {
				return _result, _err
			}

			if tea.BoolValue(util.IsUnset(headers)) {
				// endpoint is setted in product client
				request_.Headers = tea.Merge(map[string]*string{
					"host":          client.Endpoint,
					"x-acs-version": version,
					"x-acs-action":  action,
					"user-agent":    client.GetUserAgent(),
				}, globalHeaders,
					extendsHeaders)
			} else {
				request_.Headers = tea.Merge(map[string]*string{
					"host":          client.Endpoint,
					"x-acs-version": version,
					"x-acs-action":  action,
					"user-agent":    client.GetUserAgent(),
				}, globalHeaders,
					extendsHeaders,
					headers)
			}

			if !tea.BoolValue(util.IsUnset(request.Body)) {
				m, _err := util.AssertAsMap(request.Body)
				if _err != nil {
					return _result, _err
				}

				tmp := util.AnyifyMapValue(openapiutil.Query(m))
				request_.Body = tea.ToReader(util.ToFormString(tmp))
				request_.Headers["content-type"] = tea.String("application/x-www-form-urlencoded")
			}

			if !tea.BoolValue(util.EqualString(authType, tea.String("Anonymous"))) {
				credentialModel, _err := client.Credential.GetCredential()
				if _err != nil {
					return _result, _err
				}

				credentialType := credentialModel.Type
				if tea.BoolValue(util.EqualString(credentialType, tea.String("bearer"))) {
					bearerToken := credentialModel.BearerToken
					request_.Query["BearerToken"] = bearerToken
					request_.Query["SignatureType"] = tea.String("BEARERTOKEN")
				} else {
					accessKeyId := credentialModel.AccessKeyId
					accessKeySecret := credentialModel.AccessKeySecret
					securityToken := credentialModel.SecurityToken
					if !tea.BoolValue(util.Empty(securityToken)) {
						request_.Query["SecurityToken"] = securityToken
					}

					request_.Query["SignatureMethod"] = tea.String("HMAC-SHA1")
					request_.Query["SignatureVersion"] = tea.String("1.0")
					request_.Query["AccessKeyId"] = accessKeyId
					var t map[string]interface{}
					if !tea.BoolValue(util.IsUnset(request.Body)) {
						t, _err = util.AssertAsMap(request.Body)
						if _err != nil {
							return _result, _err
						}

					}

					signedParam := tea.Merge(request_.Query,
						openapiutil.Query(t))
					request_.Query["Signature"] = openapiutil.GetRPCSignature(signedParam, request_.Method, accessKeySecret)
				}

			}

			response_, _err := tea.DoRequest(request_, _runtime)
			if _err != nil {
				return _result, _err
			}
			if tea.BoolValue(util.Is4xx(response_.StatusCode)) || tea.BoolValue(util.Is5xx(response_.StatusCode)) {
				_res, _err := util.ReadAsJSON(response_.Body)
				if _err != nil {
					return _result, _err
				}

				err, _err := util.AssertAsMap(_res)
				if _err != nil {
					return _result, _err
				}

				requestId := DefaultAny(err["RequestId"], err["requestId"])
				err["statusCode"] = response_.StatusCode
				_err = tea.NewSDKError(map[string]interface{}{
					"code":               tea.ToString(DefaultAny(err["Code"], err["code"])),
					"message":            "code: " + tea.ToString(tea.IntValue(response_.StatusCode)) + ", " + tea.ToString(DefaultAny(err["Message"], err["message"])) + " request id: " + tea.ToString(requestId),
					"data":               err,
					"description":        tea.ToString(DefaultAny(err["Description"], err["description"])),
					"accessDeniedDetail": DefaultAny(err["AccessDeniedDetail"], err["accessDeniedDetail"]),
				})
				return _result, _err
			}

			if tea.BoolValue(util.EqualString(bodyType, tea.String("binary"))) {
				resp := map[string]interface{}{
					"body":       response_.Body,
					"headers":    response_.Headers,
					"statusCode": tea.IntValue(response_.StatusCode),
				}
				_result = resp
				return _result, _err
			} else if tea.BoolValue(util.EqualString(bodyType, tea.String("byte"))) {
				byt, _err := util.ReadAsBytes(response_.Body)
				if _err != nil {
					return _result, _err
				}

				_result = make(map[string]interface{})
				_err = tea.Convert(map[string]interface{}{
					"body":       byt,
					"headers":    response_.Headers,
					"statusCode": tea.IntValue(response_.StatusCode),
				}, &_result)
				return _result, _err
			} else if tea.BoolValue(util.EqualString(bodyType, tea.String("string"))) {
				str, _err := util.ReadAsString(response_.Body)
				if _err != nil {
					return _result, _err
				}

				_result = make(map[string]interface{})
				_err = tea.Convert(map[string]interface{}{
					"body":       tea.StringValue(str),
					"headers":    response_.Headers,
					"statusCode": tea.IntValue(response_.StatusCode),
				}, &_result)
				return _result, _err
			} else if tea.BoolValue(util.EqualString(bodyType, tea.String("json"))) {
				obj, _err := util.ReadAsJSON(response_.Body)
				if _err != nil {
					return _result, _err
				}

				res, _err := util.AssertAsMap(obj)
				if _err != nil {
					return _result, _err
				}

				_result = make(map[string]interface{})
				_err = tea.Convert(map[string]interface{}{
					"body":       res,
					"headers":    response_.Headers,
					"statusCode": tea.IntValue(response_.StatusCode),
				}, &_result)
				return _result, _err
			} else if tea.BoolValue(util.EqualString(bodyType, tea.String("array"))) {
				arr, _err := util.ReadAsJSON(response_.Body)
				if _err != nil {
					return _result, _err
				}

				_result = make(map[string]interface{})
				_err = tea.Convert(map[string]interface{}{
					"body":       arr,
					"headers":    response_.Headers,
					"statusCode": tea.IntValue(response_.StatusCode),
				}, &_result)
				return _result, _err
			} else {
				_result = make(map[string]interface{})
				_err = tea.Convert(map[string]interface{}{
					"headers":    response_.Headers,
					"statusCode": tea.IntValue(response_.StatusCode),
				}, &_result)
				return _result, _err
			}

		}()
		if !tea.BoolValue(tea.Retryable(_err)) {
			break
		}
	}

	return _resp, _err
}

// Description:
//
// # Encapsulate the request and invoke the network
//
// @param action - api name
//
// @param version - product version
//
// @param protocol - http or https
//
// @param method - e.g. GET
//
// @param authType - authorization type e.g. AK
//
// @param pathname - pathname of every api
//
// @param bodyType - response body type e.g. String
//
// @param request - object of OpenApiRequest
//
// @param runtime - which controls some details of call api, such as retry times
//
// @return the response
func (client *Client) DoROARequest(action *string, version *string, protocol *string, method *string, authType *string, pathname *string, bodyType *string, request *OpenApiRequest, runtime *util.RuntimeOptions) (_result map[string]interface{}, _err error) {
	_err = tea.Validate(request)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Validate(runtime)
	if _err != nil {
		return _result, _err
	}
	_runtime := map[string]interface{}{
		"timeouted":      "retry",
		"key":            tea.StringValue(util.DefaultString(runtime.Key, client.Key)),
		"cert":           tea.StringValue(util.DefaultString(runtime.Cert, client.Cert)),
		"ca":             tea.StringValue(util.DefaultString(runtime.Ca, client.Ca)),
		"readTimeout":    tea.IntValue(util.DefaultNumber(runtime.ReadTimeout, client.ReadTimeout)),
		"connectTimeout": tea.IntValue(util.DefaultNumber(runtime.ConnectTimeout, client.ConnectTimeout)),
		"httpProxy":      tea.StringValue(util.DefaultString(runtime.HttpProxy, client.HttpProxy)),
		"httpsProxy":     tea.StringValue(util.DefaultString(runtime.HttpsProxy, client.HttpsProxy)),
		"noProxy":        tea.StringValue(util.DefaultString(runtime.NoProxy, client.NoProxy)),
		"socks5Proxy":    tea.StringValue(util.DefaultString(runtime.Socks5Proxy, client.Socks5Proxy)),
		"socks5NetWork":  tea.StringValue(util.DefaultString(runtime.Socks5NetWork, client.Socks5NetWork)),
		"maxIdleConns":   tea.IntValue(util.DefaultNumber(runtime.MaxIdleConns, client.MaxIdleConns)),
		"retry": map[string]interface{}{
			"retryable":   tea.BoolValue(runtime.Autoretry),
			"maxAttempts": tea.IntValue(util.DefaultNumber(runtime.MaxAttempts, tea.Int(3))),
		},
		"backoff": map[string]interface{}{
			"policy": tea.StringValue(util.DefaultString(runtime.BackoffPolicy, tea.String("no"))),
			"period": tea.IntValue(util.DefaultNumber(runtime.BackoffPeriod, tea.Int(1))),
		},
		"ignoreSSL": tea.BoolValue(runtime.IgnoreSSL),
	}

	_resp := make(map[string]interface{})
	for _retryTimes := 0; tea.BoolValue(tea.AllowRetry(_runtime["retry"], tea.Int(_retryTimes))); _retryTimes++ {
		if _retryTimes > 0 {
			_backoffTime := tea.GetBackoffTime(_runtime["backoff"], tea.Int(_retryTimes))
			if tea.IntValue(_backoffTime) > 0 {
				tea.Sleep(_backoffTime)
			}
		}

		_resp, _err = func() (map[string]interface{}, error) {
			request_ := tea.NewRequest()
			request_.Protocol = util.DefaultString(client.Protocol, protocol)
			request_.Method = method
			request_.Pathname = pathname
			globalQueries := make(map[string]*string)
			globalHeaders := make(map[string]*string)
			if !tea.BoolValue(util.IsUnset(client.GlobalParameters)) {
				globalParams := client.GlobalParameters
				if !tea.BoolValue(util.IsUnset(globalParams.Queries)) {
					globalQueries = globalParams.Queries
				}

				if !tea.BoolValue(util.IsUnset(globalParams.Headers)) {
					globalHeaders = globalParams.Headers
				}

			}

			extendsHeaders := make(map[string]*string)
			extendsQueries := make(map[string]*string)
			if !tea.BoolValue(util.IsUnset(runtime.ExtendsParameters)) {
				extendsParameters := runtime.ExtendsParameters
				if !tea.BoolValue(util.IsUnset(extendsParameters.Headers)) {
					extendsHeaders = extendsParameters.Headers
				}

				if !tea.BoolValue(util.IsUnset(extendsParameters.Queries)) {
					extendsQueries = extendsParameters.Queries
				}

			}

			request_.Headers = tea.Merge(map[string]*string{
				"date":                    util.GetDateUTCString(),
				"host":                    client.Endpoint,
				"accept":                  tea.String("application/json"),
				"x-acs-signature-nonce":   util.GetNonce(),
				"x-acs-signature-method":  tea.String("HMAC-SHA1"),
				"x-acs-signature-version": tea.String("1.0"),
				"x-acs-version":           version,
				"x-acs-action":            action,
				"user-agent":              util.GetUserAgent(client.UserAgent),
			}, globalHeaders,
				extendsHeaders,
				request.Headers)
			if !tea.BoolValue(util.IsUnset(request.Body)) {
				request_.Body = tea.ToReader(util.ToJSONString(request.Body))
				request_.Headers["content-type"] = tea.String("application/json; charset=utf-8")
			}

			request_.Query = tea.Merge(globalQueries,
				extendsQueries)
			if !tea.BoolValue(util.IsUnset(request.Query)) {
				request_.Query = tea.Merge(request_.Query,
					request.Query)
			}

			if !tea.BoolValue(util.EqualString(authType, tea.String("Anonymous"))) {
				credentialModel, _err := client.Credential.GetCredential()
				if _err != nil {
					return _result, _err
				}

				credentialType := credentialModel.Type
				if tea.BoolValue(util.EqualString(credentialType, tea.String("bearer"))) {
					bearerToken := credentialModel.BearerToken
					request_.Headers["x-acs-bearer-token"] = bearerToken
					request_.Headers["x-acs-signature-type"] = tea.String("BEARERTOKEN")
				} else {
					accessKeyId := credentialModel.AccessKeyId
					accessKeySecret := credentialModel.AccessKeySecret
					securityToken := credentialModel.SecurityToken
					if !tea.BoolValue(util.Empty(securityToken)) {
						request_.Headers["x-acs-accesskey-id"] = accessKeyId
						request_.Headers["x-acs-security-token"] = securityToken
					}

					stringToSign := openapiutil.GetStringToSign(request_)
					request_.Headers["authorization"] = tea.String("acs " + tea.StringValue(accessKeyId) + ":" + tea.StringValue(openapiutil.GetROASignature(stringToSign, accessKeySecret)))
				}

			}

			response_, _err := tea.DoRequest(request_, _runtime)
			if _err != nil {
				return _result, _err
			}
			if tea.BoolValue(util.EqualNumber(response_.StatusCode, tea.Int(204))) {
				_result = make(map[string]interface{})
				_err = tea.Convert(map[string]map[string]*string{
					"headers": response_.Headers,
				}, &_result)
				return _result, _err
			}

			if tea.BoolValue(util.Is4xx(response_.StatusCode)) || tea.BoolValue(util.Is5xx(response_.StatusCode)) {
				_res, _err := util.ReadAsJSON(response_.Body)
				if _err != nil {
					return _result, _err
				}

				err, _err := util.AssertAsMap(_res)
				if _err != nil {
					return _result, _err
				}

				requestId := DefaultAny(err["RequestId"], err["requestId"])
				requestId = DefaultAny(requestId, err["requestid"])
				err["statusCode"] = response_.StatusCode
				_err = tea.NewSDKError(map[string]interface{}{
					"code":               tea.ToString(DefaultAny(err["Code"], err["code"])),
					"message":            "code: " + tea.ToString(tea.IntValue(response_.StatusCode)) + ", " + tea.ToString(DefaultAny(err["Message"], err["message"])) + " request id: " + tea.ToString(requestId),
					"data":               err,
					"description":        tea.ToString(DefaultAny(err["Description"], err["description"])),
					"accessDeniedDetail": DefaultAny(err["AccessDeniedDetail"], err["accessDeniedDetail"]),
				})
				return _result, _err
			}

			if tea.BoolValue(util.EqualString(bodyType, tea.String("binary"))) {
				resp := map[string]interface{}{
					"body":       response_.Body,
					"headers":    response_.Headers,
					"statusCode": tea.IntValue(response_.StatusCode),
				}
				_result = resp
				return _result, _err
			} else if tea.BoolValue(util.EqualString(bodyType, tea.String("byte"))) {
				byt, _err := util.ReadAsBytes(response_.Body)
				if _err != nil {
					return _result, _err
				}

				_result = make(map[string]interface{})
				_err = tea.Convert(map[string]interface{}{
					"body":       byt,
					"headers":    response_.Headers,
					"statusCode": tea.IntValue(response_.StatusCode),
				}, &_result)
				return _result, _err
			} else if tea.BoolValue(util.EqualString(bodyType, tea.String("string"))) {
				str, _err := util.ReadAsString(response_.Body)
				if _err != nil {
					return _result, _err
				}

				_result = make(map[string]interface{})
				_err = tea.Convert(map[string]interface{}{
					"body":       tea.StringValue(str),
					"headers":    response_.Headers,
					"statusCode": tea.IntValue(response_.StatusCode),
				}, &_result)
				return _result, _err
			} else if tea.BoolValue(util.EqualString(bodyType, tea.String("json"))) {
				obj, _err := util.ReadAsJSON(response_.Body)
				if _err != nil {
					return _result, _err
				}

				res, _err := util.AssertAsMap(obj)
				if _err != nil {
					return _result, _err
				}

				_result = make(map[string]interface{})
				_err = tea.Convert(map[string]interface{}{
					"body":       res,
					"headers":    response_.Headers,
					"statusCode": tea.IntValue(response_.StatusCode),
				}, &_result)
				return _result, _err
			} else if tea.BoolValue(util.EqualString(bodyType, tea.String("array"))) {
				arr, _err := util.ReadAsJSON(response_.Body)
				if _err != nil {
					return _result, _err
				}

				_result = make(map[string]interface{})
				_err = tea.Convert(map[string]interface{}{
					"body":       arr,
					"headers":    response_.Headers,
					"statusCode": tea.IntValue(response_.StatusCode),
				}, &_result)
				return _result, _err
			} else {
				_result = make(map[string]interface{})
				_err = tea.Convert(map[string]interface{}{
					"headers":    response_.Headers,
					"statusCode": tea.IntValue(response_.StatusCode),
				}, &_result)
				return _result, _err
			}

		}()
		if !tea.BoolValue(tea.Retryable(_err)) {
			break
		}
	}

	return _resp, _err
}

// Description:
//
// # Encapsulate the request and invoke the network with form body
//
// @param action - api name
//
// @param version - product version
//
// @param protocol - http or https
//
// @param method - e.g. GET
//
// @param authType - authorization type e.g. AK
//
// @param pathname - pathname of every api
//
// @param bodyType - response body type e.g. String
//
// @param request - object of OpenApiRequest
//
// @param runtime - which controls some details of call api, such as retry times
//
// @return the response
func (client *Client) DoROARequestWithForm(action *string, version *string, protocol *string, method *string, authType *string, pathname *string, bodyType *string, request *OpenApiRequest, runtime *util.RuntimeOptions) (_result map[string]interface{}, _err error) {
	_err = tea.Validate(request)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Validate(runtime)
	if _err != nil {
		return _result, _err
	}
	_runtime := map[string]interface{}{
		"timeouted":      "retry",
		"key":            tea.StringValue(util.DefaultString(runtime.Key, client.Key)),
		"cert":           tea.StringValue(util.DefaultString(runtime.Cert, client.Cert)),
		"ca":             tea.StringValue(util.DefaultString(runtime.Ca, client.Ca)),
		"readTimeout":    tea.IntValue(util.DefaultNumber(runtime.ReadTimeout, client.ReadTimeout)),
		"connectTimeout": tea.IntValue(util.DefaultNumber(runtime.ConnectTimeout, client.ConnectTimeout)),
		"httpProxy":      tea.StringValue(util.DefaultString(runtime.HttpProxy, client.HttpProxy)),
		"httpsProxy":     tea.StringValue(util.DefaultString(runtime.HttpsProxy, client.HttpsProxy)),
		"noProxy":        tea.StringValue(util.DefaultString(runtime.NoProxy, client.NoProxy)),
		"socks5Proxy":    tea.StringValue(util.DefaultString(runtime.Socks5Proxy, client.Socks5Proxy)),
		"socks5NetWork":  tea.StringValue(util.DefaultString(runtime.Socks5NetWork, client.Socks5NetWork)),
		"maxIdleConns":   tea.IntValue(util.DefaultNumber(runtime.MaxIdleConns, client.MaxIdleConns)),
		"retry": map[string]interface{}{
			"retryable":   tea.BoolValue(runtime.Autoretry),
			"maxAttempts": tea.IntValue(util.DefaultNumber(runtime.MaxAttempts, tea.Int(3))),
		},
		"backoff": map[string]interface{}{
			"policy": tea.StringValue(util.DefaultString(runtime.BackoffPolicy, tea.String("no"))),
			"period": tea.IntValue(util.DefaultNumber(runtime.BackoffPeriod, tea.Int(1))),
		},
		"ignoreSSL": tea.BoolValue(runtime.IgnoreSSL),
	}

	_resp := make(map[string]interface{})
	for _retryTimes := 0; tea.BoolValue(tea.AllowRetry(_runtime["retry"], tea.Int(_retryTimes))); _retryTimes++ {
		if _retryTimes > 0 {
			_backoffTime := tea.GetBackoffTime(_runtime["backoff"], tea.Int(_retryTimes))
			if tea.IntValue(_backoffTime) > 0 {
				tea.Sleep(_backoffTime)
			}
		}

		_resp, _err = func() (map[string]interface{}, error) {
			request_ := tea.NewRequest()
			request_.Protocol = util.DefaultString(client.Protocol, protocol)
			request_.Method = method
			request_.Pathname = pathname
			globalQueries := make(map[string]*string)
			globalHeaders := make(map[string]*string)
			if !tea.BoolValue(util.IsUnset(client.GlobalParameters)) {
				globalParams := client.GlobalParameters
				if !tea.BoolValue(util.IsUnset(globalParams.Queries)) {
					globalQueries = globalParams.Queries
				}

				if !tea.BoolValue(util.IsUnset(globalParams.Headers)) {
					globalHeaders = globalParams.Headers
				}

			}

			extendsHeaders := make(map[string]*string)
			extendsQueries := make(map[string]*string)
			if !tea.BoolValue(util.IsUnset(runtime.ExtendsParameters)) {
				extendsParameters := runtime.ExtendsParameters
				if !tea.BoolValue(util.IsUnset(extendsParameters.Headers)) {
					extendsHeaders = extendsParameters.Headers
				}

				if !tea.BoolValue(util.IsUnset(extendsParameters.Queries)) {
					extendsQueries = extendsParameters.Queries
				}

			}

			request_.Headers = tea.Merge(map[string]*string{
				"date":                    util.GetDateUTCString(),
				"host":                    client.Endpoint,
				"accept":                  tea.String("application/json"),
				"x-acs-signature-nonce":   util.GetNonce(),
				"x-acs-signature-method":  tea.String("HMAC-SHA1"),
				"x-acs-signature-version": tea.String("1.0"),
				"x-acs-version":           version,
				"x-acs-action":            action,
				"user-agent":              util.GetUserAgent(client.UserAgent),
			}, globalHeaders,
				extendsHeaders,
				request.Headers)
			if !tea.BoolValue(util.IsUnset(request.Body)) {
				m, _err := util.AssertAsMap(request.Body)
				if _err != nil {
					return _result, _err
				}

				request_.Body = tea.ToReader(openapiutil.ToForm(m))
				request_.Headers["content-type"] = tea.String("application/x-www-form-urlencoded")
			}

			request_.Query = tea.Merge(globalQueries,
				extendsQueries)
			if !tea.BoolValue(util.IsUnset(request.Query)) {
				request_.Query = tea.Merge(request_.Query,
					request.Query)
			}

			if !tea.BoolValue(util.EqualString(authType, tea.String("Anonymous"))) {
				credentialModel, _err := client.Credential.GetCredential()
				if _err != nil {
					return _result, _err
				}

				credentialType := credentialModel.Type
				if tea.BoolValue(util.EqualString(credentialType, tea.String("bearer"))) {
					bearerToken := credentialModel.BearerToken
					request_.Headers["x-acs-bearer-token"] = bearerToken
					request_.Headers["x-acs-signature-type"] = tea.String("BEARERTOKEN")
				} else {
					accessKeyId := credentialModel.AccessKeyId
					accessKeySecret := credentialModel.AccessKeySecret
					securityToken := credentialModel.SecurityToken
					if !tea.BoolValue(util.Empty(securityToken)) {
						request_.Headers["x-acs-accesskey-id"] = accessKeyId
						request_.Headers["x-acs-security-token"] = securityToken
					}

					stringToSign := openapiutil.GetStringToSign(request_)
					request_.Headers["authorization"] = tea.String("acs " + tea.StringValue(accessKeyId) + ":" + tea.StringValue(openapiutil.GetROASignature(stringToSign, accessKeySecret)))
				}

			}

			response_, _err := tea.DoRequest(request_, _runtime)
			if _err != nil {
				return _result, _err
			}
			if tea.BoolValue(util.EqualNumber(response_.StatusCode, tea.Int(204))) {
				_result = make(map[string]interface{})
				_err = tea.Convert(map[string]map[string]*string{
					"headers": response_.Headers,
				}, &_result)
				return _result, _err
			}

			if tea.BoolValue(util.Is4xx(response_.StatusCode)) || tea.BoolValue(util.Is5xx(response_.StatusCode)) {
				_res, _err := util.ReadAsJSON(response_.Body)
				if _err != nil {
					return _result, _err
				}

				err, _err := util.AssertAsMap(_res)
				if _err != nil {
					return _result, _err
				}

				err["statusCode"] = response_.StatusCode
				_err = tea.NewSDKError(map[string]interface{}{
					"code":               tea.ToString(DefaultAny(err["Code"], err["code"])),
					"message":            "code: " + tea.ToString(tea.IntValue(response_.StatusCode)) + ", " + tea.ToString(DefaultAny(err["Message"], err["message"])) + " request id: " + tea.ToString(DefaultAny(err["RequestId"], err["requestId"])),
					"data":               err,
					"description":        tea.ToString(DefaultAny(err["Description"], err["description"])),
					"accessDeniedDetail": DefaultAny(err["AccessDeniedDetail"], err["accessDeniedDetail"]),
				})
				return _result, _err
			}

			if tea.BoolValue(util.EqualString(bodyType, tea.String("binary"))) {
				resp := map[string]interface{}{
					"body":       response_.Body,
					"headers":    response_.Headers,
					"statusCode": tea.IntValue(response_.StatusCode),
				}
				_result = resp
				return _result, _err
			} else if tea.BoolValue(util.EqualString(bodyType, tea.String("byte"))) {
				byt, _err := util.ReadAsBytes(response_.Body)
				if _err != nil {
					return _result, _err
				}

				_result = make(map[string]interface{})
				_err = tea.Convert(map[string]interface{}{
					"body":       byt,
					"headers":    response_.Headers,
					"statusCode": tea.IntValue(response_.StatusCode),
				}, &_result)
				return _result, _err
			} else if tea.BoolValue(util.EqualString(bodyType, tea.String("string"))) {
				str, _err := util.ReadAsString(response_.Body)
				if _err != nil {
					return _result, _err
				}

				_result = make(map[string]interface{})
				_err = tea.Convert(map[string]interface{}{
					"body":       tea.StringValue(str),
					"headers":    response_.Headers,
					"statusCode": tea.IntValue(response_.StatusCode),
				}, &_result)
				return _result, _err
			} else if tea.BoolValue(util.EqualString(bodyType, tea.String("json"))) {
				obj, _err := util.ReadAsJSON(response_.Body)
				if _err != nil {
					return _result, _err
				}

				res, _err := util.AssertAsMap(obj)
				if _err != nil {
					return _result, _err
				}

				_result = make(map[string]interface{})
				_err = tea.Convert(map[string]interface{}{
					"body":       res,
					"headers":    response_.Headers,
					"statusCode": tea.IntValue(response_.StatusCode),
				}, &_result)
				return _result, _err
			} else if tea.BoolValue(util.EqualString(bodyType, tea.String("array"))) {
				arr, _err := util.ReadAsJSON(response_.Body)
				if _err != nil {
					return _result, _err
				}

				_result = make(map[string]interface{})
				_err = tea.Convert(map[string]interface{}{
					"body":       arr,
					"headers":    response_.Headers,
					"statusCode": tea.IntValue(response_.StatusCode),
				}, &_result)
				return _result, _err
			} else {
				_result = make(map[string]interface{})
				_err = tea.Convert(map[string]interface{}{
					"headers":    response_.Headers,
					"statusCode": tea.IntValue(response_.StatusCode),
				}, &_result)
				return _result, _err
			}

		}()
		if !tea.BoolValue(tea.Retryable(_err)) {
			break
		}
	}

	return _resp, _err
}

// Description:
//
// # Encapsulate the request and invoke the network
//
// @param action - api name
//
// @param version - product version
//
// @param protocol - http or https
//
// @param method - e.g. GET
//
// @param authType - authorization type e.g. AK
//
// @param bodyType - response body type e.g. String
//
// @param request - object of OpenApiRequest
//
// @param runtime - which controls some details of call api, such as retry times
//
// @return the response
func (client *Client) DoRequest(params *Params, request *OpenApiRequest, runtime *util.RuntimeOptions) (_result map[string]interface{}, _err error) {
	_err = tea.Validate(params)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Validate(request)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Validate(runtime)
	if _err != nil {
		return _result, _err
	}
	_runtime := map[string]interface{}{
		"timeouted":      "retry",
		"key":            tea.StringValue(util.DefaultString(runtime.Key, client.Key)),
		"cert":           tea.StringValue(util.DefaultString(runtime.Cert, client.Cert)),
		"ca":             tea.StringValue(util.DefaultString(runtime.Ca, client.Ca)),
		"readTimeout":    tea.IntValue(util.DefaultNumber(runtime.ReadTimeout, client.ReadTimeout)),
		"connectTimeout": tea.IntValue(util.DefaultNumber(runtime.ConnectTimeout, client.ConnectTimeout)),
		"httpProxy":      tea.StringValue(util.DefaultString(runtime.HttpProxy, client.HttpProxy)),
		"httpsProxy":     tea.StringValue(util.DefaultString(runtime.HttpsProxy, client.HttpsProxy)),
		"noProxy":        tea.StringValue(util.DefaultString(runtime.NoProxy, client.NoProxy)),
		"socks5Proxy":    tea.StringValue(util.DefaultString(runtime.Socks5Proxy, client.Socks5Proxy)),
		"socks5NetWork":  tea.StringValue(util.DefaultString(runtime.Socks5NetWork, client.Socks5NetWork)),
		"maxIdleConns":   tea.IntValue(util.DefaultNumber(runtime.MaxIdleConns, client.MaxIdleConns)),
		"retry": map[string]interface{}{
			"retryable":   tea.BoolValue(runtime.Autoretry),
			"maxAttempts": tea.IntValue(util.DefaultNumber(runtime.MaxAttempts, tea.Int(3))),
		},
		"backoff": map[string]interface{}{
			"policy": tea.StringValue(util.DefaultString(runtime.BackoffPolicy, tea.String("no"))),
			"period": tea.IntValue(util.DefaultNumber(runtime.BackoffPeriod, tea.Int(1))),
		},
		"ignoreSSL": tea.BoolValue(runtime.IgnoreSSL),
	}

	_resp := make(map[string]interface{})
	for _retryTimes := 0; tea.BoolValue(tea.AllowRetry(_runtime["retry"], tea.Int(_retryTimes))); _retryTimes++ {
		if _retryTimes > 0 {
			_backoffTime := tea.GetBackoffTime(_runtime["backoff"], tea.Int(_retryTimes))
			if tea.IntValue(_backoffTime) > 0 {
				tea.Sleep(_backoffTime)
			}
		}

		_resp, _err = func() (map[string]interface{}, error) {
			request_ := tea.NewRequest()
			request_.Protocol = util.DefaultString(client.Protocol, params.Protocol)
			request_.Method = params.Method
			request_.Pathname = params.Pathname
			globalQueries := make(map[string]*string)
			globalHeaders := make(map[string]*string)
			if !tea.BoolValue(util.IsUnset(client.GlobalParameters)) {
				globalParams := client.GlobalParameters
				if !tea.BoolValue(util.IsUnset(globalParams.Queries)) {
					globalQueries = globalParams.Queries
				}

				if !tea.BoolValue(util.IsUnset(globalParams.Headers)) {
					globalHeaders = globalParams.Headers
				}

			}

			extendsHeaders := make(map[string]*string)
			extendsQueries := make(map[string]*string)
			if !tea.BoolValue(util.IsUnset(runtime.ExtendsParameters)) {
				extendsParameters := runtime.ExtendsParameters
				if !tea.BoolValue(util.IsUnset(extendsParameters.Headers)) {
					extendsHeaders = extendsParameters.Headers
				}

				if !tea.BoolValue(util.IsUnset(extendsParameters.Queries)) {
					extendsQueries = extendsParameters.Queries
				}

			}

			request_.Query = tea.Merge(globalQueries,
				extendsQueries,
				request.Query)
			// endpoint is setted in product client
			request_.Headers = tea.Merge(map[string]*string{
				"host":                  client.Endpoint,
				"x-acs-version":         params.Version,
				"x-acs-action":          params.Action,
				"user-agent":            client.GetUserAgent(),
				"x-acs-date":            openapiutil.GetTimestamp(),
				"x-acs-signature-nonce": util.GetNonce(),
				"accept":                tea.String("application/json"),
			}, globalHeaders,
				extendsHeaders,
				request.Headers)
			if tea.BoolValue(util.EqualString(params.Style, tea.String("RPC"))) {
				headers, _err := client.GetRpcHeaders()
				if _err != nil {
					return _result, _err
				}

				if !tea.BoolValue(util.IsUnset(headers)) {
					request_.Headers = tea.Merge(request_.Headers,
						headers)
				}

			}

			signatureAlgorithm := util.DefaultString(client.SignatureAlgorithm, tea.String("ACS3-HMAC-SHA256"))
			hashedRequestPayload := openapiutil.HexEncode(openapiutil.Hash(util.ToBytes(tea.String("")), signatureAlgorithm))
			if !tea.BoolValue(util.IsUnset(request.Stream)) {
				tmp, _err := util.ReadAsBytes(request.Stream)
				if _err != nil {
					return _result, _err
				}

				hashedRequestPayload = openapiutil.HexEncode(openapiutil.Hash(tmp, signatureAlgorithm))
				request_.Body = tea.ToReader(tmp)
				request_.Headers["content-type"] = tea.String("application/octet-stream")
			} else {
				if !tea.BoolValue(util.IsUnset(request.Body)) {
					if tea.BoolValue(util.EqualString(params.ReqBodyType, tea.String("byte"))) {
						byteObj, _err := util.AssertAsBytes(request.Body)
						if _err != nil {
							return _result, _err
						}

						hashedRequestPayload = openapiutil.HexEncode(openapiutil.Hash(byteObj, signatureAlgorithm))
						request_.Body = tea.ToReader(byteObj)
					} else if tea.BoolValue(util.EqualString(params.ReqBodyType, tea.String("json"))) {
						jsonObj := util.ToJSONString(request.Body)
						hashedRequestPayload = openapiutil.HexEncode(openapiutil.Hash(util.ToBytes(jsonObj), signatureAlgorithm))
						request_.Body = tea.ToReader(jsonObj)
						request_.Headers["content-type"] = tea.String("application/json; charset=utf-8")
					} else {
						m, _err := util.AssertAsMap(request.Body)
						if _err != nil {
							return _result, _err
						}

						formObj := openapiutil.ToForm(m)
						hashedRequestPayload = openapiutil.HexEncode(openapiutil.Hash(util.ToBytes(formObj), signatureAlgorithm))
						request_.Body = tea.ToReader(formObj)
						request_.Headers["content-type"] = tea.String("application/x-www-form-urlencoded")
					}

				}

			}

			request_.Headers["x-acs-content-sha256"] = hashedRequestPayload
			if !tea.BoolValue(util.EqualString(params.AuthType, tea.String("Anonymous"))) {
				credentialModel, _err := client.Credential.GetCredential()
				if _err != nil {
					return _result, _err
				}

				authType := credentialModel.Type
				if tea.BoolValue(util.EqualString(authType, tea.String("bearer"))) {
					bearerToken := credentialModel.BearerToken
					request_.Headers["x-acs-bearer-token"] = bearerToken
					if tea.BoolValue(util.EqualString(params.Style, tea.String("RPC"))) {
						request_.Query["SignatureType"] = tea.String("BEARERTOKEN")
					} else {
						request_.Headers["x-acs-signature-type"] = tea.String("BEARERTOKEN")
					}

				} else {
					accessKeyId := credentialModel.AccessKeyId
					accessKeySecret := credentialModel.AccessKeySecret
					securityToken := credentialModel.SecurityToken
					if !tea.BoolValue(util.Empty(securityToken)) {
						request_.Headers["x-acs-accesskey-id"] = accessKeyId
						request_.Headers["x-acs-security-token"] = securityToken
					}

					request_.Headers["Authorization"] = openapiutil.GetAuthorization(request_, signatureAlgorithm, hashedRequestPayload, accessKeyId, accessKeySecret)
				}

			}

			response_, _err := tea.DoRequest(request_, _runtime)
			if _err != nil {
				return _result, _err
			}
			if tea.BoolValue(util.Is4xx(response_.StatusCode)) || tea.BoolValue(util.Is5xx(response_.StatusCode)) {
				err := map[string]interface{}{}
				if !tea.BoolValue(util.IsUnset(response_.Headers["content-type"])) && tea.BoolValue(util.EqualString(response_.Headers["content-type"], tea.String("text/xml;charset=utf-8"))) {
					_str, _err := util.ReadAsString(response_.Body)
					if _err != nil {
						return _result, _err
					}

					respMap := xml.ParseXml(_str, nil)
					err, _err = util.AssertAsMap(respMap["Error"])
					if _err != nil {
						return _result, _err
					}

				} else {
					_res, _err := util.ReadAsJSON(response_.Body)
					if _err != nil {
						return _result, _err
					}

					err, _err = util.AssertAsMap(_res)
					if _err != nil {
						return _result, _err
					}

				}

				err["statusCode"] = response_.StatusCode
				_err = tea.NewSDKError(map[string]interface{}{
					"code":               tea.ToString(DefaultAny(err["Code"], err["code"])),
					"message":            "code: " + tea.ToString(tea.IntValue(response_.StatusCode)) + ", " + tea.ToString(DefaultAny(err["Message"], err["message"])) + " request id: " + tea.ToString(DefaultAny(err["RequestId"], err["requestId"])),
					"data":               err,
					"description":        tea.ToString(DefaultAny(err["Description"], err["description"])),
					"accessDeniedDetail": DefaultAny(err["AccessDeniedDetail"], err["accessDeniedDetail"]),
				})
				return _result, _err
			}

			if tea.BoolValue(util.EqualString(params.BodyType, tea.String("binary"))) {
				resp := map[string]interface{}{
					"body":       response_.Body,
					"headers":    response_.Headers,
					"statusCode": tea.IntValue(response_.StatusCode),
				}
				_result = resp
				return _result, _err
			} else if tea.BoolValue(util.EqualString(params.BodyType, tea.String("byte"))) {
				byt, _err := util.ReadAsBytes(response_.Body)
				if _err != nil {
					return _result, _err
				}

				_result = make(map[string]interface{})
				_err = tea.Convert(map[string]interface{}{
					"body":       byt,
					"headers":    response_.Headers,
					"statusCode": tea.IntValue(response_.StatusCode),
				}, &_result)
				return _result, _err
			} else if tea.BoolValue(util.EqualString(params.BodyType, tea.String("string"))) {
				str, _err := util.ReadAsString(response_.Body)
				if _err != nil {
					return _result, _err
				}

				_result = make(map[string]interface{})
				_err = tea.Convert(map[string]interface{}{
					"body":       tea.StringValue(str),
					"headers":    response_.Headers,
					"statusCode": tea.IntValue(response_.StatusCode),
				}, &_result)
				return _result, _err
			} else if tea.BoolValue(util.EqualString(params.BodyType, tea.String("json"))) {
				obj, _err := util.ReadAsJSON(response_.Body)
				if _err != nil {
					return _result, _err
				}

				res, _err := util.AssertAsMap(obj)
				if _err != nil {
					return _result, _err
				}

				_result = make(map[string]interface{})
				_err = tea.Convert(map[string]interface{}{
					"body":       res,
					"headers":    response_.Headers,
					"statusCode": tea.IntValue(response_.StatusCode),
				}, &_result)
				return _result, _err
			} else if tea.BoolValue(util.EqualString(params.BodyType, tea.String("array"))) {
				arr, _err := util.ReadAsJSON(response_.Body)
				if _err != nil {
					return _result, _err
				}

				_result = make(map[string]interface{})
				_err = tea.Convert(map[string]interface{}{
					"body":       arr,
					"headers":    response_.Headers,
					"statusCode": tea.IntValue(response_.StatusCode),
				}, &_result)
				return _result, _err
			} else {
				anything, _err := util.ReadAsString(response_.Body)
				if _err != nil {
					return _result, _err
				}

				_result = make(map[string]interface{})
				_err = tea.Convert(map[string]interface{}{
					"body":       tea.StringValue(anything),
					"headers":    response_.Headers,
					"statusCode": tea.IntValue(response_.StatusCode),
				}, &_result)
				return _result, _err
			}

		}()
		if !tea.BoolValue(tea.Retryable(_err)) {
			break
		}
	}

	return _resp, _err
}

// Description:
//
// # Encapsulate the request and invoke the network
//
// @param action - api name
//
// @param version - product version
//
// @param protocol - http or https
//
// @param method - e.g. GET
//
// @param authType - authorization type e.g. AK
//
// @param bodyType - response body type e.g. String
//
// @param request - object of OpenApiRequest
//
// @param runtime - which controls some details of call api, such as retry times
//
// @return the response
func (client *Client) Execute(params *Params, request *OpenApiRequest, runtime *util.RuntimeOptions) (_result map[string]interface{}, _err error) {
	_err = tea.Validate(params)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Validate(request)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Validate(runtime)
	if _err != nil {
		return _result, _err
	}
	_runtime := map[string]interface{}{
		"timeouted":      "retry",
		"key":            tea.StringValue(util.DefaultString(runtime.Key, client.Key)),
		"cert":           tea.StringValue(util.DefaultString(runtime.Cert, client.Cert)),
		"ca":             tea.StringValue(util.DefaultString(runtime.Ca, client.Ca)),
		"readTimeout":    tea.IntValue(util.DefaultNumber(runtime.ReadTimeout, client.ReadTimeout)),
		"connectTimeout": tea.IntValue(util.DefaultNumber(runtime.ConnectTimeout, client.ConnectTimeout)),
		"httpProxy":      tea.StringValue(util.DefaultString(runtime.HttpProxy, client.HttpProxy)),
		"httpsProxy":     tea.StringValue(util.DefaultString(runtime.HttpsProxy, client.HttpsProxy)),
		"noProxy":        tea.StringValue(util.DefaultString(runtime.NoProxy, client.NoProxy)),
		"socks5Proxy":    tea.StringValue(util.DefaultString(runtime.Socks5Proxy, client.Socks5Proxy)),
		"socks5NetWork":  tea.StringValue(util.DefaultString(runtime.Socks5NetWork, client.Socks5NetWork)),
		"maxIdleConns":   tea.IntValue(util.DefaultNumber(runtime.MaxIdleConns, client.MaxIdleConns)),
		"retry": map[string]interface{}{
			"retryable":   tea.BoolValue(runtime.Autoretry),
			"maxAttempts": tea.IntValue(util.DefaultNumber(runtime.MaxAttempts, tea.Int(3))),
		},
		"backoff": map[string]interface{}{
			"policy": tea.StringValue(util.DefaultString(runtime.BackoffPolicy, tea.String("no"))),
			"period": tea.IntValue(util.DefaultNumber(runtime.BackoffPeriod, tea.Int(1))),
		},
		"ignoreSSL":    tea.BoolValue(runtime.IgnoreSSL),
		"disableHttp2": DefaultAny(client.DisableHttp2, tea.Bool(false)),
	}

	_resp := make(map[string]interface{})
	for _retryTimes := 0; tea.BoolValue(tea.AllowRetry(_runtime["retry"], tea.Int(_retryTimes))); _retryTimes++ {
		if _retryTimes > 0 {
			_backoffTime := tea.GetBackoffTime(_runtime["backoff"], tea.Int(_retryTimes))
			if tea.IntValue(_backoffTime) > 0 {
				tea.Sleep(_backoffTime)
			}
		}

		_resp, _err = func() (map[string]interface{}, error) {
			request_ := tea.NewRequest()
			// spi = new Gateway();//Gateway implements SPI，这一步在产品 SDK 中实例化
			headers, _err := client.GetRpcHeaders()
			if _err != nil {
				return _result, _err
			}

			globalQueries := make(map[string]*string)
			globalHeaders := make(map[string]*string)
			if !tea.BoolValue(util.IsUnset(client.GlobalParameters)) {
				globalParams := client.GlobalParameters
				if !tea.BoolValue(util.IsUnset(globalParams.Queries)) {
					globalQueries = globalParams.Queries
				}

				if !tea.BoolValue(util.IsUnset(globalParams.Headers)) {
					globalHeaders = globalParams.Headers
				}

			}

			extendsHeaders := make(map[string]*string)
			extendsQueries := make(map[string]*string)
			if !tea.BoolValue(util.IsUnset(runtime.ExtendsParameters)) {
				extendsParameters := runtime.ExtendsParameters
				if !tea.BoolValue(util.IsUnset(extendsParameters.Headers)) {
					extendsHeaders = extendsParameters.Headers
				}

				if !tea.BoolValue(util.IsUnset(extendsParameters.Queries)) {
					extendsQueries = extendsParameters.Queries
				}

			}

			requestContext := &spi.InterceptorContextRequest{
				Headers: tea.Merge(globalHeaders,
					extendsHeaders,
					request.Headers,
					headers),
				Query: tea.Merge(globalQueries,
					extendsQueries,
					request.Query),
				Body:               request.Body,
				Stream:             request.Stream,
				HostMap:            request.HostMap,
				Pathname:           params.Pathname,
				ProductId:          client.ProductId,
				Action:             params.Action,
				Version:            params.Version,
				Protocol:           util.DefaultString(client.Protocol, params.Protocol),
				Method:             util.DefaultString(client.Method, params.Method),
				AuthType:           params.AuthType,
				BodyType:           params.BodyType,
				ReqBodyType:        params.ReqBodyType,
				Style:              params.Style,
				Credential:         client.Credential,
				SignatureVersion:   client.SignatureVersion,
				SignatureAlgorithm: client.SignatureAlgorithm,
				UserAgent:          client.GetUserAgent(),
			}
			configurationContext := &spi.InterceptorContextConfiguration{
				RegionId:     client.RegionId,
				Endpoint:     util.DefaultString(request.EndpointOverride, client.Endpoint),
				EndpointRule: client.EndpointRule,
				EndpointMap:  client.EndpointMap,
				EndpointType: client.EndpointType,
				Network:      client.Network,
				Suffix:       client.Suffix,
			}
			interceptorContext := &spi.InterceptorContext{
				Request:       requestContext,
				Configuration: configurationContext,
			}
			attributeMap := &spi.AttributeMap{}
			// 1. spi.modifyConfiguration(context: SPI.InterceptorContext, attributeMap: SPI.AttributeMap);
			_err = client.Spi.ModifyConfiguration(interceptorContext, attributeMap)
			if _err != nil {
				return _result, _err
			}
			// 2. spi.modifyRequest(context: SPI.InterceptorContext, attributeMap: SPI.AttributeMap);
			_err = client.Spi.ModifyRequest(interceptorContext, attributeMap)
			if _err != nil {
				return _result, _err
			}
			request_.Protocol = interceptorContext.Request.Protocol
			request_.Method = interceptorContext.Request.Method
			request_.Pathname = interceptorContext.Request.Pathname
			request_.Query = interceptorContext.Request.Query
			request_.Body = interceptorContext.Request.Stream
			request_.Headers = interceptorContext.Request.Headers
			response_, _err := tea.DoRequest(request_, _runtime)
			if _err != nil {
				return _result, _err
			}
			responseContext := &spi.InterceptorContextResponse{
				StatusCode: response_.StatusCode,
				Headers:    response_.Headers,
				Body:       response_.Body,
			}
			interceptorContext.Response = responseContext
			// 3. spi.modifyResponse(context: SPI.InterceptorContext, attributeMap: SPI.AttributeMap);
			_err = client.Spi.ModifyResponse(interceptorContext, attributeMap)
			if _err != nil {
				return _result, _err
			}
			_result = make(map[string]interface{})
			_err = tea.Convert(map[string]interface{}{
				"headers":    interceptorContext.Response.Headers,
				"statusCode": tea.IntValue(interceptorContext.Response.StatusCode),
				"body":       interceptorContext.Response.DeserializedBody,
			}, &_result)
			return _result, _err
		}()
		if !tea.BoolValue(tea.Retryable(_err)) {
			break
		}
	}

	return _resp, _err
}

func (client *Client) CallApi(params *Params, request *OpenApiRequest, runtime *util.RuntimeOptions) (_result map[string]interface{}, _err error) {
	if tea.BoolValue(util.IsUnset(params)) {
		_err = tea.NewSDKError(map[string]interface{}{
			"code":    "ParameterMissing",
			"message": "'params' can not be unset",
		})
		return _result, _err
	}

	if tea.BoolValue(util.IsUnset(client.SignatureAlgorithm)) || !tea.BoolValue(util.EqualString(client.SignatureAlgorithm, tea.String("v2"))) {
		_result = make(map[string]interface{})
		_body, _err := client.DoRequest(params, request, runtime)
		if _err != nil {
			return _result, _err
		}
		_result = _body
		return _result, _err
	} else if tea.BoolValue(util.EqualString(params.Style, tea.String("ROA"))) && tea.BoolValue(util.EqualString(params.ReqBodyType, tea.String("json"))) {
		_result = make(map[string]interface{})
		_body, _err := client.DoROARequest(params.Action, params.Version, params.Protocol, params.Method, params.AuthType, params.Pathname, params.BodyType, request, runtime)
		if _err != nil {
			return _result, _err
		}
		_result = _body
		return _result, _err
	} else if tea.BoolValue(util.EqualString(params.Style, tea.String("ROA"))) {
		_result = make(map[string]interface{})
		_body, _err := client.DoROARequestWithForm(params.Action, params.Version, params.Protocol, params.Method, params.AuthType, params.Pathname, params.BodyType, request, runtime)
		if _err != nil {
			return _result, _err
		}
		_result = _body
		return _result, _err
	} else {
		_result = make(map[string]interface{})
		_body, _err := client.DoRPCRequest(params.Action, params.Version, params.Protocol, params.Method, params.AuthType, params.BodyType, request, runtime)
		if _err != nil {
			return _result, _err
		}
		_result = _body
		return _result, _err
	}

}

// Description:
//
// # Get user agent
//
// @return user agent
func (client *Client) GetUserAgent() (_result *string) {
	userAgent := util.GetUserAgent(client.UserAgent)
	_result = userAgent
	return _result
}

// Description:
//
// # Get accesskey id by using credential
//
// @return accesskey id
func (client *Client) GetAccessKeyId() (_result *string, _err error) {
	if tea.BoolValue(util.IsUnset(client.Credential)) {
		_result = tea.String("")
		return _result, _err
	}

	accessKeyId, _err := client.Credential.GetAccessKeyId()
	if _err != nil {
		return _result, _err
	}

	_result = accessKeyId
	return _result, _err
}

// Description:
//
// # Get accesskey secret by using credential
//
// @return accesskey secret
func (client *Client) GetAccessKeySecret() (_result *string, _err error) {
	if tea.BoolValue(util.IsUnset(client.Credential)) {
		_result = tea.String("")
		return _result, _err
	}

	secret, _err := client.Credential.GetAccessKeySecret()
	if _err != nil {
		return _result, _err
	}

	_result = secret
	return _result, _err
}

// Description:
//
// # Get security token by using credential
//
// @return security token
func (client *Client) GetSecurityToken() (_result *string, _err error) {
	if tea.BoolValue(util.IsUnset(client.Credential)) {
		_result = tea.String("")
		return _result, _err
	}

	token, _err := client.Credential.GetSecurityToken()
	if _err != nil {
		return _result, _err
	}

	_result = token
	return _result, _err
}

// Description:
//
// # Get bearer token by credential
//
// @return bearer token
func (client *Client) GetBearerToken() (_result *string, _err error) {
	if tea.BoolValue(util.IsUnset(client.Credential)) {
		_result = tea.String("")
		return _result, _err
	}

	token := client.Credential.GetBearerToken()
	_result = token
	return _result, _err
}

// Description:
//
// # Get credential type by credential
//
// @return credential type e.g. access_key
func (client *Client) GetType() (_result *string, _err error) {
	if tea.BoolValue(util.IsUnset(client.Credential)) {
		_result = tea.String("")
		return _result, _err
	}

	authType := client.Credential.GetType()
	_result = authType
	return _result, _err
}

// Description:
//
// # If inputValue is not null, return it or return defaultValue
//
// @param inputValue - users input value
//
// @param defaultValue - default value
//
// @return the final result
func DefaultAny(inputValue interface{}, defaultValue interface{}) (_result interface{}) {
	if tea.BoolValue(util.IsUnset(inputValue)) {
		_result = defaultValue
		return _result
	}

	_result = inputValue
	return _result
}

// Description:
//
// # If the endpointRule and config.endpoint are empty, throw error
//
// @param config - config contains the necessary information to create a client
func (client *Client) CheckConfig(config *Config) (_err error) {
	if tea.BoolValue(util.Empty(client.EndpointRule)) && tea.BoolValue(util.Empty(config.Endpoint)) {
		_err = tea.NewSDKError(map[string]interface{}{
			"code":    "ParameterMissing",
			"message": "'config.endpoint' can not be empty",
		})
		return _err
	}

	return _err
}

// Description:
//
// set gateway client
//
// @param spi - .
func (client *Client) SetGatewayClient(spi spi.ClientInterface) (_err error) {
	client.Spi = spi
	return _err
}

// Description:
//
// set RPC header for debug
//
// @param headers - headers for debug, this header can be used only once.
func (client *Client) SetRpcHeaders(headers map[string]*string) (_err error) {
	client.Headers = headers
	return _err
}

// Description:
//
// get RPC header for debug
func (client *Client) GetRpcHeaders() (_result map[string]*string, _err error) {
	headers := client.Headers
	client.Headers = nil
	_result = headers
	return _result, _err
}
