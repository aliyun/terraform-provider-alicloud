package alicloud

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strconv"
	"strings"

	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/denverdino/aliyungo/common"
	"github.com/hashicorp/terraform/helper/schema"
)

type InstanceNetWork string

const (
	ClassicNet = InstanceNetWork("classic")
	VpcNet     = InstanceNetWork("vpc")
)

type PayType string

const (
	PrePaid  = PayType("PrePaid")
	PostPaid = PayType("PostPaid")
	Prepaid  = PayType("Prepaid")
	Postpaid = PayType("Postpaid")
)

type NetType string

const (
	Internet = NetType("Internet")
	Intranet = NetType("Intranet")
)

type NetworkType string

const (
	Classic = NetworkType("Classic")
	Vpc     = NetworkType("Vpc")
)

type TimeType string

const (
	Hour  = TimeType("Hour")
	Day   = TimeType("Day")
	Week  = TimeType("Week")
	Month = TimeType("Month")
	Year  = TimeType("Year")
)

type Status string

const (
	Pending     = Status("Pending")
	Creating    = Status("Creating")
	Running     = Status("Running")
	Available   = Status("Available")
	Unavailable = Status("Unavailable")
	Modifying   = Status("Modifying")
	Deleting    = Status("Deleting")
	Starting    = Status("Starting")
	Stopping    = Status("Stopping")
	Stopped     = Status("Stopped")

	Associating   = Status("Associating")
	Unassociating = Status("Unassociating")
	InUse         = Status("InUse")
	DiskInUse     = Status("In_use")

	Active   = Status("Active")
	Inactive = Status("Inactive")
	Idle     = Status("Idle")

	SoldOut = Status("SoldOut")

	InService      = Status("InService")
	Removing       = Status("Removing")
	DisabledStatus = Status("Disabled")
)

type IPType string

const (
	Inner   = IPType("Inner")
	Private = IPType("Private")
	Public  = IPType("Public")
)

type ResourceType string

const (
	ResourceTypeInstance = ResourceType("Instance")
	ResourceTypeDisk     = ResourceType("Disk")
	ResourceTypeVSwitch  = ResourceType("VSwitch")
	ResourceTypeRds      = ResourceType("Rds")
	IoOptimized          = ResourceType("IoOptimized")
)

type InternetChargeType string

const (
	PayByBandwidth = InternetChargeType("PayByBandwidth")
	PayByTraffic   = InternetChargeType("PayByTraffic")
)

// timeout for common product, ecs e.g.
const DefaultTimeout = 120

const DefaultTimeoutMedium = 500

// timeout for long time progerss product, rds e.g.
const DefaultLongTimeout = 1000

const DefaultIntervalShort = 5

const DefaultIntervalMedium = 10

const DefaultIntervalLong = 20

const (
	PageSizeSmall  = 10
	PageSizeMedium = 20
	PageSizeLarge  = 50
)

func getRegion(d *schema.ResourceData, meta interface{}) common.Region {
	return meta.(*AliyunClient).Region
}

func getRegionId(d *schema.ResourceData, meta interface{}) string {
	return meta.(*AliyunClient).RegionId
}

// Protocol represents network protocol
type Protocol string

// Constants of protocol definition
const (
	Http  = Protocol("http")
	Https = Protocol("https")
	Tcp   = Protocol("tcp")
	Udp   = Protocol("udp")
	All   = Protocol("all")
	Icmp  = Protocol("icmp")
	Gre   = Protocol("gre")
)

// ValidProtocols network protocol list
var ValidProtocols = []Protocol{Http, Https, Tcp, Udp}

// simple array value check method, support string type only
func isProtocolValid(value string) bool {
	res := false
	for _, v := range ValidProtocols {
		if string(v) == value {
			res = true
		}
	}
	return res
}

// default region for all resource
const DEFAULT_REGION = "cn-beijing"

const INT_MAX = 2147483647

// symbol of multiIZ
const MULTI_IZ_SYMBOL = "MAZ"

// default connect port of db
const DB_DEFAULT_CONNECT_PORT = "3306"

const COMMA_SEPARATED = ","

const COLON_SEPARATED = ":"

const DOT_SEPARATED = "."

const LOCAL_HOST_IP = "127.0.0.1"

// Takes the result of flatmap.Expand for an array of strings
// and returns a []string
func expandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		vs = append(vs, v.(string))
	}
	return vs
}

// Convert the result for an array and returns a Json string
func convertListToJsonString(configured []interface{}) string {
	if len(configured) < 1 {
		return ""
	}
	result := "["
	for i, v := range configured {
		result += "\"" + v.(string) + "\""
		if i < len(configured)-1 {
			result += ","
		}
	}
	result += "]"
	return result
}

const ServerSideEncryptionAes256 = "AES256"

type OptimizedType string

const (
	IOOptimized   = OptimizedType("optimized")
	NoneOptimized = OptimizedType("none")
)

type TagResourceType string

const (
	TagResourceImage    = TagResourceType("image")
	TagResourceInstance = TagResourceType("instance")
	TagResourceSnapshot = TagResourceType("snapshot")
	TagResourceDisk     = TagResourceType("disk")
)

func getPagination(pageNumber, pageSize int) (pagination common.Pagination) {
	pagination.PageSize = pageSize
	pagination.PageNumber = pageNumber
	return
}

const CharityPageUrl = "http://promotion.alicdn.com/help/oss/error.html"

func (client *AliyunClient) JudgeRegionValidation(key, region string) error {
	resp, err := client.ecsconn.DescribeRegions(ecs.CreateDescribeRegionsRequest())
	if err != nil {
		return fmt.Errorf("DescribeRegions got an error: %#v", err)
	}
	if resp == nil || len(resp.Regions.Region) < 1 {
		return GetNotFoundErrorFromString("There is no any available region.")
	}

	var rs []string
	for _, v := range resp.Regions.Region {
		if v.RegionId == region {
			return nil
		}
		rs = append(rs, v.RegionId)
	}
	return fmt.Errorf("'%s' is invalid. Expected on %v.", key, strings.Join(rs, ", "))
}

func userDataHashSum(user_data string) string {
	// Check whether the user_data is not Base64 encoded.
	// Always calculate hash of base64 decoded value since we
	// check against double-encoding when setting it
	v, base64DecodeError := base64.StdEncoding.DecodeString(user_data)
	if base64DecodeError != nil {
		v = []byte(user_data)
	}
	return string(v)
}

const DBConnectionSuffix = ".mysql.rds.aliyuncs.com"

// Remove useless blank in the string.
func Trim(v string) string {
	if len(v) < 1 {
		return v
	}
	return strings.Trim(v, " ")
}

// Load endpoints from endpoints.xml or environment variables to meet specified application scenario, like private cloud.
type ServiceCode string

const (
	ECSCode     = ServiceCode("ECS")
	ESSCode     = ServiceCode("ESS")
	RAMCode     = ServiceCode("RAM")
	VPCCode     = ServiceCode("VPC")
	SLBCode     = ServiceCode("SLB")
	RDSCode     = ServiceCode("RDS")
	OSSCode     = ServiceCode("OSS")
	CONTAINCode = ServiceCode("CS")
	DOMAINCode  = ServiceCode("DOMAIN")
	CDNCode     = ServiceCode("CDN")
	CMSCode     = ServiceCode("CMS")
	KMSCode     = ServiceCode("KMS")
	OTSCode     = ServiceCode("OTS")
	LOGCode     = ServiceCode("LOG")
)

//xml
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

func LoadEndpoint(region string, serviceCode ServiceCode) string {
	endpoint := strings.TrimSpace(os.Getenv(fmt.Sprintf("%s_ENDPOINT", string(serviceCode))))
	if endpoint != "" {
		return endpoint
	}

	// Load current path endpoint file endpoints.xml, if failed, it will load from environment variables TF_ENDPOINT_PATH
	data, err := ioutil.ReadFile("./endpoints.xml")
	if err != nil || len(data) <= 0 {
		d, e := ioutil.ReadFile(os.Getenv("TF_ENDPOINT_PATH"))
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
					return product.DomainName
				}
			}
		}
	}

	return ""
}

const ApiVersion20140526 = "2014-05-26"
const ApiVersion20140828 = "2014-08-28"

type CommonRequestDomain string

const (
	ECSDomain = CommonRequestDomain("ecs.aliyuncs.com")
	ESSDomain = CommonRequestDomain("ess.aliyuncs.com")
)

func CommonRequestInit(region string, code ServiceCode, domain CommonRequestDomain) *requests.CommonRequest {
	request := requests.NewCommonRequest()
	request.Version = ApiVersion20140526
	request.Domain = string(domain)
	d := LoadEndpoint(region, code)
	if d != "" {
		request.Domain = d
	}
	return request
}

func ConvertIntegerToInt(value requests.Integer) (v int, err error) {
	if strings.TrimSpace(string(value)) == "" {
		return
	}
	v, err = strconv.Atoi(string(value))
	if err != nil {
		return v, fmt.Errorf("Converting integer %s to int got an error: %#v.", value, err)
	}
	return
}

func GetUserHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("Get current user got an error: %#v.", err)
	}
	return usr.HomeDir, nil
}

func writeToFile(filePath string, data interface{}) error {
	if strings.HasPrefix(filePath, "~") {
		home, err := GetUserHomeDir()
		if err != nil {
			return err
		}
		if home != "" {
			filePath = strings.Replace(filePath, "~", home, 1)
		}
	}

	os.Remove(filePath)

	var out string
	switch data.(type) {
	case string:
		out = data.(string)
		break
	case nil:
		return nil
	default:
		bs, err := json.MarshalIndent(data, "", "\t")
		if err != nil {
			return fmt.Errorf("MarshalIndent data %#v got an error: %#v", data, err)
		}
		out = string(bs)
	}

	ioutil.WriteFile(filePath, []byte(out), 422)
	return nil
}

type Invoker struct {
	catchers []*Catcher
}

type Catcher struct {
	Reason           string
	RetryCount       int
	RetryWaitSeconds int
}

var ClientErrorCatcher = Catcher{AliyunGoClientFailure, 10, 3}
var ServiceBusyCatcher = Catcher{"ServiceUnavailable", 10, 3}

func NewInvoker() Invoker {
	i := Invoker{}
	i.AddCatcher(ClientErrorCatcher)
	i.AddCatcher(ServiceBusyCatcher)
	return i
}

func (a *Invoker) AddCatcher(catcher Catcher) {
	a.catchers = append(a.catchers, &catcher)
}

func (a *Invoker) Run(f func() error) error {
	err := f()

	if err == nil {
		return nil
	}

	for _, catcher := range a.catchers {
		if strings.Contains(err.Error(), catcher.Reason) {
			catcher.RetryCount--

			if catcher.RetryCount <= 0 {
				return fmt.Errorf("Retry timeout and got an error: %#v.", err)
			} else {
				time.Sleep(time.Duration(catcher.RetryWaitSeconds) * time.Second)
				return a.Run(f)
			}
		}
	}
	return err
}
