package connectivity

import (
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	roa "github.com/alibabacloud-go/tea-roa/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"log"
	"regexp"
	"sync"
	"time"

	"encoding/json"
	"net/http"
	"strings"

	sts20150401 "github.com/alibabacloud-go/sts-20150401/v2/client"
	rpc "github.com/alibabacloud-go/tea-rpc/client"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	credential "github.com/aliyun/credentials-go/credentials"
	"github.com/jmespath/go-jmespath"
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
	Protocol             string
	ClientReadTimeout    int
	ClientConnectTimeout int
	SourceIp             string
	SecureTransport      string
	MaxRetryTimeout      int

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
	MseEndpoint                 string
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
	config := &openapi.Config{
		RegionId:        tea.String(c.RegionId),
		AccessKeyId:     tea.String(c.AccessKey),
		AccessKeySecret: tea.String(c.SecretKey),
		Endpoint:        tea.String(c.StsEndpoint),
		UserAgent:       tea.String(fmt.Sprintf("%s/%s %s/%s %s/%s %s/%s", Terraform, c.TerraformVersion, Provider, providerVersion, Module, c.ConfigurationSource, TerraformTraceId, c.TerraformTraceId)),
		// currently, sts endpoint only supports https
		Protocol:       tea.String("HTTPS"),
		ReadTimeout:    tea.Int(c.ClientReadTimeout),
		ConnectTimeout: tea.Int(c.ClientConnectTimeout),
		MaxIdleConns:   tea.Int(500),
	}
	if c.SecurityToken != "" {
		config.SecurityToken = tea.String(c.SecurityToken)
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
	config.GlobalParameters = param
	stsClient, err := sts20150401.NewClient(config)
	if err != nil {
		return fmt.Errorf("refreshing credential failed when building sts client. Error: %v", err)
	}

	request := &sts20150401.AssumeRoleRequest{
		RoleArn:         tea.String(c.RamRoleArn),
		RoleSessionName: tea.String(c.RamRoleSessionName),
	}
	if c.RamRolePolicy != "" {
		request.Policy = tea.String(c.RamRolePolicy)
	}
	if c.RamRoleSessionExpiration != 0 {
		request.DurationSeconds = tea.Int64(int64(c.RamRoleSessionExpiration))
	}
	if c.RamRoleExternalId != "" {
		request.ExternalId = tea.String(c.RamRoleExternalId)
	}

	runtime := &util.RuntimeOptions{}
	var response *sts20150401.AssumeRoleResponse
	maxRetries := 5
	for i := 0; i <= maxRetries; i++ {
		response, err = stsClient.AssumeRoleWithOptions(request, runtime)
		if err != nil {
			if needRetry(err) && i < maxRetries {
				time.Sleep(time.Duration(i))
				continue
			}
			return fmt.Errorf("refreshing credential failed by AssumeRole. Error: %v", err)
		}
		break
	}

	c.AccessKey, c.SecretKey, c.SecurityToken = *response.Body.Credentials.AccessKeyId, *response.Body.Credentials.AccessKeySecret, *response.Body.Credentials.SecurityToken
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
	requestUrl := securityCredURL + c.EcsRoleName
	httpRequest, err := http.NewRequest(requests.GET, requestUrl, strings.NewReader(""))
	if err != nil {
		err = fmt.Errorf("build sts requests err: %s", err.Error())
		return
	}
	httpClient := &http.Client{}
	httpResponse, err := httpClient.Do(httpRequest)
	if err != nil {
		err = fmt.Errorf("get Ecs sts token err : %s", err.Error())
		return
	}

	response := responses.NewCommonResponse()
	err = responses.Unmarshal(response, httpResponse, "")
	if err != nil {
		err = fmt.Errorf("Unmarshal Ecs sts token response err : %s", err.Error())
		return
	}

	if response.GetHttpStatus() != http.StatusOK {
		err = fmt.Errorf("get Ecs sts token err, httpStatus: %d, message = %s", response.GetHttpStatus(), response.GetHttpContentString())
		return
	}
	var data interface{}
	err = json.Unmarshal(response.GetHttpContentBytes(), &data)
	if err != nil {
		err = fmt.Errorf("refresh Ecs sts token err, json.Unmarshal fail: %s", err.Error())
		return
	}
	code, err := jmespath.Search("Code", data)
	if err != nil {
		err = fmt.Errorf("refresh Ecs sts token err, fail to get Code: %s", err.Error())
		return
	}
	if code.(string) != "Success" {
		err = fmt.Errorf("refresh Ecs sts token err, Code is not Success")
		return
	}
	accessKeyId, err := jmespath.Search("AccessKeyId", data)
	if err != nil || accessKeyId == nil {
		err = fmt.Errorf("refresh Ecs sts token err, fail to get AccessKeyId: %s", err.Error())
		return
	}
	accessKeySecret, err := jmespath.Search("AccessKeySecret", data)
	if err != nil || accessKeySecret == nil {
		err = fmt.Errorf("refresh Ecs sts token err, fail to get AccessKeySecret: %s", err.Error())
		return
	}
	securityToken, err := jmespath.Search("SecurityToken", data)
	if err != nil || securityToken == nil {
		err = fmt.Errorf("refresh Ecs sts token err, fail to get SecurityToken: %s", err.Error())
		return
	}

	c.AccessKey, c.SecretKey, c.SecurityToken = accessKeyId.(string), accessKeySecret.(string), securityToken.(string)
	return
}

// setAuthCredentialByOidc aims to access meta to get sts credential
// Actually, the job should be done by sdk, but currently not sdk support it, like alibaba-cloud-sdk-go.
// TODO：the provider can consider implementing all of credential type as an alternative to the SDK.
func (c *Config) setAuthCredentialByOidc() (err error) {
	if c.AccessKey != "" || c.AssumeRoleWithOidc == nil {
		return
	}
	config := &openapi.Config{
		RegionId:  tea.String(c.RegionId),
		Endpoint:  tea.String(c.StsEndpoint),
		UserAgent: tea.String(fmt.Sprintf("%s/%s %s/%s %s/%s %s/%s", Terraform, c.TerraformVersion, Provider, providerVersion, Module, c.ConfigurationSource, TerraformTraceId, c.TerraformTraceId)),
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
	config.GlobalParameters = param
	stsClient, err := sts20150401.NewClient(config)
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
			return fmt.Errorf("refreshing credential failed by AssumeRole. Error: %v", err)
		}
		break
	}

	c.AccessKey, c.SecretKey, c.SecurityToken = *response.Body.Credentials.AccessKeyId, *response.Body.Credentials.AccessKeySecret, *response.Body.Credentials.SecurityToken
	return nil
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
	config.SetUserAgent(fmt.Sprintf("%s/%s %s/%s %s/%s %s/%s", Terraform, c.TerraformVersion, Provider, providerVersion, Module, c.ConfigurationSource, TerraformTraceId, c.TerraformTraceId))
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
	config.SetUserAgent(fmt.Sprintf("%s/%s %s/%s %s/%s %s/%s", Terraform, c.TerraformVersion, Provider, providerVersion, Module, c.ConfigurationSource, TerraformTraceId, c.TerraformTraceId))
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
	config.SetUserAgent(fmt.Sprintf("%s/%s %s/%s %s/%s %s/%s", Terraform, c.TerraformVersion, Provider, providerVersion, Module, c.ConfigurationSource, TerraformTraceId, c.TerraformTraceId))
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
	config.SetUserAgent(fmt.Sprintf("%s/%s %s/%s %s/%s %s/%s", Terraform, c.TerraformVersion, Provider, providerVersion, Module, c.ConfigurationSource, TerraformTraceId, c.TerraformTraceId))
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
