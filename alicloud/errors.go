package alicloud

import (
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/denverdino/aliyungo/common"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
)

// The helpers below moved to the leaf packages errs, helper and logs;
// these aliases keep the flat layout compiling unchanged.
//
// WrapError/WrapErrorf must stay var aliases rather than forwarding
// functions: they read runtime.Caller(1) to attribute errors to the
// caller, and a wrapper would insert a frame and misattribute every
// error in the provider. The same holds for addDebug in common.go.
//
// ComplexError and ProviderError must be type aliases — the "=" is
// load-bearing: the type assertion in rdsErrorHasCode below must keep
// matching values wrapped by errs, or error-code exemptions silently
// stop working.
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

// Error codes moved to errs; re-exported as const aliases.
const (
	NotFound                = errs.NotFound
	ResourceNotfound        = errs.ResourceNotfound
	ServiceUnavailable      = errs.ServiceUnavailable
	InstanceNotFound        = errs.InstanceNotFound
	ForbiddenInstance       = errs.ForbiddenInstance
	MessageInstanceNotFound = errs.MessageInstanceNotFound
	Throttling              = errs.Throttling
	ThrottlingUser          = errs.ThrottlingUser

	// RAM Instance Not Found
	RamInstanceNotFound   = errs.RamInstanceNotFound
	AliyunGoClientFailure = errs.AliyunGoClientFailure

	LogClientTimeout = errs.LogClientTimeout

	InvalidFileSystemStatus_Ordering = errs.InvalidFileSystemStatus_Ordering
	NotFoundArticle                  = errs.NotFoundArticle
)

// Retryable error-code lists moved to errs. They stay var re-exports on
// purpose: alicloud_polardb_endpoint appends a product-specific code to
// OperationDeniedDBStatus, and a var re-export lets that mutation keep
// working on the root-level copy without touching the errs value.
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

	// details at: https://help.aliyun.com/document_detail/27300.html
	OtsTableIsTemporarilyUnavailable          = errs.OtsTableIsTemporarilyUnavailable
	OtsTunnelIsTemporarilyUnavailable         = errs.OtsTunnelIsTemporarilyUnavailable
	OtsSecondaryIndexIsTemporarilyUnavailable = errs.OtsSecondaryIndexIsTemporarilyUnavailable
	OtsSearchIndexIsTemporarilyUnavailable    = errs.OtsSearchIndexIsTemporarilyUnavailable
)

// rdsErrorHasCode matches structured RDS error codes and unwraps query retry metadata.
func rdsErrorHasCode(err error, codes ...string) bool {
	for {
		if e, ok := err.(*ComplexError); ok {
			err = e.Cause
			continue
		}
		if e, ok := err.(*rdsRetryableQueryError); ok {
			err = e.cause
			continue
		}
		break
	}
	var code string
	switch e := err.(type) {
	case *errors.ServerError:
		code = e.ErrorCode()
	case *tea.SDKError:
		code = tea.StringValue(e.Code)
	case *ProviderError:
		code = e.ErrorCode()
	case *common.Error:
		code = e.Code
	default:
		return false
	}
	for _, expected := range codes {
		if code == expected {
			return true
		}
	}
	return false
}

// A default message of ComplexError's Err. It is format to Resource <resource-id> <operation> Failed!!! <error source>
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
