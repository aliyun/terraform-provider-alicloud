package alicloud

import (
	"fmt"
	"strings"

	"encoding/base64"

	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ecs"
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
	VPC     = NetworkType("VPC")
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

	Associating   = Status("Associating")
	Unassociating = Status("Unassociating")
	InUse         = Status("InUse")

	Active   = Status("Active")
	Inactive = Status("Inactive")
	Idle     = Status("Idle")
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

// Protocol represents network protocol
type Protocol string

// Constants of protocol definition
const (
	Http  = Protocol("http")
	Https = Protocol("https")
	Tcp   = Protocol("tcp")
	Udp   = Protocol("udp")
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

var DefaultBusinessInfo = ecs.BusinessInfo{
	Pack: "terraform",
}

// default region for all resource
const DEFAULT_REGION = "cn-beijing"

// default security ip for db
const DEFAULT_DB_SECURITY_IP = "127.0.0.1"

// we the count of create instance is only one
const DEFAULT_INSTANCE_COUNT = 1

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

type ResourceKeyType string

const (
	ZoneKey                       = ResourceKeyType("zones")
	InstanceTypeKey               = ResourceKeyType("instanceTypes")
	OutdatedInstanceTypeKey       = ResourceKeyType("outdatedInstanceTypes")
	UpgradedInstanceTypeKey       = ResourceKeyType("upgradedInstanceTypes")
	InstanceTypeFamilyKey         = ResourceKeyType("instanceTypeFamilies")
	OutdatedInstanceTypeFamilyKey = ResourceKeyType("outdatedInstanceTypeFamilies")
	UpgradedInstanceTypeFamilyKey = ResourceKeyType("upgradedInstanceTypeFamilies")
	DiskCategoryKey               = ResourceKeyType("diskCatetories")
	OutdatedDiskCategoryKey       = ResourceKeyType("outdatedDiskCatetories")
	IoOptimizedKey                = ResourceKeyType("optimized")
)

func getPagination(pageNumber, pageSize int) (pagination common.Pagination) {
	pagination.PageSize = pageSize
	pagination.PageNumber = pageNumber
	return
}

const CharityPageUrl = "http://promotion.alicdn.com/help/oss/error.html"

func (client *AliyunClient) JudgeRegionValidation(key string, region common.Region) error {
	regions, err := client.ecsconn.DescribeRegions()
	if err != nil {
		return fmt.Errorf("DescribeRegions got an error: %#v", err)
	}

	var rs []string
	for _, v := range regions {
		if v.RegionId == region {
			return nil
		}
		rs = append(rs, string(v.RegionId))
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
