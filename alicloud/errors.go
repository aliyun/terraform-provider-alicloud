package alicloud

import (
	"strings"

	"github.com/denverdino/aliyungo/common"
)

const (
	// common
	Notfound = "Not found"
	// ecs
	InstanceNotFound        = "Instance.Notfound"
	MessageInstanceNotFound = "instance is not found"
	EcsThrottling           = "Throttling"
	EcsInternalError        = "InternalError"
	// disk
	DiskIncorrectStatus       = "IncorrectDiskStatus"
	DiskCreatingSnapshot      = "DiskCreatingSnapshot"
	InstanceLockedForSecurity = "InstanceLockedForSecurity"
	SystemDiskNotFound        = "SystemDiskNotFound"
	DiskOperationConflict     = "OperationConflict"
	DiskInternalError         = "InternalError"
	DiskInvalidOperation      = "InvalidOperation.Conflict"
	// eip
	EipIncorrectStatus      = "IncorrectEipStatus"
	InstanceIncorrectStatus = "IncorrectInstanceStatus"
	HaVipIncorrectStatus    = "IncorrectHaVipStatus"
	// slb
	LoadBalancerNotFound        = "InvalidLoadBalancerId.NotFound"
	UnsupportedProtocalPort     = "UnsupportedOperationonfixedprotocalport"
	ListenerNotFound            = "The specified resource does not exist"
	ListenerAlreadyExists       = "ListenerAlreadyExists"
	ServiceIsConfiguring        = "ServiceIsConfiguring"
	BackendServerconfiguring    = "BackendServer.configuring"
	SystemBusy                  = "SystemBusy"
	SlbOrderFailed              = "OrderFailed"
	VServerGroupNotFoundMessage = "The specified VServerGroupId does not exist"
	RspoolVipExist              = "RspoolVipExist"
	InvalidParameter            = "InvalidParameter"
	InvalidRuleIdNotFound       = "InvalidRuleId.NotFound"
	RuleDomainExist             = "DomainExist"
	// security_group
	InvalidInstanceIdAlreadyExists = "InvalidInstanceId.AlreadyExists"
	InvalidSecurityGroupIdNotFound = "InvalidSecurityGroupId.NotFound"
	SgDependencyViolation          = "DependencyViolation"

	//Nat gateway
	NatGatewayInvalidRegionId            = "Invalid.RegionId"
	DependencyViolationBandwidthPackages = "DependencyViolation.BandwidthPackages"
	NotFindSnatEntryBySnatId             = "NotFindSnatEntryBySnatId"
	NotFindForwardEntryByForwardId       = "NotFindForwardEntryByForwardId"

	// vpc
	VpcQuotaExceeded = "QuotaExceeded.Vpc"
	// vswitch
	VswitcInvalidRegionId = "InvalidRegionId.NotFound"
	//vroute entry
	IncorrectRouteEntryStatus     = "IncorrectRouteEntryStatus"
	TaskConflict                  = "TaskConflict"
	RouterEntryForbbiden          = "Forbbiden"
	RouterEntryConflictDuplicated = "RouterEntryConflict.Duplicated"

	// ess
	InvalidScalingGroupIdNotFound               = "InvalidScalingGroupId.NotFound"
	IncorrectScalingConfigurationLifecycleState = "IncorrectScalingConfigurationLifecycleState"
	IncorrectScalingGroupStatus                 = "IncorrectScalingGroupStatus"
	IncorrectCapacityMaxSize                    = "IncorrectCapacity.MaxSize"
	IncorrectCapacityMinSize                    = "IncorrectCapacity.MinSize"
	ScalingActivityInProgress                   = "ScalingActivityInProgress"

	// rds
	InvalidDBNameNotFound                  = "InvalidDBName.NotFound"
	InvalidDBInstanceNameNotFound          = "InvalidDBInstanceName.NotFound"
	InvalidCurrentConnectionStringNotFound = "InvalidCurrentConnectionString.NotFound"
	NetTypeExists                          = "NetTypeExists"
	InvalidAccountNameDuplicate            = "InvalidAccountName.Duplicate"
	InvalidAccountNameNotFound             = "InvalidAccountName.NotFound"
	OperationDeniedDBInstanceStatus        = "OperationDenied.DBInstanceStatus"
	InvalidConnectionStringDuplicate       = "InvalidConnectionString.Duplicate"
	AtLeastOneNetTypeExists                = "AtLeastOneNetTypeExists"
	// oss
	OssBucketNotFound = "NoSuchBucket"
	OssBodyNotFound   = "404 Not Found"

	// RAM Instance Not Found
	RamInstanceNotFound   = "Forbidden.InstanceNotFound"
	AliyunGoClientFailure = "AliyunGoClientFailure"

	// dns
	RecordForbiddenDNSChange    = "RecordForbidden.DNSChange"
	FobiddenNotEmptyGroup       = "Fobidden.NotEmptyGroup"
	DomainRecordNotBelongToUser = "DomainRecordNotBelongToUser"

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

	//
	InvalidRamRoleNotFound       = "InvalidRamRole.NotFound"
	RoleAttachmentUnExpectedJson = "unexpected end of JSON input"
	InvalidInstanceIdNotFound    = "InvalidInstanceId.NotFound"

	RouterInterfaceIncorrectStatus                        = "IncorrectStatus"
	DependencyViolationRouterInterfaceReferedByRouteEntry = "DependencyViolation.RouterInterfaceReferedByRouteEntry"
)

func GetNotFoundErrorFromString(str string) error {
	return &common.Error{
		ErrorResponse: common.ErrorResponse{
			Code:    InstanceNotFound,
			Message: str,
		},
		StatusCode: -1,
	}
}

func NotFoundError(err error) bool {
	if e, ok := err.(*common.Error); ok &&
		(e.Code == InstanceNotFound || e.Code == RamInstanceNotFound ||
			strings.Contains(strings.ToLower(e.Message), MessageInstanceNotFound)) {
		return true
	}

	return false
}

func IsExceptedError(err error, expectCode string) bool {
	if e, ok := err.(*common.Error); ok && (e.Code == expectCode || strings.Contains(e.Message, expectCode)) {
		return true
	}

	return false
}

func RamEntityNotExist(err error) bool {
	if e, ok := err.(*common.Error); ok && strings.Contains(e.Code, "EntityNotExist") {
		return true
	}
	return false
}
