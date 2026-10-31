package alicloud

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"gopkg.in/yaml.v2"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
)

// The utilities below moved to the leaf packages helper, convert, tags and
// logs; these aliases keep the flat layout compiling unchanged.
//
// addDebug/debugOn must stay var aliases rather than forwarding
// functions: addDebug walks runtime.Caller(1..4) relative to its own
// frame, and a wrapper would shift the reported call sites by one.
// GetFunc is safe as a plain moved function because its caller level is
// an explicit argument, so no frame is inserted between the caller and
// the implementation.
//
// Invoker, Catcher, MyMap, Protocol and TagResourceType must be type
// aliases — the "=" is load-bearing: values constructed by helper (e.g.
// NewInvoker's return) must keep matching the types named at the flat
// call sites, or method calls and type switches silently stop working.
type Invoker = helper.Invoker
type Catcher = helper.Catcher
type MyMap = helper.MyMap
type Protocol = helper.Protocol
type TagResourceType = tags.TagResourceType

var (
	addDebug         = logs.AddDebug
	debugOn          = logs.DebugOn
	buildClientToken = helper.BuildClientToken
	incrementalWait  = helper.IncrementalWait
	writeToFile      = helper.WriteToFile
	GetUserHomeDir   = helper.GetUserHomeDir

	// conversion helpers (convert/convert.go)
	expandStringList                  = convert.ExpandStringList
	convertListStringToListInterface  = convert.ConvertListStringToListInterface
	expandIntList                     = convert.ExpandIntList
	convertListToJsonString           = convert.ConvertListToJsonString
	convertJsonStringToStringList     = convert.ConvertJsonStringToStringList
	encodeToBase64String              = convert.EncodeToBase64String
	decodeFromBase64String            = convert.DecodeFromBase64String
	convertJsonStringToMap            = convert.ConvertJsonStringToMap
	convertListToCommaSeparate        = convert.ConvertListToCommaSeparate
	filterEmptyStrings                = convert.FilterEmptyStrings
	convertBoolToString               = convert.ConvertBoolToString
	convertStringToBool               = convert.ConvertStringToBool
	convertIntergerToString           = convert.ConvertIntergerToString
	convertFloat64ToString            = convert.ConvertFloat64ToString
	convertJsonStringToList           = convert.ConvertJsonStringToList
	expandArrayToMap                  = convert.ExpandArrayToMap
	convertJsonStringToObject         = convert.ConvertJsonStringToObject
	convertObjectToJsonString         = convert.ConvertObjectToJsonString
	convertMaptoJsonString            = convert.ConvertMaptoJsonString
	convertMapToJsonStringIgnoreError = convert.ConvertMapToJsonStringIgnoreError
	convertInterfaceToJsonString      = convert.ConvertInterfaceToJsonString
	convertListMapToJsonString        = convert.ConvertListMapToJsonString
	convertMapFloat64ToJsonString     = convert.ConvertMapFloat64ToJsonString
	convertArrayObjectToJsonString    = convert.ConvertArrayObjectToJsonString
	convertArrayToString              = convert.ConvertArrayToString
	convertToJsonWithoutEscapeHTML    = convert.ConvertToJsonWithoutEscapeHTML
	convertToInterfaceArray           = convert.ConvertToInterfaceArray
	convertYamlToObject               = convert.ConvertYamlToObject
	Interface2String                  = convert.Interface2String
	Interface2StrSlice                = convert.Interface2StrSlice
	Str2InterfaceSlice                = convert.Str2InterfaceSlice
	Interface2Bool                    = convert.Interface2Bool
	expandSingletonToList             = convert.ExpandSingletonToList
	convertChargeTypeToPaymentType    = convert.ConvertChargeTypeToPaymentType
	convertPaymentTypeToChargeType    = convert.ConvertPaymentTypeToChargeType
	NormalizeMap                      = convert.NormalizeMap
	normalizeValue                    = convert.NormalizeValue

	// tag converters (tags/tags.go)
	ConvertTags             = tags.ConvertTags
	ConvertLowercaseTags    = tags.ConvertLowercaseTags
	ConvertTagsForKms       = tags.ConvertTagsForKms
	expandTagsToMap         = tags.ExpandTagsToMap
	expandTagsToMapWithTags = tags.ExpandTagsToMapWithTags

	// pointer helpers (helper/pointer.go)
	StringPointer = helper.StringPointer
	BoolPointer   = helper.BoolPointer
	Int32Pointer  = helper.Int32Pointer
	Int64Pointer  = helper.Int64Pointer

	// resource-id helpers (helper/resourceid.go)
	ParseResourceIds           = helper.ParseResourceIds
	ParseResourceId            = helper.ParseResourceId
	ParseResourceIdN           = helper.ParseResourceIdN
	ParseResourceIdWithEscaped = helper.ParseResourceIdWithEscaped
	EscapeColons               = helper.EscapeColons

	// compare helpers (helper/compare.go)
	compareJsonTemplateAreEquivalent      = helper.CompareJsonTemplateAreEquivalent
	compareJsonHasSubset                  = helper.CompareJsonHasSubset
	compareArrayJsonTemplateAreEquivalent = helper.CompareArrayJsonTemplateAreEquivalent
	compareYamlTemplateAreEquivalent      = helper.CompareYamlTemplateAreEquivalent
	compareMapWithIgnoreEquivalent        = helper.CompareMapWithIgnoreEquivalent
	IsSlicesValuesEqual                   = helper.IsSlicesValuesEqual

	// state-wait helpers (helper/state.go)
	BuildStateConf      = helper.BuildStateConf
	computePeriodByUnit = helper.ComputePeriodByUnit
	checkWaitForReady   = helper.CheckWaitForReady

	// format helpers (convert/format.go)
	formatInt     = convert.FormatInt
	formatBool    = convert.FormatBool
	formatFloat64 = convert.FormatFloat64

	// invoker (helper/invoker.go); methods resolve through the aliased type
	NewInvoker         = helper.NewInvoker
	ClientErrorCatcher = helper.ClientErrorCatcher
	ServiceBusyCatcher = helper.ServiceBusyCatcher
	ThrottlingCatcher  = helper.ThrottlingCatcher

	// openapi param helpers (helper/param.go)
	rpcParam     = helper.RpcParam
	roaParam     = helper.RoaParam
	xmlParam     = helper.XmlParam
	jsonXmlParam = helper.JsonXmlParam
	xmlJsonParam = helper.XmlJsonParam

	// misc utilities (helper/misc.go)
	IntMin                       = helper.IntMin
	getPagination                = helper.GetPagination
	userDataHashSum              = helper.UserDataHashSum
	Trim                         = helper.Trim
	ConvertIntegerToInt          = helper.ConvertIntegerToInt
	getNextpageNumber            = helper.GetNextpageNumber
	terraformToAPI               = helper.TerraformToAPI
	loadFileContent              = helper.LoadFileContent
	GetFunc                      = helper.GetFunc
	Case2Camel                   = helper.Case2Camel
	FirstLower                   = helper.FirstLower
	SplitSlice                   = helper.SplitSlice
	splitMultiZoneId             = helper.SplitMultiZoneId
	mapMerge                     = helper.MapMerge
	mapSort                      = helper.MapSort
	getOneStringOrAllStringSlice = helper.GetOneStringOrAllStringSlice
	Unique                       = helper.Unique
	IsSubCollection              = helper.IsSubCollection
	MergeMaps                    = helper.MergeMaps
	InArray                      = helper.InArray
	MD5                          = helper.MD5
	bytesToTB                    = helper.BytesToTB
	compressIPv6OrCIDR           = helper.CompressIPv6OrCIDR
	randIntRange                 = helper.RandIntRange
	IsEmpty                      = helper.IsEmpty
	IsNil                        = helper.IsNil
	GetDaysBetween2Date          = helper.GetDaysBetween2Date
	isProtocolValid              = helper.IsProtocolValid
	ValidProtocols               = helper.ValidProtocols
	isPagingRequest              = helper.IsPagingRequest
	setPagingRequest             = helper.SetPagingRequest
	getWriteOnlyValue            = helper.GetWriteOnlyValue
)

// Consts moved to helper and tags; re-exported as const aliases.
const (
	DefaultIntervalMini   = helper.DefaultIntervalMini
	DefaultIntervalShort  = helper.DefaultIntervalShort
	DefaultIntervalMedium = helper.DefaultIntervalMedium
	DefaultIntervalLong   = helper.DefaultIntervalLong

	PageNumSmall   = helper.PageNumSmall
	PageSizeSmall  = helper.PageSizeSmall
	PageSizeMedium = helper.PageSizeMedium
	PageSizeLarge  = helper.PageSizeLarge
	PageSizeXLarge = helper.PageSizeXLarge

	ServerSideEncryptionAes256 = helper.ServerSideEncryptionAes256
	ServerSideEncryptionKMS    = helper.ServerSideEncryptionKMS
	ServerSideEncryptionSM4    = helper.ServerSideEncryptionSM4

	Http  = helper.Http
	Https = helper.Https
	Tcp   = helper.Tcp
	Udp   = helper.Udp
	All   = helper.All
	Icmp  = helper.Icmp
	Gre   = helper.Gre

	// TagResourceType family (tags/tags.go)
	TagResourceImage         = tags.TagResourceImage
	TagResourceInstance      = tags.TagResourceInstance
	TagResourceAcl           = tags.TagResourceAcl
	TagResourceCertificate   = tags.TagResourceCertificate
	TagResourceDisk          = tags.TagResourceDisk
	TagResourceSecurityGroup = tags.TagResourceSecurityGroup
	TagResourceCdn           = tags.TagResourceCdn
	TagResourceApp           = tags.TagResourceApp
	TagResourceTopic         = tags.TagResourceTopic
	TagResourceCluster       = tags.TagResourceCluster
)

type PayType string

const (
	PrePaid    = PayType("PrePaid")
	PostPaid   = PayType("PostPaid")
	Prepaid    = PayType("Prepaid")
	Postpaid   = PayType("Postpaid")
	Serverless = PayType("Serverless")
)

const (
	NormalMode = "normal"
	SafetyMode = "safety"
)

type DdosbgpInstanceIpType string

const (
	IPv4 = DdosbgpInstanceIpType("IPv4")
	IPv6 = DdosbgpInstanceIpType("IPv6")
)

type NetType string

const (
	Internet = NetType("Internet")
	Intranet = NetType("Intranet")
)

type NetworkType string

const (
	Classic         = NetworkType("Classic")
	Vpc             = NetworkType("Vpc")
	ClassicInternet = NetworkType("classic_internet")
	ClassicIntranet = NetworkType("classic_intranet")
	PUBLIC          = NetworkType("PUBLIC")
	PRIVATE         = NetworkType("PRIVATE")
)

type NodeType string

const (
	WORKER = NodeType("WORKER")
	KIBANA = NodeType("KIBANA")
)

type ActionType string

const (
	OPEN  = ActionType("OPEN")
	CLOSE = ActionType("CLOSE")
)

type TimeType string

const (
	Hour  = TimeType("Hour")
	Day   = TimeType("Day")
	Week  = TimeType("Week")
	Month = TimeType("Month")
	Year  = TimeType("Year")
)

type IpVersion string

const (
	IPV4 = IpVersion("ipv4")
	IPV6 = IpVersion("ipv6")
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
	Normal      = Status("Normal")
	Changing    = Status("Changing")
	Online      = Status("online")
	Configuring = Status("configuring")

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

	Init            = Status("Init")
	Provisioning    = Status("Provisioning")
	Updating        = Status("Updating")
	FinancialLocked = Status("FinancialLocked")

	PUBLISHED   = Status("Published")
	NOPUBLISHED = Status("NonPublished")

	Deleted = Status("Deleted")
	Null    = Status("Null")

	Enable = Status("Enable")
	BINDED = Status("BINDED")
)

type IPType string

const (
	Inner   = IPType("Inner")
	Private = IPType("Private")
	Public  = IPType("Public")
)

type ResourceType string

const (
	ResourceTypeInstance      = ResourceType("Instance")
	ResourceTypeDisk          = ResourceType("Disk")
	ResourceTypeVSwitch       = ResourceType("VSwitch")
	ResourceTypeRds           = ResourceType("Rds")
	ResourceTypePolarDB       = ResourceType("PolarDB")
	IoOptimized               = ResourceType("IoOptimized")
	ResourceTypeRkv           = ResourceType("KVStore")
	ResourceTypeFC            = ResourceType("FunctionCompute")
	ResourceTypeElasticsearch = ResourceType("Elasticsearch")
	ResourceTypeSlb           = ResourceType("Slb")
	ResourceTypeMongoDB       = ResourceType("MongoDB")
	ResourceTypeGpdb          = ResourceType("Gpdb")
	ResourceTypeHBase         = ResourceType("HBase")
	ResourceTypeAdb           = ResourceType("ADB")
	ResourceTypeCassandra     = ResourceType("Cassandra")
)

type InternetChargeType string

const (
	PayByBandwidth = InternetChargeType("PayByBandwidth")
	PayByTraffic   = InternetChargeType("PayByTraffic")
	PayBy95        = InternetChargeType("PayBy95")
)

type AccountSite string

const (
	DomesticSite = AccountSite("Domestic")
	IntlSite     = AccountSite("International")
)
const (
	SnapshotCreatingInProcessing = Status("progressing")
	SnapshotCreatingAccomplished = Status("accomplished")
	SnapshotCreatingFailed       = Status("failed")

	SnapshotPolicyCreating  = Status("Creating")
	SnapshotPolicyAvailable = Status("available")
	SnapshotPolicyNormal    = Status("Normal")
)

// timeout for common product, ecs e.g.
const DefaultTimeout = 120
const Timeout5Minute = 300
const DefaultTimeoutMedium = 500

// timeout for long time progerss product, rds e.g.
const DefaultLongTimeout = 1000

const (
	// HeaderEnableEBTrigger header key for enabling eventbridge trigger
	// TODO: delete the header after eventbridge trigger is totally exposed to user
	HeaderEnableEBTrigger = "x-fc-enable-eventbridge-trigger"
)

// default region for all resource
const DEFAULT_REGION = "cn-beijing"

const INT_MAX = 2147483647

// symbol of multiIZ
const MULTI_IZ_SYMBOL = helper.MultiIzSymbol

const COMMA_SEPARATED = helper.CommaSeparated

const COLON_SEPARATED = ":"

const SLASH_SEPARATED = "/"

const LOCAL_HOST_IP = "127.0.0.1"

type OptimizedType string

const (
	IOOptimized   = OptimizedType("optimized")
	NoneOptimized = OptimizedType("none")
)

type KubernetesNodeType string

const (
	KubernetesNodeMaster = ResourceType("Master")
	KubernetesNodeWorker = ResourceType("Worker")
)

const CharityPageUrl = "http://promotion.alicdn.com/help/oss/error.html"

func ParseSlbListenerId(id string) (parts []string, err error) {
	parts = strings.Split(id, ":")
	if len(parts) != 2 && len(parts) != 3 {
		err = WrapError(fmt.Errorf("Invalid alicloud_slb_listener Id %s. Expected Id format is <slb id>:<protocol>:< frontend>.", id))
	}
	return parts, err
}

func GetCenChildInstanceType(id string) (c string, e error) {
	if strings.HasPrefix(id, "vpc") {
		return ChildInstanceTypeVpc, nil
	} else if strings.HasPrefix(id, "vbr") {
		return ChildInstanceTypeVbr, nil
	} else if strings.HasPrefix(id, "ccn") {
		return ChildInstanceTypeCcn, nil
	} else {
		return c, fmt.Errorf("CEN child instance ID invalid. Now, it only supports VPC or VBR or CCN instance.")
	}
}

// Migrated to alicloud/acctest.NewInstanceDiff.
// lintignore: R001
func newInstanceDiff(resourceName string, attributes, attributesDiff map[string]interface{}, state *terraform.InstanceState) (*terraform.InstanceDiff, error) {
	p := Provider().ResourcesMap
	dOld, _ := schema.InternalMap(p[resourceName].Schema).Data(state, nil)
	dNew, _ := schema.InternalMap(p[resourceName].Schema).Data(state, nil)
	for key, value := range attributes {
		err := dOld.Set(key, value)
		if err != nil {
			return nil, WrapErrorf(err, "[ERROR] the field %s setting error.", key)
		}
	}
	for key, value := range attributesDiff {
		attributes[key] = value
	}

	for key, value := range attributes {
		err := dNew.Set(key, value)
		if err != nil {
			return nil, WrapErrorf(err, "[ERROR] the field %s setting error.", key)
		}
	}

	diff := terraform.NewInstanceDiff()
	objectKey := ""
	for _, key := range mapSort(dNew.State().Attributes) {
		newValue := dNew.State().Attributes[key]
		if objectKey != "" && !strings.HasPrefix(key, objectKey) {
			objectKey = ""
		}
		if objectKey == "" {
			for _, suffix := range []string{"#", "%"} {
				if strings.HasSuffix(key, suffix) {
					objectKey = strings.TrimSuffix(key, suffix)
					break
				}
			}
		}
		oldValue, ok := dOld.State().Attributes[key]
		if ok && oldValue == newValue {
			continue
		}
		if oldValue == "" {
			for _, suffix := range []string{"#", "%"} {
				if strings.HasSuffix(key, suffix) {
					oldValue = "0"
				}
			}
		}
		diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: oldValue, New: newValue}

		if objectKey != "" {
			for removeKey, removeValue := range dOld.State().Attributes {
				if strings.HasPrefix(removeKey, objectKey) {
					if _, ok := dNew.State().Attributes[removeKey]; !ok {
						// If the attribue has complex elements, there should remove the key, not setting it to empty
						if len(strings.Split(removeKey, ".")) > 2 {
							delete(diff.Attributes, removeKey)
						} else {
							diff.Attributes[removeKey] = &terraform.ResourceAttrDiff{Old: removeValue, New: ""}
						}
					}
				}
			}
			objectKey = ""
		}
	}
	return diff, nil
}

func compareCmsHybridMonitorFcTaskYamlConfigAreEquivalent(tem1, tem2 string) (bool, error) {
	type MetricList struct {
		MetricList []string `yaml:"metric_list"`
	}
	type Product struct {
		MetricInfo []MetricList `yaml:"metric_info"`
		Namespace  string       `yaml:"namespace"`
	}
	type Products struct {
		Products []Product
	}

	var P1 Products
	err := yaml.Unmarshal([]byte(tem1), &P1)
	if err != nil {
		return false, err
	}

	y1 := make([]string, 0)
	for _, product := range P1.Products {
		s1, _ := json.Marshal(product)
		y1 = append(y1, string(s1))
	}

	var P2 Products
	err = yaml.Unmarshal([]byte(tem2), &P2)
	if err != nil {
		return false, err
	}

	y2 := make([]string, 0)
	for _, product := range P2.Products {
		s2, _ := json.Marshal(product)
		y2 = append(y2, string(s2))
	}

	sort.Strings(y1)
	sort.Strings(y2)
	return reflect.DeepEqual(y1, y2), nil
}

// IsContain reports whether item is present in items. Compatibility copy kept in
// the root package for unmigrated root consumers (e.g. cs_kubernetes_addon)
// after the PolarDB cluster file that originally defined it was migrated to
// the service/polardb product package. Behavior is unchanged.
func IsContain(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}
