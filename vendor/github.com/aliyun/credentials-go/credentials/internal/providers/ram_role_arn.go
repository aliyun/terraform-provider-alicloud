package providers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	httputil "github.com/aliyun/credentials-go/credentials/internal/http"
	"github.com/aliyun/credentials-go/credentials/internal/utils"
)

type assumedRoleUser struct {
}

type credentials struct {
	SecurityToken   *string `json:"SecurityToken"`
	Expiration      *string `json:"Expiration"`
	AccessKeySecret *string `json:"AccessKeySecret"`
	AccessKeyId     *string `json:"AccessKeyId"`
}

type assumeRoleResponse struct {
	RequestID       *string          `json:"RequestId"`
	AssumedRoleUser *assumedRoleUser `json:"AssumedRoleUser"`
	Credentials     *credentials     `json:"Credentials"`
}

type sessionCredentials struct {
	AccessKeyId     string
	AccessKeySecret string
	SecurityToken   string
	Expiration      string
}

type HttpOptions struct {
	Proxy          string
	ConnectTimeout int
	ReadTimeout    int
}

type RAMRoleARNCredentialsProvider struct {
	credentialsProvider CredentialsProvider
	roleArn             string
	roleSessionName     string
	durationSeconds     int
	policy              string
	externalId          string
	// for sts endpoint
	stsRegionId string
	stsEndpoint string
	// for http options
	httpOptions *HttpOptions
	// inner
	expirationTimestamp int64
	lastUpdateTimestamp int64
	sessionCredentials  *sessionCredentials
}

type RAMRoleARNCredentialsProviderBuilder struct {
	provider *RAMRoleARNCredentialsProvider
}

func NewRAMRoleARNCredentialsProviderBuilder() *RAMRoleARNCredentialsProviderBuilder {
	return &RAMRoleARNCredentialsProviderBuilder{
		provider: &RAMRoleARNCredentialsProvider{},
	}
}

func (builder *RAMRoleARNCredentialsProviderBuilder) WithCredentialsProvider(credentialsProvider CredentialsProvider) *RAMRoleARNCredentialsProviderBuilder {
	builder.provider.credentialsProvider = credentialsProvider
	return builder
}

func (builder *RAMRoleARNCredentialsProviderBuilder) WithRoleArn(roleArn string) *RAMRoleARNCredentialsProviderBuilder {
	builder.provider.roleArn = roleArn
	return builder
}

func (builder *RAMRoleARNCredentialsProviderBuilder) WithStsRegionId(regionId string) *RAMRoleARNCredentialsProviderBuilder {
	builder.provider.stsRegionId = regionId
	return builder
}

func (builder *RAMRoleARNCredentialsProviderBuilder) WithStsEndpoint(endpoint string) *RAMRoleARNCredentialsProviderBuilder {
	builder.provider.stsEndpoint = endpoint
	return builder
}

func (builder *RAMRoleARNCredentialsProviderBuilder) WithRoleSessionName(roleSessionName string) *RAMRoleARNCredentialsProviderBuilder {
	builder.provider.roleSessionName = roleSessionName
	return builder
}

func (builder *RAMRoleARNCredentialsProviderBuilder) WithPolicy(policy string) *RAMRoleARNCredentialsProviderBuilder {
	builder.provider.policy = policy
	return builder
}

func (builder *RAMRoleARNCredentialsProviderBuilder) WithExternalId(externalId string) *RAMRoleARNCredentialsProviderBuilder {
	builder.provider.externalId = externalId
	return builder
}

func (builder *RAMRoleARNCredentialsProviderBuilder) WithDurationSeconds(durationSeconds int) *RAMRoleARNCredentialsProviderBuilder {
	builder.provider.durationSeconds = durationSeconds
	return builder
}

func (builder *RAMRoleARNCredentialsProviderBuilder) WithHttpOptions(httpOptions *HttpOptions) *RAMRoleARNCredentialsProviderBuilder {
	builder.provider.httpOptions = httpOptions
	return builder
}

func (builder *RAMRoleARNCredentialsProviderBuilder) Build() (provider *RAMRoleARNCredentialsProvider, err error) {
	if builder.provider.credentialsProvider == nil {
		err = errors.New("must specify a previous credentials provider to asssume role")
		return
	}

	if builder.provider.roleArn == "" {
		err = errors.New("the RoleArn is empty")
		return
	}

	if builder.provider.roleSessionName == "" {
		builder.provider.roleSessionName = "credentials-go-" + strconv.FormatInt(time.Now().UnixNano()/1000, 10)
	}

	// duration seconds
	if builder.provider.durationSeconds == 0 {
		// default to 3600
		builder.provider.durationSeconds = 3600
	}

	if builder.provider.durationSeconds < 900 {
		err = errors.New("session duration should be in the range of 900s - max session duration")
		return
	}

	// sts endpoint
	if builder.provider.stsEndpoint == "" {
		if builder.provider.stsRegionId != "" {
			builder.provider.stsEndpoint = fmt.Sprintf("sts.%s.aliyuncs.com", builder.provider.stsRegionId)
		} else {
			builder.provider.stsEndpoint = "sts.aliyuncs.com"
		}
	}

	provider = builder.provider
	return
}

func (provider *RAMRoleARNCredentialsProvider) getCredentials(cc *Credentials) (session *sessionCredentials, err error) {
	method := "POST"
	req := &httputil.Request{
		Method:   method,
		Protocol: "https",
		Host:     provider.stsEndpoint,
		Headers:  map[string]string{},
	}

	queries := make(map[string]string)
	queries["Version"] = "2015-04-01"
	queries["Action"] = "AssumeRole"
	queries["Format"] = "JSON"
	queries["Timestamp"] = utils.GetTimeInFormatISO8601()
	queries["SignatureMethod"] = "HMAC-SHA1"
	queries["SignatureVersion"] = "1.0"
	queries["SignatureNonce"] = utils.GetNonce()
	queries["AccessKeyId"] = cc.AccessKeyId

	if cc.SecurityToken != "" {
		queries["SecurityToken"] = cc.SecurityToken
	}

	bodyForm := make(map[string]string)
	bodyForm["RoleArn"] = provider.roleArn
	if provider.policy != "" {
		bodyForm["Policy"] = provider.policy
	}
	if provider.externalId != "" {
		bodyForm["ExternalId"] = provider.externalId
	}
	bodyForm["RoleSessionName"] = provider.roleSessionName
	bodyForm["DurationSeconds"] = strconv.Itoa(provider.durationSeconds)
	req.Form = bodyForm

	// caculate signature
	signParams := make(map[string]string)
	for key, value := range queries {
		signParams[key] = value
	}
	for key, value := range bodyForm {
		signParams[key] = value
	}

	stringToSign := utils.GetURLFormedMap(signParams)
	stringToSign = strings.Replace(stringToSign, "+", "%20", -1)
	stringToSign = strings.Replace(stringToSign, "*", "%2A", -1)
	stringToSign = strings.Replace(stringToSign, "%7E", "~", -1)
	stringToSign = url.QueryEscape(stringToSign)
	stringToSign = method + "&%2F&" + stringToSign
	secret := cc.AccessKeySecret + "&"
	queries["Signature"] = utils.ShaHmac1(stringToSign, secret)

	req.Queries = queries

	// set headers
	req.Headers["Accept-Encoding"] = "identity"
	req.Headers["Content-Type"] = "application/x-www-form-urlencoded"
	req.Headers["x-acs-credentials-provider"] = cc.ProviderName

	if provider.httpOptions != nil {
		req.ConnectTimeout = time.Duration(provider.httpOptions.ConnectTimeout) * time.Second
		req.ReadTimeout = time.Duration(provider.httpOptions.ReadTimeout) * time.Second
		req.Proxy = provider.httpOptions.Proxy
	}

	res, err := httpDo(req)
	if err != nil {
		return
	}

	if res.StatusCode != http.StatusOK {
		err = errors.New("refresh session token failed: " + string(res.Body))
		return
	}
	var data assumeRoleResponse
	err = json.Unmarshal(res.Body, &data)
	if err != nil {
		err = fmt.Errorf("refresh RoleArn sts token err, json.Unmarshal fail: %s", err.Error())
		return
	}
	if data.Credentials == nil {
		err = fmt.Errorf("refresh RoleArn sts token err, fail to get credentials")
		return
	}

	if data.Credentials.AccessKeyId == nil || data.Credentials.AccessKeySecret == nil || data.Credentials.SecurityToken == nil {
		err = fmt.Errorf("refresh RoleArn sts token err, fail to get credentials")
		return
	}

	session = &sessionCredentials{
		AccessKeyId:     *data.Credentials.AccessKeyId,
		AccessKeySecret: *data.Credentials.AccessKeySecret,
		SecurityToken:   *data.Credentials.SecurityToken,
		Expiration:      *data.Credentials.Expiration,
	}
	return
}

func (provider *RAMRoleARNCredentialsProvider) needUpdateCredential() (result bool) {
	if provider.expirationTimestamp == 0 {
		return true
	}

	return provider.expirationTimestamp-time.Now().Unix() <= 180
}

func (provider *RAMRoleARNCredentialsProvider) GetCredentials() (cc *Credentials, err error) {
	if provider.sessionCredentials == nil || provider.needUpdateCredential() {
		// 获取前置凭证
		previousCredentials, err1 := provider.credentialsProvider.GetCredentials()
		if err1 != nil {
			return nil, err1
		}
		sessionCredentials, err2 := provider.getCredentials(previousCredentials)
		if err2 != nil {
			return nil, err2
		}

		expirationTime, err := time.Parse("2006-01-02T15:04:05Z", sessionCredentials.Expiration)
		if err != nil {
			return nil, err
		}

		provider.expirationTimestamp = expirationTime.Unix()
		provider.lastUpdateTimestamp = time.Now().Unix()
		provider.sessionCredentials = sessionCredentials
	}

	cc = &Credentials{
		AccessKeyId:     provider.sessionCredentials.AccessKeyId,
		AccessKeySecret: provider.sessionCredentials.AccessKeySecret,
		SecurityToken:   provider.sessionCredentials.SecurityToken,
		ProviderName:    fmt.Sprintf("%s/%s", provider.GetProviderName(), provider.credentialsProvider.GetProviderName()),
	}
	return
}

func (provider *RAMRoleARNCredentialsProvider) GetProviderName() string {
	return "ram_role_arn"
}
