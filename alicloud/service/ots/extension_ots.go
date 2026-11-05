package ots

import (
	"fmt"
	"regexp"
)

// Status mirrors the root-package alicloud.Status type so the OTS wait helpers
// keep their original status comparisons after the split. It is intentionally a
// local type: helper.Status does not exist in alicloud/helper.
type Status string

const (
	Running        = Status("Running")
	Deleting       = Status("Deleting")
	DisabledStatus = Status("Disabled")
	Deleted        = Status("Deleted")
)

// Local constants that have no leaf-package counterpart in alicloud/helper.
// They are kept as plain integers to match the root-package definitions
// (WaitFor* helpers receive timeout as an int and multiply by time.Second).
const (
	INT_MAX            = 2147483647
	DefaultTimeout     = 120
	DefaultLongTimeout = 1000
	COLON_SEPARATED    = ":"
)

type PrimaryKeyTypeString string

const (
	IntegerType = PrimaryKeyTypeString("Integer")
	StringType  = PrimaryKeyTypeString("String")
	BinaryType  = PrimaryKeyTypeString("Binary")
)

type DefinedColumnTypeString string

const (
	DefinedColumnInteger = DefinedColumnTypeString("Integer")
	DefinedColumnString  = DefinedColumnTypeString("String")
	DefinedColumnBinary  = DefinedColumnTypeString("Binary")
	DefinedColumnDouble  = DefinedColumnTypeString("Double")
	DefinedColumnBoolean = DefinedColumnTypeString("Boolean")
)

type InstanceAccessedByType string

const (
	AnyNetwork   = InstanceAccessedByType("Any")
	VpcOnly      = InstanceAccessedByType("Vpc")
	VpcOrConsole = InstanceAccessedByType("ConsoleOrVpc")
)

type OtsInstanceType string

const (
	OtsCapacity        = OtsInstanceType("Capacity")
	OtsHighPerformance = OtsInstanceType("HighPerformance")
)

type OtsNetworkType string

const (
	VpcAccess      = OtsNetworkType("VPC")
	InternetAccess = OtsNetworkType("INTERNET")
	ClassicAccess  = OtsNetworkType("CLASSIC")
)

type OtsNetworkSource string

const (
	TrustProxyAccess = OtsNetworkSource("TRUST_PROXY")
)

func convertInstanceAccessedBy(accessed InstanceAccessedByType) string {
	switch accessed {
	case VpcOnly:
		return "VPC"
	case VpcOrConsole:
		return "VPC_CONSOLE"
	default:
		return "NORMAL"
	}
}

func convertInstanceAccessedByRevert(network string) InstanceAccessedByType {
	switch network {
	case "VPC":
		return VpcOnly
	case "VPC_CONSOLE":
		return VpcOrConsole
	default:
		return AnyNetwork
	}
}

func convertInstanceType(instanceType OtsInstanceType) string {
	switch instanceType {
	case OtsHighPerformance:
		return "SSD"
	default:
		return "HYBRID"
	}
}

func convertInstanceTypeRevert(instanceType string) OtsInstanceType {
	switch instanceType {
	case "SSD":
		return OtsHighPerformance
	default:
		return OtsCapacity
	}
}

func toInstanceOuterStatus(otsInstanceInnerStatus string) Status {
	switch otsInstanceInnerStatus {
	case "normal":
		return Running
	case "forbidden":
		return DisabledStatus
	case "deleting":
		return Deleting
	default:
		return Status(otsInstanceInnerStatus)
	}
}

func toInstanceInnerStatus(instanceOuterStatus Status) string {
	switch instanceOuterStatus {
	case Running:
		return "normal"
	case DisabledStatus:
		return "forbidden"
	case Deleting:
		return "deleting"
	default:
		return "INVALID"
	}
}

type TunnelTypeString string

const (
	BaseAndStreamTunnel = TunnelTypeString("BaseAndStream")
	BaseDataTunnel      = TunnelTypeString("BaseData")
	StreamTunnel        = TunnelTypeString("Stream")
)

type SseKeyTypeString string

const (
	SseKMSService = SseKeyTypeString("SSE_KMS_SERVICE")
	SseByOk       = SseKeyTypeString("SSE_BYOK")
)

type SecondaryIndexTypeString string

const (
	Local  = SecondaryIndexTypeString("Local")
	Global = SecondaryIndexTypeString("Global")
)
const (
	SearchIndexTypeHolder = "Search"
)

type OtsSearchIndexSyncPhaseString string

const (
	Full = OtsSearchIndexSyncPhaseString("Full")
	Incr = OtsSearchIndexSyncPhaseString("Incr")
)

type SearchIndexFieldTypeString string

const (
	OtsSearchTypeLong     = SearchIndexFieldTypeString("Long")
	OtsSearchTypeDouble   = SearchIndexFieldTypeString("Double")
	OtsSearchTypeBoolean  = SearchIndexFieldTypeString("Boolean")
	OtsSearchTypeKeyword  = SearchIndexFieldTypeString("Keyword")
	OtsSearchTypeText     = SearchIndexFieldTypeString("Text")
	OtsSearchTypeDate     = SearchIndexFieldTypeString("Date")
	OtsSearchTypeGeoPoint = SearchIndexFieldTypeString("GeoPoint")
	OtsSearchTypeNested   = SearchIndexFieldTypeString("Nested")
)

type SearchIndexAnalyzerTypeString string

const (
	OtsSearchSingleWord = SearchIndexAnalyzerTypeString("SingleWord")
	OtsSearchSplit      = SearchIndexAnalyzerTypeString("Split")
	OtsSearchMinWord    = SearchIndexAnalyzerTypeString("MinWord")
	OtsSearchMaxWord    = SearchIndexAnalyzerTypeString("MaxWord")
	OtsSearchFuzzy      = SearchIndexAnalyzerTypeString("Fuzzy")
)

type SearchIndexOrderTypeString string

const (
	OtsSearchSortOrderAsc  = SearchIndexOrderTypeString("Asc")
	OtsSearchSortOrderDesc = SearchIndexOrderTypeString("Desc")
)

type SearchIndexSortModeString string

const (
	OtsSearchModeMin = SearchIndexSortModeString("Min")
	OtsSearchModeMax = SearchIndexSortModeString("Max")
	OtsSearchModeAvg = SearchIndexSortModeString("Avg")
)

type SearchIndexSortFieldTypeString string

const (
	OtsSearchPrimaryKeySort = SearchIndexSortFieldTypeString("PrimaryKeySort")
	OtsSearchFieldSort      = SearchIndexSortFieldTypeString("FieldSort")
)

// otsRestTagsToMap mirrors the root-package helper so the OTS package can turn
// a REST instance's tag list into the resource-data map without importing the
// (not yet split) root alicloud package. Behavior is preserved verbatim.
func otsRestTagsToMap(tags []RestOtsTagInfo) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		result[t.Key] = t.Value
	}

	return result
}

type RestOtsInstanceInfo struct {
	InstanceStatus        string           `json:"InstanceStatus" xml:"InstanceStatus"`
	InstanceSpecification string           `json:"InstanceSpecification" xml:"InstanceSpecification"`
	Timestamp             string           `json:"Timestamp" xml:"Timestamp"`
	UserId                string           `json:"UserId" xml:"UserId"`
	ResourceGroupId       string           `json:"ResourceGroupId" xml:"ResourceGroupId"`
	InstanceName          string           `json:"InstanceName" xml:"InstanceName"`
	CreateTime            string           `json:"CreateTime" xml:"CreateTime"`
	Network               string           `json:"Network" xml:"Network"`
	NetworkTypeACL        []string         `json:"NetworkTypeACL" xml:"NetworkTypeACL"`
	NetworkSourceACL      []string         `json:"NetworkSourceACL" xml:"NetworkSourceACL"`
	Policy                string           `json:"Policy" xml:"Policy"`
	PolicyVersion         int              `json:"PolicyVersion" xml:"PolicyVersion"`
	InstanceDescription   string           `json:"InstanceDescription" xml:"InstanceDescription"`
	Quota                 RestOtsQuota     `json:"Quota" xml:"Quota"`
	Tags                  []RestOtsTagInfo `json:"Tags" xml:"Tags"`
}

type RestOtsQuota struct {
	TableQuota int `json:"TableQuota" xml:"TableQuota"`
}

type RestOtsTagInfo struct {
	Key   string `json:"Key" xml:"Key"`
	Value string `json:"Value" xml:"Value"`
}

// The instance name must be composed of a~z, A~Z, 0~9 and a hyphen (-),
// the first character must be a letter and the last character cannot be a hyphen (-),
// the legal length range is 3~16 bytes.
func validateOTSInstanceName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	reg := regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9-]{1,14}[a-zA-Z0-9]$")
	if !reg.MatchString(value) {
		errors = append(errors, fmt.Errorf("the instance name must consist of a~z, A~Z, 0~9 and a hyphen (-), "+
			"the first character must be a letter and the last character cannot be a hyphen (-), the legal length range is 3~16 bytes"))
	}
	return
}

// The table name must consist of a~z, A~Z, 0~9 and an underscore (_), the first character must be a letter or underscore (_),
// the legal length range is 1~255 bytes.
func validateOTSTableName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	reg := regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9_]{0,254}$")
	if !reg.MatchString(value) {
		errors = append(errors, fmt.Errorf("the table name must consist of a~z, A~Z, 0~9 and an underscore (_), "+
			"the first character must be a letter or underscore (_), the legal length range is 1~255 bytes"))
	}
	return
}

// The tunnel name must consist of a~z, A~Z, 0~9 and an underscore (_), the first character must be a letter or underscore (_),
// the legal length range is 1~255 bytes.
func validateOTSTunnelName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	reg := regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9_]{0,254}$")
	if !reg.MatchString(value) {
		errors = append(errors, fmt.Errorf("the tunnel name must consist of a~z, A~Z, 0~9 and an underscore (_), "+
			"the first character must be a letter or underscore (_), the legal length range is 1~255 bytes"))
	}
	return
}

// The index name must consist of a~z, A~Z, 0~9 and an underscore (_), the first character must be a letter or underscore (_),
// the legal length range is 1~255 bytes.
func validateOTSIndexName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	reg := regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9_]{0,254}$")
	if !reg.MatchString(value) {
		errors = append(errors, fmt.Errorf("the index name must consist of a~z, A~Z, 0~9 and an underscore (_), "+
			"the first character must be a letter or underscore (_), the legal length range is 1~255 bytes"))
	}
	return
}
