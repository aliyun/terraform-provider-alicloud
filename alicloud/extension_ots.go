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

// OTS instance total status: S_RUNNING = 1, S_DISABLED = 2, S_DELETING = 3
func convertOtsInstanceStatus(status Status) int {
	switch status {
	case Running:
		return 1
	case DisabledStatus:
		return 2
	case Deleting:
		return 3
	default:
		return -1
	}
}

func convertOtsInstanceStatusConvert(status int) Status {
	switch status {
	case 1:
		return Running
	case 2:
		return DisabledStatus
	case 3:
		return Deleting
	default:
		return ""
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
