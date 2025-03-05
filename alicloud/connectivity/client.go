package connectivity

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	util "github.com/alibabacloud-go/tea-utils/service"
	utilV2 "github.com/alibabacloud-go/tea-utils/v2/service"

	ossclient "github.com/alibabacloud-go/alibabacloud-gateway-oss/client"
	gatewayclient "github.com/alibabacloud-go/alibabacloud-gateway-sls/client"
	roaCS "github.com/alibabacloud-go/cs-20151215/v5/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	roa "github.com/alibabacloud-go/tea-roa/client"
	rpc "github.com/alibabacloud-go/tea-rpc/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/endpoints"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/adb"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alikafka"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/bssopenapi"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	cdn_new "github.com/aliyun/alibaba-cloud-sdk-go/services/cdn"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr_ee"
	officalCS "github.com/aliyun/alibaba-cloud-sdk-go/services/cs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ddosbgp"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ddoscoo"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dds"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/drds"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/elasticsearch"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/emr"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/gpdb"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/hbase"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/market"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/maxcompute"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ots"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	slsPop "github.com/aliyun/alibaba-cloud-sdk-go/services/sls"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/smartag"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/yundun_dbaudit"
	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	ali_mns "github.com/aliyun/aliyun-mns-go-sdk"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	otsTunnel "github.com/aliyun/aliyun-tablestore-go-sdk/tunnel"
	"github.com/aliyun/fc-go-sdk"
	"github.com/denverdino/aliyungo/cdn"
	"github.com/denverdino/aliyungo/cs"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cassandra"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dcdn"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/eci"
)

type AliyunClient struct {
	Region               Region
	RegionId             string
	SourceIp             string
	SecureTransport      string
	skipRegionValidation bool
	//In order to build ots table client, add accesskey and secretkey in aliyunclient temporarily.
	AccessKey                    string
	SecretKey                    string
	SecurityToken                string
	OtsInstanceName              string
	accountIdMutex               sync.RWMutex
	config                       *Config
	teaSdkConfig                 rpc.Config
	teaRoaSdkConfig              roa.Config
	teaRpcOpenapiConfig          openapi.Config
	teaRoaOpenapiConfig          openapi.Config
	accountId                    string
	ecsconn                      *ecs.Client
	essconn                      *ess.Client
	vpcconn                      *vpc.Client
	slbconn                      *slb.Client
	alikafkaconn                 *alikafka.Client
	ossconn                      *oss.Client
	dnsconn                      *alidns.Client
	ramconn                      *ram.Client
	csconn                       *cs.Client
	officalCSConn                *officalCS.Client
	roaCSConn                    *roaCS.Client
	cdnconn_new                  *cdn_new.Client
	crconn                       *cr.Client
	creeconn                     *cr_ee.Client
	cdnconn                      *cdn.CdnClient
	otsconn                      *ots.Client
	cmsconn                      *cms.Client
	logconn                      *sls.Client
	fcconn                       *fc.Client
	cenconn                      *cbn.Client
	logpopconn                   *slsPop.Client
	ddsconn                      *dds.Client
	gpdbconn                     *gpdb.Client
	stsconn                      *sts.Client
	rkvconn                      *r_kvstore.Client
	polarDBconn                  *polardb.Client
	dhconn                       datahub.DataHubApi
	mnsconn                      *ali_mns.MNSClient
	cloudapiconn                 *cloudapi.Client
	teaConn                      *rpc.Client
	tablestoreconnByInstanceName map[string]*tablestore.TableStoreClient
	otsTunnelConnByInstanceName  map[string]otsTunnel.TunnelClient
	csprojectconnByKey           map[string]*cs.ProjectClient
	drdsconn                     *drds.Client
	elasticsearchconn            *elasticsearch.Client
	ddoscooconn                  *ddoscoo.Client
	ddosbgpconn                  *ddosbgp.Client
	bssopenapiconn               *bssopenapi.Client
	emrconn                      *emr.Client
	sagconn                      *smartag.Client
	dbauditconn                  *yundun_dbaudit.Client
	marketconn                   *market.Client
	hbaseconn                    *hbase.Client
	adbconn                      *adb.Client
	cbnConn                      *cbn.Client
	maxcomputeconn               *maxcompute.Client
	dnsConn                      *alidns.Client
	edasconn                     *edas.Client
	bssopenapiConn               *bssopenapi.Client
	alidnsConn                   *alidns.Client
	ddoscooConn                  *ddoscoo.Client
	cassandraConn                *cassandra.Client
	eciConn                      *eci.Client
	ecsConn                      *ecs.Client
	dcdnConn                     *dcdn.Client
	cmsConn                      *cms.Client
	r_kvstoreConn                *r_kvstore.Client
	maxcomputeConn               *maxcompute.Client
}

type ApiVersion string

const (
	ApiVersion20140526 = ApiVersion("2014-05-26")
	ApiVersion20160815 = ApiVersion("2016-08-15")
	ApiVersion20140515 = ApiVersion("2014-05-15")
)

const businessInfoKey = "Terraform"

const DefaultClientRetryCountSmall = 5

const DefaultClientRetryCountMedium = 10

const DefaultClientRetryCountLarge = 15

const Terraform = "HashiCorp-Terraform"

const Provider = "Terraform-Provider"

const Module = "Terraform-Module"

const TerraformTraceId = "TerraformTraceId"

var goSdkMutex = sync.RWMutex{} // The Go SDK is not thread-safe
var loadSdkfromRemoteMutex = sync.Mutex{}
var loadSdkEndpointMutex = sync.Mutex{}

// The main version number that is being run at the moment.
var providerVersion = "1.244.0"

// Temporarily maintain map for old ecs client methods and store special endpoint information
var EndpointMap = map[string]string{
	"cn-shenzhen-su18-b01":        "ecs.aliyuncs.com",
	"cn-beijing":                  "ecs.aliyuncs.com",
	"cn-shenzhen-st4-d01":         "ecs.aliyuncs.com",
	"cn-haidian-cm12-c01":         "ecs.aliyuncs.com",
	"cn-hangzhou-internal-prod-1": "ecs.aliyuncs.com",
	"cn-qingdao":                  "ecs.aliyuncs.com",
	"cn-shanghai":                 "ecs.aliyuncs.com",
	"cn-shanghai-finance-1":       "ecs.aliyuncs.com",
	"cn-hongkong":                 "ecs.aliyuncs.com",
	"us-west-1":                   "ecs.aliyuncs.com",
	"cn-shenzhen":                 "ecs.aliyuncs.com",
	"cn-shanghai-et15-b01":        "ecs.aliyuncs.com",
	"cn-hangzhou-bj-b01":          "ecs.aliyuncs.com",
	"cn-zhangbei-na61-b01":        "ecs.aliyuncs.com",
	"cn-shenzhen-finance-1":       "ecs.aliyuncs.com",
	"cn-shanghai-et2-b01":         "ecs.aliyuncs.com",
	"ap-southeast-1":              "ecs.aliyuncs.com",
	"cn-beijing-nu16-b01":         "ecs.aliyuncs.com",
	"us-east-1":                   "ecs.aliyuncs.com",
	"cn-fujian":                   "ecs.aliyuncs.com",
	"cn-hangzhou":                 "ecs.aliyuncs.com",
}

// Client for AliyunClient
func (c *Config) Client() (*AliyunClient, error) {
	// Get the auth and region. This can fail if keys/regions were not
	// specified and we're attempting to use the environment.
	if !c.SkipRegionValidation {
		err := c.loadAndValidate()
		if err != nil {
			return nil, err
		}
	}
	loadLocalEndpoint = hasLocalEndpoint()
	if hasLocalEndpoint() {
		if err := c.loadEndpointFromLocal(); err != nil {
			return nil, err
		}
	}
	teaSdkConfig, err := c.getTeaDslSdkConfig(true)
	if err != nil {
		return nil, err
	}
	teaRoaSdkConfig, err := c.getTeaRoaDslSdkConfig(true)
	if err != nil {
		return nil, err
	}
	teaRpcOpenapiConfig, err := c.getTeaRpcOpenapiConfig(true)
	if err != nil {
		return nil, err
	}
	teaRoaOpenapiConfig, err := c.getTeaRoaOpenapiConfig(true)
	if err != nil {
		return nil, err
	}
	client := &AliyunClient{
		config:                       c,
		teaSdkConfig:                 teaSdkConfig,
		teaRoaSdkConfig:              teaRoaSdkConfig,
		teaRpcOpenapiConfig:          teaRpcOpenapiConfig,
		teaRoaOpenapiConfig:          teaRoaOpenapiConfig,
		SourceIp:                     c.SourceIp,
		Region:                       c.Region,
		RegionId:                     c.RegionId,
		AccessKey:                    c.AccessKey,
		SecretKey:                    c.SecretKey,
		SecurityToken:                c.SecurityToken,
		OtsInstanceName:              c.OtsInstanceName,
		accountId:                    c.AccountId,
		tablestoreconnByInstanceName: make(map[string]*tablestore.TableStoreClient),
		otsTunnelConnByInstanceName:  make(map[string]otsTunnel.TunnelClient),
		csprojectconnByKey:           make(map[string]*cs.ProjectClient),
		skipRegionValidation:         c.SkipRegionValidation,
	}
	if c.AccountType == "" {
		c.AccountType = client.getAccountType()
		client.config = c
	}
	log.Printf("[INFO] caller identity's account type is %s.", client.config.AccountType)
	return client, nil
}

func (client *AliyunClient) WithEcsClient(do func(*ecs.Client) (interface{}, error)) (interface{}, error) {
	if client.ecsconn != nil && !client.config.needRefreshCredential() {
		return do(client.ecsconn)
	}
	product := "ecs"
	endpoint, err := client.loadApiEndpoint(product)
	if err != nil {
		return nil, err
	}
	if endpoint != "" {
		endpoints.AddEndpointMapping(client.config.RegionId, product, endpoint)
	}
	ecsconn, err := ecs.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(time.Duration(60)*time.Second), client.config.getAuthCredential(true))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the ECS client: %#v", err)
	}
	ecs.SetClientProperty(ecsconn, "EndpointMap", map[string]string{
		client.RegionId: endpoint,
	})
	ecs.SetEndpointDataToClient(ecsconn)

	ecsconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
	ecsconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
	ecsconn.SourceIp = client.config.SourceIp
	ecsconn.SecureTransport = client.config.SecureTransport
	client.ecsconn = ecsconn
	return do(client.ecsconn)
}

func (client *AliyunClient) WithOfficalCSClient(do func(*officalCS.Client) (interface{}, error)) (interface{}, error) {
	if client.officalCSConn != nil && !client.config.needRefreshCredential() {
		return do(client.officalCSConn)
	}
	product := "cs"
	endpoint, err := client.loadApiEndpoint(product)
	if err != nil {
		return nil, err
	}
	if endpoint != "" {
		endpoints.AddEndpointMapping(client.config.RegionId, product, endpoint)
	}
	csconn, err := officalCS.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(0), client.config.getAuthCredential(true))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the CS client: %#v", err)
	}
	csconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
	csconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
	csconn.SourceIp = client.config.SourceIp
	csconn.SecureTransport = client.config.SecureTransport
	client.officalCSConn = csconn
	return do(client.officalCSConn)
}

func (client *AliyunClient) WithPolarDBClient(do func(*polardb.Client) (interface{}, error)) (interface{}, error) {
	if client.polarDBconn != nil && !client.config.needRefreshCredential() {
		return do(client.polarDBconn)
	}
	product := "polardb"
	endpoint, err := client.loadApiEndpoint(product)
	if err != nil {
		return nil, err
	}
	if endpoint != "" {
		endpoints.AddEndpointMapping(client.config.RegionId, product, endpoint)
	}
	polarDBconn, err := polardb.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(0), client.config.getAuthCredential(true))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the PolarDB client: %#v", err)

	}
	polarDBconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
	polarDBconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
	polarDBconn.SourceIp = client.config.SourceIp
	polarDBconn.SecureTransport = client.config.SecureTransport
	client.polarDBconn = polarDBconn
	return do(client.polarDBconn)
}

func (client *AliyunClient) WithSlbClient(do func(*slb.Client) (interface{}, error)) (interface{}, error) {
	if client.slbconn != nil && !client.config.needRefreshCredential() {
		return do(client.slbconn)
	}
	product := "slb"
	endpoint, err := client.loadApiEndpoint(product)
	if err != nil {
		return nil, err
	}
	if endpoint != "" {
		endpoints.AddEndpointMapping(client.config.RegionId, product, endpoint)
	}
	slbconn, err := slb.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(0), client.config.getAuthCredential(true))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the SLB client: %#v", err)
	}
	slbconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
	slbconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
	slbconn.SourceIp = client.config.SourceIp
	slbconn.SecureTransport = client.config.SecureTransport
	client.slbconn = slbconn
	return do(client.slbconn)
}

func (client *AliyunClient) WithVpcClient(do func(*vpc.Client) (interface{}, error)) (interface{}, error) {
	if client.vpcconn != nil && !client.config.needRefreshCredential() {
		return do(client.vpcconn)
	}
	product := "vpc"
	endpoint, err := client.loadApiEndpoint(product)
	if err != nil {
		return nil, err
	}
	if endpoint != "" {
		endpoints.AddEndpointMapping(client.config.RegionId, product, endpoint)
	}
	vpcconn, err := vpc.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(time.Duration(60)*time.Second), client.config.getAuthCredential(true))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the VPC client: %#v", err)
	}
	vpcconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
	vpcconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
	vpcconn.SourceIp = client.config.SourceIp
	vpcconn.SecureTransport = client.config.SecureTransport
	client.vpcconn = vpcconn
	return do(client.vpcconn)
}

func (client *AliyunClient) WithEssClient(do func(*ess.Client) (interface{}, error)) (interface{}, error) {
	if client.essconn != nil && !client.config.needRefreshCredential() {
		return do(client.essconn)
	}
	product := "ess"
	endpoint, err := client.loadApiEndpoint(product)
	if err != nil {
		return nil, err
	}
	if endpoint != "" {
		endpoints.AddEndpointMapping(client.config.RegionId, product, endpoint)
	}
	essconn, err := ess.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(0), client.config.getAuthCredential(true))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the ESS client: %#v", err)
	}
	essconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
	essconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
	essconn.SourceIp = client.config.SourceIp
	essconn.SecureTransport = client.config.SecureTransport
	client.essconn = essconn
	return do(client.essconn)
}

func (client *AliyunClient) WithOssClient(do func(*oss.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	if client.ossconn != nil && !client.config.needRefreshCredential() {
		return do(client.ossconn)
	}
	product := "oss"
	endpoint, err := client.loadApiEndpoint(product)
	if err != nil {
		return nil, err
	}
	schma := strings.ToLower(client.config.Protocol)
	if !strings.HasPrefix(endpoint, "http") {
		endpoint = fmt.Sprintf("%s://%s", schma, endpoint)
	}

	clientOptions := []oss.ClientOption{oss.UserAgent(client.config.getUserAgent())}
	proxy, err := client.getHttpProxy()
	if proxy != nil {
		skip, err := client.skipProxy(endpoint)
		if err != nil {
			return nil, err
		}
		if !skip {
			clientOptions = append(clientOptions, oss.Proxy(proxy.String()))
		}
	}

	clientOptions = append(clientOptions, oss.SetCredentialsProvider(&ossCredentialsProvider{client: client}))

	// region
	clientOptions = append(clientOptions, oss.Region(client.config.RegionId))

	// SignVersion
	if ossV, ok := client.config.SignVersion.Load("oss"); ok {
		clientOptions = append(clientOptions, oss.AuthVersion(func(v any) oss.AuthVersionType {
			switch fmt.Sprintf("%v", v) {
			case "v4":
				return oss.AuthV4
			case "v2":
				return oss.AuthV2
			}
			//default is v1
			return oss.AuthV1
		}(ossV)))
	}

	ossconn, err := oss.New(endpoint, "", "", clientOptions...)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the OSS client: %#v", err)
	}

	client.ossconn = ossconn
	return do(client.ossconn)
}

func (client *AliyunClient) WithOssBucketByName(bucketName string, do func(*oss.Bucket) (interface{}, error)) (interface{}, error) {
	return client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		bucket, err := client.ossconn.Bucket(bucketName)
		if err != nil {
			return nil, fmt.Errorf("unable to get the bucket %s: %#v", bucketName, err)
		}
		return do(bucket)
	})
}

func (client *AliyunClient) WithDnsClient(do func(*alidns.Client) (interface{}, error)) (interface{}, error) {
	if client.dnsconn != nil && !client.config.needRefreshCredential() {
		return do(client.dnsconn)
	}
	product := "alidns"
	endpoint, err := client.loadApiEndpoint(product)
	if err != nil {
		return nil, err
	}
	if endpoint != "" {
		endpoints.AddEndpointMapping(client.config.RegionId, string(DNSCode), endpoint)
	}

	dnsconn, err := alidns.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(0), client.config.getAuthCredential(true))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the DNS client: %#v", err)
	}
	dnsconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
	dnsconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
	dnsconn.SourceIp = client.config.SourceIp
	dnsconn.SecureTransport = client.config.SecureTransport
	client.dnsconn = dnsconn
	return do(client.dnsconn)
}

func (client *AliyunClient) WithRamClient(do func(*ram.Client) (interface{}, error)) (interface{}, error) {
	if client.ramconn != nil && !client.config.needRefreshCredential() {
		return do(client.ramconn)
	}
	product := "ram"
	endpoint, err := client.loadApiEndpoint(product)
	if err != nil {
		return nil, err
	}
	if endpoint != "" {
		endpoints.AddEndpointMapping(client.config.RegionId, product, endpoint)
	}
	ramconn, err := ram.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(0), client.config.getAuthCredential(true))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the RAM client: %#v", err)
	}
	ramconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
	ramconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
	ramconn.SourceIp = client.config.SourceIp
	ramconn.SecureTransport = client.config.SecureTransport
	client.ramconn = ramconn

	return do(client.ramconn)
}

func (client *AliyunClient) WithCsClient(do func(*cs.Client) (interface{}, error)) (interface{}, error) {
	if client.csconn != nil && !client.config.needRefreshCredential() {
		return do(client.csconn)
	}
	product := "cs"
	endpoint, err := client.loadApiEndpoint(product)
	if err != nil {
		return nil, err
	}
	endpoint = fmt.Sprintf("https://%s", endpoint)
	accessKey, secretKey, stsToken := client.config.AccessKey, client.config.SecretKey, client.config.SecurityToken
	credential, err := client.config.Credential.GetCredential()
	if err != nil || credential == nil {
		log.Printf("[WARN] get credential failed. Error: %#v", err)
	} else {
		accessKey, secretKey, stsToken = *credential.AccessKeyId, *credential.AccessKeySecret, *credential.SecurityToken
	}
	csconn := cs.NewClientForAussumeRole(accessKey, secretKey, stsToken)
	csconn.SetUserAgent(client.config.getUserAgent())
	csconn.SetEndpoint(endpoint)
	csconn.SetSourceIp(client.config.SourceIp)
	csconn.SetSecureTransport(client.config.SecureTransport)
	client.csconn = csconn
	return do(client.csconn)
}

func (client *AliyunClient) NewRoaCsClient() (*roaCS.Client, error) {
	product := "cs"
	endpoint, err := client.loadApiEndpoint(product)
	if err != nil {
		return nil, err
	}
	accessKey, secretKey, stsToken := client.config.AccessKey, client.config.SecretKey, client.config.SecurityToken
	credential, err := client.config.Credential.GetCredential()
	if err != nil || credential == nil {
		log.Printf("[WARN] get credential failed. Error: %#v", err)
	} else {
		accessKey, secretKey, stsToken = *credential.AccessKeyId, *credential.AccessKeySecret, *credential.SecurityToken
	}
	header := map[string]*string{
		"x-acs-source-ip":        tea.String(client.config.SourceIp),
		"x-acs-secure-transport": tea.String(client.config.SecureTransport),
	}
	param := &openapi.GlobalParameters{Headers: header}
	// Initialize the CS client if necessary
	roaCSConn, err := roaCS.NewClient(&openapi.Config{
		AccessKeyId:      tea.String(accessKey),
		AccessKeySecret:  tea.String(secretKey),
		SecurityToken:    tea.String(stsToken),
		RegionId:         tea.String(client.config.RegionId),
		UserAgent:        tea.String(client.config.getUserAgent()),
		Endpoint:         tea.String(endpoint),
		ReadTimeout:      tea.Int(client.config.ClientReadTimeout),
		ConnectTimeout:   tea.Int(client.config.ClientConnectTimeout),
		GlobalParameters: param,
	})
	if err != nil {
		return nil, err
	}

	return roaCSConn, nil
}

func (client *AliyunClient) WithCrClient(do func(*cr.Client) (interface{}, error)) (interface{}, error) {
	if client.crconn != nil && !client.config.needRefreshCredential() {
		return do(client.crconn)
	}
	product := "cr"
	endpoint, err := client.loadApiEndpoint(product)
	if err != nil {
		return nil, err
	}
	if endpoint != "" {
		endpoints.AddEndpointMapping(client.config.RegionId, string(CRCode), endpoint)
	}
	crconn, err := cr.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(0), client.config.getAuthCredential(true))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the CR client: %#v", err)
	}
	crconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
	crconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
	crconn.SourceIp = client.config.SourceIp
	crconn.SecureTransport = client.config.SecureTransport
	client.crconn = crconn

	return do(client.crconn)
}

func (client *AliyunClient) WithCrEEClient(do func(*cr_ee.Client) (interface{}, error)) (interface{}, error) {
	if client.creeconn != nil && !client.config.needRefreshCredential() {
		return do(client.creeconn)
	}
	product := "cr"
	endpoint, err := client.loadApiEndpoint(product)
	if err != nil {
		return nil, err
	}
	if endpoint != "" {
		endpoints.AddEndpointMapping(client.config.RegionId, string(CRCode), endpoint)
	}
	creeconn, err := cr_ee.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(0), client.config.getAuthCredential(true))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the CR EE client: %#v", err)
	}
	creeconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
	creeconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
	creeconn.SourceIp = client.config.SourceIp
	creeconn.SecureTransport = client.config.SecureTransport
	client.creeconn = creeconn

	return do(client.creeconn)
}

func (client *AliyunClient) WithCdnClient(do func(*cdn.CdnClient) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the CDN client if necessary
	if client.cdnconn == nil {
		cdnconn := cdn.NewClient(client.config.AccessKey, client.config.SecretKey)
		cdnconn.SetBusinessInfo(businessInfoKey)
		cdnconn.SetUserAgent(client.getUserAgent())
		cdnconn.SetSecurityToken(client.config.SecurityToken)
		endpoint := client.config.CdnEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, CDNCode)
		}
		if endpoint != "" && !strings.HasPrefix(endpoint, "http") {
			cdnconn.SetEndpoint(fmt.Sprintf("https://%s", strings.TrimPrefix(endpoint, "://")))
		}
		client.cdnconn = cdnconn
	}
	return do(client.cdnconn)
}

func (client *AliyunClient) WithCdnClient_new(do func(*cdn_new.Client) (interface{}, error)) (interface{}, error) {
	if client.cdnconn_new != nil && !client.config.needRefreshCredential() {
		return do(client.cdnconn_new)
	}
	product := "cdn"
	endpoint, err := client.loadApiEndpoint(product)
	if err != nil {
		return nil, err
	}
	if endpoint != "" {
		endpoints.AddEndpointMapping(client.config.RegionId, string(CDNCode), endpoint)
	}
	cdnconn, err := cdn_new.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(0), client.config.getAuthCredential(true))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the CDN client: %#v", err)
	}
	cdnconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
	cdnconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
	cdnconn.SourceIp = client.config.SourceIp
	cdnconn.SecureTransport = client.config.SecureTransport
	client.cdnconn_new = cdnconn

	return do(client.cdnconn_new)
}

// WithOtsClient init ots openapi publish sdk client(if necessary), and exec do func by client
func (client *AliyunClient) WithOtsClient(do func(*ots.Client) (interface{}, error)) (interface{}, error) {
	if client.otsconn != nil && !client.config.needRefreshCredential() {
		return do(client.otsconn)
	}
	product := "ots"
	endpoint, err := client.loadApiEndpoint(product)
	if err != nil {
		return nil, err
	}
	if endpoint != "" {
		endpoints.AddEndpointMapping(client.config.RegionId, product, endpoint)
	}
	otsconn, err := ots.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(0), client.config.getAuthCredential(true))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the OTS client: %#v", err)
	}

	otsconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
	otsconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
	otsconn.SourceIp = client.config.SourceIp
	otsconn.SecureTransport = client.config.SecureTransport
	client.otsconn = otsconn

	return do(client.otsconn)
}

// NewOtsRoaClient rpc client for common sdk
func (client *AliyunClient) NewOtsRoaClient(productCode string) (*roa.Client, error) {
	// first, load endpoint by user setting
	if v, ok := client.config.Endpoints.Load(productCode); !ok || v.(string) == "" {
		// second,  load endpoint by serverside rule
		if err := client.loadEndpoint(productCode); err != nil {
			return nil, err
		}
	}
	// set endpoint
	endpoint := ""
	if v, ok := client.config.Endpoints.Load(productCode); ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint", productCode)
	}
	sdkConfig := client.teaRoaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	accessKey, secretKey, securityToken := client.config.GetRefreshCredential()
	sdkConfig.SetAccessKeyId(accessKey)
	sdkConfig.SetAccessKeySecret(secretKey)
	sdkConfig.SetSecurityToken(securityToken)

	conn, err := roa.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}

	return conn, nil
}

func (client *AliyunClient) WithCmsClient(do func(*cms.Client) (interface{}, error)) (interface{}, error) {
	if client.cmsconn != nil && !client.config.needRefreshCredential() {
		return do(client.cmsconn)
	}
	product := "cms"
	endpoint, err := client.loadApiEndpoint(product)
	if err != nil {
		return nil, err
	}
	if endpoint != "" {
		endpoints.AddEndpointMapping(client.config.RegionId, product, endpoint)
	}
	cmsconn, err := cms.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(0), client.config.getAuthCredential(true))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the CMS client: %#v", err)
	}
	cmsconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
	cmsconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
	cmsconn.SourceIp = client.config.SourceIp
	cmsconn.SecureTransport = client.config.SecureTransport
	client.cmsconn = cmsconn

	return do(client.cmsconn)
}

func (client *AliyunClient) WithLogPopClient(do func(*slsPop.Client) (interface{}, error)) (interface{}, error) {
	if client.logpopconn != nil && !client.config.needRefreshCredential() {
		return do(client.logpopconn)
	}
	product := "sls"
	endpoint, err := client.loadApiEndpoint(product)
	if err != nil {
		return nil, err
	}

	logpopconn, err := slsPop.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(0), client.config.getAuthCredential(true))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the sls client: %#v", err)
	}
	logpopconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
	logpopconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
	logpopconn.SourceIp = client.config.SourceIp
	logpopconn.SecureTransport = client.config.SecureTransport
	endpoint = strings.TrimPrefix(strings.TrimPrefix(endpoint, "https://"), "http://")
	logpopconn.Domain = endpoint + "/open-api"
	client.logpopconn = logpopconn

	return do(client.logpopconn)
}

func (client *AliyunClient) WithLogClient(do func(*sls.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	if client.logconn != nil && !client.config.needRefreshCredential() {
		return do(client.logconn)
	}
	product := "sls"
	endpoint, err := client.loadApiEndpoint(product)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(endpoint, "http") {
		endpoint = fmt.Sprintf("https://%s", strings.TrimPrefix(endpoint, "://"))
	}
	accessKey, secretKey, securityToken := client.config.GetRefreshCredential()
	client.logconn = &sls.Client{
		AccessKeyID:     accessKey,
		AccessKeySecret: secretKey,
		Endpoint:        endpoint,
		SecurityToken:   securityToken,
		UserAgent:       client.getUserAgent(),
	}

	return do(client.logconn)
}

func (client *AliyunClient) WithDrdsClient(do func(*drds.Client) (interface{}, error)) (interface{}, error) {
	if client.drdsconn != nil && !client.config.needRefreshCredential() {
		return do(client.drdsconn)
	}
	product := "drds"
	endpoint, err := client.loadApiEndpoint(product)
	if err != nil {
		return nil, err
	}
	if endpoint != "" {
		endpoints.AddEndpointMapping(client.config.RegionId, product, endpoint)
	}
	drdsconn, err := drds.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(0), client.config.getAuthCredential(true))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the DRDS client: %#v", err)

	}
	drdsconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
	drdsconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
	drdsconn.SourceIp = client.config.SourceIp
	drdsconn.SecureTransport = client.config.SecureTransport
	client.drdsconn = drdsconn
	return do(client.drdsconn)
}

func (client *AliyunClient) WithDdsClient(do func(*dds.Client) (interface{}, error)) (interface{}, error) {
	if client.ddsconn != nil && !client.config.needRefreshCredential() {
		return do(client.ddsconn)
	}
	product := "dds"
	endpoint, err := client.loadApiEndpoint(product)
	if err != nil {
		return nil, err
	}
	if endpoint != "" {
		endpoints.AddEndpointMapping(client.config.RegionId, product, endpoint)
	}

	ddsconn, err := dds.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(0), client.config.getAuthCredential(true))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the mongoDB client: %#v", err)
	}
	ddsconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
	ddsconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
	ddsconn.SourceIp = client.config.SourceIp
	ddsconn.SecureTransport = client.config.SecureTransport
	client.ddsconn = ddsconn

	return do(client.ddsconn)
}

func (client *AliyunClient) WithGpdbClient(do func(*gpdb.Client) (interface{}, error)) (interface{}, error) {
	if client.gpdbconn != nil && !client.config.needRefreshCredential() {
		return do(client.gpdbconn)
	}
	product := "gpdb"
	endpoint, err := client.loadApiEndpoint(product)
	if err != nil {
		return nil, err
	}
	if endpoint != "" {
		endpoints.AddEndpointMapping(client.config.RegionId, product, endpoint)
	}

	gpdbconn, err := gpdb.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(0), client.config.getAuthCredential(true))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the GPDB client: %#v", err)
	}
	gpdbconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
	gpdbconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
	gpdbconn.SourceIp = client.config.SourceIp
	gpdbconn.SecureTransport = client.config.SecureTransport
	client.gpdbconn = gpdbconn
	return do(client.gpdbconn)
}

func (client *AliyunClient) WithFcClient(do func(*fc.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()
	if client.fcconn != nil && !client.config.needRefreshCredential() {
		return do(client.fcconn)
	}
	product := "fc"
	endpoint, err := client.loadApiEndpoint(product)
	if err != nil {
		return nil, err
	}
	if strings.HasPrefix(endpoint, "http") {
		endpoint = strings.TrimPrefix(strings.TrimPrefix(endpoint, "http://"), "https://")
	}

	config := client.getSdkConfig(0)
	transport := config.HttpTransport
	// Receiving proxy config from environment
	transport.Proxy = http.ProxyFromEnvironment
	clientOptions := []fc.ClientOption{fc.WithSecurityToken(client.config.SecurityToken), fc.WithTransport(transport),
		fc.WithTimeout(30), fc.WithRetryCount(DefaultClientRetryCountSmall)}

	accessKey, secretKey, secretToken := client.config.GetRefreshCredential()
	fcconn, err := fc.NewClient(fmt.Sprintf("https://%s", endpoint), string(ApiVersion20160815), accessKey, secretKey, clientOptions...)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the FC client: %#v", err)
	}

	fcconn.Config.UserAgent = client.config.getUserAgent()
	fcconn.Config.SecurityToken = secretToken
	client.fcconn = fcconn

	return do(client.fcconn)
}

func (client *AliyunClient) WithCloudApiClient(do func(*cloudapi.Client) (interface{}, error)) (interface{}, error) {
	if client.cloudapiconn != nil && !client.config.needRefreshCredential() {
		return do(client.cloudapiconn)
	}
	product := "cloudapi"
	endpoint, err := client.loadApiEndpoint(product)
	if err != nil {
		return nil, err
	}
	if endpoint != "" {
		endpoints.AddEndpointMapping(client.RegionId, product, endpoint)
	}
	cloudapiconn, err := cloudapi.NewClientWithOptions(client.RegionId, client.getSdkConfig(0), client.config.getAuthCredential(true))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the CloudAPI client: %#v", err)
	}
	cloudapiconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
	cloudapiconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
	cloudapiconn.SourceIp = client.config.SourceIp
	cloudapiconn.SecureTransport = client.config.SecureTransport
	client.cloudapiconn = cloudapiconn

	return do(client.cloudapiconn)
}

func (client *AliyunClient) NewTeaCommonClient(endpoint string) (*rpc.Client, error) {
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	credential, err := client.config.Credential.GetCredential()
	if err == nil && credential != nil {
		sdkConfig.SetAccessKeyId(*credential.AccessKeyId)
		sdkConfig.SetAccessKeySecret(*credential.AccessKeySecret)
		sdkConfig.SetSecurityToken(client.config.SecurityToken)
	} else {
		log.Printf("get credential error: %#v", err)
	}

	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the tea client: %#v", err)
	}

	return conn, nil
}
func (client *AliyunClient) NewTeaRoaCommonClient(endpoint string) (*roa.Client, error) {
	sdkConfig := client.teaRoaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	credential, err := client.config.Credential.GetCredential()
	if err == nil && credential != nil {
		sdkConfig.SetAccessKeyId(*credential.AccessKeyId)
		sdkConfig.SetAccessKeySecret(*credential.AccessKeySecret)
		sdkConfig.SetSecurityToken(client.config.SecurityToken)
	} else {
		log.Printf("get credential error: %#v", err)
	}

	conn, err := roa.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the tea roa client: %#v", err)
	}

	return conn, nil
}

func (client *AliyunClient) WithDataHubClient(do func(api datahub.DataHubApi) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the DataHub client if necessary
	if client.dhconn == nil {
		endpoint := client.config.DatahubEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.RegionId, DATAHUBCode)
		}
		if endpoint == "" {
			if client.RegionId == string(APSouthEast1) {
				endpoint = "dh-singapore.aliyuncs.com"
			} else {
				endpoint = fmt.Sprintf("dh-%s.aliyuncs.com", client.RegionId)
			}
		}
		if !strings.HasPrefix(endpoint, "http") {
			endpoint = fmt.Sprintf("https://%s", endpoint)
		}

		account := datahub.NewStsCredential(client.config.AccessKey, client.config.SecretKey, client.config.SecurityToken)
		config := &datahub.Config{
			UserAgent: client.getUserAgent(),
		}

		client.dhconn = datahub.NewClientWithConfig(endpoint, config, account)
	}

	return do(client.dhconn)
}

func (client *AliyunClient) WithElasticsearchClient(do func(*elasticsearch.Client) (interface{}, error)) (interface{}, error) {
	if client.elasticsearchconn != nil && !client.config.needRefreshCredential() {
		return do(client.elasticsearchconn)
	}
	product := "elasticsearch"
	endpoint, err := client.loadApiEndpoint(product)
	if err != nil {
		return nil, err
	}
	if endpoint != "" {
		endpoints.AddEndpointMapping(client.config.RegionId, product, endpoint)
	}
	elasticsearchconn, err := elasticsearch.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(0), client.config.getAuthCredential(true))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the Elasticsearch client: %#v", err)
	}
	elasticsearchconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
	elasticsearchconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
	elasticsearchconn.SourceIp = client.config.SourceIp
	elasticsearchconn.SecureTransport = client.config.SecureTransport
	client.elasticsearchconn = elasticsearchconn

	return do(client.elasticsearchconn)
}

func (client *AliyunClient) WithMnsClient(do func(*ali_mns.MNSClient) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the MNS client if necessary
	if client.mnsconn == nil {
		endpoint := client.config.MnsEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, MNSCode)
			if endpoint == "" {
				endpoint = fmt.Sprintf("%s.aliyuncs.com", client.config.RegionId)
			}
		}

		accountId, err := client.AccountId()
		if err != nil {
			return nil, err
		}
		if strings.HasPrefix(endpoint, "http") {
			endpoint = strings.TrimPrefix(strings.TrimPrefix(endpoint, "http://"), "https://")
		}
		mnsUrl := fmt.Sprintf("https://%s.mns.%s", accountId, endpoint)
		accessKey, secretKey, token := client.config.GetRefreshCredential()
		mnsClient := ali_mns.NewAliMNSClientWithToken(mnsUrl, accessKey, secretKey, token)
		proxy, err := client.getHttpProxy()
		if proxy != nil {
			skip, err := client.skipProxy(endpoint)
			if err != nil {
				return nil, err
			}
			if !skip {
				mnsClient.SetProxy(proxy.String())
			}
		}
		client.mnsconn = &mnsClient
	}

	return do(client.mnsconn)
}
func (client *AliyunClient) WithMnsQueueManager(do func(ali_mns.AliQueueManager) (interface{}, error)) (interface{}, error) {
	return client.WithMnsClient(func(mnsClient *ali_mns.MNSClient) (interface{}, error) {
		queueManager := ali_mns.NewMNSQueueManager(*mnsClient)
		return do(queueManager)
	})
}

func (client *AliyunClient) WithMnsTopicManager(do func(ali_mns.AliTopicManager) (interface{}, error)) (interface{}, error) {
	return client.WithMnsClient(func(mnsClient *ali_mns.MNSClient) (interface{}, error) {
		topicManager := ali_mns.NewMNSTopicManager(*mnsClient)
		return do(topicManager)
	})
}

func (client *AliyunClient) WithMnsSubscriptionManagerByTopicName(topicName string, do func(ali_mns.AliMNSTopic) (interface{}, error)) (interface{}, error) {
	return client.WithMnsClient(func(mnsClient *ali_mns.MNSClient) (interface{}, error) {
		subscriptionManager := ali_mns.NewMNSTopic(topicName, *mnsClient)
		return do(subscriptionManager)
	})
}

func (client *AliyunClient) WithTableStoreClient(instanceName string, do func(*tablestore.TableStoreClient) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the TABLESTORE client if necessary
	tableStoreClient, ok := client.tablestoreconnByInstanceName[instanceName]
	if ok && !client.config.needRefreshCredential() {
		return do(tableStoreClient)
	}
	endpoint := client.config.OtsEndpoint
	if endpoint == "" {
		endpoint = loadEndpoint(client.RegionId, OTSCode)
	}
	if endpoint == "" {
		endpoint = fmt.Sprintf("%s.%s.ots.aliyuncs.com", instanceName, client.RegionId)
	}
	if !strings.HasPrefix(endpoint, "https") && !strings.HasPrefix(endpoint, "http") {
		endpoint = fmt.Sprintf("https://%s", endpoint)
	}
	externalHeaders := make(map[string]string)
	if client.config.SecureTransport == "false" || client.config.SecureTransport == "true" {
		externalHeaders["x-ots-issecuretransport"] = client.config.SecureTransport
	}
	if client.config.SourceIp != "" {
		externalHeaders["x-ots-sourceip"] = client.config.SourceIp
	}
	accessKey, secretKey, token := client.config.GetRefreshCredential()
	tableStoreClient = tablestore.NewClientWithExternalHeader(endpoint, instanceName, accessKey, secretKey, token, tablestore.NewDefaultTableStoreConfig(), externalHeaders)
	client.tablestoreconnByInstanceName[instanceName] = tableStoreClient

	return do(tableStoreClient)
}

func (client *AliyunClient) WithTableStoreTunnelClient(instanceName string, do func(otsTunnel.TunnelClient) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the TABLESTORE tunnel client if necessary
	tunnelClient, ok := client.otsTunnelConnByInstanceName[instanceName]
	if ok && !client.config.needRefreshCredential() {
		return do(tunnelClient)
	}
	endpoint := client.config.OtsEndpoint
	if endpoint == "" {
		endpoint = loadEndpoint(client.RegionId, OTSCode)
	}
	if endpoint == "" {
		endpoint = fmt.Sprintf("%s.%s.ots.aliyuncs.com", instanceName, client.RegionId)
	}
	if !strings.HasPrefix(endpoint, "https") && !strings.HasPrefix(endpoint, "http") {
		endpoint = fmt.Sprintf("https://%s", endpoint)
	}

	externalHeaders := make(map[string]string)
	if client.config.SecureTransport == "false" || client.config.SecureTransport == "true" {
		externalHeaders["x-ots-issecuretransport"] = client.config.SecureTransport
	}
	if client.config.SourceIp != "" {
		externalHeaders["x-ots-sourceip"] = client.config.SourceIp
	}
	accessKey, secretKey, token := client.config.GetRefreshCredential()
	tunnelClient = otsTunnel.NewTunnelClientWithConfigAndExternalHeader(endpoint, instanceName, accessKey, secretKey, token, otsTunnel.DefaultTunnelConfig, externalHeaders)
	client.otsTunnelConnByInstanceName[instanceName] = tunnelClient

	return do(tunnelClient)
}

func (client *AliyunClient) WithCsProjectClient(clusterId, endpoint string, clusterCerts cs.ClusterCerts, do func(*cs.ProjectClient) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the PROJECT client if necessary
	key := fmt.Sprintf("%s|%s|%s|%s|%s", clusterId, endpoint, clusterCerts.CA, clusterCerts.Cert, clusterCerts.Key)
	csProjectClient, ok := client.csprojectconnByKey[key]
	if !ok {
		var err error
		csProjectClient, err = cs.NewProjectClient(clusterId, endpoint, clusterCerts)
		if err != nil {
			return nil, fmt.Errorf("Getting Application Client failed by cluster id %s: %#v.", clusterCerts, err)
		}
		csProjectClient.SetDebug(false)
		csProjectClient.SetUserAgent(client.config.getUserAgent())
		client.csprojectconnByKey[key] = csProjectClient
	}

	return do(csProjectClient)
}

func (client *AliyunClient) NewCommonRequest(product, serviceCode, schema string, apiVersion ApiVersion) (*requests.CommonRequest, error) {
	endpoint := ""
	product = strings.ToLower(product)
	if _, exist := client.config.Endpoints.Load(product); !exist {
		if err := client.loadEndpoint(product); err != nil {
			return nil, err
		}
	}
	if v, exist := client.config.Endpoints.Load(product); exist && v.(string) != "" {
		endpoint = v.(string)
	}
	request := requests.NewCommonRequest()
	// Use product code to find product domain
	if endpoint != "" {
		request.Domain = endpoint
	} else {
		// When getting endpoint failed by location, using custom endpoint instead
		request.Domain = fmt.Sprintf("%s.%s.aliyuncs.com", strings.ToLower(serviceCode), client.RegionId)
	}
	request.Version = string(apiVersion)
	request.RegionId = client.RegionId
	request.Product = product
	request.Scheme = schema
	request.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
	request.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
	request.AppendUserAgent(Terraform, client.config.TerraformVersion)
	request.AppendUserAgent(Provider, providerVersion)
	request.AppendUserAgent(Module, client.config.ConfigurationSource)
	request.AppendUserAgent(TerraformTraceId, client.config.TerraformTraceId)
	return request, nil
}

func (client *AliyunClient) AccountId() (string, error) {
	if client.accountId != "" {
		return client.accountId, nil
	}

	client.accountIdMutex.Lock()
	defer client.accountIdMutex.Unlock()

	log.Printf("[DEBUG] account_id not provided, attempting to retrieve it automatically...")
	identity, err := client.GetCallerIdentity()
	if err != nil {
		return "", err
	}
	if identity.AccountId == "" {
		return "", fmt.Errorf("caller identity doesn't contain any AccountId")
	}
	client.accountId = identity.AccountId
	return client.accountId, nil
}

// getAccountType determines and returns the account type (Domestic or International) based on the client's configuration and API endpoint.
// This function first checks if the AccountType is already set in the client configuration. If so, it returns that value directly.
// Otherwise, it defaults the account type to "Domestic" and initializes a request to query available instances through the BssOpenApi API.
// It then determines whether the account is domestic or international based on the API specific errors.
// If there is a specific error, the account type should be updates, and meantime corrects the BssOpenApi endpoint.
func (client *AliyunClient) getAccountType() string {
	if client.config.AccountType != "" {
		return client.config.AccountType
	}
	// Default to Domestic
	accountType := "Domestic"
	productCode := strings.ToLower("BssOpenApi")
	request := map[string]interface{}{
		"PageSize":         "50",
		"PageNum":          1,
		"ProductCode":      "vipcloudfw",
		"ProductType":      "vipcloudfw",
		"SubscriptionType": "Subscription",
	}
	endpoint, err := client.loadApiEndpoint(productCode)
	if err != nil {
		log.Printf("[WARN] getting BssOpenApi endpoint failed. Error: %v", err)
	} else if endpoint == BssOpenAPIEndpointInternational {
		request["ProductCode"] = "cfw"
		request["ProductType"] = "cfw_pre_intl"
		accountType = "International"
	}
	wait := incrementalWait(1*time.Second, 0*time.Second)
	resource.Retry(30*time.Second, func() *resource.RetryError {
		_, err := client.RpcPost("BssOpenApi", "2017-12-14", "QueryAvailableInstances", nil, request, true)
		log.Printf("[WARN] checking caller identity's account type by invoking BssOpenApi QueryAvailableInstances failed. Error: %v", err)
		if err != nil {
			if needRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if isExpectedErrors(err, []string{"NotApplicable", "not found article by given param"}) {
				if request["ProductType"] == "vipcloudfw" {
					accountType = "International"
					client.config.Endpoints.Store(productCode, BssOpenAPIEndpointInternational)
				} else {
					accountType = "Domestic"
					client.config.Endpoints.Store(productCode, BssOpenAPIEndpointDomestic)
				}
			} else {
				accountType = ""
			}
		}
		return nil
	})
	return accountType
}
func (client *AliyunClient) IsInternationalAccount() bool {
	return client.config.AccountType == "International"
}
func (client *AliyunClient) isInternationalRegion() bool {
	if client.config.RegionId == "cn-hongkong" {
		return true
	}
	return !strings.HasPrefix(client.config.RegionId, "cn-")
}

func (client *AliyunClient) getSdkConfig(timeout time.Duration) *sdk.Config {
	if timeout == 0 {
		timeout = time.Duration(30) * time.Second
	}
	// WithUserAgent will add a prefix Extra/ for user agent value
	return sdk.NewConfig().
		WithMaxRetryTime(DefaultClientRetryCountSmall).
		WithTimeout(timeout).
		WithEnableAsync(false).
		WithGoRoutinePoolSize(100).
		WithMaxTaskQueueSize(10000).
		WithDebug(false).
		WithHttpTransport(client.getTransport()).
		WithScheme(client.config.Protocol).
		WithUserAgent(fmt.Sprintf("Terraform %s", client.config.getUserAgent()))
}

func (client *AliyunClient) getUserAgent() string {
	return fmt.Sprintf("%s/%s %s/%s %s/%s", Terraform, client.config.TerraformVersion, Provider, providerVersion, Module, client.config.ConfigurationSource)
}

func (client *AliyunClient) getTransport() *http.Transport {
	handshakeTimeout, err := strconv.Atoi(os.Getenv("TLSHandshakeTimeout"))
	if err != nil {
		handshakeTimeout = 120
	}
	transport := &http.Transport{}
	transport.TLSHandshakeTimeout = time.Duration(handshakeTimeout) * time.Second

	return transport
}

func (client *AliyunClient) getHttpProxy() (proxy *url.URL, err error) {
	if client.config.Protocol == "HTTPS" {
		if rawurl := os.Getenv("HTTPS_PROXY"); rawurl != "" {
			proxy, err = url.Parse(rawurl)
		} else if rawurl := os.Getenv("https_proxy"); rawurl != "" {
			proxy, err = url.Parse(rawurl)
		}
	} else {
		if rawurl := os.Getenv("HTTP_PROXY"); rawurl != "" {
			proxy, err = url.Parse(rawurl)
		} else if rawurl := os.Getenv("http_proxy"); rawurl != "" {
			proxy, err = url.Parse(rawurl)
		}
	}
	return proxy, err
}

func (client *AliyunClient) skipProxy(endpoint string) (bool, error) {
	var urls []string
	if rawurl := os.Getenv("NO_PROXY"); rawurl != "" {
		urls = strings.Split(rawurl, ",")
	} else if rawurl := os.Getenv("no_proxy"); rawurl != "" {
		urls = strings.Split(rawurl, ",")
	}
	for _, value := range urls {
		if strings.HasPrefix(value, "*") {
			value = fmt.Sprintf(".%s", value)
		}
		noProxyReg, err := regexp.Compile(value)
		if err != nil {
			return false, err
		}
		if noProxyReg.MatchString(endpoint) {
			return true, nil
		}
	}
	return false, nil
}

func (client *AliyunClient) GetCallerIdentity() (*sts.GetCallerIdentityResponse, error) {
	args := sts.CreateGetCallerIdentityRequest()

	endpoint := client.config.StsEndpoint
	if endpoint == "" {
		endpoint = loadEndpoint(client.config.RegionId, STSCode)
	}
	if endpoint != "" {
		endpoints.AddEndpointMapping(client.config.RegionId, string(STSCode), endpoint)
	}
	stsClient, err := sts.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(0), client.config.getAuthCredential(true))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the STS client: %#v", err)
	}

	stsClient.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
	stsClient.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
	stsClient.SourceIp = client.config.SourceIp
	stsClient.SecureTransport = client.config.SecureTransport

	identity, err := stsClient.GetCallerIdentity(args)
	if err != nil {
		return nil, err
	}
	if identity == nil {
		return nil, fmt.Errorf("caller identity not found")
	}
	return identity, err
}
func (client *AliyunClient) WithDdosbgpClient(do func(*ddosbgp.Client) (interface{}, error)) (interface{}, error) {
	if client.ddosbgpconn != nil && !client.config.needRefreshCredential() {
		return do(client.ddosbgpconn)
	}
	product := "ddosbgp"
	endpoint, err := client.loadApiEndpoint(product)
	if err != nil {
		return nil, err
	}
	if endpoint != "" {
		endpoints.AddEndpointMapping(client.config.RegionId, product, endpoint)
	}
	ddosbgpconn, err := ddosbgp.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(0), client.config.getAuthCredential(true))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the DDOSBGP client: %#v", err)
	}
	ddosbgpconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
	ddosbgpconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
	ddosbgpconn.SourceIp = client.config.SourceIp
	ddosbgpconn.SecureTransport = client.config.SecureTransport
	client.ddosbgpconn = ddosbgpconn

	return do(client.ddosbgpconn)
}
func (client *AliyunClient) WithAlikafkaClient(do func(*alikafka.Client) (interface{}, error)) (interface{}, error) {
	if client.alikafkaconn != nil && !client.config.needRefreshCredential() {
		return do(client.alikafkaconn)
	}
	product := "alikafka"
	endpoint, err := client.loadApiEndpoint(product)
	if err != nil {
		return nil, err
	}
	if endpoint != "" {
		endpoints.AddEndpointMapping(client.config.RegionId, product, endpoint)
	}
	alikafkaconn, err := alikafka.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(0), client.config.getAuthCredential(true))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the ALIKAFKA client: %#v", err)
	}
	alikafkaconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
	alikafkaconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
	alikafkaconn.SourceIp = client.config.SourceIp
	alikafkaconn.SecureTransport = client.config.SecureTransport
	client.alikafkaconn = alikafkaconn

	return do(client.alikafkaconn)
}

func (client *AliyunClient) WithEmrClient(do func(*emr.Client) (interface{}, error)) (interface{}, error) {
	if client.emrconn != nil && !client.config.needRefreshCredential() {
		return do(client.emrconn)
	}
	product := "emr"
	endpoint, err := client.loadApiEndpoint(product)
	if err != nil {
		return nil, err
	}
	if endpoint != "" {
		endpoints.AddEndpointMapping(client.config.RegionId, product, endpoint)
	}
	emrConn, err := emr.NewClientWithOptions(client.RegionId, client.getSdkConfig(0), client.config.getAuthCredential(true))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the E-MapReduce client: %#v", err)
	}
	emrConn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
	emrConn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
	emrConn.SourceIp = client.config.SourceIp
	emrConn.SecureTransport = client.config.SecureTransport
	client.emrconn = emrConn

	return do(client.emrconn)
}

func (client *AliyunClient) WithSagClient(do func(*smartag.Client) (interface{}, error)) (interface{}, error) {
	if client.sagconn != nil && !client.config.needRefreshCredential() {
		return do(client.sagconn)
	}
	product := "smartag"
	endpoint, err := client.loadApiEndpoint(product)
	if err != nil {
		return nil, err
	}
	if endpoint != "" {
		endpoints.AddEndpointMapping(client.config.RegionId, product, endpoint)
	}
	sagconn, err := smartag.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(0), client.config.getAuthCredential(true))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the SAG client: %#v", err)
	}
	sagconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
	sagconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
	sagconn.SourceIp = client.config.SourceIp
	sagconn.SecureTransport = client.config.SecureTransport
	client.sagconn = sagconn

	return do(client.sagconn)
}

func (client *AliyunClient) WithDbauditClient(do func(*yundun_dbaudit.Client) (interface{}, error)) (interface{}, error) {
	if client.dbauditconn != nil && !client.config.needRefreshCredential() {
		return do(client.dbauditconn)
	}
	product := "yundun_dbaudit"
	endpoint, err := client.loadApiEndpoint(product)
	if err != nil {
		return nil, err
	}
	if endpoint != "" {
		endpoints.AddEndpointMapping(client.config.RegionId, product, endpoint)
	}
	dbauditconn, err := yundun_dbaudit.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(0), client.config.getAuthCredential(true))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the DBAUDIT client: %#v", err)
	}

	dbauditconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
	dbauditconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
	dbauditconn.SourceIp = client.config.SourceIp
	dbauditconn.SecureTransport = client.config.SecureTransport
	client.dbauditconn = dbauditconn

	return do(client.dbauditconn)
}
func (client *AliyunClient) WithMarketClient(do func(*market.Client) (interface{}, error)) (interface{}, error) {
	if client.marketconn != nil && !client.config.needRefreshCredential() {
		return do(client.marketconn)
	}
	product := "market"
	endpoint, err := client.loadApiEndpoint(product)
	if err != nil {
		return nil, err
	}
	if endpoint != "" {
		endpoints.AddEndpointMapping(client.config.RegionId, product, endpoint)
	}
	marketconn, err := market.NewClientWithOptions(client.RegionId, client.getSdkConfig(0), client.config.getAuthCredential(true))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the Market client: %#v", err)
	}

	marketconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
	marketconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
	marketconn.SourceIp = client.config.SourceIp
	marketconn.SecureTransport = client.config.SecureTransport
	client.marketconn = marketconn

	return do(client.marketconn)
}

func (client *AliyunClient) WithHbaseClient(do func(*hbase.Client) (interface{}, error)) (interface{}, error) {
	if client.hbaseconn != nil && !client.config.needRefreshCredential() {
		return do(client.hbaseconn)
	}
	product := "hbase"
	endpoint, err := client.loadApiEndpoint(product)
	if err != nil {
		return nil, err
	}
	if endpoint != "" {
		endpoints.AddEndpointMapping(client.config.RegionId, product, endpoint)
	}
	hbaseconn, err := hbase.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(0), client.config.getAuthCredential(true))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the hbase client: %#v", err)
	}
	hbaseconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
	hbaseconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
	hbaseconn.SourceIp = client.config.SourceIp
	hbaseconn.SecureTransport = client.config.SecureTransport

	client.hbaseconn = hbaseconn

	return do(client.hbaseconn)
}

func (client *AliyunClient) WithAdbClient(do func(*adb.Client) (interface{}, error)) (interface{}, error) {
	if client.adbconn != nil && !client.config.needRefreshCredential() {
		return do(client.adbconn)
	}
	product := "adb"
	endpoint, err := client.loadApiEndpoint(product)
	if err != nil {
		return nil, err
	}
	if endpoint != "" {
		endpoints.AddEndpointMapping(client.config.RegionId, product, endpoint)
	}
	adbconn, err := adb.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(0), client.config.getAuthCredential(true))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the adb client: %#v", err)

	}
	adbconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
	adbconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
	adbconn.SourceIp = client.config.SourceIp
	adbconn.SecureTransport = client.config.SecureTransport
	client.adbconn = adbconn

	return do(client.adbconn)
}
func (client *AliyunClient) WithCbnClient(do func(*cbn.Client) (interface{}, error)) (interface{}, error) {
	if client.cbnConn != nil && !client.config.needRefreshCredential() {
		return do(client.cbnConn)
	}
	product := "cbn"
	endpoint, err := client.loadApiEndpoint(product)
	if err != nil {
		return nil, err
	}
	if endpoint != "" {
		endpoints.AddEndpointMapping(client.config.RegionId, product, endpoint)
	}
	cbnConn, err := cbn.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(0), client.config.getAuthCredential(true))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the Cbnclient: %#v", err)
	}

	cbnConn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
	cbnConn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
	cbnConn.SourceIp = client.config.SourceIp
	cbnConn.SecureTransport = client.config.SecureTransport
	client.cbnConn = cbnConn

	return do(client.cbnConn)
}

func (client *AliyunClient) WithEdasClient(do func(*edas.Client) (interface{}, error)) (interface{}, error) {
	if client.edasconn != nil && !client.config.needRefreshCredential() {
		return do(client.edasconn)
	}
	product := "edas"
	endpoint, err := client.loadApiEndpoint(product)
	if err != nil {
		return nil, err
	}
	if endpoint != "" {
		endpoints.AddEndpointMapping(client.config.RegionId, product, endpoint)
	}
	edasconn, err := edas.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(time.Duration(60)*time.Second), client.config.getAuthCredential(true))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the EDAS client: %#v", err)
	}
	edasconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
	edasconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
	edasconn.SourceIp = client.config.SourceIp
	edasconn.SecureTransport = client.config.SecureTransport
	client.edasconn = edasconn

	return do(client.edasconn)
}

func (client *AliyunClient) WithAlidnsClient(do func(*alidns.Client) (interface{}, error)) (interface{}, error) {
	if client.alidnsConn != nil && !client.config.needRefreshCredential() {
		return do(client.alidnsConn)
	}
	product := "alidns"
	endpoint, err := client.loadApiEndpoint(product)
	if err != nil {
		return nil, err
	}

	if endpoint != "" {
		endpoints.AddEndpointMapping(client.config.RegionId, product, endpoint)
	}

	alidnsConn, err := alidns.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(0), client.config.getAuthCredential(true))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the Alidnsclient: %#v", err)
	}
	alidnsConn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
	alidnsConn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
	alidnsConn.SourceIp = client.config.SourceIp
	alidnsConn.SecureTransport = client.config.SecureTransport
	client.alidnsConn = alidnsConn
	return do(client.alidnsConn)
}

func (client *AliyunClient) WithCassandraClient(do func(*cassandra.Client) (interface{}, error)) (interface{}, error) {
	if client.cassandraConn != nil && !client.config.needRefreshCredential() {
		return do(client.cassandraConn)
	}
	product := "cassandra"
	endpoint, err := client.loadApiEndpoint(product)
	if err != nil {
		return nil, err
	}
	if endpoint != "" {
		endpoints.AddEndpointMapping(client.config.RegionId, product, endpoint)
	}
	cassandraConn, err := cassandra.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(0), client.config.getAuthCredential(true))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the Cassandraclient: %#v", err)
	}
	cassandraConn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
	cassandraConn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
	cassandraConn.SourceIp = client.config.SourceIp
	cassandraConn.SecureTransport = client.config.SecureTransport
	client.cassandraConn = cassandraConn
	return do(client.cassandraConn)
}

func (client *AliyunClient) WithEciClient(do func(*eci.Client) (interface{}, error)) (interface{}, error) {
	if client.eciConn != nil && !client.config.needRefreshCredential() {
		return do(client.eciConn)
	}
	product := "eci"
	endpoint, err := client.loadApiEndpoint(product)
	if err != nil {
		return nil, err
	}

	if endpoint != "" {
		endpoints.AddEndpointMapping(client.config.RegionId, string(EciCode), endpoint)
	}

	eciConn, err := eci.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(0), client.config.getAuthCredential(true))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the Eciclient: %#v", err)
	}
	eciConn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
	eciConn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
	eciConn.SourceIp = client.config.SourceIp
	eciConn.SecureTransport = client.config.SecureTransport
	client.eciConn = eciConn
	return do(client.eciConn)
}
func (client *AliyunClient) WithRKvstoreClient(do func(*r_kvstore.Client) (interface{}, error)) (interface{}, error) {
	if client.r_kvstoreConn != nil && !client.config.needRefreshCredential() {
		return do(client.r_kvstoreConn)
	}
	product := "r_kvstore"
	endpoint, err := client.loadApiEndpoint(product)
	if err != nil {
		return nil, err
	}
	if endpoint != "" {
		endpoints.AddEndpointMapping(client.config.RegionId, "r-kvstore", endpoint)
	}
	r_kvstoreConn, err := r_kvstore.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(0), client.config.getAuthCredential(true))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the RKvstoreclient: %#v", err)
	}

	r_kvstoreConn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
	r_kvstoreConn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
	r_kvstoreConn.SourceIp = client.config.SourceIp
	r_kvstoreConn.SecureTransport = client.config.SecureTransport
	client.r_kvstoreConn = r_kvstoreConn

	return do(client.r_kvstoreConn)
}

func (client *AliyunClient) NewQuotasClientV2() (*openapi.Client, error) {
	productCode := "quotas"
	endpoint := ""
	if v, ok := client.config.Endpoints.Load(productCode); !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			endpoint = "quotas.aliyuncs.com"
			client.config.Endpoints.Store(productCode, endpoint)
			log.Printf("[ERROR] loading %s endpoint got an error: %#v. Using the central endpoint %s instead.", productCode, err, endpoint)
		}
	}
	if v, ok := client.config.Endpoints.Load(productCode); ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	openapiConfig := client.teaRpcOpenapiConfig
	openapiConfig.Endpoint = tea.String(endpoint)
	openapiConfig.Protocol = client.teaRpcOpenapiConfig.Protocol
	result, err := openapi.NewClient(&openapiConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return result, nil
}
func (client *AliyunClient) loadApiEndpoint(productCode string) (string, error) {
	if v, ok := client.config.Endpoints.Load(productCode); !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			return "", fmt.Errorf("[ERROR] loading %s endpoint got an error: %#v.", productCode, err)
		}
	} else {
		return client.formatEndpointWithAccountID(productCode, v.(string))
	}
	if v, ok := client.config.Endpoints.Load(productCode); ok && v.(string) != "" {
		return client.formatEndpointWithAccountID(productCode, v.(string))
	}
	return "", fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
}

// RpcPost invoking RPC API request with POST method
// parameters:
//
//	apiProductCode: API Product code, its value equals to the gateway code of the API
//	apiVersion - API version
//	apiName - API Name
//	query - API parameters in query
//	body - API parameters in body
//	autoRetry - whether to auto retry while the runtime has a 5xx error
func (client *AliyunClient) RpcPost(apiProductCode string, apiVersion string, apiName string, query map[string]interface{}, body map[string]interface{}, autoRetry bool) (map[string]interface{}, error) {
	return client.rpcRequest("POST", apiProductCode, apiVersion, apiName, query, body, autoRetry, "")
}

// RpcPostWithEndpoint invoking RPC API request with POST method and specified endpoint
// parameters:
//
//		apiProductCode: API Product code, its value equals to the gateway code of the API
//		apiVersion - API version
//		apiName - API Name
//		query - API parameters in query
//		body - API parameters in body
//		autoRetry - whether to auto retry while the runtime has a 5xx error
//	 endpoint - The domain of invoking api
func (client *AliyunClient) RpcPostWithEndpoint(apiProductCode string, apiVersion string, apiName string, query map[string]interface{}, body map[string]interface{}, autoRetry bool, endpoint string) (map[string]interface{}, error) {
	return client.rpcRequest("POST", apiProductCode, apiVersion, apiName, query, body, autoRetry, endpoint)
}

// RpcGet invoking RPC API request with GET method
// parameters:
//
//	apiProductCode: API Product code, its value equals to the gateway code of the API
//	apiVersion - API version
//	apiName - API Name
//	query - API parameters in query
//	body - API parameters in body
func (client *AliyunClient) RpcGet(apiProductCode string, apiVersion string, apiName string, query map[string]interface{}, body map[string]interface{}) (map[string]interface{}, error) {
	return client.rpcRequest("GET", apiProductCode, apiVersion, apiName, query, body, true, "")
}

func (client *AliyunClient) rpcRequest(method string, apiProductCode string, apiVersion string, apiName string, query map[string]interface{}, body map[string]interface{}, autoRetry bool, endpoint string) (map[string]interface{}, error) {
	var err error
	if endpoint == "" {
		apiProductCode = strings.ToLower(ConvertKebabToSnake(apiProductCode))
		endpoint, err = client.loadApiEndpoint(apiProductCode)
		if err != nil {
			return nil, err
		}
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	credential, err := client.config.Credential.GetCredential()
	if err != nil || credential == nil {
		return nil, fmt.Errorf("get credential failed. Error: %#v", err)
	}
	sdkConfig.SetAccessKeyId(*credential.AccessKeyId)
	sdkConfig.SetAccessKeySecret(*credential.AccessKeySecret)
	sdkConfig.SetSecurityToken(*credential.SecurityToken)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s api client: %#v", apiProductCode, err)
	}
	runtime := &util.RuntimeOptions{}
	runtime.SetAutoretry(autoRetry)
	response, err := conn.DoRequest(tea.String(apiName), nil, tea.String(method), tea.String(apiVersion), tea.String("AK"), query, body, runtime)
	return response, formatError(response, err)
}

// RoaPost invoking ROA API request with POST method
// parameters:
//
//	apiProductCode: API Product code, its value equals to the gateway code of the API
//	apiVersion - API version
//	pathName - Request path name
//	query - API parameters in query
//	headers - API parameters in headers
//	body - API parameters in body
//	autoRetry - whether to auto retry while the runtime has a 5xx error
func (client *AliyunClient) RoaPost(apiProductCode string, apiVersion string, pathName string, query map[string]*string, headers map[string]*string, body interface{}, autoRetry bool) (map[string]interface{}, error) {
	return client.roaRequest("POST", apiProductCode, apiVersion, "", pathName, query, headers, body, autoRetry)
}

// RoaPut invoking ROA API request with PUT method
// parameters:
//
//	apiProductCode: API Product code, its value equals to the gateway code of the API
//	apiVersion - API version
//	pathName - Request path name
//	query - API parameters in query
//	headers - API parameters in headers
//	body - API parameters in body
//	autoRetry - whether to auto retry while the runtime has a 5xx error
func (client *AliyunClient) RoaPut(apiProductCode string, apiVersion string, pathName string, query map[string]*string, headers map[string]*string, body interface{}, autoRetry bool) (map[string]interface{}, error) {
	return client.roaRequest("PUT", apiProductCode, apiVersion, "", pathName, query, headers, body, autoRetry)
}

// RoaGet invoking ROA API request with GET method
// parameters:
//
//	apiProductCode: API Product code, its value equals to the gateway code of the API
//	apiVersion - API version
//	pathName - Request path name
//	query - API parameters in query
//	headers - API parameters in headers
//	body - API parameters in body
func (client *AliyunClient) RoaGet(apiProductCode string, apiVersion string, pathName string, query map[string]*string, headers map[string]*string, body interface{}) (map[string]interface{}, error) {
	return client.roaRequest("GET", apiProductCode, apiVersion, "", pathName, query, headers, body, true)
}

// RoaDelete invoking ROA API request with DELETE method
// parameters:
//
//	apiProductCode: API Product code, its value equals to the gateway code of the API
//	apiVersion - API version
//	pathName - Request path name
//	query - API parameters in query
//	headers - API parameters in headers
//	body - API parameters in body
//	autoRetry - whether to auto retry while the runtime has a 5xx error
func (client *AliyunClient) RoaDelete(apiProductCode string, apiVersion string, pathName string, query map[string]*string, headers map[string]*string, body interface{}, autoRetry bool) (map[string]interface{}, error) {
	return client.roaRequest("DELETE", apiProductCode, apiVersion, "", pathName, query, headers, body, autoRetry)
}

// RoaPatch invoking ROA API request with PATCH method
// parameters:
//
//	apiProductCode: API Product code, its value equals to the gateway code of the API
//	apiVersion - API version
//	pathName - Request path name
//	query - API parameters in query
//	headers - API parameters in headers
//	body - API parameters in body
//	autoRetry - whether to auto retry while the runtime has a 5xx error
func (client *AliyunClient) RoaPatch(apiProductCode string, apiVersion string, pathName string, query map[string]*string, headers map[string]*string, body interface{}, autoRetry bool) (map[string]interface{}, error) {
	return client.roaRequest("PATCH", apiProductCode, apiVersion, "", pathName, query, headers, body, autoRetry)
}

// RoaPostWithApiName invoking ROA API request with POST method
// parameters:
//
//	apiProductCode: API Product code, its value equals to the gateway code of the API
//	apiVersion - API version
//	apiName - Request path name
//	pathName - Request path name
//	query - API parameters in query
//	headers - API parameters in headers
//	body - API parameters in body
//	autoRetry - whether to auto retry while the runtime has a 5xx error
func (client *AliyunClient) RoaPostWithApiName(apiProductCode string, apiVersion string, apiName string, pathName string, query map[string]*string, headers map[string]*string, body interface{}, autoRetry bool) (map[string]interface{}, error) {
	return client.roaRequest("POST", apiProductCode, apiVersion, apiName, pathName, query, headers, body, autoRetry)
}

// RoaPutWithApiName invoking ROA API request with PUT method
// parameters:
//
//	apiProductCode: API Product code, its value equals to the gateway code of the API
//	apiVersion - API version
//	apiName - Request path name
//	pathName - Request path name
//	query - API parameters in query
//	headers - API parameters in headers
//	body - API parameters in body
//	autoRetry - whether to auto retry while the runtime has a 5xx error
func (client *AliyunClient) RoaPutWithApiName(apiProductCode string, apiVersion string, apiName string, pathName string, query map[string]*string, headers map[string]*string, body interface{}, autoRetry bool) (map[string]interface{}, error) {
	return client.roaRequest("PUT", apiProductCode, apiVersion, apiName, pathName, query, headers, body, autoRetry)
}

// RoaGetWithApiName invoking ROA API request with GET method
// parameters:
//
//	apiProductCode: API Product code, its value equals to the gateway code of the API
//	apiVersion - API version
//	apiName - Request path name
//	pathName - Request path name
//	query - API parameters in query
//	headers - API parameters in headers
//	body - API parameters in body
func (client *AliyunClient) RoaGetWithApiName(apiProductCode string, apiVersion string, apiName string, pathName string, query map[string]*string, headers map[string]*string, body interface{}) (map[string]interface{}, error) {
	return client.roaRequest("GET", apiProductCode, apiVersion, apiName, pathName, query, headers, body, true)
}

// RoaDeleteWithApiName invoking ROA API request with DELETE method
// parameters:
//
//	apiProductCode: API Product code, its value equals to the gateway code of the API
//	apiVersion - API version
//	apiName - Request path name
//	pathName - Request path name
//	query - API parameters in query
//	headers - API parameters in headers
//	body - API parameters in body
//	autoRetry - whether to auto retry while the runtime has a 5xx error
func (client *AliyunClient) RoaDeleteWithApiName(apiProductCode string, apiVersion string, apiName string, pathName string, query map[string]*string, headers map[string]*string, body interface{}, autoRetry bool) (map[string]interface{}, error) {
	return client.roaRequest("DELETE", apiProductCode, apiVersion, apiName, pathName, query, headers, body, autoRetry)
}

// RoaPatchWithApiName invoking ROA API request with PATCH method
// parameters:
//
//	apiProductCode: API Product code, its value equals to the gateway code of the API
//	apiVersion - API version
//	apiName - Request path name
//	pathName - Request path name
//	query - API parameters in query
//	headers - API parameters in headers
//	body - API parameters in body
//	autoRetry - whether to auto retry while the runtime has a 5xx error
func (client *AliyunClient) RoaPatchWithApiName(apiProductCode string, apiVersion string, apiName string, pathName string, query map[string]*string, headers map[string]*string, body interface{}, autoRetry bool) (map[string]interface{}, error) {
	return client.roaRequest("PATCH", apiProductCode, apiVersion, apiName, pathName, query, headers, body, autoRetry)
}

func (client *AliyunClient) roaRequest(method string, apiProductCode string, apiVersion string, apiName string, pathName string, query map[string]*string, headers map[string]*string, body interface{}, autoRetry bool) (map[string]interface{}, error) {
	apiProductCode = strings.ToLower(ConvertKebabToSnake(apiProductCode))
	endpoint, err := client.loadApiEndpoint(apiProductCode)
	if err != nil {
		return nil, err
	}
	sdkConfig := client.teaRoaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	credential, err := client.config.Credential.GetCredential()
	if err != nil || credential == nil {
		return nil, fmt.Errorf("get credential failed. Error: %#v", err)
	}
	sdkConfig.SetAccessKeyId(*credential.AccessKeyId)
	sdkConfig.SetAccessKeySecret(*credential.AccessKeySecret)
	sdkConfig.SetSecurityToken(*credential.SecurityToken)
	conn, err := roa.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s api client: %#v", apiProductCode, err)
	}
	var response map[string]interface{}
	runtime := &util.RuntimeOptions{}
	runtime.SetAutoretry(autoRetry)
	if apiName != "" {
		response, err = conn.DoRequestWithAction(tea.String(apiName), tea.String(apiVersion), nil, tea.String(method), tea.String("AK"), tea.String(pathName), query, headers, body, runtime)
	} else {
		response, err = conn.DoRequest(tea.String(apiVersion), nil, tea.String(method), tea.String("AK"), tea.String(pathName), query, headers, body, runtime)
	}
	if respBody, isExist := response["body"]; isExist && respBody != nil {
		response = respBody.(map[string]interface{})
	}
	return response, formatError(response, err)
}

// Do invoking API request with SDK v2
// parameters:
//
//	apiProductCode: API Product code, its value equals to the gateway code of the API
//	apiParams - API params
//	query - API parameters in query
//	body - API parameters in body
//	headers - API parameters in headers
//	hostMap - API parameters in hostMap
//	autoRetry - whether to auto retry while the runtime has a 5xx error
func (client *AliyunClient) Do(apiProductCode string, apiParams *openapi.Params, query map[string]*string, body interface{}, headers map[string]*string, hostMap map[string]*string, autoRetry bool) (map[string]interface{}, error) {
	apiProductCode = strings.ToLower(ConvertKebabToSnake(apiProductCode))
	endpoint, err := client.loadApiEndpoint(apiProductCode)
	if err != nil {
		return nil, err
	}

	sdkConfig := client.teaRoaOpenapiConfig
	if apiParams.Style != nil && *apiParams.Style == "RPC" {
		sdkConfig = client.teaRpcOpenapiConfig
	}
	if apiParams.Protocol == nil || *apiParams.Protocol == "" {
		apiParams.Protocol = tea.String(client.config.Protocol)
	}
	sdkConfig.SetEndpoint(endpoint)
	credential, err := client.config.Credential.GetCredential()
	if err != nil || credential == nil {
		return nil, fmt.Errorf("get credential failed. Error: %#v", err)
	}
	sdkConfig.SetAccessKeyId(*credential.AccessKeyId)
	sdkConfig.SetAccessKeySecret(*credential.AccessKeySecret)
	sdkConfig.SetSecurityToken(*credential.SecurityToken)
	sdkConfig.SetUserAgent(client.config.getUserAgent())
	openapiClient, err := openapi.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s api client: %#v", apiProductCode, err)
	}
	if apiProductCode == "oss" {
		openapiClient.Spi, err = ossclient.NewClient()
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the %s api client: %#v", apiProductCode, err)
		}
		// SignVersion
		if ossV, ok := client.config.SignVersion.Load("oss"); ok {
			openapiClient.SignatureVersion = tea.String(ossV.(string))
		}
	}
	if apiProductCode == "sls" {
		openapiClient.Spi, err = gatewayclient.NewClient()
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the %s api client: %#v", apiProductCode, err)
		}
		openapiClient.Protocol = tea.String(client.config.Protocol)
	}
	var response map[string]interface{}
	runtime := &utilV2.RuntimeOptions{}
	runtime.SetAutoretry(autoRetry)
	if apiParams.Style != nil && *apiParams.Style == "RPC" {
		response, err = openapiClient.CallApi(apiParams, &openapi.OpenApiRequest{Query: query, Body: body, Headers: headers, HostMap: hostMap}, runtime)
	} else {
		response, err = openapiClient.Execute(apiParams, &openapi.OpenApiRequest{Query: query, Body: body, Headers: headers, HostMap: hostMap}, runtime)
	}
	if respBody, isExist := response["body"]; isExist && respBody != nil {
		if v, ok := respBody.(map[string]interface{}); ok {
			response = v
		}
	}
	return response, formatError(response, err)
}
func formatError(response map[string]interface{}, err error) error {
	if err != nil {
		return err
	}
	code, ok1 := response["Code"]
	if !ok1 {
		code, ok1 = response["code"]
	}
	// There is a design in some product api that the response has code with a map type, like FC
	if ok1 && !(isString(code) || isInteger(code)) {
		return err
	}
	// There is a design in some product api that the request is success but its message is empty and the code is 0 or other string
	// like ENS, eflo, apig and so on.
	if ok1 && (strings.ToLower(fmt.Sprint(code)) == "success" ||
		strings.ToLower(fmt.Sprint(code)) == "ok" ||
		fmt.Sprint(code) == "200" ||
		fmt.Sprint(code) == "0") {
		return err
	}
	success, ok2 := response["Success"]
	if !ok2 {
		success, ok2 = response["success"]
	}
	if ok2 && fmt.Sprint(success) == "true" {
		return err
	}
	message, ok3 := response["Message"]
	if !ok3 {
		message, ok3 = response["message"]
	}
	if ok3 && (message == nil || fmt.Sprint(message) == "") {
		return err
	}
	if ok1 || ok2 {
		statusCode := 200
		if v, ok := response["StatusCode"]; ok {
			statusCode = tea.IntValue(v.(*int))
		}
		return tea.NewSDKError(map[string]interface{}{
			"statusCode": statusCode,
			"code":       tea.ToString(code),
			"message":    tea.ToString(message),
			"data":       response,
		})
	}
	return err
}

func (client *AliyunClient) formatEndpointWithAccountID(productCode string, endpoint string) (string, error) {
	if endpoint == "" {
		return endpoint, nil
	}
	switch productCode {
	case "fc_open", "fc":
		accountId, err := client.AccountId()
		if err != nil {
			return endpoint, err
		}
		if strings.HasPrefix(endpoint, accountId) {
			return endpoint, nil
		}
		return fmt.Sprintf("%s.%s", accountId, endpoint), nil
	}
	return endpoint, nil
}

type ossCredentials struct {
	client *AliyunClient
}

func (defCre *ossCredentials) GetAccessKeyID() string {
	value, _, _ := defCre.client.config.GetRefreshCredential()
	return value
}

func (defCre *ossCredentials) GetAccessKeySecret() string {
	_, value, _ := defCre.client.config.GetRefreshCredential()
	return value
}

func (defCre *ossCredentials) GetSecurityToken() string {
	_, _, value := defCre.client.config.GetRefreshCredential()
	return value
}

type ossCredentialsProvider struct {
	client *AliyunClient
}

func (defBuild *ossCredentialsProvider) GetCredentials() oss.Credentials {
	return &ossCredentials{client: defBuild.client}
}

func (client *AliyunClient) GetRetryTimeout(defaultTimeout time.Duration) time.Duration {

	maxRetryTimeout := client.config.MaxRetryTimeout
	if maxRetryTimeout != 0 {
		return time.Duration(maxRetryTimeout) * time.Second
	}

	return defaultTimeout
}
