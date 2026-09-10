package alicloud

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
