package connectivity

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// ServiceCode Load endpoints from endpoints.xml or environment variables to meet specified application scenario, like private cloud.
type ServiceCode string

const (
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
