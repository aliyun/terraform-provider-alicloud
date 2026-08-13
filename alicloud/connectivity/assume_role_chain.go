package connectivity

import (
	"fmt"
	"log"
	"sync"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	sts20150401 "github.com/alibabacloud-go/sts-20150401/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	credential "github.com/aliyun/credentials-go/credentials"
)

// defaultStsEndpoint is the STS service domain used when no explicit STS
// endpoint is configured on the provider. It mirrors the default used by the
// credentials-go SDK for the ram_role_arn credential type.
const defaultStsEndpoint = "sts.aliyuncs.com"

// chainedAssumeRoleCredential implements credential.Credential by performing
// an ordered chain of STS AssumeRole calls. The first hop uses the initial
// caller credential captured at provider configuration time; every subsequent
// hop uses the temporary credential returned by the previous AssumeRole call
// as its caller credential. The final temporary credential is cached and
// refreshed shortly before it expires, so downstream callers that periodically
// re-resolve via GetCredential continue to receive a live session.
type chainedAssumeRoleCredential struct {
	// initial caller credential, captured when the provider was configured
	initialAccessKey     string
	initialSecretKey     string
	initialSecurityToken string

	chain []*AssumeRoleConfig

	regionId        string
	stsEndpoint     string
	userAgent       string
	sourceIp        string
	secureTransport string
	readTimeout     int
	connectTimeout  int

	mu     sync.Mutex
	cached *chainedSession
}

// chainedSession is the cached result of a fully resolved chain.
type chainedSession struct {
	accessKeyId     string
	accessKeySecret string
	securityToken   string
	expiration      time.Time
}

// setAuthByChainedAssumeRole builds a chainedAssumeRoleCredential from the
// configured AssumeRoleChain, installs it as the provider's credential source,
// and eagerly resolves the chain once so that Config.AccessKey/SecretKey/
// SecurityToken hold the final session credential for downstream callers that
// read those fields directly.
func (c *Config) setAuthByChainedAssumeRole() (err error) {
	if len(c.AssumeRoleChain) == 0 {
		return
	}
	// The chain needs an initial caller credential. The provider configuration
	// guarantees AccessKey/SecretKey are populated by the time
	// RefreshAuthCredential runs; without them there is nothing to chain from.
	if c.AccessKey == "" {
		return fmt.Errorf("assume_role chain requires a caller credential: access_key is empty")
	}

	stsEndpoint := c.StsEndpoint
	if stsEndpoint == "" {
		stsEndpoint = defaultStsEndpoint
	}

	provider := &chainedAssumeRoleCredential{
		initialAccessKey:     c.AccessKey,
		initialSecretKey:     c.SecretKey,
		initialSecurityToken: c.SecurityToken,
		chain:                c.AssumeRoleChain,
		regionId:             c.RegionId,
		stsEndpoint:          stsEndpoint,
		userAgent:            c.getUserAgent(),
		sourceIp:             c.SourceIp,
		secureTransport:      c.SecureTransport,
		readTimeout:          c.ClientReadTimeout,
		connectTimeout:       c.ClientConnectTimeout,
	}
	c.Credential = provider
	cred, err := provider.GetCredential()
	if err != nil || cred == nil {
		return fmt.Errorf("refresh chained AssumeRole credential failed. Error: %v", err)
	}
	c.AccessKey, c.SecretKey, c.SecurityToken = *cred.AccessKeyId, *cred.AccessKeySecret, *cred.SecurityToken
	return nil
}

// GetCredential resolves the chain, returning the cached final session when it
// is still valid, otherwise re-running the whole chain from the initial caller
// credential.
func (p *chainedAssumeRoleCredential) GetCredential() (*credential.CredentialModel, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.cached != nil && time.Now().Before(p.cached.expiration.Add(-time.Minute)) {
		return p.model(), nil
	}

	ak := p.initialAccessKey
	sk := p.initialSecretKey
	st := p.initialSecurityToken
	var lastExpiration string
	for i, hop := range p.chain {
		creds, expiration, err := p.assumeRole(ak, sk, st, hop)
		if err != nil {
			return nil, fmt.Errorf("chained AssumeRole failed at hop %d (RoleArn: %s): %v", i, hop.RoleArn, err)
		}
		ak = creds.AccessKeyId
		sk = creds.AccessKeySecret
		st = creds.SecurityToken
		lastExpiration = expiration
		log.Printf("[INFO] chained AssumeRole hop %d/%d assumed (RoleArn: %s)", i+1, len(p.chain), hop.RoleArn)
	}

	expiration := time.Now().Add(time.Hour)
	if lastExpiration != "" {
		if parsed, err := time.Parse("2006-01-02T15:04:05Z", lastExpiration); err == nil {
			expiration = parsed
		}
	}
	p.cached = &chainedSession{
		accessKeyId:     ak,
		accessKeySecret: sk,
		securityToken:   st,
		expiration:      expiration,
	}
	return p.model(), nil
}

func (p *chainedAssumeRoleCredential) model() *credential.CredentialModel {
	return &credential.CredentialModel{
		AccessKeyId:     tea.String(p.cached.accessKeyId),
		AccessKeySecret: tea.String(p.cached.accessKeySecret),
		SecurityToken:   tea.String(p.cached.securityToken),
		Type:            tea.String("ram_role_arn"),
		ProviderName:    tea.String("chained_assume_role"),
	}
}

// assumeRole performs a single STS AssumeRole call using the supplied caller
// credential and returns the temporary credential and expiration string.
type stsCredential struct {
	AccessKeyId     string
	AccessKeySecret string
	SecurityToken   string
}

func (p *chainedAssumeRoleCredential) assumeRole(accessKey, secretKey, securityToken string, hop *AssumeRoleConfig) (*stsCredential, string, error) {
	conf := &openapi.Config{
		AccessKeyId:     tea.String(accessKey),
		AccessKeySecret: tea.String(secretKey),
		RegionId:        tea.String(p.regionId),
		Endpoint:        tea.String(p.stsEndpoint),
		UserAgent:       tea.String(p.userAgent),
		// STS only supports HTTPS.
		Protocol:       tea.String("HTTPS"),
		ReadTimeout:    tea.Int(p.readTimeout),
		ConnectTimeout: tea.Int(p.connectTimeout),
		MaxIdleConns:   tea.Int(500),
	}
	if securityToken != "" {
		conf.SecurityToken = tea.String(securityToken)
	}
	query := map[string]*string{
		"AcceptLanguage": tea.String("en-US"),
	}
	if p.sourceIp != "" {
		query["SourceIp"] = tea.String(p.sourceIp)
	}
	if p.secureTransport != "" {
		query["SecureTransport"] = tea.String(p.secureTransport)
	}
	conf.GlobalParameters = &openapi.GlobalParameters{Queries: query}

	stsClient, err := sts20150401.NewClient(conf)
	if err != nil {
		return nil, "", fmt.Errorf("build sts client failed: %v", err)
	}

	request := &sts20150401.AssumeRoleRequest{
		RoleArn:         tea.String(hop.RoleArn),
		RoleSessionName: tea.String(hop.RoleSessionName),
	}
	if hop.Policy != "" {
		request.Policy = tea.String(hop.Policy)
	}
	if hop.ExternalId != "" {
		request.ExternalId = tea.String(hop.ExternalId)
	}
	if hop.SessionExpiration > 0 {
		request.DurationSeconds = tea.Int64(int64(hop.SessionExpiration))
	}

	runtime := &util.RuntimeOptions{}
	maxRetries := 5
	var response *sts20150401.AssumeRoleResponse
	for i := 0; i <= maxRetries; i++ {
		response, err = stsClient.AssumeRoleWithOptions(request, runtime)
		if err != nil {
			if needRetry(err) && i < maxRetries {
				time.Sleep(time.Duration(i) * time.Second)
				continue
			}
			return nil, "", err
		}
		break
	}
	if response == nil || response.Body == nil || response.Body.Credentials == nil {
		return nil, "", fmt.Errorf("AssumeRole returned empty credentials")
	}
	creds := response.Body.Credentials
	return &stsCredential{
		AccessKeyId:     tea.StringValue(creds.AccessKeyId),
		AccessKeySecret: tea.StringValue(creds.AccessKeySecret),
		SecurityToken:   tea.StringValue(creds.SecurityToken),
	}, tea.StringValue(creds.Expiration), nil
}

// GetAccessKeyId, GetAccessKeySecret, GetSecurityToken, GetBearerToken and
// GetType implement the deprecated portions of the credential.Credential
// interface by delegating to GetCredential so the chain stays consistent.

func (p *chainedAssumeRoleCredential) GetAccessKeyId() (*string, error) {
	c, err := p.GetCredential()
	if err != nil {
		return nil, err
	}
	return c.AccessKeyId, nil
}

func (p *chainedAssumeRoleCredential) GetAccessKeySecret() (*string, error) {
	c, err := p.GetCredential()
	if err != nil {
		return nil, err
	}
	return c.AccessKeySecret, nil
}

func (p *chainedAssumeRoleCredential) GetSecurityToken() (*string, error) {
	c, err := p.GetCredential()
	if err != nil {
		return nil, err
	}
	return c.SecurityToken, nil
}

func (p *chainedAssumeRoleCredential) GetBearerToken() *string {
	return tea.String("")
}

func (p *chainedAssumeRoleCredential) GetType() *string {
	return tea.String("ram_role_arn")
}
