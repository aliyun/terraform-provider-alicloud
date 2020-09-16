package connectivity

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/location"
	"github.com/yalp/jsonpath"

	"encoding/json"
)

// ServiceCode Load endpoints from endpoints.xml or environment variables to meet specified application scenario, like private cloud.
type ServiceCode string

const (
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
	CENCode             = ServiceCode("CEN")
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

func (client *AliyunClient) loadEndpoint(productCode string) (string, error) {
	productCodeUp := strings.ToUpper(productCode)
	productCodeLow := strings.ToLower(productCode)
	config := client.config
	if config.Endpoints[productCodeLow] != nil && config.Endpoints[productCodeLow].(string) != "" {
		return config.Endpoints[productCodeLow].(string), nil
	}
	endpoint := strings.TrimSpace(os.Getenv(fmt.Sprintf("%s_ENDPOINT", productCodeUp)))
	if endpoint != "" {
		config.Endpoints[productCodeLow] = endpoint
		return endpoint, nil
	}

	// Load current path endpoint file endpoints.xml, if failed, it will load from environment variables TF_ENDPOINT_PATH
	if loadLocalEndpoint {
		data, err := ioutil.ReadFile(localEndpointPath)
		if err != nil || len(data) <= 0 {
			d, e := ioutil.ReadFile(os.Getenv(localEndpointPathEnv))
			if e != nil {
				return "", e
			}
			data = d
		}
		var endpoints Endpoints
		err = xml.Unmarshal(data, &endpoints)
		if err != nil {
			return "", err
		}
		for _, endpoint := range endpoints.Endpoint {
			if endpoint.RegionIds.RegionId == string(config.RegionId) {
				for _, product := range endpoint.Products.Product {
					if strings.ToLower(product.ProductName) == productCodeLow {
						config.Endpoints[productCodeLow] = strings.TrimSpace(product.DomainName)
						return strings.TrimSpace(product.DomainName), nil
					}
				}
			}
		}
	}

	// if not, get an endpoint by regional rule
	endpoint, err := loadEndpointFromSdk(client.config, productCodeLow)
	if err != nil {
		log.Fatalf("[ERROR] loadEndpoint from Sdk api got an error:%s", err)
		serviceCode := serviceCodeMapping[productCodeLow]
		if serviceCode == "" {
			serviceCode = productCodeLow
		}
		endpoint, err = client.describeEndpointForService(serviceCode)
		if err != nil {
			return "", err
		}
	}
	client.config.Endpoints[productCodeLow] = endpoint
	return endpoint, nil
}

func loadEndpointFromSdk(config *Config, productCode string) (string, error) {
	response, err := http.Post(fmt.Sprintf("http://sdk.aliyun-inc.com/api/get/release/endpoint/info?product_id=%s", productCode), "", nil)
	if err != nil {
		log.Fatalf("[ERROR] http.post got an error: %s", err)
		return "", err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("[ERROR] read http.post response got an error: %s", err)
		return "", err
	}

	var jsonBody interface{}
	if err := json.Unmarshal(body, &jsonBody); err != nil {
		log.Fatalf("[ERROR] json.Unmarshal http.post response body: %s got an error: %s", body, err)
		return "", err
	}
	ok, err := jsonpath.Read(jsonBody, "$.ok")
	if err != nil {
		return "", err
	}
	if ok.(bool) {
		endpointRegional, err := jsonpath.Read(jsonBody, "$.data.endpoint_regional")
		if err != nil {
			return "", err
		}
		if endpointRegional == "regional" {
			var err error
			for _, pathKey := range []string{"endpoint_map", "standard"} {
				endpoint, e := jsonpath.Read(jsonBody, fmt.Sprintf("$.data.endpoint_data.%s[\"%s\"]", pathKey, config.RegionId))
				if e != nil {
					log.Fatalf("[ERROR] jsonpath.Read endpoint got an error: %s", e)
					err = e
				}
				if endpoint != nil && endpoint.(string) != "" {
					return endpoint.(string), nil
				}
			}
			return "", err
		}
	}
	return "", nil
}
func (client *AliyunClient) describeEndpointForService(serviceCode string) (string, error) {
	args := location.CreateDescribeEndpointsRequest()
	args.ServiceCode = serviceCode
	args.Id = client.config.RegionId
	args.Domain = client.config.LocationEndpoint
	if args.Domain == "" {
		args.Domain = loadEndpoint(client.RegionId, LOCATIONCode)
	}
	if args.Domain == "" {
		args.Domain = "location-readonly.aliyuncs.com"
	}

	locationClient, err := location.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
	if err != nil {
		return "", fmt.Errorf("Unable to initialize the location client: %#v", err)

	}
	locationClient.AppendUserAgent(Terraform, terraformVersion)
	locationClient.AppendUserAgent(Provider, providerVersion)
	locationClient.AppendUserAgent(Module, client.config.ConfigurationSource)
	endpointsResponse, err := locationClient.DescribeEndpoints(args)
	if err != nil {
		return "", fmt.Errorf("Describe %s endpoint using region: %#v got an error: %#v.", serviceCode, client.RegionId, err)
	}
	if endpointsResponse != nil && len(endpointsResponse.Endpoints.Endpoint) > 0 {
		for _, e := range endpointsResponse.Endpoints.Endpoint {
			if e.Type == "openAPI" {
				return e.Endpoint, nil
			}
		}
	}
	return "", fmt.Errorf("There is no any available endpoint for %s in region %s.", serviceCode, client.RegionId)
}

var serviceCodeMapping = map[string]string{
	"cloudapi": "apigateway",
}

const (
	OpenApiGatewayService = "apigateway.cn-hangzhou.aliyuncs.com"
	OpenSlsService        = "sls.aliyuncs.com"
	OpenOtsService        = "ots.cn-hangzhou.aliyuncs.com"
	OpenOssService        = "oss-admin.aliyuncs.com"
	OpenNasService        = "nas.cn-hangzhou.aliyuncs.com"
)
