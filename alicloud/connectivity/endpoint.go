package connectivity

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/location"
)

// ServiceCode Load endpoints from endpoints.xml or environment variables to meet specified application scenario, like private cloud.
type ServiceCode string

const (
	MaxcomputeCode      = ServiceCode("MAXCOMPUTE")
	CmsCode             = ServiceCode("CMS")
	RKvstoreCode        = ServiceCode("RKVSTORE")
	OnsCode             = ServiceCode("ONS")
	DcdnCode            = ServiceCode("DCDN")
	MseCode             = ServiceCode("MSE")
	ActiontrailCode     = ServiceCode("ACTIONTRAIL")
	OosCode             = ServiceCode("OOS")
	EcsCode             = ServiceCode("ECS")
	NasCode             = ServiceCode("NAS")
	EciCode             = ServiceCode("ECI")
	DdoscooCode         = ServiceCode("DDOSCOO")
	BssopenapiCode      = ServiceCode("BSSOPENAPI")
	AlidnsCode          = ServiceCode("ALIDNS")
	ResourcemanagerCode = ServiceCode("RESOURCEMANAGER")
	WafOpenapiCode      = ServiceCode("WAFOPENAPI")
	DmsEnterpriseCode   = ServiceCode("DMSENTERPRISE")
	DnsCode             = ServiceCode("DNS")
	KmsCode             = ServiceCode("KMS")
	CbnCode             = ServiceCode("CBN")
	ECSCode             = ServiceCode("ECS")
	ESSCode             = ServiceCode("ESS")
	RAMCode             = ServiceCode("RAM")
	VPCCode             = ServiceCode("VPC")
	SLBCode             = ServiceCode("SLB")
	RDSCode             = ServiceCode("RDS")
	OSSCode             = ServiceCode("OSS")
	ONSCode             = ServiceCode("ONS")
	ALIKAFKACode        = ServiceCode("ALIKAFKA")
	CONTAINCode         = ServiceCode("CS")
	CRCode              = ServiceCode("CR")
	CDNCode             = ServiceCode("CDN")
	CMSCode             = ServiceCode("CMS")
	KMSCode             = ServiceCode("KMS")
	OTSCode             = ServiceCode("OTS")
	DNSCode             = ServiceCode("DNS")
	PVTZCode            = ServiceCode("PVTZ")
	LOGCode             = ServiceCode("LOG")
	FCCode              = ServiceCode("FC")
	DDSCode             = ServiceCode("DDS")
	GPDBCode            = ServiceCode("GPDB")
	STSCode             = ServiceCode("STS")
	KVSTORECode         = ServiceCode("KVSTORE")
	POLARDBCode         = ServiceCode("POLARDB")
	DATAHUBCode         = ServiceCode("DATAHUB")
	MNSCode             = ServiceCode("MNS")
	CLOUDAPICode        = ServiceCode("APIGATEWAY")
	DRDSCode            = ServiceCode("DRDS")
	LOCATIONCode        = ServiceCode("LOCATION")
	ELASTICSEARCHCode   = ServiceCode("ELASTICSEARCH")
	BSSOPENAPICode      = ServiceCode("BSSOPENAPI")
	DDOSCOOCode         = ServiceCode("DDOSCOO")
	DDOSBGPCode         = ServiceCode("DDOSBGP")
	SAGCode             = ServiceCode("SAG")
	EMRCode             = ServiceCode("EMR")
	CasCode             = ServiceCode("CAS")
	YUNDUNDBAUDITCode   = ServiceCode("YUNDUNDBAUDIT")
	MARKETCode          = ServiceCode("MARKET")
	HBASECode           = ServiceCode("HBASE")
	ADBCode             = ServiceCode("ADB")
	MAXCOMPUTECode      = ServiceCode("MAXCOMPUTE")
	EDASCode            = ServiceCode("EDAS")
	CassandraCode       = ServiceCode("CASSANDRA")
)

type Endpoints struct {
	Endpoint []Endpoint `xml:"Endpoint"`
}

type Endpoint struct {
	Name      string    `xml:"name,attr"`
	RegionIds RegionIds `xml:"RegionIds"`
	Products  Products  `xml:"Products"`
}

type RegionIds struct {
	RegionId string `xml:"RegionId"`
}

type Products struct {
	Product []Product `xml:"Product"`
}

type Product struct {
	ProductName string `xml:"ProductName"`
	DomainName  string `xml:"DomainName"`
}

var localEndpointPath = "./endpoints.xml"
var localEndpointPathEnv = "TF_ENDPOINT_PATH"
var loadLocalEndpoint = false

func hasLocalEndpoint() bool {
	data, err := ioutil.ReadFile(localEndpointPath)
	if err != nil || len(data) <= 0 {
		d, e := ioutil.ReadFile(os.Getenv(localEndpointPathEnv))
		if e != nil {
			return false
		}
		data = d
	}
	return len(data) > 0
}

func LoadRegionalEndpoint(region string, serviceCode string) string {
	if region == "" || serviceCode == "" {
		return ""
	}
	return fmt.Sprintf("%s.%s.aliyuncs.com", serviceCode, region)
}

func loadEndpoint(region string, serviceCode ServiceCode) string {
	endpoint := strings.TrimSpace(os.Getenv(fmt.Sprintf("%s_ENDPOINT", string(serviceCode))))
	if endpoint != "" {
		return endpoint
	}

	// Load current path endpoint file endpoints.xml, if failed, it will load from environment variables TF_ENDPOINT_PATH
	if !loadLocalEndpoint {
		return ""
	}
	data, err := ioutil.ReadFile(localEndpointPath)
	if err != nil || len(data) <= 0 {
		d, e := ioutil.ReadFile(os.Getenv(localEndpointPathEnv))
		if e != nil {
			return ""
		}
		data = d
	}
	var endpoints Endpoints
	err = xml.Unmarshal(data, &endpoints)
	if err != nil {
		return ""
	}
	for _, endpoint := range endpoints.Endpoint {
		if endpoint.RegionIds.RegionId == string(region) {
			for _, product := range endpoint.Products.Product {
				if strings.ToLower(product.ProductName) == strings.ToLower(string(serviceCode)) {
					return strings.TrimSpace(product.DomainName)
				}
			}
		}
	}

	return ""
}

// productCodeToLocationCode records all products' code mapping to location
// Key: product code, its value equals to the gateway code of the API after converting it to lowercase and using underscores
// Value: location code
var productCodeToLocationCode = map[string]string{
	"ecs":                  "ecs",     // ECS
	"adb":                  "ads",     //ADB
	"ess":                  "ess",     //AutoScaling
	"cs":                   "cs",      // ACK
	"polardb":              "polardb", // PolarDB
	"cr":                   "acr",     // CR
	"dds":                  "dds",     //MongoDB
	"gpdb":                 "gpdb",    //GPDB
	"fc_open":              "fc",      // FC, FCV2
	"fc":                   "fc",      // FCV3
	"cloudapi":             "apigateway",
	"apig":                 "nativeapigw",     // APIG
	"datahub":              "datahub",         // DataHub
	"mns_open":             "mns",             // MessageService
	"elasticsearch":        "elasticsearch",   // Elasticsearch
	"ddoscoo":              "ddoscoo",         // DdosCoo
	"ddosbgp":              "ddosbgp",         // DdosBgp
	"antiddos_public":      "ddosbasic",       // DdosBasic
	"bssopenapi":           "bssopenapi",      //BssOpenApi
	"alikafka":             "alikafka",        //AliKafka
	"emr":                  "emr",             //EMR
	"smartag":              "smartag",         // Smartag
	"yundun_dbaudit":       "dbaudit",         //DBAudit
	"yundun_bastionhost":   "bastionhost",     //Bastionhost
	"hbase":                "hbase",           //HBase
	"edas":                 "edas",            // EDAS
	"alidns":               "alidns",          //Alidns
	"cassandra":            "cds",             //Cassandra
	"eci":                  "eci",             // ECI
	"dcdn":                 "dcdn",            // DCDN
	"r_kvstore":            "redisa",          // Redis
	"ons":                  "ons",             //Ons
	"config":               "config",          //Config
	"fnf":                  "fnf",             // FnF
	"ros":                  "ros",             // ROS
	"mse":                  "mse",             // MSE
	"pvtz":                 "pvtz",            //PrivateZone
	"privatelink":          "privatelink",     // PrivateLink
	"maxcompute":           "odps",            //MaxCompute
	"resourcesharing":      "ressharing",      // ResourceManager
	"ga":                   "gaplus",          // Ga
	"actiontrail":          "actiontrail",     //ActionTrail
	"hitsdb":               "hitsdb",          //Lindorm
	"brain_industrial":     "aistudio",        //BrainIndustrial
	"eipanycast":           "eipanycast",      // Eipanycast
	"oos":                  "oos",             // OOS
	"ims":                  "ims",             //IMS
	"resourcemanager":      "resourcemanager", // ResourceManager
	"nas":                  "nas",             //NAS
	"dms_enterprise":       "dmsenterprise",   //DMSEnterprise
	"sgw":                  "hcs_sgw",         // CloudStorageGateway
	"slb":                  "slb",             // SLB
	"kms":                  "kms",             //KMS
	"dm":                   "dm",              //DirectMail
	"eventbridge":          "eventbridge",     // EventBridge
	"hbr":                  "hbr",             //HBR
	"cas":                  "cas",             //SSLCertificatesService
	"arms":                 "arms",            // ARMS
	"cloudfw":              "cloudfirewall",   //CloudFirewall
	"sae":                  "serverless",      //SAE
	"alb":                  "alb",             // ALB
	"ecd":                  "gwsecd",          // ECD
	"cloudphone":           "cloudphone",      // ECP
	"scdn":                 "scdn",            //SCDN
	"dataworks_public":     "dide",            //DataWorks
	"cdn":                  "cdn",             // CDN
	"cddc":                 "cddc",            // CDDC
	"mscopensubscription":  "mscsub",          //MscSub
	"sddp":                 "sddp",            // SDDP
	"sas":                  "sas",             // ThreatDetection
	"ehpc":                 "ehs",             // Ehpc
	"ens":                  "ens",             // ENS
	"iot":                  "iot",             // Iot
	"imm":                  "imm",             // IMM
	"clickhouse":           "clickhouse",      // ClickHouse
	"selectdb":             "selectdb",        //SelectDB
	"dts":                  "dts",             // DTS
	"dg":                   "dg",              // DatabaseGateway
	"cloudsso":             "cloudsso",        // CloudSSO
	"swas_open":            "swas",            // SimpleApplicationServer
	"vs":                   "vs",              // VideoSurveillanceSystem
	"quickbi_public":       "quickbi",         // QuickBI
	"devops_rdc":           "rdcdevops",       // RDC
	"vod":                  "vod",             // VOD
	"opensearch":           "opensearch",      // OpenSearch
	"gdb":                  "gds",             // GraphDatabase
	"dbfs":                 "dbfs",            // DBFS
	"eais":                 "eais",            // EAIS
	"cloudauth":            "cloudauth",       // Cloudauth
	"imp":                  "imp",             // IMP
	"mhub":                 "emas",            // MHUB
	"servicemesh":          "servicemesh",     // ServiceMesh
	"eds_user":             "edsuser",         // ECD
	"tag":                  "tag",             // Tag
	"schedulerx2":          "edasschedulerx",  // Schedulerx
	"dysmsapi":             "dysms",           // SMS
	"vpcpeer":              "vpcpeer",         // VpcPeer
	"dbs":                  "cbs",             // DBS
	"nlb":                  "nlb",             // NLB
	"ebs":                  "ebs",             // EBS
	"bpstudio":             "bpstudio",        // BPStudio
	"das":                  "hdm",             // DAS
	"servicecatalog":       "srvcatalog",      // ServiceCatalog
	"eflo":                 "eflo",            //Eflo
	"oceanbasepro":         "oceanbase",       // OceanBase
	"chatbot":              "beebot",          // Chatbot
	"computenest":          "computenest",     // ComputeNest
	"drds":                 "drds",            // DRDS
	"polardbx":             "polardbx",        // DRDS
	"adcp":                 "adcp",            // AckOne
	"sls":                  "sls",             // SLS
	"rocketmq":             "rmq",             // RocketMQ
	"resourcecenter":       "",                // ResourceManager
	"hologram":             "hologram",        // Hologram
	"foasconsole":          "foasconsole",     // RealtimeCompute
	"vpc":                  "vpc",             // VPC, VPNGateway,ExpressConnect, CBWP, EIP
	"oss":                  "oss",             // OSS
	"cms":                  "cms",             // CloudMonitorService
	"waf_openapi":          "waf",             //WAFV3,WAF
	"dfs":                  "alidfs",          //DFS
	"amqp":                 "onsproxy",        // Amqp
	"amqp_open":            "onsproxy",        // Amqp
	"cbn":                  "cbn",             // CEN
	"expressconnectrouter": "ecr",             // ExpressConnectRouter
	"green":                "green",           // Aligreen
	"governance":           "governance",      // Governance
	"ots":                  "ots",             // OTS
	"tablestore":           "ots",             // OTS
	"ram":                  "ram",             //RAM
	"quotas":               "quotas",          //Quotas
	"market":               "market",          //Market
	"aiworkspace":          "paiworkspace",    //PAIWorkspace
	"vpcipam":              "vpcipam",         //VpcIpam
	"gwlb":                 "gwlb",            // GWLB
	"esa":                  "dcdnservices",    // ESA
	"live":                 "live",            // Live
	"eds_aic":              "wycloudphone",    // CloudPhone
	"cloudcontrol":         "cloudcontrol",    // CloudControl
}

// irregularProductEndpoint specially records those product codes that
// cannot be parsed out by the location service.
// Key: product code, its value equals to the gateway code of the API after converting it to lowercase and using underscores
// Value: product endpoint
// The priority of this configuration is higher than location service, lower than user environment variable configuration
var irregularProductEndpoint = map[string]string{
	"tablestore":       "tablestore.%s.aliyuncs.com",
	"ram":              "ram.aliyuncs.com",
	"brain_industrial": "brain-industrial.cn-hangzhou.aliyuncs.com",
	"cassandra":        "cassandra.aliyuncs.com",
	"cloudfw":          "cloudfw.aliyuncs.com",
	"scdn":             "scdn.aliyuncs.com",
	"vpcpeer":          "vpcpeer.aliyuncs.com",
	"resourcecenter":   "resourcecenter.aliyuncs.com",
	"market":           "market.aliyuncs.com",
	"bssopenapi":       BssOpenAPIEndpointDomestic,
	"esa":              "esa.cn-hangzhou.aliyuncs.com",
	"cas":              "cas.aliyuncs.com",
	"sas":              "tds.aliyuncs.com",
	"ros":              "ros.aliyuncs.com",
	"eds_aic":          "eds-aic.cn-shanghai.aliyuncs.com",
}

// irregularProductEndpointForIntlRegion specially records those product codes that
// cannot be parsed out by the location service and sensitive to region.
// These products adapt to international region, and conflict with irregularProductEndpointForIntlAccount
// Key: product code, its value equals to the gateway code of the API after converting it to lowercase and using underscores
// Value: product endpoint
// The priority of this configuration is higher than location service, lower than user environment variable configuration
var irregularProductEndpointForIntlRegion = map[string]string{
	"sas": SaSOpenAPIEndpointInternational,
}

// irregularProductEndpointForIntlAccount specially records those product codes that
// cannot be parsed out by the location service and sensitive to account type.
// These products adapt to international account.
// Key: product code, its value equals to the gateway code of the API after converting it to lowercase and using underscores
// Value: product endpoint
// The priority of this configuration is higher than location service, lower than user environment variable configuration
var irregularProductEndpointForIntlAccount = map[string]string{
	"cloudfw":        "cloudfw.ap-southeast-1.aliyuncs.com",
	"resourcecenter": "resourcecenter-intl.aliyuncs.com",
	"bssopenapi":     BssOpenAPIEndpointInternational,
	"esa":            "esa.ap-southeast-1.aliyuncs.com",
	"eds_aic":        "eds-aic.ap-southeast-1.aliyuncs.com",
	"ros":            "ros-intl.aliyuncs.com",
}

// irregularProductEndpointForIntlAccountIntlRegion specially records those product codes that
// cannot be parsed out by the location service and sensitive to account type and region.
// These products adapt to international account.
// Key: product code, its value equals to the gateway code of the API after converting it to lowercase and using underscores
// Value: product endpoint
// The priority of this configuration is higher than location service, lower than user environment variable configuration
var irregularProductEndpointForIntlAccountIntlRegion = map[string]string{
	"cas": "cas.ap-southeast-1.aliyuncs.com",
}

// regularProductEndpoint specially records those product codes that have been confirmed to be
// regional or central endpoints.
// Key: product code, its value equals to the gateway code of the API after converting it to lowercase and using underscores
// Value: product endpoint
// The priority of this configuration is lower than location service, and as a backup endpoint
var regularProductEndpoint = map[string]string{
	"ecs":                  "ecs.%s.aliyuncs.com",
	"mse":                  "mse.%s.aliyuncs.com",
	"vpc":                  "vpc.%s.aliyuncs.com",
	"oss":                  "oss-%s.aliyuncs.com",
	"cr":                   "cr.%s.aliyuncs.com",
	"cms":                  "metrics.%s.aliyuncs.com",
	"sls":                  "%s.log.aliyuncs.com",
	"drds":                 "drds.%s.aliyuncs.com",
	"polardbx":             "polardbx.%s.aliyuncs.com",
	"fc_open":              "%s.fc.aliyuncs.com",
	"fc":                   "%s.fc.aliyuncs.com",
	"cloudapi":             "apigateway.%s.aliyuncs.com",
	"mns_open":             "mns-open.%s.aliyuncs.com",
	"elasticsearch":        "elasticsearch.%s.aliyuncs.com",
	"alikafka":             "alikafka.%s.aliyuncs.com",
	"emr":                  "emr.%s.aliyuncs.com",
	"smartag":              "smartag.%s.aliyuncs.com",
	"alidns":               "alidns.%s.aliyuncs.com",
	"eci":                  "eci.%s.aliyuncs.com",
	"ons":                  "ons.%s.aliyuncs.com",
	"pvtz":                 "pvtz.aliyuncs.com",
	"privatelink":          "privatelink.%s.aliyuncs.com",
	"maxcompute":           "maxcompute.%s.aliyuncs.com",
	"resourcesharing":      "resourcesharing.%s.aliyuncs.com",
	"actiontrail":          "actiontrail.%s.aliyuncs.com",
	"hitsdb":               "hitsdb.%s.aliyuncs.com",
	"oos":                  "oos.%s.aliyuncs.com",
	"nas":                  "nas.%s.aliyuncs.com",
	"dms_enterprise":       "dms-enterprise.%s.aliyuncs.com",
	"kms":                  "kms.%s.aliyuncs.com",
	"eventbridge":          "eventbridge-console.%s.aliyuncs.com",
	"hbr":                  "hbr.%s.aliyuncs.com",
	"arms":                 "arms.%s.aliyuncs.com",
	"sae":                  "sae.%s.aliyuncs.com",
	"alb":                  "alb.%s.aliyuncs.com",
	"ecd":                  "ecd.%s.aliyuncs.com",
	"cloudphone":           "cloudphone.%s.aliyuncs.com",
	"dataworks_public":     "dataworks.%s.aliyuncs.com",
	"ehpc":                 "ehpc.%s.aliyuncs.com",
	"ens":                  "ens.aliyuncs.com",
	"iot":                  "iot.%s.aliyuncs.com",
	"imm":                  "imm.%s.aliyuncs.com",
	"swas_open":            "swas.%s.aliyuncs.com",
	"vs":                   "vs.%s.aliyuncs.com",
	"vod":                  "vod.%s.aliyuncs.com",
	"opensearch":           "opensearch.%s.aliyuncs.com",
	"dbfs":                 "dbfs.%s.aliyuncs.com",
	"eais":                 "eais.%s.aliyuncs.com",
	"servicemesh":          "servicemesh.aliyuncs.com",
	"tag":                  "tag.%s.aliyuncs.com",
	"schedulerx2":          "schedulerx.%s.aliyuncs.com",
	"dbs":                  "dbs-api.%s.aliyuncs.com",
	"nlb":                  "nlb.%s.aliyuncs.com",
	"ebs":                  "ebs.%s.aliyuncs.com",
	"eflo":                 "eflo.%s.aliyuncs.com",
	"oceanbasepro":         "oceanbasepro.%s.aliyuncs.com",
	"adcp":                 "adcp.%s.aliyuncs.com",
	"rocketmq":             "rocketmq.%s.aliyuncs.com",
	"hologram":             "hologram.%s.aliyuncs.com",
	"foasconsole":          "foasconsole.aliyuncs.com",
	"cs":                   "cs.%s.aliyuncs.com",
	"waf_openapi":          "wafopenapi.cn-hangzhou.aliyuncs.com",
	"dfs":                  "dfs.%s.aliyuncs.com",
	"amqp":                 "amqp-open.%s.aliyuncs.com",
	"amqp_open":            "amqp-open.%s.aliyuncs.com",
	"cbn":                  "cbn.aliyuncs.com",
	"expressconnectrouter": "expressconnectrouter.cn-shanghai.aliyuncs.com",
	"green":                "green.%s.aliyuncs.com",
	"governance":           "governance.cn-hangzhou.aliyuncs.com",
	"dysmsapi":             "dysmsapi.aliyuncs.com",
	"sddp":                 "sddp.cn-zhangjiakou.aliyuncs.com",
	"ddoscoo":              "ddoscoo.cn-hangzhou.aliyuncs.com",
	"config":               "config.cn-shanghai.aliyuncs.com",
	"ga":                   "ga.cn-hangzhou.aliyuncs.com",
	"dcdn":                 "dcdn.aliyuncs.com",
	"cdn":                  "cdn.aliyuncs.com",
	"cloudauth":            "cloudauth.aliyuncs.com",
	"ims":                  "ims.aliyuncs.com",
	"mhub":                 "mhub.cn-shanghai.aliyuncs.com",
	"eds_user":             "eds-user.cn-shanghai.aliyuncs.com",
	"eipanycast":           "eipanycast.cn-hangzhou.aliyuncs.com",
	"mscopensubscription":  "mscopensubscription.aliyuncs.com",
	"resourcemanager":      "resourcemanager.aliyuncs.com",
	"quotas":               "quotas.aliyuncs.com",
	"imp":                  "imp.aliyuncs.com",
	"das":                  "das.cn-shanghai.aliyuncs.com",
	"servicecatalog":       "servicecatalog.cn-hangzhou.aliyuncs.com",
	"chatbot":              "chatbot.cn-shanghai.aliyuncs.com",
	"computenest":          "computenest.cn-hangzhou.aliyuncs.com",
	"aiworkspace":          "aiworkspace.%s.aliyuncs.com",
	"vpcipam":              "vpcipam.%s.aliyuncs.com",
	"gwlb":                 "gwlb.%s.aliyuncs.com",
	"live":                 "live.aliyuncs.com",
	"dts":                  "dts.%s.aliyuncs.com",
	"dg":                   "dg.%s.aliyuncs.com",
	"cloudsso":             "cloudsso.%s.aliyuncs.com",
	"quickbi_public":       "quickbi.%s.aliyuncs.com",
	"ddosbgp":              "ddosbgp.%s.aliyuncs.com",
	"apig":                 "apig.%s.aliyuncs.com",
	"dds":                  "mongodb.%s.aliyuncs.com",
	"cloudcontrol":         "cloudcontrol.aliyuncs.com",
}

// regularProductEndpointForIntlRegion specially records those product codes that have been confirmed to be
// regional or central endpoints. But the endpoints are sensitive to region.
// These products adapt to international region, and conflict with regularProductEndpointForIntlAccount
// Key: product code, its value equals to the gateway code of the API after converting it to lowercase and using underscores
// Value: product endpoint
// The priority of this configuration is lower than location service, and as a backup endpoint
var regularProductEndpointForIntlRegion = map[string]string{
	"ddoscoo":      "ddoscoo.ap-southeast-1.aliyuncs.com",
	"eds_user":     "eds-user.ap-southeast-1.aliyuncs.com",
	"dysmsapi":     "dysmsapi.ap-southeast-1.aliyuncs.com",
	"sddp":         "sddp.ap-southeast-1.aliyuncs.com",
	"governance":   "governance.ap-southeast-1.aliyuncs.com",
	"waf_openapi":  "wafopenapi.ap-southeast-1.aliyuncs.com",
	"cloudcontrol": "cloudcontrol.ap-southeast-1.aliyuncs.com",
}

// regularProductEndpointForIntlAccount specially records those product codes that have been confirmed to be
// regional or central endpoints. But the endpoints are sensitive to account type.
// These products adapt to international account.
// Key: product code, its value equals to the gateway code of the API after converting it to lowercase and using underscores
// Value: product endpoint
// The priority of this configuration is lower than location service, and as a backup endpoint
var regularProductEndpointForIntlAccount = map[string]string{
	"config":              "config.ap-southeast-1.aliyuncs.com",
	"mscopensubscription": "mscopensubscription.ap-southeast-1.aliyuncs.com",
}

// regularProductEndpointForIntlAccountIntlRegion specially records those product codes that have been confirmed to be
// regional or central endpoints. But the endpoints are sensitive to account type and region.
// These products adapt to international account.
// Key: product code, its value equals to the gateway code of the API after converting it to lowercase and using underscores
// Value: product endpoint
// The priority of this configuration is lower than location service, and as a backup endpoint
var regularProductEndpointForIntlAccountIntlRegion = map[string]string{}

// regularProductEndpointReplace specially records some endpoints need to be replaced results from different reasons.
// Key: source endpoint
// Value: replaced endpoint
// The priority of this configuration is lower than other endpoints mapping rules, and as a backup endpoint
var regularProductEndpointReplace = map[string]string{
	// ecs, ecs.aliyuncs.com has more higher speed
	"ecs-cn-hangzhou.aliyuncs.com": "ecs.aliyuncs.com",
}

// NOTE: The productCode must be lowed.
func (client *AliyunClient) loadEndpoint(productCode string) error {
	// Firstly, load endpoint from environment variables
	endpoint := strings.TrimSpace(os.Getenv(fmt.Sprintf("ALIBABA_CLOUD_ENDPOINT_%s", strings.ToUpper(productCode))))
	if endpoint == "" {
		// Compatible with the previous implementation method
		endpoint = strings.TrimSpace(os.Getenv(fmt.Sprintf("%s_ENDPOINT", strings.ToUpper(productCode))))
	}
	if endpoint != "" {
		client.config.Endpoints.Store(productCode, endpoint)
		return nil
	}

	// Secondly, load endpoint from known rules
	if endpointFmt, ok := irregularProductEndpoint[productCode]; ok {
		if v, ok := irregularProductEndpointForIntlRegion[productCode]; ok && client.isInternationalRegion() {
			endpointFmt = v
		}
		if v, ok := irregularProductEndpointForIntlAccount[productCode]; ok && client.IsInternationalAccount() {
			endpointFmt = v
		}
		if v, ok := irregularProductEndpointForIntlAccountIntlRegion[productCode]; ok && client.IsInternationalAccount() && client.isInternationalRegion() {
			endpointFmt = v
		}
		if strings.Contains(endpointFmt, "%s") {
			endpointFmt = fmt.Sprintf(endpointFmt, client.RegionId)
		}
		client.config.Endpoints.Store(productCode, endpointFmt)
		return nil
	}

	// Thirdly, load endpoint from location
	endpoint, err := client.describeEndpointForService(productCode)
	if err == nil {
		if v, ok := regularProductEndpointForIntlAccount[productCode]; ok && client.IsInternationalAccount() {
			endpoint = v
		}
		if v, ok := regularProductEndpointReplace[endpoint]; ok {
			endpoint = v
		}
		client.config.Endpoints.Store(strings.ToLower(productCode), endpoint)
	} else if endpointFmt, ok := regularProductEndpoint[productCode]; ok {
		if v, ok := regularProductEndpointForIntlRegion[productCode]; ok && client.isInternationalRegion() {
			endpointFmt = v
		}
		if v, ok := regularProductEndpointForIntlAccount[productCode]; ok && client.IsInternationalAccount() {
			endpointFmt = v
		}
		if v, ok := regularProductEndpointForIntlAccountIntlRegion[productCode]; ok && client.IsInternationalAccount() && client.isInternationalRegion() {
			endpointFmt = v
		}
		if strings.Contains(endpointFmt, "%s") {
			endpointFmt = fmt.Sprintf(endpointFmt, client.RegionId)
		}
		if v, ok := regularProductEndpointReplace[endpointFmt]; ok {
			endpointFmt = v
		}
		client.config.Endpoints.Store(productCode, endpointFmt)
		log.Printf("[WARN] loading %s endpoint got an error: %#v. Using the endpoint %s instead.", productCode, err, endpointFmt)
		return nil
	}
	return err
}

// Load current path endpoint file endpoints.xml, if failed, it will load from environment variables TF_ENDPOINT_PATH
func (config *Config) loadEndpointFromLocal() error {
	data, err := ioutil.ReadFile(localEndpointPath)
	if err != nil || len(data) <= 0 {
		d, e := ioutil.ReadFile(os.Getenv(localEndpointPathEnv))
		if e != nil {
			return e
		}
		data = d
	}
	var endpoints Endpoints
	err = xml.Unmarshal(data, &endpoints)
	if err != nil {
		return err
	}
	for _, endpoint := range endpoints.Endpoint {
		if endpoint.RegionIds.RegionId == string(config.RegionId) {
			for _, product := range endpoint.Products.Product {
				config.Endpoints.Store(strings.ToLower(product.ProductName), strings.TrimSpace(product.DomainName))
			}
		}
	}
	return nil
}

func FormatEndpointWithAccountID(productCode string, endpoint string, accountId string) string {
	switch productCode {
	case "fc_open", "fc":
		return fmt.Sprintf("%s.%s", accountId, endpoint)
	}
	return endpoint
}
func incrementalWait(firstDuration time.Duration, increaseDuration time.Duration) func() {
	retryCount := 1
	return func() {
		var waitTime time.Duration
		if retryCount == 1 {
			waitTime = firstDuration
		} else if retryCount > 1 {
			waitTime += increaseDuration
		}
		time.Sleep(waitTime)
		retryCount++
	}
}
func (client *AliyunClient) describeEndpointForService(productCode string) (string, error) {
	locationCode := productCodeToLocationCode[productCode]
	if locationCode == "" {
		locationCode = productCode
	}
	args := location.CreateDescribeEndpointsRequest()
	args.ServiceCode = locationCode
	args.Id = client.config.RegionId
	args.Domain = client.config.LocationEndpoint
	if args.Domain == "" {
		args.Domain = loadEndpoint(client.RegionId, LOCATIONCode)
	}
	if args.Domain == "" {
		args.Domain = "location.aliyuncs.com"
	}

	locationClient, err := location.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(time.Duration(60)*time.Second), client.config.getAuthCredential(true))
	if err != nil {
		return "", fmt.Errorf("unable to initialize the location client: %#v", err)

	}
	locationClient.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
	locationClient.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
	locationClient.SourceIp = client.config.SourceIp
	locationClient.SecureTransport = client.config.SecureTransport
	defer locationClient.Shutdown()
	wait := incrementalWait(3*time.Second, 5*time.Second)
	var endpointResult string
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		endpointsResponse, err := locationClient.DescribeEndpoints(args)
		if err != nil {
			re := regexp.MustCompile("^Post [\"]*https://.*")
			if err.Error() != "" && re.MatchString(err.Error()) {
				wait()
				args.Domain = "location-readonly.aliyuncs.com"
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if endpointsResponse != nil && len(endpointsResponse.Endpoints.Endpoint) > 0 {
			for _, e := range endpointsResponse.Endpoints.Endpoint {
				if e.Type == "openAPI" {
					endpointResult = e.Endpoint
					return nil
				}
			}
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("Describe %s endpoint using region: %#v got an error: %#v.", productCode, client.RegionId, err)
	}
	if endpointResult == "" {
		return "", fmt.Errorf("There is no any available endpoint for %s in region %s.", productCode, client.RegionId)
	}
	return endpointResult, nil
}

const (
	OpenApiGatewayService          = "apigateway.cn-hangzhou.aliyuncs.com"
	OpenOtsService                 = "ots.cn-hangzhou.aliyuncs.com"
	OpenOssService                 = "oss-admin.aliyuncs.com"
	OpenNasService                 = "nas.cn-hangzhou.aliyuncs.com"
	OpenCdnService                 = "cdn.aliyuncs.com"
	OpenKmsService                 = "kms.cn-hangzhou.aliyuncs.com"
	OpenSaeService                 = "sae.cn-hangzhou.aliyuncs.com"
	OpenCmsService                 = "metrics.cn-hangzhou.aliyuncs.com"
	OpenDatahubService             = "datahub.aliyuncs.com"
	OpenOnsService                 = "ons.cn-hangzhou.aliyuncs.com"
	OpenDcdnService                = "dcdn.aliyuncs.com"
	OpenFcService                  = "fc-open.cn-hangzhou.aliyuncs.com"
	OpenAckService                 = "cs.aliyuncs.com"
	OpenPrivateLinkService         = "privatelink.cn-hangzhou.aliyuncs.com"
	OpenBrainIndustrialService     = "brain-industrial.cn-hangzhou.aliyuncs.com"
	OpenIotService                 = "iot.aliyuncs.com"
	OpenVsService                  = "vs.cn-shanghai.aliyuncs.com"
	OpenCrService                  = "cr.cn-hangzhou.aliyuncs.com"
	OpenMaxcomputeService          = "maxcompute.aliyuncs.com"
	OpenCloudStorageGatewayService = "sgw.cn-shanghai.aliyuncs.com"
	DataWorksService               = "dataworks.aliyuncs.com"
	OpenHbrService                 = "hbr.aliyuncs.com"
)

const (
	BssOpenAPIEndpointDomestic                = "business.aliyuncs.com"
	BssOpenAPIEndpointInternational           = "business.ap-southeast-1.aliyuncs.com"
	EcdOpenAPIEndpointUser                    = "eds-user.ap-southeast-1.aliyuncs.com"
	CloudFirewallOpenAPIEndpointControlPolicy = "cloudfw.ap-southeast-1.aliyuncs.com"
	SaSOpenAPIEndpointInternational           = "tds.ap-southeast-1.aliyuncs.com"
)
