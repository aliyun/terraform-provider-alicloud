package alicloud

import (
	"strings"

	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/fc-go-sdk"
	"github.com/denverdino/aliyungo/common"
)

const (
	// common
	NotFound       = "NotFound"
	WaitForTimeout = "WaitForTimeout"
	// ecs
	InstanceNotFound        = "Instance.Notfound"
	MessageInstanceNotFound = "instance is not found"
	EcsThrottling           = "Throttling"
	EcsInternalError        = "InternalError"
	// disk
	InternalError       = "InternalError"
	DependencyViolation = "DependencyViolation"
	// eip
	EipIncorrectStatus         = "IncorrectEipStatus"
	InstanceIncorrectStatus    = "IncorrectInstanceStatus"
	HaVipIncorrectStatus       = "IncorrectHaVipStatus"
	COMMODITYINVALID_COMPONENT = "COMMODITY.INVALID_COMPONENT"
	// slb
	LoadBalancerNotFound        = "InvalidLoadBalancerId.NotFound"
	UnsupportedProtocalPort     = "UnsupportedOperationonfixedprotocalport"
	ListenerNotFound            = "The specified resource does not exist"
	ListenerAlreadyExists       = "ListenerAlreadyExists"
	SlbOrderFailed              = "OrderFailed"
	VServerGroupNotFoundMessage = "The specified VServerGroupId does not exist"
	RspoolVipExist              = "RspoolVipExist"
	InvalidParameter            = "InvalidParameter"
	InvalidRuleIdNotFound       = "InvalidRuleId.NotFound"
	RuleDomainExist             = "DomainExist"
	BackendServerConfiguring    = "BackendServer.configuring"
	// security_group
	InvalidInstanceIdAlreadyExists = "InvalidInstanceId.AlreadyExists"
	InvalidSecurityGroupIdNotFound = "InvalidSecurityGroupId.NotFound"
	SgDependencyViolation          = "DependencyViolation"

	//Nat gateway
	NatGatewayInvalidRegionId            = "Invalid.RegionId"
	DependencyViolationBandwidthPackages = "DependencyViolation.BandwidthPackages"
	VswitchStatusError                   = "VswitchStatusError"
	EIP_NOT_IN_GATEWAY                   = "EIP_NOT_IN_GATEWAY"
	InvalidNatGatewayIdNotFound          = "InvalidNatGatewayId.NotFound"
	// vpc
	VpcQuotaExceeded     = "QuotaExceeded.Vpc"
	InvalidVpcIDNotFound = "InvalidVpcID.NotFound"
	ForbiddenVpcNotFound = "Forbidden.VpcNotFound"
	Throttling           = "Throttling"

	// vswitch
	VswitcInvalidRegionId    = "InvalidRegionId.NotFound"
	InvalidVswitchIDNotFound = "InvalidVswitchID.NotFound"
	//vroute entry
	IncorrectRouteEntryStatus            = "IncorrectRouteEntryStatus"
	InvalidStatusRouteEntry              = "InvalidStatus.RouteEntry"
	TaskConflict                         = "TaskConflict"
	RouterEntryForbbiden                 = "Forbbiden"
	RouterEntryConflictDuplicated        = "RouterEntryConflict.Duplicated"
	InvalidCidrBlockOverlapped           = "InvalidCidrBlock.Overlapped"
	IncorrectOppositeInterfaceInfoNotSet = "IncorrectOppositeInterfaceInfo.NotSet"
	InvalidSnatTableIdNotFound           = "InvalidSnatTableId.NotFound"
	// Forward
	InvalidIpNotInNatgw           = "InvalidIp.NotInNatgw"
	InvalidForwardTableIdNotFound = "InvalidForwardTableId.NotFound"
	InvalidForwardEntryIdNotFound = "InvalidForwardEntryId.NotFound"

	// ess
	InvalidScalingGroupIdNotFound               = "InvalidScalingGroupId.NotFound"
	IncorrectScalingConfigurationLifecycleState = "IncorrectScalingConfigurationLifecycleState"
	IncorrectScalingGroupStatus                 = "IncorrectScalingGroupStatus"
	IncorrectCapacityMaxSize                    = "IncorrectCapacity.MaxSize"
	IncorrectCapacityMinSize                    = "IncorrectCapacity.MinSize"
	ScalingActivityInProgress                   = "ScalingActivityInProgress"
	EssThrottling                               = "Throttling"
	InvalidScalingRuleIdNotFound                = "InvalidScalingRuleId.NotFound"
	InvalidLifecycleHookIdNotFound              = "InvalidLifecycleHookId.NotExist"
	InvalidEssAlarmTaskNotFound                 = "404"

	// rds
	InvalidDBInstanceIdNotFound            = "InvalidDBInstanceId.NotFound"
	InvalidDBNameNotFound                  = "InvalidDBName.NotFound"
	InvalidDBInstanceNameNotFound          = "InvalidDBInstanceName.NotFound"
	InvalidCurrentConnectionStringNotFound = "InvalidCurrentConnectionString.NotFound"
	NetTypeExists                          = "NetTypeExists"
	InvalidAccountNameDuplicate            = "InvalidAccountName.Duplicate"
	InvalidAccountNameNotFound             = "InvalidAccountName.NotFound"
	InvalidConnectionStringDuplicate       = "InvalidConnectionString.Duplicate"
	AtLeastOneNetTypeExists                = "AtLeastOneNetTypeExists"
	ConnectionOperationDenied              = "OperationDenied"
	ConnectionConflictMessage              = "The requested resource is sold out in the specified zone; try other types of resources or other regions and zones"
	DBInternalError                        = "InternalError"
	// oss
	OssBucketNotFound          = "NoSuchBucket"
	OssBodyNotFound            = "404 Not Found"
	NoSuchCORSConfiguration    = "NoSuchCORSConfiguration"
	NoSuchWebsiteConfiguration = "NoSuchWebsiteConfiguration"

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

	// ram mfa
	DeleteConflictVirtualMFADeviceUser = "DeleteConflict.VirtualMFADevice.User"

	// ram group
	DeleteConflictGroupUser   = "DeleteConflict.Group.User"
	DeleteConflictGroupPolicy = "DeleteConflict.Group.Policy"

	// ram role
	DeleteConflictRolePolicy = "DeleteConflict.Role.Policy"

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

	// cdn
	ServiceBusy = "ServiceBusy"

	// KMS
	ForbiddenKeyNotFound = "Forbidden.KeyNotFound"
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
	// log
	ProjectNotExist      = "ProjectNotExist"
	IndexConfigNotExist  = "IndexConfigNotExist"
	IndexAlreadyExist    = "IndexAlreadyExist"
	LogStoreNotExist     = "LogStoreNotExist"
	InternalServerError  = "InternalServerError"
	GroupNotExist        = "GroupNotExist"
	MachineGroupNotExist = "MachineGroupNotExist"

	// OTS
	OTSObjectNotExist = "OTSObjectNotExist"

	// FC
	ServiceNotFound  = "ServiceNotFound"
	FunctionNotFound = "FunctionNotFound"
	TriggerNotFound  = "TriggerNotFound"
	AccessDenied     = "AccessDenied"

	// Vpn
	VpnNotFound              = "InvalidVpnGatewayInstanceId.NotFound"
	VpnForbidden             = "Forbidden"
	VpnForbiddenRelease      = "ForbiddenRelease"
	VpnForbiddenSubUser      = "Forbbiden.SubUser"
	CgwNotFound              = "InvalidCustomerGatewayInstanceId.NotFound"
	ResQuotaFull             = "Resource.QuotaFull"
	VpnConnNotFound          = "InvalidVpnConnectionInstanceId.NotFound"
	InvalidIpAddress         = "InvalidIpAddress.AlreadyExist"
	SslVpnServerNotFound     = "InvalidSslVpnServerId.NotFound"
	SslVpnClientCertNotFound = "InvalidSslVpnClientCertId.NotFound"
	VpnConfiguring           = "VpnGateway.Configuring"
	VpnInvalidSpec           = "InvalidSpec.NotFound"
	VpnEnable                = "enable"
	// CEN
	OperationBlocking                = "Operation.Blocking"
	ParameterCenInstanceIdNotExist   = "ParameterCenInstanceId"
	CenQuotaExceeded                 = "QuotaExceeded.CenCountExceeded"
	InvalidCenInstanceStatus         = "InvalidOperation.CenInstanceStatus"
	InvalidChildInstanceStatus       = "InvalidOperation.ChildInstanceStatus"
	ParameterInstanceIdNotExist      = "ParameterInstanceId"
	ForbiddenRelease                 = "Forbidden.Release"
	InvalidCenBandwidthLimitsNotZero = "InvalidOperation.CenBandwidthLimitsNotZero"
	ParameterBwpInstanceId           = "ParameterBwpInstanceId"
	InvalidBwpInstanceStatus         = "InvalidOperation.BwpInstanceStatus"
	InvalidBwpBusinessStatus         = "InvalidOperation.BwpBusinessStatus"
	// kv-store
	InvalidKVStoreInstanceIdNotFound = "InvalidInstanceId.NotFound"
	// MNS
	QueueNotExist        = "QueueNotExist"
	TopicNotExist        = "TopicNotExist"
	SubscriptionNotExist = "SubscriptionNotExist"
	//HaVip
	InvalidHaVipIdNotFound = "InvalidHaVipId.NotFound"
)

var SlbIsBusy = []string{"SystemBusy", "OperationBusy", "ServiceIsStopping", "BackendServer.configuring", "ServiceIsConfiguring"}
var EcsNotFound = []string{"InvalidInstanceId.NotFound", "Forbidden.InstanceNotFound"}
var DiskInvalidOperation = []string{"IncorrectDiskStatus", "IncorrectInstanceStatus", "OperationConflict", InternalError, "InvalidOperation.Conflict", "IncorrectDiskStatus.Initializing"}
var OperationDeniedDBStatus = []string{"OperationDenied.DBStatus", "OperationDenied.DBInstanceStatus", DBInternalError}

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
	if e, ok := err.(*common.Error); ok &&
		(e.Code == InstanceNotFound || e.Code == RamInstanceNotFound || e.Code == NotFound ||
			strings.Contains(strings.ToLower(e.Message), MessageInstanceNotFound)) {
		return true
	}

	if e, ok := err.(*errors.ServerError); ok &&
		(e.ErrorCode() == InstanceNotFound || e.ErrorCode() == RamInstanceNotFound || e.ErrorCode() == NotFound ||
			strings.Contains(strings.ToLower(e.Message()), MessageInstanceNotFound)) {
		return true
	}

	if e, ok := err.(*ProviderError); ok &&
		(e.ErrorCode() == InstanceNotFound || e.ErrorCode() == RamInstanceNotFound || e.ErrorCode() == NotFound ||
			strings.Contains(strings.ToLower(e.Message()), MessageInstanceNotFound)) {
		return true
	}

	return false
}

func IsExceptedError(err error, expectCode string) bool {
	if e, ok := err.(*common.Error); ok && (e.Code == expectCode || strings.Contains(e.Message, expectCode)) {
		return true
	}

	if e, ok := err.(*errors.ServerError); ok && (e.ErrorCode() == expectCode || strings.Contains(e.Message(), expectCode)) {
		return true
	}

	if e, ok := err.(*ProviderError); ok && (e.ErrorCode() == expectCode || strings.Contains(e.Message(), expectCode)) {
		return true
	}

	if e, ok := err.(*sls.Error); ok && (e.Code == expectCode || strings.Contains(e.Message, expectCode)) {
		return true
	}

	if e, ok := err.(oss.ServiceError); ok && (e.Code == expectCode || strings.Contains(e.Message, expectCode)) {
		return true
	}
	return false
}

func IsExceptedErrors(err error, expectCodes []string) bool {
	for _, code := range expectCodes {
		if e, ok := err.(*common.Error); ok && (e.Code == code || strings.Contains(e.Message, code)) {
			return true
		}

		if e, ok := err.(*errors.ServerError); ok && (e.ErrorCode() == code || strings.Contains(e.Message(), code)) {
			return true
		}

		if e, ok := err.(*ProviderError); ok && (e.ErrorCode() == code || strings.Contains(e.Message(), code)) {
			return true
		}
		if e, ok := err.(*sls.Error); ok && (e.Code == code || strings.Contains(e.Message, code)) {
			return true
		}
		if e, ok := err.(oss.ServiceError); ok && (e.Code == code || strings.Contains(e.Message, code)) {
			return true
		}
		if e, ok := err.(*fc.ServiceError); ok && (e.ErrorCode == code || strings.Contains(e.ErrorMessage, code)) {
			return true
		}
	}
	return false
}

func RamEntityNotExist(err error) bool {
	if e, ok := err.(*common.Error); ok && strings.Contains(e.Code, "EntityNotExist") {
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
