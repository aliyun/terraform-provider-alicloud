package alicloud

import (
	"github.com/denverdino/aliyungo/common"
	"strings"
)

const (
	// common
	Notfound = "Not found"
	// ecs
	InstanceNotfound = "Instance.Notfound"
	// disk
	DiskIncorrectStatus       = "IncorrectDiskStatus"
	DiskCreatingSnapshot      = "DiskCreatingSnapshot"
	InstanceLockedForSecurity = "InstanceLockedForSecurity"
	SystemDiskNotFound        = "SystemDiskNotFound"
	DiskOperationConflict     = "OperationConflict"
	// eip
	EipIncorrectStatus      = "IncorrectEipStatus"
	InstanceIncorrectStatus = "IncorrectInstanceStatus"
	HaVipIncorrectStatus    = "IncorrectHaVipStatus"
	// slb
	LoadBalancerNotFound    = "InvalidLoadBalancerId.NotFound"
	UnsupportedProtocalPort = "UnsupportedOperationonfixedprotocalport"

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

	// ess
	InvalidScalingGroupIdNotFound               = "InvalidScalingGroupId.NotFound"
	IncorrectScalingConfigurationLifecycleState = "IncorrectScalingConfigurationLifecycleState"

	// RAM Instance Not Found
	RamInstanceNotFound   = "Forbidden.InstanceNotFound"
	AliyunGoClientFailure = "AliyunGoClientFailure"

	//unknown Error
	UnknownError = "UnknownError"
)

func GetNotFoundErrorFromString(str string) error {
	return &common.Error{
		ErrorResponse: common.ErrorResponse{
			Code:    InstanceNotfound,
			Message: str,
		},
		StatusCode: -1,
	}
}

func IsExceptedError(err error, expectCode string) bool {
	if e, ok := err.(*common.Error); ok && (e.Code == expectCode || strings.Contains(e.Message, expectCode)) {
		return true
	}

	return false
}
