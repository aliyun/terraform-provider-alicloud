package connectivity

import (
	"fmt"
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

	ossclient "github.com/alibabacloud-go/alibabacloud-gateway-oss/client"
	gatewayclient "github.com/alibabacloud-go/alibabacloud-gateway-sls/client"
	roaCS "github.com/alibabacloud-go/cs-20151215/v5/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	roa "github.com/alibabacloud-go/tea-roa/client"
	rpc "github.com/alibabacloud-go/tea-rpc/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/endpoints"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
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
var providerVersion = "1.242.0"

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
	// Initialize the ECS client if necessary
	if client.ecsconn == nil {
		productCode := "ecs"
		endpoint := ""
		if v, ok := client.config.Endpoints.Load(productCode); !ok || v.(string) == "" {
			if err := client.loadEndpoint(productCode); err != nil {
				return nil, err
			}
		}
		if v, ok := client.config.Endpoints.Load(productCode); ok && v.(string) != "" {
			endpoint = v.(string)
			if endpoint == "ecs-cn-hangzhou.aliyuncs.com" {
				endpoint = "ecs.aliyuncs.com"
			}
		}
		if endpoint == "" {
			return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(ECSCode), endpoint)
		}
		ecsconn, err := ecs.NewClientWithOptions(client.config.RegionId, client.getSdkConfig().WithTimeout(time.Duration(60)*time.Second), client.config.getAuthCredential(true))
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
		ecsconn.AppendUserAgent(Terraform, client.config.TerraformVersion)
		ecsconn.AppendUserAgent(Provider, providerVersion)
		ecsconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		ecsconn.AppendUserAgent(TerraformTraceId, client.config.TerraformTraceId)
		client.ecsconn = ecsconn
	} else {
		err := client.ecsconn.InitWithOptions(client.config.RegionId, client.getSdkConfig().WithTimeout(time.Duration(60)*time.Second), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the ECS client: %#v", err)
		}
	}

	return do(client.ecsconn)
}

func (client *AliyunClient) WithOfficalCSClient(do func(*officalCS.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the CS client if necessary
	if client.officalCSConn == nil {
		endpoint := client.config.CsEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, CONTAINCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(CONTAINCode), endpoint)
		}
		csconn, err := officalCS.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the CS client: %#v", err)
		}

		csconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		csconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		csconn.SourceIp = client.config.SourceIp
		csconn.SecureTransport = client.config.SecureTransport
		csconn.AppendUserAgent(Terraform, client.config.TerraformVersion)
		csconn.AppendUserAgent(Provider, providerVersion)
		csconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		csconn.AppendUserAgent(TerraformTraceId, client.config.TerraformTraceId)
		client.officalCSConn = csconn
	}

	return do(client.officalCSConn)
}

func (client *AliyunClient) WithPolarDBClient(do func(*polardb.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the PolarDB client if necessary
	if client.polarDBconn == nil {
		endpoint := client.config.PolarDBEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, POLARDBCode)
			if endpoint == "" {
				endpoint = fmt.Sprintf("%s.polardb.aliyuncs.com", client.config.RegionId)
			}
		}

		polarDBconn, err := polardb.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the PolarDB client: %#v", err)

		}

		polarDBconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		polarDBconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		polarDBconn.SourceIp = client.config.SourceIp
		polarDBconn.SecureTransport = client.config.SecureTransport
		polarDBconn.AppendUserAgent(Terraform, client.config.TerraformVersion)
		polarDBconn.AppendUserAgent(Provider, providerVersion)
		polarDBconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		polarDBconn.AppendUserAgent(TerraformTraceId, client.config.TerraformTraceId)
		client.polarDBconn = polarDBconn
	} else {
		err := client.polarDBconn.InitWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the PolarDB client: %#v", err)
		}
	}

	return do(client.polarDBconn)
}

func (client *AliyunClient) WithSlbClient(do func(*slb.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the SLB client if necessary
	if client.slbconn == nil {
		endpoint := client.config.SlbEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, SLBCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(SLBCode), endpoint)
		}
		slbconn, err := slb.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the SLB client: %#v", err)
		}

		slbconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		slbconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		slbconn.SourceIp = client.config.SourceIp
		slbconn.SecureTransport = client.config.SecureTransport
		slbconn.AppendUserAgent(Terraform, client.config.TerraformVersion)
		slbconn.AppendUserAgent(Provider, providerVersion)
		slbconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		slbconn.AppendUserAgent(TerraformTraceId, client.config.TerraformTraceId)
		client.slbconn = slbconn
	} else {
		err := client.slbconn.InitWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the SLB client: %#v", err)
		}
	}

	return do(client.slbconn)
}

func (client *AliyunClient) WithVpcClient(do func(*vpc.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the VPC client if necessary
	if client.vpcconn == nil {
		endpoint := client.config.VpcEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, VPCCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(VPCCode), endpoint)
		}
		vpcconn, err := vpc.NewClientWithOptions(client.config.RegionId, client.getSdkConfig().WithTimeout(time.Duration(60)*time.Second), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the VPC client: %#v", err)
		}

		vpcconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		vpcconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		vpcconn.SourceIp = client.config.SourceIp
		vpcconn.SecureTransport = client.config.SecureTransport
		vpcconn.AppendUserAgent(Terraform, client.config.TerraformVersion)
		vpcconn.AppendUserAgent(Provider, providerVersion)
		vpcconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		vpcconn.AppendUserAgent(TerraformTraceId, client.config.TerraformTraceId)
		client.vpcconn = vpcconn
	} else {
		err := client.vpcconn.InitWithOptions(client.config.RegionId, client.getSdkConfig().WithTimeout(time.Duration(60)*time.Second), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the VPC client: %#v", err)
		}
	}

	return do(client.vpcconn)
}

func (client *AliyunClient) NewEcsClient() (*rpc.Client, error) {
	productCode := "ecs"
	endpoint := ""
	if v, ok := client.config.Endpoints.Load(productCode); !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			return nil, err
		}
	}
	if v, ok := client.config.Endpoints.Load(productCode); ok && v.(string) != "" {
		endpoint = v.(string)
		if endpoint == "ecs-cn-hangzhou.aliyuncs.com" {
			endpoint = "ecs.aliyuncs.com"
		}
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}

	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint).SetReadTimeout(60000)

	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}

	return conn, nil
}

func (client *AliyunClient) WithCenClient(do func(*cbn.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the CEN client if necessary
	if client.cenconn == nil {
		endpoint := client.config.CenEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, CbnCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(CbnCode), endpoint)
		}
		cenconn, err := cbn.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the CEN client: %#v", err)
		}

		cenconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		cenconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		cenconn.SourceIp = client.config.SourceIp
		cenconn.SecureTransport = client.config.SecureTransport
		cenconn.AppendUserAgent(Terraform, client.config.TerraformVersion)
		cenconn.AppendUserAgent(Provider, providerVersion)
		cenconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		cenconn.AppendUserAgent(TerraformTraceId, client.config.TerraformTraceId)
		client.cenconn = cenconn
	} else {
		err := client.cenconn.InitWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the CEN client: %#v", err)
		}
	}

	return do(client.cenconn)
}

func (client *AliyunClient) WithEssClient(do func(*ess.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the ESS client if necessary
	if client.essconn == nil {
		endpoint := client.config.EssEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, ESSCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(ESSCode), endpoint)
		}
		essconn, err := ess.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the ESS client: %#v", err)
		}

		essconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		essconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		essconn.SourceIp = client.config.SourceIp
		essconn.SecureTransport = client.config.SecureTransport
		essconn.AppendUserAgent(Terraform, client.config.TerraformVersion)
		essconn.AppendUserAgent(Provider, providerVersion)
		essconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		essconn.AppendUserAgent(TerraformTraceId, client.config.TerraformTraceId)
		client.essconn = essconn
	} else {
		err := client.essconn.InitWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the ESS client: %#v", err)
		}
	}

	return do(client.essconn)
}

func (client *AliyunClient) WithOssClient(do func(*oss.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the OSS client if necessary
	if client.ossconn == nil {
		schma := strings.ToLower(client.config.Protocol)
		endpoint := client.config.OssEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, OSSCode)
		}
		if endpoint == "" {
			endpointItem, err := client.describeEndpointForService(strings.ToLower(string(OSSCode)))
			if err != nil {
				log.Printf("describeEndpointForService got an error: %#v.", err)
			}
			endpoint = endpointItem
			if endpoint == "" {
				endpoint = fmt.Sprintf("oss-%s.aliyuncs.com", client.RegionId)
			}
		}
		if !strings.HasPrefix(endpoint, "http") {
			endpoint = fmt.Sprintf("%s://%s", schma, endpoint)
		}

		clientOptions := []oss.ClientOption{oss.UserAgent(client.getUserAgent())}
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
	}

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
	// Initialize the DNS client if necessary
	if client.dnsconn == nil {
		endpoint := client.config.DnsEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, DNSCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(DNSCode), endpoint)
		}

		dnsconn, err := alidns.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the DNS client: %#v", err)
		}
		dnsconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		dnsconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		dnsconn.SourceIp = client.config.SourceIp
		dnsconn.SecureTransport = client.config.SecureTransport
		dnsconn.AppendUserAgent(Terraform, client.config.TerraformVersion)
		dnsconn.AppendUserAgent(Provider, providerVersion)
		dnsconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		dnsconn.AppendUserAgent(TerraformTraceId, client.config.TerraformTraceId)
		client.dnsconn = dnsconn
	}

	return do(client.dnsconn)
}

func (client *AliyunClient) WithRamClient(do func(*ram.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the RAM client if necessary
	if client.ramconn == nil {
		endpoint := client.config.RamEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, RAMCode)
		}
		if strings.HasPrefix(endpoint, "http") {
			endpoint = fmt.Sprintf("https://%s", strings.TrimPrefix(endpoint, "http://"))
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(RAMCode), endpoint)
		}

		ramconn, err := ram.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the RAM client: %#v", err)
		}
		ramconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		ramconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		ramconn.SourceIp = client.config.SourceIp
		ramconn.SecureTransport = client.config.SecureTransport
		ramconn.AppendUserAgent(Terraform, client.config.TerraformVersion)
		ramconn.AppendUserAgent(Provider, providerVersion)
		ramconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		ramconn.AppendUserAgent(TerraformTraceId, client.config.TerraformTraceId)
		client.ramconn = ramconn
	} else {
		err := client.ramconn.InitWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the RAM client: %#v", err)
		}

	}

	return do(client.ramconn)
}

func (client *AliyunClient) WithCsClient(do func(*cs.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the CS client if necessary
	if client.csconn == nil {
		csconn := cs.NewClientForAussumeRole(client.config.AccessKey, client.config.SecretKey, client.config.SecurityToken)
		csconn.SetUserAgent(client.getUserAgent())
		endpoint := client.config.CsEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, CONTAINCode)
		}
		if endpoint != "" {
			if !strings.HasPrefix(endpoint, "http") {
				endpoint = fmt.Sprintf("https://%s", strings.TrimPrefix(endpoint, "://"))
			}
			csconn.SetEndpoint(endpoint)
		}
		csconn.SetSourceIp(client.config.SourceIp)
		csconn.SetSecureTransport(client.config.SecureTransport)
		client.csconn = csconn
	}

	return do(client.csconn)
}

func (client *AliyunClient) NewRoaCsClient() (*roaCS.Client, error) {
	endpoint := client.config.CsEndpoint
	if endpoint == "" {
		endpoint = OpenAckService
	}
	header := map[string]*string{
		"x-acs-source-ip":        tea.String(client.config.SourceIp),
		"x-acs-secure-transport": tea.String(client.config.SecureTransport),
	}
	param := &openapi.GlobalParameters{Headers: header}
	// Initialize the CS client if necessary
	roaCSConn, err := roaCS.NewClient(&openapi.Config{
		AccessKeyId:      tea.String(client.config.AccessKey),
		AccessKeySecret:  tea.String(client.config.SecretKey),
		SecurityToken:    tea.String(client.config.SecurityToken),
		RegionId:         tea.String(client.config.RegionId),
		UserAgent:        tea.String(client.getUserAgent()),
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
	// Initialize the CR client if necessary
	if client.crconn == nil {
		endpoint := client.config.CrEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, CRCode)
			if endpoint == "" {
				endpoint = fmt.Sprintf("cr.%s.aliyuncs.com", client.config.RegionId)
			}
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(CRCode), endpoint)
		}
		crconn, err := cr.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the CR client: %#v", err)
		}
		crconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		crconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		crconn.SourceIp = client.config.SourceIp
		crconn.SecureTransport = client.config.SecureTransport
		crconn.AppendUserAgent(Terraform, client.config.TerraformVersion)
		crconn.AppendUserAgent(Provider, providerVersion)
		crconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		crconn.AppendUserAgent(TerraformTraceId, client.config.TerraformTraceId)
		client.crconn = crconn
	}

	return do(client.crconn)
}

func (client *AliyunClient) WithCrEEClient(do func(*cr_ee.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the CR EE client if necessary
	if client.creeconn == nil {
		endpoint := client.config.CrEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, CRCode)
			if endpoint == "" {
				endpoint = fmt.Sprintf("cr.%s.aliyuncs.com", client.config.RegionId)
			}
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(CRCode), endpoint)
		}
		creeconn, err := cr_ee.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the CR EE client: %#v", err)
		}
		creeconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		creeconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		creeconn.SourceIp = client.config.SourceIp
		creeconn.SecureTransport = client.config.SecureTransport
		creeconn.AppendUserAgent(Terraform, client.config.TerraformVersion)
		creeconn.AppendUserAgent(Provider, providerVersion)
		creeconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		creeconn.AppendUserAgent(TerraformTraceId, client.config.TerraformTraceId)
		client.creeconn = creeconn
	}

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
	// Initialize the CDN client if necessary
	if client.cdnconn_new == nil {
		endpoint := client.config.CdnEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, CDNCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(CDNCode), endpoint)
		}
		cdnconn, err := cdn_new.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the CDN client: %#v", err)
		}
		cdnconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		cdnconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		cdnconn.SourceIp = client.config.SourceIp
		cdnconn.SecureTransport = client.config.SecureTransport
		cdnconn.AppendUserAgent(Terraform, client.config.TerraformVersion)
		cdnconn.AppendUserAgent(Provider, providerVersion)
		cdnconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		cdnconn.AppendUserAgent(TerraformTraceId, client.config.TerraformTraceId)
		client.cdnconn_new = cdnconn
	}

	return do(client.cdnconn_new)
}

// WithOtsClient init ots openapi publish sdk client(if necessary), and exec do func by client
func (client *AliyunClient) WithOtsClient(do func(*ots.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the OTS client if necessary
	if client.otsconn == nil {
		endpoint := client.config.OtsEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, OTSCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(OTSCode), endpoint)
		}
		otsconn, err := ots.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the OTS client: %#v", err)
		}

		otsconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		otsconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		otsconn.SourceIp = client.config.SourceIp
		otsconn.SecureTransport = client.config.SecureTransport
		otsconn.AppendUserAgent(Terraform, client.config.TerraformVersion)
		otsconn.AppendUserAgent(Provider, providerVersion)
		otsconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		otsconn.AppendUserAgent(TerraformTraceId, client.config.TerraformTraceId)
		client.otsconn = otsconn
	}

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

	conn, err := roa.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}

	return conn, nil
}

func (client *AliyunClient) WithCmsClient(do func(*cms.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the CMS client if necessary
	if client.cmsconn == nil {
		endpoint := client.config.CmsEndpoint
		if endpoint == "" {
			productCode := "cms"
			if v, ok := client.config.Endpoints.Load(productCode); !ok || v.(string) == "" {
				if err := client.loadEndpoint(productCode); err != nil {
					endpoint = fmt.Sprintf("metrics.%s.aliyuncs.com", client.RegionId)
					client.config.Endpoints.Store(productCode, endpoint)
					log.Printf("[ERROR] loading %s endpoint got an error: %#v. Using the central endpoint %s instead.", productCode, err, endpoint)
				}
			}
			if v, ok := client.config.Endpoints.Load(productCode); ok && v.(string) != "" {
				endpoint = v.(string)
			}
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(CMSCode), endpoint)
		}
		cmsconn, err := cms.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the CMS client: %#v", err)
		}

		cmsconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		cmsconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		cmsconn.SourceIp = client.config.SourceIp
		cmsconn.SecureTransport = client.config.SecureTransport
		cmsconn.AppendUserAgent(Terraform, client.config.TerraformVersion)
		cmsconn.AppendUserAgent(Provider, providerVersion)
		cmsconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		cmsconn.AppendUserAgent(TerraformTraceId, client.config.TerraformTraceId)
		client.cmsconn = cmsconn
	} else {
		err := client.cmsconn.InitWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the CMS client: %#v", err)
		}
	}

	return do(client.cmsconn)
}

func (client *AliyunClient) WithLogPopClient(do func(*slsPop.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the HBase client if necessary
	if client.logpopconn == nil {
		logpopconn, err := slsPop.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the sls client: %#v", err)
		}
		logpopconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		logpopconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		logpopconn.SourceIp = client.config.SourceIp
		logpopconn.SecureTransport = client.config.SecureTransport
		logpopconn.AppendUserAgent(Terraform, client.config.TerraformVersion)
		logpopconn.AppendUserAgent(Provider, providerVersion)
		logpopconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		logpopconn.AppendUserAgent(TerraformTraceId, client.config.TerraformTraceId)
		endpoint := client.config.LogEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, LOGCode)
			if endpoint == "" {
				endpoint = fmt.Sprintf("%s.log.aliyuncs.com", client.config.RegionId)
			}
		}
		endpoint = strings.TrimPrefix(endpoint, "https://")
		endpoint = strings.TrimPrefix(endpoint, "http://")
		logpopconn.Domain = endpoint + "/open-api"
		client.logpopconn = logpopconn
	}

	return do(client.logpopconn)
}

func (client *AliyunClient) WithLogClient(do func(*sls.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the LOG client if necessary
	if client.logconn == nil {
		endpoint := client.config.LogEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, LOGCode)
			if endpoint == "" {
				endpoint = fmt.Sprintf("%s.log.aliyuncs.com", client.config.RegionId)
			}
		}
		if !strings.HasPrefix(endpoint, "http") {
			endpoint = fmt.Sprintf("https://%s", strings.TrimPrefix(endpoint, "://"))
		}
		client.logconn = &sls.Client{
			AccessKeyID:     client.config.AccessKey,
			AccessKeySecret: client.config.SecretKey,
			Endpoint:        endpoint,
			SecurityToken:   client.config.SecurityToken,
			UserAgent:       client.getUserAgent(),
		}
	}

	return do(client.logconn)
}

func (client *AliyunClient) WithDrdsClient(do func(*drds.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the DRDS client if necessary
	if client.drdsconn == nil {
		endpoint := client.config.DrdsEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, DRDSCode)
			if endpoint == "" {
				endpoint = fmt.Sprintf("%s.drds.aliyuncs.com", client.config.RegionId)
			}
		}

		drdsconn, err := drds.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the DRDS client: %#v", err)

		}
		drdsconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		drdsconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		drdsconn.SourceIp = client.config.SourceIp
		drdsconn.SecureTransport = client.config.SecureTransport
		drdsconn.AppendUserAgent(Terraform, client.config.TerraformVersion)
		drdsconn.AppendUserAgent(Provider, providerVersion)
		drdsconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		drdsconn.AppendUserAgent(TerraformTraceId, client.config.TerraformTraceId)
		client.drdsconn = drdsconn
	} else {
		err := client.drdsconn.InitWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the CMS client: %#v", err)
		}
	}

	return do(client.drdsconn)
}

func (client *AliyunClient) WithDdsClient(do func(*dds.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the DDS client if necessary
	productCode := "dds"
	endpoint := ""

	if client.ddsconn == nil {
		if endpoint == "" {
			if v, ok := client.config.Endpoints.Load(productCode); !ok || v.(string) == "" {
				if err := client.loadEndpoint(productCode); err != nil {
					endpoint = fmt.Sprintf("mongodb.%s.aliyuncs.com", client.config.RegionId)
					client.config.Endpoints.Store(productCode, endpoint)
					log.Printf("[ERROR] loading %s endpoint got an error: %#v. Using the endpoint %s instead.", productCode, err, endpoint)
				}
			}

			if v, ok := client.config.Endpoints.Load(productCode); ok && v.(string) != "" {
				endpoint = v.(string)
			}
		}

		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(DDSCode), endpoint)
		} else {
			return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
		}

		ddsconn, err := dds.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the DDS client: %#v", err)
		}

		ddsconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		ddsconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		ddsconn.SourceIp = client.config.SourceIp
		ddsconn.SecureTransport = client.config.SecureTransport
		ddsconn.AppendUserAgent(Terraform, client.config.TerraformVersion)
		ddsconn.AppendUserAgent(Provider, providerVersion)
		ddsconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		ddsconn.AppendUserAgent(TerraformTraceId, client.config.TerraformTraceId)
		client.ddsconn = ddsconn
	}

	return do(client.ddsconn)
}

func (client *AliyunClient) WithGpdbClient(do func(*gpdb.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the GPDB client if necessary
	if client.gpdbconn == nil {
		endpoint := client.config.GpdbEnpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, GPDBCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(GPDBCode), endpoint)
		}
		gpdbconn, err := gpdb.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the GPDB client: %#v", err)
		}
		gpdbconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		gpdbconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		gpdbconn.SourceIp = client.config.SourceIp
		gpdbconn.SecureTransport = client.config.SecureTransport
		gpdbconn.AppendUserAgent(Terraform, client.config.TerraformVersion)
		gpdbconn.AppendUserAgent(Provider, providerVersion)
		gpdbconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		gpdbconn.AppendUserAgent(TerraformTraceId, client.config.TerraformTraceId)
		client.gpdbconn = gpdbconn
	}

	return do(client.gpdbconn)
}

func (client *AliyunClient) NewGpdbClient() (*rpc.Client, error) {
	productCode := "gpdb"
	endpoint := ""
	if v, ok := client.config.Endpoints.Load(productCode); !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			return nil, err
		}
	}
	if v, ok := client.config.Endpoints.Load(productCode); ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}

	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)

	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}

	return conn, nil
}

func (client *AliyunClient) WithFcClient(do func(*fc.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()
	endpoint := client.config.FcEndpoint
	if endpoint == "" {
		endpoint = loadEndpoint(client.config.RegionId, FCCode)
		if endpoint == "" {
			endpoint = fmt.Sprintf("%s.fc.aliyuncs.com", client.config.RegionId)
		}
	}
	if strings.HasPrefix(endpoint, "http") {
		endpoint = strings.TrimPrefix(strings.TrimPrefix(endpoint, "http://"), "https://")
	}
	accountId, err := client.AccountId()
	if err != nil {
		return nil, err
	}

	config := client.getSdkConfig()
	transport := config.HttpTransport
	// Receiving proxy config from environment
	transport.Proxy = http.ProxyFromEnvironment
	clientOptions := []fc.ClientOption{fc.WithSecurityToken(client.config.SecurityToken), fc.WithTransport(transport),
		fc.WithTimeout(30), fc.WithRetryCount(DefaultClientRetryCountSmall)}

	// Initialize the FC client if necessary
	if client.fcconn == nil {
		fcconn, err := fc.NewClient(fmt.Sprintf("https://%s.%s", accountId, endpoint), string(ApiVersion20160815), client.config.AccessKey, client.config.SecretKey, clientOptions...)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the FC client: %#v", err)
		}

		fcconn.Config.UserAgent = client.getUserAgent()
		fcconn.Config.SecurityToken = client.config.SecurityToken
		client.fcconn = fcconn
	} else {
		fcconn, err := fc.NewClient(fmt.Sprintf("https://%s.%s", accountId, endpoint), string(ApiVersion20160815), client.config.AccessKey, client.config.SecretKey, clientOptions...)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the FC client: %#v", err)
		}
		fcconn.Config.UserAgent = client.getUserAgent()
		fcconn.Config.SecurityToken = client.config.SecurityToken
		client.fcconn = fcconn
	}

	return do(client.fcconn)
}

func (client *AliyunClient) WithCloudApiClient(do func(*cloudapi.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the Cloud API client if necessary
	if client.cloudapiconn == nil {
		endpoint := client.config.ApigatewayEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.RegionId, CLOUDAPICode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.RegionId, "CLOUDAPI", endpoint)
		}
		cloudapiconn, err := cloudapi.NewClientWithOptions(client.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the CloudAPI client: %#v", err)
		}
		client.cloudapiconn = cloudapiconn
	} else {
		err := client.cloudapiconn.InitWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the CloudAPI client: %#v", err)
		}
	}
	client.cloudapiconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
	client.cloudapiconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
	client.cloudapiconn.SourceIp = client.config.SourceIp
	client.cloudapiconn.SecureTransport = client.config.SecureTransport
	client.cloudapiconn.AppendUserAgent(Terraform, client.config.TerraformVersion)
	client.cloudapiconn.AppendUserAgent(Provider, providerVersion)
	client.cloudapiconn.AppendUserAgent(Module, client.config.ConfigurationSource)
	client.cloudapiconn.AppendUserAgent(TerraformTraceId, client.config.TerraformTraceId)

	return do(client.cloudapiconn)
}

func (client *AliyunClient) NewTeaCommonClient(endpoint string) (*rpc.Client, error) {
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)

	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the tea client: %#v", err)
	}

	return conn, nil
}
func (client *AliyunClient) NewTeaRoaCommonClient(endpoint string) (*roa.Client, error) {
	sdkConfig := client.teaRoaSdkConfig
	sdkConfig.SetEndpoint(endpoint)

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

		mnsClient := ali_mns.NewAliMNSClientWithToken(mnsUrl, client.config.AccessKey, client.config.SecretKey, client.config.SecurityToken)
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

func (client *AliyunClient) WithElasticsearchClient(do func(*elasticsearch.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the Elasticsearch client if necessary
	if client.elasticsearchconn == nil {
		endpoint := client.config.ElasticsearchEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, ELASTICSEARCHCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(ELASTICSEARCHCode), endpoint)
		}
		elasticsearchconn, err := elasticsearch.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the Elasticsearch client: %#v", err)
		}
		elasticsearchconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		elasticsearchconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		elasticsearchconn.SourceIp = client.config.SourceIp
		elasticsearchconn.SecureTransport = client.config.SecureTransport
		elasticsearchconn.AppendUserAgent(Terraform, client.config.TerraformVersion)
		elasticsearchconn.AppendUserAgent(Provider, providerVersion)
		elasticsearchconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		elasticsearchconn.AppendUserAgent(TerraformTraceId, client.config.TerraformTraceId)
		client.elasticsearchconn = elasticsearchconn
	}

	return do(client.elasticsearchconn)
}

func (client *AliyunClient) NewElasticsearchClient() (*roa.Client, error) {
	productCode := "elasticsearch"
	endpoint := ""
	if v, ok := client.config.Endpoints.Load(productCode); !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			return nil, err
		}
	}

	if v, ok := client.config.Endpoints.Load(productCode); ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] misssing the product %s endpoint.", productCode)
	}
	roaSdkConfig := client.teaRoaSdkConfig
	roaSdkConfig.SetEndpoint(endpoint)

	conn, err := roa.NewClient(&roaSdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, err
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
	if !ok {
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
		tableStoreClient = tablestore.NewClientWithExternalHeader(endpoint, instanceName, client.config.AccessKey, client.config.SecretKey, client.config.SecurityToken, tablestore.NewDefaultTableStoreConfig(), externalHeaders)
		client.tablestoreconnByInstanceName[instanceName] = tableStoreClient
	}

	return do(tableStoreClient)
}

func (client *AliyunClient) WithTableStoreTunnelClient(instanceName string, do func(otsTunnel.TunnelClient) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the TABLESTORE tunnel client if necessary
	tunnelClient, ok := client.otsTunnelConnByInstanceName[instanceName]
	if !ok {
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
		tunnelClient = otsTunnel.NewTunnelClientWithConfigAndExternalHeader(endpoint, instanceName, client.config.AccessKey, client.config.SecretKey, client.config.SecurityToken, otsTunnel.DefaultTunnelConfig, externalHeaders)
		client.otsTunnelConnByInstanceName[instanceName] = tunnelClient
	}

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
		csProjectClient.SetUserAgent(client.getUserAgent())
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
	client.accountIdMutex.Lock()
	defer client.accountIdMutex.Unlock()

	if client.accountId == "" {
		log.Printf("[DEBUG] account_id not provided, attempting to retrieve it automatically...")
		identity, err := client.GetCallerIdentity()
		if err != nil {
			return "", err
		}
		if identity.AccountId == "" {
			return "", fmt.Errorf("caller identity doesn't contain any AccountId")
		}
		client.accountId = identity.AccountId
	}
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

func (client *AliyunClient) getSdkConfig() *sdk.Config {
	return sdk.NewConfig().
		WithMaxRetryTime(DefaultClientRetryCountSmall).
		WithTimeout(time.Duration(30) * time.Second).
		WithEnableAsync(false).
		WithGoRoutinePoolSize(100).
		WithMaxTaskQueueSize(10000).
		WithDebug(false).
		WithHttpTransport(client.getTransport()).
		WithScheme(client.config.Protocol)
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
	stsClient, err := sts.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the STS client: %#v", err)
	}

	stsClient.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
	stsClient.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
	stsClient.SourceIp = client.config.SourceIp
	stsClient.SecureTransport = client.config.SecureTransport
	stsClient.AppendUserAgent(Terraform, client.config.TerraformVersion)
	stsClient.AppendUserAgent(Provider, providerVersion)
	stsClient.AppendUserAgent(Module, client.config.ConfigurationSource)
	stsClient.AppendUserAgent(TerraformTraceId, client.config.TerraformTraceId)

	identity, err := stsClient.GetCallerIdentity(args)
	if err != nil {
		return nil, err
	}
	if identity == nil {
		return nil, fmt.Errorf("caller identity not found")
	}
	return identity, err
}

func (client *AliyunClient) WithDdoscooClient(do func(*ddoscoo.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the ddoscoo client if necessary
	if client.ddoscooconn == nil {
		endpoint := client.config.DdoscooEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, DDOSCOOCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(DDOSCOOCode), endpoint)
		}

		ddoscooconn, err := ddoscoo.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the DDOSCOO client: %#v", err)
		}
		ddoscooconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		ddoscooconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		ddoscooconn.SourceIp = client.config.SourceIp
		ddoscooconn.SecureTransport = client.config.SecureTransport
		ddoscooconn.AppendUserAgent(Terraform, client.config.TerraformVersion)
		ddoscooconn.AppendUserAgent(Provider, providerVersion)
		ddoscooconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		ddoscooconn.AppendUserAgent(TerraformTraceId, client.config.TerraformTraceId)
		client.ddoscooconn = ddoscooconn
	} else {
		err := client.ddoscooconn.InitWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the DdosCoo client: %#v", err)
		}
	}

	return do(client.ddoscooconn)
}

func (client *AliyunClient) WithDdosbgpClient(do func(*ddosbgp.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the ddosbgp client if necessary
	if client.ddosbgpconn == nil {
		endpoint := client.config.DdosbgpEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, DDOSBGPCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(DDOSBGPCode), endpoint)
		}

		ddosbgpconn, err := ddosbgp.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the DDOSBGP client: %#v", err)
		}
		ddosbgpconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		ddosbgpconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		ddosbgpconn.SourceIp = client.config.SourceIp
		ddosbgpconn.SecureTransport = client.config.SecureTransport
		ddosbgpconn.AppendUserAgent(Terraform, client.config.TerraformVersion)
		ddosbgpconn.AppendUserAgent(Provider, providerVersion)
		ddosbgpconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		ddosbgpconn.AppendUserAgent(TerraformTraceId, client.config.TerraformTraceId)
		client.ddosbgpconn = ddosbgpconn
	}

	return do(client.ddosbgpconn)
}

func (client *AliyunClient) WithBssopenapiClient(do func(*bssopenapi.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the bssopenapi client if necessary
	if client.bssopenapiconn == nil {
		endpoint := client.config.BssOpenApiEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, BSSOPENAPICode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(BSSOPENAPICode), endpoint)
		}

		// bss endpoint depends on the account type.
		// Domestic account is business.aliyuncs.com (region is cn-hangzhou) and International account is business.ap-southeast-1.aliyuncs.com (region is ap-southeast-1)
		bssopenapiconn, err := bssopenapi.NewClientWithOptions(string(Hangzhou), client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the BSSOPENAPI client: %#v", err)
		}
		bssopenapiconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		bssopenapiconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		bssopenapiconn.SourceIp = client.config.SourceIp
		bssopenapiconn.SecureTransport = client.config.SecureTransport
		bssopenapiconn.AppendUserAgent(Terraform, client.config.TerraformVersion)
		bssopenapiconn.AppendUserAgent(Provider, providerVersion)
		bssopenapiconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		bssopenapiconn.AppendUserAgent(TerraformTraceId, client.config.TerraformTraceId)
		client.bssopenapiconn = bssopenapiconn
	}

	return do(client.bssopenapiconn)
}

func (client *AliyunClient) WithAlikafkaClient(do func(*alikafka.Client) (interface{}, error)) (interface{}, error) {
	productCode := "alikafka"
	endpoint := client.config.AlikafkaEndpoint
	if client.alikafkaconn == nil {
		if v, ok := client.config.Endpoints.Load(productCode); !ok || v.(string) == "" {
			if err := client.loadEndpoint(productCode); err != nil {
				return nil, err
			}
		}
		if v, ok := client.config.Endpoints.Load(productCode); ok && v.(string) != "" {
			endpoint = v.(string)
		}
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, ALIKAFKACode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(ALIKAFKACode), endpoint)
		}
		alikafkaconn, err := alikafka.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the ALIKAFKA client: %#v", err)
		}
		alikafkaconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		alikafkaconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		alikafkaconn.SourceIp = client.config.SourceIp
		alikafkaconn.SecureTransport = client.config.SecureTransport
		alikafkaconn.AppendUserAgent(Terraform, client.config.TerraformVersion)
		alikafkaconn.AppendUserAgent(Provider, providerVersion)
		alikafkaconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		alikafkaconn.AppendUserAgent(TerraformTraceId, client.config.TerraformTraceId)
		client.alikafkaconn = alikafkaconn
	}

	return do(client.alikafkaconn)
}

func (client *AliyunClient) WithEmrClient(do func(*emr.Client) (interface{}, error)) (interface{}, error) {
	if client.emrconn == nil {
		endpoint := client.config.EmrEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, EMRCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(EMRCode), endpoint)
		}
		emrConn, err := emr.NewClientWithOptions(client.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the E-MapReduce client: %#v", err)
		}
		emrConn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		emrConn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		emrConn.SourceIp = client.config.SourceIp
		emrConn.SecureTransport = client.config.SecureTransport
		emrConn.AppendUserAgent(Terraform, client.config.TerraformVersion)
		emrConn.AppendUserAgent(Provider, providerVersion)
		emrConn.AppendUserAgent(Module, client.config.ConfigurationSource)
		emrConn.AppendUserAgent(TerraformTraceId, client.config.TerraformTraceId)
		client.emrconn = emrConn
	}

	return do(client.emrconn)
}

func (client *AliyunClient) WithSagClient(do func(*smartag.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the SAG client if necessary
	if client.sagconn == nil {
		endpoint := client.config.SagEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, SAGCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(SAGCode), endpoint)
		}
		sagconn, err := smartag.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the SAG client: %#v", err)
		}
		sagconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		sagconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		sagconn.SourceIp = client.config.SourceIp
		sagconn.SecureTransport = client.config.SecureTransport
		sagconn.AppendUserAgent(Terraform, client.config.TerraformVersion)
		sagconn.AppendUserAgent(Provider, providerVersion)
		sagconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		sagconn.AppendUserAgent(TerraformTraceId, client.config.TerraformTraceId)
		client.sagconn = sagconn
	}

	return do(client.sagconn)
}

func (client *AliyunClient) WithDbauditClient(do func(*yundun_dbaudit.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the ddoscoo client if necessary
	if client.dbauditconn == nil {
		dbauditconn, err := yundun_dbaudit.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the DBAUDIT client: %#v", err)
		}
		dbauditconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		dbauditconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		dbauditconn.SourceIp = client.config.SourceIp
		dbauditconn.SecureTransport = client.config.SecureTransport
		dbauditconn.AppendUserAgent(Terraform, client.config.TerraformVersion)
		dbauditconn.AppendUserAgent(Provider, providerVersion)
		dbauditconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		dbauditconn.AppendUserAgent(TerraformTraceId, client.config.TerraformTraceId)
		client.dbauditconn = dbauditconn
	}

	return do(client.dbauditconn)
}
func (client *AliyunClient) WithMarketClient(do func(*market.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the Market API client if necessary
	if client.marketconn == nil {
		endpoint := client.config.MarketEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.RegionId, MARKETCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.RegionId, "MARKET", endpoint)
		}
		marketconn, err := market.NewClientWithOptions(client.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the Market client: %#v", err)
		}
		marketconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		marketconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		marketconn.SourceIp = client.config.SourceIp
		marketconn.SecureTransport = client.config.SecureTransport
		marketconn.AppendUserAgent(Terraform, client.config.TerraformVersion)
		marketconn.AppendUserAgent(Provider, providerVersion)
		marketconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		marketconn.AppendUserAgent(TerraformTraceId, client.config.TerraformTraceId)
		client.marketconn = marketconn
	} else {
		err := client.marketconn.InitWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the Market client: %#v", err)
		}
	}

	return do(client.marketconn)
}

func (client *AliyunClient) WithHbaseClient(do func(*hbase.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the HBase client if necessary
	if client.hbaseconn == nil {
		endpoint := client.config.HBaseEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, HBASECode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(HBASECode), endpoint)
		}
		hbaseconn, err := hbase.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the hbase client: %#v", err)
		}
		hbaseconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		hbaseconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		hbaseconn.SourceIp = client.config.SourceIp
		hbaseconn.SecureTransport = client.config.SecureTransport
		hbaseconn.AppendUserAgent(Terraform, client.config.TerraformVersion)
		hbaseconn.AppendUserAgent(Provider, providerVersion)
		hbaseconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		hbaseconn.AppendUserAgent(TerraformTraceId, client.config.TerraformTraceId)
		client.hbaseconn = hbaseconn
	} else {
		err := client.hbaseconn.InitWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the HBase client: %#v", err)
		}
	}

	return do(client.hbaseconn)
}

func (client *AliyunClient) WithAdbClient(do func(*adb.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the adb client if necessary
	if client.adbconn == nil {
		endpoint := client.config.AdbEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, ADBCode)
			if endpoint == "" {
				endpoint = fmt.Sprintf("%s.adb.aliyuncs.com", client.config.RegionId)
			}
		}

		adbconn, err := adb.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the adb client: %#v", err)

		}
		adbconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		adbconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		adbconn.SourceIp = client.config.SourceIp
		adbconn.SecureTransport = client.config.SecureTransport
		adbconn.AppendUserAgent(Terraform, client.config.TerraformVersion)
		adbconn.AppendUserAgent(Provider, providerVersion)
		adbconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		adbconn.AppendUserAgent(TerraformTraceId, client.config.TerraformTraceId)
		client.adbconn = adbconn
	}

	return do(client.adbconn)
}
func (client *AliyunClient) WithCbnClient(do func(*cbn.Client) (interface{}, error)) (interface{}, error) {
	product := "cbn"
	endpoint, err := client.loadApiEndpoint(product)
	if err != nil {
		return nil, err
	}
	if endpoint != "" {
		endpoints.AddEndpointMapping(client.config.RegionId, product, endpoint)
	}

	if client.cbnConn == nil {
		cbnConn, err := cbn.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the Cbnclient: %#v", err)
		}
		client.cbnConn = cbnConn
	} else {
		err := client.cbnConn.InitWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the %s client: %#v", product, err)
		}
	}
	client.cbnConn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
	client.cbnConn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
	client.cbnConn.SourceIp = client.config.SourceIp
	client.cbnConn.SecureTransport = client.config.SecureTransport
	client.cbnConn.AppendUserAgent(Terraform, client.config.TerraformVersion)
	client.cbnConn.AppendUserAgent(Provider, providerVersion)
	client.cbnConn.AppendUserAgent(Module, client.config.ConfigurationSource)
	client.cbnConn.AppendUserAgent(TerraformTraceId, client.config.TerraformTraceId)
	return do(client.cbnConn)
}

func (client *AliyunClient) WithEdasClient(do func(*edas.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the edas client if necessary
	if client.edasconn == nil {
		endpoint := client.config.edasEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, EDASCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(EDASCode), endpoint)
		}
		edasconn, err := edas.NewClientWithOptions(client.config.RegionId, client.getSdkConfig().WithTimeout(time.Duration(60)*time.Second), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the ALIKAFKA client: %#v", err)
		}
		edasconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		edasconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		edasconn.SourceIp = client.config.SourceIp
		edasconn.SecureTransport = client.config.SecureTransport
		edasconn.AppendUserAgent(Terraform, client.config.TerraformVersion)
		edasconn.AppendUserAgent(Provider, providerVersion)
		edasconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		edasconn.AppendUserAgent(TerraformTraceId, client.config.TerraformTraceId)
		client.edasconn = edasconn
	} else {
		err := client.edasconn.InitWithOptions(client.config.RegionId, client.getSdkConfig().WithTimeout(time.Duration(60)*time.Second), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the EDAS client: %#v", err)
		}
	}

	return do(client.edasconn)
}

func (client *AliyunClient) WithAlidnsClient(do func(*alidns.Client) (interface{}, error)) (interface{}, error) {
	if client.alidnsConn == nil {
		endpoint := client.config.AlidnsEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, AlidnsCode)
		}
		if strings.HasPrefix(endpoint, "http") {
			endpoint = fmt.Sprintf("https://%s", strings.TrimPrefix(endpoint, "http://"))
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(AlidnsCode), endpoint)
		}

		alidnsConn, err := alidns.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the Alidnsclient: %#v", err)
		}
		alidnsConn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		alidnsConn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		alidnsConn.SourceIp = client.config.SourceIp
		alidnsConn.SecureTransport = client.config.SecureTransport
		alidnsConn.AppendUserAgent(Terraform, client.config.TerraformVersion)
		alidnsConn.AppendUserAgent(Provider, providerVersion)
		alidnsConn.AppendUserAgent(Module, client.config.ConfigurationSource)
		alidnsConn.AppendUserAgent(TerraformTraceId, client.config.TerraformTraceId)
		client.alidnsConn = alidnsConn
	}
	return do(client.alidnsConn)
}

func (client *AliyunClient) WithCassandraClient(do func(*cassandra.Client) (interface{}, error)) (interface{}, error) {
	if client.cassandraConn == nil {
		endpoint := client.config.CassandraEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, CassandraCode)
			endpoints.AddEndpointMapping(client.config.RegionId, string(CassandraCode), endpoint)
		}
		cassandraConn, err := cassandra.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the Cassandraclient: %#v", err)
		}
		cassandraConn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		cassandraConn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		cassandraConn.SourceIp = client.config.SourceIp
		cassandraConn.SecureTransport = client.config.SecureTransport
		cassandraConn.AppendUserAgent(Terraform, client.config.TerraformVersion)
		cassandraConn.AppendUserAgent(Provider, providerVersion)
		cassandraConn.AppendUserAgent(Module, client.config.ConfigurationSource)
		cassandraConn.AppendUserAgent(TerraformTraceId, client.config.TerraformTraceId)
		client.cassandraConn = cassandraConn
	}
	return do(client.cassandraConn)
}

func (client *AliyunClient) WithEciClient(do func(*eci.Client) (interface{}, error)) (interface{}, error) {
	if client.eciConn == nil {
		endpoint := client.config.EciEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, EciCode)
		}
		if strings.HasPrefix(endpoint, "http") {
			endpoint = fmt.Sprintf("https://%s", strings.TrimPrefix(endpoint, "http://"))
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(EciCode), endpoint)
		}

		eciConn, err := eci.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the Eciclient: %#v", err)
		}
		eciConn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		eciConn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		eciConn.SourceIp = client.config.SourceIp
		eciConn.SecureTransport = client.config.SecureTransport
		eciConn.AppendUserAgent(Terraform, client.config.TerraformVersion)
		eciConn.AppendUserAgent(Provider, providerVersion)
		eciConn.AppendUserAgent(Module, client.config.ConfigurationSource)
		eciConn.AppendUserAgent(TerraformTraceId, client.config.TerraformTraceId)
		client.eciConn = eciConn
	}
	return do(client.eciConn)
}

func (client *AliyunClient) WithDcdnClient(do func(*dcdn.Client) (interface{}, error)) (interface{}, error) {
	if client.dcdnConn == nil {
		endpoint := client.config.DcdnEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, DcdnCode)
		}
		if strings.HasPrefix(endpoint, "http") {
			endpoint = fmt.Sprintf("https://%s", strings.TrimPrefix(endpoint, "http://"))
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(DcdnCode), endpoint)
		}

		dcdnConn, err := dcdn.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the Dcdnclient: %#v", err)
		}
		dcdnConn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		dcdnConn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		dcdnConn.SourceIp = client.config.SourceIp
		dcdnConn.SecureTransport = client.config.SecureTransport
		dcdnConn.AppendUserAgent(Terraform, client.config.TerraformVersion)
		dcdnConn.AppendUserAgent(Provider, providerVersion)
		dcdnConn.AppendUserAgent(Module, client.config.ConfigurationSource)
		dcdnConn.AppendUserAgent(TerraformTraceId, client.config.TerraformTraceId)
		client.dcdnConn = dcdnConn
	}
	return do(client.dcdnConn)
}

func (client *AliyunClient) WithRKvstoreClient(do func(*r_kvstore.Client) (interface{}, error)) (interface{}, error) {
	if client.r_kvstoreConn == nil {
		productCode := "r_kvstore"
		endpoint := ""
		if v, ok := client.config.Endpoints.Load(productCode); !ok || v.(string) == "" {
			if err := client.loadEndpoint(productCode); err != nil {
				return nil, err
			}
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, "r-kvstore", endpoint)
		}

		r_kvstoreConn, err := r_kvstore.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the RKvstoreclient: %#v", err)
		}
		r_kvstoreConn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		r_kvstoreConn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		r_kvstoreConn.SourceIp = client.config.SourceIp
		r_kvstoreConn.SecureTransport = client.config.SecureTransport
		r_kvstoreConn.AppendUserAgent(Terraform, client.config.TerraformVersion)
		r_kvstoreConn.AppendUserAgent(Provider, providerVersion)
		r_kvstoreConn.AppendUserAgent(Module, client.config.ConfigurationSource)
		r_kvstoreConn.AppendUserAgent(TerraformTraceId, client.config.TerraformTraceId)
		client.r_kvstoreConn = r_kvstoreConn
	} else {
		err := client.r_kvstoreConn.InitWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the Redis client: %#v", err)
		}
	}
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

func (client *AliyunClient) NewEmrClient() (*rpc.Client, error) {
	productCode := "emr"
	endpoint := ""
	if v, ok := client.config.Endpoints.Load(productCode); !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			endpoint = fmt.Sprintf("emr.%s.aliyuncs.com", client.config.RegionId)
			client.config.Endpoints.Store(productCode, endpoint)
			log.Printf("[ERROR] loading %s endpoint got an error: %#v. Using the endpoint %s instead.", productCode, err, endpoint)
		}
	}
	if v, ok := client.config.Endpoints.Load(productCode); ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewDdsClient() (*rpc.Client, error) {
	productCode := "dds"
	endpoint := ""

	if v, ok := client.config.Endpoints.Load(productCode); !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			endpoint = fmt.Sprintf("mongodb.%s.aliyuncs.com", client.config.RegionId)
			client.config.Endpoints.Store(productCode, endpoint)
			log.Printf("[ERROR] loading %s endpoint got an error: %#v. Using the endpoint %s instead.", productCode, err, endpoint)
		}
	}

	if v, ok := client.config.Endpoints.Load(productCode); ok && v.(string) != "" {
		endpoint = v.(string)
	}

	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}

	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}

	return conn, nil
}

func (client *AliyunClient) NewBpstudioClient() (*rpc.Client, error) {
	productCode := "bpstudio"
	endpoint := ""
	if v, ok := client.config.Endpoints.Load(productCode); !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			endpoint = fmt.Sprintf("bpstudio.%s.aliyuncs.com", client.config.RegionId)
			client.config.Endpoints.Store(productCode, endpoint)
			log.Printf("[ERROR] loading %s endpoint got an error: %#v. Using the endpoint %s instead.", productCode, err, endpoint)
		}
	}
	if v, ok := client.config.Endpoints.Load(productCode); ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}
func (client *AliyunClient) NewAckoneClient() (*rpc.Client, error) {
	productCode := "adcp"
	endpoint := ""
	if v, ok := client.config.Endpoints.Load(productCode); !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			endpoint = fmt.Sprintf("adcp.%s.aliyuncs.com", client.config.RegionId)
			client.config.Endpoints.Store(productCode, endpoint)
			log.Printf("[ERROR] loading %s endpoint got an error: %#v. Using the endpoint %s instead.", productCode, err, endpoint)
		}
	}
	if v, ok := client.config.Endpoints.Load(productCode); ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}
func (client *AliyunClient) NewSlsClient() (*openapi.Client, error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(client.config.AccessKey),
		AccessKeySecret: tea.String(client.config.SecretKey),
		SecurityToken:   tea.String(client.config.SecurityToken),
		UserAgent:       tea.String(client.config.getUserAgent()),
	}

	endpoint := client.config.LogEndpoint
	if endpoint == "" {
		endpoint = loadEndpoint(client.config.RegionId, LOGCode)
		if endpoint == "" {
			endpoint = fmt.Sprintf("%s.log.aliyuncs.com", client.config.RegionId)
		}
	}
	endpoint = strings.TrimPrefix(endpoint, "https://")
	endpoint = strings.TrimPrefix(endpoint, "http://")

	config.Endpoint = tea.String(endpoint)
	openapiClient, _err := openapi.NewClient(config)
	if _err != nil {
		return nil, _err
	}
	openapiClient.Spi, _err = gatewayclient.NewClient()
	if _err != nil {
		return nil, _err
	}
	openapiClient.Protocol = tea.String(client.config.Protocol)

	return openapiClient, nil
}

func (client *AliyunClient) NewHologramClient() (*roa.Client, error) {
	productCode := "hologram"
	endpoint := ""
	if v, ok := client.config.Endpoints.Load(productCode); !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			endpoint = fmt.Sprintf("hologram.%s.aliyuncs.com", client.config.RegionId)
			client.config.Endpoints.Store(productCode, endpoint)
			log.Printf("[ERROR] loading %s endpoint got an error: %#v. Using the endpoint %s instead.", productCode, err, endpoint)
		}
	}
	if v, ok := client.config.Endpoints.Load(productCode); ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaRoaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := roa.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}
func (client *AliyunClient) NewRealtimecomputeClient() (*rpc.Client, error) {
	productCode := "foasconsole"
	endpoint := ""
	if v, ok := client.config.Endpoints.Load(productCode); !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			endpoint = "foasconsole.aliyuncs.com"
			client.config.Endpoints.Store(productCode, endpoint)
			log.Printf("[ERROR] loading %s endpoint got an error: %#v. Using the endpoint %s instead.", productCode, err, endpoint)
		}
	}
	if v, ok := client.config.Endpoints.Load(productCode); ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}
func (client *AliyunClient) NewAckClient() (*roa.Client, error) {
	productCode := "cs"
	endpoint := ""
	if v, ok := client.config.Endpoints.Load(productCode); !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			endpoint = fmt.Sprintf("cs.%s.aliyuncs.com", client.config.RegionId)
			client.config.Endpoints.Store(productCode, endpoint)
			log.Printf("[ERROR] loading %s endpoint got an error: %#v. Using the endpoint %s instead.", productCode, err, endpoint)
		}
	}
	if v, ok := client.config.Endpoints.Load(productCode); ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}

	sdkConfig := client.teaRoaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := roa.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewOssClient() (*openapi.Client, error) {
	config := &client.teaRoaOpenapiConfig

	productCode := "oss"
	endpoint := ""
	if v, ok := client.config.Endpoints.Load(productCode); !ok || v.(string) == "" {
		// Firstly, load endpoint from provider
		endpoint = client.config.OssEndpoint
		if endpoint == "" {
			// Secondly, load endpoint from environment
			endpoint = loadEndpoint(client.config.RegionId, OSSCode)
		}
		if endpoint == "" {
			// Thirdly, load endpoint from common method
			if err := client.loadEndpoint(productCode); err != nil {
				endpoint = fmt.Sprintf("oss-%s.aliyuncs.com", client.config.RegionId)
				client.config.Endpoints.Store(productCode, endpoint)
				log.Printf("[ERROR] loading %s endpoint got an error: %#v. Using the endpoint %s instead.", productCode, err, endpoint)
			}
		} else {
			client.config.Endpoints.Store(productCode, endpoint)
		}
	}
	if v, ok := client.config.Endpoints.Load(productCode); ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}

	config.Endpoint = tea.String(endpoint)
	openapiClient, _err := openapi.NewClient(config)
	if _err != nil {
		return nil, _err
	}
	openapiClient.Spi, _err = ossclient.NewClient()
	if _err != nil {
		return nil, _err
	}

	// SignVersion
	if ossV, ok := client.config.SignVersion.Load("oss"); ok {
		openapiClient.SignatureVersion = tea.String(ossV.(string))
	}

	return openapiClient, nil
}

func (client *AliyunClient) loadApiEndpoint(productCode string) (string, error) {
	accountId, err := client.AccountId()
	if err != nil {
		log.Printf("[WARN] failed to load accountId: %#v", err)
	}
	if v, ok := client.config.Endpoints.Load(productCode); !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			return "", fmt.Errorf("[ERROR] loading %s endpoint got an error: %#v.", productCode, err)
		}
	} else {
		return FormatEndpointWithAccountID(productCode, v.(string), accountId), nil
	}
	if v, ok := client.config.Endpoints.Load(productCode); ok && v.(string) != "" {
		return FormatEndpointWithAccountID(productCode, v.(string), accountId), nil
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
//	apiName - API Name
//	query - API parameters in query
//	headers - API parameters in headers
//	body - API parameters in body
//	autoRetry - whether to auto retry while the runtime has a 5xx error
func (client *AliyunClient) RoaPost(apiProductCode string, apiVersion string, apiName string, query map[string]*string, headers map[string]*string, body interface{}, autoRetry bool) (map[string]interface{}, error) {
	return client.roaRequest("POST", apiProductCode, apiVersion, apiName, query, headers, body, autoRetry)
}

// RoaPut invoking ROA API request with PUT method
// parameters:
//
//	apiProductCode: API Product code, its value equals to the gateway code of the API
//	apiVersion - API version
//	apiName - API Name
//	query - API parameters in query
//	headers - API parameters in headers
//	body - API parameters in body
//	autoRetry - whether to auto retry while the runtime has a 5xx error
func (client *AliyunClient) RoaPut(apiProductCode string, apiVersion string, apiName string, query map[string]*string, headers map[string]*string, body interface{}, autoRetry bool) (map[string]interface{}, error) {
	return client.roaRequest("PUT", apiProductCode, apiVersion, apiName, query, headers, body, autoRetry)
}

// RoaGet invoking ROA API request with GET method
// parameters:
//
//	apiProductCode: API Product code, its value equals to the gateway code of the API
//	apiVersion - API version
//	apiName - API Name
//	query - API parameters in query
//	headers - API parameters in headers
//	body - API parameters in body
func (client *AliyunClient) RoaGet(apiProductCode string, apiVersion string, apiName string, query map[string]*string, headers map[string]*string, body interface{}) (map[string]interface{}, error) {
	return client.roaRequest("GET", apiProductCode, apiVersion, apiName, query, headers, body, true)
}

// RoaDelete invoking ROA API request with DELETE method
// parameters:
//
//	apiProductCode: API Product code, its value equals to the gateway code of the API
//	apiVersion - API version
//	apiName - API Name
//	query - API parameters in query
//	headers - API parameters in headers
//	body - API parameters in body
//	autoRetry - whether to auto retry while the runtime has a 5xx error
func (client *AliyunClient) RoaDelete(apiProductCode string, apiVersion string, apiName string, query map[string]*string, headers map[string]*string, body interface{}, autoRetry bool) (map[string]interface{}, error) {
	return client.roaRequest("DELETE", apiProductCode, apiVersion, apiName, query, headers, body, autoRetry)
}

// RoaPatch invoking ROA API request with PATCH method
// parameters:
//
//	apiProductCode: API Product code, its value equals to the gateway code of the API
//	apiVersion - API version
//	apiName - API Name
//	query - API parameters in query
//	headers - API parameters in headers
//	body - API parameters in body
//	autoRetry - whether to auto retry while the runtime has a 5xx error
func (client *AliyunClient) RoaPatch(apiProductCode string, apiVersion string, apiName string, query map[string]*string, headers map[string]*string, body interface{}, autoRetry bool) (map[string]interface{}, error) {
	return client.roaRequest("PATCH", apiProductCode, apiVersion, apiName, query, headers, body, autoRetry)
}

func (client *AliyunClient) roaRequest(method string, apiProductCode string, apiVersion string, apiName string, query map[string]*string, headers map[string]*string, body interface{}, autoRetry bool) (map[string]interface{}, error) {
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
	runtime := &util.RuntimeOptions{}
	runtime.SetAutoretry(autoRetry)
	response, err := conn.DoRequest(tea.String(apiVersion), nil, tea.String(method), tea.String("AK"), tea.String(apiName), query, headers, body, runtime)
	if respBody, isExist := response["body"]; isExist && respBody != nil {
		response = respBody.(map[string]interface{})
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

type ossCredentials struct {
	client *AliyunClient
}

func (defCre *ossCredentials) GetAccessKeyID() string {
	value, err := defCre.client.teaSdkConfig.Credential.GetAccessKeyId()
	if err == nil && value != nil {
		return *value
	}
	return defCre.client.config.AccessKey
}

func (defCre *ossCredentials) GetAccessKeySecret() string {
	value, err := defCre.client.teaSdkConfig.Credential.GetAccessKeySecret()
	if err == nil && value != nil {
		return *value
	}
	return defCre.client.config.SecretKey
}

func (defCre *ossCredentials) GetSecurityToken() string {
	value, err := defCre.client.teaSdkConfig.Credential.GetSecurityToken()
	if err == nil && value != nil {
		return *value
	}
	return defCre.client.config.SecurityToken
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

func (client *AliyunClient) GenRoaParam(action, method, version, path string) *openapi.Params {
	return &openapi.Params{
		Action:      tea.String(action),
		Version:     tea.String(version),
		Protocol:    tea.String(client.config.Protocol),
		Pathname:    tea.String(path),
		Method:      tea.String(method),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
}
func (client *AliyunClient) NewCloudcontrolClient() (*roa.Client, error) {
	productCode := "CloudControl"
	endpoint := ""
	if v, ok := client.config.Endpoints.Load(productCode); !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			endpoint = fmt.Sprintf("cloudcontrol.aliyuncs.com")
			client.config.Endpoints.Store(productCode, endpoint)
			log.Printf("[ERROR] loading %s endpoint got an error: %#v. Using the endpoint %s instead.", productCode, err, endpoint)
		}
	}
	if v, ok := client.config.Endpoints.Load(productCode); ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}

	sdkConfig := client.teaRoaSdkConfig
	sdkConfig.SetEndpoint(fmt.Sprintf("%s", endpoint))

	conn, err := roa.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}
