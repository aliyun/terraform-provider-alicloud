package rds

// This file mirrors the root package's re-export aliases and local type/const
// definitions so the migrated SDKv2 RDS source files keep compiling with
// unchanged, unqualified references to symbols that previously lived in the
// flat `package alicloud` layout.
//
// WrapError/WrapErrorf/addDebug must stay var aliases rather than forwarding
// functions: they read runtime.Caller to attribute errors/debug to the caller,
// and a wrapper would insert a frame and misattribute every site. ComplexError,
// ProviderError, Catcher, Invoker, MyMap, Protocol, TagResourceType and Tag
// must be type aliases — the "=" is load-bearing for type assertions and
// method calls at the flat call sites.
//
// Package-level declarations that no migrated file references are harmless:
// Go does not report unused package-level vars/consts/types.

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
)

// --- errors.go aliases -------------------------------------------------------

type ComplexError = errs.ComplexError
type ProviderError = errs.ProviderError
type ErrorSource = errs.ErrorSource

var (
	Error                      = errs.Error
	WrapError                  = errs.WrapError
	WrapErrorf                 = errs.WrapErrorf
	WrapComplexError           = errs.WrapComplexError
	NotFoundError              = errs.NotFoundError
	NotFoundErr                = errs.NotFoundErr
	GetNotFoundErrorFromString = errs.GetNotFoundErrorFromString
	GetTimeErrorFromString     = errs.GetTimeErrorFromString
	NewProviderError           = errs.NewProviderError
	NeedRetry                  = helper.NeedRetry

	IsExpectedErrors       = errs.IsExpectedErrors
	NoCodeRegexRetry       = errs.NoCodeRegexRetry
	IsExpectedErrorCodes   = errs.IsExpectedErrorCodes
	GetNotFoundMessage     = errs.GetNotFoundMessage
	GetTimeoutMessage      = errs.GetTimeoutMessage
	GetCreateFailedMessage = errs.GetCreateFailedMessage
)

const (
	AlibabaCloudSdkGoERROR = errs.AlibabaCloudSdkGoERROR
	AliyunLogGoSdkERROR    = errs.AliyunLogGoSdkERROR
	AliyunDatahubSdkGo     = errs.AliyunDatahubSdkGo
	AliyunOssGoSdk         = errs.AliyunOssGoSdk
	FcGoSdk                = errs.FcGoSdk
	DenverdinoAliyungo     = errs.DenverdinoAliyungo
	AliyunTablestoreGoSdk  = errs.AliyunTablestoreGoSdk
	AliMnsERROR            = errs.AliMnsERROR
	ProviderERROR          = errs.ProviderERROR
	DefaultErrorMsg        = errs.DefaultErrorMsg
	DataDefaultErrorMsg    = errs.DataDefaultErrorMsg
	NotFoundWithResponse   = errs.NotFoundWithResponse
	FailedGetAttributeMsg  = errs.FailedGetAttributeMsg
	DefaultDebugMsg        = logs.DefaultDebugMsg
)

const (
	NotFound                = errs.NotFound
	ResourceNotfound        = errs.ResourceNotfound
	ServiceUnavailable      = errs.ServiceUnavailable
	InstanceNotFound        = errs.InstanceNotFound
	ForbiddenInstance       = errs.ForbiddenInstance
	MessageInstanceNotFound = errs.MessageInstanceNotFound
	Throttling              = errs.Throttling
	ThrottlingUser          = errs.ThrottlingUser

	RamInstanceNotFound   = errs.RamInstanceNotFound
	AliyunGoClientFailure = errs.AliyunGoClientFailure

	LogClientTimeout = errs.LogClientTimeout

	InvalidFileSystemStatus_Ordering = errs.InvalidFileSystemStatus_Ordering
	NotFoundArticle                  = errs.NotFoundArticle
)

var (
	SlbIsBusy                         = errs.SlbIsBusy
	EcsNotFound                       = errs.EcsNotFound
	DiskInvalidOperation              = errs.DiskInvalidOperation
	NetworkInterfaceInvalidOperations = errs.NetworkInterfaceInvalidOperations
	OperationDeniedDBStatus           = errs.OperationDeniedDBStatus
	DBReadInstanceNotReadyStatus      = errs.DBReadInstanceNotReadyStatus
	NasNotFound                       = errs.NasNotFound
	SnapshotInvalidOperations         = errs.SnapshotInvalidOperations
	DiskNotSupportOnlineChangeErrors  = errs.DiskNotSupportOnlineChangeErrors
	DBInstanceTDEErrors               = errs.DBInstanceTDEErrors

	OtsTableIsTemporarilyUnavailable          = errs.OtsTableIsTemporarilyUnavailable
	OtsTunnelIsTemporarilyUnavailable         = errs.OtsTunnelIsTemporarilyUnavailable
	OtsSecondaryIndexIsTemporarilyUnavailable = errs.OtsSecondaryIndexIsTemporarilyUnavailable
	OtsSearchIndexIsTemporarilyUnavailable    = errs.OtsSearchIndexIsTemporarilyUnavailable
)

const ResponseCodeMsg = errs.ResponseCodeMsg
const RequestIdMsg = errs.RequestIdMsg
const NotFoundMsg = errs.NotFoundMsg
const NotFoundWithError = errs.NotFoundWithError
const DefaultTimeoutMsg = errs.DefaultTimeoutMsg
const DeleteTimeoutMsg = errs.DeleteTimeoutMsg
const WaitTimeoutMsg = errs.WaitTimeoutMsg
const SweepDefaultErrorMsg = errs.SweepDefaultErrorMsg
const IdMsg = errs.IdMsg
const FailedToReachTargetStatus = errs.FailedToReachTargetStatus
const FailedToReachTargetStatusWithResponse = errs.FailedToReachTargetStatusWithResponse
const FailedToReachTargetStatusWithError = errs.FailedToReachTargetStatusWithError
const FailedToReachTargetStatusWithRequestId = errs.FailedToReachTargetStatusWithRequestId
const FailedToReachTargetAttribute = errs.FailedToReachTargetAttribute
const RequiredWhenMsg = errs.RequiredWhenMsg
const InvalidAttributeValue = errs.InvalidAttributeValue

// --- common.go aliases -------------------------------------------------------

type Invoker = helper.Invoker
type Catcher = helper.Catcher
type MyMap = helper.MyMap
type Protocol = helper.Protocol
type TagResourceType = tags.TagResourceType

var (
	addDebug           = logs.AddDebug
	debugOn            = logs.DebugOn
	buildClientToken   = helper.BuildClientToken
	incrementalWait    = helper.IncrementalWait
	writeToFile        = helper.WriteToFile
	GetUserHomeDir     = helper.GetUserHomeDir
	dataResourceIdHash = helper.DataResourceIdHash

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

	ConvertTags             = tags.ConvertTags
	ConvertLowercaseTags    = tags.ConvertLowercaseTags
	ConvertTagsForKms       = tags.ConvertTagsForKms
	expandTagsToMap         = tags.ExpandTagsToMap
	expandTagsToMapWithTags = tags.ExpandTagsToMapWithTags
	tagsSchema              = tags.TagsSchema
	ignoredTags             = tags.IgnoredTags
	parsingTags             = tags.ParsingTags
	tagsToMap               = tags.TagsToMap

	StringPointer = helper.StringPointer
	BoolPointer   = helper.BoolPointer
	Int32Pointer  = helper.Int32Pointer
	Int64Pointer  = helper.Int64Pointer

	ParseResourceIds           = helper.ParseResourceIds
	ParseResourceId            = helper.ParseResourceId
	ParseResourceIdN           = helper.ParseResourceIdN
	ParseResourceIdWithEscaped = helper.ParseResourceIdWithEscaped
	EscapeColons               = helper.EscapeColons

	compareJsonTemplateAreEquivalent      = helper.CompareJsonTemplateAreEquivalent
	compareJsonHasSubset                  = helper.CompareJsonHasSubset
	compareArrayJsonTemplateAreEquivalent = helper.CompareArrayJsonTemplateAreEquivalent
	compareYamlTemplateAreEquivalent      = helper.CompareYamlTemplateAreEquivalent
	compareMapWithIgnoreEquivalent        = helper.CompareMapWithIgnoreEquivalent
	IsSlicesValuesEqual                   = helper.IsSlicesValuesEqual

	BuildStateConf      = helper.BuildStateConf
	computePeriodByUnit = helper.ComputePeriodByUnit
	checkWaitForReady   = helper.CheckWaitForReady

	formatInt     = convert.FormatInt
	formatBool    = convert.FormatBool
	formatFloat64 = convert.FormatFloat64

	NewInvoker         = helper.NewInvoker
	ClientErrorCatcher = helper.ClientErrorCatcher
	ServiceBusyCatcher = helper.ServiceBusyCatcher
	ThrottlingCatcher  = helper.ThrottlingCatcher

	rpcParam     = helper.RpcParam
	roaParam     = helper.RoaParam
	xmlParam     = helper.XmlParam
	jsonXmlParam = helper.JsonXmlParam
	xmlJsonParam = helper.XmlJsonParam

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

type DdosbgpInsatnceType string

const (
	Enterprise   = DdosbgpInsatnceType("Enterprise")
	Professional = DdosbgpInsatnceType("Professional")
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

const DefaultTimeout = 120
const Timeout5Minute = 300
const DefaultTimeoutMedium = 500
const DefaultLongTimeout = 1000

const (
	HeaderEnableEBTrigger = "x-fc-enable-eventbridge-trigger"
)

const DEFAULT_REGION = "cn-beijing"

const INT_MAX = 2147483647

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

// --- validators.go aliases --------------------------------------------------

var (
	validateCIDRNetworkAddress       = validators.ValidateCIDRNetworkAddress
	validateVpnCIDRNetworkAddress    = validators.ValidateVpnCIDRNetworkAddress
	validateSwitchCIDRNetworkAddress = validators.ValidateSwitchCIDRNetworkAddress
	validateAllowedSplitStringValue  = validators.ValidateAllowedSplitStringValue
	validateStringConvertInt64       = validators.ValidateStringConvertInt64
	normalizeYamlString              = validators.NormalizeYamlString
	normalizeJsonString              = validators.NormalizeJsonString
	validateYamlString               = validators.ValidateYamlString
	validateNormalName               = validators.ValidateNormalName
	skipResourceSchemaValidation     = validators.SkipResourceSchemaValidation

	IntBetween                 = validators.IntBetween
	ValidateNullableIntBetween = validators.ValidateNullableIntBetween
	ValidateNullableBool       = validators.ValidateNullableBool
	IntAtLeast                 = validators.IntAtLeast
	IntAtMost                  = validators.IntAtMost
	IntInSlice                 = validators.IntInSlice
	StringInSlice              = validators.StringInSlice
	StringLenBetween           = validators.StringLenBetween
	StringMatch                = validators.StringMatch
	StringDoesNotMatch         = validators.StringDoesNotMatch
	StringDoesNotContainAny    = validators.StringDoesNotContainAny
	FloatBetween               = validators.FloatBetween
	FloatAtLeast               = validators.FloatAtLeast
	FloatAtMost                = validators.FloatAtMost
	ValidateRFC3339TimeString  = validators.ValidateRFC3339TimeString
	StringLenAtLeast           = validators.StringLenAtLeast

	UpperLowerCaseDiffSuppressFunc = validators.UpperLowerCaseDiffSuppressFunc
)

// --- extension_tags.go aliases ----------------------------------------------

type Tag = tags.Tag
