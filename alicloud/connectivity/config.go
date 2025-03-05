package connectivity

import (
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	sts20150401 "github.com/alibabacloud-go/sts-20150401/v2/client"
	roa "github.com/alibabacloud-go/tea-roa/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"

	rpc "github.com/alibabacloud-go/tea-rpc/client"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	credential "github.com/aliyun/credentials-go/credentials"
)

var securityCredURL = "http://100.100.100.200/latest/meta-data/ram/security-credentials/"

// Config of aliyun
type Config struct {
	AccessKey            string
	SecretKey            string
	EcsRoleName          string
	Region               Region
	RegionId             string
	SecurityToken        string
	OtsInstanceName      string
	AccountId            string
	AccountType          string
	Protocol             string
	ClientReadTimeout    int
	ClientConnectTimeout int
	SourceIp             string
	SecureTransport      string
	MaxRetryTimeout      int
	Credential           credential.Credential

	RamRoleArn               string
	RamRoleSessionName       string
	RamRolePolicy            string
	RamRoleExternalId        string
	RamRoleSessionExpiration int
	AssumeRoleWithOidc       *AssumeRoleWithOidc
	Endpoints                *sync.Map
	SignVersion              *sync.Map
	RKvstoreEndpoint         string
	EcsEndpoint              string
	RdsEndpoint              string
	SlbEndpoint              string
	VpcEndpoint              string
	CenEndpoint              string
	EssEndpoint              string
	OssEndpoint              string
	OnsEndpoint              string
	AlikafkaEndpoint         string
	DnsEndpoint              string
	RamEndpoint              string
	CsEndpoint               string
	CrEndpoint               string
	CdnEndpoint              string
	KmsEndpoint              string
	OtsEndpoint              string
	CmsEndpoint              string
	PvtzEndpoint             string
	StsEndpoint              string
	LogEndpoint              string
	DrdsEndpoint             string
	DdsEndpoint              string
	GpdbEnpoint              string
	KVStoreEndpoint          string
	PolarDBEndpoint          string
	FcEndpoint               string
	ApigatewayEndpoint       string
	DatahubEndpoint          string
	MnsEndpoint              string
	LocationEndpoint         string
	ElasticsearchEndpoint    string
	NasEndpoint              string
	BssOpenApiEndpoint       string
	DdoscooEndpoint          string
	DdosbgpEndpoint          string
	SagEndpoint              string
	EmrEndpoint              string
	CasEndpoint              string
	MarketEndpoint           string
	HBaseEndpoint            string
	AdbEndpoint              string
	MaxComputeEndpoint       string

	edasEndpoint                string
	SkipRegionValidation        bool
	ConfigurationSource         string
	TerraformTraceId            string
	TerraformVersion            string
	CbnEndpoint                 string
	DmsEnterpriseEndpoint       string
	WafOpenapiEndpoint          string
	ResourcemanagerEndpoint     string
	BssopenapiEndpoint          string
	AlidnsEndpoint              string
	CassandraEndpoint           string
	EciEndpoint                 string
	OosEndpoint                 string
	DcdnEndpoint                string
	ActiontrailEndpoint         string
	ConfigEndpoint              string
	FnfEndpoint                 string
	RosEndpoint                 string
	PrivatelinkEndpoint         string
	MaxcomputeEndpoint          string
	ResourcesharingEndpoint     string
	GaEndpoint                  string
	HitsdbEndpoint              string
	BrainIndustrialEndpoint     string
	EipanycastEndpoint          string
	ImsEndpoint                 string
	QuotasEndpoint              string
	SgwEndpoint                 string
	ScdnEndpoint                string
	DmEndpoint                  string
	EventbridgeEndpoint         string
	OnsproxyEndpoint            string
	CdsEndpoint                 string
	HbrEndpoint                 string
	ArmsEndpoint                string
	CloudfwEndpoint             string
	ServerlessEndpoint          string
	AlbEndpoint                 string
	RedisaEndpoint              string
	GwsecdEndpoint              string
	CloudphoneEndpoint          string
	DataworkspublicEndpoint     string
	HcsSgwEndpoint              string
	CddcEndpoint                string
	MscopensubscriptionEndpoint string
	SddpEndpoint                string
	BastionhostEndpoint         string
	SasEndpoint                 string
	AlidfsEndpoint              string
	EhpcEndpoint                string
	EnsEndpoint                 string
	IotEndpoint                 string
	ImmEndpoint                 string
	ClickhouseEndpoint          string
	SelectDBEndpoint            string
	DtsEndpoint                 string
	DgEndpoint                  string
	CloudssoEndpoint            string
	WafEndpoint                 string
	SwasEndpoint                string
	VsEndpoint                  string
	QuickbiEndpoint             string
	VodEndpoint                 string
	OpensearchEndpoint          string
	GdsEndpoint                 string
	DbfsEndpoint                string
	DevopsrdcEndpoint           string
	EaisEndpoint                string
	CloudauthEndpoint           string
	ImpEndpoint                 string
	MhubEndpoint                string
	ServicemeshEndpoint         string
	AcrEndpoint                 string
	EdsuserEndpoint             string
	GpdbEndpoint                string
	GaplusEndpoint              string
	DdosbasicEndpoint           string
	SmartagEndpoint             string
	TagEndpoint                 string
	EdasEndpoint                string
	EdasschedulerxEndpoint      string
	EhsEndpoint                 string
	DysmsEndpoint               string
	CbsEndpoint                 string
	NlbEndpoint                 string
	VpcpeerEndpoint             string
	EbsEndpoint                 string
	DmsenterpriseEndpoint       string
	BpStudioEndpoint            string
	DasEndpoint                 string
	CloudfirewallEndpoint       string
	SrvcatalogEndpoint          string
	VpcPeerEndpoint             string
	EfloEndpoint                string
	OceanbaseEndpoint           string
	BeebotEndpoint              string
	ComputeNestEndpoint         string
}

type AssumeRoleWithOidc struct {
	RoleARN         string
	DurationSeconds int
	Policy          string
	RoleSessionName string
	OIDCProviderArn string
	OIDCTokenFile   string
	OIDCToken       string
}

func (c *Config) loadAndValidate() error {
	err := c.validateRegion()
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) validateRegion() error {

	for _, valid := range ValidRegions {
		if c.Region == valid {
			return nil
		}
	}

	return fmt.Errorf("Invalid Alibaba Cloud region: %s. "+
		"You can skip checking this region by setting provider parameter 'skip_region_validation'.", c.RegionId)
}

func (c *Config) getAuthCredential(stsSupported bool) auth.Credential {
	credential, err := c.Credential.GetCredential()
	if err == nil && credential != nil {
		c.AccessKey, c.SecretKey, c.SecurityToken = *credential.AccessKeyId, *credential.AccessKeySecret, *credential.SecurityToken
	}
	if c.AccessKey != "" && c.SecretKey != "" {
		if stsSupported && c.SecurityToken != "" {
			return credentials.NewStsTokenCredential(c.AccessKey, c.SecretKey, c.SecurityToken)
		}
		if c.RamRoleArn != "" {
			log.Printf("[INFO] Assume RAM Role specified in provider block assume_role { ... }")
			if c.RamRoleExternalId != "" {
				return credentials.NewRamRoleArnWithPolicyAndExternalIdCredential(
					c.AccessKey, c.SecretKey, c.RamRoleArn,
					c.RamRoleSessionName, c.RamRolePolicy, c.RamRoleExternalId, c.RamRoleSessionExpiration)
			}
			return credentials.NewRamRoleArnWithPolicyCredential(
				c.AccessKey, c.SecretKey, c.RamRoleArn,
				c.RamRoleSessionName, c.RamRolePolicy, c.RamRoleSessionExpiration)
		}
		return credentials.NewAccessKeyCredential(c.AccessKey, c.SecretKey)
	}
	if c.EcsRoleName != "" {
		return credentials.NewEcsRamRoleCredential(c.EcsRoleName)
	}

	return credentials.NewAccessKeyCredential(c.AccessKey, c.SecretKey)
}

func (c *Config) setAuthByAssumeRole() (err error) {
	if c.AccessKey == "" || c.RamRoleArn == "" {
		return
	}

	config := new(credential.Config).
		SetType("ram_role_arn").
		SetAccessKeyId(c.AccessKey).
		SetAccessKeySecret(c.SecretKey).
		SetRoleArn(c.RamRoleArn).
		SetRoleSessionName(c.RamRoleSessionName).
		SetPolicy(c.RamRolePolicy).
		SetExternalId(c.RamRoleExternalId).
		SetRoleSessionExpiration(c.RamRoleSessionExpiration)
	if c.SecurityToken != "" {
		config.SetSecurityToken(c.SecurityToken)
	}
	if c.StsEndpoint != "" {
		config.SetSTSEndpoint(c.StsEndpoint)
	}
	if c.ClientConnectTimeout != 0 {
		config.SetConnectTimeout(c.ClientConnectTimeout)
	}
	if c.ClientReadTimeout != 0 {
		config.SetTimeout(c.ClientReadTimeout)
	}
	provider, err := credential.NewCredential(config)
	if err != nil {
		return
	}
	c.Credential = provider
	credential, err := provider.GetCredential()
	if err != nil || credential == nil {
		return fmt.Errorf("refresh Ram Role Arn credential failed. Error: %v", err)
	}
	c.AccessKey, c.SecretKey, c.SecurityToken = *credential.AccessKeyId, *credential.AccessKeySecret, *credential.SecurityToken
	return nil
}

// setAuthCredentialByEcsRoleName aims to access meta to get sts credential
// Actually, the job should be done by sdk, but currently not all resources and products support alibaba-cloud-sdk-go,
// and their go sdk does support ecs role name.
// This method is a temporary solution and it should be removed after all go sdk support ecs role name
// The related PR: https://github.com/aliyun/terraform-provider-alicloud/pull/731
func (c *Config) setAuthCredentialByEcsRoleName() (err error) {
	if c.AccessKey != "" || c.EcsRoleName == "" {
		return
	}
	config := new(credential.Config).SetType("ecs_ram_role").SetRoleName(c.EcsRoleName)
	provider, err := credential.NewCredential(config)
	if err != nil {
		return
	}
	c.Credential = provider
	credential, err := provider.GetCredential()
	if err != nil || credential == nil {
		return fmt.Errorf("refresh Ecs Ram Role credential failed. Error: %v", err)
	}

	c.AccessKey, c.SecretKey, c.SecurityToken = *credential.AccessKeyId, *credential.AccessKeySecret, *credential.SecurityToken
	return
}

// setAuthCredentialByOidc aims to access meta to get sts credential
// Actually, the job should be done by sdk, but currently not sdk support it, like alibaba-cloud-sdk-go.
// TODO：the provider can consider implementing all of credential type as an alternative to the SDK.
func (c *Config) setAuthCredentialByOidc() (err error) {
	if c.AccessKey != "" || c.AssumeRoleWithOidc == nil {
		return
	}
	credConfig := new(credential.Config)
	if c.AssumeRoleWithOidc.OIDCToken == "" && c.AssumeRoleWithOidc.OIDCTokenFile != "" {
		credConfig.SetType("oidc_role_arn").
			SetOIDCProviderArn(c.AssumeRoleWithOidc.OIDCProviderArn).
			SetOIDCTokenFilePath(c.AssumeRoleWithOidc.OIDCTokenFile).
			SetRoleSessionName(c.AssumeRoleWithOidc.RoleSessionName).
			SetPolicy(c.AssumeRoleWithOidc.Policy).
			SetRoleArn(c.AssumeRoleWithOidc.RoleARN).
			SetSessionExpiration(c.AssumeRoleWithOidc.DurationSeconds)
		if c.StsEndpoint != "" {
			credConfig.SetSTSEndpoint(c.StsEndpoint)
		}
		if c.ClientConnectTimeout != 0 {
			credConfig.SetConnectTimeout(c.ClientConnectTimeout)
		}
		if c.ClientReadTimeout != 0 {
			credConfig.SetTimeout(c.ClientReadTimeout)
		}
	} else {
		conf := &openapi.Config{
			RegionId:  tea.String(c.RegionId),
			Endpoint:  tea.String(c.StsEndpoint),
			UserAgent: tea.String(c.getUserAgent()),
			// currently, sts endpoint only supports https
			Protocol:       tea.String("HTTPS"),
			ReadTimeout:    tea.Int(c.ClientReadTimeout),
			ConnectTimeout: tea.Int(c.ClientConnectTimeout),
			MaxIdleConns:   tea.Int(500),
		}
		query := map[string]*string{
			"AcceptLanguage": tea.String("en-US"),
		}
		if c.SourceIp != "" {
			query["SourceIp"] = tea.String(c.SourceIp)
		}
		if c.SecureTransport != "" {
			query["SecureTransport"] = tea.String(c.SecureTransport)
		}

		param := &openapi.GlobalParameters{Queries: query}
		conf.GlobalParameters = param
		stsClient, err := sts20150401.NewClient(conf)
		if err != nil {
			return fmt.Errorf("refreshing credential failed when building sts client. Error: %v", err)
		}

		request := &sts20150401.AssumeRoleWithOIDCRequest{
			OIDCProviderArn: tea.String(c.AssumeRoleWithOidc.OIDCProviderArn),
			RoleArn:         tea.String(c.AssumeRoleWithOidc.RoleARN),
			OIDCToken:       tea.String(c.AssumeRoleWithOidc.OIDCToken),
			RoleSessionName: tea.String(c.AssumeRoleWithOidc.RoleSessionName),
		}
		if c.AssumeRoleWithOidc.Policy != "" {
			request.Policy = tea.String(c.AssumeRoleWithOidc.Policy)
		}
		if c.AssumeRoleWithOidc.DurationSeconds != 0 {
			request.DurationSeconds = tea.Int64(int64(c.AssumeRoleWithOidc.DurationSeconds))
		}
		runtime := &util.RuntimeOptions{}
		var response *sts20150401.AssumeRoleWithOIDCResponse
		maxRetries := 5
		for i := 0; i <= maxRetries; i++ {
			response, err = stsClient.AssumeRoleWithOIDCWithOptions(request, runtime)
			if err != nil {
				if needRetry(err) && i < maxRetries {
					time.Sleep(time.Duration(i))
					continue
				}
				return fmt.Errorf("refreshing credential failed by AssumeRoleWithOIDC. Error: %v", err)
			}
			break
		}
		credConfig.SetType("sts").
			SetAccessKeyId(*response.Body.Credentials.AccessKeyId).
			SetAccessKeySecret(*response.Body.Credentials.AccessKeySecret).
			SetSecurityToken(*response.Body.Credentials.SecurityToken)
	}
	provider, err := credential.NewCredential(credConfig)
	if err != nil {
		return
	}
	c.Credential = provider
	credential, err := provider.GetCredential()
	if err != nil || credential == nil {
		return fmt.Errorf("refresh OIDC credential failed. Error: %v", err)
	}

	c.AccessKey, c.SecretKey, c.SecurityToken = *credential.AccessKeyId, *credential.AccessKeySecret, *credential.SecurityToken
	return
}
func needRetry(err error) bool {
	postRegex := regexp.MustCompile("^Post [\"]*https://.*")
	if postRegex.MatchString(err.Error()) {
		return true
	}

	throttlingRegex := regexp.MustCompile("Throttling")
	codeRegex := regexp.MustCompile("^code: 5[\\d]{2}")

	if e, ok := err.(*tea.SDKError); ok {
		if strings.Contains(*e.Message, "Client.Timeout") {
			return true
		}
		if *e.Code == "ServiceUnavailable" || *e.Code == "Rejected.Throttling" || throttlingRegex.MatchString(*e.Code) || codeRegex.MatchString(*e.Message) {
			return true
		}
	}
	return false
}
func isExpectedErrors(err error, expectCodes []string) bool {
	if err == nil {
		return false
	}

	if e, ok := err.(*tea.SDKError); ok {
		for _, code := range expectCodes {
			// The second statement aims to match the tea sdk history bug
			if *e.Code == code || strings.HasPrefix(code, *e.Code) || strings.Contains(*e.Data, code) {
				return true
			}
		}
		return false
	}

	for _, code := range expectCodes {
		if strings.Contains(err.Error(), code) {
			return true
		}
	}
	return false
}
func (c *Config) RefreshAuthCredential() error {
	if err := c.setAuthCredentialByEcsRoleName(); err != nil {
		return err
	}
	if err := c.setAuthCredentialByOidc(); err != nil {
		return err
	}
	return c.setAuthByAssumeRole()
}

func (c *Config) getTeaDslSdkConfig(stsSupported bool) (config rpc.Config, err error) {
	config.SetRegionId(c.RegionId)
	config.SetUserAgent(c.getUserAgent())
	credential, err := credential.NewCredential(c.getCredentialConfig(stsSupported))
	config.SetCredential(credential).
		SetRegionId(c.RegionId).
		SetProtocol(c.Protocol).
		SetReadTimeout(c.ClientReadTimeout).
		SetConnectTimeout(c.ClientConnectTimeout).
		SetMaxIdleConns(500)
	if c.SourceIp != "" {
		config.SetSourceIp(c.SourceIp)
	}
	if c.SecureTransport != "" {
		config.SetSecureTransport(c.SecureTransport)
	}

	return
}
func (c *Config) getTeaRoaDslSdkConfig(stsSupported bool) (config roa.Config, err error) {
	config.SetRegionId(c.RegionId)
	config.SetUserAgent(c.getUserAgent())
	credential, err := credential.NewCredential(c.getCredentialConfig(stsSupported))
	config.SetCredential(credential).
		SetRegionId(c.RegionId).
		SetProtocol(c.Protocol).
		SetReadTimeout(c.ClientReadTimeout).
		SetConnectTimeout(c.ClientConnectTimeout).
		SetMaxIdleConns(500)
	if c.SourceIp != "" {
		config.SetSourceIp(c.SourceIp)
	}
	if c.SecureTransport != "" {
		config.SetSecureTransport(c.SecureTransport)
	}
	return
}
func (c *Config) getTeaRpcOpenapiConfig(stsSupported bool) (config openapi.Config, err error) {
	config.SetRegionId(c.RegionId)
	config.SetUserAgent(c.getUserAgent())
	credential, err := credential.NewCredential(c.getCredentialConfig(stsSupported))
	config.SetCredential(credential).
		SetRegionId(c.RegionId).
		SetProtocol(c.Protocol).
		SetReadTimeout(c.ClientReadTimeout).
		SetConnectTimeout(c.ClientConnectTimeout).
		SetMaxIdleConns(500)

	query := map[string]*string{
		"AcceptLanguage": tea.String("en-US"),
	}
	if c.SourceIp != "" {
		query["SourceIp"] = tea.String(c.SourceIp)
	}
	if c.SecureTransport != "" {
		query["SecureTransport"] = tea.String(c.SecureTransport)
	}

	param := &openapi.GlobalParameters{Queries: query}
	config.GlobalParameters = param
	return
}
func (c *Config) getTeaRoaOpenapiConfig(stsSupported bool) (config openapi.Config, err error) {
	config.SetRegionId(c.RegionId)
	config.SetUserAgent(c.getUserAgent())
	credential, err := credential.NewCredential(c.getCredentialConfig(stsSupported))
	config.SetCredential(credential).
		SetRegionId(c.RegionId).
		SetProtocol(c.Protocol).
		SetReadTimeout(c.ClientReadTimeout).
		SetConnectTimeout(c.ClientConnectTimeout).
		SetMaxIdleConns(500)

	header := make(map[string]*string)
	if c.SourceIp != "" {
		header["x-acs-source-ip"] = tea.String(c.SourceIp)
	}
	if c.SecureTransport != "" {
		header["x-acs-secure-transport"] = tea.String(c.SecureTransport)
	}

	param := &openapi.GlobalParameters{Headers: header}
	config.GlobalParameters = param
	return
}
func (c *Config) getCredentialConfig(stsSupported bool) *credential.Config {
	credentialType := ""
	credentialConfig := &credential.Config{}
	if c.AccessKey != "" && c.SecretKey != "" {
		credentialType = "access_key"
		credentialConfig.AccessKeyId = &c.AccessKey     // AccessKeyId
		credentialConfig.AccessKeySecret = &c.SecretKey // AccessKeySecret

		if stsSupported && c.SecurityToken != "" {
			credentialType = "sts"
			credentialConfig.SecurityToken = &c.SecurityToken // STS Token
		} else if c.RamRoleArn != "" {
			log.Printf("[INFO] Assume RAM Role specified in provider block assume_role { ... }")
			credentialType = "ram_role_arn"
			credentialConfig.RoleArn = &c.RamRoleArn
			credentialConfig.RoleSessionName = &c.RamRoleSessionName
			credentialConfig.RoleSessionExpiration = &c.RamRoleSessionExpiration
			credentialConfig.Policy = &c.RamRolePolicy
			if c.RamRoleExternalId != "" {
				credentialConfig.ExternalId = &c.RamRoleExternalId
			}
		}
	} else if c.EcsRoleName != "" {
		credentialType = "ecs_ram_role"
		credentialConfig.RoleName = &c.EcsRoleName
	}

	credentialConfig.Type = &credentialType
	return credentialConfig
}

func (c *Config) getUserAgent() string {
	return fmt.Sprintf("%s/%s %s/%s %s/%s %s/%s", Terraform, c.TerraformVersion, Provider, providerVersion, Module, c.ConfigurationSource, TerraformTraceId, c.TerraformTraceId)
}

func (c *Config) needRefreshCredential() bool {
	credential, err := c.Credential.GetCredential()
	if err != nil || credential == nil {
		return false
	}

	return !(*credential.Type == "sts" || *credential.Type == "access_key")
}
func (c *Config) GetRefreshCredential() (string, string, string) {
	credential, err := c.Credential.GetCredential()
	if err != nil || credential == nil {
		log.Printf("[WARN] get credential failed. Error: %#v", err)
		return c.AccessKey, c.SecretKey, c.SecurityToken
	}

	return *credential.AccessKeyId, *credential.AccessKeySecret, *credential.SecurityToken
}
