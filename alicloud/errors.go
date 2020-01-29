package alicloud

import (
	"strings"

	sls "github.com/aliyun/aliyun-log-go-sdk"

	"fmt"

	"log"
	"runtime"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/fc-go-sdk"
	"github.com/denverdino/aliyungo/common"
)

const (
	// common
	NotFound           = "NotFound"
	WaitForTimeout     = "WaitForTimeout"
	ResourceNotFound   = "ResourceNotfound"
	ServiceUnavailable = "ServiceUnavailable"
	// ecs
	InstanceNotFound        = "Instance.Notfound"
	MessageInstanceNotFound = "instance is not found"
	ImageIdNotFound         = "InvalidImageId.NotFound"
	// disk
	InternalError       = "InternalError"
	DependencyViolation = "DependencyViolation"
	// eip
	EipIncorrectStatus         = "IncorrectEipStatus"
	InstanceIncorrectStatus    = "IncorrectInstanceStatus"
	HaVipIncorrectStatus       = "IncorrectHaVipStatus"
	COMMODITYINVALID_COMPONENT = "COMMODITY.INVALID_COMPONENT"
	AllocationIdNotFound       = "InvalidAllocationId.NotFound"
	HasBeenUsedBySnatTable     = "InvalidIpStatus.HasBeenUsedBySnatTable"
	HasBeenUsedByForwardEntry  = "InvalidIpStatus.HasBeenUsedByForwardEntry"
	// slb
	LoadBalancerNotFound                  = "InvalidLoadBalancerId.NotFound"
	ListenerNotFound                      = "The specified resource does not exist"
	SlbOrderFailed                        = "OrderFailed"
	VServerGroupNotFoundMessage           = "The specified VServerGroupId does not exist"
	MasterSlaveServerGroupNotFoundMessage = "The specified MasterSlaveGroupId does not exist"
	RspoolVipExist                        = "RspoolVipExist"
	InvalidParameter                      = "InvalidParameter"
	InvalidRuleIdNotFound                 = "InvalidRuleId.NotFound"
	BackendServerConfiguring              = "BackendServer.configuring"
	InvalidDomainExtensionIdNotFound      = "InvalidParameter.DomainExtensionId"
	DomainExtensionProcessing             = "DomainExtensionProcessing"
	// slb acl
	SlbAclNotExists = "AclNotExist"

	// security_group
	InvalidSecurityGroupIdNotFound = "InvalidSecurityGroupId.NotFound"
	SgDependencyViolation          = "DependencyViolation"

	//Nat gateway
	InvalidNatGatewayIdNotFound = "InvalidNatGatewayId.NotFound"
	// vpc
	InvalidVpcIDNotFound = "InvalidVpcID.NotFound"
	ForbiddenVpcNotFound = "Forbidden.VpcNotFound"
	Throttling           = "Throttling"
	IncorrectVpcStatus   = "IncorrectVpcStatus"
	IncorrectStatus      = "IncorrectStatus"

	// NAS
	InvalidFileSystemIDNotFound = "InvalidFileSystem.NotFound"
	InvalidAccessGroupNotFound  = "InvalidAccessGroup.NotFound"
	ForbiddenNasNotFound        = "Forbidden.NasNotFound"

	//apigatway
	ApiGroupNotFound      = "NotFoundApiGroup"
	RepeatedCommit        = "RepeatedCommit"
	ApiNotFound           = "NotFoundApi"
	NotFoundApp           = "NotFoundApp"
	NotFoundAuthorization = "NotFoundAuthorization"
	NotFoundStage         = "NotFoundStage"

	// vswitch
	VswitcInvalidRegionId    = "InvalidRegionId.NotFound"
	InvalidVswitchIDNotFound = "InvalidVswitchID.NotFound"
	TokenProcessing          = "OperationFailed.IdempotentTokenProcessing"
	//route entry
	IncorrectRouteEntryStatus            = "IncorrectRouteEntryStatus"
	InvalidStatusRouteEntry              = "InvalidStatus.RouteEntry"
	TaskConflict                         = "TaskConflict"
	RouterEntryForbbiden                 = "Forbbiden"
	RouterEntryConflictDuplicated        = "RouterEntryConflict.Duplicated"
	InvalidCidrBlockOverlapped           = "InvalidCidrBlock.Overlapped"
	IncorrectOppositeInterfaceInfoNotSet = "IncorrectOppositeInterfaceInfo.NotSet"
	InvalidSnatTableIdNotFound           = "InvalidSnatTableId.NotFound"
	InvalidSnatEntryIdNotFound           = "InvalidSnatEntryId.NotFound"
	IncorretSnatEntryStatus              = "IncorretSnatEntryStatus"
	InvalidRouteEntryNotFound            = "InvalidRouteEntry.NotFound"
	// Forward
	InvalidIpNotInNatgw           = "InvalidIp.NotInNatgw"
	InvalidForwardTableIdNotFound = "InvalidForwardTableId.NotFound"
	InvalidForwardEntryIdNotFound = "InvalidForwardEntryId.NotFound"

	// ess
	InvalidScalingGroupIdNotFound               = "InvalidScalingGroupId.NotFound"
	InvalidScheduledTaskIdNotFound              = "InvalidScheduledTaskId.NotFound"
	IncorrectScalingConfigurationLifecycleState = "IncorrectScalingConfigurationLifecycleState"
	IncorrectScalingGroupStatus                 = "IncorrectScalingGroupStatus"
	IncorrectCapacityMaxSize                    = "IncorrectCapacity.MaxSize"
	IncorrectCapacityMinSize                    = "IncorrectCapacity.MinSize"
	ScalingActivityInProgress                   = "ScalingActivityInProgress"
	InvalidScalingRuleIdNotFound                = "InvalidScalingRuleId.NotFound"
	InvalidLifecycleHookIdNotFound              = "InvalidLifecycleHookId.NotExist"
	InvalidNotificationNotFound                 = "NotificationConfigurationNotExist"
	InvalidEssAlarmTaskNotFound                 = "404"
	InvalidEssScalingConfigurationIdNotFound    = "InvalidScalingConfigurationId.NotFound"
	//drds
	InvalidDRDSInstanceIdNotFound = "InvalidDrdsInstanceId.NotFound"

	//mongodb
	InvalidMongoDBInstanceIdNotFound = "InvalidDBInstanceId.NotFound"
	InvalidMongoDBNameNotFound       = "InvalidDBName.NotFound"

	//hbase
	InvalidHBaseInstanceIdNotFound = "Instance.NotFound"
	InvalidHBaseNameNotFound       = "InvalidDBName.NotFound"

	// gpdb
	InvalidGpdbNameNotFound       = "InvalidDBName.NotFound"
	InvalidGpdbInstanceIdNotFound = "InvalidDBInstanceId.NotFound"
	InvalidGpdbInstanceStatus     = "OperationDenied.DBInstanceStatus"
	InvalidGpdbConcurrentOperate  = "SYSTEM.CONCURRENT_OPERATE"

	// rds
	InvalidDBInstanceIdNotFound            = "InvalidDBInstanceId.NotFound"
	InvalidDBNameNotFound                  = "InvalidDBName.NotFound"
	InvalidDBInstanceNameNotFound          = "InvalidDBInstanceName.NotFound"
	InvalidCurrentConnectionStringNotFound = "InvalidCurrentConnectionString.NotFound"
	InvalidRwSplitNetTypeNotFound          = "InvalidRwSplitNetType.NotFound"
	InvalidAccountNameNotFound             = "InvalidAccountName.NotFound"
	AtLeastOneNetTypeExists                = "AtLeastOneNetTypeExists"
	DBOperationDeniedOutofUsage            = "OperationDenied.OutofUsage"

	// polarDB
	InvalidDBClusterIdNotFound      = "InvalidDBClusterId.NotFound"
	InvalidDBClusterNameNotFound    = "InvalidDBClusterName.NotFound"
	InvalidDBClusterStatus          = "OperationDenied.DBClusterStatus"
	EndpointStatusNotSupport        = "EndpointStatus.NotSupport"
	InvalidPolarDBClusterStatus     = "OperationDenied.PolarDBClusterStatus"
	InvalidReadPolarDBClusterStatus = "OperationDenied.ReadPolarDBClusterStatus"

	// oss
	OssBucketNotFound                 = "NoSuchBucket"
	OssBodyNotFound                   = "404 Not Found"
	NoSuchCORSConfiguration           = "NoSuchCORSConfiguration"
	NoSuchWebsiteConfiguration        = "NoSuchWebsiteConfiguration"
	InsufficientBucketPolicyException = "InsufficientBucketPolicyException"

	// RAM Instance Not Found
	RamInstanceNotFound   = "Forbidden.InstanceNotFound"
	AliyunGoClientFailure = "AliyunGoClientFailure"

	// dns
	RecordForbiddenDNSChange    = "RecordForbidden.DNSChange"
	FobiddenNotEmptyGroup       = "Fobidden.NotEmptyGroup"
	DomainRecordNotBelongToUser = "DomainRecordNotBelongToUser"
	InvalidDomainNotFound       = "InvalidDomain.NotFound"
	InvalidDomainNameNoExist    = "InvalidDomainName.NoExist"

	// ram user
	DeleteConflictUserGroup        = "DeleteConflict.User.Group"
	DeleteConflictUserAccessKey    = "DeleteConflict.User.AccessKey"
	DeleteConflictUserLoginProfile = "DeleteConflict.User.LoginProfile"
	DeleteConflictUserMFADevice    = "DeleteConflict.User.MFADevice"
	DeleteConflictUserPolicy       = "DeleteConflict.User.Policy"

	// ram group
	DeleteConflictGroupUser   = "DeleteConflict.Group.User"
	DeleteConflictGroupPolicy = "DeleteConflict.Group.Policy"

	// ram role
	DeleteConflictRolePolicy = "DeleteConflict.Role.Policy"
	EntityNotExistRole       = "EntityNotExist.Role"

	// ram policy
	DeleteConflictPolicyUser    = "DeleteConflict.Policy.User"
	DeleteConflictPolicyGroup   = "DeleteConflict.Policy.Group"
	DeleteConflictPolicyVersion = "DeleteConflict.Policy.Version"

	//unknown Error
	UnknownError = "UnknownError"

	// Keypair error
	KeyPairNotFound           = "InvalidKeyPair.NotFound"
	KeyPairServiceUnavailable = "ServiceUnavailable"

	// Container
	ErrorClusterNotFound = "ErrorClusterNotFound"

	// cr
	ErrorNamespaceNotExist = "NAMESPACE_NOT_EXIST"
	ErrorRepoNotExist      = "REPO_NOT_EXIST"

	// cdn
	ServiceBusy = "ServiceBusy"

	// RAM
	InvalidRamRoleNotFound       = "InvalidRamRole.NotFound"
	RoleAttachmentUnExpectedJson = "unexpected end of JSON input"
	InvalidInstanceIdNotFound    = "InvalidInstanceId.NotFound"

	RouterInterfaceIncorrectStatus                        = "IncorrectStatus"
	DependencyViolationRouterInterfaceReferedByRouteEntry = "DependencyViolation.RouterInterfaceReferedByRouteEntry"

	// CS
	ErrorClusterNameAlreadyExist = "ErrorClusterNameAlreadyExist"
	ApplicationNotFound          = "Not Found"
	ApplicationErrorIgnore       = "Unable to reach primary cluster manager"
	ApplicationConfirmConflict   = "Conflicts with unconfirmed updates for operation"

	// privatezone
	ZoneNotExists    = "Zone.NotExists"
	ZoneVpcNotExists = "ZoneVpc.NotExists.VpcId"
	ThrottlingUser   = "Throttling.User"
	PvtzSystemBusy   = "System.Busy"

	// log
	ProjectNotExist     = "ProjectNotExist"
	IndexConfigNotExist = "IndexConfigNotExist"
	LogStoreNotExist    = "LogStoreNotExist"
	InternalServerError = "InternalServerError"
	LogClientTimeout    = "Client.Timeout exceeded while awaiting headers"
	LogConfigNotExist   = "ConfigNotExist"
	// OTS
	OTSObjectNotExist        = "OTSObjectNotExist"
	SuffixNoSuchHost         = "no such host"
	OTSQuotaFrequentMsg      = "Too frequent table operations."
	OTSInternalServerError   = "OTSInternalServerError"
	OTSServerBusy            = "OTSServerBusy"
	OTSPartitionUnavailable  = "OTSPartitionUnavailable"
	OTSTimeout               = "OTSTimeout"
	OTSServerUnavailable     = "OTSServerUnavailable"
	OTSRowOperationConflict  = "OTSRowOperationConflict"
	OTSTableNotReady         = "OTSTableNotReady"
	OTSNotEnoughCapacityUnit = "OTSNotEnoughCapacityUnit"

	// FC
	ServiceNotFound  = "ServiceNotFound"
	FunctionNotFound = "FunctionNotFound"
	TriggerNotFound  = "TriggerNotFound"
	AccessDenied     = "AccessDenied"

	// Vpn
	VpnNotFound              = "InvalidVpnGatewayInstanceId.NotFound"
	VpnForbidden             = "Forbidden"
	CgwNotFound              = "InvalidCustomerGatewayInstanceId.NotFound"
	VpnConnNotFound          = "InvalidVpnConnectionInstanceId.NotFound"
	SslVpnServerNotFound     = "InvalidSslVpnServerId.NotFound"
	SslVpnClientCertNotFound = "InvalidSslVpnClientCertId.NotFound"
	VpnConfiguring           = "VpnGateway.Configuring"
	// CEN
	OperationBlocking              = "Operation.Blocking"
	ParameterCenInstanceIdNotExist = "ParameterCenInstanceId"
	InvalidCenInstanceStatus       = "InvalidOperation.CenInstanceStatus"
	ParameterInstanceIdNotExist    = "ParameterInstanceId"
	ParameterBwpInstanceId         = "ParameterBwpInstanceId"
	InvalidBwpInstanceStatus       = "InvalidOperation.BwpInstanceStatus"
	InvalidBwpBusinessStatus       = "InvalidOperation.BwpBusinessStatus"

	// snapshot
	SnapshotNotFound = "InvalidSnapshotId.NotFound"

	// kv-store
	InvalidKVStoreInstanceIdNotFound = "InvalidInstanceId.NotFound"
	// MNS
	QueueNotExist        = "QueueNotExist"
	TopicNotExist        = "TopicNotExist"
	SubscriptionNotExist = "SubscriptionNotExist"

	// ONS
	OnsInstanceNotExist    = "INSTANCE_NOT_FOUND"
	OnsInstanceNotEmpty    = "INSTANCE_NOT_EMPTY"
	AuthResourceOwnerError = "AUTH_RESOURCE_OWNER_ERROR"

	// ALIKAFKA
	AlikafkaFlowControl          = "ONS_SYSTEM_FLOW_CONTROL"
	AlikafkaTopicNotFound        = "BIZ_TOPIC_NOT_FOUND"
	AlikafkaSubscriptionNotFound = "BIZ_SUBSCRIPTION_NOT_FOUND"

	//HaVip
	InvalidHaVipIdNotFound = "InvalidHaVipId.NotFound"
	InvalidVipStatus       = "InvalidVip.Status"
	IncorrectHaVipStatus   = "IncorrectHaVipStatus"

	//Cas
	CertNotExist = "CertNotExist"

	InvalidPrivateIpAddressDuplicated = "InvalidPrivateIpAddress.Duplicated"

	// Elasticsearch
	InstanceActivating                 = "InstanceActivating"
	ESInstanceNotFound                 = "InstanceNotFound"
	ESMustChangeOneResource            = "MustChangeOneResource"
	ESCssCheckUpdowngradeError         = "CssCheckUpdowngradeError"
	ESConcurrencyConflictError         = "ConcurrencyUpdateInstanceConflict"
	ESNotSupportCurrentActionError     = "InstanceStatusNotSupportCurrentAction"
	ESTokenPreviousRequestProcessError = "TokenPreviousRequestProcessError"

	// Ddoscoo
	DdoscooInstanceNotFound = "InstanceNotFound"
	InvalidDdoscooInstance  = "ddos_coop3301"
	InvalidDdosbgpInstance  = "InvalidInstance"

	// Ddosbgp
	DdosbgpInstanceNotFound = "InstanceNotFound"

	//nacl
	NetworkAclNotFound = "InvalidNetworkAcl.NotFound"

	//Actiontrail
	InvalidTrailNotFound  = "TrailNotFoundException"
	TrailNeedRamAuthorize = "NeedRamAuthorize"

	//emr
	ClusterNotFound = "ClusterId.NotFound"
)

var SlbIsBusy = []string{"SystemBusy", "OperationBusy", "ServiceIsStopping", "BackendServer.configuring", "ServiceIsConfiguring"}
var EcsNotFound = []string{"InvalidInstanceId.NotFound", "Forbidden.InstanceNotFound"}
var DiskInvalidOperation = []string{"IncorrectDiskStatus", "IncorrectInstanceStatus", "OperationConflict", InternalError, "InvalidOperation.Conflict", "IncorrectDiskStatus.Initializing"}
var NetworkInterfaceInvalidOperations = []string{"InvalidOperation.InvalidEniState", "InvalidOperation.InvalidEcsState", "OperationConflict", "ServiceUnavailable", "InternalError"}
var OperationDeniedDBStatus = []string{"OperationDenied.DBStatus", "OperationDenied.DBInstanceStatus", "OperationDenied.DBClusterStatus", InternalError, DBOperationDeniedOutofUsage}
var DBReadInstanceNotReadyStatus = []string{"OperationDenied.ReadDBInstanceStatus", "OperationDenied.MasterDBInstanceState", "ReadDBInstance.Mismatch"}
var NasNotFound = []string{"InvalidMountTarget.NotFound", InvalidFileSystemIDNotFound, ForbiddenNasNotFound, "InvalidLBid.NotFound", "VolumeUnavailable"}
var SnapshotInvalidOperations = []string{"OperationConflict", "ServiceUnavailable", "InternalError", "SnapshotCreatedDisk", "SnapshotCreatedImage"}
var SnapshotPolicyInvalidOperations = []string{"OperationConflict", "ServiceUnavailable", "InternalError", "SnapshotCreatedDisk", "SnapshotCreatedImage"}
var DiskNotSupportOnlineChangeErrors = []string{"InvalidDiskCategory.NotSupported", "InvalidRegion.NotSupport", "IncorrectInstanceStatus", "IncorrectDiskStatus", "InvalidOperation.InstanceTypeNotSupport"}

// details at: https://help.aliyun.com/document_detail/27300.html
var OtsTableIsTemporarilyUnavailable = []string{SuffixNoSuchHost, OTSServerBusy, OTSPartitionUnavailable, OTSInternalServerError,
	OTSTimeout, OTSServerUnavailable, OTSRowOperationConflict, OTSTableNotReady, OTSNotEnoughCapacityUnit, OTSQuotaFrequentMsg}

// An Error represents a custom error for Terraform failure response
type ProviderError struct {
	errorCode string
	message   string
}

func (e *ProviderError) Error() string {
	return fmt.Sprintf("[ERROR] Terraform Alicloud Provider Error: Code: %s Message: %s", e.errorCode, e.message)
}

func (err *ProviderError) ErrorCode() string {
	return err.errorCode
}

func (err *ProviderError) Message() string {
	return err.message
}

func GetNotFoundErrorFromString(str string) error {
	return &ProviderError{
		errorCode: InstanceNotFound,
		message:   str,
	}
}
func NotFoundError(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*ComplexError); ok {
		if e.Err != nil && strings.HasPrefix(e.Err.Error(), ResourceNotFound) {
			return true
		}
		return NotFoundError(e.Cause)
	}
	if err == nil {
		return false
	}

	if e, ok := err.(*errors.ServerError); ok {
		return e.ErrorCode() == InstanceNotFound || e.ErrorCode() == RamInstanceNotFound || e.ErrorCode() == NotFound || strings.Contains(strings.ToLower(e.Message()), MessageInstanceNotFound)
	}

	if e, ok := err.(*ProviderError); ok {
		return e.ErrorCode() == InstanceNotFound || e.ErrorCode() == RamInstanceNotFound || e.ErrorCode() == NotFound || strings.Contains(strings.ToLower(e.Message()), MessageInstanceNotFound)
	}

	if e, ok := err.(*common.Error); ok {
		return e.Code == InstanceNotFound || e.Code == RamInstanceNotFound || e.Code == NotFound || strings.Contains(strings.ToLower(e.Message), MessageInstanceNotFound)
	}

	return false
}

func IsExpectedErrors(err error, expectCodes []string) bool {
	if err == nil {
		return false
	}

	if e, ok := err.(*ComplexError); ok {
		return IsExpectedErrors(e.Cause, expectCodes)
	}
	if err == nil {
		return false
	}

	if e, ok := err.(*errors.ServerError); ok {
		for _, code := range expectCodes {
			if e.ErrorCode() == code || strings.Contains(e.Message(), code) {
				return true
			}
		}
		return false
	}

	if e, ok := err.(*ProviderError); ok {
		for _, code := range expectCodes {
			if e.ErrorCode() == code || strings.Contains(e.Message(), code) {
				return true
			}
		}
		return false
	}

	if e, ok := err.(*common.Error); ok {
		for _, code := range expectCodes {
			if e.Code == code || strings.Contains(e.Message, code) {
				return true
			}
		}
		return false
	}

	if e, ok := err.(*sls.Error); ok {
		for _, code := range expectCodes {
			if e.Code == code || strings.Contains(e.Message, code) || strings.Contains(e.String(), code) {
				return true
			}
		}
		return false
	}

	if e, ok := err.(oss.ServiceError); ok {
		for _, code := range expectCodes {
			if e.Code == code || strings.Contains(e.Message, code) {
				return true
			}
		}
		return false
	}

	if e, ok := err.(*fc.ServiceError); ok {
		for _, code := range expectCodes {
			if e.ErrorCode == code || strings.Contains(e.ErrorMessage, code) {
				return true
			}
		}
		return false
	}

	if e, ok := err.(datahub.DatahubError); ok {
		for _, code := range expectCodes {
			if e.Code == code || strings.Contains(e.Message, code) {
				return true
			}
		}
		return false
	}

	for _, code := range expectCodes {
		if strings.Contains(err.Error(), code) {
			return true
		}
	}
	return false
}

func RamEntityNotExist(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*ComplexError); ok {
		err = e.Cause
	}
	if err == nil {
		return false
	}
	if e, ok := err.(*errors.ServerError); ok && strings.Contains(e.ErrorCode(), "EntityNotExist") {
		return true
	}
	return false
}

func GetTimeErrorFromString(str string) error {
	return &ProviderError{
		errorCode: WaitForTimeout,
		message:   str,
	}
}

func GetNotFoundMessage(product, id string) string {
	return fmt.Sprintf("The specified %s %s is not found.", product, id)
}

func GetTimeoutMessage(product, status string) string {
	return fmt.Sprintf("Waitting for %s %s is timeout.", product, status)
}

type ErrorSource string

const (
	AlibabaCloudSdkGoERROR = ErrorSource("[SDK alibaba-cloud-sdk-go ERROR]")
	AliyunLogGoSdkERROR    = ErrorSource("[SDK aliyun-log-go-sdk ERROR]")
	AliyunDatahubSdkGo     = ErrorSource("[SDK aliyun-datahub-sdk-go ERROR]")
	AliyunOssGoSdk         = ErrorSource("[SDK aliyun-oss-go-sdk ERROR]")
	FcGoSdk                = ErrorSource("[SDK fc-go-sdk ERROR]")
	DenverdinoAliyungo     = ErrorSource("[SDK denverdino/aliyungo ERROR]")
	AliyunTablestoreGoSdk  = ErrorSource("[SDK aliyun-tablestore-go-sdk ERROR]")
	AliMnsERROR            = ErrorSource("[SDK ali_mns ERROR]")
	ProviderERROR          = ErrorSource("[Provider ERROR]")
)

// ComplexError is a format error which including origin error, extra error message, error occurred file and line
// Cause: a error is a origin error that comes from SDK, some exceptions and so on
// Err: a new error is built from extra message
// Path: the file path of error occurred
// Line: the file line of error occurred
type ComplexError struct {
	Cause error
	Err   error
	Path  string
	Line  int
}

func (e ComplexError) Error() string {
	if e.Cause == nil {
		e.Cause = Error("<nil cause>")
	}
	if e.Err == nil {
		return fmt.Sprintf("[ERROR] %s:%d:\n%s", e.Path, e.Line, e.Cause.Error())
	}
	return fmt.Sprintf("[ERROR] %s:%d: %s:\n%s", e.Path, e.Line, e.Err.Error(), e.Cause.Error())
}

func Error(msg string, args ...interface{}) error {
	return fmt.Errorf(msg, args...)
}

// Return a ComplexError which including error occurred file and path
func WrapError(cause error) error {
	if cause == nil {
		return nil
	}
	_, filepath, line, ok := runtime.Caller(1)
	if !ok {
		log.Printf("[ERROR] runtime.Caller error in WrapError.")
		return WrapComplexError(cause, nil, "", -1)
	}
	parts := strings.Split(filepath, "/")
	if len(parts) > 3 {
		filepath = strings.Join(parts[len(parts)-3:], "/")
	}
	return WrapComplexError(cause, nil, filepath, line)
}

// Return a ComplexError which including extra error message, error occurred file and path
func WrapErrorf(cause error, msg string, args ...interface{}) error {
	if cause == nil && strings.TrimSpace(msg) == "" {
		return nil
	}
	_, filepath, line, ok := runtime.Caller(1)
	if !ok {
		log.Printf("[ERROR] runtime.Caller error in WrapErrorf.")
		return WrapComplexError(cause, Error(msg), "", -1)
	}
	parts := strings.Split(filepath, "/")
	if len(parts) > 3 {
		filepath = strings.Join(parts[len(parts)-3:], "/")
	}
	return WrapComplexError(cause, fmt.Errorf(msg, args...), filepath, line)
}

func WrapComplexError(cause, err error, filepath string, fileline int) error {
	return &ComplexError{
		Cause: cause,
		Err:   err,
		Path:  filepath,
		Line:  fileline,
	}
}

// A default message of ComplexError's Err. It is format to Resource <resource-id> <operation> Failed!!! <error source>
const DefaultErrorMsg = "Resource %s %s Failed!!! %s"
const NotFoundMsg = ResourceNotFound + "!!! %s"
const DefaultTimeoutMsg = "Resource %s %s Timeout!!! %s"
const DeleteTimeoutMsg = "Resource %s Still Exists. %s Timeout!!! %s"
const WaitTimeoutMsg = "Resource %s %s Timeout In %d Seconds. Got: %s Expected: %s !!! %s"
const DataDefaultErrorMsg = "Datasource %s %s Failed!!! %s"
const IdMsg = "Resource id：%s "

const DefaultDebugMsg = "\n*************** %s Response *************** \n%s\n%s******************************\n\n"
const FailedToReachTargetStatus = "Failed to reach target status. Current status is %s."
const FailedToReachTargetAttribute = "Failed to reach value for target attribute. Current value is %s."
